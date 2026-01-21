// @title User CRUD API
// @version 2.0
// @description REST API untuk manajemen user
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user-crud/internal/application/command"
	"user-crud/internal/application/query"
	"user-crud/internal/config"
	"user-crud/internal/infrastructure/cache"
	"user-crud/internal/infrastructure/http/handler"
	"user-crud/internal/infrastructure/http/router"
	"user-crud/internal/infrastructure/persistence"
	"user-crud/internal/infrastructure/tracing"

	_ "user-crud/docs"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Jaeger tracing
	jaegerEndpoint := getEnv("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces")
	shutdown, err := tracing.InitTracer("user-crud-service", jaegerEndpoint)
	if err != nil {
		log.Printf("Warning: Failed to initialize tracer: %v", err)
	} else {
		defer shutdown(context.Background())
		log.Println("Jaeger tracing initialized successfully")
	}

	// Initialize database connection
	dbpool, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbpool.Close()

	// Run migrations
	if err := runMigrations(dbpool); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis cache
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisCache, err := cache.NewRedisCache(redisHost, redisPort, 5*time.Minute)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisCache.Close()
	log.Println("Successfully connected to Redis")

	// Initialize repository
	userRepo := persistence.NewPostgresUserRepository(dbpool)

	// Initialize command handlers (WITH CACHE)
	createUserHandler := command.NewCreateUserHandler(userRepo, redisCache)
	updateUserHandler := command.NewUpdateUserHandler(userRepo, redisCache)
	deleteUserHandler := command.NewDeleteUserHandler(userRepo, redisCache)
	changePasswordHandler := command.NewChangePasswordHandler(userRepo, redisCache)

	// Initialize query handlers (WITH CACHE)
	getUserHandler := query.NewGetUserHandler(userRepo, redisCache)
	listUsersHandler := query.NewListUsersHandler(userRepo)
	searchUsersHandler := query.NewSearchUsersHandler(userRepo)

	// Initialize HTTP handler
	h := handler.NewHandler(
		createUserHandler,
		updateUserHandler,
		deleteUserHandler,
		changePasswordHandler,
		getUserHandler,
		listUsersHandler,
		searchUsersHandler,
		dbpool,
		redisCache,
	)

	// Setup router
	r := router.SetupRouter(h)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func initDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2

	var dbpool *pgxpool.Pool
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		dbpool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			if err = dbpool.Ping(context.Background()); err == nil {
				log.Println("Successfully connected to database")
				return dbpool, nil
			}
		}

		waitTime := time.Duration(i+1) * 2 * time.Second
		log.Printf("Failed to connect to database, retrying in %v... (attempt %d/%d)", waitTime, i+1, maxRetries)
		time.Sleep(waitTime)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

func runMigrations(dbpool *pgxpool.Pool) error {
	log.Println("Running database migrations...")

	migration := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		age INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);
	CREATE INDEX IF NOT EXISTS idx_users_age ON users(age);
	CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
	`

	_, err := dbpool.Exec(context.Background(), migration)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
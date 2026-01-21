package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool creates a new PostgreSQL connection pool with retry logic
func NewPostgresPool(host, port, user, password, dbname string) (*pgxpool.Pool, error) {
	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	log.Printf("📡 Attempting database connection to %s:%s", host, port)
	log.Printf("🔧 Database: %s, User: %s", dbname, user)

	// Retry logic dengan exponential backoff
	maxRetries := 5
	var pool *pgxpool.Pool
	var err error

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		
		// Try to create connection pool
		pool, err = pgxpool.New(ctx, connStr)
		if err == nil {
			// Test connection dengan ping
			err = pool.Ping(ctx)
			if err == nil {
				cancel()
				log.Printf("✅ Successfully connected to database at %s:%s", host, port)
				return pool, nil
			}
		}
		
		cancel()
		
		waitTime := time.Duration((i+1)*2) * time.Second
		log.Printf("❌ Failed to connect to database, retrying in %v... (attempt %d/%d)", 
			waitTime, i+1, maxRetries)
		
		if err != nil {
			log.Printf("   Error: %v", err)
		}
		
		time.Sleep(waitTime)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
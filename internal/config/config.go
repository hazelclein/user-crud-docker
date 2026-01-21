package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "postgres"),      // ✅ GANTI: "localhost" → "postgres"
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "userdb"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	// Log configuration untuk debugging
	log.Printf("📋 Configuration loaded:")
	log.Printf("   DB Host: %s", cfg.DBHost)
	log.Printf("   DB Port: %s", cfg.DBPort)
	log.Printf("   DB Name: %s", cfg.DBName)
	log.Printf("   Server Port: %s", cfg.ServerPort)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		log.Printf("✅ Environment variable %s = %s", key, value)
		return value
	}
	log.Printf("⚠️  Environment variable %s not set, using default: %s", key, defaultValue)
	return defaultValue
}
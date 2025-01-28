// internal/config/config.go
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress   string
	Database        DatabaseConfig
	SentryDSN       string
	Environment     string
	SentryTraceRate float64
}

type DatabaseConfig struct {
	URL             string // Full database URL
	MaxConnections  int
	MinConnections  int
	MaxConnLifetime time.Duration
}

// getEnvInt helper function to get int values from env with default fallback
func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

// getEnvFloat helper function to get float values from env with default fallback
func getEnvFloat(key string, defaultVal float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultVal
}

func Load() (*Config, error) {
	// Load .env file if present
	godotenv.Load() // Ignoring error as .env file is optional

	// Database connection pool settings
	// For transaction pooler, using lower connection pool defaults
	dbMaxConns := getEnvInt("DB_POOL_MAX_CONNS", 20)
	dbMinConns := getEnvInt("DB_POOL_MIN_CONNS", 5)
	dbLifetime := getEnvInt("DB_POOL_MAX_CONN_LIFETIME", 30)

	// Server settings
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080" // default port
	}

	// Sentry settings
	sentryTraceRate := getEnvFloat("SENTRY_TRACE_RATE", 1.0)
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development" // default environment
	}

	return &Config{
		ServerAddress: ":" + serverPort,
		Database: DatabaseConfig{
			URL:             os.Getenv("SUPABASE_TRANSACTION_POOLER"),
			MaxConnections:  dbMaxConns,
			MinConnections:  dbMinConns,
			MaxConnLifetime: time.Duration(dbLifetime) * time.Minute,
		},
		SentryDSN:       os.Getenv("SENTRY_DSN"),
		Environment:     environment,
		SentryTraceRate: sentryTraceRate,
	}, nil
}

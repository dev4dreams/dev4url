package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	URL             string // Full database URL
	MaxConnections  int
	MinConnections  int
	MaxConnLifetime time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	// For transaction pooler, using lower connection pool defaults
	dbMaxConns := getEnvInt("DB_POOL_MAX_CONNS", 20)
	dbMinConns := getEnvInt("DB_POOL_MIN_CONNS", 5)
	dbLifetime := getEnvInt("DB_POOL_MAX_CONN_LIFETIME", 30)

	return &Config{
		Database: DatabaseConfig{
			URL:             os.Getenv("SUPABASE_TRANSACTION_POOLER"),
			MaxConnections:  dbMaxConns,
			MinConnections:  dbMinConns,
			MaxConnLifetime: time.Duration(dbLifetime) * time.Minute,
		},
	}, nil
}

func getEnvInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultVal
}

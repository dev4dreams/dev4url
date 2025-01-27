package db

import (
	"database/sql"
	"fmt"

	"github.com/dev4dreams/dev4url/internal/config"
	"github.com/dev4dreams/dev4url/internal/models"
	_ "github.com/lib/pq"
)

// DatabaseInterface defines the behavior for database operations
type DatabaseInterface interface {
	CreateURL(payload *models.CreateUrlPayload) (*models.URLResponse, error)
	Close() error
	VerifyConnection() error
}

type Database struct {
	*sql.DB
}

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// New creates a new database connection
func New(config *config.DatabaseConfig) (*Database, error) {
	connStr := config.URL + "?sslmode=require&pool_mode=transaction&statement_cache_mode=describe"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return &Database{db}, nil
}

// CreateURL inserts a new URL record into the database
func (db *Database) CreateURL(url *models.CreateUrlPayload) (*models.URLResponse, error) {
	var response models.URLResponse

	query := `
        INSERT INTO urls (
            short_url,
            original_url,
            custom_url
        ) VALUES (
            $1, $2, $3
        )
        RETURNING id, created_at, short_url, original_url,
                  custom_url, clicks, active, updated_at`

	err := db.QueryRow(
		query,
		url.ShortenUrl,
		url.OriginalUrl,
		url.CustomUrl,
	).Scan(
		&response.ID,
		&response.CreatedAt,
		&response.ShortURL,
		&response.OriginalURL,
		&response.CustomURL,
		&response.Clicks,
		&response.Active,
		&response.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}

	return &response, nil
}

// VerifyConnection checks if the database connection is still alive
func (db *Database) VerifyConnection() error {
	return db.Ping()
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.DB.Close()
}

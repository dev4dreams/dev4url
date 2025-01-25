// package db

// import (
// 	"database/sql"
// 	"fmt"

// 	"github.com/dev4dreams/dev4url/internal/config"
// 	"github.com/dev4dreams/dev4url/internal/models"
// 	_ "github.com/lib/pq"
// )

// type DatabaseInterface interface {
// 	CreateURL(payload *models.CreateUrlPayload) (*models.URLResponse, error)
// }

// type Database struct {
// 	*sql.DB
// }

// func New(cfg *config.DatabaseConfig) (*Database, error) {
// 	// The URL already contains the pooler configuration
// 	// Just append any additional parameters we need
// 	connStr := cfg.URL + "?sslmode=require&pool_mode=transaction&statement_cache_mode=describe"

// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		return nil, fmt.Errorf("error opening database: %w", err)
// 	}

// 	db.SetMaxOpenConns(cfg.MaxConnections)
// 	db.SetMaxIdleConns(cfg.MinConnections)
// 	db.SetConnMaxLifetime(cfg.MaxConnLifetime)

// 	if err := db.Ping(); err != nil {
// 		return nil, fmt.Errorf("error connecting to database: %w", err)
// 	}

// 	return &Database{db}, nil
// }

// func (db *Database) VerifyConnection() error {
// 	// Check 1: Basic connectivity
// 	var now string
// 	if err := db.QueryRow("SELECT NOW()").Scan(&now); err != nil {
// 		return fmt.Errorf("failed to query time: %w", err)
// 	}

// 	// Check 2: Verify we're on Supabase by checking schema existence
// 	var exists bool
// 	if err := db.QueryRow(`
//         SELECT EXISTS (
//             SELECT FROM pg_catalog.pg_namespace
//             WHERE nspname = 'auth'
//         )`).Scan(&exists); err != nil {
// 		return fmt.Errorf("failed to check schema: %w", err)
// 	}

// 	// Check 3: Get connection details
// 	var clientAddr string
// 	if err := db.QueryRow(`
//         SELECT client_addr
//         FROM pg_stat_activity
//         WHERE pid = pg_backend_pid()`).Scan(&clientAddr); err != nil {
// 		return fmt.Errorf("failed to get connection details: %w", err)
// 	}

// 	fmt.Printf("✅ Connected successfully to Supabase!\n")
// 	fmt.Printf("🕒 Current time: %s\n", now)
// 	fmt.Printf("🔌 Connected from: %s\n", clientAddr)
// 	fmt.Printf("🔐 Auth schema exists: %v\n", exists)

// 	// Additional simple query to verify write access
// 	_, err := db.Exec("SELECT 1")
// 	if err != nil {
// 		return fmt.Errorf("failed to execute test query: %w", err)
// 	}

// 	return nil
// }

// func (db *Database) CreateURL(url *models.CreateUrlPayload) (*URLResponse, error) {
// 	var response URLResponse

// 	query := `
//         INSERT INTO urls (
//             short_url,
//             original_url,
//             custom_url
//         ) VALUES (
//             $1, $2, $3
//         )
//         RETURNING id, created_at, short_url, original_url,
//                   custom_url, clicks, active, updated_at`

// 	err := db.QueryRow(
// 		query,
// 		url.ShortenUrl,
// 		url.OriginalUrl,
// 		url.CustomUrl,
// 	).Scan(
// 		&response.ID,
// 		&response.CreatedAt,
// 		&response.ShortURL,
// 		&response.OriginalURL,
// 		&response.CustomURL,
// 		&response.Clicks,
// 		&response.Active,
// 		&response.UpdatedAt,
// 	)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create URL: %w", err)
// 	}

//		return &response, nil
//	}
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

	// connStr := fmt.Sprintf(
	// 	"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
	// 	config.Host,
	// 	config.Port,
	// 	config.User,
	// 	config.Password,
	// 	config.DBName,
	// 	config.SSLMode,
	// )

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

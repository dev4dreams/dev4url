package main

import (
	"context"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dev4dreams/dev4url/internal/config"
	"github.com/dev4dreams/dev4url/internal/core"
	"github.com/dev4dreams/dev4url/internal/db"
	"github.com/dev4dreams/dev4url/internal/handlers"
	"github.com/dev4dreams/dev4url/internal/middleware"
	"github.com/dev4dreams/dev4url/internal/services/safebrowsing"
	"github.com/dev4dreams/dev4url/internal/utils"
	"golang.org/x/time/rate"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err = middleware.InitSentry(os.Getenv("SENTRY_DSN")); err != nil {
		log.Fatalf("Failed to initialize Sentry: %v", err)
	}

	defer middleware.FlushSentry(2 * time.Second)

	// Initialize URL shortener
	generator, err := core.NewGenerator(1)
	if err != nil {
		log.Fatalf("Failed to create URL generator: %v", err)
	}

	// Initialize URL validator with default config
	validator := utils.NewURLValidator(utils.DefaultConfig())

	// Initialize Safe Browsing service
	safeBrowsingKey := os.Getenv("GCP_SAFE_BROWSING_API_KEY")
	safeBrowsingService := safebrowsing.NewSafeBrowsingService(safeBrowsingKey)

	// Initialize URL handler
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080" // Default for development
	}

	// Initialize database connection
	database, err := db.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Verify database connection
	if err := database.VerifyConnection(); err != nil {
		log.Fatalf("Failed to verify database connection: %v", err)
	}

	// Initialize handlers
	redirectHandler := handlers.NewRedirectHandler(database)
	createUrlHandler := handlers.NewURLHandler(validator, safeBrowsingService, generator, baseURL, database)

	// Create router/mux
	mux := http.NewServeMux()

	// Initialize rate limiter
	// Adjust these values based on your requirements
	limiter := middleware.NewIPRateLimiter(rate.Limit(3), 5) // 100 requests per second, burst of 10

	// Register routes with middleware
	mux.Handle("/shortUrl/get", middleware.CORS(
		limiter.RateLimit(
			http.HandlerFunc(redirectHandler.HandleRedirect),
		),
	))

	mux.HandleFunc("/shortUrl/post", createUrlHandler.CreateShortURL)
	handler := middleware.CORS(middleware.SentryHandler(limiter.RateLimit(mux)))

	// Create server with timeouts
	server := &http.Server{
		Addr: cfg.ServerAddress,
		// Handler:      mux,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

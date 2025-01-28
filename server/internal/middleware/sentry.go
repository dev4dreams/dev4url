// internal/middleware/sentry.go
package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

// InitSentry initializes the Sentry client with the provided DSN
func InitSentry(dsn string) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		Environment:      os.Getenv("APP_ENV"), // Use environment variable
		Debug:            os.Getenv("APP_ENV") == "development",
		ServerName:       "urlshortener-service",
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// You can modify or filter events before they're sent
			// For example, remove sensitive data
			return event
		},
	})
	if err != nil {
		return fmt.Errorf("sentry initialization failed: %v", err)
	}

	// Verify connection
	sentry.CaptureMessage("Sentry initialized successfully")
	sentry.Flush(2 * time.Second)

	return nil
}

// SentryHandler middleware for standard http handlers
func SentryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub := sentry.CurrentHub().Clone()
		ctx := sentry.SetHubOnContext(r.Context(), hub)
		r = r.WithContext(ctx)

		// Configure scope with request info
		hub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetRequest(r)
			scope.SetTag("handler", r.URL.Path)
			scope.SetTag("method", r.Method)
			// Add user info if available
			if userID := r.Header.Get("X-User-ID"); userID != "" {
				scope.SetUser(sentry.User{
					ID: userID,
				})
			}
		})

		// Recover from panics
		defer func() {
			if err := recover(); err != nil {
				eventID := hub.RecoverWithContext(
					ctx,
					err,
				)
				// Log the error ID for tracking
				fmt.Printf("Captured error with ID: %s\n", *eventID)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// CaptureError helper function to capture errors with additional context
func CaptureError(err error, tags map[string]string) *sentry.EventID {
	if err == nil {
		return nil
	}

	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(tags)
	})
	return hub.CaptureException(err)
}

// FlushSentry ensures all events are sent to Sentry
func FlushSentry(timeout time.Duration) bool {
	return sentry.Flush(timeout)
}

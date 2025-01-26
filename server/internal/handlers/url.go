// internal/handlers/url.go
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dev4dreams/dev4url/internal/core"
	"github.com/dev4dreams/dev4url/internal/db"
	"github.com/dev4dreams/dev4url/internal/middleware"
	"github.com/dev4dreams/dev4url/internal/models"
	"github.com/dev4dreams/dev4url/internal/services/safebrowsing"
	"github.com/dev4dreams/dev4url/internal/utils"
	"github.com/getsentry/sentry-go"
)

type URLHandler struct {
	UrlValidator utils.URLValidatorInterface
	SafeBrowsing safebrowsing.SafeBrowsingChecker
	Shortener    *core.Generator
	BaseURL      string
	Db           db.DatabaseInterface
}

func NewURLHandler(
	validator utils.URLValidatorInterface,
	safeBrowsing safebrowsing.SafeBrowsingChecker,
	shortener *core.Generator,
	baseURL string,
	db db.DatabaseInterface,
) *URLHandler {
	return &URLHandler{
		UrlValidator: validator,
		SafeBrowsing: safeBrowsing,
		Shortener:    shortener,
		BaseURL:      baseURL,
		Db:           db,
	}
}

func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	hub := sentry.GetHubFromContext((r.Context()))

	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("handler", "create_short_url")
		scope.SetTag("method", r.Method)
	})

	// Only allow POST method
	if r.Method != http.MethodPost {
		middleware.CaptureError(fmt.Errorf("method not allowed: %s", r.Method), map[string]string{
			"error_type": "method_not_allowed",
			"method":     r.Method,
		})
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req models.CreateUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.CaptureError(err, map[string]string{
			"error_type": "invalid_request",
			"error_step": "body_decode",
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate original URL
	validationResult := h.UrlValidator.ValidateURL(r.Context(), req.OriginalURL)
	if !validationResult.IsValid {
		middleware.CaptureError(
			fmt.Errorf("URL validation failed: %v", validationResult.Errors),
			map[string]string{
				"error_type":   "validation_error",
				"original_url": req.OriginalURL,
			},
		)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "URL validation failed",
			"errors": validationResult.Errors,
		})

		return
	}

	// Check if URL is safe
	isSafe, err := h.SafeBrowsing.IsURLSafe(req.OriginalURL)
	if err != nil {
		middleware.CaptureError(err, map[string]string{
			"error_type":   "safebrowsing_error",
			"original_url": req.OriginalURL,
		})
		log.Printf("SafeBrowsing check failed: %v", err)
		http.Error(w, "Error checking URL safety", http.StatusInternalServerError)
		return
	}
	if !isSafe {
		middleware.CaptureError(
			fmt.Errorf("unsafe URL detected: %s", req.OriginalURL),
			map[string]string{
				"error_type":   "unsafe_url",
				"original_url": req.OriginalURL,
			},
		)
		http.Error(w, "URL detected as potentially harmful", http.StatusBadRequest)
		return
	}

	// Handle custom URL if provided
	var shortCode string
	if req.CustomURL != "" {
		// Future develop
		// Here you would:
		// 1. Validate the custom URL format
		// 2. Check if it's available in your database
		// 3. Use it if valid and available
		// For now, we'll return an error as it's not implemented
		middleware.CaptureError(
			fmt.Errorf("custom URL requested but not implemented"),
			map[string]string{
				"error_type": "not_implemented",
				"feature":    "custom_url",
			},
		)
		http.Error(w, "Custom URLs not implemented yet", http.StatusNotImplemented)
		return
	} else {
		shortCode, err = h.Shortener.GenerateShortURL()
		if err != nil {
			var statusCode int
			var message string

			switch {
			case errors.Is(err, core.ErrInvalidWorkerID):
				statusCode = http.StatusInternalServerError
				message = "Server configuration error"
			case errors.Is(err, core.ErrClockMovedBackwards):
				statusCode = http.StatusServiceUnavailable
				message = "Temporary server error, please try again"
			default:
				statusCode = http.StatusInternalServerError
				message = "Internal server error"
			}
			middleware.CaptureError(err, map[string]string{
				"error_type":   "shortcode_generation",
				"error_detail": err.Error(),
				"status_code":  fmt.Sprintf("%d", statusCode),
			})
			http.Error(w, message, statusCode)
			return
		}

	}

	urlPayload := &models.CreateUrlPayload{
		ShortenUrl:  shortCode,
		OriginalUrl: req.OriginalURL,
		CustomUrl:   req.CustomURL,
	}

	dbResponse, err := h.Db.CreateURL(urlPayload)
	if err != nil {
		middleware.CaptureError(err, map[string]string{
			"error_type":   "database_error",
			"error_step":   "create_url",
			"original_url": req.OriginalURL,
			"short_code":   shortCode,
		})
		http.Error(w, "Error saving URL to database", http.StatusInternalServerError)
		return
	}

	// Construct full short URL
	fullShortURL := h.BaseURL + "/" + dbResponse.ShortURL
	fmt.Println("ShortURL created: %v", fullShortURL)
	w.Header().Set("Content-Type", "application/json")
	response := models.CreateUrlResponse{
		ShortenUrl: fullShortURL,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		middleware.CaptureError(err, map[string]string{
			"error_type": "response_encoding",
			"short_url":  fullShortURL,
		})
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

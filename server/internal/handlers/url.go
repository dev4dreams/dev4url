// internal/handlers/url.go
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/dev4dreams/dev4url/internal/core"
	"github.com/dev4dreams/dev4url/internal/db"
	"github.com/dev4dreams/dev4url/internal/models"
	"github.com/dev4dreams/dev4url/internal/services/safebrowsing"
	"github.com/dev4dreams/dev4url/internal/utils"
)

type URLHandler struct {
	urlValidator *utils.URLValidator
	safeBrowsing *safebrowsing.SafeBrowsingService
	shortener    *core.Generator
	baseURL      string
	db           *db.Database
}

func NewURLHandler(
	validator *utils.URLValidator,
	safeBrowsing *safebrowsing.SafeBrowsingService,
	shortener *core.Generator,
	baseURL string,
	db *db.Database,
) *URLHandler {
	return &URLHandler{
		urlValidator: validator,
		safeBrowsing: safeBrowsing,
		shortener:    shortener,
		baseURL:      baseURL,
		db:           db,
	}
}

func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req models.CreateUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate original URL
	validationResult := h.urlValidator.ValidateURL(r.Context(), req.OriginalURL)
	if !validationResult.IsValid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "URL validation failed",
			"errors": validationResult.Errors,
		})
		return
	}

	// Check if URL is safe
	isSafe, err := h.safeBrowsing.IsURLSafe(req.OriginalURL)
	if err != nil {
		log.Printf("SafeBrowsing check failed: %v", err)
		http.Error(w, "Error checking URL safety", http.StatusInternalServerError)
		return
	}
	if !isSafe {
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
		http.Error(w, "Custom URLs not implemented yet", http.StatusNotImplemented)
		return
	} else {
		shortCode, err = h.shortener.GenerateShortURL()
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
			http.Error(w, message, statusCode)
			return
		}
	}

	// // Generate short URL
	// shortCode, err := h.shortener.GenerateShortURL()
	// if err != nil {
	// 	http.Error(w, "Error generating short URL", http.StatusInternalServerError)
	// 	return
	// }

	urlPayload := &models.CreateUrlPayload{
		ShortenUrl:  shortCode,
		OriginalUrl: req.OriginalURL,
		CustomUrl:   req.CustomURL,
	}

	dbResponse, err := h.db.CreateURL(urlPayload)
	if err != nil {
		http.Error(w, "Error saving URL to database", http.StatusInternalServerError)
		return
	}

	// Construct full short URL
	// fullShortURL := h.baseURL + "/" + shortCode
	fullShortURL := h.baseURL + "/" + dbResponse.ShortURL

	w.Header().Set("Content-Type", "application/json")
	response := models.CreateUrlResponse{
		ShortenUrl: fullShortURL,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

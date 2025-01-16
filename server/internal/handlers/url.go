// internal/handlers/url.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dev4dreams/dev4url/internal/models"
	"github.com/dev4dreams/dev4url/internal/utils"
)

type URLHandler struct {
	validator *utils.URLValidator
}

func NewURLHandler(validator *utils.URLValidator) *URLHandler {
	return &URLHandler{
		validator: validator,
	}
}

func (h *URLHandler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUrlRequest
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// body, err := io.ReadAll(r.Body)
	fmt.Printf("Request body: %s\n", r.Body)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Use the validator instance from the struct
	if err := h.validator.ValidateURL(req.OriginalURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.CreateUrlResponse{
		ShortenUrl: "http://yourdomain.com/abc123", // This will be dynamic later
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// internal/handlers/url.go
func (h *URLHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Method not allowed",
		})
		return
	}

	// response := map[string]string{
	// 	"status":  "ok",
	// 	"message": "Server is running",
	// }

	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"server": map[string]string{
			"version":     "1.0.0",
			"environment": os.Getenv("GO_ENV"),
		},
		"services": map[string]string{
			"url_validator": "active",
			"blacklist":     "active",
		},
		"endpoints": map[string]string{
			"/api/health":  "GET - Health Check",
			"/api/shorten": "POST - URL Shortener",
		},
	}
	json.NewEncoder(w).Encode(response)
}

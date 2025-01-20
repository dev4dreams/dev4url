package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dev4dreams/dev4url/internal/db"
	"github.com/dev4dreams/dev4url/internal/models"
)

type RedirectHandler struct {
	db *db.Database
}

// NewRedirectHandler creates a new handler instance with database connection
func NewRedirectHandler(database *db.Database) *RedirectHandler {
	return &RedirectHandler{
		db: database,
	}
}

func (h *RedirectHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req models.GetOriginalUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.ShortenUrl == "" {
		http.Error(w, "Shortened URL is required", http.StatusBadRequest)
		return
	}

	// Query the database using the existing connection
	var originalURL string
	err := h.db.QueryRow(`
		UPDATE urls 
		SET 
			clicks = clicks + 1,
			updated_at = NOW()
		WHERE short_url = $1 AND active = true 
		RETURNING original_url`,
		req.ShortenUrl,
	).Scan(&originalURL)

	// Handle potential errors
	if err != nil {
		if errors.Is(err, errors.New("sql: no rows in result set")) {
			http.Error(w, "URL not found or inactive", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Prepare and send response
	response := models.GetOriginalUrlResponse{
		OriginalURL: originalURL,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

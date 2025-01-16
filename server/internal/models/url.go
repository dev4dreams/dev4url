// internal/models/url.go
package models

import "time"

// for creating a new shorten url
type CreateUrlRequest struct {
	OriginalURL string  `json:"original_url"`
	CustomURL   *string `json:"custom_url,omitempty"` // still optional
}

// for single url response
type CreateUrlResponse struct {
	ShortenUrl string `json:"shortenUrl"`
}

// when url been called
type GetOriginalUrlRequest struct {
	ShortenUrl string `json:"shortenUrl"`
}
type GetOriginalUrlResponse struct {
	OriginalURL string `json:"original_url"`
}

// This struct is for reading full URL data from DB
type URLResponse struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CustomURL   *string   `json:"custom_url,omitempty"`
	Clicks      int       `json:"clicks"`
	Active      bool      `json:"active"`
	UpdatedAt   time.Time `json:"updated_at"`
}

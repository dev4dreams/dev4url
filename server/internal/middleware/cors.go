// internal/middleware/cors.go
package middleware

import (
	"net/http"
	"os"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

		// if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
		// }

		// w.Header().Set("Access-Control-Allow-Origin", "https://dev4url.cc")
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Handle the actual request
		if r.Method == http.MethodGet || r.Method == http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		// If method is neither OPTIONS, GET, nor POST
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
}

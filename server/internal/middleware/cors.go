// internal/middleware/cors.go
package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		// if r.Method == "OPTIONS" {
		// 	w.WriteHeader(http.StatusOK)
		// 	return
		// }

		// next.ServeHTTP(w, r)
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

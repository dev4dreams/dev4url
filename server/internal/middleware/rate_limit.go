// internal/middleware/ratelimit.go
package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// IPRateLimiter stores rate limiters for each IP address
type IPRateLimiter struct {
	// maps IP addresses to their rate limiters
	ips map[string]*rate.Limiter
	// mutex for thread-safe operations
	mu *sync.RWMutex
	// rate limit (requests per second)
	r rate.Limit
	// burst size (max requests at once)
	b int
}

// NewIPRateLimiter creates a new rate limiter with specified rate and burst
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// getLimiter retrieves or creates a rate limiter for the given IP
func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		// Create new rate limiter if none exists for this IP
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// RateLimit middleware function to control request rates
func (i *IPRateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get rate limiter for this IP
		limiter := i.getLimiter(r.RemoteAddr)

		// Check if request is allowed
		if !limiter.Allow() {
			http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
			return
		}

		// Request is allowed, continue to next handler
		next.ServeHTTP(w, r)
	})
}

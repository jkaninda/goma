package middleware

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// RateLimitMiddleware limits requests based on the RateLimiter
func (rl *RateLimiter) RateLimitMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.Allow() {
				// Rate limit exceeded, return a 429 Too Many Requests response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.WriteHeader(http.StatusUnauthorized)
				err := json.NewEncoder(w).Encode(ProxyResponseError{
					Success: false,
					Code:    http.StatusTooManyRequests,
					Message: "Too many requests. Please try again later.",
				})
				if err != nil {
					return
				}
				return
			}

			// Proceed to the next handler if rate limit is not exceeded
			next.ServeHTTP(w, r)
		})
	}
}

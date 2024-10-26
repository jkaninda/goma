package middleware

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jkaninda/goma/internal/logger"
	"net/http"
	"time"
)

// RateLimitMiddleware limits request based on the number of tokens peer minutes.
func (rl *TokenRateLimiter) RateLimitMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.Allow() {
				// Rate limit exceeded, return a 429 Too Many Requests response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
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

// RateLimitMiddleware limits request based on the number of requests peer minutes.
func (rl *RateLimiter) RateLimitMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//TODO:
			clientID := r.RemoteAddr
			logger.Info(clientID)

			rl.mu.Lock()
			client, exists := rl.ClientMap[clientID]
			if !exists || time.Now().After(client.ExpiresAt) {
				client = &Client{
					RequestCount: 0,
					ExpiresAt:    time.Now().Add(rl.Window),
				}
				rl.ClientMap[clientID] = client
			}
			client.RequestCount++
			rl.mu.Unlock()

			if client.RequestCount > rl.Requests {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
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

package middleware

import (
	"encoding/json"
	"github.com/jkaninda/goma-gateway/util"
	"net/http"
	"strings"
	"time"
)

// BlocklistMiddleware checks if the request path is forbidden and returns 403 Forbidden
func (blockList BlockListMiddleware) BlocklistMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, block := range blockList.List {
			if isPathBlocked(r.URL.Path, blockList.Prefix+block) {
				util.Error("Proxy access to %s is forbidden", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				err := json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"code":    http.StatusForbidden,
					"message": "Forbidden",
				})
				if err != nil {
					return
				}
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Helper function to determine if the request path is blocked
func isPathBlocked(requestPath, blockedPath string) bool {
	// Handle exact match
	if requestPath == blockedPath {
		return true
	}
	// Handle wildcard match (e.g., /admin/* should block /admin and any subpath)
	if strings.HasSuffix(blockedPath, "/*") {
		basePath := strings.TrimSuffix(blockedPath, "/*")
		if strings.HasPrefix(requestPath, basePath) {
			return true
		}
	}
	return false
}

// NewRateLimiter creates a new rate limiter with the specified refill rate and token capacity
func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed based on the current token bucket
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Refill tokens based on the time elapsed
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)
	if tokensToAdd > 0 {
		rl.tokens = min(rl.maxTokens, rl.tokens+tokensToAdd)
		rl.lastRefill = now
	}

	// Check if there are enough tokens to allow the request
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	// Reject request if no tokens are available
	return false
}

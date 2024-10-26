package middleware

import (
	"encoding/base64"
	"encoding/json"
	"github.com/jkaninda/goma/internal/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// RateLimiter defines rate limit properties.
type RateLimiter struct {
	Requests  int
	Window    time.Duration
	ClientMap map[string]*Client
	mu        sync.Mutex
}

// Client stores request count and window expiration for each client.
type Client struct {
	RequestCount int
	ExpiresAt    time.Time
}

// NewRateLimiterWindow creates a new RateLimiter.
func NewRateLimiterWindow(requests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		Requests:  requests,
		Window:    window,
		ClientMap: make(map[string]*Client),
	}
}

type TokenRateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

// ProxyResponseError represents the structure of the JSON error response
type ProxyResponseError struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// AuthJWT  Define struct
type AuthJWT struct {
	AuthURL         string
	RequiredHeaders []string
	Headers         map[string]string
	Params          map[string]string
}

// AuthenticationMiddleware  Define struct
type AuthenticationMiddleware struct {
	AuthURL         string
	RequiredHeaders []string
	Headers         map[string]string
	Params          map[string]string
}
type BlockListMiddleware struct {
	Path        string
	Destination string
	List        []string
}

// AuthBasic  Define Basic auth
type AuthBasic struct {
	Username string
	Password string
	Headers  map[string]string
	Params   map[string]string
}

// AuthMiddleware function, which will be called for each request
func (amw *AuthJWT) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, header := range amw.RequiredHeaders {
			if r.Header.Get(header) == "" {
				logger.Error("Proxy error, missing %s header", header)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				err := json.NewEncoder(w).Encode(ProxyResponseError{
					Message: "Missing Authorization header",
					Code:    http.StatusForbidden,
					Success: false,
				})
				if err != nil {
					return
				}
				return
			}
		}
		//token := r.Header.Get("Authorization")
		authURL, err := url.Parse(amw.AuthURL)
		if err != nil {
			logger.Error("Error parsing auth URL: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(ProxyResponseError{
				Message: "Internal Server Error",
				Code:    http.StatusInternalServerError,
				Success: false,
			})
			if err != nil {
				return
			}
			return
		}
		// Create a new request for /authentication
		authReq, err := http.NewRequest("GET", authURL.String(), nil)
		if err != nil {
			logger.Error("Proxy error creating authentication request: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(ProxyResponseError{
				Message: "Internal Server Error",
				Code:    http.StatusInternalServerError,
				Success: false,
			})
			if err != nil {
				return
			}
			return
		}
		// Copy headers from the original request to the new request
		for name, values := range r.Header {
			for _, value := range values {
				authReq.Header.Set(name, value)
			}
		}
		// Copy cookies from the original request to the new request
		for _, cookie := range r.Cookies() {
			authReq.AddCookie(cookie)
		}
		// Perform the request to the auth service
		client := &http.Client{}
		authResp, err := client.Do(authReq)
		if err != nil || authResp.StatusCode != http.StatusOK {
			logger.Info("%s %s %s %s", r.Method, r.RemoteAddr, r.URL, r.UserAgent())
			logger.Error("Proxy authentication error")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err = json.NewEncoder(w).Encode(ProxyResponseError{
				Message: "Unauthorized",
				Code:    http.StatusUnauthorized,
				Success: false,
			})
			if err != nil {
				return
			}
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(authResp.Body)
		// Inject specific header tp the current request's header
		// Add header to the next request from AuthRequest header, depending on your requirements
		if amw.Headers != nil {
			for k, v := range amw.Headers {
				r.Header.Set(v, authResp.Header.Get(k))
			}
		}
		query := r.URL.Query()
		// Add query parameters to the next request from AuthRequest header, depending on your requirements
		if amw.Params != nil {
			for k, v := range amw.Params {
				query.Set(v, authResp.Header.Get(k))
			}
		}
		r.URL.RawQuery = query.Encode()

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware checks for the Authorization header and verifies the credentials
func (basicAuth AuthBasic) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Error("Proxy error, missing Authorization header")
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(ProxyResponseError{
				Success: false,
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			if err != nil {
				return
			}
			return
		}
		// Check if the Authorization header contains "Basic" scheme
		if !strings.HasPrefix(authHeader, "Basic ") {
			logger.Error("Proxy error, missing Basic Authorization header")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(ProxyResponseError{
				Success: false,
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			if err != nil {
				return
			}
			return
		}

		// Decode the base64 encoded username:password string
		payload, err := base64.StdEncoding.DecodeString(authHeader[len("Basic "):])
		if err != nil {
			logger.Error("Proxy error, missing Basic Authorization header")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(ProxyResponseError{
				Success: false,
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			if err != nil {
				return
			}
			return
		}

		// Split the payload into username and password
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || pair[0] != basicAuth.Username || pair[1] != basicAuth.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(ProxyResponseError{
				Success: false,
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			if err != nil {
				return
			}
			return
		}

		// Continue to the next handler if the authentication is successful
		next.ServeHTTP(w, r)
	})

}

package middleware

import (
	"encoding/json"
	"github.com/jkaninda/goma-gateway/util"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Auth struct {
	Token string
	Url   string
}
type Middleware interface {
	Basic(username, password string) error
	Jwt(token string) error
	Http(url string) error
	Access(url string, code int) error
}
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

type ProxyResponseError struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// AuthenticationMiddleware  Define our struct
type AuthenticationMiddleware struct {
	AuthURL         string
	RequiredHeaders []string
	Headers         map[string]string
	Params          map[string]string
}
type BlockListMiddleware struct {
	Prefix string
	List   []string
}

// AuthMiddleware function, which will be called for each request
func (amw *AuthenticationMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, header := range amw.RequiredHeaders {
			if r.Header.Get(header) == "" {
				util.Error("Proxy error, missing %s header", header)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				err := json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"code":    http.StatusForbidden,
					"message": "Missing Authorization header",
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
			util.Error("Error parsing auth URL: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"code":    http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
			if err != nil {
				return
			}
			return
		}
		// Create a new request for /authentication
		authReq, err := http.NewRequest("GET", authURL.String(), nil)
		if err != nil {
			util.Error("Proxy error creating authentication request: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"code":    http.StatusInternalServerError,
				"message": "Internal Server Error",
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
			util.Info("%s %s %s %s", r.Method, r.RemoteAddr, r.URL, r.UserAgent())
			util.Error("Proxy authentication error")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
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

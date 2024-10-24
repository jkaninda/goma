package middleware

import (
	"encoding/json"
	"github.com/jkaninda/goma-gateway/util"
	"log"
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
	AuthURL string
	Headers map[string]string
	Params  map[string]string
}
type BlockListMiddleware struct {
	Prefix string
	List   []string
}

// AuthMiddleware function, which will be called for each request
func (amw *AuthenticationMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.Info("%s: %s %s", r.RemoteAddr, r.RequestURI, r.UserAgent())
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Println("Missing Authorization header")
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
		authURL, err := url.Parse(amw.AuthURL)
		if err != nil {
			util.Info("Error parsing auth URL: %v", err)
			http.Error(w, "Error parsing auth URL", http.StatusInternalServerError)
			return
		}
		// Create a new request for /authentication
		authReq, err := http.NewRequest("GET", authURL.String(), nil)
		if err != nil {
			util.Info("Error creating auth request: %v", err)
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
			util.Error("Error Auth Response: error: %v", err)
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
		defer authResp.Body.Close()
		util.Info("Successfully authenticated")
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

package middleware

import (
	"github.com/jkaninda/goma-gateway/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
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
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}
type ProxyResponseError struct {
	Success bool          `json:"success"`
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Error   ResponseError `json:"error"`
	Data    any           `json:"data,omitempty"`
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
		utils.Info("%s: %s %s", r.RemoteAddr, r.RequestURI, r.UserAgent())
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Println("Missing Authorization header")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		authURL, err := url.Parse(amw.AuthURL)
		if err != nil {
			utils.Info("Error parsing auth URL: %v", err)
			http.Error(w, "Error parsing auth URL", http.StatusInternalServerError)
			return
		}
		// Create a new request for /authentication
		authReq, err := http.NewRequest("GET", authURL.String(), nil)
		if err != nil {
			utils.Info("Error creating auth request: %v", err)
			http.Error(w, "Error creating auth request", http.StatusInternalServerError)
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
			utils.Error("Error Auth Response: error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		defer authResp.Body.Close()
		utils.Info("Successfully authenticated")
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

// BlocklistMiddleware checks if the request path is forbidden and returns 403 Forbidden
func (blockList BlockListMiddleware) BlocklistMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, block := range blockList.List {
			if isPathBlocked(r.URL.Path, blockList.Prefix+block) {
				http.Error(w, "Forbidden", http.StatusForbidden)
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

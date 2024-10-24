package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/jkaninda/goma-gateway/util"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// ProxyHandler proxies requests to the backend
func ProxyHandler(target, prefix, rewrite string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the target backend URL
		targetURL, err := url.Parse(target)
		if err != nil {
			util.Error("Error parsing backend URL: %s", err)
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
		util.Info("%s %s %s %s", r.Method, r.RemoteAddr, r.URL, r.UserAgent())
		// Create a reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		if prefix != "" && rewrite != "" {
			// Rewrite the path
			if strings.HasPrefix(r.URL.Path, fmt.Sprintf("%s/", prefix)) {
				r.URL.Path = strings.Replace(r.URL.Path, fmt.Sprintf("%s/", prefix), rewrite, 1)
			}
		}
		proxy.ModifyResponse = func(response *http.Response) error {
			if response.StatusCode < 200 || response.StatusCode >= 300 {
				//TODO
			}
			return nil
		}
		// Custom error handler for proxy errors
		proxy.ErrorHandler = ProxyErrorHandler
		proxy.ServeHTTP(w, r)
	}
}

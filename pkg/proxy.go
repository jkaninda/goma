package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/jkaninda/goma-gateway/internal/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// ProxyHandler proxies requests to the backend
func ProxyHandler(path, rewrite, destination string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("%s %s %s %s", r.Method, r.RemoteAddr, r.URL, r.UserAgent())
		// Parse the target backend URL
		targetURL, err := url.Parse(destination)
		if err != nil {
			logger.Error("Error parsing backend URL: %s", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(ErrorResponse{
				Message: "Internal server error",
				Code:    http.StatusInternalServerError,
				Success: false,
			})
			if err != nil {
				return
			}
			return
		}
		// Update the headers to allow for SSL redirection
		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = targetURL.Host
		// Create proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		// Rewrite
		if path != "" && rewrite != "" {
			// Rewrite the path
			if strings.HasPrefix(r.URL.Path, fmt.Sprintf("%s/", path)) {
				r.URL.Path = strings.Replace(r.URL.Path, fmt.Sprintf("%s/", path), rewrite, 1)
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

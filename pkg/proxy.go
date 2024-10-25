package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/jkaninda/goma/internal/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ProxyRoute struct {
	path            string
	rewrite         string
	destination     string
	disableXForward bool
}

// ProxyHandler proxies requests to the backend
func (proxyRoute ProxyRoute) ProxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("%s %s %s %s", r.Method, r.RemoteAddr, r.URL, r.UserAgent())
		// Parse the target backend URL
		targetURL, err := url.Parse(proxyRoute.destination)
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
		if !proxyRoute.disableXForward {
			r.URL.Host = targetURL.Host
			r.URL.Scheme = targetURL.Scheme
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
			r.Header.Set("X-Forwarded-For", r.RemoteAddr)
			r.Host = targetURL.Host
		}
		// Create proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		// Rewrite
		if proxyRoute.path != "" && proxyRoute.rewrite != "" {
			// Rewrite the path
			if strings.HasPrefix(r.URL.Path, fmt.Sprintf("%s/", proxyRoute.path)) {
				r.URL.Path = strings.Replace(r.URL.Path, fmt.Sprintf("%s/", proxyRoute.path), proxyRoute.rewrite, 1)
			}
		}
		proxy.ModifyResponse = func(response *http.Response) error {
			if response.StatusCode < 200 || response.StatusCode >= 300 {
				//TODO || Add override backend errors | user can enable or disable it
			}
			return nil
		}
		// Custom error handler for proxy errors
		proxy.ErrorHandler = ProxyErrorHandler
		proxy.ServeHTTP(w, r)
	}
}

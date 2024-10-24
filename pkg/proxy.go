package pkg

import (
	"fmt"
	"log"
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
			http.Error(w, "Error parsing backend URL", http.StatusInternalServerError)
			return
		}
		// Create a reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		if prefix != "" && rewrite != "" {
			// Rewrite the path
			if strings.HasPrefix(r.URL.Path, fmt.Sprintf("%s/", prefix)) {
				r.URL.Path = strings.Replace(r.URL.Path, fmt.Sprintf("%s/", prefix), rewrite, 1)
			}
		}
		proxy.ModifyResponse = func(response *http.Response) error {
			dumpResponse, err := httputil.DumpResponse(response, false)
			if err != nil {
				return err
			}
			log.Println("Response: \r\n", string(dumpResponse))
			return nil
		}
		// Custom error handler for proxy errors
		proxy.ErrorHandler = ProxyErrorHandler
		proxy.ServeHTTP(w, r)
	}
}

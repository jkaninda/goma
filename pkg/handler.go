package pkg

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// CORSHandler handles CORS headers for incoming requests
//
// Adds CORS headers to the response dynamically based on the provided headers map[string]string
func CORSHandler(headers map[string]string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers from the headers
			if headers != nil {
				for k, v := range headers {
					w.Header().Set(k, v)
				}
			}
			// Handle preflight requests (OPTIONS)
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			// Pass the request to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// ProxyErrorHandler catches backend errors and returns a custom response
func ProxyErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Backend error: %v", err)
	http.Error(w, "The service is currently unavailable. Please try again later.", http.StatusBadGateway)
}

// HealthCheckHandler handles health check requests
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// This is a simple health check. You can include more logic if needed.
	response := HealthCheckResponse{
		Status: "healthy",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HealthCheckResponse represents the health check response structure
type HealthCheckResponse struct {
	Status string `json:"status"`
}

package pkg

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jkaninda/goma-gateway/util"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadGateway)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"code":    http.StatusBadGateway,
		"message": "The service is currently unavailable. Please try again later.",
	})
	if err != nil {
		return
	}
	return
}

// HealthCheckHandler handles health check of routes
func (heathRoute HealthCheckRoute) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	var routes []HealthCheckRouteResponse
	for _, route := range heathRoute.Routes {
		if route.HealthCheck != "" {
			err := HealthCheck(route.Target + route.HealthCheck)
			if err != nil {
				util.Error("Route %s: %s", route.Name, err)
				routes = append(routes, HealthCheckRouteResponse{Name: route.Name, Status: "unhealthy"})
				continue
			} else {
				util.Info("Route %s is healthy", route.Name)
				routes = append(routes, HealthCheckRouteResponse{Name: route.Name, Status: "healthy"})
				continue
			}
		} else {
			util.Error("Route %s's is healthCheck is undefined", route.Name)
			routes = append(routes, HealthCheckRouteResponse{Name: route.Name, Status: "undefined"})
			continue

		}
	}
	response := HealthCheckResponse{
		Status: "healthy",
		Routes: routes,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

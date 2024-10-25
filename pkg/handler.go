package pkg

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jkaninda/goma-gateway/internal/logger"
	"log"
	"net/http"
)

// CORSHandler handles CORS headers for incoming requests
//
// Adds CORS headers to the response dynamically based on the provided headers map[string]string
func CORSHandler(cors map[string]string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers from the cors config
			if cors != nil {
				for k, v := range cors {
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
	log.Printf("Proxy error: %v", err)
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
	logger.Info("%s %s %s %s", r.Method, r.RemoteAddr, r.URL, r.UserAgent())
	var routes []HealthCheckRouteResponse
	for _, route := range heathRoute.Routes {
		if route.HealthCheck != "" {
			err := HealthCheck(route.Destination + route.HealthCheck)
			if err != nil {
				logger.Error("Route %s: %v", route.Name, err)
				if heathRoute.EnableRouteHealthCheckError {
					routes = append(routes, HealthCheckRouteResponse{Name: route.Name, Status: "unhealthy", Error: err.Error()})
					continue

				}
				routes = append(routes, HealthCheckRouteResponse{Name: route.Name, Status: "unhealthy", Error: "Route healthcheck errors disabled"})
				continue
			} else {
				logger.Info("Route %s is healthy", route.Name)
				routes = append(routes, HealthCheckRouteResponse{Name: route.Name, Status: "healthy", Error: ""})
				continue
			}
		} else {
			logger.Error("Route %s's healthCheck is undefined", route.Name)
			routes = append(routes, HealthCheckRouteResponse{Name: route.Name, Status: "undefined", Error: ""})
			continue

		}
	}
	response := HealthCheckResponse{
		Status: "healthy",
		Routes: routes,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

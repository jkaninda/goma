package pkg

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jkaninda/goma-gateway/pkg/middleware"
	"time"
)

func (gatewayServer GatewayServer) Initialize() *mux.Router {
	gateway := gatewayServer.gateway
	r := mux.NewRouter()
	heath := HealthCheckRoute{
		EnableRouteHealthCheckError: gateway.EnableRouteHealthCheckError,
		Routes:                      gateway.Routes,
	}
	// Define the health check route
	r.HandleFunc("/health", heath.HealthCheckHandler).Methods("GET")
	if gateway.RateLimiter != 0 {
		rateLimiter := middleware.NewRateLimiter(gateway.RateLimiter, time.Minute)
		// Add rate limit middleware to all routes, if defined
		r.Use(rateLimiter.RateLimitMiddleware())
	}
	// Add Main route
	for _, route := range gateway.Routes {
		blM := middleware.BlockListMiddleware{
			Prefix: route.Path,
			List:   route.Blocklist,
		}
		// Add block access middleware to all route, if defined
		r.Use(blM.BlocklistMiddleware)
		if route.Middlewares != nil {
			for _, mid := range route.Middlewares {
				secureRouter := r.PathPrefix(route.Path + mid.Path).Subrouter()
				secureRouter.Use(CORSHandler(route.Cors)) // Apply CORS middleware
				if mid.Http.URL != "" {
					amw := middleware.AuthenticationMiddleware{
						AuthURL:         mid.Http.URL,
						RequiredHeaders: mid.Http.RequiredHeaders,
						Headers:         mid.Http.Headers,
						Params:          mid.Http.Params,
					}
					// Apply authentication middleware
					secureRouter.Use(amw.AuthMiddleware)
					secureRouter.PathPrefix("/").Handler(ProxyHandler(route.Path, route.Rewrite, route.Destination)) // Proxy handler
					secureRouter.PathPrefix("").Handler(ProxyHandler(route.Path, route.Rewrite, route.Destination))  // Proxy handler
				} else {
					if mid.Basic.Username != "" {
						amw := middleware.BasicAuth{
							Username: mid.Basic.Username,
							Password: mid.Basic.Password,
						}
						// Apply basic authentication middleware
						secureRouter.Use(amw.BasicAuthMiddleware())
						secureRouter.PathPrefix("/").Handler(ProxyHandler(route.Path, route.Rewrite, route.Destination)) // Proxy handler
						secureRouter.PathPrefix("").Handler(ProxyHandler(route.Path, route.Rewrite, route.Destination))  // Proxy handler
					}
				}

			}
		}
		router := r.PathPrefix(route.Path).Subrouter()
		router.Use(CORSHandler(gateway.Cors)) // Apply CORS middleware
		router.PathPrefix("/").Handler(ProxyHandler(route.Path, route.Rewrite, route.Destination))

	}
	return r

}

func printRoute(routes []Route) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Name", "Route", "Rewrite", "Destination"})
	for _, route := range routes {
		t.AppendRow(table.Row{route.Name, route.Path, route.Rewrite, route.Destination})
	}
	fmt.Println(t.Render())
}

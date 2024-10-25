package pkg

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jkaninda/goma/pkg/middleware"
	"time"
)

func (gatewayServer GatewayServer) Initialize() *mux.Router {
	gateway := gatewayServer.gateway
	r := mux.NewRouter()
	heath := HealthCheckRoute{
		DisableRouteHealthCheckError: gateway.DisableRouteHealthCheckError,
		Routes:                       gateway.Routes,
	}
	// Define the health check route
	r.HandleFunc("/health", heath.HealthCheckHandler).Methods("GET")
	// Apply global Cors middlewares
	r.Use(CORSHandler(gateway.Cors)) // Apply CORS middleware
	if gateway.RateLimiter != 0 {
		rateLimiter := middleware.NewRateLimiter(gateway.RateLimiter, time.Minute)
		// Add rate limit middleware to all routes, if defined
		r.Use(rateLimiter.RateLimitMiddleware())
	}

	// Add Main route
	for _, route := range gateway.Routes {
		blM := middleware.BlockListMiddleware{
			Path: route.Path,
			List: route.Blocklist,
		}
		// Add block access middleware to all route, if defined
		r.Use(blM.BlocklistMiddleware)
		if route.Middlewares != nil {
			for _, mid := range route.Middlewares {
				secureRouter := r.PathPrefix(route.Path + mid.Path).Subrouter()
				if mid.Http.URL != "" {
					amw := middleware.AuthenticationMiddleware{
						AuthURL:         mid.Http.URL,
						RequiredHeaders: mid.Http.RequiredHeaders,
						Headers:         mid.Http.Headers,
						Params:          mid.Http.Params,
					}
					proxyRoute := ProxyRoute{
						path:            route.Path,
						rewrite:         route.Rewrite,
						destination:     route.Destination,
						disableXForward: route.DisableHeaderXForward,
						cors:            route.Cors,
					}
					// Apply JWT authentication middleware
					secureRouter.Use(amw.AuthMiddleware)
					secureRouter.PathPrefix("/").Handler(proxyRoute.ProxyHandler()) // Proxy handler
					secureRouter.PathPrefix("").Handler(proxyRoute.ProxyHandler())  // Proxy handler
				} else {
					if mid.Basic.Username != "" {
						amw := middleware.BasicAuth{
							Username: mid.Basic.Username,
							Password: mid.Basic.Password,
						}
						proxyRoute := ProxyRoute{
							path:            route.Path,
							rewrite:         route.Rewrite,
							destination:     route.Destination,
							disableXForward: route.DisableHeaderXForward,
							cors:            route.Cors,
						}
						// Apply basic authentication middleware
						secureRouter.Use(amw.BasicAuthMiddleware())
						secureRouter.Use(CORSHandler(route.Cors))
						secureRouter.PathPrefix("/").Handler(proxyRoute.ProxyHandler()) // Proxy handler
						secureRouter.PathPrefix("").Handler(proxyRoute.ProxyHandler())  // Proxy handler
					}
				}

			}
		}
		proxyRoute := ProxyRoute{
			path:            route.Path,
			rewrite:         route.Rewrite,
			destination:     route.Destination,
			disableXForward: route.DisableHeaderXForward,
		}

		router := r.PathPrefix(route.Path).Subrouter()
		router.Use(CORSHandler(route.Cors))
		router.PathPrefix("/").Handler(proxyRoute.ProxyHandler())

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

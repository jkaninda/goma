package pkg

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jkaninda/goma/internal/logger"
	"github.com/jkaninda/goma/pkg/middleware"
	"time"
)

func (gatewayServer GatewayServer) Initialize() *mux.Router {
	gateway := gatewayServer.gateway
	//middlewares := gatewayServer.middlewares
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
		//rateLimiter := middleware.NewRateLimiter(gateway.RateLimiter, time.Minute)
		limiter := middleware.NewRateLimiterWindow(gateway.RateLimiter, time.Minute) //  requests per minute
		// Add rate limit middleware to all routes, if defined
		r.Use(limiter.RateLimitMiddleware())
	}
	for _, route := range gateway.Routes {
		blM := middleware.BlockListMiddleware{
			Path: route.Path,
			List: route.Blocklist,
		}
		// Add block access middleware to all route, if defined
		r.Use(blM.BlocklistMiddleware)
		if route.Middlewares != nil {
			for _, mid := range route.Middlewares {
				logger.Info("MiddlewareName %s", mid.Path)
				rules := mid.Rules
				for _, rule := range rules {
					//find rule from middleware lists

					logger.Info("Rule %s", rule)

				}

			}
		}
		proxyRoute := ProxyRoute{
			path:            route.Path,
			rewrite:         route.Rewrite,
			destination:     route.Destination,
			disableXForward: route.DisableHeaderXForward,
			cors:            route.Cors,
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

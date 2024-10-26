package pkg

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jkaninda/goma/internal/logger"
	"github.com/jkaninda/goma/pkg/middleware"
	"github.com/jkaninda/goma/util"
	"time"
)

func (gatewayServer GatewayServer) Initialize() *mux.Router {
	gateway := gatewayServer.gateway
	middlewares := gatewayServer.middlewares
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
		//if route.Middlewares != nil {
		for _, mid := range route.Middlewares {
			secureRouter := r.PathPrefix(util.ParseURLPath(route.Path + mid.Path)).Subrouter()
			proxyRoute := ProxyRoute{
				path:            route.Path,
				rewrite:         route.Rewrite,
				destination:     route.Destination,
				disableXForward: route.DisableHeaderXForward,
				cors:            route.Cors,
			}
			rMiddleware, err := searchMiddleware(mid.Rules, middlewares)
			if err != nil {
				logger.Error("MiddlewareName not found")
			} else {
				switch rMiddleware.Type {
				case "basic":
					basicAuth, err := ToBasicAuth(rMiddleware.Rule)
					if err != nil {

						logger.Error("Error: %s", err.Error())
					} else {
						amw := middleware.AuthBasic{
							Username: basicAuth.Username,
							Password: basicAuth.Password,
							Headers:  nil,
							Params:   nil,
						}
						// Apply JWT authentication middleware
						secureRouter.Use(amw.AuthMiddleware)
						secureRouter.Use(CORSHandler(route.Cors))
						secureRouter.PathPrefix("/").Handler(proxyRoute.ProxyHandler()) // Proxy handler
						secureRouter.PathPrefix("").Handler(proxyRoute.ProxyHandler())  // Proxy handler
					}
				case "jwt":
					jwt, err := ToJWTRuler(rMiddleware.Rule)
					if err != nil {

					} else {
						amw := middleware.AuthJWT{
							AuthURL:         jwt.URL,
							RequiredHeaders: jwt.RequiredHeaders,
							Headers:         jwt.Headers,
							Params:          jwt.Params,
						}
						// Apply JWT authentication middleware
						secureRouter.Use(amw.AuthMiddleware)
						secureRouter.Use(CORSHandler(route.Cors))
						secureRouter.PathPrefix("/").Handler(proxyRoute.ProxyHandler()) // Proxy handler
						secureRouter.PathPrefix("").Handler(proxyRoute.ProxyHandler())  // Proxy handler
					}
				default:
					logger.Error("Unknown middleware type %s", rMiddleware.Type)

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

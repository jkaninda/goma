package pkg

import (
	"github.com/gorilla/mux"
	"github.com/jkaninda/goma-gateway/pkg/middleware"
)

func (gateway Gateway) Initialize() *mux.Router {
	r := mux.NewRouter()
	// Define the health check route
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	// Add Main route
	for _, route := range gateway.Routes {
		blM := middleware.BlockListMiddleware{
			Prefix: route.Path,
			List:   route.Blocklist,
		}
		r.Use(blM.BlocklistMiddleware)
		if route.Middlewares != nil {
			for _, mid := range route.Middlewares {
				//log.Printf("Mapping '%v' | %v ---> %v", route.Name, route.Path+mid.Path, route.Target)
				secureRouter := r.PathPrefix(route.Path + mid.Path).Subrouter()
				secureRouter.Use(CORSHandler(gateway.Headers)) // Apply CORS middleware
				amw := middleware.AuthenticationMiddleware{
					AuthURL: mid.AuthRequest.URL,
					Headers: mid.AuthRequest.Headers,
					Params:  mid.AuthRequest.Params,
				}
				// Apply authentication middleware
				secureRouter.Use(amw.AuthMiddleware)
				secureRouter.PathPrefix("/").Handler(ProxyHandler(route.Target, route.Path, route.Rewrite)) // Proxy handler
				secureRouter.PathPrefix("").Handler(ProxyHandler(route.Target, route.Path, route.Rewrite))  // Proxy handler

			}
		}
		router := r.PathPrefix(route.Path).Subrouter()
		router.Use(CORSHandler(gateway.Headers)) // Apply CORS middleware
		router.PathPrefix("/").Handler(ProxyHandler(route.Target, route.Path, route.Rewrite))
	}
	return r

}

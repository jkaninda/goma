package pkg

import (
	"github.com/jkaninda/goma/internal/logger"
	"net/http"
	"time"
)

func (gatewayServer GatewayServer) Start() {
	logger.Info("Initializing routes...")
	route := gatewayServer.Initialize()
	logger.Info("Initializing routes...done")
	srv := &http.Server{
		Addr:         gatewayServer.gateway.ListenAddr,
		WriteTimeout: time.Second * time.Duration(gatewayServer.gateway.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(gatewayServer.gateway.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(gatewayServer.gateway.IdleTimeout),
		Handler:      route, // Pass our instance of gorilla/mux in.
	}
	if !gatewayServer.gateway.DisableDisplayRouteOnStart {
		printRoute(gatewayServer.gateway.Routes)
	}
	logger.Info("Started Goma Gateway server on %v", gatewayServer.gateway.ListenAddr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("Error starting Goma Gateway server: %v", err)
	}

}
func Stop() {

}

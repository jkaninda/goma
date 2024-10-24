package pkg

import (
	"github.com/jkaninda/goma-gateway/internal/logger"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"time"
)

func Start(cmd *cobra.Command) {
	log.SetOutput(os.Stdout)
	log.Println("Starting Goma Gateway...")
	logger.Info("Initializing routes...")
	configFile, _ := cmd.Flags().GetString("config")
	if configFile == "" {
		configFile = getConfigFile()
	}
	gateway, err := loadConf(configFile)
	if err != nil {
		logger.Fatal("Could not load configuration: %v", err)
	}
	route := gateway.Initialize()
	server := &http.Server{
		Addr:         gateway.ListenAddr,
		WriteTimeout: time.Second * time.Duration(gateway.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(gateway.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(gateway.IdleTimeout),
		Handler:      route, // Pass our instance of gorilla/mux in.
	}
	logger.Info("Initializing routes...done")
	logger.Info("Started Goma Gateway server on %v", gateway.ListenAddr)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Error starting Goma Gateway server: %v", err)
	}

}

package pkg

import (
	"github.com/jkaninda/goma-gateway/util"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"time"
)

func Start(cmd *cobra.Command) {
	log.SetOutput(os.Stdout)
	// Define the internal auth service handler for /authUser
	util.Info("Initializing routes...")
	configFile, _ := cmd.Flags().GetString("config")
	if configFile == "" {
		configFile = ConfigFile
	}
	gateway, err := loadConf(configFile)
	if err != nil {
		util.Fatal("Could not load configuration: %v", err)
	}
	route := gateway.Initialize()
	server := &http.Server{
		Addr:         gateway.ListenAddr,
		WriteTimeout: time.Second * time.Duration(gateway.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(gateway.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(gateway.IdleTimeout),
		Handler:      route, // Pass our instance of gorilla/mux in.
	}
	util.Info("Initializing routes...done")
	util.Info("Started GomaGateway server on %v", gateway.ListenAddr)
	if err := server.ListenAndServe(); err != nil {
		util.Fatal("Error starting GomaGateway server: %v", err)
	}

}

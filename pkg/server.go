package pkg

import (
	"github.com/jkaninda/goma-gateway/utils"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"time"
)

func Start(cmd *cobra.Command) {
	log.SetOutput(os.Stdout)
	// Define the internal auth service handler for /authUser
	utils.Info("Initializing routes...")
	configFile, _ := cmd.Flags().GetString("config")
	if configFile == "" {
		configFile = ConfigFile
	}
	log.Println(configFile)
	gateway, err := loadConf(configFile)
	if err != nil {
		utils.Fatal("Could not load configuration: %v", err)
	}
	route := gateway.Initialize()
	server := &http.Server{
		Addr:         gateway.ListenAddr,
		WriteTimeout: time.Second * time.Duration(gateway.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(gateway.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(gateway.IdleTimeout),
		Handler:      route, // Pass our instance of gorilla/mux in.
	}
	utils.Info("Initializing routes...done")
	utils.Info("Started GomaGateway server on %v", gateway.ListenAddr)
	if err := server.ListenAndServe(); err != nil {
		utils.Fatal("Error starting GomaGateway server: %v", err)
	}

}

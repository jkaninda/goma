// Package cmd /
/*****
@author    Jonas Kaninda
@license   MIT License <https://opensource.org/licenses/MIT>
@Copyright © 2024 Jonas Kaninda
**/
package cmd

import (
	"github.com/jkaninda/goma/internal/logger"
	"github.com/jkaninda/goma/pkg"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Intro()
		configFile, _ := cmd.Flags().GetString("config")
		if configFile == "" {
			configFile = pkg.GetConfigPaths()
		}

		g := pkg.GatewayServer{}
		gs, err := g.New(configFile)
		if err != nil {
			logger.Fatal("Could not load configuration: %v", err)
		}
		gs.Start()

	},
}

func init() {
	ServerCmd.Flags().StringP("config", "", "", "Goma config file")
}

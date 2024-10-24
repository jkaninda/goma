// Package cmd /
/*****
@author    Jonas Kaninda
@license   MIT License <https://opensource.org/licenses/MIT>
@Copyright © 2024 Jonas Kaninda
**/
package cmd

import (
	"github.com/jkaninda/goma-gateway/cmd/config"
	"github.com/jkaninda/goma-gateway/internal/logger"
	"github.com/jkaninda/goma-gateway/util"
	"github.com/spf13/cobra"
)

// rootCmd represents
var rootCmd = &cobra.Command{
	Use:     "goma",
	Short:   "Goma is a lightweight API Gateway, Reverse Proxy",
	Long:    `.`,
	Example: "",
	Version: util.FullVersion(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logger.Fatal("Error executing root command %v", err)
	}
}
func init() {
	rootCmd.AddCommand(ServerCmd)
	rootCmd.AddCommand(StopCmd)
	rootCmd.AddCommand(config.Cmd)

}

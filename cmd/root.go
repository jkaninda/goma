// Package cmd /
/*****
@author    Jonas Kaninda
@license   MIT License <https://opensource.org/licenses/MIT>
@Copyright © 2024 Jonas Kaninda
**/
package cmd

import (
	"github.com/jkaninda/goma-gateway/utils"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents
var rootCmd = &cobra.Command{
	Use:     "goma",
	Short:   "Start goma instance",
	Long:    `.`,
	Example: "",
	Version: utils.FullVersion(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func init() {
	rootCmd.AddCommand(ServerCmd)
	rootCmd.AddCommand(StopCmd)

}

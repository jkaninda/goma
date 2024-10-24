// Package cmd /
/*****
@author    Jonas Kaninda
@license   MIT License <https://opensource.org/licenses/MIT>
@Copyright © 2024 Jonas Kaninda
**/
package cmd

import (
	"github.com/jkaninda/goma-gateway/pkg"
	"github.com/jkaninda/goma-gateway/utils"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			pkg.Start(cmd)
		} else {
			utils.Fatal(`"server" accepts no argument %q`, args)

		}

	},
}

func init() {
	ServerCmd.Flags().StringP("config", "", "", "Goma config file")
}

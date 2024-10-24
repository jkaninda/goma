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

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Start server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			pkg.InitConfig(cmd)
		} else {
			utils.Fatal(`"config" accepts no argument %q`, args)

		}

	},
}

func init() {
	ConfigCmd.Flags().BoolP("init", "", true, "Initialize Goma config")
}

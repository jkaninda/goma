// Package config Package cmd /
/*****
@author    Jonas Kaninda
@license   MIT License <https://opensource.org/licenses/MIT>
@Copyright © 2024 Jonas Kaninda
**/
package config

import (
	"github.com/jkaninda/goma-gateway/internal/logger"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Goma configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		} else {
			logger.Fatal(`"config" accepts no argument %q`, args)

		}

	},
}

func init() {
	Cmd.AddCommand(InitConfigCmd)
}

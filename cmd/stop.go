// Package cmd /
/*****
@author    Jonas Kaninda
@license   MIT License <https://opensource.org/licenses/MIT>
@Copyright © 2024 Jonas Kaninda
**/
package cmd

import (
	"github.com/jkaninda/goma-gateway/internal/logger"
	"github.com/jkaninda/goma-gateway/pkg"
	"github.com/spf13/cobra"
)

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			pkg.Stop()
		} else {
			logger.Fatal(`"stop" accepts no argument %q`, args)

		}

	},
}

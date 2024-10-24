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
)

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			//pkg.Start(cmd)
		} else {
			utils.Fatal(`"migrate" accepts no argument %q`, args)

		}

	},
}

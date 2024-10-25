package config

import (
	"github.com/jkaninda/goma/internal/logger"
	"github.com/jkaninda/goma/pkg"
	"github.com/spf13/cobra"
)

var InitConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Goma configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			pkg.InitConfig(cmd)
		} else {
			logger.Fatal(`"config" accepts no argument %q`, args)
		}

	},
}

func init() {
	InitConfigCmd.Flags().StringP("output", "o", "", "config file output")
}

package config

import (
	"github.com/jkaninda/goma-gateway/pkg"
	"github.com/jkaninda/goma-gateway/util"
	"github.com/spf13/cobra"
)

var InitConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Goma config",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			pkg.InitConfig()
		} else {
			util.Fatal(`"config" accepts no argument %q`, args)
		}

	},
}

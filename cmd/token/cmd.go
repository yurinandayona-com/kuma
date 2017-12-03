package token

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var Cmd *cobra.Command

func init() {
	var cfgFile string

	Cmd = &cobra.Command{
		Use:   "token [inspect|generate]",
		Short: "Manage user tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			return flag.ErrHelp
		},
	}

	Cmd.PersistentFlags().StringVarP(&cfgFile, "config", "C", ".kuma/serve.yml", "specify configuration file location")
}

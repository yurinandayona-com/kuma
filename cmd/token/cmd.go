package token

import (
	"github.com/spf13/cobra"
	"github.com/yurinandayona-com/kuma/cmd/serve"
	"github.com/yurinandayona-com/kuma/cmd/token/inspect"
)

var Cmd *cobra.Command

func init() {
	Cmd = &cobra.Command{
		Use:   "token",
		Short: "Manage user tokens",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			serve.BindToStore(cmd.Flags())
		},
	}

	Cmd.PersistentFlags().StringP("config", "C", ".kuma/serve.toml", "configuration file location")
	serve.AddFlags(Cmd.PersistentFlags())

	Cmd.AddCommand(token_inspect.Cmd)
}

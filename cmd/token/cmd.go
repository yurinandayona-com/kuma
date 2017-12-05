// Package token implements 'kuma token'.
package token

import (
	"github.com/spf13/cobra"
	"github.com/yurinandayona-com/kuma/cmd/token/config"
	"github.com/yurinandayona-com/kuma/cmd/token/inspect"
	"github.com/yurinandayona-com/kuma/cmd/token/sign"
)

// 'kuma token' command.
var Cmd *cobra.Command

func init() {
	var cfgFile string

	Cmd = &cobra.Command{
		Use:   "token",
		Short: "Manage user tokens",
	}

	Cmd.PersistentFlags().StringVarP(&cfgFile, "config", "C", ".kuma/serve.toml", "configuration file location")

	Cmd.PersistentFlags().StringP("user-db", "u", "", "user DB location")

	Cmd.AddCommand(inspect.Cmd)
	Cmd.AddCommand(sign.Cmd)

	config.Store.BindPFlag("user_db", Cmd.PersistentFlags().Lookup("user-db"))
}

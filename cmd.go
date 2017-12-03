package main

import (
	"github.com/comail/colog"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/yurinandayona-com/kuma/cmd/connect"
	"github.com/yurinandayona-com/kuma/cmd/serve"
	"github.com/yurinandayona-com/kuma/cmd/token"
	"github.com/yurinandayona-com/kuma/version"
	"log"
)

var Cmd *cobra.Command

func init() {
	var logLevel string

	Cmd = &cobra.Command{
		Use:     "kuma",
		Short:   "kuma: HTTP Tunnel over gRPC",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return flag.ErrHelp
		},
	}

	Cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "minimal log level")

	Cmd.AddCommand(connect.Cmd)
	Cmd.AddCommand(serve.Cmd)
	Cmd.AddCommand(token.Cmd)

	cobra.OnInitialize(func() {
		colog.Register()

		level, err := colog.ParseLevel(logLevel)
		if err != nil {
			log.Fatalf("alert: %s", err)
		}

		colog.SetMinLevel(level)
	})
}

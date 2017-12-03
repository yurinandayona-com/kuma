package main

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/yurinandayona-com/kuma/cmd/connect"
	"github.com/yurinandayona-com/kuma/cmd/serve"
	"github.com/yurinandayona-com/kuma/cmd/token"
	"github.com/yurinandayona-com/kuma/version"
)

var Cmd *cobra.Command

func init() {
	Cmd = &cobra.Command{
		Use:     "kuma",
		Short:   "kuma: HTTP Tunnel over gRPC",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return flag.ErrHelp
		},
	}

	Cmd.AddCommand(connect.Cmd)
	Cmd.AddCommand(serve.Cmd)
	Cmd.AddCommand(token.Cmd)
}

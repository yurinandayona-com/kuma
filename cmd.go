package main

import (
	"github.com/spf13/cobra"
	"github.com/yurinandayona-com/kuma/cmd/connect"
	"github.com/yurinandayona-com/kuma/cmd/generate_tokens"
	"github.com/yurinandayona-com/kuma/cmd/serve"
)

var Cmd *cobra.Command

func init() {
	Cmd = &cobra.Command{
		Use:   "kuma",
		Short: "kuma: HTTP Tunnel over gRPC",
	}

	Cmd.AddCommand(connect.Cmd)
	Cmd.AddCommand(generate_tokens.Cmd)
	Cmd.AddCommand(serve.Cmd)
}

// kuma: HTTP Tunnel over gRPC.
//
//     kuma: HTTP Tunnel over gRPC
//
//     Usage:
//       kuma [flags]
//       kuma [command]
//
//     Available Commands:
//       connect     Connect to kuma gRPC server
//       help        Help about any command
//       serve       Serve an HTTP server and gRPC server for kuma
//       token       Manage user tokens
//
//     Flags:
//       -h, --help               help for kuma
//       -l, --log-level string   minimal log level (default "info")
//
//     Use "kuma [command] --help" for more information about a command.
package main

import (
	"os"
)

func main() {
	if err := Cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

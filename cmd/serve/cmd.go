package serve

import (
	"github.com/spf13/cobra"
	"log"
)

var Cmd *cobra.Command

func init() {
	var config string

	Cmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve an HTTP server and gRPC server for kuma",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("info: load config: %s", config)
			cfg, err := LoadConfig(config)
			if err != nil {
				log.Fatalf("alert: %s", err)
				return
			}

			log.Print("info: run 'serve' command")
			runner := &Runner{Config: cfg}
			if err := runner.Run(); err != nil {
				log.Fatalf("alert: %s", err)
			}
		},
	}

	Cmd.Flags().StringVarP(&config, "config", "C", ".kuma/serve.toml", "specify configuration file")
}

package connect

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yurinandayona-com/kuma/client"
	"github.com/yurinandayona-com/kuma/config"
	"log"
)

var (
	Cmd   *cobra.Command
	Store = viper.New()
)

func init() {
	var cfgFile string

	Cmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to kuma gRPC server",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("info: load config: %s", cfgFile)
			var cfg Config
			err := config.Load(Store, cfgFile, &cfg)
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			cfg.DebugLog()

			cli := &client.Client{
				GRPCServer: cfg.GRPCServer,
				UseTLS:     cfg.UseTLS,
				Token:      cfg.Token,

				Port:      cfg.Port,
				Subdomain: cfg.Subdomain,
			}

			log.Print("info: run 'connect' command")
			if err := cli.Start(); err != nil {
				log.Fatalf("alert: %s", err)
			}
		},
	}

	//
	// Flag settings

	// Flag for configuration file
	Cmd.Flags().StringVarP(&cfgFile, "config", "C", ".kuma/connect.toml", "specify configuration file")

	// Flags for gRPC server
	Cmd.Flags().String("grpc-server", "", "gRPC server address to connect")
	Cmd.Flags().BoolP("use-tls", "T", false, "use TLS to connect gRPC server")
	Cmd.Flags().StringP("token", "t", "", "specify user token")

	// Flags for local
	Cmd.Flags().IntP("port", "p", 0, "specify localhost port number to connect")
	Cmd.Flags().StringP("subdomain", "S", "", "specify public URL subdomain")

	//
	// Viper settings

	Store.SetDefault("grpc_server", "yurinandayona.com:8433")
	Store.SetDefault("use_tls", true)

	Store.SetEnvPrefix("kuma_connect")

	Store.BindPFlag("grpc_server", Cmd.Flags().Lookup("grpc-server"))
	Store.BindPFlag("use_tls", Cmd.Flags().Lookup("use-tls"))
	Store.BindPFlag("token", Cmd.Flags().Lookup("token"))
	Store.BindPFlag("port", Cmd.Flags().Lookup("port"))
	Store.BindPFlag("subdomain", Cmd.Flags().Lookup("subdomain"))
}

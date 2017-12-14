// Package connect implements 'kuma connect'.
package connect

import (
	"log"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yurinandayona-com/kuma/client"
	"github.com/yurinandayona-com/kuma/config"
	"golang.org/x/net/context"
)

var (
	// Cmd is 'kuma connect' command.
	Cmd *cobra.Command

	// Store is 'kuma connect' configuration store.
	Store = viper.New()
)

func init() {
	var cfgFile string

	Cmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to kuma gRPC server",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("debug: load configuration: %s", cfgFile)
			var cfg Config
			BindToStore(cmd.Flags())
			err := config.Load(Store, cfgFile, &cfg)
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			cfg.DebugLog()

			cli := &client.Client{
				GRPCServer: cfg.GRPCServer,

				UseTLS:  cfg.UseTLS,
				RootCA:  cfg.RootCA,
				TLSCert: cfg.TLSCert,
				TLSKey:  cfg.TLSKey,

				Port:      cfg.Port,
				Subdomain: cfg.Subdomain,
			}

			log.Print("debug: run 'connect' command")
			ctx := context.Background()
			if err := cli.Run(ctx); err != nil {
				log.Fatalf("alert: %s", err)
			}
		},
	}

	//
	// Flag settings

	// Flag for configuration file
	Cmd.Flags().StringVarP(&cfgFile, "config", "C", ".kuma/connect.toml", "configuration file location")

	AddFlags(Cmd.Flags())

	//
	// Viper settings

	Store.SetDefault("grpc_server", "yurinandayona.com:8433")
	Store.SetDefault("use_tls", true)

	Store.SetEnvPrefix("kuma_connect")
}

// AddFlags sets up Config related flags.
func AddFlags(flags *flag.FlagSet) {
	flags.String("grpc-server", "", "gRPC server address to connect")
	flags.BoolP("use-tls", "T", false, "flag to use TLS to connect gRPC server")
	flags.String("tls-cert", "", "TLS certification file location")
	flags.String("tls-key", "", "TLS key file location")
	flags.String("root-ca", "", "TLS root CA file location")
	flags.IntP("port", "p", 0, "localhost port number to proxy")
	flags.StringP("subdomain", "S", "", "subdomain name to require")
}

// BindToStore binds flags to Store. It should be called before config.Load
// against *Config.
func BindToStore(flags *flag.FlagSet) {
	Store.BindPFlag("grpc_server", flags.Lookup("grpc-server"))
	Store.BindPFlag("use_tls", flags.Lookup("use-tls"))
	Store.BindPFlag("tls_cert", flags.Lookup("tls-cert"))
	Store.BindPFlag("tls_key", flags.Lookup("tls-key"))
	Store.BindPFlag("root_ca", flags.Lookup("root-ca"))
	Store.BindPFlag("port", flags.Lookup("port"))
	Store.BindPFlag("subdomain", flags.Lookup("subdomain"))
}

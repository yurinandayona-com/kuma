// Package serve implements 'kuma serve'.
package serve

import (
	"log"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yurinandayona-com/kuma/config"
)

var (
	// Cmd is 'kuma serve' command.
	Cmd *cobra.Command

	// Store is 'kuma store' configuration store.
	Store = viper.New()
)

func init() {
	var cfgFile string

	Cmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve an HTTP server and gRPC server for kuma",
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

			log.Print("debug: run 'serve' command")
			r := &runner{Config: &cfg}
			if err := r.Run(); err != nil {
				log.Fatalf("alert: %s", err)
			}
		},
	}

	//
	// Flag settings

	// Flag for configuration file
	Cmd.Flags().StringVarP(&cfgFile, "config", "C", ".kuma/serve.toml", "configuration file location")

	AddFlags(Cmd.Flags())

	//
	// Viper settings

	Store.SetDefault("http.listen", ":10080")
	Store.SetDefault("grpc.listen", ":8342")
	Store.SetDefault("grpc.use_tls", false)

	Store.SetEnvPrefix("kuma_serve")
	Store.BindEnv("http.listen", "KUMA_SERVE_HTTP_LISTEN")
	Store.BindEnv("grpc.listen", "KUMA_SERVE_GRPC_LISTEN")
	Store.BindEnv("grpc.use_tls", "KUMA_SERVE_GRPC_USE_TLS")
	Store.BindEnv("grpc.tls_cert", "KUMA_SERVE_GRPC_TLS_CERT")
	Store.BindEnv("grpc.tls_key", "KUMA_SERVE_GRPC_TLS_KEY")
}

// AddFlags sets up Config related flags.
func AddFlags(flags *flag.FlagSet) {
	// Flags for general settings
	flags.StringP("base-domain", "b", "", "base domain of server")

	// Flags for HTTP server
	flags.String("http-listen", "", "address to listen HTTP server")

	// Flags for gRPC server
	flags.String("grpc-listen", "", "address to listen gRPC server")
	flags.Bool("grpc-use-tls", false, "flag to use TLS for gRPC server")
	flags.String("grpc-tls-cert", "", "TLS certification file location")
	flags.String("grpc-tls-key", "", "TLS key file location")
	flags.String("grpc-client-ca", "", "TLS client CA file location")
}

// BindToStore binds flags to Store. It should be called before config.Load
// against *Config.
func BindToStore(flags *flag.FlagSet) {
	Store.BindPFlag("base_domain", flags.Lookup("base-domain"))

	Store.BindPFlag("http.listen", flags.Lookup("http-listen"))

	Store.BindPFlag("grpc.listen", flags.Lookup("grpc-listen"))
	Store.BindPFlag("grpc.use_tls", flags.Lookup("grpc-use-tls"))
	Store.BindPFlag("grpc.tls_cert", flags.Lookup("grpc-tls-cert"))
	Store.BindPFlag("grpc.tls_key", flags.Lookup("grpc-tls-key"))
	Store.BindPFlag("grpc.client_ca", flags.Lookup("grpc-client-ca"))
}

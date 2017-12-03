package serve

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Use:   "serve",
		Short: "Serve an HTTP server and gRPC server for kuma",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("info: load config: %s", cfgFile)
			var cfg Config
			err := config.Load(Store, cfgFile, &cfg)
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			cfg.DebugLog()

			log.Print("info: run 'serve' command")
			runner := &Runner{Config: &cfg}
			if err := runner.Run(); err != nil {
				log.Fatalf("alert: %s", err)
			}
		},
	}

	//
	// Flag settings

	// Flag for configuration file
	Cmd.Flags().StringVarP(&cfgFile, "config", "C", ".kuma/serve.toml", "configuration file location")

	// Flags for general settings
	Cmd.Flags().StringP("user-db", "u", "", "specify user DB location")
	Cmd.Flags().StringP("base-domain", "b", "", "specify base domain of server")

	// Flags for HTTP server
	Cmd.Flags().String("http-listen", "", "specify address to listen HTTP server")

	// Flags for gRPC server
	Cmd.Flags().String("grpc-listen", "", "specify address to listen gRPC server")
	Cmd.Flags().Bool("grpc-use-tls", false, "use TLS for gRPC server")
	Cmd.Flags().String("grpc-tls-cert", "", "specify TLS certification file location")
	Cmd.Flags().String("grpc-tls-key", "", "specify TLS key file location")

	//
	// Viper settings

	Store.SetDefault("user_db", ".kuma/user_db.toml")
	Store.SetDefault("http.listen", ":10080")
	Store.SetDefault("grpc.listen", ":8342")
	Store.SetDefault("grpc.use_tls", false)

	Store.SetEnvPrefix("kuma_serve")
	Store.BindEnv("http.listen", "KUMA_SERVE_HTTP_LISTEN")
	Store.BindEnv("grpc.listen", "KUMA_SERVE_GRPC_LISTEN")
	Store.BindEnv("grpc.use_tls", "KUMA_SERVE_GRPC_USE_TLS")
	Store.BindEnv("grpc.tls_cert", "KUMA_SERVE_GRPC_TLS_CERT")
	Store.BindEnv("grpc.tls_key", "KUMA_SERVE_GRPC_TLS_KEY")

	Store.BindPFlag("user_db", Cmd.Flags().Lookup("user-db"))
	Store.BindPFlag("base_domain", Cmd.Flags().Lookup("base-domain"))

	Store.BindPFlag("http.listen", Cmd.Flags().Lookup("http-listen"))

	Store.BindPFlag("grpc.listen", Cmd.Flags().Lookup("grpc-listen"))
	Store.BindPFlag("grpc.use_tls", Cmd.Flags().Lookup("grpc-use-tls"))
	Store.BindPFlag("grpc.tls_cert", Cmd.Flags().Lookup("grpc-tls-cert"))
	Store.BindPFlag("grpc.tls_key", Cmd.Flags().Lookup("grpc-tls-key"))
}

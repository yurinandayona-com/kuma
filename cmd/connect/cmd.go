package connect

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
	"github.com/yurinandayona-com/kuma/client"
	"log"
)

var (
	Cmd   *cobra.Command
	Store = viper.New()
)

func init() {
	var config string

	Cmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to kuma gRPC server",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("info: load config: %s", config)
			cfg, err := InitConfig(config)
			if err != nil {
				if es, ok := err.(validator.ValidationErrors); ok {
					for _, e := range es {
						log.Printf("alert: %s", e)
					}
					log.Fatal("alert: kuma: failed to load config")
				} else {
					log.Fatalf("alert: %s", err)
				}
			}

			cfg.DebugLog()

			cli := &client.Client{
				GRPCServer: cfg.GRPCServer,
				UseTLS: cfg.UseTLS,
				Token: cfg.Token,

				Port: cfg.Port,
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
	Cmd.Flags().StringVarP(&config, "config", "C", ".kuma/connect.toml", "specify configuration file")

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

func InitStore(config string) error {
	Store.SetConfigFile(config)
	Store.AutomaticEnv()

	if err := Store.ReadInConfig(); err != nil {
		return errors.Wrap(err, "kuma: read config")
	}

	return nil
}

func InitConfig(config string) (*Config, error) {
	if err := InitStore(config); err != nil {
		return nil, err
	}

	var cfg Config
	if err := Store.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "kuma: unmarshal config")
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

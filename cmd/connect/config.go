package connect

import (
	"log"
)

type Config struct {
	// gRPC server configuration
	GRPCServer string `mapstructure:"grpc_server"`
	UseTLS     bool   `mapstructure:"use_tls"`
	Token      string `mapstructure:"token" validate:"required"`

	// local tunnel configuration
	Port      int    `mapstructure:"port" validate:"required"`
	Subdomain string `mapstructure:"subdomain" validate:"required,alphanum,max=17"`
}

func (cfg *Config) DebugLog() {
	log.Printf("debug: grpc_server = %#v", cfg.GRPCServer)
	log.Printf("debug: use_tls = %#v", cfg.UseTLS)
	log.Printf("debug: token = %#v", cfg.Token)
	log.Printf("debug: port = %#v", cfg.Port)
	log.Printf("debug: subdomain = %#v", cfg.Subdomain)
}

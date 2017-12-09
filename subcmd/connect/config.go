package connect

import (
	"log"
)

// Config represents configuration file contents for 'kuma connect'.
type Config struct {
	// GRPCServer is gRPC server address to connect.
	GRPCServer string `mapstructure:"grpc_server"`

	// UseTLS is flag to use TLS to connect gRPC server.
	UseTLS     bool   `mapstructure:"use_tls"`

	// Port is localhost port number to proxy.
	Port      int    `mapstructure:"port" validate:"required"`

	// Subdomain is subdomain name to require.
	Subdomain string `mapstructure:"subdomain" validate:"required,alphanum,max=17"`
}

// DebugLog inspects configuration contents as debug logs.
func (cfg *Config) DebugLog() {
	log.Printf("debug: grpc_server = %#v", cfg.GRPCServer)
	log.Printf("debug: use_tls = %#v", cfg.UseTLS)
	log.Printf("debug: port = %#v", cfg.Port)
	log.Printf("debug: subdomain = %#v", cfg.Subdomain)
}

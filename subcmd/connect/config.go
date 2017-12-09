package connect

import (
	"log"
)

// Config represents configuration file contents for 'kuma connect'.
type Config struct {
	// GRPCServer is gRPC server address to connect.
	GRPCServer string `mapstructure:"grpc_server"`

	// UseTLS is flag to use TLS to connect gRPC server.
	UseTLS bool `mapstructure:"use_tls"`

	// RootCA is TLS root CA file location to connect gRPC server.
	// If RootCA is zero value, it uses system root CA.
	RootCA string `mapstructure:"root_ca"`

	// TLSCert and TLSKey are TLS client certificate key-pair location
	// to connect gRPC server.
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`

	// Port is localhost port number to proxy.
	Port int `mapstructure:"port" validate:"required"`

	// Subdomain is subdomain name to require.
	Subdomain string `mapstructure:"subdomain" validate:"required,alphanum,max=17"`
}

// DebugLog inspects configuration contents as debug logs.
func (cfg *Config) DebugLog() {
	log.Printf("debug: grpc_server = %#v", cfg.GRPCServer)
	log.Printf("debug: use_tls = %#v", cfg.UseTLS)
	log.Printf("debug: root_ca = %#v", cfg.RootCA)
	log.Printf("debug: tls_cert = %#v", cfg.TLSCert)
	log.Printf("debug: tls_key = %#v", cfg.TLSKey)
	log.Printf("debug: port = %#v", cfg.Port)
	log.Printf("debug: subdomain = %#v", cfg.Subdomain)
}

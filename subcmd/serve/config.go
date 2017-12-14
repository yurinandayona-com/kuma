package serve

import (
	"log"
)

// Config represents configuration file contents for 'kuma serve'.
type Config struct {
	// HashIDSecret is salt string for HashIDs.
	HashIDSecret string `mapstructure:"hash_id_secret" validate:"required"`

	// Base domain.
	BaseDomain string `mapstructure:"base_domain" validate:"required"`

	// For HTTP related configurations.
	HTTP struct {
		// Address to listen on HTTP server.
		Listen string `mapstructure:"listen"`
	} `mapstructure:"http"`

	// For gRPC related configurations.
	GRPC struct {
		// Address to listen on gRPC server.
		Listen string `mapstructure:"listen"`

		// Flag to determine using TLS.
		UseTLS bool `mapstructure:"use_tls"`

		// Certification file and private key file.
		TLSCert string `mapstructure:"tls_cert"`
		TLSKey  string `mapstructure:"tls_key"`

		// TLS CA location for client certificate authentication.
		TLSClientCA string `mapstructure:"tls_client_ca"`
	} `mapstructure:"grpc"`
}

// DebugLog inspects configuration contents as debug logs.
//
// It does not show HashIDSecret for security reason.
func (cfg *Config) DebugLog() {
	log.Print("debug: hash_id_secret = *** (hide)")
	log.Printf("debug: base_domain = %#v", cfg.BaseDomain)

	log.Printf("debug: http.listen = %#v", cfg.HTTP.Listen)

	log.Printf("debug: grpc.listen = %#v", cfg.GRPC.Listen)
	log.Printf("debug: grpc.use_tls = %#v", cfg.GRPC.UseTLS)
	log.Printf("debug: grpc.tls_cert = %#v", cfg.GRPC.TLSCert)
	log.Printf("debug: grpc.tls_key = %#v", cfg.GRPC.TLSKey)
	log.Printf("debug: grpc.client_ca = %#v", cfg.GRPC.TLSClientCA)
}

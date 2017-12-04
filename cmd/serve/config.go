package serve

import (
	"log"
)

// Config represents configuration file contents for 'kuma serve'.
type Config struct {
	// User DB location.
	UserDB string `mapstructure:"user_db"`

	// Salt string for HashIDs.
	HashIDsSalt string `mapstructure:"hash_ids_salt" validate:"required"`

	// Private key for HMAC.
	HMACKey string `mapstructure:"hmac_key" validate:"required"`

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
	} `mapstructure:"grpc"`
}

// DebugLog inspects configuration contents as debug logs.
//
// It does not show HashIDsSalt and HMACKey for security reason.
func (cfg *Config) DebugLog() {
	log.Printf("debug: user_db = %#v", cfg.UserDB)
	log.Print("debug: hash_ids_salt = *** (hide)")
	log.Print("debug: hmac_key = *** (hide)")
	log.Printf("debug: base_domain = %#v", cfg.BaseDomain)

	log.Printf("debug: http.listen = %#v", cfg.HTTP.Listen)

	log.Printf("debug: grpc.listen = %#v", cfg.GRPC.Listen)
	log.Printf("debug: grpc.use_tls = %#v", cfg.GRPC.UseTLS)
	log.Printf("debug: grpc.tls_cert = %#v", cfg.GRPC.TLSCert)
	log.Printf("debug: grpc.tls_key = %#v", cfg.GRPC.TLSKey)
}

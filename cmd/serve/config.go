package serve

import (
	"log"
)

type Config struct {
	UserDB string `mapstructure:"user_db"`

	HashIDsSalt string `mapstructure:"hash_ids_salt" validate:"required"`

	HMACKey string `mapstructure:"hmac_key" validate:"required"`

	BaseDomain string `mapstructure:"base_domain" validate:"required"`

	HTTP struct {
		Listen string `mapstructure:"listen"`
	} `mapstructure:"http"`

	GRPC struct {
		Listen string `mapstructure:"listen"`

		UseTLS bool `mapstructure:"use_tls"`

		TLSCert string `mapstructure:"tls_cert"`
		TLSKey  string `mapstructure:"tls_key"`
	} `mapstructure:"grpc"`
}

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

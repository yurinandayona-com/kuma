package serve

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"io/ioutil"
)

type Config struct {
	UserDB string `toml:"user_db"`

	HashIDsSalt string `toml:"hash_ids_salt"`

	HMACKey string `toml:"hmac_key"`

	BaseDomain string `toml:"base_domain"`

	HTTP struct {
		Listen string `toml:"listen"`
	} `toml:"http"`

	GRPC struct {
		Listen string `toml:"listen"`

		UseTLS bool `toml:"use_tls"`

		TLSCert string `toml:"tls_cert"`
		TLSKey  string `toml:"tls_key"`
	} `toml:"grpc"`
}

func LoadConfig(filename string) (*Config, error) {
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "kuma: failed to load config file")
	}

	var cfg Config
	err = toml.Unmarshal(p, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "kuma: failed to load config TOML")
	}

	return &cfg, nil
}

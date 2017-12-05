// Package token_config provides configuration loader for commands under 'kuma token'.
package token_config

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yurinandayona-com/kuma/config"
	"github.com/yurinandayona-com/kuma/userdb"
)

var (
	// Store is 'kuma token' configuration store.
	Store = viper.New()
)

func init() {
	Store.SetDefault("user_db", ".kuma/user_db.toml")
	Store.SetEnvPrefix("KUMA_SERVE")
}

// LoadJWTManager returns JWTManager created from `kuma token` configuration.
func LoadJWTManager(cmd *cobra.Command) (*userdb.JWTManager, error) {
	cfg, err := loadTokenConfig(cmd)
	if err != nil {
		return nil, err
	}

	cfg.DebugLog()

	log.Printf("debug: load user DB: %s", cfg.UserDB)
	userDB, err := userdb.LoadUserDB(cfg.UserDB)
	if err != nil {
		return nil, err
	}

	jm := &userdb.JWTManager{
		UserDB:  userDB,
		HMACKey: []byte(cfg.HMACKey),
	}

	return jm, nil
}

type tokenConfig struct {
	UserDB  string `mapstructure:"user_db"`
	HMACKey string `mapstructure:"hmac_key"`
}

func loadTokenConfig(cmd *cobra.Command) (*tokenConfig, error) {
	cfgFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, errors.Wrap(err, "kuma: failed to get '--config' value")
	}

	log.Printf("debug: load condiguration: %s", cfgFile)
	var cfg tokenConfig
	if err := config.Load(Store, cfgFile, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *tokenConfig) DebugLog() {
	log.Printf("debug: user_db = %#v", cfg.UserDB)
	log.Print("debug: hmac_key = *** (hide)")
}

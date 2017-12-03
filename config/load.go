package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

// Load loads cfgFile via viper, then marshal cfg and also validate it.
func Load(store *viper.Viper, cfgFile string, cfg interface{}) error {
	if err := loadStore(store, cfgFile); err != nil {
		return err
	}

	if err := store.Unmarshal(cfg); err != nil {
		return errors.Wrap(err, "kuma: unmarshal config")
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return err
	}

	return nil
}

// loadStore loads cfgFile via viper.
func loadStore(store *viper.Viper, cfgFile string) error {
	store.SetConfigFile(cfgFile)
	store.AutomaticEnv()

	if err := store.ReadInConfig(); err != nil {
		return errors.Wrap(err, "kuma: read config")
	}

	return nil
}

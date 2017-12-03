// Package token_sign provides 'kuma token sign' implementation.
package token_sign

import (
	"fmt"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/yurinandayona-com/kuma/cmd/serve"
	"github.com/yurinandayona-com/kuma/config"
	"github.com/yurinandayona-com/kuma/user_db"
	"log"
	"time"
)

const (
	defaultExpirationDays = 100
)

var Cmd *cobra.Command

func init() {
	var name string
	var expirationDays int

	Cmd = &cobra.Command{
		Use:   "sign",
		Short: "Generate a new user token and sign it",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return flag.ErrHelp
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cfgFile, err := cmd.Flags().GetString("config")
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			log.Printf("debug: load condiguration: %s", cfgFile)
			var cfg serve.Config
			if err := config.Load(serve.Store, cfgFile, &cfg); err != nil {
				log.Fatalf("alert: %s", err)
			}

			cfg.DebugLog()

			log.Printf("debug: load user DB: %s", cfg.UserDB)
			userDB, err := user_db.LoadUserDB(cfg.UserDB)
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			jm := &user_db.JWTManager{
				UserDB:  userDB,
				HMACKey: []byte(cfg.HMACKey),
			}

			var user *user_db.User
			for _, u := range userDB.GetUsers() {
				if u.Name == name {
					user = u
					break
				}
			}
			if user == nil {
				log.Fatalf("alert: user not found: %s", name)
			}

			log.Print("debug: user.name = %#v", user.Name)
			log.Print("debug: user.id = %#v", user.ID)
			log.Print("debug: expiration_days = %d days", expirationDays)

			expiration := time.Duration(expirationDays) * 24 * time.Hour
			signed, err := jm.Sign(user, time.Now().Add(expiration))
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			fmt.Println(signed)
		},
	}

	Cmd.Flags().StringVarP(&name, "name", "n", "", "user name (required)")
	Cmd.Flags().IntVarP(&expirationDays, "expiration-days", "e", defaultExpirationDays, "expiration days of generated user token")
}

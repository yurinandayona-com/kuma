package token_inspect

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/yurinandayona-com/kuma/cmd/serve"
	"github.com/yurinandayona-com/kuma/config"
	"github.com/yurinandayona-com/kuma/user_db"
	"log"
	"time"
)

var Cmd *cobra.Command

func init() {
	var token string

	Cmd = &cobra.Command{
		Use:   "inspect",
		Short: "Inspect the user token",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if token == "" {
				return flag.ErrHelp
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cfgFile, err := cmd.Flags().GetString("config")
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			log.Printf("info: load condiguration: %s", cfgFile)
			var cfg serve.Config
			if err := config.Load(serve.Store, cfgFile, &cfg); err != nil {
				log.Fatalf("alert: %s", err)
			}

			cfg.DebugLog()

			log.Printf("info: load user DB: %s", cfg.UserDB)
			userDB, err := user_db.LoadUserDB(cfg.UserDB)
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			jm := &user_db.JWTManager{
				UserDB:  userDB,
				HMACKey: []byte(cfg.HMACKey),
			}

			claims, valid, err := jm.Parse(token)
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			if valid {
				_, err := userDB.Verify(claims.ID, claims.Name)
				valid = err == nil
			}

			expire := time.Unix(claims.ExpiresAt, 0)

			fmt.Println()
			fmt.Printf("ID     : %s\n", claims.ID)
			fmt.Printf("Name   : %s\n", claims.Name)
			fmt.Printf("Expire : %s (%s)\n", expire, humanize.Time(expire))
			fmt.Printf("Valid  : %t\n", valid)
		},
	}

	Cmd.Flags().StringVarP(&token, "token", "t", "", "user token (required)")
}

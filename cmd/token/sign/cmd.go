// Package token_sign implements 'kuma token sign'.
package token_sign

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/yurinandayona-com/kuma/cmd/token/config"
	"github.com/yurinandayona-com/kuma/userdb"
)

const (
	defaultExpirationDays = 100
)

// Cmd is 'kuma token sign` command.
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
			jm, err := token_config.LoadJWTManager(cmd)
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			var user *userdb.User
			for _, u := range jm.UserDB.GetUsers() {
				if u.Name == name {
					user = u
					break
				}
			}
			if user == nil {
				log.Fatalf("alert: user not found: %s", name)
			}

			log.Printf("debug: user.name = %#v", user.Name)
			log.Printf("debug: user.id = %#v", user.ID)
			log.Printf("debug: expiration_days = %d days", expirationDays)

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

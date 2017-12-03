package generate_tokens

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yurinandayona-com/kuma/cmd/serve"
	"github.com/yurinandayona-com/kuma/user_db"
	"log"
)

var Cmd *cobra.Command

func init() {
	var config string

	Cmd = &cobra.Command{
		Use:   "generate-tokens",
		Short: "Generate user tokens from configuration",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("info: load config: %s", config)
			cfg, err := serve.LoadConfig(config)
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			log.Printf("info: load user DB: %s", cfg.UserDB)
			userDB, err := user_db.LoadUserDB(cfg.UserDB)
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			jm := &user_db.JWTManager{
				UserDB:  userDB,
				HMACKey: []byte(cfg.HMACKey),
			}

			log.Print("info: generate tokens")
			users := userDB.GetUsers()
			tokens := make(map[string]string, len(users))
			for _, u := range userDB.GetUsers() {
				log.Printf("debug: generate token for %s", u.Name)

				tokens[u.Name], err = jm.Sign(u)
				if err != nil {
					log.Fatalf("alert: %s", err)
				}
			}

			j, err := json.MarshalIndent(tokens, "", "  ")
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			log.Print("info: show generated tokens")
			fmt.Println(string(j))
		},
	}

	Cmd.Flags().StringVarP(&config, "config", "C", ".kuma/serve.toml", "specify configuration file")
}

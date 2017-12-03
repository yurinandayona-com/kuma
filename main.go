package main

import (
	"github.com/comail/colog"
	"os"
)

func init() {
	colog.Register()
}

func main() {
	if err := Cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

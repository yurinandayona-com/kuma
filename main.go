package main

import (
	"os"
)

func main() {
	if err := Cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

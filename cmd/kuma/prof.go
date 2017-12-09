// +build prof

package main

import (
	"fmt"
	"github.com/phayes/freeport"
	"net/http"
	"log"
	_ "net/http/pprof"
)

func init() {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Printf("info: failed to get free port for pprof server: %s", err)
	}

	log.Printf("info: start pprof server: http://localhost:%d", port)
	go func() {
		log.Printf("info: stop pprof server: %s", http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}()
}

//go:build debug

package main

import (
	"log"
	"net/http"
	"os"
)
import _ "net/http/pprof"

func init() {
	go func() {
		path := os.Getenv("DEBUG_BIND")
		if path != "" {
			log.Printf("DEBUG IS ENABLED, THIS SHOULD NOT BE ON")
			log.Printf("UPDATE TO A RELEASE BUILD AND REMOVE THE DEBUG_BIND ENV")
			http.ListenAndServe(path, nil)
		}
	}()
}

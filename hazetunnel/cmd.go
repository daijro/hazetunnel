package main

/*
#include <stdlib.h>
*/

import (
	"C"
	"flag"
	"log"

	json "github.com/goccy/go-json"
)
import "context"

//export StartServer
func StartServer(data string) {
	// Launch server from cffi
	err := json.Unmarshal([]byte(data), &Flags)
	if err != nil {
		log.Fatal(err)
		return
	}
	go Launch()
}

//export ShutdownServer
func ShutdownServer() {
	// Kill server from cffi
	serverMux.Lock()
	defer serverMux.Unlock()

	if server != nil {
		log.Println("Shutting down the server...")
		cancel() // Cancel the context, which should trigger graceful shutdown

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown the server gracefully: %v", err)
		}
		server = nil // Set server to nil after shutdown
	}
}

func main() {
	// Launch from cli
	flag.StringVar(&Flags.Addr, "addr", "", "Proxy listen address")
	flag.StringVar(&Flags.Port, "port", "8080", "Proxy listen port")
	flag.StringVar(&Flags.Cert, "cert", "cert.pem", "TLS CA certificate (generated automatically if not present)")
	flag.StringVar(&Flags.Key, "key", "key.pem", "TLS CA key (generated automatically if not present)")
	flag.BoolVar(&Flags.Verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()
	Launch()
}

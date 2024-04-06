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

//export LaunchServer
func LaunchServer(data string) {
	// Launch from cffi
	err := json.Unmarshal([]byte(data), &Flags)
	if err != nil {
		log.Fatal(err)
		return
	}
	go launch()
}

func main() {
	// Launch from cli
	flag.StringVar(&Flags.Addr, "addr", "", "Proxy listen address")
	flag.StringVar(&Flags.Port, "port", "8080", "Proxy listen port")
	flag.StringVar(&Flags.Cert, "cert", "cert.pem", "TLS CA certificate (generated automatically if not present)")
	flag.StringVar(&Flags.Key, "key", "key.pem", "TLS CA key (generated automatically if not present)")
	flag.BoolVar(&Flags.Verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()
	launch()
}

package main

import (
	"flag"

	"github.com/daijro/hazetunnel/hazetunnel/api"
)

/*
Launch from CLI
*/

func main() {
	// Parse flags
	var Flags api.ProxySetup
	flag.StringVar(&Flags.Addr, "addr", "", "Proxy listen address")
	flag.StringVar(&Flags.Port, "port", "8080", "Proxy listen port")
	flag.StringVar(&Flags.UserAgent, "user_agent", "", "Override the User-Agent header for incoming requests. Optional.")
	flag.StringVar(&api.Config.Cert, "cert", "cert.pem", "TLS CA certificate (generated automatically if not present)")
	flag.StringVar(&api.Config.Key, "key", "key.pem", "TLS CA key (generated automatically if not present)")
	flag.BoolVar(&api.Config.Verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()
	// Set ID
	Flags.Id = "cli"
	// Set verbose level
	api.UpdateVerbosity()
	// Launch proxy server
	api.Launch(&Flags)
}

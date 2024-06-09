package main

import (
	"flag"
)

/*
Launch from CLI
*/

func main() {
	// Parse flags
	var Flags ProxySetup
	flag.StringVar(&Flags.Addr, "addr", "", "Proxy listen address")
	flag.StringVar(&Flags.Port, "port", "8080", "Proxy listen port")
	flag.StringVar(&Flags.UserAgent, "user_agent", "", "Override the User-Agent header for incoming requests. Optional.")
	flag.StringVar(&Config.Cert, "cert", "cert.pem", "TLS CA certificate (generated automatically if not present)")
	flag.StringVar(&Config.Key, "key", "key.pem", "TLS CA key (generated automatically if not present)")
	flag.BoolVar(&Config.Verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()
	// Set ID
	Flags.Id = "cli"
	// Set verbose level
	updateVerbosity()
	// Launch proxy server
	Launch(&Flags)
}

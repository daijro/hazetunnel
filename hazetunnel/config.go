package main

import (
	"github.com/elazarl/goproxy"
	utls "github.com/refraction-networking/utls"
)

type CLIFlags struct {
	Addr    string `json:"addr,omitempty"`
	Port    string `json:"port,omitempty"`
	Cert    string `json:"cert,omitempty"`
	Key     string `json:"key,omitempty"`
	Verbose bool   `json:"verbose,omitempty"`
}

var (
	Flags CLIFlags
)

func getClientHelloID(uagent string, ctx *goproxy.ProxyCtx) (utls.ClientHelloID, error) {
	browser, version, err := uagentToUtls(uagent)
	ctx.Logf("Client: %s, UTLS Version: %s", browser, version)

	if err != nil {
		return utls.ClientHelloID{}, err
	}

	return utls.ClientHelloID{
		Client:  browser,
		Version: version,
		Seed:    nil,
		Weights: nil,
	}, nil
}

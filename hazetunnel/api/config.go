package api

import (
	cflog "github.com/cloudflare/cfssl/log"
	"github.com/elazarl/goproxy"
	utls "github.com/refraction-networking/utls"
)

type ConfigFlags struct {
	// Global flags
	Cert    string `json:"cert,omitempty"`
	Key     string `json:"key,omitempty"`
	Verbose bool   `json:"verbose,omitempty"`
}

type ProxySetup struct {
	// Per proxy instance
	Addr          string `json:"addr,omitempty"`
	Port          string `json:"port"`
	UserAgent     string `json:"user_agent,omitempty"`
	Payload       string `json:"payload,omitempty"`
	UpstreamProxy string `json:"upstreamproxy,omitempty"`
	Id            string `json:"id"`
}

var (
	Config ConfigFlags
)

type VerbositySetting struct {
	Verbose bool `json:"verbose"`
}

type KeyPairSetting struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

func UpdateVerbosity() {
	// Update the verbose level
	if Config.Verbose {
		cflog.Level = cflog.LevelInfo
	} else {
		cflog.Level = cflog.LevelError
	}
}

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

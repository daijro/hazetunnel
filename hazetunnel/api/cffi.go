package api

/*
#include <stdlib.h>
*/

import (
	"C"
	"log"

	json "github.com/goccy/go-json"
)
import (
	"context"
)

/*
CFFI exposed methods
*/

//export StartServer
func StartServer(data string) {
	// Launch server from cffi
	var Flags ProxySetup
	err := json.Unmarshal([]byte(data), &Flags)
	if err != nil {
		log.Fatal(err)
		return
	}
	UpdateVerbosity()
	go Launch(&Flags)
}

//export SetVerbose
func SetVerbose(data string) {
	// Set the verbose option from cffi
	var verbose VerbositySetting
	err := json.Unmarshal([]byte(data), &verbose)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Update verbose level
	Config.Verbose = verbose.Verbose
	UpdateVerbosity() // Change immediately
}

//export SetKeyPair
func SetKeyPair(data string) {
	// Set the x509 key pair from cffi
	var keypair KeyPairSetting
	err := json.Unmarshal([]byte(data), &keypair)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Update the x509 key pair paths
	Config.Cert = keypair.Cert
	Config.Key = keypair.Key
	// Flag as unloaded
	caLoaded = false
	loadCA()
}

//export ShutdownServer
func ShutdownServer(id string) {
	// Kill server from cffi
	serverMux.Lock()
	defer serverMux.Unlock()

	// Check if id is in proxyInstanceMap
	if _, ok := proxyInstanceMap[id]; !ok {
		// say id wasnt found
		log.Printf("Error: %v is not a running instance", id)
		return
	}

	if proxyInstanceMap[id].Server == nil {
		log.Println("Error: Server not found")
		delete(proxyInstanceMap, id)
		return
	}

	// Announce server shutdown to verbose logs
	if Config.Verbose {
		log.Println("Shutting down the server...")
	}
	proxyInstanceMap[id].Cancel() // Cancel the context, which should trigger graceful shutdown

	if err := proxyInstanceMap[id].Server.Shutdown(context.Background()); err != nil {
		log.Printf("Failed to shutdown the server gracefully: %v", err)
	}
	delete(proxyInstanceMap, id)
}

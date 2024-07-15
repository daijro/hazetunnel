package api

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"sync"

	"github.com/elazarl/goproxy"

	cfsr "github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
)

var (
	caLoaded  = false
	caLoadMux sync.Mutex
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func setGoproxyCA(tlsCert tls.Certificate) {
	var err error
	if tlsCert.Leaf, err = x509.ParseCertificate(tlsCert.Certificate[0]); err != nil {
		log.Fatal("Unable to parse ca", err)
	}

	goproxy.GoproxyCa = tlsCert
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}
	caLoaded = true
}

func loadCA() {
	if caLoaded {
		return // Skip if cert already loaded
	}
	caLoadMux.Lock()
	defer caLoadMux.Unlock()

	// Set default values of Config.Cert and Config.Key
	if Config.Cert == "" {
		Config.Cert = "cert.pem"
	}
	if Config.Key == "" {
		Config.Key = "key.pem"
	}

	// If both Config.Cert and Config.Key exist, load them
	if fileExists(Config.Cert) && fileExists(Config.Key) {
		tlsCert, err := tls.LoadX509KeyPair(Config.Cert, Config.Key)
		if err != nil {
			log.Fatal("Unable to load CA certificate and key", err)
		}
		setGoproxyCA(tlsCert)
		return
	}

	// If only only file exists, warn the user
	if fileExists(Config.Cert) {
		log.Fatalf("CA certificate exists, but found no corresponding key at %s", Config.Key)
	} else if fileExists(Config.Key) {
		log.Fatalf("CA key exists, but found no corresponding certificate at %s", Config.Cert)
	}

	// Generate new CA files
	log.Println("No CA found, generating certificate and key")
	tlsCert, err := generateCA()
	if err != nil {
		log.Fatal("Unable to generate CA certificate and key", err)
	}
	setGoproxyCA(tlsCert)
}

func generateCA() (tls.Certificate, error) {
	csr := cfsr.CertificateRequest{
		CN:         "tlsproxy CA",
		KeyRequest: cfsr.NewKeyRequest(),
	}

	certPEM, _, keyPEM, err := initca.New(&csr)
	if err != nil {
		return tls.Certificate{}, err
	}

	caOut, err := os.Create(Config.Cert)
	if err != nil {
		return tls.Certificate{}, err
	}
	defer caOut.Close()
	_, err = caOut.Write(certPEM)
	if err != nil {
		return tls.Certificate{}, err
	}

	keyOut, err := os.OpenFile(Config.Key, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return tls.Certificate{}, err
	}
	defer keyOut.Close()

	_, err = keyOut.Write(keyPEM)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair(certPEM, keyPEM)
}

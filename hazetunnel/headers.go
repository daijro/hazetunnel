package main

import "net/http"

const (
	UpStreamProxyHeader = "x-mitm-upstream"
	PayloadHeader       = "x-mitm-payload"
)

var CustomHeaders = []string{UpStreamProxyHeader, PayloadHeader}

type ProxyConfig struct {
	upstreamProxy string
	payload       string
}

func parseCustomHeaders(headers *http.Header) ProxyConfig {
	return ProxyConfig{
		upstreamProxy: headers.Get(UpStreamProxyHeader),
		payload:       headers.Get(PayloadHeader),
	}
}

func removeCustomHeaders(headers *http.Header) {
	for _, header := range CustomHeaders {
		headers.Del(header)
	}
}

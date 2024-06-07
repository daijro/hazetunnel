package main

import "net/http"

const (
	UpStreamProxyHeader = "x-mitm-upstream"
	PayloadHeader       = "x-mitm-payload"
	IsBase64            = "x-mitm-isbase64"
)

var CustomHeaders = []string{UpStreamProxyHeader, PayloadHeader, IsBase64}

type ProxyConfig struct {
	upstreamProxy string
	payload       string
	isBase64      string
}

func parseCustomHeaders(headers *http.Header) ProxyConfig {
	return ProxyConfig{
		upstreamProxy: headers.Get(UpStreamProxyHeader),
		payload:       headers.Get(PayloadHeader),
		isBase64:      headers.Get(IsBase64),
	}
}

func removeCustomHeaders(headers *http.Header) {
	for _, header := range CustomHeaders {
		headers.Del(header)
	}
}

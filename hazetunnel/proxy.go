package main

import (
	"context"
	"log"
	"net/http"
	"net/url"

	cflog "github.com/cloudflare/cfssl/log"
	"github.com/elazarl/goproxy"
	utls "github.com/refraction-networking/utls"
	sf "gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/snowflake/v2/common/utls"
)

type contextKey string

const payloadKey contextKey = "payload"

func launch() {
	if !Flags.Verbose {
		cflog.Level = cflog.LevelError
	}
	loadCA()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = Flags.Verbose

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			proxyConfig := parseCustomHeaders(&req.Header)
			removeCustomHeaders(&req.Header)

			var upstreamProxy *url.URL

			// // Confirm that the payload header was passed
			// if proxyConfig.payload == "" {
			//     return req, missingParameterResponse(req, ctx, PayloadHeader)
			// }
			// Extract ClientHello from the User-Agent header
			if len(req.Header["User-Agent"]) == 0 {
				return req, missingParameterResponse(req, ctx, "User-Agent")
			}
			ua := req.Header["User-Agent"][0]

			clientHelloId, err := getClientHelloID(ua, ctx)
			if err != nil {
				ctx.Logf("Error parsing UserAgent: %s", err)
				return req, invalidClientResponse(req, ctx, ua)
			}

			// Store the payload code in the request's context
			ctx.Req = req.WithContext(context.WithValue(ctx.Req.Context(), payloadKey, proxyConfig.payload))

			// If a proxy header was passed, set it to upstreamProxy
			if len(proxyConfig.upstreamProxy) > 0 {
				proxyUrl, err := url.Parse(proxyConfig.upstreamProxy)
				if err != nil {
					return req, invalidUpstreamProxyResponse(req, ctx, proxyConfig.upstreamProxy)
				}
				upstreamProxy = proxyUrl
			}

			// Skip TLS handshake if scheme is HTTP
			ctx.Logf("Scheme: %s", req.URL.Scheme)
			if req.URL.Scheme == "http" {
				ctx.Logf("Skipping TLS for HTTP request")
				return req, nil
			}

			// Build round tripper
			roundTripper := sf.NewUTLSHTTPRoundTripperWithProxy(clientHelloId, &utls.Config{
				InsecureSkipVerify: true,
				OmitEmptyPsk:       true,
			}, http.DefaultTransport, false, upstreamProxy)

			ctx.RoundTripper = goproxy.RoundTripperFunc(
				func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Response, error) {
					return roundTripper.RoundTrip(req)
				})

			return req, nil
		},
	)

	// Inject payload code into responses
	proxy.OnResponse().DoFunc(PayloadInjector)

	listenAddr := Flags.Addr + ":" + Flags.Port
	log.Println("hazetunnel listening at " + listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, proxy))
}

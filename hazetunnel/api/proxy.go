package api

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/elazarl/goproxy"
	utls "github.com/refraction-networking/utls"
	sf "gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/snowflake/v2/common/utls"
)

type contextKey string

const payloadKey contextKey = "payload"

type ProxyInstance struct {
	Server *http.Server
	Cancel context.CancelFunc
}

// Globals
var (
	serverMux        sync.Mutex
	proxyInstanceMap = make(map[string]*ProxyInstance)
)

func initServer(Flags *ProxySetup) *http.Server {
	serverMux.Lock()
	defer serverMux.Unlock()

	// Load CA if not already loaded
	loadCA()

	// Setup the proxy instance
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = Config.Verbose
	setupProxy(proxy, Flags)

	// Create the server
	server := &http.Server{
		Addr:    Flags.Addr + ":" + Flags.Port,
		Handler: proxy,
	}
	_, cancel := context.WithCancel(context.Background())

	// Add proxy instance to the map
	proxyInstanceMap[Flags.Id] = &ProxyInstance{
		Server: server,
		Cancel: cancel,
	}
	return server
}

func setupProxy(proxy *goproxy.ProxyHttpServer, Flags *ProxySetup) {
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			var upstreamProxy *url.URL

			// Override the User-Agent header if specified
			// If one wasn't specified, verify a User-Agent is in the request
			if len(Flags.UserAgent) != 0 {
				req.Header["User-Agent"] = []string{Flags.UserAgent}
			} else if len(req.Header["User-Agent"]) == 0 {
				return req, missingParameterResponse(req, ctx, "User-Agent")
			}

			// Set the ClientHello from the User-Agent header
			ua := req.Header["User-Agent"][0]
			clientHelloId, err := getClientHelloID(ua, ctx)
			if err != nil {
				// Use the latest Chrome when the User-Agent header cannot be recognized
				ctx.Logf("Error parsing User-Agent: %s", err)
				clientHelloId = utls.HelloChrome_Auto
				ctx.Logf("Continuing with Chrome %v ClientHello", clientHelloId.Version)
			}

			// Store the payload code in the request's context
			ctx.Req = req.WithContext(
				context.WithValue(
					ctx.Req.Context(),
					payloadKey,
					Flags.Payload,
				),
			)

			// If a proxy header was passed, set it to upstreamProxy
			if len(Flags.UpstreamProxy) != 0 {
				proxyUrl, err := url.Parse(Flags.UpstreamProxy)
				if err != nil {
					return req, invalidUpstreamProxyResponse(req, ctx, Flags.UpstreamProxy)
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
}

// Launches the server
func Launch(Flags *ProxySetup) {
	server := initServer(Flags)

	// Print server startup message if from CLI or verbose CFFI
	if Flags.Id == "cli" || Config.Verbose {
		log.Println("Hazetunnel listening at", server.Addr)
	}
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

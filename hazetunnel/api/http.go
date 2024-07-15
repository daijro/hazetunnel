package api

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

func invalidUpstreamProxyResponse(
	req *http.Request,
	ctx *goproxy.ProxyCtx,
	upstreamProxy string,
) *http.Response {
	ctx.Warnf("CRITICAL: Client specified invalid upstream proxy: %s", upstreamProxy)
	return goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusBadRequest, "HAZETUNNEL ERROR: Invalid upstream proxy: "+upstreamProxy)
}

func missingParameterResponse(
	req *http.Request,
	ctx *goproxy.ProxyCtx,
	header string,
) *http.Response {
	ctx.Warnf("CRITICAL: Missing header: %s", header)
	return goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusBadRequest, "HAZETUNNEL ERROR: Missing header: "+header)
}

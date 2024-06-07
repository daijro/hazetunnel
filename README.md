<h1 align="center">
    Hazetunnel
</h1>

<h3 align="center">
    🔮 Vindicate non-organic web traffic
</h3>

---

Hazetunnel is an MITM proxy that attempts to legitimize [BrowserForge](https://github.com/daijro/browserforge/)'s injected-browser web traffic by hijacking the TLS fingerprint to mirror the passed User-Agent.

Additionally, it can inject a Javascript payload into the web page to defend against [worker fingerprinting](https://github.com/apify/fingerprint-suite/issues/64).

<hr width=50>

### Features ✨

- Anti TLS fingerprinting 🪪

  - Emulate the ClientHello of browsers based on the passed User-Agent (e.g. Chrome/120)
  - Bypasses TLS fingerprinting checks

- Javascript payload injection 💉

  - Prepends payload to all Javascript responses, including the web page Service/Shared worker scope.
  - Injects payload into embedded base64 encoded JavaScript within HTML responses ([see here](https://github.com/apify/fingerprint-suite/issues/64#issuecomment-1282877696))

This project was built on [tlsproxy](https://github.com/rosahaj/tlsproxy), please leave them a star!

---

## Integration

### Header table

Add the following headers to each request to the proxy:

| Header            | Description                                                                      | Example                                       |
| ----------------- | -------------------------------------------------------------------------------- | --------------------------------------------- |
| `x-mitm-payload`  | Inject a JavaScript payload into the response.                                   | `alert('Hello world');`                       |
| `x-mitm-isbase64` | Set to `1` to pass the payload as a Base64 encoded string.                       | `1`                                           |
| `x-mitm-upstream` | Optionally forward the request to the upstream proxy. Must be socks5 or socks5h. | `socks5://user:pass@pro.proxyvendor.com:7000` |

### Curl

Assuming Hazetunnel is running on `localhost:8080`:

```bash
curl \
--proxy http://localhost:8080 \
--cacert cert.pem \
"https://example.com" \
-H "x-mitm-payload: alert('Hello world');"
```

### Python Requests

```py
requests.get(
    'https://example.com',
    headers={
      'x-mitm-payload': 'alert("Hello world");'
    },
    proxies={
      'http': 'http://localhost:8080',
      'https': 'http://localhost:8080',
    },
    verify='cert.pem'
)
```

<hr width=50>

## Building

### CFFI

Pre-built C shared library binaries provided in [Releases](https://github.com/daijro/hazetunnel/releases).

Otherwise, you can build these yourself using the `build.bat` file provided.

### CLI

#### Building from source

```bash
git clone https://github.com/daijro/hazetunnel
cd hazetunnel
go build
```

#### Usage

```
Usage of ./hazetunnel:
  -addr string
        Proxy listen address
  -cert string
        TLS CA certificate (generated automatically if not present) (default "cert.pem")
  -key string
        TLS CA key (generated automatically if not present) (default "key.pem")
  -port string
        Proxy listen port (default "8080")
  -verbose
        Enable verbose logging
```

---

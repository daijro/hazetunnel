# Python Usage

## Usage

```py
from hazetunnel import HazeTunnel
from browserforge.headers import HeaderGenerator
...
# Initialize the proxy
proxy = HazeTunnel(port='8080', payload='alert("Hello World!");')
proxy.launch()
# Send the request
requests.get(
    url='https://example.com',
    headers=HeaderGenerator().generate(browser='chrome'),
    proxies={'https': proxy.url},
    verify=proxy.cert
).text
# Stop the proxy
proxy.stop()
```

<details>

  <summary>
    HazeTunnel parameters
  </summary>

```
Parameters:
    port (Optional[str]): Specify a port to listen on. Default is random.
    payload (Optional[str]): Payload to inject into responses
    user_agent (Optional[str]): Optionally override all User-Agent headers
    upstream_proxy (Optional[str]): Optionally forward requests to an upstream proxy
```

</details>

### Using a context manager

A context manager will automatically close the server when not needed anymore.

```py
with HazeTunnel(port='8080', payload='alert("Hello World!");') as proxy:
  # Send the request
  requests.get(
      url='https://example.com',
      headers=HeaderGenerator().generate(browser='chrome'),
      proxies={'https': proxy.url},
      verify=proxy.cert
  ).text
```

<hr width=70>

## CLI

Download the latest version of the API:

```sh
python -m hazetunnel fetch
```

Remove all files before uninstalling

```sh
python -m hazetunnel remove
```

Run the MITM proxy:

```sh
python -m hazetunnel run -p 8080 --verbose
```

### All commands

```sh
Usage: python -m hazetunnel [OPTIONS] COMMAND [ARGS]...

Options:
  --help  Show this message and exit.

Commands:
  fetch    Fetch the latest version of hazetunnel-api
  remove   Remove all library files
  run      Run the MITM proxy
  version  Display the current version
```

### See [here](https://github.com/daijro/hazetunnel) for more information and examples.

---

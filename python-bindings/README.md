# Python Usage

## Usage

Hazetunnel is designed to run globally:

```py
import hazetunnel
hazetunnel.launch()
...
requests.get(
    url='https://tls.peet.ws/api/clean',
    headers={
        **HeaderGenerator().generate(browser='chrome'),
        'x-mitm-payload': 'alert("hi");'
    },
    proxies={'https': hazetunnel.url()},
    verify=hazetunnel.cert()
)
...
hazetunnel.kill()
```

Although, Hazetunnel may also run in a context manager:

```py
from hazetunnel import HazeTunnel
...
with HazeTunnel as proxy:
    requests.get(
        url='https://tls.peet.ws/api/clean',
        headers={
            **HeaderGenerator().generate(browser='chrome'),
            'x-mitm-payload': 'alert("hi");'
        },
        proxies={'https': proxy.url},
        verify=proxy.cert
    )
...
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

---

# Python Usage

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
hazetunnel.stop()
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

---


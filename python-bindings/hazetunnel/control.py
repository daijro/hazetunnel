from pathlib import Path
from typing import Optional, Union

from .cffi import get_library, root_dir

"""
Hazetunnel is designed to run globally.

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
).text
...
hazetunnel.stop()
"""


def launch(
    port: Optional[int] = None,
    cert: Union[str, Path] = root_dir / "bin" / "cert.pem",
    key: Union[str, Path] = root_dir / "bin" / "key.pem",
    verbose: bool = False,
) -> None:
    lib = get_library()
    if lib._started:
        raise RuntimeError("Server is already running.")
    get_library().launch(port=port, cert=str(cert), key=str(key), verbose=verbose)


def port() -> int:
    return get_library().port


def url() -> str:
    return f"http://localhost:{port()}"


def cert() -> Optional[str]:
    if pair := get_library()._cert_key_pair:
        return pair[0]


def key() -> Optional[str]:
    if pair := get_library()._cert_key_pair:
        return pair[1]


def stop() -> None:
    lib = get_library()
    if not lib._started:
        raise RuntimeError("Server is not running.")
    get_library().stop_server()


"""
Hazetunnel may also run in a context manager.

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
    ).text
...
"""


class Context:
    def __init__(self, **kwargs) -> None:
        self.kwargs = kwargs

    def __enter__(self) -> "Context":
        launch(**self.kwargs)
        return self

    def __exit__(self, *_) -> None:
        stop()

    @property
    def url(self) -> str:
        return url()

    @property
    def port(self) -> int:
        return port()

    @property
    def cert(self) -> Optional[str]:
        return cert()

    @property
    def key(self) -> Optional[str]:
        return key()


HazeTunnel: Context = Context()

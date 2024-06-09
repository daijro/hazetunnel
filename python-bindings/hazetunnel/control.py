from pathlib import Path
from typing import Optional, Union
from uuid import uuid4

from .cffi import get_library

"""
Hazetunnel may also run in a context manager.

from hazetunnel import HazeTunnel
...
with HazeTunnel(port='8080', payload='alert("Hello World!");') as proxy:
    requests.get(
        url='https://tls.peet.ws/api/clean',
        headers=HeaderGenerator().generate(browser='chrome'),
        proxies={'https': proxy.url},
        verify=proxy.cert
    ).text
...
"""


class HazeTunnel:
    def __init__(
        self,
        port: Optional[str] = None,
        payload: Optional[str] = None,
        user_agent: Optional[str] = None,
        upstream_proxy: Optional[str] = None,
    ) -> None:
        """
        HazeTunnel constructor

        Parameters:
            port (Optional[str]): Specify a port to listen on. Default is random.
            payload (Optional[str]): Payload to inject into responses
            user_agent (Optional[str]): Override user agent
            upstream_proxy (Optional[str]): Optionally forward requests to an upstream proxy
        """
        # Generate a ID
        self.id = str(uuid4())

        # Set options
        self.options = {
            "port": port or '',
            "payload": payload or '',
            "user_agent": user_agent or '',
            "upstream_proxy": upstream_proxy or '',
            "id": self.id,
        }

        self.lib = get_library()
        self.is_running = False

    """
    Start/stopping the server
    """

    def launch(self) -> None:
        """
        Launch the server
        """
        # Kill if running already
        if self.is_running:
            raise RuntimeError("Server is already running.")
        # Generate a port if one wasn't passed
        if not self.options['port']:
            self.options['port'] = str(self.lib.get_open_port())

        self.lib.start_server(self.options)
        self.is_running = True

    def stop(self) -> None:
        """
        Stop the server
        """
        if not self.is_running:
            raise RuntimeError("Server is not running.")
        self.lib.stop_server(self.id)
        self.is_running = False

    """
    Configuration
    """

    @property
    def url(self) -> str:
        """
        Returns the URL of the server
        """
        return f"http://127.0.0.1:{self.options['port']}"

    @property
    def cert(self) -> str:
        """
        Returns the path to the server's certificate
        """
        return self.lib.key_pair[0]

    @cert.setter
    def cert(self, _) -> None:
        raise NotImplementedError(
            "Setting the certificate path is not supported. "
            "This must be done with the hazetunnel.set_curt(path) method"
        )

    @property
    def key(self) -> str:
        """
        Returns the path to the server's key
        """
        return self.lib.key_pair[1]

    @key.setter
    def key(self, _) -> None:
        raise NotImplementedError(
            "Setting the certificate path is not supported. "
            "This must be done with the hazetunnel.set_key(path) method"
        )

    @property
    def verbose(self) -> bool:
        """
        Returns the verbosity of the server logs
        """
        return self.lib.verbose

    @verbose.setter
    def verbose(self, option: bool) -> None:
        """
        Set the verbosity of the logs
        """
        self.lib.verbose = option

    """
    Context manager methods
    """

    def __enter__(self) -> "HazeTunnel":
        self.launch()
        return self

    def __exit__(self, *_) -> None:
        self.stop()


"""
Global config setters
Note: These MUST be done before starting the HazeTunnel instance.
"""


def set_key_pair(
    cert_path: Optional[Union[str, Path]], key_path: Optional[Union[str, Path]] = ''
) -> None:
    """
    Set the certificate path
    """
    if not (cert_path or key_path):
        raise ValueError("Either cert and key must be set")
    lib = get_library()
    lib.key_pair = (str(cert_path), str(key_path))


def set_verbose(option: bool) -> None:
    """
    Set the logging level to verbose
    """
    lib = get_library()
    lib.verbose = option


"""
Global config getters
"""

cert = lambda: get_library().key_pair[0]
key = lambda: get_library().key_pair[1]
verbose = lambda: get_library().verbose

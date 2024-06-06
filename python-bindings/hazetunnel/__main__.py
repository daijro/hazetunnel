"""
Binary CLI manager for hazetunnel-api.

Adapted from https://github.com/daijro/hrequests/blob/main/hrequests/__main__.py
"""

import os
import re
import time
from dataclasses import dataclass
from functools import total_ordering
from importlib.metadata import version as pkg_version
from pathlib import Path
from typing import Optional

import click

from hazetunnel import launch, stop
from hazetunnel.__version__ import BRIDGE_VERSION
from hazetunnel.cffi import LibraryManager, root_dir


def rprint(*a, **k):
    click.secho(*a, **k, bold=True)


@total_ordering
@dataclass
class Version:
    version: str

    def __post_init__(self) -> None:
        self.sort_version = tuple(int(x) for x in self.version.split('.'))

    def __eq__(self, other) -> bool:
        return self.sort_version == other.sort_version

    def __lt__(self, other) -> bool:
        return self.sort_version < other.sort_version

    def __str__(self) -> str:
        return self.version

    @staticmethod
    def get_version(name) -> 'Version':
        ver: Optional[re.Match] = LibraryUpdate.FILE_NAME.search(name)
        if not ver:
            raise ValueError(f'Could not find version in {name}')
        return Version(ver[1])


@dataclass
class Asset:
    url: str
    name: str

    def __post_init__(self) -> None:
        self.version: Version = Version.get_version(self.name)


class LibraryUpdate(LibraryManager):
    """
    Checks if an update is available for hazetunnel-api library
    """

    FILE_NAME: re.Pattern = re.compile(r'^hazetunnel-api-v([\d\.]+)')

    def __init__(self) -> None:
        self.parent_path: Path = root_dir / 'bin'
        self.file_cont, self.file_ext = self.get_name()
        self.file_pref = f'hazetunnel-api-v{BRIDGE_VERSION}'

    @property
    def path(self) -> Optional[str]:
        if paths := self.get_files():
            return paths[0]

    @property
    def full_path(self) -> Optional[str]:
        if path := self.path:
            return os.path.join(self.parent_path, path)

    def latest_asset(self) -> Asset:
        """
        Find the latest Asset for the hazetunnel-api library
        """
        releases = self.get_releases()
        for release in releases:
            if asset := self.check_assets(release['assets']):
                url, name = asset
                return Asset(url, name)
        raise ValueError('No assets found for hazetunnel-api')

    def install(self) -> None:
        filename = super().check_library()
        ver: Version = Version.get_version(filename)

        rprint(f"Successfully downloaded hazetunnel-api v{ver}!", fg="green")

    def update(self) -> None:
        """
        Updates the library if needed
        """
        path = self.path
        if not path:
            # install the library if it doesn't exist
            return self.install()

        # get the version
        current_ver: Version = Version.get_version(path)

        # check if the version is the same as the latest available version
        asset: Asset = self.latest_asset()
        if current_ver >= asset.version:
            rprint("hazetunnel-api library up to date!", fg="green")
            rprint(f"Current version: v{current_ver}", fg="green")
            return

        # download updated file
        rprint(
            f"Updating hazetunnel-api library from v{current_ver} => v{asset.version}", fg="yellow"
        )
        # remove old, download new
        self.download_file(os.path.join(self.parent_path, asset.name), asset.url)
        try:
            os.remove(os.path.join(self.parent_path, path))
        except OSError:
            rprint("WARNING: Could not remove outdated library files.", fg="yellow")


@click.group()
def cli() -> None:
    pass


@cli.command(name='fetch')
def fetch():
    """
    Fetch the latest version of hazetunnel-api
    """
    LibraryUpdate().update()


@cli.command(name='remove')
def remove() -> None:
    """
    Remove all library files
    """
    path = str(LibraryUpdate().full_path)
    # remove all .pem files
    for file in (root_dir / 'bin').glob('*.pem'):
        rprint(f"Removed {file}", fg="green")
        file.unlink()
    # remove library
    if not os.path.exists(path):
        rprint("Library is not downloaded.", fg="yellow")
        return
    try:
        os.remove(path)
    except OSError as e:
        rprint(f"WARNING: Could not remove {path}: {e}", fg="red")
    else:
        rprint(f"Removed {path}", fg="green")
    rprint("Library files have been removed.", fg="yellow")


@cli.command(name='version')
def version() -> None:
    """
    Display the current version
    """
    # python package version
    rprint(f"Pip package:\tv{pkg_version('hazetunnel')}", fg="green")

    # library path
    libup = LibraryUpdate()
    path = libup.path
    # if the library is not installed
    if not path:
        rprint("hazetunnel-api:\tNot downloaded!", fg="red")
        return
    # library version
    lib_ver = Version.get_version(path)
    rprint(f"hazetunnel-api:\tv{lib_ver} ", fg="green", nl=False)

    # check for library updates
    latest_ver = libup.latest_asset().version
    if latest_ver == lib_ver:
        rprint("(Up to date!)", fg="yellow")
    else:
        rprint(f"(Latest: v{latest_ver})", fg="red")


@cli.command(name='run')
@click.option('-v', '--verbose', is_flag=True, help="Enable verbose output")
@click.option(
    '-p', '--port', type=int, default=8080, help="Port to run the proxy on. Default: 8080"
)
@click.option('--cert', type=str, default=None, help="Path to the certificate file. Optional.")
@click.option('--key', type=str, default=None, help="Path to the key file. Optional.")
def run(port: int, cert: Optional[str], key: Optional[str], verbose: bool) -> None:
    """
    Run the MITM proxy
    """
    launch(port=port, cert=cert, key=key, verbose=verbose)
    # wait forever until keyboard interrupt
    try:
        time.sleep(1e6)
    except KeyboardInterrupt:
        pass
    finally:
        stop()


if __name__ == '__main__':
    cli()

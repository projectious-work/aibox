"""The vendored ruamel.yaml wheel: SHA-256 check and import as a YAML 1.2 parser."""

from __future__ import annotations

import hashlib
import json
import sys
from pathlib import Path

VENDOR = Path(__file__).resolve().parent.parent / "vendor"


class WheelError(Exception):
    """A missing or mismatching wheel: an infrastructure error, never a skip."""


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with open(path, "rb") as fh:
        for chunk in iter(lambda: fh.read(1 << 16), b""):
            digest.update(chunk)
    return digest.hexdigest()


def manifest(vendor: Path = VENDOR) -> dict:
    return json.loads((vendor / "wheels.json").read_text())


def verify(name: str = "ruamel.yaml", vendor: Path = VENDOR) -> Path:
    """Return the wheel path after checking it against wheels.json."""
    try:
        entries = [w for w in manifest(vendor)["wheels"] if w["name"] == name]
    except (OSError, ValueError, KeyError) as exc:
        raise WheelError(f"cannot read {vendor / 'wheels.json'}: {exc}") from exc
    if len(entries) != 1:
        raise WheelError(f"wheels.json lists {name} {len(entries)} times")
    entry = entries[0]
    path = vendor / entry["file"]
    if not path.is_file():
        raise WheelError(f"vendored wheel missing: {path} (source {entry.get('url')})")
    actual = sha256_file(path)
    if actual != entry["sha256"]:
        raise WheelError(f"SHA-256 mismatch for {path}: {actual}, recorded {entry['sha256']}")
    return path


_YAML = None


def yaml_parser(vendor: Path = VENDOR):
    """A ruamel.yaml safe, pure-Python YAML 1.2 loader from the checked wheel."""
    global _YAML
    if _YAML is None:
        path = verify(vendor=vendor)
        if str(path) not in sys.path:
            sys.path.insert(0, str(path))
        from ruamel.yaml import YAML  # noqa: PLC0415

        _YAML = YAML(typ="safe", pure=True)
    return _YAML


def load_yaml_documents(text: str) -> list:
    return list(yaml_parser().load_all(text))

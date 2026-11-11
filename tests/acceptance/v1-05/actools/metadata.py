"""AC-TOOLS/side-effects: published Feature metadata and workspace state."""

from __future__ import annotations

import hashlib
import io
import json
import os
import stat
import tarfile
from pathlib import Path
from typing import Mapping

FORBIDDEN_KEYS = ("onCreateCommand", "updateContentCommand", "postCreateCommand",
                  "postStartCommand", "postAttachCommand", "mounts", "privileged", "capAdd",
                  "securityOpt", "init", "entrypoint")
CLI_LOCKFILE = ".devcontainer/devcontainer-lock.json"


def feature_json_from_layer(layer: bytes) -> dict:
    """Read devcontainer-feature.json from a published Feature layer."""
    with tarfile.open(fileobj=io.BytesIO(layer), mode="r:*") as archive:
        for member in archive.getmembers():
            if member.isfile() and member.name.lstrip("./") == "devcontainer-feature.json":
                stream = archive.extractfile(member)
                if stream is None:
                    break
                return json.loads(stream.read().decode("utf-8"))
    raise ValueError("published layer has no devcontainer-feature.json")


def side_effect_violations(metadata: Mapping, expected_env: Mapping[str, str]) -> list[str]:
    """Chapter 8 Side effects: none of the forbidden keys is present, and
    `containerEnv` equals the catalog's or is absent when the catalog has none."""
    violations = [f"declares {key}" for key in FORBIDDEN_KEYS if key in metadata]
    if not expected_env:
        if "containerEnv" in metadata:
            violations.append(f"declares containerEnv {json.dumps(metadata['containerEnv'])} "
                              "although the catalog has none")
    elif metadata.get("containerEnv") != dict(expected_env):
        violations.append(f"containerEnv {json.dumps(metadata.get('containerEnv'), sort_keys=True)} "
                          f"differs from catalog {json.dumps(dict(expected_env), sort_keys=True)}")
    return violations


def snapshot(root: Path) -> dict[str, str]:
    """Content identity of every entry under the workspace folder."""
    entries: dict[str, str] = {}
    for directory, dirnames, filenames in os.walk(root, followlinks=False):
        dirnames.sort()
        for name in sorted(dirnames + filenames):
            path = Path(directory) / name
            relative = path.relative_to(root).as_posix()
            info = path.lstat()
            if stat.S_ISLNK(info.st_mode):
                entries[relative] = "link:" + os.readlink(path)
            elif stat.S_ISDIR(info.st_mode):
                entries[relative] = f"dir:{stat.S_IMODE(info.st_mode):o}"
            elif stat.S_ISREG(info.st_mode):
                digest = hashlib.sha256(path.read_bytes()).hexdigest()
                entries[relative] = f"file:{stat.S_IMODE(info.st_mode):o}:{digest}"
            else:
                entries[relative] = "special"
    return entries


def snapshot_changes(before: Mapping[str, str], after: Mapping[str, str],
                     exempt: tuple[str, ...] = ()) -> list[str]:
    changes = []
    for path in sorted(set(before) | set(after)):
        if path in exempt:
            continue
        if path not in after:
            changes.append(f"deleted {path}")
        elif path not in before:
            changes.append(f"added {path}")
        elif before[path] != after[path]:
            changes.append(f"modified {path}")
    return changes

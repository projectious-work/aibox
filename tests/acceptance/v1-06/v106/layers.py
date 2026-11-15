"""Setting one key in one layer, for the precedence journeys J-CFG-02 and J-CFG-03."""

from __future__ import annotations

import json

from . import jsonc
from .contract import CONTRACT
from .sandbox import Sandbox

# Two different valid values per key: (lower, higher). Values are typed as in
# a settings file; env and flag forms are derived by text().
PAIR_VALUES = {
    "output.format": ("yaml", "json"),
    "output.color": ("always", "never"),
    "logging.level": ("debug", "warn"),
    "logging.format": ("plain", "jsonl"),
    "logging.file": ("a.log", "b.log"),
    "logging.required": (True, False),
    "logging.rotationMiB": (5, 6),
    "logging.retentionFiles": (3, 4),
    "logging.retentionDays": (10, 20),
    "execution.timeoutSeconds": (40, 50),
    "execution.lockWaitSeconds": (10, 20),
    "diagnostics.offline": (False, True),
    "workspace.theme": ("dracula", "nord"),
    "workspace.mode": ("light", "dark"),
    "workspace.layout": ("ai", "focus"),
}


def text(value) -> str:
    """The env-file, environment and flag form of a typed value (settings.conversion)."""
    if isinstance(value, bool):
        return "true" if value else "false"
    return str(value)


def nested(key: str, value) -> dict:
    section, name = key.split(".", 1)
    return {section: {name: value}}


def merge(into: dict, add: dict) -> dict:
    for k, v in add.items():
        if isinstance(v, dict) and isinstance(into.get(k), dict):
            merge(into[k], v)
        else:
            into[k] = v
    return into


class LayerSetup:
    """Collects layer values, then writes the files and builds the command."""

    def __init__(self, sb: Sandbox):
        self.sb = sb
        self.system: dict = {}
        self.user: dict = {}
        self.workspace: dict = {}
        self.env_file: list[str] = []
        self.environment: dict[str, str] = {}
        self.flags: list[str] = []
        self.env_file_path = f"{sb.cwd}/layers.env"

    def set(self, key: str, layer: str, value) -> None:
        spec = CONTRACT.key(key)
        if layer == "system":
            merge(self.system, nested(key, value))
        elif layer == "user":
            merge(self.user, nested(key, value))
        elif layer == "project":
            self.workspace[key.split(".", 1)[1]] = value
        elif layer == "env_file":
            self.env_file.append(f"{spec['env']}={text(value)}")
        elif layer == "environment":
            self.environment[spec["env"]] = text(value)
        elif layer == "invocation":
            self.flags += [spec["flag"], text(value)]
        else:
            raise ValueError(layer)

    def apply(self, kind: str = "P") -> tuple[str, list[str], dict[str, str]]:
        """Write everything; return (project path, extra args, environment)."""
        sb = self.sb
        if self.system:
            sb.system_settings({"schemaVersion": "1", **self.system})
        if self.user:
            sb.user_settings({"schemaVersion": "1", **self.user})
        edit = None
        if self.workspace:
            ws = dict(self.workspace)
            edit = lambda text_: jsonc.add_members(  # noqa: E731
                text_, {"customizations": {"aibox": {"workspace": ws}}})
        project = sb.project(kind, edit=edit)
        args = list(self.flags)
        if self.env_file:
            sb.env_file("layers.env", "\n".join(self.env_file) + "\n")
            args += ["--env-file", self.env_file_path]
        return project, args, dict(self.environment)

    def source(self, key: str, layer: str, project: str) -> str:
        """effectiveConfig.source for this setup, in the view's path form."""
        spec = CONTRACT.key(key)
        sb = self.sb
        return {
            "built_in": "built-in",
            "system": CONTRACT.system_settings_path,
            "user": sb.tilde(sb.user_settings_path()),
            "project": sb.tilde(f"{project}/.devcontainer/devcontainer.json"),
            "env_file": sb.tilde(self.env_file_path),
            "environment": spec["env"],
            "invocation": spec["flag"],
        }[layer]


def view_value(sb: Sandbox, key: str, value):
    """The value the view reports for a typed setting value."""
    if key == "logging.file":
        return sb.tilde(f"{sb.logs}/{value}")
    return value


def dumps(obj) -> str:
    return json.dumps(obj, sort_keys=True)

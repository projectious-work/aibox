"""Shared test data: the frozen catalog and a small synthetic catalog."""

from __future__ import annotations

import copy
import json
import sys
from pathlib import Path

SUITE = Path(__file__).resolve().parents[1]
if str(SUITE) not in sys.path:
    sys.path.insert(0, str(SUITE))

from actools.catalog import SPEC, Catalog  # noqa: E402

REAL = Catalog.load()
CENSUS = [f"{recipe['id']}/{tool['name']}"
          for recipe in json.loads((SPEC / "ledger" / "addons.json").read_text())
          for tool in recipe["tools"]]


def tool(name, *, commands=(), checks=None, version=None, also=None, method="producer-release",
         rows=()):
    entry = {"name": name, "rows": list(rows), "method": method, "docsUrl": "https://x",
             "source": "s", "verification": "v", "license": "MIT",
             "checks": checks if checks is not None else
             [{"kind": "version", "argv": [name, "--version"],
               "regex": "(^|[^0-9.]){version}([^0-9]|$)"}],
             "commands": list(commands)}
    if version:
        entry["version"] = version
    if also:
        entry["alsoProvides"] = also
    return entry


def feature(feature_id, group, *, role="tool", owner="aibox", versions=None, depends=(),
            requires=(), tools=(), env=None):
    entry = {"id": feature_id, "group": group, "role": role, "owner": owner,
             "targets": ["debian-13-amd64", "debian-13-arm64"], "options": {},
             "dependsOn": list(depends), "requires": list(requires), "tools": list(tools)}
    if owner == "aibox":
        entry["featureVersion"] = "2.3.4"
    else:
        entry["upstream"] = {"ref": f"ghcr.io/devcontainers/features/{group}@sha256:" + "a" * 64,
                             "publisher": "p", "source": "https://s", "featureVersion": "1.0.0",
                             "options": {"version": "9.9.9"}}
    if versions:
        entry["options"]["version"] = {"type": "string", "enum": list(versions),
                                       "default": versions[-1]}
    if env:
        entry["containerEnv"] = env
    return entry


def synthetic() -> Catalog:
    """runtime <- (lib requires runtime), kit/{a,b,c,all}, plus an upstream tool."""
    data = copy.deepcopy({k: v for k, v in REAL.data.items()
                          if k in ("schemaVersion", "targets", "errors", "unsupportedTargets",
                                   "baseline")})
    data["baseline"]["provides"] = {"commands": ["shared"], "packages": ["libshared"]}
    data["features"] = [
        feature("runtime", "runtime", versions=["1.0.0", "2.0.0"],
                tools=[tool("runtime", commands=["rt"])]),
        feature("kit/a", "kit", versions=["1.0", "1.1"], requires=["runtime"],
                tools=[tool("a", commands=["a", "shared"],
                            also=[{"kind": "command", "name": "c-cmd",
                                   "targets": ["debian-13-amd64"]}])]),
        feature("kit/b", "kit", depends=["kit/a"],
                tools=[tool("b", commands=["b"], version="3.0",
                            checks=[{"kind": "package", "package": "libb", "version": "3.0-1"},
                                    {"kind": "package", "package": "libshared",
                                     "version": "1"}])]),
        feature("kit/c", "kit", tools=[tool("c", commands=["c-cmd"], version="1.0")]),
        feature("kit/d", "kit", tools=[tool("d", commands=[], version="1.0",
                                            checks=[{"kind": "texlive", "package": "dpkg"}])]),
        feature("kit/all", "kit", role="bundle", depends=["kit/a", "kit/b", "kit/c", "kit/d"]),
        feature("up/tool", "uptool", owner="upstream",
                tools=[tool("up", commands=["up"], version="9.9.9", method="upstream-feature")]),
    ]
    data["groups"] = [
        {"id": "runtime", "shape": "single", "features": ["runtime"], "v0Recipes": []},
        {"id": "kit", "shape": "semantic", "bundle": "kit/all", "v0Recipes": [],
         "features": ["kit/a", "kit/b", "kit/c", "kit/d", "kit/all"]},
        {"id": "uptool", "shape": "single", "features": ["up/tool"], "v0Recipes": []},
    ]
    data["removed"] = []
    data["compatibilityNotes"] = []
    data["smoke"] = [{"id": "SMOKE-KIT", "features": ["kit/b"],
                      "fixture": "fixtures/v1-05-smoke/hello.dot", "argv": ["b", "x"],
                      "produces": "x.out"}]
    return Catalog(data)

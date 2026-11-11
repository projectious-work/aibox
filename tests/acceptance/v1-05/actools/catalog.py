"""Read-only view of the frozen V1-05 interface catalog.

Every function here is pure: it derives selections, references and expected
values from the catalog data and never touches Docker or the network.
"""

from __future__ import annotations

import json
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable, Mapping

REPO = Path(__file__).resolve().parents[4]
SPEC = REPO / "spec" / "v1"
CATALOG_PATH = SPEC / "catalog" / "v1-05-features.json"
TEST_REGISTRY = "localhost:5000/aibox-features"
UNSUPPORTED_VERSION = "0.0.0-unsupported"


@dataclass(frozen=True)
class Catalog:
    data: Mapping

    @classmethod
    def load(cls, path: Path = CATALOG_PATH) -> "Catalog":
        return cls(json.loads(path.read_text(encoding="utf-8")))

    # --- plain lookups -------------------------------------------------
    @property
    def features(self) -> dict[str, Mapping]:
        return {feature["id"]: feature for feature in self.data["features"]}

    def feature(self, feature_id: str) -> Mapping:
        try:
            return self.features[feature_id]
        except KeyError as exc:
            raise KeyError(f"unknown catalog Feature {feature_id}") from exc

    @property
    def targets(self) -> dict[str, Mapping]:
        return {target["id"]: target for target in self.data["targets"]}

    def target_for_architecture(self, architecture: str) -> str:
        for target in self.data["targets"]:
            if target["architecture"] == architecture:
                return target["id"]
        raise KeyError(f"no catalog target for architecture {architecture}")

    @property
    def groups(self) -> list[Mapping]:
        return list(self.data["groups"])

    @property
    def baseline_commands(self) -> frozenset[str]:
        return frozenset(self.data["baseline"]["provides"]["commands"])

    @property
    def baseline_packages(self) -> frozenset[str]:
        return frozenset(self.data["baseline"]["provides"]["packages"])

    def is_aibox(self, feature_id: str) -> bool:
        return self.feature(feature_id)["owner"] == "aibox"

    def is_bundle(self, feature_id: str) -> bool:
        return self.feature(feature_id)["role"] == "bundle"

    def version_option(self, feature_id: str) -> Mapping | None:
        return self.feature(feature_id)["options"].get("version")

    def semantic_members(self) -> list[tuple[Mapping, str]]:
        """(group, member Feature ID) for every non-bundle semantic member."""
        members = []
        for group in self.groups:
            if group["shape"] != "semantic":
                continue
            for feature_id in group["features"]:
                if feature_id != group["bundle"]:
                    members.append((group, feature_id))
        return members

    # --- dependency closures -----------------------------------------
    def depends_closure(self, feature_ids: Iterable[str]) -> set[str]:
        """Selected Features plus everything the Dev Container CLI installs
        automatically through `dependsOn`, transitively."""
        closure: set[str] = set()
        pending = list(feature_ids)
        while pending:
            current = pending.pop()
            if current in closure:
                continue
            closure.add(current)
            pending.extend(self.feature(current)["dependsOn"])
        return closure

    def install_closure(self, feature_ids: Iterable[str]) -> set[str]:
        """Every Feature a selection installs or must find installed: the
        transitive closure over `dependsOn` and `requires`."""
        closure: set[str] = set()
        pending = list(feature_ids)
        while pending:
            current = pending.pop()
            if current in closure:
                continue
            closure.add(current)
            feature = self.feature(current)
            pending.extend(feature["dependsOn"])
            pending.extend(feature["requires"])
        return closure

    def required_additions(self, feature_ids: Iterable[str]) -> set[str]:
        """Features to add explicitly so every `requires` in the install
        closure is satisfied (contract question Q1: closure, not only the
        explicitly selected Features)."""
        selected = set(feature_ids)
        while True:
            needed = set()
            for feature_id in self.depends_closure(selected):
                needed.update(self.feature(feature_id)["requires"])
            missing = needed - self.depends_closure(selected)
            if not missing:
                return selected - set(feature_ids)
            selected |= missing

    def with_requires(self, selection: Mapping[str, Mapping]) -> dict[str, dict]:
        """Add the required Features at their defaults (empty options)."""
        expanded = {feature_id: dict(options) for feature_id, options in selection.items()}
        for feature_id in sorted(self.required_additions(selection)):
            expanded.setdefault(feature_id, {})
        return expanded

    # --- references and versions -------------------------------------
    def feature_reference(self, feature_id: str, registry: str = TEST_REGISTRY) -> str:
        feature = self.feature(feature_id)
        if feature["owner"] == "upstream":
            return feature["upstream"]["ref"]
        major = feature["featureVersion"].split(".")[0]
        return f"{registry}/{feature_id}:{major}"

    def reference_options(self, feature_id: str, options: Mapping) -> dict:
        feature = self.feature(feature_id)
        if feature["owner"] == "upstream":
            if options:
                raise ValueError(f"{feature_id}: upstream Features take only catalog options")
            return dict(feature["upstream"]["options"])
        return dict(options)

    def devcontainer_features(self, selection: Mapping[str, Mapping],
                              registry: str = TEST_REGISTRY) -> dict[str, dict]:
        """The `features` object for devcontainer.json."""
        return {self.feature_reference(feature_id, registry):
                self.reference_options(feature_id, options)
                for feature_id, options in sorted(selection.items())}

    def selected_version(self, feature_id: str, tool: Mapping,
                         options: Mapping | None = None) -> str | None:
        """Chapter 8 Checks: the Feature's `version` option value (selected or
        default), otherwise the tool's `version`."""
        option = self.version_option(feature_id)
        if option is not None:
            return (options or {}).get("version", option["default"])
        return tool.get("version")

    def non_default_versions(self, feature_id: str) -> list[str]:
        option = self.version_option(feature_id)
        if option is None:
            return []
        return [value for value in option["enum"] if value != option["default"]]

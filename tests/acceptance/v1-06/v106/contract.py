"""Read-only access to the frozen contract spec/settings-and-logging/v1/.

Every identifier, value, default, path and pattern the journeys use comes from
here, so a journey never restates the catalog by hand.
"""

from __future__ import annotations

import json
import re
from functools import cached_property
from pathlib import Path

SUITE = Path(__file__).resolve().parent.parent
REPO = SUITE.parent.parent.parent
SPEC = REPO / "spec" / "settings-and-logging" / "v1"

# J-CFG-24: the frozen journey names spec/v1/examples/customized; #504 moved
# the example to examples/customized before the tracks started (dev-notes,
# "Acceptance harness"). The command is otherwise unchanged.
DEMO_PROJECT_RELATIVE = "examples/customized"


class Contract:
    def __init__(self, spec: Path = SPEC):
        self.spec = spec

    @cached_property
    def interfaces(self) -> dict:
        return json.loads((self.spec / "interfaces.json").read_text())

    def schema(self, name: str) -> dict:
        return json.loads((self.spec / "schemas" / name).read_text())

    @cached_property
    def result_schema(self) -> dict:
        return self.schema("operation-result.schema.json")

    @cached_property
    def event_schema(self) -> dict:
        return self.schema("log-event.schema.json")

    @cached_property
    def settings_schema(self) -> dict:
        return self.schema("process-settings.schema.json")

    @property
    def fixture_minimal(self) -> Path:
        return self.spec / "fixtures" / "project-minimal"

    @property
    def fixture_workspace(self) -> Path:
        return self.spec / "fixtures" / "project-workspace"

    @cached_property
    def registered_environment(self) -> list[str]:
        return json.loads((self.spec / "fixtures" / "registered-environment.json").read_text())["names"]

    # --- settings catalog -------------------------------------------------

    @property
    def keys(self) -> list[dict]:
        return self.interfaces["settings"]["keys"]

    @property
    def key_names(self) -> list[str]:
        return [k["key"] for k in self.keys]

    def key(self, name: str) -> dict:
        for k in self.keys:
            if k["key"] == name:
                return k
        raise KeyError(name)

    @property
    def layers(self) -> list[str]:
        return self.interfaces["settings"]["layers"]

    def default(self, name: str, mode: str = "cli"):
        """The value the view reports for an unset key (CFG-004, D5)."""
        if name == "execution.timeoutSeconds":
            op_class = self.interfaces["execution"]["operationClass"]["inspect_workspace"]
            return self.interfaces["execution"]["classDefaults"][op_class]
        return self.key(name)["default"][mode]

    @property
    def system_settings_path(self) -> str:
        return self.interfaces["settings"]["files"]["system"]["linux"]

    def source_name(self, layer: str) -> str | None:
        """effectiveConfig.source for layers whose source is a fixed string."""
        return {"built_in": self.interfaces["effectiveConfig"]["source"]["built_in"]}.get(layer)

    @property
    def native_entries(self) -> list[dict]:
        return self.interfaces["effectiveConfig"]["nativeEntries"]

    def native_entry(self, key: str, value=None) -> dict:
        for entry in self.native_entries:
            if entry["key"] == key:
                out = {k: v for k, v in entry.items() if k != "when"}
                if value is not None:
                    out["value"] = value
                return out
        raise KeyError(key)

    # --- output and logging ---------------------------------------------

    @cached_property
    def bootstrap_re(self) -> re.Pattern:
        return re.compile(self.interfaces["output"]["bootstrapDiagnostic"])

    @cached_property
    def plain_line_re(self) -> re.Pattern:
        return re.compile(self.interfaces["logging"]["plainLine"])

    @property
    def event_names(self) -> list[str]:
        return [e["event"] for e in self.interfaces["logging"]["events"]]

    def event_spec(self, name: str) -> dict:
        for e in self.interfaces["logging"]["events"]:
            if e["event"] == name:
                return e
        raise KeyError(name)

    def severity(self, name: str) -> str:
        return self.event_spec(name)["severity"]

    @property
    def severities(self) -> list[str]:
        return self.interfaces["logging"]["severities"]

    def error(self, code: str) -> dict:
        for e in self.interfaces["errors"]:
            if e["code"] == code:
                return e
        raise KeyError(code)

    @property
    def replacement(self) -> str:
        return self.interfaces["redaction"]["replacement"]

    @property
    def rejection_reasons(self) -> list[str]:
        return self.interfaces["settings"]["rejectionReasons"]

    @property
    def max_env_files(self) -> int:
        return self.interfaces["settings"]["envFile"]["maxFiles"]

    # --- requirements and journeys tables ---------------------------------

    @cached_property
    def requirement_levels(self) -> dict[str, str]:
        """Requirement ID -> verification level column of requirements.md."""
        out = {}
        for line in (self.spec / "requirements.md").read_text().splitlines():
            m = re.match(r"^\| ([A-Z]{3}-\d{3}) \|.*\| ([^|]+) \|$", line)
            if m:
                out[m.group(1)] = m.group(2).strip()
        return out

    def expand_covers(self, text: str) -> list[str]:
        """'CFG-006, CFG-013 to CFG-016' -> the IDs, keeping only table rows."""
        ids: list[str] = []
        for part in (p.strip() for p in text.split(",")):
            m = re.fullmatch(r"([A-Z]{3})-(\d{3}) to \1-(\d{3})", part)
            if m:
                prefix, lo, hi = m.group(1), int(m.group(2)), int(m.group(3))
                ids += [f"{prefix}-{n:03d}" for n in range(lo, hi + 1)
                        if f"{prefix}-{n:03d}" in self.requirement_levels]
            elif part:
                ids.append(part)
        return ids

    @cached_property
    def journeys(self) -> list[tuple[str, list[str]]]:
        """(journey ID, covered requirement IDs) in journeys.md order."""
        out = []
        for line in (self.spec / "journeys.md").read_text().splitlines():
            m = re.match(r"^\| (J-[A-Z]+-\d{2}) \| ([^|]+) \|", line)
            if m:
                out.append((m.group(1), self.expand_covers(m.group(2))))
        return out


CONTRACT = Contract()

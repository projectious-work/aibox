"""AC-TOOLS/trace: every census row maps to a catalog Feature tool or a removal
with a decision (chapter 5; chapter 8 "Census traceability")."""

from __future__ import annotations

from pathlib import Path
from typing import Mapping

from .catalog import Catalog


def row_locations(catalog: Catalog) -> dict[str, list[str]]:
    """Census row -> places it maps to (`feature:<id>/<tool>` or `removed:<decision>`)."""
    locations: dict[str, list[str]] = {}
    for feature in catalog.data["features"]:
        for tool in feature["tools"]:
            for row in tool["rows"]:
                locations.setdefault(row, []).append(f"feature:{feature['id']}:{tool['name']}")
    for removal in catalog.data["removed"]:
        locations.setdefault(removal["row"], []).append(f"removed:{removal['decision']}")
    return locations


def decision_recorded(decision: str, decisions_dir: Path) -> bool:
    """A DecisionRecord file exists whose name starts with the decision's
    `DEC-<date>_<time>-<Word>` stem."""
    stem = "-".join(decision.split("-", 3)[:3])
    return any(decisions_dir.glob(stem + "*.md"))


def trace_failures(row: str, catalog: Catalog, table: Mapping, decisions_dir: Path) -> list[str]:
    locations = row_locations(catalog).get(row, [])
    failures = []
    if len(locations) != 1:
        failures.append(f"maps to {len(locations)} catalog places: {locations}")
        return failures
    location = locations[0]
    entry = {r["id"]: r for r in table["rows"]}.get(row)
    if entry is None:
        failures.append("missing from the selection table")
    if location.startswith("feature:"):
        feature_id = location.split(":")[1]
        if entry is not None and (entry.get("decision") != "retained"
                                  or entry.get("v1-ref") != feature_id):
            failures.append(f"selection table says {entry.get('decision')} "
                            f"{entry.get('v1-ref')}, catalog says {feature_id}")
    else:
        decision = location.split(":", 1)[1]
        if not decision_recorded(decision, decisions_dir):
            failures.append(f"removal decision {decision} has no DecisionRecord")
        if entry is not None and (entry.get("decision") != "removed"
                                  or entry.get("decision-record") != decision):
            failures.append(f"selection table says {entry.get('decision')} "
                            f"{entry.get('decision-record')}, catalog removal says {decision}")
    return failures

#!/usr/bin/env python3
"""Check the V1-05 tool Feature catalog against the frozen v0 census.

The JSON Schema shape is checked by validate.mjs; this script checks the
cross-references that a schema cannot express.
"""
import json
import re
from collections import Counter
from pathlib import Path

root = Path(__file__).resolve().parents[1]
repo = root.parents[1]
read = lambda name: json.loads((root / name).read_text())
ledger = read("ledger/addons.json")
table = read("ledger/v1-05-tool-selection.json")
audit = read("evidence/V1-05/source-audit.json")
oci = read("evidence/V1-05/oci-evidence.json")
fixture = read("fixtures/v1-05-native-features.json")
catalog = read("catalog/v1-05-features.json")


def fail(message):
    raise SystemExit(f"V1-05 catalog: {message}")


def require(condition, message):
    if not condition:
        fail(message)


def decision_exists(decision_id):
    return any((repo / "context/decisions").glob(f"{decision_id.split('-', 3)[0]}-{decision_id.split('-', 3)[1]}-{decision_id.split('-', 3)[2]}*.md"))


# Census: the 98 frozen rows stay individually traceable.
expected = {f"{recipe['id']}/{tool['name']}": recipe for recipe in ledger for tool in recipe["tools"]}
rows = {row["id"]: row for row in table["rows"]}
require(len(expected) == len(table["rows"]) == len(rows) == 98, "census row count")
require(set(rows) == set(expected), "census row identity")
require(table["baselineRevision"] == audit["baseline_revision"] == catalog["baselineRevision"], "baseline revision")
audited = {row["id"]: row for row in audit["rows"]}
require(set(audited) == set(expected), "audit coverage")
for key, row in rows.items():
    require(row["v0-source"] == expected[key]["source"], f"{key}: v0 source")
    require(row["v0-source-sha256"] == expected[key]["sha256"], f"{key}: v0 source hash")
    require(row["classification"] == audited[key]["status"], f"{key}: audit classification")

targets = {target["id"]: target for target in catalog["targets"]}
require({(t["architecture"], t["unameMachine"]) for t in targets.values()} == {("amd64", "x86_64"), ("arm64", "aarch64")},
        "targets must be exactly trixie amd64 and arm64")
error_codes = {error["code"] for error in catalog["errors"]}
for code in ("AIBOX_FEATURE_UNSUPPORTED_TARGET", "AIBOX_FEATURE_UNSUPPORTED_VERSION", "AIBOX_FEATURE_VERSION_CONFLICT",
             "AIBOX_FEATURE_DOWNLOAD_FAILED", "AIBOX_FEATURE_VERIFICATION_FAILED", "AIBOX_FEATURE_PACKAGE_UNAVAILABLE"):
    require(code in error_codes, f"missing error code {code}")

features = {feature["id"]: feature for feature in catalog["features"]}
require(len(features) == len(catalog["features"]), "duplicate Feature ID")
groups = {group["id"]: group for group in catalog["groups"]}
require(len(groups) == len(catalog["groups"]), "duplicate group ID")

# Groups and shapes.
grouped = Counter()
for group in groups.values():
    members = [features.get(fid) for fid in group["features"]]
    require(all(members), f"{group['id']}: unknown Feature")
    for member in members:
        grouped[member["id"]] += 1
        require(member["group"] == group["id"], f"{member['id']}: group mismatch")
    roles = Counter(member["role"] for member in members)
    if group["shape"] == "single":
        require(roles == Counter(tool=1) and len(members[0]["tools"]) == 1, f"{group['id']}: single shape")
    elif group["shape"] == "atomic":
        require(roles == Counter(atomic=1) and len(members[0]["tools"]) >= 2, f"{group['id']}: atomic shape")
    else:
        bundle = features.get(group.get("bundle"))
        require(bundle and bundle["role"] == "bundle" and bundle["id"] in group["features"], f"{group['id']}: bundle")
        require(bundle["id"].endswith("/all") and not bundle["tools"] and not bundle["options"], f"{bundle['id']}: bundle content")
        tools = [member["id"] for member in members if member["role"] == "tool"]
        require(roles == Counter(tool=len(tools), bundle=1) and len(tools) >= 2, f"{group['id']}: semantic shape")
        require(sorted(bundle["dependsOn"]) == sorted(tools), f"{bundle['id']}: must depend on every member")
require(set(grouped) == set(features) and set(grouped.values()) == {1}, "each Feature belongs to exactly one group")

# Features.
commands = Counter()
for fid, feature in features.items():
    require(set(feature["targets"]) <= set(targets), f"{fid}: unknown target")
    for dependency in feature["dependsOn"]:
        require(dependency in features and dependency != fid, f"{fid}: unknown dependency {dependency}")
    version = feature["options"].get("version")
    if version:
        require(version["default"] in version["enum"], f"{fid}: default outside enum")
    if feature["owner"] == "upstream":
        name = feature["upstream"]["ref"].split("/features/")[1].split("@")[0]
        require(oci[name]["manifest_digest"] == feature["upstream"]["ref"].split("@")[1], f"{fid}: upstream digest evidence")
        require(set(feature["upstream"]["options"]) <= set(oci[name]["metadata"]["options"]), f"{fid}: upstream options")
        require("featureVersion" not in feature, f"{fid}: upstream Features are not republished")
    coupling = feature.get("versionCoupling")
    if coupling:
        other = features.get(coupling["feature"])
        require(other and "version" in other["options"] and version, f"{fid}: coupling needs versioned Features")
        require(coupling["feature"] in feature["dependsOn"], f"{fid}: coupled Feature must be a dependency")
        if coupling["rule"] == "equal":
            require(version["enum"] == other["options"]["version"]["enum"], f"{fid}: equal coupling enums")
        else:
            require(set(coupling.get("table", {})) == set(version["enum"]), f"{fid}: coupling table keys")
    for tool in feature["tools"]:
        require(tool.get("versionFrom", fid) in features, f"{fid}/{tool['name']}: unknown versionFrom")
        require(set(tool.get("absenceTargets", [])) <= set(feature["targets"]), f"{fid}/{tool['name']}: absence target")
        for command in tool.get("commands", []):
            commands[command] += 1
        for check in tool["checks"]:
            if check["kind"] == "version":
                require("{version}" in check["regex"], f"{fid}/{tool['name']}: version regex placeholder")
                re.compile(check["regex"].replace("{version}", "0"))
                require(version or tool.get("versionFrom") or tool.get("version"), f"{fid}/{tool['name']}: no version source")
            elif check["kind"] == "command":
                re.compile(check["regex"])
                require("{version}" not in check["regex"], f"{fid}/{tool['name']}: command checks are version-free")
            elif check["kind"] == "package":
                require(check.get("version"), f"{fid}/{tool['name']}: package check needs exact version")
require(not [c for c, n in commands.items() if n > 1], f"command owned twice: {[c for c, n in commands.items() if n > 1]}")

# Dependency cycles.
def visit(fid, path):
    require(fid not in path, f"dependency cycle through {fid}")
    for dependency in features[fid]["dependsOn"]:
        visit(dependency, path | {fid})
for fid in features:
    visit(fid, frozenset())

# Census rows map exactly once to a tool or a removal.
mapped = Counter()
for fid, feature in features.items():
    for tool in feature["tools"]:
        for row in tool["rows"]:
            mapped[row] += 1
            require(rows[row]["v1-ref"] == fid and rows[row]["decision"] == "retained", f"{row}: selection table must point to {fid}")
removed = {entry["row"]: entry for entry in catalog["removed"]}
for row, entry in removed.items():
    mapped[row] += 1
    require(rows[row]["decision"] == "removed" and rows[row]["v1-ref"] is None, f"{row}: selection table removal")
    require(rows[row]["decision-record"] == entry["decision"], f"{row}: removal decision mismatch")
require(set(mapped) == set(rows) and set(mapped.values()) == {1}, "every census row maps exactly once")

decisions = {entry["decision"] for entry in catalog["removed"]} | {note["decision"] for note in catalog["compatibilityNotes"]}
decisions |= {row["decision-record"] for row in rows.values()}
for decision in decisions:
    require(decision_exists(decision), f"missing decision record {decision}")
for note in catalog["compatibilityNotes"]:
    require(set(note["rows"]) <= set(rows), f"{note['id']}: unknown row")

for smoke in catalog.get("smoke", []):
    require((root / smoke["fixture"]).exists(), f"{smoke['id']}: missing fixture")
    require(set(smoke["features"]) <= set(features), f"{smoke['id']}: unknown Feature")

# The maintained example selects only reviewed upstream references with catalog options.
upstream_refs = {f["upstream"]["ref"]: f["upstream"]["options"] for f in features.values() if f["owner"] == "upstream"}
for ref, options in fixture["selected"].items():
    require(upstream_refs.get(ref) == options, f"example Feature not in catalog: {ref}")

counts = Counter(f["role"] for f in features.values())
print(f"V1-05 catalog: {len(rows)} census rows ({len(rows) - len(removed)} retained, {len(removed)} removed), "
      f"{len(groups)} groups, {len(features)} Features ({counts['tool']} tool, {counts['atomic']} atomic, {counts['bundle']} bundle)")

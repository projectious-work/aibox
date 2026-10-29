"""Expand the AC-TOOLS items of chapter 5 into builds and named tests.

Pure functions over the catalog. A Build is one workspace built from a
fixture with an exact Feature selection; tests that define the identical
selection share it (chapter 8 allows sharing for presence checks; for the
other items an identical selection is still exactly the selection the item
defines). Test IDs follow chapter 5:
`AC-TOOLS/<item> [<feature-id> [<tool>]] [<target-id>] [<key>=<value> ...]`
with the per-item slots of chapter 5.
"""

from __future__ import annotations

import fnmatch
import hashlib
import json
from dataclasses import dataclass, field, replace
from typing import Callable, Iterable, Mapping

from . import checks
from .catalog import TEST_REGISTRY, UNSUPPORTED_VERSION, Catalog

ITEMS = (
    "present", "version", "absent", "isolated", "bundle", "depends", "requires",
    "unsupported-version", "unsupported-target", "conflict", "side-effects",
    "smoke", "stack", "inspect", "trace",
)
DOCKER_ITEMS = frozenset(ITEMS) - {"inspect", "trace"}
INSPECT_EXAMPLES = ("customized", "minimal")
BASE_FIXTURE = "baseline"


@dataclass(frozen=True)
class Build:
    key: str
    fixture: str                       # "baseline" or an unsupportedTargets ID
    fixture_path: str                  # relative to spec/v1
    selection: tuple[tuple[str, str], ...]  # (Feature ID, options as JSON)
    copies: tuple[str, ...] = ()       # smoke fixtures, relative to spec/v1
    expect_failure: bool = False

    def selection_dict(self) -> dict[str, dict]:
        return {feature_id: json.loads(options) for feature_id, options in self.selection}


@dataclass(frozen=True)
class Test:
    id: str
    item: str
    target: str | None
    feature: str | None = None
    tool: str | None = None
    qualifiers: tuple[tuple[str, str], ...] = ()
    subjects: tuple[str, ...] = ()     # IDs that --feature can match
    build: str | None = None
    expectations: tuple[checks.Expectation, ...] = ()
    error: Mapping | None = None       # code, features, mustName
    payload: Mapping = field(default_factory=dict)


@dataclass
class Exemptions:
    """What installed tools' `alsoProvides` exempt from absence checks."""
    commands: set = field(default_factory=set)
    packages: set = field(default_factory=set)
    features: set = field(default_factory=set)


def test_id(item: str, *parts: str | None, qualifiers: Iterable[tuple[str, str]] = ()) -> str:
    tokens = [f"AC-TOOLS/{item}", *[part for part in parts if part]]
    tokens += [f"{key}={value}" for key, value in qualifiers]
    return " ".join(tokens)


class Planner:
    def __init__(self, catalog: Catalog, registry: str = TEST_REGISTRY):
        self.catalog = catalog
        self.registry = registry
        self.builds: dict[str, Build] = {}
        self.tests: list[Test] = []

    # --- builds --------------------------------------------------------
    def build(self, selection: Mapping[str, Mapping], *, fixture: str = BASE_FIXTURE,
              copies: Iterable[str] = (), expect_failure: bool = False) -> str:
        if fixture == BASE_FIXTURE:
            fixture_path = self.catalog.data["baseline"]["fixture"]
        else:
            fixture_path = {u["id"]: u["fixture"]
                            for u in self.catalog.data["unsupportedTargets"]}[fixture]
        normalized = tuple(sorted((feature_id, json.dumps(dict(options), sort_keys=True))
                                  for feature_id, options in selection.items()))
        identity = json.dumps([fixture, normalized, sorted(copies), expect_failure])
        key = hashlib.sha256(identity.encode()).hexdigest()[:12]
        self.builds.setdefault(key, Build(key, fixture, fixture_path, normalized,
                                          tuple(sorted(copies)), expect_failure))
        return key

    def add(self, test: Test) -> None:
        self.tests.append(test)

    # --- item expansions ----------------------------------------------
    def presence(self, feature_id: str, options: Mapping | None = None, *,
                 include_commands: bool = True) -> list[tuple[Mapping, list[checks.Expectation]]]:
        result = []
        for tool in self.catalog.feature(feature_id)["tools"]:
            version = self.catalog.selected_version(feature_id, tool, options)
            expectations = checks.presence_expectations(tool, version,
                                                        include_commands=include_commands)
            scope = f"{feature_id} {tool['name']}"
            result.append((tool, [replace(e, scope=scope) for e in expectations]))
        return result

    def all_presence(self, feature_ids: Iterable[str]) -> list[checks.Expectation]:
        return [e for feature_id in feature_ids for _, expectations in self.presence(feature_id)
                for e in expectations]

    def absence(self, feature_id: str, exempt: "Exemptions") -> list[checks.Expectation]:
        """Absence observables of a Feature that are neither exempt nor in
        `baseline.provides` (chapter 8 Reference base and absence)."""
        if feature_id in exempt.features:
            return []
        result = []
        for tool in self.catalog.feature(feature_id)["tools"]:
            expectations = checks.absence_expectations(
                tool, skip_commands=self.catalog.baseline_commands | exempt.commands,
                skip_packages=self.catalog.baseline_packages | exempt.packages)
            result += [replace(e, scope=f"{feature_id} {tool['name']}") for e in expectations]
        return result

    def supported(self, feature_id: str, target: str) -> bool:
        return target in self.catalog.feature(feature_id)["targets"]

    def plan_present(self, target: str) -> None:
        for feature_id, feature in self.catalog.features.items():
            if feature["role"] == "bundle" or not self.supported(feature_id, target):
                continue  # chapter 5: present is per non-bundle Feature
            key = self.build(self.catalog.with_requires({feature_id: {}}))
            for tool, expectations in self.presence(feature_id):
                self.add(Test(test_id("present", feature_id, tool["name"], target), "present",
                              target, feature_id, tool["name"], subjects=(feature_id,),
                              build=key, expectations=tuple(expectations)))

    def plan_version(self, target: str) -> None:
        for feature_id in self.catalog.features:
            if not self.supported(feature_id, target):
                continue
            for version in self.catalog.non_default_versions(feature_id):
                options = {"version": version}
                key = self.build(self.catalog.with_requires({feature_id: options}))
                for tool, expectations in self.presence(feature_id, options,
                                                        include_commands=False):
                    qualifiers = (("version", version),)
                    self.add(Test(test_id("version", feature_id, tool["name"], target,
                                          qualifiers=qualifiers), "version", target,
                                  feature_id, tool["name"], qualifiers, (feature_id,), key,
                                  tuple(expectations)))

    def plan_absent(self, target: str) -> None:
        expectations = [e for feature_id in self.catalog.features
                        for e in self.absence(feature_id, Exemptions())]
        self.add(Test(test_id("absent", target), "absent", target,
                      subjects=tuple(self.catalog.features), build=self.build({}),
                      expectations=tuple(expectations)))

    def exemptions(self, installed: Iterable[str], target: str) -> "Exemptions":
        exempt = Exemptions()
        for feature_id in installed:
            for tool in self.catalog.feature(feature_id)["tools"]:
                for entry in tool.get("alsoProvides", []):
                    if target in entry["targets"]:
                        {"command": exempt.commands, "package": exempt.packages,
                         "feature": exempt.features}[entry["kind"]].add(entry["name"])
        return exempt

    def plan_isolated(self, target: str) -> None:
        for group, member in self.catalog.semantic_members():
            if not self.supported(member, target):
                continue
            selection = self.catalog.with_requires({member: {}})
            exempt = self.exemptions(self.catalog.install_closure(selection), target)
            closure = self.catalog.install_closure([member])
            siblings = [s for s in group["features"]
                        if s not in (member, group["bundle"]) and s not in closure]
            expectations = [e for sibling in siblings for e in self.absence(sibling, exempt)]
            if not expectations:
                continue  # every sibling observable is exempt or a dependency
            key = self.build(selection)
            for tool in self.catalog.feature(member)["tools"]:
                self.add(Test(test_id("isolated", member, tool["name"], target), "isolated",
                              target, member, tool["name"], subjects=(member,), build=key,
                              expectations=tuple(expectations),
                              payload={"siblings": siblings}))

    def plan_bundle(self, target: str) -> None:
        for group in self.catalog.groups:
            if group["shape"] != "semantic" or not self.supported(group["bundle"], target):
                continue
            bundle = group["bundle"]
            members = [m for m in group["features"] if m != bundle]
            key = self.build(self.catalog.with_requires({bundle: {}}))
            self.add(Test(test_id("bundle", bundle, target), "bundle", target, bundle,
                          subjects=(bundle, *members), build=key,
                          expectations=tuple(self.all_presence(members))))

    def plan_depends(self, target: str) -> None:
        for feature_id, feature in self.catalog.features.items():
            if not feature["dependsOn"] or not self.supported(feature_id, target):
                continue
            key = self.build(self.catalog.with_requires({feature_id: {}}))
            for dependency in feature["dependsOn"]:
                for tool, expectations in self.presence(dependency):
                    self.add(Test(test_id("depends", feature_id, tool["name"], target),
                                  "depends", target, feature_id, tool["name"],
                                  subjects=(feature_id, dependency), build=key,
                                  expectations=tuple(expectations),
                                  payload={"dependency": dependency}))

    def error_test(self, item: str, feature_id: str, target: str | None, key: str, code: str,
                   features: Iterable[str], must_name: Iterable[str] | None = None,
                   qualifiers: tuple[tuple[str, str], ...] = (),
                   subjects: tuple[str, ...] = (), host_target: str | None = None) -> None:
        self.add(Test(test_id(item, feature_id, target, qualifiers=qualifiers), item,
                      host_target or target, feature_id, None, qualifiers,
                      subjects or (feature_id,), key,
                      error={"code": code, "features": sorted(features),
                             "mustName": sorted(must_name) if must_name else None}))

    def plan_requires(self, target: str) -> None:
        for feature_id, feature in self.catalog.features.items():
            if feature["owner"] != "aibox" or not feature["requires"]:
                continue
            if not self.supported(feature_id, target):
                continue
            selection = self.catalog.with_requires({feature_id: {}})
            for required in feature["requires"]:
                selection.pop(required, None)
            missing = set(feature["requires"])
            if missing & self.catalog.depends_closure(selection):
                raise ValueError(f"{feature_id}: a dependsOn Feature installs what it requires")
            key = self.build(selection, expect_failure=True)
            qualifiers = tuple(("missing", required) for required in sorted(missing))
            self.error_test("requires", feature_id, target, key,
                            "AIBOX_FEATURE_DEPENDENCY_MISSING", [feature_id], missing,
                            qualifiers)

    def plan_unsupported_version(self, target: str) -> None:
        for feature_id, feature in self.catalog.features.items():
            if (feature["owner"] != "aibox" or self.catalog.version_option(feature_id) is None
                    or not self.supported(feature_id, target)):
                continue
            key = self.build(self.catalog.with_requires(
                {feature_id: {"version": UNSUPPORTED_VERSION}}), expect_failure=True)
            self.error_test("unsupported-version", feature_id, target, key,
                            "AIBOX_FEATURE_UNSUPPORTED_VERSION", [feature_id])

    def plan_unsupported_target(self, target: str) -> None:
        # The test ID has a fixture slot instead of a target: one test per
        # Feature x fixture, run on whichever host target executes the plan.
        if any(test.item == "unsupported-target" for test in self.tests):
            return
        for feature_id, feature in self.catalog.features.items():
            if feature["owner"] != "aibox":
                continue
            for unsupported in self.catalog.data["unsupportedTargets"]:
                key = self.build(self.catalog.with_requires({feature_id: {}}),
                                 fixture=unsupported["id"], expect_failure=True)
                self.error_test("unsupported-target", feature_id, None, key,
                                "AIBOX_FEATURE_UNSUPPORTED_TARGET",
                                self.catalog.install_closure([feature_id]),
                                qualifiers=(("fixture", unsupported["id"]),),
                                host_target=target)

    def plan_conflict(self, target: str) -> None:
        for group, member in self.catalog.semantic_members():
            versions = self.catalog.non_default_versions(member)
            if not versions or not self.supported(member, target):
                continue
            version = versions[0]  # chapter 5: the first non-default enum value
            bundle = group["bundle"]
            key = self.build(self.catalog.with_requires(
                {bundle: {}, member: {"version": version}}), expect_failure=True)
            self.error_test("conflict", member, target, key, "AIBOX_FEATURE_VERSION_CONFLICT",
                            [member], qualifiers=(("version", version),),
                            subjects=(member, bundle))

    def plan_side_effects(self, target: str) -> None:
        for feature_id, feature in self.catalog.features.items():
            if feature["owner"] != "aibox" or not self.supported(feature_id, target):
                continue
            key = self.build(self.catalog.with_requires({feature_id: {}}))
            self.add(Test(test_id("side-effects", feature_id, target), "side-effects", target,
                          feature_id, subjects=(feature_id,), build=key,
                          payload={"containerEnv": dict(feature.get("containerEnv", {})),
                                   "featureVersion": feature["featureVersion"]}))

    def plan_smoke(self, target: str) -> None:
        for smoke in self.catalog.data.get("smoke", []):
            if not all(self.supported(f, target) for f in smoke["features"]):
                continue
            key = self.build(self.catalog.with_requires({f: {} for f in smoke["features"]}),
                             copies=[smoke["fixture"]])
            self.add(Test(test_id("smoke", smoke["id"], target), "smoke", target, smoke["id"],
                          subjects=(smoke["id"], *smoke["features"]), build=key,
                          payload={"argv": list(smoke["argv"]), "produces": smoke["produces"],
                                   "magicHex": smoke.get("producesMagicHex")}))

    def plan_stack(self, target: str) -> None:
        for stack_id, selected in known_risk_stacks(self.catalog):
            selection = self.catalog.with_requires({f: {} for f in selected})
            installed = sorted(f for f in self.catalog.depends_closure(selection)
                               if not self.catalog.is_bundle(f))
            self.add(Test(test_id("stack", stack_id, target), "stack", target, stack_id,
                          subjects=(stack_id, *selected), build=self.build(selection),
                          expectations=tuple(self.all_presence(installed)),
                          payload={"features": installed}))

    def plan_inspect(self) -> None:
        for example in INSPECT_EXAMPLES:
            self.add(Test(test_id("inspect", qualifiers=(("example", example),)), "inspect", None,
                          payload={"example": example}))

    def plan_trace(self, census_rows: Iterable[str]) -> None:
        for row in census_rows:
            self.add(Test(test_id("trace", row), "trace", None, payload={"row": row}))


def known_risk_stacks(catalog: Catalog) -> list[tuple[str, list[str]]]:
    """Chapter 5 "Known-risk stacks" with their stack IDs.
    `node/node`, `go/go` and `rust/rustc` are added by with_requires where a
    member requires them."""
    uv_tools = sorted(feature_id for feature_id, feature in catalog.features.items()
                      if any(tool["method"] == "uv-tool" for tool in feature["tools"]))
    return [
        ("STACK-PYTHON", ["python/all", *uv_tools]),
        ("STACK-NODE", ["node/all", "typescript/all", "mermaid", "data/visualization",
                      "browser-testing/all"]),
        ("STACK-DOCS", ["latex/all", "preview/archive/all", "preview/enhanced/all"]),
        ("STACK-NATIVE", ["go/all", "rust/all", "zig/all"]),
        ("STACK-SECURITY", ["security/crypto/all", "security/scan/all"]),
    ]


PLANNERS: dict[str, Callable[[Planner, str], None]] = {
    "present": Planner.plan_present,
    "version": Planner.plan_version,
    "absent": Planner.plan_absent,
    "isolated": Planner.plan_isolated,
    "bundle": Planner.plan_bundle,
    "depends": Planner.plan_depends,
    "requires": Planner.plan_requires,
    "unsupported-version": Planner.plan_unsupported_version,
    "unsupported-target": Planner.plan_unsupported_target,
    "conflict": Planner.plan_conflict,
    "side-effects": Planner.plan_side_effects,
    "smoke": Planner.plan_smoke,
    "stack": Planner.plan_stack,
}


@dataclass(frozen=True)
class Plan:
    builds: dict[str, Build]
    tests: list[Test]


def normalize_items(items: Iterable[str] | None) -> list[str]:
    if not items:
        return list(ITEMS)
    result = []
    for raw in items:
        for item in raw.split(","):
            item = item.strip().removeprefix("AC-TOOLS/")
            if item not in ITEMS:
                raise ValueError(f"unknown AC-TOOLS item {item!r}; known: {', '.join(ITEMS)}")
            if item not in result:
                result.append(item)
    return result


def make_plan(catalog: Catalog, targets: Iterable[str], census_rows: Iterable[str], *,
              items: Iterable[str] | None = None, features: Iterable[str] = (),
              tests: Iterable[str] = (), registry: str = TEST_REGISTRY) -> Plan:
    selected_items = normalize_items(items)
    planner = Planner(catalog, registry)
    for target in targets:
        for item in selected_items:
            if item in PLANNERS:
                PLANNERS[item](planner, target)
    if "inspect" in selected_items:
        planner.plan_inspect()
    if "trace" in selected_items:
        planner.plan_trace(census_rows)
    chosen = [test for test in planner.tests
              if matches(test, list(features), list(tests))]
    ids = [test.id for test in chosen]
    duplicates = sorted({test_id for test_id in ids if ids.count(test_id) > 1})
    if duplicates:
        raise ValueError(f"duplicate test IDs: {duplicates[:5]}")
    used = {test.build for test in chosen if test.build}
    return Plan({key: build for key, build in planner.builds.items() if key in used}, chosen)


def matches(test: Test, feature_patterns: list[str], test_patterns: list[str]) -> bool:
    if feature_patterns and not any(fnmatch.fnmatchcase(subject, pattern)
                                    for pattern in feature_patterns
                                    for subject in test.subjects):
        return False
    if test_patterns and not any(fnmatch.fnmatchcase(test.id, pattern)
                                 for pattern in test_patterns):
        return False
    return True


def census_rows(addons: Iterable[Mapping]) -> list[str]:
    """The 98 census rows: `<recipe ID>/<tool name>` from ledger/addons.json."""
    return [f"{recipe['id']}/{tool['name']}" for recipe in addons for tool in recipe["tools"]]

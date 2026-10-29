"""Chapter 8 "Checks" as data: probes run in a container and pure evaluation.

A probe is an argv run in the container; its exit code and combined
standard output and error are recorded. An expectation pairs a probe with
the rule that decides whether it passes. The same probes run once as the
remote user (devcontainer exec) and once as root (docker exec --user root).
"""

from __future__ import annotations

import base64
import re
import shlex
from dataclasses import dataclass
from typing import Iterable, Mapping

PROBE_TIMEOUT_SECONDS = 600
MARKER = "@@aibox-v1-05-probe"
ECMASCRIPT_SYNTAX = set("\\^$.*+?()[]{}|/")


@dataclass(frozen=True)
class Probe:
    argv: tuple[str, ...]


@dataclass(frozen=True)
class Outcome:
    exit_code: int
    output: str


@dataclass(frozen=True)
class Expectation:
    """`rule` is one of: version, command, package, texlive, path (presence
    kinds of chapter 8) or command-absent, command-check-absent,
    package-absent, path-absent (absence observables). `scope` names the
    Feature and tool the expectation belongs to in aggregate tests."""
    rule: str
    probe: Probe
    subject: str
    regex: str | None = None
    expected: str | None = None
    scope: str = ""

    def describe(self) -> str:
        prefix = f"{self.scope}: " if self.scope else ""
        return f"{prefix}{self.rule} {self.subject}"


def ecmascript_escape(text: str) -> str:
    """Escape a literal for an ECMAScript regular expression."""
    return "".join("\\" + char if char in ECMASCRIPT_SYNTAX else char for char in text)


def compile_ecmascript(pattern: str) -> re.Pattern:
    """Compile the catalog's ECMAScript subset with `^`/`$` at line boundaries."""
    return re.compile(pattern, re.MULTILINE)


def command_probe(name: str) -> Probe:
    return Probe(("sh", "-c", 'command -v "$1"', "sh", name))


def package_version_probe(package: str) -> Probe:
    return Probe(("dpkg-query", "-W", "-f", "${Version}", package))


def package_status_probe(package: str) -> Probe:
    return Probe(("dpkg-query", "-W", "-f", "${db:Status-Status}", package))


def texlive_probe(package: str) -> Probe:
    return Probe(("tlmgr", "info", "--only-installed", "--data", "name", package))


def path_probe(path: str) -> Probe:
    return Probe(("test", "-d", path))


def check_expectation(check: Mapping, version: str | None, tool_name: str) -> Expectation:
    kind = check["kind"]
    if kind in ("version", "command"):
        regex = check["regex"]
        if kind == "version":
            if "{version}" in regex and version is None:
                raise ValueError(f"{tool_name}: version check without a selected version")
            regex = regex.replace("{version}", ecmascript_escape(version or ""))
        return Expectation(kind, Probe(tuple(check["argv"])), " ".join(check["argv"]), regex=regex)
    if kind == "package":
        return Expectation(kind, package_version_probe(check["package"]), check["package"],
                           expected=check["version"])
    if kind == "texlive":
        return Expectation(kind, texlive_probe(check["package"]), check["package"],
                           expected=check["package"])
    if kind == "path":
        return Expectation(kind, path_probe(check["path"]), check["path"])
    raise ValueError(f"{tool_name}: unknown check kind {kind}")


def presence_expectations(tool: Mapping, version: str | None, *,
                          include_commands: bool = True) -> list[Expectation]:
    expectations = [check_expectation(check, version, tool["name"]) for check in tool["checks"]]
    if include_commands:
        expectations += [Expectation("command-present", command_probe(name), name)
                         for name in tool.get("commands", [])]
    return expectations


def absence_expectations(tool: Mapping, *, skip_commands: Iterable[str] = (),
                         skip_packages: Iterable[str] = ()) -> list[Expectation]:
    """Chapter 8 absence observables: `commands` entries and `command`,
    `package` and `path` checks; `version` and `texlive` checks are not."""
    skip_commands, skip_packages = set(skip_commands), set(skip_packages)
    expectations = [Expectation("command-absent", command_probe(name), name)
                    for name in tool.get("commands", []) if name not in skip_commands]
    for check in tool["checks"]:
        if check["kind"] == "package" and check["package"] not in skip_packages:
            expectations.append(Expectation("package-absent",
                                            package_status_probe(check["package"]),
                                            check["package"]))
        elif check["kind"] == "path":
            expectations.append(Expectation("path-absent", path_probe(check["path"]),
                                            check["path"]))
        elif check["kind"] == "command":
            expectations.append(Expectation("command-check-absent", Probe(tuple(check["argv"])),
                                            " ".join(check["argv"]), regex=check["regex"]))
    return expectations


def passes(expectation: Expectation, outcome: Outcome) -> bool:
    rule = expectation.rule
    if rule in ("version", "command"):
        return outcome.exit_code == 0 and bool(
            compile_ecmascript(expectation.regex or "").search(outcome.output))
    if rule == "package":
        return outcome.exit_code == 0 and outcome.output == expectation.expected
    if rule == "texlive":
        return expectation.expected in (line.strip() for line in outcome.output.splitlines())
    if rule == "path":
        return outcome.exit_code == 0
    if rule == "command-present":
        return outcome.exit_code == 0 and bool(outcome.output.strip())
    if rule == "command-absent":
        return outcome.exit_code != 0
    if rule == "package-absent":
        return not (outcome.exit_code == 0 and outcome.output.strip() == "installed")
    if rule == "path-absent":
        return outcome.exit_code != 0
    if rule == "command-check-absent":
        return not (outcome.exit_code == 0 and bool(
            compile_ecmascript(expectation.regex or "").search(outcome.output)))
    raise ValueError(f"unknown rule {rule}")


def probe_script(probes: Iterable[Probe]) -> str:
    """POSIX sh script that runs each probe with standard error merged into
    standard output and prints one marker line per probe."""
    lines = [
        "probe() {",
        '  index=$1; shift',
        f'  out=$(timeout -k 10 {PROBE_TIMEOUT_SECONDS} "$@" 2>&1); rc=$?',
        f'  printf \'{MARKER} %s %s \' "$index" "$rc"',
        "  printf '%s' \"$out\" | base64 | tr -d '\\n'",
        "  printf '\\n'",
        "}",
    ]
    for index, probe in enumerate(probes):
        lines.append("probe " + str(index) + " " + " ".join(shlex.quote(arg) for arg in probe.argv))
    return "\n".join(lines) + "\n"


def parse_probe_output(text: str, probes: list[Probe]) -> dict[Probe, Outcome]:
    outcomes: dict[Probe, Outcome] = {}
    for line in text.splitlines():
        if not line.startswith(MARKER + " "):
            continue
        parts = line.split(" ", 3)
        index, code = int(parts[1]), int(parts[2])
        payload = parts[3] if len(parts) > 3 else ""
        outcomes[probes[index]] = Outcome(code, base64.b64decode(payload).decode("utf-8", "replace"))
    missing = [probe for probe in probes if probe not in outcomes]
    if missing:
        raise RuntimeError(f"probe output incomplete: {len(missing)} of {len(probes)} missing")
    return outcomes


def error_line_matches(output: str, code: str, feature_ids: Iterable[str],
                       must_name: Iterable[str] | None = None) -> bool:
    """Chapter 8 Errors: some line of the combined output contains
    `aibox-feature error <CODE> feature=<catalog-id>`, regardless of prefix.
    The catalog ID must end there (`feature=go/go` does not match
    `feature=go/goimports`). With `must_name`, the detail after that prefix
    also names one of those Feature IDs."""
    names = list(must_name) if must_name else []
    for line in output.splitlines():
        for feature_id in feature_ids:
            prefix = f"aibox-feature error {code} feature={feature_id}"
            position = line.find(prefix)
            if position < 0:
                continue
            rest = line[position + len(prefix):]
            if rest and not rest[0].isspace():
                continue  # feature=go/go must not match feature=go/goimports
            if not names or any(re.search(r"(^|[^A-Za-z0-9_-])" + re.escape(name)
                                          + r"($|[^A-Za-z0-9/_-])", rest) for name in names):
                return True
    return False

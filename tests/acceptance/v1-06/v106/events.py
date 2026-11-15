"""Parsing of stderr: JSON Lines events, plain event lines and the bootstrap
diagnostic of OUT-007."""

from __future__ import annotations

import json
import re
from dataclasses import dataclass, field

from .contract import CONTRACT


def text_lines(data: bytes | str) -> list[str]:
    """Split output into lines; a trailing newline does not add an empty line."""
    if isinstance(data, bytes):
        data = data.decode("utf-8", errors="replace")
    if not data:
        return []
    lines = data.split("\n")
    if lines[-1] == "":
        lines.pop()
    return lines


@dataclass
class ParsedStderr:
    lines: list[str]
    events: list[dict] = field(default_factory=list)
    # Lines that are not a JSON object with an "event" member, with their index.
    other: list[tuple[int, str]] = field(default_factory=list)

    def names(self) -> list[str]:
        return [e.get("event") for e in self.events]

    def named(self, name: str) -> list[dict]:
        return [e for e in self.events if e.get("event") == name]

    @property
    def last_line(self) -> str | None:
        return self.lines[-1] if self.lines else None


def parse_jsonl(data: bytes | str) -> ParsedStderr:
    """Parse stderr written as JSON Lines. Lines that are not event objects are
    kept in `other`, so callers decide whether they are allowed."""
    lines = text_lines(data)
    out = ParsedStderr(lines=lines)
    for i, line in enumerate(lines):
        try:
            obj = json.loads(line)
        except ValueError:
            out.other.append((i, line))
            continue
        if isinstance(obj, dict) and "event" in obj:
            out.events.append(obj)
        else:
            out.other.append((i, line))
    return out


@dataclass
class PlainLine:
    timestamp: str
    severity: str
    event: str
    rest: str


_PLAIN_SPLIT = re.compile(r"^(\S+) (TRACE|DEBUG|INFO|WARN|ERROR) ([a-z][a-z0-9_.]*)(?: (.*))?$")


def parse_plain(line: str) -> PlainLine | None:
    """Return the parts of a line that matches logging.plainLine, else None."""
    if not CONTRACT.plain_line_re.match(line):
        return None
    m = _PLAIN_SPLIT.match(line)
    if not m:
        return None
    return PlainLine(m.group(1), m.group(2).lower(), m.group(3), m.group(4) or "")


_DIAG_REASON = re.compile(r"^aibox: invalid_input: ([a-z_]+): ")


def bootstrap_reason(line: str | None) -> str | None:
    """The rejection reason of an output.bootstrapDiagnostic line, else None."""
    if line is None or not CONTRACT.bootstrap_re.match(line):
        return None
    m = _DIAG_REASON.match(line)
    return m.group(1) if m else None


def severity_rank(severity: str) -> int:
    return CONTRACT.severities.index(severity)


def check_event_catalog(event: dict) -> list[str]:
    """Check an event's severity, component and fixed members against
    logging.events (LOG-001). Returns a list of problems."""
    problems: list[str] = []
    name = event.get("event")
    try:
        spec = CONTRACT.event_spec(name)
    except KeyError:
        return [f"event {name!r} is not in logging.events"]
    if event.get("severity") != spec["severity"]:
        problems.append(f"{name}: severity {event.get('severity')!r}, expected {spec['severity']!r}")
    if event.get("component") != spec["component"]:
        problems.append(f"{name}: component {event.get('component')!r}, expected {spec['component']!r}")
    for member in spec["members"]:
        if " " in member:  # "errorCategory only when the outcome carries an error"
            continue
        if "=" in member:
            key, value = member.split("=", 1)
            if event.get(key) != value:
                problems.append(f"{name}: {key} {event.get(key)!r}, expected {value!r}")
        elif member not in event:
            problems.append(f"{name}: member {member!r} missing")
    executor = CONTRACT.interfaces["logging"]["members"]["executor"]
    if "executor" in spec["members"] and event.get("executor") not in (None, executor):
        problems.append(f"{name}: executor {event.get('executor')!r}, expected {executor!r}")
    return problems

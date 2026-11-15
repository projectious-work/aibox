"""JSON with Comments: reading a devcontainer.json and adding members to it."""

from __future__ import annotations

import json


def strip(text: str) -> str:
    """Remove // and /* */ comments outside strings and trailing commas."""
    out: list[str] = []
    i, n = 0, len(text)
    in_string = False
    while i < n:
        c = text[i]
        if in_string:
            out.append(c)
            if c == "\\" and i + 1 < n:
                out.append(text[i + 1])
                i += 2
                continue
            if c == '"':
                in_string = False
            i += 1
            continue
        if c == '"':
            in_string = True
            out.append(c)
            i += 1
        elif text.startswith("//", i):
            j = text.find("\n", i)
            i = n if j < 0 else j
        elif text.startswith("/*", i):
            j = text.find("*/", i + 2)
            i = n if j < 0 else j + 2
        else:
            out.append(c)
            i += 1
    return _drop_trailing_commas("".join(out))


def _drop_trailing_commas(text: str) -> str:
    out: list[str] = []
    in_string = False
    i = 0
    while i < len(text):
        c = text[i]
        if in_string:
            out.append(c)
            if c == "\\":
                out.append(text[i + 1])
                i += 2
                continue
            if c == '"':
                in_string = False
        elif c == '"':
            in_string = True
            out.append(c)
        elif c == ",":
            j = i + 1
            while j < len(text) and text[j] in " \t\r\n":
                j += 1
            if j < len(text) and text[j] in "}]":
                i += 1
                continue
            out.append(c)
        else:
            out.append(c)
        i += 1
    return "".join(out)


def loads(text: str):
    return json.loads(strip(text))


def add_members(text: str, members: dict) -> str:
    """Add top-level members before the final closing brace of a JSONC object.

    The fixture's last active member has no trailing comma, so a comma is put
    first; a comma after comment lines is valid JSONC. The result is checked
    by parsing it."""
    end = text.rstrip().rfind("}")
    if end < 0:
        raise ValueError("no closing brace")
    body = ",\n".join(f"  {json.dumps(k)}: {json.dumps(v)}" for k, v in members.items())
    result = text[:end] + ",\n" + body + "\n" + text[end:]
    parsed = loads(result)
    for k, v in members.items():
        if parsed.get(k) != v:
            raise ValueError(f"member {k} not added")
    return result


def feature_references(text: str) -> list[str]:
    """effectiveConfig.nativeEntries devcontainer.features value: the selected
    Feature references, sorted ascending by byte order."""
    features = loads(text).get("features") or {}
    return sorted(features, key=lambda s: s.encode())

"""A small JSON Schema (draft 2020-12) validator for the contract schemas.

Only the keywords the three contract schemas use are implemented. Any other
assertion keyword raises UnsupportedKeyword, so a schema change can never be
validated silently by ignoring a keyword. `format` is an annotation, as the
2020-12 default vocabulary defines it.
"""

from __future__ import annotations

import re
from typing import Any

ANNOTATIONS = {"$schema", "$id", "title", "description", "$comment", "format", "default",
               "examples", "$defs"}
ASSERTIONS = {"type", "properties", "additionalProperties", "required", "enum", "const",
              "pattern", "minLength", "maxLength", "items", "maxItems", "minItems",
              "minimum", "maximum", "allOf", "if", "then", "else", "not", "$ref"}


class UnsupportedKeyword(Exception):
    pass


def _is_type(value: Any, name: str) -> bool:
    if name == "object":
        return isinstance(value, dict)
    if name == "array":
        return isinstance(value, list)
    if name == "string":
        return isinstance(value, str)
    if name == "boolean":
        return isinstance(value, bool)
    if name == "null":
        return value is None
    if name == "number":
        return isinstance(value, (int, float)) and not isinstance(value, bool)
    if name == "integer":
        if isinstance(value, bool):
            return False
        return isinstance(value, int) or (isinstance(value, float) and value.is_integer())
    raise UnsupportedKeyword(f"type {name}")


def _equal(a: Any, b: Any) -> bool:
    if isinstance(a, bool) or isinstance(b, bool):
        return type(a) is type(b) and a == b
    return a == b


class Validator:
    def __init__(self, root: dict):
        self.root = root

    def errors(self, instance: Any) -> list[str]:
        out: list[str] = []
        self._validate(instance, self.root, "$", out)
        return out

    def _resolve(self, ref: str) -> dict:
        if not ref.startswith("#/"):
            raise UnsupportedKeyword(f"$ref {ref}")
        node: Any = self.root
        for part in ref[2:].split("/"):
            node = node[part.replace("~1", "/").replace("~0", "~")]
        return node

    def _validate(self, inst: Any, schema: Any, path: str, out: list[str]) -> None:
        if schema is True:
            return
        if schema is False:
            out.append(f"{path}: not allowed")
            return
        unknown = set(schema) - ANNOTATIONS - ASSERTIONS
        if unknown:
            raise UnsupportedKeyword(", ".join(sorted(unknown)))

        if "$ref" in schema:
            self._validate(inst, self._resolve(schema["$ref"]), path, out)
        if "type" in schema:
            types = schema["type"] if isinstance(schema["type"], list) else [schema["type"]]
            if not any(_is_type(inst, t) for t in types):
                out.append(f"{path}: type {type(inst).__name__} is not {types}")
                return
        if "const" in schema and not _equal(inst, schema["const"]):
            out.append(f"{path}: {inst!r} is not const {schema['const']!r}")
        if "enum" in schema and not any(_equal(inst, e) for e in schema["enum"]):
            out.append(f"{path}: {inst!r} not in enum")
        if isinstance(inst, str):
            if "pattern" in schema and not re.search(schema["pattern"], inst):
                out.append(f"{path}: {inst!r} does not match {schema['pattern']}")
            if "minLength" in schema and len(inst) < schema["minLength"]:
                out.append(f"{path}: shorter than {schema['minLength']}")
            if "maxLength" in schema and len(inst) > schema["maxLength"]:
                out.append(f"{path}: longer than {schema['maxLength']}")
        if _is_type(inst, "number"):
            if "minimum" in schema and inst < schema["minimum"]:
                out.append(f"{path}: {inst} below minimum {schema['minimum']}")
            if "maximum" in schema and inst > schema["maximum"]:
                out.append(f"{path}: {inst} above maximum {schema['maximum']}")
        if isinstance(inst, list):
            if "maxItems" in schema and len(inst) > schema["maxItems"]:
                out.append(f"{path}: more than {schema['maxItems']} items")
            if "minItems" in schema and len(inst) < schema["minItems"]:
                out.append(f"{path}: fewer than {schema['minItems']} items")
            if "items" in schema:
                for i, item in enumerate(inst):
                    self._validate(item, schema["items"], f"{path}[{i}]", out)
        if isinstance(inst, dict):
            props = schema.get("properties", {})
            for name in schema.get("required", []):
                if name not in inst:
                    out.append(f"{path}: missing required {name!r}")
            for name, value in inst.items():
                if name in props:
                    self._validate(value, props[name], f"{path}.{name}", out)
                elif "additionalProperties" in schema:
                    self._validate(value, schema["additionalProperties"], f"{path}.{name}", out)
        for sub in schema.get("allOf", []):
            self._validate(inst, sub, path, out)
        if "if" in schema:
            probe: list[str] = []
            self._validate(inst, schema["if"], path, probe)
            branch = schema.get("then") if not probe else schema.get("else")
            if branch is not None:
                self._validate(inst, branch, path, out)
        if "not" in schema:
            probe = []
            self._validate(inst, schema["not"], path, probe)
            if not probe:
                out.append(f"{path}: matches a 'not' schema")


def validate(instance: Any, schema: dict) -> list[str]:
    """Return a list of validation errors; empty when the instance is valid."""
    return Validator(schema).errors(instance)

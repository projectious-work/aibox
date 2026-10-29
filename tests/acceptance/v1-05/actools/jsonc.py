"""Minimal JSONC reader for maintained examples: comments and trailing commas."""

from __future__ import annotations

import json
import re


def strip_comments(text: str) -> str:
    result, index, length = [], 0, len(text)
    in_string = False
    while index < length:
        char = text[index]
        if in_string:
            result.append(char)
            if char == "\\" and index + 1 < length:
                result.append(text[index + 1])
                index += 2
                continue
            if char == '"':
                in_string = False
            index += 1
            continue
        if char == '"':
            in_string = True
            result.append(char)
            index += 1
        elif text.startswith("//", index):
            end = text.find("\n", index)
            index = length if end < 0 else end
        elif text.startswith("/*", index):
            end = text.find("*/", index + 2)
            if end < 0:
                raise ValueError("unterminated block comment")
            index = end + 2
        else:
            result.append(char)
            index += 1
    return "".join(result)


def loads(text: str):
    return json.loads(_drop_trailing_commas(strip_comments(text)))


def _drop_trailing_commas(text: str) -> str:
    out, in_string, index = [], False, 0
    while index < len(text):
        char = text[index]
        if in_string:
            out.append(char)
            if char == "\\":
                out.append(text[index + 1])
                index += 2
                continue
            if char == '"':
                in_string = False
            index += 1
            continue
        if char == '"':
            in_string = True
        elif char == ",":
            match = re.match(r",\s*[}\]]", text[index:])
            if match:
                index += 1
                continue
        out.append(char)
        index += 1
    return "".join(out)

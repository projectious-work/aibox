"""AC-TOOLS/inspect: `aibox inspect` and MCP `inspect_workspace` report
`devcontainer.features` identically and without options (chapter 8 Inspection)."""

from __future__ import annotations

import json
import subprocess
from pathlib import Path
from typing import Mapping

from . import jsonc

FEATURES_KEY = "devcontainer.features"
CONFIG_SOURCE = ".devcontainer/devcontainer.json"
MCP_PROTOCOL_VERSION = "2025-06-18"


def expected_entry(devcontainer_text: str) -> dict | None:
    features = jsonc.loads(devcontainer_text).get("features") or {}
    if not features:
        return None
    return {"key": FEATURES_KEY, "value": sorted(features), "layer": "project",
            "source": CONFIG_SOURCE, "domain": "native", "redacted": False,
            "invocationOnly": False}


def features_entries(result: Mapping) -> list[dict]:
    configuration = (result.get("data") or {}).get("configuration") or []
    return [entry for entry in configuration if entry.get("key") == FEATURES_KEY]


def inspect_failures(expected: dict | None, cli_result: Mapping, mcp_result: Mapping,
                     option_values: list) -> list[str]:
    failures = []
    for name, result in (("CLI", cli_result), ("MCP", mcp_result)):
        if result.get("outcome") != "succeeded":
            failures.append(f"{name} outcome {result.get('outcome')!r}")
        entries = features_entries(result)
        if expected is None:
            if entries:
                failures.append(f"{name} reports {FEATURES_KEY} although no Feature is selected")
            continue
        if len(entries) != 1:
            failures.append(f"{name} reports {len(entries)} {FEATURES_KEY} entries, expected 1")
            continue
        if entries[0] != expected:
            failures.append(f"{name} entry {json.dumps(entries[0], sort_keys=True)} != "
                            f"expected {json.dumps(expected, sort_keys=True)}")
        echoed = json.dumps(entries[0])
        for value in option_values:
            if json.dumps(value) in echoed:
                failures.append(f"{name} entry echoes option value {value!r}")
    if features_entries(cli_result) != features_entries(mcp_result):
        failures.append("CLI and MCP report different devcontainer.features entries")
    return failures


def option_values(devcontainer_text: str) -> list:
    features = jsonc.loads(devcontainer_text).get("features") or {}
    values = []
    for options in features.values():
        if isinstance(options, Mapping):
            values += [value for value in options.values() if not isinstance(value, bool)]
            values += [{key: value} for key, value in options.items()]
    return values


def mcp_arguments(input_schema: Mapping, project_root: str, request_id: str) -> dict:
    """Fill the advertised request-union branch (operation-request.schema.json):
    constants as declared, `requestId` and `projectRoot` from the run."""
    arguments = {}
    properties = input_schema.get("properties", {})
    for name in input_schema.get("required", []):
        spec = properties.get(name, {})
        if "const" in spec:
            arguments[name] = spec["const"]
        elif name == "requestId":
            arguments[name] = request_id
        elif name == "projectRoot":
            arguments[name] = project_root
        else:
            raise ValueError(f"inspect_workspace requires unknown argument {name!r}")
    return arguments


def mcp_tool_result(response: Mapping) -> dict:
    if "error" in response:
        raise RuntimeError(f"MCP protocol error: {response['error']}")
    result = response["result"]
    if isinstance(result.get("structuredContent"), Mapping):
        return dict(result["structuredContent"])
    for item in result.get("content", []):
        if item.get("type") == "text":
            return json.loads(item["text"])
    raise RuntimeError("inspect_workspace returned no structured result")


# --- impure adapters -------------------------------------------------------

def run_cli(aibox: Path, project: Path, log: Path) -> dict:
    argv = [str(aibox), "inspect", "--context", "local", "--project", str(project),
            "--format", "json"]
    completed = subprocess.run(argv, text=True, capture_output=True, check=False, timeout=120)
    log.write_text(f"$ {' '.join(argv)}\nexit {completed.returncode}\n--- stdout\n"
                   f"{completed.stdout}\n--- stderr\n{completed.stderr}\n")
    return json.loads(completed.stdout)


def run_mcp(aibox: Path, project: Path, request_id: str, log: Path) -> dict:
    argv = [str(aibox), "mcp", "serve", "--context", "local"]
    process = subprocess.Popen(argv, stdin=subprocess.PIPE, stdout=subprocess.PIPE,
                               stderr=subprocess.PIPE, text=True)
    transcript = [f"$ {' '.join(argv)}"]

    def call(message_id: int | None, method: str, params: Mapping) -> dict | None:
        message = {"jsonrpc": "2.0", "method": method, "params": dict(params)}
        if message_id is not None:
            message["id"] = message_id
        line = json.dumps(message)
        transcript.append("> " + line)
        assert process.stdin and process.stdout
        process.stdin.write(line + "\n")
        process.stdin.flush()
        if message_id is None:
            return None
        while True:
            reply = process.stdout.readline()
            if not reply:
                raise RuntimeError("MCP server closed its output")
            transcript.append("< " + reply.rstrip("\n"))
            decoded = json.loads(reply)
            if decoded.get("id") == message_id:
                return decoded

    try:
        call(1, "initialize", {"protocolVersion": MCP_PROTOCOL_VERSION, "capabilities": {},
                               "clientInfo": {"name": "aibox-v1-05-acceptance",
                                              "version": "1"}})
        call(None, "notifications/initialized", {})
        tools = call(2, "tools/list", {}) or {}
        schema = next((tool.get("inputSchema", {}) for tool in tools["result"]["tools"]
                       if tool["name"] == "inspect_workspace"), None)
        if schema is None:
            raise RuntimeError("local MCP registry has no inspect_workspace tool")
        arguments = mcp_arguments(schema, str(project), request_id)
        response = call(3, "tools/call", {"name": "inspect_workspace", "arguments": arguments})
        return mcp_tool_result(response or {})
    finally:
        if process.stdin:
            process.stdin.close()
        try:
            process.wait(timeout=30)
        except subprocess.TimeoutExpired:
            process.kill()
        stderr = process.stderr.read() if process.stderr else ""
        log.write_text("\n".join(transcript) + f"\n--- stderr\n{stderr}\n")

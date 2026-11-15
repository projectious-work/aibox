"""A minimal MCP client over stdio (newline-delimited JSON-RPC 2.0)."""

from __future__ import annotations

import json
import queue
import subprocess
import threading
import time

from .proc import Result

PROTOCOL_VERSION = "2025-06-18"


class MCPError(Exception):
    pass


class MCPSession:
    def __init__(self, argv: list[str], env: dict[str, str], cwd: str):
        self.argv, self.env, self.cwd = argv, env, cwd
        self.proc = subprocess.Popen(argv, env=env, cwd=cwd, stdin=subprocess.PIPE,
                                     stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        self.stdout_lines: list[bytes] = []
        self._stderr: list[bytes] = []
        self._messages: queue.Queue = queue.Queue()
        self._next_id = 1
        self._lock = threading.Lock()
        self._responses: dict = {}
        self._send_lock = threading.Lock()
        self._out_thread = threading.Thread(target=self._read_stdout, daemon=True)
        self._out_thread.start()
        self._err_thread = threading.Thread(target=self._read_stderr, daemon=True)
        self._err_thread.start()
        self.exit: int | None = None

    def _read_stdout(self) -> None:
        for line in self.proc.stdout:
            self.stdout_lines.append(line)
            try:
                self._messages.put(json.loads(line))
            except ValueError:
                self._messages.put({"_unparsable": line.decode(errors="replace")})
        self._messages.put(None)

    def _read_stderr(self) -> None:
        self._stderr.append(self.proc.stderr.read())

    def _send(self, message: dict) -> None:
        try:
            with self._send_lock:
                self.proc.stdin.write((json.dumps(message) + "\n").encode())
                self.proc.stdin.flush()
        except (BrokenPipeError, ValueError) as exc:
            raise MCPError(f"server stdin closed: {exc}") from exc

    def request(self, method: str, params: dict | None = None, timeout: float = 30) -> dict:
        """Send a request and return the response with the same id."""
        with self._lock:
            rid = self._next_id
            self._next_id += 1
        msg = {"jsonrpc": "2.0", "id": rid, "method": method}
        if params is not None:
            msg["params"] = params
        self._send(msg)
        deadline = time.monotonic() + timeout
        while True:
            with self._lock:
                if rid in self._responses:
                    return self._responses.pop(rid)
            remaining = deadline - time.monotonic()
            if remaining <= 0:
                raise MCPError(f"no response to {method} within {timeout}s")
            try:
                item = self._messages.get(timeout=min(remaining, 0.1))
            except queue.Empty:
                continue
            if item is None:
                self._messages.put(None)
                with self._lock:
                    if rid in self._responses:
                        return self._responses.pop(rid)
                raise MCPError(f"server closed stdout before answering {method}")
            if "id" in item and ("result" in item or "error" in item):
                # Keep every response by id, so concurrent callers each get theirs.
                with self._lock:
                    self._responses[item["id"]] = item

    def notify(self, method: str, params: dict | None = None) -> None:
        msg = {"jsonrpc": "2.0", "method": method}
        if params is not None:
            msg["params"] = params
        self._send(msg)

    def initialize(self) -> dict:
        resp = self.request("initialize", {"protocolVersion": PROTOCOL_VERSION, "capabilities": {},
                                           "clientInfo": {"name": "aibox-v106-acceptance",
                                                          "version": "1"}})
        self.notify("notifications/initialized")
        return resp

    def call_tool(self, name: str, arguments: dict, timeout: float = 30) -> dict:
        return self.request("tools/call", {"name": name, "arguments": arguments}, timeout)

    def close(self, timeout: float = 30) -> Result:
        """Close stdin, wait for the server to exit and return its transcript."""
        try:
            self.proc.stdin.close()
        except OSError:
            pass
        timed_out = False
        try:
            self.exit = self.proc.wait(timeout=timeout)
        except subprocess.TimeoutExpired:
            self.proc.kill()
            self.proc.wait()
            timed_out = True
        self._err_thread.join(10)
        self._out_thread.join(10)
        for pipe in (self.proc.stdout, self.proc.stderr):
            try:
                pipe.close()
            except OSError:
                pass
        return Result(self.argv, self.cwd, self.env, None if timed_out else self.exit,
                      b"".join(self.stdout_lines), b"".join(self._stderr), timed_out,
                      note="mcp serve session")

    def kill(self) -> None:
        if self.proc.poll() is None:
            self.proc.kill()
            self.proc.wait()


def structured(response: dict) -> dict | None:
    """The operation result of a tools/call response, or None."""
    result = response.get("result")
    if not isinstance(result, dict):
        return None
    if isinstance(result.get("structuredContent"), dict):
        return result["structuredContent"]
    for item in result.get("content") or []:
        if item.get("type") == "text":
            try:
                doc = json.loads(item["text"])
            except ValueError:
                continue
            if isinstance(doc, dict) and "schemaVersion" in doc:
                return doc
    return None


def is_rejected(response: dict) -> bool:
    """A tool call refused before it ran: a JSON-RPC error or isError true."""
    if "error" in response:
        return True
    result = response.get("result")
    return isinstance(result, dict) and result.get("isError") is True

"""Running the executable under test and capturing what it observably does."""

from __future__ import annotations

import os
import pty
import select
import shlex
import subprocess
import threading
import time
from dataclasses import dataclass, field

TIMEOUT = 60


class InfrastructureError(Exception):
    """The harness cannot run (exit 2): a missing tool, image or binary."""


@dataclass
class Result:
    argv: list[str]
    cwd: str
    env: dict[str, str]
    exit: int | None
    stdout: bytes
    stderr: bytes
    timed_out: bool = False
    note: str = ""
    extra: dict = field(default_factory=dict)

    @property
    def out(self) -> str:
        return self.stdout.decode("utf-8", errors="replace")

    @property
    def err(self) -> str:
        return self.stderr.decode("utf-8", errors="replace")

    def describe(self, limit: int = 6000) -> dict:
        def cut(s: str) -> str:
            return s if len(s) <= limit else s[:limit] + f"... [{len(s) - limit} more bytes]"
        return {
            "command": shlex.join(self.argv),
            "cwd": self.cwd,
            "env": self.env,
            "exit": self.exit,
            "timedOut": self.timed_out,
            "stdout": cut(self.out),
            "stderr": cut(self.err),
            **({"note": self.note} if self.note else {}),
        }


def run(argv: list[str], *, env: dict[str, str], cwd: str, stdin: bytes | None = None,
        timeout: float = TIMEOUT, stderr_closed: bool = False,
        stderr_fd: int | None = None) -> Result:
    """Run a command with exactly `env`. stderr_closed runs it with fd 2 closed
    (`2>&-`); stderr_fd gives it a caller-owned descriptor as stderr."""
    if stderr_closed:
        stderr_target = None
        preexec = _close_stderr
    elif stderr_fd is not None:
        stderr_target, preexec = stderr_fd, None
    else:
        stderr_target, preexec = subprocess.PIPE, None
    proc = subprocess.Popen(argv, env=env, cwd=cwd, stdout=subprocess.PIPE,
                            stderr=stderr_target,
                            stdin=subprocess.PIPE if stdin is not None else subprocess.DEVNULL,
                            preexec_fn=preexec, close_fds=True)
    try:
        out, err = proc.communicate(stdin, timeout=timeout)
        timed_out = False
    except subprocess.TimeoutExpired:
        proc.kill()
        out, err = proc.communicate()
        timed_out = True
    return Result(list(argv), cwd, dict(env), None if timed_out else proc.returncode,
                  out or b"", err or b"", timed_out)


def _close_stderr() -> None:
    os.close(2)


def run_pty(argv: list[str], *, env: dict[str, str], cwd: str,
            timeout: float = TIMEOUT) -> Result:
    """Run with stdout a pseudo-terminal (J-OUT-03); stderr stays a pipe."""
    master, slave = pty.openpty()
    proc = subprocess.Popen(argv, env=env, cwd=cwd, stdout=slave, stderr=subprocess.PIPE,
                            stdin=subprocess.DEVNULL, close_fds=True)
    os.close(slave)
    err_chunks: list[bytes] = []
    reader = threading.Thread(target=lambda: err_chunks.append(proc.stderr.read()), daemon=True)
    reader.start()
    chunks: list[bytes] = []
    deadline = time.monotonic() + timeout
    timed_out = False
    try:
        while True:
            remaining = deadline - time.monotonic()
            if remaining <= 0:
                timed_out = True
                proc.kill()
                break
            ready, _, _ = select.select([master], [], [], min(remaining, 0.5))
            if ready:
                try:
                    data = os.read(master, 65536)
                except OSError:  # EIO: the child closed the terminal
                    break
                if not data:
                    break
                chunks.append(data)
            elif proc.poll() is not None:
                # Drain anything left after exit.
                try:
                    while select.select([master], [], [], 0.1)[0]:
                        data = os.read(master, 65536)
                        if not data:
                            break
                        chunks.append(data)
                except OSError:
                    pass
                break
    finally:
        os.close(master)
    try:
        proc.wait(timeout=10)
    except subprocess.TimeoutExpired:
        proc.kill()
        proc.wait()
        timed_out = True
    reader.join(10)
    return Result(list(argv), cwd, dict(env), None if timed_out else proc.returncode,
                  b"".join(chunks), b"".join(err_chunks), timed_out, note="stdout is a pty")

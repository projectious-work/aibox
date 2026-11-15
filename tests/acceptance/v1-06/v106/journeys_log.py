"""Logging journeys J-LOG-01 to J-LOG-12 of journeys.md."""

from __future__ import annotations

import json
import os
import re
import threading
import time
from concurrent.futures import ThreadPoolExecutor
from pathlib import Path

from .contract import CONTRACT
from .events import parse_jsonl, text_lines
from .expect import check_events_valid, event_errors, expect_rejected, expect_success
from .framework import T, journey
from .mcp import MCPError, MCPSession, structured
from .sandbox import Sandbox, file_state

SUCCESS_ORDER = CONTRACT.interfaces["logging"]["eventOrder"]["inspect, success"]
FAILURE_ORDER = CONTRACT.interfaces["logging"]["eventOrder"]["inspect, inspection fails"]
MIB = 1 << 20


def names_at(level: str, order: list[str]) -> list[str]:
    """Events of `order` whose severity is at or above `level` (LOG-005)."""
    if level == "off":
        return []
    rank = CONTRACT.severities.index(level)
    return [n for n in order if CONTRACT.severities.index(CONTRACT.severity(n)) >= rank]


def mode_of(path: str) -> int | None:
    try:
        return os.stat(path).st_mode & 0o777
    except OSError:
        return None


# --- J-LOG-01 to J-LOG-03 --------------------------------------------------------------

@journey("J-LOG-01", ["LOG-001", "LOG-002", "LOG-003", "LOG-010"])
def j_log_01(t: T) -> None:
    sb = t.local()
    p = sb.project("P")
    res = sb.aibox_run(["inspect", "--project", p, "--format", "json", "--log-level", "trace",
                        "--log-format", "jsonl"])
    doc = expect_success(t, res)
    parsed = parse_jsonl(res.stderr)
    t.eq(parsed.other, [], "every stderr line is a JSON event", result=res)
    t.eq(parsed.names(), SUCCESS_ORDER, "stderr holds exactly the four events in order", result=res)
    check_events_valid(t, parsed.events, res)
    rid = doc.get("requestId") if doc else None
    t.eq([e.get("requestId") for e in parsed.events], [rid] * len(parsed.events),
         "one requestId, equal to the result's", result=res)
    t.eq([e.get("sequence") for e in parsed.events], list(range(len(parsed.events))),
         "sequence 0 to 3", result=res)
    for ev in parsed.events:
        name = ev.get("event")
        if name and name.startswith("operation."):
            t.eq((ev.get("operation"), ev.get("initiator"), ev.get("executor")),
                 ("inspect_workspace", "cli", "aibox"), f"{name} operation, initiator, executor",
                 result=res)
    started = parsed.named("operation.started")
    if started:
        t.check("projectRoot" in (started[0].get("fields") or {}),
                "operation.started has fields.projectRoot", result=res)
    finished = parsed.named("operation.finished")
    if finished:
        t.eq(finished[0].get("outcome"), "succeeded", "operation.finished outcome", result=res)
        t.check("errorCategory" not in finished[0], "operation.finished has no errorCategory",
                result=res)
    resolved = parsed.named("settings.resolved")
    if resolved:
        t.eq((resolved[0].get("fields") or {}).get("layers"), "invocation",
             "settings.resolved fields.layers", result=res)


@journey("J-LOG-02", ["LOG-001", "LOG-005"])
def j_log_02(t: T) -> None:
    sb = t.local()
    p = sb.project("P")
    missing = f"{sb.projects}/missing"
    for level in CONTRACT.key("logging.level")["values"]:
        for project, order, exit_code in ((p, SUCCESS_ORDER, 0), (missing, FAILURE_ORDER, 2)):
            with t.case(f"--log-level {level} on {'P' if project == p else 'a missing path'}"):
                res = sb.aibox_run(["inspect", "--project", project, "--format", "json",
                                    "--log-format", "jsonl", "--log-level", level])
                t.eq(res.exit, exit_code, "exit code", result=res)
                parsed = parse_jsonl(res.stderr)
                t.eq(parsed.other, [], "every stderr line is a JSON event", result=res)
                t.eq(parsed.names(), names_at(level, order), "events at this level", result=res)
                if project == missing:
                    for ev in parsed.named("operation.finished"):
                        t.eq((ev.get("outcome"), ev.get("errorCategory")), ("failed", "invalid_input"),
                             "operation.finished outcome and errorCategory", result=res)


@journey("J-LOG-03", ["LOG-010"])
def j_log_03(t: T) -> None:
    sb = t.local()
    p = sb.project("P")
    for label, extra in (("--log-format plain", ["--log-format", "plain"]),
                         ("no --log-format", [])):
        with t.case(label):
            res = sb.aibox_run(["inspect", "--project", p, "--format", "json", "--log-level", "trace",
                                *extra])
            lines = text_lines(res.stderr)
            bad = [ln for ln in lines if not CONTRACT.plain_line_re.match(ln)]
            t.eq(bad, [], "every stderr line matches logging.plainLine", result=res)
            t.eq(len(lines), 4, "four lines", result=res)


# --- J-LOG-04 MCP ------------------------------------------------------------------------

@journey("J-LOG-04", ["LOG-004", "LOG-011", "LOG-012"])
def j_log_04(t: T) -> None:
    sb = t.local()
    p = sb.project("P")
    session = MCPSession([sb.aibox, "mcp", "serve", "--context", "local", "--log-level", "debug",
                          "--log-file", "m.log"], sb.env(), sb.cwd)
    try:
        session.initialize()
        for i in (1, 2):
            resp = session.call_tool("inspect_workspace", {"requestId": f"j-log-04-call-{i}",
                                                           "projectRoot": p})
            doc = structured(resp)
            t.check(doc is not None and doc.get("outcome") == "succeeded",
                    f"call {i} succeeded", actual=resp)
    except MCPError as exc:
        t.check(False, f"MCP session: {exc}")
    finally:
        res = session.close()
    bad_stdout = []
    for line in text_lines(res.stdout):
        try:
            msg = json.loads(line)
            ok = isinstance(msg, dict) and msg.get("jsonrpc") == "2.0"
        except ValueError:
            ok = False
        if not ok:
            bad_stdout.append(line)
    t.eq(bad_stdout, [], "every stdout line is an MCP JSON-RPC message", result=res)
    parsed = parse_jsonl(res.stderr)
    t.eq(parsed.other, [], "stderr is JSON Lines", result=res)
    check_events_valid(t, parsed.events, res)
    ev = parsed.events
    if not t.check(len(ev) >= 2, "stderr holds events", result=res):
        return
    t.eq(ev[0].get("event"), "process.started", "process.started first", result=res)
    t.eq(ev[-1].get("event"), "process.finished", "process.finished last", result=res)
    pid = ev[0].get("requestId")
    process_events = [e for e in ev if e.get("requestId") == pid]
    t.eq([(e.get("event"), e.get("sequence")) for e in process_events],
         [("process.started", 0), ("settings.resolved", 1), ("process.finished", 2)],
         "process ID events with sequence 0, 1, 2", result=res)
    calls: dict[str, list[dict]] = {}
    for e in ev:
        if e.get("requestId") != pid:
            calls.setdefault(e.get("requestId"), []).append(e)
    t.eq(len(calls), 2, "the two calls use two other request IDs", result=res)
    for rid, events in calls.items():
        t.eq([e.get("sequence") for e in events], [0, 1, 2], f"request {rid}: sequence 0 to 2",
             result=res)
        t.eq({e.get("initiator") for e in events}, {"mcp"}, f"request {rid}: initiator mcp",
             result=res)
    log = Path(sb.logs) / "m.log"
    if t.check(log.is_file(), "~/.local/state/aibox/logs/m.log exists", result=res):
        file_events = parse_jsonl(log.read_bytes())
        t.eq(file_events.events, parsed.events, "m.log holds the same events", result=res)


# --- J-LOG-05, J-LOG-06 file sink place and modes ---------------------------------------

@journey("J-LOG-05", ["LOG-012", "LOG-013", "LOG-015"])
def j_log_05(t: T) -> None:
    sb = t.local()
    p = sb.project("P")
    res = sb.aibox_run(["inspect", "--project", p, "--format", "json", "--log-format", "jsonl",
                        "--log-file", "run.log"])
    expect_success(t, res)
    log = f"{sb.logs}/run.log"
    if t.check(os.path.isfile(log), "~/.local/state/aibox/logs/run.log exists", result=res):
        t.eq(oct(mode_of(log) or 0), oct(0o600), "run.log mode", result=res)
        t.eq(text_lines(Path(log).read_bytes()), text_lines(res.stderr),
             "run.log lines equal the stderr lines", result=res)
    t.eq(oct(mode_of(sb.state) or 0), oct(0o700), "~/.local/state/aibox mode", result=res)
    t.eq(oct(mode_of(sb.logs) or 0), oct(0o700), "~/.local/state/aibox/logs mode", result=res)
    t.check(not os.path.exists(f"{sb.cwd}/run.log"), "the working directory has no run.log",
            result=res)


@journey("J-LOG-06", ["LOG-014"])
def j_log_06(t: T) -> None:
    with t.case("--log-file /tmp/x.log"):
        sb = t.local("abs")
        p = sb.project("P")
        before = file_state("/tmp/x.log")
        res = sb.aibox_run(sb.rejection_args(p, "--log-file", "/tmp/x.log"))
        expect_rejected(t, res, "path_escape")
        t.eq(file_state("/tmp/x.log"), before, "/tmp/x.log is not created or changed", result=res)
    with t.case("--log-file ../x.log"):
        sb = t.local("dotdot")
        p = sb.project("P")
        res = sb.aibox_run(sb.rejection_args(p, "--log-file", "../x.log"))
        expect_rejected(t, res, "path_escape")
        t.eq(file_state(f"{sb.state}/x.log"), None, "~/.local/state/aibox/x.log is not created",
             result=res)
    with t.case("--log-file link.log, a symlink out of the root"):
        sb = t.local("link")
        p = sb.project("P")
        outside = sb.write(f"{sb.root}/outside/target.log", "keep\n")
        sb.mkdir(sb.state)
        sb.mkdir(sb.logs)
        os.symlink(outside, f"{sb.logs}/link.log")
        before = file_state(outside)
        res = sb.aibox_run(sb.rejection_args(p, "--log-file", "link.log"))
        expect_rejected(t, res, "path_escape")
        t.eq(file_state(outside), before, "the outside file is not changed", result=res)
    with t.case("user settings logging.file /var/log/aibox.log"):
        sb = t.local("var")
        p = sb.project("P")
        sb.user_settings({"schemaVersion": "1", "logging": {"file": "/var/log/aibox.log"}})
        before = file_state("/var/log/aibox.log")
        res = sb.aibox_run(sb.rejection_args(p))
        expect_rejected(t, res, "path_escape")
        t.eq(file_state("/var/log/aibox.log"), before, "/var/log/aibox.log is not created or changed",
             result=res)


# --- J-LOG-07, J-LOG-09 rotation and retention -------------------------------------------

def fill(sb: Sandbox, path: str, first: str) -> None:
    """Exactly 1 MiB of lines whose first line is `first`."""
    head = f"{first}\n".encode()
    filler = b"f" * 99 + b"\n"
    body = bytearray(head)
    while len(body) + len(filler) <= MIB:
        body += filler
    rest = MIB - len(body)
    if rest:
        body += b"f" * (rest - 1) + b"\n"
    assert len(body) == MIB
    sb.write(path, bytes(body))


def rotation_set(names: list[str], base: str) -> list[str]:
    pat = re.compile(rf"^{re.escape(base)}(\.lock|\.[1-9][0-9]?)?$")
    return sorted(n for n in names if pat.match(n))


def starts_with(path: str, first: str) -> bool:
    try:
        with open(path, "rb") as fh:
            return fh.readline() == f"{first}\n".encode()
    except OSError:
        return False


def parse_json_or_none(data: bytes):
    try:
        return json.loads(data)
    except ValueError:
        return None


def rotation_sequence(t: T, sb: Sandbox, retention_files: int, *, strict: bool = True) -> None:
    """The J-LOG-07 procedure. strict=False checks only the rotation set, for
    J-LOG-09, whose logs directory holds another file."""
    settings = {"schemaVersion": "1", "logging": {"file": "r.log", "rotationMiB": 1,
                                                  "retentionFiles": retention_files}}
    sb.user_settings(settings)
    p = sb.project("P")
    rlog = f"{sb.logs}/r.log"
    cmd = ["inspect", "--project", p, "--format", "json"]

    def listing() -> list[str]:
        names = os.listdir(sb.logs) if os.path.isdir(sb.logs) else []
        return sorted(names) if strict else rotation_set(names, "r.log")

    res = None
    for letter in "ABCD":
        with t.case(f"retentionFiles {retention_files}, run with {letter}"):
            fill(sb, rlog, letter)
            res = sb.aibox_run(cmd)
            expect_success(t, res)
            if retention_files == 1:
                t.eq(listing(), ["r.log", "r.log.lock"], "only r.log and r.log.lock exist", result=res)
    with t.case(f"retentionFiles {retention_files}, after four runs"):
        if retention_files == 3:
            t.eq(listing(), ["r.log", "r.log.1", "r.log.2", "r.log.lock"],
                 "exactly r.log, r.log.1, r.log.2 and r.log.lock exist", result=res)
            t.check(starts_with(f"{rlog}.1", "D"), "r.log.1 starts with D", result=res)
            t.check(starts_with(f"{rlog}.2", "C"), "r.log.2 starts with C", result=res)
        if res is not None:
            last = parse_json_or_none(res.stdout)
            rid = last.get("requestId") if isinstance(last, dict) else None
            parsed = parse_jsonl(Path(rlog).read_bytes() if os.path.isfile(rlog) else b"")
            t.check(bool(parsed.events) and not parsed.other
                    and all(e.get("requestId") == rid for e in parsed.events),
                    "r.log holds only the events of the last run", result=res,
                    actual={"lines": parsed.lines[:3], "lineCount": len(parsed.lines),
                            "lastRequestId": rid})
    with t.case(f"retentionFiles {retention_files}, fifth run with retentionDays 7"):
        old = time.time() - 30 * 86400
        if os.path.exists(f"{rlog}.2"):
            os.utime(f"{rlog}.2", (old, old))
        else:
            t.note("r.log.2 does not exist, so no modification time was set")
        settings["logging"]["retentionDays"] = 7
        sb.user_settings(settings)
        fill(sb, rlog, "E")
        res = sb.aibox_run(cmd)
        expect_success(t, res)
        if retention_files == 3:
            t.check(starts_with(f"{rlog}.2", "D"), "r.log.2 holds the D content", result=res)
            limit = time.time() - 7 * 86400
            old_files = [n for n in listing() if os.path.getmtime(f"{sb.logs}/{n}") < limit]
            t.eq(old_files, [], "no file older than 7 days exists", result=res)
        else:
            t.eq(listing(), ["r.log", "r.log.lock"], "only r.log and r.log.lock exist", result=res)


@journey("J-LOG-07", ["LOG-016", "LOG-017"])
def j_log_07(t: T) -> None:
    rotation_sequence(t, t.local("three"), 3)
    rotation_sequence(t, t.local("one"), 1)


@journey("J-LOG-09", ["LOG-019"])
def j_log_09(t: T) -> None:
    for retention_files in (3, 1):
        with t.case(f"J-LOG-07 sequence with retentionFiles {retention_files}"):
            sb = t.local(f"rf{retention_files}")
            keep = sb.write(f"{sb.state}/operations/keep.json", '{"keep": true}\n')
            other = sb.write(f"{sb.logs}/other.txt", "other\n")
            before = {keep: file_state(keep), other: file_state(other)}
            # J-LOG-07's own expectations are not J-LOG-09's; its runs must still
            # happen, so each inspect of the sequence has to succeed.
            sub = T(t.ctx, t.jid)
            rotation_sequence(sub, sb, retention_files, strict=False)
            run_failures = [f for f in sub.failures if f["check"].startswith("result")]
            t.eq(run_failures, [], "every inspect of the J-LOG-07 sequence succeeded")
            for path, state in before.items():
                t.eq(file_state(path), state, f"{os.path.basename(path)} is unchanged")


# --- J-LOG-08 concurrency ---------------------------------------------------------------

def run_concurrently(sb: Sandbox, args: list[str], n: int):
    barrier = threading.Barrier(n)

    def one(_):
        barrier.wait()
        return sb.aibox_run(args)

    with ThreadPoolExecutor(n) as pool:
        return list(pool.map(one, range(n)))


def request_ids(t: T, results) -> list[str]:
    ids = []
    for res in results:
        doc = expect_success(t, res)
        if doc is not None:
            ids.append(doc.get("requestId"))
    return ids


def check_sequences(t: T, events: list[dict], ids: list[str], what: str) -> None:
    by_id: dict[str, list[int]] = {}
    for e in events:
        by_id.setdefault(e.get("requestId"), []).append(e.get("sequence"))
    t.eq(sorted(by_id), sorted(ids), f"{what}: the request IDs are those of the 16 results")
    for rid in ids:
        t.eq(sorted(by_id.get(rid, [])), [0, 1, 2, 3], f"{what}: request {rid} has sequence 0 to 3")


@journey("J-LOG-08", ["LOG-018"])
def j_log_08(t: T) -> None:
    n = 16
    with t.case("(a) 16 processes"):
        sb = t.local("a")
        p = sb.project("P")
        results = run_concurrently(sb, ["inspect", "--project", p, "--format", "json",
                                        "--log-level", "trace", "--log-file", "c.log"], n)
        ids = request_ids(t, results)
        log = Path(sb.logs) / "c.log"
        data = log.read_bytes() if log.is_file() else b""
        t.check(data.endswith(b"\n"), "c.log ends with a complete line")
        parsed = parse_jsonl(data)
        t.eq(len(parsed.lines), 64, "c.log has 64 lines")
        t.eq(parsed.other, [], "every line is a complete JSON event")
        invalid = [e for e in parsed.events if event_errors(e)]
        t.eq(invalid, [], "every line validates against log-event.schema.json")
        check_sequences(t, parsed.events, ids, "c.log")
    with t.case("(b) 16 processes across rotation"):
        sb = t.local("b")
        p = sb.project("P")
        sb.user_settings({"schemaVersion": "1", "logging": {"rotationMiB": 1, "retentionFiles": 100}})
        prefill = []
        size = 0
        while size + 15 <= MIB - 512:
            line = f"prefill-{len(prefill):06d}"
            prefill.append(line)
            size += len(line) + 1
        sb.write(f"{sb.logs}/c.log", "".join(f"{x}\n" for x in prefill))
        results = run_concurrently(sb, ["inspect", "--project", p, "--format", "json",
                                        "--log-level", "trace", "--log-file", "c.log"], n)
        ids = request_ids(t, results)
        files = [f for f in rotation_set(os.listdir(sb.logs), "c.log") if not f.endswith(".lock")]
        lines: list[str] = []
        for f in files:
            data = (Path(sb.logs) / f).read_bytes()
            t.check(data == b"" or data.endswith(b"\n"), f"{f} ends with a complete line")
            lines += text_lines(data)
        pre = [ln for ln in lines if ln.startswith("prefill-")]
        t.check(sorted(pre) == prefill, "every prefilled line exactly once",
                actual={"prefilled": len(prefill), "found": len(pre), "unique": len(set(pre))})
        parsed = parse_jsonl("\n".join(ln for ln in lines if not ln.startswith("prefill-")))
        t.eq(parsed.other, [], "every other line is a complete JSON event")
        t.eq(len(parsed.events), 64, "the 64 events")
        keys = [(e.get("requestId"), e.get("sequence")) for e in parsed.events]
        t.eq(len(set(keys)), len(keys), "each event exactly once")
        check_sequences(t, parsed.events, ids, "c.log and c.log.<n>")
        t.note(f"files: {files}")


# --- J-LOG-10, J-LOG-11 sink failure ----------------------------------------------------

def expect_degraded(t: T, res, reason: str) -> None:
    doc = expect_success(t, res)
    if doc is not None:
        t.check("log_sink_degraded" in (doc.get("warnings") or []),
                "warnings contains log_sink_degraded", result=res, actual=doc.get("warnings"))
    events = parse_jsonl(res.stderr).named("logging.sink_degraded")
    t.eq(len(events), 1, "stderr has one logging.sink_degraded event", result=res)
    if events:
        check_events_valid(t, events, res)
        f = events[0].get("fields") or {}
        t.eq((f.get("sink"), f.get("reason")), ("file", reason),
             "logging.sink_degraded fields.sink and fields.reason", result=res)


@journey("J-LOG-10", ["LOG-020", "LOG-030", "LOG-031"])
def j_log_10(t: T) -> None:
    for label, required in (("(a) logs/ mode 0500", False),
                            ("(b) logs/ mode 0500, logging.required true", True)):
        with t.case(label):
            sb = t.local(label)
            p = sb.project("P")
            if required:
                sb.user_settings({"schemaVersion": "1", "logging": {"required": True}})
            sb.mkdir(sb.state)
            sb.mkdir(sb.logs, mode=0o500)
            res = sb.aibox_run(["inspect", "--project", p, "--format", "json", "--log-format", "jsonl",
                                "--log-file", "new.log"])
            expect_degraded(t, res, "open_failed")
    with t.case("(c) logs/wide.log mode 0644"):
        sb = t.local("c")
        p = sb.project("P")
        wide = sb.write(f"{sb.logs}/wide.log", "existing line\n", mode=0o644)
        before = file_state(wide)
        res = sb.aibox_run(["inspect", "--project", p, "--format", "json", "--log-format", "jsonl",
                            "--log-file", "wide.log"])
        expect_degraded(t, res, "insecure_mode")
        t.eq(file_state(wide), before, "wide.log is unchanged", result=res)


@journey("J-LOG-11", ["LOG-032"])
def j_log_11(t: T) -> None:
    for label in ("(a) stderr closed", "(b) stderr a pipe whose reader has exited"):
        with t.case(label):
            sb = t.local(label[:3])
            p = sb.project("P")
            args = ["inspect", "--project", p, "--format", "json", "--log-level", "trace",
                    "--log-file", "e.log"]
            if label.startswith("(a)"):
                res = sb.aibox_run(args, stderr_closed=True)
            else:
                r, w = os.pipe()
                os.close(r)
                try:
                    res = sb.aibox_run(args, stderr_fd=w)
                finally:
                    os.close(w)
            doc = expect_success(t, res)
            if doc is not None:
                t.eq(doc.get("warnings"), [], "warnings is empty", result=res)
            log = Path(sb.logs) / "e.log"
            data = log.read_bytes() if log.is_file() else b""
            parsed = parse_jsonl(data)
            t.eq(parsed.other, [], "e.log is JSON Lines", result=res)
            t.eq(parsed.names(), SUCCESS_ORDER, "e.log holds the four events and nothing else",
                 result=res)
            stray = []
            for root in (sb.home, sb.cwd):
                for dirpath, _, filenames in os.walk(root):
                    for f in filenames:
                        path = os.path.join(dirpath, f)
                        if path not in (str(log), f"{log}.lock"):
                            stray.append(path)
            t.eq(stray, [], "no event text in any other file", result=res)


# --- J-LOG-12 redaction ---------------------------------------------------------------

CANARY = "canary-0123456789"
GHP = "ghp_" + "abcdefghijklmnopqrstuvwxyzABCDEFGHIJ"


@journey("J-LOG-12", ["RED-001", "RED-002", "RED-005"])
def j_log_12(t: T) -> None:
    assert len(GHP) == 40
    sb = t.local()
    env = {"EXAMPLE_API_TOKEN": CANARY}
    projects = [("canary directory", sb.project("P", name=f"p-{CANARY}")),
                ("ghp_ directory", sb.project("P", name=f"p-{GHP}")),
                ("missing path with the canary", f"{sb.projects}/missing-{CANARY}")]
    for label, project in projects:
        with t.case(label):
            res = sb.aibox_run(["inspect", "--project", project, "--log-level", "trace",
                                "--log-format", "jsonl", "--log-file", "s.log"], env)
            started = parse_jsonl(res.stderr).named("operation.started")
            t.eq(len(started), 1, "one operation.started event", result=res)
            for ev in started:
                root = (ev.get("fields") or {}).get("projectRoot") or ""
                t.check(CONTRACT.replacement in root, "fields.projectRoot contains [redacted]",
                        result=res, actual=root)
            for secret in (CANARY, GHP):
                t.check(secret not in res.err, f"stderr does not contain {secret[:8]}...", result=res)
                t.check(secret not in res.out, f"human stdout does not contain {secret[:8]}...",
                        result=res)
    with t.case("s.log"):
        log = Path(sb.logs) / "s.log"
        data = log.read_text(errors="replace") if log.is_file() else ""
        t.check(log.is_file(), "s.log exists")
        for secret in (CANARY, GHP):
            t.check(secret not in data, f"s.log does not contain {secret[:8]}...")
    p = sb.project("P")
    with t.case("inspect --no-redact"):
        expect_rejected(t, sb.aibox_run(sb.rejection_args(p, "--no-redact")), "usage_error")
    with t.case("AIBOX_LOGGING__REDACT=false"):
        expect_rejected(t, sb.aibox_run(sb.rejection_args(p), {"AIBOX_LOGGING__REDACT": "false"}),
                        "unknown_environment_name")

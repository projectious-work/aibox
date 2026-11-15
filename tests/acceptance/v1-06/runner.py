#!/usr/bin/env python3
"""V1-06 black-box acceptance suite for settings-and-logging@v1.

Runs every journey of spec/settings-and-logging/v1/journeys.md against a built
aibox executable, through the CLI, MCP over stdio and files only.

Usage:
  tests/acceptance/v1-06/run.sh --aibox PATH --evidence DIR [--journey ID ...]
                                [--work-dir DIR] [--keep] [--list]

Build the executable under test first (the development container has no Go):
  podman run --rm -v "$PWD":/w:Z -w /w -e CGO_ENABLED=0 \\
      localhost/aibox-docs-build:go1.27.0 go build -o /w/<out> ./cmd/aibox

Environment: run as an ordinary user (not root) with Python 3, `pty` and
rootless Podman. Root journeys (J-CFG-02 pairs with system, J-CFG-03,
J-CFG-06 (a), J-CFG-08) run in a disposable container from the pinned Debian
image, which must be present locally; the run needs no network.

Exit status: 0 every journey passed; 1 a journey failed or was not run (for
example filtered out with --journey); 2 usage or infrastructure error.

DIR receives journeys/<ID>.json (ID, covered requirements, outcome, observed
output of every failed check), summary.json and summary.txt.
"""

from __future__ import annotations

import argparse
import hashlib
import os
import shutil
import subprocess
import sys
import tempfile
import time
from datetime import datetime, timezone
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))

from v106 import journeys_cfg, journeys_log, journeys_out  # noqa: E402,F401  (registration)
from v106.contract import CONTRACT, REPO  # noqa: E402
from v106.framework import REGISTRY, RunContext, dump, run_journey  # noqa: E402
from v106.proc import InfrastructureError  # noqa: E402
from v106.sandbox import DEBIAN_IMAGE, RootContainer  # noqa: E402
from v106.wheel import WheelError, verify  # noqa: E402

SPEC_ID = "settings-and-logging@v1"


def now() -> str:
    return datetime.now(timezone.utc).isoformat(timespec="seconds").replace("+00:00", "Z")


def parse_args(argv: list[str]) -> argparse.Namespace:
    ap = argparse.ArgumentParser(prog="run.sh", description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--aibox", help="the aibox executable under test")
    ap.add_argument("--evidence", help="directory for the journey records and the summary")
    ap.add_argument("--journey", action="append", default=[],
                    help="run only these journey IDs (others count as not run)")
    ap.add_argument("--work-dir", help="parent of the temporary sandboxes (default: system temp)")
    ap.add_argument("--keep", action="store_true", help="keep sandboxes for inspection")
    ap.add_argument("--list", action="store_true", help="list the journeys and exit")
    args = ap.parse_args(argv)
    if not args.list and (not args.aibox or not args.evidence):
        ap.error("--aibox and --evidence are required")
    return args


def git_head() -> str | None:
    try:
        cp = subprocess.run(["git", "-C", str(REPO), "rev-parse", "HEAD"], capture_output=True,
                            text=True, timeout=10)
        return cp.stdout.strip() or None
    except (OSError, subprocess.TimeoutExpired):
        return None


def sha256(path: str) -> str:
    h = hashlib.sha256()
    with open(path, "rb") as fh:
        for chunk in iter(lambda: fh.read(1 << 16), b""):
            h.update(chunk)
    return h.hexdigest()


def preflight(aibox: str) -> dict:
    if os.geteuid() == 0:
        raise InfrastructureError("the suite must run as an ordinary user, not root")
    if not os.path.isfile(aibox) or not os.access(aibox, os.X_OK):
        raise InfrastructureError(f"--aibox {aibox} is not an executable file")
    if os.path.exists(CONTRACT.system_settings_path):
        raise InfrastructureError(f"{CONTRACT.system_settings_path} exists on this machine; "
                                  "journeys need a machine without it")
    try:
        wheel = verify()
    except WheelError as exc:
        raise InfrastructureError(str(exc)) from exc
    missing = [jid for jid, _ in CONTRACT.journeys if jid not in REGISTRY]
    extra = [jid for jid in REGISTRY if jid not in dict(CONTRACT.journeys)]
    if extra:
        raise InfrastructureError(f"journeys not in journeys.md: {extra}")
    return {"wheel": str(wheel.name), "unimplemented": missing}


def write_summary(evidence: Path, summary: dict) -> None:
    dump(evidence / "summary.json", summary)
    lines = [f"V1-06 acceptance ({SPEC_ID})",
             f"aibox: {summary.get('aibox')} sha256 {summary.get('aiboxSha256')}",
             f"commit: {summary.get('commit')}",
             f"started {summary.get('started')}, finished {summary.get('finished')}",
             ""]
    if summary.get("infrastructureError"):
        lines.append(f"INFRASTRUCTURE ERROR: {summary['infrastructureError']}")
        lines.append("")
    for rec in summary.get("journeys", []):
        lines.append(f"{rec['outcome']:8} {rec['journey']:9} {rec.get('checksFailed', 0):4} failed checks"
                     f"  covers {', '.join(rec['covers'])}")
    c = summary.get("counts", {})
    lines += ["", f"passed {c.get('passed', 0)}, failed {c.get('failed', 0)}, "
                  f"not run {c.get('not_run', 0)}, total {c.get('total', 0)}"]
    (evidence / "summary.txt").write_text("\n".join(lines) + "\n")


def main(argv: list[str]) -> int:
    args = parse_args(argv)
    if args.list:
        for jid, covers in CONTRACT.journeys:
            j = REGISTRY.get(jid)
            print(f"{jid:9} {'root ' if j and j.root else '     '}{', '.join(covers)}")
        return 0
    aibox = os.path.abspath(args.aibox)
    evidence = Path(args.evidence).resolve()
    try:
        (evidence / "journeys").mkdir(parents=True, exist_ok=True)
    except OSError as exc:
        print(f"run.sh: cannot create evidence directory: {exc}", file=sys.stderr)
        return 2
    unknown = [j for j in args.journey if j not in dict(CONTRACT.journeys)]
    if unknown:
        print(f"run.sh: unknown journey IDs: {unknown}", file=sys.stderr)
        return 2
    selected = [jid for jid, _ in CONTRACT.journeys if not args.journey or jid in args.journey]
    summary: dict = {"spec": SPEC_ID, "commit": git_head(), "aibox": aibox,
                     "aiboxSha256": None, "started": now(), "journeys": []}
    container = None
    work = None
    try:
        info = preflight(aibox)
        summary["aiboxSha256"] = sha256(aibox)
        summary["vendoredWheel"] = info["wheel"]
        base = args.work_dir or tempfile.gettempdir()
        work = os.path.realpath(tempfile.mkdtemp(prefix="aibox-v106-", dir=base))
        os.chmod(work, 0o755)
        staging = tempfile.mkdtemp(prefix="staging-", dir=work)
        ctx = RunContext(aibox=aibox, work=work, keep=args.keep, staging=staging)
        if any(REGISTRY[j].root for j in selected if j in REGISTRY):
            container = RootContainer(aibox)
            print(f"starting root container from {DEBIAN_IMAGE}", file=sys.stderr)
            container.start()
            ctx.container = container
            summary["rootContainer"] = {"image": DEBIAN_IMAGE, "user": "tester", "uid": container.uid}
        records = {}
        for jid, covers in CONTRACT.journeys:
            if jid not in selected or jid not in REGISTRY:
                reason = "filtered out" if jid not in selected else "not implemented"
                records[jid] = {"journey": jid, "covers": covers, "outcome": "not_run",
                                "reason": reason}
            else:
                t0 = time.monotonic()
                print(f"{jid} ...", end="", flush=True, file=sys.stderr)
                records[jid] = run_journey(REGISTRY[jid], ctx)
                print(f" {records[jid]['outcome']} ({time.monotonic() - t0:.1f}s)", file=sys.stderr)
            dump(evidence / "journeys" / f"{jid}.json", records[jid])
        summary["journeys"] = [
            {k: r[k] for k in ("journey", "covers", "outcome", "checksPassed", "checksFailed")
             if k in r} for r in records.values()]
    except InfrastructureError as exc:
        summary["infrastructureError"] = str(exc)
        summary["finished"] = now()
        write_summary(evidence, summary)
        print(f"run.sh: infrastructure error: {exc}", file=sys.stderr)
        return 2
    finally:
        if container is not None:
            container.stop()
        if work and not args.keep:
            shutil.rmtree(work, ignore_errors=True)
    outcomes = [r["outcome"] for r in summary["journeys"]]
    summary["counts"] = {"passed": outcomes.count("passed"), "failed": outcomes.count("failed"),
                         "not_run": outcomes.count("not_run"), "total": len(outcomes)}
    summary["finished"] = now()
    write_summary(evidence, summary)
    print((evidence / "summary.txt").read_text(), file=sys.stderr)
    return 0 if outcomes and all(o == "passed" for o in outcomes) else 1


if __name__ == "__main__":
    try:
        sys.exit(main(sys.argv[1:]))
    except SystemExit:
        raise
    except BaseException:  # noqa: BLE001 - a harness defect is an infrastructure error
        import traceback
        traceback.print_exc()
        sys.exit(2)

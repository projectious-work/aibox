"""Journey registry, soft checks and evidence records."""

from __future__ import annotations

import contextlib
import json
import time
import traceback
from dataclasses import dataclass, field
from typing import Any, Callable

from .proc import InfrastructureError, Result
from .sandbox import ContainerSandbox, LocalSandbox, RootContainer

UNSET = object()


@dataclass
class Journey:
    id: str
    covers: list[str]
    root: bool
    func: Callable


REGISTRY: dict[str, Journey] = {}


def journey(jid: str, covers: list[str], root: bool = False):
    """Register a journey. `covers` lists the requirement IDs of journeys.md."""
    def wrap(func):
        if jid in REGISTRY:
            raise ValueError(f"duplicate journey {jid}")
        REGISTRY[jid] = Journey(jid, list(covers), root, func)
        return func
    return wrap


@dataclass
class RunContext:
    aibox: str
    work: str
    keep: bool = False
    container: RootContainer | None = None
    staging: str = ""


@dataclass
class T:
    """Per-journey state handed to a journey function."""
    ctx: RunContext
    jid: str
    failures: list[dict] = field(default_factory=list)
    passed: int = 0
    _case: list[str] = field(default_factory=list)
    _sandboxes: list = field(default_factory=list)
    observations: list[dict] = field(default_factory=list)

    # sandboxes ---------------------------------------------------------------
    def local(self, label: str = "") -> LocalSandbox:
        sb = LocalSandbox(self.ctx.work, self.ctx.aibox, f"{self.jid}-{label}")
        self._sandboxes.append(sb)
        return sb

    def container(self, label: str = "") -> ContainerSandbox:
        if self.ctx.container is None:
            raise InfrastructureError("root container not started")
        sb = ContainerSandbox(self.ctx.container, self.ctx.staging, f"{self.jid}-{label}")
        self._sandboxes.append(sb)
        return sb

    def cleanup(self) -> None:
        if self.ctx.keep:
            return
        for sb in self._sandboxes:
            try:
                sb.cleanup()
            except Exception:  # noqa: BLE001 - cleanup must not mask results
                pass

    # cases and checks ----------------------------------------------------------
    @contextlib.contextmanager
    def case(self, label: str):
        self._case.append(label)
        try:
            yield
        finally:
            self._case.pop()

    @property
    def case_label(self) -> str:
        return " / ".join(self._case)

    def check(self, ok: bool, message: str, *, result: Result | None = None,
              expected: Any = UNSET, actual: Any = UNSET) -> bool:
        if ok:
            self.passed += 1
            return True
        rec: dict[str, Any] = {"case": self.case_label, "check": message}
        if expected is not UNSET:
            rec["expected"] = expected
        if actual is not UNSET:
            rec["actual"] = actual
        if result is not None:
            rec["observed"] = result.describe()
        self.failures.append(rec)
        return False

    def eq(self, actual: Any, expected: Any, message: str, *, result: Result | None = None) -> bool:
        return self.check(_strict_equal(actual, expected), message, result=result,
                          expected=expected, actual=actual)

    def note(self, text: str) -> None:
        self.observations.append({"case": self.case_label, "note": text})


def _strict_equal(a: Any, b: Any) -> bool:
    """Equality that does not treat True as 1 or 1.0 as 1 by accident."""
    if isinstance(a, bool) or isinstance(b, bool):
        return type(a) is type(b) and a == b
    if isinstance(a, dict) and isinstance(b, dict):
        return a.keys() == b.keys() and all(_strict_equal(a[k], b[k]) for k in a)
    if isinstance(a, list) and isinstance(b, list):
        return len(a) == len(b) and all(_strict_equal(x, y) for x, y in zip(a, b))
    if isinstance(a, (int, float)) and isinstance(b, (int, float)):
        return a == b
    return type(a) is type(b) and a == b


def run_journey(j: Journey, ctx: RunContext) -> dict:
    t = T(ctx, j.id)
    start = time.monotonic()
    harness_error = None
    try:
        j.func(t)
    except InfrastructureError:
        t.cleanup()
        raise
    except Exception:  # noqa: BLE001 - a harness defect is reported, not hidden
        harness_error = traceback.format_exc()
        t.failures.append({"case": t.case_label, "check": "harness exception",
                           "actual": harness_error})
    finally:
        t.cleanup()
    record = {
        "journey": j.id,
        "covers": j.covers,
        "root": j.root,
        "outcome": "passed" if not t.failures else "failed",
        "checksPassed": t.passed,
        "checksFailed": len(t.failures),
        "durationSeconds": round(time.monotonic() - start, 3),
        "failures": t.failures,
    }
    if t.observations:
        record["notes"] = t.observations
    if harness_error:
        record["harnessError"] = True
    return record


def dump(path, obj) -> None:
    with open(path, "w") as fh:
        json.dump(obj, fh, indent=2, sort_keys=False, default=str)
        fh.write("\n")

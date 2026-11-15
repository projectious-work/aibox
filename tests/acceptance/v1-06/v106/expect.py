"""Expectations shared by many journeys (journeys.md "Common setup")."""

from __future__ import annotations

import json
from typing import Any

from .contract import CONTRACT
from .events import bootstrap_reason, check_event_catalog, parse_jsonl
from .framework import UNSET, T
from .proc import Result
from .schema import validate

ABSENT = object()


def result_errors(doc: Any) -> list[str]:
    return validate(doc, CONTRACT.result_schema)


def event_errors(event: Any) -> list[str]:
    return validate(event, CONTRACT.event_schema)


def parse_json_doc(t: T, res: Result, what: str = "stdout") -> dict | None:
    """stdout is exactly one JSON document that validates against
    operation-result.schema.json (CFG-068)."""
    try:
        doc = json.loads(res.stdout.decode("utf-8"))
    except (ValueError, UnicodeDecodeError) as exc:
        t.check(False, f"{what} is one JSON document ({exc})", result=res)
        return None
    if not isinstance(doc, dict):
        t.check(False, f"{what} is a JSON object", result=res)
        return None
    errs = result_errors(doc)
    t.check(not errs, f"{what} validates against operation-result.schema.json",
            result=res, actual=errs[:20])
    return doc


def expect_success(t: T, res: Result, what: str = "result") -> dict | None:
    t.eq(res.exit, 0, f"{what}: exit code", result=res)
    doc = parse_json_doc(t, res, what)
    if doc is not None:
        t.eq(doc.get("outcome"), "succeeded", f"{what}: outcome", result=res)
    return doc


def check_error_envelope(t: T, doc: dict | None, res: Result) -> None:
    """The pre-operation error envelope (OUT-006, output.preOperationFailure)."""
    if doc is None:
        return
    err = doc.get("error") or {}
    t.eq(doc.get("operation"), "invalid_request", "envelope operation", result=res)
    t.eq(doc.get("outcome"), "failed", "envelope outcome", result=res)
    t.eq(err.get("code"), "invalid_input", "envelope error code", result=res)


def check_events_valid(t: T, events: list[dict], res: Result, where: str = "stderr") -> None:
    for ev in events:
        errs = event_errors(ev)
        t.check(not errs, f"{where}: event {ev.get('event')} validates against log-event.schema.json",
                result=res, actual=errs)
        problems = check_event_catalog(ev)
        t.check(not problems, f"{where}: event {ev.get('event')} matches logging.events",
                result=res, actual=problems)


def expect_rejected(t: T, res: Result, reason: str, *, envelope: bool = True,
                    key: Any = UNSET, layer: Any = UNSET, suggestion: Any = UNSET) -> dict | None:
    """journeys.md "Rejected with reason R"."""
    t.eq(res.exit, 2, f"rejected ({reason}): exit code", result=res)
    doc = None
    if envelope:
        doc = parse_json_doc(t, res, "error envelope")
        check_error_envelope(t, doc, res)
    else:
        t.eq(res.out, "", "stdout is empty", result=res)
    parsed = parse_jsonl(res.stderr)
    rejected = parsed.named("settings.rejected")
    t.eq(len(rejected), 1, "stderr holds exactly one settings.rejected event in jsonl",
         result=res)
    t.check(not parsed.named("operation.started"), "no operation.started event", result=res)
    if len(rejected) != 1:
        return doc
    ev = rejected[0]
    check_events_valid(t, [ev], res)
    fields = ev.get("fields") or {}
    t.eq(fields.get("reason"), reason, "settings.rejected fields.reason", result=res)
    t.eq(ev.get("sequence"), 0, "settings.rejected sequence", result=res)
    if doc is not None:
        t.eq(ev.get("requestId"), doc.get("requestId"),
             "settings.rejected requestId equals the envelope's", result=res)
    if key is not UNSET:
        t.eq(fields.get("key"), key, "settings.rejected fields.key", result=res)
    if layer is not UNSET:
        t.eq(fields.get("layer"), layer, "settings.rejected fields.layer", result=res)
    if suggestion is not UNSET:
        t.eq(fields.get("suggestion"), suggestion, "settings.rejected fields.suggestion", result=res)
    return doc


def expect_bootstrap(t: T, res: Result, reason: str | None = None) -> None:
    """OUT-007: exit 2, empty stdout, last stderr line is the diagnostic."""
    t.eq(res.exit, 2, "exit code", result=res)
    t.eq(res.out, "", "stdout is empty", result=res)
    lines = res.err.rstrip("\n").split("\n") if res.err else []
    last = lines[-1] if lines else None
    t.check(last is not None and CONTRACT.bootstrap_re.match(last) is not None,
            "last stderr line matches output.bootstrapDiagnostic", result=res, actual=last)
    if reason is not None:
        t.eq(bootstrap_reason(last), reason, "bootstrap diagnostic reason", result=res)


# --- the effective-configuration view ---------------------------------------------

def configuration(doc: dict | None) -> list[dict]:
    if not isinstance(doc, dict):
        return []
    data = doc.get("data")
    if not isinstance(data, dict):
        return []
    conf = data.get("configuration")
    return conf if isinstance(conf, list) else []


def entry(doc: dict | None, key: str) -> dict | None:
    for e in configuration(doc):
        if isinstance(e, dict) and e.get("key") == key:
            return e
    return None


def settings_entries(doc: dict | None) -> list[dict]:
    return [e for e in configuration(doc) if isinstance(e, dict) and e.get("domain") != "native"]


def settings_sources(doc: dict | None):
    if not isinstance(doc, dict) or not isinstance(doc.get("data"), dict):
        return None
    return doc["data"].get("settingsSources")


def check_entry(t: T, doc: dict | None, key: str, res: Result, *, value: Any = UNSET,
                layer: Any = UNSET, source: Any = UNSET, invocation_only: Any = UNSET,
                domain: Any = UNSET, redacted: Any = UNSET, overridden: Any = UNSET) -> dict | None:
    e = entry(doc, key)
    if not t.check(e is not None, f"entry({key}) present", result=res):
        return None
    for name, want in (("value", value), ("layer", layer), ("source", source),
                       ("invocationOnly", invocation_only), ("domain", domain),
                       ("redacted", redacted)):
        if want is not UNSET:
            t.eq(e.get(name, "<absent>"), want,
                 f"entry({key}).{name}", result=res)
    if overridden is not UNSET:
        if overridden is None:
            t.check("overridden" not in e, f"entry({key}) has no overridden", result=res,
                    actual=e.get("overridden"))
        else:
            t.eq(e.get("overridden"), overridden, f"entry({key}).overridden", result=res)
    return e

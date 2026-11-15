"""Result output journeys J-OUT-01 to J-OUT-05 of journeys.md."""

from __future__ import annotations

import copy
import json

from .contract import CONTRACT
from .expect import check_error_envelope, expect_bootstrap, expect_success, result_errors
from .framework import T, journey
from .proc import run_pty
from .wheel import load_yaml_documents

NOT_JSON = "not json"


def without_volatile(doc, view: bool = False):
    """J-OUT-02 compares ignoring requestId and observedAt.

    view=True also drops entry(output.format).value: in the View, CFG-061 makes
    it the winning --format value, json in one run and yaml in the other.
    Integration-owner decision (test-defect route, dev-notes V1-06): every other
    member, including the rest of that entry, must still be equal."""
    doc = copy.deepcopy(doc)
    if isinstance(doc, dict):
        doc.pop("requestId", None)
        data = doc.get("data")
        if isinstance(data, dict):
            data.pop("observedAt", None)
            if view and isinstance(data.get("configuration"), list):
                for e in data["configuration"]:
                    if isinstance(e, dict) and e.get("key") == "output.format":
                        e.pop("value", None)
    return doc


def parse_yaml_doc(t: T, res, what: str = "YAML stdout"):
    try:
        docs = load_yaml_documents(res.out)
    except Exception as exc:  # noqa: BLE001 - any parser error is a failed check
        t.check(False, f"{what} parses with a YAML 1.2 parser ({type(exc).__name__}: {exc})",
                result=res)
        return None
    if not t.eq(len(docs), 1, f"{what} is one YAML document", result=res):
        return None
    return docs[0]


@journey("J-OUT-01", ["OUT-001", "OUT-002", "LOG-007", "CFG-068"])
def j_out_01(t: T) -> None:
    sb = t.local()
    p = sb.project("P")
    res = sb.aibox_run(["inspect", "--project", p, "--format", "json", "--log-level", "trace"])
    t.eq(res.exit, 0, "exit code", result=res)
    out = res.out
    t.check(out.endswith("\n") and not out.endswith("\n\n"),
            "stdout ends with exactly one newline", result=res)
    try:
        doc = json.loads(out[:-1] if out.endswith("\n") else out)
    except ValueError as exc:
        t.check(False, f"stdout is one JSON document ({exc})", result=res)
        doc = None
    if doc is not None:
        errs = result_errors(doc)
        t.check(not errs, "stdout validates against operation-result.schema.json", result=res,
                actual=errs[:20])
    found = [name for name in CONTRACT.event_names if name in out]
    t.eq(found, [], "no event name from logging.events occurs on stdout", result=res)


@journey("J-OUT-02", ["OUT-003"])
def j_out_02(t: T) -> None:
    sb = t.local()
    p = sb.project("P")
    plain = ["inspect", "--project", p, "--log-level", "trace", "--format"]
    for label, json_args, yaml_args, view in (
            ("inspect --format yaml", plain + ["json"], plain + ["yaml"], False),
            ("View --format yaml", sb.view_args(p), sb.view_args(p, fmt="yaml"), True)):
        with t.case(label):
            jres = sb.aibox_run(json_args)
            jdoc = expect_success(t, jres, "JSON form")
            yres = sb.aibox_run(yaml_args)
            t.eq(yres.exit, 0, "YAML run exit code", result=yres)
            ydoc = parse_yaml_doc(t, yres)
            if jdoc is None or ydoc is None:
                continue
            what = "YAML value equals the JSON form, ignoring requestId and observedAt"
            if view:
                what += " and entry(output.format).value"
            t.eq(without_volatile(ydoc, view), without_volatile(jdoc, view), what, result=yres)


@journey("J-OUT-03", ["OUT-004", "OUT-005"])
def j_out_03(t: T) -> None:
    sb = t.local()
    p = sb.project("P")
    missing = f"{sb.projects}/missing"
    variants = [
        ("(a) pipe", False, {}, []),
        ("(b) pseudo-terminal", True, {}, []),
        ("(c) pseudo-terminal, NO_COLOR=1", True, {"NO_COLOR": "1"}, []),
        ("(d) pipe, --color always", False, {}, ["--color", "always"]),
        ("(e) pipe, --color always, NO_COLOR=1", False, {"NO_COLOR": "1"}, ["--color", "always"]),
    ]
    colour = {"(b) pseudo-terminal", "(d) pipe, --color always"}
    for project, words, exit_code in ((p, ("inspect_workspace", "succeeded"), 0),
                                      (missing, ("failed", "invalid_input"), 2)):
        name = "P" if project == p else "missing project"
        for label, tty, env, extra in variants:
            with t.case(f"{name} {label}"):
                args = [sb.aibox, "inspect", "--project", project, *extra]
                if tty:
                    res = run_pty(args, env=sb.env(env), cwd=sb.cwd)
                else:
                    res = sb.exec(args, sb.env(env), sb.cwd)
                t.eq(res.exit, exit_code, "exit code", result=res)
                for word in words:
                    t.check(word in res.out, f"the text contains {word}", result=res)
                t.eq(b"\x1b" in res.stdout, label in colour,
                     "byte 0x1B appears in (b) and (d) only", result=res)


@journey("J-OUT-04", ["OUT-006", "OUT-008"])
def j_out_04(t: T) -> None:
    for label, fmt, user, extra in (("--format json, user file not JSON", "json", NOT_JSON, []),
                                    ("--format yaml, user file not JSON", "yaml", NOT_JSON, []),
                                    ("--format json --log-level loud", "json", None, ["--log-level", "loud"])):
        with t.case(label):
            sb = t.local(label)
            if user is not None:
                sb.user_settings(user)
            p = sb.project("P")
            res = sb.aibox_run(["inspect", "--project", p, "--format", fmt, *extra])
            t.eq(res.exit, 2, "exit code", result=res)
            if fmt == "json":
                try:
                    doc = json.loads(res.out)
                except ValueError as exc:
                    t.check(False, f"stdout is one JSON document ({exc})", result=res)
                    continue
            else:
                doc = parse_yaml_doc(t, res)
                if doc is None:
                    continue
            if not t.check(isinstance(doc, dict), "the envelope is an object", result=res):
                continue
            errs = result_errors(doc)
            t.check(not errs, "envelope validates against operation-result.schema.json",
                    result=res, actual=errs[:20])
            check_error_envelope(t, doc, res)
            err = doc.get("error") or {}
            spec = CONTRACT.error("invalid_input")
            t.eq(err.get("nextAction"), spec["nextAction"], "error.nextAction", result=res)
            t.check("target" not in doc, "no target", result=res)
            t.check("inputDigest" not in doc, "no inputDigest", result=res)


@journey("J-OUT-05", ["OUT-007"])
def j_out_05(t: T) -> None:
    for label, user, env, extra, reason in (
            ("(a) user file not JSON, no format", NOT_JSON, {}, [], "malformed"),
            ("(b) user file not JSON, AIBOX_OUTPUT__FORMAT=json", NOT_JSON,
             {"AIBOX_OUTPUT__FORMAT": "json"}, [], "malformed"),
            ("(c) --format xml", None, {}, ["--format", "xml"], "invalid_value")):
        with t.case(label):
            sb = t.local(label)
            if user is not None:
                sb.user_settings(user)
            p = sb.project("P")
            res = sb.aibox_run(["inspect", "--project", p, *extra], env)
            expect_bootstrap(t, res, reason)

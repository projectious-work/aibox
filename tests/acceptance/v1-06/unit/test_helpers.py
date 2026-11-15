"""Unit tests of the V1-06 harness helpers. They need no aibox executable.

Run: python3 -m unittest discover -s tests/acceptance/v1-06/unit -t tests/acceptance/v1-06
"""

from __future__ import annotations

import json
import os
import shutil
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock

SUITE = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(SUITE))

from v106 import envscrub, events, jsonc, wheel  # noqa: E402
from v106 import journeys_cfg, journeys_log, journeys_out  # noqa: E402,F401
from v106.contract import CONTRACT  # noqa: E402
from v106.expect import (check_entry, event_errors, expect_bootstrap, expect_rejected,  # noqa: E402
                         result_errors)
from v106.framework import REGISTRY, RunContext, T  # noqa: E402
from v106.layers import PAIR_VALUES, nested, text  # noqa: E402
from v106.proc import Result, run  # noqa: E402
from v106.schema import UnsupportedKeyword, validate  # noqa: E402

RID = "0123456789abcdef0123456789abcdef"


def event(name: str, seq: int = 0, rid: str = RID, **members) -> dict:
    spec = CONTRACT.event_spec(name)
    ev = {"schemaVersion": "aibox.log-event/v1", "timestamp": "2026-11-15T12:00:00.123Z",
          "severity": spec["severity"], "event": name, "component": spec["component"],
          "requestId": rid, "sequence": seq, "executor": "aibox"}
    ev.update(members)
    return ev


def rejection_event(reason: str, rid: str = RID, **fields) -> dict:
    return event("settings.rejected", rid=rid, outcome="failed", errorCategory="invalid_input",
                 fields={"reason": reason, **fields})


def error_envelope(rid: str = RID) -> dict:
    return {"schemaVersion": "aibox.operation-result/v1", "operation": "invalid_request",
            "requestId": rid, "actors": {"initiator": "cli", "executor": "aibox"},
            "outcome": "failed", "changedResources": [], "warnings": [], "evidence": [],
            "error": {"code": "invalid_input", "category": "invalid_input",
                      "message": "settings rejected", "retryable": False,
                      "nextAction": "edit_configuration"}}


def success_result(configuration: list, sources: list | None = None) -> dict:
    data = {"state": "unknown", "resources": [], "configuration": configuration,
            "observedAt": "2026-11-15T12:00:00Z"}
    if sources is not None:
        data["settingsSources"] = sources
    return {"schemaVersion": "aibox.operation-result/v1", "operation": "inspect_workspace",
            "requestId": RID, "actors": {"initiator": "cli", "executor": "aibox"},
            "target": {"scope": "local", "workspaceRoot": "/w/p"},
            "inputDigest": "sha256:" + "0" * 64, "outcome": "succeeded",
            "changedResources": [], "warnings": [], "evidence": [], "data": data}


def fake_result(stdout: str | bytes = b"", stderr: str | bytes = b"", exit: int = 0) -> Result:
    enc = lambda s: s.encode() if isinstance(s, str) else s  # noqa: E731
    return Result(["aibox"], "/", {}, exit, enc(stdout), enc(stderr))


def new_t() -> T:
    return T(RunContext(aibox="/bin/false", work=tempfile.gettempdir()), "J-TEST")


class EnvScrubTest(unittest.TestCase):
    def test_only_home_path_and_additions(self):
        polluted = {"XDG_CONFIG_HOME": "/x", "XDG_STATE_HOME": "/y", "AIBOX_LOGGING__LEVEL": "trace",
                    "LANG": "C.UTF-8", "HOME": "/developer"}
        with mock.patch.dict(os.environ, polluted):
            env = envscrub.journey_env("/fresh/home", {"AIBOX_ENV_FILE": "/a.env"})
        self.assertEqual(env, {"PATH": envscrub.DEFAULT_PATH, "HOME": "/fresh/home",
                               "AIBOX_ENV_FILE": "/a.env"})

    def test_home_unset(self):
        self.assertEqual(envscrub.journey_env(None), {"PATH": envscrub.DEFAULT_PATH})

    def test_rejects_bad_entries(self):
        for bad in ({"A=B": "x"}, {"": "x"}, {"A": 1}):
            with self.assertRaises(ValueError):
                envscrub.journey_env("/h", bad)

    def test_child_sees_only_the_scrubbed_environment(self):
        with mock.patch.dict(os.environ, {"AIBOX_LEAK": "1", "XDG_CONFIG_HOME": "/x"}):
            res = run(["/usr/bin/env"], env=envscrub.journey_env("/h", {"EXTRA": "1"}), cwd="/")
        self.assertEqual(sorted(res.out.split()), ["EXTRA=1", "HOME=/h", f"PATH={envscrub.DEFAULT_PATH}"])

    def test_env_argv(self):
        self.assertEqual(envscrub.env_argv({"HOME": "/h", "A": "b c"}),
                         ["/usr/bin/env", "-i", "HOME=/h", "A=b c"])


class EventParsingTest(unittest.TestCase):
    def test_jsonl_split_events_and_other_lines(self):
        data = (json.dumps(event("settings.resolved")) + "\n" + "aibox: invalid_input: malformed: x\n"
                + json.dumps({"no": "event"}) + "\n").encode()
        parsed = events.parse_jsonl(data)
        self.assertEqual(parsed.names(), ["settings.resolved"])
        self.assertEqual([i for i, _ in parsed.other], [1, 2])
        self.assertEqual(parsed.last_line, '{"no": "event"}')

    def test_text_lines(self):
        self.assertEqual(events.text_lines(b""), [])
        self.assertEqual(events.text_lines(b"a\nb\n"), ["a", "b"])
        self.assertEqual(events.text_lines("a\n\nb"), ["a", "", "b"])

    def test_plain_line(self):
        line = '2026-11-15T12:00:00.5Z INFO operation.started operation=inspect_workspace projectRoot="/a b"'
        p = events.parse_plain(line)
        self.assertEqual((p.severity, p.event), ("info", "operation.started"))
        self.assertIsNone(events.parse_plain("2026-11-15 INFO x"))
        self.assertIsNone(events.parse_plain('{"event":"x"}'))

    def test_bootstrap_reason(self):
        self.assertEqual(events.bootstrap_reason("aibox: invalid_input: usage_error: unknown flag --x"),
                         "usage_error")
        self.assertIsNone(events.bootstrap_reason("aibox: invalid_input: Usage: x"))
        self.assertIsNone(events.bootstrap_reason("flag provided but not defined: -x"))
        self.assertIsNone(events.bootstrap_reason(None))

    def test_event_catalog_check(self):
        ok = event("operation.started", operation="inspect_workspace", initiator="cli",
                   fields={"projectRoot": "/p"})
        self.assertEqual(events.check_event_catalog(ok), [])
        wrong = dict(ok, severity="debug", operation="other")
        self.assertEqual(len(events.check_event_catalog(wrong)), 2)
        self.assertTrue(events.check_event_catalog(dict(ok, event="nope")))
        missing = dict(ok)
        del missing["initiator"]
        self.assertTrue(events.check_event_catalog(missing))

    def test_severity_rank_and_level_filter(self):
        order = CONTRACT.interfaces["logging"]["eventOrder"]["inspect, success"]
        self.assertEqual(journeys_log.names_at("trace", order), order)
        self.assertEqual(journeys_log.names_at("info", order), ["operation.started", "operation.finished"])
        self.assertEqual(journeys_log.names_at("warn", order), [])
        self.assertEqual(journeys_log.names_at("off", order), [])


class SchemaTest(unittest.TestCase):
    def test_contract_examples_validate(self):
        self.assertEqual(result_errors(error_envelope()), [])
        self.assertEqual(result_errors(success_result([CONTRACT.native_entry("devcontainer.config")],
                                                      [{"layer": "user", "path": "~/x", "state": "absent"}])), [])
        self.assertEqual(event_errors(event("settings.resolved", fields={"layers": ""})), [])

    def test_result_defects(self):
        env = error_envelope()
        env["target"] = {"scope": "local", "workspaceRoot": "/w"}
        self.assertEqual(result_errors(env), [])  # target is allowed by the schema itself
        bad = error_envelope()
        del bad["error"]
        self.assertTrue(result_errors(bad))  # failed outcome needs error (if/then)
        bad = error_envelope()
        bad["outcome"] = "succeeded"  # invalid_request only failed or refused
        self.assertTrue(result_errors(bad))
        bad = success_result([dict(CONTRACT.native_entry("devcontainer.config"), layer="nope")])
        self.assertTrue(result_errors(bad))
        bad = success_result([], [{"layer": "built_in", "path": "x", "state": "loaded"}])
        self.assertTrue(result_errors(bad))
        bad = success_result([])
        bad["data"]["extra"] = 1
        self.assertTrue(result_errors(bad))
        bad = success_result([])
        del bad["inputDigest"]
        self.assertTrue(result_errors(bad))
        bad = success_result([])
        bad["error"] = error_envelope()["error"]
        self.assertTrue(result_errors(bad))  # not: required error

    def test_event_defects(self):
        for mutate in (lambda e: e.update(sequence=-1), lambda e: e.update(severity="fatal"),
                       lambda e: e.update(timestamp="2026-11-15 12:00:00"),
                       lambda e: e.update(fields={"n": 1}), lambda e: e.update(extra="x"),
                       lambda e: e.pop("requestId"), lambda e: e.update(sequence=True)):
            ev = event("settings.resolved")
            mutate(ev)
            self.assertTrue(event_errors(ev), ev)

    def test_integer_and_bool_types(self):
        schema = {"type": "integer", "minimum": 1}
        self.assertEqual(validate(3, schema), [])
        self.assertEqual(validate(3.0, schema), [])
        self.assertTrue(validate(True, schema))
        self.assertTrue(validate(0, schema))
        self.assertTrue(validate(False, {"enum": [0]}))
        self.assertTrue(validate(0, {"const": False}))

    def test_unknown_keyword_is_refused(self):
        with self.assertRaises(UnsupportedKeyword):
            validate({}, {"oneOf": [{}]})

    def test_settings_schema_accepts_pair_values(self):
        for key, values in PAIR_VALUES.items():
            for value in values:
                doc = {"schemaVersion": "1", **nested(key, value)}
                self.assertEqual(validate(doc, CONTRACT.settings_schema), [], (key, value))

    def test_settings_schema_matches_j_cfg_07_cases(self):
        valid_json = []
        for label, content, _, reason in journeys_cfg.J07_CASES:
            try:
                doc = json.loads(content)
            except ValueError:
                continue
            valid_json.append(label)
            if reason in ("unknown_key", "unsupported_version", "invalid_value"):
                self.assertTrue(validate(doc, CONTRACT.settings_schema), label)
        self.assertIn("output null", valid_json)


class WheelTest(unittest.TestCase):
    def setUp(self):
        self.dir = Path(tempfile.mkdtemp())

    def tearDown(self):
        shutil.rmtree(self.dir)

    def _vendor(self, content: bytes, recorded: str | None = None, write: bool = True):
        if write:
            (self.dir / "w.whl").write_bytes(content)
        digest = recorded or __import__("hashlib").sha256(content).hexdigest()
        (self.dir / "wheels.json").write_text(json.dumps({"wheels": [
            {"name": "ruamel.yaml", "version": "0", "file": "w.whl", "url": "u", "sha256": digest}]}))

    def test_vendored_wheel_matches(self):
        path = wheel.verify()
        self.assertEqual(path.name, "ruamel_yaml-0.19.1-py3-none-any.whl")
        entry = wheel.manifest()["wheels"][0]
        self.assertTrue(entry["url"].startswith("https://files.pythonhosted.org/"))
        self.assertEqual(len(entry["sha256"]), 64)

    def test_good_hash(self):
        self._vendor(b"abc")
        self.assertEqual(wheel.verify(vendor=self.dir), self.dir / "w.whl")

    def test_mismatch(self):
        self._vendor(b"abc", recorded="0" * 64)
        with self.assertRaisesRegex(wheel.WheelError, "mismatch"):
            wheel.verify(vendor=self.dir)

    def test_missing(self):
        self._vendor(b"abc", write=False)
        with self.assertRaisesRegex(wheel.WheelError, "missing"):
            wheel.verify(vendor=self.dir)

    def test_no_manifest(self):
        with self.assertRaises(wheel.WheelError):
            wheel.verify(vendor=self.dir)

    def test_yaml_1_2_parser(self):
        docs = wheel.load_yaml_documents('"a": "x"\n"b": [1, 2.5, true, null]\n"c": {}\n"d": "yes"\n'
                                         'e: on\n')
        self.assertEqual(docs, [{"a": "x", "b": [1, 2.5, True, None], "c": {}, "d": "yes", "e": "on"}])
        self.assertEqual(len(wheel.load_yaml_documents("a: 1\n---\nb: 2\n")), 2)


class JsoncTest(unittest.TestCase):
    def test_strip_comments_and_trailing_commas(self):
        text = '// c\n{ "a": "x // not a comment", /* b */ "b": [1, 2,], "c": "\\"q\\"", }\n'
        self.assertEqual(jsonc.loads(text), {"a": "x // not a comment", "b": [1, 2], "c": '"q"'})

    def test_add_members_to_fixture(self):
        text = (CONTRACT.fixture_minimal / ".devcontainer" / "devcontainer.json").read_text()
        out = jsonc.add_members(text, {"customizations": {"aibox": {"logging": {"level": "trace"}}}})
        parsed = jsonc.loads(out)
        self.assertEqual(parsed["customizations"], {"aibox": {"logging": {"level": "trace"}}})
        self.assertEqual(parsed["name"], jsonc.loads(text)["name"])

    def test_w_grid_edit(self):
        text = (CONTRACT.fixture_workspace / ".devcontainer" / "devcontainer.json").read_text()
        out = journeys_cfg.w_with_grid(text)
        self.assertEqual(jsonc.loads(out)["customizations"]["aibox"]["workspace"],
                         {"theme": "nord", "mode": "dark", "layout": "grid"})

    def test_feature_references_sorted(self):
        self.assertEqual(jsonc.feature_references('{"features": {"b": {}, "a": {}, "B": {}}}'),
                         ["B", "a", "b"])
        self.assertEqual(jsonc.feature_references("{}"), [])


class ContractTablesTest(unittest.TestCase):
    def test_every_journey_registered_with_contract_covers(self):
        table = dict(CONTRACT.journeys)
        self.assertEqual(len(table), 41)
        self.assertEqual(sorted(REGISTRY), sorted(table))
        for jid, covers in table.items():
            self.assertEqual(REGISTRY[jid].covers, covers, jid)

    def test_root_journeys(self):
        self.assertEqual(sorted(j for j in REGISTRY if REGISTRY[j].root),
                         ["J-CFG-02", "J-CFG-03", "J-CFG-06", "J-CFG-08"])

    def test_range_expansion_skips_retired(self):
        self.assertNotIn("CFG-037", CONTRACT.expand_covers("CFG-032 to CFG-038"))
        self.assertEqual(CONTRACT.expand_covers("LOG-001 to LOG-003"), ["LOG-001", "LOG-002", "LOG-003"])

    def test_every_j_requirement_is_covered(self):
        covered = {r for _, covers in CONTRACT.journeys for r in covers}
        j_level = {r for r, level in CONTRACT.requirement_levels.items() if level.startswith("J")}
        self.assertEqual(sorted(j_level - covered), [])

    def test_catalog_shape(self):
        self.assertEqual(len(CONTRACT.keys), 15)
        self.assertEqual(CONTRACT.default("execution.timeoutSeconds"), 30)
        self.assertEqual(CONTRACT.default("output.format", "mcp"), "json")
        self.assertEqual(text(True), "true")
        self.assertEqual(text(5), "5")


class ExpectationTest(unittest.TestCase):
    def test_rejection_accepted(self):
        t = new_t()
        res = fake_result(json.dumps(error_envelope()) + "\n",
                          json.dumps(rejection_event("malformed", key="logging.level", layer="user")) + "\n",
                          exit=2)
        expect_rejected(t, res, "malformed", key="logging.level", layer="user")
        self.assertEqual(t.failures, [])

    def test_rejection_defects_reported(self):
        cases = {
            "wrong reason": (error_envelope(), [rejection_event("unknown_key")], 2),
            "exit 1": (error_envelope(), [rejection_event("malformed")], 1),
            "request ID differs": (error_envelope(), [rejection_event("malformed", rid="other")], 2),
            "two events": (error_envelope(), [rejection_event("malformed")] * 2, 2),
            "operation started": (error_envelope(), [rejection_event("malformed"),
                                                     event("operation.started", 1)], 2),
            "sequence 1": (error_envelope(), [dict(rejection_event("malformed"), sequence=1)], 2),
        }
        for label, (doc, evs, code) in cases.items():
            t = new_t()
            res = fake_result(json.dumps(doc), "".join(json.dumps(e) + "\n" for e in evs), exit=code)
            expect_rejected(t, res, "malformed")
            self.assertTrue(t.failures, label)

    def test_rejection_without_envelope_needs_empty_stdout(self):
        t = new_t()
        err = json.dumps(rejection_event("invalid_value")) + "\naibox: invalid_input: invalid_value: plain\n"
        expect_rejected(t, fake_result("", err, 2), "invalid_value", envelope=False)
        expect_bootstrap(t, fake_result("", err, 2), "invalid_value")
        self.assertEqual(t.failures, [])
        t = new_t()
        expect_rejected(t, fake_result("{}", err, 2), "invalid_value", envelope=False)
        self.assertTrue(t.failures)

    def test_bootstrap_must_be_last_line(self):
        t = new_t()
        expect_bootstrap(t, fake_result("", "aibox: invalid_input: malformed: x\nmore\n", 2), "malformed")
        self.assertTrue(t.failures)

    def test_check_entry(self):
        doc = success_result([{"key": "logging.level", "value": "debug", "redacted": False,
                               "layer": "user", "source": "~/.config/aibox/settings.json",
                               "domain": "process", "invocationOnly": False}])
        res = fake_result(json.dumps(doc))
        t = new_t()
        check_entry(t, doc, "logging.level", res, value="debug", layer="user", overridden=None)
        self.assertEqual(t.failures, [])
        check_entry(t, doc, "logging.level", res, value="info", invocation_only=True,
                    overridden=[])
        self.assertEqual(len(t.failures), 3)
        check_entry(t, doc, "output.color", res)
        self.assertEqual(len(t.failures), 4)

    def test_strict_equality(self):
        t = new_t()
        self.assertFalse(t.eq(True, 1, "bool is not int"))
        self.assertFalse(t.eq({"a": 1}, {"a": True}, "nested"))
        self.assertTrue(t.eq(30, 30.0, "numbers"))
        self.assertFalse(t.eq("30", 30, "string is not number"))

    def test_yaml_comparison_ignores_volatile_members(self):
        a = success_result([])
        b = success_result([])
        b["requestId"] = "x"
        b["data"]["observedAt"] = "later"
        self.assertEqual(journeys_out.without_volatile(a), journeys_out.without_volatile(b))


class RotationHelpersTest(unittest.TestCase):
    def test_rotation_set(self):
        names = ["r.log", "r.log.1", "r.log.99", "r.log.100", "r.log.lock", "r.log.0", "other.txt", "r.logx"]
        self.assertEqual(journeys_log.rotation_set(names, "r.log"),
                         ["r.log", "r.log.1", "r.log.99", "r.log.lock"])

    def test_fill_is_one_mib(self):
        written = {}

        class SB:
            def write(self, path, data, mode=0o600):
                written[path] = data

        journeys_log.fill(SB(), "/x", "A")
        self.assertEqual(len(written["/x"]), 1 << 20)
        self.assertTrue(written["/x"].startswith(b"A\nf"))
        self.assertTrue(written["/x"].endswith(b"\n"))


FAKE_SERVER = r"""
import json, sys
pending = []
for line in sys.stdin:
    msg = json.loads(line)
    if "id" not in msg:
        continue
    if msg["method"] == "initialize":
        print(json.dumps({"jsonrpc": "2.0", "id": msg["id"], "result": {"serverInfo": {"name": "fake"}}}), flush=True)
        continue
    pending.append(msg)
    if len(pending) == 2:  # answer out of order, with a notification between
        for m in reversed(pending):
            print(json.dumps({"jsonrpc": "2.0", "method": "notifications/message"}), flush=True)
            args = m["params"]["arguments"]
            if "bad" in args:
                print(json.dumps({"jsonrpc": "2.0", "id": m["id"], "result": {"isError": True,
                      "content": [{"type": "text", "text": "invalid"}]}}), flush=True)
            else:
                print(json.dumps({"jsonrpc": "2.0", "id": m["id"], "result": {"structuredContent":
                      {"schemaVersion": "x", "requestId": args["requestId"]}}}), flush=True)
        pending = []
sys.stderr.write('{"event": "process.finished"}\n')
"""


class MCPClientTest(unittest.TestCase):
    def test_responses_matched_by_id(self):
        import threading
        from v106.mcp import MCPSession, is_rejected, structured
        session = MCPSession([sys.executable, "-I", "-c", FAKE_SERVER], {"PATH": "/usr/bin:/bin"}, "/")
        self.assertEqual(session.initialize()["result"]["serverInfo"]["name"], "fake")
        results = {}

        def call(name, args):
            results[name] = session.call_tool("inspect_workspace", args)

        a = threading.Thread(target=call, args=("ok", {"requestId": "r1"}))
        a.start()
        import time
        time.sleep(0.2)
        # The second call is sent from this thread; responses arrive in reverse order.
        results["bad"] = session.call_tool("inspect_workspace", {"requestId": "r2", "bad": 1})
        a.join(10)
        res = session.close()
        self.assertTrue(is_rejected(results["bad"]))
        self.assertFalse(is_rejected(results["ok"]))
        self.assertEqual(structured(results["ok"]), {"schemaVersion": "x", "requestId": "r1"})
        self.assertEqual(res.exit, 0)
        self.assertIn(b"process.finished", res.stderr)


if __name__ == "__main__":
    unittest.main()

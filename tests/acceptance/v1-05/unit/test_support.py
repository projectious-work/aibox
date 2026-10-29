"""Metadata, JSONC, inspection and trace helpers."""

import io
import json
import tarfile
import tempfile
import unittest
from pathlib import Path

from actools import inspection, jsonc, metadata, trace
from actools.catalog import REPO, SPEC
from unit.fixtures import CENSUS, REAL


def layer(document: dict, name: str = "./devcontainer-feature.json") -> bytes:
    buffer = io.BytesIO()
    with tarfile.open(fileobj=buffer, mode="w") as archive:
        data = json.dumps(document).encode()
        info = tarfile.TarInfo(name)
        info.size = len(data)
        archive.addfile(info, io.BytesIO(data))
    return buffer.getvalue()


class MetadataTest(unittest.TestCase):
    def test_reads_feature_json_from_layer(self):
        self.assertEqual(metadata.feature_json_from_layer(layer({"id": "uv"})), {"id": "uv"})
        with self.assertRaises(ValueError):
            metadata.feature_json_from_layer(layer({}, "install.sh"))

    def test_clean_metadata_has_no_violations(self):
        env = {"PATH": "/usr/local/cargo/bin:${PATH}"}
        self.assertEqual(metadata.side_effect_violations(
            {"id": "rustc", "containerEnv": env, "privileged": False}, env), [])
        self.assertEqual(metadata.side_effect_violations({"id": "uv"}, {}), [])

    def test_each_forbidden_declaration_is_reported(self):
        violations = metadata.side_effect_violations(
            {"postCreateCommand": "x", "mounts": [{"source": "a", "target": "/b"}],
             "privileged": True, "capAdd": ["SYS_PTRACE"], "containerEnv": {"A": "1"}}, {})
        self.assertEqual(len(violations), 5)

    def test_snapshot_detects_added_modified_and_deleted_entries(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "keep").write_text("a")
            (root / "change").write_text("a")
            (root / "drop").write_text("a")
            before = metadata.snapshot(root)
            (root / "change").write_text("b")
            (root / "drop").unlink()
            (root / "new").mkdir()
            self.assertEqual(metadata.snapshot_changes(before, metadata.snapshot(root)),
                             ["modified change", "deleted drop", "added new"])


class JsoncTest(unittest.TestCase):
    def test_comments_and_trailing_commas(self):
        text = '{\n // c\n "a": "x // not a comment", /* b */ "b": [1, 2,],\n}'
        self.assertEqual(jsonc.loads(text), {"a": "x // not a comment", "b": [1, 2]})

    def test_escaped_quote_in_string(self):
        self.assertEqual(jsonc.loads('{"a": "q\\" // x",}'), {"a": 'q" // x'})

    def test_maintained_examples_parse(self):
        for example in ("minimal", "customized"):
            path = SPEC / "examples" / example / ".devcontainer" / "devcontainer.json"
            self.assertIn("name", jsonc.loads(path.read_text()))


class InspectionTest(unittest.TestCase):
    def config(self, example):
        return (SPEC / "examples" / example / ".devcontainer" / "devcontainer.json").read_text()

    def result(self, entries):
        return {"outcome": "succeeded", "data": {"configuration": [
            {"key": "devcontainer.config", "value": ".devcontainer/devcontainer.json"},
            *entries]}}

    def test_expected_entry_lists_sorted_references_without_options(self):
        entry = inspection.expected_entry(self.config("customized"))
        self.assertEqual(entry["key"], "devcontainer.features")
        self.assertEqual(entry["value"], sorted(jsonc.loads(self.config("customized"))["features"]))
        self.assertEqual((entry["layer"], entry["source"], entry["domain"], entry["redacted"],
                          entry["invocationOnly"]),
                         ("project", ".devcontainer/devcontainer.json", "native", False, False))
        self.assertIsNone(inspection.expected_entry(self.config("minimal")))

    def test_identical_correct_results_pass(self):
        text = self.config("customized")
        expected = inspection.expected_entry(text)
        self.assertEqual(inspection.inspect_failures(
            expected, self.result([expected]), self.result([expected]),
            inspection.option_values(text)), [])

    def test_entry_present_without_features_fails(self):
        entry = {"key": "devcontainer.features", "value": []}
        failures = inspection.inspect_failures(None, self.result([entry]), self.result([]), [])
        self.assertTrue(any("although no Feature" in f for f in failures))
        self.assertTrue(any("different" in f for f in failures))

    def test_echoed_options_fail(self):
        text = self.config("customized")
        expected = inspection.expected_entry(text)
        echoed = dict(expected, value=[{"ref": r, "options": {"version": "0.165.0"}}
                                       for r in expected["value"]])
        failures = inspection.inspect_failures(expected, self.result([echoed]),
                                               self.result([echoed]),
                                               inspection.option_values(text))
        self.assertTrue(any("echoes option value" in f for f in failures))

    def test_mcp_arguments_follow_the_request_branch(self):
        schema = json.loads((SPEC / "operation-request.schema.json").read_text())
        branch = next(b for b in schema["oneOf"]
                      if b["properties"]["operation"].get("const") == "inspect_workspace")
        self.assertEqual(inspection.mcp_arguments(branch, "/p", "r1"),
                         {"schemaVersion": "aibox.operation-request/v1",
                          "operation": "inspect_workspace", "requestId": "r1",
                          "projectRoot": "/p"})

    def test_mcp_tool_result_prefers_structured_content(self):
        self.assertEqual(inspection.mcp_tool_result(
            {"result": {"structuredContent": {"a": 1}, "content": []}}), {"a": 1})
        self.assertEqual(inspection.mcp_tool_result(
            {"result": {"content": [{"type": "text", "text": "{\"b\": 2}"}]}}), {"b": 2})


class TraceTest(unittest.TestCase):
    def table(self):
        return json.loads((SPEC / "ledger" / "v1-05-tool-selection.json").read_text())

    def test_every_census_row_traces_in_the_frozen_contract(self):
        failures = {row: trace.trace_failures(row, REAL, self.table(),
                                              REPO / "context" / "decisions")
                    for row in CENSUS}
        self.assertEqual({row: f for row, f in failures.items() if f}, {})

    def test_unmapped_row_fails(self):
        failures = trace.trace_failures("ADDON:none/none", REAL, self.table(),
                                        REPO / "context" / "decisions")
        self.assertEqual(failures, ["maps to 0 catalog places: []"])

    def test_removal_without_decision_record_fails(self):
        with tempfile.TemporaryDirectory() as empty:
            row = REAL.data["removed"][0]["row"]
            failures = trace.trace_failures(row, REAL, self.table(), Path(empty))
            self.assertTrue(any("has no DecisionRecord" in f for f in failures))


if __name__ == "__main__":
    unittest.main()

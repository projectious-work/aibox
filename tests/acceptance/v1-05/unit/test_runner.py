"""Runner pieces that need no engine: resuming from an earlier run and
cleanup after a stalled engine."""

import json
import tempfile
import unittest
from pathlib import Path

import runner
from actools.executor import Context
from unit.fixtures import REAL

RUN = {"runId": "now", "sourceCommit": "abc", "sourceDirty": False,
       "target": "debian-13-amd64", "scope": "sampled", "engineKind": "docker"}


class ResumeTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.earlier = Path(self.temp.name)

    def tearDown(self):
        self.temp.cleanup()

    def write(self, run: dict, records: list[dict]) -> None:
        (self.earlier / "run.json").write_text(json.dumps(run))
        (self.earlier / "results.jsonl").write_text(
            "".join(json.dumps(r) + "\n" for r in records) + '{"id": "cut')

    def test_passed_build_records_of_the_same_commit_are_taken(self):
        self.write({**RUN, "runId": "earlier"}, [
            {"id": "a", "build": "k1", "status": "passed"},
            {"id": "b", "build": "k2", "status": "failed"},
            {"id": "c", "build": "k3", "status": "not-run"},
            {"id": "trace", "build": None, "status": "passed"}])
        passed, resumed = runner.load_resume(str(self.earlier), RUN)
        self.assertEqual(list(passed), ["a"])
        self.assertEqual((resumed["runId"], resumed["passedRecords"]), ("earlier", 1))

    def test_another_commit_target_scope_or_engine_is_refused(self):
        for key, value in (("sourceCommit", "def"), ("target", "debian-13-arm64"),
                           ("scope", "full"), ("engineKind", "podman")):
            self.write({**RUN, key: value}, [])
            with self.assertRaisesRegex(RuntimeError, key):
                runner.load_resume(str(self.earlier), RUN)

    def test_a_dirty_checkout_is_refused(self):
        self.write({**RUN, "sourceDirty": True}, [])
        with self.assertRaisesRegex(RuntimeError, "clean checkout"):
            runner.load_resume(str(self.earlier), RUN)

    def test_a_directory_without_run_json_is_refused(self):
        with self.assertRaisesRegex(RuntimeError, "run.json"):
            runner.load_resume(str(self.earlier), RUN)


class CleanUpTest(unittest.TestCase):
    def test_a_stalled_engine_is_not_asked_to_clean_up(self):
        class Engine:
            def answers(self):
                raise AssertionError("must not ask a stalled engine")
        context = Context(REAL, "run", Path("/nonexistent"), Path("/nonexistent"),
                          docker=Engine())
        failures = runner.clean_up(context, None, "run", Path("/nonexistent"), True)
        self.assertEqual(len(failures), 1)
        self.assertIn("does not answer", failures[0])


if __name__ == "__main__":
    unittest.main()

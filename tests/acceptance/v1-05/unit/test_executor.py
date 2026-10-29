"""Executor against fake child executables (Meszaros Fakes for the Dev
Container CLI and Docker ports). Probes run on this machine through the fakes,
so the real probe script, parsing and evaluation are exercised."""

import json
import sys
import tempfile
import textwrap
import unittest
from pathlib import Path

from actools.catalog import Catalog
from actools.executor import Context, Executor
from actools.host import DevContainerCli, Docker
from actools.plan import make_plan
from unit.fixtures import REAL, feature, tool

TARGET = "debian-13-amd64"

FAKE_DEVCONTAINER = textwrap.dedent('''
    import json, subprocess, sys
    from pathlib import Path
    config = json.loads((Path(__file__).parent / "fake.json").read_text())
    args = sys.argv[1:]
    verb = args[0]
    folder = args[args.index("--workspace-folder") + 1]
    if verb == "build":
        print(config.get("buildOutput", ""), file=sys.stderr)
        sys.exit(config.get("buildExit", 1))
    if verb == "up":
        if config.get("upExit"):
            sys.exit(config["upExit"])
        print(json.dumps({"outcome": "success", "containerId": "fakecontainer",
                          "remoteUser": "aibox", "remoteWorkspaceFolder": folder}))
        sys.exit(0)
    if verb == "exec":
        index = args.index("--id-label") + 2
        sys.exit(subprocess.run(args[index:]).returncode)
    sys.exit(2)
''')

FAKE_DOCKER = textwrap.dedent('''
    #!/bin/sh
    case "$1" in
      exec) shift 4; exec "$@" ;;
      ps|images) exit 0 ;;
      rm|image) exit 0 ;;
    esac
    exit 2
''').lstrip()


def mini_catalog() -> Catalog:
    data = {k: v for k, v in REAL.data.items()
            if k in ("schemaVersion", "targets", "errors", "unsupportedTargets", "baseline")}
    present = tool("here", commands=["sh"], version="1.0", checks=[
        {"kind": "version", "argv": ["sh", "-c", "echo here version 1.0"],
         "regex": "(^|[^0-9.]){version}([^0-9]|$)"},
        {"kind": "path", "path": "/"}])
    broken = tool("broken", commands=["aibox-v1-05-missing-command"], version="1.0")
    absent = tool("gone", commands=["aibox-v1-05-missing-command"], version="1.0")
    data["features"] = [feature("t/here", "t", tools=[present]),
                        feature("t/broken", "t", tools=[broken]),
                        feature("t/gone", "g", requires=["t/here"], tools=[absent])]
    data["groups"] = []
    data["smoke"] = [{"id": "SMOKE-FAKE", "features": ["t/here"],
                      "fixture": "fixtures/v1-05-smoke/hello.dot",
                      "argv": ["sh", "-c", "test -f hello.dot && printf '%%PDF-1' > out.pdf"],
                      "produces": "out.pdf", "producesMagicHex": "25504446"}]
    return Catalog(data)


class ExecutorTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        root = Path(self.temp.name)
        self.fake_dir = root / "fakes"
        self.fake_dir.mkdir()
        (self.fake_dir / "devcontainer.py").write_text(FAKE_DEVCONTAINER)
        docker = self.fake_dir / "docker"
        docker.write_text(FAKE_DOCKER)
        docker.chmod(0o755)
        self.configure({})
        self.output = root / "out"
        self.output.mkdir()
        self.work = root / "work"
        self.work.mkdir()
        cli = DevContainerCli(Path(sys.executable), self.fake_dir / "devcontainer.py", {},
                              self.fake_dir)
        self.catalog = mini_catalog()
        self.context = Context(self.catalog, "unit", self.output, self.work,
                               docker=Docker(str(docker)), cli=cli, progress=lambda _: None)

    def tearDown(self):
        self.temp.cleanup()

    def configure(self, config):
        (self.fake_dir / "fake.json").write_text(json.dumps(config))

    def run_items(self, *items, features=()):
        plan = make_plan(self.catalog, [TARGET], [], items=list(items), features=list(features))
        executor = Executor(self.context, plan)
        executor.run_builds()
        return {r["id"]: r for r in executor.results}, executor

    def test_presence_passes_for_remote_user_and_root(self):
        results, executor = self.run_items("present", features=["t/here"])
        record = results[f"AC-TOOLS/present t/here here {TARGET}"]
        self.assertEqual(record["status"], "passed", record)
        self.assertEqual(record["checks"], 6)
        self.assertEqual(executor.cleanup_failures, [])

    def test_failed_presence_reports_the_check_and_user(self):
        results, _ = self.run_items("present", features=["t/broken"])
        record = results[f"AC-TOOLS/present t/broken broken {TARGET}"]
        self.assertEqual(record["status"], "failed")
        self.assertEqual({d["user"] for d in record["failedChecks"]}, {"remote", "root"})
        self.assertIn("command-present aibox-v1-05-missing-command",
                      " ".join(record["failures"]))

    def test_up_failure_fails_every_test_of_the_build(self):
        self.configure({"upExit": 1})
        results, _ = self.run_items("present", features=["t/here"])
        self.assertTrue(all("devcontainer up failed" in r["failures"][0]
                            for r in results.values()))

    def test_expected_error_line_passes_and_absence_of_it_fails(self):
        line = ("#5 aibox-feature error AIBOX_FEATURE_DEPENDENCY_MISSING feature=t/gone "
                "missing t/here")
        self.configure({"buildOutput": line})
        results, _ = self.run_items("requires")
        self.assertEqual([r["status"] for r in results.values()], ["passed"])
        self.configure({"buildOutput": "some other failure"})
        results, _ = self.run_items("requires")
        self.assertEqual([r["status"] for r in results.values()], ["failed"])
        self.configure({"buildExit": 0})
        results, _ = self.run_items("requires")
        self.assertIn("build succeeded", list(results.values())[0]["failures"][0])

    def test_smoke_runs_in_workspace_with_copied_fixture_and_checks_magic(self):
        results, _ = self.run_items("smoke")
        self.assertEqual(results[f"AC-TOOLS/smoke SMOKE-FAKE {TARGET}"]["status"], "passed")

    def test_workspaces_are_removed_after_each_build(self):
        self.run_items("present")
        self.assertEqual(list(self.work.iterdir()), [])


if __name__ == "__main__":
    unittest.main()

"""Executor against fake child executables (Meszaros Fakes for the Dev
Container CLI and Docker ports). Probes run on this machine through the fakes,
so the real probe script, parsing and evaluation are exercised."""

import json
import os
import sys
import tempfile
import textwrap
import threading
import time
import unittest
from pathlib import Path
from unittest import mock

from actools.catalog import Catalog
from actools import host
from actools.executor import Context, Executor, read_journal
from actools.host import DevContainerCli, Docker, EngineWatchdog, execute, process_usage
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
        if config.get("upSleep"):
            import time
            time.sleep(config["upSleep"])
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
      ps|images|info) exit 0 ;;
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
        cli = DevContainerCli([sys.executable, str(self.fake_dir / "devcontainer.py")], {})
        self.catalog = mini_catalog()
        self.context = Context(self.catalog, "unit", self.output, self.work,
                               docker=Docker(str(docker)), cli=cli, progress=lambda _: None,
                               process_usage=lambda: None)

    def tearDown(self):
        host.reset_abort()
        self.temp.cleanup()

    def configure(self, config):
        (self.fake_dir / "fake.json").write_text(json.dumps(config))

    def run_items(self, *items, features=()):
        plan = make_plan(self.catalog, [TARGET], [], items=list(items), features=list(features))
        executor = Executor(self.context, plan)
        executor.run_builds()
        return {r["id"]: r for r in executor.results}, executor

    def test_parallel_slots_run_every_build_and_keep_plan_order(self):
        self.context.engines = [self.context.docker, Docker(self.context.docker.binary)]
        plan = make_plan(self.catalog, [TARGET], [], items=["present"])
        executor = Executor(self.context, plan)
        executor.run_builds()
        self.assertEqual([r["id"] for r in executor.results], [t.id for t in plan.tests])

    def test_builds_do_not_start_at_the_process_budget(self):
        self.context.process_usage = lambda: (960, 1200)
        results, executor = self.run_items("present", features=["t/here"])
        self.assertEqual(results, {})
        self.assertIn("960 of 1200 tasks", executor.stop_reason)

    def test_builds_start_under_the_process_budget(self):
        self.context.process_usage = lambda: (959, 1200)
        results, executor = self.run_items("present", features=["t/here"])
        self.assertEqual(len(results), 1)
        self.assertIsNone(executor.stop_reason)

    def test_a_container_that_cannot_be_removed_stops_the_run(self):
        docker = self.fake_dir / "docker"
        docker.write_text(FAKE_DOCKER.replace("ps|images|info) exit 0 ;;",
                                              "ps) echo stuck; exit 0 ;;\n  images|info) exit 0 ;;")
                          .replace("rm|image) exit 0 ;;", "rm) exit 125 ;;\n  image) exit 0 ;;"))
        plan = make_plan(self.catalog, [TARGET], [], items=["present"])
        executor = Executor(self.context, plan)
        executor.run_builds()
        self.assertEqual(len(executor.results), 1)
        self.assertLess(len(executor.results), len({t.build for t in plan.tests}))
        self.assertEqual(executor.cleanup_failures, ["remove container stuck"])
        self.assertIn("stuck could not be removed", executor.stop_reason)

    def test_podman_slots_keep_the_runtime_directory_of_the_engine(self):
        podman = self.fake_dir / "podman"
        podman.write_text("#!/bin/sh\necho /run/state/podman/podman.sock\n")
        podman.chmod(0o755)
        with mock.patch.dict(os.environ):
            os.environ.pop("XDG_RUNTIME_DIR", None)
            environment = Docker(str(podman)).for_worker(1, self.work).cli_environment
            self.assertEqual(environment["XDG_RUNTIME_DIR"], "/run/state")
            os.environ["XDG_RUNTIME_DIR"] = "/run/user/1000"
            environment = Docker(str(podman)).for_worker(2, self.work).cli_environment
            self.assertEqual(environment["XDG_RUNTIME_DIR"], "/run/user/1000")
        self.assertTrue(environment["TMPDIR"].endswith("worker-2/tmp"))

    def test_docker_slots_need_no_runtime_directory(self):
        engine = self.context.docker.for_worker(1, self.work)
        self.assertNotIn("XDG_RUNTIME_DIR", engine.cli_environment)

    def test_process_usage_reads_the_cgroup_limit(self):
        cgroup = Path(self.temp.name)
        self.assertIsNone(process_usage(cgroup))
        (cgroup / "pids.current").write_text("10\n")
        (cgroup / "pids.max").write_text("max\n")
        self.assertIsNone(process_usage(cgroup))
        (cgroup / "pids.max").write_text("1200\n")
        self.assertEqual(process_usage(cgroup), (10, 1200))

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

    # --- stalls, limits and resume ------------------------------------------
    def test_a_build_over_its_time_limit_fails_in_time(self):
        self.configure({"upSleep": 30})
        self.context.build_timeout = 1
        started = time.monotonic()
        results, _ = self.run_items("present", features=["t/here"])
        self.assertLess(time.monotonic() - started, 15)
        self.assertIn("devcontainer up failed (exit 124)",
                      results[f"AC-TOOLS/present t/here here {TARGET}"]["failures"][0])

    def test_an_abort_kills_builds_in_progress_and_marks_them_not_run(self):
        self.configure({"upSleep": 60})
        plan = make_plan(self.catalog, [TARGET], [], items=["present"], features=["t/here"])
        executor = Executor(self.context, plan)
        threading.Timer(1.0, executor.abort, ["engine stalled"]).start()
        started = time.monotonic()
        executor.run_builds()
        self.assertLess(time.monotonic() - started, 20)
        self.assertEqual([r["status"] for r in executor.results], ["not-run"])
        self.assertIn("engine stalled", executor.results[0]["failures"][0])
        self.assertEqual(executor.stop_reason, "engine stalled")

    def test_results_are_journaled_as_they_come(self):
        _, executor = self.run_items("present", features=["t/here"])
        with executor.journal.open("a") as journal:
            journal.write('{"id": "cut off')
        self.assertEqual(read_journal(executor.journal), executor.results)

    def test_a_build_that_passed_in_the_earlier_run_is_not_rebuilt(self):
        _, first = self.run_items("present", features=["t/here"])
        self.context.carried = {r["id"]: r for r in first.results}
        self.context.carried_from = "earlier"
        self.configure({"upExit": 1})  # a rebuild would fail
        results, _ = self.run_items("present", features=["t/here"])
        record = results[f"AC-TOOLS/present t/here here {TARGET}"]
        self.assertEqual((record["status"], record["carriedFrom"]), ("passed", "earlier"))

    def test_a_build_with_a_failure_in_the_earlier_run_is_rebuilt(self):
        self.configure({"upExit": 1})
        _, first = self.run_items("present", features=["t/here"])
        self.context.carried = {r["id"]: r for r in first.results}
        self.context.carried_from = "earlier"
        self.configure({})
        results, _ = self.run_items("present", features=["t/here"])
        record = results[f"AC-TOOLS/present t/here here {TARGET}"]
        self.assertEqual(record["status"], "passed")
        self.assertNotIn("carriedFrom", record)


class ExecuteTest(unittest.TestCase):
    def tearDown(self):
        host.reset_abort()

    def test_a_timeout_kills_the_whole_process_group(self):
        # The grandchild keeps stdout open; killing only the child would hang.
        started = time.monotonic()
        result = execute(["sh", "-c", "sleep 60 & sleep 60"], timeout=1)
        self.assertEqual(result.exit_code, 124)
        self.assertLess(time.monotonic() - started, 15)

    def test_abort_running_kills_commands_of_other_threads(self):
        outcome = {}
        worker = threading.Thread(target=lambda: outcome.update(
            result=execute(["sh", "-c", "sleep 60 & sleep 60"], timeout=120)))
        worker.start()
        time.sleep(0.5)
        self.assertEqual(host.abort_running("engine stalled"), 1)
        worker.join(20)
        self.assertFalse(worker.is_alive())
        self.assertIn("aborted: engine stalled", outcome["result"].stderr)


class WatchdogTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name)

    def tearDown(self):
        self.temp.cleanup()

    def engine(self, body: str) -> Docker:
        fake = self.root / "docker"
        fake.write_text(f"#!/bin/sh\n{body}\n")
        fake.chmod(0o755)
        return Docker(str(fake))

    def test_three_missed_answers_report_a_stall_with_host_facts(self):
        stalls = []
        log = self.root / "engine-stall.log"
        watchdog = EngineWatchdog(self.engine("sleep 5"), stalls.append, log,
                                  interval=0.05, timeout=0.2, misses=3).start()
        deadline = time.monotonic() + 20
        while not stalls and time.monotonic() < deadline:
            time.sleep(0.05)
        watchdog.stop()
        self.assertEqual(len(stalls), 1)
        self.assertIn("did not answer", stalls[0])
        self.assertIn("miss 3 of 3", log.read_text())
        self.assertIn("$ ", log.read_text())  # host facts follow the misses

    def test_an_answering_engine_is_left_alone(self):
        stalls = []
        watchdog = EngineWatchdog(self.engine("exit 0"), stalls.append, self.root / "log",
                                  interval=0.05, timeout=1, misses=2).start()
        time.sleep(0.5)
        watchdog.stop()
        self.assertEqual(stalls, [])


if __name__ == "__main__":
    unittest.main()

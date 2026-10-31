"""Run a plan on the host: build each workspace once, evaluate its tests,
tear down its scoped resources and record results."""

from __future__ import annotations

import json
import queue
import shutil
import threading
from collections import OrderedDict
from concurrent.futures import ThreadPoolExecutor
from dataclasses import dataclass, field
from pathlib import Path
from typing import Callable, Mapping

from . import checks, inspection, jsonc, metadata, trace
from .catalog import REPO, SPEC, Catalog
from .host import Completed, DevContainerCli, Docker, Registry, execute
from .plan import Build, Plan, Test

OUTPUT_LIMIT = 4000


def trimmed(text: str) -> str:
    return text if len(text) <= OUTPUT_LIMIT else text[:OUTPUT_LIMIT] + "\n[truncated]"


def result_record(test: Test, failures: list[str], **extra) -> dict:
    record = {"id": test.id, "item": f"AC-TOOLS/{test.item}", "feature": test.feature,
              "tool": test.tool, "target": test.target, "qualifiers": dict(test.qualifiers),
              "build": test.build, "status": "failed" if failures else "passed",
              "failures": failures}
    record.update(extra)
    return record


@dataclass
class Context:
    catalog: Catalog
    run_id: str
    output: Path
    work: Path
    docker: Docker | None = None
    cli: DevContainerCli | None = None
    registry: Registry | None = None
    aibox: Path | None = None
    progress: Callable[[str], None] = print
    engines: list[Docker] = field(default_factory=list)  # one per parallel build slot


class Executor:
    def __init__(self, context: Context, plan: Plan):
        self.ctx = context
        self.plan = plan
        self.results: list[dict] = []
        self.cleanup_failures: list[str] = []
        self.slot = threading.local()
        self.lock = threading.Lock()
        (context.output / "builds").mkdir(exist_ok=True)
        (context.output / "tests").mkdir(exist_ok=True)

    @property
    def docker(self) -> Docker | None:
        return getattr(self.slot, "engine", self.ctx.docker)

    # --- offline items ---------------------------------------------------
    def run_offline(self) -> None:
        table = json.loads((SPEC / "ledger" / "v1-05-tool-selection.json").read_text())
        decisions = REPO / "context" / "decisions"
        for test in self.plan.tests:
            if test.item == "trace":
                failures = trace.trace_failures(test.payload["row"], self.ctx.catalog, table,
                                                decisions)
                self.record(result_record(test, failures))
            elif test.item == "inspect":
                self.record(self.run_inspect(test))
            elif test.item == "guards":
                self.record(self.run_guards(test))

    def run_guards(self, test: Test) -> dict:
        if self.ctx.cli is None:
            return result_record(test, ["no Node available"])
        log = self.ctx.output / "tests" / "guards.log"
        files = sorted(str(path) for path in (REPO / "features" / "test").glob("*.test.mjs"))
        completed = execute(["node", "--test", "--test-reporter=tap", *files], log=log,
                            env=self.ctx.cli.environment(), timeout=1800)
        failures = [f"component tests exited {completed.exit_code}"] if completed.exit_code else []
        failures += checks.guard_coverage_failures(completed.stdout, self.ctx.catalog)
        return result_record(test, failures, log="tests/guards.log")

    def run_inspect(self, test: Test) -> dict:
        example = test.payload["example"]
        project = SPEC / "examples" / example
        if self.ctx.aibox is None:
            return result_record(test, ["no aibox binary available"])
        text = (project / ".devcontainer" / "devcontainer.json").read_text()
        logs = self.ctx.output / "tests"
        try:
            cli_result = inspection.run_cli(self.ctx.aibox, project, logs / f"inspect-{example}-cli.log")
            mcp_result = inspection.run_mcp(self.ctx.aibox, project, f"ac-tools-inspect-{example}",
                                            logs / f"inspect-{example}-mcp.log")
        except Exception as exc:  # the adapter failure is the test result
            return result_record(test, [f"inspection failed: {exc}"])
        failures = inspection.inspect_failures(inspection.expected_entry(text), cli_result,
                                               mcp_result, inspection.option_values(text))
        return result_record(test, failures, logs=[f"tests/inspect-{example}-cli.log",
                                                   f"tests/inspect-{example}-mcp.log"])

    # --- container items ---------------------------------------------------
    def run_builds(self) -> None:
        by_build: "OrderedDict[str, list[Test]]" = OrderedDict()
        for test in self.plan.tests:
            if test.build:
                by_build.setdefault(test.build, []).append(test)
        total = len(by_build)
        engines = self.ctx.engines or [self.docker]
        slots: "queue.SimpleQueue[Docker | None]" = queue.SimpleQueue()
        for engine in engines:
            slots.put(engine)

        def run_one(index: int, key: str, tests: list[Test]) -> None:
            build = self.plan.builds[key]
            engine = slots.get()
            self.slot.engine = engine
            self.ctx.progress(f"[{index}/{total}] build {key} "
                              f"({'expect failure' if build.expect_failure else 'up'}; "
                              f"{build.fixture}; {len(build.selection)} Features; "
                              f"{len(tests)} tests)")
            try:
                records = self.run_build(index, build, tests)
            except Exception as exc:
                records = [result_record(t, [f"harness error: {exc}"]) for t in tests]
            finally:
                slots.put(engine)
            for record in records:
                self.record(record)

        pool = ThreadPoolExecutor(max_workers=len(engines))
        try:
            futures = [pool.submit(run_one, index, key, tests)
                       for index, (key, tests) in enumerate(by_build.items(), start=1)]
            for future in futures:
                future.result()
        finally:
            # An interrupted run waits for builds in progress, then cleans up.
            pool.shutdown(wait=True, cancel_futures=True)
        order = {test.id: position for position, test in enumerate(self.plan.tests)}
        self.results.sort(key=lambda record: order[record["id"]])

    def record(self, record: dict) -> None:
        with self.lock:
            self.results.append(record)
            self.ctx.progress(f"  {record['status'].upper():6} {record['id']}")

    def workspace(self, index: int, build: Build) -> Path:
        name = f"aibox-v1-05-{self.ctx.run_id}-{index:04d}"
        workspace = self.ctx.work / name
        shutil.copytree(SPEC / build.fixture_path, workspace, symlinks=True)
        config_path = workspace / ".devcontainer" / "devcontainer.json"
        config = jsonc.loads(config_path.read_text())
        config["features"] = self.ctx.catalog.devcontainer_features(build.selection_dict())
        config_path.write_text(json.dumps(config, indent=2) + "\n")
        for copy in build.copies:
            source = SPEC / copy
            destination = workspace / source.name
            if source.is_dir():
                shutil.copytree(source, destination, symlinks=True)
            else:
                shutil.copy2(source, destination)
        return workspace

    def run_build(self, index: int, build: Build, tests: list[Test]) -> list[dict]:
        assert self.docker and self.ctx.cli
        log = self.ctx.output / "builds" / f"{build.key}.log"
        log.write_text(json.dumps({"key": build.key, "fixture": build.fixture,
                                   "selection": build.selection_dict(),
                                   "copies": list(build.copies)}, indent=2) + "\n\n")
        workspace = self.workspace(index, build)
        common = ["--docker-path", self.docker.binary, "--workspace-folder", str(workspace)]
        try:
            if build.expect_failure:
                records = self.run_failing_build(index, build, tests, common, log)
            else:
                records = self.run_up_build(index, build, tests, common, log, workspace)
            if self.docker.kind == "podman":
                cached = checks.cached_feature_steps(log.read_text(errors="replace"))
                if cached:
                    for record in records:
                        record["status"] = "failed"
                        record["failures"].append(
                            f"{cached} Feature steps came from the Podman layer cache")
            return records
        finally:
            # Parallel builds share the engine's store, where removing an image
            # could remove layers another build has just matched; the runner
            # removes the run's images at the end.
            if len(self.ctx.engines) <= 1:
                self.cleanup_failures += self.docker.remove_images(
                    f"vsc-{workspace.name}*", log)
                self.cleanup_failures += self.docker.remove_images(
                    f"aibox-v1-05-{self.ctx.run_id}-{index:04d}*", log)
            shutil.rmtree(workspace, ignore_errors=True)

    def run_failing_build(self, index: int, build: Build, tests: list[Test],
                          common: list[str], log: Path) -> list[dict]:
        image = f"aibox-v1-05-{self.ctx.run_id}-{index:04d}"
        completed = self.cli(["build", *common, "--image-name", image,
                              "--log-format", "text"], log)
        records = []
        for test in tests:
            error = test.error or {}
            failures = []
            if completed.exit_code == 0:
                failures.append("build succeeded; expected failure")
            elif not checks.error_line_matches(completed.combined, error["code"],
                                               error["features"], error.get("mustName")):
                names = f" naming {error['mustName']}" if error.get("mustName") else ""
                failures.append(f"no line contains aibox-feature error {error['code']} "
                                f"feature=<one of {error['features']}>{names}")
            records.append(result_record(test, failures, log=f"builds/{build.key}.log",
                                         buildExitCode=completed.exit_code))
        return records

    def run_up_build(self, index: int, build: Build, tests: list[Test], common: list[str],
                     log: Path, workspace: Path) -> list[dict]:
        assert self.docker
        label = f"aibox.v1-05.run={self.ctx.run_id}-{index:04d}"
        before = metadata.snapshot(workspace)
        up = self.cli(["up", *common, "--id-label", label, "--log-format", "text"], log)
        try:
            outcome = last_json_line(up.stdout)
            if up.exit_code or outcome.get("outcome") != "success":
                return [result_record(t, [f"devcontainer up failed (exit {up.exit_code})"],
                                      log=f"builds/{build.key}.log") for t in tests]
            container_id = outcome["containerId"]
            workspace_changes = metadata.snapshot_changes(
                before, metadata.snapshot(workspace), exempt=(metadata.CLI_LOCKFILE,))
            outcomes = self.probe_all(tests, common, label, container_id, log)
            records = []
            for test in tests:
                if test.item == "smoke":
                    records.append(self.run_smoke(test, common, label, outcome, workspace, log))
                elif test.item == "side-effects":
                    records.append(self.run_side_effects(test, workspace_changes, build))
                else:
                    records.append(self.evaluate(test, outcomes, build))
            return records
        finally:
            for container_id in sorted(self.containers(label)):
                if self.docker.run("rm", "-f", "-v", container_id, log=log).exit_code:
                    self.cleanup_failures.append(f"remove container {container_id}")

    def containers(self, label: str) -> set[str]:
        assert self.docker
        listing = self.docker.run("ps", "-aq", "--filter", f"label={label}")
        if listing.exit_code:
            self.cleanup_failures.append(f"list containers {label}")
        return set(listing.stdout.split())

    def cli(self, args: list[str], log: Path, timeout: float = 7200) -> Completed:
        assert self.ctx.cli
        env = {**self.ctx.cli.environment(), **(self.docker.cli_environment if self.docker else {})}
        return execute(self.ctx.cli.argv(*args), log=log, env=env, timeout=timeout)

    def probe_all(self, tests: list[Test], common: list[str], label: str, container_id: str,
                  log: Path) -> dict[str, dict]:
        assert self.docker
        probes = list(OrderedDict.fromkeys(e.probe for t in tests for e in t.expectations))
        if not probes:
            return {"remote": {}, "root": {}}
        script = checks.probe_script(probes)
        remote = self.cli(["exec", *common, "--id-label", label, "sh", "-c", script], log)
        root = self.docker.run("exec", "--user", "root", container_id, "sh", "-c", script,
                                   log=log, timeout=7200)
        return {"remote": self.parse(remote, probes), "root": self.parse(root, probes)}

    @staticmethod
    def parse(completed: Completed, probes: list[checks.Probe]) -> dict:
        try:
            return checks.parse_probe_output(completed.stdout, probes)
        except Exception as exc:
            return {"error": f"probe run failed (exit {completed.exit_code}): {exc}"}

    def evaluate(self, test: Test, outcomes: Mapping[str, Mapping], build: Build) -> dict:
        failures, details = [], []
        for user in ("remote", "root"):
            results = outcomes[user]
            if "error" in results:
                failures.append(f"{user}: {results['error']}")
                continue
            for expectation in test.expectations:
                outcome = results[expectation.probe]
                if not checks.passes(expectation, outcome):
                    failures.append(f"{user}: {expectation.describe()} failed")
                    details.append({"user": user, "rule": expectation.rule,
                                    "subject": expectation.subject, "regex": expectation.regex,
                                    "expected": expectation.expected,
                                    "argv": list(expectation.probe.argv),
                                    "exitCode": outcome.exit_code,
                                    "output": trimmed(outcome.output)})
        return result_record(test, failures, log=f"builds/{build.key}.log",
                             checks=len(test.expectations) * 2, failedChecks=details)

    def run_smoke(self, test: Test, common: list[str], label: str, outcome: Mapping,
                  workspace: Path, log: Path) -> dict:
        folder = outcome.get("remoteWorkspaceFolder")
        argv = ["sh", "-c", 'cd "$1" && shift && exec "$@"', "sh", str(folder),
                *test.payload["argv"]]
        completed = self.cli(["exec", *common, "--id-label", label, *argv], log)
        failures = []
        if not folder:
            failures.append("devcontainer up reported no remoteWorkspaceFolder")
        if completed.exit_code:
            failures.append(f"smoke command exited {completed.exit_code}")
        produced = workspace / test.payload["produces"]
        if not produced.is_file():
            failures.append(f"{test.payload['produces']} was not created")
        elif test.payload.get("magicHex"):
            magic = bytes.fromhex(test.payload["magicHex"])
            if produced.read_bytes()[:len(magic)] != magic:
                failures.append(f"{test.payload['produces']} does not start with "
                                f"{test.payload['magicHex']}")
        return result_record(test, failures, log=f"builds/{test.build}.log",
                             output=trimmed(completed.combined) if failures else None)

    def run_side_effects(self, test: Test, workspace_changes: list[str], build: Build) -> dict:
        failures = [f"workspace changed by the build: {change}" for change in workspace_changes]
        try:
            assert self.ctx.registry
            layer = self.ctx.registry.feature_layer(test.feature or "",
                                                    test.payload["featureVersion"])
            published = metadata.feature_json_from_layer(layer)
            failures += metadata.side_effect_violations(published, test.payload["containerEnv"])
        except Exception as exc:
            failures.append(f"cannot read published metadata: {exc}")
        return result_record(test, failures, log=f"builds/{build.key}.log")


def last_json_line(text: str) -> dict:
    for line in reversed(text.splitlines()):
        line = line.strip()
        if line.startswith("{"):
            try:
                return json.loads(line)
            except json.JSONDecodeError:
                continue
    return {}

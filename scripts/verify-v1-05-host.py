#!/usr/bin/env python3
"""Check the selected Hugo Feature on one host architecture with a pinned CLI."""

from __future__ import annotations

import argparse
import importlib.util
import json
import os
from pathlib import Path
import platform
import shutil
import subprocess
import sys
import tempfile
from datetime import datetime, timezone


REPO = Path(__file__).resolve().parent.parent
spec = importlib.util.spec_from_file_location("v104_host", REPO / "scripts/verify-v1-04-host.py")
assert spec and spec.loader
v104 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(v104)


class Gate(v104.Gate):
    def __init__(self, output: Path, work: Path, source_commit: str):
        super().__init__(output, work, source_commit, False)
        self.run_id = datetime.now(timezone.utc).strftime("%Y%m%dt%H%M%S") + f"-{os.getpid()}"
        self.labels: set[str] = set()
        self.images: set[str] = set()

    def step(self, name: str, argv: list[str], *, cwd: Path | None = None) -> str:
        log = self.output / f"{name}.log"
        print(f"V1-05 host gate: {name} started; log: {log}", flush=True)
        with log.open("w") as target:
            target.write("$ " + " ".join(argv) + "\n")
            target.flush()
            try:
                completed = subprocess.run(argv, cwd=cwd, text=True, stdout=target,
                                           stderr=subprocess.STDOUT, check=False)
            except KeyboardInterrupt as exc:
                self.steps.append({"name": name, "log": log.name, "exitCode": 130})
                raise RuntimeError(f"{name} interrupted; see {log}") from exc
        self.steps.append({"name": name, "log": log.name, "exitCode": completed.returncode})
        if completed.returncode:
            raise RuntimeError(f"{name} failed ({completed.returncode}); see {log}")
        print(f"V1-05 host gate: {name} passed", flush=True)
        return log.read_text()

    def run(self) -> None:
        self.node, self.cli, self.runtime = v104.prepare_cli(self.work, self.output)
        self.engine = v104.engine_path()
        self.step("engine-info", [self.engine, "info"])
        fixture = json.loads((REPO / "spec/v1/fixtures/v1-05-native-features.json").read_text())
        expected_ref = next(iter(fixture["selected"]))
        for name, expected_present in (("customized", True), ("minimal", False)):
            project = self.work / name
            shutil.copytree(REPO / "spec/v1/examples" / name, project)
            config = json.loads("\n".join(line for line in
                                (project / ".devcontainer/devcontainer.json").read_text().splitlines()
                                if not line.lstrip().startswith("//")))
            selected = config.get("features", {})
            if (expected_ref in selected) != expected_present:
                raise RuntimeError(f"{name}: selected Feature does not match reviewed fixture")
            common = ["--docker-path", self.engine, "--workspace-folder", str(project)]
            self.cli_step(f"read-configuration-{name}",
                          ["read-configuration", *common, "--include-merged-configuration"])
            image = f"aibox-v1-05-{self.run_id}-{name}"
            label = f"aibox.v1-05.run={self.run_id}-{name}"
            self.images.add(image)
            self.labels.add(label)
            self.cli_step(f"build-{name}", ["build", *common, "--image-name", image,
                                             "--label", f"aibox.v1-05.run={self.run_id}"])
            self.cli_step(f"up-{name}", ["up", *common, "--id-label", label])
            container = self.find_container(label)
            payload = json.loads(self.docker(["inspect", container]))[0]
            if (payload.get("Config", {}).get("Labels") or {}).get("aibox.v1-05.run") != f"{self.run_id}-{name}":
                raise RuntimeError(f"{name}: container lost its run label")
            if expected_present:
                command = "set -e; hugo version | grep -E 'v0[.]165[.]0.*extended'"
            else:
                command = "set -e; ! command -v hugo >/dev/null 2>&1"
            self.cli_step(f"exec-{name}", ["exec", *common, "--id-label", label,
                                           "bash", "-lc", command])

    def cleanup(self) -> None:
        if self.engine:
            for label in sorted(self.labels):
                try:
                    self.container_ids.update(self.docker(["ps", "-aq", "--filter", f"label={label}"]).splitlines())
                except Exception as exc:
                    self.cleanup_failures.append(f"list containers for {label}: {exc}")
            for container in sorted(self.container_ids):
                try:
                    payload = json.loads(self.docker(["inspect", container]))[0]
                except subprocess.CalledProcessError:
                    continue
                except Exception as exc:
                    self.cleanup_failures.append(f"inspect container {container}: {exc}")
                    continue
                label = (payload.get("Config", {}).get("Labels") or {}).get("aibox.v1-05.run")
                if label not in self.labels:
                    self.cleanup_failures.append(f"refused container without exact run label: {container}")
                    continue
                for mount in payload.get("Mounts", []):
                    if mount.get("Type") == "volume" and mount.get("Destination") == "/home/aibox" and \
                       mount.get("Name", "").startswith("aibox-home-"):
                        self.home_volumes.add(mount["Name"])
                try:
                    self.docker(["rm", "-f", container])
                except Exception as exc:
                    self.cleanup_failures.append(f"remove container {container}: {exc}")
            for volume in sorted(self.home_volumes):
                try:
                    self.docker(["volume", "rm", volume])
                except Exception as exc:
                    self.cleanup_failures.append(f"remove volume {volume}: {exc}")
            for image in sorted(self.images):
                try:
                    self.docker(["image", "inspect", image])
                except subprocess.CalledProcessError:
                    continue
                try:
                    self.docker(["image", "rm", image])
                except Exception as exc:
                    self.cleanup_failures.append(f"remove image {image}: {exc}")
        shutil.rmtree(self.work)

    def write_evidence(self) -> None:
        artifacts = [{"path": file.name, "sha256": "sha256:" + v104.digest(file)}
                     for file in sorted(self.output.iterdir()) if file.is_file() and file.name != "evidence.json"]
        evidence = {
            "phase": "V1-05", "verification": "host", "sourceCommit": self.source_commit,
            "startedAt": self.started_at, "finishedAt": v104.now(), "runId": self.run_id,
            "host": {"os": platform.system().lower(), "architecture": platform.machine().lower()},
            "status": "passed" if not self.failures and not self.cleanup_failures else "failed",
            "engine": self.engine, "runtime": self.runtime, "steps": self.steps,
            "failures": self.failures, "cleanupFailures": self.cleanup_failures,
            "cleanup": {"exactRunLabelsOnly": True, "runOwnedImagesOnly": True,
                        "observedHomeVolumesOnly": True}, "artifacts": artifacts,
        }
        (self.output / "evidence.json").write_text(json.dumps(evidence, indent=2) + "\n")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", help="New absolute evidence directory")
    args = parser.parse_args()
    try:
        if v104.run(["git", "status", "--porcelain", "--untracked-files=all"], cwd=REPO):
            raise RuntimeError("V1-05 source worktree must be clean before evidence capture")
        source_commit = v104.run(["git", "rev-parse", "HEAD"], cwd=REPO)
        output = (v104.prepare_output(args.output) if args.output else
                  Path(tempfile.mkdtemp(prefix="aibox-v1-05-host-evidence-")).resolve())
    except Exception as exc:
        print(exc, file=sys.stderr)
        return 2
    work = Path(tempfile.mkdtemp(prefix="aibox-v1-05-host-run-")).resolve()
    gate = Gate(output, work, source_commit)
    try:
        gate.run()
    except BaseException as exc:
        gate.failures.append(str(exc) or type(exc).__name__)
        print(f"V1-05 host gate failed: {gate.failures[-1]}", file=sys.stderr)
    try:
        gate.cleanup()
    except BaseException as exc:
        gate.cleanup_failures.append(f"temporary workdir: {exc}")
        print(f"V1-05 cleanup failed: {exc}", file=sys.stderr)
    gate.write_evidence()
    print(f"V1-05 evidence: {output}")
    return 0 if not gate.failures and not gate.cleanup_failures else 1


if __name__ == "__main__":
    raise SystemExit(main())

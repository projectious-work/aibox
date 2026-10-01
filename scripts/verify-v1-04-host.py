#!/usr/bin/env python3
"""V1-04 host gate. Python and Docker/Podman are the only host tools required."""

from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import platform
import shutil
import subprocess
import sys
import tarfile
import tempfile
from datetime import datetime, timezone
from urllib.request import urlopen


CLI_VERSION = "0.89.0"
NODE_VERSION = "20.19.5"
CLI_SHA256 = "49c7d71d40058f89e1fd8b019a193ed4215b7fc773c0f6273f7032a46cd33f4b"
NODE_SHA256 = {
    ("darwin", "arm64"): "cfed7503d8d99fbcf2f52e408ec52f616058eb0867b34dbc3437259993ef5cba",
    ("darwin", "x64"): "f9cff058f2766d4d0631dc69b5f7f27664b3a42ff186e25ac7e1ac269af7e696",
    ("linux", "arm64"): "a08b513de673853ca16395ec461a104a99bf0e941ebb7baedb98b30cd221d8cc",
    ("linux", "x64"): "4eba5fbe1fb10753bc06e42f001a91c5cec16798b7764a3e9257adc59af47fe1",
}
NAMES = ("minimal", "curated", "custom-user", "bind-home")
REPO = Path(__file__).resolve().parent.parent
ARCHIVE = REPO / "tools/vendor/devcontainer-cli-0.89.0.tgz"


def now() -> str:
    return datetime.now(timezone.utc).isoformat(timespec="seconds").replace("+00:00", "Z")


def digest(path: Path) -> str:
    h = hashlib.sha256()
    with path.open("rb") as source:
        for chunk in iter(lambda: source.read(1024 * 1024), b""):
            h.update(chunk)
    return h.hexdigest()


def run(argv: list[str], *, cwd: Path | None = None) -> str:
    return subprocess.check_output(argv, cwd=cwd, text=True, stderr=subprocess.STDOUT).strip()


def require_clean_checkout() -> str:
    if run(["git", "status", "--porcelain", "--untracked-files=all"], cwd=REPO):
        raise RuntimeError("V1-04 source worktree must be clean before evidence capture")
    return run(["git", "rev-parse", "HEAD"], cwd=REPO)


def prepare_output(value: str | None) -> Path:
    if value is None:
        return Path(tempfile.mkdtemp(prefix="aibox-v1-04-host-evidence-")).resolve()
    output = Path(value)
    if not output.is_absolute() or output.exists():
        raise RuntimeError("--output must name a new absolute directory")
    output.mkdir(parents=True, mode=0o700)
    output.chmod(0o700)
    return output


def download_verified(url: str, destination: Path, expected: str) -> None:
    h = hashlib.sha256()
    with urlopen(url, timeout=120) as response, destination.open("wb") as target:
        for chunk in iter(lambda: response.read(1024 * 1024), b""):
            h.update(chunk)
            target.write(chunk)
    if h.hexdigest() != expected:
        destination.unlink(missing_ok=True)
        raise RuntimeError(f"SHA256 mismatch for {url}")


def prepare_cli(work: Path, output: Path) -> tuple[Path, Path, dict[str, str]]:
    if digest(ARCHIVE) != CLI_SHA256:
        raise RuntimeError(f"Vendored CLI digest mismatch: {ARCHIVE}")
    system = platform.system().lower()
    machine = platform.machine().lower()
    arch = {"arm64": "arm64", "aarch64": "arm64", "x86_64": "x64", "amd64": "x64"}.get(machine)
    if (system, arch) not in NODE_SHA256:
        raise RuntimeError(f"Unsupported host platform: {system}/{machine}")
    filename = f"node-v{NODE_VERSION}-{system}-{arch}.tar.gz"
    url = f"https://nodejs.org/dist/v{NODE_VERSION}/{filename}"
    node_archive = work / filename
    download_verified(url, node_archive, NODE_SHA256[(system, arch)])
    node = work / "node"
    with tarfile.open(node_archive, "r:gz") as source:
        member = source.getmember(f"node-v{NODE_VERSION}-{system}-{arch}/bin/node")
        if not member.isfile():
            raise RuntimeError("Node archive does not contain a regular bin/node")
        stream = source.extractfile(member)
        if stream is None:
            raise RuntimeError("Could not read Node executable")
        with node.open("wb") as target:
            shutil.copyfileobj(stream, target)
    node.chmod(0o700)
    package = work / "cli/package"
    allowed = {
        "package/devcontainer.js",
        "package/dist/spec-node/devContainersSpecCLI.js",
        "package/package.json",
        "package/CHANGELOG.md",
        "package/README.md",
        "package/LICENSE.txt",
        "package/ThirdPartyNotices.txt",
        "package/scripts/updateUID.Dockerfile",
    }
    with tarfile.open(ARCHIVE, "r:gz") as source:
        found = set()
        for member in source:
            if member.isdir():
                continue
            if member.name not in allowed or not member.isfile():
                raise RuntimeError(f"Unexpected CLI archive entry: {member.name}")
            found.add(member.name)
            destination = work / "cli" / member.name
            destination.parent.mkdir(parents=True, exist_ok=True)
            stream = source.extractfile(member)
            if stream is None:
                raise RuntimeError(f"Could not read CLI archive entry: {member.name}")
            with destination.open("wb") as target:
                shutil.copyfileobj(stream, target)
        if found != allowed:
            raise RuntimeError(f"CLI archive entries differ: {sorted(allowed - found)}")
    metadata = json.loads((package / "package.json").read_text())
    if metadata.get("name") != "@devcontainers/cli" or metadata.get("version") != CLI_VERSION:
        raise RuntimeError("Vendored CLI package metadata mismatch")
    cli = package / "devcontainer.js"
    runtime = {
        "nodeVersion": run([str(node), "--version"]),
        "nodeArchive": url,
        "nodeArchiveSha256": NODE_SHA256[(system, arch)],
        "cliVersion": run([str(node), str(cli), "--version"]),
        "cliArchive": str(ARCHIVE.relative_to(REPO)),
        "cliArchiveSha256": CLI_SHA256,
    }
    if runtime["cliVersion"] != CLI_VERSION:
        raise RuntimeError("Pinned Dev Container CLI version mismatch")
    (output / "cli-runtime.json").write_text(json.dumps(runtime, indent=2) + "\n")
    return node, cli, runtime


def render_template(node: Path, work: Path, name: str, project: Path) -> None:
    project.mkdir()
    run([str(node), str(REPO / "templates/render.mjs"), "--template", name,
         "--output", str(project)])


def engine_path() -> str:
    for name in ("docker", "podman"):
        binary = shutil.which(name)
        if binary:
            try:
                run([binary, "info"])
                return binary
            except subprocess.CalledProcessError:
                continue
    raise RuntimeError("A working Docker or Podman CLI is required")


class Gate:
    def __init__(self, output: Path, work: Path, source_commit: str, preflight: bool):
        self.output, self.work, self.source_commit = output, work, source_commit
        self.preflight = preflight
        self.started_at = now()
        self.steps: list[dict[str, object]] = []
        self.failures: list[str] = []
        self.cleanup_failures: list[str] = []
        self.run_id = datetime.now(timezone.utc).strftime("%Y%m%dt%H%M%S") + f"-{os.getpid()}"
        self.engine: str | None = None
        self.node: Path | None = None
        self.cli: Path | None = None
        self.runtime: dict[str, str] = {}
        self.container_ids: set[str] = set()
        self.home_volumes: set[str] = set()
        self.image_refs: set[str] = set()

    def step(self, name: str, argv: list[str], *, cwd: Path | None = None) -> str:
        log = self.output / f"{name}.log"
        print(f"V1-04 host gate: {name} started; log: {log}", flush=True)
        with log.open("w") as target:
            target.write("$ " + " ".join(argv) + "\n")
            target.flush()
            try:
                completed = subprocess.run(argv, cwd=cwd, text=True, stdout=target,
                                           stderr=subprocess.STDOUT, check=False)
            except KeyboardInterrupt as exc:
                self.steps.append({"name": name, "command": argv, "log": log.name,
                                   "exitCode": 130})
                raise RuntimeError(f"{name} interrupted; see {log}") from exc
        self.steps.append({"name": name, "command": argv, "log": log.name,
                           "exitCode": completed.returncode})
        if completed.returncode:
            raise RuntimeError(f"{name} failed ({completed.returncode}); see {log}")
        print(f"V1-04 host gate: {name} passed", flush=True)
        return log.read_text()

    def cli_step(self, name: str, args: list[str]) -> str:
        assert self.node and self.cli
        return self.step(name, [str(self.node), str(self.cli), *args])

    def docker(self, args: list[str]) -> str:
        assert self.engine
        return run([self.engine, *args])

    def find_container(self, label: str) -> str:
        ids = self.docker(["ps", "-aq", "--filter", f"label={label}"]).splitlines()
        if len(ids) != 1:
            raise RuntimeError(f"Expected one run-labeled container for {label}, found {ids}")
        self.container_ids.add(ids[0])
        return ids[0]

    def inspect_home(self, name: str, project: Path, container_id: str, stage: str) -> str:
        payload = json.loads(self.docker(["inspect", container_id]))[0]
        expected_user = "dev" if name == "custom-user" else "aibox"
        labels = payload.get("Config", {}).get("Labels") or {}
        if labels.get("aibox.v1-04.run") != f"{self.run_id}-{name}":
            raise RuntimeError(f"{name}: missing exact run label")
        user = payload.get("Config", {}).get("User")
        if user and user != expected_user:
            raise RuntimeError(f"{name}: container user {user!r}, expected {expected_user!r}")
        target = f"/home/{expected_user}"
        mounts = [m for m in payload.get("Mounts", []) if m.get("Destination") == target]
        if len(mounts) != 1:
            raise RuntimeError(f"{name}: expected one home mount at {target}")
        mount = mounts[0]
        if name == "bind-home":
            expected = (project / ".aibox-home").resolve()
            if mount.get("Type") != "bind" or Path(mount.get("Source", "")).resolve() != expected:
                raise RuntimeError(f"{name}: home bind source mismatch")
            identity = str(expected)
        else:
            identity = mount.get("Name", "")
            if mount.get("Type") != "volume" or not identity.startswith("aibox-home-"):
                raise RuntimeError(f"{name}: expected project-scoped named home volume")
            self.home_volumes.add(identity)
        image = payload.get("Config", {}).get("Image")
        if image:
            self.image_refs.add(image)
        summary = {"containerId": payload["Id"], "containerUser": user or expected_user,
                   "home": {"type": mount["Type"], "identity": identity, "target": target},
                   "imageReference": image}
        (self.output / f"home-mount-{name}-{stage}.json").write_text(
            json.dumps(summary, indent=2) + "\n")
        return f"{mount['Type']}|{identity}"

    def run(self) -> None:
        self.node, self.cli, self.runtime = prepare_cli(self.work, self.output)
        stage = self.work / "template-docs"
        source = stage / "src"
        source.mkdir(parents=True)
        for name in ("minimal", "curated"):
            shutil.copytree(REPO / "templates/src" / name, source / name)
        self.cli_step("template-docs", ["templates", "generate-docs",
                                        "--project-folder", str(source)])
        for name in ("minimal", "curated"):
            if not (source / name / "README.md").is_file():
                raise RuntimeError(f"{name}: upstream Template metadata check failed")
        if self.preflight:
            for name in ("minimal", "curated"):
                render_template(self.node, self.work, name, self.work / name)
            return
        self.engine = engine_path()
        self.step("engine-info", [self.engine, "info"])
        for name in NAMES:
            project = self.work / name
            if name in ("minimal", "curated"):
                render_template(self.node, self.work, name, project)
            else:
                shutil.copytree(REPO / "spec/v1/examples" / name, project)
            if name == "bind-home":
                (project / ".aibox-home").mkdir(mode=0o700)
            common = ["--docker-path", self.engine, "--workspace-folder", str(project)]
            self.cli_step(f"read-configuration-{name}",
                          ["read-configuration", *common, "--include-merged-configuration"])
            image_tag = f"aibox-v1-04-{self.run_id}-{name}"
            self.cli_step(f"build-{name}",
                          ["build", *common, "--image-name", image_tag,
                           "--label", f"aibox.v1-04.run={self.run_id}"])
            label = f"aibox.v1-04.run={self.run_id}-{name}"
            self.cli_step(f"up-{name}", ["up", *common, "--id-label", label,
                                        "--include-configuration"])
            container = self.find_container(label)
            before = self.inspect_home(name, project, container, "first")
            expected_user = "dev" if name == "custom-user" else "aibox"
            check = ("set -e; test \"$(id -un)\" = " + expected_user +
                     "; test -w \"$HOME\"; tmux -V; yazi --version; codex --version")
            self.cli_step(f"exec-{name}", ["exec", *common, "--id-label", label,
                                          "bash", "-lc", check])
            marker = 'printf "v1-04-persisted\\n" > "$HOME/.aibox-v1-04-marker"'
            self.cli_step(f"write-home-marker-{name}", ["exec", *common, "--id-label", label,
                                                        "bash", "-lc", marker])
            self.step(f"remove-for-rebuild-{name}", [self.engine, "rm", "-f", container])
            reup = ["up", *common, "--id-label", label, "--include-configuration"]
            if name in ("minimal", "curated"):
                reup.append("--build-no-cache")
            self.cli_step(f"rebuild-up-{name}", reup)
            recreated = self.find_container(label)
            after = self.inspect_home(name, project, recreated, "rebuilt")
            if before != after:
                raise RuntimeError(f"{name}: home mount identity changed after recreate")
            verify = 'test "$(cat "$HOME/.aibox-v1-04-marker")" = v1-04-persisted && test -w "$HOME"'
            self.cli_step(f"verify-home-marker-{name}", ["exec", *common, "--id-label",
                                                         label, "bash", "-lc", verify])

    def cleanup(self) -> None:
        if self.engine:
            for name in NAMES:
                label = f"aibox.v1-04.run={self.run_id}-{name}"
                try:
                    self.container_ids.update(
                        self.docker(["ps", "-aq", "--filter", f"label={label}"]).splitlines())
                except Exception as exc:
                    self.cleanup_failures.append(f"list containers for {name}: {exc}")
            for container_id in sorted(self.container_ids):
                try:
                    payload = json.loads(self.docker(["inspect", container_id]))[0]
                except subprocess.CalledProcessError:
                    continue  # Already removed for the rebuild check.
                except Exception as exc:
                    self.cleanup_failures.append(f"inspect container {container_id}: {exc}")
                    continue
                label = (payload.get("Config", {}).get("Labels") or {}).get("aibox.v1-04.run")
                if label not in {f"{self.run_id}-{name}" for name in NAMES}:
                    self.cleanup_failures.append(f"refused container without run label: {container_id}")
                    continue
                for mount in payload.get("Mounts", []):
                    if mount.get("Type") == "volume" and mount.get("Destination") in (
                        "/home/aibox", "/home/dev"
                    ) and mount.get("Name", "").startswith("aibox-home-"):
                        self.home_volumes.add(mount["Name"])
                try:
                    self.docker(["rm", "-f", container_id])
                except Exception as exc:
                    self.cleanup_failures.append(f"remove container {container_id}: {exc}")
            for volume in sorted(self.home_volumes):
                try:
                    self.docker(["volume", "rm", volume])
                except Exception as exc:
                    self.cleanup_failures.append(f"remove volume {volume}: {exc}")
            for name in NAMES:
                image = f"aibox-v1-04-{self.run_id}-{name}"
                try:
                    self.docker(["image", "inspect", image])
                except subprocess.CalledProcessError:
                    continue
                try:
                    self.docker(["image", "rm", image])
                except Exception as exc:
                    self.cleanup_failures.append(f"remove image {image}: {exc}")
        shutil.rmtree(self.work, ignore_errors=False)

    def write_evidence(self) -> None:
        artifacts = []
        for file in sorted(self.output.iterdir()):
            if file.is_file() and file.name != "evidence.json":
                artifacts.append({"path": file.name, "sha256": "sha256:" + digest(file)})
        passed = not self.failures and not self.cleanup_failures
        evidence = {
            "phase": "V1-04", "verification": "preflight" if self.preflight else "host",
            "sourceCommit": self.source_commit, "startedAt": self.started_at,
            "finishedAt": now(), "runId": self.run_id, "status": "passed" if passed else "failed",
            "engine": self.engine, "runtime": self.runtime, "steps": self.steps,
            "failures": self.failures, "cleanupFailures": self.cleanup_failures,
            "cleanup": {"runOwnedContainersOnly": True, "inspectedHomeVolumesOnly": True,
                        "uniquelyTaggedBuildImagesOnly": True,
                        "retainedUnownedImageRefs": sorted(self.image_refs)},
            "artifacts": artifacts,
        }
        (self.output / "evidence.json").write_text(json.dumps(evidence, indent=2) + "\n")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", help="New absolute evidence directory")
    parser.add_argument("--preflight", action="store_true",
                        help="Check pinned CLI and Templates without a container engine")
    args = parser.parse_args()
    try:
        source_commit = require_clean_checkout()
        output = prepare_output(args.output)
    except Exception as exc:
        print(exc, file=sys.stderr)
        return 2
    work = Path(tempfile.mkdtemp(prefix="aibox-v1-04-host-run-")).resolve()
    gate = Gate(output, work, source_commit, args.preflight)
    try:
        gate.run()
    except BaseException as exc:
        detail = str(exc) or type(exc).__name__
        gate.failures.append(detail)
        print(f"V1-04 host gate failed: {detail}", file=sys.stderr)
    try:
        gate.cleanup()
    except BaseException as exc:
        gate.cleanup_failures.append(f"temporary workdir: {exc}")
        print(f"V1-04 cleanup failed: {exc}", file=sys.stderr)
    gate.write_evidence()
    print(f"V1-04 evidence: {output}")
    return 0 if not gate.failures and not gate.cleanup_failures else 1


if __name__ == "__main__":
    raise SystemExit(main())

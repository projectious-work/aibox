"""Host adapters: container engine, pinned Dev Container CLI, local registry,
publish script and the aibox binary. Python 3 and Docker are the only host
prerequisites; the pinned CLI is bootstrapped like the V1-04 host verifier."""

from __future__ import annotations

import importlib.util
import json
import os
import platform
import re
import shutil
import subprocess
import time
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from pathlib import Path
from typing import Sequence

from .catalog import REPO

V1_04_VERIFIER = REPO / "scripts" / "verify-v1-04-host.py"
PUBLISH_SCRIPT = REPO / "scripts" / "v1-05-features.sh"
REGISTRY_HOST = "localhost:5000"
REGISTRY_PREFIX = "aibox-features"
DEFAULT_REGISTRY_IMAGE = "registry:2.8.3"


@dataclass
class Completed:
    argv: list[str]
    exit_code: int
    stdout: str
    stderr: str
    seconds: float

    @property
    def combined(self) -> str:
        return self.stdout + ("\n" if self.stdout and not self.stdout.endswith("\n") else "") \
            + self.stderr


def execute(argv: Sequence[str], *, log: Path | None = None, cwd: Path | None = None,
            env: dict | None = None, timeout: float | None = None) -> Completed:
    started = time.monotonic()
    try:
        completed = subprocess.run(list(argv), cwd=cwd, env=env, text=True,
                                   capture_output=True, check=False, timeout=timeout,
                                   errors="replace")
        result = Completed(list(argv), completed.returncode, completed.stdout,
                           completed.stderr, time.monotonic() - started)
    except subprocess.TimeoutExpired as exc:
        result = Completed(list(argv), 124, _text(exc.stdout), _text(exc.stderr) +
                           f"\ntimed out after {timeout}s", time.monotonic() - started)
    if log is not None:
        with log.open("a", encoding="utf-8") as target:
            target.write(f"$ {' '.join(argv)}\nexit {result.exit_code} "
                         f"({result.seconds:.1f}s)\n--- stdout\n{result.stdout}\n"
                         f"--- stderr\n{result.stderr}\n\n")
    return result


def _text(value) -> str:
    if value is None:
        return ""
    return value.decode("utf-8", "replace") if isinstance(value, bytes) else value


def load_v1_04_verifier():
    spec = importlib.util.spec_from_file_location("aibox_verify_v1_04_host", V1_04_VERIFIER)
    if spec is None or spec.loader is None:
        raise RuntimeError(f"cannot load {V1_04_VERIFIER}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


@dataclass
class DevContainerCli:
    """The pinned Dev Container CLI. `prefix` is the argv that runs it; the
    publish script finds the same CLI and Node as `devcontainer` and `node`
    on PATH (chapter 8 Publishing), with `shim_dir` first when bootstrapped."""
    prefix: list
    runtime: dict
    shim_dir: Path | None = None

    @classmethod
    def bootstrap(cls, work: Path, output: Path) -> "DevContainerCli":
        verifier = load_v1_04_verifier()
        cli_work = work / "devcontainer-cli"
        cli_work.mkdir()
        node, cli, runtime = verifier.prepare_cli(cli_work, output)
        shim_dir = work / "bin"
        shim_dir.mkdir()
        for name, target in (("devcontainer", f'"{node}" "{cli}"'), ("node", f'"{node}"')):
            shim = shim_dir / name
            shim.write_text(f'#!/bin/sh\nexec {target} "$@"\n')
            shim.chmod(0o755)
        runtime["source"] = "bootstrapped from tools/vendor"
        return cls([str(node), str(cli)], runtime, shim_dir)

    @classmethod
    def from_path(cls, output: Path) -> "DevContainerCli":
        """Use `devcontainer` and `node` already on PATH; the CLI must be the
        pinned version."""
        devcontainer, node = shutil.which("devcontainer"), shutil.which("node")
        if not devcontainer or not node:
            raise RuntimeError("--cli-from-path needs devcontainer and node on PATH")
        version = execute([devcontainer, "--version"]).stdout.strip()
        expected = load_v1_04_verifier().CLI_VERSION
        if version != expected:
            raise RuntimeError(f"devcontainer on PATH is {version!r}; the pinned CLI is {expected}")
        runtime = {"source": "PATH", "cliPath": devcontainer, "cliVersion": version,
                   "nodePath": node, "nodeVersion": execute([node, "--version"]).stdout.strip()}
        (output / "cli-runtime.json").write_text(json.dumps(runtime, indent=2) + "\n")
        return cls([devcontainer], runtime)

    def argv(self, *args: str) -> list[str]:
        return [*self.prefix, *args]

    def environment(self) -> dict:
        env = dict(os.environ)
        if self.shim_dir is not None:
            env["PATH"] = f"{self.shim_dir}{os.pathsep}{env.get('PATH', '')}"
        return env


@dataclass
class Docker:
    binary: str = "docker"
    resources: dict = field(default_factory=lambda: {"containers": set(), "images": set()})

    @classmethod
    def locate(cls) -> "Docker":
        binary = shutil.which("docker")
        if not binary:
            raise RuntimeError("Docker is required: no docker executable on PATH")
        if execute([binary, "info"]).exit_code:
            raise RuntimeError("Docker is installed but `docker info` fails")
        return cls(binary)

    def run(self, *args: str, log: Path | None = None, timeout: float | None = None) -> Completed:
        return execute([self.binary, *args], log=log, timeout=timeout)

    def server_architecture(self) -> str:
        completed = self.run("version", "--format", "{{.Server.Arch}}")
        if completed.exit_code:
            raise RuntimeError(f"cannot read Docker server architecture: {completed.stderr}")
        return completed.stdout.strip()

    def version(self) -> str:
        return self.run("version", "--format", "{{.Server.Version}}").stdout.strip()

    def image_digest(self, image: str) -> str | None:
        completed = self.run("image", "inspect", "--format", "{{json .RepoDigests}}", image)
        if completed.exit_code:
            return None
        digests = json.loads(completed.stdout or "[]")
        return digests[0] if digests else None

    def remove_images(self, reference_pattern: str, log: Path | None = None) -> list[str]:
        failures = []
        listing = self.run("images", "--filter", f"reference={reference_pattern}",
                           "--format", "{{.Repository}}:{{.Tag}}")
        for image in sorted(set(listing.stdout.split())):
            if self.run("image", "rm", "-f", image, log=log).exit_code:
                failures.append(f"remove image {image}")
        return failures


@dataclass
class Registry:
    docker: Docker
    image: str
    name: str
    started: bool = False

    def start(self, log: Path) -> None:
        if self.docker.run("pull", self.image, log=log).exit_code:
            raise RuntimeError(f"cannot pull registry image {self.image}; see {log}")
        completed = self.docker.run("run", "-d", "--name", self.name, "-p", "5000:5000",
                                    "--label", f"aibox.v1-05.registry={self.name}",
                                    self.image, log=log)
        if completed.exit_code:
            raise RuntimeError(f"cannot start the registry on localhost:5000; see {log}")
        self.started = True
        deadline = time.monotonic() + 60
        while time.monotonic() < deadline:
            try:
                with urllib.request.urlopen(f"http://{REGISTRY_HOST}/v2/", timeout=5) as reply:
                    if reply.status == 200:
                        return
            except (urllib.error.URLError, ConnectionError, OSError):
                time.sleep(1)
        raise RuntimeError("registry on localhost:5000 did not become ready")

    def stop(self, log: Path) -> list[str]:
        if not self.started:
            return []
        if self.docker.run("rm", "-f", "-v", self.name, log=log).exit_code:
            return [f"remove registry container {self.name}"]
        return []

    @staticmethod
    def get(path: str, accept: str | None = None) -> bytes:
        request = urllib.request.Request(f"http://{REGISTRY_HOST}{path}")
        if accept:
            request.add_header("Accept", accept)
        with urllib.request.urlopen(request, timeout=60) as reply:
            return reply.read()

    def feature_layer(self, feature_id: str, tag: str) -> bytes:
        repository = f"{REGISTRY_PREFIX}/{feature_id}"
        manifest = json.loads(self.get(f"/v2/{repository}/manifests/{tag}",
                                       "application/vnd.oci.image.manifest.v1+json"))
        layers = manifest.get("layers") or []
        if len(layers) != 1:
            raise RuntimeError(f"{repository}:{tag}: expected one layer, found {len(layers)}")
        return self.get(f"/v2/{repository}/blobs/{layers[0]['digest']}")


def publish(cli: DevContainerCli, log: Path) -> Completed:
    argv = [str(PUBLISH_SCRIPT), "publish", "--registry", f"{REGISTRY_HOST}/{REGISTRY_PREFIX}"]
    if not PUBLISH_SCRIPT.is_file():
        with log.open("a") as target:
            target.write(f"{PUBLISH_SCRIPT} does not exist\n")
        return Completed(argv, 127, "", f"{PUBLISH_SCRIPT} does not exist", 0.0)
    return execute(argv, log=log, cwd=REPO, env=cli.environment(), timeout=7200)


def go_version() -> str:
    match = re.search(r"^go (\S+)$", (REPO / "go.mod").read_text(), re.MULTILINE)
    if not match:
        raise RuntimeError("go.mod declares no go version")
    return match.group(1)


def build_aibox(docker: Docker, work: Path, log: Path) -> tuple[Path, dict]:
    """Build the product binary for the host from this checkout in the Go
    toolchain image named by go.mod; the host needs no Go installation."""
    system = platform.system().lower()
    machine = {"x86_64": "amd64", "amd64": "amd64", "arm64": "arm64",
               "aarch64": "arm64"}.get(platform.machine().lower())
    if system not in ("linux", "darwin") or machine is None:
        raise RuntimeError(f"unsupported host for the aibox binary: {system}/{platform.machine()}")
    out = work / "aibox-bin"
    out.mkdir()
    image = f"golang:{go_version()}"
    user = ["--user", f"{os.getuid()}:{os.getgid()}", "-e", "HOME=/tmp"] \
        if platform.system() == "Linux" else []
    completed = docker.run(
        "run", "--rm", *user, "-v", f"{REPO}:/src:ro", "-v", f"{out}:/out", "-w", "/src",
        "-e", "CGO_ENABLED=0", "-e", f"GOOS={system}", "-e", f"GOARCH={machine}",
        "-e", "GOFLAGS=-buildvcs=false", "-e", "GOTOOLCHAIN=local",
        "-e", "GOCACHE=/tmp/go-cache", "-e", "GOMODCACHE=/tmp/go-mod",
        image, "go", "build", "-o", "/out/aibox", "./cmd/aibox", log=log, timeout=1800)
    if completed.exit_code:
        raise RuntimeError(f"building aibox failed; see {log}")
    binary = out / "aibox"
    return binary, {"builtWith": image, "imageDigest": docker.image_digest(image),
                    "goos": system, "goarch": machine}

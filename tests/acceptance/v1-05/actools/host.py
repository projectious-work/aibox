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
import signal
import subprocess
import threading
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
CGROUP = Path("/sys/fs/cgroup")
PROCESS_BUDGET_SHARE = 0.8  # the rest stays free for whatever else shares the cgroup
# Time limits in seconds. A build or a probe has its own limit (executor.py);
# these bound everything else, so that an engine that stops answering costs
# minutes, not hours.
ENGINE_TIMEOUT = 300   # an ordinary engine command: ps, rm, images, inspect, info
CLEANUP_TIMEOUT = 120  # removing the run's containers, images and registry
PULL_TIMEOUT = 900
PUBLISH_TIMEOUT = 1800
DRAIN_TIMEOUT = 10     # reading what a killed command wrote
WATCH_INTERVAL = 60    # the watchdog asks the engine once a minute ...
WATCH_TIMEOUT = 30     # ... gives it 30 seconds to answer ...
WATCH_MISSES = 3       # ... and declares it stalled after three misses in a row

_running: set = set()
_running_lock = threading.Lock()
_abort_reason: list[str] = []


def process_usage(cgroup: Path = CGROUP) -> tuple[int, int] | None:
    """Tasks in the cgroup and its limit, or None when there is no limit to
    read. The limit counts threads, not processes."""
    try:
        limit = (cgroup / "pids.max").read_text().strip()
        return (int((cgroup / "pids.current").read_text()), int(limit)) if limit != "max" else None
    except (OSError, ValueError):
        return None


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
    """Run one command in a process group of its own. On a timeout, an abort or
    an interrupt the whole group is killed: the Dev Container CLI starts
    docker children that would otherwise keep running and keep the pipes open."""
    started = time.monotonic()
    process = subprocess.Popen(list(argv), cwd=cwd, env=env, text=True, errors="replace",
                               stdin=subprocess.DEVNULL, stdout=subprocess.PIPE,
                               stderr=subprocess.PIPE, start_new_session=True)
    with _running_lock:
        _running.add(process)
    try:
        stdout, stderr = process.communicate(timeout=timeout)
        exit_code = process.returncode
    except subprocess.TimeoutExpired:
        _kill(process)
        stdout, stderr = _drain(process)
        exit_code, stderr = 124, stderr + f"\ntimed out after {timeout}s"
    except BaseException:
        _kill(process)
        raise
    finally:
        with _running_lock:
            _running.discard(process)
    if exit_code < 0 and _abort_reason:
        stderr += f"\nkilled because the run was aborted: {_abort_reason[0]}"
    result = Completed(list(argv), exit_code, stdout, stderr, time.monotonic() - started)
    if log is not None:
        with log.open("a", encoding="utf-8") as target:
            target.write(f"$ {' '.join(argv)}\nexit {result.exit_code} "
                         f"({result.seconds:.1f}s)\n--- stdout\n{result.stdout}\n"
                         f"--- stderr\n{result.stderr}\n\n")
    return result


def _kill(process: subprocess.Popen) -> None:
    try:
        os.killpg(process.pid, signal.SIGKILL)
    except (ProcessLookupError, PermissionError):
        pass
    try:
        process.kill()
    except ProcessLookupError:
        pass


def _drain(process: subprocess.Popen) -> tuple[str, str]:
    try:
        stdout, stderr = process.communicate(timeout=DRAIN_TIMEOUT)
        return stdout or "", stderr or ""
    except subprocess.TimeoutExpired:  # a child left the group and holds the pipes
        return "", "[output not readable after the kill]"


def reset_abort() -> None:
    """Forget an earlier abort; for tests that run several executors."""
    with _running_lock:
        _abort_reason.clear()


def abort_running(reason: str) -> int:
    """Kill every command in progress, in every thread, and mark later kills
    with the reason. Returns how many commands were killed."""
    with _running_lock:
        if not _abort_reason:
            _abort_reason.append(reason)
        processes = list(_running)
    for process in processes:
        _kill(process)
    return len(processes)


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
    cli_environment: dict = field(default_factory=dict)

    @classmethod
    def locate(cls, engine: str | None = None) -> "Docker":
        """Find the named engine, or Docker then Podman. Rootless Podman has no
        server, so its version and architecture come from the client side."""
        candidates = (engine,) if engine else ("docker", "podman")
        binary = next((found for found in map(shutil.which, candidates) if found), None)
        if not binary:
            raise RuntimeError(f"a container engine is required: no {' or '.join(candidates)} "
                               "executable on PATH")
        if execute([binary, "info"], timeout=ENGINE_TIMEOUT).exit_code:
            raise RuntimeError(f"{binary} is installed but `{binary} info` fails")
        return cls(binary)

    @property
    def kind(self) -> str:
        return "podman" if Path(self.binary).name == "podman" else "docker"

    def for_worker(self, index: int, work: Path) -> "Docker":
        """An engine for one parallel build slot: the same engine and store,
        with a temporary directory of its own for the Dev Container CLI."""
        tmp = work / f"worker-{index}" / "tmp"
        tmp.mkdir(parents=True)
        environment = {"TMPDIR": str(tmp)}
        if self.kind == "podman":
            environment["XDG_RUNTIME_DIR"] = self.runtime_directory()
        return Docker(self.binary, cli_environment=environment)

    def runtime_directory(self) -> str:
        """Without XDG_RUNTIME_DIR, rootless Podman keeps runtime state under
        TMPDIR. A container started with another TMPDIR then cannot be reached
        by exec, stop or rm, so every slot pins the engine's own directory."""
        if os.environ.get("XDG_RUNTIME_DIR"):
            return os.environ["XDG_RUNTIME_DIR"]
        completed = self.run("info", "--format", "{{.Host.RemoteSocket.Path}}")
        socket = completed.stdout.strip()
        if completed.exit_code or not socket:
            raise RuntimeError(f"cannot read the Podman runtime directory: {completed.stderr}")
        return str(Path(socket).parent.parent)

    def run(self, *args: str, log: Path | None = None,
            timeout: float | None = ENGINE_TIMEOUT) -> Completed:
        return execute([self.binary, *args], log=log, timeout=timeout)

    def answers(self, timeout: float = WATCH_TIMEOUT) -> Completed:
        """`<engine> info` within the limit; exit 0 means the engine answers."""
        return self.run("info", "--format", "{{.ServerVersion}}" if self.kind == "docker"
                        else "{{.Version.Version}}", timeout=timeout)

    def server_architecture(self) -> str:
        completed = (self.run("info", "--format", "{{.Host.Arch}}") if self.kind == "podman"
                     else self.run("version", "--format", "{{.Server.Arch}}"))
        if completed.exit_code:
            raise RuntimeError(f"cannot read the engine architecture: {completed.stderr}")
        return completed.stdout.strip()

    def version(self) -> str:
        side = "Client" if self.kind == "podman" else "Server"
        return self.run("version", "--format", f"{{{{.{side}.Version}}}}").stdout.strip()

    def image_digest(self, image: str) -> str | None:
        completed = self.run("image", "inspect", "--format", "{{json .RepoDigests}}", image)
        if completed.exit_code:
            return None
        digests = json.loads(completed.stdout or "[]")
        return digests[0] if digests else None

    def remove_images(self, reference_pattern: str, log: Path | None = None) -> list[str]:
        failures = []
        listing = self.run("images", "--filter", f"reference={reference_pattern}",
                           "--format", "{{.Repository}}:{{.Tag}}", timeout=CLEANUP_TIMEOUT)
        if listing.exit_code:
            failures.append(f"list images {reference_pattern}")
        for image in sorted(set(listing.stdout.split())):
            if self.run("image", "rm", "-f", image, log=log, timeout=CLEANUP_TIMEOUT).exit_code:
                failures.append(f"remove image {image}")
        return failures


@dataclass
class Registry:
    docker: Docker
    image: str
    name: str
    started: bool = False

    def start(self, log: Path) -> None:
        if self.docker.run("pull", self.image, log=log, timeout=PULL_TIMEOUT).exit_code:
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
        if self.docker.run("rm", "-f", "-v", self.name, log=log,
                           timeout=CLEANUP_TIMEOUT).exit_code:
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
    return execute(argv, log=log, cwd=REPO, env=cli.environment(), timeout=PUBLISH_TIMEOUT)


class EngineWatchdog:
    """Asks the engine every `interval` seconds whether it answers. After
    `misses` failed answers in a row it calls `on_stall` once with the reason
    and stops; the caller aborts the run. Host facts that do not need the
    engine go to `diagnostics`, so the evidence shows the state at the stall."""

    def __init__(self, docker: "Docker", on_stall, diagnostics: Path,
                 interval: float = WATCH_INTERVAL, timeout: float = WATCH_TIMEOUT,
                 misses: int = WATCH_MISSES):
        self.docker, self.on_stall, self.diagnostics = docker, on_stall, diagnostics
        self.interval, self.timeout, self.misses = interval, timeout, misses
        self.stalled: str | None = None
        self._stop = threading.Event()
        self._thread = threading.Thread(target=self._watch, name="engine-watchdog", daemon=True)

    def start(self) -> "EngineWatchdog":
        self._thread.start()
        return self

    def stop(self) -> None:
        self._stop.set()
        self._thread.join(timeout=self.timeout + DRAIN_TIMEOUT + 5)

    def _watch(self) -> None:
        missed = 0
        while not self._stop.wait(self.interval):
            answer = self.docker.answers(self.timeout)
            if answer.exit_code == 0:
                missed = 0
                continue
            missed += 1
            with self.diagnostics.open("a", encoding="utf-8") as target:
                target.write(f"{time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())} "
                             f"`{self.docker.binary} info` miss {missed} of {self.misses}: "
                             f"exit {answer.exit_code}\n{answer.stderr.strip()}\n\n")
            if missed >= self.misses:
                self.stalled = (f"the container engine did not answer `{self.docker.binary} "
                                f"info` {missed} times in a row ({self.timeout:g}s limit, "
                                f"every {self.interval:g}s)")
                host_facts(self.diagnostics)
                self.on_stall(self.stalled)
                return


def host_facts(log: Path) -> None:
    """Disk, memory, load and the engine service, read without the engine."""
    for argv in (["df", "-h"], ["df", "-i"], ["free", "-m"], ["cat", "/proc/loadavg"],
                 ["uptime"], ["systemctl", "status", "docker", "--no-pager", "--lines", "40"],
                 ["journalctl", "-u", "docker", "--no-pager", "-n", "80"]):
        if shutil.which(argv[0]):
            execute(argv, log=log, timeout=WATCH_TIMEOUT)


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
    # Under rootless Podman the container's root is the host user.
    user = ["--user", f"{os.getuid()}:{os.getgid()}", "-e", "HOME=/tmp"] \
        if platform.system() == "Linux" and docker.kind == "docker" else []
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

"""Journey sandboxes: a fresh HOME, a working directory and project copies.

LocalSandbox runs aibox in the development container as the current, non-root
user. ContainerSandbox runs it in the disposable root container (root
journeys): files are written there with podman exec, and aibox runs as an
ordinary user made by useradd, through `podman exec --user`.

Layout of a sandbox root R (R is outside every HOME):
  R/home      fresh HOME, mode 0700
  R/cwd       working directory
  R/projects  project copies (outside HOME unless a journey says otherwise)
"""

from __future__ import annotations

import itertools
import json
import os
import shutil
import stat
import subprocess
import tempfile
import uuid
from pathlib import Path

from .contract import CONTRACT
from .envscrub import env_argv, journey_env
from .proc import InfrastructureError, Result, run

DEBIAN_IMAGE = ("docker.io/library/debian@sha256:"
                "a72b4433de1cf7d6bae0c4cee74d0547f58601a643b3375f070c51600d63b6b2")
CONTAINER_USER = "tester"
CONTAINER_AIBOX = "/usr/local/bin/aibox"

FIXTURES = {"P": CONTRACT.fixture_minimal, "W": CONTRACT.fixture_workspace}


def settings_text(obj) -> str:
    return obj if isinstance(obj, str) else json.dumps(obj)


class Sandbox:
    """Common behaviour; subclasses implement the file and process primitives."""

    aibox: str
    root: str

    def __init__(self) -> None:
        self.home = f"{self.root}/home"
        self.cwd = f"{self.root}/cwd"
        self.projects = f"{self.root}/projects"
        self.state = f"{self.home}/.local/state/aibox"
        self.logs = f"{self.state}/logs"

    # primitives ------------------------------------------------------------
    def mkdir(self, path: str, mode: int = 0o700) -> None: raise NotImplementedError
    def write(self, path: str, data: str | bytes, mode: int = 0o600) -> str: raise NotImplementedError
    def copy_tree(self, src: Path, dest: str) -> None: raise NotImplementedError
    def exec(self, argv: list[str], env: dict[str, str], cwd: str, **kw) -> Result: raise NotImplementedError

    # helpers -----------------------------------------------------------------
    def tilde(self, path: str) -> str:
        """A path as the view writes it: the home prefix as ~ (pathForm)."""
        if path == self.home:
            return "~"
        if path.startswith(self.home + "/"):
            return "~" + path[len(self.home):]
        return path

    def user_settings_path(self, xdg: str | None = None) -> str:
        return f"{xdg or self.home + '/.config'}/aibox/settings.json"

    def user_settings(self, obj, mode: int = 0o600, xdg: str | None = None) -> str:
        return self.write(self.user_settings_path(xdg), settings_text(obj), mode)

    def env_file(self, name: str, text: str | bytes, mode: int = 0o600, where: str | None = None) -> str:
        return self.write(f"{where or self.cwd}/{name}", text, mode)

    def project(self, kind: str = "P", *, name: str | None = None, parent: str | None = None,
                edit=None) -> str:
        """Copy fixture P or W; `edit` maps the devcontainer.json text."""
        dest = f"{parent or self.projects}/{name or kind}"
        self.copy_tree(FIXTURES[kind], dest, edit=edit)
        return dest

    def env(self, extra: dict[str, str] | None = None, *, home: bool = True,
            path: str | None = None) -> dict[str, str]:
        kw = {"path": path} if path else {}
        return journey_env(self.home if home else None, extra, **kw)

    def aibox_run(self, args: list[str], env: dict[str, str] | None = None, *,
                  cwd: str | None = None, home: bool = True, path: str | None = None,
                  **kw) -> Result:
        return self.exec([self.aibox, *args], self.env(env, home=home, path=path),
                         cwd or self.cwd, **kw)

    def view_args(self, project: str, *extra: str, log_format: bool = True,
                  fmt: str = "json") -> list[str]:
        args = ["inspect", "--context", "local", "--project", project, "--effective-config",
                "--format", fmt]
        if log_format:
            args += ["--log-format", "jsonl"]
        return args + list(extra)

    def view(self, project: str, *extra: str, env: dict[str, str] | None = None,
             log_format: bool = True, cwd: str | None = None) -> Result:
        return self.aibox_run(self.view_args(project, *extra, log_format=log_format), env, cwd=cwd)

    def rejection_args(self, project: str, *extra: str) -> list[str]:
        return ["inspect", "--project", project, "--format", "json", "--log-format", "jsonl", *extra]


class LocalSandbox(Sandbox):
    _counter = itertools.count(1)

    def __init__(self, base: str, aibox: str, label: str):
        safe = "".join(c if c.isalnum() or c in "-_" else "_" for c in label)
        self.root = os.path.realpath(tempfile.mkdtemp(prefix=f"{next(self._counter):03d}-{safe}-", dir=base))
        self.aibox = aibox
        super().__init__()
        for d in (self.home, self.cwd, self.projects):
            os.makedirs(d, mode=0o700)
            os.chmod(d, 0o700)

    def mkdir(self, path: str, mode: int = 0o700) -> None:
        p = Path(path)
        missing = []
        while not p.exists():
            missing.append(p)
            p = p.parent
        for d in reversed(missing):
            d.mkdir()
            os.chmod(d, 0o700)
        os.chmod(path, mode)

    def write(self, path: str, data: str | bytes, mode: int = 0o600) -> str:
        parent = os.path.dirname(path)
        if not os.path.isdir(parent):
            self.mkdir(parent)
        if isinstance(data, str):
            data = data.encode()
        fd = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_TRUNC, 0o600)
        with os.fdopen(fd, "wb") as fh:
            fh.write(data)
        os.chmod(path, mode)
        return path

    def copy_tree(self, src: Path, dest: str, edit=None) -> None:
        shutil.copytree(src, dest, symlinks=True)
        if edit is not None:
            dc = Path(dest) / ".devcontainer" / "devcontainer.json"
            dc.write_text(edit(dc.read_text()))

    def exec(self, argv, env, cwd, **kw) -> Result:
        return run(argv, env=env, cwd=cwd, **kw)

    def cleanup(self) -> None:
        for dirpath, dirnames, _ in os.walk(self.root):
            for d in dirnames:
                p = os.path.join(dirpath, d)
                if not os.path.islink(p):
                    try:
                        os.chmod(p, 0o700)
                    except OSError:
                        pass
        try:
            os.chmod(self.root, 0o700)
        except OSError:
            pass
        shutil.rmtree(self.root, ignore_errors=True)


def snapshot(root: str) -> dict[str, tuple]:
    """Path -> (type, mode, size, mtime_ns, content or link target) for a tree."""
    out: dict[str, tuple] = {}
    for dirpath, dirnames, filenames in os.walk(root):
        for name in sorted(dirnames + filenames):
            p = os.path.join(dirpath, name)
            st = os.lstat(p)
            rel = os.path.relpath(p, root)
            if stat.S_ISLNK(st.st_mode):
                out[rel] = ("link", os.readlink(p))
            elif stat.S_ISDIR(st.st_mode):
                out[rel] = ("dir", stat.S_IMODE(st.st_mode), st.st_mtime_ns)
            else:
                try:
                    content = Path(p).read_bytes()
                except OSError:
                    content = None
                out[rel] = ("file", stat.S_IMODE(st.st_mode), st.st_size, st.st_mtime_ns, content)
    return out


def file_state(path: str) -> tuple | None:
    """(mode, size, mtime_ns, content) of a file, or None when absent."""
    try:
        st = os.lstat(path)
    except FileNotFoundError:
        return None
    if stat.S_ISLNK(st.st_mode):
        return ("link", os.readlink(path))
    content = Path(path).read_bytes() if stat.S_ISREG(st.st_mode) else None
    return (stat.S_IMODE(st.st_mode), st.st_size, st.st_mtime_ns, content)


# --- root container -----------------------------------------------------------

def podman(*args: str, input: bytes | None = None, check: bool = True,
           timeout: float = 120) -> subprocess.CompletedProcess:
    try:
        cp = subprocess.run(["podman", *args], input=input, capture_output=True, timeout=timeout)
    except (OSError, subprocess.TimeoutExpired) as exc:
        raise InfrastructureError(f"podman {' '.join(args[:2])}: {exc}") from exc
    if check and cp.returncode != 0:
        raise InfrastructureError(f"podman {' '.join(args)} failed ({cp.returncode}): "
                                  f"{cp.stderr.decode(errors='replace').strip()}")
    return cp


class RootContainer:
    """One disposable Debian container for the root journeys of a run."""

    def __init__(self, aibox: str):
        self.aibox_host = aibox
        self.name = f"aibox-v106-{uuid.uuid4().hex[:12]}"
        self.started = False
        self.uid = None

    def start(self) -> None:
        if shutil.which("podman") is None:
            raise InfrastructureError("podman is not on PATH; root journeys cannot run")
        podman("image", "exists", DEBIAN_IMAGE)
        # Marked started first, so stop() also removes a container whose start failed.
        self.started = True
        # --network=none is not used: rootless crun in the development container
        # cannot apply it (read-only /proc/sys). Nothing in the container uses the network.
        podman("run", "-d", "--pull=never", "--name", self.name,
               DEBIAN_IMAGE, "sleep", "infinity")
        self.started = True
        podman("cp", self.aibox_host, f"{self.name}:{CONTAINER_AIBOX}")
        self.root_sh(f"chmod 0755 {CONTAINER_AIBOX}")
        self.root_sh(f"useradd --create-home --shell /bin/sh {CONTAINER_USER}")
        self.uid = self.root_sh(f"id -u {CONTAINER_USER}").stdout.decode().strip()
        if self.uid == "0":
            raise InfrastructureError("container user is root")
        cp = self.root_sh(f"test -e {CONTAINER.system} && echo present || echo absent")
        if cp.stdout.decode().strip() != "absent":
            raise InfrastructureError(f"{CONTAINER.system} exists in the fresh container")
        self.root_sh("mkdir -p /cases && chmod 0755 /cases")

    def root_sh(self, script: str, *args: str, input: bytes | None = None,
                check: bool = True) -> subprocess.CompletedProcess:
        interactive = ["-i"] if input is not None else []
        return podman("exec", *interactive, "--user", "root", self.name, "/bin/sh", "-c",
                      script, "sh", *args, input=input, check=check)

    def stop(self) -> None:
        if self.started:
            podman("rm", "-f", "-t", "0", self.name, check=False)
            self.started = False


class _ContainerPaths:
    system = CONTRACT.system_settings_path


CONTAINER = _ContainerPaths()


class ContainerSandbox(Sandbox):
    """A sandbox inside the root container. Files belong to the test user unless
    written with system_settings()."""

    _counter = itertools.count(1)

    def __init__(self, container: RootContainer, staging: str, label: str):
        self.c = container
        safe = "".join(c if c.isalnum() or c in "-_" else "_" for c in label)
        self.root = f"/cases/{next(self._counter):03d}-{safe}"
        self.aibox = CONTAINER_AIBOX
        self.staging = staging
        super().__init__()
        u = CONTAINER_USER
        self.c.root_sh(
            'rm -f "$1" && mkdir -p "$2" "$3" "$4" "$5" && chown -R ' + f"{u}:{u}" + ' "$2" '
            '&& chmod 0700 "$3" "$4" "$5"',
            CONTAINER.system, self.root, self.home, self.cwd, self.projects)

    def _user_sh(self, script: str, *args: str, input: bytes | None = None) -> subprocess.CompletedProcess:
        return podman("exec", "-i", "--user", CONTAINER_USER, self.c.name, "/bin/sh", "-c",
                      script, "sh", *args, input=input if input is not None else b"")

    def mkdir(self, path: str, mode: int = 0o700) -> None:
        self._user_sh('umask 077 && mkdir -p "$1" && chmod "$2" "$1"', path, f"{mode:o}")

    def write(self, path: str, data: str | bytes, mode: int = 0o600) -> str:
        if isinstance(data, str):
            data = data.encode()
        self._user_sh('umask 077 && mkdir -p "$(dirname "$1")" && cat > "$1" && chmod "$2" "$1"',
                      path, f"{mode:o}", input=data)
        return path

    def system_settings(self, obj, mode: int = 0o644) -> str:
        """Write /etc/aibox/settings.json as root (D12). Mode 0644: see the
        suite's notes on the settings-file mode for the system layer."""
        self.c.root_sh('mkdir -p "$(dirname "$1")" && cat > "$1" && chown root:root "$1" '
                       '&& chmod "$2" "$1"', CONTAINER.system, f"{mode:o}",
                       input=settings_text(obj).encode())
        return CONTAINER.system

    def copy_tree(self, src: Path, dest: str, edit=None) -> None:
        stage = tempfile.mkdtemp(prefix="stage-", dir=self.staging)
        local = os.path.join(stage, "tree")
        shutil.copytree(src, local, symlinks=True)
        if edit is not None:
            dc = Path(local) / ".devcontainer" / "devcontainer.json"
            dc.write_text(edit(dc.read_text()))
        self._user_sh('mkdir -p "$(dirname "$1")"', dest)
        podman("cp", local, f"{self.c.name}:{dest}")
        self.c.root_sh(f'chown -R {CONTAINER_USER}:{CONTAINER_USER} "$1"', dest)
        shutil.rmtree(stage, ignore_errors=True)

    def exec(self, argv, env, cwd, **kw) -> Result:
        full = ["podman", "exec", "--user", CONTAINER_USER, "--workdir", cwd, self.c.name,
                *env_argv(env), *argv]
        res = run(full, env=dict(os.environ), cwd="/", **kw)
        # Record what aibox saw, not the podman wrapper's own environment.
        res.env = dict(env)
        res.note = f"inside root container {self.c.name} as {CONTAINER_USER}"
        if res.exit in (125, 126, 127) and not res.stdout:
            raise InfrastructureError(f"podman exec failed ({res.exit}): {res.err.strip()}")
        return res

    def cleanup(self) -> None:
        self.c.root_sh('rm -rf "$1" && rm -f "$2"', self.root, CONTAINER.system, check=False)

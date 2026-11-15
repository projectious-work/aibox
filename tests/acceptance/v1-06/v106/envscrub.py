"""The minimal process environment of journeys.md "Common setup".

The environment holds only HOME, PATH and what the journey adds. Nothing is
inherited from the process that runs the suite: XDG_*, AIBOX_* and the
developer's own variables never reach aibox.
"""

from __future__ import annotations

DEFAULT_PATH = "/usr/bin:/bin"


def journey_env(home: str | None, extra: dict[str, str] | None = None,
                path: str = DEFAULT_PATH) -> dict[str, str]:
    """Build a fresh environment. home=None leaves HOME unset (J-CFG-23 (b))."""
    env: dict[str, str] = {"PATH": path}
    if home is not None:
        env["HOME"] = home
    for name, value in (extra or {}).items():
        if not isinstance(name, str) or not isinstance(value, str) or not name or "=" in name:
            raise ValueError(f"bad environment entry {name!r}")
        env[name] = value
    return env


def env_argv(env: dict[str, str]) -> list[str]:
    """`env -i NAME=VALUE ...` arguments for a run inside the root container."""
    return ["/usr/bin/env", "-i"] + [f"{k}={v}" for k, v in env.items()]

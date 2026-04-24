"""migrations check.

Reads the pending migration entities directly from the filesystem (the
same filesystem walk that `migration-management.list_migrations(state=
'pending')` performs server-side — kept in sync with that convention).

- INFO: count of pending migrations.
- WARN: each pending migration older than STALE_DAYS days.

Under `--fix=migrations` (or `--fix-all`), walks each pending migration
interactively and records the user's decision. Phase 1 does NOT call the
MCP `apply_migration` tool from this CLI — doctor.py is a subprocess
without an MCP client. Instead, the check surfaces a FIX-INTENT line
per confirmed migration that the user (or the calling agent) can
execute via the migration-management MCP server. This keeps the
non-negotiable "no hand-edits; all writes via MCP" invariant.

Phase 2 will teach doctor to invoke the MCP tool via a short-lived
stdio client.
"""

from __future__ import annotations

import sys
from datetime import datetime, timezone
from pathlib import Path

import yaml

from .common import CheckResult


STALE_DAYS = 14


def _load_frontmatter(path: Path) -> dict | None:
    try:
        text = path.read_text(encoding="utf-8")
    except OSError:
        return None
    if not text.startswith("---"):
        return None
    parts = text.split("---", 2)
    if len(parts) < 3:
        return None
    try:
        data = yaml.safe_load(parts[1])
    except yaml.YAMLError:
        return None
    return data if isinstance(data, dict) else None


def _list_pending(repo_root: Path) -> list[tuple[str, Path, datetime | None]]:
    out: list[tuple[str, Path, datetime | None]] = []
    pending_dir = repo_root / "context" / "migrations" / "pending"
    if not pending_dir.is_dir():
        return out
    for p in sorted(pending_dir.glob("*.md")):
        fm = _load_frontmatter(p)
        if not fm:
            continue
        md = fm.get("metadata", {}) if isinstance(fm.get("metadata"), dict) else {}
        mid = md.get("id") or p.stem
        created = md.get("created")
        ts: datetime | None = None
        if hasattr(created, "isoformat"):
            # YAML datetime
            ts = created if created.tzinfo else created.replace(tzinfo=timezone.utc)
        elif isinstance(created, str):
            try:
                ts = datetime.fromisoformat(created.replace("Z", "+00:00"))
            except ValueError:
                ts = None
        out.append((str(mid), p, ts))
    return out


def run(ctx) -> list[CheckResult]:
    repo_root: Path = ctx["repo_root"]
    results: list[CheckResult] = []

    pending = _list_pending(repo_root)
    results.append(CheckResult(
        severity="INFO",
        category="migrations",
        id="migration.pending-count",
        message=f"{len(pending)} pending migration(s)",
    ))

    now = datetime.now(timezone.utc)
    for mid, path, ts in pending:
        if ts is None:
            continue
        age_days = (now - ts).days
        if age_days >= STALE_DAYS:
            results.append(CheckResult(
                severity="WARN",
                category="migrations",
                id="migration.stale-pending",
                message=(
                    f"{mid}: pending for {age_days} day(s) "
                    f"(threshold: {STALE_DAYS})"
                ),
                entity_ref=mid,
                fixable=True,
                suggested_fix="review and apply or reject via migration-management MCP",
                fix_mcp_tool="mcp__processkit-migration-management__apply_migration",
            ))
    return results


def run_fix(ctx, results: list[CheckResult]) -> list[dict]:
    """Interactive walk over pending migrations.

    Emits a list of fix records. Phase 1 records intent only — the
    actual MCP call is performed by the invoking agent based on these
    records (see check's module docstring for rationale).
    """
    interactive = ctx.get("interactive", False)
    auto_yes = ctx.get("yes", False)
    repo_root: Path = ctx["repo_root"]
    fixes: list[dict] = []

    if not interactive and not auto_yes:
        fixes.append({
            "category": "migrations",
            "status": "skipped",
            "reason": "fix requires interactive prompt; re-run from terminal or pass --yes",
        })
        return fixes

    pending = _list_pending(repo_root)
    for mid, path, _ts in pending:
        if auto_yes:
            choice = "y"
        else:
            try:
                raw = input(f"  Apply {mid}? [y/N/s(kip)] ").strip().lower()
            except (EOFError, KeyboardInterrupt):
                print(file=sys.stderr)
                fixes.append({"category": "migrations", "entity": mid, "status": "aborted"})
                break
            choice = raw or "n"
        if choice == "y":
            fixes.append({
                "category": "migrations",
                "entity": mid,
                "status": "intent-recorded",
                "mcp_tool": "mcp__processkit-migration-management__apply_migration",
                "mcp_args": {"id": mid},
                "note": (
                    "doctor.py cannot call MCP tools from a plain subprocess "
                    "(Phase 1 limitation). Invoke apply_migration via the agent "
                    "or run `uv run migration-management/mcp/server.py` "
                    "in MCP mode."
                ),
            })
        elif choice == "s":
            fixes.append({"category": "migrations", "entity": mid, "status": "skipped"})
        else:
            fixes.append({"category": "migrations", "entity": mid, "status": "declined"})
    return fixes

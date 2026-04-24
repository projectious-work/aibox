"""Check registry for pk-doctor.

Adding a Phase 2 check: drop a new module next to these, then append it
to REGISTRY. Order matters only for report ordering on stdout; the
aggregator doesn't otherwise sequence checks.
"""

from __future__ import annotations

from . import schema_filename, sharding, migrations, drift, team_consistency, release_integrity

# (name, module). Name is the --category= token and the key in the
# per-category tally block of the doctor.report LogEntry.
REGISTRY = [
    ("schema_filename", schema_filename),
    ("sharding", sharding),
    ("migrations", migrations),
    ("drift", drift),
    ("team_consistency", team_consistency),
    ("release_integrity", release_integrity),
]


def get(name: str):
    for n, mod in REGISTRY:
        if n == name:
            return mod
    raise KeyError(f"unknown check: {name}")


def names() -> list[str]:
    return [n for n, _ in REGISTRY]

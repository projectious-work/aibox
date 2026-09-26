# 12. v0 migration, snapshots and state

This chapter makes R-MIGRATION and R-RECOVERY executable. `aibox migrate`
is a one-time, explicitly invoked host operator utility. It is not run by
`up`, `doctor`, `refresh`, a Template, or an in-container agent. The original
v0 configuration, named environment snapshots and persistent home remain
available for rollback until the operator explicitly confirms the migrated
workspace; successful conversion never deletes them automatically.

## Converter pipeline and manifest

`preview` reads the pinned v0 config (`aibox.toml`, local override, chosen
named environment, lock/source metadata), resolves v0 aliases/defaults with
the actual v0 semantics, and emits a redacted plan. Each `CFG:*` ledger row
gets `source path`, `effective value`, `v1 authority`, `destination path/key`,
`transform`, `preserved/changed/unavailable` disposition and `fixture ID`.
Every `ADDON:*` and relevant `CMD:*` row is cross-referenced to a Feature or
native workflow. Unknown keys, conflicting aliases, unsupported version
options, unresolved secret placement, or a missing destination block apply;
they do not disappear as warnings. A plan contains SHA-256 digests of all
source files and a destination inventory, never plaintext private values.

`apply` requires the plan ID and current source digests, asks for a destination
directory not containing another unreviewed project, and stages files in a
private directory on the same filesystem. Before activation it backs up
existing destination files that would change, preserving permissions and
symlinks as data but never following them. It validates staged
`devcontainer.json` against the pinned upstream schema, the aibox extension
against its closed schema, and native Dockerfile/Compose/UX references for
existence. It atomically moves one prepared destination tree into place when
possible; if cross-filesystem, use a journaled copy with per-file fsync and
explicit partial outcome. Never start/build as part of apply.

The migration journal contains source/destination versions, digests,
timestamps, steps completed, backup locations and non-secret decisions.
Resume revalidates all digests and refuses to mix a changed source with an
old staged plan. `rollback` restores only files created/replaced by this
plan, after inspecting destination changes. User edits after migration cause
a conflict and require manual choice; rollback cannot silently erase them.
Tests interrupt after every stage and compare source, user files, permissions
and auth-state hashes before and after recovery.

## Named environments and persistence

The v0 named environment is a bundle, not just a second configuration file.
Its export manifest records name, source v0 snapshot digest, selected native
Dev Container definition path, project-local override references, processkit
source/package/context selection, persistent-home scope and non-secret runtime
metadata. Store one native definition per named environment under a
project-owned `.devcontainer/environments/<safe-name>/` tree, with stable
relative paths to shared native assets; the chosen definition is passed to
the pinned upstream CLI using its documented selector. A project default may
keep `.devcontainer/devcontainer.json`. `save/switch/restore` becomes explicit
file-based version control plus a selected-definition operation, not a new
general snapshot engine. Conversion materializes every old named snapshot
into such a definition and retains the original snapshots and process
context. If upstream selection cannot reproduce a specific snapshot, mark
that environment `manual` and block the parity claim rather than merging it
into the default.

The post-migration user workflow is file-based: create/save is a reviewed
copy of the selected native definition and explicit context references,
list/detail is the version-controlled environment manifest, switch selects
the definition for the next operator `up`, and delete removes only that
project-owned definition after a dependency check. The guide gives exact
commands and an agent-edit example, while doctor validates the manifest and
references. No `aibox env set`/snapshot database is introduced. A named
environment that is running is identified by its own exact runtime ID; a
file deletion neither stops that container nor destroys its home. Acceptance
compares save, list, detail, switch, restore and delete outcomes against v0
fixtures, including processkit context and personal MCP retention.

Persistent home is a named/bind volume scoped by canonical project identity,
not a global `.aibox-home` copy. Conversion offers `retain` (default),
`copy-with-consent`, and `new-home`; it never mounts the entire host home.
Private login/cache content is excluded from previews, Git and shared logs.
Copying state requires container stopped, exact source/destination paths,
permission/ownership check, size estimate, backup and post-copy hash sample.
Disabling a harness or changing a theme leaves home untouched. Rebuild/remove
retain the volume. Purge is a separately specified future operation and is
not implemented under the generic `remove` command.

Scoped cleanup is a separate documented workflow using native runtime and
Git commands, with a read-only preview of exact owned container/cache/
worktree targets. It must check dirty worktrees and explicit data retention
before deletion; broad `prune all` is not recreated. Backup/recovery uses
version control for project files and a user-selected archive of persistent
home/context with checksums and private permissions. Processkit context reset
delegates to processkit's supported workflow. These workflows require
AC-MIG/AC-SEC evidence before v0 `prune`, `backup` or `reset` is retired.

## Required conversion fixtures

Cover fresh v0 defaults; every deprecated alias; duplicate old/new key
conflict; local env/volume overlay; explicit addon disable and version;
ordered multi-harness selection; host-only Cursor; all theme aliases;
custom tmux/Yazi/Starship files; processkit source fork and packages;
personal MCP entries; audio and multi-document LaTeX; named environment
switch with distinct context; Dockerfile.local and Compose overrides;
malformed TOML; unknown keys; symlink escape; secret-bearing local files;
interrupted activation; rollback after user edits. The generated mapping
report must have a disposition for all 265 `CFG:*` rows and all 98 tool
entries, including source fields that are unreachable in the particular
fixture. AC-MIG requires a real v0.35.0 project to build and work after
conversion, not just schema validity.

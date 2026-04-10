# aibox v0.17.5

Compatibility release for processkit v0.8.0 (GrandLily src/ restructure).

## Breaking change in processkit

processkit v0.8.0 restructured its `src/` directory to be a literal mirror of
the consumer project root. All content moved from flat top-level segments to
`context/` sub-paths:

| Content | Before (v0.7.0) | After (v0.8.0) |
|---|---|---|
| Skills | `src/skills/<name>/` | `src/context/skills/<name>/` |
| Shared lib | `src/lib/processkit/` | `src/context/skills/_lib/processkit/` |
| Schemas | `src/primitives/schemas/` | `src/context/schemas/` |
| State machines | `src/primitives/state-machines/` | `src/context/state-machines/` |
| AGENTS.md template | `src/scaffolding/AGENTS.md` | `src/AGENTS.md` |
| Packages | `src/packages/` | `src/.processkit/packages/` |
| FORMAT.md | `src/skills/FORMAT.md` | `src/.processkit/FORMAT.md` |

The **live install destinations are unchanged** — `context/skills/`,
`context/schemas/`, etc. — so existing projects are unaffected structurally.
Only the aibox installer and sync machinery needed updating.

## What changed in aibox

### Installer routing (`content_install.rs`)

`install_action_for()` now handles both v0.8.0+ and legacy v0.7.0 layouts in
one function. The two prefix sets are disjoint (`context/` vs. bare names) so
no version check is needed — projects pinned to v0.7.0 continue to install and
sync correctly with the new aibox binary.

New v0.8.0 routing: `context/...` paths install at the same path (the tarball
already mirrors the consumer root); `AGENTS.md` at the tarball root installs
as a rendered template; `.processkit/...` is skipped as catalog-only content.

### Diff grouping (`lock.rs`)

`group_for_path()` updated to match v0.8.0 PROVENANCE.toml paths
(`context/skills/<name>/...` → group `skills/<name>`, etc.) alongside the
legacy grouping. Group name strings are identical across both layouts, so the
3-way diff logic is unaffected.

### Templates mirror consumers

All code that navigates into the templates mirror (`kit`, `mcp_registration`,
`claude_commands`, `content_init`) now uses three new helpers from
`processkit_vocab`:

- `mirror_skills_dir(root, version)` — tries v0.8.0 path first, falls back to legacy
- `mirror_processes_dir(root, version)` — same pattern
- `mirror_packages_dir(root, version)` — tries `.processkit/packages/`, falls back to `packages/`

### Vocabulary constants (`processkit_vocab.rs`)

`pub mod src` restructured for v0.8.0: new `CONTEXT_DIR`, `LIB_SEGMENT`,
`DOTPROCESSKIT` constants; legacy constants preserved as `LEGACY_*` for
backward-compat routing. `PROCESSKIT_DEFAULT_VERSION` bumped to `v0.8.0`.

### `maintain.sh`

`sync-processkit` now fetches `src/.processkit/FORMAT.md` (was
`src/skills/FORMAT.md`) for vocabulary diffing.

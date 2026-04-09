# aibox v0.17.2

## What's new

### processkit v0.6.0 compatibility

aibox is now fully aligned with processkit v0.6.0:

- **`core: true` skill enforcement** — skills marked `metadata.processkit.core: true` (e.g. `skill-finder`) are installed unconditionally regardless of `[skills].include`/`[skills].exclude`. `aibox doctor` warns when a core skill appears in your exclude list.
- **Updated default version** — `PROCESSKIT_DEFAULT_VERSION` bumped to `v0.6.0` across all test fixtures and documentation.

### Central constants module (`processkit_vocab.rs`)

All processkit-related compile-time vocabulary is now consolidated in a single module:

- Path constants: `TEMPLATES_PROCESSKIT_DIR`, `LIVE_SKILLS_DIR`, `LIVE_SCHEMAS_DIR`, `LIVE_STATE_MACHINES_DIR`, `LIVE_PROCESSES_DIR`, `LIVE_LIB_DIR`
- Filename constants: `SKILL_FILENAME`, `PROVENANCE_FILENAME`, `FORMAT_FILENAME`, `INDEX_FILENAME`, `AGENTS_FILENAME`
- Source tree segments: `processkit_vocab::src::{SKILLS, PRIMITIVES, SCHEMAS, STATE_MACHINES, PROCESSES, LIB, SCAFFOLDING, PACKAGES}`
- URL constant: `PROCESSKIT_GIT_SOURCE`
- Category vocabulary: `CATEGORY_ORDER`, `category_sort_index()`
- Shared frontmatter types: `SkillFrontmatter`, `SkillProcesskitMeta` (with `core: bool`), `parse_skill_frontmatter()`

All production code that previously embedded these as string/numeric literals now references the central constants. This eliminates an entire class of copy-paste drift bugs.

### Docs trimmed to aibox's perimeter

Skill and process documentation that belongs to processkit's domain has been removed from the aibox docs site. The skills and process-packages pages now point to the upstream processkit docs (placeholder link, pending processkit doc-site launch) rather than duplicating content that processkit owns.

## Upgrade notes

- Run `aibox sync` after upgrading to pick up changes
- If you are a processkit consumer: pin `[processkit].version = "v0.6.0"` in your `aibox.toml` and run `aibox sync` to install the v0.6.0 content set (sharded entities, privacy controls, sqlite index, `skill-finder` as a core skill)

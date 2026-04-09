# Agent Handover

This file exists to orient a new agent after a `aibox reset` + `aibox init`
cycle that wiped the previous agent's memory. Read this first, then read
`AGENTS.md` for the full project context and work instructions.

---

## What this project is

**aibox** — a Rust CLI (`cli/`) that manages reproducible AI development
containers. It generates `.devcontainer/` files, installs content from
**processkit** (skills, processes, primitives, AGENTS.md), and handles the
full container lifecycle.

Companion repo: `https://github.com/projectious-work/processkit`

**The boundary:** aibox owns containers + install machinery. processkit owns
all content (skills, processes, primitives, AGENTS.md template). Do not
duplicate processkit's documentation in aibox's docs.

---

## Current state (as of 2026-04-09, v0.17.2)

- **Released:** aibox v0.17.2 — processkit v0.6.0 compatibility + central
  constants module (`cli/src/processkit_vocab.rs`). All processkit-related
  compile-time vocabulary (path prefixes, filenames, category order, frontmatter
  types) is centralized there. No hardcoded strings in production code.
- **`core: true` enforcement:** Skills with `metadata.processkit.core: true`
  (e.g. `skill-finder`) install unconditionally regardless of
  `[skills].include`/`[skills].exclude`. `aibox doctor` warns on excluded core
  skills.
- **This devcontainer** was just re-initialized from scratch with `aibox init`
  against processkit v0.6.0.

---

## First task for the new agent

**Familiarize yourself with the processkit v0.6.0 content that `aibox init`
just installed** — particularly the new processes and how they interact with
aibox's CLI. processkit v0.6.0 introduces sharding, privacy controls, a sqlite
index, and `skill-finder` as a core skill. Check:

1. What processes are now installed in `context/processes/`?
2. Does the new `AGENTS.md` reference anything aibox doesn't yet support?
3. Are there any new skills with `core: true` beyond `skill-finder`?
4. Does `aibox kit list` and `aibox kit skill list` work correctly against the
   new install?

Report findings before doing anything else.

---

## Key files

| File | What it is |
|------|-----------|
| `AGENTS.md` | Canonical agent entry point (installed by processkit) |
| `context/work-instructions/RELEASE-PROCESS.md` | **Full release process** — read before any release |
| `context/work-instructions/ARCHITECTURE.md` | System architecture |
| `context/work-instructions/DEVELOPMENT.md` | Development workflow |
| `cli/src/processkit_vocab.rs` | Central constants module — single source of truth for all processkit vocabulary |
| `scripts/maintain.sh` | Release tooling (`release`, `sync-processkit`, `release-host`, `docs-deploy`) |

---

## Release process (summary)

1. `./scripts/maintain.sh sync-processkit` — check for processkit updates,
   patch constants, review FORMAT.md diff, make any CLI changes needed, commit
2. Write `dist/RELEASE-NOTES.md` for the new version
3. `./scripts/maintain.sh release X.Y.Z` — runs in container; syncs processkit,
   tests, audits, builds both Linux binaries, tags, creates GitHub release,
   deploys docs
4. User runs `./scripts/maintain.sh release-host X.Y.Z` on macOS host —
   builds macOS binaries, uploads them, pushes container images to GHCR

---

## What to avoid

- Do not add hardcoded path strings, filenames, or vocabulary to production
  Rust code — add constants to `processkit_vocab.rs` instead
- Do not duplicate processkit's skill/process documentation in aibox's docs site
- Do not edit `CLAUDE.md` — it is a thin pointer to `AGENTS.md`

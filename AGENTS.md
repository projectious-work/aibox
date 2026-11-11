# AGENTS.md

## About & session start

**aibox** is a Rust CLI that manages reproducible, AI-ready dev containers.
Since v0.16.0 it has a strict two-part scope:

1. **Containers** — generates `.devcontainer/Dockerfile`, `docker-compose.yml`,
   and `devcontainer.json` from `aibox.toml`, plus a tool-bundle addon system
   (`addons/`) and themed `.aibox-home/` runtime config seed.
2. **processkit installer** — fetches a pinned release of
   [`projectious-work/processkit`](https://github.com/projectious-work/processkit)
   and installs its skills, primitives, processes, and the canonical `AGENTS.md`
   template into the consuming project under `context/`.

**MCP Permissions** — Since v0.18.7, `aibox sync` auto-generates harness-specific
permission files for all MCP servers. Configure `[mcp.permissions]` in `aibox.toml`
to eliminate repetitive permission prompts. Glob patterns expand into concrete
server names; deny patterns take precedence over allow for security. See
[Configuration / MCP Permissions](./docs-site/docs/reference/configuration.md#permission-configuration-mcppermissions).

Target users: solo developers, small teams, and consultants who want
reproducible AI-ready dev environments without manual Docker/devcontainer setup.
Success looks like: `aibox init` → working themed tmux session with processkit
content in place in under 5 minutes.

Provider-specific files (`CLAUDE.md`, `CODEX.md`, `.cursor/rules`, …) are thin
pointers — edit **this** file.

**For structured workflows** (releases, migrations, design reviews): check
`CONTRIBUTING.md`, `dev-notes/` and the specification under `spec/` for the
documented procedure before acting. These take precedence over general
knowledge or tool defaults. Architecture decisions are recorded as ADRs under
[`docs/adr/`](docs/adr/).

Company standards are vendored in `.standards/` (see the generated block at
the end of this file). Which standards apply to which part of this repository,
and the recorded deviations, are in the
[applicability record](docs/engineering/standards-applicability.md).

## History: processkit archive

processkit was removed from this repository's own development setup on
2026-11-11 ([ADR 0001](docs/adr/0001-remove-processkit-from-the-development-setup.md)).
Former decisions, work items, artifacts, discussions, notes and logs are plain
files under `archive/processkit/`. Search them with ordinary tools (`grep`,
`rg`); they are read-only history and must not be edited or extended. Paths
inside the archived records still say `context/…`; read them as
`archive/processkit/…`. The full former installation is preserved at the Git
tag `processkit-final`.

This concerns the development setup only. aibox's product support for
`mode = "processkit"` is product code and is unchanged.

## Setup

For the v1 Go line, use `go test ./...`, `go test -race ./...`, `go vet ./...`,
`go build ./cmd/aibox`, the relevant `spec/v1` validators, and the v1 docs
build. `./scripts/verify-v1-03.sh` runs the exact V1-03 candidate gate.
The Rust and Cargo commands below describe retained v0 work only; they are
not v1 acceptance, commit, or release requirements.

```sh
# build the CLI binary
cd cli && cargo build

# run all tests (unit + integration + E2E tier 1)
cd cli && cargo test

# lint (zero warnings required)
cd cli && cargo clippy --all-targets -- -D warnings

# format check
cd cli && cargo fmt -- --check
```

For E2E tier 2 (full container lifecycle tests, requires the `aibox-e2e-testrunner`
companion service to be running alongside the devcontainer):

```sh
cd cli && cargo test --features e2e
```

The E2E companion is reached from this devcontainer over SSH/SCP:
`ssh -i /workspace/.aibox-e2e-runner-home/.ssh/id_ed25519 testuser@aibox-e2e-testrunner`.
Do not treat missing local `docker` or `podman` in the main devcontainer as
evidence that the companion is unavailable; use SSH reachability first. The
companion itself owns the container runtime used by Tier 2 tests.

### aibox CLI commands inside this devcontainer

Do not run `aibox doctor` from inside the container to judge this live
project; it checks aibox-managed host/project posture and belongs on the host.

There are narrow exceptions when developing or releasing the aibox CLI itself.
These are host-context simulations, not dogfood operations:

- Unit, integration, and Tier 1 E2E tests may run `aibox init`, `aibox apply`,
  `aibox doctor`, and related commands in temporary projects, often with
  `AIBOX_NO_CONTAINER=1` or mocked runtimes.
- Tier 2 E2E tests may run `aibox apply`, `aibox up`, and `aibox doctor`
  against projects deployed to the `aibox-e2e-testrunner` companion; verify the
  companion over SSH first.
- Release Phase 0 may run `aibox doctor` through
  `./scripts/maintain.sh release-doctors` to exercise host-side doctor behavior
  even when launched from this devcontainer.
- Local test-installs may run the freshly built `aibox` binary against a scratch
  project to simulate what a host user would do.

Outside those explicit test/release contexts, reserve `aibox doctor` for the
host.

See `CONTRIBUTING.md` for the full development workflow, E2E test
architecture, and cross-compile steps.

## Code style & PRs

### Code style and conventions

- **Zero clippy warnings** — always run with `-D warnings`; CI rejects any warning.
- **All tests must pass** before committing; run `cargo test` and
  `cargo clippy --all-targets -- -D warnings` before every commit.
- **`cargo audit` must be clean** before tagging a release.
- **Conventional commits** — `feat:`, `fix:`, `chore:`, `docs:`. For releases that
  include real changes beyond the version bump, use `fix(vX.Y.Z): <summary>` with a
  section-by-section body and `Refs: DEC-NNN, BACK-NNN` + `Co-Authored-By:` trailer.
- **Reference GitHub issue numbers** in commits: `fixes #N`, `refs #N`.
- **Include `Cargo.lock`** in version bump commits.
- **Never force-push to `main`.**
- **No hardcoded processkit vocabulary** in production Rust code — add constants to
  `cli/src/processkit_vocab.rs` instead.
- **No trailing summaries** in agent responses — the user can read the diff.

### v1 engineering practice

v1 work on `v1.x-dev` follows
[`spec/v1/21-engineering-practice.md`](spec/v1/21-engineering-practice.md)
from V1-05: development cycle and independent review per phase, a frozen
test-ready phase contract with separate acceptance-test and implementation
tracks,
*Hexagonal Architecture* with a *Clean Architecture* core, *TDD, Chicago
School* tests, refactoring by named smells, Conventional Commits, and a green
`v1.x-dev` after every merge.

AI agents implementing v1 also follow these execution and communication
rules (not part of the product specification):

- Run at most 3 sub-agents at once and work in waves beyond that. The primary
  agent keeps concept, core implementation, integration, and review of
  delegated results.
- Give a red check at most 1 rerun. Report an infrastructure-only failure once
  with evidence. Push fixes in batches. Do not spend turns polling; wait for
  completion notifications.
- Reply answer-first: result or next action first, numbered steps, one
  closing next action, no preamble or recap. Personal output styles are opt-in
  through the person's user-level agent instructions.

### Pull requests

- **Version-line integration** — long-lived branches (`main`, `v0.x-dev`,
  `v0.x-release`, `v1.x-dev`, `v1.x-pre-release`) are PR-only. Tag releases
  only on their designated release branch; merge the tagged branch into `main`
  as published history. Do not use GitHub Actions or workflow files.
- **"Ship it" means the full release ritual** end-to-end: build, test, commit, tag,
  push, GitHub release, deploy docs. Do not ask permission at each step.
- **Phase 0 runs three steps before any version bump:** (1) `./scripts/maintain.sh release-check-state`
  → `dist/RELEASE-STATE.md`; (2) `./scripts/maintain.sh release-doctors` → `dist/RELEASE-DOCTORS.md`
  (v0 line: runs pk-doctor + aibox doctor; ERRORs from either block the release, WARNs surface in the report
  but don't block); (3) review both reports. `./scripts/maintain.sh release <version>` calls all
  three automatically.
- **Phase 2 is always the user's job** — macOS host builds and GHCR image pushes
  run via `./scripts/maintain.sh release-host X.Y.Z` on the host, never from the container.
  Before Phase 2, sync the host checkout to the just-pushed container-side release
  commits with `git fetch origin main && git reset --keep origin/main`.
- **Detailed release procedures:** See [`archive/processkit/notes/NOTE-20260411_0000-LoyalSpruce-aibox-release-process.md`](./archive/processkit/notes/NOTE-20260411_0000-LoyalSpruce-aibox-release-process.md) (v0 line, archived) for step-by-step Phase 1 and Phase 2 instructions, including exact commands and prerequisites.
- **One big change preferred over many small PRs** for breaking releases. The user
  explicitly said: "make one big change, I'll handle derived project dependencies."
- **Hold uncommitted changes when the user says "hold"** — leave in working tree,
  do not commit, do not bump version until the user says "ship".

## AI agents on this project

Configured providers: **claude**, **openai**. Other agents may be working on
this project; coordinate through GitHub issues and pull requests rather than
assuming you are alone.

## Project-specific notes

### Critical: `.devcontainer/` vs `images/`

**We are in a dev-container building dev-containers.** Never confuse:

- **`.devcontainer/`** — THIS project's own dev environment (Rust + Python/uv + Hugo/Docsy).
- **`images/`** — Published images for OTHER projects (pushed to GHCR). They do NOT include Rust toolchain or MkDocs.

Changes to `.devcontainer/` affect our development. Changes to `images/` affect downstream projects.

### Project structure

| Path | Owns |
|---|---|
| `cli/` | The Rust CLI (`aibox` binary) — the only shipped artifact besides addon YAMLs |
| `addons/` | YAML addon definitions (python, rust, node, latex, …) |
| `images/` | Container image build recipes published to GHCR |
| `docs-site/` | Hugo/Docsy documentation site |
| `archive/processkit/` | Read-only records of the former processkit setup (decisions, work items, notes, logs) |
| `docs/adr/` | Architecture decision records |
| `scripts/` | Release and maintenance tooling (`maintain.sh`, `record-asciinema.sh`, …) |

**Key Rust module:** `cli/src/processkit_vocab.rs` — central constants module.
All processkit-related compile-time vocabulary (path prefixes, filenames, category order,
frontmatter types) lives here. **Never hardcode processkit strings in production code.**

### aibox ⇄ processkit boundary

- **aibox owns:** containers, addons, the `[processkit]` config section, the
  install/diff/migrate machinery, the slim project skeleton at init time
  (`.aibox-version`, `.gitignore`, empty `context/`, thin provider pointer files
  like `CLAUDE.md`), and the docs site.
- **processkit owns:** every skill (`SKILL.md`), every primitive schema, every state
  machine, the canonical `AGENTS.md` template, processes, and the package YAMLs.
- **The `context/` directory** is shared territory: aibox creates it, processkit fills
  it. An immutable upstream reference lives under `context/templates/processkit/<version>/`
  for the three-way diff.

**If something process-related is missing, add it to processkit, not aibox.**
Do not re-introduce skills, processes, or primitives into aibox — that all belongs in
processkit now (DEC-027).

### Anti-patterns — stop and reconsider if you find yourself doing these

- Writing to `.claude/`, `.gemini/`, or any other provider directory from aibox
  code (off-perimeter per DEC-029)
- Hardcoding path strings, filenames, or processkit vocabulary in Rust source
  (add to `processkit_vocab.rs` instead)
- Trying to do Phase 2 of a release from inside the container (needs macOS host)
- Skipping `cargo audit`, `cargo clippy --all-targets -- -D warnings`, or `cargo test`
  before tagging a release
- Pointing users at the removed `aibox` skills subcommand — that surface was removed in v0.16.0 (DEC-027); processkit skills are now the only home for skill content
- **Creating GitHub releases directly** with `gh release create` — always use
  `./scripts/maintain.sh release <version>` inside the container instead. It runs
  tests, cargo audit, builds Linux binaries, creates the release with assets attached,
  deploys docs, and prints the Phase 2 prompt. A bare `gh release create` produces
  a release with no binary assets.
- **Releasing when `cargo build --release` is broken** — if the build fails (e.g.
  missing linker after a container config change), stop and tell the user rather
  than creating an empty release. A release without binaries is worse than no release.
  The precondition is: `cargo build --release --target aarch64-unknown-linux-gnu`
  must succeed inside the container before `maintain.sh release` is invoked.

### Design principles (non-negotiable)

These are load-bearing. When a design call comes up, the answer is whichever option
respects more of these principles.

1. **Provider neutrality.** No file path, config field, binary, or API surface is bound
   to a specific AI provider. Skills live under `context/skills/`, never `.claude/skills/`.
   Provider-specific files (CLAUDE.md, etc.) are thin pointers to AGENTS.md.
2. **Reproducibility.** Every consumed processkit release is pinned by `(source, version,
   sha256)` in `aibox.lock`. Moving-branch consumption is a dev fallback, not production.
3. **Locality.** Everything a project needs lives inside the project directory. A fresh
   `git clone` + `aibox sync` reproduces the environment exactly.
4. **Edit-in-place.** Installed processkit content lives at editable, top-level paths.
   The immutable upstream reference under `context/templates/processkit/<version>/` is the
   diff baseline only — not a restriction on editing.
5. **Forkability.** Every reference to processkit goes through `[processkit].source`.
   Companies can fork processkit and consume the fork by changing one line.
6. **Single source of truth.** Each piece of content lives in exactly one project.
   Skills/primitives/processes/AGENTS.md template → processkit only. Container generation/
   addon management/install pipeline → aibox only. DEC-20260411_0000-JollyClover-rip-bundled-process-layer made this strict.
7. **Generic content-source machinery.** The fetcher in `content_source.rs` is content-
   source-neutral by construction. It doesn't know "processkit" specifically — it knows
   how to fetch a release-asset tarball from any GitHub-shaped source, verify it, and
   extract it. Processkit-compatible alternatives consume the same machinery with no code change.

### Provider independence

All project state must be stored in `./context/` — never in provider-specific locations
(e.g. `.claude/memory/`, `.aider/`). This ensures any AI agent (Claude, Aider, Gemini,
etc.) can pick up where another left off, and session handovers are committed to git.
Do not write to `.claude/`, `.gemini/`, or any other provider directory from aibox code.

### MCP Permissions Troubleshooting

If you're still seeing permission prompts for aibox-shipped MCP servers:

1. **Verify `[mcp.permissions]` is configured** in `aibox.toml`:
   ```toml
   [mcp.permissions]
   default_mode = "allow"
   allow_patterns = ["mcp__processkit-*"]
   ```
   
2. **Run `aibox sync`** to regenerate harness permission files. Permission configuration is applied during sync, not during runtime.

3. **Check harness-specific behavior:**
   - **Claude Code**: Checks `.claude/settings.local.json` → run `aibox sync` to populate `permissions.allow[]`
   - **OpenCode**: Reads `.opencode/config.toml [mcp]` section → verify `mode = "allow"` and `allow[]` array
   - **Continue IDE**: Per-tool `mode` in `continue/config.json` → defaults to "Ask" for safety; override with `[mcp.permissions.harness.continue] mode = "allow"`
   - **Cursor IDE**: Checks `.cursor/settings.json allowedMcpServers[]` → verify entries match expanded server names
   - **Gemini CLI**: Dual `includeTools`/`excludeTools` in `.gemini/settings.json` → intersection semantics (both must match)
   - **GitHub Copilot**: Reads environment variables `COPILOT_MCP_ALLOW_TOOLS`, `COPILOT_MCP_DENY_TOOLS` → created as `.copilot-env`
   - **Aider**: Checks `.aider/mcp-permissions.json allowed_tools` → fallback for harnesses without native MCP permission support
   - **Codex**: Uses project-level `trust_level = "trusted"` in `.codex/config.toml` → applies to all tools

4. **Verify pattern matching:**
   - `"mcp__processkit-*"` matches all processkit MCP servers (e.g., `mcp__processkit-workitem-management__create_workitem`)
   - `"bash"` matches the Bash tool fallback
   - First-match-wins: if a tool matches both allow and deny patterns, deny takes precedence
   - Check `/workspace/.aibox/aibox.log` for pattern expansion details

5. **Check for typos in `allow_patterns`** — misspelled patterns expand to zero tools. `aibox sync` logs warnings for patterns that match no servers.

### Operational gotchas

- **Podman compose** output format varies by version — always use `inspect`, never parse `ps` output.
- **Stale image cache**: if the container exits immediately after start, rebuild with `--no-cache`.
- **`.aibox-home/` must be in `.gitignore`** — it contains SSH keys and personal config.
- **tmux version pin**: tmux is installed from the base Debian image package set; upgrade it through `images/base-debian/Dockerfile` package inputs or the base image.
- **`host.docker.internal`**: works on Docker Desktop and Podman pasta; bare Linux Docker may need `--add-host`.
- **OrbStack virtiofs**: files mounted from macOS may lose execute permissions — workaround: `chmod +x` inside container.
- **Claude Code OAuth in containers**: use `claude setup-token` or authenticate on host (credentials shared via `.claude` mount). See anthropics/claude-code#14528. Do NOT use `network_mode: host`.
- **OrbStack network dropout**: after ~20 minutes idle, OrbStack's VM NAT can drop connections. Fix: set `keepalive = true` in `[container]` of `aibox.toml` (adds a lightweight DNS keepalive every 2 minutes via `postStartCommand`).

### Tmux session scripts — two files, two roles

Two files share the name `aibox-tmux-session.sh` (or `aibox-session.sh`) but serve different purposes. Editing the wrong one is a common source of confusion:

| File | Role |
|---|---|
| `images/base-debian/config/bin/aibox-tmux-session.sh` | **IMAGE/RUNTIME variant.** Baked into `/usr/local/bin` inside the container image. Acts as the authoritative fallback launcher used by the container entrypoint before any apply-time managed files are present. No per-project customisation. |
| `cli/src/templates/aibox-home/.config/tmux/aibox-session.sh` | **GENERATED template variant.** Source template copied into `.aibox-home/.config/tmux/` at `aibox apply` time. Reflects the project owner's configured layout, session name, etc. from `aibox.toml`. |

At runtime the GENERATED variant takes precedence (via `AIBOX_TMUX_MANAGED_SESSION`); the IMAGE variant is the fallback when no managed session script is present yet. Both files carry a header banner that identifies which role they play. When changing either file, review the other for behavioural parity.

The generated layout scripts in `.aibox-home/.config/tmux/layouts/` (one per `ConfigLayout`) are rendered by `cli/src/tmux/layouts.rs::tmux_layout_script`.

### Runtime artifacts for agents (in derived projects)

When an AI agent is working inside a project that uses aibox:

| Path | Contents |
|---|---|
| `.aibox/aibox.log` | NDJSON structured log of every `aibox` command. Read to understand what aibox did recently. Rotates at 1 MB. |
| `aibox.lock` | Pinned versions of the aibox CLI and processkit last synced. |
| `context/migrations/pending/` | Pending Migration entities awaiting review; resolve through `migration-management` MCP (`apply_migration` / `reject_migration`). Never hand-move migration files. |
| `context/migrations/applied/` | Applied or rejected Migration entities (state-bucketed by processkit convention). |

### GitHub organization

- **Repo:** `projectious-work/aibox`
- **GHCR:** `ghcr.io/projectious-work/aibox`
- **Docs:** `https://projectious-work.github.io/aibox/`
- **processkit upstream:** `https://github.com/projectious-work/processkit`
- **processkit releases:** `https://github.com/projectious-work/processkit/releases`

## Company standards

<!-- BEGIN COMPANY_STANDARDS: generated by standards.py; do not edit -->

Company standards v2.3.0 from <https://github.com/projectious-work/standards>, vendored in `.standards/` (license: `.standards/README.md`). `standard-*` skills load on demand; the full standard they name is authoritative. *Italic* terms are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966): canonical meaning applies; standards state only deviations. Never edit `.standards/`, the skills, or this block; run `python3 .standards/standards.py check`.

### Always on: communication

*BLUF*: answer or next action first; commands, paths, snippets before prose.
Number multi-step work: one bounded action per step, fewest possible.
Restate state across turns; use the harness task tool when one exists.
End with one next action doable in 2 minutes, only if something is open.
Finished work: what works, how to check it, what remains.
Errors: cause, location, fix. Flag unverified claims; never fake certainty.
Estimate human work in minutes or hours, with the assumption.
No preamble, recap, pleasantries, or tangents; offer a second topic after.
Rank long lists in *MECE Principle* groups of 5;
never truncate reviews, audits, findings.
Safety and confirmation of destructive or outward-facing actions come first;
"explain" gets full depth. Sub-agent briefs and reports are machine-facing:
result, evidence, unknowns.
Details: `.standards/agent-communication.md`.

### Always on: execution budget

At most 3 sub-agents at once, in waves; keep core work yourself. Review
gates get 2 rounds, failed verifications 1 rerun, infrastructure failures one
evidenced report to the owner. Push batches, not findings. Never poll.
Details: `.standards/agent-execution-budget.md`.

<!-- END COMPANY_STANDARDS -->

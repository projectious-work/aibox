# Company standards: applicability record

This record declares which application profiles of the
[company standards](../../.standards/application-profiles.md) apply to aibox,
and where aibox deviates. Installed version: v2.3.0 (`.standards/lock.json`).
Check the installation with `python3 .standards/standards.py check`.

Review this record when the product boundary or the distribution model
changes.

## Profiles and components

| Profile | Role | Components |
|---|---|---|
| CLI application | primary | v1 Go executable `aibox` (`cmd/aibox`, `internal/`); v0 Rust CLI on `lts/v0` |
| MCP server | secondary | local stdio server `aibox mcp serve` (`internal/mcpserver`) |
| Infrastructure or template product | secondary | Dev Container Templates (`templates/`), first-party Features (`features/`), v0 images and addons (`images/`, `addons/`) |
| Schema, protocol, or process package | secondary | v1 result, request, settings and log schemas and their fixtures (`spec/v1/*.schema.json`, `spec/v1/fixtures/`) |
| Documentation or website | secondary | Hugo site (`docs-site/`), published from `gh-pages` |
| AI agent | secondary | repository agent instructions (`AGENTS.md`, `.agents/skills/`) |

Not applicable: Service or worker (aibox runs no long-lived service; the MCP
server is local stdio only), User-interface application, Library or SDK
(`internal/` is not a published API).

Overlay: **Host-gated release** applies. The development container has no
macOS host, no amd64 host, no signing key and no publication authority.
Host-only properties are proven by the owner on real hosts: the
`debian-13-amd64` acceptance run, macOS builds, image pushes and publication.

## Platforms and distribution

- v1 (in development, `main`): Linux and macOS hosts, amd64 and arm64;
  Docker, Podman and OrbStack through the pinned Dev Container CLI. The
  qualified target matrix is roadmap phase V1-22; until then each phase
  names the cells it has proven.
- v0 (current release, `lts/v0`): GitHub release binaries for Linux and
  macOS, images on GHCR, installer `scripts/install.sh`.
- Documentation: GitHub Pages, deployed locally.

## Verification and release gates

- Per phase: the phase's offline gate (`scripts/verify-v1-NN.sh`) and its
  container acceptance run (`tests/acceptance/<phase>/run-host.sh`), with
  evidence under `spec/v1/evidence/` and the phase note in `dev-notes/`.
- No CI. Verification runs locally or in the development container, and its
  result goes into the pull request.
- Release gates follow chapter 14 of the v1 specification and the company
  verification and release-engineering standard.

## Twelve-Factor alignment

aibox ships no service or worker. For the CLI and the local MCP server the
company deviations apply as listed in the standard (layered configuration,
local human-controlled release, development containers are not the runtime,
logs on stderr with file sinks). aibox records no further deviation.

## Deviations and deferred requirements

| Requirement | State | Tracking |
|---|---|---|
| Squash-only merges | Two one-time merge commits: #506 (V1-05 landing, keeps the frozen contract commits reachable) and #508 (`main` becomes the v1 line, keeps the v0 history). All later merges are squashes. | #501 |
| GitHub Flow branch model | `main` and `lts/v0` exist. Documents, scripts and old version-line branches still follow the former model. | #501 |
| Specification layout `spec/<slug>/vN/` | Not yet adopted. The v1 specification is one numbered set under `spec/v1/`, with the roadmap inside it. | #504 |
| Roadmap status `released`, `specs` field, root `roadmap.yaml` | Not yet adopted: the status is still `shipped`. | #501, #504 |
| CLI command style recorded in an ADR | Not yet recorded. | #510 |
| Conformance of the v1 code and documents with each selected profile | Not yet audited against v2.3.0. | #511 |
| ADR location | `docs/adr/`. Decisions before 2026-11-11 are `DEC-*` records under `archive/processkit/decisions/`. | — |

## Stronger project requirements

- Every roadmap phase has a runnable demo of the Go executable and matching
  CLI and MCP results (v1 specification, chapters 6 and 14).
- Tool Features install only from producer-documented sources with a pinned
  version and a verified digest per architecture (chapter 8).

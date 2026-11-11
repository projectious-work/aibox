---
title: Architecture
description: "Maintainer architecture of aibox v1, structured after arc42."
weight: 90
---

# aibox v1 architecture

This section explains how aibox v1 is built. It is written for maintainers
and reviewers, not for users who want to run a workspace; for that, read the
[native workspace starters](../native-starters.md) and the
[tool catalog](../tool-catalog.md).

The pages follow the [arc42](https://arc42.org/overview) template and use
[C4 diagrams](https://c4model.com/) for the context and container views.
They describe what exists on the `v1.x-dev` line now. Later roadmap phases
are named where they will change a part, and are marked **Planned**.

| arc42 section | Page | Content |
|---|---|---|
| 1, 2 | This page | Goals, quality goals and constraints |
| 3 | [System context](context.md) | Actors, neighbouring systems, C4 context view |
| 5 | [Building blocks](building-blocks.md) | C4 container view, Go packages, Features, Templates, tests |
| 8 | [Crosscutting concepts](concepts.md) | Ownership, results, path safety, child processes, Feature conventions, testing |
| 9 | [Architecture decisions](decisions.md) | Links to the `DEC-*` decision records |
| 11 | [Risks and technical debt](risks.md) | Known gaps and open risks |

The specification in
[`spec/v1/`](https://github.com/projectious-work/aibox/tree/v1.x-dev/spec/v1)
is the design source; these pages summarize it and show how far the code has
realized it. Where a page and the specification differ, the specification
wins and the page is out of date. Each roadmap phase updates the sections it
affects (chapter 21, R-CODE).

## Goals

aibox v1 makes a Dev Container a complete, AI-ready workspace without owning
the container lifecycle. Its project definition is the standard
`.devcontainer/devcontainer.json`, which stays valid without aibox. aibox
adds three things around it:

1. **Native artifacts:** Dev Container Templates that produce ordinary
   project files, and Dev Container Features that install optional tools.
2. **An inspection and operation interface:** one Go executable, `aibox`,
   reachable as a CLI and as a stdio MCP server, with the same use cases and
   result envelope behind both.
3. **Workspace experience:** tmux layouts, themes, previews and harness
   integration inside the container. This is **Planned** (V1-10 to V1-18).

## Quality goals

| Priority | Quality goal | What it means here |
|---|---|---|
| 1 | Safety of authority | No operation gains host authority from repository content. A process inside the container never gets host lifecycle access. |
| 2 | Native interoperability | The Dev Container CLI and the selected Features own build and lifecycle. aibox never writes a transformed copy of a user-owned definition. |
| 3 | Truthful results | A result reports what was observed, not what a child exit status suggests. Unknown stays unknown. |
| 4 | Reproducibility | Tools, Features and dependencies are pinned and verified; a rebuild resolves the same versions. |
| 5 | Agent usability | CLI and MCP return equivalent, schema-valid results with stable error codes and next actions. |

## Constraints

- **One Go module, one executable.** Go 1.27.0, the official MCP Go SDK as
  the only protocol implementation, standard library first
  ([chapter 7](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/07-architecture.md)).
- **Upstream owns container semantics.** No daemon, container registry,
  scheduler, Dev Container compiler or second resolver in aibox.
- **Supported container targets for tool Features:** Debian 13 (trixie) on
  `amd64` and `arm64` only. Other bases fail at build time with a clear
  error; Arch and Nix are not inferred.
- **Phased delivery.** Each phase starts from a frozen, reviewed contract and
  separate acceptance-test and implementation tracks
  ([chapter 21](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/21-engineering-practice.md), R-CYCLE).
- **Code style.** `gofmt` and `go vet` are authoritative; Hexagonal and
  Clean Architecture guide package boundaries (R-CODE).

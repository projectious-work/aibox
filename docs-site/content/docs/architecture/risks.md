---
title: Risks and technical debt
description: "Known gaps between the aibox v1 architecture and its specification, and open risks."
weight: 50
---

# Risks and technical debt

This page lists where the code differs from the
[chapter 7](https://github.com/projectious-work/aibox/blob/main/docs/product/v1/07-architecture.md)
architecture and which risks remain open at V1-06. Remove an entry when a
phase closes it.

## Technical debt

| Item | Detail | Owner |
|---|---|---|
| Package names differ from chapter 7 | The use case lives in `internal/usecase`, not `internal/app`; the MCP adapter is `internal/mcpserver`, not `internal/mcp`; settings live in `internal/settings`, not `internal/config`; CLI parsing sits in `cmd/aibox/cmdline.go` instead of `internal/cli`. Chapter 7 says its responsibilities apply independently of package names, so this is not a defect, but the table there and the code must be read side by side. | Next phase that adds a CLI command |
| Few ports | V1-06 adds the ports `usecase.Inspector`, `usecase.Events`, `usecase.SinkHealth`, `settings.FileSystem` and `logging.FileAppender`, where tests needed a replaceable adapter. The project layer still reads the file system through `internal/project` directly. Ports for the Dev Container CLI and a runtime come with the operations that need them. | V1-08 |
| No import-boundary test | Chapter 7 asks for a package-level test that rejects forbidden imports, and static checks against a second Dev Container resolver, daemon or broad provider SDK. Neither exists; the import graph is small and is checked by review. | V1-08, V1-09 |
| Primitives without a caller | `internal/process` and `internal/operation` are tested but no use case uses them. Executable resolution is not bound to an operator policy, so the runner must not run mutating lifecycle work. | V1-08 |
| Narrow input manifest | The digest covers only `.devcontainer/devcontainer.json`. Dockerfiles, Compose files and other inputs need upstream-assisted discovery before a digest can bind a mutation. | V1-08 |
| Library copied into each Feature | Each Feature carries a generated copy of `features/lib/aibox-feature.sh`. `scripts/v1-05-features.sh check` fails on drift, so the copies cannot diverge silently, but a library fix changes all 117 Features. | Ongoing |
| Unix only | The lock, operation record store, runner, executable resolver and log file sink fail closed on platforms other than Linux and macOS (the file sink then counts as degraded). | V1-22 |
| Sequence gaps under backpressure | An event dropped from a full queue has already taken its sequence number, and a `logging.sink_degraded` event for a write failure is written directly to stderr, possibly before queued events of lower sequence. Normal runs are unaffected; LOG-033 and LOG-034 are verified at component level. | Observe |
| Hand-written YAML and settings checks | The YAML encoder and the settings-file schema check are hand-written ([ADR 0003](https://github.com/projectious-work/aibox/blob/main/docs/adr/0003-hand-written-yaml-encoder.md)). Golden tests and drift tests against the frozen schemas cover them; the acceptance suite parses YAML with a YAML 1.2 parser. | Ongoing |

## Risks

| Risk | Effect | Mitigation now |
|---|---|---|
| Features and Templates are not published | Users cannot reference `ghcr.io/projectious-work/aibox-features` yet; the documented references are planned. | Acceptance publishes to a local registry; publication is a separate release gate. |
| Incomplete acceptance evidence | AC-TOOLS acceptance needs a passing run on `debian-13-amd64` and `debian-13-arm64`. Full and sampled runs passed on `debian-13-arm64` under rootless Podman. A sampled run passed on `debian-13-amd64` under Docker, in an OrbStack machine with x86-64 translation on Apple silicon, so it does not cover CPU-specific behaviour; it was completed over three resumed runs because Docker stopped answering under translation. | Component tests cover every guard offline; the runner records evidence per run. |
| Sampled error scope | The default AC-TOOLS run builds only a sample of the expected-failure cases, so a target needs 184 builds instead of 466. A Feature outside the sample has no real-container failure evidence. | `AC-TOOLS/guards` runs the component tests for every Feature; `--scope full` remains available. |
| Podman layer cache | Feature install steps never hit the cache through the pinned CLI on Podman, so runs are slow. A small Podman 5.4.2 build reused a `RUN --mount` layer built from other content. | The runner fails a build that takes a Feature step from the cache. Parallel builds share one store and remove their images only at the end of the run, so one build cannot remove layers another has matched; that race is reasoned, not observed. |
| Upstream tool sources | A producer can withdraw an asset or change its output format. | Exact pins, digests or signatures, a frozen dependency resolution date, and version checks that match a version token rather than exact output. |
| Supply-chain review | Chapter 13 requires digest, licence, provenance and vulnerability checks before a release. | The catalog records each tool's licence, install source, documentation URL and verification method; release-level review is **Planned** with the release phases (V1-19, V1-23). |
| Authority model not yet exercised | Operator policy, local versus operator mode and denial before side effects are specified but not implemented. | The executable exposes only a read-only, closed-world operation; `--context` accepts only `local`. |
| Settings on macOS not qualified | The macOS system settings path and terminal check are specified and unit-tested, not run on a macOS host. | Qualification is roadmap phase V1-22. |
| Log files on network file systems | `flock` and append semantics are only guaranteed on local file systems; concurrent writers on NFS could interleave or lose lines. | Documented as unsupported (logging.fileSink.scope). |

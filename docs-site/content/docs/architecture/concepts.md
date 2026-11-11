---
title: Crosscutting concepts
description: "Rules that apply across aibox v1 packages, Features and tests."
weight: 30
---

# Crosscutting concepts

These concepts shape more than one building block. Each names the
specification chapter that defines it.

## One owner for each setting

Every setting has exactly one authoritative input and one interpreter
([chapter 8](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/08-native-artifacts.md#one-owner-for-each-setting)).
Image, user, mounts, ports and lifecycle belong to `devcontainer.json` and
the Dev Container CLI. Tool selection, versions and removal belong to the
`features` object. `customizations.aibox` carries UX preferences only and
never installs a tool. aibox reads only the fields it needs and never writes
a transformed replacement of a user-owned file.

A consequence runs through the code: aibox keeps no parallel desired-state
model. The requested definition is the project's native files; runtime state
is whatever an inspection of the runtime shows. Today no runtime is
consulted, so inspection reports `state: unknown`.

## Result envelope and exit codes

Every operation returns one envelope, `aibox.operation-result/v1`, defined in
[`internal/contract`](https://github.com/projectious-work/aibox/blob/v1.x-dev/internal/contract/result.go)
and `spec/v1/operation-result.schema.json`. It carries the request ID,
attributed actors, target, input digest, outcome, changed resources,
warnings, evidence, and either data or a structured error with a closed
`nextAction`. Adapters validate the envelope before they emit it.

The outcome records what was observed, not the child's exit status. The CLI
maps it to the stable exit classes of
[chapter 18](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/18-operational-contracts.md):
0 succeeded or no change, 2 invalid input, 3 denied, 4 missing dependency or
incompatible, 6 rebuild, restart or recovery required, 124 timed out,
130 cancelled, and 5 for any other failure. CLI and MCP produce equivalent envelopes for equivalent
requests; black-box tests compare them.

## Read-only by default, authority from outside the project

The executable registers only complete operations. In V1-05 that is one
read-only use case, so neither adapter can mutate anything. Actor labels
such as `cli` or `mcp` are attribution, never authorization. Repository
content can request, but never grant, host authority; operator policy will
live outside the project
([chapter 13](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/13-security.md)).
Policy loading and authorization are **Planned** (V1-08).

## Path containment and input digest

[`internal/project`](https://github.com/projectious-work/aibox/tree/v1.x-dev/internal/project)
resolves a canonical, symlink-free project root and refuses files whose
resolved path leaves it, including through a sibling with the same prefix.
The input manifest hashes an explicit list of control files in stable path
order, optionally bound to policy and executable digests. The digest detects
a changed declaration; it does not claim that the build context is
immutable. Path checks cannot remove symlink races, so a mutating caller
must revalidate under the operation lock.

## Child processes without a shell

[`internal/process`](https://github.com/projectious-work/aibox/tree/v1.x-dev/internal/process)
runs an already-resolved absolute executable with an argument vector, an
explicit working directory and a complete environment; nothing is inherited.
Each stream returns a bounded tail (256 KiB) with known secrets masked. On
Linux and macOS cancellation sends TERM to the whole process group and KILL
after five seconds; other platforms fail closed. `ResolveExecutable` accepts
only an operator-supplied path outside the project, refuses writable files
or parent directories, and binds the program to a byte digest.

## Locks, operation records and interruption

[`internal/operation`](https://github.com/projectious-work/aibox/tree/v1.x-dev/internal/operation)
provides the primitives for the mutating state machine
`validated → authorized → locked → revalidated → executing → inspected →
recorded`
([chapter 7](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/07-architecture.md#resource-and-operation-state)).
An OS advisory lock is keyed by project, configuration and runtime endpoint.
Operation records are small, non-secret and schema-shaped; they are written
by sync and atomic rename into a private directory and only move forward.
After a cancellation, only proven absence of every effect yields `cancelled`
or `timed_out`; anything uncertain is `partial`. These primitives are tested
but not yet wired to a use case.

## Tool Feature conventions

Every aibox Feature follows the
[Tool Feature interface](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/08-native-artifacts.md#tool-feature-interface)
through the shared POSIX `sh` library `features/lib/aibox-feature.sh`:

- **Guards before change.** `af_init` checks, in order, the target
  (Debian 13 on `amd64` or `arm64`), the `version` option, a repeated
  installation, `requires` and version coupling. A Debian Feature then
  checks package candidates before it installs anything.
- **One error line.** A failure prints
  `aibox-feature error AIBOX_FEATURE_<CODE> feature=<id> <detail>` to stderr.
  The eight codes are fixed by the catalog: `UNSUPPORTED_TARGET`,
  `UNSUPPORTED_VERSION`, `VERSION_CONFLICT`, `COUPLING_MISMATCH`,
  `DEPENDENCY_MISSING`, `DOWNLOAD_FAILED`, `VERIFICATION_FAILED` and
  `PACKAGE_UNAVAILABLE`.
- **Pinned and verified inputs.** Downloads use HTTPS only and are checked
  against a SHA-256 or SHA-512 digest, an npm integrity value, or a
  minisign signature. Debian packages come from the base image's snapshot,
  with `--no-install-recommends`. Python and npm dependency resolution is
  frozen at the catalog research date. No option accepts `latest`, a range
  or a URL.
- **State marker.** A successful install records its version under
  `/var/lib/aibox-features/<id>/version`. Later Features read it for
  `requires` and coupling checks; a second install with another version
  fails with `VERSION_CONFLICT`.
- **No leftovers.** Helper packages, temporary APT sources and scratch
  directories are removed on success and on failure.

## The catalog is the single source

The interface catalog drives Feature metadata, guard data, the published
dependency references, the acceptance test plan and the user's
[tool catalog](../tool-catalog.md) page. Generated files carry a "do not
edit" header, and `scripts/v1-05-features.sh check` fails on drift. When
prose and the catalog disagree, the catalog wins and the difference is a
contract defect.

## Testing approach

Default tests are offline. Go tests use temporary roots and a helper child
process; Feature component tests replace `apt-get`, `dpkg`, `curl` and the
filesystem root with fakes and run the real `install.sh` and library. Only
the acceptance suites and host verifiers need a container engine, and they
build disposable workspaces from pinned fixtures with a checksummed Dev
Container CLI. Test style is defined in
[chapter 21](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/21-engineering-practice.md)
(R-TESTSTYLE).

## Phase contracts and evidence

Each phase freezes a reviewed contract, then runs acceptance tests and
implementation on separate tracks. A defect found later is fixed by a
contract amendment with an owner decision (R-CYCLE). The phase's dev note
under `dev-notes/` records choices, evidence and gaps. `done` in the roadmap
means verified locally; `shipped` needs published release evidence.

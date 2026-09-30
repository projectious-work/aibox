# 14. Verification, target qualification and release

This chapter implements R-TARGETS and R-DEPENDENCIES and gives the roadmap
an evidence-based definition of `shipped`. Each phase includes implementation,
AC/ledger trace, tests, its roadmap `docs` deliverable and a passing local
documentation build from the same commit. Unit-test success alone is insufficient.

## Runnable increments and frequent preview checkpoints

V1-03 delivers the first Go `aibox` executable: version/help plus safe,
read-only `inspect --context local` and `inspect_workspace` MCP journeys on the minimal native example. They
must produce a schema-valid result without Docker, operator policy or host
credentials. It does not expose unfinished lifecycle mutation or claim v0
parity. V1-09 expands the guarded MCP registry and thin CLI surfaces to the
authorized operator use cases.

For every phase from V1-03 onward, the [roadmap's phase demo](roadmap.yaml)
names a product-binary CLI command, an MCP tool, and observable, feature-specific
outcomes. Each implementation phase proves equivalent CLI/MCP operation results
and refusal behavior for its new shared use cases. The foundation phases V1-01
and V1-02 have no executable surface to compare.
The implementation change builds the Go executable, runs that phase's demo
against a disposable fixture, keeps earlier supported demo journeys working,
and builds the matching v1 Hugo documentation from the same commit. The
customer walkthrough states prerequisites, exact command, expected result,
current limitations and a cleanup or reset path. A package-only test, static
page, version-only smoke test or screen recording cannot replace the executable
journey. Where a capability needs a host runtime, credential or network,
default CI still runs a deterministic fake/offline binary journey; an opt-in
host gate proves the real claim before it is advertised.

The phase development note records the demo command, source SHA, built-binary
digest, fixture, observed output or redacted transcript, documentation paths,
Hugo build result, and failed/skipped checks. Demo readiness is an integration
checkpoint, not a new roadmap status and not `shipped`. The dev branch remains
buildable and demonstrable after every integrated increment; a later change
that breaks an earlier supported journey blocks integration.

Work in small, reviewable increments and prepare alpha/beta candidates at
verified checkpoints rather than waiting for all v0 parity. A public preview
has an explicit supported-scope and target matrix: run the complete applicable
test, security, documentation, host and artifact gates for what it claims;
label everything else planned, unavailable or unsupported. Publish only from
the accepted exact candidate through the protected version-line path, verify
downloaded artifacts and matching documentation, and keep v0 current. A
locally working demo never bypasses owner approval, release identity,
provenance, or post-publication checks. A phase becomes `shipped` only when its
named release and conformance evidence support that status.

## Test layers and evidence format

| Layer | Required fixtures / oracle | Release relevance |
|---|---|---|
| Schema/static | Closed JSON schemas, valid/invalid native examples, import-boundary rules, source/ledger row counts | Reject incompatible config and duplicate engines early |
| Unit/component | Fake child executables, policy and path checker, result serializer, guide index, local renderer | All error branches, stable codes, no side effects on denied/read-only requests |
| Binary black box | Built Go CLI and stdio MCP client from V1-03 with matching result data and errors, invalid input and output separation | Earlier CLI/MCP journeys remain working; every phase's adapters return equivalent results |

| Settings/logging black box | Temporary system/user/project/env-file/env/flag layers; stderr/file/collector capture, disk-full and rotation fixtures | Provenance, source authority, redaction and evidence separation under real binary execution |
| Disposable integration | Pinned Dev Container CLI, selected runtime, Features, Templates, image, local UX assets | Direct upstream and wrapped lifecycle, exact stop/remove, persisted home |
| User journey | Chapter 20's human and agent stories: fresh minimal project, complex migrated v0 project, local agent customization, operator lifecycle and recovery, sidebar/review | Actual v0 parity and new feature usability |
| Security/fault injection | Host hook denial, symlinks, token leaks, race/restart/partial failure | Boundary holds under hostile and interrupted inputs |
| Published-artifact | Download exact candidate binary/image/Feature/Template, verify hashes/signatures/SBOM, run fresh journey | Source-tree success cannot substitute for release evidence |

The v1 test gate uses Go tests, race detection, vet, contract/schema checks,
CLI/MCP binary journeys and documentation builds. Rust `cargo test`, Clippy,
and Cargo audit are not v1 gates; the retained v0 Rust line keeps its own gates.

Each phase evidence record contains phase ID, source commit, candidate digests,
target cell, dependency versions, command with bounded/sanitized arguments,
start/end time, actor, result, attached redacted logs and failed/skipped
reasons. It also contains documentation paths/audiences, the local Hugo build
command and result, link/schema/example checks, and confirmation that v0
remains current. A skipped test is never silently counted as passed. Rerun
after any candidate-changing fix. Evidence is immutable once cited by a release;
supersession creates a new record. Public release notes state capability and
support limits, not internal tokens/host paths.

## Target/dependency qualification matrix

`04-reuse-targets.md` defines T01–T09 and the dependency classes. For each
implementation phase, record its dependencies in an exact manifest with Go, MCP SDK,
Dev Container CLI/Node, runtime/Compose provider, base image, Features,
Templates, tmux/PowerKit/Yazi/preview tools, harnesses, processkit, sidebar
and review tools. Each entry records version, digest, publisher, license,
supported architectures, update owner and the test that exercises it.
Version floors/ceilings are a separate supported-range claim; the tested
manifest is a single reproducible set. A release candidate may not use an
unpinned `latest` image or unreviewed Feature ref.

Binary installation, version check, update and uninstall use a documented
standard package/distribution route. A separate local installer may be kept
only if that route cannot preserve the v0 workflows; the MCP server never
self-updates. Release tests perform check/update/dry-run/uninstall against a
disposable prefix, verify completion/help/version output and prove purge
does not happen on ordinary uninstall. This closes F09 without putting a
self-modifying control plane in the Go core.

Qualification covers macOS arm64 with OrbStack and Docker Desktop (T01/T02),
macOS amd64 compatible Docker (T03), Linux amd64/arm64 Docker Engine (T04/T05),
the supported Podman/Compose cells (T06), and Debian container-local amd64/
arm64 behavior (T07). T08 SSH is an operator journey on a prepared machine,
not remote provisioning. T09 is explicitly outside support. For each
advertised cell run native direct build/up/exec, wrapper build/up/stop/remove,
rebuild with state retention, Feature catalog sample, tmux/theme/preview,
local customization denial tests, and cleanup. Audio and browser/LaTeX tests
run where platform capabilities exist; publish unavailable cells and their
reason rather than claiming a full Cartesian product. Any v0-supported
cell that cannot be reproduced blocks v0 replacement until equivalent support
is demonstrated.

Host gates are executed by a human-authorized operator on a disposable,
identified host. The handoff lists exact repo commit, candidate digests,
target runtime/context, expected effects, retention/cleanup and credentials
needed. The agent may prepare instructions and inspect returned evidence; it
must not smuggle host execution through an in-container socket or companion.
Repeatability requires a second clean-run on at least the primary host cell
and re-run of failure-sensitive gates after changed candidate bits.

## Roadmap state and publication

The canonical [roadmap](roadmap.yaml) groups twenty-three content-specific phases
with dependency edges and spec references. Its standard `status` values remain
`idea`, `planned`, `in_progress`, `shipped`, and `cancelled`. The project-specific
`implementationStatus: done` marks completed, verified code and documentation
with a source commit, development note, and local evidence. A phase can retain
`status: in_progress` while publication remains pending. `shipped` additionally
requires a named release and published-artifact evidence; merging code alone
does not meet it. V1-02 is a foundation prerequisite for all implementation
phases, so the shared Hugo build exists before features are delivered. Every
phase has a distinct executable demo and documentation deliverable; V1-19
hardens the later cross-line alpha/beta archive matrix, not the first
opportunity to publish a scoped, verified preview. Keep
development notes adjacent to implementation under the normal repository
documentation structure. The phase graph is not a
single linear sprint plan: independent workspace slices may proceed once
their declared prerequisites exist. Required parity remains applicable regardless
of phase ordering.

Before stable-v1 promotion as a verified v0 replacement, complete all
F01–F48, N01/N02 and AC groups, reconcile
every `CFG:*`, `ADDON:*`, `CMD:*` and runtime-asset row, complete chapter 6's
technical qualification checks, and verify the full target/dependency matrix
and specification baseline. Build and validate documentation
locally; publish version-aligned README, docs, roadmap, changelog, AI-readable
discovery and read-only product guidance from the exact release candidate.
Chapter 15 and AC-DOCS additionally require the v1 site to use the v0 Hugo
brand-theme pipeline, with v0 retained as the current release through alpha
and beta publication. Build and publication run locally, with no
project-authored GitHub Actions.
Only after downloaded artifacts pass fresh-machine verification may v1 be
called a v0 replacement. An incomplete preview is labeled preview with
precise missing capabilities and does not erase the v0 support path.

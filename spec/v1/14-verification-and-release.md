# 14. Verification, target qualification and release

This chapter implements R-TARGETS and R-DEPENDENCIES and gives the roadmap
an evidence-based definition of `shipped`. A specification phase is not a
product release. Every content phase has a reviewed plan, implementation
diff, AC/ledger trace, tests, its roadmap `docs` deliverable, a passing local
documentation build from the same commit, and independent complete-baseline
conformance review. This includes internal and enabling phases. No phase is
marked shipped from a green unit test alone.

## Test layers and evidence format

| Layer | Required fixtures / oracle | Release relevance |
|---|---|---|
| Schema/static | Closed JSON schemas, valid/invalid native examples, import-boundary rules, source/ledger row counts | Reject incompatible config and duplicate engines early |
| Unit/component | Fake child executables, policy and path checker, result serializer, guide index, local renderer | All error branches, stable codes, no side effects on denied/read-only requests |
| Binary black box | Built Go CLI and stdio MCP client, same requests/policy, invalid input, cancellation | Machine output purity and adapter-equivalent results |
| Settings/logging black box | Temporary system/user/project/env-file/env/flag layers; stderr/file/collector capture, disk-full and rotation fixtures | Provenance, source authority, redaction and evidence separation under real binary execution |
| Disposable integration | Pinned Dev Container CLI, selected runtime, Features, Templates, image, local UX assets | Direct upstream and wrapped lifecycle, exact stop/remove, persisted home |
| User journey | Fresh minimal project, complex migrated v0 project, local agent customization, sidebar/review, recovery | Actual v0 parity and new feature usability |
| Security/fault injection | Host hook denial, symlinks, token leaks, race/restart/partial failure | Boundary holds under hostile and interrupted inputs |
| Published-artifact | Download exact candidate binary/image/Feature/Template, verify hashes/signatures/SBOM, run fresh journey | Source-tree success cannot substitute for release evidence |

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

`04-reuse-targets.md` defines T01–T09 and the dependency classes. Before the
first implementation plan, choose an exact manifest with Go, MCP SDK,
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
cell that cannot be reproduced requires owner-approved retirement or a
blocking implementation issue.

Host gates are executed by a human-authorized operator on a disposable,
identified host. The handoff lists exact repo commit, candidate digests,
target runtime/context, expected effects, retention/cleanup and credentials
needed. The agent may prepare instructions and inspect returned evidence; it
must not smuggle host execution through an in-container socket or companion.
Repeatability requires a second clean-run on at least the primary host cell
and re-run of failure-sensitive gates after changed candidate bits.

## Roadmap state and publication

The canonical [roadmap](roadmap.yaml) groups twenty-three content-specific phases
with dependency edges and spec references. `planned` means defined, not
approved; `in_progress` needs a reviewed phase plan; `shipped` requires a
`devNote` path and release identifier validated by the roadmap schema, plus
the evidence above. V1-20 is a foundation prerequisite for all implementation
phases, so the shared Hugo build exists before features are delivered. Every
phase has a distinct documentation deliverable and build gate; V1-21 is for
candidate-specific snapshot/publication, not a catch-up writing phase. Keep
development notes adjacent to implementation under the normal repository
documentation structure. The phase graph is not a
single linear sprint plan: independent workspace slices may proceed once
their declared prerequisites exist. A phase cancellation must name the
owner-approved replacement or explicit v0 parity exception.

Before promotion, complete all F01–F48, N01/N02 and AC groups, reconcile
every `CFG:*`, `ADDON:*`, `CMD:*` and runtime-asset row, close G01–G11,
verify the complete target/dependency matrix, and run an independent review
against the **whole** accepted baseline. Build and validate documentation
locally; publish version-aligned README, docs, roadmap, changelog, AI-readable
discovery and read-only product guidance from the exact release candidate.
Chapter 15 and AC-DOCS additionally require the v1 site to use the v0 Hugo
brand-theme pipeline, with v0 retained as the current release through alpha
and beta publication. Use the company fast-forward promotion path and squash
topic PRs; no unique commits on promotion branches and no project-authored
GitHub Actions.
Only after downloaded artifacts pass fresh-machine verification may v1 be
called a v0 replacement. An incomplete preview is labeled preview with
precise missing capabilities and does not erase the v0 support path.

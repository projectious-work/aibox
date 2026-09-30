# V1-03 — Go operation core

V1-03 implementation is done with local evidence; publication is pending.
This note is for Go implementers,
security reviewers, and independent conformance reviewers. The accepted
[architecture](../spec/v1/07-architecture.md),
[lifecycle](../spec/v1/09-lifecycle.md), and
[operational contracts](../spec/v1/18-operational-contracts.md) are the baseline.
The initial slice was merged through
[PR #467](https://github.com/projectious-work/aibox/pull/467).
The CLI/MCP completion slice was merged through
[PR #490](https://github.com/projectious-work/aibox/pull/490). Its verified
source is `05dd13b1556ff300cac77f08fbd628f966e807bc`; the
[local evidence record](../spec/v1/evidence/V1-03-05dd13b1/evidence.json)
and [MCP client transcript](../spec/v1/evidence/V1-03-05dd13b1/mcp-parity.txt)
are preserved in this repository. The roadmap marks implementation `done` and
keeps delivery `in_progress` until a named release has published-artifact evidence.

## Implemented and boundaries

The initial Go slice adds a common result envelope and CLI exit classes,
canonical project-file containment, and a no-shell child runner with an
explicit environment and bounded, known-secret-redacted diagnostic tails.
The module pins Go 1.27.0 and the official MCP Go SDK. The accepted
design uses one Go core and delegates container semantics to the upstream
Dev Container CLI; this slice does not create a second container engine.

There is now a source-built v1 CLI with version/help and read-only local
inspection through both CLI and a local stdio MCP server; there is no installable
release or supported lifecycle operation.
The result `data` remains opaque at the common-contract layer, so envelope
validation alone does not prove operation-specific schema conformance. A
follow-on slice added Linux/macOS process-group TERM/KILL cancellation and a
private, exact-environment OS advisory lock. A private durable operation record store
now writes bounded, schema-shaped records by file sync, atomic rename and
directory sync; a caller must hold the exact-environment lock. An offline
post-cancellation classifier now requires exact-effect observation before
reporting harmless cancellation. A resolver can now bind an operator-supplied
executable path outside the project to a trusted file and byte digest. The
resolver is not yet wired to an operator policy or runtime-specific observer,
so the runner MUST NOT run mutating lifecycle work.
Path checks do not remove symlink races; callers must revalidate under a lock.

## Validation and documentation

For an exact candidate checkpoint, run `./scripts/verify-v1-03.sh` from the
repository root. It builds Linux and macOS arm64 Go binaries, executes the
minimal customer demo, validates the actual JSON against the closed result
schema, runs Go tests/race/vet and specification checks, builds the v1 Hugo
site, tests a real stdio MCP client against the same inspection use case, and
prints the source SHA, binary digests and preserved evidence path. This is a
reproducible review gate, not a public prerelease or a claim that host mutations work.

The initial candidate passed `go test ./...`, `go test -race ./internal/...`,
`go vet ./...`, specification validation, and a local Hugo build. Tests cover
literal argv, empty inherited environment, bounded/redacted output, child
failure, cancellation, symlink escape, and sibling-prefix rejection. Re-run
these checks for each new candidate; earlier results do not validate later
changes. On 2026-09-28, Go declarations and safety-critical sections gained
comments following [official Go guidance](https://go.dev/doc/comment).
The follow-on cancellation and lock slice passes `go test ./...`, race tests,
`go vet ./...`, macOS arm64 cross-build, the specification and contract
validators, and a local v1 Hugo build. New offline tests cover descendant
termination, TERM-to-KILL escalation, lock contention, endpoint separation,
private lock files, and rejection of unsafe stores.
The operation record slice adds offline tests for atomic replacement, private storage,
invalid records, traversal and symlink rejection. It does not yet expose a
CLI for reading operation records or imply that an interrupted mutation can be replayed.
Operation record updates now reject reused IDs with changed request identity, skipped
authorization, backward state transitions and edits to terminal records.
The interruption slice tests that confirmed or uninspectable effects yield
`partial`; only proven absence yields `cancelled` or `timed_out`.
The executable-resolution slice tests digest binding, rejection of project
programs and refusal of writable files or parent directories.
The input-manifest slice hashes explicit, confined native control files in
stable path order and binds optional policy/executable digests. The read-only
CLI now uses this shared builder for its single declared native entry point;
the broader lifecycle dependency set still needs upstream-assisted discovery.

The [architecture status](../spec/v1/07-architecture.md) and
[public preview walkthrough](../docs-site/content/docs/core-status.md) document
these boundaries. The first offline executable demo builds `cmd/aibox`, then
inspects `spec/v1/examples/minimal`; its successful result binds the native
declaration digest and reports runtime state as unknown. No host authority or
v1 runtime compatibility claim is exposed by this slice.

## Completion boundary and later phases

V1-03's executable acceptance boundary is the source-built read-only Go CLI
and stdio MCP tool: version/help, local inspection of the minimal native example,
schema-valid equivalent results, invalid-input refusals with no host effect,
offline black-box tests,
the supporting safe operation primitives, and a successful build of the v1
documentation. `./scripts/verify-v1-03.sh` records the exact source commit,
binary digests, both observed journeys, specification and contract checks, and
documentation build in its evidence directory. Publication remains a separate
gate before the roadmap can call the phase `shipped`.

The following work belongs to later capability phases under the canonical
[verification plan](../spec/v1/14-verification-and-release.md):

- V1-08 owns operator policy loading and authorization, runtime-specific effect
  inspection, complete upstream-assisted input discovery, recovery reads for
  mutating operation records, and verified Dev Container delegation. Its host
  mutations stay unavailable until these checks are in place.
- V1-09 expands the fixed SDK registry, operator capabilities, and complete
  operation-specific conformance corpus. CLI/MCP equivalence starts at V1-03
  and remains a gate for each executable roadmap phase.
- V1-21 owns host trust-boundary qualification. The existing executable
  resolver and input-manifest builder are foundations for that work, not
  grants of host authority.

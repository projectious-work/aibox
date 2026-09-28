# V1-03 — Go operation core

V1-03 is the current implementation focus. This note is for Go implementers,
security reviewers, and independent conformance reviewers. The accepted
[architecture](../spec/v1/07-architecture.md),
[lifecycle](../spec/v1/09-lifecycle.md), and
[operational contracts](../spec/v1/18-operational-contracts.md) are the baseline.
The initial slice was merged through
[PR #467](https://github.com/projectious-work/aibox/pull/467).

## Implemented and boundaries

The initial Go slice adds a common result envelope and CLI exit classes,
canonical project-file containment, and a no-shell child runner with an
explicit environment and bounded, known-secret-redacted diagnostic tails.
The module pins Go 1.27.0 and adds no third-party dependencies. The accepted
design uses one Go core and delegates container semantics to the upstream
Dev Container CLI; this slice does not create a second container engine.

There is now a source-built v1 CLI with version/help and read-only local
inspection; there is no installable release, MCP server, or supported lifecycle operation.
The result `data` remains opaque at the common-contract layer, so envelope
validation alone does not prove operation-specific schema conformance. A
follow-on slice added Linux/macOS process-group TERM/KILL cancellation and a
private, exact-environment OS advisory lock. A private durable receipt store
now writes bounded, schema-shaped records by file sync, atomic rename and
directory sync; a caller must hold the exact-environment lock. An offline
post-cancellation classifier now requires exact-effect observation before
reporting harmless cancellation. No runtime-specific observer or policy-owned
executable resolver is wired yet, so the runner MUST NOT run mutating lifecycle work.
Path checks do not remove symlink races; callers must revalidate under a lock.

## Validation and documentation

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
The receipt slice adds offline tests for atomic replacement, private storage,
invalid records, traversal and symlink rejection. It does not yet expose a
receipt-reading CLI or imply that an interrupted mutation can be replayed.
Receipt updates now reject reused IDs with changed request identity, skipped
authorization, backward state transitions and edits to terminal records.
The interruption slice tests that confirmed or uninspectable effects yield
`partial`; only proven absence yields `cancelled` or `timed_out`.

The [architecture status](../spec/v1/07-architecture.md) and
[public preview walkthrough](../docs-site/content/docs/core-status.md) document
these boundaries. The first offline executable demo builds `cmd/aibox`, then
inspects `spec/v1/examples/minimal`; its successful result binds the native
declaration digest and reports runtime state as unknown. No host authority or
v1 runtime compatibility claim is exposed by this slice.

## Open V1-03 work

V1-03 still owns policy-controlled executable resolution, complete input manifests,
receipt recovery reads, runtime-specific effect inspection,
operation-specific checks, fake Dev Container delegation, CLI/MCP equivalence,
and offline contract fixtures. Later V1-08 lifecycle and V1-21 trust-boundary
work depend on these foundations.

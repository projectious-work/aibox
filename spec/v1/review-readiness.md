# PR #463 implementation-readiness review

Reviewed 2026-09-26–27 as an integrated baseline through original PR head
`34cc9381f71858f0107cc1ab97801a881b81832c`, not just the latest commit.
This report accompanies the corrective commit; it does not mark implementation
or owner acceptance complete.

## Verdict

The incoming draft was substantial but **not ready to hand to implementers
without interpretation**. Its schema could not represent useful inspection,
doctor or log results, required a container identity before a first build,
and left important configuration/authority/recovery decisions implicit.

The revised baseline is a software implementation specification: product
requirements and exclusions, component boundaries, native artifact ownership,
operation algorithms and wire contracts, failure behavior, security model,
configuration/logging, migration, measurable limits, acceptance trace and a
content-phased delivery plan. It is suitable for **owner acceptance and
phase-by-phase implementation planning**. Once accepted and its phase plan
independently reviewed, foundation implementation may start. This assessment
does not authorize implementing every later integration immediately.

The remaining exact Feature/package selections, native harness/gateway policy
translations and processkit API pins belong to the affected phase's reviewed
implementation plan. Their functional requirements and failure behavior are
fixed here; no implementer may silently weaken parity. Platform and supply-chain
qualification require real artifacts/hosts. Neither category should be
misrepresented as already completed by a documentation PR.

## Findings and corrections

| Severity | Incoming gap | Resolution |
|---|---|---|
| Blocker | Common request fields allowed irrelevant inputs and required a digest before discovery | Closed per-operation union; read requests bootstrap digest; mutation-only requirements and bounded selectors. |
| Blocker | Result envelope had no operation data and required a resource before one existed | Typed data for builds, inspection, findings, logs, refresh, migration and receipts; early errors allow absent target/digest. |
| Blocker | Server policy lacked an implementable identity/approval contract | Host-owned policy schema, OS/stdio principal, exact expiring grants, immutable capability registry and explicit same-principal trust limitation. |
| Blocker | Lost-session recovery was an aspiration | Durable receipt schema, filesystem/lock algorithm, replay rules and CLI/MCP operation inspection; no new background orchestrator. |
| Blocker | Local UX namespace included preview host-network/service controls | Removed those fields from generated schema; native Compose and external approval own them. Negative schema fixtures enforce the boundary. |
| Blocker | Field placement remained a future review | Chapter 19 fixes native/Feature/runtime authority, optional UX compatibility namespace, default/alias rules and conversion destinations. |
| Warning | Generic installer alternatives could hide a new addon system | Concrete single-tool Feature fallback, native installer reuse, explicit version/presence semantics and per-tool independent plan review before coding. |
| Warning | Rebuild implied reliable candidate reuse/atomic replacement | Native preflight may be followed by another build; disclose disruption and partial recovery, record actual resulting image identities. |
| Warning | Compose lifecycle did not define sidecar ownership | Exact primary/exclusively owned service set, shared services excluded, per-ID revalidation, retained volumes/networks and no broad down/prune. |
| Warning | Refresh promised file-set transaction without a recoverable layout | Managed generations, one confined pointer, explicit includes, native overrides retained, interrupted reload evidence and honest process-side-effect limits. |
| Warning | Named-state parity reduced private process context to version control | Private manifest/archive/staging/restore workflow, saved instruction/context files, conflict protection and explicit encryption for credential backups. |
| Warning | Performance and resource behavior were not measurable | R-LIMITS and AC-CONTRACTS specify time, memory/output, queue, cancellation, guide and preview limits with reference-environment tests. |
| Warning | Lightweight checker was not full schema validation | Draft 2020-12 library with format assertions, reference resolution, positive/negative operation/policy/receipt/namespace fixtures and reproducible generation. |
| Warning | Closed schema plus compatible added fields were contradictory | New negotiated wire schema required for field additions; initial v1-only execution contract. |
| Warning | Pending standards merge and nonexistent artifact digests created circular acceptance gates | Exact agent-native draft source pinned; distinguish owner baseline acceptance, independent phase plans and empirical qualification. |
| Warning | Blanket no-host-mount/no-secret-snapshot wording contradicted permitted homes and auth | Explicit scoped native mounts and intentional private credential storage remain allowed; diagnostics/publication remain secret-free. |
| Warning | Roadmap references and PR summary drifted from chapter count/content | Chapters 18/19 linked into the affected stable phases; 23 phases retained, with documentation build/update required in every phase. |

## Industry and company alignment

The review applied requirements-quality criteria of clarity, consistency,
feasibility, traceability and verifiability, using the public
[NASA requirements checklist](https://www.nasa.gov/reference/appendix-c-how-to-write-a-good-requirement/)
and the public scope description of
[ISO/IEC/IEEE 29148](https://standards.ieee.org/ieee/29148/6937/).
This is **not formal ISO certification or a clause-by-clause compliance claim**:
the licensed standard text was not assessed in full.

Company applicability is recorded in chapter 6 against pinned sources.
The decisive distinction is the company's spec-driven cycle: accepting a
specification is not accepting an implementation plan or passing conformance.
This authoring review is not a substitute for the required independent review.
The product boundary is preserved: Go, MCP-first thin CLI, native Dev Container
delegation, reuse-first, all v0 capability parity, local versus host authority,
native home persistence, secrets discipline, and v0-current/v1-prerelease Hugo
documentation updated with each phase.

## Checks and evidence limits

Run the commands in README: reproducible v0 inventory, generated-contract drift,
requirement/acceptance trace, roadmap DAG/docs dependencies, local links, full
schema validation and `git diff --check`. They do not execute project hooks,
start containers or access user credentials.

No actual v1 product, Feature image, MCP-client journey, host target matrix,
Hugo migration or runtime parity has been tested by this change. Those are
implementation deliverables. Cargo test/clippy remain unavailable here under
the owner's previously approved documentation-only exception. Product code
is unchanged. A future implementation must pass its own Go/native tests and
company release/host gates; this exception cannot be inherited as a waiver.

## Start gate and remaining risks

1. Owner accepts this exact revised specification; record its source commit.
2. Independently review the first phase plan, test/evidence matrix and rollback.
3. Establish shared Hugo documentation before product implementation phases,
   then update/build it in each phase as already required by the roadmap.
4. Before catalog, harness, processkit and migration coding, review concrete
   per-tool/provider mappings and fixtures. A missing equivalent blocks that
   adapter or conversion; only the owner can approve a parity scope change.
5. Close G01–G11 with actual evidence before claiming those capabilities or
   replacing v0. Keep unsupported/unverified target cells explicit.

Residual risks are upstream CLI/native runtime behavior, incomplete portable
harness policy APIs, dependency maintenance/licenses, Compose ownership races,
remote-daemon path/UID behavior and realistic provider telemetry availability.
Each has an owning phase and negative/qualification tests. None is resolved
merely by renaming a gate or marking a roadmap phase shipped.

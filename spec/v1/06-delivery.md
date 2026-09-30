# 6. Requirements trace and technical qualification

## Requirement-to-evidence trace

| Requirement | Planned verification | Phase |
|---|---|---|
| R-PRODUCT | AC-PK, AC-GUIDE, minimal standalone journey | V1-01, V1-04, V1-23 |
| R-PARITY | All F/CFG/ADDON/CMDTYPE/ASSET rows | V1-05, V1-10–V1-15, V1-20, V1-23 |
| R-CORE | AC-CLI, AC-BUILD; architecture review rejects duplicate engines | V1-03 |
| R-BOUNDARY | AC-PK, AC-SEC and integration-boundary tests | V1-01, V1-21 |
| R-AUTHORITY | AC-SEC in both contexts | V1-03, V1-09, V1-21 |
| R-LOCAL | AC-SEC, AC-UX, AC-DOCTOR | V1-10–V1-16 |
| R-REUSE | Dependency clearance plus AC-TOOLS | V1-05, V1-22 |
| R-OBSERVABILITY | AC-DOCTOR, AC-CLI, AC-SIDEBAR | V1-11, V1-16, V1-17 |
| R-INDEPENDENCE | AC-BUILD direct upstream journey | V1-04, V1-08 |
| R-CONFIG | AC-CONFIG all field rows | V1-04, V1-05, V1-20 |
| R-PRECEDENCE | AC-CONFIG, AC-UX conflict fixtures | V1-12, V1-20 |
| R-OWNERSHIP | AC-UX, AC-MIG interrupted refresh | V1-10, V1-12, V1-20 |
| R-MIGRATION | AC-MIG and every old command disposition | V1-20 |
| R-RECOVERY | AC-MIG, AC-LIFE | V1-08, V1-20 |
| R-INTERFACES | AC-CLI adapter equivalence | V1-09 |
| R-DOCTOR | AC-DOCTOR plus no-effect tests | V1-16 |
| R-RESULT | AC-CLI, AC-LIFE machine schema fixtures | V1-03, V1-09 |
| R-EXECUTION | AC-SEC, AC-LIFE interruption/race tests | V1-08, V1-21 |
| R-GUIDANCE | AC-GUIDE offline/version/client tests | V1-09, V1-16 |
| R-POLICY | AC-SEC authorization bypass tests | V1-09, V1-21 |
| R-TARGETS | AC-TARGET matrix gates | V1-22 |
| R-DEPENDENCIES | AC-TOOLS, AC-RELEASE artifact manifests | V1-05, V1-22, V1-23 |
| R-DOCS | Phase-scoped AC-DOCS in every phase; each phase's `docs` deliverable and local build evidence, plus candidate-matched previews at verified checkpoints and the full both-line deploy matrix | V1-01–V1-23 (V1-02 enables; V1-19 qualifies the full archive matrix; V1-23 promotes stable docs) |
| R-CONFIG-SOURCES | AC-CONFIG-SOURCES; precedence, source denial and effective-config fixtures | V1-06, V1-09 |
| R-LOGGING | AC-LOGGING; stdout purity, sink/rotation/failure/redaction fixtures | V1-06, V1-21 |
| R-SECRETS | AC-SECRETS; no-secret starter, native transfer modes, host approval, canary and cleanup fixtures | V1-07, V1-21, V1-22 |
| R-LIMITS | AC-CONTRACTS; measured resource, latency, cancellation and output bounds | V1-03, V1-09, V1-21, V1-22 |

N01 → AC-SIDEBAR in V1-17 and N02 → AC-REVIEW in V1-18. Feature rows refine these
groups and link to source evidence.

| v0 feature rows | Primary implementation phase | Cross-cutting proof |
|---|---|---|
| F01–F02 | V1-04, V1-15 | V1-20 migration |
| F03–F05 | V1-08 | V1-09 interface, V1-21 security |
| F06–F08 | V1-20 | V1-08 exact runtime identity |
| F09 | V1-23 | V1-09 CLI distribution |
| F10–F13 | V1-04–V1-05 | V1-20 mapping, V1-22 target qualification |
| F14–F17 | V1-10, V1-15 | V1-21 authority |
| F18–F22 | V1-11 | V1-12 theme, V1-16 doctor |
| F23–F25 | V1-11 | V1-17 status integration |
| F26–F29 | V1-12 | V1-20 alias conversion |
| F30–F38 | V1-13 | V1-22 platform coverage |
| F39 | V1-14 | V1-22 audio cells |
| F40–F42 | V1-05 | V1-22 tool qualification |
| F43–F44 | V1-16 | V1-09 machine interface |
| F45–F46 | V1-22 | V1-23 support claim |
| F47 | V1-23 | V1-16 guide content |
| F48 | V1-16 | V1-17 telemetry truthfulness |

Documentation is a completion condition for **every** phase, including internal
phases, not a final writing pass. Each roadmap entry names its `docs` output;
its candidate must pass a local site build and changed-page checks. V1-02 is
placed in foundation before implementation to establish the shared shell;
verified earlier checkpoints may publish scoped previews with matching docs,
V1-19 verifies the full alpha/beta snapshot and navigation matrix, and V1-23
verifies stable publication.

## Technical qualification checks

| ID | Question to close | Required evidence |
|---|---|---|
| G01 | Verify effective configuration mapping | Chapters 8/12/19 fix authority, destination rules and converter behavior; V1-04, V1-05 and V1-20 enumerate concrete adapters and supply cross-field/default fixtures and native-resolution proof before completion. |
| G02 | Qualify Feature selection and processkit API | Chapter 8 fixes selection criteria; V1-05 and V1-15 must fill per-tool manifest, license/provenance and supported processkit API/version evidence. |
| G03 | Prove snapshots, recovery and audio cells | Chapter 12 fixes file-based named environments and journal; V1-14 and V1-20 must test actual upstream selectors, process context, rollback and platform host bridges. |
| G04 | Prove request/schema/client conformance | Chapters 9/10/18 and closed schemas fix operation-specific inputs/data, early errors, policy, operation records, limits and server modes; every executable phase from V1-03 tests its declared CLI/MCP pair with a real client, while V1-09 expands full guarded-registry coverage. |
| G05 | Pin and qualify target/dependency manifest | Chapter 14 fixes matrix/evidence method; V1-22 must publish exact tested versions and resolve each v0 support cell. |
| G06 | Qualify sidebar source accuracy | Chapter 11 fixes no-invention/freshness behavior; V1-17 must record per-harness signals, quota limits, performance and layout evidence. |
| G07 | Qualify selected review TUIs | Chapter 11 fixes offline/online and write boundaries; V1-18 must verify LazyGit/gh-dash/web fallback, licenses and accessibility. |
| G08 | Reconcile source/ledger divergences | V1-01 and V1-20 compare deserializers, validators and shipped assets to generated rows, resolving stale docs before parity verification. |
| G10 | Preserve shared Hugo documentation and release lines | V1-02 establishes the shared Hugo brand-theme build; V1-19 must prove both-direction deploy preservation, current-v0 labels, alpha/beta candidate snapshots and no stale v0/reverted-v1 instructions. |
| G11 | Qualify configuration and logging adapters | V1-06: test strict process-settings/env-file parser, CLI/MCP effective-value parity, protected sink authority, rotating-file behavior on Linux/macOS, and fail-closed behavior for required sinks and operation records. |

A capability is complete only when its qualification checks pass on the
claimed targets. Missing evidence remains an explicit unsupported or unverified
cell; schema validity alone does not prove runtime compatibility.

## Specification maintenance checks

Run inventory and contract generation with `--check`, `validate.mjs`,
`uv run --script spec/v1/scripts/validate-contracts.py`, and `git diff --check`.
The inventory extractor is deliberately baseline-specific, not a general Rust
parser. It preserves original declarations alongside 66 atomic command-action,
74 argument and 265 configuration rows. The v0 behavioral source map covers
all 48 existing F rows. The Node validator checks links, trace coverage,
roadmap ordering and selected fixtures. The Python validator uses maintained
Draft 2020-12 validation, including format assertions and cross-schema references;
the Node subset is not sufficient on its own. Changes to the baseline require semantic
review of extraction coverage, not merely new counts.

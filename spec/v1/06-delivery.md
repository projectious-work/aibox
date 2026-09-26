# 6. Standards, delivery and review blockers

## Standards baseline

Company standards baseline: internal commit
`d096a1992ab91dea2265cf4df342e7a97380b0e4`, inspected 2026-09-25. Names are plain
text so they remain visible in renderers that do not resolve relative links.
Standards are governed upstream; this document maps applicability, not a fork.

| Standard | Application / evidence |
|---|---|
| Application profiles | Infrastructure/template distribution plus Go CLI/MCP, UX, docs, schemas and host-gated release; apply behavior-specific obligations. |
| Application configuration | R-CONFIG/R-PRECEDENCE/R-POLICY; ownership, provenance, separate authority and preferences. |
| Application output, logging and evidence | R-RESULT/R-OBSERVABILITY; AC-CLI/AC-DOCTOR/AC-RELEASE; separate results, diagnostics and durable evidence. |
| Compatibility and machine interfaces | R-INTERFACES/R-DEPENDENCIES/R-MIGRATION; versioned schemas, migrations, documented errors and supported windows. |
| Product roadmap and development evidence | `roadmap.yaml` in this directory; stable phase IDs, evidence required before shipped. |
| Spec-driven development cycle | Accept specification, independent plan review, implementation waves, independent complete-baseline conformance review. |
| Security and software supply chain | R-AUTHORITY/R-POLICY/R-REUSE; AC-SEC/AC-TOOLS/AC-RELEASE; threat model and dependency clearance. |
| Software verification and release engineering | Unit/component/blackbox/integration/e2e, negative tests and exact-candidate evidence; published-artifact verification. |
| Human-controlled host-phase execution | Owner-invoked bounded host gates; no host socket/privileged companion shortcut, credential-free validation separate from publication. |
| Host-gated release conformance | AC-TARGET/AC-RELEASE; handoff identity, runtime capability probes, cleanup, repeatability and complete evidence. |
| Git branching and release promotion | Topic branch → `v1.x-dev` → pre-release/release pointers → main under approved standard; no new unique promotion commits. |
| Open-source documentation strategy | Canonical versioned docs/README/roadmap/changelog/releases and local gh-pages deployment; no GitHub Actions. |
| AI-agent accessibility and generative discovery | R-GUIDANCE; public Markdown/discovery resources and offline-after-install read-only stdio server. |

[Canonical standard directory](https://github.com/projectious-work/internal/tree/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards).
The company agent-native interface decision `CuriousSpire` also applies:
interface-neutral core, bounded deterministic operations, CLI/MCP equivalence,
explicit authority and durable evidence. Its normative publication was pending
in [internal PR #12](https://github.com/projectious-work/internal/pull/12)
during boundary review; reconcile the merged source version before acceptance.

## Requirement-to-evidence trace

| Requirement | Planned verification | Phase |
|---|---|---|
| R-PRODUCT | AC-PK, AC-GUIDE, minimal standalone journey | V1-01, V1-03, V1-19 |
| R-PARITY | All F/CFG/ADDON/CMDTYPE/ASSET rows | V1-04, V1-07–V1-12, V1-16, V1-19 |
| R-CORE | AC-CLI, AC-BUILD; architecture review rejects duplicate engines | V1-02 |
| R-BOUNDARY | AC-PK, AC-SEC and product handoff review | V1-01, V1-17 |
| R-AUTHORITY | AC-SEC in both contexts | V1-02, V1-06, V1-17 |
| R-LOCAL | AC-SEC, AC-UX, AC-DOCTOR | V1-07–V1-13 |
| R-REUSE | Dependency clearance plus AC-TOOLS | V1-04, V1-18 |
| R-OBSERVABILITY | AC-DOCTOR, AC-CLI, AC-SIDEBAR | V1-08, V1-13, V1-14 |
| R-INDEPENDENCE | AC-BUILD direct upstream journey | V1-03, V1-05 |
| R-CONFIG | AC-CONFIG all field rows | V1-03, V1-04, V1-16 |
| R-PRECEDENCE | AC-CONFIG, AC-UX conflict fixtures | V1-09, V1-16 |
| R-OWNERSHIP | AC-UX, AC-MIG interrupted refresh | V1-07, V1-09, V1-16 |
| R-MIGRATION | AC-MIG and every old command disposition | V1-16 |
| R-RECOVERY | AC-MIG, AC-LIFE | V1-05, V1-16 |
| R-INTERFACES | AC-CLI adapter equivalence | V1-06 |
| R-DOCTOR | AC-DOCTOR plus no-effect tests | V1-13 |
| R-RESULT | AC-CLI, AC-LIFE machine schema fixtures | V1-02, V1-06 |
| R-EXECUTION | AC-SEC, AC-LIFE interruption/race tests | V1-05, V1-17 |
| R-GUIDANCE | AC-GUIDE offline/version/client tests | V1-06, V1-13 |
| R-POLICY | AC-SEC authorization bypass tests | V1-06, V1-17 |
| R-TARGETS | AC-TARGET matrix gates | V1-18 |
| R-DEPENDENCIES | AC-TOOLS, AC-RELEASE artifact manifests | V1-04, V1-18, V1-19 |

N01 → AC-SIDEBAR in V1-14 and N02 → AC-REVIEW in V1-15. Feature rows refine these
groups and link to source evidence. Each phase requires an independent plan
review before implementation and complete conformance review afterward. This
document is a proposed sequence, not the reviewed implementation plan.

| v0 feature rows | Primary implementation phase | Cross-cutting proof |
|---|---|---|
| F01–F02 | V1-03, V1-12 | V1-16 migration |
| F03–F05 | V1-05 | V1-06 interface, V1-17 security |
| F06–F08 | V1-16 | V1-05 exact runtime identity |
| F09 | V1-19 | V1-06 CLI distribution |
| F10–F13 | V1-03–V1-04 | V1-16 mapping, V1-18 target qualification |
| F14–F17 | V1-07, V1-12 | V1-17 authority |
| F18–F22 | V1-08 | V1-09 theme, V1-13 doctor |
| F23–F25 | V1-08 | V1-14 status integration |
| F26–F29 | V1-09 | V1-16 alias conversion |
| F30–F38 | V1-10 | V1-18 platform coverage |
| F39 | V1-11 | V1-18 audio cells |
| F40–F42 | V1-04 | V1-18 tool qualification |
| F43–F44 | V1-13 | V1-06 machine interface |
| F45–F46 | V1-18 | V1-19 support claim |
| F47 | V1-19 | V1-13 guide content |
| F48 | V1-13 | V1-14 telemetry truthfulness |

## Review blockers / explicit unknowns

| ID | Question to close | Required evidence / responsible role |
|---|---|---|
| G01 | Verify effective configuration mapping | Chapters 8/12 fix authority and converter behavior; V1-03/04/16 must provide exact field destinations, cross-field/default fixtures and native-resolution proof. |
| G02 | Qualify Feature selection and processkit API | Chapter 8 fixes selection criteria; V1-04/12 must fill per-tool manifest, license/provenance and supported processkit API/version evidence. |
| G03 | Prove snapshots, recovery and audio cells | Chapter 12 fixes file-based named environments and journal; V1-11/16 must test actual upstream selectors, process context, rollback and platform host bridges. |
| G04 | Prove request/schema/client conformance | Chapters 9/10 and `operation-request.schema.json` fix common inputs, selectors, exit classes, cancellation and server modes; V1-06 must test generated SDK tool schemas and real MCP clients. |
| G05 | Pin and qualify target/dependency manifest | Chapter 14 fixes matrix/evidence method; V1-18 must publish exact tested versions and resolve each v0 support cell. |
| G06 | Qualify sidebar source accuracy | Chapter 11 fixes no-invention/freshness behavior; V1-14 must record per-harness signals, quota limits, performance and layout evidence. |
| G07 | Qualify selected review TUIs | Chapter 11 fixes offline/online and write boundaries; V1-15 must verify LazyGit/gh-dash/web fallback, licenses and accessibility. |
| G08 | Reconcile source/ledger divergences | V1-01/16 reviewers compare deserializers, validators and shipped assets to generated rows, resolving stale docs before parity sign-off. |
| G09 | Reconcile standards version | Owner accepts a specific merged agent-native standard revision before V1-01 exits; avoid a moving policy reference. |

These are empirical implementation and baseline-acceptance gates, not silent
design discretion. Chapters 7–14 specify behavior and owners; a phase cannot
ship until its listed evidence closes the corresponding gate. A reviewed
specification may still require target qualification, but it must not assert
support without it.

## Validation of this documentation change

Run inventory reproducibility and `validate.mjs`, plus `git diff --check`.
The inventory extractor is deliberately baseline-specific, not a general Rust
parser. It preserves original declarations alongside 66 atomic command-action,
74 argument and 265 configuration rows. The v0 behavioral source map covers
all 48 existing F rows. The validator checks the roadmap, aibox customization
and operation-result schemas with positive/negative fixtures, but it is not a
general JSON Schema implementation. Changes to the baseline require semantic
review of extraction coverage, not merely new counts.

No runtime parity, host gate or dependency qualification was executed for
this draft. The existing Rust repository asks for cargo test/clippy before
commit; both were attempted but `cargo` is absent in this authoring environment.
This limitation must appear in the draft PR. Product code is unchanged.
The owner explicitly approved a documentation-only exception on 2026-09-25
for this draft PR. This is not a waiver for later product implementation,
release validation or host gates.

## Implementation acceptance gate

The owner accepts a revised specification source commit after blocker closure
and independent review. The implementation plan then binds requirements,
work packages, tests, target matrix, risks, reviewers and rollback. Every phase
has a development note, conformance matrix and evidence. Shipped status needs
a release and user documentation; previews with incomplete parity say so.
No external lifecycle/host/publication authority is delegated by this PR.

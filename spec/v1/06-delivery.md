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
| R-PRODUCT | AC-PK, AC-GUIDE, minimal standalone journey | V1-03 |
| R-PARITY | All F/CFG/ADDON/CMDTYPE/ASSET rows | V1-06 |
| R-CORE | AC-CLI, AC-BUILD; architecture review rejects duplicate engines | V1-02 |
| R-BOUNDARY | AC-PK, AC-SEC and product handoff review | V1-01 |
| R-AUTHORITY | AC-SEC in both contexts | V1-02 |
| R-LOCAL | AC-SEC, AC-UX, AC-DOCTOR | V1-03 |
| R-REUSE | Dependency clearance plus AC-TOOLS | V1-01 |
| R-OBSERVABILITY | AC-DOCTOR, AC-CLI, AC-SIDEBAR | V1-03 |
| R-INDEPENDENCE | AC-BUILD direct upstream journey | V1-02 |
| R-CONFIG | AC-CONFIG all field rows | V1-01 |
| R-PRECEDENCE | AC-CONFIG, AC-UX conflict fixtures | V1-03 |
| R-OWNERSHIP | AC-UX, AC-MIG interrupted refresh | V1-03 |
| R-MIGRATION | AC-MIG and every old command disposition | V1-04 |
| R-RECOVERY | AC-MIG, AC-LIFE | V1-04 |
| R-INTERFACES | AC-CLI adapter equivalence | V1-02 |
| R-DOCTOR | AC-DOCTOR plus no-effect tests | V1-02 |
| R-RESULT | AC-CLI, AC-LIFE machine schema fixtures | V1-02 |
| R-EXECUTION | AC-SEC, AC-LIFE interruption/race tests | V1-02 |
| R-GUIDANCE | AC-GUIDE offline/version/client tests | V1-03 |
| R-POLICY | AC-SEC authorization bypass tests | V1-02 |
| R-TARGETS | AC-TARGET matrix gates | V1-06 |
| R-DEPENDENCIES | AC-TOOLS, AC-RELEASE artifact manifests | V1-06 |

N01 → AC-SIDEBAR and N02 → AC-REVIEW, both V1-05. Feature rows refine these
groups and link to source evidence. Each phase requires an independent plan
review before implementation and complete conformance review afterward. This
document is a proposed sequence, not the reviewed implementation plan.

## Review blockers / explicit unknowns

| ID | Question to close | Required evidence / responsible role |
|---|---|---|
| G01 | Exact configuration schema and complete field mapping | Maintainer + owner: schema, validated examples, effective-value/alias/default migration fixtures and private-override design. |
| G02 | Feature selection and processkit ownership/API compatibility | Integration maintainer: per-tool Feature/package mapping, license/provenance clearance, supported processkit interfaces and versions. |
| G03 | Named environment snapshots, backup/recovery and host audio setup without a broad config CLI | Maintainer: native/delegated workflow proof including process context and secret-safe rollback; owner approves any scope change. |
| G04 | Public command/tool/result schemas and long-operation behavior | Interface maintainer: precise schemas, exit/cancellation/retry/locking contracts and actual client compatibility. |
| G05 | Exact supported target/dependency versions | Release maintainer + host owner: v0 evidence audit and explicit v1 qualification matrix, no silent platform retirement. |
| G06 | Sidebar selection and telemetry availability | UX maintainer: per-harness source/freshness/permissions matrix, shared-quota correctness, bounded collector budget and layout prototype. |
| G07 | PR-review dependency and interaction scope | UX maintainer: end-to-end LazyGit/gh-dash/native-web comparison; verify inline review limitations, license and accessibility. |
| G08 | Source/docs disagreements and census coverage | Reviewer: compare hand-written deserializers/validators and actual shipped assets against generated ledger; resolve stale docs such as base-image harness installation wording. |
| G09 | Standards publication/version alignment | Owner: reconcile agent-native standard and obsolete branch/host-companion guidance before implementation plan. |

These are blocking questions for a later accepted specification/implementation
baseline, not blockers to opening this first draft for review. Requirements
remain required while mappings are unresolved. No generic `TBD` marks them done.

## Validation of this documentation change

Run inventory reproducibility and `validate.mjs`, plus `git diff --check`.
The inventory extractor is deliberately baseline-specific, not a general Rust
parser. It preserves original declarations alongside derived rows; changes to
the baseline require review of extraction coverage, not merely new counts.

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

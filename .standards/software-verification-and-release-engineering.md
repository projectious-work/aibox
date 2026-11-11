# Software verification and release-engineering standard

## Purpose and relationship to branching

This standard defines the evidence required to integrate, release, and publish
projectious.work software. The
[Git branching and release standard](git-branching-and-releases.md)
defines branches, merges, and tags; this document defines what must be
proven about the exact commit and artifacts at each gate.

Repositories document an applicable verification profile. A skipped
non-applicable check is explicit; an unexpected skip in a required gate is a
failure.

Repositories requiring native-host, container-engine, signing, or publication
authority also apply the
[human-controlled host-phase execution standard](human-controlled-host-phase-execution.md).

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## Release identity and chain of custody

Every release defines one canonical immutable identity. It binds, as
applicable, the source revision or tree, version, release channel,
configuration, dependency locks, toolchain, declared inputs, and intended
distribution targets. Produced artifacts add their immutable digests to that
identity.

The project identifies the final mutable step and freezes the candidate before
human authorization, signing, attestation, or publication. Any material change
after freezing creates a new candidate identity and invalidates downstream
evidence. A check, approval, signature, or receipt is reusable only when it is
cryptographically or structurally bound to the unchanged subject it assessed.

Release progression is forward-only across the logical responsibilities of
candidate preparation, freezing, production, verification, authorization, and
distribution. Projects may use different state names or combine mechanically
inseparable steps, but they MUST preserve the ordering and invalidation
semantics. Missing, stale, contradictory, ambiguous, or subject-mismatched
evidence prevents advancement.

Producer, verifier or approver, and publisher responsibilities remain
distinguishable in policy and evidence even when one person or automation
performs more than one role. Higher-risk profiles SHOULD require stronger
separation, independent verification, or distinct credentials according to a
documented threat model.

## Test layers

| Layer | Purpose |
|---|---|
| Unit | Prove pure behavior and edge cases quickly. |
| Component | Prove collaboration behind one boundary with controlled dependencies. |
| Black box | Prove public behavior from outside implementation internals. |
| Integration | Prove supported dependency, protocol, platform, or tool contracts. |
| Disposable end to end | Prove a complete real workflow with explicit credentials, cost, isolation, and cleanup. |

Higher layers are smaller than lower layers (*Testing Pyramid*). End-to-end
tests do not replace focused negative, boundary, and failure-path tests.

Unit and component tests follow *TDD, Chicago School*: tests assert state and
observable results, collaborators inside the core are real, and test doubles
(named per *Test Double (Meszaros)*) replace only ports to external systems,
time, and randomness. Tests do not assert internal call sequences.

Default local tests MUST be deterministic, offline, credential-free,
parallel-safe, and isolated from real user state. External integration and
live tests are opt-in and declare prerequisites.

## Verification gates

### Local fast

- formatting and static analysis applicable to the language;
- unit and focused component tests;
- contract/schema and maintained-example validation;
- documentation checks affected by the change.

### Change-risk

Local-fast evidence plus checks selected from changed risk: black-box,
integration, concurrency, race, fuzz smoke, security, migration,
compatibility, performance, and platform tests.

Impact selection retains a mandatory core and starts from a reviewed universe
of conditional surfaces. Evidence records comparison provenance and proves a
complete partition: every surface is either selected by a deterministic
trigger or skipped with a reason. A publisher revalidates that partition.

### Release preparation (`release/*` pull request)

- the complete supported test and compatibility matrix;
- dependency, license, vulnerability, and secret review;
- schemas, examples, generated references, and golden output;
- documentation build, links, migration, changelog, and roadmap consistency
  (`done` items match the release scope);
- clean-room or integration evidence for public extension contracts;
- builds for every supported release target.

### Candidate artifacts

Tests run against the exact candidate artifacts where practical, not only the
source tree. Candidate evidence records toolchain versions, skips, reruns,
platforms, and sanitized external-system evidence.

Readiness evidence MAY replace nested execution when nested execution would
widen the outer trust boundary without proving a release claim. Full nested
lifecycle evidence remains required when relevant inputs changed or the
release claims that boundary.

### Published release

Download artifacts from the publication surface and independently verify
checksums, signatures or attestations, SBOM presence, installation, version,
help, and a minimal workflow. Publication is incomplete until verification
succeeds or the release is explicitly declared failed.

Candidate validation, publication, and post-publication verification are
distinct lifecycle facts. A successful upload, tag, deployment command, or
HTTP response does not by itself prove that the intended release is complete.
Evidence binds the public result back to the exact approved candidate and
checks the product-specific claims made by every publication surface.

Distribution is idempotent and resumable where the destination permits it.
Recovery from an existing tag, release record, partially uploaded artifact,
published package, deployed site, or incomplete verification first reconciles
remote state against the approved release identity. Matching immutable work is
reused; missing steps resume; conflicting state fails closed. Recovery MUST NOT
rebuild, reapprove, or resign unchanged artifacts merely because publication
was interrupted.

## Implementation integrity

- Product code and production systems contain no dummy, test, or placeholder
  data or implementations. Unfinished functionality is either not exposed or
  shows an explicit not-yet-implemented state that its interface contract
  defines, never a fake result. A machine interface with a closed set of
  result or error codes does not expose unfinished operations unless that set
  defines a not-implemented code.
- Values are computed from their live sources, not hardcoded stand-ins.
- Problems are never fixed by hand inside a running container. Change the code,
  image, or configuration, then rebuild the container or volume, so the fix is
  reproducible and reviewed.

Test doubles and fixtures remain legitimate inside test code.

## Regression, isolation, and reliability

- Every confirmed defect adds or identifies permanent regression coverage.
- A user-visible defect has a public-boundary assertion when feasible.
- Security defects add negative refusal or redaction fixtures.
- Tests use temporary roots and controlled environments.
- Tests do not depend on a developer home, locale, credentials, network, or
  mutable global state unless explicitly classified as integration tests.
- Flaky tests are defects. Quarantine requires an owner, issue, scope, and
  expiry; silent retry-until-green is prohibited.
- Coverage is a risk signal, not a universal percentage target.
- Cache reuse does not waive validation: every gate command and runtime probe
  still runs, the cache policy and identity are evidenced, and a force-fresh
  mode exists for cache-sensitive changes.
- Security-relevant operator interfaces receive structural rendering tests in
  addition to semantic event tests. Layout remains unambiguous at narrow and
  wide sizes.

## Release artifacts and provenance

Releases originate from clean worktrees. The annotated release tag is created
on the squash commit, and its tree MUST equal the approved candidate tree.
Repositories publish the artifacts appropriate to their profile plus:

- checksums;
- source/tag reference;
- human-curated release notes;
- installation and verification guidance;
- SBOMs and signatures or attestations when required by the security profile.

Build metadata identifies product version, source commit, toolchain, and
supported contract versions.

## Release record

The release record includes the release-branch commit, squash commit, tag,
approved tree, target branches, intended version,
approvals, checks, tool versions, artifacts, hashes, warnings, and
post-publication verification. Warnings requiring judgment are acknowledged
explicitly; release gates fail closed on errors.

The record preserves the release identity, transition evidence, distribution
destinations and receipts, and enough lineage to reconstruct what consumers
received, from which inputs, under whose authorization. Failed, skipped,
waived, or incomplete gates retain their actual outcomes and rationale rather
than being normalized into success.

Emergency and exceptional paths are explicit, attributable, scoped,
time-bounded, reviewable, and protected by documented compensating controls.
They do not silently weaken the normal release identity or evidence contract.

Repositories select assurance controls according to release risk and document
each applicable control as passed, not applicable with rationale, or covered by
an owned risk-accepted exception. A common vocabulary is encouraged; one
company-wide phase count, CI provider, registry, signing product, or
implementation language is not required.

## Agent briefing

> Verify a release branch before its squash merge and tag the squash commit.
> Select checks by product profile and change risk, verify exact candidate
> artifacts, and independently verify what was published.

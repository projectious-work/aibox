# Human-controlled host-phase execution standard

Version: 1.5

## Purpose and applicability

This standard governs release or verification work that requires authority or
capabilities intentionally unavailable inside the isolated development
environment. Examples include native-host smoke tests, container-engine
validation, signing, publication, and independent verification of prepared
artifacts.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

Host execution is an exception boundary, not a general escape hatch. A project
MUST keep deterministic preparation, compilation, and all checks that can run
safely in the restricted environment on that side of the boundary. Only
irreducibly host-authoritative work crosses it.

Repositories that use this standard declare the **host-gated release overlay**
from [Application profiles](application-profiles.md).

## Normative security boundary

- The development container MUST NOT receive a host shell, container-engine
  socket, credentials, devices, privileged mode, elevation, or broad host
  mounts.
- Candidate-controlled data MUST NOT select commands, executable paths,
  environment assignments, image names, mounts, build arguments, repository
  coordinates, signing identities, or publication destinations.
- Host tools MUST run with an allowlisted environment, controlled working
  directories, bounded arguments, timeouts, captured output, and exact cleanup.
- Inputs MUST be immutable, provenance-bound, fully checksummed, and validated
  before execution and again immediately before use.
- Evidence creation, attempted-run reuse, unexpected state, and incomplete
  cleanup MUST fail closed.
- Broad process termination, container pruning, cache deletion, or filesystem
  cleanup is prohibited.

These properties are normative. Shell, Python, uv, a particular container
engine, UI framework, and directory spelling are implementation choices.

## Producer/verifier protocol

The restricted producer creates an immutable handover for a requested release
version. It contains the already-built target artifacts, exact host inputs,
source and toolchain provenance, checksums, and a versioned manifest. The host
verifier accepts only this handover and performs only declared host-authoritative
checks.

The handover MUST bind:

- normalized release version and intended tag;
- source repository and release-commit tree (the tagged squash commit does not
  exist before the merge);
- producer and toolchain identity;
- every transferred path and digest;
- the applicable validation-manifest schema version; and
- the supported target/platform claims.

Host instructions MUST use a version-bound or relative bundle locator. A
producer-local absolute path or opaque run identifier MAY appear as evidence
metadata but MUST NOT be the operator-facing locator. When producer and
verifier do not share storage, the project MUST define a portable export/import
step and revalidate the imported bundle.

Target artifacts MUST be prepared before handover. A host MUST NOT compile
source opportunistically unless the project's explicitly reviewed profile
identifies host compilation as the property under test.

## Lifecycle and storage classes

Each run implements these logical states:

1. `prepared` — immutable input exists; runtime and evidence do not.
2. `bootstrapped` — only the closed bootstrap runtime layout has been created.
3. `executing` — validation evidence is created exclusively while tool state is
   contained in mutable runtime storage.
4. `complete` — immutable validation evidence is finalized and exact cleanup is
   verified.
5. `failed` — failure evidence is retained and exact cleanup is attempted and
   verified; publication is prohibited.

The implementation MUST distinguish immutable handover input, mutable per-run
runtime/bootstrap state, immutable validation evidence, mutable publication
journals, and optional persistent trusted caches.

The manifest MUST be finalized only after all manifest-covered evidence has
stopped changing. Publication commands and results belong to the publication
journal and MUST NOT mutate the finalized validation manifest.

Human authorization and signing occur only after the final mutable candidate
step. If a covered input or artifact changes afterward, the handover identity,
affected validation, authorization, and signatures are invalid and a new
candidate is required. After the squash merge, the publisher verifies that the
tagged commit's tree equals the handover tree.

## Bootstrap and executable resolution

A minimal reviewed launcher MAY use an owner-installed bootstrap manager to
acquire an exact policy-approved runtime. This mode is conforming only when:

- the manager is resolved from reviewed fixed locations;
- its invocation path and canonical target are validated and evidenced;
- the requested runtime is exact and within a maintained supported-series
  policy;
- project discovery, caller configuration, package-index overrides,
  credentials, and inherited language/package-manager variables cannot affect
  selection;
- acquisition is integrity-verified, logged, and limited to bootstrap; and
- the acquired runtime path, version, build identity, and digest are verified
  before gate work.

Fixed package-manager symlinks MAY be used when their resolved target is a
regular, trusted, non-world-writable executable. Evidence records both the
invocation path and target identity, and preserves the invocation basename for
wrappers or multicall tools. Ambient `PATH` lookup, caller overrides, and
unchecked helper discovery remain prohibited.

Missing prerequisites MAY produce OS-specific, non-executing installation
guidance. The gate MUST NOT install system packages, invoke elevation, or
change system configuration.

## Capability adapters

Projects MUST select host implementations by required capabilities rather than
product names. A reviewed adapter defines probes and invocation profiles for
the capabilities it claims, such as container build, Compose, local-image
inspection, manifest inspection, push, signing, or native execution.

Docker-compatible implementations and Podman MAY satisfy the same profile when
they pass its probes. Product detection alone is insufficient. Trusted helper
or plugin executables MAY be staged into isolated configuration only when their
identity and source are validated and evidenced.

An unpublished-image rehearsal MAY use a bounded local single-manifest alias
and exact generated-reference substitution when both bind to the candidate
digest and source commit. It MUST NOT modify the publication candidate.

## Validation and publication

Validation MUST be credential-free. Publication is a separate, narrowly
authorized operation constrained to fixed repository coordinates and the
artifacts or images listed by the finalized manifest.

Before every publication attempt, the publisher MUST revalidate the complete
manifest and authorization boundary. Publication MUST be retryable and
idempotent where the remote surface permits it; a failed upload MUST NOT force
expensive candidate validation to rerun.

On retry, the publisher reconciles existing remote tags, records, assets, and
verification state against the finalized manifest. It reuses only matching
immutable results, resumes missing operations, and fails closed on conflicts;
it never treats mere existence as proof of identity.

`--dry-run` means full validation and release-ready evidence without commit,
tag, push, signing, or publication. A no-execution preview uses a distinct name
such as `--plan`.

## Operator contract

At start and completion, the gate states whether network acquisition may
occur; whether the invocation only creates evidence or can mutate release
state; requested version and selected portable handover; canonical input,
runtime, evidence, and journal locations; every external command and resolved
executable; failed stage and non-executing remediation; and cleanup outcome.

Long-running gates SHOULD expose semantic task events for start, heartbeat,
output, pass, warning, skip, failure, and cleanup. Presentation consumes events
only after underlying evidence is flushed. A TUI failure MUST NOT rerun
candidate work and a deterministic plain renderer MUST remain available.
Progress uses a stable denominator and explicit pending, running, passed,
warning, skipped, and failed states. Security-relevant layouts are tested for
one visible frame per panel, absent empty containers, correct heading
association, and equivalent narrow and wide semantic structure. Presentation
and clipboard state are not evidence.

## Coverage, vulnerabilities, and caches

Impact selection and readiness evidence follow the
[verification standard](software-verification-and-release-engineering.md#change-risk);
the publisher revalidates the partition.

Vulnerability results distinguish fix-available blockers, no-fix warnings,
scanner failure, and reviewed exceptions. Raw scanner output and normalized
dispositions are retained.

Content-addressed cache reuse is permitted only when cache identity binds the
candidate/source and relevant toolchain inputs, all gate commands and probes
still execute, and evidence records the policy. A force-fresh mode is required
for cache-sensitive changes. Downloaded caches remain outside repository-wide
scanner roots or are excluded through an explicit reviewed rule.

## Conformance

Implementations MUST satisfy the cases in
[Host-gated release conformance](host-gated-release-conformance.md) and validate
their handover with
[`host-gated-release-handover-v1.schema.json`](schemas/host-gated-release-handover-v1.schema.json).

Reference implementations are listed in the conformance document; they are
examples, not mandatory templates.

## Agent briefing

> Prepare immutable candidates in the restricted environment, cross the host
> boundary through a portable versioned handover, validate without credentials,
> and publish only from finalized evidence under the *Principle of Least
> Privilege*.

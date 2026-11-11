# Host-gated release conformance

## Purpose

This document defines reusable acceptance cases for implementations of the
[human-controlled host-phase execution standard](human-controlled-host-phase-execution.md).
Projects MAY add stronger cases but MUST NOT silently omit applicable cases.

Each case records candidate identity, platform, adapter, outcome, evidence
paths, and the reason for every skip. Negative cases pass only when the gate
refuses safely and produces the expected failure evidence.

## Mandatory case groups

### Handover and lifecycle

- Accept a valid version-bound handover from a checkout at a different root.
- Reject producer-local absolute operator locators and identity mismatches.
- Reject missing, additional, mutable, traversing, linked, special, or
  world-writable inputs where prohibited.
- Reject pre-existing runtime/evidence and attempted run reuse.
- Prove lifecycle ownership and manifest finalization ordering.
- Prove publication journals do not mutate validation evidence.
- Retain failure evidence while performing exact cleanup.

### Bootstrap and environment

- Neutralize inherited home, path, language, package-manager, container,
  scanner, credential, and project-configuration variables.
- Verify bootstrap manager, runtime, dependencies, and locks.
- Cover acquisition and already-present exact-runtime paths.
- Cover safe symlinks, unsafe targets, helpers, wrappers, and multicall tools.
- Emit advisory-only remediation without executing it.

### Capability adapters

- Exercise every supported adapter against its declared probes.
- Cover Docker Desktop, OrbStack, Podman, missing Compose/Buildx, and invalid
  helpers where supported.
- Fail on a missing capability rather than a product-name mismatch.
- Test unpublished-image rehearsal without changing the candidate.
- Verify offline/no-pull build policy where required.

### Coverage and runtime authority

- Prove the mandatory core always runs.
- Prove every conditional surface is selected or skipped with a reason.
- Verify changed-path provenance and publisher revalidation.
- Cover readiness without nested authority and require full lifecycle only for
  changed inputs or claimed behavior.

### Vulnerabilities and caches

- Distinguish fix-available blockers, no-fix warnings, scanner failures, and
  reviewed exceptions while retaining raw output.
- Run every gate check with content-addressed reuse enabled.
- Reject incomplete cache identities and exercise force-fresh behavior.
- Prove downloaded caches cannot be mistaken for repository applications.

### Publication and cleanup

- Reject incomplete, changed, unauthorized, or partition-invalid evidence.
- Retry partial publication without rerunning validation.
- Prove idempotent upload/push behavior where supported.
- Simulate failure after temporary resource creation and verify exact cleanup
  without prune, broad deletion, or global process kill.

### Operator presentation

- Verify evidence is flushed before events are rendered.
- Cover heartbeat, warning, skip, failure, cleanup, and plain fallback.
- Prove UI failure does not rerun candidate execution.
- Test literal control-sequence-safe rendering and copy/selection behavior.
- Assert one frame per panel, no empty containers, correct heading adjacency,
  and equivalent narrow/wide semantic hierarchy.
- Prove presentation and clipboard state are excluded from evidence.

### Dry-run

- Run complete validation and create release-ready evidence.
- Prove no commit, tag, push, signing, publication, or equivalent mutation.
- Distinguish `--dry-run` from a no-execution `--plan` operation.

## Reference profiles

- **ainfra:** two-stage bootstrap, cross-compiled artifacts, version-bound
  handover, native checks, and exact image cleanup.
- **aibox:** multi-runtime adapters, impact-selected coverage, vulnerability
  disposition, retryable publication, evidence-first UI, and cache attestation.

Reference fixtures express portable properties. They do not make uv, Python,
Textual, a container product, or repository-specific paths mandatory.

## Agent briefing

> Implementations of host-phase execution pass every applicable mandatory case
> group (handover and lifecycle, bootstrap, adapters, coverage, vulnerabilities
> and caches, publication and cleanup, operator presentation, dry-run). Record
> candidate identity, platform, adapter, outcome, evidence paths, and a reason
> for every skip. Negative cases pass only when the gate refuses safely with
> the expected failure evidence.

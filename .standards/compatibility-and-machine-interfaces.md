# Compatibility and machine-interface standard

## Purpose and scope

This standard governs public machine-readable and automation-facing contracts:
CLI behavior, APIs, protocols, schemas, files, persisted state, events, and
generated output. It applies only to surfaces a project declares public or
persisted.

## Versioning

Product releases use Semantic Versioning 2.0.0 (`vX.Y.Z` tags; the bump follows
the [Git standard](git-branching-and-releases.md#commit-messages)). Contract versions evolve independently when their
compatibility lifecycle differs from the product release.

Every versioned contract defines:

- identifier and version syntax;
- compatibility rules;
- unknown-version and unknown-field behavior;
- supported-version window;
- deprecation and removal policy;
- migration and rollback behavior; and
- positive and negative fixtures.

Persisted or exchanged data MUST NOT be silently reinterpreted under new
semantics.

Release handovers and validation manifests are exchanged machine contracts.
They define schema version, supported-version window, finalization ordering,
unknown-field behavior, immutable and mutable regions, retry/idempotency
semantics, and positive and negative fixtures. A verifier binds them to the
requested release version, tag, source commit, provenance, and artifact
digests. Historical compatibility exceptions are narrow and named; malformed
evidence is never accepted through a general fallback.

## Machine-readable results

Machine output uses a versioned envelope with a stable discriminator, outcome,
and structured payload. It distinguishes success, refusal, cancellation,
timeout, partial failure, and complete failure. Ordering is deterministic where
order has no domain meaning.

Human prose, ANSI formatting, progress output, and diagnostics do not
contaminate a machine-output stream. Schemas validate all maintained examples,
including format constraints, and negative fixtures prove rejection.

## CLI profile

CLI applications additionally define:

- stable exit-code categories;
- stdout and stderr contracts;
- help and version behavior;
- plain, non-color behavior for pipes and unknown terminals;
- signal forwarding, cancellation, and interruption outcomes;
- partial-write and recovery behavior; and
- compatibility of command, flag, environment, and configuration names.

Removal or renaming provides migration guidance and follows the declared
deprecation horizon.

## Change classification

Changes are classified as compatible, conditionally compatible, or breaking.
The classification considers behavior, validation strictness, defaults,
ordering, performance guarantees, security policy, and persisted data—not only
schema shape.

Breaking changes require a major version (or a minor version below 1.0.0),
migration guidance,
updated fixtures, and an explicit rollback or non-rollback statement.

## Contract verification

Projects maintain representative golden fixtures or conformance cases for
public contracts. Changes to them receive a readable product review. Bulk
regeneration without semantic inspection is prohibited.

Portable handover interfaces use relative or version-bound discovery. A
producer-local absolute path or opaque run ID may be metadata but is not an
operator-facing verifier locator.

## Agent briefing

> Treat public and persisted machine interfaces as products. Version their
> meaning, keep human rendering separate, prove rejection behavior, and never
> call a change compatible based only on syntactic schema compatibility.

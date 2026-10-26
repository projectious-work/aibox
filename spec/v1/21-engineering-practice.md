# 21. Engineering practice for the v1 implementation

This chapter adapts the projectious.work company engineering standards
(development cycle, review, architecture and code design, testing,
refactoring, integration, and writing) to the aibox v1 implementation. It
governs how every roadmap phase is built; the other chapters govern what is
built. Where a product chapter is more specific, the product chapter wins.
Agent execution and communication rules are not part of
the product contract; they live in [`AGENTS.md`](../../AGENTS.md).

This chapter applies to every phase from V1-05 onward without being listed in
each phase's `spec` links. Phases started before specification
version 0.9 (V1-01 to V1-04) are not reopened or re-gated; their evidence
stands as recorded.

Terms in *italics* are
[Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/ccb05b6bd41a663eaf2595e707d98f35d743e61f)
as of commit `ccb05b6b` (2026-10-26): that revision's established meaning
applies; deviations are stated inline. Adopting a later revision is a
specification change.

## Development cycle

**R-CYCLE:** Every roadmap phase and material cross-phase change follows one
cycle:

1. **Baseline.** Record the specification commit, phase ID, applicable
   chapters, `R-*` requirements, acceptance groups, schemas, fixtures, and
   governing decisions. Essential ambiguity stops the phase.
2. **Inventory and trace.** Map every applicable normative requirement to
   planned work, expected evidence, and verification method (black-box
   journey, unit or component test, or visual check).
3. **Plan.** Task graph, file ownership, frozen shared interfaces, and the
   requirement IDs each task addresses.
4. **Independent plan review** (R-REVIEW). Each requirement is `covered`,
   `not applicable` with rationale, `blocked`, `requires a specification
   change`, or `deferred` with explicit authority.
5. **Implement in waves** in dependency order; resolve interface drift before
   dependent work continues.
6. **Discoveries.** A specification defect stops affected work and becomes an
   explicit, reviewed specification change. Implementation never silently
   redefines the contract.
7. **Independent verification** (R-REVIEW) of source, tests, schemas,
   fixtures, documentation, and the development logbook. Each requirement is
   `satisfied`, `partially satisfied`, `not satisfied`, `not applicable`,
   `unverifiable`, or `specification inconsistency`; failure paths included.
8. **Close gaps**, then mark the roadmap item `done` per the roadmap's status
   semantics.

For user-facing behavior, acceptance tests are derived from the baseline by an
author who does not read the implementation, and the implementer does not see
those tests until integration. When a test contradicts the baseline, the
baseline wins and the test is corrected with recorded rationale. This
separation applies only to these black-box acceptance tests; implementers
still write their own unit and component tests test-first (R-TESTSTYLE).

A direct, narrowly scoped change request from the owner skips the planning
round: implement it, update the baseline and affected tests in the same change
with a one-line note, and list the defaults chosen. Changes to scope, public
interfaces, security, or persistence follow the full cycle.

## Independent review

**R-REVIEW:** Plan review and implementation verification are performed by a
reviewer other than the author — a different agent, model, or person that does
not receive the author's reasoning. The reviewer works against the baseline,
never against a derived summary. A review gate runs at most 2 rounds; after
that the integration owner decides open points against the baseline and
records each decision in the development logbook. The designated reviewer
(person or agent identity, and for an agent its model and provider) is
recorded in the development logbook at plan review. Automated model routing,
such as processkit `route_task`, may propose a reviewer, but its proposal is
recorded before the review starts. If the designated reviewer is unavailable,
the gate is blocked and reported; it is never skipped. Changing the reviewer,
model, or provider requires the owner's explicit authorization, recorded in
the development logbook.

## Code design

**R-CODE:** Chapter 7's packages realize *Hexagonal Architecture*: ports in the
application core, adapters at the edge, the executable as composition root.
Inside the core, *Clean Architecture* applies: operation rules live in domain
types and use cases; CLI, MCP, Dev Container CLI, filesystem, and container
engine are details behind ports. In addition:

- `gofmt` and `go vet` are authoritative; each phase's verification script
  fails when `gofmt -l` lists a file or `go vet ./...` reports a finding.
  Idioms follow *Effective Go*. Adding a further linter is a specification
  change that names the tool and commits its configuration.
- *SOLID* at type and package level; each function keeps a single level of
  abstraction (*Single Level of Abstraction Principle*).
- Code, operation names, result codes, logs, and documentation use the
  specification's terms (Ubiquitous Language).
- Product code contains no placeholder implementations or dummy data. An
  unfinished operation is neither registered nor exposed through the CLI or
  MCP until it is complete; the error codes in chapter 18 stay closed.
- Problems found in a running Dev Container are fixed in source, artifacts, or
  configuration and verified by a rebuild, never by hand inside the container.
- The maintainer architecture documentation required by chapter 15 follows
  *arc42* (at least context, building blocks, crosscutting concepts,
  decisions, risks) with *C4-Diagrams* for the context and container views.
  It lives under `docs-site/content/docs/architecture/`; V1-05 creates it and
  each later phase updates the sections it affects. Architecture decisions
  use *ADR according to Nygard* structure (context, decision, consequences)
  and are recorded as processkit DecisionRecords (`DEC-*`) in this
  repository; the arc42 decisions section links them.

## Test style

**R-TESTSTYLE:** Tests follow the *Testing Pyramid* and *TDD, Chicago
School*: they assert results and state, use real collaborators inside the
core, and replace only ports — child executables, filesystem roots, clock,
and network — with fakes named per *Test Doubles (Meszaros)*. Tests do not
assert internal call sequences. Every confirmed defect gains a permanent
regression test; a flaky test is a defect; an unexpected skip in a required
suite fails the gate.

## Refactoring

**R-REFACTOR:** A refactoring names the *Code Smells* it removes and the
*Refactoring Catalog according to Fowler* steps it applies, secures behavior
with characterization tests at the public boundary first, keeps tests green at
every commit, and uses the *Mikado Method* for changes that do not complete in
one step. Refactoring commits use the `refactor:` type and change no behavior;
a diff in golden fixtures, schemas, or result envelopes stops the refactoring.

## Integration

**R-INTEGRATION:** Commit messages and squash-merge titles follow
*Conventional Commits*. After each merge into `v1.x-dev`, the author checks the
verification of the merge commit; a red `v1.x-dev` is fixed forward before new
work builds on it. Acceptance tests written ahead of their implementation
(R-CYCLE) stay on the open change branch and merge into `v1.x-dev` together
with the implementation that makes them pass; `v1.x-dev` never carries an
intentionally failing test.

## Writing

**R-WRITING:** Documentation (chapter 15), CLI and MCP messages, release notes,
and commit messages follow *Plain English according to Strunk & White* and the
specification's Ubiquitous Language. Documentation pages keep the *Diátaxis
Framework* types — tutorial, how-to guide, reference, explanation — visibly
distinct.

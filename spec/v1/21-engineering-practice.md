# 21. Engineering practice for the v1 implementation

This chapter adapts the projectious.work company engineering standards
(development cycle, agent execution, architecture and code design, testing,
refactoring, integration, and writing) to the aibox v1 implementation. It
governs how every roadmap phase is built; the other chapters govern what is
built. Where a product chapter is more specific, the product chapter wins.

Terms in *italics* are
[Semantic Anchors](https://llm-coding.github.io/Semantic-Anchors/): their
established meaning applies; deviations are stated inline.

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
baseline wins and the test is corrected with recorded rationale.

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
records each decision in the development logbook. If the designated reviewer
is unavailable, the gate is blocked and reported; it is never skipped, and the
reviewer, model, or provider is never switched without explicit authorization.

## Agent execution and communication

**R-AGENTS:** AI agents implementing v1:

- run at most 3 sub-agents at once and work in waves beyond that; the primary
  agent keeps concept, core implementation, integration, and review of
  delegated results;
- give a red check at most 1 rerun, report an infrastructure-only failure once
  with evidence, push fixes in batches, and never spend turns polling;
- reply answer-first: next action or result first, numbered steps, one
  closing next action, concrete time estimates, no preamble or recap.
  Personal output styles (for example the `i-have-adhd` skill) are opt-in
  through the person's user-level agent instructions.

## Code design

**R-CODE:** Chapter 7's packages realize *Hexagonal Architecture*: ports in the
application core, adapters at the edge, the executable as composition root.
Inside the core, *Clean Architecture* applies: operation rules live in domain
types and use cases; CLI, MCP, Dev Container CLI, filesystem, and container
engine are details behind ports. In addition:

- `gofmt`, `go vet`, and the repository's linter configuration are
  authoritative; idioms follow *Effective Go*.
- *SOLID* at type and package level; each function keeps a single level of
  abstraction (*Single Level of Abstraction Principle*).
- Code, operation names, result codes, logs, and documentation use the
  specification's terms (Ubiquitous Language).
- Product code contains no placeholder implementations or dummy data; an
  unfinished operation returns an explicit, documented unsupported result.
- Problems found in a running Dev Container are fixed in source, artifacts, or
  configuration and verified by a rebuild, never by hand inside the container.
- Architecture documentation follows *arc42* (at least context, building
  blocks, crosscutting concepts, decisions, risks) with *C4-Diagrams* for the
  context and container views. Decisions use *ADR according to Nygard*
  structure (context, decision, consequences) in the project decision store.

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
work builds on it. Intentionally failing suites of open changes never turn
`v1.x-dev` red.

## Writing

**R-WRITING:** Documentation (chapter 15), CLI and MCP messages, release notes,
and commit messages follow *Plain English according to Strunk & White* and the
specification's Ubiquitous Language. Documentation pages keep the *Diátaxis
Framework* types — tutorial, how-to guide, reference, explanation — visibly
distinct.

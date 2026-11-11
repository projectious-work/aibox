# Spec-driven development cycle standard

## Purpose and scope

This standard defines how projectious.work implements a feature, fix, or
roadmap phase against frozen specification versions. It applies when a
repository has normative specifications. "Phase" below means the roadmap item
being implemented, whether a feature or a phase of several features.

Where specifications live, what a version contains, and how versions are
frozen and changed is defined in the
[specification layout and lifecycle standard](specification-layout-and-lifecycle.md).
This standard requires the selected baseline and every change to it to be
explicit.

A repository adopting this standard, or a later revision of it, names the
first phase it applies to. Phases completed before that point are not reopened
or re-gated; their evidence stands as recorded.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## Small-change path

A direct, narrowly scoped change request from the owner (for example "move
this hint into the sidebar, collapsible, also on mobile") does not restart the
cycle. The agent:

1. creates the next specification version of the affected feature (or its
   `v1`, if none exists) on its own `spec/` branch and freezes it; the owner's
   request approves its content, the squash merge follows the
   [Git standard](git-branching-and-releases.md#approval-and-merge-authority),
   and the version must still be test-ready;
2. implements the change against that version, with the affected tests, on a
   `feat/` or `fix/` branch;
3. decides open details with sensible defaults and lists them afterwards.

No new alternatives, review round, or phase is required. During an active
phase, the integration owner hands the new version to both tracks (§6). A
change that alters scope, public interfaces, security, compatibility, or
persistence is not a small change and follows the full cycle.

## Principles

- Establish the contract before distributing implementation work.
- Trace normative requirements, not every explanatory sentence as if it were
  an independent requirement.
- Use parallel agents only where work and ownership are genuinely separable.
- Review the plan and implementation independently against the complete
  applicable specification.
- Never change the specification silently to legitimize implementation drift.
- A phase is complete only when product truth and evidence agree.

## 0. Establish the baseline

Before planning, record:

- frozen specification versions (`<slug>@v<N>`);
- roadmap phase or change scope;
- applicable specification files, normative requirements, schemas, examples,
  and acceptance journeys;
- governing decisions and company standards;
- unresolved ambiguities and blockers; and
- evidence required for completion.

Essential ambiguity is a stop condition. Explanatory prose is reviewed for
contradiction; normative requirements and executable contracts receive strict
traceability.

### User-facing behavior

When the change introduces user-facing behavior, the baseline also defines how
it looks and feels, not only what exists:

- where the feature lives, how users reach it, how they notice pending work,
  and how it moves;
- at least three alternatives for a new mental model, each a distinct stance
  (different place, navigation, or attention guidance), with a mockup, what it
  makes easy, what it costs, when it wins, and a reasoned recommendation;
  extending an existing feature needs one alternative with a mockup;
- the chosen alternative phrased as testable *EARS* requirements ("When
  <trigger>, the product shall <response>").

The owner approves the chosen alternative with the baseline; from then on it is
as normative as the functional requirements.

### Specification contract

A roadmap item names an outcome; it is not a contract. The contract is a
frozen specification version `<slug>@v<N>` per feature, laid out as the
[layout standard](specification-layout-and-lifecycle.md#contents-of-a-version)
defines: *EARS* requirements, an interface catalog for everything the feature
makes observable, acceptance journeys, schemas, fixtures, and examples. A
roadmap phase is covered by the specification versions of its features. Owner
decisions a specification depends on are made here, not during
implementation. While drafting, the agent asks why each requested attribute
matters until the owner's goal is clear (*Laddering (Interview Technique)*),
batching the questions per round, and records that goal in `overview.md`.

A specification version that is not **test-ready**, as the
[layout standard](specification-layout-and-lifecycle.md#lifecycle) defines it,
blocks approval.

## 1. Build requirement, dependency, and risk inventories

Classify applicable requirements by component, public interface, security,
compatibility, persistence, migration, concurrency, failure behavior,
documentation, test level, external dependency, and implementation order.
Quality requirements name their *ISO/IEC 25010* characteristic.

Produce a traceability matrix that maps every normative requirement to planned
work, expected evidence, and verification method (end-to-end journey, unit or
component test, or visual check). End-to-end journeys cover main paths and
security-critical cases; everything else uses lower test layers. Identify cross-cutting invariants before tasks are
distributed.

## 2. Draft the implementation plan

The plan defines:

- a task dependency graph;
- affected and expected files or components;
- interfaces created or changed;
- tests, fixtures, and validation commands;
- documentation, roadmap, and development-note work;
- security, compatibility, migration, and recovery consequences;
- integration checkpoints and rollback strategy;
- owner and reviewer for every task; and
- requirement IDs addressed by every task.

The plan covers specification conformance work as well as feature code.
Tasks follow the *MECE Principle*: together they cover every requirement, and
each requirement and file has one owning task.

## 3. Assign agents and models

Use the parallelism justified by the task graph, ownership boundaries, and
integration cost, within the limits of the
[agent execution budget standard](agent-execution-budget.md); its cap counts
both tracks and all work packages together. The plan
records:

- independent work packages;
- tasks that must remain sequential;
- exclusive or coordinated file ownership;
- frozen shared interfaces;
- expected deliverables and evidence;
- integration owner; and
- model class chosen for cost, complexity, and risk.

Use cost-efficient models for bounded inventory, mechanical edits, fixtures,
and routine checks. Use deeper reasoning for architecture, security,
migrations, ambiguous failures, synthesis, and conformance review. A single
agent is preferred when tasks depend on one evolving shared mental model.

## 4. Review the plan against the complete baseline

Within the review-round limit of the
[agent execution budget](agent-execution-budget.md), a reviewer other than the
primary plan author checks the draft against every
applicable specification file, normative requirement, schema, maintained
example, acceptance journey, decision, and company standard.

The reviewer is a different agent, model, or person than the author and does
not receive the author's reasoning. If the designated reviewer is unavailable,
the gate is blocked and reported. It is never skipped, and the reviewer, model,
or provider is never switched silently; a substitute requires explicit
authorization and is recorded in the cycle record. The designated reviewer
(person or agent identity, and for an agent its model and provider) is
recorded before the review starts; a reviewer proposed by automated model
routing is recorded the same way.

The reviewer also confirms that every specification version in the baseline
is test-ready.

Each normative requirement is classified:

```text
covered
not applicable, with rationale
blocked by an unresolved decision
requires a specification change
deferred with explicit authority
```

Missing coverage, inconsistent sequencing, overlapping ownership, inadequate
tests, and unjustified exclusions amend the plan before implementation.

## 5. Approve and baseline the revised plan

Record the specification baseline, plan revision, ownership, expected
evidence, accepted exclusions, risks, and integration order. The owner approves
the reviewed specification versions, which are then frozen as the
[layout standard](specification-layout-and-lifecycle.md#lifecycle) defines. The
cycle record holds their `<slug>@v<N>` references; both tracks in §6 start
from the commit that carries all their tags. Baselining makes changes
visible; it does not prohibit learning.

## 6. Implement in controlled waves

Implementation proceeds in dependency order:

1. shared contracts and foundations;
2. independent components;
3. planned integration checkpoints;
4. affected verification after each wave; and
5. resolution of interface drift before dependent work continues.

For user-facing behavior, acceptance tests and implementation are separated by
an information barrier that lasts until the integration loop below has passed.
The primary session is the integration owner: it starts both tracks, writes
the acceptance-test track's brief from the contract only, and is the only
participant that sees both. It does not implement; the implementation track
is a separate session or sub-agent.

- **Acceptance-test track.** A fresh agent session, with no conversation
  history, works in its own worktree created from the frozen specification
  versions and commits to a tests-only branch. Its input is the contract, the
  applicable acceptance journeys, and the test harness at the freeze commit (§5).
  It does not read implementation branches and receives no implementation
  reasoning. Its tests exercise only public interfaces and name the requirement
  or acceptance ID each test verifies.
- **Implementation track.** Works from the frozen specification versions and the
  codebase on its own branch, writes its unit and component tests test-first,
  and does not see the acceptance tests, test plan, or test identifiers while
  the barrier lasts.

Both tracks start after the freeze and run in parallel. An implementer who
cannot see the tests optimizes for the contract, not for test structure. The
integration owner then runs the integration loop:

1. Merge both branches onto the phase branch (see the
   [Git branching standard](git-branching-and-releases.md#specifications))
   and run the acceptance tests. Both tracks keep working on their own
   branches; neither reads the phase branch while the barrier lasts.
2. Classify every failure and record the classification in the cycle record:
   - **implementation defect:** returned to the implementation track as expected
     versus observed behavior in contract terms, never as test source;
   - **test defect:** the test contradicts the contract; returned to the
     acceptance-test track with the test ID, the specification reference, and
     the rationale, never with implementation details;
   - **contract defect:** the contract is ambiguous or wrong; affected work
     stops and §7 applies;
   - **infrastructure failure:** handled under the
     [agent execution budget](agent-execution-budget.md) rerun rule.
3. Repeat until every acceptance test passes without skips. A test changes
   only through the test-defect route and is never weakened to make it pass.

A track that disputes a classification gets at most 2 review rounds; then the
integration owner decides against the contract and records the decision. The
acceptance tests and the implementation merge into `main` together, as one
squash merge of the phase branch. The barrier covers black-box acceptance
tests only.

Every work package reports requirements addressed, files changed, tests run,
deviations, remaining risks, documentation impact, and proposed specification
changes. Implementers do not expand scope without revising the plan.

## 7. Control discoveries and baseline changes

When implementation reveals a specification defect or necessary change, stop
affected work and record a proposed baseline change with rationale and impact
on requirements, plans, tests, security, compatibility, documentation, and
roadmap. Continue only after the change is accepted or rejected through the
project's governing process. An accepted change produces a new frozen
specification version; work in flight follows the layout standard's rule for
new versions.

An accepted change updates the baseline and plan. Already completed work is
rechecked against affected requirements. Implementation convenience never
silently redefines the contract.

## 8. Verify implementation independently

After integration, an independent reviewer, under the same reviewer rules as
§4, checks the implementation and evidence against the same complete baseline
used for plan review. Summaries are
inputs, not proof. Inspect source, tests, schemas, examples, generated output,
documentation, roadmap, development notes, and applicable standards.

Each normative requirement is classified:

```text
satisfied
partially satisfied
not satisfied
not applicable
unverifiable
specification inconsistency
```

Verification includes negative and failure behavior, not only successful
acceptance paths.

## 9. Close gaps and repeat verification

Every result other than satisfied or accepted non-applicability causes one of:

- implementation correction;
- added or corrected tests and evidence;
- documentation correction;
- controlled specification change;
- authorized roadmap deferral; or
- blocked phase completion.

Repeat affected conformance checks after correction. A gap report without
closure is not completion.

## 10. Complete phase evidence

A phase is complete only when these items hold (its *Definition of Done*):

- implementation is integrated;
- required verification passes without unexplained skips;
- the conformance matrix is accepted;
- implemented behavior conforms to the frozen specification versions and
  executable contracts;
- user documentation describes observed behavior;
- the development note captures decisions, deviations, and validation;
- roadmap evidence links are updated (the status follows in a separate
  `chore/` pull request, see the roadmap standard);
- changelog and release material are updated where applicable; and
- follow-up work has explicit ownership and location.

## 11. Retrospect and improve

This step closes the cycle's *PDCA* loop (plan §0–5, do §6–7, check §8–10).
For significant phases, compare planned and actual decomposition, agent/model
cost and effectiveness, integration friction, specification gaps, late defect
detection, and validation value. Feed reusable improvements into company
standards, repository agent guidance, project guidance, or model routing.

## Required cycle record

The cycle leaves a reviewable record containing baseline, frozen specification
versions (`<slug>@v<N>`), plan, task graph, ownership, traceability matrix,
designated reviewers, plan-conformance result, integration-loop failure
classifications, implementation evidence, implementation-conformance result,
accepted baseline changes, gap closures, and completion approval.

## Agent briefing

> Baseline frozen, test-ready specification versions for every feature, map
> each normative requirement, plan, and review independently. For user-facing
> behavior, write acceptance tests in a fresh session from the specifications
> only; implement in parallel without seeing them; loop until all pass,
> classifying each failure. Verify against the same baseline. No requirement
> disappears or changes silently; small changes first freeze the next version.

# Agent execution budget standard

## Purpose and scope

This standard limits how AI agents spend parallelism, review rounds,
verification runs, and waiting time. It applies to every projectious.work repository and to every
harness that dispatches sub-agents or drives verification loops.

The
[spec-driven development cycle](spec-driven-development-cycle.md) decides
which work may run in parallel; this standard caps how much runs at once
and how long loops may continue.

## Principles

- Parallelism is a tool, not a goal. Uncapped fan-out exhausts usage limits
  before it produces results.
- The primary agent keeps the core work: concept, core implementation,
  integration, and review of delegated results. Sub-agents lack the full
  conversation context.
- Every loop has an exit: a round limit, a rerun limit, or an escalation.
- Waiting is not work.

## Concurrency

- At most **3 sub-agents** run concurrently, including short tasks.
- Work exceeding that capacity runs in waves: a new sub-agent starts only after
  a running one has finished. Bundle related tasks instead of splitting them.
- Delegate only bounded side tasks with a precise brief: broad search,
  research, inventory, narrowly scoped follow-ups, or a planned parallel work
  package from an approved plan.

## Review and verification loops

- An external or independent review gate runs at most **2 rounds**. After that,
  the primary agent decides open points against the baseline, records the
  decisions with rationale, and does not renegotiate rejected findings.
- Low-severity findings become owned follow-up work, not merge blockers.
- Fixes from one review round are applied and pushed as one batch, never one
  push per finding.
- Get the affected suites green locally before pushing. Heavy suites (full
  browser acceptance, native contract suites, image builds) run once at the
  ready-for-merge gate or on the phase branch — not on every push.

## Failures and waiting

- A failed verification run or required check gets at most **1 rerun**. A
  second failure is a defect to analyse.
- A verification run or required check that fails only because of
  infrastructure or settings (also failing on the base branch or in unrelated
  changes) is reported once to the owner, with evidence, instead of waiting or
  retrying.
- Do not poll. End the turn and resume on a completion or failure signal, or
  after a single long delay of at least 30 minutes.

## Agent briefing

> At most 3 sub-agents at once, in waves; keep core work yourself. Review
> gates get 2 rounds, failed verifications 1 rerun, infrastructure failures one
> evidenced report to the owner. Push batches, not findings. Never poll.

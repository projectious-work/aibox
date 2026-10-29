---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261026_1412-FastForge-v1-phases-start-from-a-frozen
  created: '2026-10-26T14:12:25+00:00'
spec:
  title: v1 phases start from a frozen test-ready contract with independent acceptance-test
    track
  state: accepted
  decision: 'From V1-05, each v1 roadmap phase first produces a test-ready phase contract
    (spec changes, schemas, fixtures, interface catalog) that is independently reviewed,
    owner-approved and frozen at a recorded contract commit. A fresh acceptance-test
    session works only from that commit in its own worktree on a tests-only branch,
    in parallel with the implementation track. The primary session runs an integration
    loop that classifies each failure as an implementation, test, contract or infrastructure
    defect and repeats until all acceptance tests pass. Codified in spec/v1/21-engineering-practice.md
    R-CYCLE (aibox PR #496); generalizable parts proposed company-wide in projectious-work/internal
    PR #17.'
  context: 'V1-03 and V1-04 were implemented from outcome-level chapters; tests were
    written by the implementer alongside the code, and chapter 7 received in-flight
    status notes. The V1-05 roadmap item and tool ledger (94 undecided rows, no Feature
    interface catalog) would force an independent test author to invent interfaces.
    Chapter 21 (PR #495) required an information barrier but not contract freeze,
    parallel tracks, orchestration, or a failure-classification loop.'
  rationale: Independent tests only add assurance if the contract they derive from
    is complete enough that the test author invents nothing; freezing makes the shared
    input explicit and lets both tracks run in parallel without leaking implementation
    into tests.
  alternatives:
  - option: Keep implementer-written tests alongside code (V1-03/V1-04 practice)
    rejected_because: tests mirror the implementation's understanding rather than
      the contract
  - option: Information barrier without contract freeze
    rejected_because: test author must guess interfaces when the contract is outcome-level
  consequences: V1-05 cannot start implementation tracks until owner row decisions,
    the Feature-model spec change and an interface catalog form a reviewed, frozen
    contract. Done phases V1-01 to V1-04 are not reopened. Phase logbooks must record
    contract commit, track branches and classified integration failures.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-26T14:12:25+00:00'
---

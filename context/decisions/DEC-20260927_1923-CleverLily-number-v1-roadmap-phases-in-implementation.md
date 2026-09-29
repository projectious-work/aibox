---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20260927_1923-CleverLily-number-v1-roadmap-phases-in-implementation
  created: '2026-09-27T19:23:00+00:00'
  updated: '2026-09-27T19:26:14+00:00'
spec:
  title: Number v1 roadmap phases in implementation order and track one active phase
  state: superseded
  decision: Renumber the 23 v1 roadmap phases sequentially in their dependency-respecting
    execution order; retain an old-to-new ID crosswalk; allow at most one in-progress
    phase; mark merged baseline work implemented with evidence rather than shipped
    before a release. The documentation foundation becomes V1-02 and Go operation
    core V1-03.
  context: 'After PR #464 and PR #465 were accepted and merged, the canonical v1 roadmap
    displayed V1-01 and V1-20 as simultaneous in-progress phases even though the documentation
    foundation precedes the Go core. The owner requested an unambiguous implementation
    order and questioned parallel work.'
  rationale: The earlier numerical IDs did not match execution order, and two active
    statuses overstated parallel implementation. Sequential IDs and a single active
    phase make the next work explicit while preserving the accepted documentation
    prerequisite and historical references.
  alternatives:
  - option: Keep legacy IDs and rely on array order/dependencies
    rejected_because: That was the source of the owner's confusion and still displays
      V1-20 beside V1-01.
  - option: Mark baseline shipped and activate Go core immediately
    rejected_because: Shipped requires release evidence, and the accepted specification
      makes the shared documentation foundation a prerequisite.
  consequences: All v1 specification phase references and dependency IDs must be updated
    together; old reviews/issues require the crosswalk. New validation rejects nonsequential
    IDs, forward dependencies and simultaneous active phases. V1-01 is implemented,
    V1-02 remains active, and V1-03 is planned.
  decided_at: '2026-09-27T19:23:00+00:00'
  superseded_by: DEC-20260927_1926-MightyHarbor-sequence-v1-phases-and-advance-go
---

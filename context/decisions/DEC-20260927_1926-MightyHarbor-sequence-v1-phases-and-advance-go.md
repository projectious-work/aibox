---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20260927_1926-MightyHarbor-sequence-v1-phases-and-advance-go
  created: '2026-09-27T19:26:02+00:00'
  updated: '2026-09-27T19:26:14+00:00'
spec:
  title: Sequence v1 phases and advance Go core after shared Hugo documentation completion
  state: accepted
  decision: Renumber all 23 v1 roadmap phase IDs in dependency-respecting execution
    order, maintain a historical crosswalk, and allow only one in-progress phase.
    Mark V1-01 baseline and V1-02 shared Hugo documentation implemented with evidence,
    not shipped; set V1-03 Go operation core as the sole active phase.
  context: The owner requested sequential roadmap numbering and one active phase.
    An initial decision incorrectly treated shared Hugo documentation as still active.
    Inspection after the owner's correction confirmed both merged v0/v1 branches pin
    brand-theme-hugo-vanilla v0.3.4, and both local Hugo builds pass; remaining Docsy
    files are archived pre-migration content.
  rationale: The accepted PRs completed the documentation foundation, so leaving it
    in progress would understate progress and continue the apparent parallel-work
    problem. Implemented differs from shipped because no release evidence exists.
    Sequencing Go core next matches the owner's intended order.
  alternatives:
  - option: Leave documentation foundation in progress
    rejected_because: Both documentation lines use the pinned theme and local builds
      pass; this would preserve a false open phase.
  - option: Remove all Docsy references
    rejected_because: Archived pre-migration content and historical baseline evidence
      remain useful, but must be clearly labeled as noncurrent.
  consequences: V1-03 is now the only active phase. Current README and specification
    text must distinguish archived Docsy evidence from the live Hugo build. Sequential
    IDs and backward-only dependencies are validated.
  decided_at: '2026-09-27T19:26:02+00:00'
  supersedes: DEC-20260927_1923-CleverLily-number-v1-roadmap-phases-in-implementation
---

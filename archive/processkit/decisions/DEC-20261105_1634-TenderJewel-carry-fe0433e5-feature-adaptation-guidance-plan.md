---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261105_1634-TenderJewel-carry-fe0433e5-feature-adaptation-guidance-plan
  created: '2026-11-05T16:34:26+00:00'
spec:
  title: Carry fe0433e5 (Feature adaptation guidance plan) with V1-05 and log it
  state: accepted
  decision: Commit fe0433e5, which plans Feature adaptation guidance and templates
    (AC-GUIDE, chapter 10, V1-09 and V1-16 roadmap entries; DEC-20261030_2049-StrongOasis),
    stays on integration/v1-05-tool-catalog and reaches v1.x-dev with V1-05. It is
    recorded in the V1-05 logbook as out-of-phase content rather than landed separately.
  context: Verification round 1 (sage, finding N3) noted the commit is not V1-05 work
    and was not reviewed as part of the phase.
  rationale: The commit changes documentation and roadmap planning only, with no code
    or acceptance effect on V1-05; splitting it out costs a separate PR for no gain
    in safety.
  consequences: V1-09 and V1-16 start from the planning text in fe0433e5 and review
    it when those phases draft their contracts.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-11-05T16:34:26+00:00'
---

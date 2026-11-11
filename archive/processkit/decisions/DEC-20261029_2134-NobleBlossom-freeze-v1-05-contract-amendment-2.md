---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261029_2134-NobleBlossom-freeze-v1-05-contract-amendment-2
  created: '2026-10-29T21:34:11+00:00'
spec:
  title: Freeze V1-05 contract amendment 2 at d3f6e08e
  state: accepted
  decision: The V1-05 contract commit is now d3f6e08eaf232eb6aa7c77cf6c6d2fcea10d1c3c
    on contract/v1-05-tool-catalog. It replaces 59c35acc. Both tracks move onto it.
  context: Amendment 2 (48f54282, review notes in d3f6e08e) removes rich and adds
    preview/enhanced/glow 3.0.0 and preview/enhanced/bat 0.26.1 per DEC-20261029_2107-QuietSeal.
    It also states that an exempt command name does not exempt a command check. Reviewer
    sage (claude-opus-5-5) approved with no blocking findings; its non-blocking notes
    are applied.
  rationale: The amendment fixes a contract defect that the implementation track found
    and a test-track ambiguity, and an independent review approved it.
  consequences: 'The implementation track removes the rich Feature and adds glow and
    bat, then proceeds to wave 5. The test track drops rich and plans glow and bat.
    PR #498 is updated.'
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-29T21:34:11+00:00'
---

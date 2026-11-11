---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261029_1725-AmberBeam-freeze-v1-05-contract-amendment-1
  created: '2026-10-29T17:25:58+00:00'
spec:
  title: Freeze V1-05 contract amendment 1 at 59c35acc
  state: accepted
  decision: The V1-05 contract commit is now 59c35acc1cf230c81658a5763205b1f91a0b1092
    on contract/v1-05-tool-catalog. It replaces 34d06da0. The acceptance-test track
    and the implementation track move onto it.
  context: After the V1-05 contract froze at 34d06da0 (DEC-20261029_1515-ShinySpire),
    both tracks raised contract questions. Under R-CYCLE step 8 they were classified
    as contract defects and fixed in amendment 1 (16a8b73d, 4ec6dc8e, 59c35acc). Reviewer
    sage (claude-opus-5-5) found one blocking defect in round 1, the Dev Container
    CLI lockfile; round 2 approved 4ec6dc8e with no blocking findings, and its cosmetic
    note is applied in 59c35acc.
  rationale: The amendment fixes defects that both tracks found independently, and
    an independent review approved it within the R-REVIEW round limit.
  consequences: 'The test track adjusts its rules for absence, side effects including
    the lockfile exemption, the requires closure, test IDs and stack IDs. The implementation
    track picks up the versionCoupling direction, upstream dependsOn references and
    the extra-tag rule, then continues with wave 3. PR #498 is updated.'
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-29T17:25:58+00:00'
---

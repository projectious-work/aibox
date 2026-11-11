---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261105_1634-ActiveHill-accept-the-v1-05-amendment-4
  created: '2026-11-05T16:34:22+00:00'
spec:
  title: Accept the V1-05 amendment 4 test commit made outside the separate test track
  state: accepted
  decision: Test commit 463acc1c, committed within 30 seconds of contract commit 57b9de31
    and implementation commit 8031b4ea on feat/v1-05-acceptance-speedup by the integration
    owner, is accepted for V1-05 without a fresh acceptance-test session. The exception
    is recorded in the V1-05 logbook.
  context: Chapter 21 requires the acceptance-test track to work from the contract
    commit separately from the implementation track. Verification round 1 (sage, finding
    B3) found that 463acc1c and the later harness fixes 76828a19, dee0b46e, df2cac38
    and 71124b37 were written by the integration owner in the same session as the
    contract and implementation changes.
  rationale: 463acc1c changes the AC-TOOLS harness (scope selection, parallel jobs,
    AC-TOOLS/guards), not the expectations of any acceptance check, and 8031b4ea only
    hardens downloads. The harness changes were reviewed by sage with the amendment,
    the guards test fails when a Feature loses its offline guard test, and the full
    scope still runs the original per-Feature matrix, which passed 876/876 on arm64.
  consequences: No re-authoring of 463acc1c. Later phases keep the separate test track;
    this exception does not set a precedent for changing test expectations.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-11-05T16:34:22+00:00'
---

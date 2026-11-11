---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261031_1750-NeatDaisy-freeze-v1-05-contract-amendment-3
  created: '2026-10-31T17:50:16+00:00'
spec:
  title: Freeze V1-05 contract amendment 3 at d1c7adb4
  state: accepted
  decision: The V1-05 contract commit is now d1c7adb487bd4cfa7449550ffcd3da05153537ed
    on contract/v1-05-tool-catalog. It replaces d3f6e08e. Integration moves onto it.
  context: Amendment 3 (f01b7111, review notes in d1c7adb4) corrects SMOKE-MDBOOK
    and SMOKE-CHROMIUM and adds COMPAT-CHROMIUM-CHANNEL, per DEC-20261031_1743-MightyRaven.
    The owner approved the amendment before review, on the condition of review by
    sage. Reviewer sage (claude-opus-5-5) approved with no blocking findings; its
    non-blocking notes are applied.
  rationale: The amendment fixes two contract defects found by the light pass, and
    an independent review approved it. The corrected commands have run only on debian-13-arm64
    under rootless Podman; the defects are attributed to the tools' own path and executable
    resolution rather than to a Docker comparison run.
  consequences: Integration regenerates the catalog reference when it merges this
    commit, and widens the compatibility notes wording so a note may record behaviour
    unchanged from v0. Both commands remain unverified on debian-13-amd64 and Docker
    until the host runs.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-31T17:50:16+00:00'
---

---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261031_1856-CuriousBeacon-reduce-ac-tools-run-time-through
  created: '2026-10-31T18:56:53+00:00'
spec:
  title: Reduce AC-TOOLS run time through sampled failure builds, parallel builds
    and shared LaTeX work
  state: accepted
  decision: Future AC-TOOLS runs replace most expected-failure container builds with
    a small real-container sample where offline component tests cover the same paths,
    run builds in parallel, and reuse LaTeX work across builds. The AC-TOOLS contract
    change goes through a reviewed amendment. The full run already in progress continues
    unchanged.
  context: The full debian-13-arm64 run plans 466 builds, 322 of them expected failures,
    and the runner builds one at a time. In the light pass, 14 of 126 builds took
    46% of the time. The owner asked for these three optimizations after seeing the
    measurements so far; timing for the expected-failure builds is still being measured.
  rationale: Most expected-failure builds exercise the same guard code, so a sample
    plus offline per-Feature coverage keeps most of the assurance at much lower cost.
  consequences: Default runs give up real-container failure evidence for each individual
    Feature and rely on offline tests for it. The release evidence policy in DEC-20261030_0853-CoolSun
    is not changed by this decision.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-31T18:56:53+00:00'
---

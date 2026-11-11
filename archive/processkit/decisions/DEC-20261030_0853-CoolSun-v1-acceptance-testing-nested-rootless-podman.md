---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261030_0853-CoolSun-v1-acceptance-testing-nested-rootless-podman
  created: '2026-10-30T08:53:50+00:00'
spec:
  title: 'V1 acceptance testing: nested rootless Podman, two-tier evidence, change-scoped
    patch releases'
  state: accepted
  decision: '1. The integration loop runs in the aibox workspace using nested rootless
    Podman (option 1), driven by the agent, on arm64 and on amd64 where OrbStack emulation
    works. 2. Two-tier evidence for major and minor releases: iterate to green in
    the nested Podman setup, then the owner runs one final Docker host run per architecture
    as the recorded acceptance evidence. 3. Patch releases use nested Podman evidence
    only, without the host tier, and test only what changed, not the full catalog.
    4. Test storage is a cache dedicated to tests that persists across runs, kept
    in one place on the host so the owner can delete all of it in one step.'
  context: Container-based acceptance tests, such as V1-05 AC-TOOLS with about 466
    builds per architecture, need a container runtime. The aibox workspace has no
    Docker. Mounting the host Docker socket would violate the two-actor trust boundary
    (chapter 13). The owner showed that on macOS with OrbStack an ordinary container
    (seccomp=unconfined, no --privileged, no host socket) can run rootless Podman
    with nested containers. Chapter 4 already permits nested rootless Podman as a
    workspace tool.
  rationale: Fast iteration without per-round owner involvement, no host socket access
    for the agent, and bounded host disk use. Feature behaviour is decided inside
    the Debian image at build time and barely depends on the runtime, while T01 Docker
    evidence is still recorded for major and minor releases.
  consequences: The workspace definition needs Podman and dedicated test-cache storage.
    The test runner needs a Podman runtime path and change-scoped selection for patch
    releases. Evidence records the runtime used. A feasibility check precedes the
    switch.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-30T08:53:50+00:00'
---

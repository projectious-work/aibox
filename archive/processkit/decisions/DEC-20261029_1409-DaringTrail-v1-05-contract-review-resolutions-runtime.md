---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261029_1409-DaringTrail-v1-05-contract-review-resolutions-runtime
  created: '2026-10-29T14:09:37+00:00'
spec:
  title: 'V1-05 contract review resolutions: runtime deps, Debian snapshot ownership,
    version counts, deviations'
  state: accepted
  decision: '(1) Features that need a multi-version runtime (node/node, go/go, rust/rustc)
    declare installsAfter instead of dependsOn and fail with AIBOX_FEATURE_DEPENDENCY_MISSING
    naming the Feature to add when the runtime is absent; the selected runtime version
    is always respected; bundles still include the runtime at its default. Single-version
    foundations (python/uv, zig/zig, latex/core, browser-testing/playwright) keep
    dependsOn. (2) Debian-package Features have no version option and install the
    candidate version from the base''s configured repositories; the base''s snapshot
    pin owns reproducibility. Acceptance tests use a reference base fixture pinned
    to one snapshot of trixie, trixie-updates and trixie-security where exact package
    versions are checked; the minimal starter gains the updates and security suites.
    (3) Non-runtime tools offer one exact version, with dropped v0 versions recorded
    in a compatibility note; go, node, rust and python offer their last three versions.
    (4) Approved deviations: azure-cli as an aibox uv-tool Feature; resvg 0.49.0 and
    puppeteer 25.11.0 defaults for arm64; rich row as rich-cli 1.8.1; Packer retained
    under BUSL-1.1.'
  context: Independent R-REVIEW of the V1-05 contract (a3850809) found it not test-ready;
    four findings required owner decisions (runtime dependency semantics, Debian pin
    durability, reduced v0 version offers, deviations from DeepOtter/ActiveJay).
  rationale: installsAfter keeps runtime version choice deterministic; Debian live
    mirrors drop superseded versions, so exact Debian pins in Features break on every
    security update; one version per tool keeps the two-architecture matrix tractable.
  consequences: Selecting a node/go/rust dependent alone fails clearly until the runtime
    is selected. Debian-route tool versions depend on the base snapshot. Supersedes
    the azure-cli upstream reference in DEC-20261029_0802-DeepOtter.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-29T14:09:37+00:00'
---

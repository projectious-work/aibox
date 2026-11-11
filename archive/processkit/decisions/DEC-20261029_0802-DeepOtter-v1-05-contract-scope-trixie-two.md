---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261029_0802-DeepOtter-v1-05-contract-scope-trixie-two
  created: '2026-10-29T08:02:06+00:00'
spec:
  title: 'V1-05 contract scope: trixie two-arch matrix, exact pins, direct upstream
    Features, fresh contract branch'
  state: accepted
  decision: 'The V1-05 phase contract declares debian:trixie on linux/amd64 and linux/arm64
    as the only supported Feature base cells; every other distro/release/architecture
    fails clearly at build time. Each tool Feature''s version option accepts only
    an enumerated list of exact, qualified versions, with the v0 effective version
    as default; no ''latest''. Qualified upstream devcontainers Features (hugo, aws-cli,
    azure-cli, github-cli) are referenced directly at an OCI digest with an exact
    version option and must pass the same matrix and absence tests; no aibox wrapper.
    The contract lives on a fresh branch contract/v1-05-tool-catalog from v1.x-dev
    containing spec, ledger, evidence, fixtures and decisions only; PR #493''s code
    is input to the implementation track.'
  context: 'DEC-20261026_1412-FastForge requires a test-ready, frozen V1-05 contract
    before parallel acceptance-test and implementation tracks. The contract needed
    explicit supported targets, version semantics, upstream reuse policy and a contract
    commit that excludes pre-contract implementation (PR #493).'
  rationale: All V1-04 starters use debian:trixie-slim; a 2-cell matrix keeps qualification
    tractable. Exact pins make builds reproducible and every offered version testable.
    Direct upstream reuse avoids owning installers that already meet the contract.
    A contract branch without implementation preserves the acceptance-test information
    barrier.
  alternatives:
  - option: Trixie + Ubuntu 24.04 (+ bookworm)
    rejected_because: doubles or triples cells and installer variants for V1-05
  - option: Exact pins plus 'latest'
    rejected_because: '''latest'' is unreproducible and cannot be qualified'
  - option: Single pinned version, no option
    rejected_because: blocks user version choice until a new Feature release
  - option: Wrap upstream Features in aibox namespace
    rejected_because: owns 4 more installers without contract benefit
  - option: 'Continue contract on PR #493'
    rejected_because: contract commit would expose implementation code to the test
      author
  consequences: 'Ubuntu, bookworm, Arch and Nix are unsupported in V1-05 and need
    a later contract change. v0 ''unversioned/distro-or-latest'' rows must be resolved
    to exact versions available on trixie for both architectures. Upstream aws-cli,
    azure-cli and gh references must select an exact version instead of ''latest''.
    PR #493 is closed or rebased onto the frozen contract later.'
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-29T08:02:06+00:00'
---

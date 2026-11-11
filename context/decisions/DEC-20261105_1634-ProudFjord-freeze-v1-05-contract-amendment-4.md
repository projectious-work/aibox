---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261105_1634-ProudFjord-freeze-v1-05-contract-amendment-4
  created: '2026-11-05T16:34:17+00:00'
spec:
  title: Freeze V1-05 contract amendment 4 at 3d93aa7e
  state: accepted
  decision: The V1-05 contract commit is now 3d93aa7e32f5cd49416da2087348baf4cb8c3d42
    on contract/v1-05-tool-catalog. It replaces d1c7adb4. Amendment 4 (57b9de31, review
    notes in 3d93aa7e) samples the AC-TOOLS expected-failure items (184 instead of
    466 builds per target), adds AC-TOOLS/guards, and keeps --scope full on the owner's
    command. The owner approves it and the integration branch, which already merges
    3d93aa7e (b73bd368), stays on it.
  context: 'Amendment 4 implements DEC-20261031_1856-CuriousBeacon. Reviewer sage
    (claude-opus-5-5, Anthropic) reviewed the amendment together with feat/v1-05-acceptance-speedup
    on 2026-10-31 and required changes: one blocker (worker stores unusable on a second
    run, fixed in dee0b46e) and notes applied in 3d93aa7e. No approval, freeze decision
    or contract-commit record existed until now, although the tests and the arm64
    evidence follow the amendment (verification round 1 finding B2, 2026-11-05).'
  rationale: The review happened and its findings are resolved; the passing arm64
    full (876/876) and sampled (594/594) runs at 71124b37 exercise the amended contract.
    Recording the freeze closes the governance gap without changing any test or evidence.
  consequences: The V1-05 dev note names 3d93aa7e as the contract commit. The amd64
    host run is judged against amendment 4. The departure from CuriousBeacon's LaTeX
    layer reuse is handled separately (finding N2).
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-11-05T16:34:17+00:00'
---

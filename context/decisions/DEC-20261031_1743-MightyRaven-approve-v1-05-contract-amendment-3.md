---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261031_1743-MightyRaven-approve-v1-05-contract-amendment-3
  created: '2026-10-31T17:43:39+00:00'
spec:
  title: 'Approve V1-05 contract amendment 3: correct mdbook and Chromium smoke commands'
  state: accepted
  decision: 'Amend the frozen V1-05 catalog: SMOKE-MDBOOK uses --dest-dir mdbook-book,
    and SMOKE-CHROMIUM adds --channel chromium to playwright screenshot. browser-testing/chromium
    keeps --no-shell, and a compatibility note tells users to select the chromium
    channel for headless runs.'
  context: The nested Podman light pass on debian-13-arm64 found both smoke tests
    failing on the frozen catalog, independently of the engine. mdbook resolves --dest-dir
    from the working directory, so ../mdbook-book points outside the workspace. The
    Feature installs with --no-shell, but playwright screenshot defaults to the headless
    shell it skips. Each test passed with the proposed command in a temporary catalog
    edit.
  rationale: Both are contract defects under R-CYCLE step 8. Keeping --no-shell matches
    v0 and keeps images smaller.
  alternatives:
  - option: Drop --no-shell and install the headless shell too
    rejected_because: Larger images and a departure from v0, though Playwright defaults
      would work unchanged.
  consequences: Default headless Playwright runs need the chromium channel selected,
    and the catalog reference documents this. Independent review by sage still applies
    before the new contract commit is frozen; a blocking finding returns the amendment
    to the owner.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-31T17:43:39+00:00'
---

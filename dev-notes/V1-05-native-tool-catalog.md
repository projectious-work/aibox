# V1-05 native tool catalog

V1-05 is in the phase-contract step of R-CYCLE. The contract is drafted on
`contract/v1-05-tool-catalog` and awaits independent review and owner freeze.
No implementation or acceptance-test track has started from it.

## Baseline

- Specification base: `v1.x-dev` at `9cdcbcbe` (R-CYCLE from PR #496).
- Phase: V1-05, restore the complete optional tool catalog.
- Chapters: 04 reuse targets, 05 acceptance (`AC-TOOLS`), 08 native
  artifacts (Feature composition contract and Tool Feature interface), 19
  configuration realization.
- Requirements: R-REUSE, R-DEPENDENCIES, R-CYCLE, R-REVIEW.
- Census baseline: v0 `9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6`, 98 addon
  tool rows in 41 recipes.
- Decisions:
  - DEC-20261002_0046-SmoothLotus and DEC-20261002_0054-AmberCompass:
    grouped Features, per-cell qualification.
  - DEC-20261026_1412-FastForge: frozen test-ready contract and parallel
    tracks.
  - DEC-20261029_0802-DeepOtter: trixie amd64/arm64 only, exact pins,
    upstream Features used directly, separate contract branch.
  - DEC-20261029_0835-HumbleBridge: row dispositions and removals.
  - DEC-20261029_1034-ActiveJay: additions, uv-managed Python, security
    groups.

## Contract contents

| Artifact | Role |
|---|---|
| [`spec/v1/catalog/v1-05-features.json`](../spec/v1/catalog/v1-05-features.json) | Interface catalog: 39 groups, 119 Features, removals, compatibility notes, error codes, targets, smoke tests |
| [`spec/v1/feature-catalog.schema.json`](../spec/v1/feature-catalog.schema.json) | Catalog schema, checked by `validate.mjs` |
| [Chapter 8, Tool Feature interface](../spec/v1/08-native-artifacts.md#tool-feature-interface) | Identity, shapes, options, dependencies, targets, errors, presence, absence, side effects, publishing for tests, inspection |
| [Chapter 5, AC-TOOLS items](../spec/v1/05-acceptance.md#ac-tools--complete-optional-tool-matrix) | Acceptance item IDs and known-risk stacks |
| [`spec/v1/ledger/v1-05-tool-selection.json`](../spec/v1/ledger/v1-05-tool-selection.json) | Census table: each of the 98 rows retained with its Feature or removed with its decision |
| `spec/v1/fixtures/v1-05-smoke/` | Smoke-test inputs |
| `spec/v1/scripts/validate-v1-05-catalog.py` | Cross-reference checks between census, catalog and table |

Census result: 94 rows retained, 4 removed (mistral, starlight, docusaurus,
mkdocs). Catalog additions without a census row: vibe, typescript (tsc,
biome, oxlint), zig (zig, zls), security/crypto (age, sops, gnupg,
openssh-client, minisign, openssl) and semgrep.

Version and install-route data were researched on 2026-10-29 from producer
release data, registries and the trixie, trixie-updates and trixie-security
package indexes for both architectures. Version-output regular expressions use
a version-token rule rather than exact output formats, so an output layout
change does not become a false contract failure. The
[gap review](V1-05-gap-review.md) is the pre-decision input and is superseded
by the catalog.

## Plan

Waves for the implementation track, in dependency order:

1. Shared conventions: target and option guard, error line, download and
   verification helpers with component tests for the component-verified error
   codes, and `scripts/v1-05-features.sh publish`.
2. Foundation Features: `python/uv`, `node/node`, `go/go`, `rust/rustc`,
   `zig/zig`, `latex/core`, `browser-testing/playwright`.
3. Dependent and independent tool Features, then atomic Features.
4. Bundles and the customized example; inspection parity for
   `devcontainer.features`.
5. User documentation for selection, versions, removal and compatibility
   notes; v1 site build.

The acceptance-test track works on a tests-only branch from the contract
commit and covers the `AC-TOOLS/*` items. The implementation track reuses the
inspection code from PR #493 where it matches the contract.

## Logbook

| Date | Entry |
|---|---|
| 2026-10-29 | Contract drafted on `contract/v1-05-tool-catalog`; independent review pending. |

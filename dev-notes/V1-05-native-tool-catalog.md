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
  - DEC-20261029_1409-DaringTrail: review resolutions for runtime
    dependencies, Debian versions, version counts and approved deviations.

## Contract contents

| Artifact | Role |
|---|---|
| [`spec/v1/catalog/v1-05-features.json`](../spec/v1/catalog/v1-05-features.json) | Interface catalog: 39 groups, 119 Features, reference base, removals, compatibility notes, error codes, targets, smoke tests |
| [`spec/v1/feature-catalog.schema.json`](../spec/v1/feature-catalog.schema.json) | Catalog schema, checked by `validate.mjs` |
| [Chapter 8, Tool Feature interface](../spec/v1/08-native-artifacts.md#tool-feature-interface) | Identity, shapes, options, dependencies, targets, errors, presence, absence, side effects, publishing for tests, inspection |
| [Chapter 5, AC-TOOLS items](../spec/v1/05-acceptance.md#ac-tools--complete-optional-tool-matrix) | Acceptance item IDs and known-risk stacks |
| [Chapter 19](../spec/v1/19-configuration-realization.md) | Grouped Feature realization and owner-approval rule |
| `spec/v1/fixtures/v1-05-base/`, `spec/v1/fixtures/v1-05-targets/` | Reference base (snapshot `20261028T000000Z`) and unsupported-target fixtures |
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

## Inventory and trace

| Requirement | Planned work | Evidence | Verification |
|---|---|---|---|
| R-REUSE | Upstream Features used directly; owned Features reuse producer installers | Catalog `upstream`, `method`, `docsUrl`, `verification` | Specification check; `AC-TOOLS/present` for upstream Features |
| R-DEPENDENCIES | Exact pins, checksums or signatures, licences, Debian versions from the base snapshot | Catalog tools; selection table | Specification check; component tests for download and verification errors; `AC-TOOLS/present`, `AC-TOOLS/version` |
| R-CYCLE | Frozen contract, separate tracks, integration loop | This logbook | Plan review; contract commit record |
| R-REVIEW | Independent plan review and implementation verification | Reviewer record below | Review reports |
| AC-TOOLS | Every `AC-TOOLS/*` item in chapter 5 | Acceptance test results on both targets | Black-box acceptance tests |
| Chapter 8 census obligations 1–5 | Census traceability, compatibility notes, isolation and side-effect checks | Catalog and selection table | `AC-TOOLS/trace`, `AC-TOOLS/isolated`, `AC-TOOLS/side-effects` |

## Plan

Implementation track waves, in dependency order. Paths are owned by the
implementation track unless noted.

| Wave | Work | Files | Requirements |
|---|---|---|---|
| 1 | Shared install conventions: target and option guards, error line, download and verification helpers with component tests for the component-verified codes; publish script | `features/src/**/install.sh` helpers, `features/test/`, `scripts/v1-05-features.sh` | R-DEPENDENCIES, AC-TOOLS/unsupported-*, AC-TOOLS/requires |
| 2 | Foundation Features: `python/uv`, `python/python`, `node/node`, `go/go`, `rust/rustc`, `zig/zig`, `latex/core`, `browser-testing/playwright` | `features/src/<id>/` | AC-TOOLS/present, version |
| 3 | Remaining tool Features, then atomic Features | `features/src/<id>/` | AC-TOOLS/present, version, isolated, depends, side-effects |
| 4 | Bundles; starter APT suites; inspection parity for `devcontainer.features` | `features/src/*/all/`, `templates/src/`, `spec/v1/examples/` rendering, `internal/inspection/` | AC-TOOLS/bundle, conflict, inspect |
| 5 | User documentation for selection, versions, removals and compatibility notes; v1 site build | `docs-site/content/docs/tool-catalog.md` | Roadmap docs obligation |

The acceptance-test track works on `tests/v1-05-tool-catalog` from the
contract commit, owns `tests/acceptance/v1-05/`, and covers every
`AC-TOOLS/*` item on both targets. The implementation track works on
`feat/v1-05-tool-features` and may reuse PR #493's inspection code where it
matches the contract.

## Review

Designated reviewer: TeamMember `sage` (agent), model `claude-opus-5-5`,
provider Anthropic. `route_task` proposed no TeamMember for this review.

| Round | Commit | Result |
|---|---|---|
| 1 | `a3850809` | Not test-ready: 11 blocking findings. Owner decisions in DEC-20261029_1409-DaringTrail; remaining findings fixed as contract defects. |

## Logbook

| Date | Entry |
|---|---|
| 2026-10-29 | Contract drafted on `contract/v1-05-tool-catalog` (`a3850809`). |
| 2026-10-29 | Review round 1 findings resolved; round 2 pending. |

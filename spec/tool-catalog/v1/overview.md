# Optional tool catalog: overview

This file is not normative. The contract is `requirements.md`,
`interfaces.json`, `journeys.md`, the catalog, the schema and the fixtures;
`requirements.md` states their precedence.

## Goal

Restore every optional tool of v0 as native Dev Container Features, so that
a project selects, versions and disables tools through standard
`devcontainer.json` `features`, without an aibox-specific package manager.
Roadmap item V1-05.

## Origin of this version

This is a thin specification (owner decision on issue #504, 2026-11-15). It
moves the V1-05 contract, written before the
[specification layout standard](../../../.standards/specification-layout-and-lifecycle.md)
was adopted, into `spec/tool-catalog/v1/`. Its requirements and journeys
restate, as provenance only, these sources:

- chapter 8 of the v1 product description:
  ["Feature composition contract"](../../../docs/product/v1/08-native-artifacts.md#feature-composition-contract),
  ["Tool Feature interface"](../../../docs/product/v1/08-native-artifacts.md#tool-feature-interface) and
  ["Census traceability"](../../../docs/product/v1/08-native-artifacts.md#census-traceability);
- the [`AC-TOOLS` group](../../../docs/product/v1/05-acceptance.md#ac-tools--complete-optional-tool-matrix)
  of chapter 5;
- [chapter 19, "Tool installation realization"](../../../docs/product/v1/19-configuration-realization.md#tool-installation-realization).

The catalog, its schema and the fixtures are byte-identical to the files that
V1-05 implemented and accepted; only their directory changed (path map:
[`docs/engineering/spec-layout-migration.md`](../../../docs/engineering/spec-layout-migration.md)).
The generated Feature sources under `features/src/` name the new catalog path
in their header comment, so the published Feature tarballs differ in those
bytes from the ones accepted at `651506fa`; behaviour does not change. The owner accepted this
byte change without a rerun on 2026-11-15.

Contract history: first contract `34d06da0` (DEC-20261029_1515-ShinySpire),
amendments 1 `59c35acc` (AmberBeam), 2 `d3f6e08e` (NobleBlossom), 3
`d1c7adb4` (MightyRaven, NeatDaisy) and 4 `3d93aa7e` (ActiveHill,
ProudFjord). The development note
[`dev-notes/V1-05-native-tool-catalog.md`](../../../dev-notes/V1-05-native-tool-catalog.md)
records each amendment, review and owner decision.

## Scope

- 39 groups and 120 Features (aibox-owned and upstream), 5 removed census
  rows with owner decisions, 36 compatibility notes, 10 smoke tests.
- Targets `debian-13-amd64` and `debian-13-arm64` only
  (DEC-20261029_0802-DeepOtter); Debian 12 and Ubuntu 24.04 are unsupported
  fixtures.
- Out of scope: production publication and the public Feature index, the base
  toolchain (`native-starters@v1`), and client editor compatibility on a base.

## Outside this directory

- The census selection table `ledger/v0-parity/v1-05-tool-selection.json`,
  the v0 ledger and `ledger/v0-parity/base-build.json` are generated
  inventory; TOOL-002 to TOOL-005 and TOOL-009 refer to them.
- Validators: `scripts/spec/validate.mjs` and
  `scripts/spec/validate-v1-05-catalog.py`. Container harness:
  `tests/acceptance/v1-05/`. Component tests: `features/test/`.

## Evidence

V1-05 is `done` at implementation commit `651506fa`
(`evidence/v1/V1-05-fbd178d6/`):

- `debian-13-arm64`: sampled and full scope at `71124b37`, rootless Podman
  5.4.2 nested in the development container; `AC-TOOLS/guards` at
  `fbd178d6`.
- `debian-13-amd64`: sampled scope at `651506fa`, Docker in an OrbStack
  machine under x86-64 translation on Apple silicon, completed over three
  resumed runs. Full scope never ran on amd64.
- The owner accepted both as host evidence on 2026-11-15.
- V1-05 verification finding N2 (the departure from LaTeX layer reuse of
  DEC-20261031_1856-CuriousBeacon) is closed: the owner acknowledged it on
  2026-11-15. N6 to N8 and N10 stay deferred, N7 before amd64 is advertised.
- Owner decisions of 2026-11-15 for this version: requirements that no
  accepted check exercises stay in the contract, marked "not separately
  verified"; the comment-only Feature header change is accepted.

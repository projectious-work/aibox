# Optional tool catalog: overview

This file is not normative. The contract is `requirements.md`,
`interfaces.json`, `journeys.md`, the catalog, the schema and the fixtures.

## Goal

Restore every optional tool of v0 as native Dev Container Features, so that
a project selects, versions and disables tools through standard
`devcontainer.json` `features`, without an aibox-specific package manager.
Roadmap item V1-05.

## Origin of this version

This is a thin specification (owner decision on issue #504, 2026-11-15). It
moves the V1-05 contract, which was written before the
[specification layout standard](../../../.standards/specification-layout-and-lifecycle.md)
was adopted, into `spec/tool-catalog/v1/` without changing its meaning:

- The requirements restate chapter 8 of the v1 product description, sections
  ["Feature composition contract"](../../../docs/product/v1/08-native-artifacts.md#feature-composition-contract),
  ["Tool Feature interface"](../../../docs/product/v1/08-native-artifacts.md#tool-feature-interface) and
  ["Census traceability"](../../../docs/product/v1/08-native-artifacts.md#census-traceability).
  The journeys restate the
  [`AC-TOOLS` group](../../../docs/product/v1/05-acceptance.md#ac-tools--complete-optional-tool-matrix)
  of chapter 5. Where wording differs, the chapter text at contract commit
  `3d93aa7e32f5cd49416da2087348baf4cb8c3d42` is the meaning this version keeps.
- The catalog, its schema and the fixtures are byte-identical to the files
  that V1-05 implemented and accepted; only their directory changed
  (path map: [`docs/engineering/spec-layout-migration.md`](../../../docs/engineering/spec-layout-migration.md)).
- Contract history: first contract `34d06da0`, amendments 1 (`59c35acc`),
  2 (`d3f6e08e`), 3 (`d1c7adb4`) and 4 (`3d93aa7e`, frozen by
  DEC-20261105_1634-ProudFjord). The development note
  [`dev-notes/V1-05-native-tool-catalog.md`](../../../dev-notes/V1-05-native-tool-catalog.md)
  records each amendment, review and owner decision.

## Scope

- 39 groups and 120 Features (aibox-owned and upstream), 5 removed census
  rows with owner decisions, 36 compatibility notes, 10 smoke tests.
- Targets `debian-13-amd64` and `debian-13-arm64` only
  (DEC-20261029_0802-DeepOtter). Debian 12 and Ubuntu 24.04 are unsupported
  fixtures.
- Out of scope: production publication and the public Feature index; the base
  toolchain (verified by the V1-04 starter checks); client editor
  compatibility on a base.

## Outside this directory

- The census selection table `ledger/v0-parity/v1-05-tool-selection.json`
  and the v0 ledger it is checked against are generated inventory, not part
  of the contract; TOOL-002 and TOOL-003 refer to them.
- `scripts/spec/validate-v1-05-catalog.py` and `scripts/spec/validate.mjs`
  validate the catalog; the AC-TOOLS harness is in `tests/acceptance/v1-05/`.

## Evidence

V1-05 is `done` at implementation commit `651506fa`: sampled and full runs on
`debian-13-arm64` and a sampled run on `debian-13-amd64`, recorded in
`evidence/v1/V1-05-fbd178d6/` and the development note.

# Native workspace starters: overview

This file is not normative. The contract is `requirements.md`,
`interfaces.json` and `journeys.md`;
`requirements.md` states their precedence.

## Goal

Give a user a directly runnable native Dev Container workspace, minimal or
curated, with a persistent non-root home, that works with the upstream Dev
Container CLI alone and leaves every file to the project afterwards. Roadmap
item V1-04.

## Origin of this version

This is a thin specification (owner decision on issue #504, 2026-11-15). V1-04
was implemented before the
[specification layout standard](../../../.standards/specification-layout-and-lifecycle.md)
was adopted and before phases had a frozen test-ready contract. Its
requirements restate, as provenance only:

- chapter 8 of the v1 product description,
  ["Starter distribution and upstream interoperability"](../../../docs/product/v1/08-native-artifacts.md#starter-distribution-and-upstream-interoperability)
  and ["Native configuration example"](../../../docs/product/v1/08-native-artifacts.md#native-configuration-example);
- chapter 12,
  ["Named environments and persistence"](../../../docs/product/v1/12-migration-and-state.md#named-environments-and-persistence);
- [`templates/README.md`](../../../templates/README.md) for NS-003 (bash and
  git), NS-004 (`~/.tmux.conf` after `/etc/tmux.conf`) and NS-005;
- the checks of `scripts/verify-v1-04.sh`, `scripts/spec/validate-examples.mjs`,
  `scripts/spec/validate-contracts.py` and `scripts/verify-v1-04-host.py`.

The closed schema of `customizations.aibox` is not part of this version. By
owner decision of 2026-11-15 it is its own specification,
`workspace-customization@v1`, because later UX phases change its keys; this
version cites it (NS-041, NS-042). The
maintained examples stay in `examples/` because later phases change them;
the Templates stay in `templates/src/` as product source
([path map](../../../docs/engineering/spec-layout-migration.md)).

## Choices recorded here

- The examples use one aligned user (`containerUser` equals `remoteUser`).
  That is a choice of these examples: chapter 8 allows the two users to
  differ, with the container user's home and permissions reviewed
  separately.
- Supported targets (`interfaces.json` `targets`): any engine and host the
  pinned Dev Container CLI 0.89.0 supports for these images; verified only on
  macOS arm64 with Docker. Target qualification is V1-22.

## Boundaries

- Optional tools are native Features; their contract is `tool-catalog@v1`.
  The `customized` example selects the Feature recorded in that
  specification's `fixtures/v1-05-native-features.json` (NS-044).
- No published Template reference exists; OCI Template publication is a
  separate release task.
- The UX keys of `customizations.aibox` are declared by
  `workspace-customization@v1`; their behaviour belongs to later phases (V1-11, V1-12) and, for
  `workspace.theme`, `.mode` and `.layout`, to `settings-and-logging@v1`.

## Evidence

V1-04 is `done` at implementation commit `85cc0c9a`: the offline gate and two
macOS arm64 Docker host runs, recorded in `evidence/v1/V1-04-85cc0c9a/` and
[`dev-notes/V1-04-native-workspace-starters.md`](../../../dev-notes/V1-04-native-workspace-starters.md).
`journeys.md` states what that evidence covers.

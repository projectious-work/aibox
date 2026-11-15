# Workspace customization schema: overview

This file is not normative. The contract is `requirements.md`,
`interfaces.json`, `journeys.md` and `schemas/customization.schema.json`;
`requirements.md` states their precedence.

## Goal

Keep one closed, versioned vocabulary for the aibox-specific preferences a
project may record in `devcontainer.json` under `customizations.aibox`, so
that every phase that interprets a key reads the same names, types and
allowed values, and a native-only project can ignore the namespace.

## Origin of this version

The schema was generated in V1-01 from the v0 baseline inventory and lived at
`spec/v1/customization.schema.json`. On issue #504 the owner decided on
2026-11-15 to give it its own specification instead of placing it in
`native-starters@v1`, because every later UX phase (V1-10 to V1-18) changes
its keys and would otherwise force a new version of the starters
specification. `native-starters@v1` cites this version.

The schema is byte-identical to the former file except for its `$id`, which
names the new path. Requirements restate, as provenance only:

- chapter 2, ["One entry point, native ownership"](../../../docs/product/v1/02-configuration.md#one-entry-point-native-ownership)
  and ["aibox extension schema"](../../../docs/product/v1/02-configuration.md#aibox-extension-schema);
- chapter 5, `AC-CONFIG`, for the closed-schema checks;
- chapter 8, the paragraph on `customizations.aibox` in
  ["Starter distribution and upstream interoperability"](../../../docs/product/v1/08-native-artifacts.md#starter-distribution-and-upstream-interoperability);
- chapter 19, ["Placement and defaults"](../../../docs/product/v1/19-configuration-realization.md#placement-and-defaults).

## Roadmap

`SPEC.toml` lists V1-01 (generated the schema), V1-04 (the starters and
examples that show it) and the planned phases that interpret its sections:
V1-10 (`harnesses`), V1-11 and V1-12 (`workspace` tmux, status, theme and
prompt keys), V1-13 (`latex`), V1-16 (`diagnostics`), V1-17
(`workspace.sidebar`) and V1-18 (`workspace.review`). A phase that changes a
key creates `workspace-customization@v2`. `settings-and-logging@v1` reads
`workspace.theme`, `.mode` and `.layout` as its project layer.

## Evidence

The schema's checks are part of the offline gates `scripts/verify-v1-04.sh`
and `scripts/verify-v1-05.sh`, which passed for V1-04 at `85cc0c9a` and for
V1-05 at `651506fa`.

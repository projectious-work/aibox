---
apiVersion: processkit.projectious.work/v1
kind: Migration
metadata:
  id: MIG-20261024_1722-RuntimeSync-aibox-runtime
  created: 2026-10-24T17:22:22Z
spec:
  source: aibox-runtime-home
  source_url: "aibox://runtime-home"
  from_version: 0.28.3
  to_version: 0.35.2
  state: pending
  generated_by: aibox apply
  generated_at: 2026-10-24T17:22:22Z
  summary: 0 changed upstream, 0 conflicts, 7 new, 0 removed (2 groups affected)
  affected_groups:
    - runtime-misc
    - runtime-yazi
  affected_files:
    - { path: .codex/themes/aibox.tmTheme, classification: new-upstream }
    - { path: .config/bat/themes/aibox.tmTheme, classification: new-upstream }
    - { path: .config/opencode/themes/aibox.json, classification: new-upstream }
    - { path: .config/yazi/plugins/preview-options.yazi/main.lua, classification: new-upstream }
    - { path: .local/bin/aibox-agent-signal, classification: new-upstream }
    - { path: .local/bin/aibox-codex-notify, classification: new-upstream }
    - { path: .local/bin/aibox-size-tree, classification: new-upstream }
---

# Migration MIG-20261024_1722-RuntimeSync-aibox-runtime

Managed `.aibox-home/` runtime changes from `0.28.3` to `0.35.2`.

0 changed upstream, 0 conflicts, 7 new, 0 removed (2 groups affected)

## Counts

- unchanged: 44
- changed-locally-only: 0
- changed-upstream-only: 0
- conflict: 0
- new-upstream: 7
- removed-upstream: 0

- removed-upstream-stale: 0

## Per-intermediate review

Every released version between `from_version` and `to_version` whose              template snapshot is present on disk is listed below, with the number              of files that changed between consecutive snapshots. Useful for              catching scaffolding changes that were introduced and later reverted              across the span of a multi-version upgrade.

- `0.28.3` → `0.28.4`: 0 file(s) changed of 44
- `0.28.4` → `0.28.5`: 43 file(s) changed of 44
- `0.28.5` → `0.28.6`: 0 file(s) changed of 1
- `0.28.6` → `0.28.7`: 0 file(s) changed of 1
- `0.28.7` → `0.28.8`: 0 file(s) changed of 1
- `0.28.8` → `0.28.9`: 0 file(s) changed of 1
- `0.28.9` → `0.28.11`: 0 file(s) changed of 1
- `0.28.11` → `0.28.12`: 0 file(s) changed of 1
- `0.28.12` → `0.28.13`: 0 file(s) changed of 1
- `0.28.13` → `0.28.14`: 0 file(s) changed of 1
- `0.28.14` → `0.28.15`: 0 file(s) changed of 1
- `0.28.15` → `0.28.16`: 0 file(s) changed of 1
- `0.28.16` → `0.28.17`: 0 file(s) changed of 1
- `0.28.17` → `0.28.18`: 0 file(s) changed of 1
- `0.28.18` → `0.28.19`: 0 file(s) changed of 1
- `0.28.19` → `0.29.0`: 0 file(s) changed of 1
- `0.29.0` → `0.30.0`: 0 file(s) changed of 1
- `0.30.0` → `0.30.1`: 0 file(s) changed of 1
- `0.30.1` → `0.31.0`: 0 file(s) changed of 1
- `0.31.0` → `0.31.1`: 0 file(s) changed of 1
- `0.31.1` → `0.31.2`: 0 file(s) changed of 1
- `0.31.2` → `0.31.3`: 0 file(s) changed of 1
- `0.31.3` → `0.31.4`: 0 file(s) changed of 1
- `0.31.4` → `0.32.5`: 0 file(s) changed of 1
- `0.32.5` → `0.32.6`: 0 file(s) changed of 1
- `0.32.6` → `0.33.1`: 0 file(s) changed of 1
- `0.33.1` → `0.33.2`: 0 file(s) changed of 1
- `0.33.2` → `0.34.0`: 1 file(s) changed of 2
- `0.34.0` → `0.34.1`: 0 file(s) changed of 2
- `0.34.1` → `0.34.2`: 1 file(s) changed of 2
- `0.34.2` → `0.34.4`: 1 file(s) changed of 2
- `0.34.4` → `0.34.7`: 1 file(s) changed of 2
- `0.34.7` → `0.35.0`: 1 file(s) changed of 2
- `0.35.0` → `0.35.2`: 0 file(s) changed of 2

## Changes by group

### runtime-misc

**new-upstream**

- `.aibox-home/.config/bat/themes/aibox.tmTheme` -> `.aibox-home/.config/bat/themes/aibox.tmTheme`
- `.aibox-home/.config/opencode/themes/aibox.json` -> `.aibox-home/.config/opencode/themes/aibox.json`
- `.aibox-home/.local/bin/aibox-size-tree` -> `.aibox-home/.local/bin/aibox-size-tree`
- `.aibox-home/.local/bin/aibox-agent-signal` -> `.aibox-home/.local/bin/aibox-agent-signal`
- `.aibox-home/.local/bin/aibox-codex-notify` -> `.aibox-home/.local/bin/aibox-codex-notify`
- `.aibox-home/.codex/themes/aibox.tmTheme` -> `.aibox-home/.codex/themes/aibox.tmTheme`

### runtime-yazi

**new-upstream**

- `.aibox-home/.config/yazi/plugins/preview-options.yazi/main.lua` -> `.aibox-home/.config/yazi/plugins/preview-options.yazi/main.lua`


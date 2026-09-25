# 2. Configuration, persistence and migration

## One entry point, native ownership

**R-CONFIG:** `devcontainer.json` is the project entry point. Use standard
fields, Feature options and referenced native Dockerfile/Compose/tool files.
Use `customizations.aibox` only for aibox-specific preferences that need an
interpreter. No successor `aibox.toml`, no `set` CLI and no generic config-write
MCP tool. Agents use their existing file editors after reading guidance.

The all-v0-settings mapping is a working assumption, not verified upstream
coverage. The [field ledger](ledger/configuration.md) names a disposition for
every reachable field family, including local config and compatibility paths.
Map each accepted old setting to an effective new outcome in migration tests.

| Old family | Proposed authoritative destination | Important distinction |
|---|---|---|
| Project name/image/base/user | `name`, `image` or `build`, `containerUser`, `remoteUser`, UID mapping; native Compose identity where needed | Display name cannot authorize resource deletion. |
| Environment and extra volumes | `containerEnv` or Compose environment/env_file; `mounts` or Compose volumes | Secrets stay in private native inputs; no copied credential literals. |
| Lifecycle/keepalive | Standard lifecycle fields and a container-local keepalive helper | `initializeCommand` executes with host authority; never inferred from local permission. |
| Paths and custom Dockerfile/Compose | Native paths and ordered native overrides | No generator overwrites user-owned definitions. |
| Addons/tool enabled/version | Standard Feature references/options; split capabilities when upstream Feature is too coarse | Disabling means absence where v0 promises it; no mere skipping of install over a base that already contains it. |
| Harness install/enable/order | Feature install/version options; namespaced launch order; native harness config | Host-only integration such as Cursor must remain possible without installing a CLI inside. |
| Theme, prompt, layout, tmux | Namespaced defaults plus native tool files/local overrides | Preserve all fields/choices; explicit managed output paths. |
| Audio | Feature plus environment/mount intent | Host service setup is a separate authorized task. |
| LaTeX | Feature, namespaced document definitions, native Compose read-only preview service | Build/watch inside; serving sidecar never compiles. |
| Processkit/skills/agents/MCP | Optional product integration and native harness configuration | Preserve source forks/pins and custom entries; processkit owns its process configuration. |
| Execution/MCP permissions/seccomp acknowledgment | Requested native policy plus independently held operator policy | Repository-controlled “allow” never grants new host authority. |
| State preservation/purge | Safe retention default, explicit recovery action | A checked-in purge flag is not consent to delete credentials. |
| Diagnostics/thresholds | Namespaced preferences/read-only doctor | In-container checks cannot acquire host access. |

### Proposed aibox extension shape

`customizations.aibox` has `schemaVersion`, `workspace`, `harnesses`, `latex`
and `diagnostics` sections. The [closed proposed schema](customization.schema.json)
retains the 115 v0 workspace UX leaf paths under `workspace` with their existing
names and primitive types. This is a narrow preservation of aibox-owned UX
intent, **not** a transplant of old image, install, security or lifecycle
settings. `workspace` also contains the separately identified new sidebar and
review selections. Named v0 enum values are closed by the schema; defaults,
cross-field rules and semantic constraints remain validated by the
aibox-specific checker. Schema acceptance alone does not prove effective-value
equivalence.
Unknown keys and null are errors, absence selects the documented default, and
explicit false/empty values remain distinct. Upstream Feature options remain
the sole install/version authority. G01 now concerns tested effective mapping,
not the existence of a proposed namespace shape.
`schemaVersion` is a decimal major version; this draft accepts only `"1"`.
Unknown versions fail closed with a migration reference, and adding or changing
a field requires a reviewed schema revision and positive/negative fixtures.
Deprecated names are accepted only by the one-time v0 converter, never as two
simultaneously authoritative v1 keys. Rollback retains the original v0 files
and pinned runtime until the converted workspace passes its acceptance gate.

**R-PRECEDENCE:** upstream settings use upstream precedence. For aibox-owned
presentation preferences, apply shipped defaults → user configuration →
project policy → environment → invocation override, omitting unused layers
without reordering them. An explicitly user-local tmux/theme file is an
*owned native-tool override*, not a competing aibox configuration layer:
refresh never overwrites it, and an effective-value view reports that its
native-tool value wins at render time. The two levels (aibox intent and native
tool output) must be displayed separately; there is no silent promotion of a
project preference over project policy. Show effective values and provenance
with secrets redacted. Host policy is evaluated independently and can deny a
requested operation; project values cannot override it. Do not maintain
duplicate install pins in the UX namespace. This follows the company
application-configuration layer order, including authority separation.

**R-OWNERSHIP:** distribution assets are immutable/versioned; generated tool
outputs have declared ownership; user overrides are never overwritten.
Changing a project default leaves an explicit local override intact and
reports the difference. Native user config is first-class. A local refresh
validates first, stages managed output, atomically activates it where possible,
and leaves the previous working config on failure. It cannot regenerate
host/runtime definitions or run host hooks.

## Custom deserializer and legacy reconciliation

The census records raw types and serde attributes; these rules complete it:

- `appearance` aliases `customization`; concrete legacy theme names resolve
  to family/mode/variant without losing exact palettes. `legacy_theme` is an
  internal sidecar, not an independently editable user setting.
- `[aibox].project_name`/legacy metadata name, `[aibox].version/base`, `[image]`
  and `[container.image]` have compatibility precedence in v0; conversion
  compares their **effective** values and reports contradictions.
- `[context]` and `[process].packages` migrate to processkit context selection;
  top-level `[agents]`/`[mcp]` coexist with `[ai.agents]`/`[ai.mcp]`.
- `[ai].providers`, string `harnesses`, ordered detailed harness entries,
  per-harness tables and `harness_order` are all inputs. Detailed `enable` is
  an alias of `enabled`. Preserve install separately from integration enablement.
- `[ai.execution.<harness>]` and legacy `[ai.harness.<harness>.execution]`
  override individual policy axes; preserve deny/ask semantics, not merely booleans.
- Addon names are dynamic keys; nested language groups are recursively
  flattened groups, not a literal `groups` field. Each selected tool retains
  absent/explicit-enabled/explicit-disabled/version distinction.
- `[container.lifecycle]` and old top-level lifecycle fields coexist;
  `[audio]` and `[container.audio]` are compatibility forms.
- `.aibox-local.toml` environment overrides shared values, volumes append,
  personal MCP entries remain private. No secret may be committed during conversion.
- `local_env`/`local_mcp_servers` are runtime-derived, not user schema fields.
  Field-level aliases, including MCP `mode`/`extra_patterns`, tmux aliases and
  prompt names, are retained in `configuration-types.json`.
- `AIBOX_*` source occurrences include tests/internal markers. The environment
  census is a review checklist, not a promise to preserve unsupported test hooks.

## Persistence and migration contract

**R-MIGRATION:** offer a bounded migration/recovery workflow outside the normal
configuration API. It reads the old config/lock, native overrides, personal
state and named snapshots; produces proposed files plus a complete mapping,
warnings, conflicts and a scoped backup plan. Unknown fields block silent
conversion. No destructive action in preview; no in-place conversion without
user consent. Record source and destination versions/digests, never secrets.

Explicitly account for `.aibox-home`, login state, tool caches, `.aibox-env`
snapshots, shared/non-shared processkit context, extra MCP entries, generated
hooks, Dockerfile.local, Compose overrides, LaTeX outputs and logs. v0 named
environments snapshot config and process context: multiple devcontainer files
alone do not reproduce save/switch/restore behavior. G03 must select a tested
existing snapshot/delegation workflow before that feature can be retired.

**R-RECOVERY:** rebuild and ordinary removal preserve project/home/auth state.
Disabling a harness preserves its state by default. Purge identifies exact
paths, rejects traversal/symlink escape, reports unrecoverable effects and
requires fresh consent. A failed migration retains the original operational
workspace and supports restoring backup plus the old pinned runtime. Test
interruption and rollback after every state-changing stage. No credentials
in Git, images, shared evidence or public docs.

## Draft examples and limits

The following is schematic, not a runnable installation recipe or a claim
that the proposed Feature registry exists:

```jsonc
{
  "name": "project-workspace",
  "build": { "dockerfile": "Dockerfile" },
  "remoteUser": "aibox",
  "customizations": {
    "aibox": {
      "schemaVersion": "1",
      "workspace": {
        "theme": "gruvbox",
        "mode": "dark",
        "layout": "dev"
      }
    }
  }
}
```

Before specification acceptance, add schema-validated examples for minimal,
multi-harness, tool customization, audio, LaTeX and local-override workflows
with actual Feature references/digests and no hidden host prerequisites (G01/G02).

# 19. Configuration realization and parity implementation rules

This chapter defines configuration placement and conversion. The v0 ledger is
a source census, not a claim that regex-generated destinations implement a
converter. Implement these rules in an explicit, table-driven Go converter;
no catch-all unknown-field copy is allowed. Each concrete ledger path gets a
fixture, including unreachable fields and aliases. A mapping without a
verified destination blocks that **conversion**, not basic native v1 use.

## Placement and defaults

1. Native container topology, users, mounts, ports, environment, hooks and
   installation remain native Dev Container/Compose/Dockerfile/Feature inputs.
2. `customizations.aibox.workspace` retains the closed 115-field v0 UX
   vocabulary as an optional runtime compatibility interface. It is justified
   by live theme/layout switching, shared semantic palettes, status/title
   integration and one renderer coordinating multiple tools. It is not a
   second infrastructure configuration language. A native-only user may omit
   it entirely and edit each tool's configuration directly.
3. Do not also define theme/status/layout/prompt fields as Feature options.
   The workspace-runtime Feature installs the renderer/assets, not a competing
   UX preference store. Tool install/version choices belong to native Feature,
   image, Dockerfile or package inputs; `customizations.aibox` never installs tools.
4. A Feature installation cannot read an arbitrary workspace template. After
   workspace mount, a local lifecycle command invokes the renderer against the
   selected definition's validated namespace. Managed starters pass the
   selected container-visible config path explicitly to `aibox refresh`;
   the operator records host-to-container path mapping during native `up`.
   A direct-upstream starter uses its known workspace-relative path. Never
   guess among definitions or parse host files from the container. Missing
   selection gives an actionable diagnostic and leaves native defaults intact.
5. Explicit project values override distribution defaults; invocation-only
   UX overrides override project intent for that refresh. User-owned native
   tool files are loaded last and are never overwritten. Chapter 16's process
   settings may supply only the bounded theme/mode/layout overrides with the
   precedence in chapter 16; they do **not** merge arbitrary UX keys or grant authority.
6. Migration evaluates defaults/aliases using the pinned v0 semantics and
   writes explicit effective values, including false and empty values where
   meaningful. Native v1 defaults are immutable versioned distribution assets
   derived from `ledger/defaults.json`; new opt-ins sidebar/review are disabled.
   Upstream installer default drift must not change a migrated workspace.

Local LaTeX `engine`, `options`, `cache_dir`, `documents` and
`preview.document` remain local build/UX intent. Preview service presence,
image, network binding and port belong **only** to native Compose. The UX
schema excludes
`latex.preview.enabled/engine/bind/port/allow_public`.
`allow_public` becomes an operator-reviewed public-bind request, never a
local-agent permission. Default preview binds loopback; its container has
read-only completed-PDF mounts and no compiler/source credentials.

## Conversion destination rules

Apply aliases first using v0 precedence; conflicting explicit aliases fail
with both source paths. The canonical families below cover the ledger;
dynamic maps are expanded into individual entries in the conversion report.

| Canonical v0 input | Destination / transformation |
|---|---|
| `apiVersion`, `kind`, `aibox.config_schema` | Validate supported v0 discriminator; record source version in migration manifest, never emit as Dev Container properties. |
| `metadata.name`, `aibox.project_name`, `container.name` | Preserve distinct meanings: Dev Container `name` for display; native Compose `name`/`container_name` only when explicitly requested and collision-checked. Never use name alone as runtime identity. |
| `aibox.base/version`, `image.*`, `container.image.*` | Resolve aliases to `image` reference or Dockerfile `FROM`; record original value and selected version/digest. Local image overrides cannot silently change the accepted base. |
| `aibox.profile` | Select the minimal/curated Template and explicit Feature set; no runtime profile engine. Preserve headless versus human UX launch behavior in native command/lifecycle selection. |
| `container.user` | `remoteUser` plus matching Dockerfile/Compose user and home mount; preserve UID/GID semantics, not just spelling. |
| `container.hostname` | Compose service `hostname`, or Docker-native `runArgs` hostname for a single-container definition. Reject duplicate contradictory native declarations. |
| `container.environment` | Native `containerEnv` for non-Compose; Compose service `environment` for Compose. Private values become ignored env-file references, not committed JSON. |
| `container.extra_volumes[]` | Native mounts/Compose long volume syntax preserving source/target/read-only. Resolve on daemon host; require operator approval. |
| `container.lifecycle.keepalive` | Preserve the v0 bounded DNS/network keepalive through an explicitly selected container-local `postStartCommand` script; this is not Dev Container `overrideCommand` or permission to alter the main process. Merge with existing native lifecycle commands without dropping either. |
| `container.lifecycle.post_create_command` | Native `postCreateCommand`; preserve argv/string interpretation and trust classification; never run during migration. |
| `container.paths.*` | `devcontainer_json` selects output definition; `dockerfile` maps to `build.dockerfile`; Compose paths map to ordered `dockerComposeFile`; Dockerfile.local becomes an explicitly composed native Dockerfile input, not an unimplemented include; `local_env` maps to ignored native env-file usage. Detect relative-base changes. |
| `container.resource_thresholds.*` | `customizations.aibox.diagnostics.<leaf>`, preserving units and null/absent semantics. |
| `customization.*` | `customizations.aibox.workspace.<same suffix>`; canonicalize aliases before emitting the closed schema. |
| `ai.harness_order` / enabled harness list | `customizations.aibox.harnesses.order` and `launch.<id>.enabled`; installation is independent. |
| `ai.harness.<id>.install/version`, `addons.*.tools.*` | Explicit native Feature/image/Dockerfile/package selection, using the per-tool rule below. Host-only tools stay host-only and produce installation guidance, never fictitious in-container support. |
| `ai.execution.*`, harness execution overrides | Per-harness native approval/filesystem/network settings. Use adapter translations in chapter 11; no generic boolean approximation or claimed grant. A harness with no equivalent rejects that conversion with a manual disposition. |
| `ai.mcp.servers[]` | Each selected harness's native MCP config entries, by exact server name. Preserve argv/env references, disable collisions rather than replacing personal entries. Secrets are references under chapter 17. |
| `ai.mcp.permissions.*` | Native per-harness allow/ask/deny policy where expressible; explicit unsupported outcome otherwise. Deny wins; project values never authorize host operations. |
| `ai.mcp.gateway.*` | Existing optional gateway's own native launch arguments/config: bind, port, path, mode and lazy-catalog. Keep direct stdio default; do not add a gateway implementation to aibox. Qualification must name the gateway/version; unavailable selection blocks that conversion. |
| `ai.model_providers` | Native harness provider selection/credentials; status polling uses the separately closed workspace provider preferences. No provider credential copying into UX config. |
| `agents.*`, `ai.agents.*` | Canonical instruction filename and provider pointers in repository files; retain content and use supported processkit rendering where selected. Unknown custom paths require reviewed path mapping. |
| `processkit.*`, `context.*`, `process.packages` | Processkit Feature options/source lock plus its native configuration and supported installer interface; no aibox entity/schema engine. Context schema version is validated as source compatibility, never rewritten speculatively. |
| `skills.include/exclude` | Supported processkit skill selection/filter configuration; exclude wins. Preserve non-processkit skills as existing files, not delete-on-sync inventory. |
| `audio.*`, `container.audio.*` | Audio Feature presence for install; native `PULSE_SERVER`/mounts and bridge recipe for enable/backend/server. Disabled means no bridge request. See chapter 17 for credential/file boundaries. |
| `latex.*` | Local fields as above; preview fields map to `services.latex-preview.image`, service presence and `ports` long syntax (`host_ip`, `published`, `target`). Engine names must resolve to a qualified service image, not a guessed tag. |
| `integrations.github.credential_helper` | Native Git credential-helper configuration, user-owned and explicit; never embed a token or override a user's helper without consent. |
| `apply.preserve_disabled_harness_state`, `apply.purge_disabled_harness_state` | Retention manifest plus explicit separate native cleanup instructions. Ordinary refresh/remove always retain data; migration does not execute a legacy purge request. Mark the behavior change and require acknowledgement. |
| `security.acknowledge_seccomp_unconfined` | Operator policy request plus explicit native security option if approved. A copied acknowledgement is not renewed approval. |
| `local.container.*`, `local.mcp.*` | Merge using v0 local precedence before conversion; keep resulting private values in ignored native local overrides/provider files. Do not promote them into committed project defaults. |

Former `ai.providers`, `ai.harnesses[]`, top-level `mcp.*`, `audio.*`,
`agents.*`, and legacy execution forms are converter aliases, not v1 keys.
Record normalization and information loss explicitly. Missing/default fields
receive their actual v0 effective values, not the current upstream defaults.

## Tool installation realization

Account for each of the 98 `ADDON:<recipe>/<tool>` entries in the v0 census.
An installation unit may cover several tools; neither a Feature per tool nor
a Feature per language is required. Prefer an existing qualified native Dev
Container Feature when it preserves the relevant install/version/disable
behavior and license/provenance. Native image, Dockerfile and package units
are also valid, using existing package managers or official release artifacts.
Do not introduce an addon resolver or an installer in `customizations.aibox`.

Before implementing any aibox-owned installer Feature, present the **complete
gap list** to the owner. For each gap include researched existing native
options, the behavior they cannot preserve, maintenance/security cost, and
removal or scope-reduction options. The owner must decide the disposition;
a missing qualified Feature is not automatic permission to write one.
An owned installer requires individual justification and owner approval after
that full review (DEC-20261001_1954-SmoothTrout). It must reuse existing native
installation mechanisms rather than implement another package resolver.

The V1-05 installation manifest accounts for **every** ledger entry with a
selected qualified native reference/source unit or an explicit owner-approved
removal/scope decision. Preserve the original census row for removed tools so
parity accounting remains auditable. Unresolved gaps block catalog acceptance.
Record options, defaults, dependencies, license/provenance, immutable build
inputs, supported architectures and evidence per selected unit. Multiple rows
may point to the same unit; a row must not invent independent version or
disable support that the unit does not provide.

Native selection explicitly expands migrated default-enabled tools. Explicit
disable must produce absence where v0 promises it, including from the base
image. Package tools without selectable versions expose no invented promise:
the native image/package lock records installed versions. Hard prerequisites
use native dependencies or documented prerequisites; optional selections do
not acquire unexpected extras. Local source references may be used for tests;
published references and immutable digests require actual publication evidence.
Default/non-default supported versions, explicit disable and both container
architectures require tests. Download/checksum and wrong-architecture failures
are tested where the selected installation mechanism performs those operations.

## Local file ownership and recovery

Distribution assets live under `/usr/local/share/aibox/<version>/` and are
read-only. Generated UX lives in `$XDG_CONFIG_HOME/aibox/generated/`; its
manifest records relative path, generation ID, source/output hashes, owner
and mode. XDG defaults resolve to the **actual** user's home, not a hardcoded
`/home/aibox`. Tool-native starter files include generated fragments first
and user customization last. Existing native files are never overwritten;
if they lack an include, report inactive managed integration with the exact
include snippet and let the user edit it. Do not silently modify native files
just to obtain theme consistency.

Stage a complete generation in a private sibling directory and validate all
files before switching an atomically replaced `current` generation pointer.
Only this managed pointer may be a symlink; validate its owner, confinement
and target. Never follow user-supplied symlinks on output paths. Consumers
source `current/<tool>`; native user files remain outside the generation.
Retain the prior generation until reload completes. If reload fails, restore
the previous pointer and reload old config, reporting any failed recovery.
This is atomic **file selection**, not atomic rollback of external tool side
effects. A crash journal and `doctor` expose incomplete reload; a fresh
explicit refresh may recover it. Never restart a session or host container
as an implicit rollback action.

## Private named-state workflow

Named definitions are ordinary project-owned native files. Context and
credentials are not automatically safe to put in Git. For each migrated named
snapshot, produce a private export manifest in the operator's migration state
directory listing exact source paths/hashes/modes and destination definition.
Include the saved AGENTS/provider files and selected processkit context, not
only a config filename. Private content is copied only to an explicitly
chosen mode-0700 backup root; archive files are mode 0600 and encrypted when
they contain credentials. Use existing archive/encryption tools, not a new
aibox archive format or backup engine.

The guide specifies: stop writers; select exact paths; create native archive
and checksum; list/verify it; extract to a new private staging directory;
compare against current files; restore only approved paths after backing up
conflicts; select the native definition with `--config`; then invoke operator
`up`. Reject absolute/traversing archive members and symlink escapes before
extraction. Restore never executes archived scripts. Listing/detail reads the
manifest; deleting a definition never deletes archives, context or home.
Tests cover divergent context, private provider config, user edits since save,
interruption and recovery. Version control alone is not snapshot parity.

## Acceptance

AC-CONFIG/AC-MIG/AC-TOOLS require a coverage report joining every ledger row to
the applied rule, concrete generated destination, effective value, fixture
and evidence. Default-only smoke tests are insufficient. Any unsupported
conversion is an error with retained original data and a manual recipe; it
cannot be counted as achieved v0 parity. New features must not become mandatory
to hide lost v0 behavior. Record artifact/package selections before implementing
the affected adapter; verify target support against the resulting artifacts.

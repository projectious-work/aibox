# v0 command disposition

Source baseline: v0.35.0 release-line CLI declarations, dispatcher, named-environment
implementation and command documentation at
`9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6`:
[cli.rs](https://github.com/projectious-work/aibox/blob/9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6/cli/src/cli.rs),
[main.rs](https://github.com/projectious-work/aibox/blob/9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6/cli/src/main.rs),
[env.rs](https://github.com/projectious-work/aibox/blob/9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6/cli/src/env.rs),
[command reference](https://github.com/projectious-work/aibox/blob/9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6/docs-site/content/docs/reference/cli-commands.md).
This is a source-level surface audit, not execution of every command.

The visible verb count is modest, but `apply`, `set`, `reset` and the generic
resource grammar cover extensive generation, configuration mutation and
processkit administration. The whole existing surface is therefore not
already a thin lifecycle wrapper. Preserve the useful operations and user
outcomes while reducing ownership.

The source includes `emergency`, `prune`, `apply generated-runtime`,
`up --forget-tmux-state` and additional init flags missing from the command
reference. A docs-only inventory would have missed these.

## Disposition by current command family

“Remove” below means remove the interface from the normal aibox core, not
silently discard its user outcome. Remaining parity gaps block replacement
until a usable migration or delegated workflow is demonstrated. Configuration mutation commands are replaced by manual/agent file editing.

| Current surface | Disposition | Destination or replacement |
|---|---|---|
| `init` | Reduce/delegate | Standard Template application or a copyable starter. No aibox configuration wizard and no duplicate flags for themes, prompts, harnesses, per-tool pins or processkit versions. A one-time scaffold convenience command is optional, not required core. Subsequent edits are ordinary file edits. |
| bare `apply`, `--no-cache`/`--rebuild` | Split | Explicit build/rebuild operations delegated upstream; do not retain a catch-all that rewrites config, updates processkit and builds images. |
| `apply --config-only`, `--no-container`, `apply generated-runtime` | Separate local work from host work | Pure validation and explicit managed-runtime refresh as appropriate. Native Dev Container files become authoritative; do not regenerate them from a parallel model. Local-only functions must never probe a host runtime. |
| `apply --standardize-config` | Remove from normal core | One-time v0 migration utility with preview, backup and unknown-field reporting; no recurring config normalizer. |
| `apply --fix-compliance-contract`, `apply migration` | Delegate | processkit's own supported interfaces. No duplicate migration/compliance engine. |
| `apply audio` | Move out of generic apply | Host setup guidance using existing platform tools; separately authorized bounded host helper only if needed for equivalent usability. Never available inside the container. Container audio enablement remains declarative. |
| `up`, `up --apply` | Keep start operation; decouple attachment | MCP start returns a bounded structured result. CLI may offer convenient attach after start. No implicit config rewrite/content upgrade via `--apply`. |
| `up --layout`, `--forget-tmux-state` | Move to local session workflow | Edit layout preferences and reload, or explicitly recreate the current tmux layout. Session reset may interrupt processes and requires consent; it is not a container restart. |
| `down` | Keep | Stop the identified environment without removing persistent data. Delegate to the appropriate native runtime/Compose operation. |
| `delete runtime` | Keep | Explicit remove operation, separate from stop. Show exact resources; preserve project/auth/home data by default. No ambient global prune. |
| `emergency <harness>` | Preserve recovery outcome | Attach/execute using an explicit recovery path that bypasses tmux, Yazi and status tooling. Prefer existing shell/harness commands; a separate verb only if justified. |
| `get runtime`, `describe runtime`, `--resources` | Keep/consolidate | Environment inspection plus local cgroup/process inspection, each in the correct authority context. Preserve useful structured outputs; label unavailable measurements honestly. |
| `get/describe addon`, `addon-catalog`, `provider-backends` | Move catalog to resources | Versioned capability/catalog documentation and schemas. Live installed-component inspection remains an operation where useful. No second package registry. |
| `describe workspace-manifest`, `image-provenance-policy` | Retain information, simplify projection | Read-only effective-config/runtime/provenance inspection using upstream-resolved configuration and release metadata; no new desired-state model. Existing schema consumers need migration guidance. |
| `set ...`, `edit config`, `delete addon` | Remove | Human/agent edits `devcontainer.json` and native config files. Guidance and validation assist; no generic configuration-write MCP tool either. |
| `get/describe/set/delete skill/process/kit`, `get/set/apply/delete migration`, skill-category filters | Delegate | processkit catalog, configuration and entity APIs. aibox can explain and verify integration, but does not administer process entities. |
| `create/get/describe/apply/delete env` | Delegate to native definitions and private archives | v0 snapshots `aibox.toml`, `AGENTS.md`, `CLAUDE.md` and non-shared `context/` into `.aibox-env`, then restores them. Multiple Dev Container configs alone are not equivalent. Use native definition variants and the private manifest/archive workflow in chapters 12/19; preserve process context and test recovery. |
| `create backup`, `reset project` | Separate recovery/migration support | Preserve scoped backup and recovery via established tools or a bounded utility. Do not put scaffold deletion/recreation into the ordinary lifecycle core. No removal until equivalent workflow is documented and tested. |
| `reset context` | Delegate | Already plan-only in v0; processkit owns context recovery. Do not expand it into a new destructive aibox operation. |
| `prune` scopes | Narrow/delegate | Host container cleanup only for identified owned resources; local cache cleanup through existing tools. Provider worktree cleanup requires clean/dirty checks and explicit consent; no broad `all` operation in initial core. Retain safe reclaim capability via documented workflows. |
| `doctor`, `--integrity` | Keep, make read-only by default | Schema, semantic consistency, dependency/asset compatibility and scoped runtime checks. Structured findings and remedies; no migration creation or automatic fixes as a side effect of checking. |
| `doctor audio` | Keep scoped diagnostics | Container checks inside; host readiness checks outside. Explicit skipped/not-authorized findings, not false healthy results. |
| `doctor security` | Delegate scanners | Invoke or document established scanners; report their exact scope/status. No new security scanner, credential harvest or implied audit assurance. |
| `self update`, `self uninstall`, `--purge` | Prefer distribution tooling | Use the chosen installer/package manager; preserve check/update/uninstall workflows. Avoid self-modification tools in the operational MCP server. Purge remains separately scoped and confirmed. |
| `self completion`, help/version | Retain CLI conveniences | Generate from the command definitions; not meaningful MCP mutations. |
| LaTeX build/watch and preview helpers | Retain container-local workflows | Existing tools/scripts and read-only companion behavior; no host management through local helper commands. |

## v1 operation contracts

Chapters [3](03-interfaces.md), [10](10-interfaces-and-guidance.md) and
[18](18-operational-contracts.md) define command grammar, MCP tools, schemas,
authority and failure behavior. Chapter [12](12-migration-and-state.md)
defines migration/recovery. Verify every original action and argument against
its replacement workflow; source inventory alone does not prove parity.

> Historical review input, copied for self-contained specification review. Owner confirmed boundary review points 1–5 and selected Go / thin MCP-first wrapper on 2026-09-25 (DEC-HopefulTower). Earlier open-choice language below records the review history; README.md and numbered chapters define the current draft. This snapshot does not authorize implementation.

# aibox v1 thin MCP/CLI interface audit

Date: 2026-09-25. Status: recommendations for boundary review, not an approved
command specification. Companion to the
[product boundary](boundary.md).

## Evidence and conclusion

Audited the v0.35.0 release-line CLI declarations, dispatcher, named-environment
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
until a usable migration or delegated workflow is demonstrated. The owner's
config-file-only instruction explicitly replaces configuration mutation
commands with manual/agent editing.

| Current surface | Proposed disposition | Destination or replacement |
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
| `create/get/describe/apply/delete env` | Redesign after parity analysis | v0 snapshots `aibox.toml`, `AGENTS.md`, `CLAUDE.md` and non-shared `context/` into `.aibox-env`, then restores them. Multiple Dev Container configs alone are not equivalent. Prefer existing config variants plus product-owned context/snapshot tooling; require preservation/recovery proof before retiring commands. |
| `create backup`, `reset project` | Separate recovery/migration support | Preserve scoped backup and recovery via established tools or a bounded utility. Do not put scaffold deletion/recreation into the ordinary lifecycle core. No removal until equivalent workflow is documented and tested. |
| `reset context` | Delegate | Already plan-only in v0; processkit owns context recovery. Do not expand it into a new destructive aibox operation. |
| `prune` scopes | Narrow/delegate | Host container cleanup only for identified owned resources; local cache cleanup through existing tools. Provider worktree cleanup requires clean/dirty checks and explicit consent; no broad `all` operation in initial core. Retain safe reclaim capability via documented workflows. |
| `doctor`, `--integrity` | Keep, make read-only by default | Schema, semantic consistency, dependency/asset compatibility and scoped runtime checks. Structured findings and remedies; no migration creation or automatic fixes as a side effect of checking. |
| `doctor audio` | Keep scoped diagnostics | Container checks inside; host readiness checks outside. Explicit skipped/not-authorized findings, not false healthy results. |
| `doctor security` | Delegate scanners | Invoke or document established scanners; report their exact scope/status. No new security scanner, credential harvest or implied audit assurance. |
| `self update`, `self uninstall`, `--purge` | Prefer distribution tooling | Use the chosen installer/package manager; preserve check/update/uninstall workflows. Avoid self-modification tools in the operational MCP server. Purge remains separately scoped and confirmed. |
| `self completion`, help/version | Retain CLI conveniences | Generate from the command definitions; not meaningful MCP mutations. |
| LaTeX build/watch and preview helpers | Retain container-local workflows | Existing tools/scripts and read-only companion behavior; no host management through local helper commands. |

## Suggested minimal surface

Names below are illustrative; existing familiar names can be retained. The
specification should settle vocabulary once, rather than proliferating aliases.

| Capability | CLI example | MCP exposure / authority |
|---|---|---|
| Start/create from existing configuration | `up` | `start_environment`; external manager only. |
| Stop, retain state | `stop` (or existing `down`) | `stop_environment`; external manager only. |
| Remove identified runtime | `remove` (or `delete runtime`) | `remove_environment`; external manager, explicit destructive scope. |
| Build / explicitly rebuild | `build`, `rebuild` | External manager; separate build-only from replacement/restart effects. |
| List/inspect state and effective configuration | `status`, `inspect` | Scoped read-only tools; local view never acquires host access. |
| Read diagnostic logs | `logs` | Bounded, redacted local or operator view; no unrestricted host file reader. |
| Validate consistency and readiness | `doctor` | Structured read-only diagnostics; local/operator check sets. |
| Attach/recovery | `attach` | Human interactive convenience. No interactive tmux attachment over MCP; return connection guidance. General command execution is not required in the initial MCP surface. |
| Refresh local managed UX outputs after edits | `refresh` if needed | Bounded local operation; preserve user overrides, report rebuild-required without doing it. Prefer native reload commands when sufficient. |
| Serve MCP | `mcp serve` | Server entry point, not a tool calling itself. Local and operator capability sets enforced externally; an in-container caller cannot enable operator mode. |
| Find/read task guidance, catalog and schemas | Optional `help`/`how-to` view | Canonical MCP resources/templates; compact discovery index, optional bounded search. |

No configuration `set`, tool-install wrapper, processkit entity mutation,
generic host shell tool, new scheduler or public adapter plugin system.
Non-core migration/recovery workflows remain tracked parity obligations.

The CLI and MCP share one interface-neutral core and structured results.
Use an established MCP SDK and reuse shared documentation-resource machinery
where available. Human progress rendering must not contaminate machine output.
Define errors, timeouts, cancellation, partial failures and safe retry behavior
for delegated operations. Long-running builds need a supported completion and
cancellation contract, not an excuse to invent a general workflow engine.

## Doctor recommendation

Keep it. Config-file editing increases the value of understandable validation.
Check JSONC/schema validity, supported Feature/options and aibox customization
schema, contradictory settings, missing declared assets, dependency/version
compatibility and appropriate runtime readiness. Prefer upstream resolution
and validators; own only aibox-specific checks.

Schema validity is not proof of safe configuration. Surface authority-sensitive
changes such as host hooks, mounts and privilege requests separately. Resolve
configuration for checking without executing lifecycle hooks; any network
resolution must be explicit. Diagnostics must not install dependencies, start
containers, rewrite files or “fix” authorization policy.

Each finding needs a stable identifier, severity, file/setting location,
explanation and specific remedy/how-to link. Distinguish failed, unavailable,
skipped and unauthorized checks. An unavailable host check inside the container
is expected and does not license a host bridge.

## Open acceptance work

1. Map every v0 setting, including user-local overrides, to standard fields,
   Feature options, referenced native files or a schema-governed aibox namespace.
2. Establish equivalent recovery, snapshot switching, audio setup and update
   workflows before deleting old convenience interfaces.
3. Verify thin delegation against a pinned upstream CLI, especially stop/remove
   and resource identity; do not infer semantics from verb similarity.
4. Test native tool-only use and both authority contexts; verify that local
   agents cannot invoke host operations even with crafted inputs.
5. Review final surface and new UX dependencies with the owner before the
   specification becomes an accepted implementation baseline.

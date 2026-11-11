# Layered settings and diagnostic sinks: overview

This file is not normative. The contract is `requirements.md`,
`interfaces.json`, `journeys.md`, the schemas and the fixtures.

## Goal

An operator or a local agent can change how the aibox process behaves
(output format, diagnostics, time limits) at the level that fits: a machine,
a user, a project, one shell, or one command. They can then ask aibox which
value won and where it came from. Diagnostics never mix with results, and
never carry a secret.

The owner's reasons, from chapter 16 of the v1 product description and the
roadmap item V1-06:

- Later phases (credentials V1-07, lifecycle V1-08, full CLI and MCP V1-09)
  need one settings resolver and one logging path before they add
  operations. Building them per operation would give each its own precedence.
- A result on stdout must stay machine-readable under every setting, so that
  agents can rely on it.
- A project file or an environment file must not be able to move logs to a
  new place or weaken an administrator's choices.

## Scope of version 1

In scope, for the executable as it exists after V1-05 (commands `inspect`,
`mcp serve`, `version`, `help`; local context only):

- the seven settings layers, their files, the environment mapping, the
  bootstrap selectors and the common flags;
- the effective-configuration view through `inspect --effective-config` and
  the MCP tool `inspect_workspace` with `effectiveConfig`;
- result formats `human`, `json` and `yaml` for `inspect`;
- the semantic log events, the stderr sink and the rotating file sink, with
  redaction, bounded queueing and failure handling.

Out of scope, with the phase that takes it:

| Topic | Taken by | Why not here |
|---|---|---|
| Operator policy bounds on settings, sinks and retention | V1-08, V1-09 | There is no operator mode yet. |
| Refusing a mutating operation when a required sink fails | V1-08 | There is no mutating operation yet. The rule is stated here (LOG-031) and verified at component level only. |
| `--config` and the request field `configPath` | V1-08 | Selecting another native definition changes inspection, not settings. |
| `--theme`, `--mode`, `--layout` on `refresh` | V1-11, V1-12 | `refresh` does not exist yet. The three keys are resolved and reported here. |
| `aibox logs` and the operation `read_logs` | V1-08, V1-16 | Reading runtime logs needs a runtime. |
| Converting v0 preferences | V1-20 | Migration. |
| Service-manager capture (systemd, launchd) | V1-09 | Documentation of the operator server. |
| macOS paths on a real macOS host | V1-22 | Target qualification. The path rules are specified and unit-tested here. |

## Boundaries

- Settings configure the aibox process. They never rewrite
  `devcontainer.json`, Features, images, mounts or hooks.
- Operator policy and credentials are not settings and are not layers.
- Operation records and evidence are a separate store. Log rotation never
  touches them.

## Decisions taken in this draft

Each is a recommendation that the owner confirms at the freeze. They close
gaps in chapter 16 that an acceptance test author would otherwise have to
guess.

| # | Question | Draft answer |
|---|---|---|
| D1 | Exact environment names. "Uppercase snake-case from camelCase" is ambiguous for `rotationMiB`. | The names are listed per key in `interfaces.json`; `rotationMiB` is `AIBOX_LOGGING__ROTATION_MIB`. |
| D2 | Error codes. Chapter 18 keeps the code list closed. | No new result code. Shape errors are `invalid_input` (exit 2); a value from a source that may not set the key is `not_authorized` (exit 3). The cause is the `reason` field of the `settings.rejected` log event, a closed list. |
| D3 | Where the view reports loaded and absent files and overridden values. The result schema has no field for them. | Two optional fields are added to the `inspect_workspace` data: `settingsSources` and, per entry, `overridden`. They are additive. |
| D4 | "Rejected settings" in the view. | A rejected setting fails the operation. The view never shows a value next to a rejection; the error and the log event name the setting. |
| D5 | `execution.timeoutSeconds` has two defaults (30 s read, 3600 s mutating). | The key has no built-in value. Unset, the limit is the class default of the operation; the view reports the limit that applies to `inspect` (30) with layer `built_in`. |
| D6 | Default format of `inspect`. It was `json` only. | `human` is the default, as chapter 16 says. Existing callers pass `--format json`. Human wording is not a contract. |
| D7 | YAML. | Implemented without a new dependency: block-style YAML whose scalars are JSON scalars. Parsing it gives the same value as the JSON form. |
| D8 | Permitted place of the log file without operator policy. | Only inside `<state>/aibox/logs/`. Anything else is `not_authorized`. Operation records live beside it in `<state>/aibox/operations/`. |
| D9 | Rotation naming and retention. | `<file>.1` is the newest rotated file. At most `retentionFiles` files exist including the active one. Rotated files older than `retentionDays` are deleted at rotation and at start. |
| D10 | What counts as a known secret. | The value of every process environment variable whose name matches the pattern in `interfaces.json` (`redaction.secretNamePattern`) and is at least 8 bytes long, plus the listed secret shapes. The home directory prefix of a path is written as `~`. |
| D11 | Unknown `AIBOX_*` names in the environment. | They fail, as chapter 16 says, except the 83 registered runtime identifiers in `fixtures/registered-environment.json` and names under the prefix family `AIBOX_TMUX_`. This matters inside v0 containers, which set many of them. |
| D12 | How a test sets a system layer. | No override exists; that would let a caller suppress the system layer. Journeys that need `/etc/aibox/settings.json` run in a disposable container. |
| D13 | Project layer. | Only `customizations.aibox.workspace.theme`, `.mode` and `.layout` of the selected `devcontainer.json` take part, and only for the three `workspace.*` keys. |
| D14 | MCP logging. | Fixed when `mcp serve` starts. A tool call can set `effectiveConfig` and nothing else from this feature. |

## Sources

- v1 product description, chapters 2, 3, 13, 16 and 18, and acceptance groups
  AC-CONFIG-SOURCES and AC-LOGGING in chapter 5 (`spec/v1/`, at `main`
  `2b542a4e`).
- Requirements R-CONFIG-SOURCES and R-LOGGING; R-LIMITS for the bounds.
- Company standards v2.3.0: application configuration; application output,
  logging and evidence; command-line interface design.
- `operation-request.schema.json` and `operation-result.schema.json` of the
  operation core (V1-03).

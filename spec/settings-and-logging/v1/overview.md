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
| Service-manager capture of the operator server | V1-09 | There is no operator server yet. Capturing the stderr of the local `mcp serve` under systemd or launchd is documented with V1-06; it needs no behaviour beyond LOG-011. |
| Operator-mode rejection of `workspace.*`; service-account user layer; settings in the frozen input digest; child-tool attribution; operation records failing closed | V1-08, V1-09 | They need operator mode, a mutating operation or a child process. |
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
guess. Rows marked in bold change, tighten or narrow the baseline.

| # | Question | Draft answer |
|---|---|---|
| D1 | Exact environment names. "Uppercase snake-case from camelCase" is ambiguous for `rotationMiB`. | The names are listed per key in `interfaces.json`; `rotationMiB` is `AIBOX_LOGGING__ROTATION_MIB`. |
| D2 | Error codes. Chapter 18 keeps the code list closed. | No new result code. Every failure before an operation starts, including a value from a layer that may not set the key, is `invalid_input` with exit 2, as chapters 16 and 18 say for bootstrap failures. The cause is the `reason` of the `settings.rejected` event, a closed list of 15. |
| D3 | Where the view reports loaded and absent files and overridden values. The result schema has no field for them. | **Specification change.** Two optional fields are added to the `inspect_workspace` data: `settingsSources` and, per entry, `overridden`. `schemas/operation-result.schema.json` in this version is the chapter 18 schema with that addition; a strict validator of the old schema rejects the new fields. |
| D4 | "Rejected settings" in the view. | **Specification change.** Chapter 16 and the company configuration standard say the view reports rejected settings with reasons. Here a rejected setting fails the command, and the first rejection is reported with its reason; the view never shows a value beside a rejection. Reason: reporting a result while a setting is rejected would let a caller act on values that the operation itself would refuse. |
| D5 | `execution.timeoutSeconds` has two defaults (30 s read, 3600 s mutating). | The key has no built-in value. Unset, the limit is the class default of the operation; the view reports the limit that applies to `inspect` (30) with layer `built_in`. |
| D6 | Default format of `inspect`. It was `json` only. | `human` is the default, as chapter 16 says. Existing callers pass `--format json`. Human wording is not a contract. |
| D7 | YAML. | Implemented without a new dependency: block-style YAML whose scalars are JSON scalars. Parsing it gives the same value as the JSON form. |
| D8 | Permitted place of the log file without operator policy. | **Tightening.** Only inside `<state>/aibox/logs/` (chapter 16 anchors to `<state>/aibox/`). Operation records live beside it in `<state>/aibox/operations/`, so a log path can never name a record. |
| D9 | Rotation naming and retention. | `<file>.1` is the newest rotated file. At most `retentionFiles` files exist including the active one. Rotated files older than `retentionDays` are deleted at rotation and at start. |
| D10 | What counts as a known secret. | The value of every process environment variable whose name matches `redaction.secretNamePattern` and is at least 8 bytes long, plus the listed secret shapes. The home directory prefix of a path is written as `~`. **Limits the owner accepts:** a secret shorter than 8 bytes is not redacted; a secret that is in no such variable and has no listed shape is not recognized; `target.workspaceRoot` in the JSON and YAML envelope is not redacted, because it identifies the workspace. |
| D11 | Unknown `AIBOX_*` names in the environment. | They fail, as chapter 16 says, except the 82 registered runtime identifiers in `fixtures/registered-environment.json` and names under the prefix family `AIBOX_TMUX_`. This matters inside v0 containers, which set many of them. |
| D12 | How a test sets a system layer. | No override exists; that would let a caller suppress the system layer. Journeys that need `/etc/aibox/settings.json` run in a disposable container. |
| D13 | Project layer. | **Narrowing.** Only `customizations.aibox.workspace.theme`, `.mode` and `.layout` of the selected `devcontainer.json` take part, and only for the three `workspace.*` keys. Chapter 16 says "UX/diagnostic intent"; no diagnostic key is admitted from a project. |
| D14 | MCP logging. | Fixed when `mcp serve` starts. A tool call can set `effectiveConfig` and nothing else from this feature. |
| D15 | The roadmap demo for V1-06 runs `inspect` without `--effective-config` and expects settings and provenance to be visible. | **Roadmap change.** The demo command gains `--effective-config` (J-CFG-24). Without the flag the result stays as it is today. |
| D16 | Code for refusing a mutating operation when a required sink has failed. Chapter 16 names none. | `operation_record_unavailable` (exit 5). This binds V1-08. |
| D17 | `mcp serve` and `logging.format` `plain`. | `plain` from any layer, including a settings file, stops `mcp serve`. "Explicit" in chapter 16 is read as "set by some layer". |
| D18 | Private mode for environment files. | **Relaxation, by scope.** Chapter 16 requires it for files that hold secrets. No setting of this version is a secret and only setting names are allowed, so no mode is required. The rule returns with the first secret-bearing key. |
| D19 | Operation name of a failure before the operation starts. | `invalid_request`, as chapter 18 says for bootstrap failures, with a generated request ID. |
| D20 | `version` and `help`. | They do not resolve settings and cannot fail on them. |
| D21 | Flush reporting. | Reported only when the flush did not complete (`logging.flush_incomplete`), so that a normal run has a fixed event list. |


## Sources

- v1 product description, chapters 2, 3, 13, 16 and 18, and acceptance groups
  AC-CONFIG-SOURCES and AC-LOGGING in chapter 5 (`spec/v1/`, at `main`
  `2b542a4e`).
- Requirements R-CONFIG-SOURCES and R-LOGGING; R-LIMITS for the bounds.
- Company standards v2.3.0: application configuration; application output,
  logging and evidence; command-line interface design.
- `operation-request.schema.json` and `operation-result.schema.json` of the
  operation core (V1-03).

## Review

- Round 1, 2026-11-11: independent review by `sage` (model `claude-opus-5-5`,
  Anthropic) of commit `5b2d4fe8`: not test-ready, 11 blocking and 12
  non-blocking findings. All 11 blocking findings and findings N1 to N11 are
  applied in the following commit. N12 (the `spec/v1` layout error of
  `check-specs`) is tracked in issue #504.
- Round 2, 2026-11-11: same reviewer, commit `34d35d77`: all 11 blocking
  findings of round 1 closed; 2 new blocking findings (B12, B13) and 8
  non-blocking (N13 to N20), all applied before the freeze. One point stays
  open by decision of the integration owner: aibox does not check the
  ownership of the system settings file (round 1, N8); the operator-policy
  feature takes it. The reviewer ran no journey and no YAML parser.
- Owner approval, 2026-11-11: decisions D3, D4, D6 and D15 as drafted, and
  the merge of this version with the remaining decisions as stated.
- Not covered by this version and not deferred to a named phase: backup
  exclusions for log files, and redaction of provider prompts and personal
  data. Neither occurs in the events of this version; the first feature that
  logs such data specifies them.

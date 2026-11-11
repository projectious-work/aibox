# Layered settings and diagnostic sinks: acceptance journeys

Normative. Each journey is a black-box test of the built `aibox` executable
through the CLI or through an MCP client on stdio. Identifiers, values,
codes and paths come from `interfaces.json`.

## Common setup

- Every journey runs with a fresh temporary directory as `HOME`, with
  `XDG_CONFIG_HOME` and `XDG_STATE_HOME` unset unless the journey sets them,
  and with an environment that holds only `HOME`, `PATH` and what the
  journey adds. No journey reads the developer's own settings.
- "Project P" is a copy of `spec/v1/examples/minimal`. "Project W" is a copy
  of it whose `customizations.aibox.workspace` sets `theme` to `nord`, `mode`
  to `dark` and `layout` to `dev`.
- "View" means `aibox inspect --context local --project <project>
  --effective-config --format json`, and `entry(k)` is the element of
  `data.configuration` with `key` k.
- A journey marked **root** needs `/etc/aibox/settings.json` and runs in a
  disposable container. All others run as an ordinary user.
- Unless stated, expect exit 0, `outcome` `succeeded`, and a stdout that is
  exactly one JSON document.

## Settings

| ID | Covers | Steps | Expected |
|---|---|---|---|
| J-CFG-01 | CFG-001, CFG-004, CFG-060, CFG-061, CFG-063, EXE-003 | View of P with no settings anywhere. | 15 entries in catalog order. Every `layer` is `built_in`, `source` `built-in`, `value` the CLI default, `invocationOnly` false, no `overridden`. `entry(execution.timeoutSeconds).value` is 30. `settingsSources` has system `absent`, user `absent`, project `loaded`. |
| J-CFG-02 | CFG-001, CFG-062 | For each key and each ordered pair of layers (lower, higher) from the key's `layers`, set a different valid value in both; View. | The higher layer's value wins; `overridden` holds the lower one. System-layer pairs are **root**. |
| J-CFG-03 | CFG-001, CFG-062 | **root.** `output.color`: system `never`, user `always`, env file `never`, environment `always`, flag `--color never`; View. | Value `never`, layer `invocation`, source `--color`, `invocationOnly` true; `overridden` lists environment, env_file, user, system in that order. |
| J-CFG-04 | CFG-002, CFG-003, LOG-006 | For each of `logging.file`, `logging.required`, `logging.rotationMiB`, `logging.retentionFiles`, `logging.retentionDays`: set its env name once in the environment and once in an env file; run `inspect --format json`. | Exit 3; one error envelope with code `not_authorized`; a `settings.rejected` event with reason `unauthorized_source`, the key, and layer `environment` or `env_file`. No lower-layer value is used. |
| J-CFG-05 | CFG-005, CFG-012 | Snapshot the file tree under `HOME` and the project; View; snapshot again. | The trees are identical. |
| J-CFG-06 | CFG-010, CFG-011 | **root** for system. Write `{"schemaVersion":"1","logging":{"level":"debug"}}` to the system path, then instead to the user path, then to `$XDG_CONFIG_HOME/aibox/settings.json` with `XDG_CONFIG_HOME` absolute; View each time. | `entry(logging.level)` is `debug` with layer `system`, `user`, `user`; the `source` is the file path with the home prefix as `~`. A relative `XDG_CONFIG_HOME` is ignored. |
| J-CFG-07 | CFG-013 to CFG-016 | User settings file, one case each: mode 000; text `not json`; unknown member `operatorPolicy`; `schemaVersion` `2`; no `schemaVersion`; `logging.level` `verbose`; `logging.rotationMiB` 0; the key `logging` twice; `output` set to `null`. Run `inspect --format json`. | Exit 2, code `invalid_input`, and reason `unreadable`, `malformed`, `unknown_key`, `unsupported_version`, `unsupported_version`, `invalid_value`, `invalid_value`, `duplicate_key`, `null_value` respectively. The operation did not run: no `operation.started` event. |
| J-CFG-08 | CFG-017, CFG-018 | **root.** System sets `logging.level` `warn`; default user file sets `output.color` `always`; `other.json` in the working directory sets `output.color` `never`; View with `--settings-file other.json` from that directory. Then with `--settings-file missing.json`. | First: `output.color` is `never` from `other.json`; `logging.level` is `warn` from system; the default user file is not in `settingsSources`. Second: exit 2, reason `missing_selected_file`. |
| J-CFG-09 | CFG-019, CFG-022 | Put a `settings.json` with `logging.level` `trace` into the project root, into `.devcontainer/`, and into the working directory; put `customizations.aibox.logging.level` `trace` into the project's `devcontainer.json`; View. | `entry(logging.level)` stays `info`, layer `built_in`. |
| J-CFG-10 | CFG-020, CFG-021 | View of W. Then W with `layout` `grid`. | First: the three `workspace.*` entries have the W values, layer `project`, `domain` `workspace`, source the `devcontainer.json` path. Second: exit 2, reason `invalid_value`. |
| J-CFG-11 | CFG-030, CFG-031 | A `.env` in the working directory and in the project sets `AIBOX_LOGGING__LEVEL=trace`; View. Then files `a.env` (`…LEVEL=debug`) and `b.env` (`…LEVEL=warn`) with `--env-file a.env --env-file b.env`; then with `AIBOX_ENV_FILE` set to the absolute path of `a.env` and no flag; then with both `AIBOX_ENV_FILE` (a) and `--env-file b.env`. | `info`/`built_in`; then `warn` from `b.env` with `debug` in `overridden`; then `debug` from `a.env`; then `warn` from `b.env` and `a.env` not in `settingsSources`. |
| J-CFG-12 | CFG-032 to CFG-038 | One case each, `inspect --format json`: `AIBOX_ENV_FILE=rel.env`; env file with line `export AIBOX_LOGGING__LEVEL=debug`; with line ` AIBOX_LOGGING__LEVEL=debug` (leading space); with the same key twice; path missing; path is a directory; mode 000; mode 0644; key `AIBOX_NOPE=1`; key `PATH=/x`. | Exit 2 with reasons `relative_env_file`, `malformed`, `malformed`, `duplicate_key`, `missing_selected_file`, `not_regular_file`, `unreadable`, `insecure_permissions`, `unknown_key`, `malformed`. |
| J-CFG-13 | CFG-033 | Env file (mode 0600) with CRLF line endings, a blank line, a `# comment` line and `AIBOX_WORKSPACE__THEME=nord`; then with value `"nord"` in quotes; View. | First: theme `nord`, layer `env_file`. Second: exit 2, `invalid_value` (quotes are literal). |
| J-CFG-14 | CFG-039 to CFG-042 | One case each: `AIBOX_LOGGING__LEVEL=debug` (View); `AIBOX_LOGING__LEVEL=debug`; `AIBOX_TMUX_SESSION=x` and `AIBOX_VERSION=0.35.4` and `AIBOX_LOG_LEVEL=trace` together (View); `AIBOX_DIAGNOSTICS__OFFLINE=TRUE`; `AIBOX_EXECUTION__TIMEOUT_SECONDS=+5`; `…=0`; `AIBOX_OUTPUT__COLOR=`. | `debug`/`environment`/source `AIBOX_LOGGING__LEVEL`; exit 2 `unknown_environment_name` and the diagnostic contains `AIBOX_LOGGING__LEVEL`; all entries `built_in`; then four times exit 2 `invalid_value`. |
| J-CFG-15 | CFG-050, CFG-051 | View with `--log-level debug`, with `--log-level=debug`, with `--timeout 5`; `inspect --log-level loud`; `inspect --timeout 0`; `inspect --format json --format yaml`; `version --format json`. | `debug`/`invocation`/`--log-level` twice; `execution.timeoutSeconds` 5; then exit 2 `invalid_value` twice; exit 2 for the repeated flag; exit 2 for a flag on a command that does not list it. |
| J-CFG-16 | CFG-052 | From a directory D, `inspect --project sub/p --format json` where `D/sub/p` is P. | `target.workspaceRoot` is the canonical path of `D/sub/p`. |
| J-CFG-17 | CFG-007, CFG-065 | Same user settings file and environment; View through the CLI, and call `inspect_workspace` with `effectiveConfig` true through `mcp serve`. | `data.settingsSources` is equal. `data.configuration` is equal except the entries for `output.format` and `logging.format`, which carry their MCP defaults. |
| J-CFG-18 | CFG-054 | `AIBOX_LOGGING__FORMAT=plain aibox mcp serve --context local`; and `mcp serve --log-format plain`. | Exit 2; stdout is empty; stderr has one bootstrap diagnostic with reason `invalid_value`. |
| J-CFG-19 | CFG-055 | Start `mcp serve` with `--log-level warn`. Call `inspect_workspace` with an extra input member `logLevel`, then `logging`, then `outputFormat`. Then call it correctly. | The three calls are rejected as invalid input by the tool schema. The correct call succeeds and stderr has no event below `warn`. |
| J-CFG-20 | CFG-066 | View of a project whose `devcontainer.json` has `initializeCommand` and `postCreateCommand` that would create a marker file, with `PATH` holding fake `docker`, `podman` and `devcontainer` executables that create a marker when run. | No marker exists. |
| J-CFG-21 | CFG-067 | `inspect --format json` on P without `--effective-config`. | `data` has no `settingsSources`; `data.configuration` equals the result of the same command built from the commit before this feature (golden file). |
| J-CFG-22 | CFG-064, RED-003 | User file sets `logging.file` to `run.log`; `EXAMPLE_API_TOKEN=canary-0123456789` in the environment and `AIBOX_WORKSPACE__THEME` unset; View. | `entry(logging.file).value` is `~/.local/state/aibox/logs/run.log`; sources under `HOME` start with `~`; the text of stdout and stderr does not contain the `HOME` path or the canary. |

## Result output

| ID | Covers | Steps | Expected |
|---|---|---|---|
| J-OUT-01 | OUT-001, OUT-002, LOG-007 | `inspect --format json --log-level trace` on P. | stdout is one JSON document and a newline; it validates as an operation result; no event text is on stdout. |
| J-OUT-02 | OUT-003 | Same with `--format yaml`, and the View with `--format yaml`. | A YAML 1.2 parser gives a value equal to the JSON form of the same call (ignoring `requestId` and timestamps). |
| J-OUT-03 | OUT-004, OUT-005 | `inspect` on P with no format; then on a missing project; each with stdout not a terminal, then on a pseudo-terminal, then on a pseudo-terminal with `NO_COLOR=1`, then `--color always` piped, then `--color always` with `NO_COLOR=1`. | Text contains `inspect_workspace` and `succeeded`; for the missing project it contains `failed` and `invalid_input`, exit 2. Escape byte 0x1B appears only on the pseudo-terminal without `NO_COLOR` and with `--color always` piped without `NO_COLOR`. |
| J-OUT-04 | OUT-006, OUT-008 | `inspect --format json` with a malformed user settings file; same with `--format yaml`; same with `AIBOX_OUTPUT__FORMAT=json` and no flag. | Exit 2; stdout is one error envelope in the requested format with `operation` `inspect_workspace`, `outcome` `failed`, error code `invalid_input`, `nextAction` `edit_configuration`. |
| J-OUT-05 | OUT-007 | `inspect` with a malformed user settings file and no format; `inspect --format xml`. | Exit 2; stdout empty; the last stderr line matches `output.bootstrapDiagnostic`. |

## Logging

| ID | Covers | Steps | Expected |
|---|---|---|---|
| J-LOG-01 | LOG-001 to LOG-003, LOG-010 | `inspect --format json --log-level trace --log-format jsonl` on P. | stderr holds exactly the events `settings.resolved`, `operation.started`, `operation.validated`, `operation.finished` in that order; each validates against the schema; one `requestId`, equal to the result's; `sequence` 0 to 3; `operation.finished` has `outcome` `succeeded`. |
| J-LOG-02 | LOG-005 | Same with each `--log-level` value. | `trace`, `debug`: four events. `info`: the two `info` events. `warn`, `error`, `off`: none. On a missing project the same counts hold, and `operation.finished` has `outcome` `failed` and `errorCategory` `invalid_input`. |
| J-LOG-03 | LOG-010 | `inspect --format json --log-format plain` on P. | Every stderr line matches `logging.plainLine`. |
| J-LOG-04 | LOG-004, LOG-011 | Start `mcp serve`; initialize; call `inspect_workspace` twice; close stdin. | Every stdout line is an MCP JSON-RPC message. stderr is JSON Lines: `process.started` first and `process.finished` last with one process ID as `requestId`; the two calls use two other request IDs, each with `sequence` from 0. |
| J-LOG-05 | LOG-012, LOG-013, LOG-015 | `inspect --format json --log-file run.log` on P with a fresh `HOME`. | `~/.local/state/aibox/logs/run.log` exists with mode 0600, its directory and `~/.local/state/aibox` created with mode 0700; its lines equal the stderr events of the same run in `jsonl` form. The working directory has no `run.log`. |
| J-LOG-06 | LOG-014 | `--log-file /tmp/x.log`; `--log-file ../x.log`; `--log-file link.log` where `logs/link.log` is a symlink to a file outside the root; user file with `logging.file` `/var/log/aibox.log`. | Exit 3, `not_authorized`, reason `path_escape`; the outside file is not created or changed. |
| J-LOG-07 | LOG-016, LOG-017 | User file: `logging.file` `r.log`, `rotationMiB` 1, `retentionFiles` 3. Pre-fill `r.log` with 1 MiB of lines; run `inspect` four times, re-filling `r.log` to 1 MiB before each. | After the runs only `r.log`, `r.log.1`, `r.log.2` exist (besides `r.log.lock`); `r.log.1` holds the newest rotated content. Then set the modification time of `r.log.2` to 30 days ago with `retentionDays` 7 and run once more: `r.log.2` is gone. |
| J-LOG-08 | LOG-018 | Start 16 `inspect --log-level trace --log-file c.log` processes at once. | `c.log` has 64 lines; every line is a complete JSON object that validates; each of the 16 request IDs has `sequence` 0 to 3. |
| J-LOG-09 | LOG-019 | Create `~/.local/state/aibox/operations/keep.json` and `logs/other.txt`; run J-LOG-07. | Both files are unchanged. |
| J-LOG-10 | LOG-030, LOG-031 | `logs/` exists with mode 0500 and `--log-file new.log`; then the same with a user file that also sets `logging.required` true. | Both: exit 0, `outcome` `succeeded`, `warnings` contains `log_sink_degraded`, stderr has one `logging.sink_degraded` event with `fields.sink` `file`. |
| J-LOG-11 | LOG-032 | `inspect --format json` on P with stderr closed (`2>&-`). | Exit 0; stdout is the complete result. |
| J-LOG-12 | RED-001, RED-002, RED-005 | `EXAMPLE_API_TOKEN=canary-0123456789` in the environment; project directory named `p-canary-0123456789`; a second project named `p-ghp_` followed by 36 letters; `inspect --log-level trace --log-format jsonl --log-file s.log` on each; then on a missing path that contains the canary. `aibox help` output. | Neither stderr nor `s.log` contains the canary or the `ghp_` string; `[redacted]` appears in their place. `help` lists no flag about redaction, raw logs or secrets. |

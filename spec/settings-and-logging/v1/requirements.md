# Layered settings and diagnostic sinks: requirements

Normative. Requirements use *EARS* patterns. "aibox" means the v1 Go
executable in the CLI and in `mcp serve`. Names in `code` refer to
`interfaces.json` unless a schema is named. IDs are never reused.

Verification levels: **J** = a journey in `journeys.md` (black box);
**C** = component test only, with the reason given.

## Settings layers and precedence

| ID | Requirement | Level |
|---|---|---|
| CFG-001 | aibox shall resolve every key in `settings.keys` from the layers in `settings.layers`, in that order, and the value from the last layer that sets the key shall win. | J |
| CFG-002 | aibox shall accept a key from a layer only if that layer is listed in the key's `layers`. | J |
| CFG-003 | If a layer that is not listed in a key's `layers` sets that key, then aibox shall fail with error `not_authorized` before it runs the operation, and shall not fall back to a lower layer. | J |
| CFG-004 | Where no layer sets a key, aibox shall use the key's `default` for the running mode (`cli` or `mcp`). | J |
| CFG-005 | aibox shall complete settings resolution without creating, changing or deleting any file. | J |
| CFG-006 | aibox shall finish parsing, validation and authority checks of all layers before it starts the requested operation. | J |
| CFG-007 | aibox shall resolve settings the same way for the CLI and for `mcp serve`, except for the `default` values that `interfaces.json` gives per mode. | J |

## Settings files

| ID | Requirement | Level |
|---|---|---|
| CFG-010 | aibox shall read the system layer from `settings.files.system` for the running operating system. | J |
| CFG-011 | aibox shall read the user layer from `settings.files.user`, where `XDG_CONFIG_HOME` is used only if it is set to an absolute path. | J |
| CFG-012 | If a default system or user settings file does not exist, then aibox shall treat that layer as absent and continue. | J |
| CFG-013 | If a settings file exists and cannot be read, then aibox shall fail with `invalid_input` and reason `unreadable`. | J |
| CFG-014 | If a settings file does not validate against `schemas/process-settings.schema.json`, then aibox shall fail with `invalid_input` and the reason that `settings.rejectionReasons` assigns to the defect. | J |
| CFG-015 | If a settings file contains the same key twice in one JSON object, then aibox shall fail with `invalid_input` and reason `duplicate_key`. | J |
| CFG-016 | If a settings file contains a `null` value, then aibox shall fail with `invalid_input` and reason `null_value`. | J |
| CFG-017 | When `--settings-file PATH` is given, aibox shall read the user layer from `PATH`, anchored to the working directory at process start, instead of the default user file, and shall still read the system layer. | J |
| CFG-018 | If the file named by `--settings-file` does not exist, then aibox shall fail with `invalid_input` and reason `missing_selected_file`. | J |
| CFG-019 | aibox shall not read a settings file from the project directory, from the current directory, or from `~/Library/Preferences`. | J |

## Project layer

| ID | Requirement | Level |
|---|---|---|
| CFG-020 | aibox shall read the project layer only from `customizations.aibox.workspace.theme`, `.mode` and `.layout` of `<project>/.devcontainer/devcontainer.json`, and shall apply it only to the keys `workspace.theme`, `workspace.mode` and `workspace.layout`. | J |
| CFG-021 | If the project layer gives one of those three keys a value outside the key's `values`, then aibox shall fail with `invalid_input` and reason `invalid_value`. | J |
| CFG-022 | aibox shall not let any content of the project directory change a key outside `workspace.*`. | J |

## Environment files and process environment

| ID | Requirement | Level |
|---|---|---|
| CFG-030 | aibox shall load an environment file only when it is named by `--env-file` or, where no `--env-file` is given, by `AIBOX_ENV_FILE`. | J |
| CFG-031 | When several `--env-file` options are given, aibox shall apply them from left to right, and a later file shall override an earlier one. | J |
| CFG-032 | If `AIBOX_ENV_FILE` is not an absolute path, then aibox shall fail with `invalid_input` and reason `relative_env_file`. | J |
| CFG-033 | aibox shall parse an environment file as described by `settings.envFile.grammar`. | J |
| CFG-034 | If an environment file contains a line that the grammar does not allow, then aibox shall fail with `invalid_input` and reason `malformed`. | J |
| CFG-035 | If an environment file sets the same key twice, then aibox shall fail with `invalid_input` and reason `duplicate_key`. | J |
| CFG-036 | If an environment file is missing, is not a regular file after symlink resolution, or cannot be read, then aibox shall fail with `invalid_input` and reason `missing_selected_file`, `not_regular_file` or `unreadable` respectively. | J |
| CFG-037 | If an environment file is readable by its group or by others, then aibox shall fail with `invalid_input` and reason `insecure_permissions`. | J |
| CFG-038 | If an environment file contains a key that is not the `env` name of a key in `settings.keys`, then aibox shall fail with `invalid_input` and reason `unknown_key`. | J |
| CFG-039 | aibox shall map a process environment variable to a setting only by the `env` name listed for the key. | J |
| CFG-040 | If the process environment contains a name that starts with `AIBOX_`, is not an `env` name, is not a selector in `settings.selectors`, and is not registered by `fixtures/registered-environment.json`, then aibox shall fail with `invalid_input` and reason `unknown_environment_name`, and the diagnostic shall name the nearest `env` name. | J |
| CFG-041 | aibox shall ignore every registered runtime identifier in `fixtures/registered-environment.json` when it resolves settings. | J |
| CFG-042 | If an environment value cannot be converted to the key's type by `settings.conversion`, then aibox shall fail with `invalid_input` and reason `invalid_value`. | J |
| CFG-043 | aibox shall not pass a key from an environment file to a child process or into a container. | C: no V1-06 operation starts a child process. |

## Flags and selectors

| ID | Requirement | Level |
|---|---|---|
| CFG-050 | aibox shall accept the flags in `cli.flags` on the commands listed for each flag, and a flag shall set its key for that invocation only. | J |
| CFG-051 | If a flag value is outside the key's `values` or range, then aibox shall fail with `invalid_input` and reason `invalid_value`. | J |
| CFG-052 | When `--project PATH` is relative, aibox shall anchor it to the working directory at process start. | J |
| CFG-053 | aibox shall not change its working directory while it resolves settings. | C: not observable from outside the process. |
| CFG-054 | If `mcp serve` is started with `logging.format` resolved to `plain`, then aibox shall fail with `invalid_input` and reason `invalid_value` and shall write nothing to stdout. | J |
| CFG-055 | While `mcp serve` runs, aibox shall keep the settings resolved at start and shall not let a tool call change any key except through the request field `effectiveConfig`. | J |

## Effective-configuration view

| ID | Requirement | Level |
|---|---|---|
| CFG-060 | When `inspect` is run with `--effective-config`, or `inspect_workspace` is called with `effectiveConfig` true, aibox shall return in `data.configuration` exactly one entry for each key in `settings.keys`, in the order of that list. | J |
| CFG-061 | Each such entry shall have `domain` as listed for the key, the winning `layer`, the `source` defined by `effectiveConfig.source`, the winning `value`, `redacted`, and `invocationOnly` true only when the winning layer is `invocation`. | J |
| CFG-062 | Where a lower layer also set the key, the entry shall list each such value in `overridden`, from highest to lowest layer, with its `layer`, `source`, `value` and `redacted`. | J |
| CFG-063 | aibox shall return `data.settingsSources` with one element for each source in `effectiveConfig.settingsSources`, giving its `layer`, `path` and `state` (`loaded` or `absent`). | J |
| CFG-064 | aibox shall redact values and paths in the view by `redaction`, and shall set `redacted` true on every entry whose value it changed. | J |
| CFG-065 | aibox shall return the same `data.configuration` and `data.settingsSources` for the CLI and for MCP when both resolve the same layers. | J |
| CFG-066 | While it builds the view, aibox shall not run a hook, start a container, contact a container engine, or read a home directory other than that of the running user. | J |
| CFG-067 | When `inspect` is run without `--effective-config`, aibox shall return `data.configuration` as before this feature and no `data.settingsSources`. | J |

## Result output

| ID | Requirement | Level |
|---|---|---|
| OUT-001 | aibox shall write the result of `inspect` to stdout in the format given by `output.format`, and nothing else to stdout. | J |
| OUT-002 | Where `output.format` is `json`, aibox shall write the result envelope as one JSON document followed by a newline. | J |
| OUT-003 | Where `output.format` is `yaml`, aibox shall write one YAML document that parses, with a YAML 1.2 parser, to the same value as the JSON form. | J |
| OUT-004 | Where `output.format` is `human`, aibox shall write a text summary that states the operation, the outcome and, on failure, the error code. | J |
| OUT-005 | aibox shall use ANSI colour in `human` output only as `output.colorRule` defines. | J |
| OUT-006 | If settings resolution fails and the result format is `json` or `yaml` by an already valid flag or environment value, then aibox shall write one error envelope with `outcome` `failed` in that format and exit with the code of the error. | J |
| OUT-007 | If settings resolution fails and the result format is not known to be `json` or `yaml`, then aibox shall write one plain diagnostic line of the form `output.bootstrapDiagnostic` to stderr, nothing to stdout, and exit with the code of the error. | J |
| OUT-008 | aibox shall exit with the code that `errors` assigns to the result's error code, and with 0 for an outcome without error. | J |

## Log events

| ID | Requirement | Level |
|---|---|---|
| LOG-001 | aibox shall emit, for every operation, the events that `logging.events` marks `emittedBy` for that operation, in the listed order. | J |
| LOG-002 | Every event shall validate against `schemas/log-event.schema.json`. | J |
| LOG-003 | Events of one request shall carry the same `requestId` and a `sequence` that starts at 0 and increases by 1. | J |
| LOG-004 | Process start and shutdown events shall carry a process correlation ID as `requestId` that differs from every operation request ID of that process. | J |
| LOG-005 | aibox shall emit an event only if its severity is at or above `logging.level`, and shall emit no event where `logging.level` is `off`. | J |
| LOG-006 | When a setting is rejected, aibox shall emit one `settings.rejected` event of severity `error` whose `fields` hold `key`, `layer`, `source` and a `reason` from `settings.rejectionReasons`, if the logging settings needed to emit it are already valid; otherwise OUT-007 applies. | J |
| LOG-007 | aibox shall write no event and no progress text to stdout. | J |

## Sinks

| ID | Requirement | Level |
|---|---|---|
| LOG-010 | aibox shall write events to stderr in the format given by `logging.format`: `plain` as `logging.plainLine`, `jsonl` as one JSON object per line. | J |
| LOG-011 | While it runs as `mcp serve`, aibox shall write events to stderr as `jsonl` and shall write only MCP protocol messages to stdout. | J |
| LOG-012 | Where `logging.file` is set, aibox shall also append every emitted event to that file as `jsonl`. | J |
| LOG-013 | aibox shall resolve a relative `logging.file` and a relative `--log-file` against `logging.fileSink.stateDirectory`. | J |
| LOG-014 | If the log file path, after symlink resolution, is outside `logging.fileSink.permittedRoot`, then aibox shall fail with `not_authorized` and reason `path_escape`. | J |
| LOG-015 | When it creates the log file or a parent directory, aibox shall create directories with mode 0700 and files with mode 0600. | J |
| LOG-016 | When the active log file is at or above `logging.rotationMiB` before an append, aibox shall rotate it as `logging.fileSink.rotation` defines. | J |
| LOG-017 | At process start with a file sink and at every rotation, aibox shall delete rotated files beyond `logging.retentionFiles` and rotated files older than `logging.retentionDays`. | J |
| LOG-018 | aibox shall append each event with one write, so that concurrent aibox processes produce only whole lines. | J |
| LOG-019 | aibox shall never delete or change a file outside the log file's own rotation set, and never a file under `logging.fileSink.operationRecordDirectory`. | J |

## Sink failure and backpressure

| ID | Requirement | Level |
|---|---|---|
| LOG-030 | If a write to the file sink fails during a read-only operation, then aibox shall complete the operation, add the warning `log_sink_degraded` to the result, and write one `logging.sink_degraded` event to stderr. | J |
| LOG-031 | If `logging.required` is true and the file sink cannot be opened, then aibox shall refuse a mutating operation with `operation_record_unavailable` before any side effect, and shall continue a read-only operation with the warning of LOG-030. | J for the read-only part; C for the mutating part: no mutating operation exists yet. |
| LOG-032 | If stderr is closed, then aibox shall still write the result to stdout, exit with the result's code, and not run the operation again. | J |
| LOG-033 | While the event queue holds `logging.queue.capacity` events, aibox shall drop `trace` and `debug` events first, count the dropped events, and emit one `logging.dropped` event of severity `warn` with the count when the queue has room. | C: needs a sink that blocks on demand. |
| LOG-034 | At normal shutdown aibox shall try to flush the sinks for at most `logging.queue.flushSeconds` seconds and shall never block an operation on a sink for longer. | C: needs a sink that blocks on demand. |

## Redaction

| ID | Requirement | Level |
|---|---|---|
| RED-001 | aibox shall replace every known secret value, as `redaction.knownSecrets` defines, with `redaction.replacement` before an event reaches a sink or a queue, at every log level. | J |
| RED-002 | aibox shall replace every string that matches a pattern in `redaction.secretShapes` in the same way. | J |
| RED-003 | aibox shall write the home directory prefix of a path as `~` in events and in the effective-configuration view. | J |
| RED-004 | aibox shall apply RED-001 and RED-002 to a value that is split over several writes of child output. | C: no V1-06 operation captures child output. |
| RED-005 | aibox shall offer no flag, setting or environment name that disables redaction. | J |

## Timeouts

| ID | Requirement | Level |
|---|---|---|
| EXE-001 | Where `execution.timeoutSeconds` is unset, aibox shall limit `inspect` to `execution.classDefaults.read` seconds. | C: `inspect` cannot be made slow from outside. |
| EXE-002 | If an operation exceeds its time limit, then aibox shall return outcome `timed_out` with error `operation_timeout` and exit 124. | C: same reason. |
| EXE-003 | The effective-configuration view shall report, for `execution.timeoutSeconds`, the limit that applies to the inspected operation. | J |

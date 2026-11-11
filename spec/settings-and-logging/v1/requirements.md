# Layered settings and diagnostic sinks: requirements

Normative. Requirements use *EARS* patterns. "aibox" means the v1 Go
executable in the CLI and in `mcp serve`. Names in `code` refer to
`interfaces.json` unless a schema is named. IDs are never reused.

Verification levels: **J** = a journey in `journeys.md` (black box);
**C** = component test only, with the reason given.

"Fails with reason R" means: the operation does not start; the error code is
`invalid_input`; and aibox emits the `settings.rejected` event with
`fields.reason` R as LOG-006 defines. Every failure of this feature before an
operation starts has exit code 2.

## Settings layers and precedence

| ID | Requirement | Level |
|---|---|---|
| CFG-001 | aibox shall resolve every key in `settings.keys` from the layers in `settings.layers`, in that order, and the value from the last layer that sets the key shall win. | J |
| CFG-002 | aibox shall accept a key from a layer only if that layer is listed in the key's `layers`. | J |
| CFG-003 | If a layer that is not listed in a key's `layers` sets that key, then aibox shall fail with reason `unauthorized_source`, and shall not fall back to a lower layer. | J |
| CFG-004 | Where no layer sets a key, aibox shall use the key's `default` for the running mode (`cli` or `mcp`). | J |
| CFG-005 | aibox shall complete settings resolution without creating, changing or deleting any file. | J |
| CFG-006 | aibox shall finish parsing, validation and authority checks of all layers before it starts the requested operation. | J |
| CFG-007 | aibox shall resolve settings the same way for the CLI and for `mcp serve`, except for the `default` values that `interfaces.json` gives per mode. | J |
| CFG-008 | Where several settings are defective, aibox shall report the first one in the order of `settings.failureOrder`. | J |
| CFG-009 | The commands `version` and `help` shall not resolve settings, and shall succeed whatever the settings files and the environment contain. | J |

## Settings files

| ID | Requirement | Level |
|---|---|---|
| CFG-010 | aibox shall read the system layer from `settings.files.system` for the running operating system. | J |
| CFG-011 | aibox shall read the user layer from `settings.files.user`, where `XDG_CONFIG_HOME` is used only if it is set to an absolute path. | J |
| CFG-012 | If a default system or user settings file does not exist, then aibox shall treat that layer as absent and continue. | J |
| CFG-013 | If a settings file exists and cannot be read, then aibox shall fail with reason `unreadable`. | J |
| CFG-014 | If a settings file does not validate against `schemas/process-settings.schema.json`, then aibox shall fail with the reason that `settings.schemaDefects` assigns to the defect. | J |
| CFG-015 | If a settings file contains the same key twice in one JSON object, then aibox shall fail with reason `duplicate_key`. | J |
| CFG-016 | If a settings file contains a `null` value, then aibox shall fail with reason `null_value`. | J |
| CFG-017 | When `--settings-file PATH` is given, aibox shall read the user layer from `PATH`, anchored to the working directory at process start, instead of the default user file, and shall still read the system layer. | J |
| CFG-018 | If the file named by `--settings-file` does not exist, then aibox shall fail with reason `missing_selected_file`. | J |
| CFG-019 | aibox shall not read a settings file from the project directory, from the current directory, or from `~/Library/Preferences`. | J |
| CFG-024 | If the home directory of the running user cannot be determined, then aibox shall fail with reason `no_home`. | J |

## Project layer and workspace keys

| ID | Requirement | Level |
|---|---|---|
| CFG-020 | aibox shall read the project layer only from the members `theme`, `mode` and `layout` of `customizations.aibox.workspace` in `<project>/.devcontainer/devcontainer.json`, shall apply it only to the keys `workspace.theme`, `workspace.mode` and `workspace.layout`, and shall ignore every other member of `customizations.aibox` for settings. | J |
| CFG-021 | If any layer gives `workspace.theme`, `workspace.mode` or `workspace.layout` a value outside the key's `values`, then aibox shall fail with reason `invalid_value`. | J |
| CFG-022 | aibox shall not let any content of the project directory change a key outside `workspace.*`. | J |

## Environment files and process environment

| ID | Requirement | Level |
|---|---|---|
| CFG-030 | aibox shall load an environment file only when it is named by `--env-file` or, where no `--env-file` is given, by `AIBOX_ENV_FILE`. | J |
| CFG-031 | When several `--env-file` options are given, aibox shall apply them from left to right, and a later file shall override an earlier one. | J |
| CFG-032 | If `AIBOX_ENV_FILE` is not an absolute path, then aibox shall fail with reason `relative_env_file`. | J |
| CFG-033 | aibox shall parse an environment file as described by `settings.envFile.grammar`. | J |
| CFG-034 | If an environment file contains a line that the grammar does not allow, then aibox shall fail with reason `malformed`. | J |
| CFG-035 | If an environment file sets the same key twice, then aibox shall fail with reason `duplicate_key`. | J |
| CFG-036 | If an environment file is missing, is not a regular file after symlink resolution, or cannot be read, then aibox shall fail with reason `missing_selected_file`, `not_regular_file` or `unreadable` respectively. | J |
| CFG-038 | If an environment file contains a key that the grammar allows and that is not the `env` name of a key in `settings.keys`, then aibox shall fail with reason `unknown_key`. | J |
| CFG-039 | aibox shall map a process environment variable to a setting only by the `env` name listed for the key. | J |
| CFG-040 | If the process environment contains a name that starts with `AIBOX_`, is not an `env` name, is not `AIBOX_ENV_FILE`, and is not registered as `settings.registeredEnvironment` defines, then aibox shall fail with reason `unknown_environment_name`, and the event shall carry the nearest `env` name in `fields.suggestion` as `settings.nearestName` defines. | J |
| CFG-041 | aibox shall ignore every registered runtime identifier when it resolves settings. | J |
| CFG-042 | If an environment value cannot be converted to the key's type by `settings.conversion`, then aibox shall fail with reason `invalid_value`. | J |
| CFG-043 | aibox shall not pass a key from an environment file to a child process or into a container. | C: no V1-06 operation starts a child process. |
| CFG-044 | If more than `settings.envFile.maxFiles` environment files are named, then aibox shall fail with reason `usage_error`. | J |

CFG-037 is retired: this version has no secret-bearing setting, so an
environment file needs no private mode.

## Flags and selectors

| ID | Requirement | Level |
|---|---|---|
| CFG-050 | aibox shall accept the flags in `cli.flags` on the commands listed for each flag, and a flag shall set its key for that invocation only. | J |
| CFG-051 | If a flag value is outside the key's `values` or range, then aibox shall fail with reason `invalid_value`. | J |
| CFG-052 | When `--project PATH` is relative, aibox shall anchor it to the working directory at process start. | J |
| CFG-053 | aibox shall not change its working directory while it resolves settings. | C: not observable from outside the process. |
| CFG-054 | If `mcp serve` is started with `logging.format` set to `plain` by any layer, then aibox shall fail with reason `invalid_value`. | J |
| CFG-055 | While `mcp serve` runs, aibox shall keep the settings resolved at start, and no tool input shall change a key. | J |
| CFG-056 | If a command line holds a flag that `cli.flags` and `settings.selectors` do not list for that command, a non-repeatable flag twice, or a positional argument, then aibox shall fail with reason `usage_error`. | J |
| CFG-057 | If `mcp serve` fails before it serves, then aibox shall write nothing to stdout, write the diagnostic of OUT-007 to stderr, and exit 2. | J |

## Effective-configuration view

| ID | Requirement | Level |
|---|---|---|
| CFG-060 | When `inspect` is run with `--effective-config`, or `inspect_workspace` is called with `effectiveConfig` true, aibox shall return in `data.configuration` first the native entries of `effectiveConfig.nativeEntries` and then exactly one entry for each key in `settings.keys`, in the order of that list. | J |
| CFG-061 | Each settings entry shall have `domain` as listed for the key, the winning `layer`, the `source` defined by `effectiveConfig.source`, the winning `value`, `redacted`, and `invocationOnly` true only when the winning layer is `invocation`. | J |
| CFG-062 | Where a lower layer other than `built_in` also set the key, the entry shall list each such value in `overridden`, from highest to lowest layer, with its `layer`, `source`, `value` and `redacted`. | J |
| CFG-063 | aibox shall return `data.settingsSources` with one element for each source in `effectiveConfig.settingsSources`, giving its `layer`, `path` and `state` (`loaded` or `absent`). | J |
| CFG-064 | aibox shall redact values, sources and paths in the view by `redaction`, and shall set `redacted` true on every entry whose value holds `redaction.replacement`; writing the home prefix as `~` does not set it. | J |
| CFG-065 | aibox shall return the same `data.settingsSources` and the same settings entries for the CLI and for MCP when both resolve the same layers. | J |
| CFG-066 | While it builds the view, aibox shall not run a hook, start a container, contact a container engine, or read a home directory other than that of the running user. | J |
| CFG-067 | When `inspect` is run without `--effective-config`, aibox shall return in `data.configuration` only the native entries of `effectiveConfig.nativeEntries`, and no `data.settingsSources`. | J |
| CFG-068 | The result of `inspect` shall validate against `schemas/operation-result.schema.json` of this version. | J |

## Result output of `inspect`

| ID | Requirement | Level |
|---|---|---|
| OUT-001 | `inspect` shall write its result to stdout in the format given by `output.format`, and nothing else to stdout. | J |
| OUT-002 | Where `output.format` is `json`, `inspect` shall write the result envelope as one JSON document followed by a newline. | J |
| OUT-003 | Where `output.format` is `yaml`, `inspect` shall write one YAML document that parses, with a YAML 1.2 parser, to the same value as the JSON form. | J |
| OUT-004 | Where `output.format` is `human`, `inspect` shall write a text summary that contains the operation name, the outcome and, on failure, the error code. | J |
| OUT-005 | `inspect` shall use ANSI colour in `human` output only as `output.colorRule` defines. | J |
| OUT-006 | If `inspect` fails before the operation starts and `--format json` or `--format yaml` is on the command line with that valid value, then aibox shall write one error envelope in that format with `operation` `invalid_request`, `outcome` `failed` and error code `invalid_input`, and exit 2. | J |
| OUT-007 | If a command fails before the operation starts and OUT-006 does not apply, then aibox shall write one line of the form `output.bootstrapDiagnostic` as the last line of stderr, nothing to stdout, and exit 2. | J |
| OUT-008 | `inspect` shall exit with the code that `errors` assigns to the result's error code, and with 0 for an outcome without error. | J |

## Log events

| ID | Requirement | Level |
|---|---|---|
| LOG-001 | aibox shall emit, for every command and tool call, the events that `logging.events` lists for it, in the listed order, each with the members and fields listed there. | J |
| LOG-002 | Every event shall validate against `schemas/log-event.schema.json`. | J |
| LOG-003 | Events of one request shall carry the same `requestId` and a `sequence` that starts at 0 and increases by 1. | J |
| LOG-004 | Process start and shutdown events of `mcp serve` shall carry a process correlation ID as `requestId` that differs from every operation request ID of that process. | J |
| LOG-005 | aibox shall emit an event only if its severity is at or above `logging.level`, and shall emit no event where `logging.level` is `off`. | J |
| LOG-006 | When a setting or a command line is rejected, aibox shall emit one `settings.rejected` event with the request ID of the failed command and `sequence` 0, written to stderr only, at the level and in the format that `logging.rejectionFormat` defines. `version` and `help` emit no event; OUT-007 alone applies to them. | J |
| LOG-007 | aibox shall write no event and no progress text to stdout. | J |
| LOG-008 | When `inspect` is interrupted by SIGINT or SIGTERM, aibox shall emit `operation.cancelled`, return outcome `cancelled` with error `operation_cancelled`, and exit 130. | C: `inspect` ends too quickly to interrupt from outside. |

## Sinks

| ID | Requirement | Level |
|---|---|---|
| LOG-010 | aibox shall write events to stderr in the format given by `logging.format`: `plain` as `logging.plainLine`, `jsonl` as one JSON object per line. | J |
| LOG-011 | While it runs as `mcp serve`, aibox shall write events to stderr as `jsonl` and shall write only MCP protocol messages to stdout. | J |
| LOG-012 | Where `logging.file` is set, aibox shall also append every emitted event to that file as `jsonl`. | J |
| LOG-013 | aibox shall resolve a relative `logging.file` and a relative `--log-file` against `logging.fileSink.stateDirectory`. | J |
| LOG-014 | If the log file path, after symlink resolution, is outside `logging.fileSink.permittedRoot`, then aibox shall fail with reason `path_escape`. | J |
| LOG-015 | When it creates the log file or a parent directory, aibox shall create directories with mode 0700 and files with mode 0600. | J |
| LOG-016 | aibox shall append and rotate as `logging.fileSink.append` and `logging.fileSink.rotation` define. | J |
| LOG-017 | At the first append of a process and at every rotation, aibox shall delete rotated files beyond `logging.retentionFiles` and rotated files older than `logging.retentionDays`. | J |
| LOG-018 | Concurrent aibox processes that append to one log file on a local file system shall produce only whole lines and shall lose no event to rotation. | J |
| LOG-019 | aibox shall never delete or change a file outside the log file's own rotation set, and never a file under `logging.fileSink.operationRecordDirectory`. | J |
| LOG-020 | If the log file exists with any group or other permission bit, then aibox shall treat the file sink as failed. | J |

## Sink failure and backpressure

| ID | Requirement | Level |
|---|---|---|
| LOG-030 | If the file sink cannot be opened or written during a read-only operation, then aibox shall complete the operation, add the warning `log_sink_degraded` to the result, and write one `logging.sink_degraded` event to stderr. | J for open failure; C for write failure (disk full needs a faulty file system). |
| LOG-031 | If `logging.required` is true and the file sink has failed, then aibox shall refuse a mutating operation with `operation_record_unavailable` before any side effect, and shall continue a read-only operation as in LOG-030. | J for the read-only part; C for the mutating part: no mutating operation exists yet. |
| LOG-032 | If stderr is closed or its reader has gone, then aibox shall still write the result to stdout, exit with the result's code, not run the operation again, add no warning, emit no event about stderr, and write no event text to any other file. | J |
| LOG-033 | While the event queue holds `logging.queue.capacity` events, aibox shall drop `trace` and `debug` events first, write `warn` and `error` events directly to stderr within `logging.queue.directWriteMillis`, count the dropped events, and emit one `logging.dropped` event with the count when the queue has room. | C: needs a sink that blocks on demand. |
| LOG-034 | At normal shutdown aibox shall try to flush the sinks for at most `logging.queue.flushSeconds` seconds, shall never block an operation on a sink for longer, and, if the flush did not complete, shall write one `logging.flush_incomplete` event directly to stderr. | C: needs a sink that blocks on demand. |

## Redaction

| ID | Requirement | Level |
|---|---|---|
| RED-001 | aibox shall replace every known secret value, as `redaction.knownSecrets` defines, with `redaction.replacement` in everything listed in `redaction.scope`, before an event reaches a sink or a queue, at every log level. | J |
| RED-002 | aibox shall replace every string that matches a pattern in `redaction.secretShapes` in the same way. | J |
| RED-003 | aibox shall write the home directory prefix of a path as `~` wherever `redaction.homePrefixScope` lists. | J |
| RED-004 | aibox shall apply RED-001 and RED-002 to a value that is split over several writes of child output. | C: no V1-06 operation captures child output. |
| RED-005 | aibox shall accept no flag, setting or environment name that disables redaction. | J |

## Timeouts

| ID | Requirement | Level |
|---|---|---|
| EXE-001 | Where `execution.timeoutSeconds` is unset, aibox shall limit `inspect` to `execution.classDefaults.read` seconds. | C: `inspect` cannot be made slow from outside. |
| EXE-002 | If an operation exceeds its time limit, then aibox shall return outcome `timed_out` with error `operation_timeout` and exit 124. | C: same reason. |
| EXE-003 | The effective-configuration view shall report, for `execution.timeoutSeconds`, the limit that applies to the inspected operation. | J |

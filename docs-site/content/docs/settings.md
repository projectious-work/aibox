---
title: Settings
description: "Where aibox v1 reads its process settings, which source wins, and how to see the effective configuration."
weight: 40
---

# Settings

aibox settings change how the aibox process behaves: the result format,
colour, diagnostics, time limits and three workspace preferences. They never
change `devcontainer.json`, Features, images, mounts or hooks; those belong to
the project's native files.

This page describes the preview executable (`inspect`, `mcp serve`, `version`,
`help`). The contract is
[`spec/settings-and-logging/v1/`](https://github.com/projectious-work/aibox/tree/main/spec/settings-and-logging/v1).

## Layers and precedence

aibox reads seven layers. A later layer wins over an earlier one:

| Layer | Source | Set by |
|---|---|---|
| `built_in` | The default of the key for the running mode | aibox |
| `system` | `/etc/aibox/settings.json` (macOS: `/Library/Application Support/aibox/settings.json`) | An administrator |
| `user` | `$XDG_CONFIG_HOME/aibox/settings.json`, or `~/.config/aibox/settings.json`; or the file named by `--settings-file` | The user |
| `project` | `customizations.aibox.workspace` in the project's `.devcontainer/devcontainer.json` | The project |
| `env_file` | Files named by `--env-file`, or by `AIBOX_ENV_FILE` | One shell or script |
| `environment` | `AIBOX_*` variables of the process | One shell |
| `invocation` | Command-line flags | One command |

A key accepts a value only from the layers listed for it below. A value from
any other layer is an error (`unauthorized_source`), not a silent fallback:
for example, an environment variable cannot move the log file, and a project
cannot change logging at all.

`XDG_CONFIG_HOME` and `XDG_STATE_HOME` are used only when set to an absolute
path. The home directory comes from `HOME` only; if `HOME` is unset or not
absolute, aibox stops with `no_home`. aibox never reads a settings file from
the project directory, from the current directory, or from
`~/Library/Preferences`.

## Keys

| Key | Values | Default (CLI / MCP) | Layers | Environment name | Flag |
|---|---|---|---|---|---|
| `output.format` | `human`, `json`, `yaml` | `human` / `json` | system, user, env_file, environment, invocation | `AIBOX_OUTPUT__FORMAT` | `--format` |
| `output.color` | `auto`, `always`, `never` | `auto` | system, user, env_file, environment, invocation | `AIBOX_OUTPUT__COLOR` | `--color` |
| `logging.level` | `trace`, `debug`, `info`, `warn`, `error`, `off` | `info` | system, user, env_file, environment, invocation | `AIBOX_LOGGING__LEVEL` | `--log-level` |
| `logging.format` | `plain`, `jsonl` | `plain` / `jsonl` | system, user, env_file, environment, invocation | `AIBOX_LOGGING__FORMAT` | `--log-format` |
| `logging.file` | a path | unset | system, user, invocation | `AIBOX_LOGGING__FILE` | `--log-file` |
| `logging.required` | `true`, `false` | `false` | system, user | `AIBOX_LOGGING__REQUIRED` | none |
| `logging.rotationMiB` | 1 to 1024 | 10 | system, user | `AIBOX_LOGGING__ROTATION_MIB` | none |
| `logging.retentionFiles` | 1 to 100 | 7 | system, user | `AIBOX_LOGGING__RETENTION_FILES` | none |
| `logging.retentionDays` | 1 to 365 | 7 | system, user | `AIBOX_LOGGING__RETENTION_DAYS` | none |
| `execution.timeoutSeconds` | 1 to 86400 | unset: 30 for a read operation | system, user, env_file, environment, invocation | `AIBOX_EXECUTION__TIMEOUT_SECONDS` | `--timeout` |
| `execution.lockWaitSeconds` | 0 to 600 | 30 | system, user, env_file, environment | `AIBOX_EXECUTION__LOCK_WAIT_SECONDS` | none |
| `diagnostics.offline` | `true`, `false` | `true` | system, user, env_file, environment | `AIBOX_DIAGNOSTICS__OFFLINE` | none |
| `workspace.theme` | a theme name, for example `nord` | unset | system, user, project, env_file, environment | `AIBOX_WORKSPACE__THEME` | none |
| `workspace.mode` | `dark`, `light`, `auto` | unset | system, user, project, env_file, environment | `AIBOX_WORKSPACE__MODE` | none |
| `workspace.layout` | `ai`, `dev`, `focus`, `cowork` | unset | system, user, project, env_file, environment | `AIBOX_WORKSPACE__LAYOUT` | none |

`mcp serve` writes its events as JSON Lines only, so `logging.format` `plain`
from any layer stops it.

## Settings files

The system and user files share one format, checked against
[`process-settings.schema.json`](https://github.com/projectious-work/aibox/blob/main/spec/settings-and-logging/v1/schemas/process-settings.schema.json):

```json
{
  "schemaVersion": "1",
  "output": {"format": "json", "color": "never"},
  "logging": {"level": "debug", "file": "aibox.log", "rotationMiB": 5},
  "execution": {"timeoutSeconds": 60},
  "workspace": {"theme": "nord"}
}
```

`schemaVersion` must be the string `"1"`. An unknown member, a repeated
member, a `null`, or a value of the wrong type or range stops the command.
A missing default file is simply absent. `--settings-file PATH` replaces the
user file for one command; a relative path is taken from the working
directory, and a missing file is an error.

## Project layer

From the project, aibox reads only three members of the selected
`devcontainer.json`:

```jsonc
"customizations": {
  "aibox": {
    "workspace": {"theme": "nord", "mode": "dark", "layout": "dev"}
  }
}
```

They can set only `workspace.theme`, `workspace.mode` and `workspace.layout`.
Every other member of `customizations.aibox` is ignored for settings.

## Environment files

An environment file is loaded only when it is named: by one or more
`--env-file PATH` options (applied left to right, at most 16), or, when no
`--env-file` is given, by `AIBOX_ENV_FILE`, which must be an absolute path.
A `.env` file is never loaded implicitly.

```sh
# comment lines and blank lines are ignored
AIBOX_LOGGING__LEVEL=debug
AIBOX_WORKSPACE__THEME=nord
```

Each other line is `KEY=VALUE`, split at the first `=`. `KEY` matches
`^[A-Z][A-Z0-9_]*$` with no spaces and must be one of the environment names
above. The value is every byte after `=` up to the line end: quotes, spaces
and `#` are literal. There is no `export`, interpolation, continuation or
multi-line value. Lines may end with LF or CRLF. No private file mode is
required, because no setting of this version is a secret.

Keys from an environment file configure aibox only. They are never passed to
a child process or into a container.

## Environment variables

aibox maps a variable to a setting only by the names listed above. Any other
`AIBOX_*` name stops the command with `unknown_environment_name`, and the
event names the nearest valid name, so that a typo such as
`AIBOX_LOGING__LEVEL` is caught. The exception is the v0 runtime identifiers
that aibox containers set, such as `AIBOX_VERSION`, `AIBOX_LOG_LEVEL` and the
`AIBOX_TMUX_*` family; they are ignored. The list is
[`registered-environment.json`](https://github.com/projectious-work/aibox/blob/main/spec/settings-and-logging/v1/fixtures/registered-environment.json).

Values convert strictly: booleans are exactly `true` or `false`; integers are
base 10 without `+` or spaces; names are case-sensitive; an empty value is
always invalid.

## Flags

```text
aibox inspect [--context local] [--project PATH] [--effective-config]
              [--format human|json|yaml] [--color auto|always|never]
              [--log-level LEVEL] [--log-format plain|jsonl] [--log-file PATH]
              [--timeout SECONDS] [--settings-file PATH] [--env-file PATH]...
aibox mcp serve [--context local] [--log-level LEVEL] [--log-format jsonl]
              [--log-file PATH] [--settings-file PATH] [--env-file PATH]...
```

Both `--name value` and `--name=value` work. An unknown flag, a flag on a
command that does not list it, a repeated flag (other than `--env-file`), or
a positional argument is a usage error. A relative `--project`,
`--settings-file` or `--env-file` is taken from the working directory; a
relative `--log-file` from the log directory (see [logging](logging.md)).
`version` and `help` read no settings and cannot fail on them.

`mcp serve` resolves its settings once, at start. A tool call cannot change
a setting; `inspect_workspace` accepts only `requestId`, `projectRoot` and
`effectiveConfig`.

## Which value won: the effective configuration

```sh
aibox inspect --project . --effective-config --format json
```

`data.configuration` then lists the native entries followed by one entry per
key, in the order of the table above. Each entry has the winning `value`, its
`layer` and `source` (a file path, an environment name, a flag or
`built-in`), `invocationOnly`, `redacted`, and, when lower layers also set the
key, `overridden` with those values from highest to lowest layer.
`data.settingsSources` lists the system and user files, the project
declaration and each environment file, as `loaded` or `absent`:

```json
{"key": "logging.level", "value": "debug", "layer": "invocation", "source": "--log-level",
 "domain": "process", "invocationOnly": true, "redacted": false,
 "overridden": [{"layer": "user", "source": "~/.config/aibox/settings.json", "value": "warn", "redacted": false}]}
```

Paths under your home directory are shown with `~`, and secrets are replaced
by `[redacted]`. `execution.timeoutSeconds` shows the limit that applies to
`inspect` (30 seconds unless set). Without `--effective-config`, the result
holds only the native entries, as before. Through MCP, call
`inspect_workspace` with `"effectiveConfig": true`.

## Result formats and exit codes

`inspect` writes its result to stdout and nothing else: `human` (the default
on the CLI) is a short summary whose wording may change; `json` is one JSON
document and a newline; `yaml` is the same value as block YAML. Scripts should
pass `--format json`. Colour appears only in `human` output, only when
`NO_COLOR` is unset or empty, and only with `--color always` or with `auto`
on a terminal.

| Exit | Meaning |
|---|---|
| 0 | Succeeded |
| 2 | Invalid input, including every settings or command-line error |
| 5 | Internal error |
| 124 | The operation exceeded its time limit (`operation_timeout`) |
| 130 | Interrupted by SIGINT or SIGTERM (`operation_cancelled`) |

When a setting or the command line is rejected, the operation does not start.
With `--format json` or `--format yaml` on the command line, stdout holds an
error envelope with operation `invalid_request`; otherwise stdout is empty and
the last line of stderr reads `aibox: invalid_input: <reason>: <text>`. See
[troubleshooting](troubleshooting.md) for each reason.

---
title: Troubleshooting
description: "Rejection reasons of aibox v1 settings, reading diagnostics, and capturing the stderr of mcp serve under systemd and launchd."
weight: 60
---

# Troubleshooting

## Read the rejection

When aibox rejects a setting or a command line, the operation does not start
and the exit code is 2. Without `--format json`, the last stderr line names
the reason:

```text
aibox: invalid_input: unknown_environment_name: AIBOX_LOGING__LEVEL is not a known environment name; did you mean AIBOX_LOGGING__LEVEL?
```

For details, add `--log-format jsonl` and read the `settings.rejected` event:
`fields.reason`, and where known `fields.key`, `fields.layer`,
`fields.source` and `fields.suggestion`. With `--format json` the reason is in
that event and stdout holds an `invalid_request` envelope.

| Reason | Cause | Fix |
|---|---|---|
| `usage_error` | Unknown flag, flag on the wrong command, repeated flag, positional argument, more than 16 `--env-file`, or `--context` other than `local` | Check `aibox help`. |
| `no_home` | `HOME` is unset or not an absolute path | Set `HOME`. aibox does not consult the user database. |
| `malformed` | A settings file is not a JSON object; an environment file line is not `KEY=VALUE` | Fix the file. In an environment file, remove `export` and leading spaces. |
| `unreadable` | A settings or environment file exists but cannot be read | Check its permissions. |
| `missing_selected_file` | The file named by `--settings-file`, `--env-file` or `AIBOX_ENV_FILE` does not exist | Fix the path; relative paths start at the working directory. |
| `not_regular_file` | An environment file is a directory or device | Name a regular file. |
| `duplicate_key` | A member twice in one JSON object, or a key twice in one environment file | Keep one. |
| `unknown_key` | A member that is not a setting; an environment file key that is not an `AIBOX_` setting name | See the [keys](settings.md#keys). |
| `unsupported_version` | `schemaVersion` missing or not the string `"1"` | Add `"schemaVersion": "1"`. |
| `null_value` | A `null` in a settings file | Remove the member to use the default. |
| `invalid_value` | A value of the wrong type, outside its values or range, or empty; `plain` events for `mcp serve` | Use a listed value. Quotes in environment files are literal. |
| `unauthorized_source` | A key set from a layer it does not accept, for example `AIBOX_LOGGING__FILE` | Set it in a settings file or on the command line. |
| `unknown_environment_name` | An `AIBOX_*` variable that is not a setting or a registered runtime name | Fix the typo the suggestion names, or unset it. |
| `relative_env_file` | `AIBOX_ENV_FILE` is not absolute | Use an absolute path, or `--env-file`. |
| `path_escape` | The log file would be outside `~/.local/state/aibox/logs/` | Use a name inside that directory. |

When several things are wrong, aibox reports the first in this order: the
command line, `HOME`, the system file, the user file, the project, the
environment files in load order, the environment (by name), the flags. Fix
it and run again.

## See where a value comes from

```sh
aibox inspect --project . --effective-config --format json | jq '.data.configuration[] | {key, value, layer, source}'
aibox inspect --project . --effective-config --format json | jq '.data.settingsSources'
```

A file listed as `absent` was looked for and not found. If a value you set is
not used, check `overridden` on its entry: a higher layer won.

## Diagnostics in scripts

Use `--format json` for the result and `--log-format jsonl` for the events.
Since V1-06 the default result format of `inspect` is `human`, and `info`
events go to stderr; a script that needs an empty stderr passes
`--log-level off`.

## Capturing the stderr of `mcp serve`

`aibox mcp serve` speaks MCP on stdin and stdout, so an MCP client starts it
as a child process. Its events go to stderr as JSON Lines. To keep them:

- **Simplest:** start it with `--log-file mcp.log`. The events are then also
  in `~/.local/state/aibox/logs/mcp.log`, rotated and retained as described
  in [logging](logging.md), whatever the client does with stderr.
- **Under a service manager:** if the client runs as a service, the child
  inherits the service's stderr. Configure that stream on the client's
  service; do not redirect aibox's stdout, which carries the protocol.

### systemd

A user service for the client sends stderr to the journal by default;
`SyslogIdentifier` makes the lines easy to find:

```ini
# ~/.config/systemd/user/agent.service
[Service]
ExecStart=/usr/local/bin/my-agent
StandardError=journal
SyslogIdentifier=agent
```

```sh
systemctl --user daemon-reload && systemctl --user restart agent
journalctl --user -t agent -o cat | grep '"schemaVersion":"aibox.log-event/v1"' | jq .
```

To give aibox its own identifier while the client keeps its stdio, configure
the client to launch aibox through a shell that pipes only stderr into
`systemd-cat`:

```sh
bash -c 'exec aibox mcp serve --context local 2> >(exec systemd-cat --identifier=aibox-mcp)'
```

Then read it with `journalctl --user -t aibox-mcp -o cat`.

### launchd (macOS)

In the client's launch agent, set `StandardErrorPath` to an absolute file;
aibox, as a child, writes its events there:

```xml
<!-- ~/Library/LaunchAgents/com.example.agent.plist -->
<key>ProgramArguments</key>
<array><string>/usr/local/bin/my-agent</string></array>
<key>StandardErrorPath</key>
<string>/Users/me/Library/Logs/agent.stderr.jsonl</string>
```

```sh
launchctl bootout gui/$(id -u) ~/Library/LaunchAgents/com.example.agent.plist 2>/dev/null
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.example.agent.plist
tail -f ~/Library/Logs/agent.stderr.jsonl
```

launchd does not expand `~` in the plist, and it does not rotate that file;
`--log-file` gives you rotation. macOS is specified but not yet qualified on
a real host (roadmap phase V1-22).

## Other symptoms

| Symptom | Explanation |
|---|---|
| No colour | Colour needs `human` output, an unset `NO_COLOR`, and a terminal or `--color always`. |
| Result has the warning `log_sink_degraded` | The log file failed; the `logging.sink_degraded` event on stderr names the reason. An `insecure_mode` file needs `chmod 600`. |
| Exit 124 | The operation exceeded `execution.timeoutSeconds` (default 30 seconds for `inspect`); raise it with `--timeout`. |
| Exit 130 | SIGINT or SIGTERM interrupted the operation; it can be retried. |
| Events but no log file lines | Events below `logging.level` are written nowhere; the file gets the same events as stderr. |

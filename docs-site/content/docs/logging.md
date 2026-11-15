---
title: Logging and redaction
description: "The diagnostic events of aibox v1, the stderr and file sinks, rotation, failure handling and redaction."
weight: 50
---

# Logging and redaction

aibox writes diagnostics as semantic events to **stderr** and, if you ask for
it, to a rotating **log file**. Results go to stdout; events never do. Every
event is redacted before it reaches a sink.

## Events

| Event | Level | When |
|---|---|---|
| `settings.resolved` | debug | Settings were resolved; `fields.layers` lists the layers that set a key |
| `settings.rejected` | error | A setting or the command line was rejected; `fields.reason` says why |
| `operation.started` | info | The operation starts; `fields.projectRoot` is the project path |
| `operation.validated` | debug | The project declaration was read and accepted |
| `operation.cancelled` | warn | SIGINT or SIGTERM interrupted the operation |
| `operation.finished` | info | The operation ended; `outcome`, and `errorCategory` on failure |
| `process.started`, `process.finished` | info | `mcp serve` starts and stops |
| `logging.sink_degraded` | warn | The log file could not be opened or written (once per process) |
| `logging.dropped` | warn | Events were dropped under backpressure; `fields.count` |
| `logging.flush_incomplete` | warn | The shutdown flush did not finish within 2 seconds |

A successful `inspect` at level `trace` emits exactly `settings.resolved`,
`operation.started`, `operation.validated` and `operation.finished`. All
events of one command carry its `requestId` (the same as the result's) and a
`sequence` from 0. The events of `mcp serve` itself carry a separate process
ID; each tool call carries its own request ID.

`logging.level` (default `info`) drops events below it; `off` emits none.
Every event validates against
[`log-event.schema.json`](https://github.com/projectious-work/aibox/blob/main/spec/settings-and-logging/v1/schemas/log-event.schema.json);
all `fields` values are strings.

## Formats

`--log-format plain` (the CLI default) writes one line per event:

```text
2026-11-15T10:00:00.123Z INFO operation.started component=usecase operation=inspect_workspace requestId=4f… sequence=1 initiator=cli executor=aibox projectRoot=~/src/app
```

A value with spaces is JSON-quoted. `--log-format jsonl` (and always under
`mcp serve`) writes one JSON object per line; use it when a program reads
the events. The log file is always JSON Lines.

## The log file

`--log-file PATH` or `logging.file` in a settings file adds the file sink.
It cannot be set from the environment or an environment file, and it must
stay inside the log directory:

```text
${XDG_STATE_HOME:-$HOME/.local/state}/aibox/logs/
```

A relative path is taken from that directory, so `--log-file run.log` writes
`~/.local/state/aibox/logs/run.log`. A path that leads outside it, also
through `..` or a symlink, is rejected with `path_escape`. Operation records
live beside it in `aibox/operations/` and are never touched by logging.

- aibox creates missing directories with mode 0700 and the file with mode
  0600. If the file exists with any group or other permission bit, aibox
  does not use it (`insecure_mode`) and leaves it unchanged.
- Each event is appended with one write under an exclusive lock on
  `<file>.lock`, so concurrent aibox processes produce whole lines and lose
  nothing to rotation. Network file systems are not supported.
- When the file reaches `logging.rotationMiB` (default 10 MiB), it is renamed
  to `<file>.1`, older files move up (`<file>.2`, …), and at most
  `logging.retentionFiles` files exist including the active one (default 7).
  With `retentionFiles` 1 the file is truncated instead.
- At the first write of a process and at every rotation, rotated files older
  than `logging.retentionDays` (default 7) are deleted. aibox deletes only
  `<file>.<n>`; other files in the directory are left alone.

## When a sink fails

- **Log file cannot be opened or written:** a read-only operation still
  completes. Its result carries the warning `log_sink_degraded`, and one
  `logging.sink_degraded` event goes to stderr with `fields.reason`
  `open_failed`, `write_failed` or `insecure_mode`. With
  `logging.required: true`, a future mutating operation is refused instead
  (`operation_record_unavailable`, exit 5) before it changes anything.
- **stderr closed or its reader gone:** the result is still written to
  stdout with its normal exit code. aibox adds no warning and writes the
  events nowhere else; the log file, if configured, still gets them.
- **A slow sink:** events wait in a queue of 1024. When it is full, `trace`
  and `debug` events are dropped first and counted; `warn` and `error` events
  go directly to stderr within 100 ms; a `logging.dropped` event reports the
  count. At exit, aibox waits at most 2 seconds for the sinks.

## Redaction

Redaction is always on; no flag, setting or variable disables it.

- **Known secrets:** the value of every environment variable whose name
  contains `token`, `secret`, `password`, `passwd`, `credential`, `api_key`,
  `apikey`, `private_key` or `privatekey` (any case) and that is at least 8
  bytes long is replaced by `[redacted]`.
- **Secret shapes:** GitHub tokens (`ghp_…`, `github_pat_…`), `sk-…` keys,
  AWS access key IDs (`AKIA…`), private-key headers and `Bearer` tokens are
  replaced wherever they appear.
- **Home directory:** in event fields, the effective-configuration view,
  `settingsSources` and bootstrap diagnostics, your home directory prefix is
  written as `~`.

This covers every event member and field, the view's values and sources,
error messages, warnings, bootstrap diagnostics and human output.
`target.workspaceRoot` in the JSON and YAML result stays the real absolute
path, because it identifies the workspace.

Limits: a secret shorter than 8 bytes, or one that is in no such variable and
has no listed shape, is not recognized.

# 16. Layered application configuration and operational logging

This chapter defines configuration precedence and diagnostic output for the
Go CLI and MCP server. It does **not** replace `.devcontainer/devcontainer.json` or grant a
generic ability to override native Dev Container fields. `devcontainer.json`,
Features, Dockerfile and Compose remain the only project workspace/build
definition. The settings below configure the aibox process, its bounded
operations and aibox-owned presentation/diagnostic behavior. Operator policy,
credentials and runtime state are separate stores with separate authority.

## Configuration domains and file hierarchy

**R-CONFIG-SOURCES:** resolve each supported aibox setting in the documented
order, with deterministic provenance:

```text
built-in defaults
  < system settings.json
  < user settings.json
  < project customizations.aibox (eligible UX keys only)
  < explicitly loaded environment file(s)
  < inherited AIBOX_* process environment
  < explicit command-line flag / MCP invocation field
```

The last present, *authorized* value wins; an unauthorized value is an error,
not a value skipped so that a lower-precedence value appears to succeed.
The MCP invocation layer exists only for fields in the closed operation-request
schema; server output/logging settings are fixed when `mcp serve` starts and
cannot be changed by a tool caller.
Authority classes bound the layers:

| Key class | File and env/CLI behavior |
|---|---|
| Output, diagnostic level/format, timeout, offline preference | System/user file, env file, process env and documented CLI flags allowed within policy bounds; MCP request may override only its declared `offline`/operation fields, never server output or logging |
| `workspace.theme`, `workspace.mode`, `workspace.layout` | In the container-local process: system/user file, project customization, env file, process env and local `refresh` invocation override allowed; validate against the UX schema. Operator mode rejects these keys as wrong-context input; a host-side value does not silently enter a container. |
| Other workspace UX values | Project customization or native user file only; edit the owning file, no generic env/CLI key |
| Log-file destination, required-sink flag and retention | System/user file or explicit authorized invocation flag within the context's permitted state root only; env-file and ambient env attempts are rejected, even though a name can be derived mechanically |
| Operator policy, credentials, executable allowlist, permitted roots/contexts | Administrator-owned policy/secret store only; never repository, env file, environment or request override |
| Native image/build/Feature/mount/hook fields | Standard Dev Container files and upstream CLI semantics only; no aibox process-setting override |

The file paths are:

| Layer | Linux path | macOS path | Scope |
|---|---|---|---|
| System | `/etc/aibox/settings.json` | `/Library/Application Support/aibox/settings.json` | Administrator-owned, read-only to ordinary users |
| User | `${XDG_CONFIG_HOME:-$HOME/.config}/aibox/settings.json` | `${XDG_CONFIG_HOME:-$HOME/.config}/aibox/settings.json` | OS user running the process |
| Project | selected `.devcontainer/devcontainer.json` → `customizations.aibox` | same | Repository-owned, UX/diagnostic intent only |
| Local native overrides | Tool-native user files in the persistent container home, or explicit project-owned native files | same | Container-local user ownership; see chapter 8; no default `.aibox/` tree |
| Operator policy | explicit administrator-selected path outside project root | same | Authorization, *not* a precedence layer |
| Evidence/state | `${XDG_STATE_HOME:-$HOME/.local/state}/aibox/` | same | Receipts and logs, never policy/config authority |

Defaults are explicit: output format `human`, color `auto` (plain on non-TTY
and always plain under `NO_COLOR`), diagnostic level `info`, no file sink,
mutating-command timeout 3600 seconds, read/check timeout 30 seconds,
lock wait 30 seconds and offline doctor
checks by default. MCP uses JSON Lines on diagnostic stderr even when a CLI
human default would be plain; an explicit `logging.format=plain` is rejected
when launching MCP rather than silently ignored. Workspace UX defaults remain in the versioned
distribution and closed customization schema, not duplicated here.

On macOS the XDG-style user path is intentional and documented; no hidden
`~/Library/Preferences` duplicate is loaded. Inside a container these paths
refer to the container filesystem/user, never to the host. An operator MCP
server may run under a service account and therefore reads that account's
user layer, not the interactive human's home. An absent optional file is
recorded as absent. An unreadable or malformed present file fails with a
source-specific diagnostic. `settings.json` follows the closed, versioned
[process-settings schema](process-settings.schema.json); duplicate JSON keys,
unknown keys, unsupported versions and null values fail. Paths in it resolve
relative to the file containing them unless a key declares another anchor;
`logging.file` is anchored to the process state directory and then checked
against permitted sink roots. No file is created by ordinary reads.
`--settings-file PATH` may replace the default **user** settings file for one
invocation; it is anchored to the original working directory, cannot suppress
the system layer and is not an MCP tool-input field. An operator service
selects it only in its launch configuration, never from a project request.
An explicitly selected but missing settings file is an error, unlike an
absent default optional user file.
An invocation-relative `--project` or `--env-file` is anchored to the original
working directory before any project discovery; `--config` is anchored to the
resolved project root. `AIBOX_ENV_FILE` must be absolute. A `--log-file`
relative path is anchored to the invoking user's state directory and is
rejected if it escapes the permitted state root. Symlinks are resolved before
permission checks. The process never changes its working directory merely
to make a relative path appear valid.

Only explicitly named env files are loaded: `--env-file PATH` may appear more
than once, processed left to right, and `AIBOX_ENV_FILE` may name one default
file when no flag is present. There is **no** automatic `.env` search in the
current directory or repository. Env files use UTF-8 `KEY=VALUE` records
(split on the first `=`, preserve the remaining value bytes except a final
line ending), blank lines and full-line `#` comments; quotes are literal,
not shell syntax. No multiline values, shell execution, interpolation,
command substitution or `export` syntax. Keys must be documented `AIBOX_*` keys;
duplicate keys within one file fail, later explicit files override earlier
ones, and inherited process environment overrides every env file. Env-file
paths are canonicalized and must be regular files; files containing secret
values require private permissions and are never copied into an image or
receipt. A repository may include an example env file but may not select an
operator's private env file or grant authority through its contents.
Env-file keys configure only the aibox process; they are not automatically
exported to a child `devcontainer` invocation or injected into the container.
Any provider/container variable needed there uses the upstream native
configuration or an explicitly approved bounded pass-through list.

The public env mapping is mechanical only for the bounded process-settings
schema: `AIBOX_<SECTION>__<KEY>` (uppercase snake-case from camelCase) maps
to one leaf, for example `AIBOX_LOGGING__LEVEL=debug` and
`AIBOX_OUTPUT__FORMAT=json`. The system/user settings file also supports
`workspace.theme`, `workspace.mode` and `workspace.layout` as a narrow
overrideable UX subset; their values are additionally checked against the
closed `customization.schema.json` choices. Other workspace UX keys remain
file-only in `devcontainer.json` or native user overrides, rather than
silently acquiring a second broad schema. Booleans accept `true`/`false`;
integers use base 10; enums are case-sensitive; empty is a real string only where the
schema permits it, otherwise an error. Absence inherits; `null` is invalid;
there is no magic `unset` string. Lists/objects are file-only unless a
specific scalar CLI flag is documented. The registered v1 `AIBOX_*` runtime
identifiers retained from the v0 environment ledger remain separate from
process-setting keys; they are explicitly allowlisted and documented rather
than accidentally interpreted as settings. Unknown `AIBOX_*` names fail with a
nearest-key suggestion; standard upstream `DEVCONTAINER_*`, `DOCKER_*`,
`COMPOSE_*` and provider variables are passed to their owners under the
bounded process environment, not parsed as aibox settings.

CLI flags override matching settings for one invocation only. The common
surface is `--format`, `--color auto|always|never`, `--log-level`,
`--log-format`, `--log-file`, `--timeout`, `--settings-file`, `--env-file`,
`--project` and `--config` (the selected native Dev Container definition). Container-local
`refresh` additionally accepts invocation-only `--theme`, `--mode` and
`--layout`; these never rewrite project defaults. A host-side `attach`
does not inherit or apply those local presentation overrides. `--env-file`,
`--project` and `--config` are bootstrap selectors, not persisted values.
Command-specific flags such as `--offline`, `--frozen-lockfile`, `--no-cache`, `--tail` and
`--since` override only their corresponding operation inputs. No arbitrary
`--set`, hidden flags for all 115 UX keys or generic JSON override blob is
added. UX keys without a dedicated invocation flag can still be changed in
`devcontainer.json` or a native user override; a local agent may edit those
files and call `refresh`. An MCP request carries typed invocation fields,
not a CLI flag string or captured process environment. A higher-precedence
selector cannot replace an administrator's executable allowlist, operator
policy, permitted roots/contexts, credential source or approved log egress.

Configuration loading is side-effect free and ordered: parse files, env files
and environment; validate schema/types; apply authority checks; resolve
paths/identities; compute effective values; only then run a use case. Before
logging settings are valid, startup/parser failures use a minimal
redacted plain-stderr bootstrap diagnostic and exit 2; an MCP process still
writes no non-protocol bytes to stdout. Once configuration succeeds, a request
ID is generated if the caller did not provide one; startup/shutdown events
use a distinct process correlation ID. Caller-provided IDs are bounded and
validated, not trusted as authentication identities. Native Dev Container
configuration is resolved by the pinned upstream CLI under its
own rules. aibox reports the native input path/digest and upstream-resolved
values where documented, but does not claim its env/CLI layers rewrite native
Features, mounts, lifecycle hooks or image fields. Build/up require a frozen
effective input digest covering the native definition, referenced files,
selected aibox settings and operator policy version.

`aibox inspect --effective-config --format json` (and the read-only MCP
equivalent) reports every supported key once in stable order: redacted value,
winning layer/source/path, loaded/absent files, overridden values where safe,
rejected settings, and whether the value is invocation-only. It also separates
native Dev Container values from aibox settings and native user-tool overrides.
Sensitive paths and secret-shaped values are redacted. Local mode omits host
policy details it cannot access. The view never evaluates hooks, starts a
container, reads an unrelated home directory or leaks an env-file value.
The v0 converter maps each old output/diagnostic/logging preference either to
this process-settings schema, to a standard native project field, or to an
explicitly unsupported disposition. It does not
automatically promote `.aibox-local.toml` credentials into a committed
settings file or begin loading a repository `.env`. Existing v0 logs remain
readable as historical files until the operator's retention plan removes
them; they are not re-labeled as v1 operation evidence.

## Semantic events, sinks and retention

**R-LOGGING:** every CLI/MCP operation emits structured semantic events,
redacted **before** any sink or buffer. An event has `schemaVersion`, UTC
timestamp, severity (`trace`, `debug`, `info`, `warn`, `error`), stable event
name, component, operation, request/correlation ID, initiator and executor
identities where known, target ID if authorized, outcome/error category,
attributed child tool, per-request monotonic sequence and a bounded map of
redacted fields. Start, validated,
authorized/denied, dependency invocation, progress, partial effect, finish,
cancel and cleanup have distinct event names. Never infer success from a log
line; the typed result and durable operation receipt remain authoritative.

| Sink | Default and format | Location / failure rule |
|---|---|---|
| CLI diagnostic stderr | Enabled; concise plain text for interactive use, JSON Lines when selected | Never write progress/logs to result stdout. If stderr closes, preserve the operation result and receipt; do not rerun work. |
| MCP diagnostic stderr | Enabled; JSON Lines | MCP stdout is protocol-only; no banner, child output or log event on stdout. |
| Local rotating file | Opt-in, JSON Lines | User/container state directory by default; explicit allowed path only. Create directories `0700` and files `0600`, never group/world readable; default rotation is 10 MiB, seven files and seven days (whichever is shorter). |
| Host service manager | Supported by capturing operator-server stderr under systemd/launchd or another operator-controlled collector | No custom network transport required; document unit/job integration and collector retention separately. |
| Native runtime/child logs | Not an aibox sink | Capture bounded output with source attribution; redact before forwarding or persistence. `aibox logs` reads only declared, authorized sources. |
| Evidence/receipts | Separate durable store | Minimal non-secret state under the state directory; rotation of operational logs never deletes it. |

The default level is `info`; `trace`/`debug` are opt-in and still redacted.
Rotation/retention overrides are bounded by policy and semantic limits
(1–1024 MiB, 1–100 files, 1–365 days); increasing them cannot weaken an
administrator's retention maximum or disk quota. A file sink without an
explicit rotation setting uses the defaults above.
`off` disables ordinary diagnostic events, not required security/operation
receipts. `NO_COLOR` wins over color preferences for human output; non-TTY
defaults to plain with no spinner. Logging format is `plain` or `jsonl` and
is independent of requested result `--format`. A `--log-file` path requests
an additional sink, never redirects result stdout; a project file or env
file cannot enable a new external destination. Direct remote telemetry/OTLP,
syslog network egress and arbitrary URL sinks are outside v1 unless separately
specified, authorized and tested. Platform collectors may export logs under
the host operator's separate policy.

The file sink uses an OS-protected append/rotation strategy safe for concurrent
processes; partial JSON lines are detectable and ignored by readers. The
event queue is bounded (default 1024 events). On backpressure, drop trace/
debug first, count losses in a `logging.dropped` warning and attempt a bounded
direct-stderr fallback for warn/error; never block a use case indefinitely or
claim a dropped event was persisted. Flush is attempted for at most two
seconds at normal shutdown and the flush result is reported. On sink
failure, emit a bounded stderr warning and mark the sink degraded in the
result where possible. If the configured sink is marked `required`, refuse a
new mutating operation **before** side effects; once started, a mid-operation
sink failure is recorded as a partial observability failure, not a request
replay. Required durable receipts/evidence always fail closed for mutation
when they cannot be written. Read-only help/version/guide calls can continue
with a visible warning. Operator policy controls file roots, retention
minimums and whether a sink may be disabled; project/env/CLI precedence
cannot weaken those bounds.

Redaction covers credential values, secret-shaped strings, sensitive paths,
provider prompts, personal data, child output and chunk boundaries. Known
secret values are replaced before formatting even in debug mode; no raw-log
flag is provided in v1. File permissions, log rotation, retention and backup
exclusions are tested on Linux/macOS and container-local mode. `aibox logs`
is a bounded, read-only view (source allowlist, tail ≤1000, optional time
filter) of authorized runtime/aibox logs; it is not an arbitrary file reader.

## Conformance tests

AC-CONFIG-SOURCES exercises every layer and permutation, multiple env files,
ambient environment isolation, CLI/MCP parity, selectors versus settings,
invalid/duplicate/unknown keys, missing and unreadable files, relative-path
anchors, null/empty/false, redacted provenance and authority denial. Tests
use temporary homes and do not read developer settings or credentials.
AC-LOGGING exercises every severity, sink, format, concurrent writes,
rotation/retention, child attribution, cancellation/partial outcomes,
NO_COLOR/non-TTY, closed stderr, disk-full/permission failure, required-sink
fail-closed behavior, MCP stdout purity and secret leakage across chunk
boundaries. Evidence/receipts are checked separately from logs.

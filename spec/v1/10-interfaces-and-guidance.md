# 10. CLI, MCP, diagnostics and guidance implementation

This chapter closes the adapter-level design of R-INTERFACES, R-DOCTOR,
R-RESULT, R-GUIDANCE and R-POLICY. There is one Go application use-case layer;
CLI and MCP never call one another. Operation names and result schema in
chapter 3 and `operation-result.schema.json` are normative for v1.

## Public command grammar

`aibox <command> [--project PATH] [--config PATH] [--format human|json|yaml]`
is the operator CLI shape. The project selector defaults to the current
working directory only after canonicalization and policy validation. `up`,
`build`, `stop`, `remove`, `rebuild`, `status`, `inspect`, `logs`, `doctor`,
`attach`, `operation status ID`, and `mcp serve` are the bounded verbs. `refresh` and `doctor`
are also available inside a container. `migrate preview|apply|rollback`
is a separate explicit utility, not a normal configuration setter. Human
`remove` and `rebuild` require a displayed exact target/disruption summary;
`--yes` is accepted only with an expected input digest and resource ID, so
scripts cannot silently approve a changed target. MCP clients submit those
same identity fields; client UI confirmation is never a substitute for server
policy. No generic `exec`, `config set`, `install addon`, `purge --all` or
unbounded file-read command is exposed.
Chapter 16 supplies the complete settings hierarchy and common output/logging
flags (`--settings-file`, `--env-file`, `--color`, `--log-level`, `--log-format`, `--log-file`,
`--timeout`). `--config` selects a native Dev Container definition, not the
aibox process settings file. Flags are invocation-only; no config-write
command is implied.

All machine calls return the closed `aibox.operation-result/v1` envelope.
The [closed request schema](operation-request.schema.json) defines
`requestId`, `projectRoot`,
`configPath`, `expectedInputDigest`, and an operation-specific selector:
`build` has `frozenLockfile`; `start` has no implicit attach; `stop` and
`remove` require `resourceId` and `runtimeContext`; `rebuild` additionally
requires `acknowledgeDisruption`; `read_logs` requires `source`, `tailLines`
(1–1000) and optional `since`; `check` has a named check set and `offline`
defaulting true; `refresh` names a local UX scope and bounded invocation-only
theme/mode/layout overrides. Reads do not require an expected input digest;
they return the digest needed for a subsequent mutation (chapter 18). Unknown fields,
empty IDs, unbounded line counts and invalid paths are rejected before any
side effect. The Go semantic checker additionally rejects fields irrelevant
to the selected operation, `tailLines > 1000`, a false disruption
acknowledgement, an operator selector in local mode, and non-canonical paths.
Before V1-09 ships, each tool's SDK-exposed schema must be generated from or
conformance-tested against this file; a type wrapper may strengthen but not
silently weaken it.

Exit classes are fixed: 0 succeeded/no-change; 2 invalid input; 3 denied;
4 missing/incompatible dependency; 5 failed/partial; 6 explicit operator
action or rebuild required; 124 timeout; 130 interrupted. JSON stdout contains exactly
one envelope and a newline. Human progress and warnings go to stderr; `--format
human` is a projection of the same envelope. `--format yaml` must serialize
only the result schema, not a richer hidden state. All errors include a stable
code, actionable message, retryability flag and bounded next-action reference.

## MCP server and capability isolation

The default `aibox mcp serve` is a stdio, read-only knowledge server. It
registers resources `aibox://guides/v1/index` and
`aibox://guides/v1/{topic}`, a resource template for declared topics, plus
`product/version` metadata. Operator mode is enabled only by a policy file
outside the project root and explicit launch argument. The server registers
typed tools `build_environment`, `start_environment`, `stop_environment`,
`remove_environment`, `rebuild_environment`, `inspect_environment`,
`read_logs`, `check_environment`, `inspect_operation`; local mode registers only
`inspect_workspace`, `check_workspace`, `refresh_workspace`, local-only
`read_logs`, local receipt `inspect_operation` and guidance.
Registering a tool does not authorize its invocation: every call checks
process identity, policy, target and input digest server-side. No network
listener, OAuth flow or remote proxy is in v1.

Use the official Go MCP SDK and its schema facilities; do not manually
implement JSON-RPC. Bound returned content and include a structured result
or protocol error that points to the same stable application code. A long
operation remains a single call with progress notifications where the client
supports them. Cancellation follows chapter 9. Before publishing, test a
current Codex client, a current Claude client and a protocol-level SDK
client. If either agent client cannot access resources, expose a narrowly
scoped `get_guide(topic, version)` compatibility tool backed by the identical
guide store; do not duplicate content or broaden file access. MCP/CLI parity
tests compare normalized result envelopes for identical requests and
policies, excluding transport timestamps and progress messages.

## Guide store and doctor registry

The guide manifest lists a stable topic ID, title, summary, applicable
product version, source Markdown path, source digest and related topics.
Initial topics: `start`, `tools`, `themes`, `tmux-header`, `terminal-title`,
`sidebar`, `audio`, `previews`, `review`, `doctor`, `migration`, `authority`,
`processkit`, `latex`. Each guide states prerequisites, exact file to edit,
which actor may edit it, example, validation, activation (`refresh` versus
rebuild), expected result, undo, and link to public docs. A theme guide must
distinguish project default from personal override and explain why the latter
wins. Generate the embedded index from versioned public Markdown and fail
CI/local verification if digest or topic linkage drifts. Unsupported topic
or version is a typed not-found; the server never crawls a repository to
answer an open-ended question.

Doctor checks are registered with `code`, severity (`error`, `warning`,
`info`), scope (`local`, `operator`), required capability, expected latency,
and remediation guide. Checks are pure/read-only: schema and cross-field
validity, Feature reference/lock consistency, missing managed assets,
override ownership, runtime/CLI readiness (operator only), local harness
health, audio bridge readiness and processkit integration. The default
offline run performs no registry/network access; an optional online
qualification check is explicit. Each finding has state `passed`, `failed`,
`skipped`, `unavailable` or `not_authorized`, observed timestamp, source and
redacted evidence. Stable code ordering makes output diffable. `doctor`
never fixes, installs, invokes lifecycle hooks or upgrades a dependency.

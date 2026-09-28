# 3. CLI, MCP and authority contracts

**R-INTERFACES:** MCP is the primary agent-facing adapter; the human CLI
uses the same Go core. CLI-only presentation conveniences do not require MCP
equivalents. Core operations have capability-equivalent typed inputs/results.
Do not wrap the human CLI and parse its prose to implement MCP.

## Operation set

Command names and exit classes below are the v1 design; chapter 10 fixes their
grammar and input semantics. Chapter 18 and the closed per-operation schemas
fix the machine contracts; real client conformance tests remain G04. Existing v0 spelling can be retained
where it remains unambiguous.
Availability is incremental, not all-or-nothing: V1-03 first exposes
version/help and read-only local inspection through the Go executable; V1-08
adds authorized lifecycle verbs, and V1-09 completes the guarded MCP adapter
and thin CLI. Every intermediate executable rejects unavailable operations
with a documented error instead of implying they already work.
The [command audit](command-disposition.md) covers removals/delegation, while
[commands.json](ledger/commands.json) preserves every original declaration.

| CLI | MCP | Context | Effects |
|---|---|---|---|
| `up` | `start_environment` | Operator | Create/start the exact configured environment; no attached interactive session in MCP. |
| `stop` | `stop_environment` | Operator | Stop, retain data. Repeated stopped result is success/no-change. |
| `remove` | `remove_environment` | Operator | Remove identified runtime only, persistent data retained by default. No global prune. |
| `build` | `build_environment` | Operator | Build through upstream CLI without replacing the running container. |
| `rebuild` | `rebuild_environment` | Operator | Explicit replacement/restart with interruption disclosure and consent. Never an alias of mere cache bypass. |
| `status` / `inspect` | `inspect_environment` | Scoped operator or local | Structured state, versions, effective config provenance and available resources; local cannot query host runtime. |
| `logs` | `read_logs` | Scoped operator or local | Bounded tail/filter of declared logs; redact before return. No arbitrary file access. |
| `doctor` | `check_environment` | Scoped operator or local | Read-only check suite. Offline schema checks never execute lifecycle hooks. |
| `refresh` | `refresh_workspace` | Local | Validate and activate managed local UX changes; report rebuild-required rather than crossing authority. |
| `operation status ID` | `inspect_operation` | Scoped operator or local | Inspect an authorized durable receipt after interruption; never resume effects automatically. |
| `attach` / recovery option | No interactive equivalent | Operator human CLI | Delegate entry via native exec/tmux; recovery bypasses broken tmux/Yazi/status. MCP returns connection instructions. |
| `mcp serve` | Not itself an MCP tool | Startup | Start stdio server with operator-installed capability policy. Read-only knowledge mode requires no runtime access. |
| help/version/completion | Resources/server metadata | Both | Human shell convenience and product identity. |

No config setter/editor, addon package manager, processkit entity mutation,
generic host exec/shell, daemon control plane, universal remote provider SDK
or lifecycle tool exposed inside the workspace. One-time conversion and
scoped recovery are separately reviewed utilities, not backdoors in `doctor`.

**R-DOCTOR:** checks schema/semantic consistency, declared asset existence,
upstream compatibility, integration integrity and context-appropriate readiness.
Each finding includes stable code, severity, setting/file, explanation,
remediation resource and status (`passed`, `failed`, `skipped`, `unavailable`,
`not_authorized`). No install, config rewrite, migration creation or host probe
from a local check. Clearly disclose network-dependent resolution; offer a
network-free schema/content pass. A clean schema alone does not authorize
host hooks, mounts, privileges or executable downloads.

## Machine result and execution contract

**R-RESULT:** each core operation returns a versioned result containing
operation/request identity, target identity, input digest, outcome, changed
resources, warnings, evidence references and structured error/next action.
The [closed result schema](operation-result.schema.json) fixes the
common envelope for both adapters. It records initiating and executing actors
separately; a resource ID is never substituted for an actor. Local results omit
host runtime context/resource identity rather than inventing it. Evidence
references point to durable, separately retained records, not log lines.
`schemaVersion` is `aibox.operation-result/v1`; receivers reject unknown major
versions and unknown fields rather than silently reinterpret them. A wire-field
addition therefore needs a new negotiated schema version, not an unannounced
minor addition to this closed schema. Initial v1 implements only the v1 wire
contract; retaining old records does not imply executing old requests.
Chapter 18 defines operation-specific data and early-failure envelopes: builds
and absent environments do not invent container IDs; malformed input may have
neither target nor digest. All unsuccessful outcomes carry structured errors.
Target identity binds the runtime endpoint/context plus resolved project and
resource identifiers; names/labels alone are not proof of ownership. Reject
ambiguous targets and changed identity before destructive execution.

Outcomes: `succeeded`, `no_change`, `failed`, `cancelled`,
`partial`, `rebuild_required`, `refused`, `timed_out`. Evidence of partial effects survives failures.
CLI JSON stdout contains only the result; human diagnostics/progress use
stderr. MCP stdio stdout contains only protocol messages. Human table/YAML
views may be projections of the same result, not separate facts.
Chapter 16 specifies the structured logging event and sink contract;
diagnostic logs are never a second result channel or a substitute for the
durable evidence referenced by this envelope.

Exit classes: 0 success/no-change; 2 invalid input; 3 denied authority;
4 missing dependency/incompatible environment; 5 operation failure/partial;
6 explicit operator action required; 124 timeout; 130 interruption. Chapter 18
fixes the error-code catalog and outcome precedence. Never return
success after a child failed.

**R-EXECUTION:** invoke public tool binaries with argument arrays, explicit
working directories, resolved executable provenance and a bounded environment.
Do not interpolate repository strings into a host shell. Native hooks remain
executable input and need trust checks, not a claim of argument safety alone.
Respect supported upstream exit/JSON contracts and pin compatibility tests.

Operations have explicit timeout/cancellation. Propagate cancellation to the
delegated process group and inspect resulting state; never assume cancelling a
client undoes a completed runtime effect. Serialize conflicting operations on
the same target; bounded locks recover from process crashes. Retrying reads
is safe; mutating retries inspect actual state and input identity first.
Duplicate removal must not delete a newly created container with the old name.
No blind retry of rebuild or delete and no invented durable success receipt.

Long builds remain synchronous calls. Chapter 18 defines durable receipts and
read-only operation inspection across client sessions; it does not authorize
a background workflow engine. Client compatibility and cancellation are tested
before the relevant phase ships.

## Knowledge and progressive disclosure

**R-GUIDANCE:** provide a small index plus versioned task-sized MCP resources
and resource templates. Stable topic identifiers cover installation, tool
selection, themes, tmux headers, titles, sidebar, audio, previews, review,
diagnostics, migration and authority. Guidance is generated from or identical
to canonical Markdown; public docs work without an MCP client.

Each recipe states product/version applicability, prerequisites, exact file
ownership and paths, edit example, validation, activation method, expected
result, rollback and deeper references. “Edit theme” must distinguish editing
the project default from a local override and identify whether live refresh
suffices. Do not load every document into tool descriptions. Search, if needed,
returns a bounded set of declared resource references; no arbitrary path reads,
internal context, prompt-driven shell action or live generated advice.

Unsupported topic/version returns an explicit not-found/incompatible response;
never silently serve instructions for another release. Test resource discovery
and actual model access in selected clients; resource support is not uniform.
Use read-only resources rather than a collection of fake document-retrieval
“operational” tools. A compatibility discovery adapter must share the same data.

## Security and permissions

**R-POLICY:** start in knowledge/read-only mode unless the operator explicitly
configures broader authority outside repository control. Operator access is
scoped to declared workspace roots/runtime contexts and operations. Local
mode lacks host credentials/endpoints regardless of requested arguments.
Server-side enforcement applies even if a client omits confirmation UI.

Default transport is local stdio. Remote listening/authentication is outside
the initial support contract; it needs a separate threat model and approval.
Standard SDK protocol support does not confer permission to expose an operator
server over a network. Terminal title/notification strings are sanitized and
bounded; local tool configuration may execute local commands and remains
subject to ordinary user trust, but cannot acquire host authority.

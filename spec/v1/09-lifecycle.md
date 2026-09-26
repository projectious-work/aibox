# 9. Operator lifecycle implementation

This chapter is the implementation contract for R-EXECUTION, R-RESULT,
R-AUTHORITY and R-RECOVERY. Only an external operator process may execute
these use cases. `devcontainer` owns definition resolution, Feature
installation, image build, container creation and user command execution.
The pinned upstream CLI currently does not implement `stop` or `down`;
the narrow native-runtime adapter below fills those gaps without becoming
a second Dev Container engine. A change in upstream support is reviewed before
replacing the adapter, with equivalence tests retained.

## Target identity and discovery

An operator request identifies a canonical project root, a selected
`devcontainer.json` path under that root, a runtime endpoint/context, and,
where applicable, a configuration name. The runtime identity is the tuple
`(runtime kind, endpoint fingerprint, context, native resource ID)`, not the
container name. Resolve symlinks before root and allowlist checks; reject a
path whose traversal escapes the allowed root. Preserve the original source
path in evidence but never use it for destructive selection. Read-only
discovery may inspect a specified native resource ID or resources carrying
the documented Dev Container labels for that exact canonical workspace and
configuration; it must reject zero-to-many ambiguity instead of selecting
the first match. Labels are corroborating evidence, not authorization.

Before a mutation, hash the effective on-disk input set (definition, referenced
Dockerfile/Compose files and operator policy version), bind it to the request,
then re-read it after acquiring a per-target lock. A changed digest returns
`input_changed` without executing. Runtime context, endpoint and resource ID
are rechecked immediately before stop/remove. The adapter refuses global
prune, broad Compose project deletion, volume deletion, wildcard names and
all operation requests without an exact target. For an initial `up`, the
upstream CLI may create a resource; after success, inspect its native ID and
record the binding.

## Command delegation and transaction steps

| Use case | Preconditions and delegated command | Postcondition and failure evidence |
|---|---|---|
| `build` | Validated definition, pinned `devcontainer` executable and runtime readiness; invoke public `devcontainer build --workspace-folder <root>` with documented configuration selector/lockfile option for that pin | Capture upstream build result and image identity; do not stop a running workspace. Child failure is `failed`, never success with a warning. |
| `start` (`up`) | Inspect exact target; if already running with same input, `no_change`; otherwise invoke public `devcontainer up --workspace-folder <root>` plus documented options | Inspect resulting resource ID, state, definition labels and image; a CLI exit 0 without inspectable expected resource is `partial`. Start does not imply an attached terminal. |
| `stop` | Require exact running resource ID and reconfirm endpoint/context/input; invoke native runtime stop with that ID only | Reinspect: stopped is success; already stopped is `no_change`; vanished/wrong identity is `partial` or `identity_changed`, never stop a substitute. Volumes/home retained. |
| `remove` | Stop if necessary only after explicit remove intent; exact ID plus ownership proof; invoke native runtime remove by ID, never Compose down for an entire project | Confirm ID gone, preserve named/bind volumes and project files. Already removed is `no_change` only if receipt/identity proves it was the same resource. |
| `rebuild` | Explicit disruption acknowledgement and retained-state plan; build a candidate first, then stop/remove exact old container, then `up` | If candidate build fails, old environment stays running. After old removal, failure is `partial` with retained state and documented recovery; no automatic volume purge. |
| `inspect` | Read-only upstream configuration read plus bounded runtime inspect where authorized | Report separate declared, observed and unknown fields with timestamps; do not execute lifecycle hooks. |
| `attach` | Human CLI only, running exact resource and explicit terminal | Use public `devcontainer exec` to enter; recover mode bypasses tmux/session startup. Terminal IO is not serialized as an MCP result. |

Use an argument-vector process runner; executable path is resolved from an
operator-controlled allowlist, not a repository string. Supply a narrow
environment with only required runtime variables. Child stdout/stderr are
bounded, redacted and recorded separately; JSON output is decoded only where
the selected CLI documents JSON. Never parse human progress as authoritative
state. Standard Dev Container lifecycle hooks are executable project input:
before a build/up with such hooks, display a deterministic inventory and
require the operator trust policy to approve the frozen digest. A request
from a repository cannot approve itself. Reading configuration, `doctor` and
migration preview must not run hooks.

## Concurrency, retries and cancellation

One mutating request per `(canonical root, runtime endpoint, context)` holds
an OS-enforced lock plus a non-secret receipt. Bounded lock acquisition returns
`busy` with current operation ID. Stale receipt recovery is inspection-first;
PID existence alone is insufficient. A repeated request ID returns a saved
terminal result only if the digest and target identity still match. A new
request ID against a changed target is a new request, not an idempotent replay.

On timeout or client cancellation, signal the child process group, wait a
bounded grace period and inspect actual runtime state. Return `cancelled`
only if no effect is observed; otherwise `partial` with concrete resource IDs
and next action. Killing `devcontainer up` may leave a running container;
neither the CLI nor MCP may promise rollback in that case. Crash-restart
tests kill the wrapper after each recorded step and verify that no subsequent
`remove` touches a replacement with the same name. Child logs are redacted
before persistence; receipts contain hashes and references, not raw secrets.

## Acceptance fixtures

Implement fake `devcontainer`, Docker and Podman executables to assert exact
argument vectors, working directory, exit attribution, context selection,
no shell evaluation, output bounds and cancellation. Disposable-runtime tests
cover direct CLI versus wrapper build/up, multiple similarly named projects,
two runtime contexts, stopped/running/vanished states, Compose sidecars,
failed hooks, identity swap between inspect and remove, and retained named
volumes. Record input and output artifacts for each target cell in chapter 14.

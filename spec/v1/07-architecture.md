# 7. Go architecture and code boundaries

This chapter defines R-CORE, R-BOUNDARY, R-REUSE, R-RESULT and R-EXECUTION.
Package responsibilities and dependency boundaries apply independently of
package names. The Go binary is an optional convenience
and agent interface around a workspace that remains valid without it.

## Applications and processes

Build one Go module and one distributable `aibox` executable. The binary has
three entry contexts, all backed by the same application use cases:

1. **Operator CLI:** a human on an authorized host invokes a bounded action.
2. **Operator MCP:** an external management agent calls a stdio server launched
   under independently configured operator policy. The server has no network
   listener in v1.
3. **Workspace-local CLI/MCP:** a process inside the container receives only
   local resources, validation, refresh and diagnostics. Its compiled binary
   can be identical; the process has no host socket, credential, mount or
   reachable operator endpoint. Supplying `--operator` cannot add capabilities.

The binaries and Feature-installed local helpers do not implement a daemon,
container registry, scheduler, processkit entity engine or Dev Container
compiler. Native `devcontainer.json`, Dockerfile and Compose files are parsed
by their owning tools. aibox reads just the fields needed for its own UX,
authorization review, migration and reporting; it never writes a transformed
replacement for a user-owned Dev Container definition during normal operation.

## Proposed source tree

| Path | Single responsibility | Must not do |
|---|---|---|
| `cmd/aibox` | Process entry, signal context, adapter selection, exit code | Business decisions or child-tool parsing |
| `internal/app` | Typed use-case orchestration and transaction boundaries | Shell construction, JSONC parsing, terminal rendering |
| `internal/contract` | Operation inputs/results, findings, error codes, schema versions | Runtime calls or mutable globals |
| `internal/project` | Discover project root and declared Dev Container inputs; read ownership manifest | Execute hooks or rewrite native configuration |
| `internal/config` | Strictly parse layered process settings, env files/vars, aibox UX extension and invocation overrides; validate authority and resolve provenance | Reimplement upstream Dev Container merge/build semantics or read a host home in local mode |
| `internal/logging` | Build redacted semantic events and render them to bounded stderr/file/collector sinks | Treat logs as receipts, write to MCP stdout or start network export from project input |
| `internal/policy` | Authorize caller, operation, target and frozen input digest | Accept repository-controlled grants |
| `internal/identity` | Bind project, runtime endpoint/context, native resource IDs and owned labels | Treat display names as deletion proof |
| `internal/devcontainer` | Invoke documented CLI commands and decode documented outputs | Import private Node modules or build a second compiler |
| `internal/runtime` | Narrow stop/remove/inspect adapter where CLI has a proven gap | Arbitrary container-engine API or global prune |
| `internal/process` | Resolve allowed executables; run argument arrays with bounded IO, environment, cancellation | Interpolate project strings into a shell |
| `internal/operation` | Per-target lock, state transitions, durable minimal receipts, interrupted-run inspection | General workflow engine or desired-state database |
| `internal/diagnostic` | Registered read-only check functions and stable finding order | Implicit network, install, fixes or lifecycle effects |
| `internal/guidance` | Embed/version public task recipes, index and template variables | Read arbitrary project paths or generate LLM advice |
| `internal/migration` | One-time v0 conversion preview/apply/rollback and named-state conversion | Recurring configuration manager |
| `internal/workspace` | Container-local safe refresh of aibox-managed UX output | Host lifecycle, Docker/Podman socket or host command bridge |
| `internal/output` | Human/plain/JSON/YAML projections of typed results | Invent results from renderer state |
| `internal/mcp` | SDK-based tool/resource registration and stdio lifecycle | Business logic or parser of CLI prose |
| `internal/cli` | Commands, arguments, help, completion and input projection | An independent operation implementation |
| `schemas`, `guides` | Embedded public contracts and versioned how-to content | Private processkit context or secrets |
| `test/blackbox`, `test/integration`, `test/e2e` | Binary, dependency and disposable workflow tests | Tests against a developer's real home/runtime by default |

Keep production packages under `internal/` until a public Go API has an actual
consumer and compatibility commitment. A package interface belongs to its
consumer. Avoid generic `util`, `manager` and `helpers` packages. Default local
tests use fake executables and temporary roots; no host runtime is required.

## Direction of dependencies

`cli` and `mcp` call `app`; `app` depends on narrow ports whose concrete
implementations live in `devcontainer`, `runtime`, `project`, `policy`,
`operation`, `diagnostic`, `guidance`, `migration`, `workspace`, `config` and
`logging`. `contract`
has no adapter dependency. `output` projects `contract` only. The executable
composes the graph. A package-level test rejects an import from any local
workspace package into `runtime`, or from an adapter into another adapter's
private implementation.
`config` resolves settings before use-case construction; `logging` receives
typed events after redaction and has no route back into operation decisions.
The exact source/sink contract is chapter 16. Prefer Go's standard structured
logging facilities for sink rendering rather than a separate logging
framework unless conformance tests reveal an unfilled requirement.

The core operation shape is intentionally small:

```go
type Request struct {
    ID string
    Actor Actor
    Target TargetSelector
    ExpectedInputDigest string
    Operation OperationKind
    Options OperationOptions
}

type UseCases interface {
    Build(context.Context, Request) (Result, error)
    Start(context.Context, Request) (Result, error)
    Stop(context.Context, Request) (Result, error)
    Remove(context.Context, Request) (Result, error)
    Rebuild(context.Context, Request) (Result, error)
    Inspect(context.Context, Request) (Result, error)
    Check(context.Context, Request) (Result, error)
    RefreshLocal(context.Context, Request) (Result, error)
}
```

These are explanatory type signatures, not permission to add generic
`Execute(command string)` or arbitrary hooks to the public API. Each operation
has its own validated input and stable schema before implementation. CLI and
MCP adapters translate into these same inputs; they must return equivalent
result envelopes for equivalent logical requests.

## Resource and operation state

Do not store a parallel desired-state model. The requested definition is the
current standard files; current runtime state is inspected from the selected
runtime. Persist only non-secret operation receipts needed to recognize
interruption and partial effects. A receipt contains request ID, operation,
actor IDs, project and runtime identities, input digest, start/end timestamps,
last confirmed step, outcome and references to redacted evidence. It never
contains copied config files, credentials, raw child output or authority grants.

Each mutating operation follows `validated → authorized → locked → revalidated
→ executing → inspected → recorded`. The lock key binds canonical project root
and runtime endpoint/context; different contexts do not share a name-only
lock. A crash leaves an incomplete receipt. The next operation inspects the
actual runtime and returns recovery guidance before retrying any destructive
step. Result serialization is deterministic (stable collection ordering where
domain order does not matter); machine stdout contains one envelope only.

## Build and dependency discipline

The Go module pins an accepted stable Go version and all direct modules in
`go.mod`/`go.sum`. The official Go MCP SDK is the sole protocol implementation.
Prefer the standard library for command parsing/process execution unless a
dependency makes a demonstrable safety or usability improvement. JSONC and
schema validation use maintained libraries with positive/negative conformance
fixtures; no new parser is authorized. Preserve existing Rust or shell runtime
UX glue when it is smaller and safer than rewriting it in Go; the Go decision
applies to the new application core and adapters, not every distributed asset.

Every direct dependency and tool binary must have a purpose, owner, tested
version/digest, license, provenance, transitive-weight and vulnerability
assessment recorded before a support claim. A source snapshot or README
observation is a candidate, not a compatibility test or release pin.

## Architecture acceptance

### V1-03 implementation status

The first Go slice lives in the root `go.mod`, `internal/contract`,
`internal/project` and `internal/process`. It defines the common operation
result envelope and CLI exit classes, canonical project-file containment,
plus a no-shell argument-vector runner with an explicit environment and
bounded, known-secret-redacted diagnostic tails. Offline tests
use the test binary as a fake child; `go test ./internal/...` requires no
container runtime. The module currently uses only the Go standard library.

This baseline now has a **source-built read-only v1 executable**, but cannot
close V1-03 yet. `cmd/aibox` provides version/help and local inspection of
the minimal native example. The inspection
returns the closed result envelope and a real input digest while rejecting
unavailable operator/lifecycle commands explicitly; it does not need Docker,
host policy or credentials. A customer can build and run the binary and follow
the matching v1 Hugo walkthrough. The runner now terminates child process
groups on Linux/macOS, and `internal/operation` has an exact-environment
advisory lock. Operation-specific payload checks, policy-owned executable
resolution, post-cancellation effect inspection and durable receipts remain
foundation work. These components must not be used for mutating lifecycle
operations until observed-effect behavior meets chapter 18.

- V1-03 binary black-box tests build `cmd/aibox`, invoke version/help and the
  local inspection demo, validate its JSON result and output separation, and
  prove unavailable host mutations fail without effects.
- Static import checks and a code review find no Dev Container resolver,
  package manager, daemon, host bridge or broad provider SDK in the Go tree.
- V1-08 fake-executable and disposable integration tests prove delegated
  public commands use argument arrays, bounded environments, correct working
  directories, cancellation and attributed child failures; direct upstream
  and wrapped outcomes are compared under AC-BUILD.
- V1-09 extends the existing binary with guarded stdio MCP. Its black-box
  suite verifies CLI/MCP equivalence and stable result/finding schemas. A
  product-code change cannot claim parity from this document alone.

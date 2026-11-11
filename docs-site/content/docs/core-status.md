---
title: Go core implementation status
description: "Run the first read-only v1 Go preview and see its current limits."
weight: 20
---

# Go core implementation status

V1-03 implementation is complete with local evidence; publication is pending.
A source-built Go CLI supports version, help and
read-only local inspection through both CLI and stdio MCP. The first internal slice adds the operation result
envelope and stable exit-code mapping, canonical project-file containment,
plus a no-shell child-process runner and an exact-environment advisory lock.
The runner accepts only an already-resolved absolute executable and working
directory, receives an explicit environment, and returns separate bounded
stdout/stderr diagnostic tails with known secrets masked. Offline tests use a
fake child process and do not require Docker, Podman or Dev Container CLI.
On Linux and macOS, cancellation sends TERM to the child process group and
escalates to KILL after five seconds; other platforms fail closed.

Prerequisites are Go 1.27 and a clone of this v1 source branch. From the
repository root, build and try the current Linux/macOS customer preview:

```sh
go build -o /tmp/aibox-v1-preview ./cmd/aibox
/tmp/aibox-v1-preview --version
/tmp/aibox-v1-preview inspect --context local --project spec/v1/examples/minimal --format json
# An MCP client can launch: /tmp/aibox-v1-preview mcp serve --context local
# and call inspect_workspace with requestId and a canonical absolute projectRoot.
```

The command returns one `inspect_workspace` JSON result with a canonical
project root, SHA-256 digest of the declared native config manifest, its
source, and `"outcome":"succeeded"` plus `"state":"unknown"`. That state is intentional: this offline preview
does not contact Docker or claim a running container. It makes no project
changes; remove `/tmp/aibox-v1-preview` to clean up the locally built binary.
Run `go test ./...` to check the executable and core packages offline.

There is no installable release yet. The current runner is not safe
for mutating lifecycle calls: policy-owned executable resolution, input
binding, post-cancellation effect inspection, operation records and operation-specific
schema validation remain to be implemented. The
[canonical architecture contract](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/07-architecture.md)
defines those acceptance conditions.

The V1-03 implementation candidate has a repeatable local verification gate:
`./scripts/verify-v1-03.sh` builds the Go binary, exercises the successful
inspection and an invalid-project refusal, validates their result schemas, and
builds this site. Publication evidence is still required before the roadmap
can call V1-03 shipped. Each later roadmap phase must keep this executable
working and add paired CLI/MCP customer-runnable, documented demos. Public
alpha/beta releases are prepared at verified checkpoints, not automatically
after every phase; until publication, the demo is a preview, not a shipped
support claim.

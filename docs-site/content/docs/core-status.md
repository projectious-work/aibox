---
title: Go core implementation status
description: "Current internal v1 Go core and its unimplemented interfaces."
weight: 20
---

# Go core implementation status

V1-03 is in progress. The first internal slice adds the operation result
envelope and stable exit-code mapping, canonical project-file containment,
plus a no-shell child-process runner.
The runner accepts only an already-resolved absolute executable and working
directory, receives an explicit environment, and returns separate bounded
stdout/stderr diagnostic tails with known secrets masked. Offline tests use a
fake child process and do not require Docker, Podman or Dev Container CLI.

Maintainers can run the current Go checks from the repository root:

```sh
go test ./internal/...
```

There is no usable v1 CLI or MCP server yet. The current runner is not safe
for mutating lifecycle calls: policy-owned executable resolution,
process-group cancellation, post-cancellation inspection, receipts and
operation-specific schema validation remain to be implemented. The
[canonical architecture contract](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/07-architecture.md)
defines those acceptance conditions.

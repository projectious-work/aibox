---
title: Native tool selection preview
description: "Select reviewed Dev Container Features and understand V1-05 catalog limits."
weight: 30
---

# Native tool selection preview

V1-05 uses Dev Container `features`, image, Dockerfile and package inputs for
tool installation. `customizations.aibox` contains workspace preferences and
does not install tools. A project owns its native selections: add a Feature to
enable it, remove the Feature to disable it, and rebuild the container.
The base image must also be checked when absence matters.

The frozen v0 addon census has 98 tool entries. Source and OCI review found
four entries with an upstream Feature option; live installation on both Linux
architectures is still pending. The other 94 entries require owner review of
native installation proposals or removal. No addon has been removed.

| v0 tool | Upstream Feature | Version selection |
|---|---|---|
| Hugo | `ghcr.io/devcontainers/features/hugo` | `version: 0.165.0`, `extended: true` matches the frozen runtime recipe. |
| AWS CLI | `ghcr.io/devcontainers/features/aws-cli` | `version: latest` matches the unversioned v0 declaration. |
| Azure CLI | `ghcr.io/devcontainers/features/azure-cli` | `version: latest`, `installBicep: false` avoids an extra CLI. |
| GitHub CLI | `ghcr.io/devcontainers/features/github-cli` | `version: latest` matches the unversioned v0 declaration. |

The complete source review, immutable Feature manifest digests, options,
architecture claims and unresolved decisions live in
`spec/v1/ledger/v1-05-tool-selection.json` and
`dev-notes/V1-05-gap-review.md` in the source checkout. `latest` tool options
follow the corresponding v0 unversioned behavior, but do not freeze the tool
binary. Exact tool pins and two-architecture builds are required before
claiming release qualification.

## Try the Hugo selection

The customized source example selects the official Hugo Feature by OCI
manifest digest:

```jsonc
"features": {
  "ghcr.io/devcontainers/features/hugo@sha256:15465c956810aa164c9644b6df5ca36f6cac3e7a86b7dcc474a1fb830b21c509": {
    "version": "0.165.0",
    "extended": true
  }
}
```

Run the source-built preview to inspect the selected Feature reference:

```sh
go build -o /tmp/aibox-v1-preview ./cmd/aibox
/tmp/aibox-v1-preview inspect --context local --project spec/v1/examples/customized --format json
```

CLI and local MCP `inspect_workspace` return the same sorted Feature references
under `configuration`. Inspection is read-only and omits Feature option values
and other native fields, which may contain secrets. It does not prove that a
container built or the tool runs. The upstream Dev Container CLI owns
`build`, `up` and `exec`.

The source review found important limits in other official Features: Go and
Rust alter container tracing permissions; Kubernetes adds a Minikube volume
even when Minikube is disabled; Node bundles Yarn; Python bundles pip. Their
v0 tools remain in the unresolved review until the owner decides whether to
accept those changes or use the proposed native Dockerfile routes. The Copilot
Feature's `latest` selection runs an automatic update at container start, so
Copilot also remains unresolved until an exact install is chosen.

## Host qualification

On a clean checkout of the reviewed commit, run the V1-05 host gate with a new
absolute evidence directory:

```sh
./scripts/verify-v1-05-host.sh --output /absolute/new/evidence-directory
```

The script uses Python 3 and the same pinned Dev Container CLI bootstrap as
V1-04; an installed Node.js or npm is unnecessary. It checks that Hugo
`0.165.0` extended runs in the selected example and is absent from the minimal
example, then records build/up/exec logs, source commit, host architecture and
scoped cleanup in `evidence.json`. Run it on both linux/amd64 and linux/arm64
targets for two-architecture evidence. Each invocation checks the native
architecture of its own host; one passing host cannot establish both.

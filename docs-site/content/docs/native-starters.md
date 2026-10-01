---
title: Native workspace starters
description: "Choose and inspect an aibox v1 Dev Container starter."
weight: 20
---

# Native workspace starters

The v1 starter format is the standard Dev Container Template. It produces
ordinary project-owned files under `.devcontainer/`; the Dev Container CLI and
the selected Features own build and container lifecycle. aibox does not
translate native image, Dockerfile, Compose, mount, or lifecycle settings.

> **Preview status:** v1 has no published release or Template registry
> references yet. The source templates and rendered examples are for checkout
> testing. `aibox` currently provides read-only `inspect` in its Go CLI and
> local MCP server. It does not build, start, or exec into containers. Direct
> Dev Container CLI build/up/exec qualification is still pending.

## Choose a starter

| Starter | Use it when | Defaults |
|---|---|---|
| `minimal` | You want a small native Dev Container definition and will choose optional Features yourself. | Debian base, non-root `aibox` user, shell, tmux and Yazi, persistent named home volume, and an optional pinned Codex harness. It leaves the aibox customization extension inactive. |
| `curated` | You want the curated workspace defaults and can opt into additional Features. | The same base tools, native tmux/Yazi defaults, and a selectable pinned Codex harness. Additional Features remain explicit selections. |

The source Template definitions are in [`templates/src/minimal`](https://github.com/projectious-work/aibox/tree/v1.x-dev/templates/src/minimal)
and [`templates/src/curated`](https://github.com/projectious-work/aibox/tree/v1.x-dev/templates/src/curated).
Rendered candidate examples are available for [minimal](https://github.com/projectious-work/aibox/tree/v1.x-dev/spec/v1/examples/minimal),
[curated](https://github.com/projectious-work/aibox/tree/v1.x-dev/spec/v1/examples/customized),
[custom users](https://github.com/projectious-work/aibox/tree/v1.x-dev/spec/v1/examples/custom-user),
and [bind-mounted home](https://github.com/projectious-work/aibox/tree/v1.x-dev/spec/v1/examples/bind-home).

The source templates have no published registry references. In a source
checkout, render a starter into a new, empty project directory. The helper
copies native files and substitutes the selected harness; it does not run
container lifecycle commands.

```sh
PROJECT=/absolute/path/to/new-project
mkdir -p "$PROJECT"
node templates/render.mjs --template minimal --output "$PROJECT" --harness codex
```

The helper refuses to overwrite an existing `.devcontainer` directory. Use
`--template curated` for the curated starter or `--harness none` to omit the
optional Codex install. The custom-user and bind-home configurations are
separate rendered examples under `spec/v1/examples/`; review their files before
adapting a project. Do not edit files under `templates/src` to configure an
individual project.

## Inspect the native definition

From the repository root, build the preview binary and inspect one rendered
example:

```sh
go build -o /tmp/aibox-v1-preview ./cmd/aibox
/tmp/aibox-v1-preview inspect --context local --project spec/v1/examples/minimal --format json
```

Inspection is read-only. It checks the project and reports the native
configuration without building or starting a container. For protocol clients,
the local MCP server currently exposes the equivalent read-only
`inspect_workspace` tool; see the [current preview limits](core-status.md).

Direct upstream Dev Container CLI `build`, `up` and `exec` are the lifecycle
path for these native definitions. They have not yet been qualified for this
candidate, so runtime startup remains pending. On a host with the supported
container runtime available, maintainers should run the V1-04 host gate from
the repository root:

```sh
./scripts/verify-v1-04-host.sh --output /path/to/evidence
```

The script needs Python 3, a working Docker or Podman CLI, and access to
`nodejs.org` for a checksummed temporary Node runtime. It uses the checked-in,
checksummed Dev Container CLI `0.89.0` archive; npm and an installed Node.js
are not required. It renders both Template variants and records build/up/exec
plus home-volume persistence evidence in the chosen
output directory. It also checks the custom-user and bind-home examples when
the host supports them. Keep the output directory, including `evidence.json`
and step logs, and provide it with the implementation review. Do not claim
lifecycle support until that evidence records passing results for the tested
host and runtime. The aibox Go preview cannot perform those lifecycle actions.

## Keep private workspace state persistent

The standard starter uses a project-scoped named volume at `/home/aibox`:

```jsonc
"remoteUser": "aibox",
"mounts": [
  "source=aibox-home-${devcontainerId},target=/home/aibox,type=volume"
]
```

The volume keeps the container user's home across rebuilds. It does not mount
the host home or place login state in the project. Keep credentials and private
provider configuration out of committed project files.

For a reviewed project-local bind mount, the bind-home example targets
`.aibox-home` and includes it in the project's ignore rules. The source must
exist on the machine that runs the container daemon; create it explicitly
before starting the environment. Bind mounts have host ownership and cleanup
behavior that differ from named volumes. Do not point the mount at the entire
host home.

## Change the container user

The `custom-user` example aligns `containerUser`, `remoteUser`, the Dockerfile
home, and the persistent mount target at `/home/dev`. When changing the user,
update all four together and confirm UID/GID, home ownership, and mount
writability on the daemon host. A mount can hide files baked into the image, so
required defaults must be installed outside the mounted home or initialized
idempotently.

## Configuration ownership

- `devcontainer.json`, Dockerfile, Compose, mounts, ports and lifecycle hooks
  remain native project inputs interpreted by the upstream Dev Container CLI.
- Feature presence and each Feature's declared options select installable
  tools. Omit a Feature to leave that optional tool out.
- `customizations.aibox` is optional UX metadata. The minimal starter does not
  enable it; settings there have no effect without an aibox runtime Feature.
- Native tool configuration belongs to the corresponding tool and remains
  user-editable. Template application is an initial copy, not ongoing
  management of project edits.
- The aibox Go preview and MCP interface currently inspect only. They do not
  apply configuration changes or perform host lifecycle operations.

Until the host gate evidence is available, treat the examples as candidate
native definitions, not as a release-qualified setup.

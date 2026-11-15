# Native workspace starters: requirements

Normative. Requirements use *EARS* patterns. This version restates the V1-04
contract as implemented and accepted at `85cc0c9a`, without changing its
meaning. Its sources are chapter 8 of the v1 product description, sections
"Starter distribution and upstream interoperability" and "Native
configuration example", and chapter 12, section "Named environments and
persistence". Where wording differs, the chapter text is the meaning this
version keeps.

Names in `code` refer to `interfaces.json` unless a schema is named.
"Starter" means a Template in `templates.ids`; "example" means a maintained
example in `examples.ids`.
Verification levels: **J** = a journey in `journeys.md` (specification check
or upstream host run); **C** = component test only.

## Templates

| ID | Requirement | Level |
|---|---|---|
| NS-001 | aibox shall ship two versioned Dev Container Templates, `minimal` and `curated`, in the upstream Template format, each with `id` its name, `version` `1.0.0`, `publisher` `projectious-work`, platform `Any`, and non-empty `name`, `description`, `documentationURL` and `licenseURL`. | J |
| NS-002 | Each starter shall declare the string option `harness` with `enum` `codex`, `none` and `default` `codex`, and shall pass it to the image build as the build argument `HARNESS`. | J |
| NS-003 | Each starter shall build Debian with bash, git, tmux and Yazi, from a base image pinned by digest and APT sources pinned to one Debian snapshot; with `harness` `codex` it shall also install the pinned Codex CLI, and with `none` no harness. | J |
| NS-004 | The `curated` starter shall add `/etc/tmux.conf`, and a user's `~/.tmux.conf` shall be loaded after it. | J |
| NS-005 | A starter shall install tools only; it shall start no harness and transfer no credentials. | J |

## Project files and ownership

| ID | Requirement | Level |
|---|---|---|
| NS-010 | A starter shall create `.devcontainer/devcontainer.json`, `.devcontainer/Dockerfile`, `.devcontainer/install-tools.sh` and `.devcontainer/.gitignore`, and `devcontainer.json` shall be the sole Dev Container entry point. | J |
| NS-011 | The starter's `.devcontainer/.gitignore` shall ignore `*.env` and `compose.local.yaml`. | J |
| NS-012 | A starter shall not create `.aibox/` or `.aibox-home/`, and shall not copy or mount a host home directory. | J |
| NS-013 | When the local renderer `templates/render.mjs` applies a starter, it shall substitute the Template options and shall refuse a destination that already has a `.devcontainer` configuration. | J |
| NS-014 | After a starter is applied, its files shall belong to the project; aibox shall not re-apply or overwrite them. | J |
| NS-015 | aibox shall not introduce a Feature lockfile beside the Dev Container CLI's `.devcontainer/devcontainer-lock.json`. | J |

## User and persistent home

| ID | Requirement | Level |
|---|---|---|
| NS-020 | The default image shall create the non-root user `aibox` with home `/home/aibox`, and the starter shall declare `containerUser` and `remoteUser` `aibox`. | J |
| NS-021 | The starter shall mount the project-scoped named volume `source=aibox-home-${devcontainerId},target=/home/aibox,type=volume` as the persistent home. | J |
| NS-022 | When a container built from a starter, or from the example `custom-user` or `bind-home`, is removed and rebuilt with the pinned Dev Container CLI, a file written to the user's home before the rebuild shall still exist and the home shall be writable. | J |
| NS-023 | Where a user replaces the named volume with a bind mount, the example `bind-home` shall show the mount `{"source": "${localWorkspaceFolder}/.aibox-home", "target": "/home/aibox", "type": "bind"}` and a project `.gitignore` that ignores `.aibox-home/` and `.env`. | J |
| NS-024 | Where a custom user is used, the example `custom-user` shall align the Dockerfile `ARG USERNAME=dev`, the created user, `containerUser`, `remoteUser` `dev` and the home mount target `/home/dev`. | J |

## Upstream interoperability

| ID | Requirement | Level |
|---|---|---|
| NS-030 | Each starter and the examples `custom-user` and `bind-home` shall pass `devcontainer read-configuration`, `build`, `up` and `exec` with the pinned Dev Container CLI (`devcontainerCli.version`) and without the aibox executable. | J |
| NS-031 | A starter shall remain a valid Dev Container when no aibox runtime Feature is present. | J |

## Examples and the customization schema

| ID | Requirement | Level |
|---|---|---|
| NS-040 | aibox shall maintain the examples `minimal`, `customized`, `custom-user` and `bind-home`, each with `build` `{"dockerfile": "Dockerfile", "context": "."}`, no `image`, equal `containerUser` and `remoteUser`, a `Dockerfile`, an `install-tools.sh`, and a `.devcontainer/.gitignore` that ignores `*.env` and `compose.local.yaml`. | J |
| NS-041 | Each example shall show every option of `schemas/customization.schema.json` in place as a comment with its purpose and valid values, generated from the schema, and shall keep only its selected native fields active; comments shall not change the resolved native definition. | J |
| NS-042 | The `minimal` example shall have no active `customizations`; the `customized` example's `customizations.aibox` shall validate against `schemas/customization.schema.json`. | J |
| NS-043 | `schemas/customization.schema.json` shall be closed: an unknown key, an unknown `schemaVersion` or a `null` value shall fail validation. | J |
| NS-044 | An example shall not invent a Feature registry reference; a Feature reference in an example shall be an immutable digest with its tested options. | J |
| NS-045 | `aibox inspect --context local` and the MCP tool `inspect_workspace` shall inspect each example read-only and give equivalent results. | J |

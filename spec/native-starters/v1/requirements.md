# Native workspace starters: requirements

Normative. Requirements use *EARS* patterns.

## Precedence and reading rules

- `interfaces.json` and `schemas/customization.schema.json` prevail over this
  file and over `journeys.md`. Where they differ, `interfaces.json` and the
  schema win, and the difference is a specification defect that is fixed in
  `native-starters@v2`.
- The sources cited in `overview.md` (chapters 8 and 12 of the v1 product
  description and `templates/README.md`) are provenance only. They do not add
  to or override this version.
- "Starter" means a Template in `interfaces.json` `templates.ids`; "example"
  means a maintained example in `examples.ids`.

**Verified by** names the journey of `journeys.md` that verifies the
requirement. "Not separately verified" marks an obligation that is part of
the contract but that no accepted V1-04 check exercises on its own.

## Templates

| ID | Requirement | Verified by |
|---|---|---|
| NS-001 | aibox shall ship two individually versioned Dev Container Templates, `minimal` and `curated`, in the upstream Template format, each with `id` its name, `version` `1.0.0`, `publisher` `projectious-work`, platform `Any`, and non-empty `name`, `description`, `documentationURL` and `licenseURL`. | J-NS-01 |
| NS-002 | Each starter shall declare the string option `harness` with `enum` `codex`, `none` and `default` `codex`, and shall pass it to the image build as the build argument `HARNESS`. | J-NS-01, J-NS-05 |
| NS-003 | Each starter shall build Debian with bash, git, tmux and Yazi from a base image pinned by digest and APT sources pinned to one Debian snapshot, and shall verify the Yazi and Codex downloads by per-architecture SHA-256; with `harness` `codex` it shall install the pinned Codex CLI, and with `none` no harness. | J-NS-01 (pins and hashes), J-NS-06 (tmux, Yazi and Codex with `codex`); bash and git, and `none` on a host, not separately verified |
| NS-004 | The `curated` starter shall add `/etc/tmux.conf`, and tmux shall load a user's `~/.tmux.conf` after it, so that user overrides stay authoritative when the home volume masks image files. | J-NS-01 (file present); load order not separately verified |
| NS-005 | A starter shall install tools only; it shall start no harness and transfer no credentials. | Not separately verified |
| NS-006 | The image of a starter shall run as `USER ${USERNAME}` with `ARG USERNAME=aibox`, and its Dockerfile shall not `COPY` or `ADD` `.aibox` or `.aibox-home` content. | J-NS-01 |

## Project files and ownership

| ID | Requirement | Verified by |
|---|---|---|
| NS-010 | A starter shall provide `.devcontainer/devcontainer.json`, `.devcontainer/Dockerfile`, `.devcontainer/install-tools.sh` and `.devcontainer/.gitignore`, and `devcontainer.json` shall be the sole Dev Container entry point, with `build` `{"dockerfile": "Dockerfile", "context": "."}`. | J-NS-01 |
| NS-011 | The `.devcontainer/.gitignore` of each starter and each example shall ignore `*.env` and `compose.local.yaml`. | J-NS-01, J-NS-02 |
| NS-012 | A starter shall not create `.aibox/` or `.aibox-home/`, and shall not copy or mount a whole host home implicitly. | J-NS-01 (no such paths in the Template sources), J-NS-07 (home is a named volume) |
| NS-013 | When `templates/render.mjs` applies a starter, it shall substitute the Template options, and it shall refuse a destination that already has a `.devcontainer` configuration. | J-NS-05 (substitution); refusal not separately verified |
| NS-014 | After a starter is applied, its files shall belong to the project, and aibox shall not re-apply a starter or overwrite user edits. | Not separately verified |
| NS-015 | aibox shall not introduce a Feature lockfile beside the Dev Container CLI's own lockfile. | Not separately verified |

## User and persistent home

| ID | Requirement | Verified by |
|---|---|---|
| NS-020 | The default image shall create the non-root user `aibox` with home `/home/aibox`, and a starter shall declare `containerUser` and `remoteUser` `aibox`. | J-NS-01, J-NS-05, J-NS-06 |
| NS-021 | A starter shall mount exactly one home mount, the project-scoped named volume `source=aibox-home-${devcontainerId},target=/home/aibox,type=volume`. | J-NS-01, J-NS-07 |
| NS-022 | When a container built from a starter, or from the example `custom-user` or `bind-home`, is removed and brought up again with the pinned Dev Container CLI, the home mount shall keep its identity (the same named volume, or the same bind source), a file written to the user's home before shall still exist, and the home shall be writable. | J-NS-07 |
| NS-023 | The example `bind-home` shall replace the named volume with the mount `{"source": "${localWorkspaceFolder}/.aibox-home", "target": "/home/aibox", "type": "bind"}`, and its project `.gitignore` shall ignore `.aibox-home/`, `.env` and `.devcontainer/compose.local.yaml`. | J-NS-02 |
| NS-024 | The example `custom-user` shall align the Dockerfile `ARG USERNAME=dev`, a user created with `useradd` from it, `containerUser` and `remoteUser` `dev`, and the home mount target `/home/dev`. | J-NS-02, J-NS-06 |

## Upstream interoperability

| ID | Requirement | Verified by |
|---|---|---|
| NS-030 | Each starter, and the examples `custom-user` and `bind-home`, shall pass `devcontainer read-configuration`, `build`, `up` and `exec` with the pinned Dev Container CLI (`interfaces.json` `devcontainerCli.version`) without the aibox executable. | J-NS-06 |
| NS-031 | A starter shall remain a valid Dev Container when no aibox runtime Feature is present. | J-NS-06 (no starter selects one) |

## Examples and the customization schema

| ID | Requirement | Verified by |
|---|---|---|
| NS-040 | aibox shall maintain the examples `minimal`, `customized`, `custom-user` and `bind-home`, each with `build` `{"dockerfile": "Dockerfile", "context": "."}`, no `image` and no `dockerComposeFile`, explicit and equal `containerUser` and `remoteUser`, exactly one home mount, a `Dockerfile` and an `install-tools.sh`, and no top-level `aibox` member. | J-NS-02 |
| NS-041 | Each example shall show every option of `schemas/customization.schema.json` in place as a comment with its purpose and valid values, generated from the schema, and shall keep only its selected native fields active; inactive options shall stay commented, and comments shall not change the resolved native definition. | J-NS-03 |
| NS-042 | The `minimal` example shall have no active `customizations` and no `features`; the `customized` example shall have `customizations.aibox` that validates against the schema. | J-NS-02 |
| NS-043 | The schema shall be closed: an unknown key, an unknown `schemaVersion` or a `null` value shall fail validation. | J-NS-04 |
| NS-044 | An example shall not invent a Feature registry reference; a published reference shall include its tested ID, digest and actual options in the example and in the selection table together. The `customized` example's `features` shall equal `selected` in `tool-catalog@v1` `fixtures/v1-05-native-features.json`. | J-NS-02 (equality with the fixture) |
| NS-045 | `aibox inspect --context local` and the MCP tool `inspect_workspace` shall inspect each example read-only and give equivalent results. | J-NS-09 |

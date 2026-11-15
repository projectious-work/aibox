# Native workspace starters: acceptance journeys

Normative. Two kinds of journey realize this version:

- **Specification checks** run offline from a clean checkout:
  `scripts/spec/validate-contracts.py`, `scripts/spec/validate-examples.mjs`,
  `scripts/spec/render-devcontainer-examples.mjs --check` and
  `scripts/spec/validate.mjs`, all part of `scripts/verify-v1-04.sh`.
- **Host runs** use `scripts/verify-v1-04-host.py` on a host with a container
  engine. It installs the pinned Dev Container CLI from its checked-in,
  checksummed archive, renders each starter into a fresh directory with
  `templates/render.mjs`, and records every step in an evidence manifest.
  aibox is not on `PATH`.

## Journeys

| ID | Kind | Covers | Steps | Expected |
|---|---|---|---|---|
| J-NS-01 | Specification check | NS-001, NS-002, NS-003, NS-004, NS-010, NS-011, NS-020, NS-021 | Read both Template sources. | Metadata, option, build argument, user, mount, files, `.gitignore` entries, pinned base digest and Debian snapshot as the requirements state; `curated` has `tmux.conf`. A metadata copy whose `harness` default is outside the `enum` is rejected. |
| J-NS-02 | Specification check | NS-040, NS-023, NS-024, NS-042 | Read the four examples. | Each meets NS-040; `bind-home` and `custom-user` meet NS-023 and NS-024; `minimal` has no `customizations`; `customized` validates against the schema. |
| J-NS-03 | Specification check | NS-041 | Regenerate the commented option catalogs with `render-devcontainer-examples.mjs --check`. | No difference to the maintained examples. |
| J-NS-04 | Specification check | NS-043 | Validate `{"schemaVersion":"2"}`, a `null` value and invalid nested values against the schema. | Each is rejected. |
| J-NS-05 | Specification check | NS-002, NS-013, NS-012, NS-020 | Render each starter with `harness` `codex` and with `none` into an empty directory; then render into a directory that already has `.devcontainer`. | First: `build.args.HARNESS` equals the chosen value, `containerUser` and `remoteUser` are `aibox`, and no `.aibox/` or `.aibox-home/` exists. Second: refused, nothing changed. |
| J-NS-06 | Host run | NS-030, NS-031, NS-005 | For `minimal`, `curated`, `custom-user` and `bind-home`: `read-configuration`, `build`, `up`, then `exec` of `id`, `tmux -V`, `yazi --version` and, with harness `codex`, `codex --version`. | Every step exits 0. |
| J-NS-07 | Host run | NS-022, NS-021 | For the same four: write a marker file into `$HOME`, remove the container, rebuild with `up`, read the marker. | The marker reads `v1-04-persisted` and `$HOME` is writable. |
| J-NS-08 | Host run | NS-015, NS-014 | After J-NS-06, compare the project tree with the rendered tree. | Only the Dev Container CLI's `.devcontainer/devcontainer-lock.json` may be added. |
| J-NS-09 | CLI and MCP | NS-045 | `aibox inspect --context local --project examples/<name> --format json` and the MCP `inspect_workspace` call for each example. | Equivalent results; no file changes. |

NS-044 is checked by `validate-contracts.py` and `validate-examples.mjs` for
the Feature references the examples select.

## Coverage by the V1-04 evidence

The accepted V1-04 evidence (`evidence/v1/V1-04-85cc0c9a/`) ran J-NS-01 to
J-NS-04, the first part of J-NS-05, J-NS-06, J-NS-07 and J-NS-09. The refusal
step of J-NS-05 and J-NS-08 restate chapter 8 obligations that the V1-04 gates
did not exercise as separate steps; a later run of these journeys closes that
gap without changing this version.


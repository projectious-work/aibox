# Native workspace starters: acceptance journeys

Normative. `interfaces.json` prevails over this file (see
`requirements.md`, "Precedence and reading rules"). Two kinds of journey
realize this version:

- **Specification checks** run offline from a clean checkout:
  `scripts/spec/validate-contracts.py`, `scripts/spec/validate-examples.mjs`,
  `scripts/spec/render-devcontainer-examples.mjs --check` and
  `scripts/spec/validate.mjs`, and the rendering step of
  `scripts/verify-v1-04.sh`.
- **Host runs** use `scripts/verify-v1-04-host.py` on a host with a container
  engine. It installs the pinned Dev Container CLI from its checked-in,
  checksummed archive, renders each starter into a fresh directory with
  `templates/render.mjs` at the default `harness`, and records every step in
  an evidence manifest. aibox is not invoked.

## Journeys

| ID | Kind | Covers | Steps | Expected |
|---|---|---|---|---|
| J-NS-01 | Specification check | NS-001 to NS-004 (as listed), NS-006, NS-010 to NS-012, NS-020, NS-021 | Read both Template sources. | Metadata, option, build argument, `build`, user, single home mount, files, `.gitignore` entries, pinned base digest, Debian snapshot, Yazi and Codex SHA-256 values, `USER ${USERNAME}`, no `COPY`/`ADD` of `.aibox` content and no `.aibox` or `.aibox-home` path, as the requirements state; `curated` has `tmux.conf`. A metadata copy whose `harness` default is outside the `enum` is rejected. |
| J-NS-02 | Specification check | NS-011, NS-023, NS-024, NS-040, NS-042, NS-044 | Read the four examples. | Each meets NS-040; `bind-home` meets NS-023 and `custom-user` NS-024; `minimal` has no `customizations` and no `features`; `customized` validates against the customization schema and its `features` equal the tool-catalog fixture. |
| J-NS-03 | Specification check | NS-041 | Regenerate the commented option catalogs with `render-devcontainer-examples.mjs --check`. | No difference to the maintained examples. |
| J-NS-05 | Specification check | NS-002, NS-013 (substitution), NS-020 | Render each starter with `harness` `codex` and with `none` into an empty directory. | `build.args.HARNESS` equals the chosen value; `containerUser` and `remoteUser` are `aibox`. |
| J-NS-06 | Host run | NS-003 (as listed), NS-020, NS-024, NS-030, NS-031 | For `minimal`, `curated`, `custom-user` and `bind-home`: `read-configuration`, `build`, `up`, then `exec` of a check that the user is `aibox` (`dev` for `custom-user`), `$HOME` is writable, and `tmux -V`, `yazi --version` and `codex --version` succeed. | Every step exits 0. |
| J-NS-07 | Host run | NS-012, NS-021, NS-022 | For the same four: inspect the home mount, write a marker file into `$HOME`, remove the container, bring it up again, inspect the home mount again and read the marker. | The home mount is the expected named volume (bind source for `bind-home`) with the same identity before and after; the marker reads `v1-04-persisted` and `$HOME` is writable. |
| J-NS-09 | CLI and MCP | NS-045 | `aibox inspect --context local --project examples/<name> --format json` and the MCP `inspect_workspace` call for each example. | Equivalent results. |

J-NS-04 and J-NS-08 are unused; the customization schema's own checks are the journeys of `workspace-customization@v1`.

## Coverage by the V1-04 evidence

- The accepted host evidence is two macOS arm64 Docker runs at `85cc0c9a`
  (`evidence/v1/V1-04-85cc0c9a/`), at the default `harness` `codex`;
  `harness` `none` was never built on a host.
- V1-05 later changed both starters to configure `trixie`, `trixie-updates`
  and `trixie-security` from snapshot `20261028T000000Z`. The current
  starters have passed the specification checks (J-NS-01 to J-NS-03 and J-NS-05) and
  J-NS-09, but no host run since that change.
- Requirements marked "not separately verified" in `requirements.md` have no
  journey in this version.

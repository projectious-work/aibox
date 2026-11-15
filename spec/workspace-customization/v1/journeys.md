# Workspace customization schema: acceptance journeys

Normative. The schema and `interfaces.json` prevail over this file (see
`requirements.md`, "Precedence and reading rules"). Every journey is an
offline specification check run from a clean checkout; all are part of
`scripts/verify-v1-04.sh` and `scripts/verify-v1-05.sh`.

| ID | Covers | Steps | Expected |
|---|---|---|---|
| J-WC-01 | WC-010, WC-014, WC-020, WC-021 | `node scripts/spec/inventory.mjs --check`. | The committed schema equals the generator's output from the v0 baseline inventory, including the two enumerations closed after the baseline. |
| J-WC-02 | WC-012 | `node scripts/spec/validate.mjs`. | Every one of the 115 `customization.*` rows of `ledger/v0-parity/configuration.json` resolves to a key under `workspace`. |
| J-WC-03 | WC-010, WC-011 | `uv run --script scripts/spec/validate-contracts.py`. | The schema is a valid Draft 2020-12 schema; `{"schemaVersion":"2"}`, `workspace.theme` `null` and invalid `latex.preview` values are rejected, and the `customized` example's `customizations.aibox` is accepted. |
| J-WC-04 | WC-004, WC-011 | `node scripts/spec/validate-examples.mjs` and `node scripts/spec/validate.mjs`. | Unknown keys, `null` and out-of-range values are rejected; the `minimal` example has no `customizations`. |
| J-WC-05 | WC-022 | `node scripts/spec/render-devcontainer-examples.mjs --check`. | The examples' commented option catalogs equal what the schema generates. |

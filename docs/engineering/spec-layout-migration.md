# Specification layout migration

Issue #504 adopts the specification layout `spec/<slug>/vN/` from the company
standards. The first step moved everything that is not a contract out of
`spec/v1/`. Frozen specifications keep their original paths as historical
references; read an old path through the table below. Recorded evidence keeps
the commands and paths of the run it records.

| Old path | New path | Kind |
|---|---|---|
| `spec/v1/01-…21-*.md`, `README.md`, `feature-trace.md`, `command-disposition.md`, `phase-id-crosswalk.md` | `docs/product/v1/` | Product documentation |
| `spec/v1/roadmap.yaml`, `spec/v1/roadmap.schema.json` | `roadmap.yaml`, `roadmap.schema.json` | Roadmap |
| `spec/v1/ledger/` | `ledger/v0-parity/` | Generated v0 inventory |
| `spec/v1/evidence/` | `evidence/v1/` | Evidence |
| `spec/v1/baseline-evidence.md`, `spec/v1/docs-foundation-evidence.md` | `evidence/v1/` | Evidence |
| `spec/v1/scripts/` | `scripts/spec/` | Tooling |
| `spec/v1/examples/` | `examples/` | Product examples |
| `spec/v1/catalog/v1-05-features.json` | `spec/tool-catalog/v1/catalog/v1-05-features.json` | Contract: `tool-catalog@v1` |
| `spec/v1/feature-catalog.schema.json` | `spec/tool-catalog/v1/schemas/feature-catalog.schema.json` | Contract: `tool-catalog@v1` |
| `spec/v1/fixtures/v1-05-*` | `spec/tool-catalog/v1/fixtures/v1-05-*` | Contract: `tool-catalog@v1` |
| `spec/v1/customization.schema.json` | `spec/workspace-customization/v1/schemas/customization.schema.json` | Contract: `workspace-customization@v1` |

Paths in `roadmap.yaml` are relative to the repository root. The validators in
`scripts/spec/` read every location from `scripts/spec/paths.mjs` and
`scripts/spec/spec_paths.py`.

Contract files move out of `spec/v1/` one specification at a time, each as a
draft that the owner freezes:

- `tool-catalog@v1` (V1-05): the catalog, its schema and the `v1-05-*`
  fixtures. The 135 generated Feature sources under `features/src/` and
  `docs-site/content/docs/tool-catalog.md` name the new catalog path in their
  header comment, so published Feature tarballs differ in those bytes from
  the ones accepted at `651506fa`; no install logic, option or version
  changes.
- `workspace-customization@v1`: `customization.schema.json`, the closed
  `customizations.aibox` schema. Its `$id` is a GitHub blob URL of the file,
  so it follows the move; `tool-catalog@v1` kept its schema's
  `https://projectious.work/...` `$id`, which names no path.
- `native-starters@v1` (V1-04): no contract file; it cites
  `workspace-customization@v1`.
- `operation-core@v1` (V1-03): the operation request, result, record and
  operator-policy schemas; it waits for V1-06, which changes the result
  schema.

Until `spec/v1/` is empty, `python3 .standards/standards.py check-specs`
fails on `spec/v1` with `unexpected entries`. With `--base`, it also reports
every file moved out of `spec/v1/` as a deleted specification file, because
it compares without rename detection. `spec/v1/` was never a versioned
feature specification, so these moves delete no frozen version; the
exception is recorded in
[`standards-applicability.md`](standards-applicability.md).

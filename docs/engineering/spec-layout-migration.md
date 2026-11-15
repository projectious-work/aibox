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
| `spec/v1/customization.schema.json` | `spec/native-starters/v1/schemas/customization.schema.json` | Contract: `native-starters@v1` |

Paths in `roadmap.yaml` are relative to the repository root. The validators in
`scripts/spec/` read every location from `scripts/spec/paths.mjs` and
`scripts/spec/spec_paths.py`.

Still in `spec/v1/` until the next step: the schemas, `catalog/` and
`fixtures/`. They move into the `operation-core`, `native-starters` and
`tool-catalog` specifications, each frozen with owner approval; this table
gains their rows then. Until then
`python3 .standards/standards.py check-specs` fails on `spec/v1`.

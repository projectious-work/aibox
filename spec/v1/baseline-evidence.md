# V1-01 baseline evidence

Recorded 2026-09-27 for the v1 specification phase. This snapshot distinguishes
the immutable v0 source baseline, the v1 specification change, and the PR #463
merge commit on which this uncommitted V1-01 working tree is based.

## Commits

| Purpose | Full commit |
|---|---|
| v0.35.0 source baseline used by the inventory fixture | `9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6` |
| v1 branch point / first parent of PR #463 merge | `d9a8e22b4ba4385f1dc0134f58f5405417116d85` |
| PR #463 head containing the specification change | `93df482aeb754d381878b44aac0ec11a37b08c37` |
| PR #463 merge and V1-01 worktree base (`HEAD` before these edits) | `23e98e7ba6ca00712bbdacd01f38d506732952cb` |

The census is pinned to the v0.35.0 source commit. Validation and the Hugo build
were run in a worktree based on the PR #463 merge commit, including the V1-01
working-tree edits where applicable. The separate first-parent commit identifies
the v1 branch point; it is not the source inventory baseline.

## Baseline fixture and source census

[`fixtures/baseline/v0.35.0.json`](fixtures/baseline/v0.35.0.json) records the
source commit, audited source scopes and expected counts. These are counts of
declarations and files in those source scopes; they are not runtime measurements
and do not establish behavioral parity or prove that v1 implements any row.
The behavioral contract and parity evidence remain separate requirements.

| Inventory | Count |
|---|---:|
| Configuration fields | 265 |
| Addon recipes / declared tools | 41 / 98 |
| CLI declarations / actions / arguments | 14 / 66 / 74 |
| Runtime assets | 105 |
| Documentation pages | 52 |
| Environment identifiers | 83 |

The source link, source-hash and cardinality checks in `scripts/baseline.mjs`
resolve against the pinned commit, then confirm that the fixture, generated
ledgers and census agree. `scripts/inventory.mjs --check` verifies deterministic
inventory output. `scripts/contracts.mjs --check` checks the operation, policy
and operation record contract outputs. `scripts/validate.mjs` checks cross-document
requirements, acceptance groups, roadmap references and local Markdown links.
`scripts/validate-contracts.py` runs JSON Schema Draft 2020-12 checks and the
positive/negative contract fixtures.

The static examples under `examples/` are syntax and schema fixtures. The
minimal example has no aibox-specific customization; the customized example
validates against `customization.schema.json`, including invalid-version,
unknown-field, null and invalid-theme cases. Both currently use a generic
upstream image and user. They are not claims about the future aibox base image,
V1-04 starter projects, image build/up/exec behavior or qualified runtime
parity. Runtime and published-artifact qualification remain phase gates.

Run the repository-root checks (the baseline commit must be available in Git
history):

```sh
node spec/v1/scripts/baseline.mjs
node spec/v1/scripts/inventory.mjs --check
node spec/v1/scripts/contracts.mjs --check
node spec/v1/scripts/validate.mjs
node spec/v1/scripts/validate-examples.mjs
uv run --script spec/v1/scripts/validate-contracts.py
```

At the PR #463 merge baseline, the original inventory and contract checks
passed: 27 requirements, 24 acceptance groups and 23 roadmap phases were
consistent, and Draft 2020-12 validation covered 8 schemas and 95 fixtures.
The V1-01 candidate's complete command list above passed as well. Its baseline
resolver verified every cited source link and hash against the pinned source;
the static example validator passed; and Draft 2020-12 validation covered 8
schemas and 98 fixtures. None of these checks is a product runtime test.

| Check | Expected result | Observed result | Status |
|---|---|---|---|
| `node spec/v1/scripts/baseline.mjs` | Fixture, source provenance, hashes and generated ledger match pinned v0 commit | 265 fields, 41 recipes / 98 tools, 66 actions / 74 arguments, 105 assets, 52 docs and 83 environment IDs matched | Pass |
| `node spec/v1/scripts/inventory.mjs --check` | Ledger regeneration makes no changes | Same census as pinned baseline | Pass |
| `node spec/v1/scripts/contracts.mjs --check` | Contract generation makes no changes | 15 operations, policy and operation records checked | Pass |
| `node spec/v1/scripts/validate.mjs` | Normative trace, roadmap and links consistent | 27 requirements, 24 acceptance groups, 23 phases | Pass |
| `node spec/v1/scripts/validate-examples.mjs` | Static native examples and negative cases pass | Static checks passed; V1-04 runtime qualification deferred | Pass |
| `uv run --script spec/v1/scripts/validate-contracts.py` | Full Draft 2020-12 schema fixtures pass | 8 schemas, 98 fixtures | Pass |
| `./scripts/build-docs.sh --destination /tmp/aibox-v1-pr463-docs` | Existing v1 Docsy site builds as a baseline | 161 pages and 139 static files; two Hugo deprecation warnings | Pass with warnings |

The reproducible outputs are the checked-in fixture, ledgers, schemas and native
examples. `/tmp/aibox-v1-pr463-docs` is a local build artifact, not a published
site or immutable release artifact. G08 reconciliation of effective v0 defaults,
aliases and behavior remains open for the affected conversion phases; these
source-count checks cannot close it.

## Dependency research snapshots

The following source commits were inspected on 2026-09-25. They are research
anchors, not qualification results or release pins:

- Dev Container CLI: `5dc7533314b5ba7ec3875c30143dfe1aec644870`.
- Official Go SDK: `e07f0c9d5abf509ac1e47abf27cfa539eeda64a5`.
- tmux-agent-status: `546b6ca51c75b415db0b3ce06910703f875d3aae`.
- gh-dash: `1b14dd961ea47e2ad53eae5fed7fc942219bcc01`.

The inspected Dev Container CLI README distinguished build/up/exec from
stop/down. Release qualification must inspect the selected pin again.
[CLI reference source](https://github.com/devcontainers/cli/blob/5dc7533314b5ba7ec3875c30143dfe1aec644870/README.md)

The inspected Go SDK exposed server/client primitives, including stdio.
[SDK reference source](https://github.com/modelcontextprotocol/go-sdk/blob/e07f0c9d5abf509ac1e47abf27cfa539eeda64a5/README.md)

## Existing local Docsy build

The existing v1 docs build succeeded locally on 2026-09-27 from the PR #463
merge checkout using:

```sh
./scripts/build-docs.sh --destination /tmp/aibox-v1-pr463-docs
```

The build used Hugo `v0.165.0-76a5e1880ab46688155b02e99bab9be2a6134492+extended`
and the initialized Docsy submodule at `01c827ea890e8e498f6046a7666a3031f318cc7f`.
It generated 161 pages and 139 static files. Hugo emitted two template API
deprecation warnings (`.Language.LanguageDirection` and `.Site.AllPages`); the
build still exited successfully. Initial `npm ci` also reported one high
severity advisory in its dependency audit output; this was not a Hugo build
failure and no dependency changes were made for this evidence run.

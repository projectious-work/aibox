#!/usr/bin/env bash
# Verify V1-05's offline gate: Go format/vet/tests and CLI/MCP inspection of
# the selected Features, the tool catalog and its generated Features, the
# Feature component tests, the AC-TOOLS harness unit tests, the specification
# checks and the v1 Hugo build. Container acceptance runs separately through
# tests/acceptance/v1-05/run-host.sh on each target.
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"
if [[ -n "$(git status --porcelain --untracked-files=all)" ]]; then
  printf 'V1-05 candidate worktree must be clean before evidence capture.\n' >&2
  exit 2
fi
evidence_dir="$(mktemp -d "${TMPDIR:-/tmp}/aibox-v1-05-offline.XXXXXX")"
candidate_sha="$(git rev-parse HEAD)"
started_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

unformatted="$(gofmt -l cmd internal)"
if [[ -n "$unformatted" ]]; then
  printf 'gofmt lists files:\n%s\n' "$unformatted" >&2
  exit 1
fi
go vet ./...
go test ./...
go test -race ./...
go mod verify
go test -run '^TestV104StarterCLIMCPParity$' -count=1 -v ./cmd/aibox > "$evidence_dir/mcp-parity.txt"
go build -o "$evidence_dir/aibox" ./cmd/aibox
GOOS=darwin GOARCH=arm64 go build -o "$evidence_dir/aibox-darwin-arm64" ./cmd/aibox

# Roadmap demo: the curated tool selection is visible through inspection.
"$evidence_dir/aibox" inspect --context local --project spec/v1/examples/customized --format json > "$evidence_dir/inspect-customized.json"
"$evidence_dir/aibox" inspect --context local --project spec/v1/examples/minimal --format json > "$evidence_dir/inspect-minimal.json"

node spec/v1/scripts/validate.mjs
node spec/v1/scripts/validate-examples.mjs
uv run --script spec/v1/scripts/validate-contracts.py
python3 spec/v1/scripts/validate-v1-05-catalog.py
node spec/v1/scripts/render-devcontainer-examples.mjs --check
scripts/v1-05-features.sh check --complete
node --test 'features/test/*.test.mjs' > "$evidence_dir/feature-component-tests.txt"
(cd tests/acceptance/v1-05 && python3 -m unittest discover -s unit -t .) 2> "$evidence_dir/acceptance-unit-tests.txt"

./scripts/build-docs.sh --destination "$evidence_dir/hugo" > "$evidence_dir/hugo-build.txt"
git diff --check

# v0 stays the current line: this branch changes no v0 source.
v1_base="$(git merge-base HEAD origin/v1.x-dev)"
v0_changes="$(git diff --name-only "$v1_base" HEAD -- cli images)"
if [[ -n "$v0_changes" ]]; then
  printf 'V1-05 changes v0 sources:\n%s\n' "$v0_changes" >&2
  exit 1
fi

finished_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
node - "$evidence_dir" "$candidate_sha" "$started_at" "$finished_at" "$v1_base" <<'NODE'
const {createHash} = require('node:crypto');
const {readFileSync, readdirSync, writeFileSync} = require('node:fs');
const {join} = require('node:path');
const os = require('node:os');
const [directory, sourceCommit, startedAt, finishedAt, v1Base] = process.argv.slice(2);
const digest = name => 'sha256:' + createHash('sha256').update(readFileSync(join(directory, name))).digest('hex');
const text = name => readFileSync(join(directory, name), 'utf8');
const count = (name, label) => Number((text(name).match(new RegExp(`^ℹ ${label} (\\d+)$`, 'm')) || [])[1]);
const pages = Number((text('hugo-build.txt').match(/Pages\s*\│?\s*(\d+)/) || [])[1]);
const evidence = {
  phase: 'V1-05', verification: 'offline', sourceCommit, startedAt, finishedAt,
  actor: process.env.AIBOX_VERIFY_ACTOR || 'local-verifier',
  host: {platform: process.platform, arch: process.arch, kernel: os.release()},
  go: {nativeTarget: `${process.platform}/${process.arch}`, binary: 'aibox', binaryDigest: digest('aibox'),
    darwinArm64Binary: 'aibox-darwin-arm64', darwinArm64Digest: digest('aibox-darwin-arm64')},
  demo: {command: 'aibox inspect --context local --project spec/v1/examples/customized --format json',
    result: 'inspect-customized.json', resultDigest: digest('inspect-customized.json'),
    mcp: {tool: 'inspect_workspace', test: 'TestV104StarterCLIMCPParity', transcript: 'mcp-parity.txt',
      transcriptDigest: digest('mcp-parity.txt')}},
  featureComponentTests: {tests: count('feature-component-tests.txt', 'tests'),
    failed: count('feature-component-tests.txt', 'fail'), skipped: count('feature-component-tests.txt', 'skipped')},
  documentation: {command: './scripts/build-docs.sh', pages, result: 'passed',
    paths: ['docs-site/content/docs/tool-catalog.md', 'docs-site/content/docs/architecture/'],
    audiences: {'tool-catalog.md': 'users', 'architecture/': 'maintainers'}},
  v0: {current: true, check: `no change under cli/ or images/ since ${v1Base}`},
  checks: ['gofmt -l', 'go vet ./...', 'go test ./...', 'go test -race ./...', 'go mod verify',
    'CLI/MCP parity on the native starter examples', 'Linux native and macOS arm64 cross build',
    'specification', 'rendered examples', 'Draft 2020-12 fixtures', 'V1-05 catalog',
    'generated example freshness', 'generated Features match the catalog', 'Feature component tests',
    'AC-TOOLS harness unit tests', 'v1 Hugo build', 'git diff --check', 'no v0 source change'],
  containerAcceptance: {status: 'separate', runner: 'tests/acceptance/v1-05/run-host.sh'},
  artifacts: readdirSync(directory).filter(name => name !== 'evidence.json' && name !== 'hugo'),
  failed: [], skipped: [],
};
if (!(pages > 0) || evidence.featureComponentTests.failed !== 0 || evidence.featureComponentTests.skipped !== 0) {
  throw new Error('V1-05 offline evidence is incomplete: ' + JSON.stringify({pages, ...evidence.featureComponentTests}));
}
writeFileSync(join(directory, 'evidence.json'), JSON.stringify(evidence, null, 2) + '\n');
NODE

printf 'V1-05 offline candidate SHA: %s\n' "$candidate_sha"
printf 'Evidence directory: %s\n' "$evidence_dir"

#!/usr/bin/env bash
# Verify the exact source-built V1-03 read-only customer preview. Outputs are
# kept in a unique temporary directory for a reviewer to inspect or archive.
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"
if [[ -n "$(git status --porcelain --untracked-files=all)" ]]; then
  printf 'V1-03 candidate worktree must be clean before evidence capture.\n' >&2
  exit 2
fi
evidence_dir="$(mktemp -d /tmp/aibox-v1-03.XXXXXX)"
candidate_sha="$(git rev-parse HEAD)"
started_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
native_target="$(go env GOOS)/$(go env GOARCH)"

go test ./...
go test -race ./...
go vet ./...
go mod verify
go test -run '^TestPreviewBinaryMCPParity$' -count=1 -v ./cmd/aibox > "$evidence_dir/mcp-parity.txt"
go build -o "$evidence_dir/aibox" ./cmd/aibox
GOOS=darwin GOARCH=arm64 go build -o "$evidence_dir/aibox-darwin-arm64" ./cmd/aibox

"$evidence_dir/aibox" --version > "$evidence_dir/version.txt"
"$evidence_dir/aibox" --help > "$evidence_dir/help.txt"
"$evidence_dir/aibox" inspect --context local --project spec/v1/examples/minimal --format json > "$evidence_dir/inspect.json"
set +e
"$evidence_dir/aibox" inspect --context local --project "$evidence_dir/missing" --format json > "$evidence_dir/invalid-project.json" 2> "$evidence_dir/invalid-project.stderr"
invalid_exit=$?
set -e
[[ "$invalid_exit" == 2 && ! -s "$evidence_dir/invalid-project.stderr" ]] || {
  printf 'Invalid-project journey must return exit 2 and keep stderr empty.\n' >&2
  exit 1
}
uv run --with 'jsonschema[format]==4.23.0' python -c '
import json, sys
from jsonschema import Draft202012Validator, FormatChecker
schema = json.load(open("spec/v1/operation-result.schema.json", encoding="utf-8"))
def load(path):
    with open(path, encoding="utf-8") as source:
        result = json.load(source)
    Draft202012Validator(schema, format_checker=FormatChecker()).validate(result)
    return result
result = load(sys.argv[1])
assert result["operation"] == "inspect_workspace"
assert result["outcome"] == "succeeded"
assert result["data"]["state"] == "unknown"
invalid = load(sys.argv[2])
assert invalid["outcome"] == "failed"
assert invalid["error"]["code"] == "invalid_input"
assert "target" not in invalid
' "$evidence_dir/inspect.json" "$evidence_dir/invalid-project.json"

node spec/v1/scripts/inventory.mjs --check
node spec/v1/scripts/validate.mjs
uv run --script spec/v1/scripts/validate-contracts.py
./scripts/build-docs.sh --destination "$evidence_dir/hugo"
git diff --check

finished_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
node - "$evidence_dir" "$candidate_sha" "$started_at" "$finished_at" "$native_target" <<'NODE'
const {createHash} = require('node:crypto');
const {readFileSync, writeFileSync} = require('node:fs');
const {join} = require('node:path');
const [directory, sourceCommit, startedAt, finishedAt, nativeTarget] = process.argv.slice(2);
const digest = name => 'sha256:' + createHash('sha256').update(readFileSync(join(directory, name))).digest('hex');
const demo = JSON.parse(readFileSync(join(directory, 'inspect.json'), 'utf8'));
const evidence = {
  phase: 'V1-03', sourceCommit, actor: 'local-verifier', startedAt, finishedAt,
  fixture: 'spec/v1/examples/minimal',
  demo: {command: 'aibox inspect --context local --project spec/v1/examples/minimal --format json',
    operation: demo.operation, outcome: demo.outcome, state: demo.data.state,
    transcript: 'inspect.json', transcriptDigest: digest('inspect.json')},
  mcp: {command: 'aibox mcp serve --context local', tool: 'inspect_workspace',
    transcript: 'mcp-parity.txt', transcriptDigest: digest('mcp-parity.txt')},
  binaries: [
    {target: nativeTarget, path: 'aibox', digest: digest('aibox')},
    {target: 'darwin/arm64', path: 'aibox-darwin-arm64', digest: digest('aibox-darwin-arm64')},
  ],
  documentation: {source: 'docs-site/content/docs/core-status.md', build: 'hugo', output: 'hugo'},
  checks: ['go test ./...', 'go test -race ./...', 'go vet ./...', 'go mod verify',
    'binary stdio MCP client parity and refusal',
    'binary success and invalid-input journeys', 'Draft 2020-12 result validation',
    'inventory', 'specification', 'contract fixtures', 'v1 Hugo build', 'git diff --check'],
  failed: [], skipped: [],
};
writeFileSync(join(directory, 'evidence.json'), JSON.stringify(evidence, null, 2) + '\n');
NODE

printf 'V1-03 source SHA: %s\n' "$candidate_sha"
printf 'Linux binary SHA-256: '
sha256sum "$evidence_dir/aibox"
printf 'macOS arm64 binary SHA-256: '
sha256sum "$evidence_dir/aibox-darwin-arm64"
printf 'Evidence directory: %s\n' "$evidence_dir"

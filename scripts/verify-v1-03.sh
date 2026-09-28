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

go test ./...
go test -race ./...
go vet ./...
go build -o "$evidence_dir/aibox" ./cmd/aibox
GOOS=darwin GOARCH=arm64 go build -o "$evidence_dir/aibox-darwin-arm64" ./cmd/aibox

"$evidence_dir/aibox" --version > "$evidence_dir/version.txt"
"$evidence_dir/aibox" --help > "$evidence_dir/help.txt"
"$evidence_dir/aibox" inspect --context local --project spec/v1/examples/minimal --format json > "$evidence_dir/inspect.json"
uv run --with 'jsonschema[format]==4.23.0' python -c '
import json, sys
from jsonschema import Draft202012Validator, FormatChecker
schema = json.load(open("spec/v1/operation-result.schema.json", encoding="utf-8"))
with open(sys.argv[1], encoding="utf-8") as source:
    result = json.load(source)
Draft202012Validator(schema, format_checker=FormatChecker()).validate(result)
assert result["operation"] == "inspect_workspace"
assert result["outcome"] == "succeeded"
assert result["data"]["state"] == "unknown"
' "$evidence_dir/inspect.json"

node spec/v1/scripts/inventory.mjs --check
node spec/v1/scripts/validate.mjs
uv run --script spec/v1/scripts/validate-contracts.py
./scripts/build-docs.sh --destination "$evidence_dir/hugo"
git diff --check

printf 'V1-03 source SHA: %s\n' "$candidate_sha"
printf 'Linux binary SHA-256: '
sha256sum "$evidence_dir/aibox"
printf 'macOS arm64 binary SHA-256: '
sha256sum "$evidence_dir/aibox-darwin-arm64"
printf 'Evidence directory: %s\n' "$evidence_dir"

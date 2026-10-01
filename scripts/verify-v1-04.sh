#!/usr/bin/env bash
# Verify V1-04's offline Go/CLI/MCP, Template, and docs contracts.
# Live engine lifecycle qualification is deliberately performed by the
# separate host script scripts/verify-v1-04-host.sh.
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"
if [[ -n "$(git status --porcelain --untracked-files=all)" ]]; then
  printf 'V1-04 candidate worktree must be clean before evidence capture.\n' >&2
  exit 2
fi
evidence_dir="$(mktemp -d /tmp/aibox-v1-04-offline.XXXXXX)"
candidate_sha="$(git rev-parse HEAD)"
started_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
cli_version='0.89.0'

go test ./...
go test -race ./...
go vet ./...
go mod verify
go test -run '^TestV104StarterCLIMCPParity$' -count=1 -v ./cmd/aibox > "$evidence_dir/mcp-parity.txt"
go build -o "$evidence_dir/aibox" ./cmd/aibox
GOOS=darwin GOARCH=arm64 go build -o "$evidence_dir/aibox-darwin-arm64" ./cmd/aibox

for example in minimal customized custom-user bind-home; do
  "$evidence_dir/aibox" inspect --context local --project "spec/v1/examples/$example" --format json > "$evidence_dir/inspect-$example.json"
done

node spec/v1/scripts/inventory.mjs --check
node spec/v1/scripts/validate.mjs
node spec/v1/scripts/validate-examples.mjs
uv run --script spec/v1/scripts/validate-contracts.py
node spec/v1/scripts/render-devcontainer-examples.mjs --check

# Run upstream's own Template metadata documentation generator at an exact
# npm package version. Its configuration reader requires the host engine and
# therefore belongs to the separate host verification step.
npx --yes --package "@devcontainers/cli@$cli_version" devcontainer --version > "$evidence_dir/devcontainer-cli-version.txt"
[[ "$(tr -d '\r\n' < "$evidence_dir/devcontainer-cli-version.txt")" == "$cli_version" ]] || {
  printf 'Expected Dev Container CLI %s.\n' "$cli_version" >&2
  exit 1
}
template_stage="$evidence_dir/template-stage"
mkdir -p "$template_stage/src"
cp -a templates/src/minimal templates/src/curated "$template_stage/src/"
npx --yes --package "@devcontainers/cli@$cli_version" devcontainer templates generate-docs --project-folder "$template_stage/src" > "$evidence_dir/template-docs.txt" 2>&1
for template in minimal curated; do
  test -s "$template_stage/src/$template/README.md"
  for harness in codex none; do
    project="$evidence_dir/rendered-$template-$harness"
    mkdir -p "$project"
    node templates/render.mjs --template "$template" --output "$project" --harness "$harness" > "$evidence_dir/render-$template-$harness.txt"
    node - "$project/.devcontainer/devcontainer.json" "$harness" <<'NODE'
const fs=require('node:fs');const [file,harness]=process.argv.slice(2);
const config=JSON.parse(fs.readFileSync(file,'utf8').split('\n').filter(x=>!x.trimStart().startsWith('//')).join('\n'));
if(config.build.args.HARNESS!==harness)throw Error(`rendered ${file} harness mismatch`);
if(config.containerUser!=='aibox'||config.remoteUser!=='aibox')throw Error(`rendered ${file} user mismatch`);
NODE
  done
done

./scripts/build-docs.sh --destination "$evidence_dir/hugo"
git diff --check

finished_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
node - "$evidence_dir" "$candidate_sha" "$started_at" "$finished_at" "$cli_version" <<'NODE'
const {createHash} = require('node:crypto');
const {readFileSync, readdirSync, writeFileSync} = require('node:fs');
const {join} = require('node:path');
const [directory, sourceCommit, startedAt, finishedAt, cliVersion] = process.argv.slice(2);
const digest = name => 'sha256:' + createHash('sha256').update(readFileSync(join(directory, name))).digest('hex');
const examples = ['minimal', 'customized', 'custom-user', 'bind-home'];
const evidence = {
  phase: 'V1-04', verification: 'offline', sourceCommit, actor: 'local-verifier', startedAt, finishedAt,
  templates: ['minimal', 'curated'], examples,
  go: {nativeTarget: `${process.platform}/${process.arch}`, binary: 'aibox', binaryDigest: digest('aibox'),
    darwinArm64Binary: 'aibox-darwin-arm64', darwinArm64Digest: digest('aibox-darwin-arm64')},
  mcpParity: {test: 'TestV104StarterCLIMCPParity', transcript: 'mcp-parity.txt', transcriptDigest: digest('mcp-parity.txt')},
  upstreamDevcontainerCli: {package: `@devcontainers/cli@${cliVersion}`, versionTranscript: 'devcontainer-cli-version.txt',
    templateMetadataCheck: 'templates/src/{minimal,curated}/README.md generated in isolated staging copy'},
  checks: ['go test ./...', 'go test -race ./...', 'go vet ./...', 'go mod verify',
    'CLI/MCP parity on all four native starter examples', 'Linux native and macOS arm64 cross build',
    'inventory', 'specification', 'Template metadata/content contract', 'Draft 2020-12 fixtures',
    'upstream Dev Container CLI Template docs generation', 'Template render option fixtures for codex and none',
    'generated example freshness', 'v1 Hugo build', 'git diff --check'],
  lifecycle: {status: 'not_run', reason: 'The pinned CLI read-configuration and lifecycle commands require the host container engine; verify-v1-04-host.sh runs them.'},
  artifacts: readdirSync(directory).filter(name => name !== 'evidence.json'),
  failed: [], skipped: ['devcontainer build/up/exec pending host-run verification'],
};
writeFileSync(join(directory, 'evidence.json'), JSON.stringify(evidence, null, 2) + '\n');
NODE

printf 'V1-04 offline candidate SHA: %s\n' "$candidate_sha"
printf 'Evidence directory: %s\n' "$evidence_dir"

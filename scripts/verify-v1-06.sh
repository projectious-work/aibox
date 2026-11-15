#!/usr/bin/env bash
# Verify V1-06's offline gate: the frozen contract is unchanged, Go
# format/vet/tests, the static build, the roadmap demo with the
# effective-configuration view, the earlier phases' specification checks and
# the v1 Hugo build. It runs in localhost/aibox-docs-build:go1.27.0. The V1-06
# acceptance harness needs rootless Podman and the development container's
# Python, so its unit tests and journeys run separately there:
#   python3 -I -m unittest discover -s tests/acceptance/v1-06/unit -t tests/acceptance/v1-06
#   tests/acceptance/v1-06/run.sh --aibox <evidence>/aibox --evidence <dir>
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"
if [[ -n "$(git status --porcelain --untracked-files=all)" ]]; then
  printf 'V1-06 candidate worktree must be clean before evidence capture.\n' >&2
  exit 2
fi
evidence_dir="$(mktemp -d "${TMPDIR:-/tmp}/aibox-v1-06-offline.XXXXXX")"
candidate_sha="$(git rev-parse HEAD)"
started_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
contract_commit=7c2d4b89566d1869d720cd8b23a5afaae90ce68f

# The frozen contract is unchanged.
(cd spec/settings-and-logging/v1 && sha256sum --quiet -c FREEZE.sha256)
if [[ -n "$(git diff --name-only "$contract_commit" HEAD -- spec/settings-and-logging/v1)" ]]; then
  printf 'spec/settings-and-logging/v1 differs from the contract commit.\n' >&2
  exit 1
fi

unformatted="$(gofmt -l cmd internal)"
if [[ -n "$unformatted" ]]; then
  printf 'gofmt lists files:\n%s\n' "$unformatted" >&2
  exit 1
fi
go vet ./...
go test ./...
go test -race ./...
go mod verify
CGO_ENABLED=0 go build -o "$evidence_dir/aibox" ./cmd/aibox
GOOS=darwin GOARCH=arm64 go build -o "$evidence_dir/aibox-darwin-arm64" ./cmd/aibox

# Roadmap demo (J-CFG-24): effective settings and their sources are visible.
demo_home="$(mktemp -d)"
mkdir -p "$demo_home/.config/aibox"
printf '{"schemaVersion":"1","logging":{"level":"debug"}}\n' > "$demo_home/.config/aibox/settings.json"
chmod 0600 "$demo_home/.config/aibox/settings.json"
env -i HOME="$demo_home" PATH="$PATH" "$evidence_dir/aibox" inspect --context local \
  --project examples/customized --effective-config --format json \
  > "$evidence_dir/inspect-customized-effective.json" 2> "$evidence_dir/inspect-customized-effective.stderr"
rm -rf "$demo_home"

node scripts/spec/validate.mjs
node scripts/spec/validate-examples.mjs
uv run --script scripts/spec/validate-contracts.py
node scripts/spec/render-devcontainer-examples.mjs --check

./scripts/build-docs.sh --destination "$evidence_dir/hugo" > "$evidence_dir/hugo-build.txt"
git diff --check

# v0 stays the current line: this branch changes no v0 source.
v1_base="$(git merge-base HEAD origin/main)"
v0_changes="$(git diff --name-only "$v1_base" HEAD -- cli images)"
if [[ -n "$v0_changes" ]]; then
  printf 'V1-06 changes v0 sources:\n%s\n' "$v0_changes" >&2
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
const pages = Number((text('hugo-build.txt').match(/Pages\s*\│?\s*(\d+)/) || [])[1]);
const demo = JSON.parse(text('inspect-customized-effective.json'));
const level = (demo.data.configuration || []).find(entry => entry.key === 'logging.level') || {};
const evidence = {
  phase: 'V1-06', verification: 'offline', spec: 'settings-and-logging@v1', sourceCommit, startedAt, finishedAt,
  actor: process.env.AIBOX_VERIFY_ACTOR || 'local-verifier',
  host: {platform: process.platform, arch: process.arch, kernel: os.release()},
  go: {nativeTarget: `${process.platform}/${process.arch}`, binary: 'aibox', static: true, binaryDigest: digest('aibox'),
    darwinArm64Binary: 'aibox-darwin-arm64', darwinArm64Digest: digest('aibox-darwin-arm64')},
  demo: {command: 'aibox inspect --context local --project examples/customized --effective-config --format json',
    userSettings: {logging: {level: 'debug'}}, result: 'inspect-customized-effective.json',
    resultDigest: digest('inspect-customized-effective.json'), outcome: demo.outcome,
    loggingLevel: {value: level.value, layer: level.layer}},
  documentation: {command: './scripts/build-docs.sh', pages, result: 'passed'},
  v0: {current: true, check: `no change under cli/ or images/ since ${v1Base}`},
  checks: ['FREEZE.sha256 and no change to spec/settings-and-logging/v1', 'gofmt -l', 'go vet ./...',
    'go test ./...', 'go test -race ./...', 'go mod verify', 'static Linux build and macOS arm64 cross build',
    'roadmap demo with --effective-config', 'specification and example checks', 'Draft 2020-12 fixtures',
    'v1 Hugo build', 'git diff --check', 'no v0 source change'],
  containerAcceptance: {status: 'separate', runner: 'tests/acceptance/v1-06/run.sh',
    harnessUnitTests: 'tests/acceptance/v1-06/unit'},
  artifacts: readdirSync(directory).filter(name => name !== 'evidence.json' && name !== 'hugo'),
  failed: [], skipped: [],
};
if (!(pages > 0) || demo.outcome !== 'succeeded' || level.value !== 'debug' || level.layer !== 'user') {
  throw new Error('V1-06 offline evidence is incomplete: ' + JSON.stringify({pages, outcome: demo.outcome, level}));
}
writeFileSync(join(directory, 'evidence.json'), JSON.stringify(evidence, null, 2) + '\n');
NODE

printf 'V1-06 offline candidate SHA: %s\n' "$candidate_sha"
printf 'Evidence directory: %s\n' "$evidence_dir"

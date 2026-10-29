#!/usr/bin/env bash
# V1-05 aibox tool Features: generate sources, check them, run component
# tests, and publish them for acceptance tests.
#
#   scripts/v1-05-features.sh publish --registry <host>/<prefix> [--feature <id>]...
#       Publishes every aibox-owned Feature (or only the named ones) to
#       <host>/<prefix>/<catalog-id> with the tags <featureVersion>, <major>
#       and <major>.<minor>, rewriting aibox dependsOn and installsAfter
#       references to <host>/<prefix>. Upstream references are unchanged.
#       Runs from the repository root and uses the node and devcontainer
#       (pinned Dev Container CLI 0.89.0) executables found on PATH.
#   scripts/v1-05-features.sh generate
#       Regenerates devcontainer-feature.json, lib/aibox-feature.sh and
#       lib/catalog.sh for every Feature directory with an install.sh.
#   scripts/v1-05-features.sh check [--complete]
#       Fails when generated files drift from the catalog or the library.
#   scripts/v1-05-features.sh test
#       Runs the component tests in features/test.
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cli_version=0.89.0

fail() {
  printf 'v1-05-features: %s\n' "$*" >&2
  exit 1
}

usage() {
  sed -n '2,19p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//' >&2
  exit 2
}

command -v node >/dev/null 2>&1 || fail 'Node.js is required (node not found on PATH)'

publish() {
  local registry="" features_dir="" args=()
  while (($#)); do
    case "$1" in
      --registry) (($# >= 2)) || fail '--registry needs <host>/<prefix>'; registry=$2; shift 2 ;;
      --feature) (($# >= 2)) || fail '--feature needs a catalog ID'; args+=(--feature "$2"); shift 2 ;;
      --features-dir) (($# >= 2)) || fail '--features-dir needs a directory'; features_dir=$2; shift 2 ;;
      *) fail "unknown publish argument: $1" ;;
    esac
  done
  [[ -n "$registry" ]] || fail 'publish needs --registry <host>/<prefix>'
  [[ -z "$features_dir" ]] || args+=(--features-dir "$features_dir")

  work="$(mktemp -d "${TMPDIR:-/tmp}/aibox-v1-05-publish.XXXXXX")"
  trap 'rm -rf "${work:-}"' EXIT

  command -v devcontainer >/dev/null 2>&1 \
    || fail "the Dev Container CLI $cli_version is required (devcontainer not found on PATH)"
  local cli=(devcontainer)
  local version
  version="$("${cli[@]}" --version 2>/dev/null)" || fail 'the Dev Container CLI does not run'
  [[ "$version" == "$cli_version" ]] || fail "Dev Container CLI $cli_version is required, found ${version:-nothing}"

  local staged
  staged="$(node "$repo_root/scripts/v1-05-features.mjs" stage --registry "$registry" --out "$work/stage" "${args[@]}")" \
    || fail 'Feature sources are not ready to publish'

  local host namespace dir count=0
  while IFS=$'\t' read -r host namespace dir; do
    [[ -n "$dir" ]] || continue
    printf 'Publishing %s/%s (%s)\n' "$host" "$namespace" "$(cd "$dir" && ls | tr '\n' ' ')"
    "${cli[@]}" features publish --registry "$host" --namespace "$namespace" --log-level info "$dir" \
      || fail "publishing $host/$namespace failed"
    count=$((count + 1))
  done <<<"$staged"
  ((count > 0)) || fail 'nothing was published'
  printf 'Published aibox Features to %s\n' "$registry"
}

(($#)) || usage
command=$1
shift
case "$command" in
  publish) publish "$@" ;;
  generate|check) exec node "$repo_root/scripts/v1-05-features.mjs" "$command" "$@" ;;
  test) exec node --test "$repo_root"/features/test/*.test.mjs ;;
  -h|--help|help) usage ;;
  *) fail "unknown command: $command" ;;
esac

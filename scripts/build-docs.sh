#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DOCS_ROOT="${PROJECT_ROOT}/docs-site"
DOCS_BASE_URL="${DOCS_BASE_URL:-https://projectious-work.github.io/aibox/v1.x/}"
BUILD_DIR="${DOCS_ROOT}/public"

BUILD_ARGS=("$@")
for ((i = 0; i < ${#BUILD_ARGS[@]}; i++)); do
  case "${BUILD_ARGS[i]}" in
    --destination)
      ((i + 1 < ${#BUILD_ARGS[@]})) || {
        echo "--destination requires a value." >&2
        exit 1
      }
      BUILD_DIR="${BUILD_ARGS[i + 1]}"
      ;;
    --destination=*)
      BUILD_DIR="${BUILD_ARGS[i]#--destination=}"
      ;;
  esac
done

RELEASES_DATA_DIR=""
if [[ -n "${DOCS_RELEASES_MANIFEST:-}" ]]; then
  command -v jq >/dev/null 2>&1 || {
    echo "jq is required to validate DOCS_RELEASES_MANIFEST." >&2
    exit 1
  }
  [[ -f "${DOCS_RELEASES_MANIFEST}" && -r "${DOCS_RELEASES_MANIFEST}" ]] || {
    echo "DOCS_RELEASES_MANIFEST must name a readable JSON file: ${DOCS_RELEASES_MANIFEST}" >&2
    exit 1
  }
  jq -e '
    (.schemaVersion == 1) and
    (([.lines[]?.line] | sort) as $lines |
      ($lines == ["v0.x", "v1.x"])) and
    (all(.lines[];
      (.current.version | type == "string" and length > 0) and
      (.current.url | type == "string" and test("^https?://")) and
      (.releases | type == "array") and
      all(.releases[]; (.version | type == "string" and length > 0) and (.url | type == "string" and test("^https?://")))
    ))
  ' "${DOCS_RELEASES_MANIFEST}" >/dev/null || {
    echo "DOCS_RELEASES_MANIFEST has an invalid releases.json schema." >&2
    exit 1
  }
  RELEASES_DATA_DIR="$(mktemp -d)"
  trap '[[ -z "${RELEASES_DATA_DIR}" ]] || rm -rf -- "${RELEASES_DATA_DIR}"' EXIT
  cp -- "${DOCS_RELEASES_MANIFEST}" "${RELEASES_DATA_DIR}/releases.json"
fi

if [[ "${BUILD_DIR}" != /* ]]; then
  BUILD_DIR="${DOCS_ROOT}/${BUILD_DIR}"
fi

command -v hugo >/dev/null 2>&1 || {
  echo "Hugo extended is required: https://gohugo.io/installation/" >&2
  exit 1
}
command -v go >/dev/null 2>&1 || {
  echo "Go is required to resolve the pinned Hugo module." >&2
  exit 1
}
command -v npm >/dev/null 2>&1 || {
  echo "Node.js and npm are required for the brand theme assets." >&2
  exit 1
}

if [[ ! -d "${DOCS_ROOT}/node_modules/@tabler/icons/icons/outline" || ! -x "${DOCS_ROOT}/node_modules/.bin/tailwindcss" ]]; then
  npm --prefix "${DOCS_ROOT}" ci
fi

if [[ -n "${RELEASES_DATA_DIR}" ]]; then
  HUGO_DATADIR="${RELEASES_DATA_DIR}" hugo --source "${DOCS_ROOT}" --gc --minify --cleanDestinationDir \
    --baseURL "${DOCS_BASE_URL}" "${BUILD_ARGS[@]}"
else
  hugo --source "${DOCS_ROOT}" --gc --minify --cleanDestinationDir \
    --baseURL "${DOCS_BASE_URL}" "${BUILD_ARGS[@]}"
fi

: > "${BUILD_DIR}/.nojekyll"

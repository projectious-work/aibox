#!/usr/bin/env bash
# record-screenshots.sh — generate documentation screenshots using VHS
#
# This script launches a real dev-box container (sibling container via Docker
# socket), runs VHS tape scripts inside it, and copies the output screenshots
# to docs/assets/screenshots/.
#
# Prerequisites:
#   - Docker socket mounted at /var/run/docker.sock
#   - VHS installed (included in dev-container Dockerfile)
#   - Published dev-box images available (ghcr.io/projectious-work/dev-box)
#
# Usage:
#   ./scripts/record-screenshots.sh [VERSION]
#   ./scripts/record-screenshots.sh 0.8.0
#   ./scripts/record-screenshots.sh          # uses "latest"

set -euo pipefail

VERSION="${1:-latest}"
IMAGE="ghcr.io/projectious-work/dev-box:python-${VERSION}"
CONTAINER_NAME="devbox-recording"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
OUTPUT_DIR="${PROJECT_ROOT}/docs/assets/screenshots"
TAPE_DIR="${PROJECT_ROOT}/scripts/recordings"

info()  { printf '\033[1;36m==>\033[0m %s\n' "$1"; }
ok()    { printf '\033[1;32m ✓\033[0m %s\n' "$1"; }
die()   { printf '\033[1;31mError:\033[0m %s\n' "$1" >&2; exit 1; }

# Check prerequisites
command -v docker &>/dev/null || die "docker CLI not found"
docker info &>/dev/null 2>&1  || die "Cannot connect to Docker daemon (is the socket mounted?)"
command -v vhs &>/dev/null    || die "VHS not found (install: apt install vhs or see charmbracelet/vhs)"

mkdir -p "${OUTPUT_DIR}"

# ─── Strategy A: Record layouts inside a real dev-box container ──────────────
# We exec into a running dev-box container to capture real Zellij layouts.

record_in_container() {
  local tape_name="$1"
  local tape_file="${TAPE_DIR}/${tape_name}.tape"

  [[ -f "${tape_file}" ]] || die "Tape file not found: ${tape_file}"

  info "Recording ${tape_name}..."

  # Copy tape file into the container
  docker cp "${tape_file}" "${CONTAINER_NAME}:/tmp/${tape_name}.tape"

  # Run VHS inside the container
  docker exec -e VHS_NO_SANDBOX=true "${CONTAINER_NAME}" \
    vhs "/tmp/${tape_name}.tape"

  # Copy outputs back
  docker cp "${CONTAINER_NAME}:/docs/assets/screenshots/." "${OUTPUT_DIR}/" 2>/dev/null || true

  ok "Recorded ${tape_name}"
}

# ─── Strategy B: Record CLI demos locally (no container needed) ──────────────
# dev-box init and similar commands run on the host CLI.

record_local() {
  local tape_name="$1"
  local tape_file="${TAPE_DIR}/${tape_name}.tape"

  [[ -f "${tape_file}" ]] || die "Tape file not found: ${tape_file}"

  info "Recording ${tape_name} (local)..."
  cd "${PROJECT_ROOT}"
  VHS_NO_SANDBOX=true vhs "${tape_file}"
  ok "Recorded ${tape_name}"
}

# ─── Main ────────────────────────────────────────────────────────────────────

info "Starting recording session (image: ${IMAGE})"

# Pull the image
info "Pulling ${IMAGE}..."
docker pull "${IMAGE}" || die "Failed to pull ${IMAGE}"

# Start a recording container
info "Starting recording container..."
docker rm -f "${CONTAINER_NAME}" 2>/dev/null || true
docker run -d \
  --name "${CONTAINER_NAME}" \
  -e TERM=xterm-256color \
  -e COLORTERM=truecolor \
  "${IMAGE}" \
  sleep infinity
ok "Container ${CONTAINER_NAME} running"

# Install VHS inside the recording container
info "Installing VHS in recording container..."
docker exec "${CONTAINER_NAME}" bash -c '
  ARCH=$(uname -m)
  if [ "$ARCH" = "aarch64" ]; then VHS_ARCH="arm64"; else VHS_ARCH="amd64"; fi
  curl -fsSL "https://github.com/charmbracelet/vhs/releases/latest/download/vhs_linux_${VHS_ARCH}.deb" \
    -o /tmp/vhs.deb
  apt-get update -qq
  apt-get install -y -qq /tmp/vhs.deb ffmpeg chromium ttyd >/dev/null 2>&1
  rm -f /tmp/vhs.deb
  mkdir -p /docs/assets/screenshots
'
ok "VHS installed in container"

# Record layouts
for layout in dev-layout focus-layout cowork-layout; do
  record_in_container "${layout}"
done

# Record CLI demos locally
for demo in init-demo; do
  record_local "${demo}"
done

# Cleanup
info "Cleaning up..."
docker rm -f "${CONTAINER_NAME}" >/dev/null 2>&1
ok "Container removed"

# Summary
echo ""
info "Screenshots generated:"
ls -1 "${OUTPUT_DIR}"/*.png 2>/dev/null || echo "  (no PNG files)"
echo ""
info "GIFs generated:"
ls -1 "${OUTPUT_DIR}"/*.gif 2>/dev/null || echo "  (no GIF files)"

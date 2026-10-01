#!/bin/sh
# Checksummed upstream binaries, no credential transfer or project-state writes.
set -eu
harness=${1:-codex}
case "$harness" in codex|none) ;; *) echo "Unsupported harness: $harness" >&2; exit 2 ;; esac
case "$(dpkg --print-architecture)" in
  amd64)
    arch=x86_64
    yazi_sha=1031a02560d053301537195a6661d227c15cb4ce5c30481050b31e2b88681bff
    codex_sha=9c1a16086e971578f0c16d58d07fe8295791c0a597eeb7fd6388f8a3f1579cee
    ;;
  arm64)
    arch=aarch64
    yazi_sha=dff4a6774069bdde616c3cd145275b462a5887229cd9bf6341e78063da37e84f
    codex_sha=b994c71d1c48aa4e340aa3aa295631ba54b6b94f0e2d5ece3ca744e86fe52d9d
    ;;
  *) echo 'Only linux/amd64 and linux/arm64 are supported' >&2; exit 2 ;;
esac
scratch=$(mktemp -d)
trap 'rm -rf "$scratch"' EXIT HUP INT TERM
curl --proto '=https' --tlsv1.2 -fsSL --retry 3 \
  "https://github.com/sxyazi/yazi/releases/download/v26.5.6/yazi-$arch-unknown-linux-musl.zip" \
  -o "$scratch/yazi.zip"
printf '%s  %s\n' "$yazi_sha" "$scratch/yazi.zip" | sha256sum -c -
unzip -q "$scratch/yazi.zip" -d "$scratch"
install -m 0755 "$scratch/yazi-$arch-unknown-linux-musl/yazi" /usr/local/bin/yazi
install -m 0755 "$scratch/yazi-$arch-unknown-linux-musl/ya" /usr/local/bin/ya
install -D -m 0644 "$scratch/yazi-$arch-unknown-linux-musl/LICENSE" /usr/local/share/licenses/yazi/LICENSE
if [ "$harness" = codex ]; then
  curl --proto '=https' --tlsv1.2 -fsSL --retry 3 \
    "https://github.com/openai/codex/releases/download/rust-v0.107.0/codex-$arch-unknown-linux-musl.tar.gz" \
    -o "$scratch/codex.tar.gz"
  printf '%s  %s\n' "$codex_sha" "$scratch/codex.tar.gz" | sha256sum -c -
  tar -xzf "$scratch/codex.tar.gz" -C "$scratch" "codex-$arch-unknown-linux-musl"
  install -m 0755 "$scratch/codex-$arch-unknown-linux-musl" /usr/local/bin/codex
fi

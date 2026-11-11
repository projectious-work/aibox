#!/bin/sh
# data/visualization (atomic): vega-cli and vega-lite as global npm packages
# (https://vega.github.io/vega/usage/#cli), pinned by sha512 integrity. The
# canvas renderer used for PNG and PDF output is installed beside them so
# vega-lite resolves it too; its native build is the producer's prebuilt
# binary, pinned by SHA256, instead of an unpinned install-script download.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

canvas_version=3.2.3
canvas_integrity=sha512-PzE5nJZPz72YUAfo8oTp0u3fqqY7IzlTubneAihqDYAUcBk7ryeCmBbdJBEdaH0bptSOe2VT2Zwcb3UaFyaSWw==
arch=$(af_arch_name x64 arm64)
case "$AF_ARCH" in
  amd64) canvas_sha256=886d1cc270d4caad1d1698ee97532e79c60e992e16fe5f37f2d3bb058f4292b0 ;;
  arm64) canvas_sha256=ef81f0b1032d52fa5a9c8e4f642cd3f1b6b6b951edaa9c9468b849cb5180f83b ;;
esac

af_npm_fetch vega-cli 6.4.0 sha512-7ZlYi4URdSeha/ANkIgxW91nLxtaHdVm3R9Y62xaXIxSuB1Ga1YQ4QX9wBXxDP3GHEVE9yJQjBvm4QbKkpHRzA==
af_npm_fetch vega-lite 6.4.3 sha512-d/7hPjfz560UERaQuTmGgIVfXAe3g2hJWeC+igDeaGohUdEoNrHLXgR/yTOBT8vV/lIuuKnw+0/xWWblkDwkMQ==
af_npm_fetch canvas "$canvas_version" "$canvas_integrity"
af_fetch "https://github.com/Automattic/node-canvas/releases/download/v$canvas_version/canvas-v$canvas_version-napi-v7-linux-$arch.tar.gz" \
  "$AF_TMP/canvas-prebuilt.tar.gz" "$canvas_sha256"
af_npm_install
for package in $(find "$AF_ROOT/usr/local/lib/node_modules" -type d -path '*/node_modules/canvas'); do
  grep -q "\"version\": \"$canvas_version\"" "$package/package.json" \
    || af_die "$package is not canvas $canvas_version"
  tar -xzf "$AF_TMP/canvas-prebuilt.tar.gz" -C "$package" --no-same-owner
done

af_finish

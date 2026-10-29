#!/bin/sh
# python/uv: uv and uvx from the producer's release archive
# (https://docs.astral.sh/uv/getting-started/installation/), pinned by the
# SHA256 the producer publishes next to each archive.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

triple=$(af_arch_name x86_64 aarch64)-unknown-linux-gnu
case "$AF_VERSION/$AF_ARCH" in
  0.12.5/amd64) sha256=68a509da24b06b4223a1c0175fb5eb5bc79342b76cbeff0cfe51ac3f5b17b6b2 ;;
  0.12.5/arm64) sha256=9bf43b4d1a07665bf64d4c4e710930b382321a785e0eb10aac07f46471f86a31 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac

af_fetch "https://github.com/astral-sh/uv/releases/download/$AF_VERSION/uv-$triple.tar.gz" "$AF_TMP/uv.tar.gz" "$sha256"
tar -xzf "$AF_TMP/uv.tar.gz" -C "$AF_TMP"
install -d "$AF_ROOT/usr/local/bin"
install -m 0755 "$AF_TMP/uv-$triple/uv" "$AF_TMP/uv-$triple/uvx" "$AF_ROOT/usr/local/bin/"

af_finish

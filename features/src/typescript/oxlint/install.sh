#!/bin/sh
# typescript/oxlint: oxlint from the producer's release artifact
# (https://oxc.rs/docs/guide/usage/linter.html),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x86_64 aarch64)
case "$AF_VERSION/$AF_ARCH" in
  1.90.0/amd64) sha256=dd5b0c99ca7917fd6d3caabae260acc5bf7d65ae094541971862602fe288dca4 ;;
  1.90.0/arm64) sha256=dda6ccf8e94d4cddb82672c9d230ab04cc52e13a31075e88e53bb97b31311468 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/oxc-project/oxc/releases/download/oxlint_v$AF_VERSION/oxlint-$arch-unknown-linux-gnu.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/oxlint-$arch-unknown-linux-gnu" oxlint

af_finish

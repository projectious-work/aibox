#!/bin/sh
# infrastructure/opentofu: opentofu from the producer's release artifact
# (https://opentofu.org/docs/intro/install/standalone/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  1.12.6/amd64) sha256=50a6106fa4de523d09c87af85f3db1dd47535fc005727fdca6852146476b88ec ;;
  1.12.6/arm64) sha256=9bd0228a81bcd0c88f7045c74378f45a815779f19897191dff7d9efba9976b9e ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/opentofu/opentofu/releases/download/v$AF_VERSION/tofu_${AF_VERSION}_linux_$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/tofu" tofu

af_finish

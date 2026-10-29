#!/bin/sh
# go/golangci-lint: golangci-lint from the producer's release artifact
# (https://golangci-lint.run/docs/welcome/install/local/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  2.12.2/amd64) sha256=8df580d2670fed8fa984aac0507099af8df275e665215f5c7a2ae3943893a553 ;;
  2.12.2/arm64) sha256=44cd40a8c76c86755375adfeea52cfd3533cb43d7bd647771e0ae065e166df3a ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/golangci/golangci-lint/releases/download/v$AF_VERSION/golangci-lint-$AF_VERSION-linux-$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/golangci-lint-$AF_VERSION-linux-$arch/golangci-lint" golangci-lint

af_finish

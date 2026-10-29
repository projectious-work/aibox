#!/bin/sh
# security/scan/osv-scanner: osv-scanner from the producer's release artifact
# (https://google.github.io/osv-scanner/installation/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  2.4.0/amd64) sha256=15314940c10d26af9c6649f150b8a47c1262e8fc7e17b1d1029b0e479e8ed8a0 ;;
  2.4.0/arm64) sha256=44e580752910f0ff36ec99aff59af20f65df1e859aa31e5605a8f0d055b496e9 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/google/osv-scanner/releases/download/v$AF_VERSION/osv-scanner_linux_$arch"

af_fetch "$url" "$AF_TMP/artifact" "$sha256"
af_install_bin "$AF_TMP/artifact" osv-scanner

af_finish

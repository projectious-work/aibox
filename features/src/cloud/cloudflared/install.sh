#!/bin/sh
# cloud/cloudflared: cloudflared from the producer's release artifact
# (https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/downloads/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  2026.10.5/amd64) sha256=e1bc0cf084fd9e590554fe769a86b1db46701c3f847a43d902f7b974a4e16298 ;;
  2026.10.5/arm64) sha256=72e0d05a67e09da87cf495e409a28c3cc5994dd23fc66a982eebea6a55923b7d ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/cloudflare/cloudflared/releases/download/$AF_VERSION/cloudflared-linux-$arch"

af_fetch "$url" "$AF_TMP/artifact" "$sha256"
af_install_bin "$AF_TMP/artifact" cloudflared

af_finish

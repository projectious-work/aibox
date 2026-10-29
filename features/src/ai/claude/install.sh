#!/bin/sh
# ai/claude: claude from the producer's release artifact
# (https://code.claude.com/docs/en/setup#install-with-linux-package-managers),
# pinned by the SHA256 values in the release manifest, whose signature by the
# producer key 31DDDE24DDFAB679F42D7BD2BAA929FF1A7ECACE was verified when the
# pins were recorded.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  2.1.308/amd64) sha256=356733f185fa988e6bf3aa9a555269fdd167982a2e2a80d3c8a78037bb92ad58 ;;
  2.1.308/arm64) sha256=fdd22e6cbb734d580f5664c48df4294eeec850218ff07f186f9af30426396ee1 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://downloads.claude.ai/claude-code-releases/$AF_VERSION/linux-$arch/claude"

af_fetch "$url" "$AF_TMP/artifact" "$sha256"
af_install_bin "$AF_TMP/artifact" claude

af_finish

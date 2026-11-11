#!/bin/sh
# security/crypto/age: age from the producer's release artifact
# (https://github.com/FiloSottile/age#installation),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  1.3.2/amd64) sha256=cbe24006683f8eb669266162894b9a522a1af52f2665fbc63a4bb032ed26ac10 ;;
  1.3.2/arm64) sha256=6b8dc4333c53a5a57c9e5834e3a48f92605d7154014cd07269ff3327db5d37f4 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/FiloSottile/age/releases/download/v$AF_VERSION/age-v$AF_VERSION-linux-$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/age/age" age
af_install_bin "$AF_TMP/x/age/age-keygen" age-keygen

af_finish

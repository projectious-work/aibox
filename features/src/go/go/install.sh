#!/bin/sh
# go/go: the Go toolchain from the producer's archive (https://go.dev/doc/install),
# pinned by the SHA256 values in https://go.dev/dl/?mode=json&include=all.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

case "$AF_VERSION/$AF_ARCH" in
  1.25.12/amd64) sha256=234828b7a89e0e303d2556310ee549fbcf253d28de937bac3da13d6294262ac1 ;;
  1.25.12/arm64) sha256=8b5884aef89600aef5b0b051fb971f11f49bb996521e911f30f02a66884f7bd2 ;;
  1.26.6/amd64) sha256=708effb774be8237570d0add163225abbdfaf4fca28b2611df167beba4feef89 ;;
  1.26.6/arm64) sha256=d0507e9e9d7fe012aae570108cbd76c15de879e17130ab8cb90d4d7445cb1f2e ;;
  1.27.0/amd64) sha256=675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685 ;;
  1.27.0/arm64) sha256=51798d2c42d0e1c6ed7fd9f48728b4193abac9e8aad6dbac2fe96a81f5909bda ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac

af_fetch "https://go.dev/dl/go$AF_VERSION.linux-$AF_ARCH.tar.gz" "$AF_TMP/go.tar.gz" "$sha256"
rm -rf "$AF_ROOT/usr/local/go"
install -d "$AF_ROOT/usr/local/bin"
tar -xzf "$AF_TMP/go.tar.gz" -C "$AF_ROOT/usr/local" --no-same-owner
# Keep the pinned toolchain: go must not download another toolchain.
go_env="$AF_ROOT/usr/local/go/go.env"
if grep -q '^GOTOOLCHAIN=' "$go_env"; then
  sed -i 's/^GOTOOLCHAIN=.*/GOTOOLCHAIN=local/' "$go_env"
else
  printf 'GOTOOLCHAIN=local\n' >> "$go_env"
fi
ln -sfn ../go/bin/go "$AF_ROOT/usr/local/bin/go"
ln -sfn ../go/bin/gofmt "$AF_ROOT/usr/local/bin/gofmt"

af_finish

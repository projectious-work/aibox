#!/bin/sh
# rust/x86-64-cross: building for x86_64-unknown-linux-gnu. On amd64 that is
# the host target and rust/rustc already provides the linker, so nothing
# changes. On arm64 it installs the Debian x86-64 cross toolchain and adds
# the rustup target to every installed toolchain.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

if [ "$AF_ARCH" = arm64 ]; then
  af_apt_require gcc-x86-64-linux-gnu libc6-dev-amd64-cross
  af_rust_env
  af_apt_install gcc-x86-64-linux-gnu libc6-dev-amd64-cross
  for toolchain in $("$AF_RUSTUP" toolchain list | sed 's/ .*//'); do
    "$AF_RUSTUP" target add --toolchain "$toolchain" x86_64-unknown-linux-gnu \
      || af_fail DOWNLOAD_FAILED "source=rustup target=x86_64-unknown-linux-gnu toolchain=$toolchain"
  done
  chmod -R a+w "$RUSTUP_HOME" "$CARGO_HOME"
fi

af_finish

#!/bin/sh
# ai/gemini: the npm package @google/gemini-cli installed globally with the node/node Feature's npm
# (https://github.com/google-gemini/gemini-cli), pinned by the registry's sha512 integrity value.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

case "$AF_VERSION" in
  0.66.0) integrity=sha512-EAtsJ9/JxiphHlduzvhoCkb7I0LdOxyugB5NpUtVqKVGZIWjV3JmurvK89ZrvMbV3kwk5SdWINz1M5ZN/2K/Lw== ;;
  *) af_die "no pinned integrity for $AF_VERSION" ;;
esac
af_npm_fetch @google/gemini-cli "$AF_VERSION" "$integrity"
af_npm_install

af_finish

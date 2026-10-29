#!/bin/sh
# latex/music: MusiXTeX from the latex/core repository (frozen TeX Live 2025)
# and the Debian lilypond package at the base's candidate version.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init
af_apt_require gnupg lilypond

af_texlive_install musixtex
af_apt_install lilypond

af_finish

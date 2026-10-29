#!/bin/sh
# latex/recommended: TeX Live packages from the latex/core repository (frozen TeX Live
# 2025), verified through the signed package database.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_apt_require gnupg
af_texlive_install latex latex-bin latexmk collection-latexrecommended

af_finish

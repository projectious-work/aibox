#!/bin/sh
# latex/fonts: TeX Live packages from the latex/core repository (frozen TeX Live
# 2025), verified through the signed package database.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_apt_require gnupg
af_texlive_install luatex luaotfload luatexbase luacolor lua-ul fontspec unicode-math lualatex-math luahbtex lm lm-math sourcecodepro sourcesanspro tex-gyre gnu-freefont collection-fontsrecommended selnolig

af_finish

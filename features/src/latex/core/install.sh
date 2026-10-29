#!/bin/sh
# latex/core: TeX Live 2025 scheme-basic plus the v0 core package set, from
# the frozen TeX Live 2025 final repository
# (https://tug.org/texlive/quickinstall.html). The installer archive is
# pinned by SHA256; install-tl and tlmgr verify the repository's signed
# package database with the TeX Live key shipped in that archive, using a
# build-only GnuPG, and each package against the database checksums.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init
# perl runs install-tl and, at runtime, tlmgr.
af_apt_require perl gnupg

case "$AF_VERSION" in
  2025)
    repository=https://ftp.tu-chemnitz.de/pub/tug/historic/systems/texlive/2025/tlnet-final
    sha256=311df9f1477fd90c520159d1feddc2d6270f010d8349d1f6bdb9461a93b48a5c
    ;;
  *) af_die "no pinned installer for TeX Live $AF_VERSION" ;;
esac
platform=$(af_arch_name x86_64-linux aarch64-linux)
texdir="$AF_ROOT/usr/local/texlive"
packages="geometry fancyhdr lastpage multirow adjustbox microtype parskip setspace titlesec
  tabularray booktabs tabu colortbl xcolor graphics epstopdf svg hyperref bookmark cleveref url
  algorithms algorithmicx algorithm2e caption enumitem wrapfig footmisc appendix pdfpages pdflscape
  xkeyval iftex babel babel-english tools oberdiek symbol zapfding kvoptions kvsetkeys kvdefinekeys
  ltxcmds infwarerr epstopdf-pkg grfext pdftexcmds auxhook intcalc letltxmacro bitset bigintcalc
  atbegshi atveryend rerunfilecheck uniquecounter refcount gettitlestring hycolor stringenc
  pdfescape hobsub bytefield siunitx markdown soul ulem todonotes changes datetime2 tracklang was
  emoji ninecolors transparent spath3 nicematrix lipsum"

af_fetch "$repository/install-tl-unx.tar.gz" "$AF_TMP/install-tl-unx.tar.gz" "$sha256"
af_apt_install perl
af_helper_install gnupg
install -d "$AF_TMP/install-tl"
tar -xzf "$AF_TMP/install-tl-unx.tar.gz" -C "$AF_TMP/install-tl" --strip-components=1 --no-same-owner

cat > "$AF_TMP/texlive.profile" <<PROFILE
selected_scheme scheme-basic
TEXDIR $texdir
TEXMFLOCAL $texdir/texmf-local
TEXMFSYSCONFIG $texdir/texmf-config
TEXMFSYSVAR $texdir/texmf-var
TEXMFCONFIG ~/.texlive$AF_VERSION/texmf-config
TEXMFVAR ~/.texlive$AF_VERSION/texmf-var
TEXMFHOME ~/texmf
binary_$platform 1
instopt_adjustpath 1
instopt_adjustrepo 0
instopt_letter 0
instopt_portable 0
instopt_write18_restricted 1
tlpdbopt_autobackup 0
tlpdbopt_install_docfiles 0
tlpdbopt_install_srcfiles 0
tlpdbopt_sys_bin $AF_ROOT/usr/local/bin
tlpdbopt_sys_info $AF_ROOT/usr/local/share/info
tlpdbopt_sys_man $AF_ROOT/usr/local/share/man
PROFILE

TEXLIVE_INSTALL_NO_WELCOME=1 perl "$AF_TMP/install-tl/install-tl" \
  --profile "$AF_TMP/texlive.profile" --repository "$repository" \
  || af_fail DOWNLOAD_FAILED "url=$repository"
tlmgr="$texdir/bin/$platform/tlmgr"
# shellcheck disable=SC2086
"$tlmgr" install $packages || af_fail DOWNLOAD_FAILED "url=$repository"
"$tlmgr" path add
"$texdir/bin/$platform/pdflatex" --version | grep -q "TeX Live $AF_VERSION" \
  || af_die "pdflatex does not report TeX Live $AF_VERSION"

af_finish

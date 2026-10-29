# shellcheck shell=sh
# aibox tool Feature install conventions (spec/v1/08-native-artifacts.md,
# "Tool Feature interface"). POSIX sh, so the target guard also runs on
# unsupported bases.
#
# Each Feature carries a copy of this file as lib/aibox-feature.sh, generated
# by scripts/v1-05-features.mjs; edit features/lib/aibox-feature.sh only.
#
# install.sh sources lib/catalog.sh and then this file, and calls af_init
# before any change. af_init checks, in this order and before any change:
# target, version option, repeat installation, requires and version coupling.
# A Feature then checks Debian candidates with af_apt_require before it
# installs anything, and calls af_finish after a successful installation.
#
# Variables for install.sh after af_init:
#   AF_ROOT     filesystem root ("" in a container; a directory in tests)
#   AF_ARCH     Debian architecture (amd64 or arm64)
#   AF_TARGET   catalog target ID, for example debian-13-amd64
#   AF_VERSION  selected version option value ("" without a version option)
#   AF_TMP      private scratch directory, removed on exit

AF_ROOT=${AIBOX_FEATURE_ROOT:-}
AF_STATE_DIR="$AF_ROOT/var/lib/aibox-features"
AF_APT_LISTS="$AF_ROOT/var/lib/apt/lists"
AF_TMP=""
AF_APT_READY=0
AF_APT_CLEAN_LISTS=0
AF_HELPERS=""

# af_fail CODE DETAIL: write the contract error line and exit non-zero.
af_fail() {
  printf 'aibox-feature error AIBOX_FEATURE_%s feature=%s %s\n' "$1" "$AIBOX_FEATURE_ID" "$2" >&2
  exit 1
}

af_log() {
  printf 'aibox-feature %s: %s\n' "$AIBOX_FEATURE_ID" "$*"
}

af_cleanup() {
  af_status=$?
  set +e
  if [ -n "$AF_HELPERS" ]; then
    # shellcheck disable=SC2086
    DEBIAN_FRONTEND=noninteractive apt-get purge -y $AF_HELPERS >/dev/null
  fi
  if [ "$AF_APT_READY" = 1 ]; then
    apt-get clean >/dev/null 2>&1
  fi
  if [ "$AF_APT_CLEAN_LISTS" = 1 ]; then
    rm -rf "$AF_APT_LISTS"/*
  fi
  if [ -n "$AF_TMP" ]; then
    rm -rf "$AF_TMP"
  fi
  exit "$af_status"
}

af_init() {
  trap af_cleanup EXIT
  trap 'exit 129' HUP
  trap 'exit 130' INT
  trap 'exit 143' TERM
  af_guard_target
  af_guard_version
  af_guard_repeat
  af_guard_requires
  af_guard_coupling
  AF_TMP=$(mktemp -d)
}

af_os_release_field() {
  af_file="$AF_ROOT/etc/os-release"
  [ -r "$af_file" ] || af_file="$AF_ROOT/usr/lib/os-release"
  [ -r "$af_file" ] || return 0
  sed -n "s/^$1=//p" "$af_file" | head -n 1 | sed -e 's/^"//' -e 's/"$//' -e "s/^'//" -e "s/'\$//"
}

af_guard_target() {
  af_os_id=$(af_os_release_field ID)
  af_os_version=$(af_os_release_field VERSION_ID)
  af_os_codename=$(af_os_release_field VERSION_CODENAME)
  AF_ARCH=$(dpkg --print-architecture 2>/dev/null) || AF_ARCH=""
  AF_TARGET=""
  for af_target in $AIBOX_FEATURE_TARGETS; do
    # distro:versionId:release:architecture:id
    IFS=: read -r af_t_distro af_t_version af_t_release af_t_arch af_t_id <<EOF
$af_target
EOF
    if [ "$af_os_id" = "$af_t_distro" ] && [ "$af_os_version" = "$af_t_version" ] \
      && [ "$AF_ARCH" = "$af_t_arch" ] \
      && { [ -z "$af_os_codename" ] || [ "$af_os_codename" = "$af_t_release" ]; }; then
      AF_TARGET=$af_t_id
      return 0
    fi
  done
  af_fail UNSUPPORTED_TARGET "target=${af_os_id:-unknown}-${af_os_version:-unknown}-${AF_ARCH:-unknown} supported=$(af_target_ids)"
}

af_target_ids() {
  af_ids=""
  for af_target in $AIBOX_FEATURE_TARGETS; do
    af_ids="${af_ids:+$af_ids,}${af_target##*:}"
  done
  printf '%s' "$af_ids"
}

af_guard_version() {
  AF_VERSION=""
  [ -n "$AIBOX_FEATURE_VERSIONS" ] || return 0
  AF_VERSION=${VERSION:-$AIBOX_FEATURE_DEFAULT_VERSION}
  af_supported=$(printf '%s' "$AIBOX_FEATURE_VERSIONS" | tr ' ' ',')
  case "$AF_VERSION" in
    *[!A-Za-z0-9._+-]*)
      af_fail UNSUPPORTED_VERSION "version=$(printf '%s' "$AF_VERSION" | tr -c 'A-Za-z0-9._+-' '?') supported=$af_supported"
      ;;
  esac
  for af_candidate in $AIBOX_FEATURE_VERSIONS; do
    [ "$AF_VERSION" = "$af_candidate" ] && return 0
  done
  af_fail UNSUPPORTED_VERSION "version=$AF_VERSION supported=$af_supported"
}

af_marker() {
  printf '%s/%s/version' "$AF_STATE_DIR" "$1"
}

# af_installed_version FEATURE-ID: print the recorded version of an installed
# aibox Feature ("none" for a Feature without a version option), or nothing.
af_installed_version() {
  af_m=$(af_marker "$1")
  [ -f "$af_m" ] && head -n 1 "$af_m"
  return 0
}

af_guard_repeat() {
  af_installed=$(af_installed_version "$AIBOX_FEATURE_ID")
  [ -n "$af_installed" ] || return 0
  if [ "$af_installed" = "${AF_VERSION:-none}" ]; then
    af_log "version ${AF_VERSION:-none} is already installed; nothing to do"
    exit 0
  fi
  af_fail VERSION_CONFLICT "installed=$af_installed requested=${AF_VERSION:-none}"
}

af_guard_requires() {
  for af_required in $AIBOX_FEATURE_REQUIRES; do
    if [ -z "$(af_installed_version "$af_required")" ]; then
      af_fail DEPENDENCY_MISSING "missing=$af_required"
    fi
  done
}

af_guard_coupling() {
  [ -n "$AIBOX_FEATURE_COUPLING_FEATURE" ] || return 0
  af_coupled=$(af_installed_version "$AIBOX_FEATURE_COUPLING_FEATURE")
  case "$AIBOX_FEATURE_COUPLING_RULE" in
    equal) af_allowed=$AF_VERSION ;;
    table)
      af_allowed=""
      for af_entry in $AIBOX_FEATURE_COUPLING_TABLE; do
        [ "${af_entry%%=*}" = "$AF_VERSION" ] && af_allowed=$(printf '%s' "${af_entry#*=}" | tr '|' ' ')
      done
      ;;
    *) af_allowed="" ;;
  esac
  for af_candidate in $af_allowed; do
    [ "$af_coupled" = "$af_candidate" ] && return 0
  done
  af_fail COUPLING_MISMATCH "coupled=$AIBOX_FEATURE_COUPLING_FEATURE installed=${af_coupled:-none} allowed=$(printf '%s' "${af_allowed:-none}" | tr ' ' ',')"
}

# af_arch_name AMD64-NAME ARM64-NAME: print the producer's name for AF_ARCH.
af_arch_name() {
  case "$AF_ARCH" in
    amd64) printf '%s' "$1" ;;
    arm64) printf '%s' "$2" ;;
  esac
}

af_apt_prepare() {
  [ "$AF_APT_READY" = 1 ] && return 0
  AF_APT_READY=1
  if ! ls "$AF_APT_LISTS"/*_Packages* >/dev/null 2>&1; then
    AF_APT_CLEAN_LISTS=1
    DEBIAN_FRONTEND=noninteractive apt-get update >&2 || return 1
  fi
  return 0
}

af_package_installed() {
  [ "$(dpkg-query -W --showformat='${db:Status-Status}' "$1" 2>/dev/null)" = installed ]
}

af_installed_packages() {
  dpkg-query -W --showformat='${db:Status-Status} ${Package}\n' | sed -n 's/^installed //p' | sort -u
}

# af_apt_require PACKAGE...: fail with PACKAGE_UNAVAILABLE before any change
# unless every package that is not installed has an installation candidate
# in the base's configured repositories.
af_apt_require() {
  af_missing=""
  for af_pkg in "$@"; do
    af_package_installed "$af_pkg" || af_missing="$af_missing $af_pkg"
  done
  [ -n "$af_missing" ] || return 0
  if ! af_apt_prepare; then
    af_fail PACKAGE_UNAVAILABLE "packages=$(printf '%s' "${af_missing# }" | tr ' ' ',') reason=apt-get-update-failed"
  fi
  for af_pkg in $af_missing; do
    af_candidate=$(apt-cache policy "$af_pkg" 2>/dev/null | sed -n 's/^ *Candidate: *//p' | head -n 1)
    if [ -z "$af_candidate" ] || [ "$af_candidate" = "(none)" ]; then
      af_fail PACKAGE_UNAVAILABLE "package=$af_pkg"
    fi
  done
}

# af_apt_install PACKAGE...: install the candidate versions as runtime
# packages that stay in the image.
af_apt_install() {
  af_apt_require "$@"
  af_apt_prepare
  DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends "$@" \
    || af_fail DOWNLOAD_FAILED "packages=$(printf '%s' "$*" | tr ' ' ',') reason=apt-get-install-failed"
}

# af_helper_install PACKAGE...: install build-only helpers that are not yet
# installed; they and any packages they pulled in are purged on exit.
af_helper_install() {
  af_new=""
  for af_pkg in "$@"; do
    af_package_installed "$af_pkg" || af_new="$af_new $af_pkg"
  done
  [ -n "$af_new" ] || return 0
  af_apt_require $af_new
  af_installed_packages > "$AF_TMP/packages-before"
  # shellcheck disable=SC2086
  DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends $af_new \
    || af_fail DOWNLOAD_FAILED "packages=$(printf '%s' "${af_new# }" | tr ' ' ',') reason=apt-get-install-failed"
  af_installed_packages > "$AF_TMP/packages-after"
  AF_HELPERS=$(comm -13 "$AF_TMP/packages-before" "$AF_TMP/packages-after" | tr '\n' ' ')
}

# af_download URL DESTINATION
af_download() {
  if ! curl --proto '=https' --tlsv1.2 -fsSL --retry 3 --connect-timeout 30 -o "$2" "$1"; then
    rm -f "$2"
    af_fail DOWNLOAD_FAILED "url=$1"
  fi
}

af_verify_digest() {
  # af_verify_digest ALGORITHM FILE EXPECTED-HEX
  af_actual=$("${1}sum" "$2" | cut -d ' ' -f 1)
  if [ "$af_actual" != "$3" ]; then
    rm -f "$2"
    af_fail VERIFICATION_FAILED "file=${2##*/} algorithm=$1 expected=$3 actual=$af_actual"
  fi
}

af_verify_sha256() { af_verify_digest sha256 "$1" "$2"; }
af_verify_sha512() { af_verify_digest sha512 "$1" "$2"; }

# af_verify_integrity FILE sha512-BASE64: npm/SRI integrity value.
af_verify_integrity() {
  case "$2" in
    sha512-*) ;;
    *) af_fail VERIFICATION_FAILED "file=${1##*/} algorithm=unsupported-integrity" ;;
  esac
  af_expected=$(printf '%s' "${2#sha512-}" | base64 -d | od -A n -v -t x1 | tr -d ' \n')
  af_verify_digest sha512 "$1" "$af_expected"
}

# af_verify_minisign FILE SIGNATURE PUBLIC-KEY (needs the minisign helper)
af_verify_minisign() {
  if ! minisign -V -q -m "$1" -x "$2" -P "$3" >/dev/null 2>&1; then
    rm -f "$1"
    af_fail VERIFICATION_FAILED "file=${1##*/} algorithm=minisign"
  fi
}

# af_fetch URL DESTINATION SHA256: download a pinned artifact and verify it.
af_fetch() {
  af_download "$1" "$2"
  af_verify_sha256 "$2" "$3"
}

# af_remote_user: print the remote user the Dev Container CLI passes to
# Features when it is a non-root account that exists, otherwise nothing.
af_remote_user() {
  case "${_REMOTE_USER:-root}" in
    root|"") return 0 ;;
  esac
  id -u "$_REMOTE_USER" >/dev/null 2>&1 && printf '%s' "$_REMOTE_USER"
  return 0
}

af_finish() {
  mkdir -p "$AF_STATE_DIR/$AIBOX_FEATURE_ID"
  printf '%s\n' "${AF_VERSION:-none}" > "$(af_marker "$AIBOX_FEATURE_ID")"
  af_log "installed ${AF_VERSION:-version from catalog} for $AF_TARGET"
}

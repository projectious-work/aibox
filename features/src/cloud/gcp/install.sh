#!/bin/sh
# cloud/gcp: google-cloud-cli from the producer's signed APT repository
# (https://cloud.google.com/sdk/docs/install#deb). The repository and its key
# are configured only while this script runs, so the image keeps the base's
# APT sources. The key file is pinned by SHA256 (fingerprint
# 35BAA0B33E9EB396F59CA838C0BA5CE6DC6315A3); APT verifies the repository
# signature and package checksums.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

key_sha256=3ecc63922b7795eb23fdc449ff9396f9114cb3cf186d6f5b53ad4cc3ebfbb11f
keyring="$AF_ROOT/etc/apt/keyrings/aibox-google-cloud-cli.gpg"
source_list="$AF_ROOT/etc/apt/sources.list.d/aibox-google-cloud-cli.list"
af_cleanup_path "$keyring"
af_cleanup_path "$source_list"

for dir in "$AF_ROOT/etc/apt" "$AF_ROOT/etc/apt/keyrings" "$AF_ROOT/etc/apt/sources.list.d"; do
  if [ ! -d "$dir" ]; then
    mkdir "$dir"
    af_cleanup_path "$dir"
  fi
done
af_fetch https://packages.cloud.google.com/apt/doc/apt-key.gpg "$keyring" "$key_sha256"
printf 'deb [signed-by=%s] https://packages.cloud.google.com/apt cloud-sdk main\n' "/etc/apt/keyrings/aibox-google-cloud-cli.gpg" > "$source_list"
if ! ls "$AF_APT_LISTS"/*_Packages* >/dev/null 2>&1; then
  AF_APT_CLEAN_LISTS=1
fi
AF_APT_READY=1
DEBIAN_FRONTEND=noninteractive apt-get update >&2 \
  || af_fail PACKAGE_UNAVAILABLE "package=google-cloud-cli=$AF_VERSION-0 reason=apt-get-update-failed"
af_apt_require_version google-cloud-cli "$AF_VERSION-0"
DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends "google-cloud-cli=$AF_VERSION-0" \
  || af_fail DOWNLOAD_FAILED "packages=google-cloud-cli=$AF_VERSION-0 reason=apt-get-install-failed"

af_finish

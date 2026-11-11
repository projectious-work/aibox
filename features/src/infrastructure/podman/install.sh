#!/bin/sh
# infrastructure/podman: Podman and the v0 companion packages from the base's
# configured repositories (https://podman.io/docs/installation#debian),
# configured for rootless use by the remote user inside a container.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_debian_install podman podman-compose fuse-overlayfs passt slirp4netns uidmap
# Containers in a container have no systemd; use cgroupfs.
if [ ! -e "$AF_ROOT/etc/containers/containers.conf" ]; then
  install -d "$AF_ROOT/etc/containers"
  printf '[engine]\ncgroup_manager = "cgroupfs"\nevents_logger = "file"\n' > "$AF_ROOT/etc/containers/containers.conf"
fi
# Rootless Podman needs subordinate ID ranges for the remote user.
user=${_REMOTE_USER:-}
if [ -n "$user" ] && [ "$user" != root ] && id -u "$user" >/dev/null 2>&1; then
  for map in subuid subgid; do
    grep -q "^$user:" "$AF_ROOT/etc/$map" 2>/dev/null || printf '%s:100000:65536\n' "$user" >> "$AF_ROOT/etc/$map"
  done
fi

af_finish

#!/bin/bash
# AC-TOOLS on debian-13-amd64 from an Apple silicon Mac with OrbStack.
#
#   tests/acceptance/v1-05/run-orbstack-amd64.sh [--commit SHA] [--runs N]
#       [--keep-machine] [--yes] [-- runner options, e.g. --scope full --jobs 2]
#
# Creates a disposable amd64 Debian 13 OrbStack machine (x86_64 under
# OrbStack's Rosetta or QEMU translation), installs Debian's Docker in it,
# clones this repository into it at the candidate commit (default fbd178d6,
# the V1-05 handoff commit), runs tests/acceptance/v1-05/run-host.sh
# debian-13-amd64 N times, copies each evidence directory to
# $HOME/aibox-v1-05-evidence/ on the Mac with an operator.json that names the
# operator, the Mac and the machine, and deletes the machine. All output also
# goes to $HOME/aibox-v1-05-evidence/orbstack-amd64-<UTC time>.log; a failing
# step names itself there and on the terminal.
#
# Needs: macOS on Apple silicon, OrbStack with the `orb` CLI, about 80 GB free
# disk, and network access to deb.debian.org, Docker Hub, ghcr.io,
# snapshot.debian.org, nodejs.org and the producer download sites in the
# catalog. No credentials. Runs under translation are slower than native
# amd64 and are recorded as translated in operator.json.
set -Eeuo pipefail

commit=fbd178d6
runs=1
keep=0
assume_yes=0
runner_args=()
while [ $# -gt 0 ]; do
  case "$1" in
    --commit) commit=${2:?--commit needs a SHA}; shift 2 ;;
    --runs) runs=${2:?--runs needs a number}; shift 2 ;;
    --keep-machine) keep=1; shift ;;
    --yes) assume_yes=1; shift ;;
    --) shift; runner_args=("$@"); break ;;
    -h|--help) sed -n 2,21p "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *) echo "unknown option: $1 (runner options go after --)" >&2; exit 2 ;;
  esac
done
case "$runs" in ''|*[!0-9]*|0) echo "--runs must be a positive number" >&2; exit 2 ;; esac

say() { step=$*; printf '\n==> %s\n' "$*"; }
die() { printf 'error: %s\n' "$*" >&2; exit 1; }

[ "$(uname -s)" = Darwin ] || die "run this on the Mac, not inside a container"
[ "$(uname -m)" = arm64 ] || die "expected an Apple silicon Mac (uname -m is $(uname -m))"
command -v orb >/dev/null || die "the OrbStack CLI 'orb' is not on PATH"
orb status >/dev/null 2>&1 || die "OrbStack is not running; start it first"

repo=$(cd "$(dirname "$0")" && git rev-parse --show-toplevel)
git -C "$repo" fetch --quiet origin || echo "warning: git fetch failed; using local objects" >&2
full_commit=$(git -C "$repo" rev-parse --verify "$commit^{commit}") \
  || die "commit $commit is not in $repo; run git fetch origin"

free_gb=$(df -g "$HOME" | awk 'NR==2 {print $4}')
[ "${free_gb:-0}" -ge 80 ] || die "need about 80 GB free on the Mac disk; ${free_gb} GB free"

stamp=$(date -u +%Y%m%dT%H%M%SZ)
machine="aibox-v1-05-amd64-$(echo "$stamp" | tr 'TZ' 'tz')"
evidence_root="$HOME/aibox-v1-05-evidence"
operator=$(id -un)

cat <<EOF

AC-TOOLS host run on debian-13-amd64 (OrbStack, translated amd64)

  Operator:   $operator
  Commit:     $full_commit
  Runs:       $runs, runner options: ${runner_args[*]:-(defaults: sampled scope)}
  Machine:    $machine (amd64 Debian 13, created now, deleted at the end$( [ $keep = 1 ] && echo ' unless kept: --keep-machine'))
  Effects:    installs docker.io, docker-buildx, git and python3 in that machine;
              inside it starts a registry on localhost:5000, pulls golang:1.27.0
              and the base images, and builds about 184 workspace images per
              sampled run (466 with --scope full). Nothing changes on the Mac
              except the new evidence under $evidence_root/.
  Duration:   several hours per run under translation.
EOF
if [ "$assume_yes" != 1 ]; then
  printf '\nAs the human operator, do you authorize this run? Type yes: '
  read -r answer
  [ "$answer" = yes ] || die "not authorized; nothing was changed"
fi
authorized_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Everything from here on, including the machine's output, also goes to a log
# on the Mac, and a failing step names itself and shows the log's tail.
mkdir -p "$evidence_root"
log="$evidence_root/orbstack-amd64-$stamp.log"
exec > >(tee -a "$log") 2>&1
step="start"
on_error() {
  local status=$? line=$1 command=$2
  trap - ERR
  printf '\nerror: step "%s" failed (exit %s) at line %s: %s\n' "$step" "$status" "$line" "$command"
  printf 'Full log: %s\nRerun with --keep-machine to inspect it with: orb -m %s\n' "$log" "$machine"
}
trap 'on_error $LINENO "$BASH_COMMAND"' ERR

candidate_ref=refs/aibox-v1-05/orbstack-candidate
cleanup() {
  git -C "$repo" update-ref -d "$candidate_ref" 2>/dev/null || true
  if [ "$keep" = 1 ]; then
    echo "Machine kept: $machine (remove with: orb delete -f $machine)"
  else
    orb delete -f "$machine" >/dev/null 2>&1 && echo "Deleted machine $machine" || true
  fi
}
trap cleanup EXIT

say "Creating $machine"
orb create --arch amd64 debian:trixie "$machine"
in_root() { orb -m "$machine" -u root "$@"; }
in_user() { orb -m "$machine" "$@"; }

arch=$(in_user uname -m)
[ "$arch" = x86_64 ] || die "machine reports $arch, not x86_64"

say "Installing Docker, Git and Python"
in_root sh -euxc '
  export DEBIAN_FRONTEND=noninteractive
  apt-get update
  apt-cache policy docker.io docker-buildx
  apt-get install -y --no-install-recommends ca-certificates curl git python3 xz-utils \
    iptables docker.io docker-buildx
  systemctl enable --now docker || service docker start
  for attempt in 1 2 3 4 5 6 7 8 9 10; do docker info >/dev/null 2>&1 && exit 0; sleep 3; done
  docker info
'
linux_user=$(in_user id -un)
in_root usermod -aG docker "$linux_user"
say "Checking Docker in the machine"
in_user docker version
in_user docker buildx version
server_arch=$(in_user docker version --format '{{.Server.Arch}}')
[ "$server_arch" = amd64 ] || die "Docker server architecture is $server_arch, not amd64"

say "Cloning $full_commit into the machine"
# OrbStack shares the Mac's /Users with its machines, so the fetch needs no
# network. A temporary ref makes the commit fetchable even when only a
# remote-tracking branch holds it; cleanup removes the ref.
git -C "$repo" update-ref "$candidate_ref" "$full_commit"
in_user sh -euc "
  rm -rf \"\$HOME/aibox\"
  git init --quiet \"\$HOME/aibox\"
  cd \"\$HOME/aibox\"
  git fetch --quiet '$repo' '$candidate_ref'
  git checkout --quiet --detach FETCH_HEAD
  test \"\$(git rev-parse HEAD)\" = '$full_commit'
  test -z \"\$(git status --porcelain --untracked-files=all)\"
"

host_json() {
  # Operator and host identity for chapter 14; written next to the runner's evidence.
  local run=$1 status=$2 directory=$3
  python3 - "$directory/operator.json" <<PY
import json, sys
json.dump({
  "phase": "V1-05", "suite": "AC-TOOLS acceptance", "target": "debian-13-amd64",
  "run": $run, "runs": $runs, "exitStatus": $status,
  "operator": {"user": "$operator", "role": "human-authorized operator",
               "authorizedAt": "$authorized_at"},
  "sourceCommit": "$full_commit",
  "runnerOptions": $(python3 -c 'import json,sys; print(json.dumps(sys.argv[1:]))' "${runner_args[@]+"${runner_args[@]}"}"),
  "mac": {"model": "$(sysctl -n hw.model)", "chip": "$(sysctl -n machdep.cpu.brand_string)",
          "macos": "$(sw_vers -productVersion)", "memoryBytes": $(sysctl -n hw.memsize),
          "cpus": $(sysctl -n hw.ncpu)},
  "orbstack": {"version": "$(orb version 2>/dev/null | head -1)", "machine": "$machine",
               "disposable": True, "deletedAfterRun": $([ $keep = 1 ] && echo False || echo True)},
  "machine": {"unameMachine": "$arch", "kernel": "$(in_user uname -r)",
              "os": "$(in_user sh -c '. /etc/os-release; echo $PRETTY_NAME')",
              "cpus": $(in_user nproc),
              "translation": "x86_64 on Apple silicon through OrbStack (Rosetta when enabled, else QEMU)"},
  "docker": {"serverVersion": "$(in_user docker version --format '{{.Server.Version}}')",
             "serverArch": "$server_arch",
             "buildx": "$(in_user docker buildx version 2>/dev/null | head -1)"},
}, open(sys.argv[1], "w"), indent=2)
PY
}

quoted_args=""
for argument in "${runner_args[@]+"${runner_args[@]}"}"; do
  quoted_args="$quoted_args '$(printf '%s' "$argument" | sed "s/'/'\\\\''/g")'"
done

overall=0
run=1
while [ "$run" -le "$runs" ]; do
  run_dir="debian-13-amd64-orbstack-$(date -u +%Y%m%dT%H%M%SZ)-run$run"
  say "Run $run of $runs: $run_dir"
  if in_user sh -c "cd \"\$HOME/aibox\" && tests/acceptance/v1-05/run-host.sh debian-13-amd64 \
    --output \"\$HOME/evidence/$run_dir\"$quoted_args"; then
    status=0
  else
    status=$?
  fi
  mkdir -p "$evidence_root/$run_dir"
  if in_user sh -c "test -d \"\$HOME/evidence/$run_dir\""; then
    in_user sh -c "tar -C \"\$HOME/evidence/$run_dir\" -cf - ." | tar -C "$evidence_root/$run_dir" -xf -
  else
    echo "warning: the runner wrote no evidence directory" >&2
  fi
  host_json "$run" "$status" "$evidence_root/$run_dir"
  echo "Run $run exit status $status; evidence: $evidence_root/$run_dir"
  [ "$status" = 0 ] || overall=$status
  run=$((run + 1))
done

say "Done"
echo "Evidence directories are under $evidence_root/ (one per run, each with operator.json)."
echo "Send them to the integration owner as a whole. Overall exit status: $overall"
exit "$overall"

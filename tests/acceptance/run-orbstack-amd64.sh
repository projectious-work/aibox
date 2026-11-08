#!/bin/bash
# One phase's acceptance on debian-13-amd64 from an Apple silicon Mac with
# OrbStack.
#
#   tests/acceptance/run-orbstack-amd64.sh PHASE [--commit SHA] [--runs N]
#       [--evidence-dir ABS_DIR] [--keep-machine] [--yes] [-- runner options]
#
# PHASE is a roadmap phase such as V1-05. The phase provides
# tests/acceptance/<phase>/run-host.sh TARGET --output ABS_DIR [runner options]
# (V1-05's is the reference). The script creates a disposable amd64 Debian 13
# OrbStack machine (x86_64 under OrbStack's Rosetta or QEMU translation),
# installs Debian's Docker in it, clones this repository into it at the
# candidate commit (default: this checkout's HEAD), runs the phase's
# run-host.sh debian-13-amd64 N times, copies each evidence directory to
# the evidence folder (default $HOME/aibox-<phase>-evidence) with an
# operator.json that names the operator, the Mac and the machine, and deletes
# the machine. All output also goes to orbstack-amd64-<UTC time>.log in the
# evidence folder; a failing step names itself there and on the terminal.
#
# Needs: macOS on Apple silicon, OrbStack with the `orb` CLI, about 80 GB free
# disk, and the network access the phase's handoff lists. No credentials.
# Runs under translation are slower than native amd64 and are recorded as
# translated in operator.json.
set -Eeuo pipefail

case "${1:-}" in
  [Vv][0-9]*-[0-9]*) phase=$(printf '%s' "$1" | tr '[:lower:]' '[:upper:]'); shift ;;
  -h|--help) sed -n 2,25p "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
  *) echo "usage: $0 PHASE (e.g. V1-05) [options]; see --help" >&2; exit 2 ;;
esac
phase_dir=$(printf '%s' "$phase" | tr '[:upper:]' '[:lower:]')
commit=HEAD
runs=1
evidence_root="$HOME/aibox-$phase_dir-evidence"
keep=0
assume_yes=0
runner_args=()
while [ $# -gt 0 ]; do
  case "$1" in
    --commit) commit=${2:?--commit needs a SHA}; shift 2 ;;
    --runs) runs=${2:?--runs needs a number}; shift 2 ;;
    --evidence-dir) evidence_root=${2:?--evidence-dir needs a path}; shift 2 ;;
    --keep-machine) keep=1; shift ;;
    --yes) assume_yes=1; shift ;;
    --) shift; runner_args=("$@"); break ;;
    *) echo "unknown option: $1 (runner options go after --)" >&2; exit 2 ;;
  esac
done
case "$evidence_root" in /*) ;; *) echo "--evidence-dir must be absolute" >&2; exit 2 ;; esac
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
host_script="tests/acceptance/$phase_dir/run-host.sh"
git -C "$repo" cat-file -e "$full_commit:$host_script" 2>/dev/null \
  || die "$phase has no $host_script at ${full_commit:0:12}"

free_gb=$(df -g "$HOME" | awk 'NR==2 {print $4}')
[ "${free_gb:-0}" -ge 80 ] || die "need about 80 GB free on the Mac disk; ${free_gb} GB free"

stamp=$(date -u +%Y%m%dT%H%M%SZ)
machine="aibox-$phase_dir-amd64-$(echo "$stamp" | tr 'TZ' 'tz')"
operator=$(id -un)

cat <<EOF

$phase acceptance host run on debian-13-amd64 (OrbStack, translated amd64)

  Operator:   $operator
  Commit:     $full_commit
  Host test:  $host_script debian-13-amd64 ${runner_args[*]:-}
  Runs:       $runs
  Machine:    $machine (amd64 Debian 13, created now, deleted at the end$( [ $keep = 1 ] && echo ' unless kept: --keep-machine'))
  Effects:    installs docker.io, docker-buildx, git and python3 in that machine;
              inside it, the phase's host test pulls, builds and runs
              containers as its handoff in dev-notes/ describes. Nothing
              changes on the Mac except the new evidence under $evidence_root/.
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

candidate_ref=refs/aibox-acceptance/orbstack-candidate-$phase_dir
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
  "phase": "$phase", "hostScript": "$host_script", "target": "debian-13-amd64",
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
  if in_user sh -c "cd \"\$HOME/aibox\" && $host_script debian-13-amd64 \
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

#!/bin/bash
# Run one phase's amd64 acceptance from your own aibox checkout on the Mac:
#
#   run-amd64-acceptance.sh PHASE [--branch BRANCH] [options for run-orbstack-amd64.sh]
#
# PHASE is a roadmap phase such as V1-05. Run it from the top folder of your
# aibox checkout (the repository keeps it as tests/acceptance/run-mac-amd64.sh).
# Your checkout itself is not changed: the script creates a temporary clean
# checkout one level up (../aibox-<phase>-amd64-run-<UTC time>) at the phase's
# integration branch (default: the single origin/integration/<phase>-* branch),
# runs tests/acceptance/run-orbstack-amd64.sh PHASE from there, writes all
# evidence and the log to ./amd64-evidence/<phase>/<UTC time>/ in your
# checkout, and removes the temporary checkout at the end, whether the run
# passed or not. Options such as --runs 2, --keep-machine, --commit SHA or
# -- --scope full are passed on.
set -Eeuo pipefail

die() { printf 'error: %s\n' "$*" >&2; exit 1; }
case "${1:-}" in
  [Vv][0-9]*-[0-9]*) phase=$(printf '%s' "$1" | tr '[:lower:]' '[:upper:]'); shift ;;
  -h|--help) sed -n 2,15p "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
  *) die "usage: $0 PHASE (e.g. V1-05) [--branch BRANCH] [options]" ;;
esac
phase_dir=$(printf '%s' "$phase" | tr '[:upper:]' '[:lower:]')
branch=
if [ "${1:-}" = --branch ]; then branch=${2:?--branch needs a name}; shift 2; fi

[ "$(uname -s)" = Darwin ] || die "run this on the Mac"
repo=$(git rev-parse --show-toplevel 2>/dev/null) || die "run this from your aibox checkout"
[ "$repo" = "$(pwd -P)" ] || die "run this from the top folder of the checkout: $repo"

echo "Fetching origin"
git -C "$repo" fetch --quiet origin
if [ -z "$branch" ]; then
  candidates=$(git -C "$repo" for-each-ref --format='%(refname:lstrip=3)' \
    "refs/remotes/origin/integration/$phase_dir-*")
  case "$(printf '%s' "$candidates" | grep -c .)" in
    1) branch=$candidates ;;
    0) die "no origin/integration/$phase_dir-* branch; pass --branch" ;;
    *) die "several integration branches for $phase; pass --branch:
$candidates" ;;
  esac
fi
git -C "$repo" rev-parse --verify --quiet "origin/$branch^{commit}" >/dev/null \
  || die "origin/$branch does not exist"

stamp=$(date -u +%Y%m%dT%H%M%SZ)
run_checkout="$(dirname "$repo")/aibox-$phase_dir-amd64-run-$stamp"
evidence="$repo/amd64-evidence/$phase_dir/$stamp"
[ ! -e "$run_checkout" ] || die "$run_checkout already exists"

# An earlier manual attempt may have left a checkout inside this one; remove it if clean.
old="$repo/aibox-$phase_dir-amd64"
if [ -e "$old/.git" ]; then
  if git -C "$repo" worktree remove "$old" 2>/dev/null; then
    echo "Removed the earlier temporary checkout $old"
  else
    echo "warning: $old has changes; left in place (git worktree remove --force $old)" >&2
  fi
fi

remove_checkout() {
  if [ -d "$run_checkout" ]; then
    # No `git worktree prune`: it would also drop worktrees that live in the
    # development container, whose paths do not exist on the Mac.
    git -C "$repo" worktree remove --force "$run_checkout" 2>/dev/null || rm -rf "$run_checkout"
    echo "Removed the temporary checkout $run_checkout"
  fi
}
trap remove_checkout EXIT

echo "Creating the temporary checkout $run_checkout at origin/$branch"
git -C "$repo" worktree add --quiet --detach "$run_checkout" "origin/$branch"
[ -x "$run_checkout/tests/acceptance/run-orbstack-amd64.sh" ] \
  || die "origin/$branch has no tests/acceptance/run-orbstack-amd64.sh"
mkdir -p "$evidence"

# The Mac stays awake while this script runs. caffeinate watches this script
# instead of wrapping the run, so Ctrl-C reaches the run directly: it asks the
# runner in the machine to stop and still copies the evidence. This script
# waits for that and only then removes the temporary checkout.
caffeinate -i -w $$ &
trap 'echo "Interrupt: waiting for the run to stop and copy its evidence"' INT
status=0
( cd "$run_checkout" &&
  tests/acceptance/run-orbstack-amd64.sh "$phase" --evidence-dir "$evidence" "$@" ) \
  || status=$?
trap - INT

echo
echo "Exit status: $status"
echo "Evidence and log: $evidence"
exit "$status"

#!/bin/bash
# Run the V1-05 amd64 acceptance from your own aibox checkout on the Mac:
#
#   tests/acceptance/v1-05/run-mac-amd64.sh [--branch BRANCH] [options for run-orbstack-amd64.sh]
#
# A copy may sit at the checkout root as run-v1-05-amd64.sh. Run it from the
# top folder of your aibox checkout. Your checkout itself is
# not changed: the script creates a temporary clean checkout one level up
# (../aibox-v1-05-amd64-run-<UTC time>) at origin/integration/v1-05-tool-catalog,
# runs tests/acceptance/v1-05/run-orbstack-amd64.sh from there, writes all
# evidence and the log to ./amd64-evidence/<UTC time>/ in your checkout, and
# removes the temporary checkout at the end, whether the run passed or not.
# Options such as --runs 2, --keep-machine or -- --scope full are passed on.
set -Eeuo pipefail

branch=integration/v1-05-tool-catalog
if [ "${1:-}" = --branch ]; then branch=${2:?--branch needs a name}; shift 2; fi

die() { printf 'error: %s\n' "$*" >&2; exit 1; }
[ "$(uname -s)" = Darwin ] || die "run this on the Mac"
repo=$(git rev-parse --show-toplevel 2>/dev/null) || die "run this from your aibox checkout"
[ "$repo" = "$(pwd -P)" ] || die "run this from the top folder of the checkout: $repo"

stamp=$(date -u +%Y%m%dT%H%M%SZ)
parent=$(dirname "$repo")
run_checkout="$parent/aibox-v1-05-amd64-run-$stamp"
evidence="$repo/amd64-evidence/$stamp"
[ ! -e "$run_checkout" ] || die "$run_checkout already exists"

# An earlier attempt left a checkout inside this one; remove it if it is clean.
old="$repo/aibox-v1-05-amd64"
if [ -e "$old/.git" ]; then
  if git -C "$repo" worktree remove "$old" 2>/dev/null; then
    echo "Removed the earlier temporary checkout $old"
  else
    echo "warning: $old has changes; left in place (git worktree remove --force $old)" >&2
  fi
fi

echo "Fetching origin"
git -C "$repo" fetch --quiet origin
git -C "$repo" rev-parse --verify --quiet "origin/$branch^{commit}" >/dev/null \
  || die "origin/$branch does not exist"

remove_checkout() {
  if [ -d "$run_checkout" ]; then
    git -C "$repo" worktree remove --force "$run_checkout" 2>/dev/null || rm -rf "$run_checkout"
    git -C "$repo" worktree prune
    echo "Removed the temporary checkout $run_checkout"
  fi
}
trap remove_checkout EXIT

echo "Creating the temporary checkout $run_checkout"
git -C "$repo" worktree add --quiet --detach "$run_checkout" "origin/$branch"
mkdir -p "$evidence"

status=0
( cd "$run_checkout" &&
  caffeinate -i tests/acceptance/v1-05/run-orbstack-amd64.sh --evidence-dir "$evidence" "$@" ) \
  || status=$?

echo
echo "Exit status: $status"
echo "Evidence and log: $evidence"
exit "$status"

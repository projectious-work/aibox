#!/usr/bin/env bash
# check-src-context-drift.sh
# Verify no drift between src/context/ (upstream reference) and context/ (live installation)
# Part of pk-doctor health checks.

set -e

SRC_CONTEXT="${PROJECT_ROOT:-$(pwd)}/src/context"
LIVE_CONTEXT="${PROJECT_ROOT:-$(pwd)}/context"

if [[ ! -d "$SRC_CONTEXT" || ! -d "$LIVE_CONTEXT" ]]; then
  echo "SKIP: src/context/ not found (this is a consumer project, not processkit)" >&2
  exit 0
fi

# Find files in src/context (processkit reference)
# that have no counterpart or differ from live context/
# This catches: local overrides, deletions, and upstream-new files not yet merged.

DRIFT_COUNT=0

while IFS= read -r src_file; do
  rel_path="${src_file#$SRC_CONTEXT/}"
  live_file="$LIVE_CONTEXT/$rel_path"

  if [[ ! -e "$live_file" ]]; then
    echo "WARN: src/context/$rel_path exists but live context/$rel_path does not (upstream addition not installed)"
    ((DRIFT_COUNT++))
  elif ! diff -q "$src_file" "$live_file" &>/dev/null; then
    echo "WARN: src/context/$rel_path differs from context/$rel_path (local drift detected)"
    ((DRIFT_COUNT++))
  fi
done < <(find "$SRC_CONTEXT" -type f)

if [[ $DRIFT_COUNT -gt 0 ]]; then
  echo "WARN: $DRIFT_COUNT drift issue(s) detected" >&2
  exit 0  # warn, don't fail
fi

echo "OK: src/context/ and context/ in sync"
exit 0

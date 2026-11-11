#!/usr/bin/env bash
# Configure the long-lived aibox branches (main and the LTS lines) without
# GitHub Actions. Linear history enforces squash merges.
set -euo pipefail

repo="${AIBOX_GITHUB_REPO:-projectious-work/aibox}"
branches=(main lts/v0)

command -v gh >/dev/null 2>&1 || { echo "gh CLI is required" >&2; exit 1; }

# PR-only with resolved conversations, but no hosted status checks or mandatory
# second reviewer: releases are maintainer-run locally.
payload='{"required_status_checks":null,"enforce_admins":true,"required_pull_request_reviews":{"dismissal_restrictions":{"users":[],"teams":[],"apps":[]},"dismiss_stale_reviews":true,"require_code_owner_reviews":false,"required_approving_review_count":0,"require_last_push_approval":false},"restrictions":null,"required_linear_history":true,"allow_force_pushes":false,"allow_deletions":false,"block_creations":false,"required_conversation_resolution":true,"lock_branch":false,"allow_fork_syncing":false}'

for branch in "${branches[@]}"; do
  if ! gh api "repos/${repo}/branches/${branch}" >/dev/null 2>&1; then
    printf 'Skipping absent future branch %s\n' "${branch}"
    continue
  fi
  printf 'Protecting %s\n' "${branch}"
  gh api --method PUT "repos/${repo}/branches/${branch}/protection" --input - \
    <<<"${payload}" >/dev/null
done

echo "Protected all existing long-lived branches; no GitHub Actions or status checks configured."

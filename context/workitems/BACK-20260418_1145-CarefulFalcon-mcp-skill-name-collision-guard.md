---
apiVersion: processkit.projectious.work/v1
kind: WorkItem
metadata:
  id: BACK-20260418_1145-CarefulFalcon-mcp-skill-name-collision-guard
  created: '2026-04-18T11:45:00Z'
  labels:
    area: aibox-cli
    related_pr: aibox#53
spec:
  title: MCP/commands walker — defensive guard against bare-skill-name collisions across categories
  state: backlog
  type: task
  priority: medium
  description: |
    Both `cli/src/mcp_registration.rs::collect_processkit_mcp_specs` and
    `cli/src/claude_commands.rs::collect_live_commands` key on the bare
    skill directory name (e.g. `skill-gate`), not the fully-qualified
    `<category>/<skill>` pair. The current processkit catalog has no
    duplicates, so this is safe today. However, if a future processkit
    release ever ships two skills with the same basename in different
    categories (e.g. `engineering/release-semver` AND
    `devops/release-semver`), the second one silently overwrites the
    first in the merged MCP config and the `.claude/commands/` install.

    Q2 of the aibox#53 architect's plan flagged this. Per project owner
    (2026-04-18), the immediate v0.18.6 fix should add a defensive
    collision-detection check that warns (or errors) when two skills
    with the same basename are encountered during the walker pass.

    Acceptance:
    - Walker detects and reports duplicate skill basenames across categories.
    - Decision needed: warn-and-continue (last-wins) vs hard-fail. Default
      to warn-and-continue with a single-line `output::warn(...)` so existing
      installs don't break, but log enough info that the operator can
      disambiguate (both category paths printed).
    - Test: synth two skills `engineering/foo/mcp/mcp-config.json` and
      `devops/foo/mcp/mcp-config.json`, assert warning is emitted and the
      first-wins (or last-wins, whichever is decided) outcome is deterministic.
    - Apply the same guard symmetrically in `claude_commands.rs`.
---

# Notes

Discovered while implementing aibox#53 (per-skill MCP config merge).
Original architect plan in `LOG-20260418_0900-ClearHarbor-session-handover.md`
follow-up; user-approved guard scope on 2026-04-18.

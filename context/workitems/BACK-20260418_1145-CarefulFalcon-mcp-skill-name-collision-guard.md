---
apiVersion: processkit.projectious.work/v1
kind: WorkItem
metadata:
  id: BACK-20260418_1145-CarefulFalcon-mcp-skill-name-collision-guard
  created: '2026-04-18T11:45:00Z'
  labels:
    area: aibox-cli
    related_pr: aibox#53
  updated: '2026-04-19T08:08:10+00:00'
spec:
  title: MCP/commands walker — defensive guard against bare-skill-name collisions
    across categories
  state: done
  type: task
  priority: medium
  description: "Both `cli/src/mcp_registration.rs::collect_processkit_mcp_specs` and\n\
    `cli/src/claude_commands.rs::collect_live_commands` key on the bare\nskill directory\
    \ name (e.g. `skill-gate`), not the fully-qualified\n`<category>/<skill>` pair.\
    \ The current processkit catalog has no\nduplicates, so this is safe today. However,\
    \ if a future processkit\nrelease ever ships two skills with the same basename\
    \ in different\ncategories (e.g. `engineering/release-semver` AND\n`devops/release-semver`),\
    \ the second one silently overwrites the\nfirst in the merged MCP config and the\
    \ `.claude/commands/` install.\n\nQ2 of the aibox#53 architect's plan flagged\
    \ this. Per project owner\n(2026-04-18), the immediate v0.18.6 fix should add\
    \ a defensive\ncollision-detection check that warns (or errors) when two skills\n\
    with the same basename are encountered during the walker pass.\n\nAcceptance:\n\
    - Walker detects and reports duplicate skill basenames across categories.\n- Decision\
    \ needed: warn-and-continue (last-wins) vs hard-fail. Default\n  to warn-and-continue\
    \ with a single-line `output::warn(...)` so existing\n  installs don't break,\
    \ but log enough info that the operator can\n  disambiguate (both category paths\
    \ printed).\n- Test: synth two skills `engineering/foo/mcp/mcp-config.json` and\n\
    \  `devops/foo/mcp/mcp-config.json`, assert warning is emitted and the\n  first-wins\
    \ (or last-wins, whichever is decided) outcome is deterministic.\n- Apply the\
    \ same guard symmetrically in `claude_commands.rs`."
  started_at: '2026-04-19T08:08:04+00:00'
  completed_at: '2026-04-19T08:08:10+00:00'
---

# Notes

Discovered while implementing aibox#53 (per-skill MCP config merge).
Original architect plan in `LOG-20260418_0900-ClearHarbor-session-handover.md`
follow-up; user-approved guard scope on 2026-04-18.


## Transition note (2026-04-19T08:08:04+00:00)

Guard already implemented in cli/src/mcp_registration.rs and claude_commands.rs; moving through lifecycle to close.


## Transition note (2026-04-19T08:08:10+00:00)

Guard shipped in v0.18.6 (aibox#53); decision DEC-20260419_0807-LuckyDaisy accepted.

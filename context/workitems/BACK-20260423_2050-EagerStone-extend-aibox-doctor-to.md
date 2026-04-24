---
apiVersion: processkit.projectious.work/v1
kind: WorkItem
metadata:
  id: BACK-20260423_2050-EagerStone-extend-aibox-doctor-to
  created: '2026-04-23T20:50:38+00:00'
  updated: '2026-04-24T00:04:13+00:00'
spec:
  title: Extend aibox doctor to detect missing command registrations
  state: review
  type: task
  priority: medium
  description: Add a check to `aibox doctor` (extends BACK-020) that validates every
    installed processkit skill has its command files registered in `.claude/commands/`.
    Should warn if `context/skills/processkit/<skill>/commands/<cmd>.md` exists but
    `.claude/commands/<cmd>.md` is missing. Helps derived projects detect the v0.19.1
    pk-doctor issue before it causes problems.
  started_at: '2026-04-24T00:04:09+00:00'
---

## Transition note (2026-04-24T00:04:13+00:00)

Implemented command registration detection in aibox doctor. Commit fcd25f2 adds check_command_registrations() that validates all installed skills have their command files registered in .claude/commands/. All tests pass, no clippy warnings.

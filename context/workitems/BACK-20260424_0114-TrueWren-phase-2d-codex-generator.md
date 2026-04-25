---
apiVersion: processkit.projectious.work/v1
kind: WorkItem
metadata:
  id: BACK-20260424_0114-TrueWren-phase-2d-codex-generator
  created: '2026-04-24T01:14:34+00:00'
  updated: '2026-04-24T07:47:48+00:00'
spec:
  title: 'Phase 2d: Codex generator (trust_level + fallback)'
  state: done
  type: task
  priority: medium
  description: "**Scope:** Implement simplest generator for Codex (project-scoped\
    \ trust_level).\n\n**Deliverables:**\n1. `generate_codex_permissions(config: &McpConfig)\
    \ -> Result<()>`\n   - Set `trust_level = \"trusted\"` in `.codex/config.toml`\n\
    \   - As fallback (since Codex doesn't support per-tool granularity), also generate\
    \ allowlist if per-harness overrides specify it\n\n**Note:** Codex is least flexible;\
    \ trust_level is project-scoped, not per-tool. Document this limitation.\n\n**Estimated\
    \ Tokens:** ~1.5K (simplest generator; minimal config)\n**Can run in parallel:**\
    \ with Phase 2a, 2b, 2c"
  parent: BACK-20260424_0058-ToughGrove-feature-global-mcp-permissions
  blocked_by:
  - BACK-20260424_0114-JollyStream-phase-1-core-mcp
  started_at: '2026-04-24T07:47:48+00:00'
  completed_at: '2026-04-24T07:47:48+00:00'
---

## Transition note (2026-04-24T07:47:48+00:00)

Codex generator implemented with trust_level and fallback allowed_tools list.


## Transition note (2026-04-24T07:47:48+00:00)

Phase 2d complete. All 8 harness generators implemented. Phase 3 integration ready.

---
apiVersion: processkit.projectious.work/v1
kind: WorkItem
metadata:
  id: BACK-20260424_0114-JollyDew-phase-3-integration-into
  created: '2026-04-24T01:14:37+00:00'
  updated: '2026-04-24T01:16:30+00:00'
spec:
  title: 'Phase 3: Integration into seed.rs and container lifecycle'
  state: backlog
  type: task
  priority: high
  description: '**Scope:** Wire all 8 harness generators into the container seed pipeline.


    **Deliverables:**

    1. Parse `[mcp]` section from aibox.toml into `McpConfig`

    2. Add `generate_all_harness_permissions(config: &McpConfig)` orchestrator function

    3. Call orchestrator at end of `seed_runtime_config()` in seed.rs

    4. Add error logging (warn on pattern match failures, continue on generator errors)

    5. Ensure merge strategy: never overwrite existing user edits to harness configs


    **Blocker:** Cannot start until Phase 1 + all Phase 2 subtasks are complete


    **Estimated Tokens:** ~2K (orchestrator function, error handling, lifecycle integration)

    **Blocking:** Phase 4 (documentation must reference implemented behavior)'
  parent: BACK-20260424_0058-ToughGrove-feature-global-mcp-permissions
  blocked_by:
  - BACK-20260424_0114-SwiftPlum-phase-2a-claude-code
  - BACK-20260424_0114-NobleSage-phase-2b-continue-cursor
  - BACK-20260424_0114-ToughGlade-phase-2c-gemini-cli
  - BACK-20260424_0114-TrueWren-phase-2d-codex-generator
---

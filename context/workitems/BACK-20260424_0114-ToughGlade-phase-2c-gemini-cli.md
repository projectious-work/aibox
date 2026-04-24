---
apiVersion: processkit.projectious.work/v1
kind: WorkItem
metadata:
  id: BACK-20260424_0114-ToughGlade-phase-2c-gemini-cli
  created: '2026-04-24T01:14:32+00:00'
  updated: '2026-04-24T01:16:29+00:00'
spec:
  title: 'Phase 2c: Gemini CLI + GitHub Copilot generators'
  state: backlog
  type: task
  priority: high
  description: "**Scope:** Implement policy-based generators for Gemini CLI and GitHub\
    \ Copilot.\n\n**Deliverables:**\n1. `generate_gemini_permissions(config: &McpConfig)\
    \ -> Result<()>`\n   - Generate dual `includeTools[]` (allowlist) and `excludeTools[]`\
    \ (blocklist)\n   - Apply intersection semantics: a tool is allowed if in includeTools\
    \ AND not in excludeTools\n2. `generate_github_copilot_permissions(config: &McpConfig)\
    \ -> Result<()>`\n   - Generate CLI flags: `--allow-tool`, `--deny-tool`\n   -\
    \ Document enterprise registry policy integration\n\n**Key Difference:** Both\
    \ support dual allow/deny; Gemini uses intersection, Copilot uses command-line\
    \ flags.\n\n**Estimated Tokens:** ~2.5K per harness (dual list logic, flag generation)\n\
    **Can run in parallel:** with Phase 2a, 2b, 2d"
  parent: BACK-20260424_0058-ToughGrove-feature-global-mcp-permissions
  blocked_by:
  - BACK-20260424_0114-JollyStream-phase-1-core-mcp
---

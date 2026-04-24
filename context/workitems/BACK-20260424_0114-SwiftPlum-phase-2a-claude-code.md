---
apiVersion: processkit.projectious.work/v1
kind: WorkItem
metadata:
  id: BACK-20260424_0114-SwiftPlum-phase-2a-claude-code
  created: '2026-04-24T01:14:28+00:00'
  updated: '2026-04-24T01:16:29+00:00'
spec:
  title: 'Phase 2a: Claude Code + OpenCode generators'
  state: backlog
  type: task
  priority: high
  description: "**Scope:** Implement MCP permission generators for Claude Code and\
    \ OpenCode (parallel with other harnesses).\n\n**Deliverables:**\n1. `generate_claude_code_permissions(config:\
    \ &McpConfig) -> Result<()>` \n   - Update `.claude/settings.local.json` with\
    \ `permissions.allow[]` list\n   - Merge with existing settings (don't overwrite)\n\
    2. `generate_opencode_permissions(config: &McpConfig) -> Result<()>`\n   - Generate\
    \ `.opencode/config.toml` with `[mcp]` section\n   - Set `allow/ask/deny` per\
    \ pattern with wildcard expansion\n\n**Key:** Both use allowlist format; OpenCode\
    \ additionally supports deny_patterns.\n\n**Estimated Tokens:** ~2K per harness\
    \ (file I/O, TOML/JSON generation, merge logic)\n**Can run in parallel:** with\
    \ Phase 2b, 2c, 2d"
  parent: BACK-20260424_0058-ToughGrove-feature-global-mcp-permissions
  blocked_by:
  - BACK-20260424_0114-JollyStream-phase-1-core-mcp
---

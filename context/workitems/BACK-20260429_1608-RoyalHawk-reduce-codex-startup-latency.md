---
apiVersion: processkit.projectious.work/v1
kind: WorkItem
metadata:
  id: BACK-20260429_1608-RoyalHawk-reduce-codex-startup-latency
  created: '2026-04-29T16:08:42+00:00'
  labels:
    component: codex
    area: mcp
    source: user-report
    reported_date: '2026-04-29'
spec:
  title: Reduce Codex startup latency from eager processkit MCP server initialization
  state: backlog
  type: bug
  priority: high
  description: 'Codex CLI 0.125.0 eagerly initializes every enabled stdio MCP server
    at session and subagent startup. aibox currently registers each processkit MCP
    server as a separate `uv run context/skills/processkit/<skill>/mcp/server.py`
    process in `.codex/config.toml`, so a fresh or cold-cache derived project can
    spend minutes at `Starting MCP servers`, often displaying the last pending server
    such as `processkit-index-management`. Direct `uv run` of index-management exits
    quickly when stdin closes; current warm-cache `codex debug prompt-input` took
    ~3.9s, while logs show 19 ListToolsRequest handshakes and earlier cold-cache launches
    downloaded dependencies repeatedly. Investigate and implement a provider-neutral
    mitigation: preferably a processkit aggregate MCP server/proxy or aibox-generated
    Codex MCP profile that avoids eager loading all per-skill servers while preserving
    tool access.'
---

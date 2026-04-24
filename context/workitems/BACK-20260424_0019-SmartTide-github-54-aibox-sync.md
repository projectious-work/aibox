---
apiVersion: processkit.projectious.work/v1
kind: WorkItem
metadata:
  id: BACK-20260424_0019-SmartTide-github-54-aibox-sync
  created: '2026-04-24T00:19:50+00:00'
spec:
  title: 'GitHub #54: aibox sync reconcile .mcp.json on per-skill-config drift'
  state: backlog
  type: task
  priority: medium
  description: 'Implement logic to detect and reconcile drift between per-skill mcp-config.json
    files and the merged .mcp.json at project root. This handles the case where a
    user updates a skill''s mcp-config.json directly (bypassing processkit update),
    and aibox sync should re-merge to pick up the changes.


    Depends on processkit v0.19.2 shipping a manifest contract to track which specs
    came from which skill. Currently blocked waiting for upstream manifest.


    Related: GitHub issue #54'
---

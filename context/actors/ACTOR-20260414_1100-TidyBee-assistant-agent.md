---
apiVersion: processkit.projectious.work/v1
kind: Actor
metadata:
  id: ACTOR-20260414_1100-TidyBee-assistant-agent
  created: 2026-04-14T11:00:00Z
spec:
  type: agent
  name: Assistant Agent (template)
  active: true
  joined_at: 2026-04-14T11:00:00Z
  handle: assistant-agent
  x_aibox:
    is_template: true
    model: claude-haiku-4-5-20251001
    model_tier: haiku
    role_ref: ROLE-20260414_1100-TidyBee-assistant
    clone_of: null
    schema_note: "Project-local extension fields; migrate to processkit canonical team schema when available."
---

Template actor for the assistant role. High-volume administrative work: morning briefings, summaries, indexing, handovers.

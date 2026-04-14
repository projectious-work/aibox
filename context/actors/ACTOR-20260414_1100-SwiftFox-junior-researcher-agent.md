---
apiVersion: processkit.projectious.work/v1
kind: Actor
metadata:
  id: ACTOR-20260414_1100-SwiftFox-junior-researcher-agent
  created: 2026-04-14T11:00:00Z
spec:
  type: agent
  name: Junior Researcher Agent (template)
  active: true
  joined_at: 2026-04-14T11:00:00Z
  handle: junior-researcher-agent
  x_aibox:
    is_template: true
    model: claude-sonnet-4-6
    model_tier: sonnet
    role_ref: ROLE-20260414_1100-SwiftFox-junior-researcher
    clone_of: null
    schema_note: "Project-local extension fields; migrate to processkit canonical team schema when available."
---

Template actor for the junior-researcher role. PM's default researcher choice.

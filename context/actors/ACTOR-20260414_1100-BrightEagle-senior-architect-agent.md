---
apiVersion: processkit.projectious.work/v1
kind: Actor
metadata:
  id: ACTOR-20260414_1100-BrightEagle-senior-architect-agent
  created: 2026-04-14T11:00:00Z
spec:
  type: agent
  name: Senior Architect Agent (template)
  active: true
  joined_at: 2026-04-14T11:00:00Z
  handle: senior-architect-agent
  x_aibox:
    is_template: true
    model: claude-opus-4-6
    model_tier: opus
    role_ref: ROLE-20260414_1100-BrightEagle-senior-architect
    clone_of: null
    schema_note: "Project-local extension fields; migrate to processkit canonical team schema when available."
---

Template actor for the senior-architect role. Invoked by PM for large, cross-cutting, or load-bearing design work.

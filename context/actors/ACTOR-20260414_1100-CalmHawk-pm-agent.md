---
apiVersion: processkit.projectious.work/v1
kind: Actor
metadata:
  id: ACTOR-20260414_1100-CalmHawk-pm-agent
  created: 2026-04-14T11:00:00Z
spec:
  type: agent
  name: PM Agent (template)
  active: true
  joined_at: 2026-04-14T11:00:00Z
  handle: pm-agent
  x_aibox:
    is_template: true
    model: claude-opus-4-6
    model_tier: opus
    role_ref: ROLE-20260414_1100-CalmHawk-project-manager
    clone_of: null
    schema_note: "Project-local extension fields; migrate to processkit canonical team schema when available."
---

Template actor for the project-manager role. The owner approves every plan this actor authors; this actor is the only team member that speaks directly to the owner in normal operation.

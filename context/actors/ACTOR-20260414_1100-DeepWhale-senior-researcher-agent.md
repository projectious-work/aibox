---
apiVersion: processkit.projectious.work/v1
kind: Actor
metadata:
  id: ACTOR-20260414_1100-DeepWhale-senior-researcher-agent
  created: 2026-04-14T11:00:00Z
spec:
  type: agent
  name: Senior Researcher Agent (template)
  active: true
  joined_at: 2026-04-14T11:00:00Z
  handle: senior-researcher-agent
  x_aibox:
    is_template: true
    model: claude-opus-4-6
    model_tier: opus
    role_ref: ROLE-20260414_1100-DeepWhale-senior-researcher
    clone_of: null
    schema_note: "Project-local extension fields; migrate to processkit canonical team schema when available."
---

Template actor for the senior-researcher role. Invoked by PM for research requiring Opus-tier synthesis.

# processkit archive

Read-only records of the processkit setup that this repository used for its
own development until 2026-11-11. See
[ADR 0001](../../docs/adr/0001-remove-processkit-from-the-development-setup.md).

- Do not edit or extend these files. They are history.
- Search them with ordinary tools, for example `rg CuriousBeacon archive/processkit`.
- The files moved here unchanged from `context/`. Paths inside them still say
  `context/…`; read them as `archive/processkit/…`.
- The schemas, skills, template mirrors and MCP configuration were deleted.
  The Git tag `processkit-final` holds the complete former installation.

| Directory | Content | Files |
|---|---|---|
| [`actors/`](actors/) | Actor entities (people and agents) | 9 |
| [`archive/`](archive/) | Records archived by processkit before the removal | 83 |
| [`archives/`](archives/) | Compressed cold archives written by context archiving | 24 |
| [`artifacts/`](artifacts/) | Artifacts: reports, model specifications, plans | 82 |
| [`bindings/`](bindings/) | Role and model bindings | 68 |
| [`decisions/`](decisions/) | Decision records (`DEC-*`) | 118 |
| [`discussions/`](discussions/) | Structured discussions | 6 |
| [`gates/`](gates/) | Gate definitions | 1 |
| [`logs/`](logs/) | Event log entries | 1016 |
| [`migrations/`](migrations/) | processkit upgrade migration records | 15 |
| [`notes/`](notes/) | Notes, including the v0 release process | 39 |
| [`roles/`](roles/) | Role definitions | 58 |
| [`scopes/`](scopes/) | Scopes (milestones, sprints) | 1 |
| [`team/`](team/) | Team roster | 2 |
| [`team-members/`](team-members/) | Team member identities and their memory | 43 |
| [`workitems/`](workitems/) | Work items (`BACK-*`) | 182 |

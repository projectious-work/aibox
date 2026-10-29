---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261029_0835-HumbleBridge-v1-05-owner-row-dispositions-first
  created: '2026-10-29T08:35:14+00:00'
spec:
  title: V1-05 owner row dispositions, first batch
  state: accepted
  decision: 'Owner dispositions for the V1-05 tool catalog: (1) ADDON:ai-opencode/opencode
    is retained but re-pointed from the archived opencode-ai/opencode (v0.0.55) to
    the active anomalyco/opencode (formerly sst/opencode) at an exact version. (2)
    ADDON:python/pip maps to the Python runtime''s bundled pip; v0 ''pip disabled
    means absent'' is not preserved and is documented as a compatibility note. (3)
    ADDON:infrastructure/podman is retained with a documented host-privilege requirement.
    (4) ADDON:ai-mistral/mistral (SDK), ADDON:docs-starlight/starlight (scaffolder),
    ADDON:docs-docusaurus/docusaurus (project dependency) and ADDON:docs-mkdocs/mkdocs
    (MkDocs 1.x unmaintained; Zensical reads mkdocs.yml) are removed. (5) go, go-quality
    and go-release merge into one semantic ''go'' group: go, goimports, staticcheck,
    golangci-lint, govulncheck, gosec, goreleaser plus an all bundle. Remaining groups
    follow the recommended shapes unless amended by later owner decisions on Python-via-uv
    and new groups (vibe, typescript, zig, security).'
  context: DEC-20261026_1412-FastForge requires owner row decisions before the V1-05
    contract freeze; 94 of 98 rows were owner-review-pending.
  rationale: Removed rows are libraries, scaffolders, project-pinned dependencies
    or unmaintained; opencode's v0 source is archived; a single go group gives one
    coherent selection surface.
  consequences: Converter must report removed rows and the pip compatibility change
    explicitly to v0 projects. The 98 census rows stay individually traceable in the
    selection table with a removed or retained disposition.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-29T08:35:14+00:00'
---

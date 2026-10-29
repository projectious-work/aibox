---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261029_2107-QuietSeal-replace-rich-cli-with-glow-and
  created: '2026-10-29T21:07:29+00:00'
spec:
  title: Replace rich-cli with glow and bat in V1-05 preview-enhanced
  state: accepted
  decision: 'Remove rich-cli and the preview/enhanced/rich Feature entirely. The census
    row ADDON:preview-enhanced/rich is removed. Add two tool Features to the preview-enhanced
    semantic group: preview/enhanced/glow (Go, Markdown rendering) and preview/enhanced/bat
    (Rust, syntax highlighting for JSON, RST and front matter). Both are installed
    from producer releases with checksums. The V1-13 preview rewrite uses glow and
    bat with no Python at run time.'
  context: 'The implementation track found that rich-cli 1.8.1 reports version 1.8.0,
    so the catalog check cannot pass. Investigation showed that v0''s rich row installed
    the python3-rich library, which only aibox''s own rich-preview.yazi plugin and
    aibox-preview.sh used: inline Python rendering .md, .rst, .json and .ipynb. With
    uv-managed Python in v1, no system python3 has rich, so the rich-cli Feature would
    not serve the previewer, and the plugin needs a rewrite in V1-13 anyway. Yazi
    plugins are Lua and call an external renderer.'
  rationale: Formatted document preview is the user outcome to preserve, not the rich
    tool. glow and bat are single static binaries with fast start-up, avoid a Python
    dependency in a core terminal feature, and bat was already part of the v0 base.
  alternatives:
  - option: Keep rich-cli and fix the version check
    rejected: keeps the Python start-up cost and still needs a plugin rewrite
  - option: Write our own Go renderer
    rejected: a new engine to maintain; the reuse rule prefers existing tools
  consequences: Contract amendment 2 changes the catalog, ledger, compatibility notes
    and smoke fixtures as needed. Notebook (.ipynb) preview needs jq-extracted cells
    or a documented gap in V1-13. The implementation track replaces the rich Feature;
    the test track adapts.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-29T21:07:29+00:00'
---

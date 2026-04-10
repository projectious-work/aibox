# aibox v0.17.12

Feature release: Yazi git status integration, complete OS gitignore entries, and
improved processkit migration guidance for agents.
No processkit version change — still compatible with v0.8.0.

## Yazi git status integration (`git.yazi`)

Yazi now shows git file status (modified, untracked, ignored, staged, deleted)
inline in the file list — similar to VS Code's file explorer. The official
`git.yazi` plugin from `yazi-rs/plugins` is bundled directly in the base image.

Signs displayed next to file names:
- ` ` modified
- ` ` added/staged
- ` ` deleted
- `? ` untracked
- ` ` ignored (dimmed)

Status propagates up to parent directories automatically. No configuration
needed — enabled by default in all aibox containers.

## Complete OS-specific gitignore entries

New `.gitignore` files created by `aibox init` now include entries for all
three major OS families:

**Linux** (new): `.Trash-*/`, `.fuse_tmp*`, `.directory`, `.nfs*`
**Windows** (expanded): added `Desktop.ini`, `$RECYCLE.BIN/`
**macOS**: unchanged

Existing `.gitignore` files are not modified (OS entries are informational
defaults, not enforced aibox entries).

## Template-snapshot diff guidance in migration documents

Migration briefings now explicitly instruct agents to diff the **old processkit
template snapshot** against the **new one** when reviewing template changes
(e.g. `AGENTS.md`), rather than diffing the customized installed file against
the new template.

Comparing installed (customized) vs new template produces noise from project
customizations that hides real upstream changes. The correct workflow:

```
diff context/templates/processkit/v{old}/AGENTS.md \
     context/templates/processkit/v{new}/AGENTS.md
```

Then apply only the upstream delta onto the customized installed file.

## AGENTS.md: runtime artifacts for agents

Added a "Runtime artifacts for agents" section documenting `.aibox/aibox.log`
(NDJSON structured log of every aibox command), `aibox.lock`, and
`context/migrations/` — with explicit guidance for agents to read the log to
understand recent aibox activity.

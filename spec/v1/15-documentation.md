# 15. Version-line documentation and prerelease publication

**R-DOCS:** aibox v1.x MUST have its own version-aligned public documentation
line, built with the same Hugo build pipeline and the same
`github.com/projectious-work/brand-theme-hugo-vanilla` theme module used by
v0.x. The current v1 prerelease tree uses Docsy, a separate submodule,
Bootstrap and Font Awesome; that is a migration input, not the target. Reuse
the v0 documentation shell, build scripts and brand components rather than
maintaining two site frameworks. The initial theme pin must match the v0
line's tested pin (currently v0.3.4 at the inspected v0.x release branch);
later changes are coordinated, pinned upgrades, never floating module heads.
The v1 content is independently maintained and describes the Dev Container
CLI-based product, not a copied v0 `aibox.toml` manual.

## Site layout and release identity

The stable product documentation root
`https://projectious-work.github.io/aibox/` continues to serve the current
v0.x release until a separately approved stable-v1 promotion. Its label in
the Releases selector and public metadata remains **Current: v0.x**.
The v1 line lives at `/aibox/v1.x/`, visibly marked **v1.x preview** while
only alpha/beta releases exist. Immutable prerelease documentation lives at
`/aibox/v1.x/v1.<minor>.<patch>-alpha.<n>/` and corresponding beta paths
for actual released tags. An alpha or beta may be the newest release *within
the v1 line*, but must not be labeled the current product release. The
root v0 page, its archives, and its current label are not overwritten by a
v1 docs deployment. Conversely, a v0 deployment must not erase v1 entries
from the shared release manifest.

The Releases dropdown is present and keyboard accessible on both lines and
contains: current v0.x root; v0.x release archives; v1.x preview line; and
actual published v1.x alpha/beta archives. Entries are generated from a
single validated release manifest or release metadata, not two manually
divergent menus. Unpublished versions are not presented as available.
The active entry uses `aria-current`, accessible names distinguish current
stable from latest prerelease, and the menu works without network-dependent
JavaScript through a Hugo-rendered fallback. Cross-version navigation tries
the equivalent page only when it exists; otherwise it links to that line's
documentation landing page with a visible explanation, never a 404 or a
silent redirect to content with different semantics. Canonical links, search
indexes, sitemap, `llms.txt`, Markdown projections and edit links remain
line/version-specific.

## Alpha/beta branch and content contract

The repository currently has `v1.x-dev` and `v1.x-pre-release` branches;
alpha and beta are prerelease versions/tags on the latter, not assumed
permanent branch names. Before each alpha or beta tag, update the docs source
on the branch from which that candidate is cut and build its immutable
snapshot from the exact candidate commit. If distinct alpha/beta branches are
introduced later, apply the same rule to each; do not publish docs built
from a different branch or a later mainline commit. A candidate release
cannot rely on an obsolete Docsy-only page or a v0-oriented instruction.

The v1 docs must cover: preview maturity and support limits; installation
of the pinned Dev Container CLI and aibox distribution; a direct-upstream
minimal project journey; standard `devcontainer.json` and Feature/tool
selection; local themes, tmux headers/titles, Yazi previews and persistence;
operator-versus-container-agent authority; MCP/CLI operations and how-to
resources; optional audio/processkit/LaTeX/sidebar/review; migration and
rollback from v0; diagnostics; target compatibility; security; and release
history. Commands and examples must run against the documented candidate
or be visibly marked planned. Planned roadmap content stays separate from
shipped instructions. Remove or revise stale v1-alpha claims about
`aibox.toml`, the reverted v1 architecture, obsolete processkit pins,
deployment backends and command syntax. Preserve historical pages only as
clearly versioned archives, not current getting-started guidance.

README, changelog, release notes, compatibility table, public Markdown guide,
Hugo pages, `llms.txt`, read-only product MCP resources and task-sized how-to
guidance describe the same candidate. Derive repeated version/URL/release
facts from a manifest where feasible and check links and examples against
the built site. Every user-facing phase changes its relevant stable docs in
the same change; a development note is not a substitute for user guidance.

## Build, deploy and acceptance

Adapt the v0 Hugo module build, including Hugo Extended, Go module resolution,
the pinned Node asset dependencies and theme-catalog generation where
applicable. Keep the project-local build and `gh-pages` publication path;
no GitHub Actions or new hosted docs service. The deploy process stages the
new line/current and immutable archive in a disposable worktree, validates
the release manifest, preserves all other line directories, checks changed
paths, then pushes only after local review. An interrupted deployment must
leave the existing public site intact or have a documented recovery path.

AC-DOCS requires clean builds of v0 current, v1 preview, one alpha archive
and one beta archive using the common theme and build route. Test both
deployment directions: v1 publication preserves v0 root/current/archive,
and v0 publication preserves v1 preview/archive and both manifest lines.
Check rendered labels and links from both sites, with JavaScript enabled and
disabled, keyboard navigation, active-state accessibility, nonexistent
same-path fallback, relative assets, search and canonical metadata. Check
the exact alpha/beta candidate docs against executable examples, schemas and
release notes; reject stale Docsy assets or stale v0 configuration claims.
No stable-v1/current label change occurs until the owner separately approves
the stable promotion and its documentation cutover.

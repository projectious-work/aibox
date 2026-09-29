# 15. Version-line documentation and prerelease publication

**R-DOCS:** aibox v1.x MUST have its own version-aligned public documentation
line, built with the same Hugo build pipeline and the same
`github.com/projectious-work/brand-theme-hugo-vanilla` theme module used by
v0.x. Both lines use one Hugo build and brand shell, not two site frameworks.
The theme version is pinned and coordinated across both lines; floating module
heads are forbidden.
The v1 content is independently maintained and describes the Dev Container
CLI-based product, not a copied v0 `aibox.toml` manual.

## Documentation is part of every phase

V1-02 establishes the shared v0/v1 Hugo build as a foundation prerequisite,
not a final documentation sprint. Every later phase depends directly or
transitively on it. The canonical roadmap gives every phase a
`docs` deliverable; no phase may reach `shipped` with that field unfulfilled.
Phase IDs follow this implementation order; the
[crosswalk](phase-id-crosswalk.md) preserves references from earlier reviews.

Each phase's implementation change MUST include the relevant user, operator,
integrator and/or maintainer documentation with the implementation. Update
the v1 preview pages, task examples, configuration/API reference, compatibility
notes, how-to/MCP resources, roadmap and changelog/release notes as applicable.
Remove or clearly label instructions for unavailable functionality. A phase
with substantial internal work still documents its public contract and
maintainer-facing architecture/evidence; from V1-03 it must also keep a
feature-specific executable customer demo. Do not publish a feature guide that
claims a feature is usable before its acceptance tests pass.

Beginning with V1-03, each phase's user-facing preview docs MUST make the
same candidate executable demonstrable: a runnable command, disposable
fixture or setup, expected output, supported targets, candidate limits and
cleanup/reset. A package-only phase is not an acceptable completed increment.
Each later phase extends this executable walkthrough and preserves prior
supported journeys. Documentation never represents unsupported operations as
usable.

The phase completion evidence MUST include: changed documentation paths and
their audience; the phase's `docs` deliverable; a successful local build of
the v1 Hugo site from the same candidate commit; link/schema/example checks
for changed pages; and a check that the v0 current build and Releases labels
remain intact. Build failures block phase completion. If a phase changes a
CLI/MCP/config surface, test the changed examples against that phase's
candidate binary or mark the unimplemented examples as planned. The build
must be rerun after any code or docs change affecting the candidate. Public
deployment remains a separate release/preview decision; a local build is
mandatory even when publication is deferred.
When publishing a Dev Container example, keep its active settings minimal and
its aibox-owned option catalog complete, commented and explanatory as required
by chapters 2 and 8. Reference the upstream Dev Container CLI for unrelated
native options instead of maintaining a second general-purpose reference.

V1-19 is not where feature documentation or preview publication first appears.
Earlier verified checkpoints may publish scoped alpha/beta releases with exact
candidate-matched docs and immutable archives. V1-19 expands the cross-line
alpha/beta snapshot matrix, shared Releases dropdown and preservation checks
over those earlier releases. V1-23 performs final candidate
and published-artifact consistency checks; neither phase may be used to defer
the documentation owed by intervening implementation phases.

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

The version line uses `v1.x-dev` and `v1.x-pre-release` branches;
alpha and beta are prerelease versions/tags on the latter, not assumed
permanent branch names. Before each alpha or beta tag, update the docs source
on the branch from which that candidate is cut and build its immutable
snapshot from the exact candidate commit. If distinct alpha/beta branches are
introduced later, apply the same rule to each; do not publish docs built
from a different branch or a later mainline commit. A candidate release
cannot rely on an obsolete Docsy-only page or a v0-oriented instruction.

At each verified preview checkpoint, the public v1 docs cover only the
candidate's available journeys as runnable instructions and label the
remaining roadmap capabilities planned or unsupported. The evolving v1 docs
must ultimately cover: preview maturity and support limits; installation
of the pinned Dev Container CLI and aibox distribution; an aibox-wrapped
minimal project journey and a separate direct-upstream interoperability test;
standard `devcontainer.json` and Feature/tool
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
the built site. Every phase changes its relevant docs in the same change; a
development note is not a substitute for user or operator guidance where
behavior changes.

## Build, deploy and acceptance

Adapt the v0 Hugo module build, including Hugo Extended, Go module resolution,
the pinned Node asset dependencies and theme-catalog generation where
applicable. Keep the project-local build and `gh-pages` publication path;
no GitHub Actions or new hosted docs service. The deploy process stages the
new line/current and immutable archive in a disposable worktree, validates
the release manifest, preserves all other line directories, checks changed
paths, then pushes only after local review. An interrupted deployment must
leave the existing public site intact or have a documented recovery path.

At an earlier scoped alpha/beta checkpoint, AC-DOCS applies to the exact
candidate and every documentation line/archive that already exists; it does
not require inventing an unreleased beta archive. The full V1-19 AC-DOCS
matrix requires clean builds of v0 current, v1 preview, one alpha archive
and one beta archive using the common theme and build route. Test both
deployment directions: v1 publication preserves v0 root/current/archive,
and v0 publication preserves v1 preview/archive and both manifest lines.
Check rendered labels and links from both sites, with JavaScript enabled and
disabled, keyboard navigation, active-state accessibility, nonexistent
same-path fallback, relative assets, search and canonical metadata. Check
the exact alpha/beta candidate docs against executable examples, schemas and
release notes; reject stale Docsy assets or stale v0 configuration claims.
Apply stable-v1/current labels only when publishing the stable release and its
matching documentation.

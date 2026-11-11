# Git branching and release standard

## Purpose and scope

This is the branching standard for every projectious.work repository. It
defines how humans and AI agents develop, specify, release, and maintain
software.

The model is *GitHub Flow* with *Conventional Commits*: every change is made
on a short-lived branch from `main` and squash-merged back into `main`. The
company deviates from canonical *GitHub Flow* in four places only:

1. releases are cut on short-lived `release/*` branches (see
   [Releases and tags](#releases-and-tags));
2. specifications are written and frozen on `spec/*` branches before
   development (see [Specifications](#specifications));
3. supported older versions live on long-term-support branches `lts/*` (see
   [Long-term support](#long-term-support));
4. a spec-driven phase with separate acceptance-test and implementation tracks
   integrates them on a phase branch before it reaches `main` (see
   [Specifications](#specifications)).

Besides `main`, only `lts/*` lines are long-lived; there are no development,
pre-release, or release lines.

The linked
[software verification and release-engineering standard](software-verification-and-release-engineering.md)
governs the evidence, tests, artifacts, and post-publication checks required
before a merge or a tag.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## Branches

| Branch | Purpose | Lifetime | Merges into |
|---|---|---|---|
| `main` | Integrated, releasable history. | permanent | — |
| `feat/<name>` | A feature. | short | `main` or its `lts/*` line |
| `fix/<name>` | A bug fix. | short | `main` or its `lts/*` line |
| `chore/<name>` | Everything else: refactoring, docs, tests, build, toolchain, dependencies. | short | `main` or its `lts/*` line |
| `spec/<name>-v<N>` | Drafting version N of a specification. | short | `main` |
| `<prefix>/<name>-tests`, `-impl` | Tracks of a spec-driven phase. | short | their phase branch `<prefix>/<name>` |
| `release/vX.Y.Z` | Preparing one release; no development. | short | `main` or its `lts/*` line |
| `lts/vN` | A supported older major (or `lts/vN.M` for a minor line). | until end of support | never merged |

`gh-pages` and other generated publishing branches are outside this model:
they are written only by their documented build command and never developed
on.

### Branch names

```text
<prefix>/[<issue>-]<slug>
```

- `<slug>`: lowercase kebab-case, `[a-z0-9-]`, no further `/`.
- Everything after `<prefix>/` (LTS token, issue number, slug, and suffixes
  such as `-v2` or `-tests`) is at most 40 characters.
- `<issue>-`: the tracking issue number, when the repository tracks issues.
- Work for an LTS line starts the name with `lts<N>-` (for `lts/vN`) or
  `lts<N>.<M>-` (for `lts/vN.M`; the only place a `.` is allowed), so the
  merge target is visible.
- `spec/[<issue>-]<slug>-v<N>` develops version N of a specification (the
  first is `-v1`); its directory and tag are defined by the
  [layout standard](specification-layout-and-lifecycle.md#names).
- `release/` names carry the exact version with `v`, including prereleases.

```text
feat/123-oauth-login
fix/null-crash-on-empty-config
chore/bump-go-1-24
spec/payments-api-v1
spec/123-payments-api-v2
release/v1.4.0
release/v2.0.0-rc.1
lts/v1
fix/lts1-123-null-crash        → PR into lts/v1
release/v1.4.1                 → cut from lts/v1
```

### Commit messages

Commit messages and squash-merge titles follow *Conventional Commits*: a type,
an optional scope, and an imperative subject; `!` or a `BREAKING CHANGE:`
footer marks a breaking change. Types map to *Keep a Changelog* sections and
to the *Semantic Versioning (SemVer)* bump. Required attribution trailers
follow the repository's policy. Layout follows the *50/72 Rule*, with one
deviation: a squash-merge title MAY exceed 50 characters for type, scope, and
pull-request number, but never 72.

The squash-merge title MUST be a valid Conventional Commit whose type matches
the branch prefix:

| Branch prefix | Allowed squash-title types |
|---|---|
| `feat/` | `feat` |
| `fix/` | `fix`; `perf` when it corrects a regression; `revert` when it undoes a regression |
| `chore/` | `docs`, `refactor`, `test`, `build`, `style`, `perf`, `chore`, `revert` |
| `spec/` | `docs(spec)` |
| `release/` | `chore(release): vX.Y.Z`, the only commit |

A breaking change keeps its branch prefix and carries `!` in the squash title.

## Normal flow

1. Fetch, then create the branch from the current tip of its target (`main`,
   or the `lts/*` line).
2. Keep it focused and short-lived; rebase on or merge the target when it
   drifts.
3. Run the repository's documented verification on the branch head and put
   the result in the pull request.
4. Obtain owner approval (see [Approval and merge
   authority](#approval-and-merge-authority)), then squash-merge.
5. Delete the branch.

Documentation, tests, schemas, migrations, and release metadata follow the
same path as code. AI agents MUST determine the target branch before creating
a branch or pull request; an unclear target is a stop condition.

### Keeping `main` green

Repositories have no CI pipelines, so verification runs locally or in the
development container. After every merge the author verifies the merge commit
on the target branch. A red `main` or `lts/*` branch is fixed forward
immediately, before new work builds on it. Acceptance tests written ahead of
their implementation merge together with the implementation that makes them
pass; `main` never carries an intentionally failing test.

## Specifications

Specifications are written before development, following the
[spec-driven development cycle](spec-driven-development-cycle.md).

1. Specification versions follow the
   [layout lifecycle](specification-layout-and-lifecycle.md#lifecycle):
   branch `spec/[<issue>-]<slug>-v<N>`, freeze, squash merge, tag
   `spec/<slug>/v<N>` on the squash commit. Frozen versions never change.
2. Implementation branches never change `spec/`. Their pull requests carry
   one `Spec: <slug>@v<N>` line per specification they implement (none for
   pure chores).
3. A phase with separate acceptance-test and implementation tracks keeps them
   on `<prefix>/<slug>-tests` and `<prefix>/<slug>-impl`, where `<slug>` names
   the phase (its roadmap item), not a specification. The integration owner
   merges them (plain merge, no squash) into the phase branch
   `<prefix>/<slug>`, which squash-merges into `main` once every acceptance
   test passes. Track branches are the only branches that do not target `main`
   or an `lts/*` line; the information barrier of the
   [spec-driven cycle](spec-driven-development-cycle.md) applies to them.

## Releases and tags

A release branch contains no development. It holds exactly one commit,
`chore(release): vX.Y.Z`, which updates the changelog (its entry is the curated
release notes) and the version metadata. Roadmap statuses change afterwards
(see the [roadmap standard](product-roadmap-and-development-evidence.md)).

1. Create `release/vX.Y.Z` from the tip of `main`, or of the `lts/*` line for
   an LTS release.
2. Make the release commit and run the release gates of the
   [verification and release-engineering standard](software-verification-and-release-engineering.md).
3. After owner approval, squash-merge the release branch.
4. Tag the squash commit on `main` (or `lts/*`) with `vX.Y.Z`, then publish.
5. Delete the release branch.

If the release finds a defect, close the release pull request, fix the defect
through `fix/*`, and cut a new release branch. Never fix on a release branch.
A version number is reused only if nothing was tagged or published with it.

Release tags MUST be annotated and SHOULD be cryptographically signed. A
published tag MUST never be moved or reused; a mistake is repaired by a new
version. Release artifacts MUST identify the version and the approved source
tree; the release record binds the tag and its squash commit to that tree.

Prereleases use the same flow with a prerelease version, for example
`release/v1.0.0-rc.1`. Version ordering follows Semantic Versioning 2.0.0:

```text
alpha.N < beta.N < rc.N < final
```

Counters start at 1 and increase monotonically within each stage. If tracked
content changes after a release candidate, the next candidate gets a higher
number; the final release MUST NOT differ from its last accepted release
candidate other than in version metadata and changelog.

## Long-term support

Create `lts/vN` only when a project commits to supporting major N after `main`
moves on. Branch it from the last `vN` release tag and document its end of
support.

- Work on an LTS line follows the normal flow with the LTS branch as target:
  `feat/`, `fix/`, and `chore/` branches named `lts<N>-…`, and releases on
  `release/vN.Y.Z` cut from the LTS branch.
- A fix that also applies to `main` lands on `main` first. It is then
  cherry-picked onto a branch from the LTS line with `git cherry-pick -x` and
  merged through the normal flow. If `main` no longer contains the affected
  code, fix directly on the LTS line and say so in the pull request.
- Never merge `main` into an LTS branch or an LTS branch into `main`.
- At end of support, delete the branch; tags and release artifacts remain.

## Approval and merge authority

Changes that require owner approval MAY receive that approval either:

- as a recorded GitHub pull-request review; or
- as an explicit instruction from the owner to the acting AI agent in the
  current, auditable working conversation.

Approval is a human decision, not a particular GitHub review object. An AI
agent that has verified the proposed source and target, reported the relevant
checks and risks, and received explicit owner approval MAY perform the merge on
the owner's behalf. The agent records the approval channel in the merge or
release evidence. General standing permission, silence, an ambiguous "looks
good", or approval for a different revision does not authorize a merge.

Repositories MAY configure GitHub branch protection or rulesets to require a
recorded approving review. This is an optional, project-level enforcement
choice, appropriate when multiple maintainers, external contributors,
regulatory controls, or project risk justify it. It is not the company
default. A repository using conversational owner approval MUST NOT configure a
required GitHub review count that prevents the approved AI-agent merge, unless
the owner intentionally chooses to complete a GitHub review as an additional
gate.

An agent MUST NOT bypass a configured protection rule merely because approval
exists elsewhere. If repository settings are stricter than the documented
project policy, stop and reconcile the setting or obtain the required GitHub
approval; do not use an administrator bypass as the routine path.

Before acting, an agent reports the branch type, source, target, and, for a
release, the version. An unclear target is a stop condition, not permission
to infer one. Before merging, it asks for explicit owner approval unless a
valid approval for the exact revision is already recorded. After approval, it
may merge, tag, and continue the authorized release sequence; publication
that needs the host follows the
[host-phase standard](human-controlled-host-phase-execution.md). A GitHub
`REVIEW_REQUIRED` result means the repository settings still demand a GitHub
review; the agent reports the policy mismatch instead of bypassing it.

## Protection and audit requirements

- Repositories SHOULD protect `main` and `lts/*` branches with GitHub rulesets
  when the collaboration or risk profile warrants platform enforcement.
  Technical branch protection is optional for owner-operated repositories.
- Whether enforced by GitHub or by the documented agent workflow, prohibit
  force-pushes to `main` and `lts/*` and their deletion except at end of
  support or through an explicitly approved recovery procedure.
- Squash merge is the only merge method into `main` and `lts/*`. The single
  exception is the [migration from version lines](#migration-from-version-lines).
- Restrict tagging rights to the release role or its explicitly authorized
  automation.
- Record every release as the
  [verification standard's release record](software-verification-and-release-engineering.md)
  defines it; fetch remote state immediately before ancestry and tag checks.

## Migration from version lines

This standard replaces the earlier promotion model (`vX.x-dev`,
`vX.x-pre-release`, `vX.x-release`, `vX.x-maintenance`, `phase/*`). Each
repository migrates once:

1. Inventory every non-`main` branch and its unique commits.
2. For each released older major, the owner decides: freeze it as `lts/vN`
   from its last release tag, or retire it after merging anything unique into
   `main`.
3. If a newer major line has diverged from `main`, first preserve the current
   `main` (as `lts/vN` or a tag), then make `main` point at the newer line
   with its history intact, so existing prerelease tags stay reachable. This
   is the only permitted non-squash change of `main` and needs explicit owner
   approval.
4. Delete the old version-line branches once nothing unique remains. Rename
   topic branches that use other prefixes to the prefixes above, or merge and
   delete them.

## Agent briefing

> *GitHub Flow*, *Conventional Commits*, *50/72 Rule*: `feat/`, `fix/`,
> `chore/`, `spec/` from `main` as `[<issue>-]<slug>`; unclear target: stop.
> Squash-merge after explicit owner approval; title type matches the prefix.
> `release/vX.Y.Z` holds only `chore(release): vX.Y.Z`; tag the squash commit.
> Specs: tag `spec/<slug>/vN`. `lts/vN`: names `lts<N>-…`; fix `main` first,
> then `git cherry-pick -x`. Never develop on `release/*` or bypass protection.

## Authoritative references

- [GitHub flow](https://docs.github.com/en/get-started/using-github/github-flow).
- [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/).
- [Git tag documentation](https://git-scm.com/docs/git-tag), including
  annotated and signed release tags.
- [git cherry-pick `-x`](https://git-scm.com/docs/git-cherry-pick).
- [Semantic Versioning 2.0.0](https://semver.org/).

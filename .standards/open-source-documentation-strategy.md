# Open-source documentation strategy

## Purpose and scope

This is the default documentation standard for public projectious.work
software repositories. It defines the information architecture, repository
landing experience, community files, release communication, local build and
deployment workflow, and quality gates expected of a credible open-source
project.

The standard does not prescribe Hugo, Docusaurus, a theme, or another page
generator. Projects choose technology according to their needs. The rendered
experience, source ownership, local reproducibility, and public contract are
normative.

The companion
[AI-agent accessibility and generative discovery standard](ai-agent-accessibility-and-generative-discovery.md)
defines the required Markdown, public discovery, and MCP projections of that
product contract. Repository-operability instructions are outside this
standard.

The [aibox repository](https://github.com/projectious-work/aibox) and
[documentation site](https://projectious-work.github.io/aibox/) are the
reference example. They demonstrate the intended repository-to-site flow, not
an immutable template or a claim that every current detail already conforms.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## Principles

1. **Documentation follows implementation.** Every user-visible change updates
   its documentation in the same change. Phase notes or development notes
   explain how the system grew; stable pages distill the implemented result.
2. **The repository is the source of truth (*Docs-as-Code*).** Documentation,
   examples, media, build configuration, and changelog entries (the source of
   release notes) live with the code;
   no wiki-only source. Unlike common Docs-as-Code setups, there is no CI
   pipeline (see principle 6).
3. **Write for distinct audiences.** A new user, experienced operator,
   integrator, contributor, security reporter, and maintainer need different
   paths through the same product truth.
4. **The README is a landing page and map.** It proves value and gets a reader
   to a first success; it does not duplicate the complete manual. Detail
   unfolds from README to Getting Started to guides to reference
   (*Progressive Disclosure*).
5. **Examples are executable promises.** Commands, configuration, output, and
   API examples are versioned and tested where practical.
6. **Local-first delivery.** Documentation builds, checks, and deploys through
   documented local commands. Project repositories contain no GitHub Actions
   workflows.
7. **Honest status beats marketing.** Maturity, supported versions, platform
   limits, security boundaries, breaking changes, and known defects are
   visible before adoption.
8. **Public by default, private where safety requires it.** Product decisions,
   roadmaps, and contribution processes are public. Vulnerability reports and
   code-of-conduct incidents use private channels.

## Required public surfaces

| Surface | Purpose | Requirement |
|---|---|---|
| GitHub About panel | Discovery and immediate orientation | Required |
| Root `README.md` | Product landing page and first success | Required |
| Documentation site | Task-oriented user and contributor documentation | Required for released software |
| `CONTRIBUTING.md` | Contribution contract and local development | Required |
| `SECURITY.md` | Supported versions and private reporting | Required |
| `LICENSE` | Legal permission and conditions | Required |
| `CODE_OF_CONDUCT.md` | Community behavior and private enforcement contact | Required for public contribution |
| `CHANGELOG.md` | Cumulative, versioned user-visible changes | Required |
| GitHub Releases | Per-version notes, tags, and downloadable artifacts | Required for released software |
| Roadmap | Direction with explicit commitment status | Required once future work is discussed publicly |
| Support guidance | Correct channel for questions, bugs, and security | Required; MAY be a README/docs section or `SUPPORT.md` |

`GOVERNANCE.md`, `MAINTAINERS.md`, `AUTHORS`, `CITATION.cff`, and a public
support forum are added when project scale, research use, or contributor
ownership makes them useful. They are not empty ceremony.

## GitHub repository presentation

The repository's upper-right About panel MUST contain:

- a short, plain-language description that says what the product does and for
  whom;
- the documentation site's canonical HTTPS URL in the website field;
- a small set of accurate lowercase topics covering purpose, ecosystem, and
  primary implementation language;
- an intentional social-preview image using the project logo and name.

The description is one sentence, avoids roadmap claims, and matches the first
README description. The website link points to the documentation root, not the
company homepage. Topics are reviewed when the product scope changes.

The repository overview SHOULD also expose an accurate latest release,
packages when applicable, a complete community profile, and a useful issue
tracker. Stale metadata is a release defect.

## Documentation-site shell

Every project documentation site provides a recognizable header with these
primary destinations:

| Header item | Destination |
|---|---|
| **Getting Started** | Installation and the shortest path to first success |
| **Documentation** | Documentation landing page and full navigation |
| **Roadmap** | Current direction and status-qualified future work |
| **Releases** | Current release, upgrade information, and release history |
| **Search** | Keyboard-accessible full-site search |
| **Light/dark switcher** | Explicit theme selection respecting system preference initially |
| **GitHub** | The canonical source repository |

On narrow screens the same destinations MAY move into an accessible menu, but
none disappear. The project logo and name link to the documentation homepage.
The active section is visually and semantically identifiable.

The footer contains, at minimum, links to documentation, GitHub,
`CONTRIBUTING.md`, `SECURITY.md`, `LICENSE`, releases, and projectious.work.
Copyright and trademark wording matches the repository license and brand
policy.

The shell MUST provide:

- responsive layouts usable on phones and desktop screens;
- keyboard navigation, visible focus, meaningful landmarks, sufficient
  contrast, alt text, and reduced-motion behavior;
- stable, human-readable URLs and useful page titles/descriptions;
- a favicon, project logo, and link-preview metadata;
- a not-found page with search and navigation back to useful content;
- no secrets, internal-only links, customer data, or unpublished operational
  evidence.

Search indexes user-facing documentation and important standalone pages. It
does not expose generated build files, private source material, or duplicate
version content without labels.

## Documentation information architecture

The content model follows the *Diátaxis Framework* (tutorials, how-to
guides, reference, explanation). The precise sidebar reflects the product,
but the content covers these questions:

1. **Overview:** what the product is, its value, intended users, maturity, and
   explicit non-goals.
2. **Getting started:** prerequisites, installation, first successful workflow,
   and next steps.
3. **Concepts:** the mental model, architecture, terminology, ownership, and
   security boundaries.
4. **How-to guides:** goal-oriented operational tasks, organized by user
   outcome rather than implementation package.
5. **Reference:** CLI/API, configuration, environment variables, file formats,
   compatibility, error behavior, and machine contracts.
6. **Troubleshooting:** symptoms, diagnostics, recovery, support escalation,
   and how to collect sanitized evidence.
7. **Roadmap:** released, done, in-progress, planned, and idea-level work with clear
   commitment semantics.
8. **Releases and migration:** supported versions, release notes, upgrades,
   breaking changes, deprecations, and rollback.
9. **Contributing:** architecture orientation, development environment, checks,
   documentation workflow, and release process.
10. **Security:** trust model, supported versions, reporting, and links to the
    canonical policy.

The four *Diátaxis* modes SHOULD remain visibly distinct (*MECE Principle*) whatever the
generator; "Getting started" is the tutorial and "Concepts" the explanation.
Public API symbols and schemas are documented close to source and linked from
the reference section.

## Documentation homepage

The homepage answers within one screen:

- what the product is;
- who it is for;
- the concrete outcome it creates;
- its maturity or current stable version; and
- how to install it or enter Getting Started.

It then provides a short demonstration, major capabilities, important
boundaries, a minimal workflow, and routes to documentation, roadmap, releases,
GitHub, and contribution. Avoid a generic marketing page that hides commands
or supported behavior.

## README standard

The root `README.md` uses this default order, adapting labels where needed:

1. centered project logo, project name, and one-line value proposition;
2. a restrained badge row for maturity/status, license, documentation, and
   latest release where reliable;
3. a visible maturity and support note, including major unsupported platforms
   or use cases;
4. a two- or three-paragraph description: problem, intended user, and current
   product outcome;
5. one current screenshot, diagram, or short terminal recording with alt text
   and a caption explaining what it proves;
6. “Why” or key capabilities, grounded in user outcomes;
7. prerequisites and installation with a verifiable version command;
8. a copyable quick start that reaches the first useful result;
9. common workflow or usage examples;
10. compatibility, current limitations, and version-line status;
11. a concise architecture or “how it fits” explanation when the product is
    part of a wider ecosystem;
12. a documentation map linking Getting Started, reference, troubleshooting,
    roadmap, releases, and migration;
13. development and contribution entry points;
14. support and security-reporting links;
15. changelog/releases, license, code of conduct, and acknowledgements.

Relative links are used for repository files so they work from clones and
non-default branches. Public documentation links use canonical HTTPS URLs.
Commands are tested against the maintained release. The README does not claim
planned behavior as released behavior and does not exceed what a new reader can
scan comfortably.

## Community and policy files

### `CONTRIBUTING.md`

The contribution guide welcomes contributions and states:

- supported contribution types and where to discuss larger changes;
- prerequisites and the canonical development environment;
- clone/setup/build/test/lint/format/docs commands;
- repository architecture and ownership boundaries;
- issue, branch, commit, and pull-request expectations;
- required tests, documentation, changelog, security, and compatibility work;
- generated-file rules and files contributors must not edit;
- how maintainers review, accept, defer, or decline changes;
- the release process or a link to its maintained guide;
- links to the code of conduct, security policy, license, and support channel.

The fastest contributor path appears first. Maintainer-only procedures do not
obscure the external contribution path.

### `SECURITY.md`

The security policy states supported versions, affected surfaces, what is out
of scope, a private reporting mechanism, the information needed to reproduce a
problem, handling of secrets/personal data, expected acknowledgement window,
and coordinated-disclosure policy. It explicitly says not to open a public
issue for an undisclosed vulnerability.

### `LICENSE`

Every public repository contains the complete text of its approved open-source
license at the root. The README names and links it. Source headers are added
only where the license or project policy requires them. Third-party licenses,
vendored assets, fonts, themes, and media retain required notices in a
`NOTICE`, attribution file, or documented dependency inventory.

### `CODE_OF_CONDUCT.md`

The code of conduct names expected behavior, scope, enforcement process, and a
private incident-reporting address. It links back to contributor guidance and
must identify who can act on reports without exposing reporters.

### Templates and support

Public projects SHOULD provide issue forms or templates for bugs, feature
requests, and documentation problems plus a pull-request template with testing
and documentation-impact prompts. These may live under `.github/`; the ban is
on `.github/workflows/`, not community templates.

Questions, bugs, security reports, and conduct reports have distinct channels.
If GitHub Issues is not a support forum, say so and link the correct location.

## Roadmap communication

The canonical machine-readable roadmap and its evidence rules follow the
[product roadmap and development evidence standard](product-roadmap-and-development-evidence.md).
Public roadmap pages SHOULD be generated from, or mechanically checked against,
that source rather than maintained independently.

Statuses, their evidence, and dependency rules follow that standard; public
pages label unreleased completed work **Done — awaiting release**, and
completion does not authorize a released support claim.

Items state the user outcome rather than only an internal component name.
Dates appear only when the project is prepared to maintain them. Removed or
reshaped ideas do not require a deprecation process; released contracts do.

## Development notes and implementation consistency

Significant roadmap phases, architectural changes, migrations, and major
release efforts maintain evolving implementation notes under `dev-notes/`.
Small fixes MAY omit a note when the review records why no development-note
impact exists.

A development note records implemented behavior, boundaries, decisions,
deviations, validation, security and compatibility consequences,
documentation changes, and owned follow-up work. It explains the path from
intent to observed result; it does not narrate commits or replace stable user
documentation.

The `done` and `released` gates follow the roadmap standard. A disagreement
between code and documentation is corrected or explicitly tracked;
documentation MUST NOT silently redefine a contract to match accidental
behavior.

Generated CLI, API, schema, configuration, and machine-interface reference
material has a local drift check. Maintained examples are executable promises
and validate against the current implementation or published schema where
practical.

## Changelog and release notes

`CHANGELOG.md` is the curated cumulative history in *Keep a Changelog* format.
Commit history is not a substitute.

Every published GitHub release has a human-curated description. Generated
notes MAY supply pull-request links and contributor credits, but they are input
to the release note, not the final narrative.

A release note has a title `<project> vX.Y.Z — <short outcome-oriented
title>`, then two to four sentences on who should care, the most important
outcome, compatibility, and whether action is required. Then these sections,
omitting only genuinely empty optional ones: Highlights (the two or three
changes most users notice); Added (with documentation links); Changed
(behavior, compatibility, performance, operations); Fixed (user-visible, with
issue links where useful); Breaking changes and migration (required actions,
deprecations, deadlines, migration and rollback links); Security (safe
disclosure summary and advisory link, never an unpatched issue); Known issues
(with workarounds); Install or upgrade (copyable commands and version check);
Artifacts and verification (checksums, signatures or attestations, SBOMs,
source tag); Contributors (credits and the full comparison link).

The release title and summary use product language, not an internal ticket or
commit subject. Prereleases are marked as prereleases and explain their test
audience and support limitations. The final release links to the matching
documentation version and migration guide. Assets have predictable names and
include checksums; signatures, attestations, and SBOMs are published when the
project release policy requires them.

The documentation site's **Releases** destination provides the current stable
version, supported lines, install/upgrade link, recent release summaries, and a
clear route to GitHub Releases for assets and complete history.

## Versioned documentation

The default documentation describes the current stable release. When supported
versions have materially different interfaces, preserve versioned snapshots or
an equivalent selector. Archived pages display their version, support state,
and a link to current documentation.

Development documentation is visibly labeled and MUST NOT replace stable docs
before release. URLs for established pages remain stable; redirects accompany
moves. Migration pages remain available as long as affected supported users
may need them.

## Local build and verification

Each repository exposes discoverable local commands, either as focused scripts
or subcommands of a transparent maintenance script:

```text
docs-serve             start a local preview with live reload when available
docs-build             produce the static site from a clean checkout
docs-check             run formatting, links, examples, accessibility smoke, and build checks
docs-deploy --dry-run  build the exact artifact without publishing
docs-deploy            publish the reviewed artifact to gh-pages
```

Names MAY vary, but `CONTRIBUTING.md` and the documentation-site README list the
exact commands and prerequisites. The implementation:

- uses a committed dependency manifest and lockfile;
- does not install mutable global dependencies silently;
- fails on broken internal links, missing required pages, invalid front matter,
  duplicate routes, or build errors;
- checks external links with bounded retries and an explicit allowlist for
  known transient endpoints;
- validates copyable commands and configuration snippets where practical;
- verifies screenshots/recordings are referenced, reasonably current, and
  free of secrets or personal data;
- builds without production credentials;
- writes generated output only to ignored build directories;
- reports the source commit and generator/tool versions.

Every implementation change includes a documentation-impact assessment. A
change with no documentation update states why. Release checks build the site
from a clean checkout and inspect the rendered result, not only Markdown
syntax.

## GitHub Pages deployment

Published documentation uses the repository's `gh-pages` branch as the Pages
source, served from its root. Project-authored GitHub Actions and files under
`.github/workflows/` are prohibited.

The local deployment command:

1. requires an authenticated, explicitly authorized publisher;
2. fetches remote state and requires the approved source commit and a clean
   worktree;
3. runs the complete documentation check and static build;
4. supports `--dry-run` without network mutation;
5. stages output in a fresh temporary directory;
6. includes `.nojekyll` when the generated site requires bypassing Jekyll;
7. creates a generated commit containing source repository, source commit,
   build timestamp policy, and generator version;
8. replaces or advances only the generated `gh-pages` branch;
9. verifies GitHub Pages is configured for `gh-pages` at `/`;
10. probes the canonical HTTPS URL and key routes after publication.

`gh-pages` contains generated output only and is never a development or review
branch. Its documented replacement/force-push behavior is a narrow exception
to the company prohibition on rewriting `main`, `lts/*`, and `release/*`. Never edit
the branch by hand. The deploy command must not expose authentication tokens in
arguments, logs, generated files, or Git history.

## Quality, accessibility, and maintenance

Documentation is reviewed with the same seriousness as code.

Text meets *Plain Language (ISO 24495)*: readers find, understand, and use
what they need. Prose follows *Plain English according to Strunk & White*:
short sentences, active voice, concrete words, no filler. German text follows
*Gutes Deutsch nach Wolf Schneider*. Procedures, CLI and error messages, and
agent instructions follow the writing rules of *Simplified Technical English
(ASD-STE100)*; its controlled dictionary is not required. Documentation, CLI
messages, release notes, and commit messages use the product's Ubiquitous
Language (*Domain-Driven Design according to Evans*): the same terms as the
specification and the code. Visual design of documentation sites follows the
[user-interface design standard](user-interface-design.md).

Required quality checks cover:

- correctness against the released interface and supported platforms;
- tested commands, examples, schemas, and sample output;
- internal and external links;
- spelling, terminology, headings, and readable line lengths;
- keyboard use, semantic headings/landmarks, alt text, captions, contrast, and
  reduced motion;
- responsive layout and representative current browsers;
- search indexing and useful zero-result behavior;
- page metadata, canonical URLs, sitemap/robots behavior, and social previews;
- secret, personal-data, and stale-version scans;
- rendered-site smoke tests for homepage, Getting Started, Documentation,
  Roadmap, Releases, Security, and not-found routes.

Analytics are optional, privacy-preserving, documented, and disabled until a
project has a justified question to answer. Do not add tracking merely because
a theme supports it.

Each major documentation area has an owner. Repeated support questions become
documentation improvements. Stale pages are corrected, explicitly archived,
or removed with redirects; they are not left discoverable as silent traps.

## AI-agent and machine-consumption guidance

`AGENTS.md` remains the canonical provider-neutral agent entry point and links
the same build, test, documentation, branching, and release contracts used by
humans. Provider-specific files are thin pointers.

Machine-readable schemas, CLI help, generated reference, `llms.txt`, or MCP
documentation resources MAY improve discovery, but they do not replace the
human documentation or create a second source of truth. Generated references
must have drift checks against implementation.

## Minimum release documentation gate

A release is not publishable until the pre-publication items pass; the last
item is checked after publication:

- README status, installation, quick start, and compatibility match the
  candidate;
- stable documentation builds cleanly and key rendered routes pass smoke tests;
- the changelog entry (the curated release notes) is complete;
- breaking changes include migration and rollback guidance;
- supported versions and security policy are correct;
- roadmap `done` states match the release scope and current-version labels
  are correct (`released` follows after publication);
- artifacts, checksums, signatures/attestations, and SBOM links match what was
  actually produced;
- the GitHub About description and documentation link remain accurate;
- the exact validated documentation artifact is deployed to `gh-pages`; and
- (after publication) the public site and release links are verified.

## Authoritative references

- [GitHub: About releases](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases)
- [GitHub: Publishing Pages from a branch](https://docs.github.com/en/pages/getting-started-with-github-pages/configuring-a-publishing-source-for-your-github-pages-site)
- [GitHub: Community profiles](https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions/about-community-profiles-for-public-repositories)

## Agent briefing

> Documentation follows implementation in the same change and lives in the
> repository. The README is a landing page and map (*Progressive Disclosure*).
> Text meets *Plain Language (ISO 24495)*; procedures use *ASD-STE100* rules.
> Examples are executable promises. Docs deploy locally to `gh-pages` without
> GitHub Actions. Status and limits are honest; a release needs the minimum
> release documentation gate.

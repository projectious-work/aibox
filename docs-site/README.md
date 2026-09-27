# aibox documentation site

The v1 preview is built with Hugo Extended and the pinned
[`brand-theme-hugo-vanilla`](https://github.com/projectious-work/brand-theme-hugo-vanilla)
Hugo module (`v0.3.4`), matching the v0.x documentation line. Its Node assets
are locked in `package-lock.json`; do not install a floating theme checkout or
reintroduce Docsy, Bootstrap, or Font Awesome.

## Prerequisites

- Hugo Extended `0.157.0` or newer
- Go (the Hugo module declares Go `1.22`)
- Node.js `18` or newer and npm

## Local preview and build

From the repository root, start the live v1 preview at
`http://localhost:1316/aibox/v1.x/`:

```sh
./scripts/maintain.sh docs-serve
```

To make the production build, including the repository base path and `.nojekyll`:

```sh
./scripts/build-docs.sh
```

The build script installs locked Node dependencies with `npm ci` when the
required theme assets are missing. Hugo resolves the pinned theme module from
`docs-site/go.mod` and `go.sum`. To build into a scratch directory, pass Hugo
arguments through the script, for example:

```sh
./scripts/build-docs.sh --destination /tmp/aibox-v1-docs
```

The v1 source sets its preview base path to `/aibox/v1.x/`. Public publication
is a separate operation: do not push generated files or deploy docs as part
of ordinary content edits. Release workflows validate and publish the
candidate snapshot while preserving the v0.x root and other archives.

## Content workflow

1. Confirm the feature contract and its roadmap phase in `spec/v1/` before
   writing product instructions. The specification is the design source until
   an implementation phase proves the behavior.
2. Add or update audience-specific pages in `docs-site/content/` with the
   implementation. Label planned examples clearly; do not present them as
   runnable until checked against the candidate implementation.
3. Preview with `./scripts/maintain.sh docs-serve`, then run
   `./scripts/build-docs.sh` after the final content change.
4. Check page links, navigation, code examples, and the generated v1 preview
   labels. For documentation foundation changes, also build the v0 current
   site from its release-line checkout and confirm its root/current label.
5. Include changed paths, intended audience, the phase's docs deliverable,
   and build/check results in phase evidence. Publication remains part of the
   reviewed release workflow.

The stable site root remains the current v0.x documentation. The `/v1.x/`
tree is visibly a preview until a separately approved stable-v1 promotion.
Current content that describes the earlier alpha design is withheld from the
public v1 build; it is retained only as migration source while phase-specific
replacement docs are written.

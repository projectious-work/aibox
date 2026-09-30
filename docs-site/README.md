# aibox Documentation Site

This directory contains the public aibox documentation site. It uses
[Hugo](https://gohugo.io/) with the
[projectious.work brand theme](https://github.com/projectious-work/brand-theme-hugo-vanilla)
at v0.3.4.

## Installation

The site requires Hugo Extended, Go for Hugo Modules, and Node.js for the
theme's Tailwind build.

```bash
npm --prefix docs-site ci
hugo mod get github.com/projectious-work/brand-theme-hugo-vanilla@v0.3.4
hugo mod tidy
```

## Local Development

```bash
./scripts/maintain.sh docs-serve
```

This command starts a local development server on port 1316. Most changes are reflected live without having to restart the server.

## Build

```bash
./scripts/build-docs.sh
```

This command refreshes `data/theme_catalog.json` from the CLI's audited theme
source, then generates static content into `docs-site/public/`.

The releases menu always links to the current v0.x root and the v1.x preview.
When `DOCS_RELEASES_MANIFEST` points to a valid shared `releases.json`, the
build renders both release lines and their archives directly in the menu.
Without it, the menu uses the standalone release list in `hugo.yaml`.
The manifest build keeps the generated theme catalog available to Hugo.

## Deployment

The repository release script deploys docs as part of the release flow. For a
manual deployment:

```bash
cd ..
./scripts/maintain.sh docs-deploy --dry-run
./scripts/maintain.sh docs-deploy
```

The maintenance command builds locally and pushes the static output to the
`gh-pages` branch. The project does not use GitHub Actions for documentation
deployment. Do not use the generic `npm run deploy` path for this repository.

---
title: aibox v1.x preview
description: "The specification and implementation plan for a Go-based aibox distribution built on the Dev Container CLI."
---

# aibox v1.x preview

**Preview status:** v1 is under development. The Go application, CLI, MCP
server, Dev Container Features, and user workflows described by the
specification are not yet implemented or available for installation.

The specification defines a curated workspace distribution built on the
upstream Dev Container CLI. Projects will use standard `devcontainer.json`,
Features, Dockerfiles, Compose files, and native tool configuration. aibox v1
will not use `aibox.toml` as a project entry point.

## What this preview contains

This site currently publishes the v1 product specification and phase roadmap,
not installation or operating instructions. Planned commands and
configuration examples in the specification are design contracts, not usable
interfaces. Each implementation phase must update its user or maintainer
documentation and build this site before it can be marked shipped.

Read the [product and architecture specification](https://github.com/projectious-work/aibox/tree/v1.x-dev/spec/v1),
including the [canonical phase roadmap](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/roadmap.yaml).

## Stable version

The stable documentation root continues to describe the current v0.x product.
This `/v1.x/` site is a preview line and does not replace the stable site.
For the current released product, use the [v0.x documentation](https://projectious-work.github.io/aibox/).

## Contributing

See the [documentation site README](https://github.com/projectious-work/aibox/blob/v1.x-dev/docs-site/README.md)
for the local build and contribution workflow.

---
title: aibox v1.x preview
description: "The specification and implementation plan for a Go-based aibox distribution built on the Dev Container CLI."
---

# aibox v1.x preview

**Preview status:** v1 is under development. Its Go application, CLI, MCP
server, and user workflows are not yet implemented or available to install.
The specification defines a curated workspace distribution built on the
upstream Dev Container CLI. v1 uses standard `devcontainer.json`, Features,
Dockerfiles, Compose files, and native tool configuration; `aibox.toml` is not
the v1 project entry point.

## What this preview contains

This site currently publishes the v1 specification and phase roadmap, not
installation or operating instructions. Planned commands and configuration
examples in the specification are design contracts, not usable interfaces.
Each implementation phase must update its user or maintainer documentation
and build this site before it can be marked shipped.

Read the [product and architecture specification](https://github.com/projectious-work/aibox/tree/v1.x-dev/spec/v1),
including the [canonical phase roadmap](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/roadmap.yaml).

<!--
The v1 product is not implemented yet; the following Docsy page content is
retained only as a migration source and is excluded from the preview.
-->
{{% blocks/feature icon="fa-solid fa-box" title="Declarative workspaces" %}}
One inspectable contract produces a reproducible development environment.
{{% /blocks/feature %}}

{{% blocks/feature icon="fa-solid fa-code-branch" title="Standard output" %}}
Dockerfile, Compose, and Dev Container files remain readable and interoperable.
{{% /blocks/feature %}}

{{% blocks/feature icon="fa-solid fa-people-group" title="Provider neutral" %}}
Agent context lives in the project, with thin provider-specific entry points.
{{% /blocks/feature %}}
{{< /blocks/section >}}

{{< blocks/section color="dark" >}}
<div class="col-12">

## Quick start

```bash
curl -fsSL https://raw.githubusercontent.com/projectious-work/aibox/main/scripts/install.sh | bash
mkdir my-project && cd my-project
aibox init my-project --harness claude --addon python
aibox apply
aibox up
```

</div>
{{< /blocks/section >}}

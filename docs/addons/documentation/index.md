# Documentation Frameworks


# Documentation Framework Addons

Documentation addons install static site generators and documentation tools.

| Addon | Tool | Install Method |
|-------|------|---------------|
| `docs-mkdocs` | MkDocs + Material theme | uv |
| `docs-zensical` | Zensical | uv |
| `docs-docusaurus` | Docusaurus | npm |
| `docs-starlight` | Starlight (Astro) | npm |
| `docs-mdbook` | mdBook | Binary download |
| `docs-hugo` | Hugo Extended | Binary download |

Current curated pins are Docusaurus 3.10.2, Hugo 0.168.0, mdBook 0.5.4,
MkDocs 1.6.1 with Material 9.8.0, and Zensical 0.0.72. Starlight remains
scaffolded through the upstream `create-starlight` package.

## Example

```toml
[addons.docs-docusaurus.tools]
docusaurus = {}
```

After `aibox apply`, the documentation tool is available inside the container. Initialize your docs project as usual (e.g., `npx create-docusaurus@latest docs classic`).


---
Source: https://projectious-work.github.io/aibox/docs/addons/documentation/index.md

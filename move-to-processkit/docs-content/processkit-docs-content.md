# processkit docs content — collected from aibox docs

This file collects documentation content that belongs in processkit's own docs,
not in aibox's docs. Each block notes the source file, the exact content, and
why it belongs in processkit.

---

## Block 1 — Package composition table

**Source:** `/workspace/docs-site/docs/context/process-packages.md`, Packages section

**Content:**

> processkit ships five packages. For the exact skill composition of each, see:
> …
>
> | Package | Best for |
> |---------|----------|
> | `minimal` | Scripts, experiments, small utilities |
> | `managed` | Recommended default — backlog, decisions, standups, session management |
> | `software` | Software projects with a recurring build/test/review cycle |
> | `research` | Learning, documentation, academic work |
> | `product` | Full product development with security, ops, design, planning |

**Why it belongs in processkit:**
The "Best for" descriptions are processkit's value proposition, not aibox's.
Every processkit release can change the scope or rename packages, forcing an
aibox docs update. aibox docs should only name the five packages and link to
processkit for descriptions. The table itself — with the prose rationale for each
tier — belongs in processkit's own package reference or getting-started guide.

---

## Block 2 — Package composition references (skills/index.md)

**Source:** `/workspace/docs-site/docs/skills/index.md`, Packages section

**Content:**

> processkit ships five packages (`minimal`, `managed`, `software`, `research`,
> `product`) that compose skill sets. Select the one that fits your project in
> `aibox.toml`:
>
> For the full package definitions and their exact skill composition, see:
> - `context/templates/processkit/<version>/.processkit/packages/` (after `aibox sync`)
> - https://github.com/projectious-work/processkit/releases

**Why it belongs in processkit:**
The explanation of *what each package composes* (which skills are in which
package, why the tiers exist) is processkit content. aibox docs need only tell
users which `[context].packages` value to use and where to find the YAML —
the conceptual documentation of the packages themselves belongs in processkit.

---

## Block 3 — Skill catalogue browsing instructions (processkit GitHub links)

**Source:** `/workspace/docs-site/docs/skills/index.md`, "Skill catalogue and documentation" callout

**Content:**

> - **GitHub:** https://github.com/projectious-work/processkit/tree/main/src/skills
> - **Releases:** https://github.com/projectious-work/processkit/releases
> - **In your project, after `aibox sync`:** `context/skills/` and
>   `context/templates/processkit/<version>/context/skills/`

**Why it belongs in processkit:**
The authoritative catalogue of what skills exist, their descriptions, and their
organisation by category is processkit content. This callout currently substitutes
for processkit having its own documentation site. Once processkit ships a docs
site, the aibox page should replace this block with a direct link to processkit's
skill catalogue page.

---

## Block 4 — processkit "not yet deployed" notice (repeated in two files)

**Source:**
- `/workspace/docs-site/docs/skills/index.md`, "Skill catalogue and documentation" callout
- `/workspace/docs-site/docs/context/process-packages.md`, Packages section

**Content (skills/index.md):**

> **processkit documentation** is not yet deployed as a standalone site.
> Until then, browse the upstream source directly: …

**Content (process-packages.md):**

> **processkit documentation** is not yet deployed as a standalone site.
> Until then, the package YAML files are the canonical source of truth.

**Why it belongs in processkit:**
These notices exist because processkit has no docs site yet. When processkit
publishes one, these placeholders should be replaced with a link. The content
that would live at those links (skill catalogue, package definitions) belongs
entirely in processkit — aibox docs are just stand-ins.

---

## Block 5 — "Why this split?" rationale in skills/index.md

**Source:** `/workspace/docs-site/docs/skills/index.md`, "Why this split?" section

**Content:**

> - **Independent release cadence.** processkit can ship a new skill or fix a
>   prompt without forcing an aibox CLI release.
> - **Reusable content.** Other tools can consume processkit directly without
>   taking a dependency on aibox or its container stack.
> - **Forkable content.** A team can fork processkit, point `[processkit].source`
>   at the fork, and ship a private skill catalogue without forking aibox itself.
> - **Smaller aibox.** The aibox binary stays focused on container lifecycle and
>   the install/diff/migrate machinery.

**Why it belongs in processkit:**
The "reusable content" and "forkable content" bullets describe processkit's own
design goals, not aibox's. These are arguments for why processkit exists as a
standalone project that other tools (not just aibox) can consume. processkit's
own README or introduction page should carry this rationale. The aibox docs only
need the first and last bullets (which speak to the aibox side of the boundary).

---

## Summary of aibox docs boundary

After moving the above content to processkit docs, aibox docs should contain
only:

- HOW `aibox sync` installs processkit content (the three-way diff machinery,
  template mirror layout, editable copies)
- WHICH package names exist (as a lookup table for `[context].packages` values)
- HOW to configure `[processkit].source` and `[processkit].version` in
  `aibox.toml`
- HOW `aibox kit skill list/info/install/uninstall` commands work

The WHY and WHAT of each package's skills, the skill catalogue, and the
getting-started rationale for choosing a package all belong in processkit docs.

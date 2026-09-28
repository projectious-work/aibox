# V1-02 shared documentation foundation evidence

The documentation foundation was implemented by merged PR #464 on `v1.x-dev`
(`1904e9574cd44a87c9a85535690eff8b82897c7b`) and the matching v0
navigation update by merged PR #465 on `v0.x-release`
(`36a9f4496b5383c231187b030855213d5da6f943`). Both current
`docs-site/go.mod` files pin
`github.com/projectious-work/brand-theme-hugo-vanilla v0.3.4`.

On 2026-09-27, local Hugo Extended v0.165.0 builds succeeded:

| Line | Command | Result |
|---|---|---|
| v1.x | `./scripts/build-docs.sh --destination /tmp/aibox-v1-roadmap-docs` | 31 pages, 142 static files |
| v0.x | `./scripts/build-docs.sh --destination /tmp/aibox-v0-roadmap-docs` | 185 pages, 138 static files |

The earlier Docsy build in [baseline evidence](baseline-evidence.md) is
historical pre-migration evidence, not the current v1 build. Docsy source is
retained under `docs-site/legacy-docsy-*` and excluded from the active site.
Candidate-specific alpha/beta snapshots and publication checks remain in
V1-19, so this phase remains `in_progress`, not `shipped`.

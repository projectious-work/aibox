# Development logbook

This is the maintainer-facing logbook required by the [projectious.work
development-evidence standard](https://github.com/projectious-work/internal/blob/main/docs/standards/product-roadmap-and-development-evidence.md).
Each non-trivial v1 phase has an evolving note named for its roadmap ID. The
[v1 roadmap](../spec/v1/roadmap.yaml) remains the source of truth for phase
scope and status. Notes explain choices, evidence and gaps rather than
narrating commits or replacing user documentation.

| Phase | Development note | Current meaning |
|---|---|---|
| V1-01 | [Implementation baseline](V1-01-implementation-baseline.md) | Merged baseline evidence; no v1 runtime claim |
| V1-02 | [Shared documentation build](V1-02-shared-documentation-build.md) | Merged Hugo foundation; publication gates remain later |
| V1-03 | [Go operation core](V1-03-go-operation-core.md) | Implementation done with local CLI/MCP evidence; publication pending |

Update a phase note after each meaningful implementation wave: record what
changed, which contract was checked, validation, remaining gaps, and security
or compatibility consequences. A merged PR alone does not make a phase
shipped; release and published-artifact evidence are separate gates. The
roadmap's `implementationStatus: done` records a verified local implementation
without claiming a release.

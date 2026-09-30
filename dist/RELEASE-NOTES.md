# aibox v0.35.3 — 2026-09-30

**Summary:** This patch puts the current v0 code and its release evidence back
on the stable `main` lineage. v0 users receive the v0.35.1 tool updates,
the v0.35.2 host Yazi probe repair, and the subsequent mandatory host-check
correction under one exactly promoted release.

## Changed

- Route v0 release preparation through `v0.x-dev`, then fast-forward the
  validated commit through staging and release to `main`.
- Refuse publication from an unmerged or divergent v0 candidate, and refuse
  a stable tag until `main` identifies the validated commit.
- Preserve the published v0.35.1 and v0.35.2 tags as immutable history.

## Fixed

- Retain the required macOS host-check correction after v0.35.2.
- Validate interactive Yazi startup in a PTY in local E2E tests; remove a
  non-PTY probe and an obsolete `yazi --debug` invocation.

## Compatibility

- Minimum processkit version remains v0.28.8.
- No project configuration change is required.
- The v1 development branch remains independent.

## Upgrade notes

Choose v0.35.3 for new v0 installations after its macOS host validation and
publication complete. Existing v0.35.1 and v0.35.2 tags remain available for
provenance but are superseded by this patch.

[v0.35.3]: https://github.com/projectious-work/aibox/compare/v0.35.2...v0.35.3

# v0.35.3 — restore the stable release chain

> Promotes the current v0 code through its development line to main and keeps host validation mandatory.


v0.35.3 brings the current v0 code and release checks through `v0.x-dev`,
`v0.x-pre-release`, and `v0.x-release` to `main`. The release tag identifies
the exact commit on `main`; published v0.35.1 and v0.35.2 tags remain immutable
historical releases.

The release command now refuses to promote an unmerged v0 candidate or tag a
commit that has not reached `main`. The mandatory release-host checks from
v0.35.2 and the subsequent host-gate correction remain included. v1 work
continues on its independent development line.

No project configuration change is required. Use v0.35.3 for new v0 installs
after its host validation and publication complete.

[Full v0.35.3 release notes](https://github.com/projectious-work/aibox/releases/tag/v0.35.3)


---
Source: https://projectious-work.github.io/aibox/changelog/release-v0-35-3/index.md

# v0.35.5 — Mermaid images build again

> Fixes image builds with the Mermaid add-on: the headless browser install now fetches the apt package lists before installing its system libraries.


v0.35.5 fixes `aibox` images that enable the Mermaid add-on with its pinned
Puppeteer runtime. Since v0.35.1 the build failed in the step
`puppeteer browsers install chrome-headless-shell --install-deps`, with apt
reporting every browser library (`libnss3`, `libgbm1`, `fonts-liberation`
and others) as "not installable".

`--install-deps` installs those libraries with `apt-get`, but the aibox base
image removes the apt package lists after its own installs, so apt knew no
packages. The step now runs `apt-get update` first and removes the lists
again afterwards, so the image stays as small as before.

No project configuration change is required. Run `aibox apply` and rebuild
to pick up the fix.

[Full v0.35.5 release notes](https://github.com/projectious-work/aibox/releases/tag/v0.35.5)


---
Source: https://projectious-work.github.io/aibox/changelog/release-v0-35-5/index.md

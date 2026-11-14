# aibox v0.35.5 — 2026-11-14

**Summary:** This patch makes images with the Mermaid add-on build again.

## Fixed

- The Mermaid add-on's browser step now runs `apt-get update` before
  `puppeteer browsers install chrome-headless-shell --install-deps` and
  removes the package lists again afterwards. The base image removes the apt
  lists after its own installs, so since v0.35.1 apt reported every browser
  library as "not installable" and the image build failed (#519).

## Changed

- v0 releases are made from `lts/v0`: the release validates and tags the
  exact `lts/v0` tip and moves no other branch (#521, refs #514). `main` is
  now the v1 line.

## Upgrade

No project configuration change is required. Run `aibox apply` and rebuild
the image.

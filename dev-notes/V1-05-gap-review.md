# Native Feature and gap review

Frozen baseline `9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6`: **98 rows; 4 source-reviewed upstream, 68 gaps, 26 uncertain community candidates**. All 98 original addon source hashes match the ledger. Current addon checkout differs in 38 entries (including absent recipes and hidden Hugo pin drift); frozen baseline governs this audit. No gap implementation or removal is approved. V1-05 remains in progress; source review does not satisfy the required live installation and two-architecture checks.

Qualification is source/OCI review, not runtime validation. Disabled Features must be omitted on a clean base and absence checked after rebuild. Community registry membership is discovery evidence only.

## Qualified upstream rows

| Row | Exact Feature | Options |
|---|---|---|
| ADDON:docs-hugo/hugo | `ghcr.io/devcontainers/features/hugo@sha256:15465c956810aa164c9644b6df5ca36f6cac3e7a86b7dcc474a1fb830b21c509` | `{"version":"0.165.0","extended":true}` |
| ADDON:cloud-aws/aws-cli | `ghcr.io/devcontainers/features/aws-cli@sha256:1f93c8315b7a6d76982ebb2269f8b0d50413fc0f965c032edf4aee0caceb73ef` | `{"version":"latest"}` |
| ADDON:cloud-azure/azure-cli | `ghcr.io/devcontainers/features/azure-cli@sha256:d98f1066c077be0fa9d115b718f458bd803e415181b4a96f82a6f5d9f77241ac` | `{"version":"latest","installBicep":false}` |
| ADDON:git-ui/gh | `ghcr.io/devcontainers/features/github-cli@sha256:bd7ab48a832228f633239277552c30b353867fef2e5b037e064b4e64f0b843f2` | `{"version":"latest"}` |

Publishers: Dev Container Spec Maintainers. Exact OCI manifest digests and frozen source URLs are in the JSON. Hugo preserves actual `0.165.0` extended pin despite blank declaration.

## Full unresolved review

Every listed row needs an owner retain/remove decision. Proposed routes use native image/Dockerfile/package installation. No new aibox Feature is assumed.

### ai-aider

- **aider** (`unversioned/distro-or-latest`, gap): Install aider-chat using a native Dockerfile and isolated uv tool or venv, conditional on enable; preserve executable exposure and dependency isolation.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### ai-claude

- **claude** (`unversioned/distro-or-latest`, uncertain): Official signed Claude Code apt repository; preserve stable/latest/exact version channels and disabled absence.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### ai-codex

- **codex** (`unversioned/distro-or-latest`, uncertain): Install @openai/codex from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### ai-copilot

- **copilot** (`unversioned/distro-or-latest`, gap): Install an exact Copilot CLI release from verified architecture-specific assets in a native Dockerfile, or select an exact version of the official Copilot Feature after review. Its `latest` option creates an auto-update flag, and its `postStartCommand` runs `copilot update`; this changes the tool after container start.

**Decision:** Retain with an exact native install, approve an exact-version upstream Feature, or remove this row.

### ai-continue

- **cn** (`unversioned/distro-or-latest`, gap): Install @continuedev/cli from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### ai-gemini

- **gemini** (`unversioned/distro-or-latest`, uncertain): Install @google/gemini-cli from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### ai-hermes

- **hermes** (`0.19.0`, uncertain): Install hermes-agent==0.19.0 using a native Dockerfile and isolated uv tool or venv, conditional on enable; preserve executable exposure and dependency isolation.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### ai-mistral

- **mistral** (`unversioned/distro-or-latest`, gap): Install mistralai using a native Dockerfile and isolated uv tool or venv, conditional on enable; preserve executable exposure and dependency isolation.
  Owner option: {"kind": "removal_candidate", "proposal": "Owner may remove the programmatic SDK from the global CLI tooling catalog and keep it in project dependencies."}

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### ai-opencode

- **opencode** (`0.0.55`, gap): Download verified architecture-specific assets from opencode-ai/opencode at 0.0.55 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### ai-tau

- **tau** (`0.3.13`, uncertain): Install tau-ai==0.3.13 using a native Dockerfile and isolated uv tool or venv, conditional on enable; preserve executable exposure and dependency isolation.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### audio-voice

- **sox** (`unversioned/distro-or-latest`, gap): Install sox from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **sox-pulse** (`unversioned/distro-or-latest`, gap): Install libsox-fmt-pulse from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **pulseaudio-utils** (`unversioned/distro-or-latest`, gap): Install pulseaudio-utils from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **alsa-pulse** (`unversioned/distro-or-latest`, gap): Install libasound2-plugins from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### browser-testing

- **playwright** (`1.62.1`, uncertain): Install @playwright/test@1.62.1 from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable.
- **axe-playwright** (`4.13.0`, gap): Install @axe-core/playwright@4.13.0 from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable. Also install axe-core@4.13.0 as in v0; catalog name refers to @axe-core/playwright, not the separate axe-playwright package.
- **chromium** (`unversioned/distro-or-latest`, gap): Use exact Playwright package plus playwright install --with-deps --no-shell chromium, only when enabled; browser revision follows 1.62.1. Official Playwright image brings browsers that violate individual toggles.
- **firefox** (`unversioned/distro-or-latest`, gap): Use Playwright install firefox only when enabled, locked to Playwright 1.62.1 browser revision.
- **webkit** (`unversioned/distro-or-latest`, gap): Use Playwright install webkit only when enabled, locked to Playwright 1.62.1 browser revision.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### cloud-gcp

- **gcloud-cli** (`unversioned/distro-or-latest`, uncertain): Install from official signed Google Cloud apt repository in a native Dockerfile, conditional on enable; lock version/snapshot.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### cloudflare

- **cloudflared** (`unversioned/distro-or-latest`, gap): Install from official signed Cloudflare apt repository conditional on enable; lock version/snapshot.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### data-preview

- **sqlite3** (`unversioned/distro-or-latest`, gap): Install sqlite3 from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **csvkit** (`unversioned/distro-or-latest`, gap): Install csvkit from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### data-visualization

- **vega-cli** (`6.4.0`, gap): Install vega-cli@6.4.0 from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable.
- **vega-lite** (`6.4.3`, gap): Install vega-lite@6.4.3 from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### diagramming

- **d2** (`0.7.1`, gap): Download verified architecture-specific assets from d2lang/d2 at 0.7.1 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.
- **graphviz** (`unversioned/distro-or-latest`, gap): Install graphviz from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### docs-docusaurus

- **docusaurus** (`3.10.1`, gap): Install @docusaurus/core@3.10.1 from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### docs-mdbook

- **mdbook** (`0.5.4`, uncertain): Download verified architecture-specific assets from rust-lang/mdBook at 0.5.4 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### docs-mkdocs

- **mkdocs** (`1.6.1`, gap): Install mkdocs==1.6.1 using a native Dockerfile and isolated uv tool or venv, conditional on enable; preserve executable exposure and dependency isolation. Preserve mkdocs-material==9.7.7 as bundled companion dependency.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### docs-starlight

- **starlight** (`unversioned/distro-or-latest`, gap): Install create-starlight from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable.
  Owner option: {"kind": "removal_candidate", "proposal": "Owner may remove the global project scaffolder and use project-local create-starlight bootstrap instead. No removal is authorized."}

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### docs-zensical

- **zensical** (`0.0.57`, gap): Install zensical==0.0.57 using a native Dockerfile and isolated uv tool or venv, conditional on enable; preserve executable exposure and dependency isolation.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### git-ui

- **lazygit** (`unversioned/distro-or-latest`, uncertain): Install lazygit from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### go

- **go** (`1.27.0`, gap): Install official Go 1.27.0 archive with verified digest, only when enabled; preserve native config without additional container permissions.
  Owner option: Owner may explicitly accept upstream Feature added ptrace/unconfined permissions; then existing feature_options can be used. No approval assumed.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### go-quality

- **goimports** (`v0.48.0`, gap): Run GOBIN=/usr/local/bin go install golang.org/x/tools/cmd/goimports@v0.48.0 in native build stage only when enabled, then copy the binary. Official Go Feature installGoTools bundles unrelated tools and does not independently pin this tool.
- **staticcheck** (`v0.7.0`, gap): Run GOBIN=/usr/local/bin go install honnef.co/go/tools/cmd/staticcheck@v0.7.0 in native build stage only when enabled, then copy the binary. Official Go Feature installGoTools bundles unrelated tools and does not independently pin this tool.
- **golangci-lint** (`v2.12.2`, gap): Run GOBIN=/usr/local/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 in native build stage only when enabled, then copy the binary. Official Go Feature installGoTools bundles unrelated tools and does not independently pin this tool.
- **govulncheck** (`v1.6.0`, gap): Run GOBIN=/usr/local/bin go install golang.org/x/vuln/cmd/govulncheck@v1.6.0 in native build stage only when enabled, then copy the binary. Official Go Feature installGoTools bundles unrelated tools and does not independently pin this tool.
- **gosec** (`v2.28.0`, gap): Run GOBIN=/usr/local/bin go install github.com/securego/gosec/v2/cmd/gosec@v2.28.0 in native build stage only when enabled, then copy the binary. Official Go Feature installGoTools bundles unrelated tools and does not independently pin this tool.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### go-release

- **goreleaser** (`2.17.1`, uncertain): Download verified architecture-specific assets from goreleaser/goreleaser at 2.17.1 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### infrastructure

- **opentofu** (`1.12.6`, uncertain): Download verified architecture-specific assets from opentofu/opentofu at 1.12.6 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.
- **ansible** (`14.3.1`, uncertain): Install ansible==14.3.1 using a native Dockerfile and isolated uv tool or venv, conditional on enable; preserve executable exposure and dependency isolation.
- **packer** (`1.16.0`, uncertain): Download verified architecture-specific assets from hashicorp/packer at 1.16.0 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.
- **podman** (`unversioned/distro-or-latest`, gap): Install podman/buildah/skopeo/slirp4netns/fuse-overlayfs from selected distro only when enabled; separately decide required user namespaces/cgroups rather than substituting Docker Feature.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### kubernetes

- **kubectl** (`1.36.4`, gap): Install kubectl 1.36.4 from verified official release assets only when enabled, using a native Dockerfile and preserving current mounts.
  Owner option: Owner may approve the upstream Feature unconditional named Minikube volume; existing feature_options then apply.
- **helm** (`4.2.4`, gap): Install helm 4.2.4 from verified official release assets only when enabled, using a native Dockerfile and preserving current mounts.
  Owner option: Owner may approve the upstream Feature unconditional named Minikube volume; existing feature_options then apply.
- **kustomize** (`5.8.1`, uncertain): Download verified architecture-specific assets from kubernetes-sigs/kustomize at 5.8.1 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.
- **k9s** (`0.51.0`, uncertain): Download verified architecture-specific assets from derailed/k9s at 0.51.0 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### latex

- **texlive-core** (`unversioned/distro-or-latest`, gap): Preserve frozen TeX Live 2025 mirror and exact package-group selections in a native multi-stage Dockerfile. Consider one minimal TeX Live Feature only after owner approves a reusable Feature API for independent group toggles; do not implement automatically.
- **texlive-recommended** (`unversioned/distro-or-latest`, gap): Preserve frozen TeX Live 2025 mirror and exact package-group selections in a native multi-stage Dockerfile. Consider one minimal TeX Live Feature only after owner approves a reusable Feature API for independent group toggles; do not implement automatically.
- **texlive-fonts** (`unversioned/distro-or-latest`, gap): Preserve frozen TeX Live 2025 mirror and exact package-group selections in a native multi-stage Dockerfile. Consider one minimal TeX Live Feature only after owner approves a reusable Feature API for independent group toggles; do not implement automatically.
- **texlive-biber** (`unversioned/distro-or-latest`, gap): Preserve frozen TeX Live 2025 mirror and exact package-group selections in a native multi-stage Dockerfile. Consider one minimal TeX Live Feature only after owner approves a reusable Feature API for independent group toggles; do not implement automatically.
- **texlive-code** (`unversioned/distro-or-latest`, gap): Preserve frozen TeX Live 2025 mirror and exact package-group selections in a native multi-stage Dockerfile. Consider one minimal TeX Live Feature only after owner approves a reusable Feature API for independent group toggles; do not implement automatically.
- **texlive-diagrams** (`unversioned/distro-or-latest`, gap): Preserve frozen TeX Live 2025 mirror and exact package-group selections in a native multi-stage Dockerfile. Consider one minimal TeX Live Feature only after owner approves a reusable Feature API for independent group toggles; do not implement automatically.
- **texlive-math** (`unversioned/distro-or-latest`, gap): Preserve frozen TeX Live 2025 mirror and exact package-group selections in a native multi-stage Dockerfile. Consider one minimal TeX Live Feature only after owner approves a reusable Feature API for independent group toggles; do not implement automatically.
- **texlive-music** (`unversioned/distro-or-latest`, gap): Preserve frozen TeX Live 2025 mirror and exact package-group selections in a native multi-stage Dockerfile. Consider one minimal TeX Live Feature only after owner approves a reusable Feature API for independent group toggles; do not implement automatically.
- **texlive-chemistry** (`unversioned/distro-or-latest`, gap): Preserve frozen TeX Live 2025 mirror and exact package-group selections in a native multi-stage Dockerfile. Consider one minimal TeX Live Feature only after owner approves a reusable Feature API for independent group toggles; do not implement automatically.
- **texlive-linguistics** (`unversioned/distro-or-latest`, gap): Preserve frozen TeX Live 2025 mirror and exact package-group selections in a native multi-stage Dockerfile. Consider one minimal TeX Live Feature only after owner approves a reusable Feature API for independent group toggles; do not implement automatically.
- **svg-inkscape** (`unversioned/distro-or-latest`, gap): Install inkscape from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### mermaid

- **mermaid-cli** (`11.16.0`, gap): Install @mermaid-js/mermaid-cli@11.16.0 from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable.
- **puppeteer** (`25.9.0`, gap): Install puppeteer@25.9.0 from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable. Preserve puppeteer@25.9.0 plus its compatible Chrome headless shell and runtime libraries; package alone is insufficient.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### node

- **node** (`26`, gap): Official Node image or verified Node tarball; independent pnpm/yarn/bun selection. Official Node Feature always installs Yarn and cannot preserve yarn=false.
- **pnpm** (`11.22.0`, gap): Install pnpm@11.22.0 from npm in a native Dockerfile only when enabled; retain pinned transitive/browser dependencies where applicable.
- **yarn** (`4.17.0`, gap): Use Corepack exact Yarn activation only when enabled. Official Node Feature has no yarn version or disable option.
- **bun** (`1.4.0`, uncertain): Download verified architecture-specific assets from oven-sh/bun at 1.4.0 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### preview-archive

- **chafa** (`unversioned/distro-or-latest`, gap): Install chafa from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **librsvg** (`unversioned/distro-or-latest`, gap): Install librsvg2-bin from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **poppler** (`unversioned/distro-or-latest`, gap): Install poppler-utils from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **timg** (`unversioned/distro-or-latest`, gap): Install timg from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **mupdf** (`unversioned/distro-or-latest`, gap): Install mupdf-tools from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **entr** (`unversioned/distro-or-latest`, gap): Install entr from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **p7zip** (`unversioned/distro-or-latest`, gap): Install p7zip-full from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **resvg** (`0.47.0`, gap): Download verified architecture-specific assets from linebender/resvg at 0.47.0 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### preview-enhanced

- **ffmpeg** (`unversioned/distro-or-latest`, gap): Install ffmpeg from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **ghostscript** (`unversioned/distro-or-latest`, gap): Install ghostscript from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.
- **rich** (`unversioned/distro-or-latest`, gap): Install python3-rich from the selected distro in a native Dockerfile only when enabled; preserve disabled-tool purge/absence guards where v0 supplied them. Pin the base image and package repository snapshot.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### python

- **python** (`3.14`, gap): Keep uv-managed Python 3.14 from official uv image/binary and explicit pip selection in native Dockerfile; verify false means either no pip CLI or owner-approved bundled runtime pip.
  Owner option: Owner may approve bundled pip and official Python Feature runtime installation.
- **uv** (`0.12.5`, uncertain): Copy /uv and /uvx from digest-pinned ghcr.io/astral-sh/uv at 0.12.5, only where enabled policy allows; upstream documents this path.
- **pip** (`unversioned/distro-or-latest`, gap): Use distro python3-pip/python3-venv only when enabled; decide whether pip supplied by the Python runtime is acceptable when false. Official Python Feature cannot omit bundled pip.
- **poetry** (`2.4.1`, uncertain): Install poetry==2.4.1 using a native Dockerfile and isolated uv tool or venv, conditional on enable; preserve executable exposure and dependency isolation.
- **pdm** (`2.28.2`, uncertain): Install pdm==2.28.2 using a native Dockerfile and isolated uv tool or venv, conditional on enable; preserve executable exposure and dependency isolation.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### release

- **shellcheck** (`0.11.0`, uncertain): Download verified architecture-specific assets from koalaman/shellcheck at 0.11.0 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.
- **hadolint** (`2.15.1`, uncertain): Download verified architecture-specific assets from hadolint/hadolint at 2.15.1 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### rust

- **rustc** (`1.98.0`, gap): Use rustup minimal profile, pinned toolchain 1.98.0 and explicit enabled components in native Dockerfile; exclude this component when disabled. Preserve existing container permissions.
  Owner option: Owner may explicitly accept upstream Feature added ptrace/unconfined permissions; then existing feature_options can be used. No approval assumed.
- **clippy** (`unversioned/distro-or-latest`, gap): Use rustup minimal profile, pinned toolchain 1.98.0 and explicit enabled components in native Dockerfile; exclude this component when disabled. Preserve existing container permissions.
  Owner option: Owner may explicitly accept upstream Feature added ptrace/unconfined permissions; then existing feature_options can be used. No approval assumed.
- **rustfmt** (`unversioned/distro-or-latest`, gap): Use rustup minimal profile, pinned toolchain 1.98.0 and explicit enabled components in native Dockerfile; exclude this component when disabled. Preserve existing container permissions.
  Owner option: Owner may explicitly accept upstream Feature added ptrace/unconfined permissions; then existing feature_options can be used. No approval assumed.
- **cargo-audit** (`unversioned/distro-or-latest`, gap): cargo install cargo-audit --locked in a native Rust build stage only when enabled; resolve and lock a version for reproducibility.
- **x86_64-cross** (`unversioned/distro-or-latest`, gap): Rust Feature targets option alone is insufficient: install gcc-x86-64-linux-gnu/libc6-dev-amd64-cross and preserve Cargo linker config; condition all three on enable.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### supply-chain

- **gitleaks** (`8.30.1`, uncertain): Download verified architecture-specific assets from gitleaks/gitleaks at 8.30.1 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.
- **osv-scanner** (`2.4.0`, gap): Download verified architecture-specific assets from google/osv-scanner at 2.4.0 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.
- **syft** (`1.50.0`, uncertain): Download verified architecture-specific assets from anchore/syft at 1.50.0 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.
- **grype** (`0.116.1`, uncertain): Download verified architecture-specific assets from anchore/grype at 0.116.1 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.
- **cosign** (`3.1.2`, uncertain): Download verified architecture-specific assets from sigstore/cosign at 3.1.2 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

### typst

- **typst** (`0.15.0`, uncertain): Download verified architecture-specific assets from typst/typst at 0.15.0 in a native Dockerfile, conditional on enable; avoid a new Feature wrapper unless reusable installation/configuration needs justify it.

**Decision:** Choose retain with the proposed native route, a reviewed upstream exception where listed, or remove for each row in this group.

## Evidence and limitations

- [Official/community registry](https://containers.dev/features), used only for publisher attribution and candidate discovery.
- [Frozen official Feature source](https://github.com/devcontainers/features/tree/9640551520736897481d83e92082186e4d812d50/src), install scripts and metadata inspected.
- [Frozen aibox source](https://github.com/projectious-work/aibox/tree/9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6/addons), all 41 recipes inspected.
- `spec/v1/evidence/V1-05/source-audit.json`: every row, its primary sources, candidate publisher/source/ref/digest, reasons, proposals, drift.
- `spec/v1/evidence/V1-05/oci-evidence.json`: independently fetched GHCR manifests/metadata, including conditional upstream candidates.

Go and Rust Features add `SYS_PTRACE` and `seccomp=unconfined`; no suppression option was found. Kubernetes Feature always creates `minikube-config` volume at `/home/vscode/.minikube`. Node Feature always installs Yarn without independent pin/disable. Python Feature always provides pip. Exact install tests are still required before production mapping.

---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261029_1034-ActiveJay-v1-05-catalog-additions-vibe-uv
  created: '2026-10-29T10:34:59+00:00'
spec:
  title: 'V1-05 catalog additions: vibe, uv-managed Python, typescript, zig, security
    groups'
  state: accepted
  decision: 'V1-05 catalog additions and reshaping: (1) add single Feature ai/vibe
    (mistral-vibe via uv tool, depends on uv). (2) python is a semantic group managed
    entirely by uv: uv (pinned binary), python (exact CPython versions installed with
    uv python install into a system location, python/python3 on PATH, bundled pip),
    poetry and pdm via uv tool; aider, hermes, tau, vibe, zensical, csvkit and rich
    also install via uv tool so no Debian python3 is pulled in unless requested. (3)
    add semantic typescript group: tsc (typescript, depends on node), biome, oxlint.
    (4) add semantic zig group: zig from ziglang.org minisign-verified tarballs, zls
    depending on zig with matching version. (5) add security area with two semantic
    groups: security/crypto (age, sops, gnupg, openssh-client, minisign, openssl)
    and security/scan (v0 supply-chain rows gitleaks, osv-scanner, syft, grype, cosign,
    plus semgrep); trivy excluded. Feature paths are explicit catalog data, not derived
    from group names.'
  context: Owner review of the V1-05 gap list (follows DEC-20261029_0835-HumbleBridge).
  rationale: uv gives one reproducible Python foundation; TypeScript scope is limited
    to project-independent tools because project lockfiles stay authoritative; security
    split keeps light daily-use crypto tools separate from heavy audit scanners; trivy
    overlaps syft/grype and its release pipeline was compromised in March 2026 (CVE-2026-33634).
  consequences: New Features have no v0 census row and are tracked as catalog additions.
    v0 supply-chain rows remain traceable under security/scan. Path derivation by
    hyphen split is dropped.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-29T10:34:59+00:00'
---

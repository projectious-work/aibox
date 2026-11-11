# 13. Security and trust-boundary implementation

This chapter specifies R-AUTHORITY and R-POLICY. There are two distinct
agents: an external operator agent acting for a human over host lifecycle,
and a repository/workspace agent inside the Dev Container. The latter may
edit container-local tmux, themes and project files, but never gains host
lifecycle or host-file authority. “Same binary” is not “same authority.”

## Actors, assets and boundaries

| Boundary | Trusted authority | Untrusted or less-trusted input | Mandatory control |
|---|---|---|---|
| Host operator process → runtime | Human-installed policy and OS user | Project definition, repository hooks, MCP client request | Exact root/context/operation allowlist, frozen digest, hook review, least privilege |
| Container process → workspace | Container user | Repository content, agent prompts, tool output | No host socket, operator endpoint or host credentials; local-only operations |
| Workspace → host filesystem | Explicit native mounts only | Path strings in definition and migration | Realpath, no traversal/symlink escape, host policy denies forbidden mounts |
| Image/Feature → supply chain | Approved digest/lock and provenance | Registry tags, download scripts, package sources | Pin, verify, SBOM/scan/license review, reproducible test |
| CLI/MCP → user | Typed result/evidence schema | Child stdout, error text, titles and guides | Redaction, bounded output, no terminal control injection |

The operator policy lives outside the repository and is configured by the
human at server launch. It specifies allowed canonical roots, runtime
endpoints/contexts, operations, hook/trust approval mode, expiry and evidence
retention. Repository config can request but cannot grant `privileged`, host
PID/network, Docker/Podman socket, broad host mount, `initializeCommand`,
unsafe device or arbitrary host command. Such requests are denied or require
an out-of-band operator decision binding a digest; a checked-in flag is not
consent. The first-run experience is read-only knowledge mode. Even after
operator mode is enabled, `remove` and `rebuild` are separate capabilities.
Direct use of the upstream CLI by a human with host runtime credentials is
outside the wrapper's enforcement; installation guidance therefore requires
the human to review untrusted Dev Container hooks/mounts before direct
build/up as well. aibox never claims its policy mediates an independently
invoked native tool.

Runtime socket and operator credentials are never mounted into the container.
Work credentials have the separate, opt-in native transfer paths and
qualification gate in chapter 17. A persistent home is not implicit host
credential access; a read-only bind, forwarded agent or broker identity is
reviewed independently of workspace startup.
The local MCP server is configured with only local resources/tools; it does
not proxy an operator server or expose network transport. The host process
rejects calls whose actor identity is unestablished; it never trusts an
`actor` field supplied by a tool caller as proof of identity. An operator
agent can advise a human but cannot impersonate a local container process
to bypass policy. A local agent may edit native tool config in its container
and call `refresh_workspace`, which validates ownership and reloads only
container-local tools.

## Threats and negative controls

- Repository prompt injection or README instructions cannot override MCP
  policy. Tool descriptions and how-to content are versioned product assets;
  an untrusted repository cannot register new operator tools or rewrite them.
- Project paths, names, branch labels, terminal titles, filenames and child
  stderr may contain shell/control sequences or secrets. Use argument arrays,
  text escaping, bounded lengths and redaction before display/persistence.
- TOCTOU between inspect and delete is handled by lock plus revalidation of
  endpoint, context, native ID and workspace labels. Never delete by display
  name or wildcard. A replaced resource gets `identity_changed`.
- User-owned files and generated outputs are distinct. Refresh refuses a
  symlink/special-file target, unexpected owner, writable parent escape or
  unclaimed content. Migration follows chapter 12's backups and journal.
- Hook execution is an acknowledged host-risk: inventory all host-executed
  lifecycle commands and mounts, bind approval to the exact input digest,
  and re-review when the digest changes. `doctor` never executes hooks.
- Auth tokens, provider state, audio sockets and private MCP entries are not
  copied to image layers, Git, telemetry, guides or logs. Evidence stores
  hashes, resource IDs and redacted excerpts with explicit local retention.
- A dependency's official status is not a security approval. Feature refs,
  image base, Go modules, native tools and updater scripts receive digest,
  license, provenance and vulnerability checks before release.

Security tests run local/host mode on the same binary, try forged mode flags,
operator MCP config inside the container, repository-controlled policy,
malicious JSONC/hook/path strings, symlink races, changed runtime contexts,
title injection, secret-bearing child errors and cancellation during removal.
Acceptance requires denial **before** any side effect, exact attribution of
partial effects, and a repeat run that demonstrates no widened authority.
The release includes a security policy and supported vulnerability-reporting
path; disclosure handling is separate from ordinary doctor findings.

# 17. Secret and credential transfer into the Dev Container

**R-SECRETS:** Secret entry is an explicit user/operator decision, not an
automatic consequence of starting a workspace or persisting its home. The
container-local agent may use credentials made available by the user for its
repository work, but it may neither request host lifecycle credentials nor
turn a project file into permission to mount a host secret. The host operator
reviews the exact native mount/env/hook request and binds authorization to the
input digest (chapter 13). No host runtime socket, operator MCP endpoint or
operator credential enters the container.

## Recovered design and v1 boundary

The accepted company decisions
`DEC-20260821_1830-LivelyCrow-standardize-aibox-secret-requests-across-portable`
and `DEC-20260823_1428-CalmWolf-define-aibox-v1-x-rewrite-architecture`
identified several distinct mechanisms: SOPS with age as a portable encrypted
file source; OpenBao as a broker for shared/long-running environments; native
Kubernetes external-secret mounts for Kubernetes targets; and a future
attestation-gated KBS/Trustee provider for confidential workloads. They also
favored logical secret names, file-first delivery, redaction and no plaintext
secret in a committed configuration or deployment artifact.

Those decisions arose in the earlier broad deployment-compiler design, which
this Dev Container CLI rewrite does not inherit wholesale. V1 retains the
security principles and user-facing transfer options, but does **not** create
a universal secret-provider API, run OpenBao, implement Kubernetes/CSI, or
build a custom attestation broker. The native `devcontainer.json`, Compose,
Dockerfile/BuildKit, provider tools and user-approved mounts remain owners.
Future remote/Kubernetes/confidential-computing products may reuse the old
logical-request model only after a separate boundary and review. This is a
deliberate narrowing, not an assertion that SOPS or OpenBao are insecure.

## Supported transfer paths and their limits

| Path | Intended use | Required behavior |
|---|---|---|
| Persistent container home | Harness login state, user-managed tool config and caches created *inside* the container | Default project-scoped volume at `/home/aibox`; custom user/home or explicit bind is supported (chapter 8). Treat the entire volume as sensitive data. Persistence is not a host secret transfer and does not justify mounting the host's entire home. |
| Explicit read-only file bind | Selected existing host credential/config file or directory | Owner chooses exact daemon-host source and least-scope container target through native `mounts`/Compose; use `readonly` where supported. Do not symlink-hop or mount broad `.ssh`, `.config` or `$HOME` automatically. Check target user access and provider-specific file permissions. |
| Agent/socket forwarding | SSH agent or other provider's supported socket/forwarding mechanism | Prefer a short-lived agent over copying a private key when the runtime supports it. Request/approve the exact socket mount; document that a process with socket access can act as the user while connected. No Docker/Podman socket exception. |
| Native environment variable | Tools that genuinely require env input | Use an explicit user-selected native `containerEnv`, `remoteEnv`, Compose `env_file` or runtime provider path. Never put secret literals in committed JSON, templates, image metadata, CLI arguments or logs. Warn that environment values may be visible to other container processes, runtime inspection or process descendants. Prefer a file/agent where available. The aibox process `--env-file` in chapter 16 does not automatically inject values. |
| SOPS/age encrypted file | Portable, user-mediated local/headless source | Keep ciphertext in an approved location; decrypt with standard SOPS/age outside image build and only after owner approval. Deliver a narrowly scoped file or provider-supported runtime mount with private permissions, clean up host plaintext on failure/stop according to an explicit lifecycle plan, and never place decryption keys in the image or project config. Aibox does not implement cryptography or silently decrypt on `up`. |
| OpenBao or another established broker | Shared, long-running or rotating secrets where an operator already runs a broker | Use its supported CLI/agent/sidecar and workload-appropriate identity, lease/renewal and file rendering; do not pass an operator root token or mount the operator's broker credentials. Broker outage/expiry must be visible, never silently replaced by stale plaintext. Aibox may document/integrate a standard Feature/Compose fragment after qualification, but is not the broker. |
| BuildKit build secret/SSH forwarding | Build-time fetch only | Declare via native Dockerfile/BuildKit and approved build invocation. Do not use `ARG`, `ENV` or `COPY` to persist a secret in an image layer. Build secrets are not runtime credentials and are not copied to persistent home. |

Kubernetes Secrets Store CSI and KBS/Trustee remain recorded historical
options, **outside this v1 local Dev Container target**. Their inclusion in
the earlier decision does not make them delivery requirements here.

## Native configuration and authority

Each starter offers a credential-free path. No default template mounts a
host private directory or asks for a secret. Document opt-in examples for a
single read-only SSH configuration/key file, an SSH agent socket, a private
provider env file and a SOPS/age-to-file workflow, with the risk and cleanup
shown next to each. Examples must use placeholders, never test credentials.
The standard home volume can retain authentication established *inside* the
container, but moving an existing host auth directory into it requires an
explicit, scoped migration; do not seed it from a host home by default.

The proposed `customizations.aibox` schema contains no secret values,
provider locator authority, decryption key, or host source path. A project
may declare a non-secret logical need in its own tool-native documentation,
but the user/operator chooses the concrete native delivery configuration.
Checked-in `devcontainer.json`, a Feature, lifecycle hook or in-container MCP
request cannot authorize itself. The operator policy classifies secret mounts,
agent sockets, env injection and host decrypt hooks as sensitive, requires a
review of exact source/target/permissions/lifetime, and denies an unapproved
change before side effects. `doctor` reports missing/unreadable credentials as
`missing`/`not_authorized`/`unknown` without reading values or asserting that
the secret is valid. An in-container doctor is local-only.

Secret values never enter operation requests/results, receipts, logs,
diagnostic bundles, generated UX files, image layers, release artifacts or
the MCP how-to corpus. Native tools may have their own logs and caches; the
qualification gate inspects those too. Redaction is defense in depth, not a
reason to print a value. Removal of a container retains its named home and
does not revoke provider tokens; the docs give the owner an explicit revoke,
rotate, archive and secure-delete/retention checklist. A bind-mounted host
file remains owned by the host user, not by aibox cleanup.

## Acceptance and release gate

**AC-SECRETS:** Exercise each supported path above with canary credentials and
negative fixtures: no-secret starter; precise read-only file mount; socket
present/absent; env inheritance and process inspection; SOPS key missing,
decrypt failure and cleanup; broker expiry/rotation/outage; BuildKit layer
scan; custom user/home permissions; remote-daemon bind path; malicious
project-requested broad mount/hook; changed input digest; and disabled
provider. Assert no canary in Git, image/layers, generated files, CLI/MCP
stdout/stderr, logs, receipts, diagnostics, snapshots or documentation.
Verify the container-local agent cannot escalate from a delivered work
credential to host lifecycle authority. Record per-mode platform support,
owner approval and cleanup evidence. Unqualified modes are documented as
manual/native workflows, not advertised as automatic aibox features.

Primary specifications: [Dev Container Features](https://github.com/devcontainers/spec/blob/main/docs/specs/devcontainer-features.md),
[Dev Container Templates](https://github.com/devcontainers/spec/blob/main/docs/specs/devcontainer-templates.md),
[Docker bind mounts](https://docs.docker.com/engine/storage/bind-mounts/),
[Docker volumes](https://docs.docker.com/engine/storage/volumes/),
[Docker Build secrets](https://docs.docker.com/build/building/secrets/),
[SOPS](https://github.com/getsops/sops).

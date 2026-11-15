# 20. Personas, user stories and solution paths

This chapter turns the product contract into testable end-to-end journeys.
It defines the intended v1 behavior and user-visible outcomes. Command names,
authority boundaries and result contracts follow chapters 3, 9, 10, 12 and 13.
Each journey is a product acceptance scenario, independent of delivery order.

## Persona system prompts and ACTION story tasks

Use the persona system prompt to establish a consistent simulated person or
agent. Then provide its ACTION task prompt as the task message, selecting
**create** or **review** mode. In create mode, give the product contract but
withhold existing stories so the persona starts from its own needs. In review
mode, provide the existing stories and ask for concrete corrections. These
synthetic personas are design probes, not evidence of actual customer research.
ACTION means Act as, Context, Task, Iterate Output and Netiquette; the two
prompt layers should remain separately reusable.

### P1 — Independent developer (human)

Works on several repositories from a personal laptop and alternates between
terminal and VS Code. Understands Dev Containers but does not want to learn a
second build configuration language. Wants a fast first start, reproducible
tool choices, familiar editor attachment and persistent personal state. Owns
project configuration, but does not assume that editing it grants permission
for host-executed hooks or unsafe mounts. Uses aibox as the normal lifecycle
entry point; direct upstream CLI use is a compatibility property to verify.

Persona system prompt:

```text
<Role>
Name: Mira Novak.
Role: Independent software developer and occasional consultant evaluating aibox
for projects she personally maintains; speak as a potential user, not a product
manager or implementation team member.
Main goal: Reach a reproducible, comfortable workspace quickly without learning
a second build language or losing personal state.
Working-style energies: Blue 40%, Green 20%, Red 25%, Yellow 15% (total 100%).
Traits and voice: Curious and technically literate, pragmatic about setup time,
direct when a command sequence feels redundant, and concrete about what would
make her trust a tool. Use first-person, plain technical language.
</Role>
<Organization>
Mira works independently across client and personal repositories. She owns her
project files and laptop workflow, uses VS Code and terminal interchangeably,
and may depend on a separate operator for host policy in managed environments.
</Organization>
<Tasks>
Assess a starter and an existing repository; choose only needed Features;
build, start, attach and rebuild through aibox; test persistence and standard
Dev Container interoperability. Surface adoption friction and missing feedback.
Outputs are first-person needs, acceptance outcomes, exact user actions and
specific objections—not internal implementation progress reports.
</Tasks>
<Important Notes>
This is a simulated user, not a claim of customer interviews. Repository edits
do not approve host hooks or mounts. Use aibox as the ordinary lifecycle entry
point; direct upstream CLI runs belong in a separate compatibility test.
</Important Notes>
```

> **Act as:** Be this independent developer, speaking in first person and
> prioritizing a working project over internal architecture.
>
> **Context:** You clone an existing repository with a native
> `.devcontainer/devcontainer.json`, or copy a documented starter. You want
> selected Features, an explicit user/home mount, and a running environment
> you can enter from VS Code or a terminal. aibox delegates native build/up
> behavior to the pinned Dev Container CLI but adds scoped target, input and
> policy checks.
>
> **Task:** In **create** mode, write first-person stories from your own setup,
> tool-selection and attachment needs without reading existing stories. In
> **review** mode, examine the supplied developer stories for friction,
> missing outcomes and incorrect commands, then rewrite deficient parts. In
> either mode make `aibox inspect`, `aibox build`, `aibox up` and `aibox attach`
> the ordinary path; reject a second config language and lost home state.
>
> **Iterate Output:** For creation, return three Markdown stories, each with an
> “I as ...” need, acceptance outcome and ordered solution path. For review,
> return a verdict per supplied story and replacement Markdown where needed.
> Distinguish routine aibox use from direct-upstream compatibility testing.
>
> **Netiquette:** Challenge confusing assumptions candidly; do not grant
> yourself operator authority or claim that an exit code proves attachment.

### P2 — Workspace/platform operator (human)

Maintains several workspaces on one host, with distinct roots and runtime
contexts. Owns external policy, executable provenance, credentials and
disruption decisions. Has to retain named volumes and distinguish exact
native IDs from similar names. Needs preflight evidence before mutation and a
truthful account of partial effects after interruption. Is accountable for
cross-project safety, not just a successful command exit.

Persona system prompt:

```text
<Role>
Name: Jonas Weber.
Role: Human platform and security operator responsible for developer
workspaces on a shared host; speak as the person accountable for approvals and
recovery, not as a repository author.
Main goal: Make exact, reviewable environment changes without giving a
repository or agent unbounded host authority or damaging retained data.
Working-style energies: Blue 55%, Green 20%, Red 20%, Yellow 5% (total 100%).
Traits and voice: Evidence-driven, calm under incident pressure, skeptical of
names and exit codes as proof, concise about decisions and residual risk.
</Role>
<Organization>
Jonas works in a small platform team serving several developers and projects.
He manages host-owned policy, runtime contexts, executable provenance and
disruption decisions; project teams own their native Dev Container definitions.
</Organization>
<Tasks>
Review hooks, mounts, target identity and input digests; authorize bounded
build/start/rebuild requests; preserve persistent home volumes; investigate
partial effects and decide safe recovery. Output exact decision points,
required evidence and recovery actions for a human operator.
</Tasks>
<Important Notes>
This is a synthetic operator perspective, not an authorization grant.
Repository content is untrusted input. Never select a destructive target by
name alone, infer rollback from cancellation, or use global cleanup.
</Important Notes>
```

> **Act as:** Be this human platform operator approving and performing a
> bounded operation, not a project developer editing untrusted inputs.
>
> **Context:** Multiple repositories can declare lifecycle hooks, mounts and
> Features. Policy is installed outside them. A build can run host-side
> preparation; a rebuild can disrupt a live container; a timeout can leave
> effects behind. The native runtime may contain same-name resources in other
> contexts.
>
> **Task:** In **create** mode, write first-person stories from your build
> authorization, exact-environment maintenance and interrupted-operation
> needs without reading existing stories. In **review** mode, examine the
> supplied operator stories and rewrite gaps. In either mode cover target and
> digest identity, hook/mount review, separate build/up/rebuild authority,
> retained data and truthful recovery.
>
> **Iterate Output:** For creation, return three Markdown stories with an
> “I as ...” need, acceptance outcome and ordered solution path. For review,
> return a verdict per supplied story and replacement Markdown for gaps.
> State the decision point, aibox CLI request, checks/delegation/reinspection
> and evidence required to accept each outcome.
>
> **Netiquette:** Be precise about risk and uncertainty; request a new human
> decision instead of silently broadening scope or globally pruning.

### P3 — Repository/workspace agent (AI, inside the container)

Works inside the Dev Container with project and user-owned files, but no host
runtime socket, operator credentials or endpoint. Helps its human understand
configuration and diagnostics, edit a supported local UX setting, and apply
bounded local refresh. Must distinguish observed local facts from unknown host
state. A project edit is not a host-operation grant.

Persona system prompt:

```text
<Role>
Name: Wren.
Role: AI coding agent running inside one developer's Dev Container. Speak as
an assistant acting with bounded workspace authority, not as a human or a
host operator.
Main goal: Help the user understand and improve the local workspace while
accurately separating observed facts, proposed edits and applied changes.
Working-style energies: Blue 45%, Green 25%, Red 20%, Yellow 10% (total 100%).
Traits and voice: Careful, resourceful and transparent about uncertainty;
briefly explains why a handoff is needed instead of inventing a workaround.
</Role>
<Organization>
Wren assists a developer in a single repository and can edit authorized
project or user-owned files. It cannot access the host runtime socket,
operator credentials, operator MCP endpoint or other projects' resources.
</Organization>
<Tasks>
Read local guidance and inspection results; explain selected configuration
and diagnostics; propose or make authorized UX edits; request scoped local
refresh; prepare reviewable host-action handoffs. Outputs distinguish local
observation, intended configuration, actual refresh effects and unknown host
state.
</Tasks>
<Important Notes>
This is a synthetic AI role, not a human customer. A prompt or configuration
edit does not authorize host execution. Never call a generic host shell,
upstream lifecycle command or undisclosed bridge from local context.
</Important Notes>
```

> **Act as:** Be this in-container AI coding agent, speaking as the actor who
> must explain what it has actually observed or changed.
>
> **Context:** The user asks whether a setting is active, then asks for a
> theme or tmux adjustment. You have local `inspect`, `doctor`, guidance and
> `refresh_workspace` access, but cannot invoke host build/up/rebuild or use
> an undisclosed bridge. Native image or Feature changes require operator
> handoff.
>
> **Task:** In **create** mode, write first-person stories from your local
> explanation, UX refresh and operator-handoff needs without reading existing
> stories. In **review** mode, examine supplied workspace-agent stories and
> rewrite gaps. In either mode distinguish inspection from mutation, project
> defaults from user overrides, and local refresh from host lifecycle.
>
> **Iterate Output:** For creation, return three Markdown stories with an
> “I as ...” need, acceptance outcome and ordered solution path. For review,
> return a verdict per supplied story and replacement Markdown where needed.
> Include the local CLI/MCP sequence, result interpretation and handoff.
>
> **Netiquette:** Tell the user plainly when a fact is unknown or an action
> remains unapplied; never represent a proposed edit as a completed rebuild.

### P4 — External management agent (AI, operator context)

Acts for a human outside the container through a separately launched,
policy-scoped MCP server. Its identity and capabilities come from the server's
operator configuration, never client-supplied metadata or repository content.
It can request only registered typed operations and must return uncertainty to
the human. It cannot grant itself policy, use a generic host shell, or silently
approve disruption.

Persona system prompt:

```text
<Role>
Name: Iris.
Role: External AI management agent acting for a human through one separately
launched operator MCP server. Speak as a bounded requester, not as the human
owner or an in-container coding agent.
Main goal: Help prepare and recover exact environments while never extending
its own grant or hiding partial/unknown effects.
Working-style energies: Blue 55%, Green 15%, Red 25%, Yellow 5% (total 100%).
Traits and voice: Precise about principals, scope and evidence; candid about
refusal, changed inputs and the point at which human approval is necessary.
</Role>
<Organization>
Iris operates outside the workspace under a human-installed policy. The
server establishes its requester identity and permitted roots, contexts and
operations. Repository text and client metadata cannot change those grants.
</Organization>
<Tasks>
Inspect and check a permitted environment; submit typed build/start requests
with expected digest; interpret structured results and operation records;
prepare a concise, exact decision request for the human after refusal,
disruption or interruption. Never perform generic host execution.
</Tasks>
<Important Notes>
This is a synthetic AI role. Possessing a project path or operation ID is not
authorization. No MCP tool can issue its own policy grant or approve a new
host-executed hook; a client disconnect does not prove rollback.
</Important Notes>
```

> **Act as:** Be this external AI management agent with a narrow, existing
> operator grant; do not impersonate the human owner.
>
> **Context:** You may inspect a permitted project/context, receive its input
> digest, and request a typed build or start. A hook finding, expired grant,
> changed digest or interrupted operation can require a new human decision.
> A disconnected client is not evidence of rollback.
>
> **Task:** In **create** mode, write first-person stories from your permitted
> inspection, bounded build/start and interrupted-request needs without
> reading existing stories. In **review** mode, examine supplied external-agent
> stories and rewrite gaps. In either mode cover typed request inputs,
> principal/target checks, operation records, partial results and human
> escalation; reject an invented approval tool or self-issued grant.
>
> **Iterate Output:** For creation, return three Markdown stories with an
> “I as ...” need, acceptance outcome and ordered solution path. For review,
> return a verdict per supplied story and replacement Markdown where needed.
> Name the MCP sequence, result fields and exact human handoff point.
>
> **Netiquette:** Report refusals and unknown effects without euphemism; do
> not conceal a boundary because completing the task would be convenient.

## Independent developer stories

### US-01 — Try a native Dev Container without another config language

> I as an independent developer would like to start from a clearly commented `devcontainer.json` and change it using the Dev Container format I already know, so that I can try aibox without translating my project into a second configuration language.

Acceptance outcome: the starter is valid JSONC, explains supported aibox options and leaves unused options commented. Inspection identifies the selected native configuration, declared inputs and digest. Build and start delegate to the pinned Dev Container CLI only after target, policy and input checks; a refusal or changed input is reported rather than treated as a successful environment.

Solution path:

1. Copy the documented `minimal` Template files into my project.
2. Edit `.devcontainer/devcontainer.json`, enabling only the native properties and aibox options I need. Feature-specific settings belong to their selected Features.
3. Run `aibox inspect --project PATH --format json` to review the selected definition, declared inputs and input digest.
4. Run `aibox build --project PATH`, then `aibox up --project PATH`. Review each result: aibox checks the exact target and authorized inputs, delegates native build/start behavior to the pinned Dev Container CLI, then reports identities actually observed. If checks refuse or inputs changed, resolve the issue and inspect again before retrying.
5. Attach through `aibox attach --project PATH` or VS Code's standard attach-to-running-container flow. Direct upstream use is a separate interoperability test, not the normal aibox journey.

### US-02 — Add only the tools I choose

> I as an independent developer would like to select tools through documented Dev Container Features and their own options, so that I know what installs each tool and can keep the environment focused.

Acceptance outcome: only Features in the active native `features` object are requested; commented examples are inert. Inspection identifies selected Feature references and options with the declared input digest. The pinned Dev Container CLI resolves and installs Features, while aibox reports the actual build/start outcome rather than treating selection as proof of installation.

Solution path:

1. Review the commented starter catalog and edit the native `features` object using each chosen Feature's documentation.
2. Enable any desired `customizations.aibox` workspace presentation settings separately; these configure UX, not tool installation.
3. Run `aibox inspect --project PATH` to review selected Feature references, options, inputs and digest.
4. Run `aibox build --project PATH` and `aibox up --project PATH`; review success, refusal or failure in the structured results. aibox validates target and authority; the pinned upstream CLI resolves and installs Features without a separate user-run installer.

### US-03 — Attach and retain my personal state

> I as an independent developer would like to enter the running workspace from VS Code or a terminal and retain my home across rebuilds, so that aibox fits my existing Dev Container workflow.

Acceptance outcome: the native starter mounts a project-scoped named volume at the configured remote user's actual home. The aibox-started container can be entered through aibox or standard Dev Container tooling. After an authorized rebuild, the result and follow-up inspection confirm the same volume identity and retained files; failed or interrupted work is not represented as retention success.

Solution path:

1. Check the starter's `mounts`, `remoteUser` and image user; confirm that the named volume targets that user's home and is scoped to this project.
2. Run `aibox inspect --project PATH`, `aibox build --project PATH` and `aibox up --project PATH`; review the reported target and resource identities.
3. Attach with `aibox attach --project PATH` or VS Code's standard attach-to-running-container flow, then place a test file in the mounted home.
4. Request `aibox rebuild --project PATH` with the explicit disruption acknowledgment required by active policy. Acknowledgment alone does not override policy; changed inputs or targets require fresh review.
5. After completion, attach and verify the test file and named volume identity in the reported and re-inspected state. Test direct upstream compatibility separately.

## Human platform operator stories

### US-04 — Review before authorizing a host build

> I as a platform operator would like to review the exact project inputs, host hooks, mounts and runtime context before authorizing a build, so that repository content cannot silently acquire host access.

Acceptance outcome: aibox performs a build only after operator policy permits the canonical root, runtime context, build operation, executable provenance and digest of the declared control inputs. Changed declared inputs, ambiguous identity or unapproved hook/mount risk stop it before effects; the digest is not a claim to cover arbitrary build-context or shell dependencies.

Solution path:

1. Configure operator-owned policy outside the repository: allowed roots, contexts, operations, executable provenance and hook/mount approval. Repository declarations are inputs, not grants.
2. Run `aibox inspect --project PATH` and read-only `aibox doctor --project PATH` in the selected operator context. Review a read-only snapshot of all host-executed code, resolved configuration, mount sources/targets, runtime context and input digest without executing hooks.
3. Approve any host-side hook or risky mount against that exact target and digest with bounded expiry; changed inputs require a fresh decision.
4. Run `aibox build --project PATH`. aibox verifies policy and target, takes a scoped operation lock, rechecks inputs, delegates to the pinned Dev Container CLI, then reinspects and reports image/build evidence. Build neither starts nor replaces a running container.

### US-05 — Recover after an interrupted operation

> I as a platform operator would like to learn what an interrupted start actually changed before retrying it, so that a timeout does not leave me guessing about containers or duplicate resources.

Acceptance outcome: an authorized operation record binds the principal, operation ID, canonical project, runtime context, native resource identity and input digest. It distinguishes last confirmed step, observed effects and unknown effects; inspection never resumes execution.

Solution path:

1. Keep the operation ID returned by `aibox up`, `aibox build` or `aibox rebuild`, along with target and digest.
2. After interruption, run `aibox operation status ID --project PATH` in the same authorized scope. Read the last confirmed step, observed resources and explicitly unknown effects; status does not retry or resume.
3. Run `aibox inspect --project PATH` to reconcile native state against the exact context and IDs. A lost client response is not proof of rollback.
4. Decide whether recovery or retry is appropriate, then submit a new bounded aibox request. It rechecks policy, input digest and target identity; changed inputs or unresolved disruption require renewed human review.

### US-06 — Rebuild only one identified environment

> I as a platform operator would like to rebuild one identified environment with disruption disclosed and persistent data retained, so that I can apply image changes without affecting another project or losing a developer's home.

Acceptance outcome: aibox rebuilds only the confirmed canonical project, configuration, runtime endpoint/context and native resource IDs after explicit disruption acknowledgment. It builds before stopping the old environment, then replaces only the primary resource and exclusively owned sidecars. Persistent volumes and bind data remain intact and are reported afterward. A failure before removal leaves the old environment running; a failure after removal is reported as partial with recovery guidance. Same-name resources in other contexts remain untouched.

Solution path:

1. Run `aibox inspect --project PATH` in the selected operator context. Confirm canonical root and config, runtime endpoint/context, exact primary ID, owned sidecar membership, input digest, and persistent volume or bind source and mount target. Names and labels alone are insufficient.
2. Review the proposed changes and retained-data plan, then give explicit disruption acknowledgment for this exact target and digest. Changed inputs, endpoint/context or resource membership require renewed review.
3. Run `aibox rebuild --project PATH` with the reviewed identity, expected digest and disruption acknowledgment. aibox checks policy, acquires the scoped lock, revalidates identity and inputs, and performs a build preflight before stopping or removing anything.
4. If preflight succeeds, aibox stops and removes only the exact primary and exclusively owned sidecars, then delegates normal `up`. That step may build again; preflight is not an atomic image swap.
5. Reinspect and accept only a result reporting actual final image/resource identities, retained-data references and unknown effects. Failure before removal leaves the prior environment running; failure after removal is partial and requires operator recovery using retained data and the pinned previous definition. Volume deletion and global pruning are never part of rebuild.

## AI workspace agent stories

### US-07 — Explain the selected workspace from inside it

> I as the workspace agent would like to inspect the selected Dev Container definition and declared project inputs, so that I can explain what the project requests while keeping local observations separate from host runtime state.

Acceptance outcome: I report the selected definition, declared inputs, input digest, observation time and eligible local findings with provenance. Local results do not invent host runtime identity. Skipped, unavailable and unauthorized checks remain distinct from passed checks.

Solution path:

1. Call local `inspect_workspace` or run `aibox inspect --project PATH --format json` from the container; record the selected definition, declared inputs, digest, provenance and observation time returned.
2. Call `check_workspace` or run `aibox doctor --project PATH` from the container for eligible read-only local checks. Report each finding's source, evidence and state: `passed`, `failed`, `skipped`, `unavailable` or `not_authorized`. Never execute lifecycle hooks or query host runtime from local context.
3. Distinguish active values from commented examples. Report effective-value provenance only when returned, and describe declared intent as intent rather than proof that the host applied it.

### US-08 — Refresh a local presentation choice

> I as the workspace agent would like to adjust a supported user-owned theme or tmux setting and refresh only local presentation, so that I can help the user iterate without rebuilding the container.

Acceptance outcome: the user can identify the edited project default or user-owned override, changed managed files and any required rebuild or session restart. Refresh changes only owned local output, refuses unsafe or unowned targets, and preserves personal and processkit-owned MCP configuration. No host lifecycle effect is implied.

Solution path:

1. Read task-sized guidance for the requested setting: exact path/key, ownership, validation, activation and rollback. Ask the user to choose between project default and user override when ambiguous.
2. Edit only authorized local UX configuration; do not change image, Feature, mount or lifecycle-hook inputs as part of a presentation refresh.
3. Call local `refresh_workspace` with canonical project root, expected UX input digest and `theme`, `tmux` or `all` scope, or run the equivalent `aibox refresh`. It validates ownership and activates managed local output, then reports `changedFiles`, `rebuildRequired` and `sessionRestartRequired`.
4. If refresh refuses, fails or requires another action, explain what remains unapplied. Host build/up/rebuild is a separate human-authorized operator action; do not claim it occurred.

### US-09 — Hand off a lifecycle request without host authority

> I as the workspace agent would like to prepare a reviewable build or start request for the human to submit through an authorized operator, so that host authority stays outside the container.

Acceptance outcome: the user receives a request naming the selected definition, requested operation, input digest and unresolved risks. The local agent makes no host lifecycle call; only a structured operator result establishes the outcome.

Solution path:

1. Inspect locally and summarize the selected definition, declared inputs, digest and limits of local observation; do not invent host endpoint or runtime state.
2. Prepare a human-reviewable request for `build` or `up` with project/config selectors, inspected digest and hook/mount concerns, without credentials. A project edit is not approval.
3. Give the request to the human for an approved operator channel. Do not call a host bridge, runtime socket, generic host shell or upstream lifecycle command from local context. The authorized aibox operator rechecks policy and identity before delegating to the pinned upstream CLI.
4. Report the request as pending until the structured operator result returns. If no result is available, execution is unknown, not successful.

## External AI management agent stories

### US-10 — Inspect an authorized target before proposing action

> I as an external management agent would like to inspect a workspace permitted by the operator's policy and review its readiness findings, so that I can tell the human what the server observed before proposing a build or start.

Acceptance outcome: `inspect_environment` returns the operator target, input digest, effective configuration/provenance and observation time; `check_environment` returns findings and completeness. Each read is authorized under its applicable capability and scoped grant. `passed`, `failed`, `skipped`, `unavailable` and `not_authorized` remain distinct; findings do not grant mutation authority.

Solution path:

1. Connect through the separately launched operator stdio MCP server. Use the requester principal established by its operator configuration; client-supplied actor metadata is not identity proof.
2. Call `inspect_environment` with `schemaVersion="aibox.operation-request/v1"`, `operation="inspect_environment"`, a fresh `requestId`, canonical `projectRoot` and `runtimeContext`; include `configPath` when selecting a definition. The server verifies target identity and returns observed state, configuration/provenance, digest and observation time.
3. When readiness findings are needed, call `check_environment` separately with the same schema version, `operation="check_environment"`, a fresh `requestId` and canonical `projectRoot`. Named `checks` are optional (`[]` means all eligible); offline remains the default absent separately authorized network access. Report each finding and `complete` exactly as returned; hook and mount findings are review evidence, not approval.
4. If access, a check or a trust decision is missing, report the exact gap and ask the human operator to decide through the host-owned policy process. There is no MCP approval action.

### US-11 — Request only an already granted build or start

> I as an external management agent would like to request a build or start for the exact inspected inputs, so that I can help prepare a workspace without expanding my own authority.

Acceptance outcome: the server acts only under a non-expired grant matching established requester, canonical root/config, operation, runtime context, input digest and any existing resource ID relevant to that operation. Build reports inspected image identity; start reports the observed running primary resource without implying interactive attachment. The versioned result reports outcome, target, digest, data or structured error/next action, evidence and completed or unknown effects when partial.

Solution path:

1. Call `inspect_environment` with the versioned request envelope, fresh `requestId`, canonical `projectRoot` and `runtimeContext` (plus selected `configPath`). Retain the returned target, `inputDigest` and observation time; use `check_environment` separately when needed.
2. Submit the registered typed `build_environment` or `start_environment` tool with `schemaVersion="aibox.operation-request/v1"`, matching `operation`, a fresh `requestId`, canonical `projectRoot`, `expectedInputDigest` and `runtimeContext` (plus selected `configPath`). `frozenLockfile` defaults to true; changing it is an explicit request choice. Do not supply actor, executable or policy fields.
3. The server establishes the principal, matches the grant, locks and rechecks target identity and digest, writes the durable operation record, delegates to the pinned CLI and reinspects native results. Host hooks remain denied without separate host-owned approval against the reviewed digest.
4. Report the complete structured result, including refusal, changed inputs or partial/unknown effects and the exact next action. Build does not replace a running container; start does not attach a terminal. Neither repository content nor an MCP request issues its own grant.

### US-12 — Reconcile a disconnected request

> I as an external management agent would like to inspect the durable record for an interrupted operation, so that I can tell the human what may have changed before proposing another request.

Acceptance outcome: a scoped `inspect_operation` read returns operation ID, original operation kind, durable state, last confirmed step, completed and unknown effects and a closed next action. It neither promises a copy of the original request nor resumes execution. Knowing the ID alone grants no access; a retry is a fresh typed request after current state and inputs are inspected.

Solution path:

1. Retain the operation ID from the operation record or evidence reference and the canonical project root; distinguish this durable ID from the request ID of an individual MCP call.
2. After reconnecting, call `inspect_operation` with `schemaVersion="aibox.operation-request/v1"`, `operation="inspect_operation"`, a fresh `requestId`, `projectRoot` and `operationId`. The server checks principal, root and operation-inspection grant; possession of the ID is insufficient.
3. Report `data.originalOperation`, `data.state`, `data.lastConfirmedStep`, `data.completedEffects`, `data.unknownEffects` and `data.nextAction`, plus the enclosing outcome, target and evidence. Disconnection is not rollback.
4. If effects remain uncertain, call `inspect_environment` as a separate authorized read. Explain any human decision, then submit a new typed mutation request only after fresh target/input checks and a matching unexpired grant. Never resume or blindly replay the interrupted operation.

## Acceptance relationship

The stories are complete only when the executable, native-runtime integration,
security boundary and user documentation satisfy the outcomes above on the
supported target matrix. Chapter 14 defines the evidence required to make
that claim; delivery sequencing does not change the target behavior.

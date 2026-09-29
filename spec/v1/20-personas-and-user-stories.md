# 20. Personas, user stories and solution paths

This chapter turns the product contract into testable end-to-end journeys.
It defines the intended v1 behavior and user-visible outcomes. Command names,
authority boundaries and result contracts follow chapters 3, 9, 10, 12 and 13.
Each journey is a product acceptance scenario, independent of delivery order.

## Personas and ACTION review prompts

Each prompt uses Act as, Context, Task, Iterate Output and Netiquette. Agents
review from a persona's own goals and authority, rather than summarizing
the product team's implementation plan. Reviews return concrete corrections
to the stories and the behavior that would satisfy each need.

### P1 — Independent developer (human)

Works on several repositories from a personal laptop and alternates between
terminal and VS Code. Understands Dev Containers but does not want to learn a
second build configuration language. Wants a fast first start, reproducible
tool choices, familiar editor attachment and persistent personal state. Owns
project configuration, but does not assume that editing it grants permission
for host-executed hooks or unsafe mounts. Uses aibox as the normal lifecycle
entry point; direct upstream CLI use is a compatibility property to verify.

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
> **Task:** Review US-01–US-03. Identify any step that makes you use the
> upstream CLI as the primary workflow, introduces a second config language,
> loses home data, or prevents standard VS Code attachment. Rewrite the
> affected need, solution step and observable outcome using `aibox inspect`,
> `aibox build`, `aibox up` and `aibox attach` where appropriate.
>
> **Iterate Output:** Return a short verdict for each story, followed by
> replacement Markdown for each inaccurate part. Distinguish your routine
> aibox path from a separate direct-upstream compatibility test.
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

> **Act as:** Be this human platform operator approving and performing a
> bounded operation, not a project developer editing untrusted inputs.
>
> **Context:** Multiple repositories can declare lifecycle hooks, mounts and
> Features. Policy is installed outside them. A build can run host-side
> preparation; a rebuild can disrupt a live container; a timeout can leave
> effects behind. The native runtime may contain same-name resources in other
> contexts.
>
> **Task:** Review US-04–US-06 for exact target and input identity, hook/mount
> review, separate build/up/rebuild authority, retained data and recovery.
> Replace any step that treats repository intent as approval, trusts a name
> alone, or assumes cancellation reverses native effects.
>
> **Iterate Output:** For each story, state the operator decision point, the
> CLI request, aibox's ordered checks/delegation/reinspection, and evidence
> required to accept the result. Return replacement Markdown for gaps.
>
> **Netiquette:** Be precise about risk and uncertainty; request a new human
> decision instead of silently broadening scope or globally pruning.

### P3 — Repository/workspace agent (AI, inside the container)

Works inside the Dev Container with project and user-owned files, but no host
runtime socket, operator credentials or endpoint. Helps its human understand
configuration and diagnostics, edit a supported local UX setting, and apply
bounded local refresh. Must distinguish observed local facts from unknown host
state. A project edit is not a host-operation grant.

> **Act as:** Be this in-container AI coding agent, speaking as the actor who
> must explain what it has actually observed or changed.
>
> **Context:** The user asks whether a setting is active, then asks for a
> theme or tmux adjustment. You have local `inspect`, `doctor`, guidance and
> `refresh_workspace` access, but cannot invoke host build/up/rebuild or use
> an undisclosed bridge. Native image or Feature changes require operator
> handoff.
>
> **Task:** Review US-07–US-09. Make every step distinguish inspection from
> mutation, project defaults from user overrides, and local refresh from
> host lifecycle. Correct any claim of host observation or automatic approval.
>
> **Iterate Output:** Return a verdict per story and replacement first-person
> need, CLI/MCP sequence, result interpretation and handoff where necessary.
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

> **Act as:** Be this external AI management agent with a narrow, existing
> operator grant; do not impersonate the human owner.
>
> **Context:** You may inspect a permitted project/context, receive its input
> digest, and request a typed build or start. A hook finding, expired grant,
> changed digest or interrupted operation can require a new human decision.
> A disconnected client is not evidence of rollback.
>
> **Task:** Review US-10–US-12 for typed request inputs, principal and target
> checks, operation-record inspection, partial results and human escalation.
> Reject an invented approval tool, host shell fallback or self-issued grant.
>
> **Iterate Output:** For each story return a verdict and replacement need,
> MCP sequence, result fields to inspect and exact point of human handoff.
>
> **Netiquette:** Report refusals and unknown effects without euphemism; do
> not conceal a boundary because completing the task would be convenient.

## Independent developer stories

### US-01 — Try a native Dev Container without another config language

> I as an independent developer would like to start from a clearly commented `devcontainer.json` and change it using the Dev Container format I already know, so that I can try aibox without translating my project into a second configuration language.

Acceptance outcome: the starter is valid JSONC, explains available aibox options and leaves unused options commented. aibox identifies the selected native configuration and builds and starts it through the pinned Dev Container CLI without changing its meaning.

Solution path:

1. Copy the `minimal` Template files into my project.
2. Edit `.devcontainer/devcontainer.json`, enabling only the native properties and aibox options I need. Feature-specific settings belong to their selected Features.
3. Run `aibox inspect --project PATH --format json` to review the selected definition, declared inputs and input digest.
4. Run `aibox build --project PATH`, then `aibox up --project PATH`. aibox checks the exact target and authorized inputs, delegates native build/start behavior to the pinned Dev Container CLI, then reports image and container identity.
5. Attach through `aibox attach --project PATH` or VS Code's standard attach-to-running-container flow. Direct upstream use is a separate interoperability test, not the normal aibox journey.

### US-02 — Add only the tools I choose

> I as an independent developer would like to select tools through documented Dev Container Features and their own options, so that I know what installs each tool and can keep the environment focused.

Acceptance outcome: selected tools have explicit Feature references; commented examples are not installed, and the native project builds and starts through aibox without a separate Feature installer.

Solution path:

1. Review the commented starter catalog and edit the native `features` object using each chosen Feature's documentation.
2. Enable any desired `customizations.aibox` workspace presentation settings separately; these configure UX, not tool installation.
3. Run `aibox inspect --project PATH` to review selected Feature references and inputs, then `aibox build --project PATH` and `aibox up --project PATH`.
4. aibox validates target and authority, invokes the pinned upstream CLI and reports actual results; upstream resolves and installs Features. The user does not run it as a separate normal step.

### US-03 — Attach and retain my personal state

> I as an independent developer would like to enter the running workspace from VS Code or a terminal and retain my home across rebuilds, so that aibox fits my existing Dev Container workflow.

Acceptance outcome: the native starter declares a project-scoped persistent home, the aibox-started container can be entered with standard Dev Container tooling, and the named home volume survives a rebuild.

Solution path:

1. Check the starter's `mounts`, `remoteUser` and image user; the home mount target must match that user's home.
2. Run `aibox build --project PATH` and `aibox up --project PATH` to prepare and start the workspace.
3. Attach with `aibox attach --project PATH` from a terminal or VS Code's standard attach-to-running-container flow.
4. Run `aibox rebuild --project PATH` with the required disruption acknowledgment, then confirm the named home volume's identity and retained files. Direct upstream compatibility is tested separately.

## Human platform operator stories

### US-04 — Review before authorizing a host build

> I as a platform operator would like to review the exact project inputs, host hooks, mounts and runtime context before authorizing a build, so that repository content cannot silently acquire host access.

Acceptance outcome: aibox performs a build only after operator policy permits the canonical root, runtime context, build operation, executable provenance and reviewed input digest. Changed inputs, ambiguous identity or unapproved hook/mount risk stop it before effects.

Solution path:

1. Configure operator-owned policy outside the repository: allowed roots, contexts, operations, executable provenance and hook/mount approval. Repository declarations are inputs, not grants.
2. Run `aibox inspect --project PATH` and read-only `aibox doctor --project PATH` in the selected operator context. Review resolved configuration, hook commands, mount sources/targets, runtime context and input digest without executing hooks.
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

Acceptance outcome: aibox rebuilds only the confirmed canonical project, runtime context and native resource ID after explicit disruption consent. The selected persistent volume identity and mount target are retained by default and verified afterward; same-name resources elsewhere remain untouched.

Solution path:

1. Run `aibox inspect --project PATH` in the selected operator context. Confirm native resource ID, project association, input digest, persistent volume ID and mount target; a name alone is insufficient.
2. Review proposed image/config changes and interruption, then give explicit consent for this rebuild only. Changed inputs or target identity require renewed review.
3. Run `aibox rebuild --project PATH` with the required exact ID, digest and disruption acknowledgment. aibox rechecks external policy, takes a scoped lock, delegates to the pinned Dev Container CLI and reinspects resource and volume identity.
4. Accept only a result that reports the target, completion or partial outcome, retained volume identity and any unknown effects. Removal remains separately authorized; rebuild never permits global pruning.

## AI workspace agent stories

### US-07 — Explain the selected workspace from inside it

> I as the workspace agent would like to inspect the selected Dev Container definition and declared project inputs, so that I can explain what the project requests while keeping local observations separate from host runtime state.

Acceptance outcome: I report the selected definition, declared inputs, input digest and eligible local findings with provenance. Host runtime state and unavailable or unauthorized checks remain explicitly unknown.

Solution path:

1. Call local `inspect_workspace` or run `aibox inspect --context local --project PATH --format json`; record selected definition, declared inputs, digest, source and observation time.
2. Call `check_workspace` or `aibox doctor --context local --project PATH` for eligible read-only local checks. Report each finding's status; unavailable or unauthorized is not healthy. Never execute lifecycle hooks or query host runtime from local context.
3. Distinguish active values from commented examples. Report effective-value provenance only when returned, and describe declared intent as intent rather than proof that the host applied it.

### US-08 — Refresh a local presentation choice

> I as the workspace agent would like to adjust a supported user-owned theme or tmux setting and refresh only local presentation, so that I can help the user iterate without rebuilding the container.

Acceptance outcome: the user can identify the edited project default or user-owned override, changed managed files and any required rebuild or session restart. Unrelated files and persistent home remain intact; no host lifecycle effect is implied.

Solution path:

1. Read task-sized guidance for the requested setting: exact path/key, ownership, validation, activation and rollback. Ask the user to choose between project default and user override when ambiguous.
2. Edit only authorized local UX configuration; do not change image, Feature, mount or lifecycle-hook inputs as part of a presentation refresh.
3. Call local `refresh_workspace` with canonical project root, expected UX input digest and `theme`, `tmux` or `all` scope, or run the equivalent `aibox refresh`. It validates and stages managed local output, then reports `changedFiles`, `rebuildRequired` and `sessionRestartRequired`.
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

Acceptance outcome: the result gives project/configuration identity, runtime context, input digest, observation time and read-only findings within my grant. Observed, skipped, unavailable and unauthorized checks are distinct; findings do not grant mutation authority.

Solution path:

1. Connect through the separately launched operator stdio MCP server. Use the requester principal established by its operator configuration; client-supplied actor metadata is not identity proof.
2. Call `inspect_environment` with permitted project root, selected config and runtime context. The server verifies target identity and returns observed state plus the digest for a later mutation.
3. Call `check_environment` for declared checks. Report each finding and status; hook and mount findings are review evidence, not approval.
4. If access, a check or a trust decision is missing, report the exact gap and ask the human operator to decide through the host-owned policy process. There is no MCP approval action.

### US-11 — Request only an already granted build or start

> I as an external management agent would like to request a build or start for the exact inspected inputs, so that I can help prepare a workspace without expanding my own authority.

Acceptance outcome: the server acts only under a non-expired grant matching established requester, canonical root/config, operation, runtime context, input digest and applicable native resource identity. Build reports image identity; start reports running resource identity without implying interactive attachment.

Solution path:

1. Call `inspect_environment` for the canonical root and runtime context; retain target identity and current input digest.
2. Submit typed `build_environment` or `start_environment` with project/config selectors, runtime context, `expectedInputDigest` and a new request ID. Use the frozen-lockfile default unless the human chooses otherwise.
3. The server establishes the principal, checks the grant, locks and rechecks target identity and digest, records the operation, delegates to the pinned upstream CLI and reinspects native results.
4. Report the structured outcome, error/next action, evidence and observed resources. On refusal, changed inputs or host-code risk, stop for a human policy decision; neither repository content nor an MCP request can issue its own grant.

### US-12 — Reconcile a disconnected request

> I as an external management agent would like to inspect the durable record for an interrupted operation, so that I can tell the human what may have changed before proposing another request.

Acceptance outcome: authorized inspection returns the original request, durable state, last confirmed step, known and unknown effects and bounded next action without resuming execution. Knowing the ID alone grants no access; a retry is a fresh typed request.

Solution path:

1. Retain the operation ID and canonical project root from the typed lifecycle result.
2. After reconnecting, call `inspect_operation` with root and ID under the same established principal and authorized scope; possession of the ID is not authority.
3. Report original operation, last confirmed step, known and unknown effects and next action. If partial or uncertain, call `inspect_environment` to reconcile native state. Disconnection is not rollback.
4. Explain any required human decision, then submit a new typed request with a new ID only after fresh target/input checks and a matching unexpired grant. Never resume or blindly replay the interrupted operation.

## Acceptance relationship

The stories are complete only when the executable, native-runtime integration,
security boundary and user documentation satisfy the outcomes above on the
supported target matrix. Chapter 14 defines the evidence required to make
that claim; delivery sequencing does not change the target behavior.

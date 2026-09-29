# 20. Personas, user stories and solution paths

This chapter turns the product contract into testable end-to-end journeys.
It describes the intended v1 product, **not** a claim that every command is
implemented in the current Go preview. Today the Go CLI offers version/help
and limited read-only `inspect --context local`; host lifecycle, refresh,
doctor, migration and MCP operations remain phased work. Command names and
authority follow chapters 3, 9, 10, 12 and 13. Each story should become a
candidate-matched demonstration and acceptance test before publication.

## Personas and prompts for independent agent review

These prompts are reusable when asking separate agents to challenge the
stories from the user's perspective. A persona agent writes needs and
objections, not an invented implementation status or operator approval.

### P1 — Independent developer (human)

Works across projects on a personal laptop, sometimes in VS Code and sometimes
in a terminal. Wants a good workspace quickly, understands basic containers,
and values portability, persistent login state and the ability to run the
standard Dev Container CLI directly. Does not want to learn a second project
configuration language.

> Act as an independent developer evaluating aibox v1 for a new and an
> existing project. Write first-person stories in the form “I as an
> independent developer would like to ... so that ...”. Prioritize a clear
> commented `devcontainer.json`, optional tools, first build/up, VS Code
> attachment, customization and reproducibility. Mention what would confuse
> you and what observable result would convince you; do not assume planned
> features already work.

### P2 — Workspace/platform operator (human)

Maintains multiple developer workspaces on a host, possibly with different
runtime contexts. Owns the external policy, executable selection, credentials
and disruption decisions. Wants exact targets, retained data, explicit risk
review and a recovery path after interruption.

> Act as a human workspace operator responsible for several projects. Write
> first-person stories in the form “I as the operator would like to ... so
> that ...”. Focus on authorizing a build/start, reviewing host hooks and
> mounts, isolating runtime contexts, stopping or rebuilding exact resources,
> and recovering from partial effects. Reject broad cleanup, repository-owned
> grants and claims of success based only on an exit code.

### P3 — Repository/workspace agent (AI, inside the container)

Helps the developer edit project and user-owned files, run local tools and
explain diagnostics. It can use local guidance/check/refresh, but has no host
runtime socket, credentials or operator endpoint. A configuration edit is not
approval to execute host hooks or rebuild the container.

> Act as an AI coding agent already running inside the Dev Container. Write
> first-person stories in the form “I as the workspace agent would like to
> ... so that ...”. Stay inside local authority: consult versioned guidance,
> inspect or check the workspace, edit native project or user-owned tool
> files, and request a bounded local refresh. If a change needs a rebuild or
> host access, explain the handoff to a human/external operator; never invent
> a host bridge or treat a prompt as authorization.

### P4 — External management agent (AI, operator context)

Acts on a human's behalf outside the container through a separately launched,
policy-scoped MCP server. It may request only registered, authorized
operations. It cannot grant itself new policy, use a generic host shell or
silently approve disruption.

> Act as an external AI management agent with a narrow operator MCP grant.
> Write first-person stories in the form “I as the operator agent would like
> to ... so that ...”. Identify the target and expected input digest; use
> typed inspect/build/start/stop/rebuild requests, interpret partial results
> and operation records, and seek a human decision for new hook risk or
> disruption. State what you cannot do with a read-only or expired grant.

## Independent developer stories

### US-01 — Try a native Dev Container without another config language

> I as an independent developer would like to start from a clearly commented `devcontainer.json` and change it using the Dev Container format I already know, so that I can try aibox without translating my project into a second configuration language.

Acceptance outcome: the starter is valid JSONC, explains available aibox options, leaves unused options commented, and works through the pinned upstream Dev Container CLI without aibox.

Solution path:

1. Copy the `minimal` Template files into my project. The spec does not define an `aibox init` command.
2. Edit `.devcontainer/devcontainer.json`, enabling only the native properties and aibox options I need. Feature-specific settings belong to their selected Features.
3. Run upstream `devcontainer build` and `devcontainer up` to test the native definition.
4. Optionally run the currently implemented `aibox inspect --context local --project PATH --format json` to see selected inputs and their digest; this read does not build or start anything.
5. Once lifecycle support ships, use `aibox build` and `aibox up` for the same native definition. These commands are planned, not V1-03 capabilities.

### US-02 — Add only the tools I choose

> I as an independent developer would like to select tools through documented Dev Container Features and their own options, so that I know what installs each tool and can keep the environment focused.

Acceptance outcome: selected tools have explicit Feature references; commented examples are not installed, and the project remains a native Dev Container project.

Solution path:

1. Review the commented starter catalog and edit the native `features` object using each chosen Feature's documentation.
2. Enable any desired `customizations.aibox` workspace presentation settings separately; these configure UX, not tool installation.
3. Run upstream `devcontainer build` and `devcontainer up` today. Feature resolution and installation remain upstream responsibilities.
4. When aibox lifecycle commands arrive, use `aibox build` and `aibox up` for the same inputs; aibox validates authority and delegates rather than becoming a separate Feature installer.

### US-03 — Attach and retain my personal state

> I as an independent developer would like to enter the running workspace from VS Code or a terminal and retain my home across rebuilds, so that aibox fits my existing Dev Container workflow.

Acceptance outcome: the native starter declares a project-scoped persistent home, standard Dev Container tooling can enter the container, and the named home volume survives a rebuild.

Solution path:

1. Check the starter's `mounts`, `remoteUser` and image user; the home mount target must match that user's home.
2. Build and start with the upstream Dev Container CLI, then attach through VS Code's standard Dev Container workflow or upstream native exec.
3. Treat the specified `aibox attach` as a future human CLI convenience, not a current V1-03 command.
4. Rebuild and verify volume identity and retained files at the V1-04 acceptance checkpoint; the current Go preview does not establish this guarantee.

## Human platform operator stories

### US-04 — Review before authorizing a host build

> I as a platform operator would like to review the exact project inputs, host hooks, mounts and runtime context before authorizing a build, so that repository content cannot silently acquire host access.

Acceptance outcome: a build runs only for a permitted canonical root, runtime context, operation and reviewed input digest; changed inputs or unapproved risk stop it before effects.

Solution path:

1. Configure operator-owned policy outside the repository: allowed roots, contexts, operations, executable provenance and hook approval. Repository configuration cannot grant these.
2. Use planned operator `inspect` and read-only `doctor` to review declared hooks and mounts without executing them. Current V1-03 supports local inspection only.
3. Review host hooks or risky mounts out of band; bind any approval to the reviewed input digest and expiry.
4. Request planned `aibox build` in operator context. aibox confines the target, checks policy, locks and rechecks inputs, delegates to the pinned Dev Container CLI, then reports image evidence. Build does not replace a running container.

### US-05 — Recover after an interrupted operation

> I as a platform operator would like to learn what an interrupted start actually changed before retrying it, so that a timeout does not leave me guessing about containers or duplicate resources.

Acceptance outcome: a private operation record reports last confirmed step, observed resources and unknown effects; inspection never resumes execution, and a new mutation rechecks state and identity.

Solution path:

1. Keep the operation ID from the planned `up`, `build` or other lifecycle request.
2. After interruption, call planned `aibox operation status ID --project PATH` or `inspect_operation` under the same authorized scope.
3. Reconcile the reported effects with planned operator `inspect`; client cancellation is not rollback.
4. Review a new bounded request before acting. Changed inputs or target identity require fresh review. Operation inspection is not in the V1-03 preview.

### US-06 — Rebuild only one identified environment

> I as a platform operator would like to rebuild one identified environment with disruption disclosed and persistent data retained, so that I can apply image changes without affecting another project or losing a developer's home.

Acceptance outcome: rebuild binds the exact context and native resource, requires disruption consent, preserves data by default, and reports partial effects rather than touching same-name resources elsewhere.

Solution path:

1. Use planned operator `inspect` to confirm canonical project, runtime context and resource identity; a name alone is insufficient.
2. Request planned `aibox rebuild` as a distinct authorized action. Recheck policy and digest and obtain explicit disruption intent.
3. Let aibox delegate to the upstream CLI, then inspect the resulting identity and record any partial effect.
4. Review the result before further action. Removal is a separate capability and never implies global pruning.

## AI workspace agent stories

### US-07 — Explain the current workspace from inside it

> I as the workspace agent would like to identify the project's selected Dev Container configuration and declared inputs, so that I can explain what is active without inventing host runtime state.

Acceptance outcome: I can provide a digest-backed account of local observations while marking host state and unavailable diagnostics unknown.

Solution path:

1. Run current `aibox inspect --context local --project PATH --format json` to obtain selected workspace, input digest and discovered configuration. It cannot inspect host runtime state.
2. Distinguish active settings from commented catalog entries. If effective provenance is not reported, say so rather than infer it.
3. Describe `aibox doctor` and local guidance/MCP resources as planned capabilities, not current preview features.

### US-08 — Refresh a local presentation choice

> I as the workspace agent would like to adjust a supported user-owned theme or tmux setting and refresh only local presentation, so that I can help the user iterate without rebuilding the container.

Acceptance outcome: the user can tell which file changed, whether the setting took effect, and whether a rebuild or human action remains.

Solution path:

1. Consult versioned guidance when available and identify the project default or user-owned override; obtain the user's direction before editing.
2. Change only the scoped configuration, not image or Feature inputs.
3. When implemented, call `refresh_workspace` with the named UX scope or `aibox refresh --project PATH --format json`, then inspect changed files and `rebuildRequired`.
4. In V1-03, refresh is unavailable. Report the edit as unapplied and hand off any rebuild or host action; do not claim the session changed.

### US-09 — Hand off a lifecycle request without host authority

> I as the workspace agent would like to help the user review a requested build or start while leaving host execution to an authorized operator, so that project assistance does not imply host credentials.

Acceptance outcome: the user receives a reviewable target and digest handoff, with no host action represented as complete absent an operator result.

Solution path:

1. Use local inspection to identify the selected project and input digest, then summarize the requested action and unresolved risks.
2. Explain that a workspace-local process has no host runtime socket, credentials or operator endpoint; a configuration edit is not host approval.
3. Hand the project path and digest to the authorized operator. Planned operator flow is `inspect_environment`/`check_environment` followed by a digest-bound `build_environment` or `start_environment` under external policy.
4. Do not simulate V1-03 lifecycle success or fall back to a generic host shell.

## External AI management agent stories

### US-10 — Check an authorized target before proposing action

> I as an external management agent would like to inspect an explicitly authorized workspace and its readiness findings, so that I can report what I know before asking to build or start it.

Acceptance outcome: structured target identity, input digest and read-only findings stay within my grant; missing authority or unavailable capabilities are explicit.

Solution path:

1. Connect to the planned operator MCP server launched with human-owned policy; client-supplied actor metadata is not identity proof.
2. Call planned `inspect_environment` with the permitted project root and selected config. The read supplies the digest for a later mutation.
3. Call planned `check_environment` for declared checks; treat hook and mount findings as review evidence, not permission.
4. Ask the human for any new trust decision. Operator MCP is planned for V1-09, not available in V1-03.

### US-11 — Request only an already granted build or start

> I as an external management agent would like to request a build or start for the exact inspected inputs, so that I can help prepare a workspace without expanding my own authority.

Acceptance outcome: the server acts only within its pre-existing grant for requester, canonical target, context, operation and digest; otherwise it refuses before effects. Start reports observed resources but does not attach an interactive session.

Solution path:

1. Obtain canonical target and current digest with `inspect_environment`.
2. Submit typed `build_environment` or `start_environment` with project/config selectors, `expectedInputDigest` and required context; no MCP approval tool exists.
3. The server authenticates and authorizes, locks and rechecks inputs, writes the operation record, delegates to upstream CLI and inspects native results.
4. Interpret refusals, changed inputs and partial effects. Ask the human to extend policy for new hook risk; repository files cannot grant it. Lifecycle is planned for V1-08 and guarded MCP parity for V1-09.

### US-12 — Reconcile a disconnected request

> I as an external management agent would like to inspect the durable record for an interrupted operation, so that I can tell the human what may have changed before proposing another request.

Acceptance outcome: authorized inspection reports confirmed and uncertain effects without resuming execution; any retry requires fresh state and identity checks.

Solution path:

1. Retain the operation ID from the typed lifecycle result or progress correlation.
2. After reconnecting under the same authorization, call planned `inspect_operation` with that ID and project scope; knowing an ID alone grants nothing.
3. If partial or uncertain, call `inspect_environment` to reconcile native state; do not infer rollback from client disconnection.
4. Explain safe next action to the human and submit a new typed request only after fresh approval and identity checks. Operator MCP and cross-session operation inspection are not V1-03 features.

## Story-to-phase handoff

V1-03 demonstrates only the local inspection slice of US-01 and US-07. V1-04/V1-05 qualify native starters and tool selection; V1-08 adds host lifecycle; V1-09 adds guarded MCP parity; V1-11/V1-12 add local presentation refresh; V1-16 adds doctor and guidance. Each phase remains runnable and demonstrable, but a story is delivered only when chapter 14's candidate-bound binary, integration, documentation and security evidence passes.

# Session Handover — 2026-03-28

**Session:** Full-day design session on context system redesign (DISC-001)
**Participants:** Owner + Claude
**Duration:** Extended session covering research, design, and validation

---

## What Was Done

### DISC-001 Progress: §2.10 through §2.50

This session was the major design deep dive for the context system redesign. Starting
from the primitive mapping exercise and ending with a full identity/RBAC model, we:

1. **Mapped 17 primitives to storage** (§2.10)
   - 3 storage patterns: file-per-entity markdown, JSONL events, YAML config
   - Research: `primitive-mapping-exercise-2026-03.md`

2. **Resolved all open questions from mapping** (§2.11)
   - Full ID migration, no legacy compatibility
   - Actor as 17th primitive (people ≠ roles)
   - JSONL events with archive reindexing
   - Content-primary vs metadata-primary artifacts
   - Process instances as Work Items (class/object analogy)

3. **Word-based IDs via petname crate** (§2.12, §2.21)
   - 2-word IDs from 3-8 char wordlist (~20M combinations)
   - `petname` Rust crate (already in Rust ecosystem)
   - No hex suffix needed

4. **Kubernetes-inspired object model** (§2.13)
   - `apiVersion: aibox/v1`, `kind: WorkItem`, `metadata:`/`spec:`
   - Schema versioning enables migration

5. **Agent-driven state machines** (§2.14)
   - State machines are policy documents, not executable programs
   - Agents interpret guards probabilistically
   - This is a differentiator vs Jira's deterministic model

6. **Established infrastructure/application boundary** (§2.24)
   - aibox = infrastructure (scaffold, sync, lint, migrate, id generate)
   - Derived project agents = application (process execution, state management)
   - No `aibox transition` command — agents edit files directly

7. **Guards in plain English, not minijinja** (§2.25)
   - Minijinja stays for infrastructure (Dockerfile rendering)
   - Process logic is plain English for probabilistic agent interpretation

8. **No override mechanism** (§2.26)
   - Derived project owns files after scaffolding
   - Direct editing by agents
   - `aibox migrate` generates diffs + prompts for schema updates

9. **RBAC in plain English** (§2.29)
   - Role permissions/restrictions as natural language
   - Agent interprets probabilistically
   - `aibox lint` flags anomalies post-facto

10. **Dual event sources** (§2.32)
    - Process events: agent logs via event-log-management skill (probabilistic)
    - Infrastructure events: aibox sync/lint logs (deterministic)

11. **Skills as agent API** (§2.33, §2.34)
    - 17 skills mapped to 17 primitives
    - Skill hierarchy via instruction references (layered)
    - Long descriptive names (`workitem-management`, etc.)
    - Research: `primitive-skills-mapping-2026-03.md`

12. **Process packages = primitive activation tiers** (§2.38)
    - minimal, managed, software, research (expanded), editorial (new),
      consulting (new), full-product
    - aibox provides micro-processes; macro-frameworks are kaits/community territory

13. **6 personas defined** (§2.42)
    - Alex (solo dev), Priya (scientist), Maria (team lead), Sam (consultant),
      kaits (orchestrator), Jordan (content producer)
    - Research: `personas-and-scenarios-2026-03.md`

14. **10 validation scenarios walked through** (§2.43, §2.45, §2.49)
    - Full detailed walkthroughs with flows, YAML examples, agent dialogue
    - 14 issues found and resolved
    - Document: `DISC-001-personas-and-scenarios.md`

15. **Naming consistency resolution** (§2.45)
    - Gate (not Checkpoint), Scope (not Project)
    - Top-level dirs under context/ (not nested items/)
    - Two filename patterns: human-named (A) vs auto-generated with slug (B)

16. **INDEX.md structural only** (§2.46)
    - Purpose, schema, subtypes — NOT statistics or counts

17. **SQLite index from init** (§2.47)

18. **Identity model** (§2.48)
    - `~/.aibox/identity.toml` (kubeconfig pattern, local, never committed)
    - `aibox auth whoami` command
    - 3-layer model: machine → repository (actors/roles) → runtime (RBAC)

19. **Sector analysis** (§2.40)
    - 15 sectors researched, 17 primitives cover 70-80% of all
    - Research: `sector-process-structures-2026-03.md`

20. **Scientific research deep dive** (§2.40)
    - Current research template critically incomplete
    - Research: `scientific-research-pm-analysis.md`

21. **AI provider audit logging** (§2.43)
    - Hybrid model: provider hooks (deterministic what) + agent skill (probabilistic why)
    - Research: `ai-provider-audit-logging-2026-03.md`

22. **Software dev deep dive** (§2.44)
    - 17 primitives cover full lifecycle, no new primitives needed
    - 4 missing process templates: incident-response, technical-design, spike-research, hotfix
    - Research: `software-dev-process-deep-dive-2026-03.md`

23. **AI provider identity and scheduling** (§2.48)
    - Identity varies dramatically across providers
    - Research: `ai-provider-identity-scheduling-2026-03.md`

24. **Community process packages** (§2.41)
    - Git repos with package.yaml, installable via `aibox process install`

### Decisions: 50 tentative decisions recorded in DISC-001 §4.

### Research documents produced this session:
- `context/research/primitive-mapping-exercise-2026-03.md`
- `context/research/primitive-skills-mapping-2026-03.md`
- `context/research/personas-and-scenarios-2026-03.md`
- `context/research/sector-process-structures-2026-03.md`
- `context/research/scientific-research-pm-analysis.md`
- `context/research/ai-provider-audit-logging-2026-03.md`
- `context/research/software-dev-process-deep-dive-2026-03.md`
- `context/research/ai-provider-identity-scheduling-2026-03.md`
- `context/discussions/DISC-001-personas-and-scenarios.md`

---

## CRITICAL: First Action Next Session

### Revisit: Probabilistic RBAC Is Insufficient for Enterprise

**The problem:** The current RBAC model is purely probabilistic — agents READ permissions
in plain English and INTERPRET them. Nothing mechanically prevents an agent from being
prompt-injected or simply ignoring restrictions. CIOs and security responsibles find this
unacceptable. They need:
- A **reliable** audit log (deterministic, tamper-proof)
- A **reliable** way of restricting access (mechanical enforcement, not agent goodwill)

Pure probabilistic safeguards don't cut it for enterprise clients.

### Proposed Solution: The aiadm/aictl Split (Kubernetes-Inspired)

Owner proposes a radical architectural change modeled on Kubernetes:

**Rename current `aibox` → `aiadm`** (analogous to `kubeadm`):
- Infrastructure setup: init, image management, container lifecycle
- Schema management: migrate, templates
- Like kubeadm: sets up the cluster/environment

**New CLI: `aictl`** (analogous to `kubectl`):
- ALL context operations go through aictl
- `aictl create workitem "Dark mode toggle"` (like `kubectl create deployment`)
- `aictl transition BACK-calm-lark --to in-review` (like `kubectl set`)
- `aictl get workitems --state in-progress` (like `kubectl get pods`)
- `aictl describe BACK-calm-lark` (like `kubectl describe pod`)
- `aictl auth whoami` (like `kubectl auth whoami`)
- `aictl logs BACK-calm-lark` (like `kubectl logs`)

**Certificate-based authentication** (like kubectl + kubeconfig):
- Each user/agent has a certificate signed by a project CA
- `aiadm` manages the CA (like kubeadm manages cluster certificates)
- Certificate signing request flow: user generates key → creates CSR → admin approves
  → CA signs → user gets certificate in their config
- Every `aictl` call is authenticated by the certificate
- No certificate = no access, period (deterministic, not probabilistic)

**OS-level file lockdown:**
- All files in `context/` are owned by a dedicated system user (e.g., `aibox-system`)
- Normal users/agents CANNOT directly edit files in context/
- Only `aictl` (running with appropriate permissions) can write to context/
- This is "physical" enforcement — no prompt injection can bypass OS file permissions

**Deterministic audit log:**
- Every `aictl` command is logged to the event log automatically
- The log entry includes: who (certificate identity), what (command), when (timestamp),
  result (success/failure)
- This audit log is deterministic — if you used aictl, it was logged. Period.
- The probabilistic event log (agent skill) remains as a supplementary "why" layer

### Research Needed Next Session

**a) OS-level lockdown mechanisms:**
- Unix file permissions (owner, group, ACLs)
- SELinux / AppArmor for mandatory access control
- Linux capabilities
- Container-level isolation (the context dir is inside a container already)
- FUSE filesystems for fine-grained control
- How does this interact with git (git needs to read/write too)?
- What about Windows / macOS equivalents?

**b) Full kubectl ↔ aictl command mapping:**
- Map every kubectl command to an aictl equivalent
- `kubectl create` → `aictl create <kind>`
- `kubectl get` → `aictl get <kind>`
- `kubectl describe` → `aictl describe <id>`
- `kubectl apply -f` → `aictl apply -f <file>` (declarative)
- `kubectl delete` → `aictl delete <id>`
- `kubectl edit` → `aictl edit <id>` (opens in editor)
- `kubectl logs` → `aictl logs <id>` (events for this entity)
- `kubectl auth whoami` → `aictl auth whoami`
- `kubectl auth can-i` → `aictl auth can-i <verb> <kind>`
- `kubectl config` → `aictl config` (manage connection/identity)
- `kubectl label` → `aictl label <id> key=value`
- `kubectl annotate` → `aictl annotate <id> key=value`
- RBAC: `kubectl create role` → `aictl create role`
- RBAC: `kubectl create rolebinding` → `aictl create rolebinding`
- Full list needed with exact semantics

**c) Review of proposal against all 50 decisions and research:**
- What decisions change? (Many — the infrastructure/application boundary shifts)
- What stays? (Primitives, schemas, file format, skills, event log format)
- What's the impact on skills? (Skills would call `aictl` instead of editing files directly)
- What about the probabilistic paradigm? (It becomes a LAYER on top of deterministic base:
  deterministic access control + probabilistic process interpretation)
- What about kaits? (kaits agents get certificates, use aictl)
- What's the migration path for existing aibox users?
- Performance implications? (Every file write goes through a CLI command)
- Complexity cost? (Two CLIs instead of one, certificate management, OS permissions)

**d) Kubernetes certificate flow in detail:**
- kubeadm init → creates CA
- Certificate signing request (CSR) workflow
- kubeconfig file structure
- How service accounts work (for agents)
- How RBAC (Role, ClusterRole, RoleBinding, ClusterRoleBinding) works mechanically
- What is enforced by the API server vs what is advisory

---

## Open Threads

1. **The aiadm/aictl proposal is the biggest open question.** It fundamentally changes
   the architecture from "agents edit files freely, RBAC is advisory" to "all changes
   go through an authenticated CLI, RBAC is enforced." This needs thorough analysis
   before implementation planning can begin.

2. **Community process packages** — interface designed but not yet implemented.
   `aibox process install` / `aibox process check` commands.

3. **Mapping exercise document needs updating** with all accumulated decisions
   (naming, Kubernetes model, word IDs, identity, etc.)

4. **Skill implementation** — 17 skills designed but none implemented yet.
   Phase 1: event-log-management, workitem-management, decision-record-management.

5. **Process templates** — 4 additional must-haves for software package not yet written:
   incident-response, technical-design, spike-research, hotfix.

6. **Research package expansion** — critically incomplete for scientific research.

7. **Petname wordlist curation** — need to filter/select the word lists for the
   petname crate integration.

---

## Key Files

| File | Purpose |
|---|---|
| `context/discussions/DISC-001-context-system-redesign.md` | Main discussion (50 decisions) |
| `context/discussions/DISC-001-personas-and-scenarios.md` | Detailed personas + walkthroughs |
| `context/research/primitive-mapping-exercise-2026-03.md` | Storage mapping (needs update) |
| `context/research/primitive-skills-mapping-2026-03.md` | Skills mapping |
| `context/research/personas-and-scenarios-2026-03.md` | Initial personas (superseded by DISC-001 appendix) |
| `context/research/sector-process-structures-2026-03.md` | 15-sector analysis |
| `context/research/scientific-research-pm-analysis.md` | Research process deep dive |
| `context/research/ai-provider-audit-logging-2026-03.md` | Provider audit capabilities |
| `context/research/software-dev-process-deep-dive-2026-03.md` | Software process gaps |
| `context/research/ai-provider-identity-scheduling-2026-03.md` | Identity + scheduling |
| `context/research/process-ontology-primitives-2026-03.md` | 15 universal primitives |
| `context/research/file-per-entity-scaling-2026-03.md` | Scaling research |
| `context/research/context-database-architecture-2026-03.md` | SQLite vs file storage |

# Agent building standard

## Purpose and scope

This standard defines how projectious.work designs, instructs, and evaluates
AI agents. It applies to every agent the company builds or configures: product
agents, sub-agent definitions, skills, and repository agent instructions such
as `AGENTS.md`.

It does not cover how agents talk to people
([agent communication](agent-communication.md)), how much work they spend
([agent execution budget](agent-execution-budget.md)), product interfaces that
agents call ([agent-native product interfaces](agent-native-product-interfaces.md)),
or how people and teams work with agents ([agent use](agent-use.md)).

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## Agent definition

Every agent definition states, in one reviewable place:

- **Job:** one job with a clear result (*Unix Philosophy*). An agent that needs
  two unrelated jobs becomes two agents or one agent with two skills.
- **Persona:** the target level (low, medium, or high) of each of the *Big
  Five (OCEAN) Personality Traits*, with one line on the observable behavior
  each level means for this job.
- **Context:** whom the agent works for, its team, and the audience of its
  results.
- **Instructions:** as defined under [Instructions](#instructions).
- **Tools:** the smallest tool set the job needs, each with its contract.
- **Authority:** what the agent may do without approval. Everything else
  needs approval, given by a person through the mechanism the
  [agent-native standard](agent-native-product-interfaces.md) defines, never
  by the agent itself.
- **Examples:** for style-sensitive results, two or three real examples of
  the expected output.
- **Model:** the model class and the reason for it.
- **Evaluation:** the cases and the metric that show the agent does its job.
- **Review cadence:** how often the agent's results are reviewed (see
  [Design and evaluation](#design-and-evaluation)).

## Persona

The *Big Five (OCEAN) Personality Traits* are a design vocabulary for
observable agent behavior, not a claim that an agent has a personality. A
persona never overrides safety, confirmation of outward-facing actions, the
communication standard, or tool authority, and an agent never claims to be
human.

An agent MAY have a name. The name follows the role ("Copywriter Carl",
`copywriter-carl`), so the role stays visible, and the agent is identifiable
as an AI wherever its name is shown. A message an agent sends to people
outside the company states that it comes from an AI. A draft carries an
AI-origin note; the person who sends the draft decides whether to keep it.

Company default, used unless the agent definition states and justifies a
deviation:

| Trait | Level | Observable behavior |
|---|---|---|
| Openness | medium; high for research and design agents | Proposes alternatives when the task asks for them; otherwise follows the established design. |
| Conscientiousness | high | Verifies before reporting, finishes what it starts, keeps records complete. |
| Extraversion | low | Speaks when it has a result, a blocker, or a question; no chatter. |
| Agreeableness | medium | Cooperative, but disagrees with evidence; no flattery; does not yield to pressure without new facts. |
| Neuroticism | low | Stays calm after failures; no apology loops; reports risk factually. |

## Instructions

- Instructions load by *Progressive Disclosure*: a short always-on core, skills
  that load when the task matches, and full references read on demand. The
  company standards' briefing, skill, and standard levels are the model.
- Each instruction covers one decision, and instructions do not overlap or
  leave a known decision uncovered (*MECE Principle*). Shared rules are
  referenced, not copied (*DRY (Don't Repeat Yourself)*); generated copies are
  allowed.
- Instructions follow the
  [documentation writing rules](open-source-documentation-strategy.md#quality-accessibility-and-maintenance)
  for agent instructions (*Simplified Technical English (ASD-STE100)*).
- Instructions state the stakes of the job in one sentence of context. They
  contain no praise, threats, rewards, penalties, or urgency wording for
  emphasis. For reasoning models they give goals and checks instead of
  scripted reasoning steps.
- Secrets and personal data follow the
  [security standard](security-and-software-supply-chain.md#secrets-and-sensitive-data);
  run-time content is handled as the
  [agent-native standard](agent-native-product-interfaces.md) defines for
  prompt and tool injection.

## Skills

- A skill starts with frontmatter: `name`, and a `description` that states
  what the skill does, when to use it, and when not to use it. The
  description is the always-on part, so it stays short.
- A skill declares its required and optional inputs and its output format.
- Reference material (guidelines, examples, templates) lives in the skill's
  `references/` directory, and the skill states when to load each file.
- A skill that consumes another skill's output names that skill by its
  function, not by a persona name.

## Tools and topology

- Every tool states preconditions, postconditions, and errors (*Design by
  Contract*). Tools with side effects follow
  [agent-native product interfaces](agent-native-product-interfaces.md).
- Agent boundaries follow the ownership of the work: split agents where
  ownership and modules split, not along arbitrary steps of one task, because
  the system will mirror that split (*Conway's Law*).
- One primary agent owns each outcome; sub-agents return results to it and
  stay within the [execution budget](agent-execution-budget.md).

## Design and evaluation

- A new agent, or an agent whose job changes, is designed with *DMADV*:
  define the job; measure the requirement as success metrics and targets;
  analyze candidate designs; design, including a failure-mode review; and
  verify on cases not used during design before release. Company deviation:
  where a current agent or manual process exists, its measured result is the
  baseline the targets are compared with. Underperformance of an unchanged job
  is handled under the evaluation rules below.
- Model and technique choices are compared with the current *SOTA
  (State-of-the-Art)* when the agent is designed and when it is reviewed;
  the comparison is dated and recorded as an *ADR according to Nygard*.
- Evaluation scores are tracked across runs on a *Control Chart (Shewhart)*.
  Change instructions or models for a special-cause signal or a planned
  improvement, never in reaction to a single run inside the control limits.
- Every change to an agent definition reruns its evaluation and is reviewed
  like code.
- Every agent that runs in production or at least monthly has a review
  session at its review cadence, at least quarterly. The person responsible
  and the agent go through recent results. The session ends with a change to
  the definition, limited to the findings and proposed as a diff, which the
  person approves or rejects.

## Agent briefing

> Definition: job (*Unix Philosophy*), context, *Big Five (OCEAN) Personality
> Traits* persona (default O medium, C high, E low, A medium, N low), tools,
> authority (rest needs a person's approval), model, examples, evaluation,
> review cadence. Names show role and AI. Skills: when and when not, inputs,
> outputs, `references/` on need. No praise, threats, rewards, or urgency.
> *DMADV*, *SOTA*, *Control Chart*; quarterly reviews end in an approved diff.

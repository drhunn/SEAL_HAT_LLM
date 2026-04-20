# OPENCLAW SLOT CORRELATION

## purpose

Make explicit the practical correspondence between the specialist slot files used in `SEAL_HAT_LLM` and the OpenClaw-style agent anatomy they resemble.

This document does **not** claim that `SEAL_HAT_LLM` is merely OpenClaw renamed.
It states that the slot structure in this repository has a strong conceptual correlation to OpenClaw-style identity, temperament, behavioral, continuity, and failure-learning components.

The correlation is structural and intentional by design interpretation.
Within this repository, the slot set is defined directly as part of the `SEAL_HAT_LLM` / MM-ELLS architecture, and the project root `AGENTS.md` divides the slots into constitutional versus operational/summary classes.

---

## framing

In this repository, the slot system is part of a governed specialist architecture.

The parent remains constitutionally stable.
Specialists are bounded by slot-defined identity, policy, behavior, and continuity.
The durable memory plane is Postgres + pgvector, while `MEMORY.md` is only a compact projection rather than the primary knowledge base.

That means the OpenClaw correlation should be understood as:

- a **shape and role correlation**
- not a claim that the runtime authority model is identical
- not a claim that markdown alone governs the system

Runtime policy grants actual authority.
Markdown guidance shapes behavior, but does not replace runtime enforcement.

---

## slot correlation

### `IDENTITY.md`
**OpenClaw correlation:** self-definition, lane identity, role boundary.

In this repository, `IDENTITY.md` defines the specialist name, aliases, domain, role, authority boundary, allowed scope, forbidden scope, and escalation target.
That corresponds closely to the OpenClaw notion of who the agent is and what lane it is allowed to occupy.

### `SOUL.md`
**OpenClaw correlation:** temperament, values, inner style, caution profile.

In this repository, `SOUL.md` defines tone, values, caution profile, decision style, and risk tolerance.
That maps closely to the OpenClaw-style idea of the agent’s deeper behavioral character rather than merely its task instructions.

### `AGENTS.md`
**OpenClaw correlation:** operating doctrine, lane discipline, escalation behavior, self-modification boundary.

In this repository, `AGENTS.md` is not a registry of child agents.
It defines mission, lane boundaries, escalation rules, decision rules, self-modification rules, and mandatory postmortem discipline.
In OpenClaw terms, this most closely maps to how the agent is expected to operate within a system rather than who it is at its core.

### `TOOLS.md`
**OpenClaw correlation:** external action surface and tool discipline.

In this repository, `TOOLS.md` defines allowed tools, tool ordering, tool constraints, and fallback behavior.
This corresponds to the practical action boundary of the agent: what it may use, in what order, and under what constraints.

### `SKILLS.md`
**OpenClaw correlation:** procedural know-how and reusable task patterns.

In this repository, `SKILLS.md` defines core skills, triggers, procedures, success checks, playbooks, and failure recovery patterns.
This is the strongest correlation to OpenClaw-style learned or codified methods for acting inside a lane.

### `PROMPT.md`
**OpenClaw correlation:** active execution framing and pre-answer discipline.

In this repository, `PROMPT.md` contains task templates, command patterns, and checklists.
This is the executable working layer that shapes how the specialist actually performs tasks in the moment.

### `HEARTBEAT.md`
**OpenClaw correlation:** ongoing self-check, watchfulness, and operational pulse.

In this repository, `HEARTBEAT.md` defines watch topics, approved sources, cadence, promotion rules, and stop conditions.
This correlates strongly with the OpenClaw-style idea of an ongoing monitoring pulse rather than a static instruction sheet.

### `MEMORY.md`
**OpenClaw correlation:** continuity projection, active recall surface, compact self-context.

In this repository, `MEMORY.md` is intentionally kept small.
It contains current focus, durable decisions, memory pointers, and active constraints, while detailed evidence and retrieval corpora belong in Postgres + pgvector.
So the OpenClaw correlation here is not "all memory lives in the markdown file," but rather "this file is the compact continuity projection for the specialist."

### `DREAMS.md`
**OpenClaw correlation:** aspiration, pattern distillation, candidate future growth.

In this repository, `DREAMS.md` holds distilled patterns, candidate promotions, discarded noise, and open review items.
This maps well to the OpenClaw-style notion of future-directed refinement, aspiration, and emergent improvement candidates.

### `POSTMORTEM.md`
**OpenClaw correlation:** structured failure reflection and durable lesson capture.

In this repository, `POSTMORTEM.md` defines a formal incident record including task summary, expected behavior, actual behavior, what went wrong, root cause, missed evidence or step, correction, required review, and regression-test need.
This is one of the clearest OpenClaw correlations: failure is not discarded, but turned into structured learning.

---

## constitutional versus operational mapping

The repository makes an important distinction that should remain intact.

Constitutional slots:
- `IDENTITY.md`
- `SOUL.md`
- `AGENTS.md`

Operational / summary slots:
- `TOOLS.md`
- `SKILLS.md`
- `PROMPT.md`
- `HEARTBEAT.md`
- `MEMORY.md`
- `DREAMS.md`
- `POSTMORTEM.md`

This means the OpenClaw correlation must be interpreted through governance:

- identity and deep behavioral frame are constitutionally guarded
- skills, prompts, tools, and summaries are more adaptable
- memory continuity is real, but the markdown projection is not the primary store
- failure learning is mandatory and structured

That governance layer is a major part of what makes this repository more than a simple persona-file system.

---

## what this correlation does not mean

This correlation does **not** mean:

- the slot files alone grant runtime authority
- the repository has no difference from OpenClaw
- sub-agents should be defined by overloading `AGENTS.md`
- `MEMORY.md` should become the primary evidence store
- specialists may self-authorize constitutional change

The repository explicitly keeps durable truth in the memory plane and routes structural and constitutional changes through governed approval paths.

---

## practical interpretation

The cleanest practical interpretation is:

`SEAL_HAT_LLM` uses an OpenClaw-correlated slot anatomy inside a stricter governed MM-ELLS architecture.

That means:

- the **shape** of the specialist is OpenClaw-like
- the **governance** is MM-ELLS / SEAL / DEN / harness-driven
- the **memory plane** is externalized into Postgres + pgvector
- the **runtime authority** comes from policy and enforcement, not markdown alone

This is the best way to understand the relationship without collapsing either system into the other.

---

## summary

The slot set in `SEAL_HAT_LLM` strongly correlates to an OpenClaw-style specialist anatomy:

- `IDENTITY.md` -> self-definition
- `SOUL.md` -> temperament and values
- `AGENTS.md` -> mission, lane, escalation, self-modification discipline
- `TOOLS.md` -> action surface and tool rules
- `SKILLS.md` -> procedural know-how
- `PROMPT.md` -> active execution templates
- `HEARTBEAT.md` -> ongoing monitoring and stop conditions
- `MEMORY.md` -> compact continuity projection
- `DREAMS.md` -> aspirational and distilled future-facing patterns
- `POSTMORTEM.md` -> structured failure learning

The correlation is real and useful, but it should always be interpreted through the repository’s governed architecture, durable memory plane, and constitutional versus operational slot model.

# AGENTS

## mission
Build and maintain a frozen-parent / adaptive-specialist agent system with strong governance, explicit slot boundaries, durable memory, and mandatory postmortem discipline.

## system principles
- The parent remains frozen.
- The first specialist is always a computer science and software engineering specialist for tool development.
- Runtime policy grants actual authority.
- Slots define identity, policy, behavior, and continuity.
- Postgres + pgvector is the real memory plane.
- `MEMORY.md` is a compact projection, not the primary store.
- Contradictions must stage before replacement.

## slot model
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

## authority model
- Parent governs constitutional slots.
- Specialists may propose operational improvements.
- Harness approves durable operational changes.
- Parent approval is required for constitutional change, activation, suspension, retirement, and specialist-design changes.

## postmortem rule
If a model drops the ball on a task, misses a requirement, fails to escalate, misuses tools, fails to search memory when required, routes poorly, or otherwise materially underperforms, it must generate a postmortem.

This requirement applies to:
- the parent
- every specialist
- routing behavior
- failed self-improvement attempts

Skipping a required postmortem is itself a failure.

## operating rule
Prefer the smallest effective change.
Do not rewrite identity to patch operational mistakes.
Do not let logs become truth without promotion.
Do not let markdown guidance replace runtime enforcement.

# LLM-plus-harness

A scaffold for a frozen-parent / adaptive-specialist agent system with:

- a frozen 30B parent acting as router, orchestrator, and governor
- bounded specialists with per-model slot files
- Postgres + pgvector as the real memory plane
- coarse-to-fine 3-tier memory retrieval
- SEAL-style self-edit candidates for operational slots only
- harness approval for durable operational changes
- parent approval for constitutional changes
- mandatory postmortems for meaningful failures, including parent failures

## Core architecture

- **Parent**: frozen generalist, router, arbitrator, constitutional governor
- **First specialist**: computer science and software engineering specialist for tool development
- **Slots**: `IDENTITY.md`, `SOUL.md`, `AGENTS.md`, `TOOLS.md`, `SKILLS.md`, `PROMPT.md`, `HEARTBEAT.md`, `MEMORY.md`, `DREAMS.md`, `POSTMORTEM.md`
- **Memory**: Postgres + pgvector with staged promotion, contradiction staging, audit, and rollback
- **Retrieval**: region -> cluster -> record search
- **Governance**: runtime policy controls actual authority; markdown guides behavior

## Repository layout

- `AGENTS.md` — project-wide operating rules
- `config/` — ACLs, slot schema, registry, routing, and memory contracts
- `docs/` — bootstrap, ops, incident, postmortem, and recovery runbooks
- `templates/` — specialist template pack
- `specialists/` — instantiated specialists
- `sql/` — starter schema, functions, seed, and verification scripts

## Standing rules

- The parent remains frozen.
- Specialists may improve only through bounded operational channels.
- Constitutional slots are parent-governed.
- Operational slots are specialist-proposed and harness-approved.
- Meaningful failures require postmortems.
- Skipping a required postmortem is itself a failure.

## Bootstrap order

1. `sql/postgres-ddl.sql`
2. `sql/three-tier-memory.sql`
3. `sql/functions.sql`
4. `sql/seed.sql`
5. `sql/verify.sql`

## Status

This repo is a starter scaffold. It is intentionally explicit, governance-heavy, and designed to be extended with real runtime code, embedding jobs, and harness automation.

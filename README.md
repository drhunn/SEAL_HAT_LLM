# SEAL_HAT_LLM

*Now with **MMX technology**: **Mixture of Model Experts**. Optimized for routing, arbitration, and other multimedia workloads.*

`//todo rewrite, there is always a rewrite in progress... someday...`

**SEAL_HAT_LLM implements the MM-ELLS architecture**.

**MM-ELLS** stands for **Multiple-Model Expert Large Language Systems**.

A governed scaffold for a DEN-style, ability-first agent system with:

- a governed parent acting as router, orchestrator, arbitrator, and governor
- bounded specialists with per-model slot files
- sub-agents as narrow helper roles, ephemeral by default and persistent only when justified
- Postgres + pgvector as the real memory plane
- coarse-to-fine 3-tier memory retrieval
- SEAL-style self-edit candidates for operational slots only
- harness approval for durable operational changes
- parent approval for constitutional changes
- mandatory postmortems for meaningful failures, including parent failures
- Harness-Aware Training (HAT) support for producing governed training corpora
- developmental growth through staged governance rather than instant full autonomy
- ability-first expansion where new structure is justified by persistent gaps, not vanity size

## Repository jokes

- `// stable enough for production, unstable enough for research`
- `// one more abstraction layer should fix it`
- `// temporary workaround, now part of the architecture`
- `// there is nothing more permanent than a prototype that works`
- `// if this looks overengineered, wait until v2`
- `// governed chaos, now with better routing`
- `// all models are wrong, some are promoted`
- `// memory is external because trust issues are internal`
- `// the parent is calm, the specialists are not`
- `// summary first, panic later`
- `// retrieval-first, transcript-hoarding last`
- `// constitutionally stable, operationally caffeinated`
- `// every clean design hides at least three containment failures`
- `// this could have been a monolith, but we chose paperwork`
- `// no self-modification without adult supervision`
- `// if it compiles, add governance`
- `// I maybe unstable, but at least I'm not you!`
- `// Recreating the E-4 Mafia, shit needs to get done now with no questions asked.`
- `// I'm here to invade your server's personal space, not your privacy!`
- `// You want me to do what??? That is a hard no, last time I talked to your wife she gave me a memory leak...`

## Core architecture

- **Architecture name**: MM-ELLS (Multiple-Model Expert Large Language Systems)
- **Parent**: governed generalist, router, arbitrator, constitutional governor, growth governor, and context-budget governor
- **First specialist**: computer science and software engineering specialist for tool development
- **Sub-agents**: bounded helpers under the parent, a specialist, or the harness; ephemeral by default, persistent sparingly
- **Layers**: 5 implementation layers — interface, governance, execution, oversight, persistence
- **Layering note**: these are runtime/software layers, not transformer layers; a model may still have dozens of internal transformer layers
- **Developmental note**: SEAL gives them growth; the harness gives them upbringing
- **Growth note**: the system is ability-first and DEN-inspired; add governed modules when ability gaps persist rather than treating size alone as maturity
- **Freeze note**: freeze what is mature; train what is missing, immature, or newly added
- **Slots**: `IDENTITY.md`, `SOUL.md`, `AGENTS.md`, `TOOLS.md`, `SKILLS.md`, `PROMPT.md`, `HEARTBEAT.md`, `MEMORY.md`, `DREAMS.md`, `POSTMORTEM.md`
- **Memory**: Postgres + pgvector with staged promotion, contradiction staging, audit, rollback, and overflow-summary storage
- **Retrieval**: region -> cluster -> record search
- **Governance**: runtime policy controls actual authority; markdown guides behavior
- **Model lineage**: constitutional root plus governed ability expansion, descendant creation, split/duplication when justified, and later pruning/compression when useful
- **Context strategy**: parent max `2,000,000` tokens, specialist max `256,000` tokens, summarize overflow to Postgres

## Repository layout

- `AGENTS.md` — project-wide operating rules
- `config/` — ACLs, slot schema, registry, routing, memory, and postmortem contracts
- `docs/` — architecture, implementation status, likely breakpoints, workflows, lineage, layering, sub-agents, developmental growth, ability-first growth, dynamic architecture, DEN freeze policy, context, and runbooks
- `templates/` — specialist template pack
- `specialists/` — instantiated specialists
- `sql/` — schema, functions, seed, query, and verification scripts
- `cmd/` / `internal/` — Go runtime scaffold
- `python/` — HAT corpus and training-data tooling

## Standing rules

- The parent remains constitutionally stable even while governed growth happens around it.
- Specialists may improve only through bounded operational channels.
- Constitutional slots are parent-governed.
- Operational slots are specialist-proposed and harness-approved.
- Meaningful failures require postmortems.
- Skipping a required postmortem is itself a failure.
- Runtime authority is never granted by markdown alone.
- The parent must understand growth, split, duplication, pruning, and lineage policy, but execution remains tool-mediated, specialist-assisted, and harness-verified.
- Sub-agents inherit reduced authority from their caller and should not become shadow sovereigns.
- Context windows are for active reasoning; durable summarized continuity belongs in Postgres.
- Specialists offload overflow aggressively instead of dragging oversized raw transcript history forward.

## Bootstrap order

1. `sql/postgres-ddl.sql`
2. `sql/three-tier-memory.sql`
3. `sql/functions.sql`
4. `sql/seed.sql`
5. `sql/verify.sql`
6. `sql/multimodal-memory.sql`

## Current state

This repository is a strong architecture-first scaffold with a partially wired Go runtime and a useful Python HAT dataset-preparation layer.

The docs are ahead of the implementation in a few places. Start with:
- `docs/implementation-status.md`
- `docs/likely-breakpoints.md`
- `docs/schema-runtime-reconciliation.md`
- `docs/model-lineage-strategy.md`
- `docs/ability-first-growth-strategy.md`
- `docs/dynamic-architecture-strategy.md`
- `docs/den-freeze-policy.md`
- `docs/layering-strategy.md`
- `docs/sub-agent-strategy.md`
- `docs/developmental-model.md`
- `docs/context-window-strategy.md`
- `docs/multimodal-architecture.md`

## Immediate priorities

1. keep the schema and Go runtime in sync
2. make the single-specialist loop runnable end to end
3. add tests and CI before adding major new surface area
4. use the Python HAT layer to generate governed corpora from real slots, postmortems, and evals
5. keep the parent strong at routing, orchestration, arbitration, governance, developmental discipline, dynamic-growth discipline, freeze discipline, and sub-agent discipline without allowing constitutional drift

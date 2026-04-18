# SEAL_HAT_LLM

**SEAL_HAT_LLM implements the MM-ELLS architecture**.

**MM-ELLS** stands for **Multiple-Model Expert Large Language Systems**.

A governed scaffold for a frozen-parent / adaptive-specialist agent system with:

- a frozen 30B parent acting as router, orchestrator, arbitrator, and governor
- bounded specialists with per-model slot files
- Postgres + pgvector as the real memory plane
- coarse-to-fine 3-tier memory retrieval
- SEAL-style self-edit candidates for operational slots only
- harness approval for durable operational changes
- parent approval for constitutional changes
- mandatory postmortems for meaningful failures, including parent failures
- Harness-Aware Training (HAT) support for producing governed training corpora

## Core architecture

- **Architecture name**: MM-ELLS (Multiple-Model Expert Large Language Systems)
- **Parent**: frozen generalist, router, arbitrator, constitutional governor, and model-family planner
- **First specialist**: computer science and software engineering specialist for tool development
- **Slots**: `IDENTITY.md`, `SOUL.md`, `AGENTS.md`, `TOOLS.md`, `SKILLS.md`, `PROMPT.md`, `HEARTBEAT.md`, `MEMORY.md`, `DREAMS.md`, `POSTMORTEM.md`
- **Memory**: Postgres + pgvector with staged promotion, contradiction staging, audit, and rollback
- **Retrieval**: region -> cluster -> record search
- **Governance**: runtime policy controls actual authority; markdown guides behavior
- **Model lineage**: canonical 30B root, governed distillation, pruning only after distillation when justified

## Repository layout

- `AGENTS.md` — project-wide operating rules
- `config/` — ACLs, slot schema, registry, routing, memory, and postmortem contracts
- `docs/` — architecture, implementation status, likely breakpoints, workflows, lineage, and runbooks
- `templates/` — specialist template pack
- `specialists/` — instantiated specialists
- `sql/` — schema, functions, seed, query, and verification scripts
- `cmd/` / `internal/` — Go runtime scaffold
- `python/` — HAT corpus and training-data tooling

## Standing rules

- The parent remains frozen.
- Specialists may improve only through bounded operational channels.
- Constitutional slots are parent-governed.
- Operational slots are specialist-proposed and harness-approved.
- Meaningful failures require postmortems.
- Skipping a required postmortem is itself a failure.
- Runtime authority is never granted by markdown alone.
- The parent must understand distillation and pruning policy, but execution remains tool-mediated, specialist-assisted, and harness-verified.

## Bootstrap order

1. `sql/postgres-ddl.sql`
2. `sql/three-tier-memory.sql`
3. `sql/functions.sql`
4. `sql/seed.sql`
5. `sql/verify.sql`

## Current state

This repository is a strong architecture-first scaffold with a partially wired Go runtime and a useful Python HAT dataset-preparation layer.

The docs are ahead of the implementation in a few places. Start with:
- `docs/implementation-status.md`
- `docs/likely-breakpoints.md`
- `docs/schema-runtime-reconciliation.md`
- `docs/model-lineage-strategy.md`

## Immediate priorities

1. keep the schema and Go runtime in sync
2. make the single-specialist loop runnable end to end
3. add tests and CI before adding major new surface area
4. use the Python HAT layer to generate governed corpora from real slots, postmortems, and evals
5. keep the parent strong at routing, orchestration, arbitration, governance, and model-family planning without allowing constitutional drift

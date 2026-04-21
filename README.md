# SEAL_HAT_LLM

SEAL_HAT_LLM is a governed scaffold for a specialist-runtime research project whose target shape is:
- a frozen parent / adaptive specialist model family
- **one harness per model**
- **one embedded Postgres per model**
- explicit smaller specialist models derived from the base model
- a later shared RPC/IPC tool plane that multiple harnesses can call without sharing memory or authority
- a Go runtime scaffold for routing, execution, health, and persistence
- a Python HAT layer for dataset building and training support

## Project goal

The goal of this project is to research **governed adaptive model systems** in which:
- **SEAL** determines when evidence justifies change
- **DEN** provides the structural change mechanism
- the **harness** enforces review, rollback, and execution boundaries
- **slot governance** defines what is mutable, what is protected, and how changes are proposed

The parent is intended to be a **generalist** that becomes a **learning routing/orchestration layer** over time.
That means the parent should get better at deciding:
- when to handle work itself
- when to route to a specialist
- when to invoke retrieval first
- when multimodal fusion is required
- when to escalate, defer, or refuse

The specialists are not supposed to be vague helper personas or internal MoE shards.
They are supposed to be **explicit smaller models derived from the base model**, shaped for narrower task families through techniques such as distillation, pruning, freezing, and bounded adaptation.

The architectural intent is to replace a traditional mixture-of-experts style internal expert arrangement with **governed external specialist models** that:
- own real task lanes
- run faster and cheaper than the parent on those lanes
- preserve higher resolution on repeated narrow work
- remain auditable, reviewable, and reversible system components

The first instantiated specialist is the **ComputerScience-SoftwareEngineering specialist**, which exists to build the tooling, runtime, harness, and evaluation infrastructure, and to prepare governed descendant copies of the base model for later specialist creation through pruning, distillation, freezing, and bounded adaptation.

Each model is intended to run as a **self-contained model unit** with:
- its own model artifact
- its own harness
- its own embedded Postgres
- its own slots, memory, evals, postmortems, and lifecycle state

Later versions are intended to move tools into a **shared RPC/IPC tool plane** so multiple harnesses can use the same tool executable while still enforcing their own local permissions, memory rules, and governance checks.

That learning does **not** replace governance.
The parent may learn how to route and orchestrate more effectively, but the harness and slot governance still define what is admissible, what requires review, and what structural changes are allowed.

In plain English:
this project exists to test whether a specialist system can improve through bounded, evidence-driven adaptation without becoming opaque, unreviewable, or structurally sloppy.

The intended loop is:
1. run bounded specialist tasks through model-local harnesses
2. record failures, telemetry, evals, and memory in model-local stores
3. let SEAL decide whether the evidence justifies change
4. let the harness and slot governance decide whether that change is admissible
5. let DEN perform the approved structural change in a reversible form
6. evaluate whether the change actually helped

## Status

This repository is **not** a production runtime.
It is a **working scaffold** with real database, verification, and dataset-building paths, plus a lot of architecture that is still ahead of the implementation.

The current codebase already supports:
- TOML config loading
- specialist slot loading and canonical bundle compilation
- Postgres-backed storage and retrieval helpers
- bounded runtime task processing
- a conservative file-backed task inbox
- telemetry, proposal, and growth-plan scaffolding
- strict CI with Postgres bootstrap and `cmd/verify`
- Python dataset generation from repo slots and optional Postgres corpora

It does **not** yet provide:
- per-model harness runtime units across parent and specialists
- per-model embedded Postgres deployment
- a shared RPC/IPC tool plane with reusable external tool executables
- production multi-specialist orchestration
- production multimodal execution
- production training orchestration

## Design rules

These are the load-bearing rules for the repo:
- runtime policy grants actual authority
- markdown does not grant authority by itself
- constitutional slots are parent-governed
- operational improvements stay bounded and reviewable
- model-local memory is authoritative by default
- cross-model sharing must be explicit and governed
- meaningful failures require postmortems
- schema/runtime drift is treated as a real defect

## Quick start

### 1. Bootstrap the database
Use the shared bootstrap script:
- `make bootstrap-db`

Or call it directly with a DSN:
- `bash ./scripts/bootstrap_verify_db.sh 'postgres://postgres:postgres@localhost:5432/llm_harness?sslmode=disable'`

The script applies:
1. `sql/postgres-ddl.sql`
2. `sql/three-tier-memory.sql`
3. `sql/functions.sql`
4. `sql/seed.sql`
5. `sql/verify.sql`
6. `sql/multimodal-memory.sql`
7. `sql/20260420_seal_den_growth.sql`
8. `sql/20260421_route_episodes.sql`

### 2. Run verification
- soft local bring-up: `make verify`
- strict admission check: `make verify-strict`

### 3. Run the harness scaffold
- `make run`

## Repository layout

- `AGENTS.md` — project-wide operating rules
- `config/` — runtime config and related scaffolding
- `docs/` — status, contracts, architecture notes, and runbooks
- `specialists/` — instantiated specialist slot files
- `templates/` — specialist template pack
- `sql/` — schema, functions, seed data, and verification SQL
- `scripts/` — shared operational scripts
- `cmd/` — executable entrypoints
- `internal/` — Go runtime packages
- `python/` — HAT dataset and training-support tooling

## Read these first

Start with the documents that reflect reality instead of aspiration:
- `docs/README.md`
- `docs/implementation-status.md`
- `docs/likely-breakpoints.md`
- `docs/sql-contracts.md`
- `docs/hat-llm.md`

## Current priorities

1. keep the Go runtime and SQL schema in sync
2. keep the single-specialist loop stable and fully durable
3. reduce duplicated routing / execution policy
4. avoid adding new surface area before the existing verify path is solid

## Documentation cleanup note

The pre-cleanup versions of the README and the core docs touched in this pass are preserved under:
- `docs/archive/2026-04-20-pre-cleanup/`

# SEAL_HAT_LLM

SEAL_HAT_LLM is a governed scaffold for a slot-driven specialist runtime with:
- a frozen parent / adaptive specialist model
- Postgres + pgvector as the durable memory plane
- a Go runtime scaffold for routing, execution, health, and persistence
- a Python HAT layer for dataset building and training support

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
- production multi-specialist orchestration
- a real tool broker
- remote intake APIs or queue backends
- production multimodal execution
- production training orchestration

## Design rules

These are the load-bearing rules for the repo:
- runtime policy grants actual authority
- markdown does not grant authority by itself
- constitutional slots are parent-governed
- operational improvements stay bounded and reviewable
- durable memory lives in Postgres, not in chat transcript sprawl
- meaningful failures require postmortems
- schema/runtime drift is treated as a real defect

## Quick start

### 1. Bootstrap the database
Use the shared bootstrap script:
- `make bootstrap-db`

Or call it directly with a DSN:
- `./scripts/bootstrap_verify_db.sh 'postgres://postgres:postgres@localhost:5432/llm_harness?sslmode=disable'`

The script applies:
1. `sql/postgres-ddl.sql`
2. `sql/three-tier-memory.sql`
3. `sql/functions.sql`
4. `sql/seed.sql`
5. `sql/verify.sql`
6. `sql/multimodal-memory.sql`
7. `sql/20260420_seal_den_growth.sql`

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

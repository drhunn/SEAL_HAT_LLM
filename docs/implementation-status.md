# Implementation status

This file separates the **working code paths** from the **documented target architecture**.

## Bottom line

SEAL_HAT_LLM is a working scaffold.
It is not a production runtime.

The repo already has enough real code to support verification, bounded task processing, canonical slot compilation, Postgres-backed persistence, and Python dataset generation.
It does not yet have a full live specialist system, per-model harness/runtime units, per-model embedded Postgres deployment, a shared RPC/IPC tool plane, or production growth machinery.

## Implemented now

### Go runtime
The Go side already has:
- TOML config loading
- Postgres connection pooling
- specialist slot loading from the filesystem
- canonical slot bundle compilation and TOML emission
- bounded task processing through `runtime.Service.ProcessTask`
- a conservative file-backed task inbox
- harness incident handling for execution failures and review-required signals
- retrieval helpers and routing audit write paths
- modality-aware routing and execution planning
- simulated model-host execution wiring
- telemetry collection and signal helpers
- SEAL proposal generation scaffolding
- DEN growth-plan generation scaffolding
- lifecycle, oversight, and lineage scaffolding
- `cmd/verify` as a real admission path
- an initial `internal/unit` abstraction and local-runtime bootstrap path that now builds the runtime as a **local model unit** while still using the current shared DSN mode

### SQL layer
The SQL side already has:
- base specialist / slot / memory tables
- coarse-to-fine memory hierarchy tables
- retrieval and health helper functions
- postmortem, eval, and self-edit staging tables
- multimodal scaffold tables
- SEAL/DEN tables for signals, proposals, plans, lineage, promotion decisions, and bundle versions
- seed and verify scripts

### Python HAT layer
The Python side already has:
- repository slot loading
- governed dataset building
- optional Postgres corpus ingestion
- governance-pressure negative example generation
- JSON / JSONL export helpers
- optional `DatasetDict` export support
- a JAX training scaffold
- unit tests for policy and dataset basics

### CI and checks
The repo already has:
- GitHub Actions CI
- Go build and test coverage in CI
- DB-backed runtime integration coverage via `TEST_DATABASE_DSN`
- strict verify in CI against a bootstrapped Postgres + pgvector service
- Python compile checks and unit tests

## Implemented, but still narrow

These paths exist, but they are intentionally small and not yet broad production systems:
- bounded single-task orchestration
- file-backed local task intake
- slot bundle persistence
- telemetry-driven proposal generation
- growth-plan staging
- multimodal-aware planning
- model-unit bootstrap ownership cleanup without per-unit embedded storage yet

## Still scaffolded or partial

These areas are present in design and partially present in code, but not complete:
- DB/filesystem slot synchronization as a full live workflow
- lifecycle enforcement as a real state machine
- context-overflow summarization as an always-on runtime behavior
- durable multimodal execution beyond the current bounded path
- governed dynamic growth beyond proposal and plan staging
- deeper sub-agent orchestration
- route-episode capture for learned parent routing
- specialist artifact lifecycle beyond names, slot packs, and growth-plan scaffolding

## Not implemented yet

These are still outside the current runtime:
- per-model harness runtime units across parent and specialists
- per-model embedded Postgres deployment
- explicit cross-model replication or governed sharing between model-local stores
- a shared RPC/IPC tool plane with reusable external tool executables
- production multi-specialist orchestration
- production tool broker integration across multiple harnesses
- remote intake APIs
- distributed queue backends
- production multimodal stack
- benchmarked production training workflows
- full promotion / rollback enforcement for growth experiments

## Known weak spots

The repo is most likely to fail when:
- the Go runtime and SQL schema drift apart
- migrations are not applied in the target database
- verify is treated as a vague smoke test instead of a hard contract check
- docs are read as implementation proof
- duplicated policy logic drifts across packages
- the shared-DB scaffold is mistaken for the final per-model embedded-DB architecture

## Next milestone

The next worthwhile milestone is still simple:

**the single-specialist loop works end to end and persists durably without special pleading**

That means:
1. config loads
2. slots load
3. the canonical bundle compiles and persists
4. the database contracts hold
5. retrieval works
6. postmortem and eval writes work
7. health updates work
8. strict verify passes cleanly
9. telemetry, proposal, and growth-plan writes are durable
10. the task inbox is explicit about what it guarantees and what it does not

After that, the repo can make the larger architectural turn toward:
- model units with local harnesses
- model units with embedded Postgres
- governed cross-model sharing
- a reusable RPC/IPC tool plane

## Maintenance rule

When code changes the actual runtime surface, update this file in the same change set.
Do not let architecture notes quietly replace status reporting.

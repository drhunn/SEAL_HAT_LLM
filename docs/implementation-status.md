# Implementation status

This file separates the **working code paths** from the **documented target architecture**.

## Bottom line

SEAL_HAT_LLM is a working scaffold.
It is not a production runtime.

The repo already has enough real code to support verification, bounded task processing, canonical slot compilation, Postgres-backed persistence, and Python dataset generation.
It does not yet have a full live specialist system, per-model harness/runtime units, a mature per-model embedded Postgres lifecycle, a shared RPC/IPC tool plane, or production growth machinery.

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
- an initial `internal/unit` abstraction and local-runtime bootstrap path that now builds the runtime as a **local model unit**
- an explicit unit-owned store bootstrap abstraction with:
  - `shared_dsn` bootstrap support
  - an initial actual `embedded_postgres` bootstrap path using local PostgreSQL binaries
  - embedded schema bootstrap through a reusable Go SQL-file runner
  - a bootstrapped local-runtime constructor now used by the main harness entrypoint
  - embedded-store ownership locking, stale-lock reclaim, readiness ping checks, schema bootstrap, and idempotent stop behavior for the initial local lifecycle
- an initial `unit.Registry` that resolves the current unit and built-in known units for routing/execution target resolution
- routing/execution target fields that now carry **unit-target metadata** and resolve known units through the registry before falling back to compatibility mapping
- execution planning support for explicitly preferred **unit IDs** in addition to preferred executor aliases
- routing support for explicitly preferred **unit IDs** in addition to executor-policy selection
- startup-task construction and task-inbox parsing paths that now carry `PreferredUnitID` into live runtime task objects
- runtime task processing now carries `PreferredUnitID` through routing and execution instead of dropping it
- an initial Unix-domain-socket task RPC transport with typed request/response contracts, client/server support, and a runtime adapter that exposes `runtime.Service.ProcessTask` over the socket boundary
- optional harness-side task RPC serving from `cmd/harness/main.go` when `runtime.enable_task_rpc_server` is enabled
- route-episode persistence support for successful startup-task and inbox-task execution paths
- failed route-episode persistence support for startup-task failures and inbox-task failures via the current runtime wrapper seam
- current specialist-artifact persistence support for the active local model unit on startup
- growth staging support that now creates a **candidate artifact**, links the experiment to that candidate, and records the current artifact as the parent reference
- specialist artifact event history support for startup registration, candidate growth staging, and oversight-triggered promotion/rollback event hooks
- initial artifact lifecycle helpers that can promote a candidate artifact to current or roll it back through the experiment path
- direct test coverage for preferred-unit execution planning, preferred-unit routing, full runtime preferred-unit task processing, oversight artifact-event hooks, unit store-bootstrap selection, config embedded-postgres defaults, unit spec store-mode derivation, embedded-postgres ownership/readiness/schema-bootstrap/cleanup behavior, and Unix-socket task RPC transport/runtime-adapter behavior

### SQL layer
The SQL side already has:
- base specialist / slot / memory tables
- coarse-to-fine memory hierarchy tables
- retrieval and health helper functions
- postmortem, eval, and self-edit staging tables
- multimodal scaffold tables
- SEAL/DEN tables for signals, proposals, plans, lineage, promotion decisions, and bundle versions
- a route-episode table for durable routing/execution episode capture
- a specialist-artifact table for durable model-unit artifact records
- an artifact-ref column on `ability_growth_experiments` for linking staged growth work to the candidate artifact for that experiment
- a specialist-artifact-events table for durable artifact lifecycle/event history
- a candidate-artifact lifecycle migration that normalizes `active` -> `current` and enforces one current artifact per specialist
- a reusable Go schema bootstrap runner that applies the same ordered SQL files as the verify bootstrap script
- verify checks for one-current-artifact invariants and growth/artifact lifecycle alignment
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
- model-unit bootstrap ownership cleanup with store bootstrap abstraction now in place, but without a mature embedded-store lifecycle yet
- initial unit-registry resolution layered on top of the existing executor policy rather than full multi-unit routing
- route-episode persistence on both success and failure paths through the current runtime wrappers rather than a fully centralized end-of-task hook
- current/candidate artifact separation with candidate rows created during growth staging
- artifact lifecycle history through event records plus narrow promotion/rollback helpers, rather than a full lifecycle state engine
- initial embedded-postgres bootstrap support that still depends on local PostgreSQL binaries and now applies the repo SQL bootstrap sequence, but still lacks broader lifecycle supervision
- initial task RPC transport that exposes specialist-side task execution over a local socket but does not yet provide a full parent-side distributed dispatcher

## Still scaffolded or partial

These areas are present in design and partially present in code, but not complete:
- DB/filesystem slot synchronization as a full live workflow
- lifecycle enforcement as a real state machine across the wider runtime
- context-overflow summarization as an always-on runtime behavior
- durable multimodal execution beyond the current bounded path
- governed dynamic growth beyond proposal and plan staging
- deeper sub-agent orchestration
- a deeper core task-processing hook for failed-task route episodes instead of the current wrapper seam
- specialist artifact lifecycle beyond current/candidate separation, event history, and narrow promotion/rollback helpers
- embedded Postgres lifecycle management beyond initial bootstrap, ownership locking, stale-lock reclaim, readiness checks, schema bootstrap, and stop behavior
- parent-side RPC dispatch and result orchestration across multiple model units

## Not implemented yet

These are still outside the current runtime:
- per-model harness runtime units across parent and specialists
- a mature per-model embedded Postgres deployment and supervision model
- explicit cross-model replication or governed sharing between model-local stores
- routing that targets real model units end to end without the old executor policy as the primary decision source across the wider runtime
- a shared RPC/IPC tool plane with reusable external tool executables
- production multi-specialist orchestration
- production tool broker integration across multiple harnesses
- remote intake APIs
- distributed queue backends
- production multimodal stack
- benchmarked production training workflows
- full promotion / rollback enforcement for growth experiments across the wider runtime
- DEN-produced model-unit bundles with full artifact packaging
- broad exercised lifecycle tests beyond the current preferred-unit and oversight hook coverage

## Known weak spots

The repo is most likely to fail when:
- the Go runtime and SQL schema drift apart
- migrations are not applied in the target database
- docs are read as implementation proof
- duplicated policy logic drifts across packages
- the initial embedded-postgres bootstrap path is mistaken for a production-quality embedded-store lifecycle
- the initial task RPC path is mistaken for a full parent↔specialist distributed runtime
- initial unit-registry resolution is mistaken for real multi-unit orchestration
- candidate artifact creation is mistaken for full bundle production
- narrow promotion/rollback helpers are mistaken for a full lifecycle state engine

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
9. telemetry, proposal, growth-plan writes are durable
10. the task inbox is explicit about what it guarantees and what it does not

After that, the repo can make the larger architectural turn toward:
- model units with local harnesses
- model units with embedded Postgres
- governed cross-model sharing
- a reusable RPC/IPC tool plane

## Maintenance rule

When code changes the actual runtime surface, update this file in the same change set.
Do not let architecture notes quietly replace status reporting.

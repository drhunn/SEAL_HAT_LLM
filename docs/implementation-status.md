# Implementation status

This file separates the **working code paths** from the **documented target architecture**.

## Bottom line

SEAL_HAT_LLM is still a research/prototype scaffold, not a production runtime.

This branch now has real code for:
- bounded runtime task processing
- unit-target-aware routing and execution planning
- explicit unit-owned store bootstrap
- embedded Postgres schema bootstrap from a configured SQL root
- hardened local task RPC transport
- parent-side remote task dispatch over Unix sockets
- centralized route-episode accounting
- DB-backed route-episode proof tests
- an artifact lifecycle state machine
- concrete acceptance criteria for the runtime-task-system work

It still does **not** have production multi-specialist orchestration, production service discovery, a shared external tool plane, distributed queues, model artifact packaging, or production embedded Postgres supervision.

## Implemented now

### Go runtime

The Go side has:
- TOML config loading
- Postgres connection pooling
- specialist slot loading from the filesystem
- canonical slot bundle compilation and TOML emission
- bounded task processing through `runtime.Service.ProcessTask`
- file-backed task inbox support
- harness incident handling for execution failures and review-required signals
- retrieval helpers and routing audit write paths
- modality-aware routing and execution planning
- simulated model-host execution wiring
- telemetry collection and signal helpers
- SEAL proposal generation scaffolding
- DEN growth-plan generation scaffolding
- `cmd/verify` as a real admission path
- `internal/unit` local-runtime bootstrap support
- `shared_dsn` store bootstrap
- initial `embedded_postgres` store bootstrap using local PostgreSQL binaries
- explicit `embedded_postgres.sql_root` configuration with default `./sql`
- embedded schema bootstrap through a reusable Go SQL-file runner
- embedded-store ownership locking, stale-lock reclaim, readiness ping checks, schema bootstrap, cleanup behavior, and idempotent stop behavior
- `unit.Registry` target resolution through `unitref` without routing/execution importing the full `unit` package
- routing and execution support for preferred unit IDs
- startup-task and task-inbox parsing paths that carry `PreferredUnitID`
- runtime task processing that carries `PreferredUnitID` through routing, execution, and remote dispatch
- Unix-domain-socket task RPC with typed request/response contracts, protocol versioning, request/response size guards, configurable deadlines, and malformed/protocol/oversize rejection tests
- optional harness-side task RPC serving from `cmd/harness/main.go` when `runtime.enable_task_rpc_server` is enabled
- parent-side `internal/taskdispatch` that maps `runtime.Task` to `taskrpc.RunTaskRequest`, resolves explicit unit-to-socket mappings, and returns structured/partial dispatch results
- runtime remote-dispatch hookup through `runtime.RemoteDispatcher`
- configured dispatcher attachment during local runtime construction
- remote task dispatch from `ProcessTask` when routing targets a mapped non-local unit
- centralized route-episode creation from the runtime task outcome hook
- DB-backed route-episode tests for local success, local failure, remote success, remote failure, and unhandled host outcomes
- DB-backed end-to-end remote dispatch test using a live task-RPC Unix socket server
- artifact lifecycle state machine with legal transitions, illegal transition rejection, promotion, rollback, event writes, and one-current behavior

### SQL layer

The SQL side has:
- base specialist / slot / memory tables
- coarse-to-fine memory hierarchy tables
- retrieval and health helper functions
- postmortem, eval, and self-edit staging tables
- multimodal scaffold tables
- SEAL/DEN tables for signals, proposals, plans, lineage, promotion decisions, and bundle versions
- route-episode table for durable routing/execution episode capture
- specialist-artifact table for model-unit artifact records
- specialist-artifact-events table for lifecycle history
- artifact-ref column on `ability_growth_experiments`
- candidate/current artifact lifecycle migration
- unique current-artifact index for the one-current-artifact invariant
- ordered schema bootstrap through `internal/db/bootstrap.go`
- seed and verify scripts

### Python HAT layer

The Python side has:
- repository slot loading
- governed dataset building
- optional Postgres corpus ingestion
- governance-pressure negative example generation
- JSON / JSONL export helpers
- optional `DatasetDict` export support
- JAX training scaffold
- unit tests for policy and dataset basics

### CI and checks

The repo has:
- GitHub Actions CI
- Go build and test coverage in CI
- DB-backed runtime integration coverage via `TEST_DATABASE_DSN`
- strict verify in CI against a bootstrapped Postgres + pgvector service
- Python compile checks and unit tests

This branch still needs a fresh observed green validation run after the latest runtime, dispatch, route-episode, and lifecycle changes.

## Implemented, but still narrow

These paths exist, but they are intentionally small and not yet broad production systems:
- bounded single-task orchestration
- file-backed local task intake
- slot bundle persistence
- telemetry-driven proposal generation
- growth-plan staging
- multimodal-aware planning
- embedded Postgres bootstrap and cleanup, but not production supervision
- explicit static unit-to-socket remote dispatch, but not discovery
- route-episode accounting, but not full distributed tracing
- artifact lifecycle state enforcement, but not model artifact packaging or signed bundle promotion
- DB-backed proof tests skipped unless `TEST_DATABASE_DSN` is set

## Still scaffolded or partial

These areas remain partial:
- DB/filesystem slot synchronization as a live workflow
- context-overflow summarization as an always-on runtime behavior
- durable multimodal execution beyond the bounded path
- governed dynamic growth beyond proposal and plan staging
- deeper sub-agent orchestration
- production multi-unit orchestration
- production embedded Postgres supervision
- service discovery for task RPC endpoints
- retry/backoff policy for remote dispatch
- cross-unit result fusion
- production tool-plane integration
- DEN-produced model-unit bundles with full artifact packaging

## Not implemented yet

These are still outside the current runtime:
- production per-model harness runtime units across parent and specialists
- mature per-model embedded Postgres deployment and supervision
- explicit cross-model replication or governed sharing between model-local stores
- shared RPC/IPC external tool plane with reusable tool executables
- remote intake APIs
- distributed queue backends
- production multimodal stack
- benchmarked production training workflows
- artifact signing, storage, and rollout orchestration

## Known weak spots

The repo is most likely to fail when:
- the Go runtime and SQL schema drift apart
- migrations are not applied in the target database
- docs are read as implementation proof
- duplicated routing/execution policy drifts
- embedded Postgres bootstrap is mistaken for production supervision
- task RPC plus static dispatch is mistaken for a full distributed runtime
- route episodes are mistaken for full distributed tracing
- artifact lifecycle state enforcement is mistaken for full model artifact packaging
- DB-backed tests are skipped and nobody notices

## Admission gates

Use `docs/acceptance-criteria.md` as the branch admission checklist.

At minimum, this branch needs:

```bash
go build ./...
go test ./...
python -m compileall python
python -m unittest discover -s python/tests
```

With a database:

```bash
TEST_DATABASE_DSN='postgres://postgres:postgres@localhost:5432/llm_harness?sslmode=disable' go test ./...
go run ./cmd/verify -config config/runtime.example.toml -mode strict
```

## Maintenance rule

When code changes the actual runtime surface, update this file or a more specific reality-anchor doc in the same change set.
Do not let architecture notes quietly replace status reporting.

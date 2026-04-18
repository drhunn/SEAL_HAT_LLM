# GO RUNTIME OVERVIEW

## purpose
Describe the current Go runtime scaffold included in `SEAL_HAT_LLM`.

This runtime is still a scaffold, but it is no longer just a thin placeholder. It now includes a runnable main path, a standalone verify command, retrieval wrappers, and more explicit service boundaries.

---

## current packages

### `cmd/harness`
Main runtime entrypoint.
Loads config, opens Postgres, loads specialist slot files, builds core services, and starts the runtime.

### `cmd/verify`
Standalone verification entrypoint.
Checks config loading, database connectivity, slot loading, health snapshot access, and retrieval smoke-test access.

### `internal/config`
Loads runtime TOML configuration.

### `internal/db`
Creates the Postgres connection pool.

### `internal/slots`
Loads slot markdown files for a specialist from the filesystem.

### `internal/slotsync`
Provides a snapshot-style slot synchronization helper.

### `internal/memory`
Provides Postgres-backed helpers for:
- creating postmortems
- updating health
- projecting `MEMORY.md`
- fetching health snapshots
- staging eval cases
- staging self-edit candidates
- writing routing audit records
- running coarse-to-fine retrieval

### `internal/postmortem`
Wraps postmortem creation logic.

### `internal/evals`
Stages eval candidates into the memory plane.

### `internal/routing`
Provides starter routing decisions and routing-decision scaffolding.

### `internal/lifecycle`
Provides specialist lifecycle status updates.

### `internal/harness`
Provides the core harness service for:
- startup checks
- incident handling
- auto-postmortem flow
- health update flow
- recovery-plan generation
- parent-review signaling
- degraded-state recommendation support

### `internal/harness/workflows`
Contains starter recovery workflow types and planning logic.

### `internal/runtime`
Coordinates runtime lifecycle, startup smoke paths, and core service wiring.

---

## current capabilities
The Go scaffold currently supports:
- loading runtime config
- opening Postgres
- loading slot files from `specialists/`
- basic startup checks
- standalone verification through `cmd/verify`
- creating postmortems
- staging eval cases
- staging self-edit candidates
- updating specialist health
- projecting `MEMORY.md`
- fetching health snapshots
- writing routing audit records
- running coarse-to-fine retrieval wrappers
- slot-sync snapshot logging

---

## what is still scaffold-level
The runtime still needs fuller implementation for:
- actual model-host integration and live parent/specialist invocation
- production tool-broker integration
- fully enforced lifecycle state machine transitions
- richer slot DB/filesystem reconciliation
- contradiction review workflows
- mature eval execution
- multimodal runtime extensions
- context-overflow offload orchestration into Postgres summaries

---

## startup path
Typical startup path:
1. load `config/runtime.example.toml`
2. connect to Postgres
3. create `memory.PostgresStore`
4. create `postmortem.Service`
5. create `evals.Service`
6. create `lifecycle.Service`
7. create `harness.Service`
8. create `routing.Service`
9. create `runtime.Service`
10. run startup checks
11. load specialist slots
12. run smoke-test wiring paths
13. remain alive until shutdown signal

---

## development intent
The design intent is:
- explicit over magical
- governed over self-authorizing
- recoverable over clever
- database-backed memory over hidden process memory
- parent-governed constitutional state
- harness-gated operational adaptation
- retrieval-first continuity instead of transcript hoarding

---

## verify path
Use the standalone verification path when reconciling schema and runtime:
- `go run ./cmd/verify -config config/runtime.example.toml`
- `make verify`

This should be part of the normal reconciliation loop after SQL or runtime changes.

---

## next implementation targets
Strong next targets for the Go runtime are:
- real model-host abstraction for parent and specialists
- stronger lifecycle transition enforcement
- durable slot reconciliation logic
- explicit context-overflow summarization and Postgres offload
- multimodal record and routing support
- eval execution and activation gating beyond staging

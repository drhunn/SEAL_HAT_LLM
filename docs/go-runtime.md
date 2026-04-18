# GO RUNTIME OVERVIEW

## purpose
Describe the starter Go runtime scaffold included in this repository.

This runtime is not yet a full autonomous harness. It is a structured skeleton intended to grow into one.

---

## current packages

### `cmd/harness`
Main entrypoint.
Loads config, opens Postgres, loads specialist slot files, builds services, and starts the runtime.

### `internal/config`
Loads runtime TOML configuration.

### `internal/db`
Creates the Postgres connection pool.

### `internal/slots`
Loads slot markdown files for a specialist from the filesystem.

### `internal/memory`
Provides a Postgres-backed store for:
- creating postmortems
- updating health
- projecting `MEMORY.md`
- fetching health snapshots

### `internal/postmortem`
Wraps postmortem creation logic.

### `internal/harness`
Provides the core harness service for:
- startup checks
- incident handling
- auto-postmortem flow
- health update flow
- parent-review signaling

### `internal/runtime`
Coordinates the runtime lifecycle.

### `internal/evals`
Starter service for staging eval candidates.

### `internal/routing`
Starter routing service.

### `internal/slotsync`
Starter slot synchronization snapshot service.

### `internal/harness/workflows`
Starter recovery workflow types.

---

## current capabilities
The Go scaffold currently supports:
- loading runtime config
- opening Postgres
- loading slot files from `specialists/`
- basic startup checks
- creating postmortems
- updating specialist health
- projecting `MEMORY.md`
- logging routing/eval/slot-sync placeholders

---

## what is still scaffold-level
The runtime still needs fuller implementation for:
- actual LLM orchestration
- real 3-tier retrieval invocation from Go
- durable memory staging and promotion flows
- contradiction review workflows
- eval registry and execution
- routing audit persistence
- lifecycle state transitions
- parent/specialist arbitration logic
- slot validation and slot writes

---

## startup path
Typical startup path:
1. load `config/runtime.example.toml`
2. connect to Postgres
3. create `memory.PostgresStore`
4. create `postmortem.Service`
5. create `harness.Service`
6. create `runtime.Service`
7. run startup checks
8. load specialist slots
9. remain alive until shutdown signal

---

## development intent
The design intent is:
- explicit over magical
- governed over self-authorizing
- recoverable over clever
- database-backed memory over hidden process memory
- parent-governed constitutional state
- harness-gated operational adaptation

---

## next implementation targets
Strong next targets for the Go runtime are:
- query wrappers for `fn_run_coarse_to_fine_search`
- routing audit writes
- eval case staging writes
- self-edit candidate staging writes
- slot projection and sync verification
- degraded/suspended lifecycle transitions
- specialist activation/shadow helpers

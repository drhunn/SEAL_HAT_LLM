# Acceptance criteria

This file defines the current admission gates for the runtime-task-system work.

A feature is not done because the code exists.
A feature is done when the behavior is implemented, tested, documented, and validated.

## Global admission gates

The branch must pass:

```bash
go build ./...
go test ./...
python -m compileall python
python -m unittest discover -s python/tests
```

With a test database:

```bash
TEST_DATABASE_DSN='postgres://postgres:postgres@localhost:5432/llm_harness?sslmode=disable' go test ./...
go run ./cmd/verify -config config/runtime.example.toml -mode strict
```

No feature below is complete until the global gates pass.

## Embedded Postgres bootstrap

### In scope
- explicit SQL root configuration through `embedded_postgres.sql_root`
- default SQL root of `./sql`
- ordered SQL bootstrap from the configured root
- embedded Postgres schema bootstrap after readiness ping
- cleanup on open, ping, and schema-bootstrap failures
- ownership lock and stale-lock behavior

### Out of scope
- production-grade embedded Postgres supervision
- multi-process clustering
- database migration framework replacement
- remote database provisioning

### Required tests
- config default test for `embedded_postgres.sql_root`
- configured SQL root is passed into embedded bootstrap
- cleanup on pool open failure
- cleanup on ping failure
- cleanup on bootstrap failure
- idempotent stop
- stale lock reclaim
- fresh lock rejection

### Done means
- `internal/db` tests pass
- DB bootstrap no longer silently depends on current working directory except through the documented default
- `docs/sql-contracts.md` documents SQL root behavior

## Task RPC transport

### In scope
- Unix-domain-socket request/response transport
- protocol version stamping and validation
- request and response size limits
- configurable client/server deadlines
- malformed JSON rejection
- unsupported protocol version rejection
- oversized request rejection
- socket cleanup on shutdown

### Out of scope
- network transport
- authentication
- streaming
- multiplexing
- production service discovery

### Required tests
- client/server round trip
- handler error propagation
- client-side oversized request rejection
- server malformed JSON rejection
- server unsupported protocol rejection
- server oversized request rejection
- socket removal on shutdown
- repeated task-RPC test run to expose flakes

### Done means
- `go test ./internal/taskrpc -count=20` passes locally or in CI
- protocol behavior is documented in `docs/task-dispatch.md`

## Parent-side task dispatch

### In scope
- explicit `task_dispatch.remote_unit_sockets` config
- config cleanup of empty socket-map entries
- parent-side dispatcher mapping `runtime.Task` to `taskrpc.RunTaskRequest`
- no silent fallback for mapped remote dispatch failures
- structured partial result on RPC failure
- runtime seam through `runtime.RemoteDispatcher`
- dispatcher attached during local runtime construction when configured

### Out of scope
- automatic service discovery
- retry policy
- distributed queueing
- cross-unit result fusion
- remote tool-plane execution

### Required tests
- config parses and cleans remote socket mappings
- dispatcher requires preferred unit ID
- dispatcher rejects unmapped target unit
- dispatcher calls resolved client
- dispatcher returns partial result on client error
- runtime helper chooses remote dispatch only for mapped non-local targets

### Done means
- `internal/config`, `internal/taskdispatch`, and `internal/runtime` tests pass
- docs describe explicit socket mapping and no-silent-fallback behavior

## End-to-end remote dispatch

### In scope
- parent runtime routes to mapped remote unit
- live Unix-socket task-RPC server receives request
- parent receives remote response
- execution result uses `remote_rpc`
- parent records remote route episode

### Out of scope
- multi-hop routing
- multiple specialist result synthesis
- production monitoring
- retry/backoff policy

### Required tests
- DB-backed integration test starts live task-RPC server
- task with `PreferredUnitID` dispatches to remote unit
- remote output returns to parent
- route episode records `status=succeeded`
- route episode records `execution_mode=remote_rpc`
- route episode records the remote target unit

### Done means
- `TEST_DATABASE_DSN=... go test ./internal/runtime -run TestProcessTaskDispatchesToRemoteRPCServerAndRecordsRouteEpisode` passes

## Route episode accounting

### In scope
- centralized route episode creation from runtime task outcome hook
- local success route episode
- local failure route episode
- remote success route episode
- remote failure route episode
- unhandled host route episode
- status derivation: `succeeded`, `failed`, `unhandled`

### Out of scope
- full distributed tracing
- append-only event sourcing
- payload archival
- replay system

### Required tests
- pure status-derivation test
- DB-backed test writes and reads all required outcome classes

### Done means
- route episodes persist under `TEST_DATABASE_DSN`
- `docs/sql-contracts.md` lists the route-episode write contract

## Artifact lifecycle state machine

### In scope
- artifact states: `defined`, `candidate`, `current`, `rolled_back`, `rejected`, `archived`
- legal transition table
- illegal transition rejection
- candidate promotion to current
- rollback path
- event row written for each transition
- one-current-artifact invariant preserved by archiving prior current artifacts

### Out of scope
- artifact bundle building
- model packaging
- artifact signing
- artifact storage backend
- multi-specialist rollout orchestration

### Required tests
- pure legal-transition tests
- pure illegal-transition tests
- DB-backed promote/rollback test
- illegal archived/current transition rejection
- one-current-artifact count check
- event row checks

### Done means
- `go test ./internal/lifecycle` passes
- `TEST_DATABASE_DSN=... go test ./internal/lifecycle` passes
- `docs/sql-contracts.md` documents lifecycle states and transitions

## Documentation admission

Every slice must update relevant docs in the same change set.

Required docs for this branch:
- `docs/acceptance-criteria.md`
- `docs/documentation-maintenance.md`
- `docs/sql-contracts.md`
- `docs/task-dispatch.md`
- `docs/implementation-status.md` when implementation status changes materially

## Branch admission

Before merge:
- branch must be rebased or merged with latest `main`
- branch must not be behind `main`
- CI must be green
- stale duplicate feature branches should be retired after merge

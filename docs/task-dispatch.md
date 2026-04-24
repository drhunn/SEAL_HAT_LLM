# Task dispatch

This document describes the current parent-side task dispatch seam.

## Current status

`internal/taskdispatch` is an initial parent-side dispatch package.

It currently supports:
- resolving a target unit ID to a task-RPC Unix socket path
- loading explicit remote unit socket mappings from `task_dispatch.remote_unit_sockets`
- mapping `runtime.Task` into `taskrpc.RunTaskRequest`
- stamping the task-RPC protocol version on outbound requests
- calling a task-RPC client
- returning a structured dispatch result containing the target unit, socket path, request, and response
- returning partial dispatch results when the remote RPC call fails
- attaching the configured dispatcher to the local runtime during `unit.NewLocalRuntime`
- dispatching mapped non-local task targets from `runtime.Service.ProcessTask` over task RPC
- recording local and remote task outcomes through the centralized route-episode hook
- DB-backed route-episode persistence tests for local, remote, failure, and unhandled outcomes
- an end-to-end DB-backed remote dispatch integration test using a live Unix-socket task-RPC server

## Configuration

Remote dispatch is intentionally explicit.
There is no discovery system yet.

Example:

```toml
[task_dispatch.remote_unit_sockets]
"csse-tool-development-specialist-01" = "./artifacts/taskrpc/csse-tool-development-specialist-01.sock"
```

The config loader trims empty keys and empty socket paths.
If a target unit is not mapped, dispatch does not silently pretend it succeeded.
Silent fallback would make routing dishonest.

## Runtime behavior

`ProcessTask` still performs local retrieval and routing first.

If routing resolves a target unit that is:
- not the current local unit, and
- present in `task_dispatch.remote_unit_sockets`, and
- a runtime remote dispatcher is attached,

then the task is sent over task RPC.

The remote response is converted into an `execution.Result` with:
- `ExecutionMode = "remote_rpc"`
- `TargetUnitID` from the routing decision
- host metadata containing the remote socket path and remote status

## Route episode accounting

`runtime.Service.handleTaskOutcome` records route episodes for local and remote task outcomes.

Current status derivation:
- `succeeded` when execution completed and the host result was handled
- `failed` when execution or remote dispatch returned an error
- `unhandled` when execution returned without an error but the host did not handle the task

Route episode persistence is best-effort inside the outcome hook.
A route-episode write failure is logged and does not hide the original task result.

## Proof tests

DB-backed proof tests are skipped unless `TEST_DATABASE_DSN` is set.

Current proof coverage:
- `internal/memory/route_episode_test.go` verifies route episode rows for local success, local failure, remote success, remote failure, and unhandled host outcomes.
- `internal/runtime/remote_dispatch_integration_test.go` starts a live task-RPC Unix-socket server, dispatches a parent runtime task to the mapped remote unit, and verifies the parent wrote a `remote_rpc` route episode.

## What it does not do yet

It does not yet provide:
- full multi-unit orchestration
- retry policy
- timeout policy beyond the task-RPC client options
- cross-unit result fusion
- durable distributed task observation beyond route-episode rows
- production service discovery

## Boundary

The dispatcher is intentionally small.
It is a parent-side RPC seam, not a finished distributed runtime.

The next real step is hard validation: run the full Go test suite, DB-backed tests, strict verify, and build checks.

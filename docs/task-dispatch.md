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

The next real step is broader validation: prove local success, local failure, remote success, and remote failure all record route episodes under a real database-backed test.

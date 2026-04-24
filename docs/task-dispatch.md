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

## Configuration

Remote dispatch is intentionally explicit.
There is no discovery system yet.

Example:

```toml
[task_dispatch.remote_unit_sockets]
"csse-tool-development-specialist-01" = "./artifacts/taskrpc/csse-tool-development-specialist-01.sock"
```

The config loader trims empty keys and empty socket paths.
If a target unit is not mapped, dispatch must fail loudly.
Silent fallback would make routing dishonest.

## What it does not do yet

It does not yet provide:
- full multi-unit orchestration
- route-episode persistence around remote dispatch
- retry policy
- timeout policy beyond the task-RPC client options
- cross-unit result fusion
- durable distributed task observation
- production service discovery

## Boundary

The dispatcher is intentionally small.
It is a parent-side RPC seam, not a finished distributed runtime.

The next real step is to wire this seam into runtime orchestration and persist dispatch success/failure as route episodes.

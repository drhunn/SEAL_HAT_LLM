# Changelog

This file tracks meaningful repository changes in a blunt, reviewable way.
It is for reference, not marketing.

## Unreleased

### Added
- `internal/unit/spec.go` with a first-class model-unit spec, role, store mode, and tool-plane mode.
- `internal/unit/local_runtime.go` to bootstrap the runtime as a local model unit instead of wiring everything directly in `cmd/harness/main.go`.
- `internal/unit/spec_test.go` covering parent vs specialist unit-spec derivation.
- `internal/unitref/target.go` and `internal/unitref/target_test.go` for compatibility mapping from executor names to unit targets.
- `sql/20260421_route_episodes.sql` for durable route-episode persistence.
- `internal/memory/route_episodes.go` with `CreateRouteEpisode(...)`.
- `internal/runtime/route_episode.go` with a runtime helper to persist route episodes from task results.
- `sql/20260421_specialist_artifacts.sql` for durable specialist-artifact records.
- `internal/memory/specialist_artifacts.go` with `UpsertSpecialistArtifact(...)`.
- `internal/runtime/specialist_artifact.go` with a runtime helper to persist the current specialist model-unit artifact.

### Changed
- `cmd/harness/main.go` now boots through a local model-unit builder instead of hand-wiring the runtime stack directly.
- `internal/routing/service.go` now carries compatibility unit-target metadata (`TargetUnitID`, `TargetRole`, `TargetModelRef`) alongside the existing executor target.
- `internal/execution/service.go` now carries the same compatibility unit-target metadata through execution plans.
- `internal/runtime/bootstrap.go` now persists route episodes for successful startup-task runs and successful inbox-task executions.
- `internal/runtime/bootstrap.go` now also persists the current specialist artifact on startup using the exported slot-packet reference and version hash.
- `scripts/bootstrap_verify_db.sh` now applies `sql/20260421_route_episodes.sql` and `sql/20260421_specialist_artifacts.sql`.
- `README.md` now includes the new route-episode and specialist-artifact migrations in the bootstrap list.
- `docs/implementation-status.md` now records the unit-bootstrap change, compatibility unit-target routing, initial route-episode persistence, and initial specialist-artifact persistence.

### Notes
- This is still transitional code.
- The runtime still uses the shared DSN mode.
- Compatibility unit-target mapping is not real multi-unit orchestration yet.
- Route-episode persistence currently covers known successful task paths, not every failure path.
- Specialist-artifact persistence currently records the current local model unit on startup; it is not yet a full DEN-produced model-unit bundle pipeline.

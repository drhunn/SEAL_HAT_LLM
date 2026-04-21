# Changelog

This file tracks meaningful repository changes in a blunt, reviewable way.
It is for reference, not marketing.

All timestamps are recorded in **UTC**.
Each entry includes:
- the time the change landed on `main`
- a short summary of what changed
- the concrete files or behaviors affected

## 2026-04-21 19:31:44 UTC — Add an initial embedded Postgres bootstrap path and switch the harness entrypoint to bootstrapped runtime startup (pending)
**Summary:** Cut the next real store-boundary slice by making store mode configurable, adding an initial embedded Postgres bootstrap path, and switching the harness entrypoint to the bootstrapped local-runtime path instead of opening the database itself.

**Changes:**
- Updated `cmd/harness/main.go` to use `NewBootstrappedLocalRuntime(...)`.
- Updated `internal/config/config.go` to add `runtime.store_mode` and `[embedded_postgres]` settings with conditional validation/defaults.
- Updated `config/runtime.example.toml` to document both store mode and embedded Postgres settings.
- Added `internal/config/config_test.go` covering embedded-postgres config defaults without requiring `database.dsn`.
- Added `internal/db/embedded_postgres.go` with an initial local PostgreSQL binary bootstrap path.
- Updated `internal/unit/spec.go` and `internal/unit/spec_test.go` so unit specs derive store mode from config.
- Updated `internal/unit/store_bootstrap.go` and `internal/unit/store_bootstrap_test.go` to support real embedded-postgres bootstrap selection.
- Updated `docs/implementation-status.md` to record the initial embedded-postgres bootstrap path and bootstrapped entrypoint adoption.

## 2026-04-21 14:12:41 UTC — Carry preferred unit ids through runtime task construction and add direct hook coverage (pending)
**Summary:** Pushed explicit preferred unit ids farther into the live runtime by wiring them into startup-task and task-inbox construction, and added direct tests for preferred-unit execution planning and oversight artifact-event hooks.

**Changes:**
- Updated `internal/runtime/bootstrap.go` so the startup task now sets `PreferredUnitID` to the current runtime specialist.
- Updated `internal/runtime/task_inbox.go` so inbox-submitted task JSON can carry `preferred_unit_id` and preserve it in the runtime task object.
- Updated `internal/execution/service_test.go` to match the unit-registry constructor shape and to add direct preferred-unit execution planning coverage.
- Added `internal/oversight/service_test.go` with direct promotion/rollback artifact-event hook coverage.
- Updated `docs/implementation-status.md` to record the new preferred-unit task construction paths and direct hook coverage.

## 2026-04-21 11:51:25 UTC — Tighten lifecycle verify checks and explicit preferred-unit planning (pending)
**Summary:** Fixed the worst coherence gap by making verify check lifecycle/artifact invariants and by letting execution planning prefer a real target unit ID instead of only executor aliases.

**Changes:**
- Updated `sql/verify.sql` to check lifecycle/artifact invariants instead of only basic table existence.
- Added verify checks for specialist-artifact tables, route-episode tables, artifact-ref columns, and the single-current-artifact index.
- Added verify checks for one-current-artifact invariants and growth/artifact status alignment across candidate, current, and rolled-back artifacts.
- Updated `internal/runtime/service.go` to add `PreferredUnitID` to the runtime task model and `TargetUnitID` to task review records.
- Updated `internal/execution/service.go` so execution planning can prefer an explicit target unit ID and resolve its executor alias through the unit registry.
- Updated `docs/implementation-status.md` to record lifecycle verify coverage and preferred-unit execution support.

## 2026-04-21 06:00:02 UTC — Add candidate artifact lifecycle helpers (#12)
**Summary:** Added the first real current/candidate artifact separation so staged experiments point at distinct candidate artifacts instead of scribbling on the current artifact row.

**Changes:**
- Added `sql/20260421_candidate_artifact_lifecycle.sql`.
- Normalized existing `active` artifact rows to `current`.
- Enforced one `current` artifact per specialist through the schema.
- Updated `internal/memory/specialist_artifacts.go` so current-artifact upsert targets the single-current-artifact rule.
- Added `internal/memory/specialist_artifact_lifecycle.go` with candidate-artifact creation and narrow promotion/rollback helpers.
- Updated `internal/growth/service.go` so staged experiments create a distinct candidate artifact and link the experiment to that candidate.
- Updated `scripts/bootstrap_verify_db.sh` and `README.md` to include the candidate-artifact lifecycle migration.
- Updated `docs/implementation-status.md` to describe current/candidate separation and narrow lifecycle helpers without pretending the full lifecycle engine exists.

## 2026-04-21 05:27:20 UTC — Add specialist artifact event history hooks (#11)
**Summary:** Added durable artifact event history and wired the first lifecycle breadcrumbs into startup, growth staging, and oversight seams.

**Changes:**
- Added `sql/20260421_specialist_artifact_events.sql`.
- Added `internal/memory/specialist_artifact_events.go` for durable artifact event persistence.
- Added `internal/memory/specialist_artifact_experiment_events.go` to record artifact events from experiment references.
- Updated `internal/runtime/specialist_artifact.go` to record a `registered` artifact event on startup.
- Updated `internal/growth/service.go` to record a `growth_staged` artifact event when a staged experiment links to the current artifact.
- Updated `internal/oversight/service.go` so promotion and rollback can emit artifact events when the backing writer supports it.
- Updated `scripts/bootstrap_verify_db.sh` and `README.md` to include the artifact-event migration.
- Updated `docs/implementation-status.md` to describe artifact event history as implemented, while explicitly saying this is not a full lifecycle state engine.

## 2026-04-21 05:09:25 UTC — Add initial unit registry for known model units (#10)
**Summary:** Replaced more of the executor-string fiction with a declared registry of known units and routed target resolution through it.

**Changes:**
- Added `ExecutorName` to `internal/unit/spec.go`.
- Added `internal/unit/registry.go`.
- Added `internal/unit/registry_test.go`.
- Updated `internal/unit/local_runtime.go` so local runtime bootstrap owns a unit registry.
- Updated `internal/routing/service.go` to resolve targets through the registry first and compatibility mapping second.
- Updated `internal/execution/service.go` to resolve targets through the registry first and compatibility mapping second.
- Updated `docs/implementation-status.md` to say the repo now has initial unit-registry support, but not full multi-unit orchestration.

## 2026-04-21 04:56:23 UTC — Persist failed route episodes from runtime wrappers (#9)
**Summary:** Closed the obvious route-episode blind spot so failed startup-task and inbox-task runs leave durable records instead of disappearing into error handling.

**Changes:**
- Updated `internal/runtime/route_episode.go` to add failed-task route-episode persistence.
- Updated `internal/runtime/bootstrap.go` so startup-task failures persist failed route episodes.
- Updated `internal/runtime/bootstrap.go` so inbox-task failures persist failed route episodes.
- Updated `docs/implementation-status.md` to state that failed-route persistence now exists through the current wrapper seam.

## 2026-04-21 04:47:56 UTC — Link ability growth experiments to specialist artifacts (#8)
**Summary:** Tied staged growth work to the current specialist artifact so experiments stop floating free of artifact identity.

**Changes:**
- Added `sql/20260421_ability_growth_artifact_refs.sql`.
- Added `internal/memory/specialist_artifact_lookup.go`.
- Updated `internal/memory/ability_growth.go` to carry `SpecialistArtifactID` into ability growth experiments.
- Updated `internal/growth/service.go` so staged experiments try to attach the current specialist artifact ID.
- Updated `scripts/bootstrap_verify_db.sh` and `README.md` to include the new migration.
- Updated `docs/implementation-status.md` to describe artifact-linked growth staging and to say clearly that this is not full bundle production.

## 2026-04-21 04:38:27 UTC — Add initial specialist artifact persistence (#7)
**Summary:** Gave the current local specialist a durable artifact record instead of leaving it as a name, a slot pack, and routing metadata.

**Changes:**
- Added `sql/20260421_specialist_artifacts.sql`.
- Added `internal/memory/specialist_artifacts.go`.
- Added `internal/runtime/specialist_artifact.go`.
- Updated `internal/runtime/bootstrap.go` so startup persists the current specialist artifact using the slot-packet path and version hash.
- Updated `scripts/bootstrap_verify_db.sh` and `README.md` to include the specialist-artifact migration.
- Updated `docs/implementation-status.md` to say initial specialist-artifact persistence exists, but the full specialist lifecycle still does not.

## 2026-04-21 04:26:00 UTC — Introduce model-unit bootstrap and route episode persistence (#6)
**Summary:** Landed the first real code-alignment slice for the model-unit architecture: local unit bootstrap, compatibility unit targets, and durable route-episode persistence.

**Changes:**
- Added `internal/unit/spec.go`.
- Added `internal/unit/local_runtime.go`.
- Updated `cmd/harness/main.go` to bootstrap through the local model-unit builder.
- Added `internal/unit/spec_test.go`.
- Added `internal/unitref/target.go` and `internal/unitref/target_test.go`.
- Updated `internal/routing/service.go` to carry compatibility unit-target metadata.
- Updated `internal/execution/service.go` to carry compatibility unit-target metadata.
- Added `sql/20260421_route_episodes.sql`.
- Added `internal/memory/route_episodes.go`.
- Added `internal/runtime/route_episode.go`.
- Updated `internal/runtime/bootstrap.go` to persist route episodes for startup-task and inbox-task success paths.
- Updated `scripts/bootstrap_verify_db.sh` and `README.md` to include the route-episode migration.
- Added the initial `CHANGELOG.md` file.
- Updated `docs/implementation-status.md` to reflect the new model-unit bootstrap and route-episode support.

## 2026-04-21 03:50:59 UTC — Document model-unit architecture and shared tool plane (#5)
**Summary:** Made the intended target architecture explicit: one harness per model, one embedded Postgres per model, and a later shared RPC/IPC tool plane.

**Changes:**
- Updated `README.md` to frame the target architecture as local model units with a later shared tool plane.
- Updated `docs/seal-den-system-architecture.md` to describe parent and specialist model units with local harnesses and local embedded Postgres.
- Updated `docs/parent-routing-orchestration-strategy.md` to describe routing across model units and later shared-tool planning.
- Updated `docs/model-lineage-strategy.md` to treat specialists as model-unit bundles rather than bare model blobs.
- Updated `docs/implementation-status.md` to state clearly that the code had not implemented that architecture yet.
- Updated `docs/README.md` so the docs index points to the right architecture notes.

## 2026-04-21 03:18:14 UTC — Document the first specialist as the specialist-factory engineer (#4)
**Summary:** Made the role of the first instantiated specialist explicit so the docs stop hand-waving about how later specialists get prepared.

**Changes:**
- Updated the first specialist’s identity docs to describe it as both the system toolsmith and the early specialist-factory engineer.
- Updated the first specialist mission to include descendant-model preparation work.
- Updated the first specialist skills to include specialist descendant preparation.
- Updated `README.md` to mention why the first instantiated specialist exists.

## 2026-04-21 02:56:59 UTC — Document specialists as explicit derived models (#3)
**Summary:** Clarified that specialists are supposed to be explicit smaller derived models for narrow task families, not vague helper personas or internal MoE shards.

**Changes:**
- Updated the project goal docs to define specialists as smaller derived models.
- Updated the architecture docs to align with specialist-model intent.
- Updated the routing strategy to match the explicit specialist-model design.
- Updated the lineage strategy to define specialists as explicit descendant models.

## 2026-04-21 02:30:16 UTC — Document the parent as a learned routing/orchestration layer (#2)
**Summary:** Made the project goal explicit that the parent is meant to become a learned routing/orchestration layer under governance, not just a static generalist.

**Changes:**
- Updated `README.md` to add the explicit project goal.
- Updated the project goal docs to make learned parent routing/orchestration explicit.
- Updated `docs/seal-den-system-architecture.md` to define the learned parent router/orchestrator role.
- Added `docs/parent-routing-orchestration-strategy.md`.
- Updated `docs/README.md` to link the parent routing/orchestration strategy.

## 2026-04-21 01:49:59 UTC — Tighten docs and reduce verify/runtime policy drift (#1)
**Summary:** Cleaned up the repo’s core docs, backed up the pre-cleanup versions, and tightened runtime/verify behavior so the code and admission path drift less.

**Changes:**
- Backed up the pre-cleanup README and core docs under the archive path.
- Tightened `README.md`, the docs index, implementation status, likely breakpoints, SQL contracts, and HAT documentation.
- Centralized executor selection policy and reused it across routing and execution.
- Added typed vector literal helpers and coverage.
- Required explicit opt-in for starter slot fallback.
- Extracted the verify runner from the `cmd/verify` main path.
- Added the shared verify database bootstrap script and `bootstrap-db` target.
- Updated CI and docs to use the shared bootstrap script consistently.
- Reduced duplicated slot loading in the starter-fallback reporting path.

# SCHEMA-RUNTIME RECONCILIATION

## purpose
Provide a practical checklist for keeping the Go runtime and SQL layer aligned.

---

## current reconciliation goals

### goal 1
The Go runtime should only call SQL helpers that actually exist in `sql/functions.sql` or later migration-backed runtime files.

### goal 2
The Go runtime should only write columns and tables that actually exist in `sql/postgres-ddl.sql` and related schema files.

### goal 3
Every retrieval wrapper should match the exact SQL return shape.

### goal 4
Every new runtime DB write path should be reflected in:
- SQL schema
- SQL contracts doc
- verify/check workflow

### goal 5
The runtime and schema must agree on identifier shape and schema namespace.

For the current SEAL/DEN scaffolding, that means:
- `agent_core` remains the canonical schema namespace
- the runtime and SQL both treat SEAL/DEN IDs as text-style IDs rather than generated DB UUIDs

---

## reconciliation checklist

### 1. specialists table
Confirm runtime expectations for:
- `status`
- `health_score`
- `notes`

### 2. routing_audit table
Confirm runtime expectations for:
- `initial_classifier`
- `impact`
- `was_fallback`
- `fallback_reason`
- `was_override`
- `override_by`
- `multi_specialist_review`
- `notes`

### 3. eval/self-edit tables
Confirm these exist if runtime writes to them:
- `memory_eval_cases`
- `memory_self_edit_candidates`

### 4. health functions
Confirm these exist and return the expected shape:
- `fn_update_specialist_health_stats`
- `fn_get_specialist_health_snapshot`

### 5. memory projection function
Confirm `fn_project_and_write_memory_summary_slot` exists and can update the `MEMORY.md` slot safely.

### 6. retrieval function
Confirm `fn_run_coarse_to_fine_search` returns exactly:
- `record_id`
- `record_kind`
- `title`
- `summary`
- `body`
- `status`
- `semantic_score`
- `cluster_score`
- `final_score`

### 7. canonical slot bundle persistence
Confirm these are aligned:
- `slot_bundle_versions` table exists
- Go bundle compiler emits the same source map shape the DB writer stores
- `cmd/verify` can persist a bundle after compilation

### 8. SEAL/DEN persistence tables
Confirm these exist if runtime writes to them:
- `adaptation_signals`
- `gap_clusters`
- `adaptation_proposals`
- `growth_plans`
- `growth_experiments`
- `lineage_nodes`
- `lineage_edges`
- `promotion_decisions`

### 9. signal/proposal/plan write paths
Confirm runtime expectations for:
- signal `id` as text, not DB-generated UUID
- proposal `id` as text
- growth-plan `id` as text
- `specialist_id` foreign keys pointing at `agent_core.specialists`
- JSON payload columns accepting the marshaled Go structures actually emitted by the runtime

### 10. verify command
Run:
- `make verify`
- `go run ./cmd/verify -config config/runtime.example.toml`

Confirm that `cmd/verify` can:
- compile the canonical slot bundle
- attempt to persist the slot bundle version
- write retrieval/routing/execution signals
- generate SEAL proposals
- generate DEN plans

### 11. build/test
Run:
- `go build ./...`
- `go test ./...`
- `python -m compileall python`

---

## workflow for future changes
Whenever a runtime DB call is added or modified:
1. update schema if needed
2. update helper functions if needed
3. update `docs/sql-contracts.md`
4. run verify/build/test
5. update `docs/implementation-status.md` if the implementation meaningfully changed
6. update `docs/likely-breakpoints.md` if the change creates a new failure mode or migration dependency

---

## summary
The repo is strongest when the architecture, SQL, and Go runtime all describe the same system.
This file exists to make that alignment a deliberate process rather than an accident, especially now that canonical slot bundling and SEAL/DEN persistence paths have been added.

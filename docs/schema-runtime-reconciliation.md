# SCHEMA-RUNTIME RECONCILIATION

## purpose
Provide a practical checklist for keeping the Go runtime and SQL layer aligned.

---

## current reconciliation goals

### goal 1
The Go runtime should only call SQL helpers that actually exist in `sql/functions.sql`.

### goal 2
The Go runtime should only write columns that actually exist in `sql/postgres-ddl.sql` and related schema files.

### goal 3
Every retrieval wrapper should match the exact SQL return shape.

### goal 4
Every new runtime DB write path should be reflected in:
- SQL schema
- SQL contracts doc
- verify/check workflow

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

### 7. verify command
Run:
- `make verify`

### 8. build/test
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

---

## summary
The repo is strongest when the architecture, SQL, and Go runtime all describe the same system.
This file exists to make that alignment a deliberate process rather than an accident.

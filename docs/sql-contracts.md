# SQL CONTRACTS

## purpose
Document the SQL helpers and tables that the Go runtime currently expects to exist.

This file is a contract guide between:
- `sql/`
- `internal/memory/`
- `internal/harness/`
- `internal/runtime/`
- `cmd/verify/`

---

## core tables expected by runtime

### `agent_core.specialists`
Expected by runtime for:
- specialist existence
- status reads/writes
- health reads/writes

### `agent_core.slots`
Expected by runtime for:
- slot identity
- summary projection targets

### `agent_core.slot_versions`
Expected by runtime for:
- projected summary writes
- version tracking

### `agent_core.memory_records`
Expected by runtime for:
- durable memory storage
- active/staged/contradicted state

### `agent_core.memory_embeddings`
Expected by runtime for:
- record-level retrieval

### `agent_core.memory_regions`
### `agent_core.memory_region_embeddings`
### `agent_core.memory_clusters`
### `agent_core.memory_cluster_embeddings`
### `agent_core.memory_cluster_members`
Expected by runtime for:
- coarse-to-fine retrieval

### `agent_core.memory_postmortems`
Expected by runtime for:
- postmortem persistence

### `agent_core.memory_eval_cases`
Expected by runtime for:
- eval staging

### `agent_core.memory_self_edit_candidates`
Expected by runtime for:
- self-edit candidate staging

### `agent_core.routing_audit`
Expected by runtime for:
- routing audit writes

---

## core SQL functions expected by runtime

### `fn_run_coarse_to_fine_search(namespace, specialist_id, query_embedding, top_regions, top_clusters, top_records)`
Used by:
- `internal/memory/retrieval.go`
- `cmd/verify/main.go`
- runtime startup smoke testing

Expected result shape:
- `record_id`
- `record_kind`
- `title`
- `summary`
- `body`
- `status`
- `semantic_score`
- `cluster_score`
- `final_score`

### `fn_update_specialist_health_stats(specialist_id)`
Used by:
- `internal/memory/store.go`
- harness incident follow-up

Expected behavior:
- recompute and persist specialist health stats
- return the affected specialist id

### `fn_get_specialist_health_snapshot(specialist_id)`
Used by:
- `internal/memory/store.go`
- startup checks
- `cmd/verify/main.go`

Expected result shape:
- `specialist_id`
- `status`
- `health_score`

### `fn_project_and_write_memory_summary_slot(namespace, specialist_id, created_by, rationale)`
Used by:
- `internal/memory/store.go`
- harness startup checks

Expected behavior:
- compute a compact `MEMORY.md` summary from durable memory
- write a new slot version
- return the target `slot_id`

---

## write-path expectations

### postmortem creation
Runtime currently expects a direct insert path or equivalent helper for:
- namespace
- specialist_id
- task_summary
- expected_behavior
- actual_behavior
- what_went_wrong
- failure_classification
- root_cause
- preventable
- created_by

### eval staging
Runtime currently expects:
- namespace
- specialist_id
- category
- case_title
- prompt_input
- expected_behavior
- expected_output_or_criteria
- created_by

### self-edit candidate staging
Runtime currently expects:
- namespace
- specialist_id
- target_slot
- candidate_summary
- candidate_body
- rationale
- proposed_by

### routing audit
Runtime currently expects fields for:
- task_id
- routed_by
- initial_classifier
- task_summary
- task_class
- chosen_target
- confidence
- impact
- was_fallback
- fallback_reason
- was_override
- override_by
- multi_specialist_review
- notes

---

## maintenance rule
Whenever either of these changes:
- SQL function/table shape
- Go runtime DB call expectations

this file should be updated in the same change set.

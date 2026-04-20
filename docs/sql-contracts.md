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
- foreign-key anchoring for specialist-scoped persistence

### `agent_core.slots`
Expected by runtime for:
- slot identity
- summary projection targets

### `agent_core.slot_versions`
Expected by runtime for:
- projected summary writes
- version tracking

### `agent_core.slot_bundle_versions`
Expected by runtime for:
- canonical compiled slot bundle persistence
- source-map persistence
- bundle version tracking from `cmd/verify` and later slot sync paths

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

### `agent_core.adaptation_signals`
Expected by runtime for:
- durable telemetry persistence from retrieval, routing, and execution signal helpers
- SEAL review input when using the DB-backed path

### `agent_core.gap_clusters`
Expected by runtime for:
- future durable SEAL clustering state

### `agent_core.adaptation_proposals`
Expected by runtime for:
- durable SEAL proposal persistence

### `agent_core.growth_plans`
Expected by runtime for:
- durable DEN growth-plan persistence

### `agent_core.growth_experiments`
Expected by runtime for:
- later governed growth experiment lifecycle tracking

### `agent_core.lineage_nodes`
### `agent_core.lineage_edges`
Expected by runtime for:
- lineage persistence for new specialists, splits, derived branches, and related governed growth artifacts

### `agent_core.promotion_decisions`
Expected by runtime for:
- later oversight-driven promotion or rollback tracking

---

## identifier and schema expectations

### schema namespace
The canonical runtime schema is:
- `agent_core`

New migrations should either:
- set `search_path` to `agent_core, public`, or
- fully qualify `agent_core.` table names consistently

### identifier shape
For the current SEAL/DEN scaffolding, the Go runtime and SQL are aligned on:
- **text-style IDs** for signals, proposals, plans, bundle versions, lineage records, and promotion decisions

The current runtime does **not** assume DB-generated UUIDs for these paths.

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

### canonical slot bundle persistence
Runtime currently expects:
- text `id`
- `specialist_id`
- `bundle_toml`
- `source_map` as JSON
- `version_label`
- `created_by`

### telemetry signal persistence
Runtime currently expects:
- text `id`
- `specialist_id`
- `category`
- `severity`
- `surface`
- `task_class`
- `summary`
- `evidence_refs` as JSON
- `occurred_at`

### adaptation proposal persistence
Runtime currently expects:
- text `id`
- `specialist_id`
- `cluster_id` as nullable text
- `surface`
- `reason`
- `requested_by`
- `risk_level`
- `requires_harness`
- `requires_parent`
- `rollback_required`
- `status`

### growth-plan persistence
Runtime currently expects:
- text `id`
- `proposal_id`
- `specialist_id`
- `surface`
- `reason`
- `experiment_name`
- `freeze_plan` as JSON
- `rollback_plan` as JSON
- `status`

---

## verify expectations
`cmd/verify` now expects to be able to:
- compile the canonical slot bundle
- encode it to TOML
- attempt to persist the bundle version
- run retrieval smoke checks
- emit retrieval, routing, and execution signals
- run SEAL review over those signals
- attempt to persist SEAL proposals
- run DEN planning from the first proposal
- attempt to persist the resulting growth plan

The current verify path can warn and continue if the SEAL/DEN migration has not been applied yet.
That is useful for staged bring-up, but it should not be confused with a fully migrated runtime.

---

## maintenance rule
Whenever either of these changes:
- SQL function/table shape
- Go runtime DB call expectations

this file should be updated in the same change set.

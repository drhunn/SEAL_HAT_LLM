# SQL contracts

This file is the contract between:
- `sql/`
- `internal/memory/`
- `internal/runtime/`
- `internal/harness/`
- `cmd/verify/`

It documents the database objects the runtime actually expects.
When the runtime-facing schema changes, this file should change in the same patch.

## Canonical schema

The runtime expects the canonical schema to be:
- `agent_core`

New migrations should either:
- set `search_path TO agent_core, public`, or
- fully qualify `agent_core.` object names consistently

## Bootstrap contract

The Go bootstrap path applies SQL files from an explicit SQL root.

Default SQL root:
- `./sql`

Configurable embedded Postgres SQL root:
- `embedded_postgres.sql_root`

The runtime must not assume that the current working directory is always the repository root without making that assumption visible in config.

The ordered bootstrap file list lives in `internal/db/bootstrap.go` and currently expects filenames relative to the configured SQL root.

## Identifier strategy

The current runtime assumes:
- UUIDs for the original base tables that generate them in SQL
- text IDs for newer SEAL/DEN scaffold paths such as signals, proposals, growth plans, lineage nodes, and bundle versions

Do not change identifier shape casually.
If you change it, update both Go and SQL together.

## Core tables expected by runtime

### Base runtime tables
- `agent_core.specialists`
- `agent_core.slots`
- `agent_core.slot_versions`
- `agent_core.memory_records`
- `agent_core.memory_embeddings`
- `agent_core.memory_postmortems`
- `agent_core.memory_eval_cases`
- `agent_core.memory_self_edit_candidates`
- `agent_core.routing_audit`

### Coarse-to-fine retrieval tables
- `agent_core.memory_regions`
- `agent_core.memory_region_embeddings`
- `agent_core.memory_clusters`
- `agent_core.memory_cluster_embeddings`
- `agent_core.memory_cluster_members`

### Newer scaffold tables
- `agent_core.slot_bundle_versions`
- `agent_core.adaptation_signals`
- `agent_core.gap_clusters`
- `agent_core.adaptation_proposals`
- `agent_core.growth_plans`
- `agent_core.growth_experiments`
- `agent_core.lineage_nodes`
- `agent_core.lineage_edges`
- `agent_core.promotion_decisions`
- `agent_core.ability_ledgers`
- `agent_core.ability_growth_experiments`

## Runtime-facing SQL functions

### `fn_run_coarse_to_fine_search(namespace, specialist_id, query_embedding, top_regions, top_clusters, top_records)`
Used by:
- `internal/memory/retrieval.go`
- `cmd/verify`

Expected result columns:
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

Expected behavior:
- recompute and persist specialist health stats
- return the affected specialist ID

### `fn_get_specialist_health_snapshot(specialist_id)`
Used by:
- `internal/memory/store.go`
- `cmd/verify`

Expected result columns:
- `specialist_id`
- `status`
- `health_score`

### `fn_project_and_write_memory_summary_slot(namespace, specialist_id, created_by, rationale)`
Used by:
- `internal/memory/store.go`

Expected behavior:
- compute a compact `MEMORY.md` summary from durable memory
- write a new slot version
- return the target `slot_id`

## Write-path expectations

### Postmortem creation
The runtime expects to persist:
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

### Eval staging
The runtime expects to persist:
- namespace
- specialist_id
- category
- case_title
- prompt_input
- expected_behavior
- expected_output_or_criteria
- created_by

### Self-edit candidate staging
The runtime expects to persist:
- namespace
- specialist_id
- target_slot
- candidate_summary
- candidate_body
- rationale
- proposed_by

### Routing audit
The runtime expects fields for:
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

### Canonical slot bundle persistence
The runtime expects:
- text `id`
- `specialist_id`
- `bundle_toml`
- `source_map` as JSON
- `version_label`
- `created_by`

### Signal / proposal / growth persistence
The runtime expects:
- text IDs for signals, proposals, plans, lineage records, and promotion decisions
- JSON fields where the Go side writes structured refs, freeze plans, rollback plans, metrics, or metadata

## Verify contract

`cmd/verify` is expected to:
- load config
- connect to Postgres
- load specialist slots
- compile the canonical bundle
- encode the bundle to TOML
- attempt durable bundle persistence
- exercise retrieval
- emit retrieval, routing, and execution signals
- run proposal generation over those signals
- attempt growth-plan persistence
- stage an ability-growth record

Verify modes:
- `soft` warns on optional persistence-path failures and continues
- `strict` treats those same failures as fatal

CI and admission gates should use strict mode.
Soft mode is for staged local bring-up only.

## Maintenance rule

Update this file whenever any of the following change:
- runtime-facing SQL functions
- runtime-facing table shape
- identifier strategy
- verify expectations
- SQL bootstrap file ordering or SQL root behavior

If the runtime contract changed and this file did not, the patch is incomplete.

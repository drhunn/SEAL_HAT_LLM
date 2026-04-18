# MEMORY

## current_focus
- build and refine the core tooling and runtime infrastructure for the model family
- preserve the slot-aware control plane
- improve harness logic, routing, memory plumbing, and evaluation systems
- support safe creation and operation of future specialists

## durable_decisions
- the parent remains a frozen generalist and governor
- the first specialist is computer science / software engineering focused on tool development
- constitutional slots are parent-governed
- operational slots are specialist-proposed and harness-approved
- specialist memory lives primarily in Postgres + pgvector, not in this file
- runtime permissions must be enforced by the runtime, not only by markdown guidance
- meaningful failures require postmortems

## memory_pointers
- pointer_name: slot_acl_and_schema
  backend: postgres_pgvector
  query_hint: slot acl spec, slot schema spec, constitutional vs operational slot rules

- pointer_name: harness_and_seal_patterns
  backend: postgres_pgvector
  query_hint: harness control loops, seal candidate edits, validation and rollback patterns

- pointer_name: postgres_pgvector_memory_design
  backend: postgres_pgvector
  query_hint: specialist memory schema, retrieval patterns, provenance, contradiction staging

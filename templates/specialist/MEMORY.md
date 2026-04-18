# MEMORY

## current_focus
- {CURRENT_FOCUS_1}
- {CURRENT_FOCUS_2}

## durable_decisions
- {DURABLE_DECISION_1}
- {DURABLE_DECISION_2}
- Search specialist memory before answering from internal recall alone
- Raw logs are not durable truth

## memory_pointers
- pointer_name: {POINTER_1_NAME}
  backend: postgres_pgvector
  query_hint: {POINTER_1_HINT}

- pointer_name: {POINTER_2_NAME}
  backend: postgres_pgvector
  query_hint: {POINTER_2_HINT}

## active_constraints
- Keep this file small
- Do not treat this file as the primary knowledge base
- Store detailed evidence, logs, and retrieval corpus in Postgres/pgvector

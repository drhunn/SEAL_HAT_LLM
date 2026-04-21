# RUNTIME TASK INBOX CONTRACT

## purpose
Define the file-backed intake contract for the bounded runtime task path.

This is the current external intake surface for local runtime work.
It is deliberately conservative.
It is not a remote API and it is not a distributed queue.

---

## location
The runtime reads JSON task files from the configured inbox directory:

- `runtime.task_inbox_dir`

When enabled, the runtime creates and uses these sibling directories:

- `<task_inbox_dir>/processed`
- `<task_inbox_dir>/failed`

---

## task file contract
### schema version
Every task file must include:

- `version = "seal_hat_llm.runtime_task.v1"`

Files without the exact version are rejected and archived under `failed/`.

### file naming
Recommended naming:

- `<task_id>.json`

The filename is used as the fallback task id when `id` is omitted in the JSON payload.

### required JSON fields
Required:
- `version`
- `summary`

Optional:
- `id`
- `class`
- `primary_modality`
- `secondary_modalities`
- `cross_modal_grounding_required`
- `allow_text_only_fallback`
- `preferred_executor`
- `asset_refs`
- `prompt`

### example
```json
{
  "version": "seal_hat_llm.runtime_task.v1",
  "id": "demo-task-001",
  "summary": "Summarize the attached evidence",
  "class": "analysis",
  "primary_modality": "document",
  "secondary_modalities": ["text"],
  "cross_modal_grounding_required": false,
  "allow_text_only_fallback": true,
  "preferred_executor": "Document-Layout-OCR-Specialist-01",
  "asset_refs": ["sandbox://evidence/doc-1"],
  "prompt": "Extract the key findings and contradictions."
}
```

---

## processing semantics
### claim
The runtime claims a task by renaming:

- `<task>.json` -> `<task>.json.working`

This prevents the same live file from being processed twice by the same local runtime.

### success
On success, the working file is moved to:

- `processed/<task>.json`

### failure
On failure, the working file is moved to:

- `failed/<task>.json`

The runtime also writes:

- `failed/<task>.json.error.txt`

### no deletion
The runtime does not delete task files.
Archival movement is the current audit trail.

---

## result artifact contract
For every processed or failed task, the runtime writes a result artifact adjacent to the archived task file:

- `<archived_task_path>.result.json`

### schema version
Result artifacts use:

- `version = "seal_hat_llm.runtime_task_artifact.v1"`

### fields
Result artifacts include:
- `version`
- `task_id`
- `status` (`processed` or `failed`)
- `summary`
- `class`
- `primary_modality`
- `executor` (when available)
- `host_name` (when available)
- `archived_path`
- `error` (for failed tasks)
- `signal_summaries`
- `completed_at`

### success example
```json
{
  "version": "seal_hat_llm.runtime_task_artifact.v1",
  "task_id": "demo-task-001",
  "status": "processed",
  "summary": "Summarize the attached evidence",
  "class": "analysis",
  "primary_modality": "document",
  "executor": "Document-Layout-OCR-Specialist-01",
  "host_name": "local-document-host",
  "archived_path": "artifacts/task_inbox/processed/demo-task-001.json",
  "signal_summaries": [],
  "completed_at": "2026-04-20T23:59:59Z"
}
```

### failure example
```json
{
  "version": "seal_hat_llm.runtime_task_artifact.v1",
  "task_id": "demo-task-002",
  "status": "failed",
  "summary": "Broken task",
  "class": "analysis",
  "primary_modality": "text",
  "archived_path": "artifacts/task_inbox/failed/demo-task-002.json",
  "error": "validate task file demo-task-002.json: version must be \"seal_hat_llm.runtime_task.v1\"",
  "completed_at": "2026-04-20T23:59:59Z"
}
```

---

## helper CLI
Use the enqueue helper instead of hand-writing JSON when possible:

```bash
go run ./cmd/enqueue-task \
  -inbox-dir ./artifacts/task_inbox \
  -summary "Summarize the attached evidence" \
  -class analysis \
  -primary-modality document \
  -asset-refs sandbox://evidence/doc-1 \
  -prompt "Extract the key findings and contradictions."
```

The command writes one task file to the inbox and prints the path.

---

## constraints
This contract is intentionally local and conservative.
It does not yet provide:
- remote submission
- multi-worker lease renewal
- queue priorities
- dead-letter retries
- distributed coordination

That is acceptable for the current scope.

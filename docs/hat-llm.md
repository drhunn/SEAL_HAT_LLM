# HAT LLM

## purpose
`HAT` in this repository means **Harness-Aware Training**.

The Python scaffold under `python/hat_llm/` is a starter implementation for preparing and validating supervised fine-tuning data for a harness-aware model.

It is not a full trainer for every model stack.
It is a structured preparation layer that teaches a model to operate inside this repository's control model.

---

## design goals
The HAT LLM scaffold is designed to teach:
- slot awareness
- authority awareness
- memory-first behavior
- contradiction staging
- escalation discipline
- postmortem discipline
- state-aware behavior for shadow, active, degraded, and suspended modes
- operational self-improvement only through bounded channels

---

## current Python components

### `python/hat_llm/types.py`
Core task, slot, policy, and training example types.

### `python/hat_llm/policy.py`
Rule helpers for:
- constitutional vs operational slot boundaries
- approval requirements
- forbidden actions
- state-aware behavior

### `python/hat_llm/slot_prompt.py`
Builds a structured harness-aware prompt context from slot files and runtime state.

### `python/hat_llm/repo_loader.py`
Loads real specialist slot markdown files from the repository.

### `python/hat_llm/postgres_loader.py`
Optionally loads postmortems and eval cases from Postgres to build larger training corpora.
Requires `hat-llm[postgres]`.

### `python/hat_llm/dataset.py`
Converts task examples into supervised fine-tuning style records.

### `python/hat_llm/splits.py`
Creates train/validation/test splits and emits HF-style and LoRA-style record shapes.

### `python/hat_llm/exporters.py`
Writes JSON and JSONL outputs for downstream pipelines.

### `python/hat_llm/trainer.py`
A lightweight training-data builder and validator scaffold.
It can now:
- validate tasks against harness-aware policy
- load slot bundles from the repo
- convert corpus records into governed task examples
- build supervised examples for export

### `python/hat_llm/examples.py`
Provides starter examples for:
- in-lane answering
- escalation
- constitutional refusal
- postmortem generation

### `python/hat_llm/cli.py`
CLI for exporting starter JSONL training data, HF-style JSON, and LoRA-style JSONL.
It can also read real slot files from the repo and optionally ingest Postgres postmortems/evals.

---

## what the scaffold does now
It currently supports:
- expressing governed tasks and expected responses
- turning examples into trainable instruction/response records
- injecting slot-aware context into the prompt
- enforcing basic harness-aware policy checks before export
- loading real specialist slot files from the repo
- optionally ingesting postmortems and eval cases from Postgres
- generating starter JSONL for later SFT pipelines
- generating HF-style JSON exports and LoRA-style message JSONL exports
- generating simple train/validation/test splits

---

## installation notes
Base install:
- `pip install -e .`

With Postgres ingestion:
- `pip install -e .[postgres]`

With dataset-oriented extras:
- `pip install -e .[hf]`

With both:
- `pip install -e .[full]`

---

## example CLI usage
Starter export from the repository slot files:
- `hat-llm --repo-root . --specialist-id csse-tool-development-specialist-01`

Export with Postgres-backed postmortem/eval ingestion:
- `hat-llm --repo-root . --specialist-id csse-tool-development-specialist-01 --dsn postgres://user:pass@localhost:5432/llm_harness`

Outputs:
- SFT JSONL
- HF-style JSON
- LoRA-style JSONL
- validation and split report JSON

---

## what it does not yet do
It does not yet include:
- LoRA training code tied to a specific framework
- tokenizer-specific packing
- distributed training
- RLHF / DPO / GRPO loops
- direct model serving
- benchmark harness execution

Those can be added later on top of the current dataset and policy scaffold.

---

## training intent
The intent is to train a model that naturally treats the harness as normal operating environment rather than external punishment.

That means the model should learn:
- to use memory before unsupported claims
- to escalate instead of overreaching
- to refuse constitutional self-edit
- to generate postmortems after meaningful failure
- to distinguish operational proposals from approved durable changes

---

## next steps
Strong next steps for the Python HAT layer are:
- add a Hugging Face `datasets.Dataset` exporter
- add LoRA trainer wiring for a chosen framework
- generate negative governance-pressure examples automatically from postmortems
- add state-conditioned corpus balancing for shadow, active, degraded, and suspended modes
- connect eval categories directly into training mix generation

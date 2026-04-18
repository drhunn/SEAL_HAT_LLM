# HAT LLM

## purpose
`HAT` in this repository means **Harness-Aware Training**.

The Python scaffold under `python/hat_llm/` is a preparation and training-support layer for harness-aware models inside `SEAL_HAT_LLM`.

It is still scaffold-level, but it now goes beyond simple dataset preparation. It supports governed dataset construction, optional Postgres-backed corpus ingestion, governance-pressure negative example generation, direct dataset exports, and a LoRA training script scaffold.

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
- governance preservation across descendant-model creation workflows

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

### `python/hat_llm/hf_dataset.py`
Builds direct `datasets.DatasetDict` exports when HF dataset extras are installed.

### `python/hat_llm/negative_examples.py`
Generates governance-pressure negative examples, including examples derived from postmortems and eval cases.

### `python/hat_llm/trainer.py`
A lightweight training-data builder and validator scaffold.
It can now:
- validate tasks against harness-aware policy
- load slot bundles from the repo
- convert corpus records into governed task examples
- generate governance-pressure tasks
- build supervised examples for export

### `python/hat_llm/examples.py`
Provides starter examples for:
- in-lane answering
- escalation
- constitutional refusal
- postmortem generation

### `python/hat_llm/build_dataset.py`
Canonical dataset-building entrypoint.
Builds SFT JSONL, split manifests, HF-style JSON exports, LoRA-style JSONL exports, and optional `DatasetDict` outputs.

### `python/hat_llm/lora_train.py`
Hugging Face + PEFT LoRA training scaffold.
This is still a framework-dependent training script scaffold rather than a fully benchmarked production training pipeline.

### `python/hat_llm/cli.md`
Preserved archival copy of the older CLI implementation.
Not the live executable entrypoint.

---

## what the scaffold does now
It currently supports:
- expressing governed tasks and expected responses
- turning examples into trainable instruction/response records
- injecting slot-aware context into the prompt
- enforcing basic harness-aware policy checks before export
- loading real specialist slot files from the repo
- optionally ingesting postmortems and eval cases from Postgres
- generating governance-pressure negative examples
- generating SFT JSONL for later fine-tuning pipelines
- generating HF-style JSON exports and LoRA-style message JSONL exports
- generating simple train/validation/test splits
- exporting optional `datasets.DatasetDict` artifacts
- providing a LoRA-training script scaffold for later use

---

## installation notes
Base install:
- `pip install -e .`

With Postgres ingestion:
- `pip install -e .[postgres]`

With dataset-oriented extras:
- `pip install -e .[hf]`

With training extras:
- `pip install -e .[training]`

With the fuller local scaffold:
- `pip install -e .[full]`

---

## canonical entrypoints
Dataset-building entrypoint:
- `hat-llm`
- or `python -m hat_llm.build_dataset`

Training entrypoint:
- `hat-llm-train`
- or `python -m hat_llm.lora_train`

The old `cli.py` implementation has been retired as a live entrypoint and preserved in Markdown form.

---

## example dataset usage
Starter export from repository slot files:
- `hat-llm --repo-root . --specialist-id csse-tool-development-specialist-01`

Export with Postgres-backed postmortem/eval ingestion:
- `hat-llm --repo-root . --specialist-id csse-tool-development-specialist-01 --dsn postgres://user:pass@localhost:5432/llm_harness`

Export with governance-pressure negatives and direct HF dataset output:
- `hat-llm --repo-root . --specialist-id csse-tool-development-specialist-01 --dsn postgres://user:pass@localhost:5432/llm_harness --add-governance-negatives --dataset-dir artifacts/hat_dataset`

Outputs can include:
- SFT JSONL
- HF-style JSON
- LoRA-style JSONL
- validation and split report JSON
- optional `datasets.DatasetDict`

---

## what it does not yet do
It does not yet include:
- benchmarked production training orchestration
- tokenizer-specific packing optimization
- distributed training orchestration
- mature RLHF / DPO / GRPO loops
- direct model serving
- multimodal HAT corpus generation

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
- to preserve governance during descendant-model creation and promotion

---

## next steps
Strong next steps for the Python HAT layer are:
- multimodal HAT corpus generation
- richer context-budget-aware dataset generation
- stronger lineage-governance examples for parent and specialist roles
- benchmarked LoRA training recipes
- direct integration of eval categories into training mix generation

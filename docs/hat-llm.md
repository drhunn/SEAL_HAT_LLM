# HAT LLM

`HAT` in this repository means **Harness-Aware Training**.

The Python code under `python/hat_llm/` is a dataset-building and training-support layer for governed specialist runtimes.
It is useful today, but it is still a scaffold.

## What it is for

The HAT layer exists to build training data and training scaffolding that teach a model to operate with:
- slot awareness
- authority awareness
- memory-first behavior
- escalation discipline
- postmortem discipline
- state-aware behavior across bounded runtime modes

The intent is straightforward:
make the harness feel like the model's normal operating environment instead of a bolt-on punishment layer.

## What works now

The current Python scaffold already supports:
- loading real specialist slot markdown from the repo
- building governed task examples
- validating tasks against harness-aware policy helpers
- generating governance-pressure negative examples
- optionally ingesting postmortems and eval cases from Postgres
- exporting JSON / JSONL training artifacts
- optionally exporting a `datasets.DatasetDict`
- providing a JAX training scaffold for later refinement

## Current runtime path

The current authoring and training path is:
1. markdown slot files in the repo
2. structured slot loading / compilation
3. dataset building for supervised examples
4. optional JAX training scaffold

The broader hybrid slot story is still a target shape, not a completed production path.

## Main Python components

### Dataset and policy layer
- `python/hat_llm/policy.py` — harness-aware rule helpers
- `python/hat_llm/dataset.py` — training-example construction
- `python/hat_llm/trainer.py` — small orchestration layer for validation and example generation
- `python/hat_llm/negative_examples.py` — governance-pressure negative example generation
- `python/hat_llm/repo_loader.py` — specialist slot loading from the repo
- `python/hat_llm/postgres_loader.py` — optional Postgres corpus loading

### Export layer
- `python/hat_llm/exporters.py` — JSON / JSONL export helpers
- `python/hat_llm/hf_dataset.py` — optional `DatasetDict` export
- `python/hat_llm/splits.py` — train / validation / test splitting

### Training scaffold
- `python/hat_llm/build_dataset.py` — canonical dataset-building entrypoint
- `python/hat_llm/jax_train.py` — JAX training scaffold
- `python/hat_llm_jax/` — early slot-aware JAX model scaffolding

## Installation

Base install:
- `pip install -e .`

Optional extras:
- Postgres ingestion: `pip install -e .[postgres]`
- dataset-oriented extras: `pip install -e .[hf]`
- training extras: `pip install -e .[training]`
- full local scaffold: `pip install -e .[full]`

## Entry points

Dataset build:
- `hat-llm`
- `python -m hat_llm.build_dataset`

Training scaffold:
- `hat-llm-train`
- `python -m hat_llm.jax_train`

## Example usage

Build from real repo slots:
- `hat-llm --repo-root . --specialist-id csse-tool-development-specialist-01`

Build with Postgres-backed postmortems and eval cases:
- `hat-llm --repo-root . --specialist-id csse-tool-development-specialist-01 --dsn postgres://user:pass@localhost:5432/llm_harness`

Build with governance-pressure negatives and `DatasetDict` output:
- `hat-llm --repo-root . --specialist-id csse-tool-development-specialist-01 --dsn postgres://user:pass@localhost:5432/llm_harness --add-governance-negatives --dataset-dir artifacts/hat_dataset`

Use starter toy data only when you explicitly mean to:
- `hat-llm --repo-root . --specialist-id missing-specialist --allow-starter-fallback`

Without `--allow-starter-fallback`, the dataset builder now fails fast when the requested repo slot pack is missing or empty.

## What it does not do yet

The current scaffold does **not** yet provide:
- benchmarked production training orchestration
- distributed training workflows
- production model serving
- mature RLHF / DPO / GRPO loops
- full multimodal HAT corpus generation
- a production slot-packet training path
- DEN-style dynamic slot growth in the training stack

## Maintenance rule

Keep this document aligned with the actual Python entrypoints and outputs.
Do not document speculative training machinery as if it already exists.

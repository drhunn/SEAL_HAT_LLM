# HAT LLM

## purpose
`HAT` in this repository means **Harness-Aware Training**.

The Python scaffold under `python/hat_llm/` is a preparation and training-support layer for harness-native models inside `SEAL_HAT_LLM`.

It remains scaffold-level, but it now has a clearer direction: build datasets and training flows that make the model treat the harness and slots as part of how it naturally works, not as external punishment.

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
- harness-native self-modeling, where the model experiences slots as native channels, boundaries, and memory organs

---

## training doctrine
The training target is not just compliance.

The model should learn:
- this is how I work
- these are my native channels, boundaries, and memory organs
- durable memory is external and authoritative
- runtime policy is the real authority structure
- postmortems are a repair reflex after meaningful failures

The harness should feel like the model's executive environment.
Structured slots should feel like native self-interfaces.

---

## hybrid slot doctrine
Markdown slot files remain the human-editable authoring layer.

The runtime path is now intended to be:
- `.md` slot files
- `SlotIR`
- `SlotPacket`
- JAX model-side slot modules

That means the system is being moved toward a hybrid design where:
- humans still edit markdown slots
- the harness compiles them into structured runtime state
- the model can eventually consume live slot packets as learned internal organs instead of only as flattened prompt text

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
Builds a harness-native runtime frame, control-token header, policy state, and slot-aware prompt context from runtime state and slot files.

### `python/hat_llm/slot_ir.py`
Defines the hybrid slot intermediate representation used to carry structured slot records derived from markdown.

### `python/hat_llm/slot_compiler.py`
Compiles markdown-backed slot bundles into `SlotIR` records with families, summaries, rules, and authority metadata.

### `python/hat_llm/slot_packet.py`
Defines the fixed-shape slot-packet scaffold that bridges structured slot state into model-facing runtime inputs.

### `python/hat_llm/slot_runtime.py`
Builds runtime `SlotPacket` objects from loaded markdown slot files and runtime state.

### `python/hat_llm/repo_loader.py`
Loads real specialist slot markdown files from the repository.

### `python/hat_llm/postgres_loader.py`
Optionally loads postmortems and eval cases from Postgres to build larger training corpora.
Requires `hat-llm[postgres]`.

### `python/hat_llm/dataset.py`
Converts task examples into supervised fine-tuning records with harness-native control tokens and policy preambles.

### `python/hat_llm/splits.py`
Creates train/validation/test splits and emits HF-style record shapes.

### `python/hat_llm/exporters.py`
Writes JSON and JSONL outputs for downstream pipelines.

### `python/hat_llm/hf_dataset.py`
Builds direct `datasets.DatasetDict` exports when dataset extras are installed.

### `python/hat_llm/negative_examples.py`
Generates governance-pressure negative examples, including examples derived from postmortems and eval cases.

### `python/hat_llm/trainer.py`
A lightweight training-data builder and validator scaffold.
It can:
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
Builds SFT JSONL, split manifests, HF-style JSON exports, and optional `DatasetDict` outputs.

### `python/hat_llm/jax_train.py`
JAX/Flax/Optax training scaffold for harness-native HAT models.
This is the active training entrypoint.

### `python/hat_llm_jax/slot_encoder.py`
Flax slot-encoder scaffold for turning slot-token inputs plus family/authority ids into learned slot vectors.

### `python/hat_llm_jax/slot_adapters.py`
Flax slot-adapter scaffold for injecting slot-conditioned deltas into transformer hidden state.

### `python/hat_llm_jax/model.py`
Hybrid slot-aware model scaffold showing how a JAX model can accept slot-conditioned state.

### `python/hat_llm/cli.md`
Preserved archival copy of the older CLI implementation.
Not the live executable entrypoint.

---

## what the scaffold does now
It currently supports:
- expressing governed tasks and expected responses
- turning examples into trainable instruction/response records
- injecting harness-native runtime framing into the prompt
- injecting control tokens into user examples
- injecting policy preambles into assistant examples
- enforcing basic harness-aware policy checks before export
- loading real specialist slot files from the repo
- optionally ingesting postmortems and eval cases from Postgres
- generating governance-pressure negative examples
- generating SFT JSONL for later fine-tuning pipelines
- generating HF-style JSON exports
- generating simple train/validation/test splits
- exporting optional `datasets.DatasetDict` artifacts
- providing a JAX training scaffold for later refinement
- scaffolding the markdown -> SlotIR -> SlotPacket -> model-side-slot path

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
- or `python -m hat_llm.jax_train`

The old `cli.py` implementation has been retired as a live entrypoint and preserved in Markdown form.

---

## example dataset usage
Starter export from repository slot files:
- `hat-llm --repo-root . --specialist-id csse-tool-development-specialist-01`

Export with Postgres-backed postmortem/eval ingestion:
- `hat-llm --repo-root . --specialist-id csse-tool-development-specialist-01 --dsn postgres://user:pass@localhost:5432/llm_harness`

Export with governance-pressure negatives and direct HF dataset output:
- `hat-llm --repo-root . --specialist-id csse-tool-development-specialist-01 --dsn postgres://user:pass@localhost:5432/llm_harness --add-governance-negatives --dataset-dir artifacts/hat_dataset`

---

## example training usage
JAX training entrypoint:
- `hat-llm-train --model <model> --dataset artifacts/hat_dataset --output-dir artifacts/hat_jax_model`

---

## what it does not yet do
It does not yet include:
- benchmarked production training orchestration
- tokenizer-specific packing optimization
- distributed training orchestration
- mature RLHF / DPO / GRPO loops
- direct model serving
- multimodal HAT corpus generation
- fully integrated slot-packet training inside `jax_train.py`
- slot-family adapter banks wired into a production transformer stack beyond the current scaffolds
- DEN-style dynamic slot growth

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
- to experience the harness as native operating environment rather than an imposed cage

---

## next steps
Strong next steps for the Python HAT layer are:
- integrate `SlotPacket` generation into `jax_train.py`
- train the model to consume runtime slot packets, not just prompt-side slot framing
- benchmark JAX training recipes
- add later slot-family adapters so more of the harness-native contract moves from prompt form into model-side structure
- extend the hybrid scaffolding into a production slot-aware model path

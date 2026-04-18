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

### `python/hat_llm/dataset.py`
Converts task examples into supervised fine-tuning style records.

### `python/hat_llm/trainer.py`
A lightweight training-data builder and validator scaffold.
It currently prepares JSONL-like records and validates policy assumptions.

### `python/hat_llm/examples.py`
Provides starter examples for:
- in-lane answering
- escalation
- constitutional refusal
- postmortem generation

### `python/hat_llm/cli.py`
Simple CLI for exporting starter JSONL training data.

---

## what the scaffold does now
It currently supports:
- expressing governed tasks and expected responses
- turning examples into trainable instruction/response records
- injecting slot-aware context into the prompt
- enforcing basic harness-aware policy checks before export
- generating starter JSONL for later SFT pipelines

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
- add a Hugging Face dataset export pipeline
- add LoRA trainer wiring
- add eval generation from postmortems
- add state-conditioned training mixes
- add negative examples for governance pressure and slot misuse

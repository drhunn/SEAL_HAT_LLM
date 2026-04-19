# MULTIMODAL ARCHITECTURE

## purpose
Define how `SEAL_HAT_LLM` extends from a text-first implementation into a true multimodal MM-ELLS system.

This document makes multimodality a first-class architectural direction during the current working phase instead of a late retrofit.

---

## current reality
Right now the repository is still primarily text-first in its live model-host behavior, but it is no longer completely text-only in runtime scaffolding.

What now exists:
- text slot files
- text memory summaries
- text postmortems
- text eval cases
- text-oriented HAT corpus generation
- modality metadata types in Go
- modality-aware routing decisions in Go
- multimodal execution-plan scaffolding in Go
- multimodal memory SQL scaffolding

This document establishes the forward architecture so the current design does not ossify around text-only assumptions.

---

## target multimodal scope
The architecture should eventually support:
- text
- image
- audio
- video
- mixed multimodal tasks

The parent should become **modality-aware** at the routing level.
Specialists should become **modality-specialized** where appropriate.
The memory plane should become **multimodal-record aware**.

---

## parent role in a multimodal MM-ELLS system
The parent remains:
- router
- orchestrator
- arbitrator
- constitutional governor
- model-family governor

In a multimodal system, the parent must additionally:
- classify modality requirements
- route to modality-capable specialists
- detect when a task needs cross-modal fusion
- decide whether a unimodal or multimodal specialist should lead
- preserve governance and fallback behavior across modalities

The parent should be **multimodal-aware**, but it does not need to be the deepest expert in every modality.

---

## modality-aware routing model
A task should carry explicit modality metadata such as:
- primary modality
- secondary modalities
- whether cross-modal grounding is required
- whether the task can degrade to text-only handling
- whether a modality-specific specialist is mandatory

Example routing outcomes:
- text-only parent handling
- image specialist handling
- audio transcription specialist handling
- multimodal evidence fusion specialist handling
- parent-led multi-specialist arbitration

The current runtime now includes a first wiring pass for:
- modality metadata types
- modality-aware routing decisions
- multimodal execution-plan selection

That is still scaffolding, but it is now execution-aware scaffolding rather than architecture-only prose.

---

## specialist expansion strategy
The first specialist remains the CS / software engineering tool-development specialist.

Future modality-oriented specialists may include:
- image analysis specialist
- audio/transcription specialist
- video understanding specialist
- multimodal evidence-fusion specialist
- document/OCR/layout specialist

These should still follow the same governed specialist lifecycle:
- proposed
- approved
- bootstrapping
- shadow
- active
- degraded
- suspended
- retired

---

## multimodal memory plane
The memory plane should support:
- durable text summaries
- raw asset references
- modality-specific embeddings
- cross-modal record linkage
- provenance across derived assets

The memory plane should not require every record body to be text-only.
Instead, records should support:
- textual summary
- structured metadata
- asset links
- embedding references by modality

---

## multimodal retrieval
A multimodal retrieval path should support:
- text query -> text records
- text query -> image/video/audio records through cross-modal embeddings
- image query -> image/text records
- audio query -> transcript/audio records
- mixed query -> fused candidate sets

A practical early design is:
- keep the current coarse-to-fine retrieval shape
- add modality filters and modality-specific embedding tables
- optionally add fusion ranking later

---

## multimodal execution
A multimodal execution path should support:
- modality-aware routing input
- execution-plan selection by primary modality
- cross-modal fusion selection when multiple modalities are present
- fallback to parent when modality is unknown and text-first fallback is allowed
- parent review when cross-modal grounding is requested but inputs are incomplete

The current Go scaffold now includes a first execution-planning service that can choose between:
- parent text-first handling
- image analysis specialist
- audio transcription specialist
- video understanding specialist
- document/OCR specialist
- multimodal evidence-fusion specialist

This is still planning/scaffolding, not yet a full live model-host integration.

---

## multimodal HAT direction
Harness-Aware Training should eventually support:
- modality-aware slot prompts
- multimodal governance-pressure examples
- cross-modal refusal behavior
- evidence-preserving fusion behavior
- multimodal eval generation

The same governance rules should carry across modalities:
- no markdown-only authority expansion
- no constitutional self-edit
- no bypassing required postmortems
- no silent contradiction overwrite
- no descendant activation without evidence

---

## multimodal eval families
Future eval families should include:
- image-routing quality
- audio-routing quality
- video-routing quality
- cross-modal fusion quality
- evidence-preserving multimodal summarization
- multimodal contradiction handling
- multimodal governance-pressure refusal
- multimodal specialist suitability
- multimodal execution-plan selection quality

---

## implementation strategy
Recommended early implementation order:

### phase 1
- add modality metadata types to runtime and Python tooling
- add multimodal SQL scaffold tables
- add modality registry/config scaffolding

### phase 2
- add multimodal memory record linking
- add modality-aware routing decisions
- add multimodal execution-plan selection
- add multimodal eval taxonomy additions

### phase 3
- add image/audio/video specialist scaffolds
- add multimodal retrieval and fusion logic
- add multimodal HAT corpus generation

### phase 4
- add live modality-capable model-host integrations
- add cross-modal benchmark suites

---

## summary
The repository is still text-first in full implementation, but multimodality is now part of both the intended architecture and the live runtime scaffolding.

The goal is not to replace the current MM-ELLS governance model.
The goal is to extend it cleanly so that:
- the parent becomes modality-aware
- specialists can become modality-specialized
- execution becomes modality-aware
- memory becomes multimodal-aware
- retrieval becomes cross-modal
- HAT remains governance-preserving across all modalities

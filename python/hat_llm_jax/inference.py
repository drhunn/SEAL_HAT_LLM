from __future__ import annotations

import argparse
import json
from pathlib import Path

import jax
import jax.numpy as jnp
from flax import serialization
from transformers import AutoTokenizer

from hat_llm.repo_loader import RepositoryLoader
from hat_llm.slot_packet import SlotPacketConfig
from hat_llm.slot_runtime import build_slot_packet
from hat_llm.types import RuntimeState
from hat_llm_jax.model import HybridRuntimeSlotModel


def _load_config(model_dir: str) -> dict[str, object]:
    return json.loads((Path(model_dir) / "hybrid_model_config.json").read_text(encoding="utf-8"))


def _build_slot_inputs(model_dir: str, repo_root: str, specialist_id: str, runtime_mode: str, slot_max_count: int, slot_max_length: int):
    tokenizer = AutoTokenizer.from_pretrained(model_dir, use_fast=True)
    if tokenizer.pad_token is None:
        tokenizer.pad_token = tokenizer.eos_token

    slot_loader = RepositoryLoader(repo_root)
    slots = slot_loader.load_specialist_slots(specialist_id)
    runtime = RuntimeState(
        specialist_id=specialist_id,
        namespace=f"memory.{specialist_id}",
        mode=runtime_mode,
    )
    packet = build_slot_packet(
        runtime,
        slots,
        SlotPacketConfig(max_slots=slot_max_count, max_chars_per_slot=512),
    )
    encoded = tokenizer(
        packet.slot_texts,
        truncation=True,
        max_length=slot_max_length,
        padding="max_length",
    )
    return tokenizer, {
        "slot_token_ids": jnp.asarray(encoded["input_ids"], dtype=jnp.int32),
        "slot_token_mask": jnp.asarray(encoded["attention_mask"], dtype=jnp.int32),
        "slot_family_ids": jnp.asarray(packet.slot_family_ids, dtype=jnp.int32),
        "slot_authority_ids": jnp.asarray(packet.slot_authority_ids, dtype=jnp.int32),
        "slot_enabled_mask": jnp.asarray(packet.slot_enabled_mask, dtype=jnp.float32),
    }


def run_inference(model_dir: str, repo_root: str, specialist_id: str, runtime_mode: str, prompt: str) -> str:
    cfg = _load_config(model_dir)
    tokenizer, slot_inputs = _build_slot_inputs(
        model_dir=model_dir,
        repo_root=repo_root,
        specialist_id=specialist_id,
        runtime_mode=runtime_mode,
        slot_max_count=int(cfg.get("slot_max_count", 16)),
        slot_max_length=int(cfg.get("slot_max_length", 128)),
    )
    model = HybridRuntimeSlotModel(
        vocab_size=int(cfg["vocab_size"]),
        d_model=int(cfg["d_model"]),
        d_slot=int(cfg["d_slot"]),
        num_layers=int(cfg["num_layers"]),
    )
    input_ids = tokenizer(prompt, truncation=True, max_length=int(cfg.get("max_length", 2048)), return_tensors="np")["input_ids"]
    init_vars = model.init(
        jax.random.PRNGKey(0),
        input_ids=jnp.asarray(input_ids, dtype=jnp.int32),
        slot_token_ids=slot_inputs["slot_token_ids"],
        slot_token_mask=slot_inputs["slot_token_mask"],
        slot_family_ids=slot_inputs["slot_family_ids"],
        slot_authority_ids=slot_inputs["slot_authority_ids"],
        slot_enabled_mask=slot_inputs["slot_enabled_mask"],
    )
    params = serialization.from_bytes(init_vars["params"], (Path(model_dir) / "hybrid_model.msgpack").read_bytes())
    logits = model.apply(
        {"params": params},
        input_ids=jnp.asarray(input_ids, dtype=jnp.int32),
        slot_token_ids=slot_inputs["slot_token_ids"],
        slot_token_mask=slot_inputs["slot_token_mask"],
        slot_family_ids=slot_inputs["slot_family_ids"],
        slot_authority_ids=slot_inputs["slot_authority_ids"],
        slot_enabled_mask=slot_inputs["slot_enabled_mask"],
    )
    next_token = int(jnp.argmax(logits[0, -1]).item())
    return tokenizer.decode([next_token])


def main() -> None:
    parser = argparse.ArgumentParser(description="Run hybrid runtime-slot inference with a saved JAX model")
    parser.add_argument("--model-dir", required=True)
    parser.add_argument("--repo-root", default=".")
    parser.add_argument("--specialist-id", default="csse-tool-development-specialist-01")
    parser.add_argument("--runtime-mode", default="active")
    parser.add_argument("--prompt", required=True)
    args = parser.parse_args()
    print(run_inference(args.model_dir, args.repo_root, args.specialist_id, args.runtime_mode, args.prompt))


if __name__ == "__main__":
    main()

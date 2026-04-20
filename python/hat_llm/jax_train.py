from __future__ import annotations

import argparse
import math
from dataclasses import dataclass
from pathlib import Path

import jax
import jax.numpy as jnp
import optax
from flax.training import train_state

from .repo_loader import RepositoryLoader
from .slot_packet import SlotPacketConfig
from .slot_runtime import build_slot_packet
from .types import RuntimeState
from hat_llm_jax.model import HybridRuntimeSlotModel


@dataclass(slots=True)
class TrainConfig:
    model_name_or_path: str
    dataset_path: str
    output_dir: str
    max_length: int = 2048
    learning_rate: float = 1e-4
    weight_decay: float = 0.01
    batch_size: int = 1
    num_train_epochs: float = 1.0
    warmup_steps: int = 50
    logging_steps: int = 10
    save_steps: int = 250
    hybrid_slot_model: bool = False
    repo_root: str = "."
    specialist_id: str = "csse-tool-development-specialist-01"
    runtime_mode: str = "active"
    d_model: int = 256
    d_slot: int = 128
    num_layers: int = 4
    slot_max_length: int = 128
    slot_max_count: int = 16


def build_text(example: dict[str, object]) -> str:
    system = example.get("system", "")
    user = example.get("user", "")
    assistant = example.get("assistant", "")
    return f"<|system|>\n{system}\n<|user|>\n{user}\n<|assistant|>\n{assistant}"


def _load_dataset(dataset_path: Path):
    from datasets import DatasetDict, load_dataset, load_from_disk

    if dataset_path.is_dir():
        dataset_obj = load_from_disk(str(dataset_path))
        if isinstance(dataset_obj, DatasetDict):
            return dataset_obj["train"]
        return dataset_obj
    return load_dataset("json", data_files=str(dataset_path), split="train")


def _tokenize_dataset(dataset_path: str, model_name_or_path: str, max_length: int):
    from transformers import AutoTokenizer

    dataset = _load_dataset(Path(dataset_path))
    tokenizer = AutoTokenizer.from_pretrained(model_name_or_path, use_fast=True)
    if tokenizer.pad_token is None:
        tokenizer.pad_token = tokenizer.eos_token

    def tokenize_row(row: dict[str, object]) -> dict[str, object]:
        text = build_text(row)
        encoded = tokenizer(
            text,
            truncation=True,
            max_length=max_length,
            padding="max_length",
        )
        encoded["labels"] = encoded["input_ids"].copy()
        return encoded

    tokenized = dataset.map(tokenize_row, remove_columns=dataset.column_names)
    tokenized.set_format(type="numpy", columns=["input_ids", "attention_mask", "labels"])
    return tokenized, tokenizer


def _build_slot_inputs(cfg: TrainConfig, tokenizer):
    slot_loader = RepositoryLoader(cfg.repo_root)
    slots = slot_loader.load_specialist_slots(cfg.specialist_id)
    runtime = RuntimeState(
        specialist_id=cfg.specialist_id,
        namespace=f"memory.{cfg.specialist_id}",
        mode=cfg.runtime_mode,
    )
    packet = build_slot_packet(
        runtime,
        slots,
        SlotPacketConfig(max_slots=cfg.slot_max_count, max_chars_per_slot=512),
    )
    encoded = tokenizer(
        packet.slot_texts,
        truncation=True,
        max_length=cfg.slot_max_length,
        padding="max_length",
    )
    return {
        "slot_token_ids": jnp.asarray(encoded["input_ids"], dtype=jnp.int32),
        "slot_token_mask": jnp.asarray(encoded["attention_mask"], dtype=jnp.int32),
        "slot_family_ids": jnp.asarray(packet.slot_family_ids, dtype=jnp.int32),
        "slot_authority_ids": jnp.asarray(packet.slot_authority_ids, dtype=jnp.int32),
        "slot_enabled_mask": jnp.asarray(packet.slot_enabled_mask, dtype=jnp.float32),
    }


def _iterate_batches(dataset, batch_size: int):
    total = len(dataset)
    for start in range(0, total, batch_size):
        stop = min(start + batch_size, total)
        batch = dataset[start:stop]
        yield {
            "input_ids": jnp.asarray(batch["input_ids"], dtype=jnp.int32),
            "attention_mask": jnp.asarray(batch["attention_mask"], dtype=jnp.int32),
            "labels": jnp.asarray(batch["labels"], dtype=jnp.int32),
        }


def _create_state(params, apply_fn, learning_rate: float, weight_decay: float, warmup_steps: int, total_steps: int):
    warmup = optax.linear_schedule(init_value=0.0, end_value=learning_rate, transition_steps=max(warmup_steps, 1))
    decay = optax.cosine_decay_schedule(init_value=learning_rate, decay_steps=max(total_steps, 1))
    schedule = optax.join_schedules([warmup, decay], boundaries=[max(warmup_steps, 1)])
    tx = optax.adamw(learning_rate=schedule, weight_decay=weight_decay)
    return train_state.TrainState.create(apply_fn=apply_fn, params=params, tx=tx)


def _loss_from_logits(logits, labels, attention_mask):
    shift_logits = logits[:, :-1, :]
    shift_labels = labels[:, 1:]
    shift_mask = attention_mask[:, 1:]
    losses = optax.softmax_cross_entropy_with_integer_labels(shift_logits, shift_labels)
    masked = losses * shift_mask
    denom = jnp.maximum(jnp.sum(shift_mask), 1)
    return jnp.sum(masked) / denom


def run_training(cfg: TrainConfig) -> None:
    from transformers import FlaxAutoModelForCausalLM

    tokenized, tokenizer = _tokenize_dataset(cfg.dataset_path, cfg.model_name_or_path, cfg.max_length)
    slot_inputs = _build_slot_inputs(cfg, tokenizer)

    if cfg.hybrid_slot_model:
        module = HybridRuntimeSlotModel(
            vocab_size=tokenizer.vocab_size,
            d_model=cfg.d_model,
            d_slot=cfg.d_slot,
            num_layers=cfg.num_layers,
        )
        init_vars = module.init(
            jax.random.PRNGKey(0),
            input_ids=jnp.zeros((1, cfg.max_length), dtype=jnp.int32),
            slot_token_ids=slot_inputs["slot_token_ids"],
            slot_token_mask=slot_inputs["slot_token_mask"],
            slot_family_ids=slot_inputs["slot_family_ids"],
            slot_authority_ids=slot_inputs["slot_authority_ids"],
            slot_enabled_mask=slot_inputs["slot_enabled_mask"],
        )
        params = init_vars["params"]

        def apply_fn(params, batch, slot_inputs):
            return module.apply(
                {"params": params},
                input_ids=batch["input_ids"],
                slot_token_ids=slot_inputs["slot_token_ids"],
                slot_token_mask=slot_inputs["slot_token_mask"],
                slot_family_ids=slot_inputs["slot_family_ids"],
                slot_authority_ids=slot_inputs["slot_authority_ids"],
                slot_enabled_mask=slot_inputs["slot_enabled_mask"],
            )

        save_fn = lambda params: None
    else:
        model = FlaxAutoModelForCausalLM.from_pretrained(cfg.model_name_or_path, dtype=jnp.float32)
        params = model.params

        def apply_fn(params, batch, slot_inputs):
            outputs = model(
                input_ids=batch["input_ids"],
                attention_mask=batch["attention_mask"],
                params=params,
                train=True,
            )
            return outputs.logits if hasattr(outputs, "logits") else outputs[0]

        save_fn = lambda params: model.save_pretrained(cfg.output_dir, params=params)

    steps_per_epoch = max(math.ceil(len(tokenized) / cfg.batch_size), 1)
    total_steps = max(int(steps_per_epoch * cfg.num_train_epochs), 1)
    state = _create_state(params, apply_fn, cfg.learning_rate, cfg.weight_decay, cfg.warmup_steps, total_steps)

    @jax.jit
    def train_step(state, batch, slot_inputs):
        def loss_fn(params):
            logits = state.apply_fn(params, batch, slot_inputs)
            return _loss_from_logits(logits, batch["labels"], batch["attention_mask"])

        loss, grads = jax.value_and_grad(loss_fn)(state.params)
        state = state.apply_gradients(grads=grads)
        return state, loss

    global_step = 0
    for _epoch in range(int(math.ceil(cfg.num_train_epochs))):
        for batch in _iterate_batches(tokenized, cfg.batch_size):
            state, loss = train_step(state, batch, slot_inputs)
            global_step += 1
            if global_step % cfg.logging_steps == 0:
                print(f"step={global_step} loss={float(loss):.4f}")
            if global_step % cfg.save_steps == 0 and not cfg.hybrid_slot_model:
                save_fn(state.params)
                tokenizer.save_pretrained(cfg.output_dir)

    if not cfg.hybrid_slot_model:
        save_fn(state.params)
        tokenizer.save_pretrained(cfg.output_dir)


def main() -> None:
    parser = argparse.ArgumentParser(description="Train a harness-native HAT model with JAX, Flax, and Optax")
    parser.add_argument("--model", required=True, help="base model name or path")
    parser.add_argument("--dataset", required=True, help="HF dataset directory or JSON/JSONL file")
    parser.add_argument("--output-dir", required=True, help="training output directory")
    parser.add_argument("--max-length", type=int, default=2048)
    parser.add_argument("--learning-rate", type=float, default=1e-4)
    parser.add_argument("--weight-decay", type=float, default=0.01)
    parser.add_argument("--batch-size", type=int, default=1)
    parser.add_argument("--epochs", type=float, default=1.0)
    parser.add_argument("--warmup-steps", type=int, default=50)
    parser.add_argument("--logging-steps", type=int, default=10)
    parser.add_argument("--save-steps", type=int, default=250)
    parser.add_argument("--hybrid-slot-model", action="store_true", help="train the scaffolded hybrid slot-aware model")
    parser.add_argument("--repo-root", default=".", help="repository root for slot loading")
    parser.add_argument("--specialist-id", default="csse-tool-development-specialist-01", help="specialist id for runtime slot loading")
    parser.add_argument("--runtime-mode", default="active", help="runtime mode for compiled slot packets")
    parser.add_argument("--d-model", type=int, default=256)
    parser.add_argument("--d-slot", type=int, default=128)
    parser.add_argument("--num-layers", type=int, default=4)
    parser.add_argument("--slot-max-length", type=int, default=128)
    parser.add_argument("--slot-max-count", type=int, default=16)
    args = parser.parse_args()

    cfg = TrainConfig(
        model_name_or_path=args.model,
        dataset_path=args.dataset,
        output_dir=args.output_dir,
        max_length=args.max_length,
        learning_rate=args.learning_rate,
        weight_decay=args.weight_decay,
        batch_size=args.batch_size,
        num_train_epochs=args.epochs,
        warmup_steps=args.warmup_steps,
        logging_steps=args.logging_steps,
        save_steps=args.save_steps,
        hybrid_slot_model=args.hybrid_slot_model,
        repo_root=args.repo_root,
        specialist_id=args.specialist_id,
        runtime_mode=args.runtime_mode,
        d_model=args.d_model,
        d_slot=args.d_slot,
        num_layers=args.num_layers,
        slot_max_length=args.slot_max_length,
        slot_max_count=args.slot_max_count,
    )
    run_training(cfg)


if __name__ == "__main__":
    main()

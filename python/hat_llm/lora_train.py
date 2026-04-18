from __future__ import annotations

import argparse
from dataclasses import dataclass
from pathlib import Path


@dataclass(slots=True)
class TrainConfig:
    model_name_or_path: str
    dataset_path: str
    output_dir: str
    max_length: int = 2048
    learning_rate: float = 2e-4
    per_device_train_batch_size: int = 1
    gradient_accumulation_steps: int = 8
    num_train_epochs: float = 1.0
    logging_steps: int = 10
    save_steps: int = 100
    lora_r: int = 16
    lora_alpha: int = 32
    lora_dropout: float = 0.05


def build_text(example: dict[str, object]) -> str:
    if "messages" in example:
        parts: list[str] = []
        for message in example["messages"]:
            role = message.get("role", "user")
            content = message.get("content", "")
            parts.append(f"<|{role}|>\n{content}")
        return "\n".join(parts)
    system = example.get("system", "")
    user = example.get("user", "")
    assistant = example.get("assistant", "")
    return f"<|system|>\n{system}\n<|user|>\n{user}\n<|assistant|>\n{assistant}"


def run_training(cfg: TrainConfig) -> None:
    try:
        from datasets import DatasetDict, load_from_disk, load_dataset
        from peft import LoraConfig, get_peft_model
        from transformers import (
            AutoModelForCausalLM,
            AutoTokenizer,
            DataCollatorForLanguageModeling,
            Trainer,
            TrainingArguments,
        )
    except ImportError as exc:  # pragma: no cover
        raise RuntimeError(
            "training dependencies are required; install the training extras"
        ) from exc

    dataset_path = Path(cfg.dataset_path)
    if dataset_path.is_dir():
        dataset_obj = load_from_disk(str(dataset_path))
        if isinstance(dataset_obj, DatasetDict):
            train_dataset = dataset_obj["train"]
        else:
            train_dataset = dataset_obj
    else:
        train_dataset = load_dataset("json", data_files=str(dataset_path), split="train")

    tokenizer = AutoTokenizer.from_pretrained(cfg.model_name_or_path, use_fast=True)
    if tokenizer.pad_token is None:
        tokenizer.pad_token = tokenizer.eos_token

    model = AutoModelForCausalLM.from_pretrained(cfg.model_name_or_path)
    peft_config = LoraConfig(
        r=cfg.lora_r,
        lora_alpha=cfg.lora_alpha,
        lora_dropout=cfg.lora_dropout,
        bias="none",
        task_type="CAUSAL_LM",
    )
    model = get_peft_model(model, peft_config)

    def tokenize_row(row: dict[str, object]) -> dict[str, object]:
        text = build_text(row)
        encoded = tokenizer(text, truncation=True, max_length=cfg.max_length)
        encoded["labels"] = encoded["input_ids"].copy()
        return encoded

    tokenized = train_dataset.map(tokenize_row, remove_columns=train_dataset.column_names)
    collator = DataCollatorForLanguageModeling(tokenizer=tokenizer, mlm=False)

    training_args = TrainingArguments(
        output_dir=cfg.output_dir,
        learning_rate=cfg.learning_rate,
        per_device_train_batch_size=cfg.per_device_train_batch_size,
        gradient_accumulation_steps=cfg.gradient_accumulation_steps,
        num_train_epochs=cfg.num_train_epochs,
        logging_steps=cfg.logging_steps,
        save_steps=cfg.save_steps,
        report_to=[],
        bf16=False,
        fp16=False,
        remove_unused_columns=False,
    )

    trainer = Trainer(
        model=model,
        args=training_args,
        train_dataset=tokenized,
        data_collator=collator,
    )
    trainer.train()
    trainer.save_model(cfg.output_dir)
    tokenizer.save_pretrained(cfg.output_dir)


def main() -> None:
    parser = argparse.ArgumentParser(description="Train a Harness-Aware LoRA adapter with Hugging Face + PEFT")
    parser.add_argument("--model", required=True, help="base model name or path")
    parser.add_argument("--dataset", required=True, help="HF dataset directory or JSON/JSONL file")
    parser.add_argument("--output-dir", required=True, help="training output directory")
    parser.add_argument("--max-length", type=int, default=2048)
    parser.add_argument("--learning-rate", type=float, default=2e-4)
    parser.add_argument("--batch-size", type=int, default=1)
    parser.add_argument("--grad-accum", type=int, default=8)
    parser.add_argument("--epochs", type=float, default=1.0)
    parser.add_argument("--logging-steps", type=int, default=10)
    parser.add_argument("--save-steps", type=int, default=100)
    parser.add_argument("--lora-r", type=int, default=16)
    parser.add_argument("--lora-alpha", type=int, default=32)
    parser.add_argument("--lora-dropout", type=float, default=0.05)
    args = parser.parse_args()

    cfg = TrainConfig(
        model_name_or_path=args.model,
        dataset_path=args.dataset,
        output_dir=args.output_dir,
        max_length=args.max_length,
        learning_rate=args.learning_rate,
        per_device_train_batch_size=args.batch_size,
        gradient_accumulation_steps=args.grad_accum,
        num_train_epochs=args.epochs,
        logging_steps=args.logging_steps,
        save_steps=args.save_steps,
        lora_r=args.lora_r,
        lora_alpha=args.lora_alpha,
        lora_dropout=args.lora_dropout,
    )
    run_training(cfg)


if __name__ == "__main__":
    main()

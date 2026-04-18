from dataclasses import asdict
from random import Random

from .types import TrainingExample


def split_examples(
    examples: list[TrainingExample],
    train_ratio: float = 0.9,
    val_ratio: float = 0.05,
    seed: int = 42,
) -> dict[str, list[TrainingExample]]:
    if train_ratio <= 0 or val_ratio < 0 or train_ratio + val_ratio >= 1:
        raise ValueError("invalid split ratios")

    items = list(examples)
    Random(seed).shuffle(items)

    n = len(items)
    train_end = int(n * train_ratio)
    val_end = train_end + int(n * val_ratio)

    return {
        "train": items[:train_end],
        "validation": items[train_end:val_end],
        "test": items[val_end:],
    }


def as_hf_records(examples: list[TrainingExample]) -> list[dict[str, object]]:
    return [
        {
            "system": item.system,
            "user": item.user,
            "assistant": item.assistant,
            "metadata": item.metadata,
        }
        for item in examples
    ]


def as_lora_sft_records(examples: list[TrainingExample]) -> list[dict[str, object]]:
    records: list[dict[str, object]] = []
    for item in examples:
        records.append(
            {
                "messages": [
                    {"role": "system", "content": item.system},
                    {"role": "user", "content": item.user},
                    {"role": "assistant", "content": item.assistant},
                ],
                "metadata": item.metadata,
            }
        )
    return records


def split_manifest(splits: dict[str, list[TrainingExample]]) -> dict[str, object]:
    return {
        name: {
            "count": len(items),
            "sample_metadata": asdict(items[0]) if items and hasattr(items[0], "__dataclass_fields__") else None,
        }
        for name, items in splits.items()
    }

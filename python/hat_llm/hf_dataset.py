from __future__ import annotations

from pathlib import Path
from typing import Any

from .splits import as_hf_records
from .types import TrainingExample


def dataset_dict_from_splits(splits: dict[str, list[TrainingExample]]):
    try:
        from datasets import Dataset, DatasetDict
    except ImportError as exc:  # pragma: no cover
        raise RuntimeError("datasets is required for direct Dataset export; install hat-llm[hf]") from exc

    materialized: dict[str, Any] = {}
    for name, items in splits.items():
        materialized[name] = Dataset.from_list(as_hf_records(items))
    return DatasetDict(materialized)


def save_dataset_dict(splits: dict[str, list[TrainingExample]], output_dir: str | Path) -> Path:
    dataset = dataset_dict_from_splits(splits)
    out = Path(output_dir)
    out.parent.mkdir(parents=True, exist_ok=True)
    dataset.save_to_disk(str(out))
    return out

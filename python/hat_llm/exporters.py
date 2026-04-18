import json
from pathlib import Path

from .splits import as_hf_records, as_lora_sft_records
from .types import TrainingExample


def write_json(path: str | Path, payload: object) -> Path:
    out = Path(path)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(payload, indent=2, ensure_ascii=False), encoding="utf-8")
    return out


def write_jsonl(path: str | Path, records: list[dict[str, object]]) -> Path:
    out = Path(path)
    out.parent.mkdir(parents=True, exist_ok=True)
    with out.open("w", encoding="utf-8") as handle:
        for row in records:
            handle.write(json.dumps(row, ensure_ascii=False) + "\n")
    return out


def export_hf_json(path: str | Path, examples: list[TrainingExample]) -> Path:
    return write_json(path, as_hf_records(examples))


def export_lora_jsonl(path: str | Path, examples: list[TrainingExample]) -> Path:
    return write_jsonl(path, as_lora_sft_records(examples))

import argparse
import json
from pathlib import Path

from .dataset import DatasetBuilder
from .examples import starter_runtime_state, starter_slots, starter_tasks
from .exporters import export_hf_json, export_lora_jsonl
from .postgres_loader import PostgresCorpusLoader
from .splits import split_examples, split_manifest
from .trainer import HatTrainer


def main() -> None:
    parser = argparse.ArgumentParser(description="Export Harness-Aware Training datasets")
    parser.add_argument("--output", default="artifacts/hat_sft.jsonl", help="output JSONL path")
    parser.add_argument("--report", default="artifacts/hat_report.json", help="output validation report path")
    parser.add_argument("--repo-root", default=".", help="repository root for loading real specialist slots")
    parser.add_argument("--specialist-id", default="csse-tool-development-specialist-01", help="specialist id to load")
    parser.add_argument("--dsn", default="", help="optional Postgres DSN for loading postmortems/evals")
    parser.add_argument("--hf-out", default="artifacts/hf_dataset.json", help="output JSON path for Hugging Face style records")
    parser.add_argument("--lora-out", default="artifacts/lora_sft.jsonl", help="output JSONL path for LoRA-style messages")
    args = parser.parse_args()

    trainer = HatTrainer()
    runtime = starter_runtime_state()
    runtime.specialist_id = args.specialist_id
    runtime.namespace = f"memory.{args.specialist_id}"

    slots = trainer.load_slots_from_repo(args.repo_root, args.specialist_id)
    if not slots.identity.strip():
        slots = starter_slots()

    tasks = starter_tasks()

    if args.dsn:
        loader = PostgresCorpusLoader(args.dsn)
        corpus_records = loader.load_postmortems(args.specialist_id, limit=50) + loader.load_eval_cases(args.specialist_id, limit=50)
        tasks.extend(trainer.tasks_from_corpus(corpus_records, prefix="db"))

    examples, reports = trainer.build_training_examples(runtime, slots, tasks)

    builder = DatasetBuilder()
    out_path = builder.export_jsonl(examples, args.output)

    splits = split_examples(examples)
    hf_path = export_hf_json(args.hf_out, splits["train"])
    lora_path = export_lora_jsonl(args.lora_out, splits["train"])

    report_path = Path(args.report)
    report_path.parent.mkdir(parents=True, exist_ok=True)
    report_payload = {
        "reports": reports,
        "split_manifest": split_manifest(splits),
        "output": str(out_path),
        "hf_output": str(hf_path),
        "lora_output": str(lora_path),
    }
    report_path.write_text(json.dumps(report_payload, indent=2), encoding="utf-8")

    print(f"exported {len(examples)} training examples to {out_path}")
    print(f"wrote report to {report_path}")
    print(f"wrote HF-style train split to {hf_path}")
    print(f"wrote LoRA-style train split to {lora_path}")


if __name__ == "__main__":
    main()

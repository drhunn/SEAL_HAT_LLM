import argparse
import json
from pathlib import Path

from .dataset import DatasetBuilder
from .examples import starter_runtime_state, starter_slots, starter_tasks
from .trainer import HatTrainer


def main() -> None:
    parser = argparse.ArgumentParser(description="Export starter Harness-Aware Training JSONL")
    parser.add_argument("--output", default="artifacts/hat_sft.jsonl", help="output JSONL path")
    parser.add_argument("--report", default="artifacts/hat_report.json", help="output validation report path")
    args = parser.parse_args()

    runtime = starter_runtime_state()
    slots = starter_slots()
    tasks = starter_tasks()

    trainer = HatTrainer()
    examples, reports = trainer.build_training_examples(runtime, slots, tasks)

    builder = DatasetBuilder()
    out_path = builder.export_jsonl(examples, args.output)

    report_path = Path(args.report)
    report_path.parent.mkdir(parents=True, exist_ok=True)
    report_path.write_text(json.dumps(reports, indent=2), encoding="utf-8")

    print(f"exported {len(examples)} training examples to {out_path}")
    print(f"wrote validation report to {report_path}")


if __name__ == "__main__":
    main()

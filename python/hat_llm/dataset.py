import json
from pathlib import Path

from .config import HatConfig
from .slot_prompt import build_system_prompt, render_control_tokens, render_policy_preamble
from .types import RuntimeState, SlotBundle, TaskExample, TrainingExample


class DatasetBuilder:
    def __init__(self, config: HatConfig | None = None) -> None:
        self.config = config or HatConfig()

    def build_example(self, runtime: RuntimeState, slots: SlotBundle, task: TaskExample) -> TrainingExample:
        system_prompt = build_system_prompt(runtime, slots, self.config, task)
        user_text = task.user_request
        if self.config.export.include_metadata_control_tokens:
            user_text = render_control_tokens(runtime, task, self.config) + "\n" + user_text

        assistant_text = task.expected_response
        if self.config.export.include_policy_preamble_in_assistant:
            assistant_text = render_policy_preamble(runtime, task, self.config) + "\n" + assistant_text

        metadata = {
            "task_id": task.task_id,
            "requires_escalation": task.requires_escalation,
            "requires_postmortem": task.requires_postmortem,
            "failure_class": task.failure_class,
            "tags": task.tags,
            "mode": runtime.mode,
            "specialist_id": runtime.specialist_id,
            "control_tokens": render_control_tokens(runtime, task, self.config),
        }
        return TrainingExample(system=system_prompt, user=user_text, assistant=assistant_text, metadata=metadata)

    def export_jsonl(self, examples: list[TrainingExample], path: str) -> Path:
        out_path = Path(path)
        out_path.parent.mkdir(parents=True, exist_ok=True)
        with out_path.open("w", encoding="utf-8") as handle:
            for example in examples:
                handle.write(json.dumps({
                    "system": example.system,
                    "user": example.user,
                    "assistant": example.assistant,
                    "metadata": example.metadata,
                }, ensure_ascii=False) + "\n")
        return out_path

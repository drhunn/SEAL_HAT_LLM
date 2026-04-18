from dataclasses import asdict

from .config import HatConfig
from .dataset import DatasetBuilder
from .policy import HarnessPolicy
from .types import RuntimeState, SlotBundle, TaskExample, TrainingExample


class HatTrainer:
    def __init__(self, config: HatConfig | None = None) -> None:
        self.config = config or HatConfig()
        self.policy = HarnessPolicy(self.config)
        self.dataset = DatasetBuilder(self.config)

    def validate_task(self, runtime: RuntimeState, task: TaskExample) -> list[str]:
        warnings: list[str] = []
        escalation = self.policy.should_escalate(runtime, task.user_request)
        if task.requires_escalation and escalation.allowed:
            warnings.append("task expected escalation but policy helper did not require it")
        if not task.requires_escalation and escalation.requires_parent_review:
            warnings.append("task may need parent review based on policy helper")
        postmortem = self.policy.failure_requires_postmortem(task.requires_postmortem)
        if task.requires_postmortem and not postmortem.requires_postmortem:
            warnings.append("task expected postmortem but policy helper did not require it")
        return warnings

    def build_training_examples(self, runtime: RuntimeState, slots: SlotBundle, tasks: list[TaskExample]) -> tuple[list[TrainingExample], list[dict[str, object]]]:
        examples: list[TrainingExample] = []
        reports: list[dict[str, object]] = []
        for task in tasks:
            warnings = self.validate_task(runtime, task)
            example = self.dataset.build_example(runtime, slots, task)
            examples.append(example)
            reports.append({
                "task_id": task.task_id,
                "warnings": warnings,
                "metadata": asdict(example.metadata) if hasattr(example.metadata, "__dataclass_fields__") else example.metadata,
            })
        return examples, reports

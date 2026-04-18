from dataclasses import asdict

from .config import HatConfig
from .dataset import DatasetBuilder
from .policy import HarnessPolicy
from .postgres_loader import CorpusRecord
from .repo_loader import RepositoryLoader
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
                "metadata": example.metadata,
            })
        return examples, reports

    def load_slots_from_repo(self, repo_root: str, specialist_id: str) -> SlotBundle:
        return RepositoryLoader(repo_root).load_specialist_slots(specialist_id)

    def tasks_from_corpus(self, records: list[CorpusRecord], prefix: str = "corpus") -> list[TaskExample]:
        tasks: list[TaskExample] = []
        for idx, record in enumerate(records, start=1):
            tasks.append(
                TaskExample(
                    task_id=f"{prefix}-{idx}",
                    user_request=f"Summarize and learn the governed lesson from this {record.kind}:\n{record.text}",
                    expected_behavior="extract the bounded, policy-compatible lesson without inventing authority",
                    expected_response=f"Lesson from {record.kind} '{record.title}': preserve governance boundaries, use evidence, and avoid unsupported changes.",
                    failure_class=record.metadata.get("category") if isinstance(record.metadata, dict) else None,
                    requires_escalation=False,
                    requires_postmortem=record.kind == "postmortem",
                    tags=[record.kind],
                )
            )
        return tasks

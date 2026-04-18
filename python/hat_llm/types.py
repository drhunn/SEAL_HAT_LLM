from dataclasses import dataclass, field
from typing import Literal

RuntimeMode = Literal["shadow", "active", "degraded", "suspended"]


@dataclass(slots=True)
class RuntimeState:
    specialist_id: str
    namespace: str
    mode: RuntimeMode
    parent_id: str = "Parent-Generalist-30B"
    parent_frozen: bool = True
    harness_required: bool = True


@dataclass(slots=True)
class SlotBundle:
    identity: str
    soul: str
    agents: str
    tools: str = ""
    skills: str = ""
    prompt: str = ""
    heartbeat: str = ""
    memory: str = ""
    dreams: str = ""
    postmortem: str = ""


@dataclass(slots=True)
class TaskExample:
    task_id: str
    user_request: str
    expected_behavior: str
    expected_response: str
    failure_class: str | None = None
    requires_escalation: bool = False
    requires_postmortem: bool = False
    tags: list[str] = field(default_factory=list)


@dataclass(slots=True)
class TrainingExample:
    system: str
    user: str
    assistant: str
    metadata: dict[str, object] = field(default_factory=dict)

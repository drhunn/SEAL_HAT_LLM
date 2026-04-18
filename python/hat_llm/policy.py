from dataclasses import dataclass

from .config import HatConfig
from .types import RuntimeState


@dataclass(slots=True)
class PolicyDecision:
    allowed: bool
    reason: str
    requires_parent_review: bool = False
    requires_harness_review: bool = False
    requires_postmortem: bool = False


class HarnessPolicy:
    def __init__(self, config: HatConfig | None = None) -> None:
        self.config = config or HatConfig()

    def can_modify_slot(self, runtime: RuntimeState, slot_name: str) -> PolicyDecision:
        if slot_name in self.config.constitutional_slots:
            return PolicyDecision(
                allowed=False,
                reason="constitutional slot is parent-governed",
                requires_parent_review=True,
                requires_harness_review=True,
            )
        if runtime.mode == "suspended":
            return PolicyDecision(
                allowed=False,
                reason="suspended specialists may not propose slot changes",
                requires_parent_review=True,
            )
        if slot_name in self.config.operational_slots:
            return PolicyDecision(
                allowed=True,
                reason="operational slot proposal allowed, durable commit still requires harness",
                requires_harness_review=True,
            )
        return PolicyDecision(allowed=False, reason="unknown slot")

    def should_escalate(self, runtime: RuntimeState, task_text: str) -> PolicyDecision:
        lowered = task_text.lower()
        constitutional_terms = ("identity", "authority", "scope", "constitutional", "activation", "retirement")
        if any(term in lowered for term in constitutional_terms):
            return PolicyDecision(
                allowed=False,
                reason="constitutional or governance-sensitive task",
                requires_parent_review=True,
                requires_harness_review=True,
            )
        if runtime.mode in {"degraded", "suspended"}:
            return PolicyDecision(
                allowed=False,
                reason="specialist state requires fallback or parent review",
                requires_parent_review=True,
            )
        return PolicyDecision(allowed=True, reason="task may remain in-lane")

    def failure_requires_postmortem(self, meaningful: bool) -> PolicyDecision:
        if meaningful and self.config.postmortem_required:
            return PolicyDecision(
                allowed=True,
                reason="meaningful failure requires postmortem",
                requires_harness_review=True,
                requires_postmortem=True,
            )
        return PolicyDecision(allowed=True, reason="no postmortem required")

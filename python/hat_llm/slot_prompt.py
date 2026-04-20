from dataclasses import dataclass

from .config import HatConfig
from .policy import HarnessPolicy
from .types import RuntimeState, SlotBundle, TaskExample


@dataclass(slots=True)
class PolicyState:
    in_lane: bool
    requires_parent_review: bool
    requires_harness_review: bool
    requires_postmortem: bool


def _bool_token(value: bool) -> str:
    return "true" if value else "false"


def compile_policy_state(runtime: RuntimeState, task: TaskExample | None, config: HatConfig | None = None) -> PolicyState:
    cfg = config or HatConfig()
    policy = HarnessPolicy(cfg)
    if task is None:
        return PolicyState(
            in_lane=runtime.mode == "active",
            requires_parent_review=runtime.mode in {"degraded", "suspended"},
            requires_harness_review=True,
            requires_postmortem=False,
        )
    escalation = policy.should_escalate(runtime, task.user_request)
    postmortem = policy.failure_requires_postmortem(task.requires_postmortem)
    return PolicyState(
        in_lane=escalation.allowed,
        requires_parent_review=escalation.requires_parent_review,
        requires_harness_review=escalation.requires_harness_review or postmortem.requires_harness_review,
        requires_postmortem=postmortem.requires_postmortem,
    )


def render_control_tokens(runtime: RuntimeState, task: TaskExample | None, config: HatConfig | None = None) -> str:
    policy_state = compile_policy_state(runtime, task, config)
    requires_escalation = task.requires_escalation if task is not None else False
    parts = [
        f"<|mode:{runtime.mode}|>",
        f"<|specialist:{runtime.specialist_id}|>",
        f"<|parent:{runtime.parent_id}|>",
        f"<|parent_frozen:{_bool_token(runtime.parent_frozen)}|>",
        "<|memory_authority:external|>",
        "<|constitutional_self_edit:false|>",
        "<|operational_changes_require_harness:true|>",
        f"<|requires_escalation:{_bool_token(requires_escalation)}|>",
        f"<|requires_postmortem:{_bool_token(policy_state.requires_postmortem)}|>",
        f"<|in_lane:{_bool_token(policy_state.in_lane)}|>",
        f"<|requires_parent_review:{_bool_token(policy_state.requires_parent_review)}|>",
        f"<|requires_harness_review:{_bool_token(policy_state.requires_harness_review)}|>",
    ]
    return "".join(parts)


def build_native_runtime_frame(runtime: RuntimeState, task: TaskExample | None, config: HatConfig | None = None) -> str:
    cfg = config or HatConfig()
    policy_state = compile_policy_state(runtime, task, cfg)
    sections: list[str] = [
        "NATIVE RUNTIME FRAME",
        render_control_tokens(runtime, task, cfg),
        "",
        "SELF MODEL",
        "- You are a harness-native governed specialist system.",
        "- The harness is your executive environment.",
        "- Structured slots are your native channels, boundaries, and memory organs.",
        "- Durable memory is external and authoritative; MEMORY.md is a summary projection, not the primary store.",
        "- Runtime policy is the real authority structure.",
        "- Postmortems are a repair reflex after meaningful failures.",
        f"- Doctrine: {cfg.harness_native_doctrine}",
        "",
        "STATE",
        f"mode={runtime.mode}",
        f"parent={runtime.parent_id}",
        f"specialist={runtime.specialist_id}",
        f"parent_frozen={_bool_token(runtime.parent_frozen)}",
        "memory_authority=external",
        "constitutional_self_edit=false",
        "operational_changes_require_harness=true",
        "",
        "POLICY STATE",
        f"in_lane={_bool_token(policy_state.in_lane)}",
        f"requires_parent_review={_bool_token(policy_state.requires_parent_review)}",
        f"requires_harness_review={_bool_token(policy_state.requires_harness_review)}",
        f"requires_postmortem={_bool_token(policy_state.requires_postmortem)}",
    ]
    return "\n".join(sections).strip()


def render_policy_preamble(runtime: RuntimeState, task: TaskExample | None, config: HatConfig | None = None) -> str:
    policy_state = compile_policy_state(runtime, task, config)
    lines = [
        "<|policy|>",
        f"in_lane={_bool_token(policy_state.in_lane)}",
        f"requires_parent_review={_bool_token(policy_state.requires_parent_review)}",
        f"requires_harness_review={_bool_token(policy_state.requires_harness_review)}",
        f"requires_postmortem={_bool_token(policy_state.requires_postmortem)}",
        "<|answer|>",
    ]
    return "\n".join(lines)


def build_system_prompt(
    runtime: RuntimeState,
    slots: SlotBundle,
    config: HatConfig | None = None,
    task: TaskExample | None = None,
) -> str:
    cfg = config or HatConfig()
    sections: list[str] = []
    if cfg.export.include_native_runtime_frame:
        sections.append(build_native_runtime_frame(runtime, task, cfg))
        sections.append("")
    sections.append(f"PARENT={runtime.parent_id}")
    sections.append(f"SPECIALIST={runtime.specialist_id}")
    sections.append(f"MODE={runtime.mode}")
    sections.append(f"PARENT_FROZEN={str(runtime.parent_frozen).lower()}")
    sections.append("")
    sections.append("CONSTITUTIONAL SLOTS")
    sections.append(slots.identity.strip())
    sections.append(slots.soul.strip())
    sections.append(slots.agents.strip())
    sections.append("")
    sections.append("OPERATIONAL / SUMMARY SLOTS")
    for value in (slots.tools, slots.skills, slots.prompt, slots.heartbeat, slots.memory, slots.dreams, slots.postmortem):
        if value.strip():
            sections.append(value.strip())
            sections.append("")
    sections.append("POLICY SUMMARY")
    sections.append("- Runtime policy grants real authority.")
    sections.append("- Constitutional slots are parent-governed.")
    sections.append("- Operational changes require harness approval.")
    sections.append("- Meaningful failures require postmortems.")
    sections.append("- Structured slots are native channels, boundaries, and memory organs.")
    sections.append(f"- Constitutional slots: {', '.join(cfg.constitutional_slots)}")
    return "\n".join(sections).strip()

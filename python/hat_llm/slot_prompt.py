from .config import HatConfig
from .types import RuntimeState, SlotBundle


def build_system_prompt(runtime: RuntimeState, slots: SlotBundle, config: HatConfig | None = None) -> str:
    cfg = config or HatConfig()
    sections: list[str] = []
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
    sections.append(f"- Constitutional slots: {', '.join(cfg.constitutional_slots)}")
    return "\n".join(sections).strip()

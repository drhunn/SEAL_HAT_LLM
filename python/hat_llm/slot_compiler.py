from __future__ import annotations

from .slot_ir import SlotIR, SlotRecord
from .types import RuntimeState, SlotBundle


_SLOT_FAMILY_BY_NAME = {
    "IDENTITY.md": "constitutional",
    "SOUL.md": "constitutional",
    "AGENTS.md": "constitutional",
    "TOOLS.md": "tools",
    "SKILLS.md": "skills",
    "PROMPT.md": "prompt",
    "HEARTBEAT.md": "heartbeat",
    "MEMORY.md": "memory",
    "DREAMS.md": "dreams",
    "POSTMORTEM.md": "postmortem",
}


def _authority_for_family(family: str) -> str:
    if family == "constitutional":
        return "parent"
    return "harness"


def _summary_from_text(text: str) -> str:
    lines = [line.strip() for line in text.splitlines() if line.strip()]
    if not lines:
        return ""
    return " ".join(lines[:2])[:240]


def _rules_from_text(text: str) -> list[str]:
    rules: list[str] = []
    for line in text.splitlines():
        stripped = line.strip()
        if stripped.startswith("-"):
            rules.append(stripped.lstrip("- ").strip())
    return rules[:16]


def compile_slot_ir(runtime: RuntimeState, slots: SlotBundle) -> SlotIR:
    records: list[SlotRecord] = []
    for slot_name, text in (
        ("IDENTITY.md", slots.identity),
        ("SOUL.md", slots.soul),
        ("AGENTS.md", slots.agents),
        ("TOOLS.md", slots.tools),
        ("SKILLS.md", slots.skills),
        ("PROMPT.md", slots.prompt),
        ("HEARTBEAT.md", slots.heartbeat),
        ("MEMORY.md", slots.memory),
        ("DREAMS.md", slots.dreams),
        ("POSTMORTEM.md", slots.postmortem),
    ):
        if not text.strip():
            continue
        family = _SLOT_FAMILY_BY_NAME[slot_name]
        records.append(
            SlotRecord(
                slot_name=slot_name,
                family=family,
                authority=_authority_for_family(family),
                text=text,
                summary=_summary_from_text(text),
                rules=_rules_from_text(text),
                hot_reloadable=family != "constitutional",
            )
        )
    return SlotIR(
        specialist_id=runtime.specialist_id,
        mode=runtime.mode,
        parent_frozen=runtime.parent_frozen,
        records=records,
    )

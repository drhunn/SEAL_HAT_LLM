from dataclasses import dataclass, field


@dataclass(slots=True)
class ExportConfig:
    output_path: str = "artifacts/hat_sft.jsonl"
    include_slot_context: bool = True
    include_policy_summary: bool = True
    include_native_runtime_frame: bool = True
    include_metadata_control_tokens: bool = True
    include_policy_preamble_in_assistant: bool = True


@dataclass(slots=True)
class HatConfig:
    parent_id: str = "Parent-Generalist-30B"
    constitutional_slots: tuple[str, ...] = ("IDENTITY.md", "SOUL.md", "AGENTS.md")
    operational_slots: tuple[str, ...] = (
        "TOOLS.md",
        "SKILLS.md",
        "PROMPT.md",
        "HEARTBEAT.md",
        "MEMORY.md",
        "DREAMS.md",
        "POSTMORTEM.md",
    )
    postmortem_required: bool = True
    harness_native_doctrine: str = (
        "The model is a harness-native intelligence whose identity, authority, memory, "
        "and repair pathways are expressed through structured slots and governed runtime channels."
    )
    export: ExportConfig = field(default_factory=ExportConfig)

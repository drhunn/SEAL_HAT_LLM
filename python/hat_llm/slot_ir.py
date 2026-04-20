from __future__ import annotations

from dataclasses import dataclass, field


@dataclass(slots=True)
class SlotRecord:
    slot_name: str
    family: str
    authority: str
    text: str
    summary: str
    tags: list[str] = field(default_factory=list)
    rules: list[str] = field(default_factory=list)
    procedures: list[str] = field(default_factory=list)
    hot_reloadable: bool = True


@dataclass(slots=True)
class SlotIR:
    specialist_id: str
    mode: str
    parent_frozen: bool
    records: list[SlotRecord] = field(default_factory=list)

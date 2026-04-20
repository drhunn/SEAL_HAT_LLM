from __future__ import annotations

from dataclasses import dataclass


@dataclass(slots=True)
class SlotPacketConfig:
    max_slots: int = 16
    max_chars_per_slot: int = 512


@dataclass(slots=True)
class SlotPacket:
    slot_texts: list[str]
    slot_family_ids: list[int]
    slot_authority_ids: list[int]
    slot_enabled_mask: list[int]
    mode_id: int
    parent_frozen: int
    scalar_flags: list[int]

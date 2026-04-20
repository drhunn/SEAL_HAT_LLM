from __future__ import annotations

from .slot_compiler import compile_slot_ir
from .slot_ir import SlotIR
from .slot_packet import SlotPacket, SlotPacketConfig
from .types import RuntimeState, SlotBundle


_FAMILY_IDS = {
    "constitutional": 0,
    "tools": 1,
    "skills": 2,
    "prompt": 3,
    "heartbeat": 4,
    "memory": 5,
    "dreams": 6,
    "postmortem": 7,
}

_AUTHORITY_IDS = {
    "parent": 0,
    "harness": 1,
    "specialist": 2,
}

_MODE_IDS = {
    "shadow": 0,
    "active": 1,
    "degraded": 2,
    "suspended": 3,
}


def build_slot_packet_from_ir(ir: SlotIR, config: SlotPacketConfig | None = None) -> SlotPacket:
    cfg = config or SlotPacketConfig()
    records = ir.records[: cfg.max_slots]
    slot_texts = [record.text[: cfg.max_chars_per_slot] for record in records]
    slot_family_ids = [_FAMILY_IDS.get(record.family, 0) for record in records]
    slot_authority_ids = [_AUTHORITY_IDS.get(record.authority, 1) for record in records]
    slot_enabled_mask = [1 for _ in records]

    while len(slot_texts) < cfg.max_slots:
        slot_texts.append("")
        slot_family_ids.append(0)
        slot_authority_ids.append(1)
        slot_enabled_mask.append(0)

    scalar_flags = [
        1,  # memory authority is external
        0,  # constitutional self-edit forbidden
        1,  # operational changes require harness
        1,  # meaningful failures require postmortem
    ]

    return SlotPacket(
        slot_texts=slot_texts,
        slot_family_ids=slot_family_ids,
        slot_authority_ids=slot_authority_ids,
        slot_enabled_mask=slot_enabled_mask,
        mode_id=_MODE_IDS.get(ir.mode, 1),
        parent_frozen=1 if ir.parent_frozen else 0,
        scalar_flags=scalar_flags,
    )


def build_slot_packet(runtime: RuntimeState, slots: SlotBundle, config: SlotPacketConfig | None = None) -> SlotPacket:
    ir = compile_slot_ir(runtime, slots)
    return build_slot_packet_from_ir(ir, config)

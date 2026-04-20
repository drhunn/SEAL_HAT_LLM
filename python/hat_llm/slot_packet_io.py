from __future__ import annotations

import json
from pathlib import Path

import jax.numpy as jnp

from .slot_packet import SlotPacket

_MODE_IDS = {
    "shadow": 0,
    "active": 1,
    "degraded": 2,
    "suspended": 3,
}


def load_slot_packet_json(path: str | Path) -> SlotPacket:
    data = json.loads(Path(path).read_text(encoding="utf-8"))
    mode = str(data.get("mode", "active"))
    return SlotPacket(
        slot_texts=list(data.get("slot_texts", [])),
        slot_family_ids=[int(v) for v in data.get("slot_family_ids", [])],
        slot_authority_ids=[int(v) for v in data.get("slot_authority_ids", [])],
        slot_enabled_mask=[int(v) for v in data.get("slot_enabled_mask", [])],
        mode_id=int(data.get("mode_id", _MODE_IDS.get(mode, 1))),
        parent_frozen=int(data.get("parent_frozen", 1 if data.get("parent_frozen") else 0)),
        scalar_flags=[int(v) for v in data.get("scalar_flags", [1, 0, 1, 1])],
    )


def slot_packet_to_jax_arrays(packet: SlotPacket, tokenizer, slot_max_length: int, slot_max_count: int) -> dict[str, jnp.ndarray]:
    slot_texts = list(packet.slot_texts[:slot_max_count])
    slot_family_ids = list(packet.slot_family_ids[:slot_max_count])
    slot_authority_ids = list(packet.slot_authority_ids[:slot_max_count])
    slot_enabled_mask = list(packet.slot_enabled_mask[:slot_max_count])

    while len(slot_texts) < slot_max_count:
        slot_texts.append("")
        slot_family_ids.append(0)
        slot_authority_ids.append(1)
        slot_enabled_mask.append(0)

    encoded = tokenizer(
        slot_texts,
        truncation=True,
        max_length=slot_max_length,
        padding="max_length",
    )
    return {
        "slot_token_ids": jnp.asarray(encoded["input_ids"], dtype=jnp.int32),
        "slot_token_mask": jnp.asarray(encoded["attention_mask"], dtype=jnp.int32),
        "slot_family_ids": jnp.asarray(slot_family_ids, dtype=jnp.int32),
        "slot_authority_ids": jnp.asarray(slot_authority_ids, dtype=jnp.int32),
        "slot_enabled_mask": jnp.asarray(slot_enabled_mask, dtype=jnp.float32),
    }

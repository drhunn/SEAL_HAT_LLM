from __future__ import annotations

import flax.linen as nn
import jax.numpy as jnp

from .slot_adapters import SlotAdapter
from .slot_encoder import SlotEncoder


class HybridSlotAwareBlock(nn.Module):
    d_model: int
    d_slot: int

    @nn.compact
    def __call__(self, x, global_slot_state):
        attn_norm = nn.LayerNorm()(x)
        attn_out = nn.Dense(self.d_model)(attn_norm)
        x = x + attn_out
        x = x + SlotAdapter(d_model=self.d_model, d_slot=self.d_slot)(x, global_slot_state)

        mlp_norm = nn.LayerNorm()(x)
        mlp_hidden = nn.gelu(nn.Dense(self.d_model * 4)(mlp_norm))
        mlp_out = nn.Dense(self.d_model)(mlp_hidden)
        x = x + mlp_out
        x = x + SlotAdapter(d_model=self.d_model, d_slot=self.d_slot)(x, global_slot_state)
        return x


class HybridSlotAwareModel(nn.Module):
    vocab_size: int
    d_model: int
    d_slot: int
    num_layers: int

    @nn.compact
    def __call__(self, input_ids, global_slot_state):
        x = nn.Embed(num_embeddings=self.vocab_size, features=self.d_model)(input_ids)
        for _ in range(self.num_layers):
            x = HybridSlotAwareBlock(d_model=self.d_model, d_slot=self.d_slot)(x, global_slot_state)
        x = nn.LayerNorm()(x)
        return nn.Dense(self.vocab_size)(x)


class HybridRuntimeSlotModel(nn.Module):
    vocab_size: int
    d_model: int
    d_slot: int
    num_layers: int
    num_slot_families: int = 8
    num_authority_types: int = 3

    @nn.compact
    def __call__(
        self,
        input_ids,
        slot_token_ids,
        slot_token_mask,
        slot_family_ids,
        slot_authority_ids,
        slot_enabled_mask,
    ):
        slot_vecs = SlotEncoder(
            d_model=self.d_model,
            d_slot=self.d_slot,
            vocab_size=self.vocab_size,
            num_slot_families=self.num_slot_families,
            num_authority_types=self.num_authority_types,
        )(
            token_ids=slot_token_ids,
            token_mask=slot_token_mask,
            family_ids=slot_family_ids,
            authority_ids=slot_authority_ids,
        )
        slot_mask = slot_enabled_mask[:, None]
        denom = jnp.maximum(jnp.sum(slot_mask), 1)
        global_slot_state = jnp.sum(slot_vecs * slot_mask, axis=0) / denom
        return HybridSlotAwareModel(
            vocab_size=self.vocab_size,
            d_model=self.d_model,
            d_slot=self.d_slot,
            num_layers=self.num_layers,
        )(input_ids=input_ids, global_slot_state=global_slot_state)

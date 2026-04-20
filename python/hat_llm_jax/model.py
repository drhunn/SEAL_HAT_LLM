from __future__ import annotations

import flax.linen as nn
import jax.numpy as jnp

from .slot_adapters import SlotAdapter


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

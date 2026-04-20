from __future__ import annotations

import flax.linen as nn
import jax.numpy as jnp


class SlotEncoder(nn.Module):
    d_model: int
    d_slot: int
    vocab_size: int
    num_slot_families: int
    num_authority_types: int

    @nn.compact
    def __call__(self, token_ids, token_mask, family_ids, authority_ids):
        tok_emb = nn.Embed(num_embeddings=self.vocab_size, features=self.d_model)(token_ids)
        fam_emb = nn.Embed(num_embeddings=self.num_slot_families, features=self.d_model)(family_ids)[:, None, :]
        auth_emb = nn.Embed(num_embeddings=self.num_authority_types, features=self.d_model)(authority_ids)[:, None, :]
        x = tok_emb + fam_emb + auth_emb
        x = x * token_mask[..., None]
        denom = jnp.maximum(jnp.sum(token_mask, axis=1, keepdims=True), 1)
        pooled = jnp.sum(x, axis=1) / denom
        return nn.Dense(self.d_slot)(pooled)

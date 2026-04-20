from __future__ import annotations

import flax.linen as nn


class SlotAdapter(nn.Module):
    d_model: int
    d_slot: int

    @nn.compact
    def __call__(self, x, slot_state):
        gate = nn.sigmoid(nn.Dense(self.d_model)(slot_state))
        delta = nn.Dense(self.d_model)(nn.tanh(nn.Dense(self.d_slot)(x)))
        return delta * gate

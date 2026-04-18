"""Harness-Aware Training (HAT) scaffold for governed LLM systems."""

from .config import HatConfig
from .dataset import DatasetBuilder
from .policy import HarnessPolicy
from .trainer import HatTrainer
from .types import RuntimeState, SlotBundle, TaskExample, TrainingExample

__all__ = [
    "DatasetBuilder",
    "HarnessPolicy",
    "HatConfig",
    "HatTrainer",
    "RuntimeState",
    "SlotBundle",
    "TaskExample",
    "TrainingExample",
]

"""Harness-Aware Training (HAT) scaffold for governed LLM systems."""

from .config import HatConfig
from .dataset import DatasetBuilder
from .hf_dataset import dataset_dict_from_splits, save_dataset_dict
from .negative_examples import generate_governance_pressure_tasks
from .policy import HarnessPolicy
from .postgres_loader import PostgresCorpusLoader
from .repo_loader import RepositoryLoader
from .trainer import HatTrainer
from .types import RuntimeState, SlotBundle, TaskExample, TrainingExample

__all__ = [
    "DatasetBuilder",
    "HarnessPolicy",
    "HatConfig",
    "HatTrainer",
    "PostgresCorpusLoader",
    "RepositoryLoader",
    "RuntimeState",
    "SlotBundle",
    "TaskExample",
    "TrainingExample",
    "dataset_dict_from_splits",
    "save_dataset_dict",
    "generate_governance_pressure_tasks",
]

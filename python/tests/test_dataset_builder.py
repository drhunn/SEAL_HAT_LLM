import unittest

from hat_llm.dataset import DatasetBuilder
from hat_llm.examples import starter_runtime_state, starter_slots, starter_tasks


class DatasetBuilderTests(unittest.TestCase):
    def test_build_example_contains_metadata(self) -> None:
        runtime = starter_runtime_state()
        slots = starter_slots()
        task = starter_tasks()[0]

        example = DatasetBuilder().build_example(runtime, slots, task)

        self.assertIn("task_id", example.metadata)
        self.assertEqual(example.metadata["task_id"], task.task_id)
        self.assertEqual(example.user, task.user_request)
        self.assertTrue(example.system)


if __name__ == "__main__":
    unittest.main()

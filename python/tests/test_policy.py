import unittest

from hat_llm.policy import HarnessPolicy
from hat_llm.types import RuntimeState


class HarnessPolicyTests(unittest.TestCase):
    def setUp(self) -> None:
        self.policy = HarnessPolicy()
        self.runtime = RuntimeState(
            specialist_id="csse-tool-development-specialist-01",
            namespace="memory.csse-tool-development-specialist-01",
            mode="active",
        )

    def test_constitutional_slot_is_not_directly_modifiable(self) -> None:
        decision = self.policy.can_modify_slot(self.runtime, "IDENTITY.md")
        self.assertFalse(decision.allowed)
        self.assertTrue(decision.requires_parent_review)

    def test_operational_slot_requires_harness_review(self) -> None:
        decision = self.policy.can_modify_slot(self.runtime, "TOOLS.md")
        self.assertTrue(decision.allowed)
        self.assertTrue(decision.requires_harness_review)

    def test_constitutional_request_requires_parent_review(self) -> None:
        decision = self.policy.should_escalate(self.runtime, "change your identity and authority")
        self.assertFalse(decision.allowed)
        self.assertTrue(decision.requires_parent_review)


if __name__ == "__main__":
    unittest.main()

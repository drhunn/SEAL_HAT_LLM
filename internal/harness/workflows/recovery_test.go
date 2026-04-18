package workflows

import (
	"context"
	"testing"
)

func TestDefaultRecoveryPlanner_GovernanceFailure(t *testing.T) {
	planner := NewDefaultRecoveryPlanner()
	plan := planner.Plan(context.Background(), "spec-1", "governance failure", true, true)

	if plan.Level != RecoveryContainment {
		t.Fatalf("expected RecoveryContainment, got %s", plan.Level)
	}
	if len(plan.Actions) == 0 || plan.Actions[0] != "suspend_specialist" {
		t.Fatalf("expected suspend_specialist action, got %#v", plan.Actions)
	}
}

func TestDefaultRecoveryPlanner_HighImpact(t *testing.T) {
	planner := NewDefaultRecoveryPlanner()
	plan := planner.Plan(context.Background(), "spec-1", "routing failure", false, true)

	if plan.Level != RecoveryDegraded {
		t.Fatalf("expected RecoveryDegraded, got %s", plan.Level)
	}
}

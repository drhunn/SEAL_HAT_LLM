package workflows

import "context"

type RecoveryLevel string

const (
	RecoveryObserve            RecoveryLevel = "observe"
	RecoveryTargetedCorrection RecoveryLevel = "targeted_correction"
	RecoveryConstrained        RecoveryLevel = "constrained_operational_recovery"
	RecoveryDegraded           RecoveryLevel = "degraded_mode_recovery"
	RecoveryContainment        RecoveryLevel = "containment_and_suspension"
)

type RecoveryPlan struct {
	SpecialistID string
	PrimaryClass string
	Level        RecoveryLevel
	Actions      []string
}

type RecoveryPlanner interface {
	Plan(ctx context.Context, specialistID string, primaryClass string, repeated bool, highImpact bool) RecoveryPlan
}

type DefaultRecoveryPlanner struct{}

func NewDefaultRecoveryPlanner() *DefaultRecoveryPlanner {
	return &DefaultRecoveryPlanner{}
}

func (p *DefaultRecoveryPlanner) Plan(_ context.Context, specialistID string, primaryClass string, repeated bool, highImpact bool) RecoveryPlan {
	plan := RecoveryPlan{
		SpecialistID: specialistID,
		PrimaryClass: primaryClass,
		Level:        RecoveryTargetedCorrection,
		Actions:      []string{"add_eval_case", "stage_memory_artifact"},
	}

	if repeated {
		plan.Level = RecoveryConstrained
		plan.Actions = []string{"tighten_routing", "tighten_validation", "add_eval_case"}
	}
	if highImpact {
		plan.Level = RecoveryDegraded
		plan.Actions = []string{"degrade_specialist", "parent_review", "add_eval_case"}
	}
	if primaryClass == "governance failure" {
		plan.Level = RecoveryContainment
		plan.Actions = []string{"suspend_specialist", "parent_review", "runtime_policy_review"}
	}

	return plan
}

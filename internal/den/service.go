package den

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/seal"
)

type GrowthSurface string

const (
	GrowthOperationalPatch GrowthSurface = "operational_patch"
	GrowthAdapter          GrowthSurface = "adapter"
	GrowthSplitSpecialist  GrowthSurface = "split_specialist"
	GrowthNewSpecialist    GrowthSurface = "new_specialist"
	GrowthExpertBank       GrowthSurface = "expert_bank"
)

type FreezePlan struct {
	FreezeParentCore       bool
	FreezeSourceSpecialist bool
	TrainTargets           []string
}

type RollbackPlan struct {
	Strategy      string
	RestoreRefs   []string
	DisableOnFail bool
}

type GrowthPlan struct {
	ID             string
	ProposalID     string
	SpecialistID   string
	Surface        GrowthSurface
	Reason         string
	ExperimentName string
	FreezePlan     FreezePlan
	RollbackPlan   RollbackPlan
}

type GrowthPlanWriter interface {
	WriteGrowthPlan(ctx context.Context, plan GrowthPlan) error
}

type Service struct {
	writer GrowthPlanWriter
	logger *slog.Logger
}

func NewService(writer GrowthPlanWriter, logger *slog.Logger) *Service {
	return &Service{writer: writer, logger: logger}
}

func (s *Service) PlanFromProposal(ctx context.Context, proposal seal.AdaptationProposal) (*GrowthPlan, error) {
	if strings.TrimSpace(proposal.ID) == "" {
		return nil, fmt.Errorf("proposal id is required")
	}
	if strings.TrimSpace(proposal.SpecialistID) == "" {
		return nil, fmt.Errorf("specialist id is required")
	}

	plan := GrowthPlan{
		ID:             fmt.Sprintf("growth|%s", proposal.ID),
		ProposalID:     proposal.ID,
		SpecialistID:   proposal.SpecialistID,
		Surface:        growthSurfaceForProposal(proposal.Surface),
		Reason:         proposal.Reason,
		ExperimentName: fmt.Sprintf("%s-%s", proposal.SpecialistID, growthSurfaceForProposal(proposal.Surface)),
		FreezePlan:     freezePlanForProposal(proposal.Surface),
		RollbackPlan:   rollbackPlanForProposal(proposal.Surface),
	}

	if s.writer != nil {
		if err := s.writer.WriteGrowthPlan(ctx, plan); err != nil {
			return nil, err
		}
	}

	s.logger.InfoContext(ctx, "den growth plan created", "proposal_id", proposal.ID, "specialist_id", proposal.SpecialistID, "surface", plan.Surface)
	return &plan, nil
}

func growthSurfaceForProposal(surface seal.ProposalSurface) GrowthSurface {
	switch surface {
	case seal.SurfaceAdapterTuning:
		return GrowthAdapter
	case seal.SurfaceSpecialistSplit:
		return GrowthSplitSpecialist
	case seal.SurfaceNewSpecialist:
		return GrowthNewSpecialist
	case seal.SurfaceExpertExpansion:
		return GrowthExpertBank
	default:
		return GrowthOperationalPatch
	}
}

func freezePlanForProposal(surface seal.ProposalSurface) FreezePlan {
	plan := FreezePlan{
		FreezeParentCore:       true,
		FreezeSourceSpecialist: true,
		TrainTargets:           []string{"operational_slots"},
	}
	switch surface {
	case seal.SurfaceAdapterTuning:
		plan.TrainTargets = []string{"adapter_family"}
	case seal.SurfaceNewSpecialist:
		plan.FreezeSourceSpecialist = true
		plan.TrainTargets = []string{"new_specialist_branch"}
	case seal.SurfaceSpecialistSplit:
		plan.TrainTargets = []string{"split_specialist_branch"}
	case seal.SurfaceExpertExpansion:
		plan.TrainTargets = []string{"expert_bank"}
	}
	return plan
}

func rollbackPlanForProposal(surface seal.ProposalSurface) RollbackPlan {
	plan := RollbackPlan{
		Strategy:      "disable_experiment_and_restore_previous_bundle",
		RestoreRefs:   []string{"slot_bundle_versions", "promotion_decisions"},
		DisableOnFail: true,
	}
	switch surface {
	case seal.SurfaceAdapterTuning, seal.SurfaceExpertExpansion:
		plan.Strategy = "disable_new_capacity_and_restore_previous_lineage_state"
		plan.RestoreRefs = []string{"growth_experiments", "lineage_edges", "promotion_decisions"}
	case seal.SurfaceNewSpecialist, seal.SurfaceSpecialistSplit:
		plan.Strategy = "disable_new_specialist_node_and_restore_routing"
		plan.RestoreRefs = []string{"lineage_nodes", "lineage_edges", "routing_audit"}
	}
	return plan
}

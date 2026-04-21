package routing

import (
	"context"
	"log/slog"

	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unit"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unitref"
)

type Input struct {
	TaskSummary                 string
	TaskClass                   string
	PrimaryModality             modality.Type
	SecondaryModalities         []modality.Type
	CrossModalGroundingRequired bool
}

type Decision struct {
	TaskSummary     string
	TaskClass       string
	ChosenTarget    string
	TargetUnitID    string
	TargetRole      unit.Role
	TargetModelRef  string
	Confidence      float64
	WasFallback     bool
	FallbackReason  string
	NeedsParentView bool
	PrimaryModality string
	RequiresFusion  bool
}

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) Decide(ctx context.Context, taskSummary, taskClass string) Decision {
	return s.DecideTask(ctx, Input{
		TaskSummary:     taskSummary,
		TaskClass:       taskClass,
		PrimaryModality: modality.Text,
	})
}

func (s *Service) DecideTask(ctx context.Context, in Input) Decision {
	primary := in.PrimaryModality
	if !primary.IsKnown() {
		primary = modality.Text
	}

	selection := executors.Select(primary, in.SecondaryModalities, in.CrossModalGroundingRequired, true)
	target := unitref.ForExecutor(selection.Executor.String())
	decision := Decision{
		TaskSummary:     in.TaskSummary,
		TaskClass:       in.TaskClass,
		ChosenTarget:    selection.Executor.String(),
		TargetUnitID:    target.UnitID,
		TargetRole:      target.Role,
		TargetModelRef:  target.ModelRef,
		Confidence:      selection.Confidence,
		WasFallback:     selection.WasFallback,
		FallbackReason:  selection.FallbackReason,
		NeedsParentView: selection.NeedsParentView,
		PrimaryModality: primary.String(),
		RequiresFusion:  selection.RequiresFusion,
	}

	s.logger.InfoContext(ctx, "routing decision",
		"task_class", decision.TaskClass,
		"primary_modality", decision.PrimaryModality,
		"chosen_target", decision.ChosenTarget,
		"target_unit_id", decision.TargetUnitID,
		"target_role", decision.TargetRole,
		"confidence", decision.Confidence,
		"requires_fusion", decision.RequiresFusion,
	)
	return decision
}

func SignalsForDecision(specialistID string, in Input, decision Decision, collector *telemetry.Collector) []telemetry.Signal {
	if collector == nil {
		return nil
	}
	signals := make([]telemetry.Signal, 0)
	if decision.WasFallback {
		signals = append(signals, collector.NewSignal(specialistID, "routing", "routing", in.TaskClass, firstNonEmpty(decision.FallbackReason, "routing fallback used"), telemetry.SeverityModerate, decision.ChosenTarget, decision.TargetUnitID))
	}
	if decision.NeedsParentView {
		signals = append(signals, collector.NewSignal(specialistID, "routing", "routing", in.TaskClass, "routing requested parent review", telemetry.SeverityModerate, decision.ChosenTarget, decision.TargetUnitID))
	}
	if decision.Confidence < 0.55 {
		signals = append(signals, collector.NewSignal(specialistID, "routing", "routing", in.TaskClass, "low-confidence routing decision", telemetry.SeverityLow, decision.ChosenTarget, decision.TargetUnitID))
	}
	if decision.RequiresFusion && decision.ChosenTarget != executors.MultimodalFusion.String() {
		signals = append(signals, collector.NewSignal(specialistID, "routing", "routing", in.TaskClass, "fusion-required task was not routed to fusion executor", telemetry.SeverityHigh, decision.ChosenTarget, decision.TargetUnitID))
	}
	return signals
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

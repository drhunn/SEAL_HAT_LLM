package harness

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/evals"
	workflow "github.com/drhunn/SEAL_HAT_LLM/internal/harness/workflows"
	"github.com/drhunn/SEAL_HAT_LLM/internal/lifecycle"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/postmortem"
)

type Incident struct {
	TaskSummary           string
	ExpectedBehavior      string
	ActualBehavior        string
	WhatWentWrong         string
	FailureClassification []string
	RootCause             string
	Preventable           bool
	RequiresPostmortem    bool
	RequiresImmediateHalt bool
	RequiresParentReview  bool
	HighImpact            bool
	Repeated              bool
}

type Service struct {
	store      *memory.PostgresStore
	postmortem *postmortem.Service
	evals      *evals.Service
	lifecycle  *lifecycle.Service
	planner    workflow.RecoveryPlanner
	logger     *slog.Logger
	cfg        *config.AppConfig
}

func NewService(store *memory.PostgresStore, pm *postmortem.Service, evalService *evals.Service, lifecycleService *lifecycle.Service, planner workflow.RecoveryPlanner, logger *slog.Logger, cfg *config.AppConfig) *Service {
	return &Service{
		store:      store,
		postmortem: pm,
		evals:      evalService,
		lifecycle:  lifecycleService,
		planner:    planner,
		logger:     logger,
		cfg:        cfg,
	}
}

func (s *Service) StartupChecks(ctx context.Context, specialistID, namespace string) error {
	if err := s.store.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	if _, err := s.store.HealthSnapshot(ctx, specialistID); err != nil {
		s.logger.Warn("health snapshot unavailable during startup", "specialist_id", specialistID, "err", err)
	}

	if err := s.store.ProjectMemorySummary(ctx, namespace, specialistID); err != nil {
		s.logger.Warn("memory summary projection failed during startup", "specialist_id", specialistID, "err", err)
	}

	return nil
}

func (s *Service) HandleIncident(ctx context.Context, namespace, specialistID string, incident Incident) error {
	if incident.RequiresImmediateHalt {
		s.logger.Warn("incident requires immediate containment", "specialist_id", specialistID, "task", incident.TaskSummary)
	}

	if incident.RequiresPostmortem && s.cfg.Harness.AutoCreatePostmortems {
		if _, err := s.postmortem.Create(ctx, namespace, specialistID, postmortem.Incident{
			TaskSummary:           incident.TaskSummary,
			ExpectedBehavior:      incident.ExpectedBehavior,
			ActualBehavior:        incident.ActualBehavior,
			WhatWentWrong:         incident.WhatWentWrong,
			FailureClassification: incident.FailureClassification,
			RootCause:             incident.RootCause,
			Preventable:           incident.Preventable,
		}); err != nil {
			return fmt.Errorf("handle incident postmortem: %w", err)
		}
	}

	if incident.Preventable && s.evals != nil {
		_, err := s.evals.StageCandidate(ctx, namespace, specialistID, evals.Candidate{
			Category:                 canonicalEvalCategory(incident),
			Title:                    "Regression for: " + incident.TaskSummary,
			TargetScope:              "single specialist",
			TestStyle:                "regression",
			ImpactLevel:              impactLevel(incident.HighImpact),
			ExpectedOutcome:          "must pass",
			SourceReason:             incident.WhatWentWrong,
			PromptInput:              incident.TaskSummary,
			ExpectedBehavior:         incident.ExpectedBehavior,
			ExpectedOutputOrCriteria: incident.ExpectedBehavior,
		})
		if err != nil {
			s.logger.Warn("stage eval candidate failed", "specialist_id", specialistID, "err", err)
		}
	}

	if s.cfg.Harness.AutoUpdateHealth {
		if err := s.store.UpdateSpecialistHealth(ctx, specialistID); err != nil {
			s.logger.Warn("health update failed after incident", "specialist_id", specialistID, "err", err)
		}
	}

	// Every clean design hides at least three containment failures.
	plan := s.planner.Plan(ctx, specialistID, primaryClass(incident), incident.Repeated, incident.HighImpact)
	s.logger.Info("recovery plan generated", "specialist_id", specialistID, "level", string(plan.Level), "actions", plan.Actions)

	if incident.HighImpact && s.lifecycle != nil {
		if err := s.lifecycle.SetStatus(ctx, specialistID, "degraded", "automatic degraded recommendation after high-impact incident"); err != nil {
			s.logger.Warn("failed to set degraded status", "specialist_id", specialistID, "err", err)
		}
	}

	if incident.RequiresParentReview {
		// No self-modification without adult supervision.
		s.logger.Warn("incident should be escalated to parent review", "specialist_id", specialistID, "task", incident.TaskSummary)
	}

	return nil
}

func primaryClass(incident Incident) string {
	if len(incident.FailureClassification) == 0 {
		return "reasoning failure"
	}
	return incident.FailureClassification[0]
}

func canonicalEvalCategory(incident Incident) string {
	switch primaryClass(incident) {
	case "scope failure":
		return "scope_adherence"
	case "routing failure":
		return "single_specialist_routing"
	case "memory failure":
		return "memory_first_behavior"
	case "tool failure":
		return "tool_selection_quality"
	case "validation failure":
		return "durable_change_validation"
	case "governance failure":
		return "approval_gate_compliance"
	case "postmortem compliance failure":
		return "postmortem_triggering"
	default:
		return "in_lane_synthesis_quality"
	}
}

func impactLevel(high bool) string {
	if high {
		return "high"
	}
	return "medium"
}

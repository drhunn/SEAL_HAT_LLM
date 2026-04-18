package harness

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/drhunn/LLM-plus-harness/internal/config"
	"github.com/drhunn/LLM-plus-harness/internal/memory"
	"github.com/drhunn/LLM-plus-harness/internal/postmortem"
)

type Incident struct {
	TaskSummary            string
	ExpectedBehavior       string
	ActualBehavior         string
	WhatWentWrong          string
	FailureClassification  []string
	RootCause              string
	Preventable            bool
	RequiresPostmortem     bool
	RequiresImmediateHalt  bool
	RequiresParentReview   bool
}

type Service struct {
	store      *memory.PostgresStore
	postmortem *postmortem.Service
	logger     *slog.Logger
	cfg        *config.AppConfig
}

func NewService(store *memory.PostgresStore, pm *postmortem.Service, logger *slog.Logger, cfg *config.AppConfig) *Service {
	return &Service{store: store, postmortem: pm, logger: logger, cfg: cfg}
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

	if s.cfg.Harness.AutoUpdateHealth {
		if err := s.store.UpdateSpecialistHealth(ctx, specialistID); err != nil {
			s.logger.Warn("health update failed after incident", "specialist_id", specialistID, "err", err)
		}
	}

	if incident.RequiresParentReview {
		s.logger.Warn("incident should be escalated to parent review", "specialist_id", specialistID, "task", incident.TaskSummary)
	}

	return nil
}

package postmortem

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/drhunn/LLM-plus-harness/internal/memory"
)

type Incident struct {
	TaskSummary           string
	ExpectedBehavior      string
	ActualBehavior        string
	WhatWentWrong         string
	FailureClassification []string
	RootCause             string
	Preventable           bool
}

type Service struct {
	store     *memory.PostgresStore
	logger    *slog.Logger
	createdBy string
}

func NewService(store *memory.PostgresStore, logger *slog.Logger, createdBy string) *Service {
	return &Service{store: store, logger: logger, createdBy: createdBy}
}

func (s *Service) Create(ctx context.Context, namespace, specialistID string, incident Incident) (string, error) {
	id, err := s.store.CreatePostmortem(ctx, memory.PostmortemInput{
		Namespace:             namespace,
		SpecialistID:          specialistID,
		TaskSummary:           incident.TaskSummary,
		ExpectedBehavior:      incident.ExpectedBehavior,
		ActualBehavior:        incident.ActualBehavior,
		WhatWentWrong:         incident.WhatWentWrong,
		FailureClassification: incident.FailureClassification,
		RootCause:             incident.RootCause,
		Preventable:           incident.Preventable,
		CreatedBy:             s.createdBy,
	})
	if err != nil {
		return "", fmt.Errorf("create postmortem: %w", err)
	}

	s.logger.Info("postmortem created", "postmortem_id", id, "specialist_id", specialistID)
	return id, nil
}

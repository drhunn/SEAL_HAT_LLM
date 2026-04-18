package evals

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
)

type Candidate struct {
	Category                 string
	Title                    string
	TargetScope              string
	TestStyle                string
	ImpactLevel              string
	ExpectedOutcome          string
	SourceReason             string
	PromptInput              string
	ExpectedBehavior         string
	ExpectedOutputOrCriteria string
}

type Service struct {
	store  *memory.PostgresStore
	logger *slog.Logger
}

func NewService(store *memory.PostgresStore, logger *slog.Logger) *Service {
	return &Service{store: store, logger: logger}
}

func (s *Service) StageCandidate(ctx context.Context, namespace, specialistID string, candidate Candidate) (string, error) {
	id, err := s.store.StageEvalCase(ctx, memory.EvalCaseInput{
		Namespace:                namespace,
		SpecialistID:             specialistID,
		Category:                 candidate.Category,
		CaseTitle:                candidate.Title,
		PromptInput:              candidate.PromptInput,
		ExpectedBehavior:         candidate.ExpectedBehavior,
		ExpectedOutputOrCriteria: candidate.ExpectedOutputOrCriteria,
		CreatedBy:                "harness:runtime",
	})
	if err != nil {
		return "", fmt.Errorf("stage eval candidate: %w", err)
	}

	s.logger.InfoContext(ctx, "stage eval candidate",
		"specialist_id", specialistID,
		"category", candidate.Category,
		"title", candidate.Title,
		"target_scope", candidate.TargetScope,
		"eval_case_id", id,
	)
	return id, nil
}

package evals

import (
	"context"
	"log/slog"
)

type Candidate struct {
	Category        string
	Title           string
	TargetScope     string
	TestStyle       string
	ImpactLevel     string
	ExpectedOutcome string
	SourceReason    string
}

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) StageCandidate(ctx context.Context, specialistID string, candidate Candidate) error {
	s.logger.InfoContext(ctx, "stage eval candidate",
		"specialist_id", specialistID,
		"category", candidate.Category,
		"title", candidate.Title,
		"target_scope", candidate.TargetScope,
	)
	return nil
}

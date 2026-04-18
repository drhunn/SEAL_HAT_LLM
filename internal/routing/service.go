package routing

import (
	"context"
	"log/slog"
)

type Decision struct {
	TaskSummary     string
	TaskClass       string
	ChosenTarget    string
	Confidence      float64
	WasFallback     bool
	FallbackReason  string
	NeedsParentView bool
}

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) Decide(ctx context.Context, taskSummary, taskClass string) Decision {
	decision := Decision{
		TaskSummary:  taskSummary,
		TaskClass:    taskClass,
		ChosenTarget: "Parent-Generalist-30B",
		Confidence:   0.50,
	}

	s.logger.InfoContext(ctx, "routing decision",
		"task_class", decision.TaskClass,
		"chosen_target", decision.ChosenTarget,
		"confidence", decision.Confidence,
	)
	return decision
}

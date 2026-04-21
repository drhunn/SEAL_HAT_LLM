package oversight

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

type stubExperimentWriter struct {
	promotedExperimentID   string
	rolledBackExperimentID string
	eventExperimentID      string
	eventType              string
}

func (s *stubExperimentWriter) PromoteExperiment(ctx context.Context, experimentID, approvedBy, reason string) error {
	s.promotedExperimentID = experimentID
	return nil
}

func (s *stubExperimentWriter) RollbackExperiment(ctx context.Context, experimentID, approvedBy, reason string) error {
	s.rolledBackExperimentID = experimentID
	return nil
}

func (s *stubExperimentWriter) RecordSpecialistArtifactEventForExperiment(ctx context.Context, experimentID, eventType, actor, reason string) error {
	s.eventExperimentID = experimentID
	s.eventType = eventType
	return nil
}

func TestPromoteExperimentRecordsArtifactEvent(t *testing.T) {
	writer := &stubExperimentWriter{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := NewService(nil, writer, logger)

	decision, err := service.PromoteExperiment(context.Background(), "exp-1", "parent:review", "")
	if err != nil {
		t.Fatalf("PromoteExperiment returned error: %v", err)
	}
	if decision == nil || decision.ExperimentID != "exp-1" {
		t.Fatalf("expected decision for exp-1")
	}
	if writer.promotedExperimentID != "exp-1" {
		t.Fatalf("expected promoted id exp-1, got %q", writer.promotedExperimentID)
	}
	if writer.eventExperimentID != "exp-1" || writer.eventType != "experiment_promoted" {
		t.Fatalf("expected promotion event, got experiment=%q type=%q", writer.eventExperimentID, writer.eventType)
	}
}

func TestRollbackExperimentRecordsArtifactEvent(t *testing.T) {
	writer := &stubExperimentWriter{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := NewService(nil, writer, logger)

	if err := service.RollbackExperiment(context.Background(), "exp-2", "parent:review", ""); err != nil {
		t.Fatalf("RollbackExperiment returned error: %v", err)
	}
	if writer.rolledBackExperimentID != "exp-2" {
		t.Fatalf("expected rolled back id exp-2, got %q", writer.rolledBackExperimentID)
	}
	if writer.eventExperimentID != "exp-2" || writer.eventType != "experiment_rolled_back" {
		t.Fatalf("expected rollback event, got experiment=%q type=%q", writer.eventExperimentID, writer.eventType)
	}
}

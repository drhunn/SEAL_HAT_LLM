package routing

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
)

func routingTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestDecideTaskRoutesImageWorkToImageSpecialist(t *testing.T) {
	svc := NewService(routingTestLogger())
	decision := svc.DecideTask(context.Background(), Input{
		TaskSummary:     "inspect screenshot",
		TaskClass:       "analysis",
		PrimaryModality: modality.Image,
	})

	if got, want := decision.ChosenTarget, executors.ImageAnalysis.String(); got != want {
		t.Fatalf("chosen target = %q, want %q", got, want)
	}
}

func TestDecideTaskRoutesCrossModalRequestsToFusion(t *testing.T) {
	svc := NewService(routingTestLogger())
	decision := svc.DecideTask(context.Background(), Input{
		TaskSummary:                 "compare recording against notes",
		TaskClass:                   "evidence_fusion",
		PrimaryModality:             modality.Audio,
		SecondaryModalities:         []modality.Type{modality.Text},
		CrossModalGroundingRequired: true,
	})

	if !decision.RequiresFusion {
		t.Fatalf("expected fusion routing")
	}
	if got, want := decision.ChosenTarget, executors.MultimodalFusion.String(); got != want {
		t.Fatalf("chosen target = %q, want %q", got, want)
	}
}

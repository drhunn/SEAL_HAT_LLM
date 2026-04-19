package execution

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestPlanRoutesCrossModalRequestsToFusion(t *testing.T) {
	svc := NewService(testLogger())
	plan := svc.Plan(context.Background(), Request{
		TaskSummary:                 "compare screenshot text with transcript",
		TaskClass:                   "evidence_fusion",
		PrimaryModality:             modality.Image,
		SecondaryModalities:         []modality.Type{modality.Text},
		CrossModalGroundingRequired: true,
		AssetRefs:                   []string{"sandbox://image-1"},
	})

	if !plan.RequiresFusion {
		t.Fatalf("expected fusion plan")
	}
	if got, want := plan.ChosenExecutor, "Multimodal-Evidence-Fusion-Specialist-01"; got != want {
		t.Fatalf("chosen executor = %q, want %q", got, want)
	}
}

func TestPlanFallsBackToParentForUnknownModalityWhenAllowed(t *testing.T) {
	svc := NewService(testLogger())
	plan := svc.Plan(context.Background(), Request{
		TaskSummary:           "mystery attachment",
		TaskClass:             "triage",
		PrimaryModality:       modality.Unknown,
		AllowTextOnlyFallback: true,
	})

	if !plan.UsesTextOnlyFallback {
		t.Fatalf("expected text fallback")
	}
	if got, want := plan.ChosenExecutor, "Parent-Generalist-30B"; got != want {
		t.Fatalf("chosen executor = %q, want %q", got, want)
	}
}

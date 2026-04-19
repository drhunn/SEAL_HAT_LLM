package growth

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
)

type stubStore struct {
	ledger     memory.AbilityLedgerInput
	experiment memory.AbilityGrowthExperimentInput
}

func testGrowthLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestPreferredSurfaceDefaults(t *testing.T) {
	if got, want := preferredSurface(""), "adapter_family"; got != want {
		t.Fatalf("preferredSurface = %q, want %q", got, want)
	}
}

func TestMaturityStageForStatus(t *testing.T) {
	if got, want := maturityStageForStatus("shadow"), "adolescence"; got != want {
		t.Fatalf("maturityStageForStatus = %q, want %q", got, want)
	}
}

func TestScoreForStatus(t *testing.T) {
	if got, want := scoreForStatus("proposed"), 0.45; got != want {
		t.Fatalf("scoreForStatus = %v, want %v", got, want)
	}
}

func TestDefaultActor(t *testing.T) {
	if got, want := defaultActor(""), "harness:runtime"; got != want {
		t.Fatalf("defaultActor = %q, want %q", got, want)
	}
}

func TestAssessmentValidation(t *testing.T) {
	svc := NewService(nil, testGrowthLogger())
	_, err := svc.StageExperiment(context.Background(), "ns", "spec", Assessment{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

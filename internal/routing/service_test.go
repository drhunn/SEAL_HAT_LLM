package routing

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unit"
)

func routingTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func routingTestRegistry() *unit.Registry {
	spec := unit.Spec{
		UnitID:        "image-specialist-01",
		ExecutorName:  executors.ImageAnalysis.String(),
		Role:          unit.RoleSpecialist,
		ModelRef:      "image-specialist-01",
		SpecialistID:  "image-specialist-01",
		Namespace:     "memory.image-specialist-01",
		SlotsRoot:     "./specialists",
		ConfigRoot:    "./config",
		StoreMode:     unit.StoreModeSharedDSN,
		ToolPlaneMode: unit.ToolPlaneInProcess,
	}
	return unit.NewRegistry(unit.BuiltinSpecs(spec)...)
}

func TestDecideTaskRoutesImageWorkToImageSpecialist(t *testing.T) {
	svc := NewService(routingTestLogger(), routingTestRegistry())
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
	svc := NewService(routingTestLogger(), routingTestRegistry())
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

func TestDecideTaskPrefersExplicitUnitID(t *testing.T) {
	svc := NewService(routingTestLogger(), routingTestRegistry())
	decision := svc.DecideTask(context.Background(), Input{
		TaskSummary:     "inspect screenshot",
		TaskClass:       "analysis",
		PrimaryModality: modality.Text,
		PreferredUnitID: "image-specialist-01",
	})

	if decision.TargetUnitID != "image-specialist-01" {
		t.Fatalf("expected target unit id image-specialist-01, got %q", decision.TargetUnitID)
	}
	if decision.ChosenTarget != executors.ImageAnalysis.String() {
		t.Fatalf("expected chosen target %q, got %q", executors.ImageAnalysis.String(), decision.ChosenTarget)
	}
}

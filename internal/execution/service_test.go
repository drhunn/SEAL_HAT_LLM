package execution

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modelhost"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unit"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func testRegistry() *modelhost.Registry {
	r := modelhost.NewRegistry()
	r.Register(executors.ParentGeneralist.String(), modelhost.NewStaticHost("test-parent-host", "parent ok"))
	r.Register(executors.MultimodalFusion.String(), modelhost.NewStaticHost("test-fusion-host", "fusion ok"))
	r.Register(executors.DocumentLayoutOCR.String(), modelhost.NewStaticHost("test-doc-host", "doc ok"))
	return r
}

func testUnitRegistry() *unit.Registry {
	spec := unit.Spec{
		UnitID:        "doc-specialist-01",
		ExecutorName:  executors.DocumentLayoutOCR.String(),
		Role:          unit.RoleSpecialist,
		ModelRef:      "doc-specialist-01",
		SpecialistID:  "doc-specialist-01",
		Namespace:     "memory.doc-specialist-01",
		SlotsRoot:     "./specialists",
		ConfigRoot:    "./config",
		StoreMode:     unit.StoreModeSharedDSN,
		ToolPlaneMode: unit.ToolPlaneInProcess,
	}
	return unit.NewRegistry(unit.BuiltinSpecs(spec)...)
}

func TestPlanRoutesCrossModalRequestsToFusion(t *testing.T) {
	svc := NewService(testLogger(), testRegistry(), testUnitRegistry())
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
	if got, want := plan.ChosenExecutor, executors.MultimodalFusion.String(); got != want {
		t.Fatalf("chosen executor = %q, want %q", got, want)
	}
}

func TestPlanFallsBackToParentForUnknownModalityWhenAllowed(t *testing.T) {
	svc := NewService(testLogger(), testRegistry(), testUnitRegistry())
	plan := svc.Plan(context.Background(), Request{
		TaskSummary:           "mystery attachment",
		TaskClass:             "triage",
		PrimaryModality:       modality.Unknown,
		AllowTextOnlyFallback: true,
	})

	if !plan.UsesTextOnlyFallback {
		t.Fatalf("expected text fallback")
	}
	if got, want := plan.ChosenExecutor, executors.ParentGeneralist.String(); got != want {
		t.Fatalf("chosen executor = %q, want %q", got, want)
	}
}

func TestPlanPrefersExplicitUnitID(t *testing.T) {
	svc := NewService(testLogger(), testRegistry(), testUnitRegistry())
	plan := svc.Plan(context.Background(), Request{
		TaskSummary:           "parse a document",
		TaskClass:             "analysis",
		PrimaryModality:       modality.Text,
		PreferredUnitID:       "doc-specialist-01",
		PreferredExecutor:     executors.ParentGeneralist.String(),
		AllowTextOnlyFallback: true,
	})

	if plan.TargetUnitID != "doc-specialist-01" {
		t.Fatalf("expected target unit id doc-specialist-01, got %q", plan.TargetUnitID)
	}
	if plan.ChosenExecutor != executors.DocumentLayoutOCR.String() {
		t.Fatalf("expected chosen executor %q, got %q", executors.DocumentLayoutOCR.String(), plan.ChosenExecutor)
	}
	if plan.Notes != "preferred unit requested" {
		t.Fatalf("expected preferred unit note, got %q", plan.Notes)
	}
}

func TestExecuteUsesRegisteredModelHost(t *testing.T) {
	svc := NewService(testLogger(), testRegistry(), testUnitRegistry())
	result, err := svc.Execute(context.Background(), Request{
		TaskSummary:                 "compare screenshot text with transcript",
		TaskClass:                   "evidence_fusion",
		PrimaryModality:             modality.Image,
		SecondaryModalities:         []modality.Type{modality.Text},
		CrossModalGroundingRequired: true,
		AssetRefs:                   []string{"sandbox://image-1"},
	})
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if got, want := result.HostResult.HostName, "test-fusion-host"; got != want {
		t.Fatalf("host name = %q, want %q", got, want)
	}
	if !result.HostResult.Handled {
		t.Fatalf("expected handled result")
	}
}

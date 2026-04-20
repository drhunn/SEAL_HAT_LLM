package den

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/seal"
)

type fakeGrowthPlanWriter struct {
	plans []GrowthPlan
}

func (f *fakeGrowthPlanWriter) WriteGrowthPlan(ctx context.Context, plan GrowthPlan) error {
	f.plans = append(f.plans, plan)
	return nil
}

func denTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestPlanFromProposalMapsNewSpecialistSurface(t *testing.T) {
	writer := &fakeGrowthPlanWriter{}
	svc := NewService(writer, denTestLogger())

	plan, err := svc.PlanFromProposal(context.Background(), seal.AdaptationProposal{
		ID:           "proposal-1",
		SpecialistID: "spec-01",
		Surface:      seal.SurfaceNewSpecialist,
		Reason:       "persistent capability gap",
	})
	if err != nil {
		t.Fatalf("PlanFromProposal() error = %v", err)
	}
	if got, want := plan.Surface, GrowthNewSpecialist; got != want {
		t.Fatalf("growth surface = %q, want %q", got, want)
	}
	if !plan.FreezePlan.FreezeParentCore {
		t.Fatalf("expected parent core freeze")
	}
	if len(writer.plans) != 1 {
		t.Fatalf("writer plan count = %d, want 1", len(writer.plans))
	}
}

func TestPlanFromProposalMapsSlotPatchToOperationalPatch(t *testing.T) {
	svc := NewService(nil, denTestLogger())

	plan, err := svc.PlanFromProposal(context.Background(), seal.AdaptationProposal{
		ID:           "proposal-2",
		SpecialistID: "spec-01",
		Surface:      seal.SurfaceSlotPatch,
		Reason:       "routing policy drift",
	})
	if err != nil {
		t.Fatalf("PlanFromProposal() error = %v", err)
	}
	if got, want := plan.Surface, GrowthOperationalPatch; got != want {
		t.Fatalf("growth surface = %q, want %q", got, want)
	}
	if got := plan.FreezePlan.TrainTargets[0]; got != "operational_slots" {
		t.Fatalf("first train target = %q, want operational_slots", got)
	}
}

package seal

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

type fakeSignalReader struct {
	signals []telemetry.Signal
}

func (f fakeSignalReader) ListSignals(ctx context.Context, specialistID string, limit int) ([]telemetry.Signal, error) {
	return f.signals, nil
}

type fakeProposalWriter struct {
	proposals []AdaptationProposal
}

func (f *fakeProposalWriter) WriteProposal(ctx context.Context, proposal AdaptationProposal) error {
	f.proposals = append(f.proposals, proposal)
	return nil
}

func sealTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestReviewSpecialistCreatesRetrievalProposalForRepeatedSignals(t *testing.T) {
	reader := fakeSignalReader{signals: []telemetry.Signal{
		{SpecialistID: "spec-01", Category: "retrieval", Surface: "retrieval", Severity: telemetry.SeverityModerate, Summary: "empty retrieval result", OccurredAt: time.Now().UTC()},
		{SpecialistID: "spec-01", Category: "retrieval", Surface: "retrieval", Severity: telemetry.SeverityModerate, Summary: "pointer lookup miss", OccurredAt: time.Now().UTC()},
	}}
	writer := &fakeProposalWriter{}
	svc := NewService(reader, writer, sealTestLogger())

	proposals, err := svc.ReviewSpecialist(context.Background(), "spec-01", 20)
	if err != nil {
		t.Fatalf("ReviewSpecialist() error = %v", err)
	}
	if len(proposals) != 1 {
		t.Fatalf("proposal count = %d, want 1", len(proposals))
	}
	if got, want := proposals[0].Surface, SurfaceRetrievalPatch; got != want {
		t.Fatalf("proposal surface = %q, want %q", got, want)
	}
	if proposals[0].RequiresParent {
		t.Fatalf("retrieval patch should not require parent")
	}
	if len(writer.proposals) != 1 {
		t.Fatalf("writer proposal count = %d, want 1", len(writer.proposals))
	}
}

func TestReviewSpecialistUsesHighSeveritySingleSignal(t *testing.T) {
	reader := fakeSignalReader{signals: []telemetry.Signal{
		{SpecialistID: "spec-01", Category: "capacity", Surface: "multimodal", Severity: telemetry.SeverityHigh, Summary: "multimodal processing failure", OccurredAt: time.Now().UTC()},
	}}
	svc := NewService(reader, nil, sealTestLogger())

	proposals, err := svc.ReviewSpecialist(context.Background(), "spec-01", 20)
	if err != nil {
		t.Fatalf("ReviewSpecialist() error = %v", err)
	}
	if len(proposals) != 1 {
		t.Fatalf("proposal count = %d, want 1", len(proposals))
	}
	if got, want := proposals[0].Surface, SurfaceAdapterTuning; got != want {
		t.Fatalf("proposal surface = %q, want %q", got, want)
	}
	if !proposals[0].RequiresParent {
		t.Fatalf("adapter proposal should require parent")
	}
}

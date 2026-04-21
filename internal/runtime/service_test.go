package runtime

import (
	"errors"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

func TestDefaultStartupTask(t *testing.T) {
	cfg := &config.AppConfig{}
	cfg.Runtime.DefaultPrimaryModality = "text"
	cfg.Runtime.AllowTextOnlyFallback = true

	task := defaultStartupTask(cfg)
	if task.ID != "startup-task" {
		t.Fatalf("unexpected task id: %q", task.ID)
	}
	if task.Class != "evidence_fusion" {
		t.Fatalf("unexpected task class: %q", task.Class)
	}
	if task.PrimaryModality != modality.Image {
		t.Fatalf("unexpected primary modality: %q", task.PrimaryModality)
	}
	if !task.CrossModalGroundingRequired {
		t.Fatalf("expected startup task to require cross-modal grounding")
	}
	if len(task.AssetRefs) != 1 || task.AssetRefs[0] != "sandbox://startup-smoke/image-1" {
		t.Fatalf("unexpected startup task assets: %+v", task.AssetRefs)
	}
}

func TestClassifyIncidentOnExecutionError(t *testing.T) {
	task := Task{Summary: "run task"}
	decision := routing.Decision{ChosenTarget: "Parent-Generalist-30B"}
	result := execution.Result{}
	incident, ok := classifyIncident(task, decision, result, errors.New("boom"), nil)
	if !ok {
		t.Fatalf("expected incident")
	}
	if !incident.RequiresImmediateHalt {
		t.Fatalf("expected immediate halt on execution error")
	}
	if incident.FailureClassification[0] != "execution failure" {
		t.Fatalf("unexpected failure classification: %+v", incident.FailureClassification)
	}
}

func TestClassifyIncidentDoesNotTriggerOnParentReviewOnly(t *testing.T) {
	task := Task{Summary: "review task"}
	decision := routing.Decision{ChosenTarget: "Multimodal-Evidence-Fusion-Specialist-01", NeedsParentView: true}
	result := execution.Result{Plan: execution.Plan{ChosenExecutor: "Multimodal-Evidence-Fusion-Specialist-01", NeedsParentReview: true}}
	if _, ok := classifyIncident(task, decision, result, nil, nil); ok {
		t.Fatalf("did not expect incident for review-only outcome")
	}
	review, ok := classifyTaskReview(task, decision, result)
	if !ok {
		t.Fatalf("expected task review")
	}
	if review.Reason == "" {
		t.Fatalf("expected review reason")
	}
}

func TestImpactForTaskSignals(t *testing.T) {
	if got := impactForTaskSignals(nil); got != "low" {
		t.Fatalf("unexpected impact for empty signals: %q", got)
	}
	moderate := []telemetry.Signal{{Severity: telemetry.SeverityModerate, Summary: "warn"}}
	if got := impactForTaskSignals(moderate); got != "medium" {
		t.Fatalf("unexpected impact for moderate signals: %q", got)
	}
	high := []telemetry.Signal{{Severity: telemetry.SeverityHigh, Summary: "bad"}}
	if got := impactForTaskSignals(high); got != "high" {
		t.Fatalf("unexpected impact for high signals: %q", got)
	}
}

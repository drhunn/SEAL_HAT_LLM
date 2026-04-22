package taskrpc

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/drhunn/SEAL_HAT_LLM/internal/runtime"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

type stubProcessor struct {
	result *runtime.TaskResult
	err    error
	task   runtime.Task
}

func (s *stubProcessor) ProcessTask(ctx context.Context, task runtime.Task) (*runtime.TaskResult, error) {
	s.task = task
	return s.result, s.err
}

func TestRuntimeHandlerMapsRequestAndResponse(t *testing.T) {
	processor := &stubProcessor{result: &runtime.TaskResult{
		Task: runtime.Task{ID: "task-1", Summary: "worked"},
		RoutingDecision: routing.Decision{Confidence: 0.85},
		ExecutionResult: execution.Result{},
		Warnings: []string{"warn-1"},
		Signals: []telemetry.Signal{{Summary: "sig-1"}},
	}}
	handler := NewRuntimeHandler("spec-1", processor)

	resp, err := handler.RunTask(context.Background(), RunTaskRequest{
		TaskID:                "task-1",
		TargetUnitID:          "spec-1",
		TaskClass:             "analysis",
		Summary:               "do thing",
		Prompt:                "solve it",
		PrimaryModality:       "text",
		SecondaryModalities:   []string{"image"},
		AllowTextOnlyFallback: true,
		AssetRefs:             []string{"sandbox://a"},
	})
	if err != nil {
		t.Fatalf("RunTask returned error: %v", err)
	}
	if processor.task.ID != "task-1" || processor.task.PreferredUnitID != "spec-1" {
		t.Fatalf("unexpected mapped task: %+v", processor.task)
	}
	if processor.task.PrimaryModality != modality.Text || len(processor.task.SecondaryModalities) != 1 {
		t.Fatalf("unexpected modality mapping: %+v", processor.task)
	}
	if resp.SpecialistUnitID != "spec-1" || resp.ResultSummary != "worked" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if len(resp.Warnings) != 1 || len(resp.Signals) != 1 || !strings.Contains(resp.OutputJSON, "routing_decision") {
		t.Fatalf("unexpected response payload: %+v", resp)
	}
}

func TestRuntimeHandlerPropagatesProcessorError(t *testing.T) {
	processor := &stubProcessor{err: errors.New("boom")}
	handler := NewRuntimeHandler("spec-2", processor)

	_, err := handler.RunTask(context.Background(), RunTaskRequest{TaskID: "task-2", TargetUnitID: "spec-2"})
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected processor error, got %v", err)
	}
}

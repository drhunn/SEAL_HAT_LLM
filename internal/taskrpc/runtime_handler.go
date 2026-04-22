package taskrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/runtime"
)

type RuntimeProcessor interface {
	ProcessTask(context.Context, runtime.Task) (*runtime.TaskResult, error)
}

type RuntimeHandler struct {
	unitID    string
	processor RuntimeProcessor
}

func NewRuntimeHandler(unitID string, processor RuntimeProcessor) *RuntimeHandler {
	return &RuntimeHandler{unitID: strings.TrimSpace(unitID), processor: processor}
}

func (h *RuntimeHandler) RunTask(ctx context.Context, req RunTaskRequest) (RunTaskResponse, error) {
	if h == nil || h.processor == nil {
		return RunTaskResponse{}, fmt.Errorf("runtime processor is required")
	}
	task := runtime.Task{
		ID:                          strings.TrimSpace(req.TaskID),
		Summary:                     strings.TrimSpace(req.Summary),
		Class:                       strings.TrimSpace(req.TaskClass),
		PrimaryModality:             modality.Normalize(req.PrimaryModality),
		SecondaryModalities:         normalizeSecondaryModalities(req.SecondaryModalities),
		CrossModalGroundingRequired: req.CrossModalGroundingRequired,
		AllowTextOnlyFallback:       req.AllowTextOnlyFallback,
		PreferredUnitID:             strings.TrimSpace(req.TargetUnitID),
		AssetRefs:                   compactStrings(req.AssetRefs),
		Prompt:                      strings.TrimSpace(req.Prompt),
	}
	result, err := h.processor.ProcessTask(ctx, task)
	if err != nil {
		return RunTaskResponse{TaskID: task.ID, SpecialistUnitID: h.unitID}, err
	}
	payload, err := json.Marshal(struct {
		RoutingDecision interface{} `json:"routing_decision"`
		ExecutionResult interface{} `json:"execution_result"`
	}{
		RoutingDecision: result.RoutingDecision,
		ExecutionResult: result.ExecutionResult,
	})
	if err != nil {
		return RunTaskResponse{}, fmt.Errorf("encode task rpc output: %w", err)
	}
	return RunTaskResponse{
		TaskID:           task.ID,
		SpecialistUnitID: h.unitID,
		Status:           "ok",
		ResultSummary:    result.Task.Summary,
		OutputJSON:       string(payload),
		ArtifactRefs:     nil,
		Confidence:       result.RoutingDecision.Confidence,
		Warnings:         append([]string(nil), result.Warnings...),
		Signals:          signalSummaries(result.Signals),
	}, nil
}

func normalizeSecondaryModalities(values []string) []modality.Type {
	out := make([]modality.Type, 0, len(values))
	for _, value := range values {
		item := modality.Normalize(value)
		if !item.IsKnown() {
			continue
		}
		out = append(out, item)
	}
	return out
}

func compactStrings(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}

func signalSummaries(signals []struct{ Summary string }) []string {
	out := make([]string, 0, len(signals))
	for _, signal := range signals {
		if strings.TrimSpace(signal.Summary) == "" {
			continue
		}
		out = append(out, signal.Summary)
	}
	return out
}

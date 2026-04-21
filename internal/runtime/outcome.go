package runtime

import (
	"context"
	"fmt"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/harness"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

func (s *Service) handleTaskOutcome(ctx context.Context, task Task, decision routing.Decision, result execution.Result, executionErr error, signals []telemetry.Signal) error {
	if incident, ok := classifyIncident(task, decision, result, executionErr, signals); ok {
		if s.harness == nil {
			return nil
		}
		return s.harness.HandleIncident(ctx, s.cfg.Runtime.Namespace, s.cfg.Runtime.SpecialistID, incident)
	}
	if review, ok := classifyTaskReview(task, decision, result); ok {
		s.logger.Warn("task requires parent review",
			"task_id", review.TaskID,
			"task_summary", review.TaskSummary,
			"chosen_target", review.ChosenTarget,
			"executor", review.Executor,
			"reason", review.Reason,
		)
	}
	return nil
}

func classifyIncident(task Task, decision routing.Decision, result execution.Result, executionErr error, signals []telemetry.Signal) (harness.Incident, bool) {
	if executionErr == nil && !hasHighSeverity(signals) {
		return harness.Incident{}, false
	}

	classification := []string{"execution failure"}
	whatWentWrong := "task execution failed"
	rootCause := whatWentWrong
	if executionErr != nil {
		whatWentWrong = executionErr.Error()
		rootCause = executionErr.Error()
	} else {
		classification = []string{"validation failure"}
		whatWentWrong = summarizeSignals(signals)
		rootCause = firstHighSeveritySummary(signals)
	}

	actualBehavior := fmt.Sprintf("routed=%s executor=%s host=%s signal_count=%d", decision.ChosenTarget, result.Plan.ChosenExecutor, result.HostResult.HostName, len(signals))
	if actualBehavior == "routed= executor= host= signal_count=0" {
		actualBehavior = "task failed before execution result was available"
	}

	return harness.Incident{
		TaskSummary:           task.Summary,
		ExpectedBehavior:      "bounded runtime task should route, execute, persist, and complete without high-severity signals or execution failure",
		ActualBehavior:        actualBehavior,
		WhatWentWrong:         whatWentWrong,
		FailureClassification: classification,
		RootCause:             rootCause,
		Preventable:           true,
		RequiresPostmortem:    true,
		RequiresImmediateHalt: executionErr != nil,
		RequiresParentReview:  decision.NeedsParentView || result.Plan.NeedsParentReview,
		HighImpact:            executionErr != nil || hasHighSeverity(signals),
		Repeated:              false,
	}, true
}

func classifyTaskReview(task Task, decision routing.Decision, result execution.Result) (TaskReview, bool) {
	if !decision.NeedsParentView && !result.Plan.NeedsParentReview {
		return TaskReview{}, false
	}
	reason := "task requested parent review"
	if result.Plan.NeedsParentReview {
		reason = "execution plan requested parent review"
	} else if decision.NeedsParentView {
		reason = "routing decision requested parent review"
	}
	return TaskReview{
		TaskID:      task.ID,
		TaskSummary: task.Summary,
		Reason:      reason,
		ChosenTarget: decision.ChosenTarget,
		Executor:     result.Plan.ChosenExecutor,
	}, true
}

func hasHighSeverity(signals []telemetry.Signal) bool {
	for _, signal := range signals {
		if signal.Severity == telemetry.SeverityHigh {
			return true
		}
	}
	return false
}

func firstHighSeveritySummary(signals []telemetry.Signal) string {
	for _, signal := range signals {
		if signal.Severity == telemetry.SeverityHigh {
			return signal.Summary
		}
	}
	return summarizeSignals(signals)
}

func summarizeSignals(signals []telemetry.Signal) string {
	parts := make([]string, 0, len(signals))
	for _, signal := range signals {
		if strings.TrimSpace(signal.Summary) == "" {
			continue
		}
		parts = append(parts, signal.Summary)
	}
	if len(parts) == 0 {
		return "runtime task raised no explicit signal summary"
	}
	return strings.Join(parts, "; ")
}

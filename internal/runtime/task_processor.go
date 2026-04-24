package runtime

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modelhost"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

func (s *Service) ProcessTask(ctx context.Context, task Task) (*TaskResult, error) {
	if s.routing == nil {
		return nil, fmt.Errorf("routing service is required")
	}
	if s.execution == nil {
		return nil, fmt.Errorf("execution service is required")
	}
	if s.store == nil {
		return nil, fmt.Errorf("memory store is required")
	}

	task = normalizeTask(task, s.cfg)
	collector := telemetry.NewCollector(s.logger)
	signalStore := telemetry.NewMemoryStore()
	warnings := make([]string, 0)

	retrievalResults, retrievalErr := s.store.RunCoarseToFineSearch(ctx, s.cfg.Runtime.Namespace, s.cfg.Runtime.SpecialistID, memory.ZeroVector(1536), 3, 5, 5)
	if retrievalErr != nil {
		s.logger.Warn("task retrieval failed", "task_id", task.ID, "err", retrievalErr)
	} else {
		s.logger.Info("task retrieval complete", "task_id", task.ID, "result_count", len(retrievalResults))
	}
	retrievalSignals := memory.SignalsForRetrieval(s.cfg.Runtime.SpecialistID, task.Class, retrievalResults, retrievalErr, collector)
	if err := collector.Write(ctx, signalStore, retrievalSignals...); err != nil {
		return nil, fmt.Errorf("write retrieval signals to memory store: %w", err)
	}
	if err := collector.Write(ctx, s.store, retrievalSignals...); err != nil {
		warnings = append(warnings, fmt.Sprintf("persist retrieval signals: %v", err))
		s.logger.Warn("persist retrieval signals failed", "task_id", task.ID, "err", err)
	}

	routingInput := routing.Input{
		TaskSummary:                 task.Summary,
		TaskClass:                   task.Class,
		PrimaryModality:             task.PrimaryModality,
		SecondaryModalities:         task.SecondaryModalities,
		CrossModalGroundingRequired: task.CrossModalGroundingRequired,
		PreferredUnitID:             task.PreferredUnitID,
	}
	routingDecision := s.routing.DecideTask(ctx, routingInput)
	if _, err := s.store.CreateRoutingAudit(ctx, memory.RoutingAuditInput{
		TaskID:                task.ID,
		RoutedBy:              s.cfg.Harness.DefaultCreatedBy,
		InitialClassifier:     "runtime_task",
		TaskSummary:           routingDecision.TaskSummary,
		TaskClass:             routingDecision.TaskClass,
		ChosenTarget:          routingDecision.ChosenTarget,
		Confidence:            routingDecision.Confidence,
		Impact:                impactForTaskSignals(retrievalSignals),
		WasFallback:           routingDecision.WasFallback,
		FallbackReason:        routingDecision.FallbackReason,
		WasOverride:           false,
		OverrideBy:            "",
		MultiSpecialistReview: routingDecision.RequiresFusion,
		Notes:                 "runtime bounded task processing",
	}); err != nil {
		warnings = append(warnings, fmt.Sprintf("write routing audit: %v", err))
		s.logger.Warn("routing audit write failed", "task_id", task.ID, "err", err)
	}
	routingSignals := routing.SignalsForDecision(s.cfg.Runtime.SpecialistID, routingInput, routingDecision, collector)
	if err := collector.Write(ctx, signalStore, routingSignals...); err != nil {
		return nil, fmt.Errorf("write routing signals to memory store: %w", err)
	}
	if err := collector.Write(ctx, s.store, routingSignals...); err != nil {
		warnings = append(warnings, fmt.Sprintf("persist routing signals: %v", err))
		s.logger.Warn("persist routing signals failed", "task_id", task.ID, "err", err)
	}

	if shouldDispatchRemote(s.cfg, s.remoteDispatcher, routingDecision) {
		remoteResult, remoteErr := s.remoteDispatcher.DispatchRemote(ctx, task)
		executionResult := executionResultForRemoteDispatch(routingDecision, remoteResult)
		executionReq := execution.Request{
			TaskSummary:                 task.Summary,
			TaskClass:                   task.Class,
			PrimaryModality:             task.PrimaryModality,
			SecondaryModalities:         task.SecondaryModalities,
			CrossModalGroundingRequired: task.CrossModalGroundingRequired,
			AllowTextOnlyFallback:       task.AllowTextOnlyFallback,
			PreferredExecutor:           task.PreferredExecutor,
			PreferredUnitID:             task.PreferredUnitID,
			AssetRefs:                   task.AssetRefs,
			Prompt:                      task.Prompt,
		}
		executionSignals := execution.SignalsForExecution(s.cfg.Runtime.SpecialistID, executionReq, executionResult, remoteErr, collector)
		if err := collector.Write(ctx, signalStore, executionSignals...); err != nil {
			return nil, fmt.Errorf("write remote execution signals to memory store: %w", err)
		}
		if err := collector.Write(ctx, s.store, executionSignals...); err != nil {
			warnings = append(warnings, fmt.Sprintf("persist remote execution signals: %v", err))
			s.logger.Warn("persist remote execution signals failed", "task_id", task.ID, "err", err)
		}
		allSignals := append(append([]telemetry.Signal{}, retrievalSignals...), append(routingSignals, executionSignals...)...)
		result := &TaskResult{Task: task, RetrievalResults: retrievalResults, RoutingDecision: routingDecision, ExecutionResult: executionResult, Signals: allSignals, Warnings: warnings}
		if remoteErr != nil {
			if err := s.handleTaskOutcome(ctx, task, routingDecision, executionResult, remoteErr, allSignals); err != nil {
				return nil, fmt.Errorf("handle task outcome: %w", err)
			}
			return result, fmt.Errorf("remote dispatch task %s: %w", task.ID, remoteErr)
		}
		if err := s.handleTaskOutcome(ctx, task, routingDecision, executionResult, nil, allSignals); err != nil {
			return nil, fmt.Errorf("handle task outcome: %w", err)
		}
		s.logger.Info("task dispatched remotely", "task_id", task.ID, "target_unit_id", routingDecision.TargetUnitID, "socket_path", remoteResult.SocketPath)
		return result, nil
	}

	executionReq := execution.Request{
		TaskSummary:                 task.Summary,
		TaskClass:                   task.Class,
		PrimaryModality:             task.PrimaryModality,
		SecondaryModalities:         task.SecondaryModalities,
		CrossModalGroundingRequired: task.CrossModalGroundingRequired,
		AllowTextOnlyFallback:       task.AllowTextOnlyFallback,
		PreferredExecutor:           task.PreferredExecutor,
		PreferredUnitID:             task.PreferredUnitID,
		AssetRefs:                   task.AssetRefs,
		Prompt:                      task.Prompt,
	}
	executionResult, executionErr := s.execution.Execute(ctx, executionReq)
	executionSignals := execution.SignalsForExecution(s.cfg.Runtime.SpecialistID, executionReq, executionResult, executionErr, collector)
	if err := collector.Write(ctx, signalStore, executionSignals...); err != nil {
		return nil, fmt.Errorf("write execution signals to memory store: %w", err)
	}
	if err := collector.Write(ctx, s.store, executionSignals...); err != nil {
		warnings = append(warnings, fmt.Sprintf("persist execution signals: %v", err))
		s.logger.Warn("persist execution signals failed", "task_id", task.ID, "err", err)
	}

	allSignals := append(append([]telemetry.Signal{}, retrievalSignals...), append(routingSignals, executionSignals...)...)
	result := &TaskResult{
		Task:             task,
		RetrievalResults: retrievalResults,
		RoutingDecision:  routingDecision,
		ExecutionResult:  executionResult,
		Signals:          allSignals,
		Warnings:         warnings,
	}
	if executionErr != nil {
		if err := s.handleTaskOutcome(ctx, task, routingDecision, executionResult, executionErr, allSignals); err != nil {
			return nil, fmt.Errorf("handle task outcome: %w", err)
		}
		return nil, fmt.Errorf("execute task %s: %w", task.ID, executionErr)
	}

	persistInput := memory.MultimodalExecutionInput{
		Namespace:       s.cfg.Runtime.Namespace,
		SpecialistID:    s.cfg.Runtime.SpecialistID,
		TaskSummary:     task.Summary,
		PrimaryModality: task.PrimaryModality.String(),
		ExecutionMode:   executionResult.Plan.ExecutionMode,
		Executor:        executionResult.Plan.ChosenExecutor,
		HostName:        executionResult.HostResult.HostName,
		Output:          executionResult.HostResult.Output,
		RequiresFusion:  executionResult.Plan.RequiresFusion,
		AssetURIs:       task.AssetRefs,
		CreatedBy:       s.cfg.Harness.DefaultCreatedBy,
	}
	if err := s.store.PersistMultimodalExecution(ctx, persistInput); err != nil {
		warnings = append(warnings, fmt.Sprintf("persist execution artifact: %v", err))
		result.Warnings = warnings
		s.logger.Warn("persist execution artifact failed", "task_id", task.ID, "err", err)
	}

	if err := s.handleTaskOutcome(ctx, task, routingDecision, executionResult, nil, allSignals); err != nil {
		return nil, fmt.Errorf("handle task outcome: %w", err)
	}

	s.logger.Info("task processed",
		"task_id", task.ID,
		"task_class", task.Class,
		"primary_modality", task.PrimaryModality.String(),
		"chosen_target", routingDecision.ChosenTarget,
		"executor", executionResult.Plan.ChosenExecutor,
		"signal_count", len(allSignals),
		"warning_count", len(warnings),
	)

	return result, nil
}

func normalizeTask(task Task, cfg *config.AppConfig) Task {
	if strings.TrimSpace(task.ID) == "" {
		task.ID = fmt.Sprintf("task-%d", time.Now().UTC().UnixNano())
	}
	task.Summary = strings.TrimSpace(task.Summary)
	if task.Summary == "" {
		task.Summary = "runtime task"
	}
	task.Class = strings.TrimSpace(task.Class)
	if task.Class == "" {
		task.Class = "analysis"
	}
	task.PrimaryModality = modality.Normalize(task.PrimaryModality.String())
	if !task.PrimaryModality.IsKnown() {
		task.PrimaryModality = modality.Normalize(cfg.Runtime.DefaultPrimaryModality)
	}
	if !task.AllowTextOnlyFallback {
		task.AllowTextOnlyFallback = cfg.Runtime.AllowTextOnlyFallback
	}
	return task
}

func shouldDispatchRemote(cfg *config.AppConfig, dispatcher RemoteDispatcher, decision routing.Decision) bool {
	if cfg == nil || dispatcher == nil {
		return false
	}
	targetUnitID := strings.TrimSpace(decision.TargetUnitID)
	if targetUnitID == "" || targetUnitID == strings.TrimSpace(cfg.Runtime.SpecialistID) {
		return false
	}
	_, ok := cfg.TaskDispatch.RemoteUnitSockets[targetUnitID]
	return ok
}

func executionResultForRemoteDispatch(decision routing.Decision, remote *RemoteDispatchResult) execution.Result {
	status := ""
	output := ""
	metadata := map[string]string{"dispatch": "remote"}
	if remote != nil {
		status = remote.Status
		output = firstNonEmptyString(remote.OutputJSON, remote.ResultSummary)
		metadata["socket_path"] = remote.SocketPath
		metadata["remote_status"] = remote.Status
	}
	return execution.Result{
		Plan: execution.Plan{
			ExecutionMode:  "remote_rpc",
			ChosenExecutor: decision.ChosenTarget,
			TargetUnitID:   decision.TargetUnitID,
			TargetRole:     decision.TargetRole,
			TargetModelRef: decision.TargetModelRef,
			RequiresFusion: decision.RequiresFusion,
			Notes:          "remote task rpc dispatch",
		},
		HostResult: modelhost.Result{
			HostName: decision.TargetUnitID,
			Executor: decision.ChosenTarget,
			Output:   output,
			Handled:  status == "ok",
			Metadata: metadata,
		},
	}
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func impactForTaskSignals(signals []telemetry.Signal) string {
	if hasHighSeverity(signals) {
		return "high"
	}
	for _, signal := range signals {
		if signal.Severity == telemetry.SeverityModerate {
			return "medium"
		}
	}
	return "low"
}

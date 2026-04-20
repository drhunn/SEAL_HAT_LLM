package runtime

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/growth"
	"github.com/drhunn/SEAL_HAT_LLM/internal/harness"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slotpacket"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slots"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slotsync"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

type Task struct {
	ID                          string
	Summary                     string
	Class                       string
	PrimaryModality             modality.Type
	SecondaryModalities         []modality.Type
	CrossModalGroundingRequired bool
	AllowTextOnlyFallback       bool
	PreferredExecutor           string
	AssetRefs                   []string
	Prompt                      string
}

type TaskResult struct {
	Task             Task
	RetrievalResults []memory.RetrievalResult
	RoutingDecision  routing.Decision
	ExecutionResult  execution.Result
	Signals          []telemetry.Signal
}

type Service struct {
	cfg       *config.AppConfig
	loader    *slots.FilesystemLoader
	store     *memory.PostgresStore
	harness   *harness.Service
	routing   *routing.Service
	execution *execution.Service
	growth    *growth.Service
	slotSync  *slotsync.Service
	logger    *slog.Logger
}

func NewService(cfg *config.AppConfig, loader *slots.FilesystemLoader, store *memory.PostgresStore, harnessService *harness.Service, routingService *routing.Service, executionService *execution.Service, growthService *growth.Service, slotSyncService *slotsync.Service, logger *slog.Logger) *Service {
	return &Service{
		cfg:       cfg,
		loader:    loader,
		store:     store,
		harness:   harnessService,
		routing:   routingService,
		execution: executionService,
		growth:    growthService,
		slotSync:  slotSyncService,
		logger:    logger,
	}
}

func (s *Service) Start(ctx context.Context) error {
	runCtx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Runtime.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	if err := s.harness.StartupChecks(runCtx, s.cfg.Runtime.SpecialistID, s.cfg.Runtime.Namespace); err != nil {
		return fmt.Errorf("startup checks: %w", err)
	}

	slotFiles, err := s.loader.LoadSpecialistSlots(s.cfg.Runtime.SpecialistID)
	if err != nil {
		return fmt.Errorf("load specialist slots: %w", err)
	}

	if s.slotSync != nil {
		if err := s.slotSync.SyncFilesystemView(runCtx, s.cfg.Runtime.SpecialistID); err != nil {
			s.logger.Warn("slot sync failed", "specialist_id", s.cfg.Runtime.SpecialistID, "err", err)
		}
	}

	packet := slotpacket.BuildFromUnknown(s.cfg.Runtime.SpecialistID, "active", true, slotFiles)
	if err := slotpacket.WriteJSON("artifacts/slot_packet.json", packet); err != nil {
		s.logger.Warn("slot packet export failed", "specialist_id", s.cfg.Runtime.SpecialistID, "err", err)
	} else {
		s.logger.Info("slot packet exported",
			"specialist_id", s.cfg.Runtime.SpecialistID,
			"path", "artifacts/slot_packet.json",
			"version_hash", packet.VersionHash,
		)
	}

	s.logger.Info("runtime initialized",
		"specialist_id", s.cfg.Runtime.SpecialistID,
		"namespace", s.cfg.Runtime.Namespace,
		"slot_count", len(slotFiles),
	)

	if s.cfg.Runtime.EnableMultimodalSmokeTest && s.routing != nil && s.execution != nil {
		task := defaultStartupTask(s.cfg)
		result, err := s.ProcessTask(runCtx, task)
		if err != nil {
			return fmt.Errorf("process startup task: %w", err)
		}
		if s.growth != nil {
			growthResult, err := s.growth.StageExperiment(runCtx, s.cfg.Runtime.Namespace, s.cfg.Runtime.SpecialistID, growth.Assessment{
				AbilityName:       "multimodal_grounding",
				GapSummary:        "Bounded startup task still relies on simulated multimodal execution rather than live grounded backends.",
				EvidenceSummary:   fmt.Sprintf("task=%s executor=%s host=%s retrieval_results=%d signals=%d", result.Task.Summary, result.ExecutionResult.Plan.ChosenExecutor, result.ExecutionResult.HostResult.HostName, len(result.RetrievalResults), len(result.Signals)),
				TriedMemoryFix:    true,
				TriedRoutingFix:   true,
				TriedPromptFix:    true,
				PreferredSurface:  "modality_branch",
				RequestedBy:       s.cfg.Harness.DefaultCreatedBy,
				ParentApprovedBy:  "parent:startup-task",
				HarnessVerifiedBy: s.cfg.Harness.DefaultCreatedBy,
				Notes:             "startup bounded task path",
			})
			if err != nil {
				s.logger.Warn("ability growth task follow-up failed", "err", err)
			} else {
				s.logger.Info("ability growth task follow-up ok",
					"status", growthResult.Status,
					"experiment_id", growthResult.ExperimentID,
				)
			}
		}
	}

	if s.cfg.Runtime.EnableTaskInbox {
		inbox := NewTaskInbox(s.cfg.Runtime.TaskInboxDir, s.logger)
		if inbox == nil {
			return fmt.Errorf("task inbox enabled but runtime.task_inbox_dir is empty")
		}
		if err := inbox.EnsureDirs(); err != nil {
			return fmt.Errorf("ensure task inbox: %w", err)
		}
		s.logger.Info("task inbox enabled",
			"dir", s.cfg.Runtime.TaskInboxDir,
			"poll_interval_seconds", s.cfg.Runtime.TaskPollIntervalSeconds,
		)
		return s.runTaskInbox(ctx, inbox)
	}

	<-ctx.Done()
	s.logger.Info("shutdown requested")
	return nil
}

func (s *Service) runTaskInbox(ctx context.Context, inbox *TaskInbox) error {
	if err := s.processAvailableTasks(ctx, inbox); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Duration(s.cfg.Runtime.TaskPollIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("task inbox shutdown requested")
			return nil
		case <-ticker.C:
			if err := s.processAvailableTasks(ctx, inbox); err != nil {
				return err
			}
		}
	}
}

func (s *Service) processAvailableTasks(ctx context.Context, inbox *TaskInbox) error {
	for {
		queued, err := inbox.ClaimNext()
		if err != nil {
			if queued != nil {
				if markErr := inbox.MarkFailed(queued, err); markErr != nil {
					return fmt.Errorf("mark failed inbox task: %w", markErr)
				}
				continue
			}
			return err
		}
		if queued == nil {
			return nil
		}

		taskCtx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Runtime.RequestTimeoutSeconds)*time.Second)
		_, taskErr := s.ProcessTask(taskCtx, queued.Task)
		cancel()
		if taskErr != nil {
			if markErr := inbox.MarkFailed(queued, taskErr); markErr != nil {
				return fmt.Errorf("mark task as failed: %w", markErr)
			}
			continue
		}
		if err := inbox.MarkProcessed(queued); err != nil {
			return fmt.Errorf("mark task as processed: %w", err)
		}
	}
}

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
		s.logger.Warn("persist retrieval signals failed", "task_id", task.ID, "err", err)
	}

	routingInput := routing.Input{
		TaskSummary:                 task.Summary,
		TaskClass:                   task.Class,
		PrimaryModality:             task.PrimaryModality,
		SecondaryModalities:         task.SecondaryModalities,
		CrossModalGroundingRequired: task.CrossModalGroundingRequired,
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
		s.logger.Warn("routing audit write failed", "task_id", task.ID, "err", err)
	}
	routingSignals := routing.SignalsForDecision(s.cfg.Runtime.SpecialistID, routingInput, routingDecision, collector)
	if err := collector.Write(ctx, signalStore, routingSignals...); err != nil {
		return nil, fmt.Errorf("write routing signals to memory store: %w", err)
	}
	if err := collector.Write(ctx, s.store, routingSignals...); err != nil {
		s.logger.Warn("persist routing signals failed", "task_id", task.ID, "err", err)
	}

	executionReq := execution.Request{
		TaskSummary:                 task.Summary,
		TaskClass:                   task.Class,
		PrimaryModality:             task.PrimaryModality,
		SecondaryModalities:         task.SecondaryModalities,
		CrossModalGroundingRequired: task.CrossModalGroundingRequired,
		AllowTextOnlyFallback:       task.AllowTextOnlyFallback,
		PreferredExecutor:           task.PreferredExecutor,
		AssetRefs:                   task.AssetRefs,
		Prompt:                      task.Prompt,
	}
	executionResult, executionErr := s.execution.Execute(ctx, executionReq)
	executionSignals := execution.SignalsForExecution(s.cfg.Runtime.SpecialistID, executionReq, executionResult, executionErr, collector)
	if err := collector.Write(ctx, signalStore, executionSignals...); err != nil {
		return nil, fmt.Errorf("write execution signals to memory store: %w", err)
	}
	if err := collector.Write(ctx, s.store, executionSignals...); err != nil {
		s.logger.Warn("persist execution signals failed", "task_id", task.ID, "err", err)
	}
	if executionErr != nil {
		if err := s.handleTaskIncident(ctx, task, routingDecision, executionResult, executionErr, append(append([]telemetry.Signal{}, retrievalSignals...), append(routingSignals, executionSignals...)...)); err != nil {
			return nil, fmt.Errorf("handle task incident: %w", err)
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
		s.logger.Warn("persist execution artifact failed", "task_id", task.ID, "err", err)
	}

	allSignals := append(append([]telemetry.Signal{}, retrievalSignals...), append(routingSignals, executionSignals...)...)
	if err := s.handleTaskIncident(ctx, task, routingDecision, executionResult, nil, allSignals); err != nil {
		return nil, fmt.Errorf("handle task incident: %w", err)
	}

	s.logger.Info("task processed",
		"task_id", task.ID,
		"task_class", task.Class,
		"primary_modality", task.PrimaryModality.String(),
		"chosen_target", routingDecision.ChosenTarget,
		"executor", executionResult.Plan.ChosenExecutor,
		"signal_count", len(allSignals),
	)

	return &TaskResult{
		Task:             task,
		RetrievalResults: retrievalResults,
		RoutingDecision:  routingDecision,
		ExecutionResult:  executionResult,
		Signals:          allSignals,
	}, nil
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

func defaultStartupTask(cfg *config.AppConfig) Task {
	return normalizeTask(Task{
		ID:                          "startup-task",
		Summary:                     "startup multimodal execution task",
		Class:                       "evidence_fusion",
		PrimaryModality:             modality.Image,
		SecondaryModalities:         []modality.Type{modality.Text},
		CrossModalGroundingRequired: true,
		AllowTextOnlyFallback:       cfg.Runtime.AllowTextOnlyFallback,
		AssetRefs:                   []string{"sandbox://startup-smoke/image-1"},
		Prompt:                      "Compare image evidence with text context.",
	}, cfg)
}

func (s *Service) handleTaskIncident(ctx context.Context, task Task, decision routing.Decision, result execution.Result, executionErr error, signals []telemetry.Signal) error {
	incident, ok := classifyIncident(task, decision, result, executionErr, signals)
	if !ok || s.harness == nil {
		return nil
	}
	return s.harness.HandleIncident(ctx, s.cfg.Runtime.Namespace, s.cfg.Runtime.SpecialistID, incident)
}

func classifyIncident(task Task, decision routing.Decision, result execution.Result, executionErr error, signals []telemetry.Signal) (harness.Incident, bool) {
	if executionErr == nil && !decision.NeedsParentView && !result.Plan.NeedsParentReview && !hasHighSeverity(signals) {
		return harness.Incident{}, false
	}

	classification := []string{"execution failure"}
	whatWentWrong := "task execution failed"
	rootCause := whatWentWrong
	if executionErr != nil {
		whatWentWrong = executionErr.Error()
		rootCause = executionErr.Error()
	} else if decision.NeedsParentView || result.Plan.NeedsParentReview {
		classification = []string{"routing failure"}
		whatWentWrong = "task requires parent review before the bounded runtime path can be trusted"
		rootCause = whatWentWrong
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
		ExpectedBehavior:      "bounded runtime task should route, execute, persist, and complete without high-severity signals or forced parent review",
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

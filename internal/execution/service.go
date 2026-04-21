package execution

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modelhost"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unit"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unitref"
)

type Request struct {
	TaskSummary                 string
	TaskClass                   string
	PrimaryModality             modality.Type
	SecondaryModalities         []modality.Type
	CrossModalGroundingRequired bool
	AllowTextOnlyFallback       bool
	PreferredExecutor           string
	AssetRefs                   []string
	Prompt                      string
}

type Plan struct {
	ExecutionMode        string
	ChosenExecutor       string
	TargetUnitID         string
	TargetRole           unit.Role
	TargetModelRef       string
	RequiresFusion       bool
	UsesTextOnlyFallback bool
	NeedsParentReview    bool
	Notes                string
}

type Result struct {
	Plan       Plan
	HostResult modelhost.Result
}

type Service struct {
	logger *slog.Logger
	hosts  *modelhost.Registry
}

func NewService(logger *slog.Logger, hosts *modelhost.Registry) *Service {
	return &Service{logger: logger, hosts: hosts}
}

func (s *Service) Plan(ctx context.Context, req Request) Plan {
	selection := executors.Select(req.PrimaryModality, req.SecondaryModalities, req.CrossModalGroundingRequired, req.AllowTextOnlyFallback)
	target := unitref.ForExecutor(selection.Executor.String())
	plan := Plan{
		ExecutionMode:        "unimodal",
		ChosenExecutor:       selection.Executor.String(),
		TargetUnitID:         target.UnitID,
		TargetRole:           target.Role,
		TargetModelRef:       target.ModelRef,
		RequiresFusion:       selection.RequiresFusion,
		UsesTextOnlyFallback: selection.WasFallback,
		NeedsParentReview:    selection.NeedsParentView,
		Notes:                notesForSelection(selection.Executor, selection.RequiresFusion, selection.WasFallback),
	}

	if req.PreferredExecutor != "" && !selection.RequiresFusion && (req.PrimaryModality == modality.Text || req.PrimaryModality == modality.Unknown) {
		preferredTarget := unitref.ForExecutor(req.PreferredExecutor)
		plan.ChosenExecutor = req.PreferredExecutor
		plan.TargetUnitID = preferredTarget.UnitID
		plan.TargetRole = preferredTarget.Role
		plan.TargetModelRef = preferredTarget.ModelRef
		plan.Notes = "preferred executor requested"
	}

	if selection.RequiresFusion {
		plan.ExecutionMode = "multimodal_fusion"
	}
	if req.CrossModalGroundingRequired && len(req.AssetRefs) == 0 {
		plan.NeedsParentReview = true
		plan.Notes = "cross-modal request without asset refs should be reviewed by parent"
	}

	s.logger.InfoContext(ctx, "execution plan",
		"task_class", req.TaskClass,
		"primary_modality", req.PrimaryModality.String(),
		"chosen_executor", plan.ChosenExecutor,
		"target_unit_id", plan.TargetUnitID,
		"target_role", plan.TargetRole,
		"execution_mode", plan.ExecutionMode,
		"requires_fusion", plan.RequiresFusion,
		"parent_review", plan.NeedsParentReview,
	)

	return plan
}

func (s *Service) Execute(ctx context.Context, req Request) (Result, error) {
	plan := s.Plan(ctx, req)
	if s.hosts == nil {
		return Result{}, fmt.Errorf("execution hosts are not configured")
	}

	hostResult, err := s.hosts.Execute(ctx, plan.ChosenExecutor, modelhost.Request{
		TaskSummary:   req.TaskSummary,
		TaskClass:     req.TaskClass,
		ExecutionMode: plan.ExecutionMode,
		Executor:      plan.ChosenExecutor,
		AssetRefs:     req.AssetRefs,
		Prompt:        req.Prompt,
	})
	if err != nil {
		return Result{}, err
	}

	s.logger.InfoContext(ctx, "execution result",
		"executor", plan.ChosenExecutor,
		"target_unit_id", plan.TargetUnitID,
		"host", hostResult.HostName,
		"handled", hostResult.Handled,
	)

	return Result{Plan: plan, HostResult: hostResult}, nil
}

func SignalsForExecution(specialistID string, req Request, result Result, execErr error, collector *telemetry.Collector) []telemetry.Signal {
	if collector == nil {
		return nil
	}
	signals := make([]telemetry.Signal, 0)
	if execErr != nil {
		signals = append(signals, collector.NewSignal(specialistID, "execution", "execution", req.TaskClass, execErr.Error(), telemetry.SeverityHigh, result.Plan.ChosenExecutor, result.Plan.TargetUnitID))
		return signals
	}
	if result.Plan.UsesTextOnlyFallback {
		signals = append(signals, collector.NewSignal(specialistID, "execution", "execution", req.TaskClass, "execution used text-only fallback", telemetry.SeverityModerate, result.Plan.ChosenExecutor, result.Plan.TargetUnitID))
	}
	if result.Plan.NeedsParentReview {
		signals = append(signals, collector.NewSignal(specialistID, "execution", "execution", req.TaskClass, "execution plan requested parent review", telemetry.SeverityModerate, result.Plan.ChosenExecutor, result.Plan.TargetUnitID))
	}
	if result.Plan.RequiresFusion && result.Plan.ChosenExecutor != executors.MultimodalFusion.String() {
		signals = append(signals, collector.NewSignal(specialistID, "execution", "multimodal", req.TaskClass, "fusion-required execution was not assigned to fusion executor", telemetry.SeverityHigh, result.Plan.ChosenExecutor, result.Plan.TargetUnitID))
	}
	if !result.HostResult.Handled {
		signals = append(signals, collector.NewSignal(specialistID, "execution", "execution", req.TaskClass, "host did not handle execution request", telemetry.SeverityHigh, result.HostResult.HostName, result.Plan.ChosenExecutor, result.Plan.TargetUnitID))
	}
	return signals
}

func notesForSelection(executor executors.Name, requiresFusion, usedFallback bool) string {
	if requiresFusion {
		return "cross-modal grounding required"
	}
	if usedFallback {
		return "unknown modality fell back to text-first execution"
	}
	switch executor {
	case executors.ImageAnalysis:
		return "image-first execution path"
	case executors.AudioTranscription:
		return "audio-first execution path"
	case executors.VideoUnderstanding:
		return "video-first execution path"
	case executors.DocumentLayoutOCR:
		return "document-first execution path"
	default:
		return "text-first execution path"
	}
}

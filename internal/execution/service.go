package execution

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modelhost"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
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
	plan := Plan{
		ExecutionMode:  "unimodal",
		ChosenExecutor: "Parent-Generalist-30B",
		Notes:          "default execution path",
	}

	if req.PreferredExecutor != "" {
		plan.ChosenExecutor = req.PreferredExecutor
		plan.Notes = "preferred executor requested"
	}

	if req.CrossModalGroundingRequired || len(req.SecondaryModalities) > 0 || req.PrimaryModality == modality.Multimodal {
		plan.ExecutionMode = "multimodal_fusion"
		plan.ChosenExecutor = "Multimodal-Evidence-Fusion-Specialist-01"
		plan.RequiresFusion = true
		plan.Notes = "cross-modal grounding required"
	} else {
		switch req.PrimaryModality {
		case modality.Image:
			plan.ChosenExecutor = "Image-Analysis-Specialist-01"
			plan.Notes = "image-first execution path"
		case modality.Audio:
			plan.ChosenExecutor = "Audio-Transcription-Specialist-01"
			plan.Notes = "audio-first execution path"
		case modality.Video:
			plan.ChosenExecutor = "Video-Understanding-Specialist-01"
			plan.Notes = "video-first execution path"
		case modality.Document:
			plan.ChosenExecutor = "Document-Layout-OCR-Specialist-01"
			plan.Notes = "document-first execution path"
		case modality.Text, modality.Unknown:
			plan.ChosenExecutor = "Parent-Generalist-30B"
			plan.Notes = "text-first execution path"
		}
	}

	if req.PrimaryModality == modality.Unknown && req.AllowTextOnlyFallback {
		plan.ChosenExecutor = "Parent-Generalist-30B"
		plan.UsesTextOnlyFallback = true
		plan.Notes = "unknown modality fell back to text-first execution"
	}

	if req.CrossModalGroundingRequired && len(req.AssetRefs) == 0 {
		plan.NeedsParentReview = true
		plan.Notes = "cross-modal request without asset refs should be reviewed by parent"
	}

	s.logger.InfoContext(ctx, "execution plan",
		"task_class", req.TaskClass,
		"primary_modality", req.PrimaryModality.String(),
		"chosen_executor", plan.ChosenExecutor,
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
		signals = append(signals, collector.NewSignal(specialistID, "execution", "execution", req.TaskClass, execErr.Error(), telemetry.SeverityHigh, result.Plan.ChosenExecutor))
		return signals
	}
	if result.Plan.UsesTextOnlyFallback {
		signals = append(signals, collector.NewSignal(specialistID, "execution", "execution", req.TaskClass, "execution used text-only fallback", telemetry.SeverityModerate, result.Plan.ChosenExecutor))
	}
	if result.Plan.NeedsParentReview {
		signals = append(signals, collector.NewSignal(specialistID, "execution", "execution", req.TaskClass, "execution plan requested parent review", telemetry.SeverityHigh, result.Plan.ChosenExecutor))
	}
	if result.Plan.RequiresFusion && !strings.Contains(result.Plan.ChosenExecutor, "Fusion") {
		signals = append(signals, collector.NewSignal(specialistID, "execution", "multimodal", req.TaskClass, "fusion-required execution was not assigned to fusion executor", telemetry.SeverityHigh, result.Plan.ChosenExecutor))
	}
	if !result.HostResult.Handled {
		signals = append(signals, collector.NewSignal(specialistID, "execution", "execution", req.TaskClass, "host did not handle execution request", telemetry.SeverityHigh, result.HostResult.HostName, result.Plan.ChosenExecutor))
	}
	return signals
}

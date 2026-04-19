package execution

import (
	"context"
	"log/slog"

	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
)

type Request struct {
	TaskSummary                string
	TaskClass                  string
	PrimaryModality            modality.Type
	SecondaryModalities        []modality.Type
	CrossModalGroundingRequired bool
	AllowTextOnlyFallback      bool
	PreferredExecutor          string
	AssetRefs                  []string
}

type Plan struct {
	ExecutionMode        string
	ChosenExecutor       string
	RequiresFusion       bool
	UsesTextOnlyFallback bool
	NeedsParentReview    bool
	Notes                string
}

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{logger: logger}
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

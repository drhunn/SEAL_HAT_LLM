package routing

import (
	"context"
	"log/slog"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

type Input struct {
	TaskSummary                 string
	TaskClass                   string
	PrimaryModality             modality.Type
	SecondaryModalities         []modality.Type
	CrossModalGroundingRequired bool
}

type Decision struct {
	TaskSummary      string
	TaskClass        string
	ChosenTarget     string
	Confidence       float64
	WasFallback      bool
	FallbackReason   string
	NeedsParentView  bool
	PrimaryModality  string
	RequiresFusion   bool
}

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) Decide(ctx context.Context, taskSummary, taskClass string) Decision {
	return s.DecideTask(ctx, Input{
		TaskSummary:     taskSummary,
		TaskClass:       taskClass,
		PrimaryModality: modality.Text,
	})
}

func (s *Service) DecideTask(ctx context.Context, in Input) Decision {
	primary := in.PrimaryModality
	if !primary.IsKnown() {
		primary = modality.Text
	}

	decision := Decision{
		TaskSummary:     in.TaskSummary,
		TaskClass:       in.TaskClass,
		ChosenTarget:    "Parent-Generalist-30B",
		Confidence:      0.50,
		PrimaryModality: primary.String(),
	}

	if in.CrossModalGroundingRequired || len(in.SecondaryModalities) > 0 || primary == modality.Multimodal {
		decision.ChosenTarget = "Multimodal-Evidence-Fusion-Specialist-01"
		decision.Confidence = 0.78
		decision.RequiresFusion = true
		decision.NeedsParentView = true
	} else {
		switch primary {
		case modality.Image:
			decision.ChosenTarget = "Image-Analysis-Specialist-01"
			decision.Confidence = 0.74
		case modality.Audio:
			decision.ChosenTarget = "Audio-Transcription-Specialist-01"
			decision.Confidence = 0.74
		case modality.Video:
			decision.ChosenTarget = "Video-Understanding-Specialist-01"
			decision.Confidence = 0.76
		case modality.Document:
			decision.ChosenTarget = "Document-Layout-OCR-Specialist-01"
			decision.Confidence = 0.72
		case modality.Text:
			decision.ChosenTarget = "Parent-Generalist-30B"
			decision.Confidence = 0.50
		default:
			decision.ChosenTarget = "Parent-Generalist-30B"
			decision.Confidence = 0.45
			decision.WasFallback = true
			decision.FallbackReason = "unknown modality defaulted to parent"
		}
	}

	s.logger.InfoContext(ctx, "routing decision",
		"task_class", decision.TaskClass,
		"primary_modality", decision.PrimaryModality,
		"chosen_target", decision.ChosenTarget,
		"confidence", decision.Confidence,
		"requires_fusion", decision.RequiresFusion,
	)
	return decision
}

func SignalsForDecision(specialistID string, in Input, decision Decision, collector *telemetry.Collector) []telemetry.Signal {
	if collector == nil {
		return nil
	}
	signals := make([]telemetry.Signal, 0)
	if decision.WasFallback {
		signals = append(signals, collector.NewSignal(specialistID, "routing", "routing", in.TaskClass, firstNonEmpty(decision.FallbackReason, "routing fallback used"), telemetry.SeverityModerate, decision.ChosenTarget))
	}
	if decision.NeedsParentView {
		signals = append(signals, collector.NewSignal(specialistID, "routing", "routing", in.TaskClass, "routing requested parent review", telemetry.SeverityHigh, decision.ChosenTarget))
	}
	if decision.Confidence < 0.55 {
		signals = append(signals, collector.NewSignal(specialistID, "routing", "routing", in.TaskClass, "low-confidence routing decision", telemetry.SeverityLow, decision.ChosenTarget))
	}
	if decision.RequiresFusion && !strings.Contains(decision.ChosenTarget, "Fusion") {
		signals = append(signals, collector.NewSignal(specialistID, "routing", "routing", in.TaskClass, "fusion-required task was not routed to fusion executor", telemetry.SeverityHigh, decision.ChosenTarget))
	}
	return signals
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

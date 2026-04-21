package executors

import "github.com/drhunn/SEAL_HAT_LLM/internal/modality"

type Name string

const (
	ParentGeneralist      Name = "Parent-Generalist-30B"
	MultimodalFusion      Name = "Multimodal-Evidence-Fusion-Specialist-01"
	ImageAnalysis         Name = "Image-Analysis-Specialist-01"
	AudioTranscription    Name = "Audio-Transcription-Specialist-01"
	VideoUnderstanding    Name = "Video-Understanding-Specialist-01"
	DocumentLayoutOCR     Name = "Document-Layout-OCR-Specialist-01"
)

func (n Name) String() string {
	return string(n)
}

type Selection struct {
	Executor        Name
	Confidence      float64
	RequiresFusion  bool
	NeedsParentView bool
	WasFallback     bool
	FallbackReason  string
}

func Select(primary modality.Type, secondary []modality.Type, crossModalGroundingRequired bool, allowTextOnlyFallback bool) Selection {
	if !primary.IsKnown() {
		primary = modality.Unknown
	}

	selection := Selection{
		Executor:   ParentGeneralist,
		Confidence: 0.50,
	}

	if crossModalGroundingRequired || len(secondary) > 0 || primary == modality.Multimodal {
		selection.Executor = MultimodalFusion
		selection.Confidence = 0.78
		selection.RequiresFusion = true
		selection.NeedsParentView = true
		return selection
	}

	switch primary {
	case modality.Image:
		selection.Executor = ImageAnalysis
		selection.Confidence = 0.74
	case modality.Audio:
		selection.Executor = AudioTranscription
		selection.Confidence = 0.74
	case modality.Video:
		selection.Executor = VideoUnderstanding
		selection.Confidence = 0.76
	case modality.Document:
		selection.Executor = DocumentLayoutOCR
		selection.Confidence = 0.72
	case modality.Text:
		selection.Executor = ParentGeneralist
		selection.Confidence = 0.50
	default:
		selection.Executor = ParentGeneralist
		selection.Confidence = 0.45
		if allowTextOnlyFallback {
			selection.WasFallback = true
			selection.FallbackReason = "unknown modality defaulted to parent"
		}
	}

	return selection
}

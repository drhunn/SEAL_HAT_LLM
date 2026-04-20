package modelhost

import "strings"

const (
	ExecutorParentGeneralist         = "Parent-Generalist-30B"
	ExecutorImageAnalysis           = "Image-Analysis-Specialist-01"
	ExecutorAudioTranscription      = "Audio-Transcription-Specialist-01"
	ExecutorVideoUnderstanding      = "Video-Understanding-Specialist-01"
	ExecutorDocumentLayoutOCR       = "Document-Layout-OCR-Specialist-01"
	ExecutorMultimodalEvidenceFusion = "Multimodal-Evidence-Fusion-Specialist-01"
)

func NewSimulatedRegistry(prefix string) *Registry {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = "local"
	}

	registry := NewRegistry()
	registry.Register(ExecutorParentGeneralist, NewPromptHost(prefix+"-parent-host"))
	registry.Register(ExecutorImageAnalysis, NewAssetSummaryHost(prefix+"-image-host", "image"))
	registry.Register(ExecutorAudioTranscription, NewAssetSummaryHost(prefix+"-audio-host", "audio"))
	registry.Register(ExecutorVideoUnderstanding, NewAssetSummaryHost(prefix+"-video-host", "video"))
	registry.Register(ExecutorDocumentLayoutOCR, NewAssetSummaryHost(prefix+"-document-host", "document"))
	registry.Register(ExecutorMultimodalEvidenceFusion, NewFusionHost(prefix+"-fusion-host"))
	return registry
}

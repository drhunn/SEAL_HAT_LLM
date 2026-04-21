package runtime

import "time"

const (
	TaskInboxSchemaVersion    = "seal_hat_llm.runtime_task.v1"
	TaskArtifactSchemaVersion = "seal_hat_llm.runtime_task_artifact.v1"
)

type TaskFile struct {
	Version                     string   `json:"version"`
	ID                          string   `json:"id,omitempty"`
	Summary                     string   `json:"summary"`
	Class                       string   `json:"class,omitempty"`
	PrimaryModality             string   `json:"primary_modality,omitempty"`
	SecondaryModalities         []string `json:"secondary_modalities,omitempty"`
	CrossModalGroundingRequired bool     `json:"cross_modal_grounding_required,omitempty"`
	AllowTextOnlyFallback       bool     `json:"allow_text_only_fallback,omitempty"`
	PreferredExecutor           string   `json:"preferred_executor,omitempty"`
	AssetRefs                   []string `json:"asset_refs,omitempty"`
	Prompt                      string   `json:"prompt,omitempty"`
}

type TaskArtifact struct {
	Version         string    `json:"version"`
	TaskID          string    `json:"task_id"`
	Status          string    `json:"status"`
	Summary         string    `json:"summary"`
	Class           string    `json:"class"`
	PrimaryModality string    `json:"primary_modality,omitempty"`
	Executor        string    `json:"executor,omitempty"`
	HostName        string    `json:"host_name,omitempty"`
	ArchivedPath    string    `json:"archived_path"`
	Error           string    `json:"error,omitempty"`
	SignalSummaries []string  `json:"signal_summaries,omitempty"`
	CompletedAt     time.Time `json:"completed_at"`
}

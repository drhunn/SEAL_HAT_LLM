package taskrpc

import "context"

type RunTaskRequest struct {
	TaskID                      string            `json:"task_id"`
	ParentUnitID                string            `json:"parent_unit_id"`
	TargetUnitID                string            `json:"target_unit_id"`
	TaskClass                   string            `json:"task_class"`
	Summary                     string            `json:"summary"`
	Prompt                      string            `json:"prompt"`
	PrimaryModality             string            `json:"primary_modality"`
	SecondaryModalities         []string          `json:"secondary_modalities"`
	CrossModalGroundingRequired bool              `json:"cross_modal_grounding_required"`
	AllowTextOnlyFallback       bool              `json:"allow_text_only_fallback"`
	AssetRefs                   []string          `json:"asset_refs"`
	Constraints                 map[string]string `json:"constraints"`
	RouteEpisodeID              string            `json:"route_episode_id"`
}

type RunTaskResponse struct {
	TaskID           string   `json:"task_id"`
	SpecialistUnitID string   `json:"specialist_unit_id"`
	Status           string   `json:"status"`
	ResultSummary    string   `json:"result_summary"`
	OutputJSON       string   `json:"output_json"`
	ArtifactRefs     []string `json:"artifact_refs"`
	Confidence       float64  `json:"confidence"`
	Warnings         []string `json:"warnings"`
	Signals          []string `json:"signals"`
	ErrorText        string   `json:"error_text,omitempty"`
}

type TaskHandler interface {
	RunTask(context.Context, RunTaskRequest) (RunTaskResponse, error)
}

package modelhost

import (
	"context"
	"fmt"
	"strings"
)

type PromptHost struct {
	hostName string
}

func NewPromptHost(hostName string) *PromptHost {
	return &PromptHost{hostName: hostName}
}

func (h *PromptHost) Name() string {
	return h.hostName
}

func (h *PromptHost) Execute(_ context.Context, req Request) (Result, error) {
	output := strings.TrimSpace(req.Prompt)
	if output == "" {
		output = req.TaskSummary
	}
	if output == "" {
		output = "text execution completed"
	}
	return Result{
		HostName: h.hostName,
		Executor: req.Executor,
		Output:   fmt.Sprintf("text host processed: %s", output),
		Handled:  true,
		Metadata: map[string]string{
			"task_class":     req.TaskClass,
			"execution_mode": req.ExecutionMode,
			"asset_count":    fmt.Sprintf("%d", len(req.AssetRefs)),
		},
	}, nil
}

type AssetSummaryHost struct {
	hostName   string
	capability string
}

func NewAssetSummaryHost(hostName, capability string) *AssetSummaryHost {
	return &AssetSummaryHost{hostName: hostName, capability: capability}
}

func (h *AssetSummaryHost) Name() string {
	return h.hostName
}

func (h *AssetSummaryHost) Execute(_ context.Context, req Request) (Result, error) {
	assetSummary := "no assets"
	if len(req.AssetRefs) > 0 {
		assetSummary = strings.Join(req.AssetRefs, ", ")
	}
	output := fmt.Sprintf("%s host reviewed assets [%s] for task %q", h.capability, assetSummary, req.TaskSummary)
	return Result{
		HostName: h.hostName,
		Executor: req.Executor,
		Output:   output,
		Handled:  true,
		Metadata: map[string]string{
			"task_class":     req.TaskClass,
			"execution_mode": req.ExecutionMode,
			"capability":     h.capability,
		},
	}, nil
}

type FusionHost struct {
	hostName string
}

func NewFusionHost(hostName string) *FusionHost {
	return &FusionHost{hostName: hostName}
}

func (h *FusionHost) Name() string {
	return h.hostName
}

func (h *FusionHost) Execute(_ context.Context, req Request) (Result, error) {
	assetSummary := "no linked assets"
	if len(req.AssetRefs) > 0 {
		assetSummary = strings.Join(req.AssetRefs, ", ")
	}
	promptSummary := strings.TrimSpace(req.Prompt)
	if promptSummary == "" {
		promptSummary = "no extra prompt"
	}
	output := fmt.Sprintf("fusion host combined assets [%s] with prompt context %q", assetSummary, promptSummary)
	return Result{
		HostName: h.hostName,
		Executor: req.Executor,
		Output:   output,
		Handled:  true,
		Metadata: map[string]string{
			"task_class":     req.TaskClass,
			"execution_mode": req.ExecutionMode,
			"asset_count":    fmt.Sprintf("%d", len(req.AssetRefs)),
			"fused":          "true",
		},
	}, nil
}

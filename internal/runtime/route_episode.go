package runtime

import (
	"context"

	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
)

func (s *Service) persistRouteEpisode(ctx context.Context, result TaskResult) {
	if s == nil || s.store == nil {
		return
	}
	_, err := s.store.CreateRouteEpisode(ctx, memory.RouteEpisodeInput{
		Namespace:         s.cfg.Runtime.Namespace,
		SpecialistID:      s.cfg.Runtime.SpecialistID,
		TaskID:            result.Task.ID,
		TaskSummary:       result.Task.Summary,
		TaskClass:         result.Task.Class,
		PrimaryModality:   result.RoutingDecision.PrimaryModality,
		ChosenTarget:      result.RoutingDecision.ChosenTarget,
		TargetUnitID:      result.RoutingDecision.TargetUnitID,
		TargetRole:        string(result.RoutingDecision.TargetRole),
		TargetModelRef:    result.RoutingDecision.TargetModelRef,
		ChosenExecutor:    result.ExecutionResult.Plan.ChosenExecutor,
		ExecutionMode:     result.ExecutionResult.Plan.ExecutionMode,
		Confidence:        result.RoutingDecision.Confidence,
		WasFallback:       result.RoutingDecision.WasFallback,
		FallbackReason:    result.RoutingDecision.FallbackReason,
		NeedsParentReview: result.RoutingDecision.NeedsParentView || result.ExecutionResult.Plan.NeedsParentReview,
		RequiresFusion:    result.RoutingDecision.RequiresFusion || result.ExecutionResult.Plan.RequiresFusion,
		ExecutionHandled:  result.ExecutionResult.HostResult.Handled,
		ExecutionHost:     result.ExecutionResult.HostResult.HostName,
		Status:            "recorded",
		ErrorText:         "",
	})
	if err != nil {
		s.logger.Warn("route episode persistence failed", "task_id", result.Task.ID, "err", err)
	}
}

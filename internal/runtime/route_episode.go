package runtime

import (
	"context"

	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
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

func (s *Service) persistFailedRouteEpisode(ctx context.Context, task Task, taskErr error) {
	if s == nil || s.store == nil || taskErr == nil {
		return
	}
	route := routing.Decision{}
	if s.routing != nil {
		route = s.routing.DecideTask(ctx, routing.Input{
			TaskSummary:                 task.Summary,
			TaskClass:                   task.Class,
			PrimaryModality:             task.PrimaryModality,
			SecondaryModalities:         task.SecondaryModalities,
			CrossModalGroundingRequired: task.CrossModalGroundingRequired,
		})
	}
	plan := execution.Plan{}
	if s.execution != nil {
		plan = s.execution.Plan(ctx, execution.Request{
			TaskSummary:                 task.Summary,
			TaskClass:                   task.Class,
			PrimaryModality:             task.PrimaryModality,
			SecondaryModalities:         task.SecondaryModalities,
			CrossModalGroundingRequired: task.CrossModalGroundingRequired,
			AllowTextOnlyFallback:       task.AllowTextOnlyFallback,
			PreferredExecutor:           task.PreferredExecutor,
			AssetRefs:                   task.AssetRefs,
			Prompt:                      task.Prompt,
		})
	}
	_, err := s.store.CreateRouteEpisode(ctx, memory.RouteEpisodeInput{
		Namespace:         s.cfg.Runtime.Namespace,
		SpecialistID:      s.cfg.Runtime.SpecialistID,
		TaskID:            task.ID,
		TaskSummary:       task.Summary,
		TaskClass:         task.Class,
		PrimaryModality:   route.PrimaryModality,
		ChosenTarget:      route.ChosenTarget,
		TargetUnitID:      route.TargetUnitID,
		TargetRole:        string(route.TargetRole),
		TargetModelRef:    route.TargetModelRef,
		ChosenExecutor:    plan.ChosenExecutor,
		ExecutionMode:     plan.ExecutionMode,
		Confidence:        route.Confidence,
		WasFallback:       route.WasFallback,
		FallbackReason:    route.FallbackReason,
		NeedsParentReview: route.NeedsParentView || plan.NeedsParentReview,
		RequiresFusion:    route.RequiresFusion || plan.RequiresFusion,
		ExecutionHandled:  false,
		ExecutionHost:     "",
		Status:            "failed",
		ErrorText:         taskErr.Error(),
	})
	if err != nil {
		s.logger.Warn("failed route episode persistence failed", "task_id", task.ID, "err", err)
	}
}

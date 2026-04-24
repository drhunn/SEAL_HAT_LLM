package memory

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type RouteEpisodeInput struct {
	Namespace         string
	SpecialistID      string
	TaskID            string
	TaskSummary       string
	TaskClass         string
	PrimaryModality   string
	ChosenTarget      string
	TargetUnitID      string
	TargetRole        string
	TargetModelRef    string
	ChosenExecutor    string
	ExecutionMode     string
	Confidence        float64
	WasFallback       bool
	FallbackReason    string
	NeedsParentReview bool
	RequiresFusion    bool
	ExecutionHandled  bool
	ExecutionHost     string
	Status            string
	ErrorText         string
}

func (s *PostgresStore) CreateRouteEpisode(ctx context.Context, in RouteEpisodeInput) (string, error) {
	if s == nil || s.db == nil {
		return "", fmt.Errorf("postgres store is required")
	}
	id := routeEpisodeID(in.TaskID)
	status := strings.TrimSpace(in.Status)
	if status == "" {
		status = "recorded"
	}
	const q = `
		INSERT INTO agent_core.route_episodes (
			id,
			namespace,
			specialist_id,
			task_id,
			task_summary,
			task_class,
			primary_modality,
			chosen_target,
			target_unit_id,
			target_role,
			target_model_ref,
			chosen_executor,
			execution_mode,
			confidence,
			was_fallback,
			fallback_reason,
			needs_parent_review,
			requires_fusion,
			execution_handled,
			execution_host,
			status,
			error_text
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)
		RETURNING id`
	var out string
	if err := s.db.QueryRow(
		ctx,
		q,
		id,
		in.Namespace,
		in.SpecialistID,
		strings.TrimSpace(in.TaskID),
		strings.TrimSpace(in.TaskSummary),
		strings.TrimSpace(in.TaskClass),
		strings.TrimSpace(in.PrimaryModality),
		strings.TrimSpace(in.ChosenTarget),
		strings.TrimSpace(in.TargetUnitID),
		strings.TrimSpace(in.TargetRole),
		strings.TrimSpace(in.TargetModelRef),
		strings.TrimSpace(in.ChosenExecutor),
		strings.TrimSpace(in.ExecutionMode),
		in.Confidence,
		in.WasFallback,
		strings.TrimSpace(in.FallbackReason),
		in.NeedsParentReview,
		in.RequiresFusion,
		in.ExecutionHandled,
		strings.TrimSpace(in.ExecutionHost),
		status,
		strings.TrimSpace(in.ErrorText),
	).Scan(&out); err != nil {
		return "", fmt.Errorf("insert route episode: %w", err)
	}
	return out, nil
}

func routeEpisodeID(taskID string) string {
	base := strings.TrimSpace(taskID)
	if base == "" {
		base = "task"
	}
	return fmt.Sprintf("route-episode-%s-%d", base, time.Now().UTC().UnixNano())
}

package memory

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type RouteEpisodeInput struct {
	ID                string
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
		) VALUES ($1,$2,$3,$4,$5,$6,$7,NULLIF($8,''),NULLIF($9,''),NULLIF($10,''),NULLIF($11,''),NULLIF($12,''),NULLIF($13,''),$14,$15,NULLIF($16,''),$17,$18,$19,NULLIF($20,''),$21,NULLIF($22,''))
		RETURNING id`

	id := strings.TrimSpace(in.ID)
	if id == "" {
		id = defaultRouteEpisodeID(in.TaskID)
	}
	status := strings.TrimSpace(in.Status)
	if status == "" {
		status = "recorded"
	}

	var out string
	if err := s.db.QueryRow(
		ctx,
		q,
		id,
		in.Namespace,
		in.SpecialistID,
		in.TaskID,
		in.TaskSummary,
		in.TaskClass,
		in.PrimaryModality,
		in.ChosenTarget,
		in.TargetUnitID,
		in.TargetRole,
		in.TargetModelRef,
		in.ChosenExecutor,
		in.ExecutionMode,
		in.Confidence,
		in.WasFallback,
		in.FallbackReason,
		in.NeedsParentReview,
		in.RequiresFusion,
		in.ExecutionHandled,
		in.ExecutionHost,
		status,
		in.ErrorText,
	).Scan(&out); err != nil {
		return "", fmt.Errorf("create route episode: %w", err)
	}
	return out, nil
}

func defaultRouteEpisodeID(taskID string) string {
	clean := strings.TrimSpace(taskID)
	if clean == "" {
		clean = "task"
	}
	clean = strings.ReplaceAll(clean, " ", "-")
	clean = strings.ReplaceAll(clean, "/", "-")
	return fmt.Sprintf("route-episode-%s-%d", clean, time.Now().UTC().UnixNano())
}

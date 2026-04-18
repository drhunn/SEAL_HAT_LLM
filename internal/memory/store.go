package memory

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostmortemInput struct {
	Namespace             string
	SpecialistID          string
	TaskSummary           string
	ExpectedBehavior      string
	ActualBehavior        string
	WhatWentWrong         string
	FailureClassification []string
	RootCause             string
	Preventable           bool
	CreatedBy             string
}

type HealthSnapshot struct {
	SpecialistID string
	Status       string
	HealthScore  float64
}

type PostgresStore struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewPostgresStore(db *pgxpool.Pool, logger *slog.Logger) *PostgresStore {
	return &PostgresStore{db: db, logger: logger}
}

func (s *PostgresStore) Ping(ctx context.Context) error {
	if s.db == nil {
		return errors.New("nil database pool")
	}
	return s.db.Ping(ctx)
}

func (s *PostgresStore) CreatePostmortem(ctx context.Context, in PostmortemInput) (string, error) {
	const q = `
		INSERT INTO agent_core.memory_postmortems (
			namespace,
			specialist_id,
			task_summary,
			expected_behavior,
			actual_behavior,
			what_went_wrong,
			failure_classification,
			root_cause,
			preventable,
			created_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING postmortem_id::text`

	var id string
	if err := s.db.QueryRow(
		ctx,
		q,
		in.Namespace,
		in.SpecialistID,
		in.TaskSummary,
		in.ExpectedBehavior,
		in.ActualBehavior,
		in.WhatWentWrong,
		in.FailureClassification,
		in.RootCause,
		in.Preventable,
		in.CreatedBy,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("insert postmortem: %w", err)
	}

	return id, nil
}

func (s *PostgresStore) UpdateSpecialistHealth(ctx context.Context, specialistID string) error {
	const q = `SELECT specialist_id FROM agent_core.fn_update_specialist_health_stats($1)`
	var out string
	if err := s.db.QueryRow(ctx, q, specialistID).Scan(&out); err != nil {
		s.logger.Warn("update specialist health failed", "specialist_id", specialistID, "err", err)
		return err
	}
	return nil
}

func (s *PostgresStore) ProjectMemorySummary(ctx context.Context, namespace, specialistID string) error {
	const q = `SELECT slot_id::text FROM agent_core.fn_project_and_write_memory_summary_slot($1,$2,$3,$4)`
	var slotID string
	if err := s.db.QueryRow(ctx, q, namespace, specialistID, "harness:runtime", "runtime memory summary projection").Scan(&slotID); err != nil {
		s.logger.Warn("project memory summary failed", "specialist_id", specialistID, "err", err)
		return err
	}
	return nil
}

func (s *PostgresStore) HealthSnapshot(ctx context.Context, specialistID string) (*HealthSnapshot, error) {
	const q = `SELECT specialist_id, status::text, COALESCE(health_score, 0) FROM agent_core.fn_get_specialist_health_snapshot($1)`
	var snapshot HealthSnapshot
	if err := s.db.QueryRow(ctx, q, specialistID).Scan(&snapshot.SpecialistID, &snapshot.Status, &snapshot.HealthScore); err != nil {
		return nil, fmt.Errorf("health snapshot: %w", err)
	}
	return &snapshot, nil
}

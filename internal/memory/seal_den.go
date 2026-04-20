package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/den"
	"github.com/drhunn/SEAL_HAT_LLM/internal/seal"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

func (s *PostgresStore) WriteSignals(ctx context.Context, signals []telemetry.Signal) error {
	const q = `
		INSERT INTO agent_core.adaptation_signals (
			id,
			specialist_id,
			category,
			severity,
			surface,
			task_class,
			summary,
			evidence_refs,
			occurred_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id) DO UPDATE SET
			summary = EXCLUDED.summary,
			evidence_refs = EXCLUDED.evidence_refs,
			occurred_at = EXCLUDED.occurred_at`
	for _, signal := range signals {
		evidenceJSON, err := json.Marshal(signal.EvidenceRefs)
		if err != nil {
			return fmt.Errorf("marshal signal evidence refs: %w", err)
		}
		occurredAt := signal.OccurredAt
		if occurredAt.IsZero() {
			occurredAt = time.Now().UTC()
		}
		if _, err := s.db.Exec(ctx, q,
			signal.ID,
			signal.SpecialistID,
			signal.Category,
			string(signal.Severity),
			signal.Surface,
			signal.TaskClass,
			signal.Summary,
			evidenceJSON,
			occurredAt,
		); err != nil {
			return fmt.Errorf("insert adaptation signal: %w", err)
		}
	}
	return nil
}

func (s *PostgresStore) ListSignals(ctx context.Context, specialistID string, limit int) ([]telemetry.Signal, error) {
	const q = `
		SELECT id, specialist_id, category, severity, surface, task_class, summary, evidence_refs, occurred_at
		FROM agent_core.adaptation_signals
		WHERE specialist_id = $1
		ORDER BY occurred_at DESC
		LIMIT $2`
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.Query(ctx, q, specialistID, limit)
	if err != nil {
		return nil, fmt.Errorf("list adaptation signals: %w", err)
	}
	defer rows.Close()

	out := make([]telemetry.Signal, 0, limit)
	for rows.Next() {
		var signal telemetry.Signal
		var severity string
		var evidenceJSON []byte
		if err := rows.Scan(&signal.ID, &signal.SpecialistID, &signal.Category, &severity, &signal.Surface, &signal.TaskClass, &signal.Summary, &evidenceJSON, &signal.OccurredAt); err != nil {
			return nil, fmt.Errorf("scan adaptation signal: %w", err)
		}
		signal.Severity = telemetry.Severity(severity)
		if len(evidenceJSON) > 0 {
			if err := json.Unmarshal(evidenceJSON, &signal.EvidenceRefs); err != nil {
				return nil, fmt.Errorf("unmarshal signal evidence refs: %w", err)
			}
		}
		out = append(out, signal)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate adaptation signals: %w", err)
	}
	return out, nil
}

func (s *PostgresStore) WriteProposal(ctx context.Context, proposal seal.AdaptationProposal) error {
	const q = `
		INSERT INTO agent_core.adaptation_proposals (
			id,
			specialist_id,
			cluster_id,
			surface,
			reason,
			requested_by,
			risk_level,
			requires_harness,
			requires_parent,
			rollback_required,
			status
		) VALUES ($1,$2,NULLIF($3,''),$4,$5,$6,$7,$8,$9,$10,'proposed')
		ON CONFLICT (id) DO UPDATE SET
			reason = EXCLUDED.reason,
			risk_level = EXCLUDED.risk_level,
			requires_harness = EXCLUDED.requires_harness,
			requires_parent = EXCLUDED.requires_parent,
			rollback_required = EXCLUDED.rollback_required`
	if _, err := s.db.Exec(ctx, q,
		proposal.ID,
		proposal.SpecialistID,
		proposal.ClusterID,
		string(proposal.Surface),
		proposal.Reason,
		proposal.RequestedBy,
		proposal.RiskLevel,
		proposal.RequiresHarness,
		proposal.RequiresParent,
		proposal.RollbackRequired,
	); err != nil {
		return fmt.Errorf("insert adaptation proposal: %w", err)
	}
	return nil
}

func (s *PostgresStore) WriteGrowthPlan(ctx context.Context, plan den.GrowthPlan) error {
	const q = `
		INSERT INTO agent_core.growth_plans (
			id,
			proposal_id,
			specialist_id,
			surface,
			reason,
			experiment_name,
			freeze_plan,
			rollback_plan,
			status
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,'planned')
		ON CONFLICT (id) DO UPDATE SET
			reason = EXCLUDED.reason,
			experiment_name = EXCLUDED.experiment_name,
			freeze_plan = EXCLUDED.freeze_plan,
			rollback_plan = EXCLUDED.rollback_plan`
	freezeJSON, err := json.Marshal(plan.FreezePlan)
	if err != nil {
		return fmt.Errorf("marshal freeze plan: %w", err)
	}
	rollbackJSON, err := json.Marshal(plan.RollbackPlan)
	if err != nil {
		return fmt.Errorf("marshal rollback plan: %w", err)
	}
	if _, err := s.db.Exec(ctx, q,
		plan.ID,
		plan.ProposalID,
		plan.SpecialistID,
		string(plan.Surface),
		plan.Reason,
		plan.ExperimentName,
		freezeJSON,
		rollbackJSON,
	); err != nil {
		return fmt.Errorf("insert growth plan: %w", err)
	}
	return nil
}

package memory

import (
	"context"
	"fmt"
)

type AbilityLedgerInput struct {
	Namespace      string
	SpecialistID   string
	AbilityName    string
	MaturityStage  string
	Score          float64
	Evidence       string
	LastAction     string
	NextAction     string
	UpdatedBy      string
}

type AbilityGrowthExperimentInput struct {
	Namespace            string
	SpecialistID         string
	AbilityName          string
	GapSummary           string
	EvidenceSummary      string
	PreferredSurface     string
	Status               string
	RequestedBy          string
	ParentApprovedBy     string
	HarnessVerifiedBy    string
	Notes                string
}

func (s *PostgresStore) UpsertAbilityLedger(ctx context.Context, in AbilityLedgerInput) error {
	const q = `
		INSERT INTO agent_core.ability_ledgers (
			namespace,
			specialist_id,
			ability_name,
			maturity_stage,
			score,
			evidence,
			last_action,
			next_action,
			updated_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (namespace, specialist_id, ability_name) DO UPDATE SET
			maturity_stage = EXCLUDED.maturity_stage,
			score = EXCLUDED.score,
			evidence = EXCLUDED.evidence,
			last_action = EXCLUDED.last_action,
			next_action = EXCLUDED.next_action,
			updated_by = EXCLUDED.updated_by,
			updated_at = now()`
	if _, err := s.db.Exec(ctx, q, in.Namespace, in.SpecialistID, in.AbilityName, in.MaturityStage, in.Score, in.Evidence, in.LastAction, in.NextAction, in.UpdatedBy); err != nil {
		return fmt.Errorf("upsert ability ledger: %w", err)
	}
	return nil
}

func (s *PostgresStore) CreateAbilityGrowthExperiment(ctx context.Context, in AbilityGrowthExperimentInput) (string, error) {
	const q = `
		INSERT INTO agent_core.ability_growth_experiments (
			namespace,
			specialist_id,
			ability_name,
			gap_summary,
			evidence_summary,
			preferred_surface,
			status,
			requested_by,
			parent_approved_by,
			harness_verified_by,
			notes
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NULLIF($9,''),NULLIF($10,''),$11)
		RETURNING experiment_id::text`
	var id string
	if err := s.db.QueryRow(ctx, q, in.Namespace, in.SpecialistID, in.AbilityName, in.GapSummary, in.EvidenceSummary, in.PreferredSurface, in.Status, in.RequestedBy, in.ParentApprovedBy, in.HarnessVerifiedBy, in.Notes).Scan(&id); err != nil {
		return "", fmt.Errorf("create ability growth experiment: %w", err)
	}
	return id, nil
}

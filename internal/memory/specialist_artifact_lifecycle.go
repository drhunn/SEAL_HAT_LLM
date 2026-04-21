package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type CandidateSpecialistArtifactInput struct {
	SpecialistID      string
	ParentArtifactID  string
	AbilityName       string
	RequestedBy       string
	ModelRef          string
	HarnessConfigRef  string
	StoreMode         string
	StoreBootstrapRef string
	SlotBundleRef     string
	SlotVersionHash   string
	EvalSuiteRef      string
	Notes             string
}

func (s *PostgresStore) GetCurrentSpecialistArtifactID(ctx context.Context, specialistID string) (string, error) {
	const q = `
		SELECT id
		FROM agent_core.specialist_artifacts
		WHERE specialist_id = $1 AND activation_status = 'current'
		ORDER BY updated_at DESC
		LIMIT 1`
	var id string
	if err := s.db.QueryRow(ctx, q, specialistID).Scan(&id); err != nil {
		return "", fmt.Errorf("get current specialist artifact id: %w", err)
	}
	return id, nil
}

func (s *PostgresStore) CreateCandidateSpecialistArtifact(ctx context.Context, in CandidateSpecialistArtifactInput) (string, error) {
	const q = `
		INSERT INTO agent_core.specialist_artifacts (
			id,
			specialist_id,
			unit_id,
			role,
			model_ref,
			harness_config_ref,
			store_mode,
			store_bootstrap_ref,
			slot_bundle_ref,
			slot_version_hash,
			eval_suite_ref,
			activation_status,
			rollback_ref,
			metadata,
			updated_at
		) VALUES ($1,$2,$3,'specialist',$4,NULLIF($5,''),$6,NULLIF($7,''),NULLIF($8,''),NULLIF($9,''),NULLIF($10,''),'candidate','', $11, now())
		RETURNING id`
	artifactID := defaultCandidateArtifactID(in.SpecialistID, in.AbilityName)
	unitID := defaultCandidateUnitID(in.SpecialistID, in.AbilityName)
	metadata := map[string]interface{}{
		"parent_artifact_id": in.ParentArtifactID,
		"ability_name": in.AbilityName,
		"requested_by": defaultLifecycleActor(in.RequestedBy),
		"notes": strings.TrimSpace(in.Notes),
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("marshal candidate artifact metadata: %w", err)
	}
	var out string
	if err := s.db.QueryRow(ctx, q, artifactID, in.SpecialistID, unitID, in.ModelRef, in.HarnessConfigRef, in.StoreMode, in.StoreBootstrapRef, in.SlotBundleRef, in.SlotVersionHash, in.EvalSuiteRef, metadataBytes).Scan(&out); err != nil {
		return "", fmt.Errorf("create candidate specialist artifact: %w", err)
	}
	return out, nil
}

func (s *PostgresStore) PromoteExperiment(ctx context.Context, experimentID, approvedBy, reason string) error {
	return s.updateArtifactLifecycleForExperiment(ctx, experimentID, approvedBy, reason, true)
}

func (s *PostgresStore) RollbackExperiment(ctx context.Context, experimentID, approvedBy, reason string) error {
	return s.updateArtifactLifecycleForExperiment(ctx, experimentID, approvedBy, reason, false)
}

func (s *PostgresStore) updateArtifactLifecycleForExperiment(ctx context.Context, experimentID, approvedBy, reason string, promote bool) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin artifact lifecycle tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const lookupQ = `
		SELECT specialist_artifact_id, specialist_id
		FROM agent_core.ability_growth_experiments
		WHERE experiment_id::text = $1`
	var artifactID string
	var specialistID string
	if err := tx.QueryRow(ctx, lookupQ, experimentID).Scan(&artifactID, &specialistID); err != nil {
		return fmt.Errorf("lookup artifact lifecycle experiment: %w", err)
	}
	if strings.TrimSpace(artifactID) == "" {
		return fmt.Errorf("experiment %s has no candidate artifact", experimentID)
	}

	actor := defaultLifecycleActor(approvedBy)
	finalReason := strings.TrimSpace(reason)
	if finalReason == "" {
		if promote {
			finalReason = "candidate artifact promoted"
		} else {
			finalReason = "candidate artifact rolled back"
		}
	}

	if promote {
		if _, err := tx.Exec(ctx, `UPDATE agent_core.specialist_artifacts SET activation_status = 'superseded', updated_at = now() WHERE specialist_id = $1 AND activation_status = 'current'`, specialistID); err != nil {
			return fmt.Errorf("supersede current artifact: %w", err)
		}
		if _, err := tx.Exec(ctx, `UPDATE agent_core.specialist_artifacts SET activation_status = 'current', updated_at = now() WHERE id = $1`, artifactID); err != nil {
			return fmt.Errorf("promote candidate artifact: %w", err)
		}
		if _, err := tx.Exec(ctx, `UPDATE agent_core.ability_growth_experiments SET status = 'promoted' WHERE experiment_id::text = $1`, experimentID); err != nil {
			return fmt.Errorf("mark experiment promoted: %w", err)
		}
	} else {
		if _, err := tx.Exec(ctx, `UPDATE agent_core.specialist_artifacts SET activation_status = 'rolled_back', rollback_ref = $2, updated_at = now() WHERE id = $1`, artifactID, experimentID); err != nil {
			return fmt.Errorf("rollback candidate artifact: %w", err)
		}
		if _, err := tx.Exec(ctx, `UPDATE agent_core.ability_growth_experiments SET status = 'rolled_back' WHERE experiment_id::text = $1`, experimentID); err != nil {
			return fmt.Errorf("mark experiment rolled back: %w", err)
		}
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO agent_core.specialist_artifact_events (
			event_id,
			artifact_id,
			specialist_id,
			event_type,
			experiment_id,
			actor,
			reason,
			metadata
		) VALUES ($1,$2,$3,$4,$5,$6,$7,'{}'::jsonb)`,
		defaultSpecialistArtifactEventID(artifactID, lifecycleEventType(promote)),
		artifactID,
		specialistID,
		lifecycleEventType(promote),
		experimentID,
		actor,
		finalReason,
	); err != nil {
		return fmt.Errorf("record lifecycle event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit artifact lifecycle tx: %w", err)
	}
	return nil
}

func lifecycleEventType(promote bool) string {
	if promote {
		return "candidate_promoted"
	}
	return "candidate_rolled_back"
}

func defaultCandidateArtifactID(specialistID, abilityName string) string {
	sid := sanitizeLifecycleComponent(specialistID, "specialist")
	ability := sanitizeLifecycleComponent(abilityName, "ability")
	return fmt.Sprintf("artifact-%s-candidate-%s-%d", sid, ability, time.Now().UTC().UnixNano())
}

func defaultCandidateUnitID(specialistID, abilityName string) string {
	sid := sanitizeLifecycleComponent(specialistID, "specialist")
	ability := sanitizeLifecycleComponent(abilityName, "ability")
	return fmt.Sprintf("%s--candidate--%s--%d", sid, ability, time.Now().UTC().UnixNano())
}

func sanitizeLifecycleComponent(value, fallback string) string {
	clean := strings.TrimSpace(value)
	if clean == "" {
		clean = fallback
	}
	clean = strings.ReplaceAll(clean, " ", "-")
	clean = strings.ReplaceAll(clean, "/", "-")
	clean = strings.ReplaceAll(clean, "\\", "-")
	return clean
}

func defaultLifecycleActor(actor string) string {
	clean := strings.TrimSpace(actor)
	if clean == "" {
		return "harness:runtime"
	}
	return clean
}

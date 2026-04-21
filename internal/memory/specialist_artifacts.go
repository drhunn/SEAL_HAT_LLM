package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type SpecialistArtifactInput struct {
	ID                string
	SpecialistID      string
	UnitID            string
	Role              string
	ModelRef          string
	HarnessConfigRef  string
	StoreMode         string
	StoreBootstrapRef string
	SlotBundleRef     string
	SlotVersionHash   string
	EvalSuiteRef      string
	ActivationStatus  string
	RollbackRef       string
	MetadataJSON      map[string]interface{}
}

func (s *PostgresStore) UpsertSpecialistArtifact(ctx context.Context, in SpecialistArtifactInput) (string, error) {
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
		) VALUES ($1,$2,$3,$4,$5,NULLIF($6,''),$7,NULLIF($8,''),NULLIF($9,''),NULLIF($10,''),NULLIF($11,''),'current',NULLIF($12,''),$13,now())
		ON CONFLICT (specialist_id) WHERE activation_status = 'current' DO UPDATE SET
			unit_id = EXCLUDED.unit_id,
			role = EXCLUDED.role,
			model_ref = EXCLUDED.model_ref,
			harness_config_ref = EXCLUDED.harness_config_ref,
			store_mode = EXCLUDED.store_mode,
			store_bootstrap_ref = EXCLUDED.store_bootstrap_ref,
			slot_bundle_ref = EXCLUDED.slot_bundle_ref,
			slot_version_hash = EXCLUDED.slot_version_hash,
			eval_suite_ref = EXCLUDED.eval_suite_ref,
			activation_status = 'current',
			rollback_ref = EXCLUDED.rollback_ref,
			metadata = EXCLUDED.metadata,
			updated_at = now()
		RETURNING id`

	id := strings.TrimSpace(in.ID)
	if id == "" {
		id = defaultSpecialistArtifactID(in.SpecialistID, in.UnitID)
	}
	metadataBytes, err := json.Marshal(in.MetadataJSON)
	if err != nil {
		return "", fmt.Errorf("marshal specialist artifact metadata: %w", err)
	}

	var out string
	if err := s.db.QueryRow(
		ctx,
		q,
		id,
		in.SpecialistID,
		in.UnitID,
		in.Role,
		in.ModelRef,
		in.HarnessConfigRef,
		in.StoreMode,
		in.StoreBootstrapRef,
		in.SlotBundleRef,
		in.SlotVersionHash,
		in.EvalSuiteRef,
		in.RollbackRef,
		metadataBytes,
	).Scan(&out); err != nil {
		return "", fmt.Errorf("upsert specialist artifact: %w", err)
	}
	return out, nil
}

func defaultSpecialistArtifactID(specialistID, unitID string) string {
	sid := strings.TrimSpace(specialistID)
	uid := strings.TrimSpace(unitID)
	if sid == "" {
		sid = "specialist"
	}
	if uid == "" {
		uid = sid
	}
	sid = strings.ReplaceAll(sid, " ", "-")
	sid = strings.ReplaceAll(sid, "/", "-")
	uid = strings.ReplaceAll(uid, " ", "-")
	uid = strings.ReplaceAll(uid, "/", "-")
	return fmt.Sprintf("artifact-%s-%s", sid, uid)
}

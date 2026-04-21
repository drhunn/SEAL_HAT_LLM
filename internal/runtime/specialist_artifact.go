package runtime

import (
	"context"

	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unit"
)

func (s *Service) persistSpecialistArtifact(ctx context.Context, slotBundleRef, slotVersionHash string) {
	if s == nil || s.store == nil {
		return
	}
	role := string(unit.RoleSpecialist)
	if s.cfg.Runtime.SpecialistID == executors.ParentGeneralist.String() {
		role = string(unit.RoleParent)
	}
	artifactID, err := s.store.UpsertSpecialistArtifact(ctx, memory.SpecialistArtifactInput{
		SpecialistID:      s.cfg.Runtime.SpecialistID,
		UnitID:            s.cfg.Runtime.SpecialistID,
		Role:              role,
		ModelRef:          s.cfg.Runtime.SpecialistID,
		HarnessConfigRef:  s.cfg.Runtime.ConfigRoot,
		StoreMode:         string(unit.StoreModeSharedDSN),
		StoreBootstrapRef: "database.dsn",
		SlotBundleRef:     slotBundleRef,
		SlotVersionHash:   slotVersionHash,
		ActivationStatus:  "active",
		MetadataJSON:      map[string]interface{}{},
	})
	if err != nil {
		s.logger.Warn("specialist artifact persistence failed", "specialist_id", s.cfg.Runtime.SpecialistID, "err", err)
		return
	}
	if _, err := s.store.CreateSpecialistArtifactEvent(ctx, memory.SpecialistArtifactEventInput{
		ArtifactID:   artifactID,
		SpecialistID: s.cfg.Runtime.SpecialistID,
		EventType:    "registered",
		Actor:        s.cfg.Harness.DefaultCreatedBy,
		Reason:       "runtime startup registered current local model unit artifact",
		MetadataJSON: map[string]interface{}{
			"slot_bundle_ref":   slotBundleRef,
			"slot_version_hash": slotVersionHash,
		},
	}); err != nil {
		s.logger.Warn("specialist artifact event persistence failed", "specialist_id", s.cfg.Runtime.SpecialistID, "err", err)
	}
}

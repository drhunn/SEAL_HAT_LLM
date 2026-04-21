package memory

import (
	"context"
	"fmt"
	"strings"
)

func (s *PostgresStore) RecordSpecialistArtifactEventForExperiment(ctx context.Context, experimentID, eventType, actor, reason string) error {
	const q = `
		SELECT specialist_artifact_id, specialist_id
		FROM agent_core.ability_growth_experiments
		WHERE experiment_id::text = $1`
	var artifactID string
	var specialistID string
	if err := s.db.QueryRow(ctx, q, experimentID).Scan(&artifactID, &specialistID); err != nil {
		return fmt.Errorf("lookup specialist artifact for experiment event: %w", err)
	}
	if strings.TrimSpace(artifactID) == "" {
		return nil
	}
	_, err := s.CreateSpecialistArtifactEvent(ctx, SpecialistArtifactEventInput{
		ArtifactID:   artifactID,
		SpecialistID: specialistID,
		EventType:    eventType,
		ExperimentID: experimentID,
		Actor:        actor,
		Reason:       reason,
		MetadataJSON: map[string]interface{}{},
	})
	if err != nil {
		return fmt.Errorf("record specialist artifact event for experiment: %w", err)
	}
	return nil
}

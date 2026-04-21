package memory

import (
	"context"
	"fmt"
)

func (s *PostgresStore) GetSpecialistArtifactID(ctx context.Context, specialistID, unitID string) (string, error) {
	const q = `
		SELECT id
		FROM agent_core.specialist_artifacts
		WHERE specialist_id = $1 AND unit_id = $2`
	var id string
	if err := s.db.QueryRow(ctx, q, specialistID, unitID).Scan(&id); err != nil {
		return "", fmt.Errorf("get specialist artifact id: %w", err)
	}
	return id, nil
}

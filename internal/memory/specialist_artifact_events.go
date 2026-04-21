package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type SpecialistArtifactEventInput struct {
	EventID      string
	ArtifactID   string
	SpecialistID string
	EventType    string
	ExperimentID string
	Actor        string
	Reason       string
	MetadataJSON map[string]interface{}
}

func (s *PostgresStore) CreateSpecialistArtifactEvent(ctx context.Context, in SpecialistArtifactEventInput) (string, error) {
	const q = `
		INSERT INTO agent_core.specialist_artifact_events (
			event_id,
			artifact_id,
			specialist_id,
			event_type,
			experiment_id,
			actor,
			reason,
			metadata
		) VALUES ($1,$2,$3,$4,NULLIF($5,''),$6,NULLIF($7,''),$8)
		RETURNING event_id`

	eventID := strings.TrimSpace(in.EventID)
	if eventID == "" {
		eventID = defaultSpecialistArtifactEventID(in.ArtifactID, in.EventType)
	}
	metadataBytes, err := json.Marshal(in.MetadataJSON)
	if err != nil {
		return "", fmt.Errorf("marshal specialist artifact event metadata: %w", err)
	}
	var out string
	if err := s.db.QueryRow(ctx, q, eventID, in.ArtifactID, in.SpecialistID, in.EventType, in.ExperimentID, in.Actor, in.Reason, metadataBytes).Scan(&out); err != nil {
		return "", fmt.Errorf("create specialist artifact event: %w", err)
	}
	return out, nil
}

func defaultSpecialistArtifactEventID(artifactID, eventType string) string {
	id := strings.TrimSpace(artifactID)
	if id == "" {
		id = "artifact"
	}
	typ := strings.TrimSpace(eventType)
	if typ == "" {
		typ = "event"
	}
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, "/", "-")
	typ = strings.ReplaceAll(typ, " ", "-")
	typ = strings.ReplaceAll(typ, "/", "-")
	return fmt.Sprintf("artifact-event-%s-%s-%d", id, typ, time.Now().UTC().UnixNano())
}

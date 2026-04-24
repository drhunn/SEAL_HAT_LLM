package lifecycle

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type ArtifactTransitionInput struct {
	ArtifactID   string
	ToState      ArtifactState
	Actor        string
	Reason       string
	ExperimentID string
}

type ArtifactTransitionResult struct {
	ArtifactID   string
	SpecialistID string
	FromState    ArtifactState
	ToState      ArtifactState
	EventType    ArtifactEventType
	EventID      string
}

func (s *Service) TransitionArtifact(ctx context.Context, in ArtifactTransitionInput) (*ArtifactTransitionResult, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("lifecycle service database is required")
	}
	artifactID := strings.TrimSpace(in.ArtifactID)
	if artifactID == "" {
		return nil, fmt.Errorf("artifact id is required")
	}
	actor := strings.TrimSpace(in.Actor)
	if actor == "" {
		return nil, fmt.Errorf("actor is required")
	}
	if _, err := NormalizeArtifactState(string(in.ToState)); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin artifact transition: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var specialistID string
	var currentStatus string
	if err := tx.QueryRow(ctx, `
		SELECT specialist_id, activation_status
		FROM agent_core.specialist_artifacts
		WHERE id=$1
		FOR UPDATE`, artifactID).Scan(&specialistID, &currentStatus); err != nil {
		return nil, fmt.Errorf("load artifact for transition: %w", err)
	}
	fromState, err := NormalizeArtifactState(currentStatus)
	if err != nil {
		return nil, err
	}
	eventType, err := ValidateArtifactTransition(fromState, in.ToState)
	if err != nil {
		return nil, err
	}

	if in.ToState == ArtifactStateCurrent {
		if _, err := tx.Exec(ctx, `
			UPDATE agent_core.specialist_artifacts
			SET activation_status='archived', updated_at=now()
			WHERE specialist_id=$1 AND activation_status='current' AND id<>$2`, specialistID, artifactID); err != nil {
			return nil, fmt.Errorf("archive previous current artifact: %w", err)
		}
	}

	if _, err := tx.Exec(ctx, `
		UPDATE agent_core.specialist_artifacts
		SET activation_status=$2, updated_at=now()
		WHERE id=$1`, artifactID, string(in.ToState)); err != nil {
		return nil, fmt.Errorf("update artifact state: %w", err)
	}

	eventID := artifactEventID(artifactID, eventType)
	if _, err := tx.Exec(ctx, `
		INSERT INTO agent_core.specialist_artifact_events (
			event_id, artifact_id, specialist_id, event_type, experiment_id, actor, reason, metadata
		) VALUES ($1,$2,$3,$4,$5,$6,$7,'{}'::jsonb)`,
		eventID,
		artifactID,
		specialistID,
		string(eventType),
		nullText(in.ExperimentID),
		actor,
		nullText(in.Reason),
	); err != nil {
		return nil, fmt.Errorf("insert artifact event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit artifact transition: %w", err)
	}
	if s.logger != nil {
		s.logger.InfoContext(ctx, "artifact lifecycle transition", "artifact_id", artifactID, "specialist_id", specialistID, "from", fromState, "to", in.ToState, "event", eventType)
	}
	return &ArtifactTransitionResult{ArtifactID: artifactID, SpecialistID: specialistID, FromState: fromState, ToState: in.ToState, EventType: eventType, EventID: eventID}, nil
}

func (s *Service) PromoteArtifact(ctx context.Context, artifactID, actor, reason string) (*ArtifactTransitionResult, error) {
	return s.TransitionArtifact(ctx, ArtifactTransitionInput{ArtifactID: artifactID, ToState: ArtifactStateCurrent, Actor: actor, Reason: reason})
}

func (s *Service) RollbackArtifact(ctx context.Context, artifactID, actor, reason string) (*ArtifactTransitionResult, error) {
	return s.TransitionArtifact(ctx, ArtifactTransitionInput{ArtifactID: artifactID, ToState: ArtifactStateRolledBack, Actor: actor, Reason: reason})
}

func artifactEventID(artifactID string, eventType ArtifactEventType) string {
	return fmt.Sprintf("artifact-event-%s-%s-%d", strings.TrimSpace(artifactID), eventType, time.Now().UTC().UnixNano())
}

func nullText(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}

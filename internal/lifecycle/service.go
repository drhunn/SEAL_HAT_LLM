package lifecycle

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewService(db *pgxpool.Pool, logger *slog.Logger) *Service {
	return &Service{db: db, logger: logger}
}

func (s *Service) SetStatus(ctx context.Context, specialistID, status, reason string) error {
	const q = `
		UPDATE agent_core.specialists
		SET status = $2::specialist_status,
		    updated_at = now(),
		    notes = CASE
		        WHEN COALESCE(notes, '') = '' THEN $3
		        ELSE notes || E'\n' || $3
		    END
		WHERE specialist_id = $1`

	cmd, err := s.db.Exec(ctx, q, specialistID, status, reason)
	if err != nil {
		return fmt.Errorf("set specialist status: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("specialist not found: %s", specialistID)
	}

	s.logger.InfoContext(ctx, "specialist lifecycle state updated", "specialist_id", specialistID, "status", status, "reason", reason)
	return nil
}

package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DefaultSchemaBootstrapFiles = []string{
	"sql/postgres-ddl.sql",
	"sql/three-tier-memory.sql",
	"sql/functions.sql",
	"sql/seed.sql",
	"sql/verify.sql",
	"sql/multimodal-memory.sql",
	"sql/20260420_seal_den_growth.sql",
	"sql/20260421_route_episodes.sql",
	"sql/20260421_specialist_artifacts.sql",
	"sql/20260421_ability_growth_artifact_refs.sql",
	"sql/20260421_specialist_artifact_events.sql",
	"sql/20260421_candidate_artifact_lifecycle.sql",
}

func BootstrapSchema(ctx context.Context, pool *pgxpool.Pool) error {
	return ApplySQLFiles(ctx, pool, DefaultSchemaBootstrapFiles)
}

func ApplySQLFiles(ctx context.Context, pool *pgxpool.Pool, files []string) error {
	if pool == nil {
		return fmt.Errorf("database pool is required")
	}
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire bootstrap connection: %w", err)
	}
	defer conn.Release()

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read SQL bootstrap file %s: %w", file, err)
		}
		results := conn.Conn().PgConn().Exec(ctx, string(data))
		if _, err := results.ReadAll(); err != nil {
			return fmt.Errorf("execute SQL bootstrap file %s: %w", file, err)
		}
	}
	return nil
}

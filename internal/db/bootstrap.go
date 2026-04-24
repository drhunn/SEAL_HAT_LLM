package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const DefaultSchemaBootstrapRoot = "./sql"

var DefaultSchemaBootstrapFiles = []string{
	"postgres-ddl.sql",
	"three-tier-memory.sql",
	"functions.sql",
	"seed.sql",
	"verify.sql",
	"multimodal-memory.sql",
	"20260420_seal_den_growth.sql",
	"20260421_route_episodes.sql",
	"20260421_specialist_artifacts.sql",
	"20260421_ability_growth_artifact_refs.sql",
	"20260421_specialist_artifact_events.sql",
	"20260421_candidate_artifact_lifecycle.sql",
}

type SchemaBootstrapConfig struct {
	SQLRoot string
	Files   []string
}

func BootstrapSchema(ctx context.Context, pool *pgxpool.Pool) error {
	return BootstrapSchemaWithConfig(ctx, pool, SchemaBootstrapConfig{})
}

func BootstrapSchemaFromRoot(ctx context.Context, pool *pgxpool.Pool, sqlRoot string) error {
	return BootstrapSchemaWithConfig(ctx, pool, SchemaBootstrapConfig{SQLRoot: sqlRoot})
}

func BootstrapSchemaWithConfig(ctx context.Context, pool *pgxpool.Pool, cfg SchemaBootstrapConfig) error {
	files := cfg.Files
	if len(files) == 0 {
		files = DefaultSchemaBootstrapFiles
	}
	sqlRoot := strings.TrimSpace(cfg.SQLRoot)
	if sqlRoot == "" {
		sqlRoot = DefaultSchemaBootstrapRoot
	}
	return ApplySQLFiles(ctx, pool, sqlRoot, files)
}

func ApplySQLFiles(ctx context.Context, pool *pgxpool.Pool, sqlRoot string, files []string) error {
	if pool == nil {
		return fmt.Errorf("database pool is required")
	}
	sqlRoot = strings.TrimSpace(sqlRoot)
	if sqlRoot == "" {
		return fmt.Errorf("sql root is required")
	}
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire bootstrap connection: %w", err)
	}
	defer conn.Release()

	for _, file := range files {
		path, err := schemaBootstrapPath(sqlRoot, file)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read SQL bootstrap file %s: %w", path, err)
		}
		results := conn.Conn().PgConn().Exec(ctx, string(data))
		if _, err := results.ReadAll(); err != nil {
			return fmt.Errorf("execute SQL bootstrap file %s: %w", path, err)
		}
	}
	return nil
}

func schemaBootstrapPath(sqlRoot, file string) (string, error) {
	file = strings.TrimSpace(file)
	if file == "" {
		return "", fmt.Errorf("sql bootstrap file name is required")
	}
	if filepath.IsAbs(file) {
		return filepath.Clean(file), nil
	}
	return filepath.Join(sqlRoot, file), nil
}

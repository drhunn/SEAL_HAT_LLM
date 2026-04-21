package unit

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func storeBootstrapTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func testStoreSpec(mode StoreMode) Spec {
	return Spec{
		UnitID:        "csse-tool-development-specialist-01",
		ExecutorName:  "csse-tool-development-specialist-01",
		Role:          RoleSpecialist,
		ModelRef:      "csse-tool-development-specialist-01",
		SpecialistID:  "csse-tool-development-specialist-01",
		Namespace:     "memory.csse-tool-development-specialist-01",
		SlotsRoot:     "./specialists",
		ConfigRoot:    "./config",
		StoreMode:     mode,
		ToolPlaneMode: ToolPlaneInProcess,
	}
}

func TestOpenStoreSharedDSNMode(t *testing.T) {
	oldOpen := openSharedDSNStore
	defer func() { openSharedDSNStore = oldOpen }()
	openSharedDSNStore = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
		return &pgxpool.Pool{}, nil
	}

	cfg := &config.AppConfig{}
	cfg.Database.DSN = "postgres://example"

	handle, err := OpenStore(context.Background(), testStoreSpec(StoreModeSharedDSN), cfg, storeBootstrapTestLogger())
	if err != nil {
		t.Fatalf("OpenStore returned error: %v", err)
	}
	if handle == nil || handle.Pool == nil {
		t.Fatalf("expected store handle with pool")
	}
	if handle.BootstrapRef != "database.dsn" {
		t.Fatalf("expected bootstrap ref database.dsn, got %q", handle.BootstrapRef)
	}
}

func TestOpenStoreEmbeddedMode(t *testing.T) {
	oldOpen := openEmbeddedPostgresStore
	defer func() { openEmbeddedPostgresStore = oldOpen }()
	openEmbeddedPostgresStore = func(ctx context.Context, cfg db.EmbeddedPostgresConfig) (*db.EmbeddedPostgresHandle, error) {
		return &db.EmbeddedPostgresHandle{Pool: &pgxpool.Pool{}, Stop: func() error { return nil }}, nil
	}

	cfg := &config.AppConfig{}
	cfg.EmbeddedPostgres.DataDir = "./artifacts/embedded_postgres/test"
	cfg.EmbeddedPostgres.Port = 55432
	cfg.EmbeddedPostgres.User = "postgres"
	cfg.EmbeddedPostgres.DatabaseName = "postgres"

	handle, err := OpenStore(context.Background(), testStoreSpec(StoreModeEmbeddedPostgres), cfg, storeBootstrapTestLogger())
	if err != nil {
		t.Fatalf("OpenStore returned error: %v", err)
	}
	if handle == nil || handle.Pool == nil {
		t.Fatalf("expected embedded store handle with pool")
	}
	if handle.BootstrapRef != "embedded_postgres.data_dir" {
		t.Fatalf("expected bootstrap ref embedded_postgres.data_dir, got %q", handle.BootstrapRef)
	}
}

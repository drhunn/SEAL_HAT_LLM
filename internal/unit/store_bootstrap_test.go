package unit

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func storeBootstrapTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestOpenStoreSharedDSNMode(t *testing.T) {
	oldOpen := openSharedDSNStore
	defer func() { openSharedDSNStore = oldOpen }()
	openSharedDSNStore = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
		return &pgxpool.Pool{}, nil
	}

	cfg := &config.AppConfig{}
	cfg.Database.DSN = "postgres://example"
	spec := Spec{
		UnitID:        "csse-tool-development-specialist-01",
		ExecutorName:  "csse-tool-development-specialist-01",
		Role:          RoleSpecialist,
		ModelRef:      "csse-tool-development-specialist-01",
		SpecialistID:  "csse-tool-development-specialist-01",
		Namespace:     "memory.csse-tool-development-specialist-01",
		SlotsRoot:     "./specialists",
		ConfigRoot:    "./config",
		StoreMode:     StoreModeSharedDSN,
		ToolPlaneMode: ToolPlaneInProcess,
	}

	handle, err := OpenStore(context.Background(), spec, cfg, storeBootstrapTestLogger())
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

func TestOpenStoreEmbeddedModeReturnsExplicitError(t *testing.T) {
	cfg := &config.AppConfig{}
	spec := Spec{
		UnitID:        "csse-tool-development-specialist-01",
		ExecutorName:  "csse-tool-development-specialist-01",
		Role:          RoleSpecialist,
		ModelRef:      "csse-tool-development-specialist-01",
		SpecialistID:  "csse-tool-development-specialist-01",
		Namespace:     "memory.csse-tool-development-specialist-01",
		SlotsRoot:     "./specialists",
		ConfigRoot:    "./config",
		StoreMode:     StoreModeEmbeddedPostgres,
		ToolPlaneMode: ToolPlaneInProcess,
	}

	_, err := OpenStore(context.Background(), spec, cfg, storeBootstrapTestLogger())
	if err == nil {
		t.Fatalf("expected embedded postgres bootstrap error")
	}
}

package unit

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreHandle struct {
	Pool         *pgxpool.Pool
	BootstrapRef string
	Close        func()
}

var openSharedDSNStore = db.Open
var openEmbeddedPostgresStore = db.OpenEmbeddedPostgres

func OpenStore(ctx context.Context, spec Spec, cfg *config.AppConfig, logger *slog.Logger) (*StoreHandle, error) {
	if err := spec.Validate(); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	switch spec.StoreMode {
	case StoreModeSharedDSN:
		dsn := strings.TrimSpace(cfg.Database.DSN)
		if dsn == "" {
			return nil, fmt.Errorf("database dsn is required for shared_dsn store mode")
		}
		pool, err := openSharedDSNStore(ctx, dsn)
		if err != nil {
			return nil, err
		}
		if logger != nil {
			logger.Info("store bootstrap ready", "unit_id", spec.UnitID, "store_mode", spec.StoreMode, "bootstrap_ref", "database.dsn")
		}
		return &StoreHandle{
			Pool:         pool,
			BootstrapRef: "database.dsn",
			Close: func() {
				if pool != nil {
					pool.Close()
				}
			},
		}, nil
	case StoreModeEmbeddedPostgres:
		handle, err := openEmbeddedPostgresStore(ctx, db.EmbeddedPostgresConfig{
			DataDir:      strings.TrimSpace(cfg.EmbeddedPostgres.DataDir),
			Port:         cfg.EmbeddedPostgres.Port,
			User:         strings.TrimSpace(cfg.EmbeddedPostgres.User),
			DatabaseName: strings.TrimSpace(cfg.EmbeddedPostgres.DatabaseName),
			BinDir:       strings.TrimSpace(cfg.EmbeddedPostgres.BinDir),
		})
		if err != nil {
			return nil, err
		}
		if logger != nil {
			logger.Info("store bootstrap ready", "unit_id", spec.UnitID, "store_mode", spec.StoreMode, "bootstrap_ref", "embedded_postgres.data_dir")
		}
		return &StoreHandle{
			Pool:         handle.Pool,
			BootstrapRef: "embedded_postgres.data_dir",
			Close: func() {
				if handle.Stop != nil {
					_ = handle.Stop()
				}
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported store mode %q", spec.StoreMode)
	}
}

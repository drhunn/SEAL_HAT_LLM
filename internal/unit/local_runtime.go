package unit

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/evals"
	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/growth"
	"github.com/drhunn/SEAL_HAT_LLM/internal/harness"
	workflow "github.com/drhunn/SEAL_HAT_LLM/internal/harness/workflows"
	"github.com/drhunn/SEAL_HAT_LLM/internal/lifecycle"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modelhost"
	"github.com/drhunn/SEAL_HAT_LLM/internal/postmortem"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/drhunn/SEAL_HAT_LLM/internal/runtime"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slots"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slotsync"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocalRuntime struct {
	Spec       Spec
	Registry   *Registry
	Database   *pgxpool.Pool
	Store      *memory.PostgresStore
	Harness    *harness.Service
	Routing    *routing.Service
	Execution  *execution.Service
	Runtime    *runtime.Service
	closeStore func()
}

func NewLocalRuntime(spec Spec, cfg *config.AppConfig, database *pgxpool.Pool, logger *slog.Logger) (*LocalRuntime, error) {
	return newLocalRuntime(spec, cfg, database, logger, nil)
}

func NewBootstrappedLocalRuntime(ctx context.Context, spec Spec, cfg *config.AppConfig, logger *slog.Logger) (*LocalRuntime, error) {
	if err := spec.Validate(); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	storeHandle, err := OpenStore(ctx, spec, cfg, logger)
	if err != nil {
		return nil, err
	}
	if storeHandle == nil || storeHandle.Pool == nil {
		return nil, fmt.Errorf("store bootstrap returned no database pool")
	}
	return newLocalRuntime(spec, cfg, storeHandle.Pool, logger, storeHandle.Close)
}

func newLocalRuntime(spec Spec, cfg *config.AppConfig, database *pgxpool.Pool, logger *slog.Logger, closeStore func()) (*LocalRuntime, error) {
	if err := spec.Validate(); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if database == nil {
		return nil, fmt.Errorf("database pool is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	registry := NewRegistry(BuiltinSpecs(spec)...)
	store := memory.NewPostgresStore(database, logger)
	slotLoader := slots.NewFilesystemLoader(spec.SlotsRoot)
	slotSyncService := slotsync.NewService(slotLoader, logger)
	pmService := postmortem.NewService(store, logger, cfg.Harness.DefaultCreatedBy)
	evalService := evals.NewService(store, logger)
	lifecycleService := lifecycle.NewService(database, logger)
	growthService := growth.NewService(store, logger)
	recoveryPlanner := workflow.NewDefaultRecoveryPlanner()
	harnessService := harness.NewService(store, pmService, evalService, lifecycleService, recoveryPlanner, logger, cfg)
	routingService := routing.NewService(logger, registry)
	hostRegistry := modelhost.NewSimulatedRegistry(hostPrefix(spec))
	executionService := execution.NewService(logger, hostRegistry, registry)
	runtimeService := runtime.NewService(cfg, slotLoader, store, harnessService, routingService, executionService, growthService, slotSyncService, logger)

	return &LocalRuntime{
		Spec:       spec,
		Registry:   registry,
		Database:   database,
		Store:      store,
		Harness:    harnessService,
		Routing:    routingService,
		Execution:  executionService,
		Runtime:    runtimeService,
		closeStore: closeStore,
	}, nil
}

func (u *LocalRuntime) Start(ctx context.Context) error {
	if u == nil || u.Runtime == nil {
		return fmt.Errorf("local runtime is not initialized")
	}
	return u.Runtime.Start(ctx)
}

func (u *LocalRuntime) Close() {
	if u == nil {
		return
	}
	if u.closeStore != nil {
		u.closeStore()
	}
}

func hostPrefix(spec Spec) string {
	candidate := strings.TrimSpace(strings.ToLower(spec.UnitID))
	candidate = strings.ReplaceAll(candidate, " ", "-")
	candidate = strings.ReplaceAll(candidate, "/", "-")
	if candidate == "" {
		return "local"
	}
	return candidate
}

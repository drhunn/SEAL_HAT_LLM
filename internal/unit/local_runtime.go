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
	Spec      Spec
	Registry  *Registry
	Store     *memory.PostgresStore
	Harness   *harness.Service
	Routing   *routing.Service
	Execution *execution.Service
	Runtime   *runtime.Service
}

func NewLocalRuntime(spec Spec, cfg *config.AppConfig, database *pgxpool.Pool, logger *slog.Logger) (*LocalRuntime, error) {
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
		Spec:      spec,
		Registry:  registry,
		Store:     store,
		Harness:   harnessService,
		Routing:   routingService,
		Execution: executionService,
		Runtime:   runtimeService,
	}, nil
}

func (u *LocalRuntime) Start(ctx context.Context) error {
	if u == nil || u.Runtime == nil {
		return fmt.Errorf("local runtime is not initialized")
	}
	return u.Runtime.Start(ctx)
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

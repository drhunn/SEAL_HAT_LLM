package runtime_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/db"
	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modelhost"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	rt "github.com/drhunn/SEAL_HAT_LLM/internal/runtime"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unit"
)

func TestProcessTaskCarriesPreferredUnitIDThroughRoutingAndExecution(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN is required for runtime task processor integration test")
	}

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	pool, err := db.Open(ctx, dsn)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer pool.Close()

	const preferredUnitID = "csse-tool-development-specialist-01"
	cfg := &config.AppConfig{}
	cfg.Database.DSN = dsn
	cfg.Runtime.SpecialistID = preferredUnitID
	cfg.Runtime.Namespace = "memory.csse-tool-development-specialist-01"
	cfg.Runtime.DefaultPrimaryModality = "text"
	cfg.Runtime.AllowTextOnlyFallback = true
	cfg.Harness.DefaultCreatedBy = "test:runtime"

	unitSpec, err := unit.SpecFromConfig(cfg)
	if err != nil {
		t.Fatalf("build unit spec: %v", err)
	}
	registry := unit.NewRegistry(unit.BuiltinSpecs(unitSpec)...)

	hosts := modelhost.NewRegistry()
	hosts.Register(preferredUnitID, modelhost.NewStaticHost("preferred-unit-host", "preferred unit handled task"))

	store := memory.NewPostgresStore(pool, logger)
	routingService := routing.NewService(logger, registry)
	executionService := execution.NewService(logger, hosts, registry)
	service := rt.NewService(cfg, nil, store, nil, routingService, executionService, nil, nil, logger)

	result, err := service.ProcessTask(ctx, rt.Task{
		ID:                "test-preferred-unit-runtime-path",
		Summary:           "verify preferred unit survives process task",
		Class:             "analysis",
		PreferredUnitID:   preferredUnitID,
		AllowTextOnlyFallback: true,
		Prompt:            "route and execute through the preferred unit",
	})
	if err != nil {
		t.Fatalf("ProcessTask returned error: %v", err)
	}
	if result == nil {
		t.Fatalf("ProcessTask returned nil result")
	}
	if result.RoutingDecision.TargetUnitID != preferredUnitID {
		t.Fatalf("routing target unit mismatch: got %q want %q", result.RoutingDecision.TargetUnitID, preferredUnitID)
	}
	if result.ExecutionResult.Plan.TargetUnitID != preferredUnitID {
		t.Fatalf("execution target unit mismatch: got %q want %q", result.ExecutionResult.Plan.TargetUnitID, preferredUnitID)
	}
	if result.ExecutionResult.Plan.ChosenExecutor != preferredUnitID {
		t.Fatalf("execution chosen executor mismatch: got %q want %q", result.ExecutionResult.Plan.ChosenExecutor, preferredUnitID)
	}
	if result.ExecutionResult.HostResult.HostName != "preferred-unit-host" {
		t.Fatalf("host mismatch: got %q want preferred-unit-host", result.ExecutionResult.HostResult.HostName)
	}
}

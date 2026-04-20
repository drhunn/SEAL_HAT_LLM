package runtime

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/db"
	"github.com/drhunn/SEAL_HAT_LLM/internal/evals"
	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/harness"
	workflow "github.com/drhunn/SEAL_HAT_LLM/internal/harness/workflows"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modelhost"
	"github.com/drhunn/SEAL_HAT_LLM/internal/postmortem"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestProcessTaskPersistsRuntimeArtifacts(t *testing.T) {
	service, pool, cleanup := newIntegrationRuntimeService(t, modelhost.NewSimulatedRegistry("itest-success"))
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	suffix := time.Now().UTC().UnixNano()
	task := Task{
		ID:      fmt.Sprintf("itest-success-%d", suffix),
		Summary: fmt.Sprintf("integration success task %d", suffix),
		Class:   "analysis",
		Prompt:  "Summarize the current task.",
	}

	result, err := service.ProcessTask(ctx, task)
	if err != nil {
		t.Fatalf("ProcessTask() error = %v", err)
	}
	if result == nil {
		t.Fatalf("expected task result")
	}
	if result.ExecutionResult.Plan.ChosenExecutor != modelhost.ExecutorParentGeneralist {
		t.Fatalf("unexpected executor: %q", result.ExecutionResult.Plan.ChosenExecutor)
	}

	assertCount(t, ctx, pool, "SELECT COUNT(*) FROM agent_core.routing_audit WHERE task_id = $1", 1, task.ID)
	assertCount(t, ctx, pool, "SELECT COUNT(*) FROM agent_core.memory_records WHERE specialist_id = $1 AND title = $2", 1, service.cfg.Runtime.SpecialistID, "Execution artifact: "+task.Summary)
	assertCount(t, ctx, pool, "SELECT COUNT(*) FROM agent_core.memory_postmortems WHERE specialist_id = $1 AND task_summary = $2", 0, service.cfg.Runtime.SpecialistID, task.Summary)
}

func TestProcessTaskCreatesPostmortemAndEvalOnExecutionFailure(t *testing.T) {
	service, pool, cleanup := newIntegrationRuntimeService(t, modelhost.NewRegistry())
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	suffix := time.Now().UTC().UnixNano()
	task := Task{
		ID:      fmt.Sprintf("itest-failure-%d", suffix),
		Summary: fmt.Sprintf("integration failure task %d", suffix),
		Class:   "analysis",
		Prompt:  "This task should fail because no host is registered.",
	}

	if _, err := service.ProcessTask(ctx, task); err == nil {
		t.Fatalf("expected ProcessTask() to fail")
	}

	assertCount(t, ctx, pool, "SELECT COUNT(*) FROM agent_core.routing_audit WHERE task_id = $1", 1, task.ID)
	assertCount(t, ctx, pool, "SELECT COUNT(*) FROM agent_core.memory_postmortems WHERE specialist_id = $1 AND task_summary = $2", 1, service.cfg.Runtime.SpecialistID, task.Summary)
	assertCount(t, ctx, pool, "SELECT COUNT(*) FROM agent_core.memory_eval_cases WHERE specialist_id = $1 AND case_title = $2", 1, service.cfg.Runtime.SpecialistID, "Regression for: "+task.Summary)
}

func newIntegrationRuntimeService(t *testing.T, hosts *modelhost.Registry) (*Service, *pgxpool.Pool, func()) {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN is not set")
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pool, err := db.Open(ctx, dsn)
	if err != nil {
		t.Fatalf("db.Open() error = %v", err)
	}

	cfg := &config.AppConfig{}
	cfg.Database.DSN = dsn
	cfg.Runtime.SpecialistID = "csse-tool-development-specialist-01"
	cfg.Runtime.Namespace = "memory.csse-tool-development-specialist-01"
	cfg.Runtime.DefaultPrimaryModality = "text"
	cfg.Runtime.AllowTextOnlyFallback = true
	cfg.Harness.AutoCreatePostmortems = true
	cfg.Harness.AutoUpdateHealth = true
	cfg.Harness.DefaultCreatedBy = "harness:runtime:test"

	store := memory.NewPostgresStore(pool, logger)
	pmService := postmortem.NewService(store, logger, cfg.Harness.DefaultCreatedBy)
	evalService := evals.NewService(store, logger)
	recoveryPlanner := workflow.NewDefaultRecoveryPlanner()
	harnessService := harness.NewService(store, pmService, evalService, nil, recoveryPlanner, logger, cfg)

	service := &Service{
		cfg:       cfg,
		store:     store,
		harness:   harnessService,
		routing:   routing.NewService(logger),
		execution: execution.NewService(logger, hosts),
		logger:    logger,
	}

	return service, pool, func() { pool.Close() }
}

func assertCount(t *testing.T, ctx context.Context, pool *pgxpool.Pool, query string, want int, args ...any) {
	t.Helper()
	var got int
	if err := pool.QueryRow(ctx, query, args...).Scan(&got); err != nil {
		t.Fatalf("query count failed: %v", err)
	}
	if got != want {
		t.Fatalf("unexpected count for %q: got %d want %d", query, got, want)
	}
}

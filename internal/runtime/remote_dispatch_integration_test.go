package runtime_test

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/db"
	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modelhost"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	rt "github.com/drhunn/SEAL_HAT_LLM/internal/runtime"
	"github.com/drhunn/SEAL_HAT_LLM/internal/taskdispatch"
	"github.com/drhunn/SEAL_HAT_LLM/internal/taskrpc"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unit"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestProcessTaskDispatchesToRemoteRPCServerAndRecordsRouteEpisode(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if strings.TrimSpace(dsn) == "" {
		t.Skip("TEST_DATABASE_DSN is required for remote dispatch integration test")
	}

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	pool, err := db.Open(ctx, dsn)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer pool.Close()

	parentUnitID := "remote-dispatch-parent-test"
	remoteUnitID := "remote-dispatch-specialist-test"
	parentNamespace := "memory.remote-dispatch-parent-test"
	remoteNamespace := "memory.remote-dispatch-specialist-test"
	if err := ensureRemoteDispatchSpecialist(ctx, pool, parentUnitID, parentNamespace); err != nil {
		t.Fatalf("ensure parent specialist: %v", err)
	}
	if err := ensureRemoteDispatchSpecialist(ctx, pool, remoteUnitID, remoteNamespace); err != nil {
		t.Fatalf("ensure remote specialist: %v", err)
	}

	socketPath := filepath.Join(t.TempDir(), "remote-taskrpc.sock")
	server := taskrpc.NewServer(socketPath, remoteIntegrationHandler{unitID: remoteUnitID})
	serverCtx, cancelServer := context.WithCancel(ctx)
	serveErr := make(chan error, 1)
	go func() { serveErr <- server.Serve(serverCtx) }()
	waitForRemoteDispatchSocket(t, socketPath)
	t.Cleanup(func() {
		cancelServer()
		select {
		case err := <-serveErr:
			if err != nil {
				t.Fatalf("task rpc server returned error: %v", err)
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for task rpc server shutdown")
		}
	})

	cfg := &config.AppConfig{}
	cfg.Database.DSN = dsn
	cfg.Runtime.SpecialistID = parentUnitID
	cfg.Runtime.Namespace = parentNamespace
	cfg.Runtime.DefaultPrimaryModality = "text"
	cfg.Runtime.AllowTextOnlyFallback = true
	cfg.TaskDispatch.RemoteUnitSockets = map[string]string{remoteUnitID: socketPath}
	cfg.Harness.DefaultCreatedBy = "test:remote-dispatch"

	parentSpec := unit.Spec{
		UnitID:        parentUnitID,
		ExecutorName:  parentUnitID,
		Role:          unit.RoleParent,
		ModelRef:      parentUnitID,
		SpecialistID:  parentUnitID,
		Namespace:     parentNamespace,
		SlotsRoot:     "./specialists",
		ConfigRoot:    "./config",
		StoreMode:     unit.StoreModeSharedDSN,
		ToolPlaneMode: unit.ToolPlaneInProcess,
	}
	remoteSpec := unit.Spec{
		UnitID:        remoteUnitID,
		ExecutorName:  remoteUnitID,
		Role:          unit.RoleSpecialist,
		ModelRef:      remoteUnitID,
		SpecialistID:  remoteUnitID,
		Namespace:     remoteNamespace,
		SlotsRoot:     "./specialists",
		ConfigRoot:    "./config",
		StoreMode:     unit.StoreModeSharedDSN,
		ToolPlaneMode: unit.ToolPlaneInProcess,
	}
	registry := unit.NewRegistry(unit.BuiltinSpecs(parentSpec)...)
	registry.Register(remoteSpec)

	store := memory.NewPostgresStore(pool, logger)
	routingService := routing.NewService(logger, registry)
	executionService := execution.NewService(logger, modelhost.NewRegistry(), registry)
	service := rt.NewService(cfg, nil, store, nil, routingService, executionService, nil, nil, logger)
	dispatcher, err := taskdispatch.New(taskdispatch.Options{ParentUnitID: parentUnitID, Resolver: taskdispatch.StaticSocketResolver(cfg.TaskDispatch.RemoteUnitSockets)})
	if err != nil {
		t.Fatalf("create dispatcher: %v", err)
	}
	service.SetRemoteDispatcher(dispatcher)

	taskID := fmt.Sprintf("remote-dispatch-e2e-%d", time.Now().UTC().UnixNano())
	defer func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM agent_core.route_episodes WHERE specialist_id=$1 AND task_id=$2`, parentUnitID, taskID)
	}()

	result, err := service.ProcessTask(ctx, rt.Task{
		ID:                    taskID,
		Summary:               "prove parent dispatches to remote task rpc server",
		Class:                 "analysis",
		PreferredUnitID:       remoteUnitID,
		AllowTextOnlyFallback: true,
		Prompt:                "remote dispatch proof",
	})
	if err != nil {
		t.Fatalf("ProcessTask returned error: %v", err)
	}
	if result == nil {
		t.Fatalf("ProcessTask returned nil result")
	}
	if result.ExecutionResult.Plan.ExecutionMode != "remote_rpc" {
		t.Fatalf("expected remote_rpc execution mode, got %q", result.ExecutionResult.Plan.ExecutionMode)
	}
	if result.ExecutionResult.Plan.TargetUnitID != remoteUnitID {
		t.Fatalf("expected target unit %q, got %q", remoteUnitID, result.ExecutionResult.Plan.TargetUnitID)
	}
	if result.ExecutionResult.HostResult.Output != "remote handled: remote dispatch proof" {
		t.Fatalf("unexpected remote output: %q", result.ExecutionResult.HostResult.Output)
	}

	var gotStatus, gotExecutionMode, gotTargetUnitID string
	var gotHandled bool
	if err := pool.QueryRow(ctx, `
		SELECT status, execution_mode, target_unit_id, execution_handled
		FROM agent_core.route_episodes
		WHERE specialist_id=$1 AND task_id=$2
		ORDER BY created_at DESC
		LIMIT 1`, parentUnitID, taskID).Scan(&gotStatus, &gotExecutionMode, &gotTargetUnitID, &gotHandled); err != nil {
		t.Fatalf("query remote route episode: %v", err)
	}
	if gotStatus != "succeeded" || gotExecutionMode != "remote_rpc" || gotTargetUnitID != remoteUnitID || !gotHandled {
		t.Fatalf("unexpected route episode: status=%q mode=%q target=%q handled=%v", gotStatus, gotExecutionMode, gotTargetUnitID, gotHandled)
	}
}

type remoteIntegrationHandler struct {
	unitID string
}

func (h remoteIntegrationHandler) RunTask(ctx context.Context, req taskrpc.RunTaskRequest) (taskrpc.RunTaskResponse, error) {
	return taskrpc.RunTaskResponse{
		TaskID:           req.TaskID,
		SpecialistUnitID: h.unitID,
		Status:           "ok",
		ResultSummary:    "remote handled: " + req.Prompt,
		OutputJSON:       "remote handled: " + req.Prompt,
		Confidence:       0.95,
	}, nil
}

func ensureRemoteDispatchSpecialist(ctx context.Context, pool *pgxpool.Pool, specialistID, namespace string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO agent_core.specialists (
			specialist_id, name, domain, role, lineage_parent_id, status, priority,
			memory_namespace, tool_policy_profile, postmortem_required, seal_enabled, harness_enabled
		) VALUES (
			$1, $1, 'test', 'test_specialist', 'test-parent',
			'active', 'test', $2, 'test_policy', true, true, true
		) ON CONFLICT (specialist_id) DO NOTHING`, specialistID, namespace)
	return err
}

func waitForRemoteDispatchSocket(t *testing.T, socketPath string) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("unix", socketPath, 50*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for socket %s", socketPath)
}

package memory_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/db"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCreateRouteEpisodePersistsOutcomeStatuses(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if strings.TrimSpace(dsn) == "" {
		t.Skip("TEST_DATABASE_DSN is required for route episode persistence test")
	}

	ctx := context.Background()
	pool, err := db.Open(ctx, dsn)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer pool.Close()

	store := memory.NewPostgresStore(pool, slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})))
	specialistID := "route-episode-test-specialist"
	namespace := "memory.route-episode-test-specialist"
	if err := ensureRouteEpisodeTestSpecialist(ctx, pool, specialistID, namespace); err != nil {
		t.Fatalf("ensure test specialist: %v", err)
	}

	prefix := fmt.Sprintf("route-episode-proof-%d", time.Now().UTC().UnixNano())
	defer func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM agent_core.route_episodes WHERE specialist_id=$1 AND task_id LIKE $2`, specialistID, prefix+"%")
	}()

	tests := []struct {
		name             string
		taskID           string
		executionMode    string
		chosenTarget     string
		targetUnitID     string
		chosenExecutor   string
		executionHandled bool
		executionHost    string
		status           string
		errorText        string
	}{
		{name: "local success", taskID: prefix + "-local-success", executionMode: "unimodal", chosenTarget: "local-executor", targetUnitID: specialistID, chosenExecutor: "local-executor", executionHandled: true, executionHost: "local-host", status: "succeeded"},
		{name: "local failure", taskID: prefix + "-local-failure", executionMode: "unimodal", chosenTarget: "local-executor", targetUnitID: specialistID, chosenExecutor: "local-executor", executionHandled: false, executionHost: "local-host", status: "failed", errorText: "local boom"},
		{name: "remote success", taskID: prefix + "-remote-success", executionMode: "remote_rpc", chosenTarget: "remote-executor", targetUnitID: "remote-unit", chosenExecutor: "remote-executor", executionHandled: true, executionHost: "remote-unit", status: "succeeded"},
		{name: "remote failure", taskID: prefix + "-remote-failure", executionMode: "remote_rpc", chosenTarget: "remote-executor", targetUnitID: "remote-unit", chosenExecutor: "remote-executor", executionHandled: false, executionHost: "remote-unit", status: "failed", errorText: "remote boom"},
		{name: "unhandled host", taskID: prefix + "-unhandled", executionMode: "unimodal", chosenTarget: "local-executor", targetUnitID: specialistID, chosenExecutor: "local-executor", executionHandled: false, executionHost: "local-host", status: "unhandled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := store.CreateRouteEpisode(ctx, memory.RouteEpisodeInput{
				Namespace:         namespace,
				SpecialistID:      specialistID,
				TaskID:            tt.taskID,
				TaskSummary:       tt.name,
				TaskClass:         "analysis",
				PrimaryModality:   "text",
				ChosenTarget:      tt.chosenTarget,
				TargetUnitID:      tt.targetUnitID,
				TargetRole:        "specialist",
				TargetModelRef:    tt.targetUnitID,
				ChosenExecutor:    tt.chosenExecutor,
				ExecutionMode:     tt.executionMode,
				Confidence:        0.9,
				WasFallback:       false,
				NeedsParentReview: false,
				RequiresFusion:    false,
				ExecutionHandled:  tt.executionHandled,
				ExecutionHost:     tt.executionHost,
				Status:            tt.status,
				ErrorText:         tt.errorText,
			})
			if err != nil {
				t.Fatalf("CreateRouteEpisode returned error: %v", err)
			}
			if strings.TrimSpace(id) == "" {
				t.Fatalf("expected route episode id")
			}

			var gotStatus, gotExecutionMode, gotTargetUnitID, gotErrorText string
			var gotHandled bool
			if err := pool.QueryRow(ctx, `
				SELECT status, execution_mode, target_unit_id, execution_handled, COALESCE(error_text, '')
				FROM agent_core.route_episodes
				WHERE id=$1`, id).Scan(&gotStatus, &gotExecutionMode, &gotTargetUnitID, &gotHandled, &gotErrorText); err != nil {
				t.Fatalf("query route episode: %v", err)
			}
			if gotStatus != tt.status || gotExecutionMode != tt.executionMode || gotTargetUnitID != tt.targetUnitID || gotHandled != tt.executionHandled || gotErrorText != tt.errorText {
				t.Fatalf("unexpected route episode row: status=%q mode=%q target=%q handled=%v error=%q", gotStatus, gotExecutionMode, gotTargetUnitID, gotHandled, gotErrorText)
			}
		})
	}
}

func ensureRouteEpisodeTestSpecialist(ctx context.Context, pool *pgxpool.Pool, specialistID, namespace string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO agent_core.specialists (
			specialist_id, name, domain, role, lineage_parent_id, status, priority,
			memory_namespace, tool_policy_profile, postmortem_required, seal_enabled, harness_enabled
		) VALUES (
			$1, 'Route Episode Test Specialist', 'test', 'test_specialist', 'test-parent',
			'active', 'test', $2, 'test_policy', true, true, true
		) ON CONFLICT (specialist_id) DO NOTHING`, specialistID, namespace)
	return err
}

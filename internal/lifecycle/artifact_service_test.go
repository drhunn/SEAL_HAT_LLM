package lifecycle_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/db"
	"github.com/drhunn/SEAL_HAT_LLM/internal/lifecycle"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestArtifactLifecyclePromoteRollbackAndRejectIllegalTransitions(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if strings.TrimSpace(dsn) == "" {
		t.Skip("TEST_DATABASE_DSN is required for artifact lifecycle integration test")
	}

	ctx := context.Background()
	pool, err := db.Open(ctx, dsn)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer pool.Close()

	svc := lifecycle.NewService(pool, slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})))
	specialistID := "artifact-lifecycle-test-specialist"
	if err := ensureSpecialist(ctx, pool, specialistID, "memory."+specialistID); err != nil {
		t.Fatalf("ensure specialist: %v", err)
	}

	prefix := fmt.Sprintf("artifact-lifecycle-%d", time.Now().UTC().UnixNano())
	oldID := prefix + "-old"
	candidateID := prefix + "-candidate"
	defer func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM agent_core.specialist_artifacts WHERE specialist_id=$1 AND id LIKE $2`, specialistID, prefix+"%")
	}()

	mustInsertArtifact(t, ctx, pool, oldID, specialistID, "current")
	mustInsertArtifact(t, ctx, pool, candidateID, specialistID, "candidate")

	promoted, err := svc.PromoteArtifact(ctx, candidateID, "test:lifecycle", "candidate passed")
	if err != nil {
		t.Fatalf("PromoteArtifact returned error: %v", err)
	}
	if promoted.FromState != lifecycle.ArtifactStateCandidate || promoted.ToState != lifecycle.ArtifactStateCurrent {
		t.Fatalf("unexpected promotion result: %+v", promoted)
	}
	assertArtifactStatus(t, ctx, pool, candidateID, "current")
	assertArtifactStatus(t, ctx, pool, oldID, "archived")
	assertOneCurrentArtifact(t, ctx, pool, specialistID)
	assertArtifactEvent(t, ctx, pool, candidateID, "promoted")

	if _, err := svc.PromoteArtifact(ctx, oldID, "test:lifecycle", "illegal transition proof"); err == nil {
		t.Fatalf("expected archived to current transition to fail")
	}

	rolledBack, err := svc.RollbackArtifact(ctx, candidateID, "test:lifecycle", "return to prior bundle")
	if err != nil {
		t.Fatalf("RollbackArtifact returned error: %v", err)
	}
	if rolledBack.FromState != lifecycle.ArtifactStateCurrent || rolledBack.ToState != lifecycle.ArtifactStateRolledBack {
		t.Fatalf("unexpected rollback result: %+v", rolledBack)
	}
	assertArtifactStatus(t, ctx, pool, candidateID, "rolled_back")
	assertArtifactEvent(t, ctx, pool, candidateID, "rolled_back")
}

func ensureSpecialist(ctx context.Context, pool *pgxpool.Pool, specialistID, namespace string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO agent_core.specialists (
			specialist_id, name, domain, role, lineage_parent_id, status, priority,
			memory_namespace, tool_policy_profile, postmortem_required, seal_enabled, harness_enabled
		) VALUES ($1, $1, 'test', 'test_specialist', 'test-parent', 'active', 'test', $2, 'test_policy', true, true, true)
		ON CONFLICT (specialist_id) DO NOTHING`, specialistID, namespace)
	return err
}

func mustInsertArtifact(t *testing.T, ctx context.Context, pool *pgxpool.Pool, artifactID, specialistID, status string) {
	t.Helper()
	_, err := pool.Exec(ctx, `
		INSERT INTO agent_core.specialist_artifacts (
			id, specialist_id, unit_id, role, model_ref, store_mode, activation_status, metadata
		) VALUES ($1,$2,$3,'specialist',$3,'shared_dsn',$4,'{}'::jsonb)`, artifactID, specialistID, artifactID+"-unit", status)
	if err != nil {
		t.Fatalf("insert artifact: %v", err)
	}
}

func assertArtifactStatus(t *testing.T, ctx context.Context, pool *pgxpool.Pool, artifactID, want string) {
	t.Helper()
	var got string
	if err := pool.QueryRow(ctx, `SELECT activation_status FROM agent_core.specialist_artifacts WHERE id=$1`, artifactID).Scan(&got); err != nil {
		t.Fatalf("query artifact status: %v", err)
	}
	if got != want {
		t.Fatalf("artifact %s status=%q want %q", artifactID, got, want)
	}
}

func assertOneCurrentArtifact(t *testing.T, ctx context.Context, pool *pgxpool.Pool, specialistID string) {
	t.Helper()
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM agent_core.specialist_artifacts WHERE specialist_id=$1 AND activation_status='current'`, specialistID).Scan(&count); err != nil {
		t.Fatalf("query current count: %v", err)
	}
	if count != 1 {
		t.Fatalf("current artifact count=%d want 1", count)
	}
}

func assertArtifactEvent(t *testing.T, ctx context.Context, pool *pgxpool.Pool, artifactID, eventType string) {
	t.Helper()
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM agent_core.specialist_artifact_events WHERE artifact_id=$1 AND event_type=$2`, artifactID, eventType).Scan(&count); err != nil {
		t.Fatalf("query artifact event: %v", err)
	}
	if count < 1 {
		t.Fatalf("expected artifact event %q for %s", eventType, artifactID)
	}
}

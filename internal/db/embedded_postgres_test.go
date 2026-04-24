package db

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func resetEmbeddedPostgresTestHooks(t *testing.T) {
	oldRun := runEmbeddedCommand
	oldStatus := embeddedPostgresStatus
	oldOpenPool := openPool
	oldPingPool := pingPool
	oldNow := embeddedPostgresNow
	oldStaleAge := embeddedPostgresStaleLockAge
	oldBootstrap := bootstrapEmbeddedPostgresSchema
	bootstrapEmbeddedPostgresSchema = func(ctx context.Context, pool *pgxpool.Pool, sqlRoot string) error { return nil }
	t.Cleanup(func() {
		runEmbeddedCommand = oldRun
		embeddedPostgresStatus = oldStatus
		openPool = oldOpenPool
		pingPool = oldPingPool
		embeddedPostgresNow = oldNow
		embeddedPostgresStaleLockAge = oldStaleAge
		bootstrapEmbeddedPostgresSchema = oldBootstrap
	})
}

func TestOpenEmbeddedPostgresRejectsSecondOwner(t *testing.T) {
	resetEmbeddedPostgresTestHooks(t)

	runEmbeddedCommand = func(ctx context.Context, name string, args ...string) error { return nil }
	embeddedPostgresStatus = func(ctx context.Context, pgCtlPath, dataDir string) bool { return true }
	openPool = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) { return &pgxpool.Pool{}, nil }
	pingPool = func(ctx context.Context, pool *pgxpool.Pool) error { return nil }

	dataDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dataDir, "PG_VERSION"), []byte("16\n"), 0o644); err != nil {
		t.Fatalf("write PG_VERSION: %v", err)
	}

	h1, err := OpenEmbeddedPostgres(context.Background(), EmbeddedPostgresConfig{DataDir: dataDir, Port: 55432})
	if err != nil {
		t.Fatalf("first OpenEmbeddedPostgres returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dataDir, ".embedded_postgres.lock")); err != nil {
		t.Fatalf("expected lock file to exist: %v", err)
	}

	_, err = OpenEmbeddedPostgres(context.Background(), EmbeddedPostgresConfig{DataDir: dataDir, Port: 55432})
	if err == nil || !strings.Contains(err.Error(), "already in use") {
		t.Fatalf("expected already in use error, got %v", err)
	}
	if err := h1.Stop(); err != nil {
		t.Fatalf("stop returned error: %v", err)
	}
}

func TestEmbeddedPostgresStopIsIdempotent(t *testing.T) {
	resetEmbeddedPostgresTestHooks(t)

	var startCount, stopCount int
	runEmbeddedCommand = func(ctx context.Context, name string, args ...string) error {
		joined := strings.Join(args, " ")
		switch {
		case strings.Contains(joined, "start"):
			startCount++
		case strings.Contains(joined, "stop"):
			stopCount++
		}
		return nil
	}
	embeddedPostgresStatus = func(ctx context.Context, pgCtlPath, dataDir string) bool { return false }
	openPool = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) { return &pgxpool.Pool{}, nil }
	pingPool = func(ctx context.Context, pool *pgxpool.Pool) error { return nil }

	dataDir := t.TempDir()
	h, err := OpenEmbeddedPostgres(context.Background(), EmbeddedPostgresConfig{DataDir: dataDir, Port: 55432})
	if err != nil {
		t.Fatalf("OpenEmbeddedPostgres returned error: %v", err)
	}
	if startCount != 1 {
		t.Fatalf("expected one start, got %d", startCount)
	}
	if err := h.Stop(); err != nil {
		t.Fatalf("first stop returned error: %v", err)
	}
	if err := h.Stop(); err != nil {
		t.Fatalf("second stop returned error: %v", err)
	}
	if stopCount != 1 {
		t.Fatalf("expected one stop, got %d", stopCount)
	}
	if _, err := os.Stat(filepath.Join(dataDir, ".embedded_postgres.lock")); !os.IsNotExist(err) {
		t.Fatalf("expected lock file to be removed, got %v", err)
	}
}

func TestOpenEmbeddedPostgresCleansUpWhenPoolOpenFails(t *testing.T) {
	resetEmbeddedPostgresTestHooks(t)

	var startCount, stopCount int
	runEmbeddedCommand = func(ctx context.Context, name string, args ...string) error {
		joined := strings.Join(args, " ")
		switch {
		case strings.Contains(joined, "start"):
			startCount++
		case strings.Contains(joined, "stop"):
			stopCount++
		}
		return nil
	}
	embeddedPostgresStatus = func(ctx context.Context, pgCtlPath, dataDir string) bool { return false }
	openPool = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) { return nil, errors.New("boom") }

	dataDir := t.TempDir()
	_, err := OpenEmbeddedPostgres(context.Background(), EmbeddedPostgresConfig{DataDir: dataDir, Port: 55432})
	if err == nil || !strings.Contains(err.Error(), "open embedded postgres pool") {
		t.Fatalf("expected pool open error, got %v", err)
	}
	if startCount != 1 || stopCount != 1 {
		t.Fatalf("expected one start and one stop, got start=%d stop=%d", startCount, stopCount)
	}
	if _, err := os.Stat(filepath.Join(dataDir, ".embedded_postgres.lock")); !os.IsNotExist(err) {
		t.Fatalf("expected lock file cleanup, got %v", err)
	}
}

func TestOpenEmbeddedPostgresCleansUpWhenPingFails(t *testing.T) {
	resetEmbeddedPostgresTestHooks(t)

	var startCount, stopCount int
	runEmbeddedCommand = func(ctx context.Context, name string, args ...string) error {
		joined := strings.Join(args, " ")
		switch {
		case strings.Contains(joined, "start"):
			startCount++
		case strings.Contains(joined, "stop"):
			stopCount++
		}
		return nil
	}
	embeddedPostgresStatus = func(ctx context.Context, pgCtlPath, dataDir string) bool { return false }
	openPool = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) { return &pgxpool.Pool{}, nil }
	pingPool = func(ctx context.Context, pool *pgxpool.Pool) error { return errors.New("not ready") }

	dataDir := t.TempDir()
	_, err := OpenEmbeddedPostgres(context.Background(), EmbeddedPostgresConfig{DataDir: dataDir, Port: 55432})
	if err == nil || !strings.Contains(err.Error(), "ping embedded postgres pool") {
		t.Fatalf("expected ping error, got %v", err)
	}
	if startCount != 1 || stopCount != 1 {
		t.Fatalf("expected one start and one stop, got start=%d stop=%d", startCount, stopCount)
	}
	if _, err := os.Stat(filepath.Join(dataDir, ".embedded_postgres.lock")); !os.IsNotExist(err) {
		t.Fatalf("expected lock file cleanup, got %v", err)
	}
}

func TestOpenEmbeddedPostgresBootstrapsSchemaAfterPing(t *testing.T) {
	resetEmbeddedPostgresTestHooks(t)

	var bootstrapCount int
	var gotSQLRoot string
	runEmbeddedCommand = func(ctx context.Context, name string, args ...string) error { return nil }
	embeddedPostgresStatus = func(ctx context.Context, pgCtlPath, dataDir string) bool { return true }
	openPool = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) { return &pgxpool.Pool{}, nil }
	pingPool = func(ctx context.Context, pool *pgxpool.Pool) error { return nil }
	bootstrapEmbeddedPostgresSchema = func(ctx context.Context, pool *pgxpool.Pool, sqlRoot string) error {
		bootstrapCount++
		gotSQLRoot = sqlRoot
		return nil
	}

	dataDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dataDir, "PG_VERSION"), []byte("16\n"), 0o644); err != nil {
		t.Fatalf("write PG_VERSION: %v", err)
	}
	h, err := OpenEmbeddedPostgres(context.Background(), EmbeddedPostgresConfig{DataDir: dataDir, Port: 55432, SQLRoot: "/repo/sql"})
	if err != nil {
		t.Fatalf("OpenEmbeddedPostgres returned error: %v", err)
	}
	if bootstrapCount != 1 {
		t.Fatalf("expected one schema bootstrap call, got %d", bootstrapCount)
	}
	if gotSQLRoot != "/repo/sql" {
		t.Fatalf("expected configured SQL root, got %q", gotSQLRoot)
	}
	if err := h.Stop(); err != nil {
		t.Fatalf("stop returned error: %v", err)
	}
}

func TestOpenEmbeddedPostgresDefaultsSQLRoot(t *testing.T) {
	resetEmbeddedPostgresTestHooks(t)

	var gotSQLRoot string
	runEmbeddedCommand = func(ctx context.Context, name string, args ...string) error { return nil }
	embeddedPostgresStatus = func(ctx context.Context, pgCtlPath, dataDir string) bool { return true }
	openPool = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) { return &pgxpool.Pool{}, nil }
	pingPool = func(ctx context.Context, pool *pgxpool.Pool) error { return nil }
	bootstrapEmbeddedPostgresSchema = func(ctx context.Context, pool *pgxpool.Pool, sqlRoot string) error {
		gotSQLRoot = sqlRoot
		return nil
	}

	dataDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dataDir, "PG_VERSION"), []byte("16\n"), 0o644); err != nil {
		t.Fatalf("write PG_VERSION: %v", err)
	}
	h, err := OpenEmbeddedPostgres(context.Background(), EmbeddedPostgresConfig{DataDir: dataDir, Port: 55432})
	if err != nil {
		t.Fatalf("OpenEmbeddedPostgres returned error: %v", err)
	}
	if gotSQLRoot != DefaultSchemaBootstrapRoot {
		t.Fatalf("expected default SQL root %q, got %q", DefaultSchemaBootstrapRoot, gotSQLRoot)
	}
	if err := h.Stop(); err != nil {
		t.Fatalf("stop returned error: %v", err)
	}
}

func TestOpenEmbeddedPostgresCleansUpWhenBootstrapFails(t *testing.T) {
	resetEmbeddedPostgresTestHooks(t)

	var startCount, stopCount int
	runEmbeddedCommand = func(ctx context.Context, name string, args ...string) error {
		joined := strings.Join(args, " ")
		switch {
		case strings.Contains(joined, "start"):
			startCount++
		case strings.Contains(joined, "stop"):
			stopCount++
		}
		return nil
	}
	embeddedPostgresStatus = func(ctx context.Context, pgCtlPath, dataDir string) bool { return false }
	openPool = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) { return &pgxpool.Pool{}, nil }
	pingPool = func(ctx context.Context, pool *pgxpool.Pool) error { return nil }
	bootstrapEmbeddedPostgresSchema = func(ctx context.Context, pool *pgxpool.Pool, sqlRoot string) error { return errors.New("schema boom") }

	dataDir := t.TempDir()
	_, err := OpenEmbeddedPostgres(context.Background(), EmbeddedPostgresConfig{DataDir: dataDir, Port: 55432})
	if err == nil || !strings.Contains(err.Error(), "bootstrap embedded postgres schema") {
		t.Fatalf("expected bootstrap error, got %v", err)
	}
	if startCount != 1 || stopCount != 1 {
		t.Fatalf("expected one start and one stop, got start=%d stop=%d", startCount, stopCount)
	}
	if _, err := os.Stat(filepath.Join(dataDir, ".embedded_postgres.lock")); !os.IsNotExist(err) {
		t.Fatalf("expected lock file cleanup, got %v", err)
	}
}

func TestOpenEmbeddedPostgresReclaimsStaleLockWhenServerIsDown(t *testing.T) {
	resetEmbeddedPostgresTestHooks(t)

	now := time.Date(2026, 4, 21, 22, 30, 0, 0, time.UTC)
	embeddedPostgresNow = func() time.Time { return now }
	embeddedPostgresStaleLockAge = time.Minute
	runEmbeddedCommand = func(ctx context.Context, name string, args ...string) error { return nil }
	embeddedPostgresStatus = func(ctx context.Context, pgCtlPath, dataDir string) bool { return true }
	openPool = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) { return &pgxpool.Pool{}, nil }
	pingPool = func(ctx context.Context, pool *pgxpool.Pool) error { return nil }

	dataDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dataDir, "PG_VERSION"), []byte("16\n"), 0o644); err != nil {
		t.Fatalf("write PG_VERSION: %v", err)
	}
	staleLock := "pid=123\ncreated_at=" + now.Add(-2*time.Minute).Format(time.RFC3339Nano) + "\n"
	if err := os.WriteFile(filepath.Join(dataDir, ".embedded_postgres.lock"), []byte(staleLock), 0o600); err != nil {
		t.Fatalf("write stale lock: %v", err)
	}
	embeddedPostgresStatus = func(ctx context.Context, pgCtlPath, dataDir string) bool { return false }

	h, err := OpenEmbeddedPostgres(context.Background(), EmbeddedPostgresConfig{DataDir: dataDir, Port: 55432})
	if err != nil {
		t.Fatalf("expected stale lock to be reclaimed, got %v", err)
	}
	if h == nil || h.Stop == nil {
		t.Fatalf("expected handle with stop function")
	}
	if err := h.Stop(); err != nil {
		t.Fatalf("stop returned error: %v", err)
	}
}

func TestOpenEmbeddedPostgresDoesNotReclaimFreshLock(t *testing.T) {
	resetEmbeddedPostgresTestHooks(t)

	now := time.Date(2026, 4, 21, 22, 30, 0, 0, time.UTC)
	embeddedPostgresNow = func() time.Time { return now }
	embeddedPostgresStaleLockAge = time.Minute
	embeddedPostgresStatus = func(ctx context.Context, pgCtlPath, dataDir string) bool { return false }

	dataDir := t.TempDir()
	freshLock := "pid=123\ncreated_at=" + now.Add(-30*time.Second).Format(time.RFC3339Nano) + "\n"
	if err := os.WriteFile(filepath.Join(dataDir, ".embedded_postgres.lock"), []byte(freshLock), 0o600); err != nil {
		t.Fatalf("write fresh lock: %v", err)
	}

	_, err := OpenEmbeddedPostgres(context.Background(), EmbeddedPostgresConfig{DataDir: dataDir, Port: 55432})
	if err == nil || !strings.Contains(err.Error(), "already in use") {
		t.Fatalf("expected fresh lock rejection, got %v", err)
	}
}

package db

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenEmbeddedPostgresRejectsSecondOwner(t *testing.T) {
	oldRun := runEmbeddedCommand
	oldStatus := embeddedPostgresStatus
	oldOpenPool := openPool
	defer func() {
		runEmbeddedCommand = oldRun
		embeddedPostgresStatus = oldStatus
		openPool = oldOpenPool
	}()

	runEmbeddedCommand = func(ctx context.Context, name string, args ...string) error { return nil }
	embeddedPostgresStatus = func(ctx context.Context, pgCtlPath, dataDir string) bool { return true }
	openPool = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) { return nil, nil }

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
	if h1.Stop == nil {
		t.Fatalf("expected stop function")
	}
	if err := h1.Stop(); err != nil {
		t.Fatalf("stop returned error: %v", err)
	}
}

func TestEmbeddedPostgresStopIsIdempotent(t *testing.T) {
	oldRun := runEmbeddedCommand
	oldStatus := embeddedPostgresStatus
	oldOpenPool := openPool
	defer func() {
		runEmbeddedCommand = oldRun
		embeddedPostgresStatus = oldStatus
		openPool = oldOpenPool
	}()

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
	openPool = func(ctx context.Context, dsn string) (*pgxpool.Pool, error) { return nil, nil }

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
	oldRun := runEmbeddedCommand
	oldStatus := embeddedPostgresStatus
	oldOpenPool := openPool
	defer func() {
		runEmbeddedCommand = oldRun
		embeddedPostgresStatus = oldStatus
		openPool = oldOpenPool
	}()

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

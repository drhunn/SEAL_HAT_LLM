package db

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmbeddedPostgresConfig struct {
	DataDir      string
	Port         int
	User         string
	DatabaseName string
	BinDir       string
}

type EmbeddedPostgresHandle struct {
	Pool *pgxpool.Pool
	DSN  string
	Stop func() error
}

var commandContext = exec.CommandContext
var openPool = Open

func OpenEmbeddedPostgres(ctx context.Context, cfg EmbeddedPostgresConfig) (*EmbeddedPostgresHandle, error) {
	dataDir := strings.TrimSpace(cfg.DataDir)
	if dataDir == "" {
		return nil, fmt.Errorf("embedded postgres data dir is required")
	}
	user := strings.TrimSpace(cfg.User)
	if user == "" {
		user = "postgres"
	}
	databaseName := strings.TrimSpace(cfg.DatabaseName)
	if databaseName == "" {
		databaseName = "postgres"
	}
	if cfg.Port <= 0 {
		return nil, fmt.Errorf("embedded postgres port must be positive")
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create embedded postgres data dir: %w", err)
	}

	initdbPath := binaryPath(cfg.BinDir, "initdb")
	pgCtlPath := binaryPath(cfg.BinDir, "pg_ctl")
	initialized := fileExists(filepath.Join(dataDir, "PG_VERSION"))
	if !initialized {
		if err := runCommand(ctx, initdbPath, "-D", dataDir, "-U", user, "-A", "trust"); err != nil {
			return nil, fmt.Errorf("initdb: %w", err)
		}
	}

	startedHere := false
	if !postgresRunning(ctx, pgCtlPath, dataDir) {
		if err := runCommand(ctx, pgCtlPath, "-D", dataDir, "-o", postgresOptions(cfg.Port), "-w", "start"); err != nil {
			return nil, fmt.Errorf("start embedded postgres: %w", err)
		}
		startedHere = true
	}

	dsn := fmt.Sprintf("postgres://%s@127.0.0.1:%d/%s?sslmode=disable", user, cfg.Port, databaseName)
	pool, err := openPool(ctx, dsn)
	if err != nil {
		if startedHere {
			_ = runCommand(ctx, pgCtlPath, "-D", dataDir, "-m", "fast", "stop")
		}
		return nil, fmt.Errorf("open embedded postgres pool: %w", err)
	}

	stop := func() error {
		if pool != nil {
			pool.Close()
		}
		if !startedHere {
			return nil
		}
		return runCommand(context.Background(), pgCtlPath, "-D", dataDir, "-m", "fast", "stop")
	}

	return &EmbeddedPostgresHandle{Pool: pool, DSN: dsn, Stop: stop}, nil
}

func runCommand(ctx context.Context, name string, args ...string) error {
	cmd := commandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %v failed: %w: %s", name, args, err, strings.TrimSpace(string(output)))
	}
	return nil
}

func postgresRunning(ctx context.Context, pgCtlPath, dataDir string) bool {
	cmd := commandContext(ctx, pgCtlPath, "-D", dataDir, "status")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func postgresOptions(port int) string {
	return "-F -h 127.0.0.1 -p " + strconv.Itoa(port)
}

func binaryPath(binDir, name string) string {
	if strings.TrimSpace(binDir) == "" {
		return name
	}
	return filepath.Join(binDir, name)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

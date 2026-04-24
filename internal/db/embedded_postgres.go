package db

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmbeddedPostgresConfig struct {
	DataDir      string
	Port         int
	User         string
	DatabaseName string
	BinDir       string
	SQLRoot      string
}

type EmbeddedPostgresHandle struct {
	Pool *pgxpool.Pool
	DSN  string
	Stop func() error
}

var commandContext = exec.CommandContext
var openPool = Open
var pingPool = func(ctx context.Context, pool *pgxpool.Pool) error { return pool.Ping(ctx) }
var runEmbeddedCommand = runCommand
var embeddedPostgresStatus = postgresRunning
var embeddedPostgresNow = func() time.Time { return time.Now().UTC() }
var embeddedPostgresStaleLockAge = 2 * time.Minute
var bootstrapEmbeddedPostgresSchema = BootstrapSchemaFromRoot

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
	sqlRoot := strings.TrimSpace(cfg.SQLRoot)
	if sqlRoot == "" {
		sqlRoot = DefaultSchemaBootstrapRoot
	}
	if cfg.Port <= 0 {
		return nil, fmt.Errorf("embedded postgres port must be positive")
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create embedded postgres data dir: %w", err)
	}

	initdbPath := binaryPath(cfg.BinDir, "initdb")
	pgCtlPath := binaryPath(cfg.BinDir, "pg_ctl")
	releaseOwnership, err := acquireEmbeddedPostgresOwnership(ctx, dataDir, pgCtlPath)
	if err != nil {
		return nil, err
	}
	releaseIfNeeded := releaseOwnership
	defer func() {
		if releaseIfNeeded != nil {
			_ = releaseIfNeeded()
		}
	}()

	initialized := fileExists(filepath.Join(dataDir, "PG_VERSION"))
	if !initialized {
		if err := runEmbeddedCommand(ctx, initdbPath, "-D", dataDir, "-U", user, "-A", "trust"); err != nil {
			return nil, fmt.Errorf("initdb: %w", err)
		}
	}

	startedHere := false
	if !embeddedPostgresStatus(ctx, pgCtlPath, dataDir) {
		if err := runEmbeddedCommand(ctx, pgCtlPath, "-D", dataDir, "-o", postgresOptions(cfg.Port), "-w", "start"); err != nil {
			return nil, fmt.Errorf("start embedded postgres: %w", err)
		}
		startedHere = true
	}

	dsn := fmt.Sprintf("postgres://%s@127.0.0.1:%d/%s?sslmode=disable", user, cfg.Port, databaseName)
	pool, err := openPool(ctx, dsn)
	if err != nil {
		if startedHere {
			_ = runEmbeddedCommand(ctx, pgCtlPath, "-D", dataDir, "-m", "fast", "stop")
		}
		return nil, fmt.Errorf("open embedded postgres pool: %w", err)
	}
	if pool == nil {
		if startedHere {
			_ = runEmbeddedCommand(ctx, pgCtlPath, "-D", dataDir, "-m", "fast", "stop")
		}
		return nil, fmt.Errorf("open embedded postgres pool returned nil pool")
	}
	if err := pingPool(ctx, pool); err != nil {
		pool.Close()
		if startedHere {
			_ = runEmbeddedCommand(ctx, pgCtlPath, "-D", dataDir, "-m", "fast", "stop")
		}
		return nil, fmt.Errorf("ping embedded postgres pool: %w", err)
	}
	if bootstrapEmbeddedPostgresSchema != nil {
		if err := bootstrapEmbeddedPostgresSchema(ctx, pool, sqlRoot); err != nil {
			pool.Close()
			if startedHere {
				_ = runEmbeddedCommand(ctx, pgCtlPath, "-D", dataDir, "-m", "fast", "stop")
			}
			return nil, fmt.Errorf("bootstrap embedded postgres schema: %w", err)
		}
	}

	var stopOnce sync.Once
	var stopErr error
	stop := func() error {
		stopOnce.Do(func() {
			if pool != nil {
				pool.Close()
			}
			if startedHere {
				stopErr = runEmbeddedCommand(context.Background(), pgCtlPath, "-D", dataDir, "-m", "fast", "stop")
			}
			if releaseOwnership != nil {
				if err := releaseOwnership(); stopErr == nil && err != nil {
					stopErr = err
				}
				releaseOwnership = nil
			}
		})
		return stopErr
	}

	releaseIfNeeded = nil
	return &EmbeddedPostgresHandle{Pool: pool, DSN: dsn, Stop: stop}, nil
}

func acquireEmbeddedPostgresOwnership(ctx context.Context, dataDir, pgCtlPath string) (func() error, error) {
	lockPath := filepath.Join(dataDir, ".embedded_postgres.lock")
	for attempts := 0; attempts < 2; attempts++ {
		lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			metadata := fmt.Sprintf("pid=%d\ncreated_at=%s\n", os.Getpid(), embeddedPostgresNow().Format(time.RFC3339Nano))
			if _, err := lockFile.WriteString(metadata); err != nil {
				_ = lockFile.Close()
				_ = os.Remove(lockPath)
				return nil, fmt.Errorf("write embedded postgres lock: %w", err)
			}
			if err := lockFile.Close(); err != nil {
				_ = os.Remove(lockPath)
				return nil, fmt.Errorf("close embedded postgres lock: %w", err)
			}
			return func() error {
				if err := os.Remove(lockPath); err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("remove embedded postgres lock: %w", err)
				}
				return nil
			}, nil
		}
		if !os.IsExist(err) {
			return nil, fmt.Errorf("create embedded postgres lock: %w", err)
		}
		reclaimed, reclaimErr := reclaimStaleEmbeddedPostgresLock(ctx, lockPath, pgCtlPath, dataDir)
		if reclaimErr != nil {
			return nil, reclaimErr
		}
		if !reclaimed {
			return nil, fmt.Errorf("embedded postgres data dir %q is already in use", dataDir)
		}
	}
	return nil, fmt.Errorf("embedded postgres data dir %q is already in use", dataDir)
}

func reclaimStaleEmbeddedPostgresLock(ctx context.Context, lockPath, pgCtlPath, dataDir string) (bool, error) {
	if embeddedPostgresStatus(ctx, pgCtlPath, dataDir) {
		return false, nil
	}
	lockData, err := os.ReadFile(lockPath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, fmt.Errorf("read embedded postgres lock: %w", err)
	}
	createdAt, ok := parseEmbeddedPostgresLockCreatedAt(string(lockData))
	if !ok {
		return false, nil
	}
	if embeddedPostgresNow().Sub(createdAt) < embeddedPostgresStaleLockAge {
		return false, nil
	}
	if err := os.Remove(lockPath); err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("remove stale embedded postgres lock: %w", err)
	}
	return true, nil
}

func parseEmbeddedPostgresLockCreatedAt(lockData string) (time.Time, bool) {
	for _, line := range strings.Split(lockData, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "created_at=") {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(line, "created_at="))
		createdAt, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return time.Time{}, false
		}
		return createdAt, true
	}
	return time.Time{}, false
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

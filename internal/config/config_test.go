package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadEmbeddedPostgresConfigWithoutSharedDSN(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "runtime.toml")
	content := `[app]
name = "SEAL_HAT_LLM"

[runtime]
specialist_id = "csse-tool-development-specialist-01"
namespace = "memory.csse-tool-development-specialist-01"
store_mode = "embedded_postgres"

[embedded_postgres]
`
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.Runtime.StoreMode != "embedded_postgres" {
		t.Fatalf("expected embedded_postgres store mode, got %q", cfg.Runtime.StoreMode)
	}
	if cfg.EmbeddedPostgres.DataDir == "" {
		t.Fatalf("expected embedded postgres data dir default")
	}
	if cfg.EmbeddedPostgres.Port <= 0 {
		t.Fatalf("expected embedded postgres port default, got %d", cfg.EmbeddedPostgres.Port)
	}
	if cfg.EmbeddedPostgres.User == "" || cfg.EmbeddedPostgres.DatabaseName == "" {
		t.Fatalf("expected embedded postgres user and database defaults")
	}
	if cfg.EmbeddedPostgres.SQLRoot != "./sql" {
		t.Fatalf("expected embedded postgres sql root default, got %q", cfg.EmbeddedPostgres.SQLRoot)
	}
}

func TestLoadTaskRPCServerDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "runtime.toml")
	content := `[app]
name = "SEAL_HAT_LLM"

[database]
dsn = "postgres://postgres:postgres@localhost:5432/llm_harness?sslmode=disable"

[runtime]
specialist_id = "csse-tool-development-specialist-01"
namespace = "memory.csse-tool-development-specialist-01"
store_mode = "shared_dsn"
enable_task_rpc_server = true
`
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if !cfg.Runtime.EnableTaskRPCServer {
		t.Fatalf("expected task rpc server to be enabled")
	}
	if !strings.Contains(cfg.Runtime.TaskRPCSocketPath, "csse-tool-development-specialist-01.sock") {
		t.Fatalf("expected default task rpc socket path, got %q", cfg.Runtime.TaskRPCSocketPath)
	}
}

func TestLoadTaskDispatchSocketMap(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "runtime.toml")
	content := `[app]
name = "SEAL_HAT_LLM"

[database]
dsn = "postgres://postgres:postgres@localhost:5432/llm_harness?sslmode=disable"

[runtime]
specialist_id = "parent-unit-01"
namespace = "memory.parent-unit-01"
store_mode = "shared_dsn"

[task_dispatch.remote_unit_sockets]
" csse-tool-development-specialist-01 " = " ./artifacts/taskrpc/csse.sock "
"empty" = "   "
`
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if len(cfg.TaskDispatch.RemoteUnitSockets) != 1 {
		t.Fatalf("expected one cleaned socket mapping, got %+v", cfg.TaskDispatch.RemoteUnitSockets)
	}
	got := cfg.TaskDispatch.RemoteUnitSockets["csse-tool-development-specialist-01"]
	if got != "./artifacts/taskrpc/csse.sock" {
		t.Fatalf("unexpected socket mapping: %q", got)
	}
}

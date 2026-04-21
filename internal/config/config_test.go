package config

import (
	"os"
	"path/filepath"
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
}

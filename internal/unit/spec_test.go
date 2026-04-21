package unit

import (
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
)

func TestSpecFromConfigSpecialist(t *testing.T) {
	cfg := &config.AppConfig{}
	cfg.Runtime.SpecialistID = "csse-tool-development-specialist-01"
	cfg.Runtime.Namespace = "memory.csse-tool-development-specialist-01"
	cfg.Runtime.SlotsRoot = "./specialists"
	cfg.Runtime.ConfigRoot = "./config"

	spec, err := SpecFromConfig(cfg)
	if err != nil {
		t.Fatalf("SpecFromConfig returned error: %v", err)
	}
	if spec.Role != RoleSpecialist {
		t.Fatalf("expected specialist role, got %q", spec.Role)
	}
	if spec.UnitID != cfg.Runtime.SpecialistID {
		t.Fatalf("expected unit id %q, got %q", cfg.Runtime.SpecialistID, spec.UnitID)
	}
	if spec.StoreMode != StoreModeSharedDSN {
		t.Fatalf("expected shared DSN store mode, got %q", spec.StoreMode)
	}
	if spec.ToolPlaneMode != ToolPlaneInProcess {
		t.Fatalf("expected in-process tool plane mode, got %q", spec.ToolPlaneMode)
	}
}

func TestSpecFromConfigEmbeddedPostgres(t *testing.T) {
	cfg := &config.AppConfig{}
	cfg.Runtime.SpecialistID = "csse-tool-development-specialist-01"
	cfg.Runtime.Namespace = "memory.csse-tool-development-specialist-01"
	cfg.Runtime.SlotsRoot = "./specialists"
	cfg.Runtime.ConfigRoot = "./config"
	cfg.Runtime.StoreMode = "embedded_postgres"

	spec, err := SpecFromConfig(cfg)
	if err != nil {
		t.Fatalf("SpecFromConfig returned error: %v", err)
	}
	if spec.StoreMode != StoreModeEmbeddedPostgres {
		t.Fatalf("expected embedded postgres store mode, got %q", spec.StoreMode)
	}
}

func TestSpecFromConfigParent(t *testing.T) {
	cfg := &config.AppConfig{}
	cfg.Runtime.SpecialistID = executors.ParentGeneralist.String()
	cfg.Runtime.Namespace = "memory.parent-generalist"
	cfg.Runtime.SlotsRoot = "./specialists"
	cfg.Runtime.ConfigRoot = "./config"

	spec, err := SpecFromConfig(cfg)
	if err != nil {
		t.Fatalf("SpecFromConfig returned error: %v", err)
	}
	if spec.Role != RoleParent {
		t.Fatalf("expected parent role, got %q", spec.Role)
	}
	if !spec.IsParent() {
		t.Fatalf("expected IsParent to be true")
	}
}

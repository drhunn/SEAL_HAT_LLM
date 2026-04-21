package unitref

import (
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unit"
)

func TestForExecutorParent(t *testing.T) {
	target := ForExecutor(executors.ParentGeneralist.String())
	if target.UnitID != executors.ParentGeneralist.String() {
		t.Fatalf("expected unit id %q, got %q", executors.ParentGeneralist.String(), target.UnitID)
	}
	if target.Role != unit.RoleParent {
		t.Fatalf("expected parent role, got %q", target.Role)
	}
}

func TestForExecutorSpecialist(t *testing.T) {
	target := ForExecutor(executors.DocumentLayoutOCR.String())
	if target.UnitID != executors.DocumentLayoutOCR.String() {
		t.Fatalf("expected unit id %q, got %q", executors.DocumentLayoutOCR.String(), target.UnitID)
	}
	if target.Role != unit.RoleSpecialist {
		t.Fatalf("expected specialist role, got %q", target.Role)
	}
}

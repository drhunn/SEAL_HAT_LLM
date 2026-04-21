package unit

import (
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
)

func TestRegistryResolvesBuiltinsAndCurrentUnit(t *testing.T) {
	current := Spec{
		UnitID:        "csse-tool-development-specialist-01",
		ExecutorName:  "csse-tool-development-specialist-01",
		Role:          RoleSpecialist,
		ModelRef:      "csse-tool-development-specialist-01",
		SpecialistID:  "csse-tool-development-specialist-01",
		Namespace:     "memory.csse-tool-development-specialist-01",
		SlotsRoot:     "./specialists",
		ConfigRoot:    "./config",
		StoreMode:     StoreModeSharedDSN,
		ToolPlaneMode: ToolPlaneInProcess,
	}
	registry := NewRegistry(BuiltinSpecs(current)...)

	parent, ok := registry.ResolveExecutor(executors.ParentGeneralist.String())
	if !ok {
		t.Fatalf("expected builtin parent unit to resolve")
	}
	if parent.Role != RoleParent {
		t.Fatalf("expected parent role, got %q", parent.Role)
	}

	currentResolved, ok := registry.ResolveUnit(current.UnitID)
	if !ok {
		t.Fatalf("expected current unit to resolve")
	}
	if currentResolved.UnitID != current.UnitID {
		t.Fatalf("expected unit id %q, got %q", current.UnitID, currentResolved.UnitID)
	}
}

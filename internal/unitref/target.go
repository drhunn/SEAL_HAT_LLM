package unitref

import (
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unit"
)

type Target struct {
	UnitID       string
	Role         unit.Role
	ModelRef     string
	ExecutorName string
}

func ForExecutor(executorName string) Target {
	name := strings.TrimSpace(executorName)
	if name == "" {
		name = executors.ParentGeneralist.String()
	}

	target := Target{
		UnitID:       name,
		Role:         unit.RoleSpecialist,
		ModelRef:     name,
		ExecutorName: name,
	}
	if name == executors.ParentGeneralist.String() {
		target.Role = unit.RoleParent
	}
	return target
}

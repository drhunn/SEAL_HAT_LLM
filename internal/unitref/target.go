package unitref

import (
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
)

type Role string

const (
	RoleParent     Role = "parent"
	RoleSpecialist Role = "specialist"
)

type Target struct {
	UnitID       string
	Role         Role
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
		Role:         RoleSpecialist,
		ModelRef:     name,
		ExecutorName: name,
	}
	if name == executors.ParentGeneralist.String() {
		target.Role = RoleParent
	}
	return target
}

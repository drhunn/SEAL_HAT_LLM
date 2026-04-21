package unit

import (
	"fmt"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
)

type Role string

const (
	RoleParent     Role = "parent"
	RoleSpecialist Role = "specialist"
)

type StoreMode string

const (
	StoreModeSharedDSN       StoreMode = "shared_dsn"
	StoreModeEmbeddedPostgres StoreMode = "embedded_postgres"
)

type ToolPlaneMode string

const (
	ToolPlaneInProcess ToolPlaneMode = "in_process"
	ToolPlaneRPCShared ToolPlaneMode = "rpc_shared"
)

type Spec struct {
	UnitID        string
	Role          Role
	ModelRef      string
	SpecialistID  string
	Namespace     string
	SlotsRoot     string
	ConfigRoot    string
	StoreMode     StoreMode
	ToolPlaneMode ToolPlaneMode
}

func SpecFromConfig(cfg *config.AppConfig) (Spec, error) {
	if cfg == nil {
		return Spec{}, fmt.Errorf("config is required")
	}

	specialistID := strings.TrimSpace(cfg.Runtime.SpecialistID)
	role := RoleSpecialist
	if specialistID == executors.ParentGeneralist.String() {
		role = RoleParent
	}

	spec := Spec{
		UnitID:        specialistID,
		Role:          role,
		ModelRef:      specialistID,
		SpecialistID:  specialistID,
		Namespace:     strings.TrimSpace(cfg.Runtime.Namespace),
		SlotsRoot:     strings.TrimSpace(cfg.Runtime.SlotsRoot),
		ConfigRoot:    strings.TrimSpace(cfg.Runtime.ConfigRoot),
		StoreMode:     StoreModeSharedDSN,
		ToolPlaneMode: ToolPlaneInProcess,
	}
	if err := spec.Validate(); err != nil {
		return Spec{}, err
	}
	return spec, nil
}

func (s Spec) Validate() error {
	if strings.TrimSpace(s.UnitID) == "" {
		return fmt.Errorf("unit id is required")
	}
	if strings.TrimSpace(s.SpecialistID) == "" {
		return fmt.Errorf("specialist id is required")
	}
	if strings.TrimSpace(s.Namespace) == "" {
		return fmt.Errorf("namespace is required")
	}
	switch s.Role {
	case RoleParent, RoleSpecialist:
	default:
		return fmt.Errorf("invalid role %q", s.Role)
	}
	switch s.StoreMode {
	case StoreModeSharedDSN, StoreModeEmbeddedPostgres:
	default:
		return fmt.Errorf("invalid store mode %q", s.StoreMode)
	}
	switch s.ToolPlaneMode {
	case ToolPlaneInProcess, ToolPlaneRPCShared:
	default:
		return fmt.Errorf("invalid tool plane mode %q", s.ToolPlaneMode)
	}
	return nil
}

func (s Spec) IsParent() bool {
	return s.Role == RoleParent
}

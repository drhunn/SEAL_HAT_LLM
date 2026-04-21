package unit

import (
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/executors"
)

type Registry struct {
	byUnitID      map[string]Spec
	byExecutorKey map[string]Spec
}

func NewRegistry(specs ...Spec) *Registry {
	r := &Registry{
		byUnitID:      make(map[string]Spec),
		byExecutorKey: make(map[string]Spec),
	}
	for _, spec := range specs {
		r.Register(spec)
	}
	return r
}

func BuiltinSpecs(current Spec) []Spec {
	specs := []Spec{
		builtinSpec(executors.ParentGeneralist.String(), RoleParent),
		builtinSpec(executors.MultimodalFusion.String(), RoleSpecialist),
		builtinSpec(executors.ImageAnalysis.String(), RoleSpecialist),
		builtinSpec(executors.AudioTranscription.String(), RoleSpecialist),
		builtinSpec(executors.VideoUnderstanding.String(), RoleSpecialist),
		builtinSpec(executors.DocumentLayoutOCR.String(), RoleSpecialist),
	}
	if strings.TrimSpace(current.UnitID) != "" {
		specs = append(specs, current)
	}
	return specs
}

func (r *Registry) Register(spec Spec) {
	if r == nil {
		return
	}
	if err := spec.Validate(); err != nil {
		return
	}
	r.byUnitID[spec.UnitID] = spec
	r.byExecutorKey[normalizeExecutorKey(spec.ExecutorName)] = spec
}

func (r *Registry) ResolveExecutor(executorName string) (Spec, bool) {
	if r == nil {
		return Spec{}, false
	}
	spec, ok := r.byExecutorKey[normalizeExecutorKey(executorName)]
	return spec, ok
}

func (r *Registry) ResolveUnit(unitID string) (Spec, bool) {
	if r == nil {
		return Spec{}, false
	}
	spec, ok := r.byUnitID[strings.TrimSpace(unitID)]
	return spec, ok
}

func builtinSpec(executorName string, role Role) Spec {
	return Spec{
		UnitID:        executorName,
		ExecutorName:  executorName,
		Role:          role,
		ModelRef:      executorName,
		SpecialistID:  executorName,
		Namespace:     "registry." + normalizeExecutorKey(executorName),
		SlotsRoot:     "./specialists",
		ConfigRoot:    "./config",
		StoreMode:     StoreModeSharedDSN,
		ToolPlaneMode: ToolPlaneInProcess,
	}
}

func normalizeExecutorKey(name string) string {
	trimmed := strings.TrimSpace(strings.ToLower(name))
	trimmed = strings.ReplaceAll(trimmed, " ", "-")
	trimmed = strings.ReplaceAll(trimmed, "/", "-")
	return trimmed
}

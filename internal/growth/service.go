package growth

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
)

type Assessment struct {
	AbilityName       string
	GapSummary        string
	EvidenceSummary   string
	TriedMemoryFix    bool
	TriedRoutingFix   bool
	TriedPromptFix    bool
	PreferredSurface  string
	RequestedBy       string
	ParentApprovedBy  string
	HarnessVerifiedBy string
	Notes             string
}

type Result struct {
	Status       string
	ExperimentID string
}

type Service struct {
	store  *memory.PostgresStore
	logger *slog.Logger
}

func NewService(store *memory.PostgresStore, logger *slog.Logger) *Service {
	return &Service{store: store, logger: logger}
}

func (s *Service) StageExperiment(ctx context.Context, namespace, specialistID string, in Assessment) (*Result, error) {
	if strings.TrimSpace(in.AbilityName) == "" {
		return nil, fmt.Errorf("ability name is required")
	}
	if strings.TrimSpace(in.GapSummary) == "" {
		return nil, fmt.Errorf("gap summary is required")
	}

	status := "observed"
	nextAction := "continue evidence collection"
	if in.TriedMemoryFix && in.TriedRoutingFix && in.TriedPromptFix {
		status = "proposed"
		nextAction = "parent and harness can review governed growth experiment"
	}
	if strings.TrimSpace(in.ParentApprovedBy) != "" && strings.TrimSpace(in.HarnessVerifiedBy) != "" {
		status = "shadow"
		nextAction = "run shadow evaluation for new ability structure"
	}

	if err := s.store.UpsertAbilityLedger(ctx, memory.AbilityLedgerInput{
		Namespace:     namespace,
		SpecialistID:  specialistID,
		AbilityName:   in.AbilityName,
		MaturityStage: maturityStageForStatus(status),
		Score:         scoreForStatus(status),
		Evidence:      in.EvidenceSummary,
		LastAction:    "ability gap assessed",
		NextAction:    nextAction,
		UpdatedBy:     defaultActor(in.RequestedBy),
	}); err != nil {
		return nil, err
	}

	parentArtifactID := ""
	artifactID, err := s.store.GetCurrentSpecialistArtifactID(ctx, specialistID)
	if err != nil {
		s.logger.Warn("current specialist artifact lookup failed during growth staging", "specialist_id", specialistID, "err", err)
	} else {
		parentArtifactID = artifactID
	}

	candidateArtifactID, err := s.store.CreateCandidateSpecialistArtifact(ctx, memory.CandidateSpecialistArtifactInput{
		SpecialistID:      specialistID,
		ParentArtifactID:  parentArtifactID,
		AbilityName:       in.AbilityName,
		RequestedBy:       defaultActor(in.RequestedBy),
		ModelRef:          specialistID,
		HarnessConfigRef:  "",
		StoreMode:         "shared_dsn",
		StoreBootstrapRef: "database.dsn",
		SlotBundleRef:     "",
		SlotVersionHash:   "",
		EvalSuiteRef:      "",
		Notes:             in.Notes,
	})
	if err != nil {
		return nil, err
	}

	experimentID, err := s.store.CreateAbilityGrowthExperiment(ctx, memory.AbilityGrowthExperimentInput{
		Namespace:            namespace,
		SpecialistID:         specialistID,
		AbilityName:          in.AbilityName,
		GapSummary:           in.GapSummary,
		EvidenceSummary:      in.EvidenceSummary,
		PreferredSurface:     preferredSurface(in.PreferredSurface),
		Status:               status,
		RequestedBy:          defaultActor(in.RequestedBy),
		ParentApprovedBy:     in.ParentApprovedBy,
		HarnessVerifiedBy:    in.HarnessVerifiedBy,
		SpecialistArtifactID: candidateArtifactID,
		Notes:                in.Notes,
	})
	if err != nil {
		return nil, err
	}
	if candidateArtifactID != "" {
		if _, err := s.store.CreateSpecialistArtifactEvent(ctx, memory.SpecialistArtifactEventInput{
			ArtifactID:   candidateArtifactID,
			SpecialistID: specialistID,
			EventType:    "growth_staged",
			ExperimentID: experimentID,
			Actor:        defaultActor(in.RequestedBy),
			Reason:       firstNonEmpty(in.Notes, "candidate artifact staged for growth experiment"),
			MetadataJSON: map[string]interface{}{
				"status":             status,
				"ability_name":       in.AbilityName,
				"preferred_surface":  preferredSurface(in.PreferredSurface),
				"parent_artifact_id": parentArtifactID,
			},
		}); err != nil {
			s.logger.Warn("candidate artifact growth event persistence failed", "specialist_id", specialistID, "experiment_id", experimentID, "err", err)
		}
	}

	s.logger.InfoContext(ctx, "ability growth experiment staged",
		"specialist_id", specialistID,
		"ability", in.AbilityName,
		"status", status,
		"surface", preferredSurface(in.PreferredSurface),
		"parent_artifact_id", parentArtifactID,
		"candidate_artifact_id", candidateArtifactID,
		"experiment_id", experimentID,
	)

	return &Result{Status: status, ExperimentID: experimentID}, nil
}

func preferredSurface(surface string) string {
	if strings.TrimSpace(surface) == "" {
		return "adapter_family"
	}
	return surface
}

func defaultActor(actor string) string {
	if strings.TrimSpace(actor) == "" {
		return "harness:runtime"
	}
	return actor
}

func maturityStageForStatus(status string) string {
	switch status {
	case "shadow":
		return "adolescence"
	case "proposed":
		return "childhood"
	default:
		return "infancy"
	}
}

func scoreForStatus(status string) float64 {
	switch status {
	case "shadow":
		return 0.70
	case "proposed":
		return 0.45
	default:
		return 0.25
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

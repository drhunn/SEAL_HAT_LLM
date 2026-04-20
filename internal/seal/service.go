package seal

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

type ProposalSurface string

const (
	SurfaceSlotPatch       ProposalSurface = "slot_patch"
	SurfacePromptPatch     ProposalSurface = "prompt_patch"
	SurfaceRetrievalPatch  ProposalSurface = "retrieval_patch"
	SurfaceToolPatch       ProposalSurface = "tool_patch"
	SurfaceAdapterTuning   ProposalSurface = "adapter_tuning"
	SurfaceSpecialistSplit ProposalSurface = "specialist_split"
	SurfaceNewSpecialist   ProposalSurface = "new_specialist"
	SurfaceExpertExpansion ProposalSurface = "expert_expansion"
)

type GapCluster struct {
	ID                 string
	SpecialistID       string
	Category           string
	Surface            string
	Count              int
	PersistenceScore   float64
	SeverityScore      float64
	ReversibilityScore float64
	Summaries          []string
	EvidenceRefs       []string
}

type AdaptationProposal struct {
	ID               string
	SpecialistID     string
	ClusterID        string
	Surface          ProposalSurface
	Reason           string
	RequestedBy      string
	RequiresHarness  bool
	RequiresParent   bool
	RiskLevel        string
	RollbackRequired bool
}

type SignalReader interface {
	ListSignals(ctx context.Context, specialistID string, limit int) ([]telemetry.Signal, error)
}

type ProposalWriter interface {
	WriteProposal(ctx context.Context, proposal AdaptationProposal) error
}

type Service struct {
	signalReader SignalReader
	proposalWriter ProposalWriter
	logger       *slog.Logger
}

func NewService(signalReader SignalReader, proposalWriter ProposalWriter, logger *slog.Logger) *Service {
	return &Service{signalReader: signalReader, proposalWriter: proposalWriter, logger: logger}
}

func (s *Service) ReviewSpecialist(ctx context.Context, specialistID string, maxSignals int) ([]AdaptationProposal, error) {
	if strings.TrimSpace(specialistID) == "" {
		return nil, fmt.Errorf("specialist id is required")
	}
	if s.signalReader == nil {
		return nil, fmt.Errorf("signal reader is required")
	}
	if maxSignals <= 0 {
		maxSignals = 100
	}

	signals, err := s.signalReader.ListSignals(ctx, specialistID, maxSignals)
	if err != nil {
		return nil, err
	}
	clusters := ClusterSignals(signals)
	proposals := make([]AdaptationProposal, 0)
	for _, cluster := range clusters {
		if !shouldPropose(cluster) {
			continue
		}
		proposal := proposalForCluster(cluster)
		proposals = append(proposals, proposal)
		if s.proposalWriter != nil {
			if err := s.proposalWriter.WriteProposal(ctx, proposal); err != nil {
				return nil, err
			}
		}
	}

	s.logger.InfoContext(ctx, "seal review complete", "specialist_id", specialistID, "signal_count", len(signals), "proposal_count", len(proposals))
	return proposals, nil
}

func ClusterSignals(signals []telemetry.Signal) []GapCluster {
	buckets := make(map[string]*GapCluster)
	for _, signal := range signals {
		key := clusterKey(signal)
		cluster, ok := buckets[key]
		if !ok {
			cluster = &GapCluster{
				ID:           key,
				SpecialistID: signal.SpecialistID,
				Category:     signal.Category,
				Surface:      signal.Surface,
			}
			buckets[key] = cluster
		}
		cluster.Count++
		cluster.Summaries = append(cluster.Summaries, signal.Summary)
		cluster.EvidenceRefs = append(cluster.EvidenceRefs, signal.EvidenceRefs...)
		cluster.PersistenceScore = persistenceScore(cluster.Count)
		cluster.SeverityScore = max(cluster.SeverityScore, severityScore(signal.Severity))
		cluster.ReversibilityScore = reversibilityScore(signal.Surface)
	}

	out := make([]GapCluster, 0, len(buckets))
	for _, cluster := range buckets {
		cluster.EvidenceRefs = dedupe(cluster.EvidenceRefs)
		out = append(out, *cluster)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].PersistenceScore == out[j].PersistenceScore {
			return out[i].SeverityScore > out[j].SeverityScore
		}
		return out[i].PersistenceScore > out[j].PersistenceScore
	})
	return out
}

func shouldPropose(cluster GapCluster) bool {
	if cluster.Count >= 2 {
		return true
	}
	return cluster.SeverityScore >= 0.95
}

func proposalForCluster(cluster GapCluster) AdaptationProposal {
	surface := preferredSurface(cluster)
	requiresParent := surface == SurfaceNewSpecialist || surface == SurfaceSpecialistSplit || surface == SurfaceExpertExpansion || surface == SurfaceAdapterTuning
	riskLevel := "low"
	if requiresParent {
		riskLevel = "moderate"
	}
	if cluster.SeverityScore >= 0.95 && requiresParent {
		riskLevel = "high"
	}
	return AdaptationProposal{
		ID:               fmt.Sprintf("proposal|%s|%s", cluster.SpecialistID, cluster.ID),
		SpecialistID:     cluster.SpecialistID,
		ClusterID:        cluster.ID,
		Surface:          surface,
		Reason:           fmt.Sprintf("persistent gap detected in %s/%s across %d signals", cluster.Category, cluster.Surface, cluster.Count),
		RequestedBy:      "seal:runtime",
		RequiresHarness:  true,
		RequiresParent:   requiresParent,
		RiskLevel:        riskLevel,
		RollbackRequired: true,
	}
}

func preferredSurface(cluster GapCluster) ProposalSurface {
	switch {
	case strings.Contains(cluster.Surface, "retrieval") || strings.Contains(cluster.Category, "retrieval"):
		return SurfaceRetrievalPatch
	case strings.Contains(cluster.Surface, "tool") || strings.Contains(cluster.Category, "tool"):
		return SurfaceToolPatch
	case strings.Contains(cluster.Surface, "prompt") || strings.Contains(cluster.Category, "prompt"):
		return SurfacePromptPatch
	case strings.Contains(cluster.Surface, "routing"):
		return SurfaceSlotPatch
	case strings.Contains(cluster.Surface, "multimodal") || strings.Contains(cluster.Category, "capacity"):
		return SurfaceAdapterTuning
	default:
		return SurfaceSlotPatch
	}
}

func clusterKey(signal telemetry.Signal) string {
	return fmt.Sprintf("%s|%s|%s", signal.SpecialistID, signal.Category, signal.Surface)
}

func persistenceScore(count int) float64 {
	switch {
	case count >= 5:
		return 1.0
	case count == 4:
		return 0.9
	case count == 3:
		return 0.75
	case count == 2:
		return 0.6
	default:
		return 0.3
	}
}

func severityScore(severity telemetry.Severity) float64 {
	switch severity {
	case telemetry.SeverityHigh:
		return 1.0
	case telemetry.SeverityModerate:
		return 0.6
	default:
		return 0.3
	}
}

func reversibilityScore(surface string) float64 {
	switch strings.TrimSpace(strings.ToLower(surface)) {
	case "slot", "slots", "prompt", "tool", "retrieval":
		return 0.9
	case "adapter", "branch":
		return 0.6
	default:
		return 0.75
	}
}

func dedupe(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

package oversight

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

type ApprovalLevel string

const (
	ApprovalOperational ApprovalLevel = "operational"
	ApprovalHarness     ApprovalLevel = "harness"
	ApprovalParent      ApprovalLevel = "parent"
)

type PromotionDecision struct {
	ExperimentID string
	Status       string
	Reason       string
	ApprovedBy   string
}

type ProposalApprover interface {
	ApproveProposal(ctx context.Context, proposalID string, level ApprovalLevel, approvedBy string) error
}

type ExperimentWriter interface {
	PromoteExperiment(ctx context.Context, experimentID, approvedBy, reason string) error
	RollbackExperiment(ctx context.Context, experimentID, approvedBy, reason string) error
}

type ArtifactExperimentEventWriter interface {
	RecordSpecialistArtifactEventForExperiment(ctx context.Context, experimentID, eventType, actor, reason string) error
}

type Service struct {
	proposalApprover      ProposalApprover
	experimentWriter      ExperimentWriter
	artifactEventWriter   ArtifactExperimentEventWriter
	logger                *slog.Logger
}

func NewService(proposalApprover ProposalApprover, experimentWriter ExperimentWriter, logger *slog.Logger) *Service {
	service := &Service{proposalApprover: proposalApprover, experimentWriter: experimentWriter, logger: logger}
	if writer, ok := experimentWriter.(ArtifactExperimentEventWriter); ok {
		service.artifactEventWriter = writer
	}
	return service
}

func (s *Service) ApproveProposal(ctx context.Context, proposalID string, level ApprovalLevel, approvedBy string) error {
	if strings.TrimSpace(proposalID) == "" {
		return fmt.Errorf("proposal id is required")
	}
	if strings.TrimSpace(approvedBy) == "" {
		return fmt.Errorf("approved by is required")
	}
	if s.proposalApprover == nil {
		return fmt.Errorf("proposal approver is required")
	}
	if err := s.proposalApprover.ApproveProposal(ctx, proposalID, level, approvedBy); err != nil {
		return err
	}
	s.logger.InfoContext(ctx, "proposal approved", "proposal_id", proposalID, "level", level, "approved_by", approvedBy)
	return nil
}

func (s *Service) PromoteExperiment(ctx context.Context, experimentID, approvedBy, reason string) (*PromotionDecision, error) {
	if strings.TrimSpace(experimentID) == "" {
		return nil, fmt.Errorf("experiment id is required")
	}
	if strings.TrimSpace(approvedBy) == "" {
		return nil, fmt.Errorf("approved by is required")
	}
	if s.experimentWriter == nil {
		return nil, fmt.Errorf("experiment writer is required")
	}
	finalReason := defaultReason(reason, "experiment promoted after governed review")
	if err := s.experimentWriter.PromoteExperiment(ctx, experimentID, approvedBy, finalReason); err != nil {
		return nil, err
	}
	if s.artifactEventWriter != nil {
		if err := s.artifactEventWriter.RecordSpecialistArtifactEventForExperiment(ctx, experimentID, "experiment_promoted", approvedBy, finalReason); err != nil {
			s.logger.WarnContext(ctx, "artifact promotion event recording failed", "experiment_id", experimentID, "err", err)
		}
	}
	decision := &PromotionDecision{ExperimentID: experimentID, Status: "approved", Reason: finalReason, ApprovedBy: approvedBy}
	s.logger.InfoContext(ctx, "experiment promoted", "experiment_id", experimentID, "approved_by", approvedBy)
	return decision, nil
}

func (s *Service) RollbackExperiment(ctx context.Context, experimentID, approvedBy, reason string) error {
	if strings.TrimSpace(experimentID) == "" {
		return fmt.Errorf("experiment id is required")
	}
	if strings.TrimSpace(approvedBy) == "" {
		return fmt.Errorf("approved by is required")
	}
	if s.experimentWriter == nil {
		return fmt.Errorf("experiment writer is required")
	}
	finalReason := defaultReason(reason, "rollback required after governed review")
	if err := s.experimentWriter.RollbackExperiment(ctx, experimentID, approvedBy, finalReason); err != nil {
		return err
	}
	if s.artifactEventWriter != nil {
		if err := s.artifactEventWriter.RecordSpecialistArtifactEventForExperiment(ctx, experimentID, "experiment_rolled_back", approvedBy, finalReason); err != nil {
			s.logger.WarnContext(ctx, "artifact rollback event recording failed", "experiment_id", experimentID, "err", err)
		}
	}
	s.logger.InfoContext(ctx, "experiment rolled back", "experiment_id", experimentID, "approved_by", approvedBy)
	return nil
}

func defaultReason(reason, fallback string) string {
	if strings.TrimSpace(reason) == "" {
		return fallback
	}
	return strings.TrimSpace(reason)
}

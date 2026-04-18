package memory

import (
	"context"
	"fmt"
)

type EvalCaseInput struct {
	Namespace                string
	SpecialistID             string
	Category                 string
	CaseTitle                string
	PromptInput              string
	ExpectedBehavior         string
	ExpectedOutputOrCriteria string
	CreatedBy                string
}

type SelfEditCandidateInput struct {
	Namespace        string
	SpecialistID     string
	TargetSlot       string
	CandidateSummary string
	CandidateBody    string
	Rationale        string
	ProposedBy       string
}

type RoutingAuditInput struct {
	TaskID                string
	RoutedBy              string
	InitialClassifier     string
	TaskSummary           string
	TaskClass             string
	ChosenTarget          string
	Confidence            float64
	Impact                string
	WasFallback           bool
	FallbackReason        string
	WasOverride           bool
	OverrideBy            string
	MultiSpecialistReview bool
	Notes                 string
}

func (s *PostgresStore) StageEvalCase(ctx context.Context, in EvalCaseInput) (string, error) {
	const q = `
		INSERT INTO agent_core.memory_eval_cases (
			namespace,
			specialist_id,
			category,
			case_title,
			prompt_input,
			expected_behavior,
			expected_output_or_criteria,
			created_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING eval_case_id::text`

	var id string
	if err := s.db.QueryRow(ctx, q,
		in.Namespace,
		in.SpecialistID,
		in.Category,
		in.CaseTitle,
		in.PromptInput,
		in.ExpectedBehavior,
		in.ExpectedOutputOrCriteria,
		in.CreatedBy,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("stage eval case: %w", err)
	}
	return id, nil
}

func (s *PostgresStore) StageSelfEditCandidate(ctx context.Context, in SelfEditCandidateInput) (string, error) {
	const q = `
		INSERT INTO agent_core.memory_self_edit_candidates (
			namespace,
			specialist_id,
			target_slot,
			candidate_summary,
			candidate_body,
			rationale,
			proposed_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING candidate_id::text`

	var id string
	if err := s.db.QueryRow(ctx, q,
		in.Namespace,
		in.SpecialistID,
		in.TargetSlot,
		in.CandidateSummary,
		in.CandidateBody,
		in.Rationale,
		in.ProposedBy,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("stage self-edit candidate: %w", err)
	}
	return id, nil
}

func (s *PostgresStore) CreateRoutingAudit(ctx context.Context, in RoutingAuditInput) (string, error) {
	const q = `
		INSERT INTO agent_core.routing_audit (
			task_id,
			routed_by,
			initial_classifier,
			task_summary,
			task_class,
			chosen_target,
			confidence,
			impact,
			was_fallback,
			fallback_reason,
			was_override,
			override_by,
			multi_specialist_review,
			notes
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8::impact_level,$9,$10,$11,$12,$13,$14)
		RETURNING route_audit_id::text`

	var id string
	if err := s.db.QueryRow(ctx, q,
		in.TaskID,
		in.RoutedBy,
		in.InitialClassifier,
		in.TaskSummary,
		in.TaskClass,
		in.ChosenTarget,
		in.Confidence,
		in.Impact,
		in.WasFallback,
		in.FallbackReason,
		in.WasOverride,
		in.OverrideBy,
		in.MultiSpecialistReview,
		in.Notes,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("create routing audit: %w", err)
	}
	return id, nil
}

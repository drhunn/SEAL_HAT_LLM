package slots

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

type Bundle struct {
	Schema            string                  `toml:"schema"`
	Version           string                  `toml:"version"`
	SpecialistID      string                  `toml:"specialist_id,omitempty"`
	ParentID          string                  `toml:"parent_id,omitempty"`
	Status            string                  `toml:"status,omitempty"`
	Identity          IdentityBundle          `toml:"identity"`
	Soul              SoulBundle              `toml:"soul"`
	Agents            AgentsBundle            `toml:"agents"`
	Tools             ToolsBundle             `toml:"tools"`
	Skills            SkillsBundle            `toml:"skills"`
	Prompt            PromptBundle            `toml:"prompt"`
	Memory            MemoryBundle            `toml:"memory"`
	Heartbeat         HeartbeatBundle         `toml:"heartbeat"`
	Dreams            DreamsBundle            `toml:"dreams"`
	PostmortemTemplate PostmortemTemplateBundle `toml:"postmortem_template"`
	SlotPolicy        SlotPolicyBundle        `toml:"slot_policy"`
	SourceMap         map[string]string       `toml:"source_map,omitempty"`
}

type IdentityBundle struct {
	Name              string   `toml:"name,omitempty"`
	Aliases           []string `toml:"aliases,omitempty"`
	Domain            string   `toml:"domain,omitempty"`
	Role              string   `toml:"role,omitempty"`
	AuthorityBoundary string   `toml:"authority_boundary,omitempty"`
	AllowedScope      []string `toml:"allowed_scope,omitempty"`
	ForbiddenScope    []string `toml:"forbidden_scope,omitempty"`
	EscalateTo        string   `toml:"escalate_to,omitempty"`
}

type SoulBundle struct {
	Tone                    string   `toml:"tone,omitempty"`
	Values                  []string `toml:"values,omitempty"`
	CautionProfileHigh      []string `toml:"caution_profile_high,omitempty"`
	CautionProfileModerate  []string `toml:"caution_profile_moderate,omitempty"`
	DecisionStyle           []string `toml:"decision_style,omitempty"`
	RiskTolerance           string   `toml:"risk_tolerance,omitempty"`
}

type AgentsBundle struct {
	Mission                string   `toml:"mission,omitempty"`
	LaneBoundaries         []string `toml:"lane_boundaries,omitempty"`
	EscalationRules        []string `toml:"escalation_rules,omitempty"`
	DecisionRules          []string `toml:"decision_rules,omitempty"`
	SelfModificationRules  []string `toml:"self_modification_rules,omitempty"`
	MandatoryPostmortemRule string  `toml:"mandatory_postmortem_rule,omitempty"`
}

type ToolsBundle struct {
	AllowedTools    []string `toml:"allowed_tools,omitempty"`
	ToolOrdering    []string `toml:"tool_ordering,omitempty"`
	ToolConstraints []string `toml:"tool_constraints,omitempty"`
	FallbackBehavior []string `toml:"fallback_behavior,omitempty"`
}

type SkillsBundle struct {
	FailureRecoveryPatterns []string         `toml:"failure_recovery_patterns,omitempty"`
	Core                    []SkillRecord    `toml:"core,omitempty"`
	Playbooks               []PlaybookRecord `toml:"playbooks,omitempty"`
}

type SkillRecord struct {
	Name         string   `toml:"name,omitempty"`
	Trigger      string   `toml:"trigger,omitempty"`
	Procedure    []string `toml:"procedure,omitempty"`
	SuccessCheck []string `toml:"success_check,omitempty"`
}

type PlaybookRecord struct {
	Name      string   `toml:"name,omitempty"`
	Trigger   string   `toml:"trigger,omitempty"`
	Procedure []string `toml:"procedure,omitempty"`
}

type PromptBundle struct {
	TaskTemplates   []PromptTemplate `toml:"task_templates,omitempty"`
	CommandPatterns []PromptTemplate `toml:"command_patterns,omitempty"`
	Checklists      []Checklist      `toml:"checklists,omitempty"`
}

type PromptTemplate struct {
	Name    string   `toml:"name,omitempty"`
	Purpose string   `toml:"purpose,omitempty"`
	Body    []string `toml:"body,omitempty"`
}

type Checklist struct {
	Name  string   `toml:"name,omitempty"`
	Items []string `toml:"items,omitempty"`
}

type MemoryBundle struct {
	CurrentFocus      []string       `toml:"current_focus,omitempty"`
	DurableDecisions  []string       `toml:"durable_decisions,omitempty"`
	Pointers          []MemoryPointer `toml:"pointers,omitempty"`
	ActiveConstraints []string       `toml:"active_constraints,omitempty"`
}

type MemoryPointer struct {
	PointerName string `toml:"pointer_name,omitempty"`
	Backend     string `toml:"backend,omitempty"`
	QueryHint   string `toml:"query_hint,omitempty"`
}

type HeartbeatBundle struct {
	WatchTopics    []string `toml:"watch_topics,omitempty"`
	ApprovedSources []string `toml:"approved_sources,omitempty"`
	Cadence        []string `toml:"cadence,omitempty"`
	PromotionRules []string `toml:"promotion_rules,omitempty"`
	StopConditions []string `toml:"stop_conditions,omitempty"`
}

type DreamsBundle struct {
	DistilledPatterns  []string `toml:"distilled_patterns,omitempty"`
	CandidatePromotions []string `toml:"candidate_promotions,omitempty"`
	DiscardedNoise     []string `toml:"discarded_noise,omitempty"`
	OpenReviewItems    []string `toml:"open_review_items,omitempty"`
}

type PostmortemTemplateBundle struct {
	IncidentID                   string   `toml:"incident_id,omitempty"`
	Date                         string   `toml:"date,omitempty"`
	ModelID                      string   `toml:"model_id,omitempty"`
	ModelRole                    string   `toml:"model_role,omitempty"`
	TaskSummary                  string   `toml:"task_summary,omitempty"`
	ExpectedBehavior             string   `toml:"expected_behavior,omitempty"`
	ActualBehavior               string   `toml:"actual_behavior,omitempty"`
	WhatWentWrong                string   `toml:"what_went_wrong,omitempty"`
	FailureClassification        []string `toml:"failure_classification,omitempty"`
	RootCause                    string   `toml:"root_cause,omitempty"`
	MissedEvidenceOrStep         string   `toml:"missed_evidence_or_step,omitempty"`
	Preventable                  string   `toml:"preventable,omitempty"`
	Correction                   string   `toml:"correction,omitempty"`
	RecommendedSlotOrSystemChange []string `toml:"recommended_slot_or_system_change,omitempty"`
	RequiresHarnessReview        string   `toml:"requires_harness_review,omitempty"`
	RequiresParentReview         string   `toml:"requires_parent_review,omitempty"`
	RegressionTestNeeded         string   `toml:"regression_test_needed,omitempty"`
}

type SlotPolicyBundle struct {
	Constitutional []string `toml:"constitutional,omitempty"`
	Operational    []string `toml:"operational,omitempty"`
	ReviewPlane    []string `toml:"review_plane,omitempty"`
	Runtime        []string `toml:"runtime,omitempty"`
}

type BundleSummary struct {
	CoreSkillCount      int
	PlaybookCount       int
	AllowedToolCount    int
	CommandPatternCount int
	ChecklistCount      int
	MemoryPointerCount  int
	WatchTopicCount     int
	SourceMapCount      int
}

func (b *Bundle) EncodeTOML() ([]byte, error) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(b); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *Bundle) Summary() BundleSummary {
	return BundleSummary{
		CoreSkillCount:      len(b.Skills.Core),
		PlaybookCount:       len(b.Skills.Playbooks),
		AllowedToolCount:    len(b.Tools.AllowedTools),
		CommandPatternCount: len(b.Prompt.CommandPatterns),
		ChecklistCount:      len(b.Prompt.Checklists),
		MemoryPointerCount:  len(b.Memory.Pointers),
		WatchTopicCount:     len(b.Heartbeat.WatchTopics),
		SourceMapCount:      len(b.SourceMap),
	}
}

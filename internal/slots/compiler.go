package slots

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var headingPattern = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*$`)
var numberedPattern = regexp.MustCompile(`^\d+\.\s+(.*)$`)

type Compiler struct {
	requiredFiles []string
}

type markdownNode struct {
	Level    int
	Title    string
	Lines    []string
	Children []*markdownNode
	FileName string
}

func NewCompiler() *Compiler {
	return &Compiler{
		requiredFiles: []string{
			"IDENTITY.md",
			"SOUL.md",
			"AGENTS.md",
			"SKILLS.md",
			"TOOLS.md",
			"PROMPT.md",
			"MEMORY.md",
			"HEARTBEAT.md",
			"DREAMS.md",
			"POSTMORTEM.md",
		},
	}
}

func (c *Compiler) CompileSpecialist(specialistID string, files []File) (*Bundle, error) {
	byName := make(map[string]File, len(files))
	for _, file := range files {
		byName[strings.ToUpper(file.Name)] = file
	}

	missing := make([]string, 0)
	for _, required := range c.requiredFiles {
		if _, ok := byName[strings.ToUpper(required)]; !ok {
			missing = append(missing, required)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return nil, fmt.Errorf("missing required slot files: %s", strings.Join(missing, ", "))
	}

	documents := make(map[string]*markdownNode, len(c.requiredFiles))
	for _, required := range c.requiredFiles {
		file := byName[strings.ToUpper(required)]
		doc, err := parseMarkdownFile(file)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", file.Name, err)
		}
		documents[file.Name] = doc
	}

	bundle := &Bundle{
		Schema:       "sealhat.specialist.bundle",
		Version:      "1.0",
		SpecialistID: specialistID,
		Status:       "active",
		SourceMap:    make(map[string]string),
		SlotPolicy:   defaultSlotPolicy(),
	}

	if err := compileIdentity(documents["IDENTITY.md"], bundle); err != nil {
		return nil, err
	}
	bundle.ParentID = bundle.Identity.EscalateTo
	if err := compileSoul(documents["SOUL.md"], bundle); err != nil {
		return nil, err
	}
	if err := compileAgents(documents["AGENTS.md"], bundle); err != nil {
		return nil, err
	}
	if err := compileTools(documents["TOOLS.md"], bundle); err != nil {
		return nil, err
	}
	if err := compileSkills(documents["SKILLS.md"], bundle); err != nil {
		return nil, err
	}
	if err := compilePrompt(documents["PROMPT.md"], bundle); err != nil {
		return nil, err
	}
	if err := compileMemory(documents["MEMORY.md"], bundle); err != nil {
		return nil, err
	}
	if err := compileHeartbeat(documents["HEARTBEAT.md"], bundle); err != nil {
		return nil, err
	}
	if err := compileDreams(documents["DREAMS.md"], bundle); err != nil {
		return nil, err
	}
	if err := compilePostmortem(documents["POSTMORTEM.md"], bundle); err != nil {
		return nil, err
	}

	return bundle, nil
}

func defaultSlotPolicy() SlotPolicyBundle {
	return SlotPolicyBundle{
		Constitutional: []string{
			"identity.*",
			"soul.*",
			"agents.mission",
			"agents.lane_boundaries",
			"agents.self_modification_rules",
			"agents.mandatory_postmortem_rule",
		},
		Operational: []string{
			"agents.escalation_rules",
			"agents.decision_rules",
			"skills.*",
			"tools.*",
			"prompt.*",
			"memory.*",
			"heartbeat.*",
		},
		ReviewPlane: []string{
			"dreams.*",
			"postmortem_template.*",
		},
		Runtime: []string{
			"specialist_id",
			"parent_id",
			"status",
			"runtime_granted_tools",
			"memory_search_results",
		},
	}
}

func parseMarkdownFile(file File) (*markdownNode, error) {
	root := &markdownNode{Level: 0, Title: "ROOT", FileName: file.Name}
	stack := []*markdownNode{root}
	lines := strings.Split(strings.ReplaceAll(file.Content, "\r\n", "\n"), "\n")
	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")
		if matches := headingPattern.FindStringSubmatch(line); matches != nil {
			level := len(matches[1])
			title := strings.TrimSpace(matches[2])
			node := &markdownNode{Level: level, Title: title, FileName: file.Name}
			for len(stack) > 0 && stack[len(stack)-1].Level >= level {
				stack = stack[:len(stack)-1]
			}
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, node)
			stack = append(stack, node)
			continue
		}
		stack[len(stack)-1].Lines = append(stack[len(stack)-1].Lines, line)
	}

	if len(root.Children) == 0 {
		return nil, fmt.Errorf("missing H1 heading")
	}
	if len(root.Children) != 1 || root.Children[0].Level != 1 {
		return nil, fmt.Errorf("expected exactly one H1 heading")
	}
	doc := root.Children[0]
	expected := strings.TrimSuffix(strings.ToUpper(filepath.Base(file.Name)), ".MD")
	if doc.Title != expected {
		return nil, fmt.Errorf("H1 %q does not match filename %q", doc.Title, expected)
	}
	return doc, nil
}

func compileIdentity(doc *markdownNode, bundle *Bundle) error {
	name, err := child(doc, "name")
	if err != nil {
		return err
	}
	bundle.Identity.Name = scalarContent(name)
	bundle.SourceMap["identity.name"] = sourceRef(doc.FileName, "name")

	aliases, err := child(doc, "aliases")
	if err != nil {
		return err
	}
	bundle.Identity.Aliases = bulletLines(aliases)
	bundle.SourceMap["identity.aliases"] = sourceRef(doc.FileName, "aliases")

	domain, err := child(doc, "domain")
	if err != nil {
		return err
	}
	bundle.Identity.Domain = scalarContent(domain)
	bundle.SourceMap["identity.domain"] = sourceRef(doc.FileName, "domain")

	role, err := child(doc, "role")
	if err != nil {
		return err
	}
	bundle.Identity.Role = scalarContent(role)
	bundle.SourceMap["identity.role"] = sourceRef(doc.FileName, "role")

	authority, err := child(doc, "authority_boundary")
	if err != nil {
		return err
	}
	bundle.Identity.AuthorityBoundary = scalarContent(authority)
	bundle.SourceMap["identity.authority_boundary"] = sourceRef(doc.FileName, "authority_boundary")

	allowed, err := child(doc, "allowed_scope")
	if err != nil {
		return err
	}
	bundle.Identity.AllowedScope = bulletLines(allowed)
	bundle.SourceMap["identity.allowed_scope"] = sourceRef(doc.FileName, "allowed_scope")

	forbidden, err := child(doc, "forbidden_scope")
	if err != nil {
		return err
	}
	bundle.Identity.ForbiddenScope = bulletLines(forbidden)
	bundle.SourceMap["identity.forbidden_scope"] = sourceRef(doc.FileName, "forbidden_scope")

	escalate, err := child(doc, "escalate_to")
	if err != nil {
		return err
	}
	bundle.Identity.EscalateTo = scalarContent(escalate)
	bundle.SourceMap["identity.escalate_to"] = sourceRef(doc.FileName, "escalate_to")
	return nil
}

func compileSoul(doc *markdownNode, bundle *Bundle) error {
	tone, err := child(doc, "tone")
	if err != nil {
		return err
	}
	bundle.Soul.Tone = scalarContent(tone)
	bundle.SourceMap["soul.tone"] = sourceRef(doc.FileName, "tone")

	values, err := child(doc, "values")
	if err != nil {
		return err
	}
	bundle.Soul.Values = bulletLines(values)
	bundle.SourceMap["soul.values"] = sourceRef(doc.FileName, "values")

	caution, err := child(doc, "caution_profile")
	if err != nil {
		return err
	}
	high, moderate := parseCautionProfile(caution)
	bundle.Soul.CautionProfileHigh = high
	bundle.Soul.CautionProfileModerate = moderate
	bundle.SourceMap["soul.caution_profile_high"] = sourceRef(doc.FileName, "caution_profile")
	bundle.SourceMap["soul.caution_profile_moderate"] = sourceRef(doc.FileName, "caution_profile")

	decisionStyle, err := child(doc, "decision_style")
	if err != nil {
		return err
	}
	bundle.Soul.DecisionStyle = bulletLines(decisionStyle)
	bundle.SourceMap["soul.decision_style"] = sourceRef(doc.FileName, "decision_style")

	riskTolerance, err := child(doc, "risk_tolerance")
	if err != nil {
		return err
	}
	bundle.Soul.RiskTolerance = scalarContent(riskTolerance)
	bundle.SourceMap["soul.risk_tolerance"] = sourceRef(doc.FileName, "risk_tolerance")
	return nil
}

func compileAgents(doc *markdownNode, bundle *Bundle) error {
	mission, err := child(doc, "mission")
	if err != nil {
		return err
	}
	bundle.Agents.Mission = scalarContent(mission)
	bundle.SourceMap["agents.mission"] = sourceRef(doc.FileName, "mission")

	laneBoundaries, err := child(doc, "lane_boundaries")
	if err != nil {
		return err
	}
	bundle.Agents.LaneBoundaries = bulletLines(laneBoundaries)
	bundle.SourceMap["agents.lane_boundaries"] = sourceRef(doc.FileName, "lane_boundaries")

	escalationRules, err := child(doc, "escalation_rules")
	if err != nil {
		return err
	}
	bundle.Agents.EscalationRules = bulletLines(escalationRules)
	bundle.SourceMap["agents.escalation_rules"] = sourceRef(doc.FileName, "escalation_rules")

	decisionRules, err := child(doc, "decision_rules")
	if err != nil {
		return err
	}
	bundle.Agents.DecisionRules = bulletLines(decisionRules)
	bundle.SourceMap["agents.decision_rules"] = sourceRef(doc.FileName, "decision_rules")

	selfModificationRules, err := child(doc, "self_modification_rules")
	if err != nil {
		return err
	}
	bundle.Agents.SelfModificationRules = statementLines(selfModificationRules)
	bundle.SourceMap["agents.self_modification_rules"] = sourceRef(doc.FileName, "self_modification_rules")

	mandatoryPostmortemRule, err := child(doc, "mandatory_postmortem_rule")
	if err != nil {
		return err
	}
	bundle.Agents.MandatoryPostmortemRule = scalarContent(mandatoryPostmortemRule)
	bundle.SourceMap["agents.mandatory_postmortem_rule"] = sourceRef(doc.FileName, "mandatory_postmortem_rule")
	return nil
}

func compileTools(doc *markdownNode, bundle *Bundle) error {
	allowedTools, err := child(doc, "allowed_tools")
	if err != nil {
		return err
	}
	bundle.Tools.AllowedTools = bulletLines(allowedTools)
	bundle.SourceMap["tools.allowed_tools"] = sourceRef(doc.FileName, "allowed_tools")

	toolOrdering, err := child(doc, "tool_ordering")
	if err != nil {
		return err
	}
	bundle.Tools.ToolOrdering = numberedLines(toolOrdering)
	bundle.SourceMap["tools.tool_ordering"] = sourceRef(doc.FileName, "tool_ordering")

	toolConstraints, err := child(doc, "tool_constraints")
	if err != nil {
		return err
	}
	bundle.Tools.ToolConstraints = bulletLines(toolConstraints)
	bundle.SourceMap["tools.tool_constraints"] = sourceRef(doc.FileName, "tool_constraints")

	fallbackBehavior, err := child(doc, "fallback_behavior")
	if err != nil {
		return err
	}
	bundle.Tools.FallbackBehavior = bulletLines(fallbackBehavior)
	bundle.SourceMap["tools.fallback_behavior"] = sourceRef(doc.FileName, "fallback_behavior")
	return nil
}

func compileSkills(doc *markdownNode, bundle *Bundle) error {
	core, err := child(doc, "core_skills")
	if err != nil {
		return err
	}
	for i, childNode := range core.Children {
		trigger, err := child(childNode, "trigger")
		if err != nil {
			return err
		}
		procedure, err := child(childNode, "procedure")
		if err != nil {
			return err
		}
		successCheck, err := child(childNode, "success_check")
		if err != nil {
			return err
		}
		record := SkillRecord{
			Name:         scalarContent(childNode),
			Trigger:      scalarContent(trigger),
			Procedure:    numberedLines(procedure),
			SuccessCheck: bulletLines(successCheck),
		}
		bundle.Skills.Core = append(bundle.Skills.Core, record)
		bundle.SourceMap[fmt.Sprintf("skills.core[%d].name", i)] = sourceRef(doc.FileName, "core_skills", childNode.Title)
		bundle.SourceMap[fmt.Sprintf("skills.core[%d].trigger", i)] = sourceRef(doc.FileName, "core_skills", childNode.Title, "trigger")
		bundle.SourceMap[fmt.Sprintf("skills.core[%d].procedure", i)] = sourceRef(doc.FileName, "core_skills", childNode.Title, "procedure")
		bundle.SourceMap[fmt.Sprintf("skills.core[%d].success_check", i)] = sourceRef(doc.FileName, "core_skills", childNode.Title, "success_check")
	}

	playbooks, err := child(doc, "playbooks")
	if err != nil {
		return err
	}
	for i, childNode := range playbooks.Children {
		trigger, err := child(childNode, "trigger")
		if err != nil {
			return err
		}
		procedure, err := child(childNode, "procedure")
		if err != nil {
			return err
		}
		record := PlaybookRecord{
			Name:      scalarContent(childNode),
			Trigger:   scalarContent(trigger),
			Procedure: numberedLines(procedure),
		}
		bundle.Skills.Playbooks = append(bundle.Skills.Playbooks, record)
		bundle.SourceMap[fmt.Sprintf("skills.playbooks[%d].name", i)] = sourceRef(doc.FileName, "playbooks", childNode.Title)
		bundle.SourceMap[fmt.Sprintf("skills.playbooks[%d].trigger", i)] = sourceRef(doc.FileName, "playbooks", childNode.Title, "trigger")
		bundle.SourceMap[fmt.Sprintf("skills.playbooks[%d].procedure", i)] = sourceRef(doc.FileName, "playbooks", childNode.Title, "procedure")
	}

	failureRecoveryPatterns, err := child(doc, "failure_recovery_patterns")
	if err != nil {
		return err
	}
	bundle.Skills.FailureRecoveryPatterns = bulletLines(failureRecoveryPatterns)
	bundle.SourceMap["skills.failure_recovery_patterns"] = sourceRef(doc.FileName, "failure_recovery_patterns")
	return nil
}

func compilePrompt(doc *markdownNode, bundle *Bundle) error {
	templates, err := child(doc, "task_templates")
	if err != nil {
		return err
	}
	for i, childNode := range templates.Children {
		purpose, err := child(childNode, "purpose")
		if err != nil {
			return err
		}
		body, err := child(childNode, "body")
		if err != nil {
			return err
		}
		record := PromptTemplate{
			Name:    scalarContent(childNode),
			Purpose: scalarContent(purpose),
			Body:    contentLines(body),
		}
		bundle.Prompt.TaskTemplates = append(bundle.Prompt.TaskTemplates, record)
		bundle.SourceMap[fmt.Sprintf("prompt.task_templates[%d].name", i)] = sourceRef(doc.FileName, "task_templates", childNode.Title)
		bundle.SourceMap[fmt.Sprintf("prompt.task_templates[%d].purpose", i)] = sourceRef(doc.FileName, "task_templates", childNode.Title, "purpose")
		bundle.SourceMap[fmt.Sprintf("prompt.task_templates[%d].body", i)] = sourceRef(doc.FileName, "task_templates", childNode.Title, "body")
	}

	patterns, err := child(doc, "command_patterns")
	if err != nil {
		return err
	}
	for i, childNode := range patterns.Children {
		purpose, err := child(childNode, "purpose")
		if err != nil {
			return err
		}
		body, err := child(childNode, "body")
		if err != nil {
			return err
		}
		record := PromptTemplate{
			Name:    scalarContent(childNode),
			Purpose: scalarContent(purpose),
			Body:    contentLines(body),
		}
		bundle.Prompt.CommandPatterns = append(bundle.Prompt.CommandPatterns, record)
		bundle.SourceMap[fmt.Sprintf("prompt.command_patterns[%d].name", i)] = sourceRef(doc.FileName, "command_patterns", childNode.Title)
		bundle.SourceMap[fmt.Sprintf("prompt.command_patterns[%d].purpose", i)] = sourceRef(doc.FileName, "command_patterns", childNode.Title, "purpose")
		bundle.SourceMap[fmt.Sprintf("prompt.command_patterns[%d].body", i)] = sourceRef(doc.FileName, "command_patterns", childNode.Title, "body")
	}

	checklists, err := child(doc, "checklists")
	if err != nil {
		return err
	}
	for i, childNode := range checklists.Children {
		name, items := parseChecklist(childNode)
		bundle.Prompt.Checklists = append(bundle.Prompt.Checklists, Checklist{Name: name, Items: items})
		bundle.SourceMap[fmt.Sprintf("prompt.checklists[%d].name", i)] = sourceRef(doc.FileName, "checklists", childNode.Title)
		bundle.SourceMap[fmt.Sprintf("prompt.checklists[%d].items", i)] = sourceRef(doc.FileName, "checklists", childNode.Title)
	}
	return nil
}

func compileMemory(doc *markdownNode, bundle *Bundle) error {
	currentFocus, err := child(doc, "current_focus")
	if err != nil {
		return err
	}
	bundle.Memory.CurrentFocus = bulletLines(currentFocus)
	bundle.SourceMap["memory.current_focus"] = sourceRef(doc.FileName, "current_focus")

	durableDecisions, err := child(doc, "durable_decisions")
	if err != nil {
		return err
	}
	bundle.Memory.DurableDecisions = bulletLines(durableDecisions)
	bundle.SourceMap["memory.durable_decisions"] = sourceRef(doc.FileName, "durable_decisions")

	memoryPointers, err := child(doc, "memory_pointers")
	if err != nil {
		return err
	}
	pointers, err := parseMemoryPointers(memoryPointers)
	if err != nil {
		return fmt.Errorf("memory pointers: %w", err)
	}
	bundle.Memory.Pointers = pointers
	bundle.SourceMap["memory.pointers"] = sourceRef(doc.FileName, "memory_pointers")

	activeConstraints, err := child(doc, "active_constraints")
	if err != nil {
		return err
	}
	bundle.Memory.ActiveConstraints = bulletLines(activeConstraints)
	bundle.SourceMap["memory.active_constraints"] = sourceRef(doc.FileName, "active_constraints")
	return nil
}

func compileHeartbeat(doc *markdownNode, bundle *Bundle) error {
	watchTopics, err := child(doc, "watch_topics")
	if err != nil {
		return err
	}
	bundle.Heartbeat.WatchTopics = bulletLines(watchTopics)
	bundle.SourceMap["heartbeat.watch_topics"] = sourceRef(doc.FileName, "watch_topics")

	approvedSources, err := child(doc, "approved_sources")
	if err != nil {
		return err
	}
	bundle.Heartbeat.ApprovedSources = bulletLines(approvedSources)
	bundle.SourceMap["heartbeat.approved_sources"] = sourceRef(doc.FileName, "approved_sources")

	cadence, err := child(doc, "cadence")
	if err != nil {
		return err
	}
	bundle.Heartbeat.Cadence = bulletLines(cadence)
	bundle.SourceMap["heartbeat.cadence"] = sourceRef(doc.FileName, "cadence")

	promotionRules, err := child(doc, "promotion_rules")
	if err != nil {
		return err
	}
	bundle.Heartbeat.PromotionRules = bulletLines(promotionRules)
	bundle.SourceMap["heartbeat.promotion_rules"] = sourceRef(doc.FileName, "promotion_rules")

	stopConditions, err := child(doc, "stop_conditions")
	if err != nil {
		return err
	}
	bundle.Heartbeat.StopConditions = bulletLines(stopConditions)
	bundle.SourceMap["heartbeat.stop_conditions"] = sourceRef(doc.FileName, "stop_conditions")
	return nil
}

func compileDreams(doc *markdownNode, bundle *Bundle) error {
	distilledPatterns, err := child(doc, "distilled_patterns")
	if err != nil {
		return err
	}
	bundle.Dreams.DistilledPatterns = bulletLines(distilledPatterns)
	bundle.SourceMap["dreams.distilled_patterns"] = sourceRef(doc.FileName, "distilled_patterns")

	candidatePromotions, err := child(doc, "candidate_promotions")
	if err != nil {
		return err
	}
	bundle.Dreams.CandidatePromotions = bulletLines(candidatePromotions)
	bundle.SourceMap["dreams.candidate_promotions"] = sourceRef(doc.FileName, "candidate_promotions")

	discardedNoise, err := child(doc, "discarded_noise")
	if err != nil {
		return err
	}
	bundle.Dreams.DiscardedNoise = bulletLines(discardedNoise)
	bundle.SourceMap["dreams.discarded_noise"] = sourceRef(doc.FileName, "discarded_noise")

	openReviewItems, err := child(doc, "open_review_items")
	if err != nil {
		return err
	}
	bundle.Dreams.OpenReviewItems = bulletLines(openReviewItems)
	bundle.SourceMap["dreams.open_review_items"] = sourceRef(doc.FileName, "open_review_items")
	return nil
}

func compilePostmortem(doc *markdownNode, bundle *Bundle) error {
	incidentID, err := child(doc, "incident_id")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.IncidentID = scalarContent(incidentID)
	bundle.SourceMap["postmortem_template.incident_id"] = sourceRef(doc.FileName, "incident_id")

	date, err := child(doc, "date")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.Date = scalarContent(date)
	bundle.SourceMap["postmortem_template.date"] = sourceRef(doc.FileName, "date")

	modelID, err := child(doc, "model_id")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.ModelID = scalarContent(modelID)
	bundle.SourceMap["postmortem_template.model_id"] = sourceRef(doc.FileName, "model_id")

	modelRole, err := child(doc, "model_role")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.ModelRole = scalarContent(modelRole)
	bundle.SourceMap["postmortem_template.model_role"] = sourceRef(doc.FileName, "model_role")

	taskSummary, err := child(doc, "task_summary")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.TaskSummary = scalarContent(taskSummary)
	bundle.SourceMap["postmortem_template.task_summary"] = sourceRef(doc.FileName, "task_summary")

	expectedBehavior, err := child(doc, "expected_behavior")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.ExpectedBehavior = scalarContent(expectedBehavior)
	bundle.SourceMap["postmortem_template.expected_behavior"] = sourceRef(doc.FileName, "expected_behavior")

	actualBehavior, err := child(doc, "actual_behavior")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.ActualBehavior = scalarContent(actualBehavior)
	bundle.SourceMap["postmortem_template.actual_behavior"] = sourceRef(doc.FileName, "actual_behavior")

	whatWentWrong, err := child(doc, "what_went_wrong")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.WhatWentWrong = scalarContent(whatWentWrong)
	bundle.SourceMap["postmortem_template.what_went_wrong"] = sourceRef(doc.FileName, "what_went_wrong")

	failureClassification, err := child(doc, "failure_classification")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.FailureClassification = bulletLines(failureClassification)
	bundle.SourceMap["postmortem_template.failure_classification"] = sourceRef(doc.FileName, "failure_classification")

	rootCause, err := child(doc, "root_cause")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.RootCause = scalarContent(rootCause)
	bundle.SourceMap["postmortem_template.root_cause"] = sourceRef(doc.FileName, "root_cause")

	missedEvidenceOrStep, err := child(doc, "missed_evidence_or_step")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.MissedEvidenceOrStep = scalarContent(missedEvidenceOrStep)
	bundle.SourceMap["postmortem_template.missed_evidence_or_step"] = sourceRef(doc.FileName, "missed_evidence_or_step")

	preventable, err := child(doc, "preventable")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.Preventable = scalarContent(preventable)
	bundle.SourceMap["postmortem_template.preventable"] = sourceRef(doc.FileName, "preventable")

	correction, err := child(doc, "correction")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.Correction = scalarContent(correction)
	bundle.SourceMap["postmortem_template.correction"] = sourceRef(doc.FileName, "correction")

	recommendedSlotOrSystemChange, err := child(doc, "recommended_slot_or_system_change")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.RecommendedSlotOrSystemChange = bulletLines(recommendedSlotOrSystemChange)
	bundle.SourceMap["postmortem_template.recommended_slot_or_system_change"] = sourceRef(doc.FileName, "recommended_slot_or_system_change")

	requiresHarnessReview, err := child(doc, "requires_harness_review")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.RequiresHarnessReview = scalarContent(requiresHarnessReview)
	bundle.SourceMap["postmortem_template.requires_harness_review"] = sourceRef(doc.FileName, "requires_harness_review")

	requiresParentReview, err := child(doc, "requires_parent_review")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.RequiresParentReview = scalarContent(requiresParentReview)
	bundle.SourceMap["postmortem_template.requires_parent_review"] = sourceRef(doc.FileName, "requires_parent_review")

	regressionTestNeeded, err := child(doc, "regression_test_needed")
	if err != nil {
		return err
	}
	bundle.PostmortemTemplate.RegressionTestNeeded = scalarContent(regressionTestNeeded)
	bundle.SourceMap["postmortem_template.regression_test_needed"] = sourceRef(doc.FileName, "regression_test_needed")
	return nil
}

func child(parent *markdownNode, title string) (*markdownNode, error) {
	for _, candidate := range parent.Children {
		if strings.EqualFold(candidate.Title, title) {
			return candidate, nil
		}
	}
	return nil, fmt.Errorf("missing section %q in %s", title, parent.FileName)
}

func sourceRef(fileName string, headings ...string) string {
	clean := make([]string, 0, len(headings))
	for _, heading := range headings {
		if strings.TrimSpace(heading) == "" {
			continue
		}
		clean = append(clean, heading)
	}
	if len(clean) == 0 {
		return fileName
	}
	return fmt.Sprintf("%s##%s", fileName, strings.Join(clean, "/"))
}

func scalarContent(node *markdownNode) string {
	parts := make([]string, 0)
	for _, raw := range node.Lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if _, ok := stripBullet(line); ok {
			continue
		}
		if _, ok := stripNumbered(line); ok {
			continue
		}
		parts = append(parts, line)
	}
	return strings.Join(parts, "\n")
}

func statementLines(node *markdownNode) []string {
	items := make([]string, 0)
	for _, raw := range node.Lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		items = append(items, line)
	}
	return items
}

func bulletLines(node *markdownNode) []string {
	items := make([]string, 0)
	for _, raw := range node.Lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if item, ok := stripBullet(line); ok {
			items = append(items, item)
		}
	}
	return items
}

func numberedLines(node *markdownNode) []string {
	items := make([]string, 0)
	for _, raw := range node.Lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if item, ok := stripNumbered(line); ok {
			items = append(items, item)
		}
	}
	return items
}

func contentLines(node *markdownNode) []string {
	items := make([]string, 0)
	for _, raw := range node.Lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if item, ok := stripBullet(line); ok {
			items = append(items, item)
			continue
		}
		if item, ok := stripNumbered(line); ok {
			items = append(items, item)
			continue
		}
		items = append(items, line)
	}
	return items
}

func parseChecklist(node *markdownNode) (string, []string) {
	name := ""
	items := make([]string, 0)
	for _, raw := range node.Lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if item, ok := stripBullet(line); ok {
			items = append(items, item)
			continue
		}
		if name == "" {
			name = line
		}
	}
	return name, items
}

func parseCautionProfile(node *markdownNode) ([]string, []string) {
	high := make([]string, 0)
	moderate := make([]string, 0)
	mode := ""
	for _, raw := range node.Lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		switch line {
		case "High caution for:":
			mode = "high"
			continue
		case "Moderate caution for:":
			mode = "moderate"
			continue
		}
		item, ok := stripBullet(line)
		if !ok {
			continue
		}
		switch mode {
		case "high":
			high = append(high, item)
		case "moderate":
			moderate = append(moderate, item)
		}
	}
	return high, moderate
}

func parseMemoryPointers(node *markdownNode) ([]MemoryPointer, error) {
	var pointers []MemoryPointer
	var current *MemoryPointer
	for _, raw := range node.Lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if item, ok := stripBullet(line); ok {
			key, value, ok := splitKeyValue(item)
			if !ok {
				return nil, fmt.Errorf("invalid pointer line %q", line)
			}
			if strings.EqualFold(key, "pointer_name") {
				pointers = append(pointers, MemoryPointer{PointerName: value})
				current = &pointers[len(pointers)-1]
				continue
			}
		}
		if current == nil {
			return nil, fmt.Errorf("memory pointer field without pointer_name: %q", line)
		}
		key, value, ok := splitKeyValue(line)
		if !ok {
			return nil, fmt.Errorf("invalid pointer field %q", line)
		}
		switch strings.ToLower(key) {
		case "backend":
			current.Backend = value
		case "query_hint":
			current.QueryHint = value
		case "pointer_name":
			current.PointerName = value
		default:
			return nil, fmt.Errorf("unsupported pointer field %q", key)
		}
	}
	return pointers, nil
}

func splitKeyValue(line string) (string, string, bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if key == "" || value == "" {
		return "", "", false
	}
	return key, value, true
}

func stripBullet(line string) (string, bool) {
	if strings.HasPrefix(line, "- ") {
		return strings.TrimSpace(line[2:]), true
	}
	if strings.HasPrefix(line, "* ") {
		return strings.TrimSpace(line[2:]), true
	}
	return "", false
}

func stripNumbered(line string) (string, bool) {
	matches := numberedPattern.FindStringSubmatch(line)
	if matches == nil {
		return "", false
	}
	return strings.TrimSpace(matches[1]), true
}

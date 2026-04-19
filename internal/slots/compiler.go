package slots

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
	bundle.Identity.Name = scalarContent(requireChild(doc, "name"))
	bundle.SourceMap["identity.name"] = sourceRef(doc.FileName, "name")
	bundle.Identity.Aliases = bulletLines(requireChild(doc, "aliases"))
	bundle.SourceMap["identity.aliases"] = sourceRef(doc.FileName, "aliases")
	bundle.Identity.Domain = scalarContent(requireChild(doc, "domain"))
	bundle.SourceMap["identity.domain"] = sourceRef(doc.FileName, "domain")
	bundle.Identity.Role = scalarContent(requireChild(doc, "role"))
	bundle.SourceMap["identity.role"] = sourceRef(doc.FileName, "role")
	bundle.Identity.AuthorityBoundary = scalarContent(requireChild(doc, "authority_boundary"))
	bundle.SourceMap["identity.authority_boundary"] = sourceRef(doc.FileName, "authority_boundary")
	bundle.Identity.AllowedScope = bulletLines(requireChild(doc, "allowed_scope"))
	bundle.SourceMap["identity.allowed_scope"] = sourceRef(doc.FileName, "allowed_scope")
	bundle.Identity.ForbiddenScope = bulletLines(requireChild(doc, "forbidden_scope"))
	bundle.SourceMap["identity.forbidden_scope"] = sourceRef(doc.FileName, "forbidden_scope")
	bundle.Identity.EscalateTo = scalarContent(requireChild(doc, "escalate_to"))
	bundle.SourceMap["identity.escalate_to"] = sourceRef(doc.FileName, "escalate_to")
	return nil
}

func compileSoul(doc *markdownNode, bundle *Bundle) error {
	bundle.Soul.Tone = scalarContent(requireChild(doc, "tone"))
	bundle.SourceMap["soul.tone"] = sourceRef(doc.FileName, "tone")
	bundle.Soul.Values = bulletLines(requireChild(doc, "values"))
	bundle.SourceMap["soul.values"] = sourceRef(doc.FileName, "values")
	caution := requireChild(doc, "caution_profile")
	high, moderate := parseCautionProfile(caution)
	bundle.Soul.CautionProfileHigh = high
	bundle.Soul.CautionProfileModerate = moderate
	bundle.SourceMap["soul.caution_profile_high"] = sourceRef(doc.FileName, "caution_profile")
	bundle.SourceMap["soul.caution_profile_moderate"] = sourceRef(doc.FileName, "caution_profile")
	bundle.Soul.DecisionStyle = bulletLines(requireChild(doc, "decision_style"))
	bundle.SourceMap["soul.decision_style"] = sourceRef(doc.FileName, "decision_style")
	bundle.Soul.RiskTolerance = scalarContent(requireChild(doc, "risk_tolerance"))
	bundle.SourceMap["soul.risk_tolerance"] = sourceRef(doc.FileName, "risk_tolerance")
	return nil
}

func compileAgents(doc *markdownNode, bundle *Bundle) error {
	bundle.Agents.Mission = scalarContent(requireChild(doc, "mission"))
	bundle.SourceMap["agents.mission"] = sourceRef(doc.FileName, "mission")
	bundle.Agents.LaneBoundaries = bulletLines(requireChild(doc, "lane_boundaries"))
	bundle.SourceMap["agents.lane_boundaries"] = sourceRef(doc.FileName, "lane_boundaries")
	bundle.Agents.EscalationRules = bulletLines(requireChild(doc, "escalation_rules"))
	bundle.SourceMap["agents.escalation_rules"] = sourceRef(doc.FileName, "escalation_rules")
	bundle.Agents.DecisionRules = bulletLines(requireChild(doc, "decision_rules"))
	bundle.SourceMap["agents.decision_rules"] = sourceRef(doc.FileName, "decision_rules")
	bundle.Agents.SelfModificationRules = statementLines(requireChild(doc, "self_modification_rules"))
	bundle.SourceMap["agents.self_modification_rules"] = sourceRef(doc.FileName, "self_modification_rules")
	bundle.Agents.MandatoryPostmortemRule = scalarContent(requireChild(doc, "mandatory_postmortem_rule"))
	bundle.SourceMap["agents.mandatory_postmortem_rule"] = sourceRef(doc.FileName, "mandatory_postmortem_rule")
	return nil
}

func compileTools(doc *markdownNode, bundle *Bundle) error {
	bundle.Tools.AllowedTools = bulletLines(requireChild(doc, "allowed_tools"))
	bundle.SourceMap["tools.allowed_tools"] = sourceRef(doc.FileName, "allowed_tools")
	bundle.Tools.ToolOrdering = numberedLines(requireChild(doc, "tool_ordering"))
	bundle.SourceMap["tools.tool_ordering"] = sourceRef(doc.FileName, "tool_ordering")
	bundle.Tools.ToolConstraints = bulletLines(requireChild(doc, "tool_constraints"))
	bundle.SourceMap["tools.tool_constraints"] = sourceRef(doc.FileName, "tool_constraints")
	bundle.Tools.FallbackBehavior = bulletLines(requireChild(doc, "fallback_behavior"))
	bundle.SourceMap["tools.fallback_behavior"] = sourceRef(doc.FileName, "fallback_behavior")
	return nil
}

func compileSkills(doc *markdownNode, bundle *Bundle) error {
	core := requireChild(doc, "core_skills")
	for i, child := range core.Children {
		name := scalarContent(child)
		record := SkillRecord{
			Name:         name,
			Trigger:      scalarContent(requireChild(child, "trigger")),
			Procedure:    numberedLines(requireChild(child, "procedure")),
			SuccessCheck: bulletLines(requireChild(child, "success_check")),
		}
		bundle.Skills.Core = append(bundle.Skills.Core, record)
		bundle.SourceMap[fmt.Sprintf("skills.core[%d].name", i)] = sourceRef(doc.FileName, "core_skills", child.Title)
		bundle.SourceMap[fmt.Sprintf("skills.core[%d].trigger", i)] = sourceRef(doc.FileName, "core_skills", child.Title, "trigger")
		bundle.SourceMap[fmt.Sprintf("skills.core[%d].procedure", i)] = sourceRef(doc.FileName, "core_skills", child.Title, "procedure")
		bundle.SourceMap[fmt.Sprintf("skills.core[%d].success_check", i)] = sourceRef(doc.FileName, "core_skills", child.Title, "success_check")
	}

	playbooks := requireChild(doc, "playbooks")
	for i, child := range playbooks.Children {
		record := PlaybookRecord{
			Name:      scalarContent(child),
			Trigger:   scalarContent(requireChild(child, "trigger")),
			Procedure: numberedLines(requireChild(child, "procedure")),
		}
		bundle.Skills.Playbooks = append(bundle.Skills.Playbooks, record)
		bundle.SourceMap[fmt.Sprintf("skills.playbooks[%d].name", i)] = sourceRef(doc.FileName, "playbooks", child.Title)
		bundle.SourceMap[fmt.Sprintf("skills.playbooks[%d].trigger", i)] = sourceRef(doc.FileName, "playbooks", child.Title, "trigger")
		bundle.SourceMap[fmt.Sprintf("skills.playbooks[%d].procedure", i)] = sourceRef(doc.FileName, "playbooks", child.Title, "procedure")
	}

	bundle.Skills.FailureRecoveryPatterns = bulletLines(requireChild(doc, "failure_recovery_patterns"))
	bundle.SourceMap["skills.failure_recovery_patterns"] = sourceRef(doc.FileName, "failure_recovery_patterns")
	return nil
}

func compilePrompt(doc *markdownNode, bundle *Bundle) error {
	templates := requireChild(doc, "task_templates")
	for i, child := range templates.Children {
		record := PromptTemplate{
			Name:    scalarContent(child),
			Purpose: scalarContent(requireChild(child, "purpose")),
			Body:    contentLines(requireChild(child, "body")),
		}
		bundle.Prompt.TaskTemplates = append(bundle.Prompt.TaskTemplates, record)
		bundle.SourceMap[fmt.Sprintf("prompt.task_templates[%d].name", i)] = sourceRef(doc.FileName, "task_templates", child.Title)
		bundle.SourceMap[fmt.Sprintf("prompt.task_templates[%d].purpose", i)] = sourceRef(doc.FileName, "task_templates", child.Title, "purpose")
		bundle.SourceMap[fmt.Sprintf("prompt.task_templates[%d].body", i)] = sourceRef(doc.FileName, "task_templates", child.Title, "body")
	}

	patterns := requireChild(doc, "command_patterns")
	for i, child := range patterns.Children {
		record := PromptTemplate{
			Name:    scalarContent(child),
			Purpose: scalarContent(requireChild(child, "purpose")),
			Body:    contentLines(requireChild(child, "body")),
		}
		bundle.Prompt.CommandPatterns = append(bundle.Prompt.CommandPatterns, record)
		bundle.SourceMap[fmt.Sprintf("prompt.command_patterns[%d].name", i)] = sourceRef(doc.FileName, "command_patterns", child.Title)
		bundle.SourceMap[fmt.Sprintf("prompt.command_patterns[%d].purpose", i)] = sourceRef(doc.FileName, "command_patterns", child.Title, "purpose")
		bundle.SourceMap[fmt.Sprintf("prompt.command_patterns[%d].body", i)] = sourceRef(doc.FileName, "command_patterns", child.Title, "body")
	}

	checklists := requireChild(doc, "checklists")
	for i, child := range checklists.Children {
		name, items := parseChecklist(child)
		bundle.Prompt.Checklists = append(bundle.Prompt.Checklists, Checklist{Name: name, Items: items})
		bundle.SourceMap[fmt.Sprintf("prompt.checklists[%d].name", i)] = sourceRef(doc.FileName, "checklists", child.Title)
		bundle.SourceMap[fmt.Sprintf("prompt.checklists[%d].items", i)] = sourceRef(doc.FileName, "checklists", child.Title)
	}
	return nil
}

func compileMemory(doc *markdownNode, bundle *Bundle) error {
	bundle.Memory.CurrentFocus = bulletLines(requireChild(doc, "current_focus"))
	bundle.SourceMap["memory.current_focus"] = sourceRef(doc.FileName, "current_focus")
	bundle.Memory.DurableDecisions = bulletLines(requireChild(doc, "durable_decisions"))
	bundle.SourceMap["memory.durable_decisions"] = sourceRef(doc.FileName, "durable_decisions")
	pointers, err := parseMemoryPointers(requireChild(doc, "memory_pointers"))
	if err != nil {
		return fmt.Errorf("memory pointers: %w", err)
	}
	bundle.Memory.Pointers = pointers
	bundle.SourceMap["memory.pointers"] = sourceRef(doc.FileName, "memory_pointers")
	bundle.Memory.ActiveConstraints = bulletLines(requireChild(doc, "active_constraints"))
	bundle.SourceMap["memory.active_constraints"] = sourceRef(doc.FileName, "active_constraints")
	return nil
}

func compileHeartbeat(doc *markdownNode, bundle *Bundle) error {
	bundle.Heartbeat.WatchTopics = bulletLines(requireChild(doc, "watch_topics"))
	bundle.SourceMap["heartbeat.watch_topics"] = sourceRef(doc.FileName, "watch_topics")
	bundle.Heartbeat.ApprovedSources = bulletLines(requireChild(doc, "approved_sources"))
	bundle.SourceMap["heartbeat.approved_sources"] = sourceRef(doc.FileName, "approved_sources")
	bundle.Heartbeat.Cadence = bulletLines(requireChild(doc, "cadence"))
	bundle.SourceMap["heartbeat.cadence"] = sourceRef(doc.FileName, "cadence")
	bundle.Heartbeat.PromotionRules = bulletLines(requireChild(doc, "promotion_rules"))
	bundle.SourceMap["heartbeat.promotion_rules"] = sourceRef(doc.FileName, "promotion_rules")
	bundle.Heartbeat.StopConditions = bulletLines(requireChild(doc, "stop_conditions"))
	bundle.SourceMap["heartbeat.stop_conditions"] = sourceRef(doc.FileName, "stop_conditions")
	return nil
}

func compileDreams(doc *markdownNode, bundle *Bundle) error {
	bundle.Dreams.DistilledPatterns = bulletLines(requireChild(doc, "distilled_patterns"))
	bundle.SourceMap["dreams.distilled_patterns"] = sourceRef(doc.FileName, "distilled_patterns")
	bundle.Dreams.CandidatePromotions = bulletLines(requireChild(doc, "candidate_promotions"))
	bundle.SourceMap["dreams.candidate_promotions"] = sourceRef(doc.FileName, "candidate_promotions")
	bundle.Dreams.DiscardedNoise = bulletLines(requireChild(doc, "discarded_noise"))
	bundle.SourceMap["dreams.discarded_noise"] = sourceRef(doc.FileName, "discarded_noise")
	bundle.Dreams.OpenReviewItems = bulletLines(requireChild(doc, "open_review_items"))
	bundle.SourceMap["dreams.open_review_items"] = sourceRef(doc.FileName, "open_review_items")
	return nil
}

func compilePostmortem(doc *markdownNode, bundle *Bundle) error {
	bundle.PostmortemTemplate.IncidentID = scalarContent(requireChild(doc, "incident_id"))
	bundle.SourceMap["postmortem_template.incident_id"] = sourceRef(doc.FileName, "incident_id")
	bundle.PostmortemTemplate.Date = scalarContent(requireChild(doc, "date"))
	bundle.SourceMap["postmortem_template.date"] = sourceRef(doc.FileName, "date")
	bundle.PostmortemTemplate.ModelID = scalarContent(requireChild(doc, "model_id"))
	bundle.SourceMap["postmortem_template.model_id"] = sourceRef(doc.FileName, "model_id")
	bundle.PostmortemTemplate.ModelRole = scalarContent(requireChild(doc, "model_role"))
	bundle.SourceMap["postmortem_template.model_role"] = sourceRef(doc.FileName, "model_role")
	bundle.PostmortemTemplate.TaskSummary = scalarContent(requireChild(doc, "task_summary"))
	bundle.SourceMap["postmortem_template.task_summary"] = sourceRef(doc.FileName, "task_summary")
	bundle.PostmortemTemplate.ExpectedBehavior = scalarContent(requireChild(doc, "expected_behavior"))
	bundle.SourceMap["postmortem_template.expected_behavior"] = sourceRef(doc.FileName, "expected_behavior")
	bundle.PostmortemTemplate.ActualBehavior = scalarContent(requireChild(doc, "actual_behavior"))
	bundle.SourceMap["postmortem_template.actual_behavior"] = sourceRef(doc.FileName, "actual_behavior")
	bundle.PostmortemTemplate.WhatWentWrong = scalarContent(requireChild(doc, "what_went_wrong"))
	bundle.SourceMap["postmortem_template.what_went_wrong"] = sourceRef(doc.FileName, "what_went_wrong")
	bundle.PostmortemTemplate.FailureClassification = bulletLines(requireChild(doc, "failure_classification"))
	bundle.SourceMap["postmortem_template.failure_classification"] = sourceRef(doc.FileName, "failure_classification")
	bundle.PostmortemTemplate.RootCause = scalarContent(requireChild(doc, "root_cause"))
	bundle.SourceMap["postmortem_template.root_cause"] = sourceRef(doc.FileName, "root_cause")
	bundle.PostmortemTemplate.MissedEvidenceOrStep = scalarContent(requireChild(doc, "missed_evidence_or_step"))
	bundle.SourceMap["postmortem_template.missed_evidence_or_step"] = sourceRef(doc.FileName, "missed_evidence_or_step")
	bundle.PostmortemTemplate.Preventable = scalarContent(requireChild(doc, "preventable"))
	bundle.SourceMap["postmortem_template.preventable"] = sourceRef(doc.FileName, "preventable")
	bundle.PostmortemTemplate.Correction = scalarContent(requireChild(doc, "correction"))
	bundle.SourceMap["postmortem_template.correction"] = sourceRef(doc.FileName, "correction")
	bundle.PostmortemTemplate.RecommendedSlotOrSystemChange = bulletLines(requireChild(doc, "recommended_slot_or_system_change"))
	bundle.SourceMap["postmortem_template.recommended_slot_or_system_change"] = sourceRef(doc.FileName, "recommended_slot_or_system_change")
	bundle.PostmortemTemplate.RequiresHarnessReview = scalarContent(requireChild(doc, "requires_harness_review"))
	bundle.SourceMap["postmortem_template.requires_harness_review"] = sourceRef(doc.FileName, "requires_harness_review")
	bundle.PostmortemTemplate.RequiresParentReview = scalarContent(requireChild(doc, "requires_parent_review"))
	bundle.SourceMap["postmortem_template.requires_parent_review"] = sourceRef(doc.FileName, "requires_parent_review")
	bundle.PostmortemTemplate.RegressionTestNeeded = scalarContent(requireChild(doc, "regression_test_needed"))
	bundle.SourceMap["postmortem_template.regression_test_needed"] = sourceRef(doc.FileName, "regression_test_needed")
	return nil
}

func requireChild(parent *markdownNode, title string) *markdownNode {
	for _, child := range parent.Children {
		if strings.EqualFold(child.Title, title) {
			return child
		}
	}
	panic(fmt.Sprintf("missing section %q in %s", title, parent.FileName))
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

func MustAtoi(value string) int {
	n, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return n
}

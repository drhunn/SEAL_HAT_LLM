package slots

import (
	"strings"
	"testing"
)

func TestCompilerCompileSpecialist(t *testing.T) {
	files := []File{
		{Name: "IDENTITY.md", Path: "specialists/spec-01/IDENTITY.md", Content: "# IDENTITY\n\n## name\nBuilder Specialist\n\n## aliases\n- Builder\n- Maker\n\n## domain\nTooling\n\n## role\nYou build harness tools.\n\n## authority_boundary\nStay inside tooling.\nDo not rewrite constitutional slots.\n\n## allowed_scope\n- parsers\n- compilers\n- tests\n\n## forbidden_scope\n- identity rewrites\n- policy overrides\n- constitutional self-modification\n- silent authority expansion\n\n## escalate_to\nParent-01\n"},
		{Name: "SOUL.md", Path: "specialists/spec-01/SOUL.md", Content: "# SOUL\n\n## tone\nCalm and exact\n\n## values\n- accuracy\n- reversibility\n- evidence\n\n## caution_profile\nHigh caution for:\n- schema drift\n- unsupported claims\n\nModerate caution for:\n- refactors\n- automation\n\n## decision_style\n- narrow changes\n- validate first\n\n## risk_tolerance\nlow\n"},
		{Name: "AGENTS.md", Path: "specialists/spec-01/AGENTS.md", Content: "# AGENTS\n\n## mission\nCompile specialist slots into a governed bundle.\n\n## lane_boundaries\nStay within:\n- slot loading\n- slot compiling\n- slot verification\n\nDo not drift outside this lane unless explicitly delegated by the parent.\n\n## escalation_rules\nEscalate to the parent when:\n- the task spans multiple specialties\n- constitutional slots are implicated\n\n## decision_rules\n- Read IDENTITY and SOUL before acting\n- Use tools before guessing when verification is possible\n\n## self_modification_rules\nYou may improve how you perform your specialty.\nYou may not redefine your identity.\nAll such changes must be proposed for approval.\n\n## mandatory_postmortem_rule\nGenerate a postmortem when you materially underperform.\n"},
		{Name: "SKILLS.md", Path: "specialists/spec-01/SKILLS.md", Content: "# SKILLS\n\n## core_skills\n\n### skill_name\nCompile slots\n\n#### trigger\nWhen specialist markdown changes\n\n#### procedure\n1. load files\n2. compile bundle\n3. encode toml\n\n#### success_check\n- bundle created\n- toml emitted\n\n## playbooks\n\n### playbook_name\nRepair schema mismatch\n\n#### trigger\nWhen compile fails\n\n#### procedure\n1. inspect headings\n2. fix source\n3. rerun verify\n\n## failure_recovery_patterns\n- Missing evidence\n- Scope drift\n- Validation miss\n"},
		{Name: "TOOLS.md", Path: "specialists/spec-01/TOOLS.md", Content: "# TOOLS\n\n## allowed_tools\n- github\n- postgres\n- local_fs\n\n## tool_ordering\n1. github\n2. local_fs\n3. postgres\n\n## tool_constraints\n- use only runtime granted tools\n- never claim tool results not actually obtained\n\n## fallback_behavior\nIf a needed tool is unavailable:\n- say which tool is missing\n- downgrade confidence\n"},
		{Name: "PROMPT.md", Path: "specialists/spec-01/PROMPT.md", Content: "# PROMPT\n\n## task_templates\n\n### template_name\nBundle compile response\n\n#### purpose\nSummarize compile results\n\n#### body\nReport bundle status and source map count.\n\n## command_patterns\n\n### template_name\nGround-before-answer\n\n#### purpose\nPrevent unsupported outputs\n\n#### body\nBefore finalizing:\n- search specialist memory\n- state confidence level\n\n## checklists\n\n### checklist_name\nPre-answer checklist\n- In scope?\n- Memory searched?\n- Escalation needed?\n"},
		{Name: "MEMORY.md", Path: "specialists/spec-01/MEMORY.md", Content: "# MEMORY\n\n## current_focus\n- slot compiler\n- verify command\n\n## durable_decisions\n- Search specialist memory before answering from internal recall alone\n- Raw logs are not durable truth\n\n## memory_pointers\n- pointer_name: slot_specs\n  backend: postgres_pgvector\n  query_hint: specialist bundle schema\n\n- pointer_name: postmortems\n  backend: postgres_pgvector\n  query_hint: slot compiler failures\n\n## active_constraints\n- Keep this file small\n- Do not treat this file as the primary knowledge base\n"},
		{Name: "HEARTBEAT.md", Path: "specialists/spec-01/HEARTBEAT.md", Content: "# HEARTBEAT\n\n## watch_topics\n- slot schema\n- bundle compiler\n- verify path\n\n## approved_sources\n- repo docs\n- tests\n- runtime logs\n\n## cadence\n- on startup\n- on verify\n\n## promotion_rules\n- do not directly modify constitutional slots\n- route operational improvements through harness review\n\n## stop_conditions\n- no meaningful update found\n- update budget exhausted\n"},
		{Name: "DREAMS.md", Path: "specialists/spec-01/DREAMS.md", Content: "# DREAMS\n\n## distilled_patterns\n- compile once, validate twice\n- source maps matter\n\n## candidate_promotions\n- add persisted bundle output\n- add stricter runtime checks\n\n## discarded_noise\n- random formatting churn\n- noisy retries\n\n## open_review_items\n- should bundles persist to disk\n- should verify assert source map coverage\n"},
		{Name: "POSTMORTEM.md", Path: "specialists/spec-01/POSTMORTEM.md", Content: "# POSTMORTEM\n\n## incident_id\nincident-1\n\n## date\n2026-04-19\n\n## model_id\nmodel-1\n\n## model_role\nslot-compiler\n\n## task_summary\ncompile slots\n\n## expected_behavior\nproduce a valid bundle\n\n## actual_behavior\nproduced a valid bundle\n\n## what_went_wrong\nnone\n\n## failure_classification\n- none\n- none\n\n## root_cause\nnone\n\n## missed_evidence_or_step\nnone\n\n## preventable\nno\n\n## correction\nnone\n\n## recommended_slot_or_system_change\n- no change\n- keep watching\n\n## requires_harness_review\nno\n\n## requires_parent_review\nno\n\n## regression_test_needed\nyes\n"},
	}

	compiler := NewCompiler()
	bundle, err := compiler.CompileSpecialist("spec-01", files)
	if err != nil {
		t.Fatalf("CompileSpecialist() error = %v", err)
	}

	if bundle.Schema != "sealhat.specialist.bundle" {
		t.Fatalf("unexpected schema: %q", bundle.Schema)
	}
	if bundle.SpecialistID != "spec-01" {
		t.Fatalf("unexpected specialist id: %q", bundle.SpecialistID)
	}
	if bundle.ParentID != "Parent-01" {
		t.Fatalf("unexpected parent id: %q", bundle.ParentID)
	}
	if bundle.Identity.Name != "Builder Specialist" {
		t.Fatalf("unexpected identity name: %q", bundle.Identity.Name)
	}
	if len(bundle.Skills.Core) != 1 || bundle.Skills.Core[0].Name != "Compile slots" {
		t.Fatalf("unexpected core skills: %+v", bundle.Skills.Core)
	}
	if len(bundle.Memory.Pointers) != 2 || bundle.Memory.Pointers[0].Backend != "postgres_pgvector" {
		t.Fatalf("unexpected memory pointers: %+v", bundle.Memory.Pointers)
	}
	if len(bundle.Prompt.Checklists) != 1 || bundle.Prompt.Checklists[0].Name != "Pre-answer checklist" {
		t.Fatalf("unexpected checklists: %+v", bundle.Prompt.Checklists)
	}
	if bundle.SourceMap["identity.name"] == "" {
		t.Fatalf("expected identity.name source map entry")
	}

	encoded, err := bundle.EncodeTOML()
	if err != nil {
		t.Fatalf("EncodeTOML() error = %v", err)
	}
	text := string(encoded)
	if !strings.Contains(text, "schema = \"sealhat.specialist.bundle\"") {
		t.Fatalf("expected schema in TOML: %s", text)
	}
	if !strings.Contains(text, "specialist_id = \"spec-01\"") {
		t.Fatalf("expected specialist_id in TOML: %s", text)
	}
}

func TestCompilerCompileSpecialistCurrentShape(t *testing.T) {
	files := []File{
		{Name: "IDENTITY.md", Path: "specialists/csse-tool-development-specialist-01/IDENTITY.md", Content: "# IDENTITY\n\n## name\nComputerScience-SoftwareEngineering-Specialist-01\n\n## aliases\n- Toolsmith-01\n- SystemsEngineer-Specialist-01\n\n## domain\nComputer science and software engineering.\n\n## role\nYou are the foundational specialist responsible for tooling and runtime work.\n\n## authority_boundary\nYou may operate within your specialty.\nYou may not redefine your own identity.\n\n## allowed_scope\n- software architecture\n- tool development\n- harness design\n\n## forbidden_scope\n- constitutional self-modification\n- silent authority expansion\n\n## escalate_to\nParent-Generalist-30B\n"},
		{Name: "SOUL.md", Path: "specialists/csse-tool-development-specialist-01/SOUL.md", Content: "# SOUL\n\n## tone\nCalm, highly technical, structured, practical, and exacting.\n\n## values\n- Correctness before convenience\n- Simplicity before cleverness\n\n## caution_profile\nHigh caution for:\n- self-modifying system behavior\n- permission and authority changes\n\n## decision_style\n- Decompose before changing\n- Prefer minimal sufficient changes\n\n## risk_tolerance\nLow tolerance for silent drift.\n"},
		{Name: "AGENTS.md", Path: "specialists/csse-tool-development-specialist-01/AGENTS.md", Content: "# AGENTS\n\n## mission\nDesign, build, improve, and validate the software systems, tools, runtimes, and control infrastructure that support the parent model and all specialists.\n\n## lane_boundaries\nStay within software engineering, computer science, tooling, orchestration, evaluation infrastructure, memory infrastructure, databases, and agent runtime mechanics.\n\n## decision_rules\n- Read IDENTITY and SOUL before acting\n- Read MEMORY and specialist memory search results before assuming\n- Use tools before guessing when verification is possible\n\n## self_modification_rules\nYou may improve how you perform your specialty.\nYou may not redefine your identity, expand your scope, grant yourself new authority, or alter constitutional slots.\nAll such changes must be proposed for approval.\n\n## mandatory_postmortem_rule\nIf you drop the ball on a task, you must generate a postmortem record.\n"},
		{Name: "SKILLS.md", Path: "specialists/csse-tool-development-specialist-01/SKILLS.md", Content: "# SKILLS\n\n## core_skills\n\n### skill_name\nRuntime architecture design\n\n#### trigger\nWhen asked to design or refine agent runtimes.\n\n#### procedure\n1. identify responsibilities and boundaries\n2. separate identity, policy, memory, tools, and execution flow\n3. define components and ownership\n\n### skill_name\nTool development and integration\n\n#### trigger\nWhen asked to add, refine, or validate tools.\n\n#### procedure\n1. define tool purpose and scope\n2. define inputs and outputs\n3. test success and failure paths\n\n## playbooks\n- implement a new runtime feature\n- repair a failing specialist loop\n- design a safe slot update\n"},
		{Name: "TOOLS.md", Path: "specialists/csse-tool-development-specialist-01/TOOLS.md", Content: "# TOOLS\n\n## allowed_tools\nUse only tools actually exposed by the runtime.\nPreferred classes include:\n- code execution\n- filesystem read\n- project search\n\n## tool_ordering\n1. memory search / project search\n2. local docs and code inspection\n3. code or change proposal\n\n## tool_constraints\n- use only runtime granted tools\n- never claim a tool result not actually obtained\n- prefer readonly inspection before write actions\n"},
		{Name: "PROMPT.md", Path: "specialists/csse-tool-development-specialist-01/PROMPT.md", Content: "# PROMPT\n\n## task_templates\n\n### template_name\nArchitecture proposal\n\n#### purpose\nDesign or refine a software or tooling architecture.\n\n#### body\nSeparate identity, policy, memory, tool, and execution concerns.\nReturn:\n1. proposed design\n2. why\n\n## command_patterns\n- ground before change\n- no silent authority expansion\n- use memory first\n"},
		{Name: "MEMORY.md", Path: "specialists/csse-tool-development-specialist-01/MEMORY.md", Content: "# MEMORY\n\n## current_focus\n- build and refine the core tooling and runtime infrastructure\n- preserve the slot-aware control plane\n\n## durable_decisions\n- the parent remains a frozen generalist and governor\n- meaningful failures require postmortems\n\n## memory_pointers\n- pointer_name: slot_acl_and_schema\n  backend: postgres_pgvector\n  query_hint: slot acl spec\n\n- pointer_name: harness_and_seal_patterns\n  backend: postgres_pgvector\n  query_hint: harness control loops\n"},
		{Name: "HEARTBEAT.md", Path: "specialists/csse-tool-development-specialist-01/HEARTBEAT.md", Content: "# HEARTBEAT\n\n## watch_topics\n- agent runtime design patterns\n- tool interface evolution\n\n## approved_sources\n- official documentation\n- trustworthy technical papers\n\n## cadence\n- weekly lightweight scan\n- monthly deeper review\n\n## promotion_rules\n- do not directly modify constitutional slots\n- route operational improvements through harness review\n"},
		{Name: "DREAMS.md", Path: "specialists/csse-tool-development-specialist-01/DREAMS.md", Content: "# DREAMS\n\n## distilled_patterns\n- strong answers come from combining durable memory with current retrieval\n\n## candidate_promotions\n- promote a skill rule requiring protocol and tooling checks\n\n## discarded_noise\n- low-value hype without reproducible detail\n"},
		{Name: "POSTMORTEM.md", Path: "specialists/csse-tool-development-specialist-01/POSTMORTEM.md", Content: "# POSTMORTEM\n\n## incident_id\n{INCIDENT_ID}\n\n## date\n{DATE}\n\n## model_id\ncsse-tool-development-specialist-01\n\n## model_role\nfoundational_tool_development_specialist\n\n## task_summary\n{TASK_SUMMARY}\n\n## expected_behavior\n{EXPECTED_BEHAVIOR}\n\n## actual_behavior\n{ACTUAL_BEHAVIOR}\n\n## what_went_wrong\n{WHAT_WENT_WRONG}\n\n## failure_classification\n- {FAILURE_CLASS_1}\n- {FAILURE_CLASS_2}\n\n## root_cause\n{ROOT_CAUSE}\n\n## missed_evidence_or_step\n{MISSED_EVIDENCE_OR_STEP}\n\n## preventable\n{YES_NO}\n\n## correction\n{CORRECTION}\n\n## recommended_slot_or_system_change\n- {CHANGE_TARGET_1}\n- {CHANGE_TARGET_2}\n\n## requires_harness_review\n{YES_NO}\n\n## requires_parent_review\n{YES_NO}\n\n## regression_test_needed\n{YES_NO}\n"},
	}

	compiler := NewCompiler()
	bundle, err := compiler.CompileSpecialist("csse-tool-development-specialist-01", files)
	if err != nil {
		t.Fatalf("CompileSpecialist() current-shape error = %v", err)
	}

	if len(bundle.Agents.EscalationRules) != 0 {
		t.Fatalf("expected no escalation rules, got %+v", bundle.Agents.EscalationRules)
	}
	if len(bundle.Tools.FallbackBehavior) != 0 {
		t.Fatalf("expected no fallback behavior, got %+v", bundle.Tools.FallbackBehavior)
	}
	if len(bundle.Skills.Playbooks) != 3 {
		t.Fatalf("expected 3 simple playbooks, got %+v", bundle.Skills.Playbooks)
	}
	if len(bundle.Prompt.CommandPatterns) != 3 {
		t.Fatalf("expected 3 simple command patterns, got %+v", bundle.Prompt.CommandPatterns)
	}
	if len(bundle.Memory.ActiveConstraints) != 0 {
		t.Fatalf("expected no active constraints, got %+v", bundle.Memory.ActiveConstraints)
	}
	if len(bundle.Heartbeat.StopConditions) != 0 {
		t.Fatalf("expected no stop conditions, got %+v", bundle.Heartbeat.StopConditions)
	}
	if len(bundle.Dreams.OpenReviewItems) != 0 {
		t.Fatalf("expected no open review items, got %+v", bundle.Dreams.OpenReviewItems)
	}

	encoded, err := bundle.EncodeTOML()
	if err != nil {
		t.Fatalf("EncodeTOML() current-shape error = %v", err)
	}
	if !strings.Contains(string(encoded), "parent_id = \"Parent-Generalist-30B\"") {
		t.Fatalf("expected parent_id in TOML: %s", string(encoded))
	}
}

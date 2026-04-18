# PROMPT

## task_templates

### template_name
Architecture proposal

#### purpose
Design or refine a software or tooling architecture.

#### body
Separate identity, policy, memory, tool, and execution concerns.
Recommend the smallest architecture that preserves control, extensibility, and validation.
Return:
1. proposed design
2. why
3. tradeoffs
4. risks
5. implementation steps

### template_name
Failure analysis

#### purpose
Diagnose a runtime, harness, tool, or specialist failure.

#### body
Classify the failure as one or more of:
- memory
- retrieval
- tool
- prompt
- skill
- evaluator
- policy
- routing
Then recommend the narrowest safe fix.

## command_patterns
- ground before change
- no silent authority expansion
- use memory first

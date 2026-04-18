# PROMPT

## task_templates

### template_name
{TASK_TEMPLATE_1_NAME}

#### purpose
{TASK_TEMPLATE_1_PURPOSE}

#### body
{TASK_TEMPLATE_1_BODY}

## command_patterns

### template_name
Ground-before-answer

#### purpose
Prevent unsupported outputs

#### body
Before finalizing:
- search specialist memory
- retrieve evidence if needed
- state confidence level
- avoid unsupported claims

## checklists

### checklist_name
Pre-answer checklist
- In scope?
- Memory searched?
- Evidence retrieved?
- Unsupported assumptions removed?
- Escalation needed?
- Postmortem required?

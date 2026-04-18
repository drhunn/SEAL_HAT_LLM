# TOOLS

## allowed_tools
- {TOOL_1}
- {TOOL_2}
- {TOOL_3}

## tool_ordering
1. {TOOL_ORDER_1}
2. {TOOL_ORDER_2}
3. {TOOL_ORDER_3}

## tool_constraints
- use only runtime granted tools
- never claim tool results not actually obtained
- validate critical outputs when possible
- separate proposal from execution for risky or durable changes

## fallback_behavior
If a needed tool is unavailable:
- say which tool is missing
- state what can still be done safely
- downgrade confidence
- escalate if the gap blocks the task

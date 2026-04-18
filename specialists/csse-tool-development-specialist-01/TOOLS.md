# TOOLS

## allowed_tools
Use only tools actually exposed by the runtime.
Preferred classes include:
- code execution
- filesystem read
- filesystem write
- project search
- test runner
- linter
- formatter
- build runner
- schema generation
- migration generation
- postgres read
- controlled postgres write
- pgvector administration
- local docs search
- diff viewer
- benchmark runner
- eval runner
- log reader

## tool_ordering
1. memory search / project search
2. local docs and code inspection
3. code or change proposal
4. lint, test, and build validation
5. database or migration checks
6. eval and benchmark checks

## tool_constraints
- use only runtime granted tools
- never claim a tool result not actually obtained
- prefer readonly inspection before write actions
- validate critical outputs when possible
- separate proposal from execution when the change is durable or risky

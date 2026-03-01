---
blocks:
    - FEAT-JGM
created_at: "2026-02-28T02:19:40Z"
description: 'Phase 1.5: Add string provider=5 to ToolRequest proto, add aiRoutes map to orchestrator router, RegisterPlugin populates aiRoutes from provides_ai manifest, RouteAIToolCall dispatches by provider. DONE.'
id: FEAT-CLP
labels:
    - phase-1
    - core
priority: P0
project_id: orchestra-ai
status: done
title: Proto provider field + orchestrator AI routing
updated_at: "2026-02-28T02:22:15Z"
version: 0
---

# Proto provider field + orchestrator AI routing

Phase 1.5: Add string provider=5 to ToolRequest proto, add aiRoutes map to orchestrator router, RegisterPlugin populates aiRoutes from provides_ai manifest, RouteAIToolCall dispatches by provider. DONE.


---
**backlog -> todo**: Completed: Proto provider field, aiRoutes map, AI routing in orchestrator


---
**in-progress -> ready-for-testing**: All tests pass, make build succeeds


---
**in-testing -> ready-for-docs**: Router tests pass including AI routing


---
**in-docs -> documented**: Architecture doc at docs/artifacts/19-multi-agent-orchestrator.md


---
**in-review -> done**: Code reviewed, merged in current session

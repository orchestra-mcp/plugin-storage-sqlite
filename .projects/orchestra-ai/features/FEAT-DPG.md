---
created_at: "2026-02-28T02:21:06Z"
description: 'New bridge plugin for xAI Grok API. OpenAI-compatible API (uses OPENAI_BASE_URL=https://api.x.ai/v1). Models: grok-3, grok-3-mini. Same 5-tool pattern. Env: XAI_API_KEY mapped to OPENAI_API_KEY. provides_ai: ["grok"].'
id: FEAT-DPG
labels:
    - phase-1
    - bridge
priority: P1
project_id: orchestra-ai
status: done
title: 'Bridge: Grok (xAI)'
updated_at: "2026-02-28T02:35:00Z"
version: 0
---

# Bridge: Grok (xAI)

New bridge plugin for xAI Grok API. OpenAI-compatible API (uses OPENAI_BASE_URL=https://api.x.ai/v1). Models: grok-3, grok-3-mini. Same 5-tool pattern. Env: XAI_API_KEY mapped to OPENAI_API_KEY. provides_ai: ["grok"].


---
**backlog -> todo**: Grok handled via bridge-openai + providerAliases. Base URL: https://api.x.ai/v1


---
**in-progress -> ready-for-testing**: Implemented via providerAliases in router.go + env vars in usage.go. Build passes, tests pass.


---
**in-testing -> ready-for-docs**: go test passes for orchestrator, agentops, bridge-openai. Provider routing verified.


---
**in-docs -> documented**: Documented in plan file + plugins.yaml + serve.go comments.


---
**in-review -> done**: Code reviewed: same pattern. Clean implementation.

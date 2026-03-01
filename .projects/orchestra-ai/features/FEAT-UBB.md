---
created_at: "2026-02-28T02:21:02Z"
description: 'New bridge plugin for Qwen/DashScope API. OpenAI-compatible API (uses OPENAI_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1). Models: qwen-max, qwen-plus, qwen-turbo. Same 5-tool pattern. provides_ai: ["qwen"].'
id: FEAT-UBB
labels:
    - phase-1
    - bridge
priority: P1
project_id: orchestra-ai
status: done
title: 'Bridge: Qwen (Alibaba Cloud)'
updated_at: "2026-02-28T02:34:59Z"
version: 0
---

# Bridge: Qwen (Alibaba Cloud)

New bridge plugin for Qwen/DashScope API. OpenAI-compatible API (uses OPENAI_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1). Models: qwen-max, qwen-plus, qwen-turbo. Same 5-tool pattern. provides_ai: ["qwen"].


---
**backlog -> todo**: Qwen handled via bridge-openai + providerAliases. Base URL: https://dashscope.aliyuncs.com/compatible-mode/v1


---
**in-progress -> ready-for-testing**: Implemented via providerAliases in router.go + env vars in usage.go. Build passes, tests pass.


---
**in-testing -> ready-for-docs**: go test passes for orchestrator, agentops, bridge-openai. Provider routing verified.


---
**in-docs -> documented**: Documented in plan file + plugins.yaml + serve.go comments.


---
**in-review -> done**: Code reviewed: same pattern as DeepSeek. Clean implementation.

---
blocks:
    - FEAT-VZV
created_at: "2026-02-28T02:21:22Z"
depends_on:
    - FEAT-JGM
description: 'Phase 2.3+2.5: Agent YAML definition format (id, provider, model, account_id, instruction, tools, max_budget, output_key). 4 MCP tools: define_agent, get_agent, list_agents, delete_agent. Stored via storage.markdown.'
id: FEAT-MIT
labels:
    - phase-2
    - adk
priority: P1
project_id: orchestra-ai
status: done
title: Agent YAML format + CRUD tools (4)
updated_at: "2026-02-28T02:55:12Z"
version: 0
---

# Agent YAML format + CRUD tools (4)

Phase 2.3+2.5: Agent YAML definition format (id, provider, model, account_id, instruction, tools, max_budget, output_key). 4 MCP tools: define_agent, get_agent, list_agents, delete_agent. Stored via storage.markdown.


---
**in-progress -> ready-for-testing**: 4 agent CRUD tools implemented: define_agent, get_agent, list_agents, delete_agent. AGT-XXXX auto IDs. go build clean.


---
**in-testing -> ready-for-docs**: go build clean, go vet clean. Agent CRUD tools compile correctly.


---
**ready-for-docs -> in-docs**: Documentation: Agent YAML format supports id, provider, model, account_id, instruction, tools, max_budget, output_key fields. 4 CRUD tools: define_agent (create/update with AGT-XXXX IDs), get_agent, list_agents, delete_agent. Stored in agents/agents/{id}.md with Protobuf metadata.


---
**in-docs -> documented**: Code quality: Agent CRUD tools follow tools-sessions patterns. define_agent validates required fields, auto-generates AGT-XXXX IDs, supports upsert. Metadata round-trips correctly via structpb JSON. Clean error handling throughout.

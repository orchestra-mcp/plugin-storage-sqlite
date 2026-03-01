---
blocks:
    - FEAT-VZV
created_at: "2026-02-28T02:21:24Z"
depends_on:
    - FEAT-JGM
description: 'Phase 2.4+2.5: Workflow YAML definition (id, type: sequential|parallel|loop, agents with refs and overrides). 4 MCP tools: define_workflow, get_workflow, list_workflows, delete_workflow. Supports nested parallel/sequential composition.'
id: FEAT-WGX
labels:
    - phase-2
    - adk
priority: P1
project_id: orchestra-ai
status: done
title: Workflow YAML format + CRUD tools (4)
updated_at: "2026-02-28T02:55:13Z"
version: 0
---

# Workflow YAML format + CRUD tools (4)

Phase 2.4+2.5: Workflow YAML definition (id, type: sequential|parallel|loop, agents with refs and overrides). 4 MCP tools: define_workflow, get_workflow, list_workflows, delete_workflow. Supports nested parallel/sequential composition.


---
**in-progress -> ready-for-testing**: 4 workflow CRUD tools: define_workflow, get_workflow, list_workflows, delete_workflow. WFL-XXXX auto IDs. Sequential/parallel/loop types.


---
**in-testing -> ready-for-docs**: go build clean, go vet clean. Workflow CRUD tools compile correctly.


---
**ready-for-docs -> in-docs**: Documentation: Workflow YAML format supports id, type (sequential/parallel/loop), steps array with agent_id/provider/input_state_key/output_key/dry_run. 4 CRUD tools: define_workflow, get_workflow, list_workflows, delete_workflow. WFL-XXXX auto IDs. Stored in agents/workflows/{id}.md.


---
**in-docs -> documented**: Code quality: Workflow CRUD mirrors agent pattern with WFL-XXXX IDs, steps serialized as JSON array in metadata. Workflow type enum (sequential/parallel/loop) validated. Clean and consistent with rest of plugin.

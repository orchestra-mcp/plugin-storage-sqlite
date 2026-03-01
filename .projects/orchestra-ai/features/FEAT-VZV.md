---
blocks:
    - FEAT-KNF
created_at: "2026-02-28T02:21:27Z"
depends_on:
    - FEAT-MIT
    - FEAT-WGX
description: 'Phase 2.5: 5 MCP tools for execution: run_agent (single agent), run_workflow (full pipeline), get_run_status, list_runs, cancel_run. Tracks run state, token usage, cost per step. Async execution with polling.'
id: FEAT-VZV
labels:
    - phase-2
    - adk
priority: P1
project_id: orchestra-ai
status: done
title: 'Execution tools: run_agent, run_workflow, status, cancel (5)'
updated_at: "2026-02-28T02:55:13Z"
version: 0
---

# Execution tools: run_agent, run_workflow, status, cancel (5)

Phase 2.5: 5 MCP tools for execution: run_agent (single agent), run_workflow (full pipeline), get_run_status, list_runs, cancel_run. Tracks run state, token usage, cost per step. Async execution with polling.


---
**in-progress -> ready-for-testing**: 5 execution tools: run_agent, run_workflow, get_run_status, list_runs, cancel_run. Provider routing, dry_run mode, usage reporting.


---
**in-testing -> ready-for-docs**: go build clean, go vet clean. Execution tools compile correctly with provider routing.


---
**ready-for-docs -> in-docs**: Documentation: 5 execution tools — run_agent (calls ai_prompt with provider routing, records run, reports usage), run_workflow (sequential steps with state passing via output_key), get_run_status, list_runs, cancel_run. All support dry_run=true for mock execution. Run IDs: RUN-XXXX format. Stored in agents/runs/{id}.md.


---
**in-docs -> documented**: Code quality: Execution tools correctly chain cross-plugin calls. run_workflow passes output_key state between sequential steps. Run records stored atomically with status tracking (running/completed/failed/cancelled). dry_run mode returns canned responses for CI use.

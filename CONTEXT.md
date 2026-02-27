# CONTEXT.md

Project context managed by Orchestra MCP.

## Orchestra MCP Integration

This project uses Orchestra MCP for AI-driven project management with 85 tools covering:

- **Project hierarchy**: Epics → Stories → Tasks (full CRUD)
- **13-state workflow**: backlog through done with gated testing, docs, and review stages
- **Multi-audience PRD**: 4 audience types (business/product/technical/qa) with conditional follow-up questions, validation, agent briefings, auto-backlog generation
- **PRD templates**: Save/load reusable PRD templates
- **Sprint management**: Create/start/end sprints with auto task promotion, velocity tracking
- **Scrum ceremonies**: Standup summaries, burndown charts, retrospectives, backlog reordering
- **WIP limits**: Configurable max in-progress tasks, enforced on set_current_task
- **Dependencies**: Task dependency graph with blocker/blocked-by relationships
- **Metadata**: Assignment, labels, estimates, external links
- **Memory system**: Persistent project knowledge base across sessions
- **Session tracking**: What happened, what's next
- **Usage tracking**: Token and cost monitoring

## Workflow States

```
backlog → todo → in-progress → blocked
                             → ready-for-testing → in-testing
                             → ready-for-docs → in-docs → documented
                             → in-review → done / rejected / cancelled
```

4 gated transitions require evidence: in-progress→ready-for-testing, in-testing→ready-for-docs, in-docs→documented, in-review→done.

## Data Storage

- `.projects/` — Project data (epics, stories, tasks) in markdown frontmatter format
- `.projects/{name}/.memory/` — Memory chunks and session logs
- `.projects/{name}/.plans/` — Saved implementation plans
- `.projects/{name}/sprints/` — Sprint data and retrospectives
- `.projects/{name}/templates/` — PRD and issue templates
- `.projects/.events/` — Hook event log
- `.mcp.json` — MCP server configuration

## Skills & Agents

- `.claude/skills/` — Domain-specific skill patterns (auto-activated by context)
- `.claude/agents/` — Specialized agent configurations
- `.claude/hooks/` — Claude Code hook scripts

Use `list_skills` and `list_agents` MCP tools to discover what's installed.

## Session Protocol

1. Start: `get_project_status` → `get_standup_summary` → `check_wip_limit` → `get_next_task`
2. Work: `set_current_task` → write code → `advance_task` (with evidence at gates)
3. End: `save_session` → `save_memory`

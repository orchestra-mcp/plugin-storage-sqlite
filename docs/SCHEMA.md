# Database Schema

## Overview

The `storage.sqlite` plugin uses two databases:

1. **Workspace DB** at `~/.orchestra/db/<hash>.db` — Project data scoped to a workspace
2. **Global DB** at `~/.orchestra/db/global.db` — Cross-workspace config (accounts, workspaces, current user)

The workspace hash is the first 16 hex chars of SHA-256(absolute workspace path).

## Workspace Database

### projects

| Column | Type | Description |
|--------|------|-------------|
| slug | TEXT PK | URL-safe project identifier |
| name | TEXT | Display name |
| description | TEXT | Project description |
| metadata | TEXT | JSON object for flexible data |
| body | TEXT | Markdown body |
| version | INTEGER | CAS version |
| created_at | TEXT | ISO 8601 timestamp |
| updated_at | TEXT | ISO 8601 timestamp |

### features

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | Feature ID (FEAT-XXX) |
| project_id | TEXT FK | References projects(slug) |
| title | TEXT | Feature title |
| description | TEXT | Short description |
| status | TEXT | Workflow status (backlog, todo, in-progress, ...) |
| priority | TEXT | P0-P3 |
| kind | TEXT | feature, bug, hotfix, chore, testcase |
| assignee | TEXT | Person name or ID |
| estimate | TEXT | S, M, L, XL |
| labels | TEXT | JSON array of strings |
| depends_on | TEXT | JSON array of feature IDs |
| body | TEXT | Markdown body (notes, gate evidence) |
| version | INTEGER | CAS version |
| created_at | TEXT | ISO 8601 timestamp |
| updated_at | TEXT | ISO 8601 timestamp |

**Indexes:** `idx_features_project`, `idx_features_status`, `idx_features_assignee`

### persons

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | Person ID (PERS-XXX) |
| project_id | TEXT FK | References projects(slug) |
| name | TEXT | Full name |
| email | TEXT | Email address |
| role | TEXT | developer, qa, reviewer, lead |
| status | TEXT | active, inactive |
| bio | TEXT | Short bio |
| github_email | TEXT | GitHub email for commits |
| integrations | TEXT | JSON object (jira_email, slack_id, timezone, etc.) |
| labels | TEXT | JSON array of strings |
| body | TEXT | Markdown body |
| version | INTEGER | CAS version |
| created_at | TEXT | ISO 8601 timestamp |
| updated_at | TEXT | ISO 8601 timestamp |

### plans

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | Plan ID (PLAN-XXX) |
| project_id | TEXT | Project slug |
| title | TEXT | Plan title |
| description | TEXT | Plan description |
| status | TEXT | draft, approved, in-progress, completed |
| features | TEXT | JSON array of feature IDs |
| body | TEXT | Markdown body |
| version | INTEGER | CAS version |
| created_at | TEXT | ISO 8601 timestamp |
| updated_at | TEXT | ISO 8601 timestamp |

### requests

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | Request ID (REQ-XXX) |
| project_id | TEXT | Project slug |
| title | TEXT | Request title |
| description | TEXT | Request description |
| kind | TEXT | feature, hotfix, bug |
| status | TEXT | pending, converted, dismissed |
| priority | TEXT | P0-P3 |
| body | TEXT | Markdown body |
| version | INTEGER | CAS version |
| created_at | TEXT | ISO 8601 timestamp |
| updated_at | TEXT | ISO 8601 timestamp |

### assignment_rules

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | Rule ID (RULE-XXX) |
| project_id | TEXT | Project slug |
| kind | TEXT | Feature kind to match |
| person_id | TEXT | Person to auto-assign |
| body | TEXT | Markdown body |
| version | INTEGER | CAS version |
| created_at | TEXT | ISO 8601 timestamp |
| updated_at | TEXT | ISO 8601 timestamp |

### notes

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | Note ID |
| project_id | TEXT | Project slug |
| title | TEXT | Note title |
| body | TEXT | Markdown body |
| pinned | INTEGER | 0 or 1 |
| deleted | INTEGER | 0 or 1 (soft delete) |
| tags | TEXT | JSON array of strings |
| icon | TEXT | Icon name |
| color | TEXT | Color label |
| version | INTEGER | CAS version |
| created_at | TEXT | ISO 8601 timestamp |
| updated_at | TEXT | ISO 8601 timestamp |

### wip_config

| Column | Type | Description |
|--------|------|-------------|
| project_id | TEXT PK | Project slug |
| max_in_progress | INTEGER | Max features in-progress |
| version | INTEGER | CAS version |
| updated_at | TEXT | ISO 8601 timestamp |

### sessions

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | Session UUID |
| account_id | TEXT | AI account ID |
| name | TEXT | Session name |
| workspace | TEXT | Workspace path |
| model | TEXT | AI model name |
| permission_mode | TEXT | plan, auto-edit, etc. |
| allowed_tools | TEXT | JSON array of tool names |
| max_budget | REAL | Max budget in USD |
| system_prompt | TEXT | Custom system prompt |
| status | TEXT | active, paused, completed |
| message_count | INTEGER | Total messages |
| total_tokens_in | INTEGER | Total input tokens |
| total_tokens_out | INTEGER | Total output tokens |
| total_cost_usd | REAL | Total cost |
| claude_session_id | TEXT | CLI session ID |
| last_message_at | TEXT | ISO 8601 timestamp |
| body | TEXT | Markdown body |
| created_at | TEXT | ISO 8601 timestamp |

### session_turns

| Column | Type | Description |
|--------|------|-------------|
| session_id | TEXT PK | References sessions(id) |
| turn_number | INTEGER PK | Turn sequence number |
| user_prompt | TEXT | User's message |
| response | TEXT | AI response |
| tokens_in | INTEGER | Input tokens for this turn |
| tokens_out | INTEGER | Output tokens for this turn |
| cost_usd | REAL | Cost for this turn |
| model | TEXT | Model used |
| duration_ms | INTEGER | Response time |
| timestamp | TEXT | ISO 8601 timestamp |

### packs

| Column | Type | Description |
|--------|------|-------------|
| name | TEXT PK | Pack name |
| version | TEXT | Installed version |
| repo | TEXT | GitHub repo URL |
| installed_at | TEXT | ISO 8601 timestamp |
| metadata | TEXT | JSON object |

### stacks

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER PK | Always 1 (singleton) |
| stacks | TEXT | JSON array of stack names |
| version | INTEGER | CAS version |

### kv_store

Generic key-value fallback for paths that don't match a typed table.

| Column | Type | Description |
|--------|------|-------------|
| path | TEXT PK | Original storage path |
| content | BLOB | Raw content |
| metadata | TEXT | JSON object |
| version | INTEGER | CAS version |
| created_at | TEXT | ISO 8601 timestamp |
| updated_at | TEXT | ISO 8601 timestamp |

## SQLite Pragmas

```sql
PRAGMA journal_mode=WAL;
PRAGMA busy_timeout=5000;
PRAGMA foreign_keys=ON;
PRAGMA synchronous=NORMAL;
PRAGMA cache_size=-20000;   -- 20MB cache
PRAGMA temp_store=MEMORY;
```

## Path Routing

Storage paths are parsed and routed to the corresponding table:

```
"my-app/features/FEAT-ABC.md"  →  features  (project=my-app, id=FEAT-ABC)
"bridge/sessions/<uuid>.md"    →  sessions  (id=<uuid>)
".packs/registry.json"         →  packs
"anything/else.txt"            →  kv_store  (path=anything/else.txt)
```

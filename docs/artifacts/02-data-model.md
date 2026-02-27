# Data Model — Orchestra Reference

> Extracted from `orch-ref/`. Every entity type, its fields, relationships, and storage location.

---

## Table of Contents

1. [File-Based Entities (TOON/YAML in .projects/)](#file-based-entities-toonyaml-in-projects)
2. [Notification System Types](#notification-system-types)
3. [MCP Protocol Types](#mcp-protocol-types)
4. [GitHub Integration Types](#github-integration-types)
5. [Linear Integration Types](#linear-integration-types)
6. [Jira Integration Types](#jira-integration-types)
7. [Usage & Hook Tracking Types](#usage--hook-tracking-types)
8. [PostgreSQL Entities (Cloud)](#postgresql-entities-cloud)
9. [SQLite Entities (Local Desktop)](#sqlite-entities-local-desktop)
10. [Sync Protocol Types](#sync-protocol-types)
11. [gRPC/Proto Messages](#grpcproto-messages)
12. [Workflow States](#workflow-states)

---

## File-Based Entities (TOON/YAML in .projects/)

These entities are stored as YAML frontmatter + Markdown body files within the `.projects/` directory tree. The MCP plugin reads/writes them directly on the local filesystem.

### ProjectStatus

Root tracking file for a project.

**Source:** `orch-ref/app/types/data.go`
**Storage:** `.projects/{slug}/status.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Project | `string` | `project` | `project` | Yes | Human-readable project name |
| Slug | `string` | `slug` | `slug` | Yes | URL-safe project identifier |
| Status | `string` | `status` | `status` | Yes | Current project status |
| Description | `string` | `description,omitempty` | `description,omitempty` | No | Project description |
| CreatedAt | `string` | `created_at` | `created_at` | Yes | ISO 8601 creation timestamp |
| UpdatedAt | `string` | `updated_at,omitempty` | `updated_at,omitempty` | No | ISO 8601 last-update timestamp |
| NextID | `int` | `next_id,omitempty` | `next_id,omitempty` | No | Auto-increment counter for issue IDs |
| Epics | `[]IssueEntry` | `epics,omitempty` | `epics,omitempty` | No | Summary list of epics |
| Stories | `[]IssueEntry` | `stories,omitempty` | `stories,omitempty` | No | Summary list of stories |
| Tasks | `[]IssueEntry` | `tasks,omitempty` | `tasks,omitempty` | No | Summary list of tasks |
| Sprints | `[]SprintEntry` | `sprints,omitempty` | `sprints,omitempty` | No | Summary list of sprints |

### IssueEntry

Summary row in project status referencing an epic, story, or task.

**Source:** `orch-ref/app/types/data.go`
**Storage:** Inline within `ProjectStatus`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| ID | `string` | `id` | `id` | Yes | Issue identifier (e.g., "TASK-1") |
| Title | `string` | `title` | `title` | Yes | Issue title |
| Status | `string` | `status` | `status` | Yes | Current workflow status |

### SprintEntry

Summary row for a sprint in project status.

**Source:** `orch-ref/app/types/data.go`
**Storage:** Inline within `ProjectStatus`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| ID | `string` | `id` | `id` | Yes | Sprint identifier |
| Name | `string` | `name` | `name` | Yes | Sprint name |
| Status | `string` | `status` | `status` | Yes | Sprint status (planned/active/completed) |

### IssueData

Full data for any issue (epic, story, task, or bug). Stored as YAML frontmatter with Markdown body for description.

**Source:** `orch-ref/app/types/data.go`
**Storage:** `.projects/{slug}/epics/{id}.md`, `.projects/{slug}/stories/{id}.md`, `.projects/{slug}/tasks/{id}.md`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| ID | `string` | `id` | `id` | Yes | Issue identifier |
| Title | `string` | `title` | `title` | Yes | Issue title |
| Type | `string` | `type` | `type` | Yes | Issue type: "epic", "story", "task", "bug" |
| Status | `string` | `status` | `status` | Yes | Current workflow status (one of 13 states) |
| Description | `string` | `-` (body) | `description,omitempty` | No | Markdown body (not in YAML frontmatter) |
| Priority | `string` | `priority,omitempty` | `priority,omitempty` | No | Priority level |
| CreatedAt | `string` | `created_at` | `created_at` | Yes | ISO 8601 creation timestamp |
| UpdatedAt | `string` | `updated_at,omitempty` | `updated_at,omitempty` | No | ISO 8601 last-update timestamp |
| Children | `[]IssueChild` | `children,omitempty` | `children,omitempty` | No | Child issue references |
| DependsOn | `[]string` | `depends_on,omitempty` | `depends_on,omitempty` | No | IDs of blocking issues |
| AssignedTo | `string` | `assigned_to,omitempty` | `assigned_to,omitempty` | No | Assignee identifier |
| Labels | `[]string` | `labels,omitempty` | `labels,omitempty` | No | Classification labels |
| Estimate | `float64` | `estimate,omitempty` | `estimate,omitempty` | No | Story points or hours estimate |
| StartedAt | `string` | `started_at,omitempty` | `started_at,omitempty` | No | When work began |
| CompletedAt | `string` | `completed_at,omitempty` | `completed_at,omitempty` | No | When work completed |
| Links | `[]IssueLink` | `links,omitempty` | `links,omitempty` | No | References to PRs, commits, URLs |
| Evidence | `[]EvidenceEntry` | `evidence,omitempty` | `evidence,omitempty` | No | Gate transition evidence records |
| SprintID | `string` | `sprint_id,omitempty` | `sprint_id,omitempty` | No | Associated sprint |
| Rank | `int` | `rank,omitempty` | `rank,omitempty` | No | Backlog ordering position |

### IssueChild

Child reference stored on a parent issue.

**Source:** `orch-ref/app/types/data.go`
**Storage:** Inline within `IssueData.Children`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| ID | `string` | `id` | `id` | Yes | Child issue ID |
| Title | `string` | `title` | `title` | Yes | Child issue title |
| Status | `string` | `status` | `status` | Yes | Child issue status |

### IssueLink

Reference to a PR, commit, URL, or related issue.

**Source:** `orch-ref/app/types/data.go`
**Storage:** Inline within `IssueData.Links`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Type | `string` | `type` | `type` | Yes | Link type: "pr", "commit", "url", "issue" |
| URL | `string` | `url` | `url` | Yes | Link URL |
| Title | `string` | `title,omitempty` | `title,omitempty` | No | Display title |

### EvidenceEntry

Evidence submitted at a gate transition.

**Source:** `orch-ref/app/types/data.go`
**Storage:** Inline within `IssueData.Evidence`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Gate | `string` | `gate` | `gate` | Yes | Gate name (e.g., "testing", "docs", "review") |
| Description | `string` | `description` | `description` | Yes | Evidence description |
| Files | `[]string` | `files,omitempty` | `files,omitempty` | No | File paths related to the evidence |
| Verified | `[]string` | `verified,omitempty` | `verified,omitempty` | No | Verification checks passed |
| Timestamp | `string` | `timestamp` | `timestamp` | Yes | When evidence was submitted |

### Sprint

Timeboxed iteration.

**Source:** `orch-ref/app/types/data.go`
**Storage:** `.projects/{slug}/sprints/{id}.md` (via `SprintData`)

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| ID | `string` | `id` | `id` | Yes | Sprint identifier |
| Name | `string` | `name` | `name` | Yes | Sprint name |
| Status | `string` | `status` | `status` | Yes | "planned", "active", or "completed" |
| StartDate | `string` | `start_date,omitempty` | `start_date,omitempty` | No | Sprint start date |
| EndDate | `string` | `end_date,omitempty` | `end_date,omitempty` | No | Sprint end date |
| Goal | `string` | `goal,omitempty` | `goal,omitempty` | No | Sprint goal statement |
| TaskIDs | `[]string` | `task_ids,omitempty` | `task_ids,omitempty` | No | IDs of tasks in this sprint |
| Velocity | `float64` | `velocity,omitempty` | `velocity,omitempty` | No | Calculated on completion (story points done) |
| CreatedAt | `string` | `created_at` | `created_at` | Yes | ISO 8601 creation timestamp |

### SprintData

Stored wrapper for Sprint in markdown files. Inlines all Sprint fields.

**Source:** `orch-ref/app/types/data.go`
**Storage:** `.projects/{slug}/sprints/{id}.md`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| (Sprint) | `Sprint` | `,inline` | (all Sprint fields) | Yes | All Sprint fields inlined |

### Template

Reusable issue template.

**Source:** `orch-ref/app/types/data.go`
**Storage:** `.projects/{slug}/templates/{name}.md`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Name | `string` | `name` | `name` | Yes | Template name |
| Type | `string` | `type` | `type` | Yes | Issue type: "epic", "story", "task", "bug" |
| Title | `string` | `title,omitempty` | `title,omitempty` | No | Default title |
| Priority | `string` | `priority,omitempty` | `priority,omitempty` | No | Default priority |
| Labels | `[]string` | `labels,omitempty` | `labels,omitempty` | No | Default labels |
| Description | `string` | `-` (body) | `description,omitempty` | No | Markdown body (not in YAML frontmatter) |

### Retrospective

Sprint retrospective data.

**Source:** `orch-ref/app/types/scrum.go`
**Storage:** `.projects/{slug}/sprints/{id}-retro.md`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| SprintID | `string` | `sprint_id` | `sprint_id` | Yes | Associated sprint ID |
| SprintName | `string` | `sprint_name` | `sprint_name` | Yes | Sprint name |
| Date | `string` | `date` | `date` | Yes | Retrospective date |
| WentWell | `[]string` | `went_well,omitempty` | `went_well,omitempty` | No | Things that went well |
| DidntGoWell | `[]string` | `didnt_go_well,omitempty` | `didnt_go_well,omitempty` | No | Things that did not go well |
| ActionItems | `[]ActionItem` | `action_items,omitempty` | `action_items,omitempty` | No | Follow-up action items |

### ActionItem

Retro action item.

**Source:** `orch-ref/app/types/scrum.go`
**Storage:** Inline within `Retrospective.ActionItems`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Action | `string` | `action` | `action` | Yes | Action description |
| Owner | `string` | `owner,omitempty` | `owner,omitempty` | No | Person responsible |
| DueDate | `string` | `due_date,omitempty` | `due_date,omitempty` | No | Due date |
| Status | `string` | `status` | `status` | Yes | Action status |

### ProjectSettings

Project-level configuration.

**Source:** `orch-ref/app/types/scrum.go`
**Storage:** `.projects/{slug}/settings.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| WipLimits | `WipLimits` | `wip_limits,omitempty` | `wip_limits,omitempty` | No | Work-in-progress limit configuration |

### WipLimits

Maximum in-progress items configuration.

**Source:** `orch-ref/app/types/scrum.go`
**Storage:** Inline within `ProjectSettings.WipLimits`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| MaxInProgress | `int` | `max_in_progress,omitempty` | `max_in_progress,omitempty` | No | Global max in-progress tasks |
| MaxPerAssignee | `int` | `max_per_assignee,omitempty` | `max_per_assignee,omitempty` | No | Max in-progress per assignee |
| MaxPerSprint | `int` | `max_per_sprint,omitempty` | `max_per_sprint,omitempty` | No | Max tasks per sprint |

### MemoryChunk

Piece of project context stored for RAG retrieval.

**Source:** `orch-ref/app/types/memory.go`
**Storage:** `.projects/{slug}/memory/index.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| ID | `string` | `id` | `id` | Yes | Chunk identifier |
| Project | `string` | `project` | `project` | Yes | Project slug |
| Source | `string` | `source` | `source` | Yes | Source type: "task", "prd", "session", "user" |
| SourceID | `string` | `source_id` | `source_id` | Yes | Source entity ID |
| Summary | `string` | `summary` | `summary` | Yes | Short summary |
| Content | `string` | `content` | `content` | Yes | Full content |
| Tags | `[]string` | `tags,omitempty` | `tags,omitempty` | No | Searchable tags |
| CreatedAt | `string` | `created_at` | `created_at` | Yes | ISO 8601 creation timestamp |

### MemoryIndex

All memory chunks for a project.

**Source:** `orch-ref/app/types/memory.go`
**Storage:** `.projects/{slug}/memory/index.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Chunks | `[]MemoryChunk` | `chunks` | `chunks` | Yes | All stored memory chunks |

### SessionLog

Record of a Claude Code session for a project.

**Source:** `orch-ref/app/types/memory.go`
**Storage:** `.projects/{slug}/sessions/index.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| SessionID | `string` | `session_id` | `session_id` | Yes | Session identifier |
| Project | `string` | `project` | `project` | Yes | Project slug |
| Summary | `string` | `summary` | `summary` | Yes | Session summary |
| Events | `[]SessionEvent` | `events,omitempty` | `events,omitempty` | No | Session event log |
| StartedAt | `string` | `started_at` | `started_at` | Yes | Session start time |
| EndedAt | `string` | `ended_at,omitempty` | `ended_at,omitempty` | No | Session end time |

### SessionEvent

Single event within a session log.

**Source:** `orch-ref/app/types/memory.go`
**Storage:** Inline within `SessionLog.Events`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Type | `string` | `type` | `type` | Yes | Event type: "tool_call", "decision", "output" |
| Summary | `string` | `summary` | `summary` | Yes | Event summary |
| Timestamp | `string` | `timestamp` | `timestamp` | Yes | ISO 8601 timestamp |

### SessionIndex

All session logs for a project.

**Source:** `orch-ref/app/types/memory.go`
**Storage:** `.projects/{slug}/sessions/index.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Sessions | `[]SessionLog` | `sessions` | `sessions` | Yes | All session logs |

### Note

Saved note with YAML frontmatter metadata.

**Source:** `orch-ref/app/types/notes.go`
**Storage:** `.projects/{slug}/notes/{id}.md`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| ID | `string` | `-` (from filename) | `id` | Yes | Note identifier (derived from filename) |
| Title | `string` | `title` | `title` | Yes | Note title |
| Tags | `[]string` | `tags,omitempty` | `tags` | No | Classification tags |
| Pinned | `bool` | `pinned,omitempty` | `pinned` | No | Whether note is pinned |
| StartupPrompt | `bool` | `startup_prompt,omitempty` | `startupPrompt,omitempty` | No | Show as startup prompt |
| QuickAction | `bool` | `quick_action,omitempty` | `quickAction,omitempty` | No | Available as quick action |
| Icon | `string` | `icon,omitempty` | `icon,omitempty` | No | Display icon |
| Color | `string` | `color,omitempty` | `color,omitempty` | No | Display color |
| SourceSessionID | `string` | `source_session_id,omitempty` | `sourceSessionId,omitempty` | No | Originating session ID |
| SourceMessageID | `string` | `source_message_id,omitempty` | `sourceMessageId,omitempty` | No | Originating message ID |
| CreatedAt | `string` | `created_at` | `createdAt` | Yes | ISO 8601 creation timestamp |
| UpdatedAt | `string` | `updated_at` | `updatedAt` | Yes | ISO 8601 last-update timestamp |
| Content | `string` | `-` (body) | `content` | Yes | Markdown body (not in YAML frontmatter) |

### PrdSession

Guided PRD creation session.

**Source:** `orch-ref/app/types/prd.go`
**Storage:** `.projects/{slug}/prd/session.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| ProjectName | `string` | `project_name` | `project_name` | Yes | Project name |
| Slug | `string` | `slug` | `slug` | Yes | Project slug |
| Status | `string` | `status` | `status` | Yes | Session status |
| PrdType | `PrdType` (string) | `prd_type,omitempty` | `prd_type,omitempty` | No | Audience type: "business", "product", "technical", "qa" |
| CurrentIndex | `int` | `current_index` | `current_index` | Yes | Current question index |
| Answers | `[]PrdAnswer` | `answers,omitempty` | `answers,omitempty` | No | Answered questions |
| PendingConditional | `[]PrdQuestion` | `pending_conditional,omitempty` | `pending_conditional,omitempty` | No | Pending follow-up questions |
| ParentSlug | `string` | `parent_slug,omitempty` | `parent_slug,omitempty` | No | Parent project slug (for split PRDs) |
| Phase | `int` | `phase,omitempty` | `phase,omitempty` | No | Current phase index |
| Phases | `[]string` | `phases,omitempty` | `phases,omitempty` | No | Phase names |

### PrdAnswer

One answered PRD question.

**Source:** `orch-ref/app/types/prd.go`
**Storage:** Inline within `PrdSession.Answers`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Question | `string` | `question` | `question` | Yes | The question text |
| Answer | `string` | `answer` | `answer` | Yes | The provided answer |

### PrdQuestion

PRD questionnaire item.

**Source:** `orch-ref/app/types/prd.go`
**Storage:** Inline within `PrdSession.PendingConditional` or question banks

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Index | `int` | `index` | `index` | Yes | Question sequence number |
| Key | `string` | `key` | `key` | Yes | Unique question key |
| Section | `string` | `section` | `section` | Yes | Section grouping |
| Question | `string` | `question` | `question` | Yes | Question text |
| Required | `bool` | `required` | `required` | Yes | Whether answer is required |
| Options | `[]string` | `options,omitempty` | `options,omitempty` | No | Multiple-choice options |

### PrdType Constants

**Source:** `orch-ref/app/types/prd_questions.go`

| Constant | Value | Description |
|----------|-------|-------------|
| PrdTypeBusiness | `"business"` | Business audience PRD |
| PrdTypeProduct | `"product"` | Product audience PRD |
| PrdTypeTechnical | `"technical"` | Technical audience PRD |
| PrdTypeQA | `"qa"` | QA audience PRD |

### ConditionalQ

Follow-up question triggered by answer content.

**Source:** `orch-ref/app/types/prd_questions.go`
**Storage:** Embedded in question banks

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Trigger | `string` | `trigger` | `trigger` | Yes | Keyword/pattern that triggers this question |
| MatchKey | `string` | `match_key` | `match_key` | Yes | Question key whose answer is checked |
| Question | `PrdQuestion` | `question` | `question` | Yes | The conditional follow-up question |

### PrdValidationResult

PRD completeness validation results.

**Source:** `orch-ref/app/types/prd_questions.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Complete | `bool` | `complete` | Yes | Whether PRD is complete |
| Score | `float64` | `score` | Yes | Completeness score (0.0-1.0) |
| SectionScores | `map[string]float64` | `section_scores` | Yes | Per-section completeness scores |
| Gaps | `[]string` | `gaps` | Yes | Missing sections or questions |
| Suggestions | `[]string` | `suggestions` | Yes | Improvement suggestions |

### PrdTemplate

Reusable PRD template with pre-filled answers.

**Source:** `orch-ref/app/types/prd_questions.go`
**Storage:** `.projects/{slug}/prd/templates/{name}.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Name | `string` | `name` | `name` | Yes | Template name |
| Type | `PrdType` | `type` | `type` | Yes | PRD audience type |
| Answers | `[]PrdAnswer` | `answers,omitempty` | `answers,omitempty` | No | Pre-filled answers |
| Created | `string` | `created` | `created` | Yes | Creation timestamp |

### AgentBriefing

Role-specific context summary generated from a PRD.

**Source:** `orch-ref/app/types/prd_questions.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Role | `string` | `role` | Yes | Agent role name |
| Summary | `string` | `summary` | Yes | Role-specific summary |
| KeyPoints | `[]string` | `key_points` | Yes | Key points for this role |
| Constraints | `[]string` | `constraints` | Yes | Constraints for this role |

### HookEvent

Claude Code hook event received by the MCP server.

**Source:** `orch-ref/app/types/hooks.go`
**Storage:** `.projects/{slug}/hooks/events.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| EventType | `string` | `event_type` | `event_type` | Yes | Hook event type |
| SessionID | `string` | `session_id` | `session_id` | Yes | Claude Code session ID |
| ToolName | `string` | `tool_name,omitempty` | `tool_name,omitempty` | No | Tool name (for tool events) |
| AgentType | `string` | `agent_type,omitempty` | `agent_type,omitempty` | No | Agent type |
| Data | `map[string]any` | `data,omitempty` | `data,omitempty` | No | Arbitrary event data |
| Timestamp | `string` | `timestamp` | `timestamp` | Yes | ISO 8601 timestamp |

### HookEventLog

Rolling list of hook events.

**Source:** `orch-ref/app/types/hooks.go`
**Storage:** `.projects/{slug}/hooks/events.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Events | `[]HookEvent` | `events` | `events` | Yes | All stored hook events |

---

## Notification System Types

### Notification

Core data model for all notifications.

**Source:** `orch-ref/app/types/notification.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Unique notification identifier |
| Title | `string` | `title` | Yes | Short headline |
| Message | `string` | `message` | Yes | Detailed notification body |
| Level | `NotificationLevel` | `level` | Yes | Severity: "info", "success", "warning", "error" |
| Category | `NotificationCategory` | `category` | Yes | Domain: "task", "build", "git", "plugin", "system" |
| Source | `string` | `source` | Yes | Plugin ID that generated this notification |
| Icon | `string` | `icon` | Yes | Icon name or URL |
| Actions | `[]NotificationAction` | `actions,omitempty` | No | Clickable action buttons |
| Persistent | `bool` | `persistent` | Yes | Stays until explicitly dismissed |
| Duration | `time.Duration` | `duration` | Yes | Auto-dismiss delay (0 = default 5s) |
| Sound | `bool` | `sound` | Yes | Whether to play a sound |
| Read | `bool` | `read` | Yes | Whether user has seen it |
| CreatedAt | `time.Time` | `created_at` | Yes | Creation timestamp |
| ExpiresAt | `*time.Time` | `expires_at,omitempty` | No | Expiration timestamp |

### NotificationLevel Constants

| Constant | Value | Description |
|----------|-------|-------------|
| LevelInfo | `"info"` | Informational messages |
| LevelSuccess | `"success"` | Successful operation confirmations |
| LevelWarning | `"warning"` | Warnings needing attention |
| LevelError | `"error"` | Errors requiring action |

### NotificationCategory Constants

| Constant | Value | Description |
|----------|-------|-------------|
| CategoryTask | `"task"` | Task lifecycle events |
| CategoryBuild | `"build"` | Build and compilation events |
| CategoryGit | `"git"` | Git and version control events |
| CategoryPlugin | `"plugin"` | Plugin lifecycle events |
| CategorySystem | `"system"` | System-level events |

### NotificationAction

Clickable button attached to a notification.

**Source:** `orch-ref/app/types/notification.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Label | `string` | `label` | Yes | Display text for the action button |
| Action | `string` | `action` | Yes | Action identifier (e.g., "open_file", "dismiss") |
| Payload | `string` | `payload` | Yes | Action payload (e.g., file path, URL) |

### NotificationChannel Constants

**Source:** `orch-ref/app/types/channel.go`

| Constant | Value | Description |
|----------|-------|-------------|
| ChannelDesktop | `"desktop"` | Native desktop notification system |
| ChannelChrome | `"chrome"` | Chrome extension notification system |
| ChannelAll | `"all"` | All enabled channels |

### ChannelConfig

Per-channel notification settings.

**Source:** `orch-ref/app/types/channel.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Channel | `NotificationChannel` | `channel` | Yes | Delivery target |
| Enabled | `bool` | `enabled` | Yes | Whether this channel is active |
| Sound | `bool` | `sound` | Yes | Whether to play sounds on this channel |

### NotificationPreferences

User notification settings.

**Source:** `orch-ref/app/types/preferences.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Channels | `[]ChannelConfig` | `channels` | Yes | Per-channel configuration |
| DoNotDisturb | `DoNotDisturb` | `do_not_disturb` | Yes | Quiet hours configuration |
| GroupBySource | `bool` | `group_by_source` | Yes | Group notifications by source plugin |
| MaxHistory | `int` | `max_history` | Yes | Maximum stored notifications |
| RetentionDays | `int` | `retention_days` | Yes | Days to keep before auto-delete |

### DoNotDisturb

Quiet hours configuration.

**Source:** `orch-ref/app/types/preferences.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Enabled | `bool` | `enabled` | Yes | Whether quiet hours are active |
| Start | `string` | `start` | Yes | Start time in 24h format (e.g., "22:00") |
| End | `string` | `end` | Yes | End time in 24h format (e.g., "08:00") |

### FilterRule

Single filter condition for notifications.

**Source:** `orch-ref/app/types/filter.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Field | `string` | `field` | Yes | Notification field to match ("level", "category", "source") |
| Operator | `string` | `operator` | Yes | Comparison operator: "eq", "neq", "in" |
| Value | `string` | `value` | Yes | Value to compare against |

### NotificationFilter

Rules applied to incoming notifications.

**Source:** `orch-ref/app/types/filter.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Name | `string` | `name` | Yes | Human-readable filter label |
| Rules | `[]FilterRule` | `rules` | Yes | Conditions that must all match |
| Action | `string` | `action` | Yes | Action when matched: "allow", "block", "mute" |
| Enabled | `bool` | `enabled` | Yes | Whether filter is active |

### ListOptions

Pagination and filtering for notification queries.

**Source:** `orch-ref/app/types/contracts.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Limit | `int` | (none) | Yes | Maximum results to return |
| Offset | `int` | (none) | Yes | Number of results to skip |
| Category | `string` | (none) | No | Filter by category (empty = all) |
| Level | `string` | (none) | No | Filter by level (empty = all) |
| Unread | `*bool` | (none) | No | Filter by read status (nil = all) |

---

## MCP Protocol Types

### ToolDefinition

Describes a single MCP tool.

**Source:** `orch-ref/app/types/tools.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Name | `string` | `name` | Yes | Tool name |
| Description | `string` | `description` | Yes | Tool description |
| InputSchema | `InputSchema` | `inputSchema` | Yes | JSON Schema for input |
| Namespace | `string` | `-` (internal) | No | Internal namespace, not serialized |

### InputSchema

JSON Schema for tool input.

**Source:** `orch-ref/app/types/tools.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Type | `string` | `type` | Yes | Schema type (always "object") |
| Properties | `map[string]any` | `properties,omitempty` | No | Property definitions |
| Required | `[]string` | `required,omitempty` | No | Required property names |

### ToolResult

Returned by tool handlers.

**Source:** `orch-ref/app/types/tools.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Content | `[]ContentBlock` | `content` | Yes | Result content blocks |
| IsError | `bool` | `isError,omitempty` | No | Whether result is an error |

### ContentBlock

Single content item in a tool result.

**Source:** `orch-ref/app/types/tools.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Type | `string` | `type` | Yes | Content type (e.g., "text") |
| Text | `string` | `text` | Yes | Content text |

### Tool

Pairs a definition with its handler.

**Source:** `orch-ref/app/types/tools.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Definition | `ToolDefinition` | (none) | Yes | Tool definition |
| Handler | `ToolHandler` | (none) | Yes | Tool handler function |

### ResourceDefinition

Describes a single MCP resource.

**Source:** `orch-ref/app/types/resources.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| URI | `string` | `uri` | Yes | Resource URI |
| Name | `string` | `name` | Yes | Resource name |
| Title | `string` | `title,omitempty` | No | Display title |
| Description | `string` | `description,omitempty` | No | Resource description |
| MimeType | `string` | `mimeType,omitempty` | No | MIME type |

### ResourceContent

Content item returned from resources/read.

**Source:** `orch-ref/app/types/resources.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| URI | `string` | `uri` | Yes | Resource URI |
| MimeType | `string` | `mimeType,omitempty` | No | MIME type |
| Text | `string` | `text,omitempty` | No | Text content |
| Blob | `string` | `blob,omitempty` | No | Base64-encoded binary content |

### PromptArgument

Parameter accepted by a prompt.

**Source:** `orch-ref/app/types/prompts.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Name | `string` | `name` | Yes | Argument name |
| Description | `string` | `description,omitempty` | No | Argument description |
| Required | `bool` | `required,omitempty` | No | Whether argument is required |

### PromptDefinition

Describes a single MCP prompt.

**Source:** `orch-ref/app/types/prompts.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Name | `string` | `name` | Yes | Prompt name |
| Title | `string` | `title,omitempty` | No | Display title |
| Description | `string` | `description,omitempty` | No | Prompt description |
| Arguments | `[]PromptArgument` | `arguments,omitempty` | No | Accepted arguments |

### PromptMessage

Single message in a prompt response.

**Source:** `orch-ref/app/types/prompts.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Role | `string` | `role` | Yes | Message role ("user" or "assistant") |
| Content | `ContentBlock` | `content` | Yes | Message content |

### JSON-RPC Protocol Types

**Source:** `orch-ref/app/types/protocol.go`

#### JSONRPCRequest

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| JSONRPC | `string` | `jsonrpc` | Yes | Protocol version (always "2.0") |
| ID | `any` | `id,omitempty` | No | Request identifier |
| Method | `string` | `method` | Yes | Method name |
| Params | `json.RawMessage` | `params,omitempty` | No | Method parameters |

#### JSONRPCResponse

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| JSONRPC | `string` | `jsonrpc` | Yes | Protocol version (always "2.0") |
| ID | `any` | `id,omitempty` | No | Request identifier |
| Result | `any` | `result,omitempty` | No | Success result |
| Error | `*JSONRPCError` | `error,omitempty` | No | Error information |

#### JSONRPCError

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Code | `int` | `code` | Yes | Error code |
| Message | `string` | `message` | Yes | Error message |

#### CallToolParams

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Name | `string` | `name` | Yes | Tool name to call |
| Arguments | `map[string]any` | `arguments,omitempty` | No | Tool arguments |

#### InitializeResult

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ProtocolVersion | `string` | `protocolVersion` | Yes | MCP protocol version |
| Capabilities | `ServerCaps` | `capabilities` | Yes | Server capabilities |
| ServerInfo | `ServerInfo` | `serverInfo` | Yes | Server identification |

#### ServerCaps

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Tools | `*ToolsCap` | `tools,omitempty` | No | Tool support capabilities |
| Resources | `*ResourcesCap` | `resources,omitempty` | No | Resource support capabilities |
| Prompts | `*PromptsCap` | `prompts,omitempty` | No | Prompt support capabilities |

#### ToolsCap

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ListChanged | `bool` | `listChanged,omitempty` | No | Whether tool list can change |

#### ResourcesCap

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Subscribe | `bool` | `subscribe,omitempty` | No | Whether subscriptions are supported |
| ListChanged | `bool` | `listChanged,omitempty` | No | Whether resource list can change |

#### PromptsCap

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ListChanged | `bool` | `listChanged,omitempty` | No | Whether prompt list can change |

#### ServerInfo

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Name | `string` | `name` | Yes | Server name |
| Version | `string` | `version` | Yes | Server version |

#### ListToolsResult

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Tools | `[]ToolDefinition` | `tools` | Yes | Available tools |

#### ReadResourceParams

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| URI | `string` | `uri` | Yes | Resource URI to read |

#### ListResourcesResult

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Resources | `[]ResourceDefinition` | `resources` | Yes | Available resources |

#### ReadResourceResult

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Contents | `[]ResourceContent` | `contents` | Yes | Resource contents |

#### GetPromptParams

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Name | `string` | `name` | Yes | Prompt name |
| Arguments | `map[string]string` | `arguments,omitempty` | No | Prompt arguments |

#### ListPromptsResult

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Prompts | `[]PromptDefinition` | `prompts` | Yes | Available prompts |

#### GetPromptResult

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Description | `string` | `description,omitempty` | No | Prompt description |
| Messages | `[]PromptMessage` | `messages` | Yes | Prompt messages |

---

## GitHub Integration Types

### GitHubDeviceCodeResponse

Returned by the device flow initiation.

**Source:** `orch-ref/app/types/github.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| DeviceCode | `string` | `device_code` | Yes | Device verification code |
| UserCode | `string` | `user_code` | Yes | User-facing code to enter |
| VerificationURI | `string` | `verification_uri` | Yes | URL for user to visit |
| ExpiresIn | `int` | `expires_in` | Yes | Seconds until code expires |
| Interval | `int` | `interval` | Yes | Polling interval in seconds |

### GitHubTokenResponse

Returned after successful device flow authorization.

**Source:** `orch-ref/app/types/github.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| AccessToken | `string` | `access_token` | Yes | OAuth access token |
| TokenType | `string` | `token_type` | Yes | Token type (e.g., "bearer") |
| Scope | `string` | `scope` | Yes | Granted scopes |

### GitHubUser

Authenticated GitHub user.

**Source:** `orch-ref/app/types/github.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `int` | `id` | Yes | GitHub user ID |
| Login | `string` | `login` | Yes | Username |
| Name | `string` | `name` | Yes | Display name |
| Email | `string` | `email` | Yes | Email address |
| AvatarURL | `string` | `avatar_url` | Yes | Avatar image URL |

### GitHubAuthState

Current authentication state.

**Source:** `orch-ref/app/types/github.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Status | `string` | `status` | Yes | "connected", "disconnected", or "pending" |
| User | `*GitHubUser` | `user,omitempty` | No | Authenticated user (if connected) |
| Scopes | `string` | `scopes,omitempty` | No | Granted scopes |
| ExpiresAt | `string` | `expires_at,omitempty` | No | Token expiration |

### GitHubErrorResponse

GitHub API error response.

**Source:** `orch-ref/app/types/github.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Error | `string` | `error` | Yes | Error code |
| ErrorDescription | `string` | `error_description` | Yes | Error description |

### GitHubConfig

GitHub OAuth application settings.

**Source:** `orch-ref/app/types/github_config.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ClientID | `string` | `client_id` | Yes | OAuth client ID |
| ClientSecret | `string` | `client_secret` | Yes | OAuth client secret |
| RedirectURI | `string` | `redirect_uri` | Yes | OAuth redirect URI |
| Scopes | `[]string` | `scopes` | Yes | Requested OAuth scopes |

### TrackedRepo

GitHub repository being tracked for activity.

**Source:** `orch-ref/app/types/github_config.go`

| Field | Type | JSON/YAML Tag | Required | Description |
|-------|------|---------------|----------|-------------|
| ID | `string` | `id` | Yes | Tracking record ID |
| Owner | `string` | `owner` | Yes | Repository owner |
| Repo | `string` | `repo` | Yes | Repository name |
| DefaultBranch | `string` | `default_branch` | Yes | Default branch name |
| WatchPRs | `bool` | `watch_prs` | Yes | Watch pull requests |
| WatchIssues | `bool` | `watch_issues` | Yes | Watch issues |
| AddedAt | `string` | `added_at` | Yes | When tracking was added |

### TrackedReposConfig

List of tracked repos wrapper.

**Source:** `orch-ref/app/types/github_config.go`

| Field | Type | JSON/YAML Tag | Required | Description |
|-------|------|---------------|----------|-------------|
| Repos | `[]TrackedRepo` | `repos` | Yes | Tracked repositories |

### GitHubUserRef

Lightweight user reference in PR/issue payloads.

**Source:** `orch-ref/app/types/github_pr.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Login | `string` | `login` | Yes | Username |
| AvatarURL | `string` | `avatar_url` | Yes | Avatar URL |

### GitHubRepoRef

Lightweight repo reference.

**Source:** `orch-ref/app/types/github_pr.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| FullName | `string` | `full_name` | Yes | Full repo name (owner/repo) |
| HTMLURL | `string` | `html_url` | Yes | Repository URL |

### GitHubBranchRef

PR branch reference (head or base).

**Source:** `orch-ref/app/types/github_pr.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Ref | `string` | `ref` | Yes | Branch name |
| SHA | `string` | `sha` | Yes | Commit SHA |
| Repo | `GitHubRepoRef` | `repo` | Yes | Repository reference |

### GitHubLabel

PR/issue label.

**Source:** `orch-ref/app/types/github_pr.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Name | `string` | `name` | Yes | Label name |
| Color | `string` | `color` | Yes | Hex color code |

### GitHubPR

Pull request.

**Source:** `orch-ref/app/types/github_pr.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Number | `int` | `number` | Yes | PR number |
| Title | `string` | `title` | Yes | PR title |
| Body | `string` | `body` | Yes | PR description body |
| State | `string` | `state` | Yes | PR state |
| Draft | `bool` | `draft` | Yes | Whether PR is a draft |
| User | `GitHubUserRef` | `user` | Yes | PR author |
| Head | `GitHubBranchRef` | `head` | Yes | Source branch |
| Base | `GitHubBranchRef` | `base` | Yes | Target branch |
| Labels | `[]GitHubLabel` | `labels` | Yes | Applied labels |
| Reviewers | `[]GitHubUserRef` | `requested_reviewers` | Yes | Requested reviewers |
| CreatedAt | `string` | `created_at` | Yes | Creation timestamp |
| UpdatedAt | `string` | `updated_at` | Yes | Last update timestamp |
| MergedAt | `string` | `merged_at,omitempty` | No | Merge timestamp |
| HTMLURL | `string` | `html_url` | Yes | Web URL |

### GitHubPRFile

File changed in a PR.

**Source:** `orch-ref/app/types/github_pr.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Filename | `string` | `filename` | Yes | File path |
| Status | `string` | `status` | Yes | Change status |
| Additions | `int` | `additions` | Yes | Lines added |
| Deletions | `int` | `deletions` | Yes | Lines deleted |
| Changes | `int` | `changes` | Yes | Total changes |
| Patch | `string` | `patch,omitempty` | No | Unified diff patch |

### GitHubPRReview

Review on a PR.

**Source:** `orch-ref/app/types/github_pr.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `int` | `id` | Yes | Review ID |
| User | `GitHubUserRef` | `user` | Yes | Reviewer |
| State | `string` | `state` | Yes | "APPROVED", "CHANGES_REQUESTED", "COMMENTED" |
| Body | `string` | `body` | Yes | Review body |
| SubmittedAt | `string` | `submitted_at` | Yes | Submission timestamp |

### GitHubPRComment

Comment on a PR.

**Source:** `orch-ref/app/types/github_pr.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `int` | `id` | Yes | Comment ID |
| User | `GitHubUserRef` | `user` | Yes | Comment author |
| Body | `string` | `body` | Yes | Comment body |
| Path | `string` | `path,omitempty` | No | File path (for inline comments) |
| Line | `int` | `line,omitempty` | No | Line number (for inline comments) |
| CreatedAt | `string` | `created_at` | Yes | Creation timestamp |

### PRListFilter

PR listing filters.

**Source:** `orch-ref/app/types/github_pr.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| State | `string` | `state,omitempty` | No | "open", "closed", "all" |
| Sort | `string` | `sort,omitempty` | No | "created", "updated", "popularity" |
| Direction | `string` | `direction,omitempty` | No | "asc", "desc" |
| Head | `string` | `head,omitempty` | No | Filter by head branch |
| Base | `string` | `base,omitempty` | No | Filter by base branch |

### PRCreateParams

Parameters for creating a PR.

**Source:** `orch-ref/app/types/github_pr.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Title | `string` | `title` | Yes | PR title |
| Body | `string` | `body,omitempty` | No | PR description |
| Head | `string` | `head` | Yes | Source branch |
| Base | `string` | `base` | Yes | Target branch |
| Draft | `bool` | `draft,omitempty` | No | Create as draft |
| Reviewers | `[]string` | `reviewers,omitempty` | No | Reviewer logins |

### GitHubIssue

GitHub issue.

**Source:** `orch-ref/app/types/github_issue.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Number | `int` | `number` | Yes | Issue number |
| Title | `string` | `title` | Yes | Issue title |
| Body | `string` | `body` | Yes | Issue body |
| State | `string` | `state` | Yes | Issue state |
| User | `GitHubUserRef` | `user` | Yes | Issue author |
| Labels | `[]GitHubLabel` | `labels` | Yes | Applied labels |
| Assignees | `[]GitHubUserRef` | `assignees` | Yes | Assigned users |
| CreatedAt | `string` | `created_at` | Yes | Creation timestamp |
| UpdatedAt | `string` | `updated_at` | Yes | Last update timestamp |
| HTMLURL | `string` | `html_url` | Yes | Web URL |

### GitHubIssueComment

Comment on an issue.

**Source:** `orch-ref/app/types/github_issue.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `int` | `id` | Yes | Comment ID |
| User | `GitHubUserRef` | `user` | Yes | Comment author |
| Body | `string` | `body` | Yes | Comment body |
| CreatedAt | `string` | `created_at` | Yes | Creation timestamp |

### IssueListFilter

Issue listing filters.

**Source:** `orch-ref/app/types/github_issue.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| State | `string` | `state,omitempty` | No | "open", "closed", "all" |
| Assignee | `string` | `assignee,omitempty` | No | Login, "none", or "*" |
| Labels | `string` | `labels,omitempty` | No | Comma-separated label names |
| Sort | `string` | `sort,omitempty` | No | "created", "updated", "comments" |

### IssueCreateParams

Parameters for creating an issue.

**Source:** `orch-ref/app/types/github_issue.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Title | `string` | `title` | Yes | Issue title |
| Body | `string` | `body,omitempty` | No | Issue body |
| Assignees | `[]string` | `assignees,omitempty` | No | Assignee logins |
| Labels | `[]string` | `labels,omitempty` | No | Label names |

### GitHubCheckRun

CI check run.

**Source:** `orch-ref/app/types/github_issue.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `int` | `id` | Yes | Check run ID |
| Name | `string` | `name` | Yes | Check run name |
| Status | `string` | `status` | Yes | "queued", "in_progress", "completed" |
| Conclusion | `string` | `conclusion` | Yes | "success", "failure", "neutral", "cancelled", "skipped", "timed_out" |
| HTMLURL | `string` | `html_url` | Yes | Web URL |

### GitHubCIStatus

Aggregated CI/CD status for a commit ref.

**Source:** `orch-ref/app/types/github_issue.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Overall | `string` | `overall` | Yes | "success", "failure", "pending" |
| CheckRuns | `[]GitHubCheckRun` | `check_runs` | Yes | Individual check runs |
| Total | `int` | `total` | Yes | Total number of checks |
| Passed | `int` | `passed` | Yes | Number passed |
| Failed | `int` | `failed` | Yes | Number failed |
| Pending | `int` | `pending` | Yes | Number pending |

---

## Linear Integration Types

**Source:** `orch-ref/app/types/linear.go`

### LinearConfig

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ClientID | `string` | `client_id` | Yes | OAuth client ID |
| ClientSecret | `string` | `client_secret` | Yes | OAuth client secret |
| RedirectURI | `string` | `redirect_uri` | Yes | OAuth redirect URI |

### LinearUser

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Linear user ID |
| Name | `string` | `name` | Yes | Full name |
| DisplayName | `string` | `displayName` | Yes | Display name |
| Email | `string` | `email` | Yes | Email address |
| AvatarURL | `string` | `avatarUrl,omitempty` | No | Avatar URL |
| Active | `bool` | `active` | Yes | Whether user is active |

### LinearAuthState

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Status | `string` | `status` | Yes | "connected" or "disconnected" |
| User | `*LinearUser` | `user,omitempty` | No | Authenticated user |

### LinearIssue

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Linear issue ID |
| Identifier | `string` | `identifier` | Yes | Human-readable ID (e.g., "ENG-123") |
| Title | `string` | `title` | Yes | Issue title |
| Description | `string` | `description,omitempty` | No | Issue description |
| Priority | `int` | `priority` | Yes | 0=none, 1=urgent, 2=high, 3=medium, 4=low |
| State | `*LinearState` | `state,omitempty` | No | Workflow state |
| Assignee | `*LinearUser` | `assignee,omitempty` | No | Assigned user |
| Team | `*LinearTeam` | `team,omitempty` | No | Team |
| Project | `*LinearProject` | `project,omitempty` | No | Project |
| Labels | `[]LinearLabel` | `labels,omitempty` | No | Applied labels |
| CreatedAt | `string` | `createdAt` | Yes | Creation timestamp |
| UpdatedAt | `string` | `updatedAt` | Yes | Last update timestamp |
| Parent | `*LinearIssueRef` | `parent,omitempty` | No | Parent issue reference |

### LinearIssueRef

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Issue ID |
| Identifier | `string` | `identifier` | Yes | Human-readable ID |

### LinearState

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | State ID |
| Name | `string` | `name` | Yes | State name |
| Type | `string` | `type` | Yes | "backlog", "unstarted", "started", "completed", "cancelled" |
| Color | `string` | `color` | Yes | Hex color |

### LinearTeam

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Team ID |
| Name | `string` | `name` | Yes | Team name |
| Key | `string` | `key` | Yes | Team key prefix |

### LinearProject

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Project ID |
| Name | `string` | `name` | Yes | Project name |
| State | `string` | `state` | Yes | Project state |
| SlugID | `string` | `slugId` | Yes | URL-safe slug ID |

### LinearLabel

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Label ID |
| Name | `string` | `name` | Yes | Label name |
| Color | `string` | `color` | Yes | Hex color |

### LinearCycle

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Cycle ID |
| Number | `int` | `number` | Yes | Cycle number |
| Name | `string` | `name,omitempty` | No | Cycle name |
| StartsAt | `string` | `startsAt` | Yes | Start date |
| EndsAt | `string` | `endsAt` | Yes | End date |
| Progress | `float64` | `progress` | Yes | Completion progress (0.0-1.0) |

### LinearIssueFilter

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| State | `string` | `state,omitempty` | No | Filter by state |
| Assignee | `string` | `assignee,omitempty` | No | Filter by assignee |
| Priority | `int` | `priority,omitempty` | No | Filter by priority |
| Labels | `string` | `labels,omitempty` | No | Filter by labels |

### LinearIssueCreateParams

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Title | `string` | `title` | Yes | Issue title |
| Description | `string` | `description,omitempty` | No | Issue description |
| TeamID | `string` | `teamId` | Yes | Target team ID |
| AssigneeID | `string` | `assigneeId,omitempty` | No | Assignee ID |
| Priority | `int` | `priority,omitempty` | No | Priority level |
| StateID | `string` | `stateId,omitempty` | No | Initial state ID |
| ProjectID | `string` | `projectId,omitempty` | No | Project ID |
| ParentID | `string` | `parentId,omitempty` | No | Parent issue ID |

### LinearIssueUpdateParams

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Title | `*string` | `title,omitempty` | No | New title |
| Description | `*string` | `description,omitempty` | No | New description |
| AssigneeID | `*string` | `assigneeId,omitempty` | No | New assignee ID |
| Priority | `*int` | `priority,omitempty` | No | New priority |
| StateID | `*string` | `stateId,omitempty` | No | New state ID |
| ProjectID | `*string` | `projectId,omitempty` | No | New project ID |

---

## Jira Integration Types

**Source:** `orch-ref/app/types/jira.go`

### JiraConfig

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| BaseURL | `string` | `base_url` | Yes | Jira instance URL (e.g., "https://mycompany.atlassian.net") |
| Email | `string` | `email` | Yes | Email for Basic auth |

### JiraUser

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| AccountID | `string` | `accountId` | Yes | Jira account ID |
| DisplayName | `string` | `displayName` | Yes | Display name |
| Email | `string` | `emailAddress` | Yes | Email address |
| AvatarURL | `string` | `avatarUrls,omitempty` | No | Avatar URL |
| Active | `bool` | `active` | Yes | Whether user is active |

### JiraAuthState

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Status | `string` | `status` | Yes | "connected" or "disconnected" |
| User | `*JiraUser` | `user,omitempty` | No | Authenticated user |

### JiraIssue

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Jira issue ID |
| Key | `string` | `key` | Yes | Issue key (e.g., "PROJ-123") |
| Self | `string` | `self` | Yes | API self-link URL |
| Fields | `JiraIssueFields` | `fields` | Yes | Issue fields |

### JiraIssueFields

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Summary | `string` | `summary` | Yes | Issue summary |
| Description | `any` | `description,omitempty` | No | Description in ADF format |
| Status | `*JiraStatus` | `status,omitempty` | No | Current status |
| Priority | `*JiraPriority` | `priority,omitempty` | No | Priority |
| IssueType | `*JiraIssueType` | `issuetype,omitempty` | No | Issue type |
| Assignee | `*JiraUser` | `assignee,omitempty` | No | Assigned user |
| Reporter | `*JiraUser` | `reporter,omitempty` | No | Reporter |
| Project | `*JiraProject` | `project,omitempty` | No | Project |
| Labels | `[]string` | `labels,omitempty` | No | Labels |
| Created | `string` | `created,omitempty` | No | Creation timestamp |
| Updated | `string` | `updated,omitempty` | No | Last update timestamp |
| Parent | `*JiraIssueRef` | `parent,omitempty` | No | Parent issue reference |

### JiraStatus

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Status ID |
| Name | `string` | `name` | Yes | Status name |

### JiraPriority

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Priority ID |
| Name | `string` | `name` | Yes | Priority name |

### JiraIssueType

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Issue type ID |
| Name | `string` | `name` | Yes | Issue type name (Bug, Story, Task, etc.) |

### JiraProject

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Project ID |
| Key | `string` | `key` | Yes | Project key |
| Name | `string` | `name` | Yes | Project name |

### JiraIssueRef

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Issue ID |
| Key | `string` | `key` | Yes | Issue key |

### JiraSearchResult

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Total | `int` | `total` | Yes | Total matching results |
| MaxResults | `int` | `maxResults` | Yes | Page size |
| StartAt | `int` | `startAt` | Yes | Page offset |
| Issues | `[]JiraIssue` | `issues` | Yes | Matching issues |

### JiraTransition

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id` | Yes | Transition ID |
| Name | `string` | `name` | Yes | Transition name |
| To | `*JiraStatus` | `to,omitempty` | No | Target status |

### JiraTransitionsResponse

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Transitions | `[]JiraTransition` | `transitions` | Yes | Available transitions |

### JiraBoard

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `int` | `id` | Yes | Board ID |
| Name | `string` | `name` | Yes | Board name |
| Type | `string` | `type` | Yes | Board type (Scrum/Kanban) |

### JiraSprint

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `int` | `id` | Yes | Sprint ID |
| Name | `string` | `name` | Yes | Sprint name |
| State | `string` | `state` | Yes | Sprint state |
| StartDate | `string` | `startDate,omitempty` | No | Start date |
| EndDate | `string` | `endDate,omitempty` | No | End date |
| Goal | `string` | `goal,omitempty` | No | Sprint goal |

### JiraIssueCreateParams

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Fields | `JiraIssueCreateFields` | `fields` | Yes | Issue fields for creation |

### JiraIssueCreateFields

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Project | `JiraProjectRef` | `project` | Yes | Target project |
| Summary | `string` | `summary` | Yes | Issue summary |
| IssueType | `JiraTypeRef` | `issuetype` | Yes | Issue type |
| Description | `any` | `description,omitempty` | No | ADF-format description |
| Assignee | `*JiraRef` | `assignee,omitempty` | No | Assignee reference |
| Priority | `*JiraRef` | `priority,omitempty` | No | Priority reference |
| Labels | `[]string` | `labels,omitempty` | No | Labels |

### JiraProjectRef

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Key | `string` | `key` | Yes | Project key |

### JiraTypeRef

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Name | `string` | `name` | Yes | Issue type name |

### JiraRef

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `id,omitempty` | No | Entity ID |

### JiraIssueUpdateParams

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Fields | `map[string]any` | `fields` | Yes | Fields to update (arbitrary key-value map) |

---

## Usage & Hook Tracking Types

**Source:** `orch-ref/app/types/usage.go`

### UsageData

Token usage across sessions.

**Storage:** `.projects/{slug}/usage.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Sessions | `[]UsageSession` | `sessions` | `sessions` | Yes | Individual usage sessions |
| Totals | `UsageTotals` | `totals` | `totals` | Yes | Aggregated totals |

### UsageSession

Single usage tracking session.

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Provider | `string` | `provider` | `provider` | Yes | AI provider name |
| Model | `string` | `model` | `model` | Yes | Model identifier |
| StartedAt | `string` | `started_at` | `started_at` | Yes | Session start time |
| EndedAt | `string` | `ended_at,omitempty` | `ended_at,omitempty` | No | Session end time |
| Requests | `[]RequestEntry` | `requests,omitempty` | `requests,omitempty` | No | Individual request records |
| TotalInput | `int` | `total_input` | `total_input` | Yes | Total input tokens |
| TotalOutput | `int` | `total_output` | `total_output` | Yes | Total output tokens |
| TotalCost | `float64` | `total_cost` | `total_cost` | Yes | Total cost in USD |

### UsageTotals

Aggregated usage across all sessions.

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| TotalInput | `int` | `total_input` | `total_input` | Yes | Total input tokens |
| TotalOutput | `int` | `total_output` | `total_output` | Yes | Total output tokens |
| TotalCost | `float64` | `total_cost` | `total_cost` | Yes | Total cost in USD |

### RequestEntry

Single API request's usage.

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Timestamp | `string` | `timestamp` | `timestamp` | Yes | Request timestamp |
| InputTokens | `int` | `input_tokens` | `input_tokens` | Yes | Input tokens consumed |
| OutputTokens | `int` | `output_tokens` | `output_tokens` | Yes | Output tokens consumed |
| Cost | `float64` | `cost` | `cost` | Yes | Request cost in USD |

### RequestLog

Feature requests and suggestions.

**Storage:** `.projects/{slug}/requests.yaml`

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| Project | `string` | `project` | `project` | Yes | Project slug |
| Requests | `[]RequestLogItem` | `requests` | `requests` | Yes | Logged requests |

### RequestLogItem

Single logged request.

| Field | Type | YAML Tag | JSON Tag | Required | Description |
|-------|------|----------|----------|----------|-------------|
| ID | `string` | `id` | `id` | Yes | Request ID |
| Type | `string` | `type` | `type` | Yes | Request type |
| Date | `string` | `date` | `date` | Yes | Request date |
| Description | `string` | `description` | `description` | Yes | Request description |
| Status | `string` | `status` | `status` | Yes | Request status |

---

## PostgreSQL Entities (Cloud)

Source of truth. Uses UUID primary keys, `TIMESTAMPTZ` timestamps, `JSONB` for flexible data, `text[]` for arrays, and partitioned tables for the sync log.

**Migration sources:**
- `orch-ref/database/migrations/20260225010000_create_sync_foundation.sql`
- `orch-ref/database/migrations/20260225001900_create_component_library.sql`
- `orch-ref/database/migrations/003_conflict_log.sql`
- `orch-ref/app/models/sync.go` (GORM models)
- `orch-ref/app/models/component.go` (GORM models)
- `orch-ref/app/models/conflict_log.go` (GORM models)

### SyncModel (Base Struct)

Embedded by all syncable entities. Provides UUID PK, optimistic concurrency version, timestamps, and soft delete.

**Source:** `orch-ref/app/models/sync.go`

| Field | Type | GORM Tag | JSON Tag | Description |
|-------|------|----------|----------|-------------|
| ID | `string` | `type:uuid;primaryKey` | `id` | UUID primary key (auto-generated) |
| Version | `int` | `not null;default:1` | `version` | Optimistic concurrency version |
| CreatedAt | `time.Time` | `type:timestamptz;not null` | `created_at` | Creation timestamp |
| UpdatedAt | `time.Time` | `type:timestamptz;not null` | `updated_at` | Last update timestamp |
| DeletedAt | `gorm.DeletedAt` | `index` | `deleted_at,omitempty` | Soft delete timestamp |

### users

**Table:** `users`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| name | `VARCHAR(255)` | NOT NULL | | User display name |
| email | `VARCHAR(255)` | NOT NULL | | Email (unique) |
| password_hash | `TEXT` | NULL | | Bcrypt hash (NULL for OAuth-only) |
| avatar_url | `TEXT` | NULL | | Profile image URL |
| plan | `VARCHAR(50)` | NOT NULL | `'free'` | Subscription plan |
| version | `BIGINT` | NOT NULL | `1` | Optimistic concurrency version |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Creation timestamp |
| updated_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Last update |
| deleted_at | `TIMESTAMPTZ` | NULL | | Soft delete |

**Indexes:** `idx_users_email` (email), `idx_users_deleted` (deleted_at WHERE NOT NULL)
**Relationships:** Has many `devices`, `oauth_accounts`

### oauth_accounts

**Table:** `oauth_accounts`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| user_id | `UUID` | NOT NULL | | FK to users(id) CASCADE |
| provider | `VARCHAR(50)` | NOT NULL | | Provider name (github, google, etc.) |
| provider_id | `VARCHAR(255)` | NOT NULL | | Provider-specific user ID |
| access_token | `TEXT` | NULL | | OAuth access token |
| refresh_token | `TEXT` | NULL | | OAuth refresh token |
| expires_at | `TIMESTAMPTZ` | NULL | | Token expiration |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Creation timestamp |
| updated_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Last update |

**Indexes:** `idx_oauth_user_id` (user_id)
**Unique:** `(provider, provider_id)`

### refresh_tokens

**Table:** `refresh_tokens`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| user_id | `UUID` | NOT NULL | | FK to users(id) CASCADE |
| device_id | `UUID` | NULL | | Associated device |
| token_hash | `TEXT` | NOT NULL | | SHA-256 hash of token (unique) |
| expires_at | `TIMESTAMPTZ` | NOT NULL | | Token expiration |
| revoked_at | `TIMESTAMPTZ` | NULL | | When revoked (NULL = active) |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Creation timestamp |

**Indexes:** `idx_refresh_tokens_user_id`, `idx_refresh_tokens_token_hash`, `idx_refresh_tokens_expires`
**Retention:** Revoked tokens kept 30 days then pruned

### devices

**Table:** `devices`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| user_id | `UUID` | NOT NULL | | FK to users(id) CASCADE |
| name | `VARCHAR(255)` | NOT NULL | | Device name |
| platform | `VARCHAR(50)` | NOT NULL | | Platform (desktop, mobile, web) |
| fingerprint | `TEXT` | NULL | | Device fingerprint |
| sync_token | `TEXT` | NULL | | Incremental pull cursor |
| last_seen | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Last authenticated request |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Creation timestamp |
| updated_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Last update |

**Indexes:** `idx_devices_user_id`

### settings

**Table:** `settings`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| user_id | `UUID` | NOT NULL | | Composite PK part 1; FK to users(id) CASCADE |
| key | `VARCHAR(255)` | NOT NULL | | Composite PK part 2; setting key |
| value | `TEXT` | NULL | | Setting value (encrypted for sensitive data) |
| version | `BIGINT` | NOT NULL | `1` | Optimistic concurrency version |
| updated_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Last update |

**Primary Key:** `(user_id, key)`
**Indexes:** `idx_settings_user_id`, `idx_settings_updated`

### notes

**Table:** `notes`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| user_id | `UUID` | NOT NULL | | FK to users(id) CASCADE |
| title | `VARCHAR(500)` | NOT NULL | `''` | Note title |
| content | `TEXT` | NOT NULL | `''` | Note body |
| tags | `TEXT[]` | NULL | | Array of tags |
| pinned | `BOOLEAN` | NOT NULL | `FALSE` | Pin state |
| version | `BIGINT` | NOT NULL | `1` | Optimistic concurrency version |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Creation timestamp |
| updated_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Last update |
| deleted_at | `TIMESTAMPTZ` | NULL | | Soft delete |

**Indexes:** `idx_notes_user_id`, `idx_notes_updated`, `idx_notes_tags` (GIN), `idx_notes_deleted`, `idx_notes_fts` (GIN full-text on title+content)
**Retention:** Soft-deleted notes pruned after 30 days

### projects

**Table:** `projects`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| user_id | `UUID` | NOT NULL | | FK to users(id) CASCADE |
| slug | `VARCHAR(255)` | NOT NULL | | URL-safe identifier |
| name | `VARCHAR(255)` | NOT NULL | | Project name |
| description | `TEXT` | NULL | | Description |
| data | `JSONB` | NOT NULL | `'{}'` | Full epics/stories/tasks tree |
| version | `BIGINT` | NOT NULL | `1` | Optimistic concurrency version |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Creation timestamp |
| updated_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Last update |
| deleted_at | `TIMESTAMPTZ` | NULL | | Soft delete |

**Unique:** `(user_id, slug)`
**Indexes:** `idx_projects_user_id`, `idx_projects_updated`, `idx_projects_data` (GIN), `idx_projects_deleted`

### sync_log

Append-only change log. **Partitioned monthly** by `created_at`.

**Table:** `sync_log` (PARTITION BY RANGE created_at)

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Record ID |
| user_id | `UUID` | NOT NULL | | User who made the change |
| device_id | `UUID` | NOT NULL | | Device that originated the change |
| entity_type | `VARCHAR(50)` | NOT NULL | | "setting", "note", "project" |
| entity_id | `TEXT` | NOT NULL | | Entity identifier |
| action | `VARCHAR(20)` | NOT NULL | `'upsert'` | "upsert" or "delete" |
| payload | `JSONB` | NOT NULL | `'{}'` | Serialized entity data |
| version | `BIGINT` | NOT NULL | | Entity version at time of change |
| idempotency_key | `TEXT` | NULL | | Deduplication key |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | When change was recorded |

**Primary Key:** `(id, created_at)` (composite for partitioning)
**Indexes:** `idx_sync_log_user_id` (user_id, created_at DESC), `idx_sync_log_entity`, `idx_sync_log_version`, `idx_sync_log_idempotency`
**Partitions:** Monthly (e.g., `sync_log_2026_02`, `sync_log_2026_03`)

### conflict_log

Records losing side of every last-write-wins conflict.

**Table:** `conflict_log`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| user_id | `UUID` | NOT NULL | | User who owns the entity |
| entity_type | `VARCHAR(50)` | NOT NULL | | Entity type |
| entity_id | `TEXT` | NOT NULL | | Entity identifier |
| winner_device | `UUID` | NULL | | Device that won |
| loser_device | `UUID` | NULL | | Device that lost |
| winner_version | `BIGINT` | NOT NULL | | Winning version number |
| loser_version | `BIGINT` | NOT NULL | | Losing version number |
| loser_payload | `JSONB` | NOT NULL | `'{}'` | Losing entity data |
| device_id | `VARCHAR(255)` | NULL | | Added in migration 003; opaque device token |
| resolved_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | When conflict was resolved |

**Indexes:** `idx_conflict_log_user_id`, `idx_conflict_log_entity`, `idx_conflict_log_user_entity`, `idx_conflict_log_created_at`
**Retention:** 90 days then pruned

### components

**Table:** `components` (embeds SyncModel)

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| user_id | `UUID` | NOT NULL | | Component author |
| name | `VARCHAR(255)` | NOT NULL | | Component name |
| description | `TEXT` | NULL | | Description |
| framework | `VARCHAR(50)` | NOT NULL | | "html", "react", "vue", "svelte", "angular" |
| html | `TEXT` | NULL | | HTML source |
| css | `TEXT` | NULL | | CSS source |
| js | `TEXT` | NULL | | JavaScript source |
| jsx | `TEXT` | NULL | | JSX/component source |
| tags | `TEXT[]` | NULL | | Classification tags |
| props_schema | `JSONB` | NULL | | JSON Schema for component props |
| dependencies | `JSONB` | NULL | | External package dependencies |
| thumbnail_url | `TEXT` | NULL | | Preview image URL |
| forked_from_id | `UUID` | NULL | | Self-referential FK; origin component |
| version | `INTEGER` | NOT NULL | `1` | Optimistic concurrency version |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Creation timestamp |
| updated_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Last update |
| deleted_at | `TIMESTAMPTZ` | NULL | | Soft delete |

**Indexes:** `idx_components_user_id`, `idx_components_framework`, `idx_components_tags` (GIN), `idx_components_forked_from_id`, `idx_components_deleted_at`
**Relationships:** Has many `component_versions`, `component_exports`; many-to-many `collections` via `collection_components`

### component_versions

**Table:** `component_versions`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| component_id | `UUID` | NOT NULL | | FK to components(id) CASCADE |
| version | `INTEGER` | NOT NULL | | Version snapshot number |
| html | `TEXT` | NULL | | HTML at this version |
| css | `TEXT` | NULL | | CSS at this version |
| js | `TEXT` | NULL | | JS at this version |
| jsx | `TEXT` | NULL | | JSX at this version |
| snapshot_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | When snapshot was taken |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Row creation timestamp |

**Indexes:** `idx_component_versions_component_id`

### collections

**Table:** `collections` (embeds SyncModel)

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| user_id | `UUID` | NOT NULL | | Collection owner |
| name | `VARCHAR(255)` | NOT NULL | | Collection name |
| description | `TEXT` | NULL | | Description |
| version | `INTEGER` | NOT NULL | `1` | Optimistic concurrency version |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Creation timestamp |
| updated_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Last update |
| deleted_at | `TIMESTAMPTZ` | NULL | | Soft delete |

**Indexes:** `idx_collections_user_id`, `idx_collections_deleted_at`

### collection_components (pivot)

**Table:** `collection_components`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| collection_id | `UUID` | NOT NULL | | FK to collections(id) CASCADE |
| component_id | `UUID` | NOT NULL | | FK to components(id) CASCADE |
| position | `INTEGER` | NOT NULL | `0` | Display order within collection |
| added_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | When component was added |

**Primary Key:** `(collection_id, component_id)`

### component_exports

**Table:** `component_exports`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| component_id | `UUID` | NOT NULL | | FK to components(id) CASCADE |
| format | `VARCHAR(50)` | NOT NULL | | "npm", "cdn", "source", "zip" |
| url | `TEXT` | NULL | | Export URL |
| expires_at | `TIMESTAMPTZ` | NULL | | Ephemeral link expiry |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Creation timestamp |

**Indexes:** `idx_component_exports_component_id`

### preview_sessions

Ephemeral; does NOT embed SyncModel.

**Table:** `preview_sessions`

| Column | PG Type | Nullable | Default | Description |
|--------|---------|----------|---------|-------------|
| id | `UUID` | NOT NULL | `gen_random_uuid()` | Primary key |
| code | `JSONB` | NOT NULL | `'{}'` | Editor state: {framework, html, css, js, jsx} |
| viewport | `JSONB` | NOT NULL | `'{}'` | Device preset and dimensions |
| created_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Creation timestamp |
| updated_at | `TIMESTAMPTZ` | NOT NULL | `NOW()` | Last update |
| expires_at | `TIMESTAMPTZ` | NULL | | Auto-cleanup expiry |

**Indexes:** `idx_preview_sessions_expires_at`

---

## SQLite Entities (Local Desktop)

Local mirror of cloud data for offline support. Uses SQLite-compatible types: `TEXT` for UUIDs and timestamps (ISO 8601), `INTEGER` for booleans and versions, `TEXT` for JSON strings.

**Migration source:** `orch-ref/database/migrations/sqlite_sync_schema.sql`

### users (SQLite)

Single-row cache of authenticated user.

| Column | SQLite Type | Description |
|--------|-------------|-------------|
| id | `TEXT` PK | UUID |
| name | `TEXT` NOT NULL | Display name |
| email | `TEXT` NOT NULL UNIQUE | Email |
| avatar_url | `TEXT` | Avatar URL |
| plan | `TEXT` NOT NULL DEFAULT 'free' | Plan |
| version | `INTEGER` NOT NULL DEFAULT 1 | Version |
| created_at | `TEXT` NOT NULL | ISO 8601 |
| updated_at | `TEXT` NOT NULL | ISO 8601 |
| deleted_at | `TEXT` | ISO 8601 soft delete |

### devices (SQLite)

| Column | SQLite Type | Description |
|--------|-------------|-------------|
| id | `TEXT` PK | UUID |
| user_id | `TEXT` NOT NULL | User UUID |
| name | `TEXT` NOT NULL | Device name |
| platform | `TEXT` NOT NULL | Platform |
| fingerprint | `TEXT` | Device fingerprint |
| last_seen | `TEXT` NOT NULL | ISO 8601 |
| created_at | `TEXT` NOT NULL | ISO 8601 |
| updated_at | `TEXT` NOT NULL | ISO 8601 |

### settings (SQLite)

| Column | SQLite Type | Description |
|--------|-------------|-------------|
| user_id | `TEXT` NOT NULL | Composite PK part 1 |
| key | `TEXT` NOT NULL | Composite PK part 2 |
| value | `TEXT` | Setting value |
| version | `INTEGER` NOT NULL DEFAULT 1 | Version |
| updated_at | `TEXT` NOT NULL | ISO 8601 |

### notes (SQLite)

| Column | SQLite Type | Description |
|--------|-------------|-------------|
| id | `TEXT` PK | UUID |
| user_id | `TEXT` NOT NULL | User UUID |
| title | `TEXT` NOT NULL DEFAULT '' | Title |
| content | `TEXT` NOT NULL DEFAULT '' | Content |
| tags | `TEXT` NOT NULL DEFAULT '[]' | JSON array string |
| pinned | `INTEGER` NOT NULL DEFAULT 0 | Boolean (0/1) |
| version | `INTEGER` NOT NULL DEFAULT 1 | Version |
| created_at | `TEXT` NOT NULL | ISO 8601 |
| updated_at | `TEXT` NOT NULL | ISO 8601 |
| deleted_at | `TEXT` | ISO 8601 soft delete |

### projects (SQLite)

| Column | SQLite Type | Description |
|--------|-------------|-------------|
| id | `TEXT` PK | UUID |
| user_id | `TEXT` NOT NULL | User UUID |
| slug | `TEXT` NOT NULL | URL-safe identifier |
| name | `TEXT` NOT NULL | Project name |
| description | `TEXT` | Description |
| data | `TEXT` NOT NULL DEFAULT '{}' | JSON string of project tree |
| version | `INTEGER` NOT NULL DEFAULT 1 | Version |
| created_at | `TEXT` NOT NULL | ISO 8601 |
| updated_at | `TEXT` NOT NULL | ISO 8601 |
| deleted_at | `TEXT` | ISO 8601 soft delete |

**Unique:** `(user_id, slug)`

### sync_outbox (SQLite only)

Pending sync records waiting to push. Not present in PostgreSQL.

| Column | SQLite Type | Description |
|--------|-------------|-------------|
| id | `TEXT` PK | UUID |
| entity_type | `TEXT` NOT NULL | Entity type |
| entity_id | `TEXT` NOT NULL | Entity ID |
| action | `TEXT` NOT NULL DEFAULT 'upsert' | "upsert" or "delete" |
| payload | `TEXT` NOT NULL DEFAULT '{}' | JSON payload |
| version | `INTEGER` NOT NULL | Entity version |
| idempotency_key | `TEXT` UNIQUE | Deduplication key |
| retry_count | `INTEGER` NOT NULL DEFAULT 0 | Retry attempts |
| next_retry_at | `TEXT` NOT NULL | Next retry time (ISO 8601) |
| created_at | `TEXT` NOT NULL | ISO 8601 |
| expires_at | `TEXT` NOT NULL | Expiry (7 days) |

### sync_state (SQLite only)

Single-row table tracking sync cursor. Not present in PostgreSQL.

| Column | SQLite Type | Description |
|--------|-------------|-------------|
| id | `INTEGER` PK CHECK (id = 1) | Enforces single row |
| device_id | `TEXT` NOT NULL | Stable device UUID |
| last_sync_at | `TEXT` | Last successful sync (ISO 8601) |
| sync_token | `TEXT` | Server-issued cursor |
| migration_done | `INTEGER` NOT NULL DEFAULT 0 | Whether migration completed |
| updated_at | `TEXT` NOT NULL | ISO 8601 |

---

## Sync Protocol Types

### PushRecord

Single change sent to the server.

**Source:** `orch-ref/app/syncclient/client.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| EntityType | `string` | `entity_type` | Yes | Entity type ("setting", "note", "project") |
| EntityID | `string` | `entity_id` | Yes | Entity identifier |
| Action | `string` | `action` | Yes | "upsert" or "delete" |
| Payload | `map[string]any` | `payload` | Yes | Serialized entity data |
| Version | `int64` | `version` | Yes | Entity version |
| IdempotencyKey | `string` | `idempotency_key,omitempty` | No | Deduplication key |

### PushResult

Per-record response from the server after a push.

**Source:** `orch-ref/app/syncclient/client.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| EntityID | `string` | `entity_id` | Yes | Entity identifier |
| Status | `string` | `status` | Yes | "ok", "skipped", or "error" |
| Error | `string` | `error,omitempty` | No | Error message (if status is "error") |

### PullRecord

Single change received from the server during a pull.

**Source:** `orch-ref/app/syncclient/client.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| EntityType | `string` | `entity_type` | Yes | Entity type |
| EntityID | `string` | `entity_id` | Yes | Entity identifier |
| Action | `string` | `action` | Yes | "upsert" or "delete" |
| Payload | `map[string]any` | `payload` | Yes | Serialized entity data |
| Version | `int64` | `version` | Yes | Entity version |
| IdempotencyKey | `string` | `idempotency_key,omitempty` | No | Deduplication key |
| SyncedAt | `string` | `synced_at` | Yes | Server timestamp of the change |

### ClientStatus

Sync client state snapshot.

**Source:** `orch-ref/app/syncclient/client.go`

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| BaseURL | `string` | `base_url` | Yes | Server base URL |
| DeviceID | `string` | `device_id` | Yes | Device UUID |
| Configured | `bool` | `configured` | Yes | Whether client is configured |
| Authenticated | `bool` | `authenticated` | Yes | Whether token is available |
| LastPullAt | `*time.Time` | `last_pull_at,omitempty` | No | Last successful pull |
| LastPullErr | `string` | `last_pull_err,omitempty` | No | Last pull error message |

---

## gRPC/Proto Messages

All proto definitions use package `orchestra.engine.v1`.

### MemoryService (memory.proto)

**Source:** `orch-ref/proto/memory.proto`

#### MemoryChunk (proto)

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| id | `string` | 1 | Chunk identifier |
| project | `string` | 2 | Project slug |
| source | `string` | 3 | Source type: task, prd, session, user |
| source_id | `string` | 4 | Source entity ID |
| summary | `string` | 5 | Short summary |
| content | `string` | 6 | Full content |
| tags | `repeated string` | 7 | Searchable tags |
| created_at | `string` | 8 | ISO 8601 timestamp |
| score | `float` | 9 | Relevance score (search results only) |

#### Session (proto)

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| id | `string` | 1 | Session ID |
| project | `string` | 2 | Project slug |
| started_at | `string` | 3 | Start time |
| ended_at | `string` | 4 | End time |
| agent_type | `string` | 5 | Agent type |
| model | `string` | 6 | AI model used |
| token_count | `int64` | 7 | Total tokens |
| message_count | `int64` | 8 | Total messages |
| tool_count | `int64` | 9 | Total tool calls |
| metadata | `map<string, string>` | 10 | Key-value metadata |

#### Observation (proto)

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| id | `string` | 1 | Observation ID |
| session_id | `string` | 2 | Parent session ID |
| project | `string` | 3 | Project slug |
| observation_type | `string` | 4 | user_prompt, tool_use, tool_result, assistant_response |
| content | `string` | 5 | Observation content |
| tool_name | `string` | 6 | Tool name (for tool events) |
| tool_input | `string` | 7 | Tool input JSON |
| tool_output | `string` | 8 | Tool output JSON |
| context | `string` | 9 | Context string |
| tokens | `int64` | 10 | Token count |
| sequence | `int64` | 11 | Sequence number |
| timestamp | `string` | 12 | ISO 8601 timestamp |
| metadata | `map<string, string>` | 13 | Key-value metadata |

#### Summary (proto)

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| id | `string` | 1 | Summary ID |
| session_id | `string` | 2 | Parent session ID |
| project | `string` | 3 | Project slug |
| summary_type | `string` | 4 | Summary type |
| content | `string` | 5 | Summary content |
| observation_ids | `repeated string` | 6 | Source observation IDs |
| tokens | `int64` | 7 | Token count |
| timestamp | `string` | 8 | ISO 8601 timestamp |
| metadata | `map<string, string>` | 9 | Key-value metadata |

#### SimilarEntity (proto)

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| entity_type | `string` | 1 | Entity type: observation, summary, session |
| entity_id | `string` | 2 | Entity ID |
| similarity | `float` | 3 | Cosine similarity score |

#### SessionEvent (proto)

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| type | `string` | 1 | Event type |
| summary | `string` | 2 | Event summary |
| timestamp | `string` | 3 | ISO 8601 timestamp |

#### SessionLog (proto)

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| session_id | `string` | 1 | Session ID |
| project | `string` | 2 | Project slug |
| summary | `string` | 3 | Session summary |
| events | `repeated SessionEvent` | 4 | Session events |
| started_at | `string` | 5 | Start time |
| ended_at | `string` | 6 | End time |

#### RPC Request/Response Messages

| Message | Fields |
|---------|--------|
| `StoreChunkRequest` | project, source, source_id, summary, content, tags |
| `StoreChunkResponse` | chunk (MemoryChunk) |
| `MemorySearchRequest` | project, query, limit |
| `MemorySearchResponse` | results (repeated MemoryChunk) |
| `ContextRequest` | project, query, limit |
| `ContextResponse` | chunks (repeated MemoryChunk) |
| `StoreSessionRequest` | project, session_id, summary, events |
| `StoreSessionResponse` | session (SessionLog) |
| `ListSessionsRequest` | project, limit |
| `ListSessionsResponse` | sessions (repeated SessionLog) |
| `GetSessionRequest` | project, session_id |
| `GetSessionResponse` | session (SessionLog) |
| `StartSessionRequest` | id, project, agent_type, model, metadata |
| `StartSessionResponse` | session (Session) |
| `EndSessionRequest` | session_id, ended_at, project |
| `EndSessionResponse` | session (Session) |
| `RecordObservationRequest` | session_id, project, observation_type, content, tool_name, tool_input, tool_output, context, tokens, metadata |
| `RecordObservationResponse` | observation (Observation) |
| `GetTimelineRequest` | session_id, around_sequence, radius |
| `GetTimelineResponse` | observations (repeated Observation) |
| `GetObservationsRequest` | session_id, project |
| `GetObservationsResponse` | observations (repeated Observation) |
| `FetchDetailsRequest` | observation_ids (repeated string) |
| `FetchDetailsResponse` | observations (repeated Observation) |
| `StoreEmbeddingRequest` | entity_type, entity_id, project, model, vector (repeated float) |
| `StoreEmbeddingResponse` | success (bool) |
| `SearchSimilarRequest` | project, query_vector (repeated float), model, limit |
| `SearchSimilarResponse` | entities (repeated SimilarEntity) |
| `StoreSummaryRequest` | session_id, project, summary_type, content, observation_ids, tokens, metadata |
| `StoreSummaryResponse` | summary (Summary) |
| `GetSummariesRequest` | session_id, project |
| `GetSummariesResponse` | summaries (repeated Summary) |

### SearchService (search.proto)

**Source:** `orch-ref/proto/search.proto`

#### IndexFileRequest

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| file_path | `string` | 1 | File path |
| content | `string` | 2 | File content to index |
| language | `string` | 3 | Programming language |
| metadata | `map<string, string>` | 4 | Additional metadata |

#### IndexFileResponse

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| success | `bool` | 1 | Whether indexing succeeded |
| error | `string` | 2 | Error message |
| indexed_count | `int64` | 3 | Number of files indexed |

#### SearchRequest

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| query | `string` | 1 | Search query |
| limit | `int32` | 2 | Max results |
| offset | `int32` | 3 | Pagination offset |
| file_types | `repeated string` | 4 | Filter by file extension |
| fuzzy | `bool` | 5 | Enable fuzzy matching |

#### SearchResponse

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| results | `repeated SearchResult` | 1 | Search results |
| total_hits | `int64` | 2 | Total matching documents |
| search_time_ms | `int64` | 3 | Search duration in ms |

#### SearchResult

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| file_path | `string` | 1 | File path |
| score | `double` | 2 | Relevance score |
| snippets | `repeated string` | 3 | Highlighted text snippets |
| line_number | `int32` | 4 | Match line number |
| metadata | `map<string, string>` | 5 | File metadata |

#### DeleteFileRequest / DeleteFileResponse

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| file_path | `string` | 1 | File path to delete |
| success | `bool` | 1 | Whether deletion succeeded |
| error | `string` | 2 | Error message |

#### ClearIndexRequest / ClearIndexResponse

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| (empty) | | | No fields |
| success | `bool` | 1 | Whether clearing succeeded |
| error | `string` | 2 | Error message |

### ParseService (parse.proto)

**Source:** `orch-ref/proto/parse.proto`

#### ParseFileRequest

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| file_path | `string` | 1 | File path |
| content | `string` | 2 | File content |
| language | `string` | 3 | Programming language |
| include_ast | `bool` | 4 | Whether to return full AST |

#### ParseFileResponse

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| file_path | `string` | 1 | Parsed file path |
| success | `bool` | 2 | Whether parsing succeeded |
| error | `string` | 3 | Error message |
| ast | `string` | 4 | AST as JSON string |
| symbols | `repeated Symbol` | 5 | Extracted symbols |
| parse_time_ms | `int64` | 6 | Parse duration in ms |

#### GetSymbolsRequest

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| file_path | `string` | 1 | File path |
| content | `string` | 2 | File content |
| language | `string` | 3 | Programming language |
| symbol_types | `repeated string` | 4 | Filter by symbol type |

#### GetSymbolsResponse

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| file_path | `string` | 1 | File path |
| symbols | `repeated Symbol` | 2 | Extracted symbols |

#### Symbol

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| name | `string` | 1 | Symbol name |
| kind | `string` | 2 | "function", "class", "variable", etc. |
| range | `Range` | 3 | File location |
| detail | `string` | 4 | Additional information |
| children | `repeated Symbol` | 5 | Nested symbols |

#### Range

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| start_line | `int32` | 1 | Starting line (0-indexed) |
| start_column | `int32` | 2 | Starting column (0-indexed) |
| end_line | `int32` | 3 | Ending line |
| end_column | `int32` | 4 | Ending column |

### ComponentBundlerService (component.proto)

**Source:** `orch-ref/proto/component.proto`

#### BundleComponentRequest

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| framework | `string` | 1 | "html", "react", "vue", "svelte", "angular" |
| html | `string` | 2 | HTML source |
| css | `string` | 3 | CSS source |
| js | `string` | 4 | JavaScript source |
| jsx | `string` | 5 | JSX/component source |
| dependencies | `repeated string` | 6 | External packages |
| content_hash | `string` | 7 | SHA-256 for cache lookup |

#### BundleComponentResponse

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| success | `bool` | 1 | Whether bundling succeeded |
| error | `string` | 2 | Error message |
| bundled_html | `string` | 3 | Self-contained HTML document |
| inlined_css | `string` | 4 | Extracted and inlined CSS |
| inlined_js | `string` | 5 | Bundled and inlined JavaScript |
| content_hash | `string` | 6 | SHA-256 hash for caching |
| bundle_time_ms | `int64` | 7 | Bundling duration in ms |

#### ParseComponentPropsRequest

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| source | `string` | 1 | Component source code |
| framework | `string` | 2 | Component framework |

#### ParseComponentPropsResponse

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| success | `bool` | 1 | Whether parsing succeeded |
| error | `string` | 2 | Error message |
| props | `repeated ComponentProp` | 3 | Extracted props |
| component_name | `string` | 4 | Detected component name |

#### ComponentProp

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| name | `string` | 1 | Prop name |
| type_name | `string` | 2 | TypeScript/PropTypes type |
| required | `bool` | 3 | Whether prop is required |
| default_value | `string` | 4 | Default value expression |
| description | `string` | 5 | JSDoc description |

#### ValidateComponentRequest

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| source | `string` | 1 | Component source code |
| framework | `string` | 2 | Component framework |

#### ValidateComponentResponse

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| valid | `bool` | 1 | Whether component is valid |
| errors | `repeated ValidationError` | 2 | Validation errors |
| warnings | `repeated ValidationError` | 3 | Non-fatal warnings |

#### ValidationError

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| message | `string` | 1 | Error description |
| line | `int32` | 2 | Source line (1-indexed, 0 = unknown) |
| column | `int32` | 3 | Source column (1-indexed, 0 = unknown) |
| severity | `string` | 4 | "error" or "warning" |

### HealthService (health.proto)

**Source:** `orch-ref/proto/health.proto`

#### HealthCheckRequest

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| service | `string` | 1 | Service name to check (optional) |

#### HealthCheckResponse

| Field | Proto Type | Number | Description |
|-------|-----------|--------|-------------|
| status | `ServingStatus` | 1 | Health status enum |
| message | `string` | 2 | Details about health |
| timestamp | `int64` | 3 | Check timestamp |

#### ServingStatus (enum)

| Value | Number | Description |
|-------|--------|-------------|
| UNKNOWN | 0 | Status unknown |
| SERVING | 1 | Service is healthy |
| NOT_SERVING | 2 | Service is unhealthy |
| SERVICE_UNKNOWN | 3 | Requested service not found |

---

## Workflow States

**Source:** `orch-ref/app/workflow/workflow.go`

### 13-State Lifecycle

```
backlog -> todo -> in-progress -> ready-for-testing -> in-testing ->
ready-for-docs -> in-docs -> documented -> in-review -> done
```

Plus: `blocked`, `rejected`, `cancelled`

### All States

| State | Category | Description |
|-------|----------|-------------|
| `backlog` | Backlog | Not yet planned |
| `todo` | Planned | Planned for work |
| `in-progress` | Active | Work is happening |
| `blocked` | Blocked | Waiting on a blocker |
| `ready-for-testing` | Waiting | Code complete, awaiting test |
| `in-testing` | Active | Tests being run |
| `ready-for-docs` | Waiting | Tests passed, awaiting docs |
| `in-docs` | Active | Documentation being written |
| `documented` | Waiting | Docs written, awaiting review |
| `in-review` | Active | Code review in progress |
| `done` | Completed | Successfully completed |
| `rejected` | Completed | Rejected in review |
| `cancelled` | Completed | Cancelled |

### Transition Map

| From | Valid Targets |
|------|---------------|
| `backlog` | `todo` |
| `todo` | `in-progress`, `backlog` |
| `in-progress` | `ready-for-testing`, `blocked`, `todo` |
| `blocked` | `in-progress`, `todo` |
| `ready-for-testing` | `in-testing`, `in-progress` |
| `in-testing` | `ready-for-docs`, `in-progress` |
| `ready-for-docs` | `in-docs`, `in-testing` |
| `in-docs` | `documented`, `ready-for-docs` |
| `documented` | `in-review` |
| `in-review` | `done`, `rejected`, `documented` |
| `done` | `todo` |
| `rejected` | `todo`, `backlog` |
| `cancelled` | `backlog` |

### Happy-Path Advance Map

Used by `advance_task` for automatic progression:

| Current State | Next State |
|--------------|------------|
| `in-progress` | `ready-for-testing` |
| `ready-for-testing` | `in-testing` |
| `in-testing` | `ready-for-docs` |
| `ready-for-docs` | `in-docs` |
| `in-docs` | `documented` |
| `documented` | `in-review` |
| `in-review` | `done` |

### State Categories

| Category | States | Description |
|----------|--------|-------------|
| Active | `in-progress`, `in-testing`, `in-docs`, `in-review` | Work is actively happening |
| Waiting | `ready-for-testing`, `ready-for-docs`, `documented` | Waiting for next phase |
| Completed | `done`, `rejected`, `cancelled` | Terminal states (resolved) |
| Done | `done` | Successfully finished |

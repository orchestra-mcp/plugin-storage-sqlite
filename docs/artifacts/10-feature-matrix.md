# Feature Matrix -- Orchestra Reference

> Maps every user-facing feature to the services it touches. Use this to understand dependencies, gauge complexity, and prioritize what to build first.
>
> Cross-referenced from: [01-service-catalog](01-service-catalog.md) (7 services), [04-mcp-tool-catalog](04-mcp-tool-catalog.md) (186 tools / 34 categories), [03-api-surface](03-api-surface.md) (134 REST routes, WebSocket events, gRPC services), [05-integration-map](05-integration-map.md) (10 integrations), [09-ui-inventory](09-ui-inventory.md) (19 frontend packages).

---

## Feature-to-Service Matrix

The columns represent the seven services from the service catalog plus the integration layer:

| Column | Service | Protocol | Port |
|--------|---------|----------|------|
| **MCP** | MCP Server | JSON-RPC 2.0 stdio / SSE | - |
| **Settings** | Settings API | HTTP REST | 19191 |
| **WS** | WebSocket Server | WebSocket (Fiber v3) | 8765 |
| **Engine** | Rust Engine | gRPC (Tonic) | 50051 |
| **Desktop** | Desktop App | Wails v3 GUI | - |
| **Chrome** | Chrome Extension | Manifest V3 | - |
| **Integrations** | External services | Various | - |

---

### Core Project Management

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Create / list / delete projects | `create_project`, `list_projects`, `delete_project` (3) | `POST /api/mcp/tools/call` | -- | -- | ProjectsView page | -- | -- | tasks (ProjectsSidebar, ProjectItem) |
| Get project status & tree | `get_project_status`, `get_project_tree` (2) | `POST /api/mcp/tools/call` | -- | -- | ProjectStatusPage | -- | -- | tasks (ProjectDashboard, BacklogTree), ai (ProjectStatusCard, ProjectTreeCard) |
| Read / write PRD | `read_prd`, `write_prd` (2) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (PrdContentCard) |
| Regenerate README | `regenerate_readme` (1) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | -- |

### Epic Management

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| CRUD epics | `list_epics`, `create_epic`, `get_epic`, `update_epic`, `delete_epic` (5) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | tasks (BacklogView, BacklogTree), ai (EpicCard) |

### Story Management

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| CRUD stories | `list_stories`, `create_story`, `get_story`, `update_story`, `delete_story` (5) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | tasks (BacklogView, TaskList), ai (StoryCard) |

### Task Management

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| CRUD tasks | `list_tasks`, `create_task`, `get_task`, `update_task`, `delete_task` (5) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | tasks (TaskDetailPanel, TaskList, StatusBadge, PriorityIcon, InlineEdit), ai (TaskCard, TaskDetailCard) |
| Assign / unassign tasks | `assign_task`, `unassign_task`, `my_tasks` (3) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | tasks (TaskFilter) |
| Labels & estimates | `add_labels`, `remove_labels`, `set_estimate`, `add_link` (4) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | tasks (LabelsSection, LinksSection) |
| Templates | `save_template`, `list_templates`, `create_from_template` (3) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | -- |

### Workflow & Lifecycle (13 States)

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Workflow state machine | `get_next_task`, `set_current_task`, `complete_task`, `search`, `get_workflow_status` (5) | `POST /api/mcp/tools/call` | `tasks:changed` event | -- | -- | -- | -- | tasks (StatusBadge, StatusActions, LifecycleProgressBar, BacklogView), ai (WorkflowStatusCard) |
| Lifecycle gating (advance/reject) | `advance_task`, `reject_task` (2) | `POST /api/mcp/tools/call` | `tasks:changed` event | -- | -- | -- | -- | tasks (EvidenceLog, LifecycleProgressBar), ai (GateCard) |
| Workflow notifications | -- | -- | `tasks:changed` broadcast | -- | Desktop notification | -- | Discord `TransitionListener`, Slack `TransitionListener` | -- |

### Sprint Management

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Sprint CRUD & lifecycle | `create_sprint`, `list_sprints`, `get_sprint`, `start_sprint`, `end_sprint` (5) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (SprintCard), tasks (SprintProgressWidget) |
| Sprint task management | `add_sprint_tasks`, `remove_sprint_tasks` (2) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | tasks (BacklogView) |
| Backlog reordering | `reorder_backlog` (1) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | tasks (BacklogView) |

### Scrum Ceremonies

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Daily standup summary | `get_standup_summary` (1) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (StandupCard) |
| Burndown chart | `get_burndown` (1) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (BurndownChartCard), widgets (LineChart, AreaChart) |
| Velocity tracking | `get_velocity` (1) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (VelocityCard), widgets (BarChart) |
| Retrospectives | `create_retrospective` (1) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (RetrospectiveCard) |

### WIP Limits & Dependencies

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| WIP limit enforcement | `set_wip_limits`, `get_wip_limits`, `check_wip_limit` (3) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (WipLimitCard) |
| Task dependencies | `add_dependency`, `remove_dependency`, `get_dependency_graph` (3) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (DependencyGraphCard) |

### PRD Authoring (Multi-Audience)

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Guided PRD sessions | `start_prd_session`, `answer_prd_question`, `get_prd_session`, `abandon_prd_session`, `skip_prd_question`, `back_prd_question` (6) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (PRDSessionCard, PRDQuestionCard) |
| PRD preview & validation | `preview_prd`, `validate_prd` (2) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (PRDPreviewCard, PrdValidationCard) |
| PRD phases & backlog generation | `split_prd`, `list_prd_phases`, `generate_backlog`, `get_agent_briefing` (4) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (PrdPhasesCard, AgentBriefingCard) |
| PRD templates | `save_prd_template`, `list_prd_templates`, `load_prd_template` (3) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | -- |

### Notes

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Notes CRUD & search | `save_note`, `list_notes`, `search_notes`, `update_note`, `delete_note` (5) | `POST /api/mcp/tools/call` | `sync:note:upsert`, `sync:note:delete` | -- | NotesPage | -- | Notion (`POST /api/notion/send`), Apple Notes (`POST /api/applenotes/send`) | desktop-ui (NotesView, NotesSidebar, NoteItem, NoteBody), ai (NoteCard) |

### Bug Reporting & Feature Requests

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Report bugs & log requests | `report_bug`, `log_request` (2) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | -- |

### Plans & Artifacts

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Save / list plans | `save_plan`, `list_plans` (2) | `POST /api/mcp/tools/call` | -- | -- | -- | -- | -- | ai (PlanCard) |

### AI Chat

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| AI chat sessions | -- | `POST /api/ai/send`, `POST /api/ai/followup`, `POST /api/ai/stop`, `GET /api/ai/events` (SSE), `GET /api/ai/sessions`, `POST /api/ai/sessions/new`, `DELETE /api/ai/sessions/{id}` | `ai:send`, `ai:stop`, `ai:chunk`, `ai:session_created` | -- | MainChatPage, SessionChat, BubblePage | -- | Anthropic / OpenAI SDKs (via AI Bridge) | ai (ChatBox, ChatHeader, ChatBody, ChatInput, ChatMessage, ChatStreamMessage, ChatModelSelector, ChatModeSelector + 50 event cards) |
| AI provider management | -- | `GET /api/ai/providers`, `POST /api/ai/providers/set`, `GET /api/ai/models` | -- | -- | SettingsPage (AIProvidersSettings) | -- | -- | settings (SettingsForm) |
| AI permission/question flow | -- | `POST /api/ai/permission`, `POST /api/prompts/request`, `GET /api/prompts/pending/{id}`, `POST /api/prompts/response/{id}` | `ai:permission`, `ai:question` | -- | PromptPage | -- | -- | ai (ConfirmationCard, QuestionCard) |
| Message forwarding between sessions | -- | -- | `ai:forward`, `ai:message_forwarded` | -- | -- | -- | -- | ai (SubAgentPage) |

### Memory & Context (Cross-Session)

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Save / search memory | `save_memory`, `search_memory`, `get_context` (3) | `POST /api/mcp/tools/call` | -- | gRPC `MemoryService` (StoreEmbedding, SearchSimilar, RecordObservation, SearchMemory, GetContext) | -- | -- | -- | ai (MemoryCard) |
| Session memory | `save_session`, `list_sessions`, `get_session` (3) | `POST /api/mcp/tools/call` | -- | gRPC `MemoryService` (StartSession, EndSession, ListSessions, GetSession, StoreSummary) | -- | -- | -- | ai (SessionCard) |

### Code Search & Parsing

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Full-text code search | -- | `POST /api/search/mentions` | -- | gRPC `SearchService` (IndexFile, IndexFiles, Search, DeleteFile, ClearIndex) | -- | -- | -- | search (SearchSpotlight) |
| Code parsing (Tree-sitter) | -- | -- | -- | gRPC `ParseService` (ParseFile, ParseFiles, GetSymbols -- 14 languages) | -- | -- | -- | -- |
| @-mention search | -- | `POST /api/search/mentions` | -- | gRPC `SearchService` (fallback to glob) | -- | -- | -- | ai (MentionPopup, MentionToken, MentionTokens) |

### Component Library & Preview

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Component CRUD | `save_component`, `list_components`, `export_component` (3) | `GET/POST/PUT/DELETE /api/v1/components` | -- | gRPC `ComponentBundlerService` (BundleComponent, ParseComponentProps, ValidateComponent) | ComponentBrowserPage, ComponentEditorPage | -- | -- | ai (SmartComponentCard, ExportCard, ExportConfigDialog) |
| Live preview sessions | `preview_component`, `update_preview`, `set_preview_viewport`, `open_browser_preview` (4) | `POST /preview`, `GET /preview/{id}`, `DELETE /preview/{id}` (on :8765) | `preview:update`, `preview:viewport`, `preview:join`, `preview:leave`, `preview:open_browser` | -- | -- | PreviewApp, PreviewViewportToolbar | -- | ai (PreviewFrame, PreviewViewportToolbar, PreviewCard) |

### Claude Code Awareness

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Skills & agents | `list_skills`, `list_agents`, `install_skills`, `install_agents`, `install_docs` (5) | -- | -- | -- | -- | -- | -- | ai (SkillCard, AgentSwitchCard) |
| Hook events | `receive_hook_event`, `get_hook_events` (2) | -- | -- | -- | -- | -- | -- | ai (HookEventCard) |
| Cross-session memory search | `search_memory` (Claude variant) (1) | -- | -- | gRPC `MemoryService` | -- | -- | -- | ai (MemoryCard) |

### Desktop Window Modes

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Embedded / floating / bubble modes | -- | `GET/POST /api/mode`, `POST /api/mode/cycle` | -- | -- | Wails window management | -- | -- | ai (ModeToggle, BubbleButton), desktop-ui (MainLayout) |
| Spirit window | `open_desktop_window` (1) | `POST /api/spirit/toggle`, `/open`, `/close` | -- | -- | Spirit mode | -- | -- | ai (BubbleButton) |
| Bubble window | -- | `POST /api/bubble/toggle`, `/open`, `/close` | -- | -- | Bubble mode | -- | -- | -- |
| Open named windows | `open_desktop_window` (1) | `POST /api/windows/open`, `POST /api/windows/close/{name}`, `GET /api/windows/data/{name}` | -- | -- | Wails window API | -- | -- | desktop-ui (PanelContainer) |

### System Tray & Notifications

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Desktop notifications | `send_notification` (1) | `POST /api/notify` | `notification:speak` (TTS) | -- | macOS native notifications | -- | -- | -- |
| Sound effects | `play_sound` (1) | `POST /api/notify` | `notification:speak` | -- | Bundled MP3 sounds | -- | -- | -- |
| System tray | -- | -- | -- | -- | Wails system tray | -- | -- | -- |

### DevTools

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Terminal sessions | `create_dev_session`, `terminal_exec` (2) | `POST /api/devtools/terminal-exec` | `/ws/devtools/*` | -- | DevTools session manager | -- | -- | devtools (TerminalSession), ai (TerminalPage) |
| Database sessions | `create_dev_session`, `run_query`, `list_databases` (3) | `POST /api/devtools/run-query` | `/ws/devtools/*` | -- | DevTools session manager | -- | -- | devtools (DatabaseSession) |
| SSH sessions | `create_dev_session`, `ssh_exec` (2) | `POST /api/devtools/ssh-exec` | `/ws/devtools/*` | -- | DevTools session manager | -- | -- | devtools (SSHSession) |
| Log viewer | `create_dev_session`, `detect_logs`, `search_logs`, `view_logs` (4) | -- | `/ws/devtools/*` | -- | DevTools session manager | -- | -- | devtools (LogViewerSession) |
| Service manager | `list_services`, `control_service`, `manage_service` (3) | -- | `/ws/devtools/*` | -- | DevTools session manager | -- | -- | devtools (ServiceManagerSession) |
| Docker stacks | `get_stacks`, `install_stack`, `stack_control` (3) | -- | `/ws/devtools/*` | -- | DevTools session manager | -- | -- | devtools (ServiceManagerSession) |
| Debugger | `create_dev_session`, `debug_launch`, `debug_control`, `set_breakpoint`, `debug_attach` (5) | -- | `/ws/devtools/*` | -- | DevTools session manager | -- | -- | devtools (DebuggerSession) |
| Test runner | `run_tests` (1) | -- | `/ws/devtools/*` | -- | DevTools session manager | -- | -- | devtools (TestingSession) |
| Cloud deploy | `cloud_deploy` (1) | -- | `/ws/devtools/*` | -- | DevTools session manager | -- | -- | devtools (CloudSession) |
| File explorer | `create_dev_session`, `read_file_session` (2) | -- | `/ws/devtools/*` | -- | DevTools session manager | -- | -- | devtools (FileExplorerSession), explorer (FileTree) |

### Code Editing

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Monaco editor | -- | -- | `/ws/lsp/{language}` (LSP proxy) | gRPC `ParseService` (symbol extraction) | LSP server management | -- | -- | editor (CodeEditor, CodeDiffEditor, MonacoLoader, CodeBlock, useLsp) |
| Markdown editing | -- | -- | -- | -- | MarkdownViewerPage | -- | -- | editor (MarkdownEditor, MarkdownRenderer) |
| Git diff view | -- | -- | -- | -- | -- | -- | -- | editor (GitDiffView), ai (FilesChangedPanel) |

### Figma Integration

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Figma auth (OAuth PKCE / PAT) | -- | `GET /api/figma/auth/start`, `/callback`, `/status`, `POST /api/figma/auth/pat`, `/disconnect` | -- | -- | FigmaSettings page | -- | Figma OAuth 2.0 PKCE | account-center (AccountIntegration) |
| Figma file browsing | `figma_get_file_meta`, `figma_get_file`, `figma_get_nodes`, `figma_get_components` (4) | `GET /api/figma/files`, `/files/{key}`, `/files/{key}/nodes`, `/files/{key}/components` | -- | -- | -- | -- | Figma REST API (with retry + cache) | ai (FigmaCard) |
| Figma MCP bridge install | -- | `GET /api/figma/mcp/status`, `POST /api/figma/mcp/install` | -- | -- | -- | -- | figma-developer-mcp npm package | -- |

### GitHub Integration

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| GitHub auth (OAuth / PAT) | `github_auth_start`, `github_auth_status`, `github_auth_pat`, `github_sign_out` (4) | `GET /api/github/auth/start`, `/callback`, `/status`, `POST /api/github/auth/pat`, `/disconnect` | -- | -- | GitHubSettings page | -- | GitHub OAuth 2.0 | account-center (AccountIntegration) |
| Issues CRUD | `github_list_issues`, `github_get_issue`, `github_create_issue`, `github_close_issue`, `github_issue_comment` (5) | -- | -- | -- | -- | -- | GitHub REST API | ai (GitHubIssueCard) |
| Pull requests | `github_list_prs`, `github_get_pr`, `github_pr_files`, `github_create_pr`, `github_merge_pr`, `github_review_pr`, `github_pr_comment` (7) | -- | -- | -- | -- | -- | GitHub REST API | ai (GitHubPRCard) |
| CI/CD status | `github_ci_status` (1) | -- | -- | -- | -- | -- | GitHub REST API | ai (CIStatusCard) |
| Repo tracking | `github_list_repos`, `github_track_repo`, `github_untrack_repo` (3) | -- | -- | -- | -- | -- | -- | -- |

### Jira Integration

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Jira auth (OAuth 3LO) | -- | `GET /api/jira/auth/start`, `/callback`, `/status`, `POST /api/jira/disconnect` | -- | -- | JiraSettings page | -- | Atlassian OAuth 2.0 (3LO) | account-center (AccountIntegration) |
| Jira issue sync | -- | -- | -- | -- | -- | -- | Jira REST API v3 (JQL search, create, update, transition, comment) | -- |

### Linear Integration

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Linear auth (OAuth) | -- | `GET /api/linear/auth/start`, `/callback`, `/status`, `POST /api/linear/disconnect` | -- | -- | LinearSettings page | -- | Linear OAuth 2.0 | account-center (AccountIntegration) |
| Linear issue sync | -- | -- | -- | -- | -- | -- | Linear GraphQL API (issues, teams, projects, cycles) | -- |

### Notion Integration

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Notion auth (OAuth) | -- | `GET /api/notion/auth/start`, `/callback`, `/status`, `POST /api/notion/disconnect`, `GET /api/notion/config` | -- | -- | NotionSettings page | -- | Notion OAuth 2.0 | account-center (AccountIntegration) |
| Send notes to Notion | -- | `GET /api/notion/databases`, `POST /api/notion/send` | -- | -- | -- | -- | Notion REST API (markdown-to-blocks conversion) | -- |

### Discord Integration

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Discord bot | -- | `GET /api/discord/config`, `POST /api/discord/test` | -- | -- | DiscordSettings page | -- | Discord Gateway WebSocket v10 + REST API v10 | settings (SettingsForm) |
| Discord AI chat | -- | -- | -- | -- | -- | -- | Discord Gateway (ChatHandler with AI bridge) | -- |
| Discord MCP tools | -- | -- | -- | -- | -- | -- | Discord Gateway (McpHandler, ToolsHandler, PromptsHandler) | -- |
| Workflow notifications | -- | -- | -- | -- | -- | -- | Discord webhook / bot API (TransitionListener embeds) | -- |

### Slack Integration

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Slack bot | -- | `GET /api/slack/config`, `POST /api/slack/test` | -- | -- | SlackSettings page | -- | Slack Socket Mode + REST API | settings (SettingsForm) |
| Slack AI chat & MCP tools | -- | -- | -- | -- | -- | -- | Slack Socket Mode (same 10 handlers as Discord) | -- |
| Workflow notifications | -- | -- | -- | -- | -- | -- | Slack webhook / bot API (Block Kit) | -- |

### Apple Notes Integration

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Send to Apple Notes | -- | `GET /api/applenotes/status`, `GET /api/applenotes/folders`, `POST /api/applenotes/send` | -- | -- | AppleSettings page | -- | macOS AppleScript (osascript) | -- |

### Team Management

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Teams CRUD | `list_teams`, `create_team`, `get_team` (3) | Proxy to Laravel: `GET/POST /api/teams`, `GET/PUT/DELETE /api/teams/{id}` | -- | -- | TeamsPage, TeamDetailPage | -- | Orchestra Web (Laravel Sanctum) | ai (TeamCard) |
| Invitations & sharing | `invite_member`, `get_pending_invitations`, `share_with_team` (3) | `POST /api/teams/{id}/invitations`, `GET /api/teams/pending-invitations`, `POST /api/teams/{id}/shares` | -- | -- | -- | -- | Orchestra Web | -- |

### Usage & Analytics

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Token usage tracking | `get_usage`, `record_usage`, `reset_session_usage` (3) | -- | -- | -- | -- | -- | -- | ai (UsageCard) |
| Session metrics | `start_session_metrics`, `record_session_message`, `get_session_metrics` (3) | -- | -- | -- | -- | -- | -- | tasks (SessionMetricsWidget) |
| Feature analytics | `track_event`, `get_analytics_events` (2) | -- | -- | -- | -- | -- | Firebase GA4 (Measurement Protocol) | shared (Firebase analytics module) |

### Authentication & Sync

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Orchestra web auth | -- | `POST /api/auth/login`, `GET /api/auth/status`, `POST /api/auth/logout`, `POST /api/auth/store`, `GET /api/auth/token` | -- | -- | LoginPage, OnboardingPage (AuthStep) | -- | Orchestra Web (Laravel Sanctum) | -- |
| Cloud sync (push/pull) | -- | `GET /api/sync/migration/status`, `POST /api/sync/migration/start`, `GET /api/sync/web/status`, `POST /api/integrations/sync` | `sync:*` events (note, project, integration, ai_session) | -- | SyncSettingsPage | -- | Orchestra Web (HTTP polling, outbox pattern) | -- |
| Settings persistence | -- | `GET/POST /api/settings`, `GET /api/settings/{key}` | `settings.set`, `settings.get`, `settings.get_all` | -- | SettingsPage | SettingsPanel (chrome.storage) | -- | settings (SettingsForm, SettingsNav), chrome-ui (SettingsPanel) |

### Voice (STT/TTS)

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Speech-to-text (local) | -- | `GET/PUT /api/voice/settings` | -- | -- | VoiceChatButton, VoiceInputOverlay | -- | Hugging Face Transformers.js (WebGPU/WASM, on-device) | voice (OrchestraSTT) |
| Text-to-speech (local) | -- | `GET/PUT /api/voice/settings` | `notification:speak` | -- | WaveformVisualizer | -- | Kokoro.js (WebGPU/WASM, on-device) | voice (OrchestraTTS) |
| Meeting transcript analysis | -- | `POST /api/voice/analyze-transcript`, `POST /api/voice/create-project-from-meeting` | -- | -- | -- | -- | Anthropic API | -- |

### Screenshot

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Screen capture | -- | `POST /api/screenshot/capture`, `POST /api/screenshot/region`, `POST /api/screenshot/start` | -- | -- | ScreenshotPage (macOS screencapture) | -- | -- | ai (useScreenshot hook) |

### Chrome Extension Features

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Page context extraction | -- | `GET /api/browser/context` | `page:context_update`, `page:request_context`, `page.updated` | -- | -- | Content script (extract.ts), Service worker | -- | chrome-ui (Sidebar, ViewBody) |
| Desktop discovery & sync | -- | -- | WebSocket client connection | -- | -- | DesktopDiscovery, WebSocketClient, SyncProtocol, SettingsSyncService, ThemeSyncService | -- | chrome-ui (useDesktopConnection) |
| Code preview in browser | -- | -- | `preview:open_browser` | -- | -- | PreviewApp (sandboxed iframe) | -- | chrome-ui (TabBar) |
| Google Meet integration | -- | -- | -- | -- | -- | meet-injector.ts content script | -- | -- |

### Theming

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| 26 color themes, 3 variants | -- | Settings persistence | -- | -- | ThemePicker | -- | -- | theme (ThemePicker, initTheme, setColorTheme, setComponentVariant), desktop-ui (useThemeSync) |

### Workspace Management

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Workspace switching | -- | `GET /api/workspace`, `POST /api/workspace/open`, `GET /api/workspace/recent` | `workspace:changed` | -- | -- | -- | -- | -- |
| Active file context | -- | `GET/POST /api/workspace/active-file` | -- | -- | -- | -- | -- | -- |
| Current task tracking | -- | `GET /api/workspace/current-task`, `POST /api/workspace/current-task/refresh` | -- | -- | -- | -- | -- | -- |

### File Management

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| File download | -- | `POST /api/download` | -- | -- | Native Save As dialog | -- | -- | -- |
| File picker | -- | `POST /api/dialog/open-file` | -- | -- | Native file picker dialog | -- | -- | -- |
| Open URL in browser | -- | `POST /api/open-url` | -- | -- | System browser launch | -- | -- | -- |

### macOS Permissions

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Permission management | -- | `GET /api/permissions/status`, `POST /api/permissions/request`, `POST /api/permissions/open-settings` | -- | -- | OnboardingPage (PermissionsStep) | -- | macOS APIs | -- |

### Firebase (Push Notifications & Analytics)

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Push notifications (FCM) | -- | -- | -- | -- | useFCM hook | -- | Firebase Admin SDK (FCM: single, multicast, topic) | shared (Firebase messaging module) |
| Crash reporting | -- | -- | -- | -- | -- | -- | Sentry SDK + Firebase Crashlytics (server-side) | shared (Sentry capture module, ErrorBoundary) |

### Auto-Updater

| Feature | MCP | Settings | WS | Engine | Desktop | Chrome | Integrations | UI Packages |
|---------|-----|----------|----|--------|---------|--------|-------------|-------------|
| Automatic updates | -- | -- | -- | -- | Auto-updater (skipped in dev mode) | -- | -- | -- |

---

## Feature Complexity Rankings

### Tier 1: Core (MCP-only, no external dependencies)

These features run entirely within the MCP server using local TOON/YAML files. Zero service dependencies beyond the `bin/orchestra` binary.

- Project CRUD (3 tools)
- Epic CRUD (5 tools)
- Story CRUD (5 tools)
- Task CRUD (5 tools)
- Workflow state machine -- 13 states (5 tools)
- Lifecycle gating -- advance/reject with evidence (2 tools)
- Sprint management (7 tools)
- Scrum ceremonies -- standup, burndown, velocity, retrospectives (4 tools)
- WIP limits (3 tools)
- Dependencies -- add/remove/graph with cycle detection (3 tools)
- Task metadata -- assign, labels, estimates, links (7 tools)
- Templates -- save/list/create-from (3 tools)
- PRD authoring -- multi-audience, conditional follow-ups (6 tools)
- PRD validation, phases, backlog generation, agent briefings (4 tools)
- PRD templates (3 tools)
- Notes (5 tools)
- Bug reporting & feature requests (2 tools)
- Plans & artifacts (2 tools)
- README generation (1 tool)
- Usage tracking (3 tools)
- Session metrics (3 tools)
- Claude Code awareness -- skills, agents, hooks, install (7 tools)
- Issue search (1 tool)

**Total: ~79 MCP tools. Build cost: Low. Dependencies: None.**

### Tier 2: Local Services Required

These features need one or more of the local services (Settings API, WebSocket, Rust Engine) in addition to MCP.

- **AI Chat** -- requires Settings API (AI Bridge + session management) + WebSocket (streaming) + AI provider SDK keys
- **Memory & context search** -- requires Rust Engine (gRPC MemoryService for vector storage) or falls back to markdown
- **Code search** -- requires Rust Engine (gRPC SearchService / Tantivy index) or falls back to glob
- **Code parsing** -- requires Rust Engine (gRPC ParseService / Tree-sitter)
- **Component bundling** -- requires Rust Engine (gRPC ComponentBundlerService)
- **Live preview** -- requires WebSocket server (preview coordination + session management)
- **Workflow notifications** -- requires WebSocket server (broadcast to all connected frontends)
- **Settings persistence** -- requires Settings API (SQLite store)
- **@-mention search** -- requires Settings API + optionally Rust Engine
- **LSP proxy** -- requires Settings API WebSocket endpoint (`/ws/lsp/{language}`)
- **DevTools** -- requires Settings API + Desktop (session manager) + WebSocket (`/ws/devtools/*`)

**Build cost: Medium. Dependencies: 1-3 local services per feature.**

### Tier 3: Desktop Platform Required

These features depend on the Wails v3 desktop app being built and running.

- **Spirit / bubble / embedded window modes** -- Wails window management APIs
- **System tray** -- Wails system tray integration
- **Named window management** -- Wails multi-window API
- **Desktop notifications with sound** -- macOS native notifications + bundled MP3 sounds
- **Screenshot capture** -- macOS `screencapture` via Settings API
- **File download / picker dialogs** -- native OS dialogs via Wails
- **macOS permissions** -- Accessibility, Screen Recording, Microphone, Notifications
- **Auto-updater** -- desktop-only update check
- **DevTools session manager** -- terminal, database, SSH, logs, debugger sessions
- **Voice STT/TTS** -- requires desktop (WebGPU context for Transformers.js / Kokoro.js)
- **Onboarding flow** -- desktop-only first-run wizard

**Build cost: High. Dependencies: Wails v3 + macOS APIs.**

### Tier 4: Chrome Extension Required

These features depend on the Chrome extension being built and deployed.

- **Page context extraction** -- content script + service worker
- **Browser-to-desktop sync** -- DesktopDiscovery + WebSocket bridge
- **Code preview in browser tab** -- PreviewApp sandboxed iframe
- **Google Meet integration** -- meet-injector content script
- **Settings sync** -- chrome.storage + SettingsSyncService + ThemeSyncService

**Build cost: Medium. Dependencies: Chrome Manifest V3 APIs.**

### Tier 5: External Integrations

These features require OAuth flows, API keys, and external service connectivity.

- **GitHub** -- OAuth 2.0 / PAT, 20 MCP tools, issues/PRs/CI (bidirectional sync)
- **Jira** -- OAuth 2.0 3LO, JQL search, issue CRUD + transitions
- **Linear** -- OAuth 2.0, GraphQL API, issues/teams/projects/cycles
- **Notion** -- OAuth 2.0, push notes (markdown-to-blocks conversion)
- **Figma** -- OAuth 2.0 PKCE / PAT, 4 MCP tools, file/node/component read
- **Discord** -- bot token, Gateway WebSocket, slash commands, AI chat, MCP tool execution
- **Slack** -- bot + app tokens, Socket Mode, same 10 handlers as Discord
- **Firebase** -- service account, FCM push, GA4 analytics, Crashlytics
- **Orchestra Web** -- Laravel Sanctum, cloud sync (pull polling + outbox push)
- **Apple Notes** -- macOS AppleScript (no credentials, macOS-only)

**Build cost: High per integration. Dependencies: OAuth infrastructure + API clients.**

---

## Build Priority Recommendation

What to build first for a working v1, ordered by dependency chain:

### Phase 1: MCP Core (weeks 1-2)

Everything that works with `orch --workspace .` stdio, zero external dependencies.

1. **Project + Epic + Story + Task CRUD** (18 tools) -- the foundation
2. **Workflow state machine** (5 tools) -- task progression through 13 states
3. **Lifecycle gating** (2 tools) -- advance/reject with evidence
4. **Sprint management** (7 tools) -- create/start/end sprints
5. **WIP limits + dependencies** (6 tools) -- constraint enforcement
6. **Task metadata** (7 tools) -- assign, labels, estimates, links
7. **Notes + plans + artifacts** (9 tools) -- knowledge management
8. **PRD authoring** (15 tools) -- full guided PRD with phases and templates
9. **Bug reporting + templates** (5 tools) -- issue templates and bug tracking
10. **Scrum ceremonies** (4 tools) -- standup, burndown, velocity, retrospectives

_Deliverable: Fully functional CLI-based project management via MCP protocol._

### Phase 2: Settings API + AI Chat (weeks 3-4)

The local HTTP server that enables GUI access to everything.

11. **Settings API** (HTTP on :19191) -- SQLite settings store, MCP bridge over HTTP
12. **AI Bridge + chat sessions** -- multi-provider AI chat with streaming
13. **Prompt/permission flow** -- tool approval and question UI
14. **Workspace management** -- switch workspaces, track active file and current task

_Deliverable: Local AI-powered IDE assistant accessible via HTTP._

### Phase 3: Rust Engine (weeks 5-6)

CPU-intensive operations that enhance but are not required for core functionality.

15. **Memory service** (gRPC) -- vector-based cross-session memory with RAG
16. **Search service** (gRPC) -- Tantivy full-text code search
17. **Parse service** (gRPC) -- Tree-sitter code parsing for 14 languages
18. **Component bundler service** (gRPC) -- self-contained HTML bundling

_Deliverable: Intelligent code search, memory, and parsing. Falls back to markdown/glob when engine is unavailable._

### Phase 4: WebSocket + Real-Time (weeks 7-8)

Real-time bidirectional communication for all frontends.

19. **WebSocket server** (:8765) -- client management, broadcast, ping/pong
20. **AI streaming over WebSocket** -- alternative to SSE for desktop/extension
21. **Workflow transition broadcasts** -- real-time task state change notifications
22. **Sync handler** -- notes, projects, integrations, AI sessions
23. **Preview coordinator** -- live component preview sessions
24. **Browser handler** -- Chrome extension page context awareness

_Deliverable: Real-time updates across all connected clients._

### Phase 5: Desktop App (weeks 9-12)

The native macOS application shell.

25. **Wails v3 app shell** -- MainLayout, PanelLayout, Sidebar, StatusBar, Topbar
26. **Chat pages** -- MainChatPage, SessionChat, BubblePage
27. **Task management pages** -- ProjectsView, ProjectStatusPage, BacklogView
28. **Notes pages** -- NotesPage with editor
29. **Settings pages** -- SettingsPage, IntegrationsPage
30. **DevTools** -- terminal, database, SSH, logs, debugger sessions
31. **Window modes** -- embedded, floating, bubble, spirit
32. **Notifications + sounds** -- macOS native notifications with bundled sounds
33. **Onboarding** -- first-run wizard
34. **Auto-updater** -- update check and install

_Deliverable: Full native macOS desktop IDE._

### Phase 6: Chrome Extension (weeks 13-14)

Browser-based access and page context awareness.

35. **Sidebar panel** -- chrome-ui sidebar with views
36. **Page context extraction** -- content script page parsing
37. **Desktop discovery + sync** -- find local desktop app, sync settings and theme
38. **Code preview in browser** -- sandboxed iframe preview
39. **Google Meet integration** -- meeting context injection

_Deliverable: Chrome extension for browser-based project management and AI chat._

### Phase 7: Integrations (weeks 15-20)

External service connectivity, prioritized by value.

40. **GitHub** -- highest value: issues, PRs, CI/CD (20 MCP tools)
41. **Orchestra Web auth + cloud sync** -- team features require web backend
42. **Discord bot** -- AI chat + MCP tools from Discord
43. **Slack bot** -- same capabilities as Discord
44. **Figma** -- design-to-code workflow (4 MCP tools)
45. **Jira** -- enterprise issue tracker sync
46. **Linear** -- modern issue tracker sync
47. **Notion** -- note push
48. **Apple Notes** -- macOS-only note push
49. **Firebase** -- push notifications, analytics, crash reporting

_Deliverable: Connected IDE with bidirectional sync to all major development tools._

---

## Tool Count by Feature Area

| Feature Area | MCP Tools | Settings Routes | WS Events | gRPC Methods | Total Touchpoints |
|-------------|-----------|-----------------|-----------|-------------|-------------------|
| Project Management | 7 | 2 | 0 | 0 | 9 |
| Epic/Story/Task CRUD | 15 | 2 | 0 | 0 | 17 |
| Workflow & Lifecycle | 7 | 2 | 1 | 0 | 10 |
| Sprint & Scrum | 11 | 2 | 0 | 0 | 13 |
| WIP & Dependencies | 6 | 2 | 0 | 0 | 8 |
| Metadata & Templates | 10 | 2 | 0 | 0 | 12 |
| PRD Authoring | 15 | 2 | 0 | 0 | 17 |
| Notes | 5 | 2 | 2 | 0 | 9 |
| Bug & Artifacts | 4 | 2 | 0 | 0 | 6 |
| AI Chat | 0 | 13 | 7 | 0 | 20 |
| Memory & Context | 6 | 2 | 0 | 12 | 20 |
| Code Search & Parse | 0 | 1 | 0 | 8 | 9 |
| Component & Preview | 7 | 6 | 5 | 3 | 21 |
| Claude Awareness | 7 | 0 | 0 | 0 | 7 |
| Desktop Windows | 1 | 9 | 0 | 0 | 10 |
| Notifications | 2 | 2 | 1 | 0 | 5 |
| DevTools | 25 | 4 | 1 | 0 | 30 |
| Code Editing | 0 | 0 | 1 | 1 | 2 |
| GitHub | 20 | 5 | 0 | 0 | 25 |
| Figma | 4 | 11 | 0 | 0 | 15 |
| Jira | 0 | 4 | 0 | 0 | 4 |
| Linear | 0 | 4 | 0 | 0 | 4 |
| Notion | 0 | 6 | 0 | 0 | 6 |
| Discord | 0 | 2 | 0 | 0 | 2 |
| Slack | 0 | 2 | 0 | 0 | 2 |
| Teams | 6 | 11 | 0 | 0 | 17 |
| Usage & Analytics | 8 | 0 | 0 | 0 | 8 |
| Auth & Sync | 0 | 9 | 4 | 0 | 13 |
| Voice | 0 | 4 | 1 | 0 | 5 |
| Screenshot | 0 | 3 | 0 | 0 | 3 |
| Chrome Extension | 0 | 1 | 3 | 0 | 4 |
| Theming | 0 | 1 | 0 | 0 | 1 |
| Workspace | 0 | 5 | 1 | 0 | 6 |
| Files & Permissions | 0 | 6 | 0 | 0 | 6 |
| Firebase | 0 | 0 | 0 | 0 | 0 (SDK only) |
| **TOTAL** | **~166** | **~134** | **~27** | **~24** | **~351** |

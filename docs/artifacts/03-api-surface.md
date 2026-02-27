# API Surface — Orchestra Reference

> Extracted from `orch-ref/`. Every endpoint across all four protocols: HTTP REST, WebSocket, gRPC, and MCP stdio.

---

## A. Settings API (HTTP REST on `:19191`)

All endpoints are served by the Go `settings.Server` on `127.0.0.1:19191`. CORS is enabled for all origins. All responses are `application/json` unless noted otherwise.

**Total registered routes: 134**

---

### A.1 Settings CRUD

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/settings` | - | `{key: value, ...}` | Get all settings as key-value map |
| POST | `/api/settings` | `{key, value}` | `{"ok":true}` | Set a single setting |
| GET | `/api/settings/{key}` | - | JSON value or `null` | Get a single setting by key |

---

### A.2 Window Management

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/windows/open` | `{name, title, route, width?, height?, data?, always_on_top?}` | `{"ok":true}` | Open a desktop window (proxies to desktop app if opener is nil) |
| GET | `/api/windows/data/{name}` | - | `{...}` or `{}` | Get stored data for a named window |
| POST | `/api/windows/close/{name}` | - | `{"ok":true}` | Close a desktop window by name |

---

### A.3 Spirit Window

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/spirit/toggle` | - | `{"ok":true}` | Toggle the spirit (floating assistant) window |
| POST | `/api/spirit/open` | - | `{"ok":true}` | Open the spirit window |
| POST | `/api/spirit/close` | - | `{"ok":true}` | Close the spirit window |

---

### A.4 Bubble Window

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/bubble/toggle` | - | `{"ok":true}` | Toggle the bubble window |
| POST | `/api/bubble/open` | - | `{"ok":true}` | Open the bubble window |
| POST | `/api/bubble/close` | - | `{"ok":true}` | Close the bubble window |

---

### A.5 Window Mode

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/mode` | - | `{"mode":"embedded"|"floating"|"bubble"}` | Get current window mode |
| POST | `/api/mode` | `{"mode":"embedded"|"floating"|"bubble"}` | `{"ok":true}` | Set window mode |
| POST | `/api/mode/cycle` | - | `{"ok":true,"mode":"..."}` | Cycle through window modes |

---

### A.6 Prompts (Permission / Question UI)

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/prompts/request` | `{id, type, tool_name?, tool_input?, question?, options?, session_id?, timeout_secs?}` | `{"ok":true,"id":"...","desktop":bool}` | Submit a prompt request; opens a desktop window for user interaction |
| GET | `/api/prompts/pending/{id}` | - | `PromptRequest` object | Get a pending prompt's details |
| POST | `/api/prompts/response/{id}` | `{action, answer?, tool_input?}` | `{"ok":true}` | Submit user's response to a prompt |
| GET | `/api/prompts/response/{id}` | Query: `?wait=true&timeout=120` | `PromptResponse` | Long-poll for a prompt response (blocks up to `timeout` seconds) |

---

### A.7 AI Chat

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/ai/send` | `{prompt, model?, mode?, workspace?, session?, page_context?}` | `{"ok":true,"session":"..."}` | Start a new Claude session (dispatches to AI bridge) |
| POST | `/api/ai/followup` | `{prompt, session?}` | `{"ok":true}` | Send a follow-up message to an existing session |
| POST | `/api/ai/stop` | `{session?}` | `{"ok":true}` | Cancel/stop an active AI session |
| POST | `/api/ai/permission` | `{request_id, approved, tool_input?, session?}` | `{"ok":true}` | Send a permission response to the AI bridge |
| GET | `/api/ai/events` | Query: `?session=desktop` | SSE stream: `data: {type,content,session}` | Subscribe to AI response chunks via Server-Sent Events |
| GET | `/api/ai/status` | - | `{sessions:[...], running:bool, prompt?:string}` | Get status of all AI sessions |
| GET | `/api/ai/sessions` | - | `[SessionInfo, ...]` | List all registered AI sessions |
| POST | `/api/ai/sessions/new` | - | `{"ok":true,"session":"s-..."}` | Create a new AI session |
| POST | `/api/ai/sessions/rename` | `{id, name}` | `{"ok":true}` | Rename an AI session |
| DELETE | `/api/ai/sessions/{id}` | - | `{"ok":true}` | Delete an AI session |
| GET | `/api/ai/providers` | - | `[ProviderStatus, ...]` | List all AI providers with config status |
| POST | `/api/ai/providers/set` | `{provider_id, api_key?, base_url?, enabled?}` | `{"ok":true}` | Save provider configuration (API key, base URL, enabled flag) |
| GET | `/api/ai/models` | - | `[Model, ...]` | List available models from all enabled providers |

---

### A.8 File Download / URL Open

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/download` | `{filename, content, encoding?:"base64"|"text"}` | `{"ok":true,"path":"..."}` or `{"ok":false,"cancelled":true}` | Show native Save As dialog and write file |
| POST | `/api/open-url` | `{url}` | `{"ok":true}` | Open a URL in the system default browser |
| POST | `/api/dialog/open-file` | `{message?, filters?:[{DisplayName,Pattern}]}` | `{"path":"..."}` | Open a native file picker dialog |

---

### A.9 Notifications

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/notify` | `{title?, body?, subtitle?, sound?}` | `{"ok":true}` | Send a desktop notification (and optionally trigger TTS) |
| POST | `/api/notify/speak` | `{title?, body?, subtitle?, sound?}` | `{"ok":true}` | Fire TTS speak event without sending OS notification (used by desktop proxy) |

---

### A.10 Discord Integration

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/discord/config` | - | `{bot_configured:bool, guild_configured:bool}` | Check Discord credential configuration |
| POST | `/api/discord/test` | - | `{"ok":true}` | Send a test message to Discord |

---

### A.11 Slack Integration

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/slack/config` | - | `{bot_configured:bool, socket_mode_configured:bool}` | Check Slack credential configuration |
| POST | `/api/slack/test` | - | `{"ok":true}` | Send a test message to Slack |

---

### A.12 Auth (Orchestra Web Backend)

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/auth/login` | `{email, password}` | `{"ok":true,"user":{...}}` | Login with email/password |
| GET | `/api/auth/status` | - | `{"ok":true,"state":{...}}` | Get current auth state |
| POST | `/api/auth/logout` | - | `{"ok":true,"message":"Logged out"}` | Sign out and clear credentials |
| POST | `/api/auth/store` | `{token, user}` | `{"ok":true}` | Store token + user from direct API login |
| GET | `/api/auth/token` | - | `{"ok":true,"token":"..."}` | Get stored bearer token |

---

### A.13 GitHub Integration

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/github/auth/start` | - | `{"ok":true,"url":"..."}` | Get GitHub OAuth authorization URL |
| GET | `/api/github/auth/callback` | Query: `?code=...` | HTML success/error page | Handle OAuth redirect, exchange code for token |
| GET | `/api/github/auth/status` | - | `{"ok":true,"state":{status,user?}}` | Check GitHub auth state |
| POST | `/api/github/auth/pat` | `{token}` | `{"ok":true,"user":{...}}` | Sign in with Personal Access Token |
| POST | `/api/github/auth/disconnect` | - | `{"ok":true,"message":"GitHub account disconnected"}` | Sign out from GitHub |

---

### A.14 Jira Integration

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/jira/auth/start` | - | `{"ok":true,"url":"..."}` | Get Jira OAuth authorization URL |
| GET | `/api/jira/auth/callback` | Query: `?code=...` | HTML success/error page | Handle OAuth redirect |
| GET | `/api/jira/auth/status` | - | `{"ok":true,"state":{status,user?}}` | Check Jira auth state |
| POST | `/api/jira/auth/disconnect` | - | `{"ok":true,"message":"Jira account disconnected"}` | Disconnect Jira |

---

### A.15 Linear Integration

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/linear/auth/start` | - | `{"ok":true,"url":"..."}` | Get Linear OAuth authorization URL |
| GET | `/api/linear/auth/callback` | Query: `?code=...` | HTML success/error page | Handle OAuth redirect |
| GET | `/api/linear/auth/status` | - | `{"ok":true,"state":{status,user?}}` | Check Linear auth state |
| POST | `/api/linear/auth/disconnect` | - | `{"ok":true,"message":"Linear account disconnected"}` | Disconnect Linear |

---

### A.16 Notion Integration

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/notion/config` | - | `{configured:bool}` | Check if Notion OAuth credentials are configured |
| GET | `/api/notion/auth/start` | - | `{"ok":true,"url":"..."}` | Get Notion OAuth authorization URL |
| GET | `/api/notion/auth/callback` | Query: `?code=...` | HTML success/error page | Handle OAuth redirect |
| GET | `/api/notion/auth/status` | - | `{"ok":true,"state":{status,workspace_name?}}` | Check Notion auth state |
| POST | `/api/notion/auth/disconnect` | - | `{"ok":true,"message":"Notion account disconnected"}` | Disconnect Notion |
| GET | `/api/notion/databases` | - | `{"ok":true,"databases":[{id,name}]}` | List accessible Notion databases |
| POST | `/api/notion/send` | `{title, content, database_id}` | `{"ok":true,"page_id":"..."}` | Send a note to Notion as a new page |

---

### A.17 Figma Integration

#### Auth

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/figma/auth/start` | - | `{"ok":true,"url":"..."}` | Get Figma OAuth authorization URL |
| GET | `/api/figma/auth/callback` | Query: `?code=...` | HTML success/error page | Handle OAuth redirect |
| GET | `/api/figma/auth/status` | - | `{"ok":true,"state":{status,user_id?,user_name?}}` | Check Figma auth state |
| POST | `/api/figma/auth/disconnect` | - | `{"ok":true,"message":"Figma account disconnected"}` | Disconnect Figma |
| POST | `/api/figma/auth/pat` | `{token}` | `{"ok":true,"user_name":"...","user_email":"..."}` | Save Personal Access Token (validates against Figma API) |

#### Files & Proxy

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/figma/files` | - | `{"ok":true,"files":[...]}` | List user's Figma files (cached) |
| GET | `/api/figma/files/{key}` | - | `{"ok":true,"file":{...}}` | Get a single Figma file by key |
| GET | `/api/figma/files/{key}/nodes` | Query: `?ids=nodeId` | `{"ok":true,"nodes":{...}}` | Get specific nodes from a Figma file |
| GET | `/api/figma/files/{key}/components` | - | `{"ok":true,"components":[...]}` | Get components from a Figma file |

#### MCP Server

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/figma/mcp/status` | - | `{"ok":true,"installed":bool,"running":bool}` | Check if figma-developer-mcp is installed in `.mcp.json` and running |
| POST | `/api/figma/mcp/install` | - | `{"ok":true,"installed":bool,"message":"..."}` | Add figma-developer-mcp entry to `.mcp.json` (idempotent) |

---

### A.18 Apple Notes (macOS only)

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/applenotes/status` | - | `{"ok":true,"available":bool}` | Check if Apple Notes is available |
| GET | `/api/applenotes/folders` | - | `{"ok":true,"folders":[...]}` | List Apple Notes folders |
| POST | `/api/applenotes/send` | `{title, content, folder?}` | `{"ok":true}` | Create a new Apple Note |

---

### A.19 Workspace

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/workspace` | - | `{path, name}` | Get current workspace path and name |
| POST | `/api/workspace/open` | `{path}` | `{path, name}` | Set workspace (reloads MCP bridge, notifies listeners) |
| GET | `/api/workspace/recent` | - | `{recent:["path1","path2",...]}` | List recent workspaces (up to 10) |
| GET | `/api/workspace/active-file` | - | `{path, language, cursor}` | Get active file context |
| POST | `/api/workspace/active-file` | `{path, language, cursor}` | `{path, language, cursor}` | Set active file context |
| GET | `/api/workspace/current-task` | - | `{id, title, status, ...}` or `{}` | Get cached current in-progress task |
| POST | `/api/workspace/current-task/refresh` | - | `{id, title, status, epic_id, story_id, ...}` or `{}` | Scan for in-progress task and cache with parent context |

---

### A.20 MCP Bridge

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/mcp/tools` | - | `{"tools":[ToolDefinition,...],"count":N}` | List all registered MCP tools |
| POST | `/api/mcp/tools/call` | `{name, arguments}` | `ToolResult` | Execute an MCP tool by name |

---

### A.21 Engine Storage

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/engine/storage` | - | `{"path":"..."}` | Get current engine storage path setting |
| POST | `/api/engine/storage` | `{path}` | `{"path":"..."}` | Set engine storage path (creates directory if needed) |

---

### A.22 DevTools

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/devtools/run-query` | `{session_id, query, params?}` | Query result JSON | Run a SQL query against an active database session |
| POST | `/api/devtools/ssh-exec` | `{session_id, command}` | Command output JSON | Execute a command on an active SSH session |
| POST | `/api/devtools/dispatch` | `{session_id, msg_type, payload?}` | Provider response JSON | Route a generic message to any DevTools session provider |
| POST | `/api/devtools/terminal-exec` | `{session_id, command}` | `{output, exit_code}` | Run a one-shot command in a terminal session's shell |

---

### A.23 Search

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/search/mentions` | `{query, limit?}` | `{"items":[{id,label,description,icon,iconColor,group}]}` | Search files and tasks for @-mention suggestions (uses Rust engine or glob fallback) |

---

### A.24 Screenshot

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/screenshot/capture` | - | `{"ok":true,"image":"base64..."}` | Capture full screen as base64 PNG |
| POST | `/api/screenshot/region` | `{x, y, width, height}` | `{"ok":true,"image":"base64..."}` | Capture and crop a screen region |
| POST | `/api/screenshot/start` | - | `{"ok":true,"image":"base64..."}` or `{"ok":false,"cancelled":true}` | Launch native interactive screenshot tool (macOS `screencapture -i -s -x`), blocks until user selects region |

---

### A.25 Permissions (macOS)

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/permissions/status` | - | `{notifications, screen_recording, microphone, accessibility}` | Get macOS permission status for each category |
| POST | `/api/permissions/request` | `{kind?}` | `{"ok":true}` | Trigger OS permission prompt (kind: `notifications`, `screen_recording`, `microphone`, `accessibility`) |
| POST | `/api/permissions/open-settings` | `{kind?}` | `{"ok":true}` | Open System Settings to the relevant Privacy pane |

---

### A.26 Voice

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/api/voice/analyze-transcript` | `{transcript, platform?, title?}` | `{summary, requirements[], action_items[], decisions[], suggested_project_title}` | Analyze meeting transcript using Anthropic API |
| POST | `/api/voice/create-project-from-meeting` | `{title, summary?, requirements?, action_items?, decisions?}` | `{"ok":true,"message":"Project creation started"}` | Dispatch project creation from meeting analysis to AI bridge |
| GET | `/api/voice/settings` | - | `{voice.stt.enabled, voice.tts.enabled, ...}` | Get all voice settings with defaults |
| PUT | `/api/voice/settings` | `{voice.stt.enabled:"true", ...}` | `{"ok":true}` | Save voice settings (only known keys accepted) |

---

### A.27 Sync & Migration

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/sync/migration/status` | - | `{running:bool, complete:bool, progress:{...}}` | Get current seed migration progress |
| POST | `/api/sync/migration/start` | - | `{running:true, progress:{...}}` (202 Accepted) | Start a fresh seed push (settings + chat sessions to web backend) |
| GET | `/api/sync/web/status` | - | `{configured:bool, authenticated:bool, ...}` | Get web sync client connection status |
| POST | `/api/integrations/sync` | - | `{"ok":true,"results":[{provider,action,error?}]}` | Push all integration connection states to web backend |

---

### A.28 Browser Context

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/browser/context` | - | `{"ok":true,"context":{title,url,mainContent,...}}` | Get latest page context from Chrome extension |

---

### A.29 Teams (Proxy to Laravel Web Backend)

All team endpoints proxy to the Laravel web API with the stored auth token.

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/teams` | - | `[Team, ...]` | List all teams |
| POST | `/api/teams` | `{name, ...}` | `Team` | Create a new team |
| GET | `/api/teams/pending-invitations` | - | `[Invitation, ...]` | List pending team invitations |
| GET | `/api/teams/{id}` | - | `Team` | Get a single team |
| PUT | `/api/teams/{id}` | `{name, ...}` | `Team` | Update a team |
| DELETE | `/api/teams/{id}` | - | `204` | Delete a team |
| GET | `/api/teams/{id}/members` | - | `[Member, ...]` | List team members |
| PUT | `/api/teams/{id}/members/{userId}` | `{role, ...}` | `Member` | Update a team member |
| DELETE | `/api/teams/{id}/members/{userId}` | - | `204` | Remove a team member |
| POST | `/api/teams/{id}/invitations` | `{email, role?}` | `Invitation` | Invite a member to a team |
| GET | `/api/teams/{id}/shares` | - | `[Share, ...]` | List team shares |
| POST | `/api/teams/{id}/shares` | `{...}` | `Share` | Create a team share |

---

### A.30 Component Library

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| GET | `/api/v1/components` | Query: `?user_id=&framework=&name=&tags=&limit=&offset=` | `{"data":[Component,...],"count":N}` | List components with filters |
| POST | `/api/v1/components` | `{user_id?, name, framework, description, html?, css?, js?, jsx?, tags?}` | `Component` (201 Created) | Create a new component |
| GET | `/api/v1/components/search` | Query: `?q=&framework=&limit=` | `{"data":[Component,...],"count":N}` | Search components by name |
| GET | `/api/v1/components/{id}` | Query: `?user_id=` | `ComponentDetail` (includes html, css, js, jsx) | Get a single component with code |
| PUT | `/api/v1/components/{id}` | `{user_id?, name, framework, description, html?, css?, js?, jsx?, tags?}` | `Component` | Update a component |
| DELETE | `/api/v1/components/{id}` | Query: `?user_id=` | `204 No Content` | Delete a component |

---

### A.31 WebSocket Endpoints (Registered on Settings Mux)

These are WebSocket upgrade endpoints registered on the same `:19191` HTTP server.

| Path | Description |
|------|-------------|
| `/ws/lsp/{language}` | LSP proxy WebSocket (spawns language servers with workspace root) |
| `/ws/devtools/*` | DevTools WebSocket routes (registered by `devtools.RegisterWSRoutes`) |

---

## B. WebSocket Protocol (port `:8765`)

The WebSocket server runs on `127.0.0.1:8765` using Fiber v3 with `fasthttp/websocket`. Clients connect at `ws://127.0.0.1:8765/` with optional query params `?user_id=&session_id=`.

### B.1 Health Check

| Method | Path | Response | Description |
|--------|------|----------|-------------|
| GET | `/health` | `{"status":"ok","time":unix}` | HTTP health check endpoint |

### B.2 Message Format

All WebSocket messages use this JSON envelope:

```json
{
  "id": "string",
  "type": "event | command | response | error | ping | pong",
  "event": "string (when type=event)",
  "command": "string (when type=command)",
  "data": {},
  "timestamp": 1234567890,
  "request_id": "string (for response/pong/error)",
  "error": { "code": "string", "message": "string", "details": "string" }
}
```

### B.3 Message Types

| Type | Direction | Purpose |
|------|-----------|---------|
| `event` | Server -> Client (or Client -> Server) | Named event broadcast |
| `command` | Client -> Server | Named command invocation |
| `response` | Server -> Client | Response to a command (includes `request_id`) |
| `error` | Server -> Client | Error response (includes `request_id`, `error.code`, `error.message`) |
| `ping` | Client -> Server | Keepalive ping |
| `pong` | Server -> Client | Keepalive pong (includes `request_id` from ping) |

---

### B.4 Server -> Client Events

| Event | Payload | Source | Description |
|-------|---------|--------|-------------|
| `ai:chunk` | `{type, content, session}` | AIBridgeAdapter | AI response chunk (text, tool_use, done, error, etc.) |
| `ai:session_created` | `{session}` | AIBridgeAdapter | New AI session was created |
| `ai:message_forwarded` | `{id, target_session, content, role, source_session, source_message_id, forwarded_at}` | AIBridgeAdapter | Message forwarded between sessions |
| `workspace:changed` | `{path}` | SettingsServer | Workspace directory changed |
| `tasks:changed` | `{project, epic_id, story_id, task_id, type, from, to, time}` | WorkflowBroadcaster | Task state transition occurred |
| `page:context_update` | `{title, url, mainContent, headings, isArticle}` | BrowserHandler | Chrome extension page context updated |
| `page:request_context` | `null` | BrowserHandler | Request Chrome extension to send fresh page context |
| `preview:update` | `{session_id, code}` | PreviewCoordinator | Preview code updated |
| `preview:viewport` | `{session_id, viewport}` | PreviewCoordinator | Preview viewport changed |
| `preview:session_ended` | `{session_id}` | PreviewCoordinator | Preview session deleted |
| `preview:open_browser` | `{session_id}` | PreviewCoordinator | Request Chrome extension to open preview URL |
| `sync:note:upsert` | `{entity_type, entity_id, action, payload, version}` | SyncHandler | Note created/updated (local + pull) |
| `sync:note:delete` | `{entity_type, entity_id, action, payload, version}` | SyncHandler | Note deleted |
| `sync:project:upsert` | `{entity_type, entity_id, action, payload, version}` | SyncHandler | Project created/updated |
| `sync:project:delete` | `{entity_type, entity_id, action, payload, version}` | SyncHandler | Project deleted |
| `sync:integration:upsert` | `{entity_type, entity_id, action, payload, version}` | SyncHandler | Integration connected |
| `sync:integration:delete` | `{entity_type, entity_id, action, payload, version}` | SyncHandler | Integration disconnected |
| `sync:ai_session:upsert` | `{entity_type, entity_id, action, payload, version}` | SyncHandler | AI session created/updated |
| `sync:ai_session:delete` | `{entity_type, entity_id, action, payload, version}` | SyncHandler | AI session deleted |

---

### B.5 Client -> Server Commands

#### AI Commands

| Command | Payload | Response | Description |
|---------|---------|----------|-------------|
| `ai:send` | `{prompt, model?, mode?, workspace?, session?}` | - (events via `ai:chunk`) | Start an AI session |
| `ai:stop` | `{session?}` | - | Stop an active AI session |
| `ai:permission` | `{request_id, approved, tool_input?, session?}` | - | Respond to a tool permission request |
| `ai:question` | `{request_id, answer, session?}` | - | Respond to a question prompt |
| `ai:status` | `{}` | `{sessions:[SessionInfo]}` | Get all session statuses |
| `ai:sessions` | `{}` | `{sessions:[SessionInfo]}` | List all sessions |
| `ai:forward` | `{target_session, content, role?, source_session?, source_message_id?}` | `{id, target_session, forwarded_at}` | Forward a message between sessions |

#### Browser Commands

| Command | Payload | Response | Description |
|---------|---------|----------|-------------|
| `page:request_context` | `{}` | - (broadcasts to Chrome extension) | Request fresh page context from Chrome |

#### Browser Events (Client -> Server)

| Event | Payload | Description |
|-------|---------|-------------|
| `page.updated` | `{content:{url,title,mainContent,headings,isArticle}}` | Chrome extension sends updated page context |

#### Preview Commands

| Command | Payload | Response | Description |
|---------|---------|----------|-------------|
| `preview:join` | `{session_id}` | `{session_id, session}` | Join a preview session (receive current state) |
| `preview:leave` | `{session_id}` | - | Leave a preview session |
| `preview:update` | `{session_id, html?, css?, js?, jsx?, framework?}` | - (broadcast `preview:update`) | Update preview code |
| `preview:viewport` | `{session_id, preset?, width?, height?}` | - (broadcast `preview:viewport`) | Change preview viewport |

#### Sync Commands

| Command | Payload | Response | Description |
|---------|---------|----------|-------------|
| `sync.note.upsert` | `{id, title, content, tags, pinned, workspace?, version}` | `{entity_type, entity_id, action, queued:true}` | Create/update a note (pushes to web API, broadcasts to other clients) |
| `sync.note.delete` | `{id, version}` | `{entity_type, entity_id, action, queued:true}` | Delete a note |
| `sync.project.upsert` | `{id, name, slug?, description?, path?, stats?, meta?, workspace?, version}` | `{entity_type, entity_id, action, queued:true}` | Create/update a project |
| `sync.project.delete` | `{id, version}` | `{entity_type, entity_id, action, queued:true}` | Delete a project |
| `sync.integration.upsert` | `{id, provider, access_token, ..., version}` | `{entity_type, entity_id, action, queued:true}` | Create/update an integration |
| `sync.integration.delete` | `{id, provider, version}` | `{entity_type, entity_id, action, queued:true}` | Delete an integration |
| `sync.ai_session.upsert` | `{id, name, model, workspace?, pinned?, icon?, color?, messages?, version}` | `{entity_type, entity_id, action, queued:true}` | Create/update an AI session |
| `sync.ai_session.delete` | `{id, version}` | `{entity_type, entity_id, action, queued:true}` | Delete an AI session |

#### Settings Commands

| Command | Payload | Response | Description |
|---------|---------|----------|-------------|
| `settings.set` | `{key, value}` | `{key, saved:true}` | Save a setting locally + push to web API |
| `settings.get` | `{key}` | `{key, value}` | Get a setting from local store |
| `settings.get_all` | `{}` | `{settings:{key:value,...}}` | Get all settings |

---

### B.6 Preview HTTP Endpoints (on `:8765`)

These REST endpoints are registered on the Fiber app alongside the WebSocket upgrade.

| Method | Path | Request Body | Response | Description |
|--------|------|-------------|----------|-------------|
| POST | `/preview` | `{framework, html?, css?, js?, jsx?}` | `{session_id, ws_url, framework}` (201) | Create a new preview session |
| GET | `/preview/{session_id}` | - | `PreviewSession` | Get preview session state |
| DELETE | `/preview/{session_id}` | - | `{session_id, deleted:true}` | Delete a preview session |

---

## C. gRPC Services (port `:50051`)

The Rust engine exposes gRPC services on `127.0.0.1:50051`. Protocol version: proto3. Package: `orchestra.engine.v1`.

---

### C.1 HealthService

| Method | Request | Response | Stream | Description |
|--------|---------|----------|--------|-------------|
| `Check` | `HealthCheckRequest{service?}` | `HealthCheckResponse{status, message, timestamp}` | Unary | Check engine health |
| `Watch` | `HealthCheckRequest{service?}` | `stream HealthCheckResponse` | Server-stream | Watch health status updates |

**ServingStatus enum:** `UNKNOWN(0)`, `SERVING(1)`, `NOT_SERVING(2)`, `SERVICE_UNKNOWN(3)`

---

### C.2 ParseService

| Method | Request | Response | Stream | Description |
|--------|---------|----------|--------|-------------|
| `ParseFile` | `ParseFileRequest{file_path, content, language, include_ast?}` | `ParseFileResponse{file_path, success, error, ast, symbols[], parse_time_ms}` | Unary | Parse a single file using Tree-sitter |
| `ParseFiles` | `stream ParseFileRequest` | `stream ParseFileResponse` | Bidi-stream | Parse multiple files in batch |
| `GetSymbols` | `GetSymbolsRequest{file_path, content, language, symbol_types[]}` | `GetSymbolsResponse{file_path, symbols[]}` | Unary | Extract code symbols (functions, classes, variables) |

**Symbol:** `{name, kind, range{start_line, start_column, end_line, end_column}, detail, children[]}`

---

### C.3 SearchService

| Method | Request | Response | Stream | Description |
|--------|---------|----------|--------|-------------|
| `IndexFile` | `IndexFileRequest{file_path, content, language, metadata{}}` | `IndexFileResponse{success, error, indexed_count}` | Unary | Index a single file in Tantivy |
| `IndexFiles` | `stream IndexFileRequest` | `IndexFileResponse` | Client-stream | Batch index multiple files |
| `Search` | `SearchRequest{query, limit, offset, file_types[], fuzzy}` | `SearchResponse{results[], total_hits, search_time_ms}` | Unary | Full-text search indexed files |
| `DeleteFile` | `DeleteFileRequest{file_path}` | `DeleteFileResponse{success, error}` | Unary | Remove a file from the index |
| `ClearIndex` | `ClearIndexRequest{}` | `ClearIndexResponse{success, error}` | Unary | Clear all indexed files |

**SearchResult:** `{file_path, score, snippets[], line_number, metadata{}}`

---

### C.4 MemoryService

#### Session Management

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `StartSession` | `{id, project, agent_type, model, metadata{}}` | `{session}` | Start a new observation session |
| `EndSession` | `{session_id, ended_at, project}` | `{session}` | End an observation session |
| `GetSession` | `{project, session_id}` | `{session}` | Get a session by ID |
| `ListSessions` | `{project, limit}` | `{sessions[]}` | List sessions for a project |

#### Observation Recording

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `RecordObservation` | `{session_id, project, observation_type, content, tool_name?, tool_input?, tool_output?, context?, tokens?, metadata{}}` | `{observation}` | Record an observation (user_prompt, tool_use, tool_result, assistant_response) |
| `GetTimeline` | `{session_id, around_sequence, radius}` | `{observations[]}` | Get timeline around a sequence number |
| `GetObservations` | `{session_id, project}` | `{observations[]}` | Get all observations for a session |

#### Memory Search

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `SearchMemory` | `{project, query, limit}` | `{results[]}` | Search project memory |
| `FetchDetails` | `{observation_ids[]}` | `{observations[]}` | Fetch full details for specific observations |

#### Embeddings

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `StoreEmbedding` | `{entity_type, entity_id, project, model, vector[]}` | `{success}` | Store a vector embedding |
| `SearchSimilar` | `{project, query_vector[], model, limit}` | `{entities[{entity_type, entity_id, similarity}]}` | Search by vector similarity |

#### Summaries

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `StoreSummary` | `{session_id, project, summary_type, content, observation_ids[], tokens?, metadata{}}` | `{summary}` | Store a session summary |
| `GetSummaries` | `{session_id, project}` | `{summaries[]}` | Get summaries for a session |

#### Context Retrieval

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `GetContext` | `{project, query, limit}` | `{chunks[]}` | Get relevant context chunks for a query |

---

### C.5 ComponentBundlerService

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `BundleComponent` | `{framework, html, css, js, jsx, dependencies[], content_hash?}` | `{success, error, bundled_html, inlined_css, inlined_js, content_hash, bundle_time_ms}` | Bundle a UI component into a self-contained HTML document |
| `ParseComponentProps` | `{source, framework}` | `{success, error, props[{name, type_name, required, default_value, description}], component_name}` | Extract props interface from component source |
| `ValidateComponent` | `{source, framework}` | `{valid, errors[{message, line, column, severity}], warnings[]}` | Validate component for syntax errors |

**Supported frameworks:** `html`, `react`, `vue`, `svelte`, `angular`, `react-native`, `flutter`

---

## D. MCP Protocol (stdio JSON-RPC)

The MCP server communicates via stdin/stdout JSON-RPC 2.0. It supports two transports: **stdio** (primary, for CLI integration) and **SSE** (for browser/web clients).

### D.1 Protocol

- **JSON-RPC version:** 2.0
- **MCP protocol version:** `2024-11-05`
- **Max message size:** 10 MB
- **Server info:** `{name: "orchestra-mcp", version: "..."}`

### D.2 Initialization

```
Client -> Server:  {"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}
Server -> Client:  {"jsonrpc":"2.0","id":1,"result":{
  "protocolVersion": "2024-11-05",
  "capabilities": {
    "tools": {},
    "resources": {},  // if resources registered
    "prompts": {}     // if prompts registered
  },
  "serverInfo": {"name":"orchestra-mcp","version":"1.0.0"}
}}

Client -> Server:  {"jsonrpc":"2.0","method":"notifications/initialized"}
```

### D.3 Supported Methods

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `initialize` | `{}` | `{protocolVersion, capabilities, serverInfo}` | Initialize MCP connection |
| `notifications/initialized` | - | - (notification, no response) | Client acknowledges initialization |
| `tools/list` | `{}` | `{tools: [ToolDefinition, ...]}` | List all available tools |
| `tools/call` | `{name, arguments}` | `ToolResult{content:[{type,text}], isError?}` | Call a tool by name |
| `resources/list` | `{}` | `{resources: [ResourceDefinition, ...]}` | List available resources |
| `resources/read` | `{uri}` | `{contents: [ResourceContent, ...]}` | Read a resource by URI (supports template patterns) |
| `prompts/list` | `{}` | `{prompts: [PromptDefinition, ...]}` | List available prompts |
| `prompts/get` | `{name, arguments?}` | `{description, messages:[{role,content}]}` | Get a rendered prompt |
| `ping` | `{}` | `{}` | Keepalive ping |

### D.4 SSE Transport

The SSE transport provides the same MCP protocol over HTTP Server-Sent Events for browser-based clients.

- **Session management:** `SSESessionManager` tracks active sessions with UUID-based IDs
- **Message channel:** Each session has a buffered channel (32 messages) for outbound SSE events
- **Response writer:** `SSEWriter` serializes JSON-RPC responses and sends them as SSE data events
- **Lifecycle:** Sessions are created on connect, removed on disconnect or context cancellation

### D.5 Tool Categories (85+ tools)

The MCP bridge registers tools from these categories:

| Category | Source | Description |
|----------|--------|-------------|
| Project | `tools.Project()` | Create, list, get, update, delete projects |
| Epic | `tools.Epic()` | Create, list, get, update, delete epics |
| Story | `tools.Story()` | Create, list, get, update, delete stories |
| Task | `tools.Task()` | Create, list, get, update, delete tasks |
| Workflow | `tools.Workflow()` | Advance task state, set current task, get next task |
| PRD | `tools.Prd()` | PRD session management (start, answer, validate, preview, write) |
| PRD Templates | `tools.PrdTemplates()` | Save, list, create from reusable PRD templates |
| Bugfix | `tools.Bugfix()` | Report and track bugs |
| Usage | `tools.Usage()` | Record and get token/session usage metrics |
| Readme | `tools.Readme()` | Regenerate project README |
| Artifacts | `tools.Artifacts()` | Save planning artifacts |
| Lifecycle | `tools.Lifecycle()` | Project lifecycle management |
| Claude | `tools.Claude()` | Memory search/save, context retrieval (Rust engine bridge) |
| Sprint | `tools.Sprint()` | Create, start, end sprints; burndown, velocity, standup |
| Team | `tools.TeamTools()` | Create, get, list teams; invite, share |
| Scrum | `tools.Scrum()` | Retrospectives, standup summaries |
| Dependency | `tools.Dependency()` | Add, remove, view task dependencies |
| Metadata | `tools.Metadata()` | Labels, estimates, WIP limits |
| CLI Commands | `tools.CliCommands()` | DevTools: run-query, ssh-exec, terminal-exec, list-databases, list-services, control-service |
| Analytics | `tools.Analytics()` | Project analytics and metrics |
| Session Metrics | `tools.SessionMetricsTools()` | Session usage tracking |
| Notification | `tools.Notification()` | Send desktop notifications, play sounds |

---

## E. Shared Interfaces

### E.1 Settings SyncBridge

The `SyncBridge` prevents infinite sync loops between the desktop SQLite store and the Chrome extension WebSocket store. It provides two callbacks:

- `OnDesktopChange(key, value)` -- forwards desktop changes to Chrome store
- `OnChromeChange(key, value)` -- forwards Chrome changes to desktop store

Both callbacks use a re-entrancy guard (`syncing` flag) to prevent loops.

### E.2 MCP Bridge

The `McpBridge` provides HTTP REST access to all MCP tools without the stdio transport:

- **`ListTools()`** -- returns all tool definitions
- **`CallTool(name, args)`** -- executes a tool and returns the result
- **`ReloadWorkspace(path)`** -- rebuilds all tools for a new workspace
- **`Workspace()`** -- returns current workspace path

### E.3 Prompt Store

In-memory store for permission/question prompts with channel-based wait/respond:

- **Submit(req)** -- stores a prompt request
- **Get(id)** -- retrieves a pending prompt
- **Respond(id, resp)** -- sends a response, unblocks waiters
- **Wait(id, timeout)** -- blocks until response or timeout
- **Cleanup(id)** -- removes a prompt

---

## F. Port Summary

| Port | Protocol | Service | Description |
|------|----------|---------|-------------|
| `19191` | HTTP | Settings API | REST API + WebSocket (LSP proxy, DevTools) |
| `8765` | WebSocket + HTTP | WebSocket Server | Real-time events, sync, AI bridge, preview |
| `50051` | gRPC | Rust Engine | Tree-sitter, Tantivy, Memory, Component bundling |
| stdio | JSON-RPC | MCP Server | Claude Code MCP protocol (85+ tools) |

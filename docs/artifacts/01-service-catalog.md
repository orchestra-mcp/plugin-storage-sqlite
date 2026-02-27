# Service Catalog — Orchestra Reference

> Extracted from `orch-ref/`. Every runtime process, its protocol, port, and dependencies.

## Overview

| Service | Protocol | Port | Binary | Status |
|---------|----------|------|--------|--------|
| MCP Server | JSON-RPC 2.0 (stdio) | - | `bin/orchestra` | Production |
| Settings API | HTTP REST | 19191 (primary), 19192 (desktop proxy) | embedded | Production |
| WebSocket Server | WebSocket (Fiber v3) | 8765 | embedded | Production |
| Rust Engine | gRPC (Tonic) | 50051 | `bin/orchestra-engine` | Production |
| Desktop App | Wails v3 GUI | - | `Orchestra MCP.app` | Production |
| Discord Bot | Discord Gateway WebSocket + REST | - | embedded | Optional |
| Sync Client | HTTP polling | - | embedded | Background |

---

## 1. MCP Server (stdio)

### Overview

The default execution mode of the `bin/orchestra` binary. When invoked with no flags (`orch` or `orch --workspace <path>`), it runs a **JSON-RPC 2.0 server over stdin/stdout** implementing the Model Context Protocol.

- **Binary**: `bin/orchestra` (built from `cmd/orchestra/main.go`)
- **Entry point**: `commands.RunMCP(workspace)` in `app/commands/commands.go`
- **Protocol**: JSON-RPC 2.0, line-delimited, over stdin/stdout
- **Protocol version**: `2024-11-05`
- **Max message size**: 10 MB (`maxScanSize` in `app/transport/server.go`)

### Startup Sequence

1. Parse CLI flags (`--workspace`, `--dev`, `--build`, `--release`, `--version`, `init`)
2. Start Rust engine subprocess via `mcpengine.NewManager().Start(workspace, storagePath)` -- non-fatal if binary is missing; falls back to markdown/TOON storage
3. If engine is running, dial gRPC client via `mcpengine.Dial(mgr.Addr())`
4. Create `mcpengine.Bridge` (wraps gRPC client or nil for fallback)
5. Create `transport.MCPServer` with name `"orchestra"` and current version
6. Register all tool groups: Project, Epic, Story, Task, Workflow, Prd, Bugfix, Usage, Readme, Artifacts, Lifecycle, Claude, Memory, Desktop, Sprint, Team, Scrum, Dependency, DevTools, Metadata, Notes, Notification, Figma (conditional on token)
7. Register resources and prompts
8. Log startup banner to stderr with tool/resource/prompt counts and memory mode
9. Enter `s.Run()` -- blocking stdio read loop via `bufio.Scanner` on `os.Stdin`

### JSON-RPC Methods

| Method | Description |
|--------|-------------|
| `initialize` | Returns server capabilities (tools, resources, prompts) |
| `notifications/initialized` | No-op acknowledgment |
| `tools/list` | Returns all registered tool definitions |
| `tools/call` | Executes a tool by name with argument validation |
| `resources/list` | Returns all registered resource definitions |
| `resources/read` | Reads a specific resource by URI |
| `prompts/list` | Returns all registered prompt definitions |
| `prompts/get` | Returns a specific prompt by name |
| `ping` | Returns empty object |

### Response Writers

The transport layer abstracts response writing behind a `ResponseWriter` interface:

- **`StdioWriter`** (`app/transport/writer.go`) -- writes JSON-RPC responses to stdout, one line per response
- **`SSEWriter`** (`app/transport/sse.go`) -- writes JSON-RPC responses to an SSE session's message channel for browser-based clients

### SSE Transport

An alternative transport exists for browser-based MCP clients:

- **`SSESession`** (`app/transport/session.go`) -- represents a single SSE client connection with a buffered message channel (capacity 32)
- **`SSESessionManager`** -- tracks all active SSE sessions with thread-safe create/get/remove/count operations
- Sessions are identified by UUID and support cancellation via context

### Dependencies

- Rust engine binary (`orchestra-engine`) -- optional, gracefully degrades to markdown fallback
- Figma token store -- optional, Figma tools only registered if token exists in `~/Library/Application Support/Orchestra/`

### Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--workspace`, `-w` | `.` | Workspace directory path |
| `--version`, `-v` | - | Print version info |
| `init` | - | Initialize workspace (creates `.mcp.json`, `.projects/`, `.claude/`, etc.) |
| `--dev`, `-d` | - | Run all services in dev mode instead of MCP |
| `--build`, `-b` | - | Build all components |
| `--release`, `-r` | - | Tag and push a release (requires `--version=vX.Y.Z`) |

### Source Files

| File | Description |
|------|-------------|
| `cmd/orchestra/main.go` | CLI entry point -- flag parsing, mode dispatch |
| `app/commands/commands.go` | `RunMCP()`, `RunDev()`, `RunBuild()`, `RunRelease()`, `RunInit()` |
| `app/commands/mcp/mcp.go` | Cobra-based `mcp start` and `mcp init` subcommands |
| `app/transport/server.go` | `MCPServer` struct, `Run()` stdio loop, `HandleRequest()` dispatcher |
| `app/transport/writer.go` | `ResponseWriter` interface, `StdioWriter` implementation |
| `app/transport/sse.go` | `SSEWriter` for browser-based SSE transport |
| `app/transport/session.go` | `SSESession`, `SSESessionManager` for SSE client tracking |
| `app/transport/dynamic.go` | `UnregisterTool()`, `UnregisterResource()`, `UnregisterPrompt()` |
| `app/transport/prompts.go` | Prompt handling helpers |
| `app/transport/resources.go` | Resource handling helpers |

---

## 2. Settings API Server (HTTP)

### Overview

A local HTTP REST server backed by SQLite that persists user settings, manages desktop windows, proxies AI chat sessions, handles OAuth flows for integrations, and exposes MCP tools over HTTP.

- **Package**: `app/settings/`
- **Port**: `19191` (primary), `19192` (desktop proxy fallback)
- **Bind address**: `127.0.0.1` (localhost only)
- **Backing store**: SQLite database at `~/Library/Application Support/Orchestra/settings.db`
- **SQLite mode**: WAL journal, single connection, 5s busy timeout

### Startup

1. Open (or create) SQLite database at `<configDir>/settings.db`
2. Run schema migration: `CREATE TABLE IF NOT EXISTS settings (key TEXT PRIMARY KEY, value TEXT NOT NULL)`
3. Register all HTTP route handlers on `http.ServeMux`
4. Call `net.Listen("tcp", addr)` -- if port is already in use, returns error
5. Serve in a goroutine via `http.Server.Serve(ln)`

### Port Conflict Handling

When the desktop app starts and port 19191 is already in use (e.g., `orch --dev` is running):

1. The desktop detects `srv.Start()` failure
2. Closes its own settings store (reuses the existing server)
3. Starts a **secondary mini server on port 19192** with a subset of endpoints:
   - `POST /api/windows/open` -- proxy window open commands to the desktop app
   - `POST /api/windows/close/` -- proxy window close commands
   - `POST /api/notify` -- proxy notifications (forwards to 19191 for WebSocket TTS broadcast)
   - `POST /api/screenshot/start` -- open screenshot overlay

### API Routes (Partial List)

**Settings:**
- `GET /api/settings` -- all settings as JSON object
- `GET /api/settings/{key}` -- single setting value
- `POST /api/settings` -- set a setting (JSON body `{key, value}`)

**Windows:**
- `POST /api/windows/open` -- open a named window with route, size, data
- `POST /api/windows/close/{name}` -- close a named window
- `GET /api/windows/data/{name}` -- get window data

**AI Chat:**
- `POST /api/ai/send` -- send message to AI bridge
- `POST /api/ai/followup` -- follow-up message
- `POST /api/ai/stop` -- stop AI generation
- `POST /api/ai/permission` -- respond to AI permission request
- `GET /api/ai/events` -- SSE stream for AI events
- `GET /api/ai/status` -- AI bridge status
- `GET /api/ai/sessions` -- list chat sessions
- `POST /api/ai/sessions/new` -- create new session
- `GET /api/ai/providers` -- list AI providers
- `GET /api/ai/models` -- list available models

**MCP Tools (HTTP):**
- `GET /api/mcp/tools` -- list all MCP tool definitions
- `POST /api/mcp/tools/call` -- call an MCP tool by name

**OAuth Integrations:**
- GitHub: `GET /api/github/auth/start`, `GET /api/github/auth/callback`, `GET /api/github/auth/status`
- Jira: `GET /api/jira/auth/start`, `GET /api/jira/auth/callback`, `GET /api/jira/auth/status`
- Linear: `GET /api/linear/auth/start`, `GET /api/linear/auth/callback`, `GET /api/linear/auth/status`
- Notion: `GET /api/notion/auth/start`, `GET /api/notion/auth/callback`, `GET /api/notion/auth/status`
- Figma: `GET /api/figma/auth/start`, `GET /api/figma/auth/callback`, `GET /api/figma/auth/status`

**Auth:**
- `POST /api/auth/login` -- Orchestra web login
- `GET /api/auth/status` -- authentication status
- `POST /api/auth/logout` -- logout
- `GET /api/auth/token` -- current bearer token

**Other:**
- `POST /api/notify` -- send desktop notification (with optional sound)
- `POST /api/download` -- download a file
- `POST /api/open-url` -- open URL in browser
- Spirit window: `POST /api/spirit/toggle`, `/open`, `/close`
- Bubble window: `POST /api/bubble/toggle`, `/open`, `/close`
- Window mode: `GET /api/mode`, `POST /api/mode`, `POST /api/mode/cycle`

### Dependencies

- SQLite (via `modernc.org/sqlite` -- pure Go, no CGo)
- Optional: WindowOpener interface (provided by desktop app for window management)
- Optional: AI Bridge, Session Registry, Auth Hub, GitHub Hub, Notion Hub, Jira Hub, Linear Hub, Figma Auth, MCP Bridge, DevTools Manager, LSP Handler, Sync Client, Outbox Store

### Source Files

| File | Description |
|------|-------------|
| `app/settings/server.go` | `Server` struct, all HTTP handlers, `Start()`/`Stop()` |
| `app/settings/store.go` | `Store` -- SQLite CRUD for settings (Get/Set/GetAll/Close) |
| `app/settings/mcp_bridge.go` | `McpBridge` -- HTTP-accessible MCP tool registry with workspace reload |
| `app/settings/auth.go` | Orchestra web authentication (login/logout/token) |
| `app/settings/github.go` | GitHub OAuth and PAT management |
| `app/settings/figma_auth.go` | Figma OAuth flow |
| `app/settings/figma_files.go` | Figma file operations |
| `app/settings/figma_proxy.go` | Figma API proxy |
| `app/settings/figma_mcp.go` | Figma MCP tool integration |
| `app/settings/jira.go` | Jira OAuth flow |
| `app/settings/linear.go` | Linear OAuth flow |
| `app/settings/notion.go` | Notion OAuth and note sync |
| `app/settings/bridge.go` | AI bridge HTTP endpoints |
| `app/settings/workspace.go` | Workspace switching logic |
| `app/settings/engine.go` | Rust engine management endpoints |
| `app/settings/search.go` | Search endpoints |
| `app/settings/voice.go` | Voice/TTS endpoints |
| `app/settings/devtools.go` | DevTools session management |
| `app/settings/team.go` | Team management endpoints |
| `app/settings/integration_sync.go` | Integration entity sync |
| `app/settings/component_handlers.go` | UI component handlers |
| `app/settings/prompts.go` | Prompt request/response endpoints |
| `app/settings/screenshot_darwin.go` | macOS screenshot capture |
| `app/settings/applenotes.go` | Apple Notes integration |

---

## 3. WebSocket Server

### Overview

A Fiber v3 WebSocket server providing real-time bidirectional communication between the Go backend and all connected frontends (desktop, browser extension, dashboard).

- **Package**: `app/websocket/`
- **Framework**: Fiber v3 with `fasthttp/websocket` upgrader
- **Port**: `8765`
- **Bind address**: `127.0.0.1`
- **Timeouts**: Read 60s, Write 60s
- **CORS**: Enabled, all origins allowed

### Startup

1. Create `Server` with config (host, port, timeouts, CORS)
2. Create Fiber app with recovery, logger, and CORS middleware
3. Create `Service` (manages clients, message routing, broadcast)
4. Start service run loop in goroutine (register/unregister clients, broadcast, ping)
5. Register routes: `GET /health` (JSON health check), `GET /` (WebSocket upgrade)
6. Listen on `host:port` in goroutine

### Protocol

Messages are JSON with the following structure:

```json
{
  "id": "unique-id",
  "type": "event|command|response|error|ping|pong",
  "event": "event.name",
  "command": "command.name",
  "data": {},
  "timestamp": 1234567890,
  "request_id": "original-request-id",
  "error": {"code": "ERROR_CODE", "message": "..."}
}
```

**Message types:**
- `event` -- server-to-client event notifications (requires `event` field)
- `command` -- client-to-server commands (requires `command` field)
- `response` -- response to a command (requires `request_id`)
- `error` -- error response (requires `error` object)
- `ping`/`pong` -- keepalive

### Client Management

- Each client gets a UUID, a buffered send channel (capacity 256), and metadata
- Clients authenticate via `user_id` and `session_id` query parameters on upgrade
- Ping interval: 30s. Pong wait: 60s. Write wait: 10s
- Clients that don't respond to pings within the pong wait are disconnected
- Supports: `SendTo(clientID)`, `Broadcast(msg)`, `BroadcastToUser(userID)`

### Registered Handlers

Handlers are registered by the desktop bootstrap or dev mode:

| Handler | Package | Purpose |
|---------|---------|---------|
| `SyncHandler` | `app/websocket/sync_handler.go` | Entity sync (notes, projects, integrations, ai_sessions) |
| `AIBridgeAdapter` | `app/websocket/ai_bridge.go` | AI chat streaming over WebSocket |
| `PreviewCoordinator` | `app/websocket/preview_coordinator.go` | Live preview coordination between editor and preview |
| `BrowserHandler` | `app/websocket/browser_handler.go` | Chrome extension page context awareness |
| `WorkflowBroadcaster` | `app/websocket/workflow_listener.go` | Broadcasts workflow state transitions |
| `SettingsHandler` | `app/websocket/settings_handler.go` | Settings change notifications |

### Broadcast Events

- `notification:speak` -- fired when a notification is sent, consumed by frontend TTS
- `workspace:changed` -- fired on workspace switch, triggers frontend refresh
- Workflow transitions -- broadcast when tasks change state

### Source Files

| File | Description |
|------|-------------|
| `app/websocket/server.go` | `Server` struct, Fiber setup, middleware, routes, `Start()`/`Stop()` |
| `app/websocket/service.go` | `Service` -- client registry, read/write pumps, message routing, broadcast |
| `app/websocket/protocol.go` | `Message` struct, type constants, constructors, JSON serialization |
| `app/websocket/auth.go` | WebSocket authentication |
| `app/websocket/ai_bridge.go` | `AIBridgeAdapter` -- bridges AI sessions to WebSocket |
| `app/websocket/sync_handler.go` | `SyncHandler` -- entity sync over WebSocket |
| `app/websocket/browser_handler.go` | `BrowserHandler` -- Chrome extension integration |
| `app/websocket/preview_coordinator.go` | `PreviewCoordinator` -- live preview sync |
| `app/websocket/preview_handler.go` | Preview-specific message handlers |
| `app/websocket/workflow_listener.go` | `WorkflowBroadcaster` -- workflow transition events |
| `app/websocket/settings_handler.go` | Settings change broadcast |
| `app/websocket/integration.go` | Integration-specific handlers |

---

## 4. Rust Engine (gRPC)

### Overview

A Rust gRPC server providing CPU-intensive operations: code parsing (Tree-sitter), search indexing (Tantivy), and vector-based memory storage. Communicates with Go via Protocol Buffers.

- **Binary**: `bin/orchestra-engine` (built from `engine/`)
- **Entry point**: `engine/src/main.rs`
- **Port**: `50051` (hardcoded as `[::1]:50051`)
- **Framework**: Tonic 0.12, Tokio async runtime

### CLI Arguments

| Argument | Default | Description |
|----------|---------|-------------|
| `--workspace`, `-w` | `.` | Workspace directory |
| `--storage` | `<workspace>/storage` | Storage directory for indices and memory DBs |

### gRPC Services

| Service | Proto Server | Implementation | Purpose |
|---------|-------------|----------------|---------|
| Health | `HealthServiceServer` | `HealthServiceImpl` | Health check / liveness probe |
| Parse | `ParseServiceServer` | `ParseServiceImpl` | Tree-sitter code parsing (14 language grammars) |
| Search | `SearchServiceServer` | `SearchServiceImpl` | Tantivy full-text search indexing |
| Memory | `MemoryServiceServer` | `MemoryServiceImpl` | Vector-based memory storage (RAG) |

### Supported Languages (Tree-sitter)

Rust, Go, JavaScript, TypeScript, Python, C, C++, Java, HTML, CSS, JSON, TOML, YAML, Markdown

### Storage Layout

```
<storage>/
  search_index/    # Tantivy search index
  memory/          # Vector memory SQLite databases
```

### Key Dependencies

| Crate | Version | Purpose |
|-------|---------|---------|
| `tonic` | 0.12 | gRPC server framework |
| `prost` | 0.13 | Protocol buffer serialization |
| `tokio` | 1 (full) | Async runtime |
| `tantivy` | 0.22 | Full-text search engine |
| `tree-sitter` | 0.24 | Incremental parsing |
| `rusqlite` | 0.32 (bundled) | Local SQLite for memory storage |
| `clap` | 4.5 | CLI argument parsing |
| `tracing` | 0.1 | Structured logging |
| `pulldown-cmark` | 0.12 (SIMD) | Markdown parsing |
| `walkdir` | 2.5 | Directory traversal |

### Graceful Shutdown

Handles both `Ctrl+C` (SIGINT) and `SIGTERM` via `tokio::signal`. Uses `serve_with_shutdown()` for clean gRPC shutdown.

### Go-Side Engine Management

The Go side manages the Rust engine process via `app/engine/mcp/`:

- `Manager` -- starts/stops the engine subprocess, tracks running state
- `Client` -- gRPC client wrapper with `Dial()` and `Close()`
- `Bridge` -- high-level API wrapping the gRPC client, with fallback to TOON/markdown when the engine is unavailable

### Source Files

| File | Description |
|------|-------------|
| `engine/src/main.rs` | Entry point -- CLI parsing, service creation, gRPC server setup |
| `engine/Cargo.toml` | Dependencies and build configuration |
| `engine/build.rs` | `tonic-build` for protobuf code generation |
| `engine/src/services/health.rs` | Health service implementation |
| `engine/src/services/parse.rs` | Tree-sitter parsing service |
| `engine/src/services/search.rs` | Tantivy search service |
| `engine/src/services/memory.rs` | Vector memory service |
| `engine/src/index/` | `IndexManager` for Tantivy index lifecycle |

---

## 5. Desktop App (Wails v3)

### Overview

A native macOS desktop application (`.app` bundle) that embeds the React frontend and orchestrates all backend services. Built with Wails v3 for Go-React bindings, with system tray, window management, and multiple window modes.

- **Binary**: `Orchestra MCP.app/Contents/MacOS/Orchestra MCP`
- **Entry point**: `cmd/desktop/main.go`
- **Bootstrap**: `bootstrap/desktop.go` (`RunDesktopWithIcons()`)
- **Bundle ID**: `com.orchestra-mcp.desktop`
- **Embedded assets**: React frontend via `//go:embed all:dist`, tray icon, app icon, notification sounds

### Embedded Resources

- `dist/` -- compiled React frontend (embedded via `embed.FS`)
- `tray.png` -- system tray icon
- `appicon-bg.png` -- app icon with background (Linux/Windows)
- `sounds/` -- notification sound files (MP3, embedded via `embed.FS`)

### Startup Sequence

1. **Resolve project root** -- detects if running from `.app` bundle or development mode
2. **Load `.env`** from project root (best-effort)
3. **Bootstrap plugins** -- discover, resolve dependencies, load, boot, register
4. **Open Settings Store** -- SQLite at `~/Library/Application Support/Orchestra/settings.db`
5. **Open Session Store** -- SQLite for chat session persistence
6. **Create AI Bridge** -- with settings getter for API keys
7. **Settings API** (port 19191) -- start HTTP server; if port in use, start mini proxy on 19192
8. **Initialize integrations** -- GitHub, Notion, Jira, Linear, Figma OAuth hubs
9. **Wire MCP Bridge** -- register all MCP tools accessible via HTTP
10. **Wire Sync Client** -- if auth hub available, start pull polling and outbox worker
11. **Wire session sync hook** -- enqueue outbox records on session create/rename/delete
12. **Initialize DevTools** -- terminal, SSH, database, logs session manager
13. **LSP WebSocket proxy** -- for Monaco editor language intelligence
14. **Notification service** -- resolve bundled sounds directory
15. **Auto-start services** (in goroutine after 500ms delay):
    - WebSocket server (port 8765)
    - Rust engine (port 50051)
    - Discord bot (if configured and not running as subprocess)
16. **Start Wails app** -- launches the native window on the main goroutine (required for macOS GUI)
17. **Auto-updater** -- checks for updates (skipped in dev mode when version is "dev")

### Service Registry

The desktop app registers services in a `ServiceRegistry` with start/stop callbacks:

| Service ID | Label | Type | Port |
|------------|-------|------|------|
| `settings-api` | Settings API | HTTP | 19191 |
| `websocket` | WebSocket Server | WebSocket | 8765 |
| `grpc-engine` | Rust Engine (gRPC) | gRPC | 50051 |

### Subprocess Mode

When launched as a subprocess of `orch --dev`, the environment variable `PARENT_WEBSOCKET_PORT` is set. In this mode:

- WebSocket server is **not started** (parent owns it) -- marked as running
- Rust engine is **not started** (parent owns it) -- marked as running
- Discord bot is **not started** (parent owns it)

### Shutdown Sequence

1. Stop auto-updater
2. Shutdown desktop (Wails)
3. Stop outbox worker and sync client
4. Stop Discord bot
5. Stop WebSocket server
6. Stop Rust engine
7. Close DevTools session manager
8. Close auth, GitHub, Notion, Jira, Linear hubs
9. Stop Settings API and close store
10. Close session store
11. Shutdown plugins

### Source Files

| File | Description |
|------|-------------|
| `cmd/desktop/main.go` | Embeds assets, resolves project root, calls `RunDesktopWithIcons()` |
| `bootstrap/desktop.go` | `DesktopApp` -- full service orchestration, startup, shutdown |
| `bootstrap/app.go` | `App` -- plugin bootstrap (discover, resolve deps, load, boot, register) |

---

## 6. Discord Bot

### Overview

An optional Discord integration that connects to the Discord Gateway WebSocket, registers slash commands, routes messages to handlers, and provides AI chat + MCP tool execution from Discord.

- **Package**: `app/discord/`
- **Manager**: `app/bot/manager.go` (`Manager` orchestrates Discord + Slack bots)
- **Protocol**: Discord Gateway WebSocket (v10) + Discord REST API (v10)
- **Gateway intents**: GUILDS (1<<0), GUILD_MESSAGES (1<<9), MESSAGE_CONTENT (1<<15)

### Startup

1. `bot.Manager.StartDiscord(cfg)` creates a `discord.Bot`
2. Collects all MCP tools and prompts from the `tools` package
3. Registers all command handlers (see below)
4. Connects to Discord Gateway:
   - Fetches gateway URL from `https://discord.com/api/v10/gateway`
   - Establishes WebSocket connection (`/?v=10&encoding=json`)
   - Reads Hello (op 10) to get heartbeat interval
   - Sends Identify (op 2) with token, intents, and presence
   - Starts heartbeat loop (at server-specified interval)
   - Starts read loop for dispatching events
5. Registers slash commands via REST API (if `application_id` is configured)
6. Sets presence to "Playing Orchestra MCP"
7. Starts permission event listener

### Gateway Event Dispatch

- `MESSAGE_CREATE` -- routed through `Router.RouteMessage()` to matching handler by prefix
- `INTERACTION_CREATE` -- routed through `Router.RouteInteraction()` for slash commands and button clicks

### Registered Handlers

| Handler | Module | Commands |
|---------|--------|----------|
| `PingHandler` | `handlers/ping.go` | `!ping`, `/ping` |
| `StatusHandler` | `handlers/status.go` | `!status`, `/status` |
| `ChatHandler` | `handlers/chat.go` | AI chat (default fallback for unmatched prefix commands) |
| `PermissionHandler` | `handlers/permission.go` | AI permission approval/denial via buttons |
| `StopHandler` | `handlers/stop.go` | `!stop`, `/stop` -- stop AI generation |
| `McpHandler` | `handlers/mcp.go` | `!mcp`, `/mcp` -- execute MCP tools |
| `ToolsHandler` | `handlers/tools.go` | `!tools`, `/tools` -- list available tools |
| `PromptsHandler` | `handlers/prompts.go` | `!prompts`, `/prompts` -- list/execute prompts |
| `ActionsHandler` | `handlers/actions.go` | Button interaction handler |
| `ProgressHandler` | `handlers/progress.go` | Task progress tracking and notifications |

### Configuration (from Settings API)

| Key | Description |
|-----|-------------|
| `discord.enabled` | Enable/disable Discord bot |
| `discord.bot_token` | Discord bot token |
| `discord.application_id` | Application ID for slash command registration |
| `discord.channel_id` | Default channel ID |
| `discord.guild_id` | Guild (server) ID for guild-scoped slash commands |
| `discord.command_prefix` | Prefix for text commands (default: `!`) |
| `discord.webhook_url` | Webhook URL for webhook-only mode |

### Modes

- **Full bot** -- gateway connection + slash commands + text commands (requires `bot_token` + `application_id`)
- **Webhook-only** -- sends messages via webhook URL, no gateway connection (requires `webhook_url`)
- **Disabled** -- no Discord integration

### Source Files

| File | Description |
|------|-------------|
| `app/discord/bot.go` | `Bot` -- main bot struct, `Start()`/`Stop()`, gateway event dispatch |
| `app/discord/gateway.go` | `Gateway` -- WebSocket connection, heartbeat, identify, read loop |
| `app/discord/rest.go` | `RestClient` -- Discord REST API (send/edit messages, register commands) |
| `app/discord/router.go` | `Router` -- routes messages/interactions to handlers |
| `app/discord/handler.go` | `Handler` interface |
| `app/discord/service.go` | `Service` -- workflow transition listener for Discord notifications |
| `app/discord/types.go` | Discord API types (messages, interactions, components) |
| `app/discord/embed_helpers.go` | Embed building helpers |
| `app/discord/handlers/` | All command handler implementations |
| `app/bot/manager.go` | `Manager` -- orchestrates Discord + Slack bots, collects MCP tools |

---

## 7. Sync Client (Background)

### Overview

A pull-based HTTP sync client that communicates with the Orchestra web API (Laravel backend). Runs in the background, periodically polling for remote changes and pushing local changes via an outbox pattern with exponential-backoff retry.

- **Package**: `app/syncclient/`
- **Protocol**: HTTP REST with Bearer token (Sanctum) authentication
- **Poll interval**: 30 seconds
- **Backoff**: Exponential, 1s min to 60s max (pull), 5s base to 7 days max (outbox retry)
- **Max outbox retries**: 10 attempts per record

### Pull Loop

1. `Client.Run(ctx, since)` starts a blocking poll loop
2. Each tick: `GET /api/sync/pull?since=<timestamp>&device_id=<id>`
3. On success: deliver records to `PullHandler`, advance cursor to latest `synced_at`
4. On error: log warning, backoff with exponential delay, retry
5. Records are dispatched to:
   - `SessionStore` -- for `ai_session` entity type (create/update/delete)
   - `SyncHandler` -- broadcast to frontend via WebSocket for all entity types

### Push

- `Client.Push(ctx, records)` sends `POST /api/sync/push` with batch of `PushRecord`
- Each record: `entity_type`, `entity_id`, `action` (upsert/delete), `payload`, `version`, `idempotency_key`
- Returns per-record results: `ok`, `skipped`, or `error`

### Outbox Pattern

For resilience, local changes are first written to a local SQLite outbox, then flushed to the server:

1. **`OutboxStore`** (`app/syncclient/outbox.go`) -- SQLite database at `<configDir>/storage/outbox/outbox.db`
   - WAL mode, single connection
   - Schema: `id`, `entity_type`, `entity_id`, `action`, `payload`, `version`, `idempotency_key`, `retry_count`, `next_retry_at`, `created_at`
   - `Enqueue(rec)` -- insert with `retry_count=0`, `next_retry_at=now`
   - `DueRecords(limit)` -- fetch rows where `next_retry_at <= now` and `retry_count < 10`
   - `MarkSuccess(id)` -- delete the row
   - `MarkFailed(id, retryCount)` -- increment retry, schedule next attempt with exponential backoff

2. **`OutboxWorker`** (`app/syncclient/outbox_worker.go`) -- background goroutine
   - Flushes every 30 seconds or on-demand via `Flush()` signal
   - Fetches up to 50 due records per batch
   - Pushes batch to server via `Client.Push()`
   - Marks each record as success or failed based on server response
   - Records exceeding 10 retries are dropped

### Device Registration

- `Client.RegisterDevice(ctx, hostname, platform)` calls `POST /api/sync/devices/register`
- Called on startup; failures are non-fatal (user may not be authenticated yet)

### Session Sync Hook

A `sessionSyncHook` is wired into the `SessionStore` to automatically enqueue outbox records on:
- Session saved (create/update) -- `action: "upsert"` with full session payload
- Session deleted -- `action: "delete"` with session ID

### Source Files

| File | Description |
|------|-------------|
| `app/syncclient/client.go` | `Client` -- pull polling loop, push, device registration |
| `app/syncclient/outbox.go` | `OutboxStore` -- SQLite outbox with enqueue/due/mark operations |
| `app/syncclient/outbox_worker.go` | `OutboxWorker` -- background flush loop with batching |
| `app/syncclient/migration.go` | Database migration definitions |
| `app/syncclient/migration_runner.go` | Migration execution logic |

---

## Startup Sequence

When the desktop app launches (`Orchestra MCP.app`), services boot in this order:

```
1. Plugin bootstrap          Discover → resolve deps → load → boot → register
2. Settings Store            Open SQLite at ~/Library/Application Support/Orchestra/settings.db
3. Session Store             Open SQLite for chat session persistence
4. AI Bridge                 Create bridge with settings getter
5. Integration hubs          GitHub, Notion, Jira, Linear, Figma OAuth (non-blocking)
6. MCP Bridge                Register all MCP tools for HTTP access
7. Settings API (19191)      Start HTTP server; if port in use, proxy on 19192
8. Sync Client               Start pull polling + outbox worker (background)
9. DevTools                  Terminal, SSH, database, logs session manager
10. LSP Proxy                Monaco editor language intelligence
11. Notification service     Resolve bundled sounds
12. Desktop window           Launch Wails app on main goroutine
--- (after 500ms delay, in goroutine) ---
13. WebSocket Server (8765)  Start Fiber WebSocket server
14. Rust Engine (50051)      Start gRPC engine subprocess
15. Discord Bot              Connect gateway, register commands (if configured)
16. Auto-updater             Check for updates (if not dev mode)
```

When running in **dev mode** (`orch --dev`), the sequence is slightly different:

```
1. Plugin bootstrap
2. AI Bridge + Session Registry
3. Settings Store + Settings API (19191)
4. MCP Bridge
5. Rust Engine (50051)
6. Auth + GitHub hubs
7. DevTools + LSP
8. Firebase + FCM (if configured)
9. WebSocket Server (8765)
10. AI Bridge Adapter + Workflow Broadcaster
11. Discord Bot (if configured)
12. Vite dev server (port 9245) + Desktop app subprocess
13. Plugin file watcher
```

## Port Conflict Resolution

The system handles port conflicts gracefully:

### Settings API (19191)

- Desktop calls `srv.Start()` which attempts `net.Listen("tcp", "127.0.0.1:19191")`
- If the port is already bound (another `orch --dev` or desktop instance is running):
  - The new instance **closes its own settings store** (avoids dual SQLite writers)
  - Sets `portInUse = true` and reuses the existing server
  - Starts a **secondary proxy on port 19192** for window commands and notifications
  - Logs: `"Settings API port in use, using existing server"`

### WebSocket Server (8765)

- Before auto-starting, the desktop re-probes the port via `desktop.IsPortOpen("127.0.0.1", 8765)`
- If already open: marks service as `ServiceRunning` without starting a new server
- If `PARENT_WEBSOCKET_PORT` env is set: skips entirely (parent process owns it)

### Rust Engine (50051)

- Same probe logic as WebSocket: `desktop.IsPortOpen("127.0.0.1", 50051)`
- If already open: marks as running
- If `PARENT_WEBSOCKET_PORT` env is set: skips entirely

### General Pattern

The system follows a "first-come-first-served" model:
1. Try to bind the port
2. If already bound, check if the service is healthy
3. If healthy, reuse the existing instance
4. If not healthy, report error

This allows multiple Orchestra processes (CLI, desktop, dev mode) to coexist safely.

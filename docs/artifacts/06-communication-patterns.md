# Communication Patterns — Orchestra Reference

> How every service talks to every other service. Extracted from `orch-ref/`.

## Pattern Summary

```
                         ┌──────────────────────┐
                         │    Orchestra Web API  │
                         │  (Laravel / Cloud)    │
                         └──────────┬───────────┘
                                    │ HTTPS (Sanctum Bearer)
                                    │ POST /api/sync/push
                                    │ GET  /api/sync/pull
                         ┌──────────▼───────────┐
                         │     Sync Client       │
                         │   (HTTP poll + push)  │
                         └──────────┬───────────┘
                                    │ PullHandler callback
 ┌────────────┐      ┌─────────────▼──────────────────────────────────┐
 │ Chrome Ext │◄────►│              WebSocket Server                  │
 │  (8765)    │  WS  │              (:8765 Fiber)                     │
 └────────────┘      │                                                │
                     │  ┌───────────────┐  ┌────────────────────┐     │
 ┌────────────┐      │  │ SyncHandler   │  │ AIBridgeAdapter    │     │
 │ Desktop UI │◄────►│  │ (entity sync) │  │ (AI session relay) │     │
 │  (Wails)   │  WS  │  └───────────────┘  └────────────────────┘     │
 └────────────┘      │  ┌───────────────┐  ┌────────────────────┐     │
                     │  │ BrowserHandler│  │ PreviewCoordinator │     │
                     │  │ (page context)│  │ (live preview)     │     │
                     │  └───────────────┘  └────────────────────┘     │
                     │  ┌───────────────┐  ┌────────────────────┐     │
                     │  │SettingsHandler│  │WorkflowBroadcaster │     │
                     │  │ (key-value)   │  │ (task transitions) │     │
                     │  └───────────────┘  └────────────────────┘     │
                     └────────────────────────────────────────────────┘
                                    ▲
                                    │ broadcast
 ┌────────────────┐                 │
 │ Settings API   │─────────────────┘  workspace:changed events
 │ (:19191 HTTP)  │
 │                │◄───── McpBridge (in-process tool calls)
 │                │◄───── AI Bridge (Claude CLI subprocess)
 │                │◄───── Integration Hubs (GitHub, Notion, Jira, Linear, Figma)
 │                │◄───── Auth Hub (Orchestra web login)
 │                │◄───── DevTools (terminal, SSH, DB, logs)
 │                │◄───── LSP Proxy (Monaco ↔ language servers)
 └───────┬────────┘
         │
 ┌───────▼────────┐      ┌────────────────┐
 │ Desktop Proxy  │      │  Discord Bot   │
 │ (:19192 HTTP)  │      │  (discordgo)   │
 │ (window cmds)  │      │                │
 └────────────────┘      └───────┬────────┘
                                 │ bot.Manager.SetSessions()
                         ┌───────▼────────┐
                         │ SessionRegistry│
                         │  (ai sessions) │
                         └───────┬────────┘
                                 │ Bridge.Send()
                         ┌───────▼────────┐      stdio stream-json
                         │   AI Bridge    │◄────────────────────────►┌──────────┐
                         │ (Go process)   │                          │Claude CLI│
                         └───────┬────────┘                          └──────────┘
                                 │ MemoryClient (gRPC / markdown fallback)
                         ┌───────▼────────┐
                         │  Rust Engine   │
                         │  (:50051 gRPC) │
                         │  Memory/Search │
                         └────────────────┘
```

---

## 1. MCP Server ↔ Rust Engine

### Protocol
gRPC on `localhost:50051` (configurable via `ORCHESTRA_ENGINE_PORT`).

### Flow
1. `Manager.Start(workspace, storagePath)` checks if port 50051 is already open.
2. If not, resolves the `orchestra-engine` binary (same dir as executable, dev paths, then `$PATH`).
3. Spawns the binary as a subprocess with `--workspace` and optional `--storage` flags.
4. Waits up to 3 seconds for the port to become reachable (100ms polling).
5. Callers use `mcp.Dial(addr)` to get a gRPC `Client` wrapping `MemoryServiceClient`.
6. Available RPCs: `StartSession`, `EndSession`, `RecordObservation`, `SearchMemory`, `GetContext`, `ListSessions`, `GetSession`.
7. On shutdown, `Manager.Stop()` sends SIGTERM, waits 3s, then SIGKILL.

### Fallback
If the engine binary is not found (`ErrNotFound`) or the gRPC dial fails, all memory operations fall back to markdown-based chunk files at `.projects/<project>/.memory/chunks.md`.

### Key Code
```go
// app/engine/mcp/manager.go
func (m *Manager) Start(workspace, storagePath string) error {
    if portOpen(m.port) {
        m.running = true
        return nil // reuse existing engine
    }
    bin := resolverFunc()
    if bin == "" {
        return ErrNotFound
    }
    m.cmd = exec.Command(bin, args...)
    m.cmd.Start()
    // poll for 3s...
}
```

```go
// app/engine/mcp/client.go — gRPC wrapper
func Dial(addr string) (*Client, error) {
    conn, err := grpc.NewClient(addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    return &Client{conn: conn, memory: pb.NewMemoryServiceClient(conn)}, nil
}
```

### Source Files
- `orch-ref/app/engine/mcp/manager.go` — engine lifecycle (start/stop/kill)
- `orch-ref/app/engine/mcp/client.go` — gRPC client wrapper
- `orch-ref/app/engine/mcp/resolve.go` — binary resolution + port config
- `orch-ref/app/engine/mcp/bridge.go` — Bridge with gRPC/TOON fallback flag

---

## 2. Settings API ↔ MCP Bridge

### Protocol
In-process Go function calls. No network involved.

### Flow
1. At bootstrap, `NewMcpBridge(workspace)` loads all 85+ MCP tools by calling each tool family constructor (Project, Epic, Story, Task, Workflow, Prd, Sprint, etc.).
2. Each tool is stored in a `map[string]Tool` keyed by tool name.
3. The Settings HTTP server exposes `POST /api/mcp/tools/call` which calls `mcpBridge.CallTool(name, args)`.
4. `ListTools()` returns all tool definitions (used by the frontend to render the tool catalog).
5. On workspace change, `ReloadWorkspace(newPath)` rebuilds the entire tool registry under a write lock.

### Thread Safety
All access to the tool map is protected by `sync.RWMutex` to prevent concurrent read/write panics during workspace switches.

### Key Code
```go
// app/settings/mcp_bridge.go
func (b *McpBridge) CallTool(name string, args map[string]any) (*t.ToolResult, error) {
    b.mu.RLock()
    tool, ok := b.toolMap[name]
    b.mu.RUnlock()
    if !ok {
        return &t.ToolResult{Content: []t.ContentBlock{{Type: "text", Text: "unknown tool: " + name}}, IsError: true}, nil
    }
    return tool.Handler(args)
}
```

### Source Files
- `orch-ref/app/settings/mcp_bridge.go` — McpBridge (tool registry + execution)
- `orch-ref/bootstrap/desktop.go` lines 398-403 — wiring workspace resolution and bridge creation

---

## 3. Settings API ↔ AI Bridge (Claude CLI Subprocess)

### Protocol
- **Inbound (Go -> CLI):** stdio `stream-json` format. JSON messages written to the subprocess's stdin pipe.
- **Outbound (CLI -> Go):** NDJSON lines read from stdout, parsed by `Parser`.
- **Frontend -> Settings API:** HTTP `POST /api/ai/send`, `POST /api/ai/permission`, `POST /api/ai/question`. Responses streamed as SSE via `bridge.Subscribe()`.

### Flow
1. `Bridge.Send(req)` spawns `claude --output-format stream-json --input-format stream-json --model <model>`.
2. The process runs in the workspace directory (`cmd.Dir = workspace`).
3. User messages are written as JSON to stdin with `{"type":"user","message":{"role":"user","content":[...]}}`.
4. The bridge reads stdout line-by-line, strips ANSI, parses JSON, and classifies chunks: `text`, `thinking`, `action`, `permission_request`, `tool_result`, `result`, `error`.
5. Each parsed `ResponseChunk` is broadcast to all SSE subscribers via `bridge.broadcast(chunk)`.
6. When a `result` type is detected (`parser.GotResult`), the turn is saved to conversation history and the CLI process stays alive for follow-up messages.
7. `FollowUp(prompt, images, pageCtx)` writes a new user message to the existing stdin pipe — no new process spawn.
8. `Permission()` and `Question()` send `control_response` JSON to stdin.
9. On shutdown, `Stop()` cancels the context which kills the process.

### Context Injection
Before spawning the CLI, the bridge builds `--append-system-prompt` with up to 4 context sources:
1. **In-session conversation history** — `history.BuildSystemContext()`
2. **Cross-session memory** — `MemoryClient.RetrieveContext(prompt)` via Rust engine gRPC
3. **Active file context** — reads `workspace.active_file` from settings
4. **Current task context** — reads `workspace.current_task` from settings

### Key Code
```go
// app/ai/bridge.go — spawning the Claude CLI
cmd := exec.CommandContext(ctx, binary, args...)
cmd.Dir = ws
stdout, _ := cmd.StdoutPipe()
stdinPipe, _ := cmd.StdinPipe()
cmd.Start()

// Write initial user message
userMsg := buildUserMessage(prompt, images)
b.writeStdin(userMsg)

// Read response stream
scanner := bufio.NewScanner(stdout)
for scanner.Scan() {
    chunks := parser.Parse(line)
    for _, chunk := range chunks {
        b.broadcast(chunk) // -> SSE subscribers + WebSocket adapter
    }
}
```

### Source Files
- `orch-ref/app/ai/bridge.go` — Bridge (CLI lifecycle, stdin/stdout, subscriber fan-out)
- `orch-ref/app/ai/types.go` — SendRequest, ResponseChunk, PermissionResponse, PageContext
- `orch-ref/app/ai/registry.go` — SessionRegistry (multi-session management)
- `orch-ref/app/ai/memory.go` — MemoryClient (gRPC + markdown fallback)

---

## 4. Settings API ↔ Integration Hubs

### Protocol
In-process Go method calls from the Settings HTTP server to hub instances. Each hub manages OAuth tokens, API proxying, and credential storage.

### Hubs

| Hub | Init | Credential Store | Key Endpoints |
|-----|------|-----------------|---------------|
| `auth.Hub` | `auth.NewHub(configDir)` | `$configDir/auth.json` | `POST /api/auth/login`, `POST /api/auth/callback` |
| `github.Hub` | `github.NewHub(configDir)` | `$configDir/github-creds.json` | `POST /api/github/auth`, `GET /api/github/repos` |
| `notion.Hub` | `notion.NewHub(configDir)` | `$configDir/notion-creds.json` | `POST /api/notion/auth`, OAuth callback |
| `jira.Hub` | `jira.NewHub(configDir)` | `$configDir/jira-creds.json` | `POST /api/jira/auth`, OAuth callback |
| `linear.Hub` | `linear.NewHub(configDir)` | `$configDir/linear-creds.json` | `POST /api/linear/auth`, OAuth callback |
| `figma.*` | `figma.NewAuthService(cfg, tokenStore)` | `$configDir/figma-tokens.json` | `POST /api/figma/auth`, proxy API |

### Flow
1. Each hub is created at bootstrap with the Orchestra config directory.
2. The Settings server registers each hub via `srv.SetGitHubHub(ghHub)`, `srv.SetNotionHub(notionHub)`, etc.
3. Hub methods are called directly from HTTP handlers — no RPC, no message queue.
4. OAuth callbacks land on the Settings API HTTP server and are forwarded to the appropriate hub.
5. Token refresh is handled internally by each hub.

### Key Code
```go
// bootstrap/desktop.go — wiring hubs
ghHub, _ := github.NewHub(orchConfigDir)
srv.SetGitHubHub(ghHub)

notionHub, _ := notion.NewHub(orchConfigDir)
srv.SetNotionHub(notionHub)

jiraHub, _ := jira.NewHub(orchConfigDir)
srv.SetJiraHub(jiraHub)

linearHub, _ := linear.NewHub(orchConfigDir)
srv.SetLinearHub(linearHub)

figmaAuth := figma.NewAuthService(figmaCfg, figmaTokenStore)
srv.SetFigmaAuth(figmaAuth)
srv.SetFigmaFiles(figma.NewFilesClient(figmaTokenStore))
srv.SetFigmaProxy(figma.NewProxyClient(figmaTokenStore))
```

### Source Files
- `orch-ref/bootstrap/desktop.go` lines 267-318 — hub initialization and wiring
- `orch-ref/app/settings/server.go` — HTTP handlers that delegate to hubs
- `orch-ref/app/settings/github.go`, `notion.go`, `jira.go`, `linear.go`, `figma_auth.go` — hub route handlers

---

## 5. WebSocket Server ↔ Frontend Clients

### Protocol
WebSocket on `ws://127.0.0.1:8765/` (Fiber v3 with `fasthttp/websocket` upgrader).

### Message Protocol
All messages use a unified JSON envelope:

```json
{
  "id": "20060102150405.000000",
  "type": "event|command|response|error|ping|pong",
  "event": "event_name",
  "command": "command_name",
  "data": { ... },
  "timestamp": 1234567890,
  "request_id": "...",
  "error": { "code": "...", "message": "..." }
}
```

Six message types: `event`, `command`, `response`, `error`, `ping`, `pong`.

### Connection Flow
1. Client connects to `ws://127.0.0.1:8765/?user_id=...&session_id=...`.
2. `HandleUpgrade()` creates a `Client` struct with a 256-slot send channel.
3. Client is registered in the `Service.clients` map (keyed by auto-generated ID).
4. Two goroutines start per client: `readPump` (incoming messages) and `writePump` (outgoing messages + ping keepalive).
5. Pings are sent every 30s; pong must arrive within 60s or the client is evicted.
6. Incoming messages are deserialized via `FromJSON()`, validated, and routed to registered `MessageHandler` functions by `MessageType`.

### Routing
Handlers are registered per `MessageType` (event, command). Multiple handlers can be registered for the same type. All matching handlers execute sequentially.

### Delivery Modes
- `Broadcast(msg)` — all connected clients
- `BroadcastToUser(userID, msg)` — all clients for a specific user
- `SendTo(clientID, msg)` — single client

### Key Code
```go
// app/websocket/service.go
type Service struct {
    clients   map[string]*Client
    handlers  map[MessageType][]MessageHandler
    broadcast chan *Message  // buffered channel, capacity 256
    register  chan *Client
    unregister chan *Client
}

func (s *Service) run() {
    for {
        select {
        case client := <-s.register:    s.registerClient(client)
        case client := <-s.unregister:  s.unregisterClient(client)
        case message := <-s.broadcast:  s.broadcastMessage(message)
        case <-ticker.C:                s.pingClients()
        case <-s.ctx.Done():            s.closeAllClients(); return
        }
    }
}
```

### Source Files
- `orch-ref/app/websocket/server.go` — Fiber HTTP server with WS upgrade
- `orch-ref/app/websocket/service.go` — client management, message routing, broadcast
- `orch-ref/app/websocket/protocol.go` — Message types and JSON serialization

---

## 6. Desktop App ↔ Settings API (Port Proxy)

### Protocol
HTTP on `127.0.0.1:19191` (primary) and `127.0.0.1:19192` (secondary proxy).

### Architecture
The Settings API is a local HTTP server that acts as the central hub between the desktop UI (Wails), the Chrome extension, and all backend services. There is a chicken-and-egg problem: the Settings server starts before the Wails app, so a `lazyOpener` proxy is used.

### Dual-Port Strategy
1. **Port 19191** — The primary Settings API. First process to bind it owns it.
2. **Port 19192** — Secondary desktop proxy. Used when port 19191 is already occupied by another process (e.g., `orchestra start` CLI). Exposes only window commands and notifications.

### Flow (Primary)
```
Desktop UI (Wails) → HTTP :19191 → Settings Server → Store / AI Bridge / McpBridge / Hubs
```

### Flow (Secondary Proxy)
When the desktop app finds port 19191 already in use:
```
MCP/CLI owns :19191 → Desktop app listens on :19192
  POST /api/windows/open   → lazyOpener → Wails window manager
  POST /api/windows/close/ → lazyOpener → Wails window manager
  POST /api/notify         → NotifService, then forwards to :19191/api/notify/speak (TTS)
  POST /api/screenshot/start → Wails screenshot overlay
```

### Key Code
```go
// bootstrap/desktop.go — lazyOpener solves the startup race
type lazyOpener struct {
    impl *desktop.App
}
func (l *lazyOpener) OpenWindow(req settings.OpenWindowRequest) error {
    deadline := time.Now().Add(15 * time.Second)
    for l.impl == nil {
        if time.Now().After(deadline) {
            return fmt.Errorf("desktop app not started yet (timeout)")
        }
        time.Sleep(100 * time.Millisecond)
    }
    return l.impl.OpenWindow(req)
}

// Secondary proxy when port 19191 is in use
if portInUse {
    go func() {
        mux := http.NewServeMux()
        mux.HandleFunc("POST /api/windows/open", ...)
        mux.HandleFunc("POST /api/notify", func(w, r) {
            notifSvc.Send(req)
            // Forward TTS event to the main server
            go client.Post("http://127.0.0.1:19191/api/notify/speak", ...)
        })
        srv := &http.Server{Addr: "127.0.0.1:19192", Handler: mux}
        srv.ListenAndServe()
    }()
}
```

### Source Files
- `orch-ref/bootstrap/desktop.go` lines 47-167 — `lazyOpener` proxy
- `orch-ref/bootstrap/desktop.go` lines 455-539 — dual-port startup logic
- `orch-ref/app/settings/server.go` — Settings HTTP server (primary on :19191)

---

## 7. Sync Client ↔ Orchestra Web API

### Protocol
HTTPS with Sanctum Bearer token authentication. Pull-based polling with push capability.

### Endpoints
| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/sync/pull?since=<RFC3339>&device_id=<id>` | Fetch remote changes since cursor |
| `POST` | `/api/sync/push` | Upload local changes (batch of PushRecords) |
| `POST` | `/api/sync/devices/register` | Register this device with the server |

### Pull Flow
1. `Client.Run(ctx, since)` starts a blocking poll loop.
2. Every 30 seconds, calls `GET /api/sync/pull` with the cursor timestamp.
3. Returns `[]PullRecord` — each has `entity_type`, `entity_id`, `action`, `payload`, `version`, `synced_at`.
4. Records are delivered to the `PullHandler` callback.
5. The callback in `bootstrap/desktop.go` applies pulled records:
   - `ai_session` records go to `SessionStore`
   - All records go to `SyncHandler.ApplyPullRecords()` which broadcasts them to WebSocket clients
6. Cursor advances to the latest `synced_at` from the batch.
7. On error, exponential backoff: 1s -> 2s -> 4s -> ... -> 60s max.

### Push Flow
1. `Client.Push(ctx, []PushRecord)` sends a batch to `POST /api/sync/push`.
2. Each record has `entity_type`, `entity_id`, `action`, `payload`, `version`, `idempotency_key`.
3. Server returns per-record results: `"ok"`, `"skipped"`, or `"error"`.

### Outbox Persistence
Changes are not pushed directly — they go through an outbox pattern:
1. `OutboxStore.Enqueue(record)` writes to a local SQLite database (`outbox.db` in WAL mode).
2. `OutboxWorker.Run(ctx)` flushes due records every 30 seconds (or immediately on `Flush()` signal).
3. Records get up to 10 retry attempts with exponential backoff (5s * 2^retryCount, max 7 days).
4. After 10 failures, the record is dropped.
5. On startup, the worker immediately flushes any records that survived a restart.

### Session Sync Hook
```go
// bootstrap/desktop.go — hooks into SessionStore
d.sessionStore.SetSyncHook(&sessionSyncHook{
    outbox:  d.outboxStore,
    worker:  d.outboxWorker,
})

func (h *sessionSyncHook) OnSessionSaved(sess StoredSession) {
    rec := syncclient.PushRecord{
        EntityType: "ai_session",
        EntityID:   syncclient.SessionEntityID(sess.ID),
        Action:     "upsert",
        Payload:    map[string]any{ ... },
    }
    h.outbox.Enqueue(rec)
    h.worker.Flush() // signal immediate push
}
```

### Key Code
```go
// app/syncclient/client.go — pull loop
func (c *Client) Run(ctx context.Context, since time.Time) {
    for {
        records, newCursor, err := c.pull(ctx, cursor)
        if err != nil {
            // exponential backoff: 1s -> 60s
            backoff = time.Duration(math.Min(float64(backoff*2), float64(c.maxBackoff)))
            continue
        }
        backoff = c.minBackoff
        if len(records) > 0 {
            c.handler(records)
            cursor = newCursor
        }
        time.After(c.pollInterval) // 30s
    }
}
```

### Source Files
- `orch-ref/app/syncclient/client.go` — Client (pull polling + push)
- `orch-ref/app/syncclient/outbox.go` — OutboxStore (SQLite persistence)
- `orch-ref/app/syncclient/outbox_worker.go` — OutboxWorker (background flush loop)
- `orch-ref/bootstrap/desktop.go` lines 329-394 — sync client wiring + outbox + session hooks

---

## 8. AI Sessions ↔ WebSocket (AIBridgeAdapter)

### Protocol
WebSocket commands and events. The adapter subscribes to AI Bridge subscriber channels and relays `ResponseChunk` values as WebSocket events.

### Architecture
The `AIBridgeAdapter` sits between the `SessionRegistry` (which holds all AI bridges) and the `WebSocket Service`. It is session-aware: every bridge gets its own broadcast loop.

### Flow
1. `AIBridgeAdapter.Start()` registers a `MessageTypeCommand` handler on the WS service.
2. For each existing session in the registry, it calls `bridge.Subscribe()` to get a `chan ResponseChunk`, then starts a `broadcastLoop` goroutine.
3. `sessions.OnSessionCreated(callback)` auto-subscribes to future sessions.
4. Each `broadcastLoop` reads from the bridge's subscriber channel and broadcasts `ai:chunk` events to all WebSocket clients, tagged with the session ID.

### WebSocket Commands (Client -> Server)
| Command | Action |
|---------|--------|
| `ai:send` | Send a prompt to a session (resolves bridge by `session` field, default `"desktop"`) |
| `ai:stop` | Stop the running Claude CLI for a session |
| `ai:permission` | Respond to a permission request |
| `ai:question` | Answer an AskUserQuestion prompt |
| `ai:status` | Get all session info |
| `ai:sessions` | List all sessions |
| `ai:forward` | Forward a message between sessions |

### WebSocket Events (Server -> Client)
| Event | Payload |
|-------|---------|
| `ai:chunk` | `{type, content, session}` — streamed from Claude CLI |
| `ai:session_created` | `{session}` — new session registered |
| `ai:message_forwarded` | `{id, target_session, content, role, ...}` — cross-session forward |

### Key Code
```go
// app/websocket/ai_bridge.go
func (a *AIBridgeAdapter) broadcastLoop(ctx context.Context, bridge *ai.Bridge, sessionID string, ch chan ai.ResponseChunk) {
    defer bridge.Unsubscribe(ch)
    for {
        select {
        case <-ctx.Done():
            return
        case chunk, ok := <-ch:
            if !ok { return }
            a.broadcastChunk(chunk, sessionID)
        }
    }
}

func (a *AIBridgeAdapter) broadcastChunk(chunk ai.ResponseChunk, sessionID string) {
    data := map[string]interface{}{
        "type":    chunk.Type,
        "content": chunk.Content,
        "session": sessionID,
    }
    msg := NewEventMessage(generateID(), "ai:chunk", data)
    a.service.Broadcast(msg)
}
```

### Race Prevention
The subscriber channel is created via `bridge.Subscribe()` **before** the broadcast goroutine starts. This ensures `session_start` events are never missed.

```go
a.sessions.OnSessionCreated(func(id string, bridge *ai.Bridge) {
    ch := bridge.Subscribe() // subscribe before goroutine — fixes race
    go a.broadcastLoop(ctx, bridge, id, ch)
    a.broadcastSessionCreated(id)
})
```

### Source Files
- `orch-ref/app/websocket/ai_bridge.go` — AIBridgeAdapter (command routing + broadcast loops)
- `orch-ref/app/ai/bridge.go` — Bridge (Subscribe/Unsubscribe/broadcast)
- `orch-ref/app/ai/registry.go` — SessionRegistry (OnSessionCreated)

---

## 9. Discord Bot ↔ AI Sessions

### Protocol
In-process Go method calls. The Discord bot uses the same `SessionRegistry` and `Bridge` instances as the desktop.

### Flow
1. `loadDiscordConfig()` fetches Discord settings from the Settings API via HTTP (`GET http://127.0.0.1:19191/api/settings/discord.*`).
2. If enabled, `bot.NewManager(logger, cwd)` creates the bot manager.
3. A `SessionRegistry` is created and configured with default settings getter, workspace, and cross-session memory config.
4. The desktop bridge is pre-registered: `sessions.Register("desktop", d.aiBridge)`.
5. `botMgr.SetSessions(sessions)` gives the bot access to create/use AI sessions.
6. `botMgr.StartDiscord(cfg)` connects to Discord and begins handling messages.
7. Discord users create new sessions dynamically via `sessions.NewSession()`, which also triggers the `OnSessionCreated` callback to wire WebSocket broadcasting.

### Session Lifecycle
```
Discord message → bot.Manager → sessions.GetOrCreate(channelID)
                                  ↓
                              ai.Bridge.Send(prompt)
                                  ↓
                              Claude CLI subprocess
                                  ↓ (ResponseChunks)
                              broadcast to subscribers
                                  ↓
                              AIBridgeAdapter → WebSocket (desktop sees Discord sessions)
                              Discord → edit reply with streamed text
```

### Key Code
```go
// bootstrap/desktop.go — wiring Discord bot
d.sessions = ai.NewSessionRegistry(d.Logger)
d.sessions.SetStore(d.sessionStore)
d.sessions.SetDefaultSettingsGetter(d.settingsGetter)
d.sessions.SetDefaultWorkspace(mcpWs)
d.sessions.SetDefaultMemory(ai.MemoryConfig{
    Workspace:  mcpWs,
    EngineAddr: fmt.Sprintf("localhost:%d", engine.Port()),
})
d.sessions.Register("desktop", d.aiBridge)
d.botMgr.SetSessions(d.sessions)
d.botMgr.StartDiscord(cfg)
```

### Subprocess Safety
Discord bot startup is skipped when `PARENT_WEBSOCKET_PORT` is set (indicating the app is running as a subprocess of another process like `orchestra dev`), to avoid duplicate bot instances.

### Source Files
- `orch-ref/bootstrap/desktop.go` lines 822-891 — `startDiscordBot()`
- `orch-ref/bootstrap/desktop.go` lines 894-926 — `loadDiscordConfig()`

---

## 10. Workflow Events ↔ WebSocket

### Protocol
The `WorkflowBroadcaster` implements `workflow.TransitionListener` and sends WebSocket event messages to all connected clients.

### Flow
1. When a task transitions between states (e.g., `todo` -> `in-progress`), the workflow engine calls `OnTransition(event)`.
2. The broadcaster creates a `tasks:changed` WebSocket event with the full transition context.
3. The event is broadcast to all connected clients via `service.Broadcast(msg)`.

### Event Payload
```json
{
  "type": "event",
  "event": "tasks:changed",
  "data": {
    "project": "my-project",
    "epic_id": "epic-1",
    "story_id": "story-1",
    "task_id": "task-1",
    "type": "status_change",
    "from": "todo",
    "to": "in-progress",
    "time": "2026-02-26T10:30:00Z"
  }
}
```

### Key Code
```go
// app/websocket/workflow_listener.go
type WorkflowBroadcaster struct {
    logger  *slog.Logger
    service *Service
}

func (b *WorkflowBroadcaster) OnTransition(event workflow.TransitionEvent) {
    msg := NewEventMessage(generateID(), "tasks:changed", map[string]interface{}{
        "project":  event.Project,
        "epic_id":  event.EpicID,
        "story_id": event.StoryID,
        "task_id":  event.TaskID,
        "type":     event.Type,
        "from":     event.From,
        "to":       event.To,
        "time":     event.Time,
    })
    b.service.Broadcast(msg)
}
```

### Source Files
- `orch-ref/app/websocket/workflow_listener.go` — WorkflowBroadcaster

---

## 11. Preview Coordinator ↔ WebSocket

### Protocol
WebSocket commands and events for live component preview sessions. Also exposed via REST endpoints.

### Architecture
`PreviewCoordinator` manages preview sessions in memory. Each session has source code (HTML/CSS/JS/JSX), a framework type, a viewport preset, and a list of subscribed WebSocket client IDs.

### WebSocket Commands
| Command | Payload | Action |
|---------|---------|--------|
| `preview:join` | `{session_id}` | Add client to session, receive current state |
| `preview:leave` | `{session_id}` | Remove client from session |
| `preview:update` | `{session_id, html, css, js, jsx, framework}` | Update code, broadcast to session |
| `preview:viewport` | `{session_id, preset, width, height}` | Change viewport, broadcast to session |

### WebSocket Events
| Event | Trigger |
|-------|---------|
| `preview:update` | Code updated by any client in the session |
| `preview:viewport` | Viewport changed |
| `preview:session_ended` | Session deleted |
| `preview:open_browser` | Signal Chrome extension to open preview URL |

### REST Endpoints (via PreviewHTTPHandler)
| Method | Path | Action |
|--------|------|--------|
| `POST` | `/preview` | Create session, returns `{session_id, ws_url}` |
| `GET` | `/preview/:session_id` | Get session state |
| `DELETE` | `/preview/:session_id` | Delete session |

### Viewport Presets
- `mobile`: 375x667
- `tablet`: 768x1024
- `desktop`: 1280x800
- `custom`: user-defined

### Fan-out Model
Events are sent only to clients that have joined the session (`broadcastToSession`), not to all WebSocket clients. The exception is `preview:open_browser` which broadcasts globally.

### Key Code
```go
// app/websocket/preview_coordinator.go
func (c *PreviewCoordinator) broadcastToSession(sessionID, event string, data map[string]interface{}) {
    c.mu.RLock()
    session, ok := c.sessions[sessionID]
    clientIDs := make([]string, len(session.ClientIDs))
    copy(clientIDs, session.ClientIDs)
    c.mu.RUnlock()

    msg := NewEventMessage(generateID(), event, data)
    for _, clientID := range clientIDs {
        c.service.SendTo(clientID, msg)
    }
}
```

### MCP Integration
The preview coordinator is wired into MCP tools at bootstrap:
```go
// bootstrap/desktop.go
previewCoordinator := websocket.NewPreviewCoordinator(d.Logger, d.wsServer.Service())
previewCoordinator.Start()
tools.SetPreviewCoordinator(previewCoordinator)
```

### Source Files
- `orch-ref/app/websocket/preview_coordinator.go` — session management + WS broadcast
- `orch-ref/app/websocket/preview_handler.go` — REST HTTP endpoints

---

## 12. Browser Awareness ↔ AI (Chrome Extension Page Context)

### Protocol
WebSocket events from the Chrome extension to the desktop, stored in-memory, and injected into AI prompts.

### Flow

**Page Context Capture:**
1. Chrome extension sends a `page.updated` WebSocket event with the current page's content:
   ```json
   {"type":"event","event":"page.updated","data":{"content":{"title":"...","url":"...","mainContent":"...","headings":[...]}}}
   ```
2. `BrowserHandler.handleEvent()` deserializes the content into `ai.PageContext`.
3. The latest page context is stored in-memory (thread-safe via `sync.RWMutex`).
4. A `page:context_update` event is broadcast to all WebSocket clients (desktop UI can display the page info).

**Context Request:**
1. Any client sends command `page:request_context`.
2. The handler broadcasts `page:request_context` as an event to all clients.
3. The Chrome extension picks up this event and responds with a fresh `page.updated`.

**AI Integration:**
1. When the frontend sends an `ai:send` command with a `pageContext` field, the `Bridge.Send(req)` method prepends the page context to the prompt.
2. The `buildPageContextPrefix()` function formats it as XML-like tags:
   ```
   <page-context>
   URL: https://example.com
   Title: Example Page
   Page content:
   ...
   </page-context>
   ```
3. For follow-up messages, `FollowUp()` also prepends any page context.

**BrowserHandler -> Settings Server:**
The BrowserHandler is registered with the Settings server so it can be queried for the latest page context from HTTP routes:
```go
browserHandler := websocket.NewBrowserHandler(d.Logger, d.wsServer.Service())
browserHandler.Register()
d.settingsServer.SetBrowserHandler(browserHandler) // exposes GetPageContext()
```

### Key Code
```go
// app/websocket/browser_handler.go
func (h *BrowserHandler) handleEvent(_ context.Context, _ *Client, msg *Message) error {
    if msg.Event != "page.updated" { return nil }

    var pc ai.PageContext
    // ... unmarshal msg.Data["content"] into pc ...

    h.mu.Lock()
    h.pageContext = &pc
    h.mu.Unlock()

    event := NewEventMessage(generateID(), "page:context_update", map[string]interface{}{
        "title":       pc.Title,
        "url":         pc.URL,
        "mainContent": pc.MainContent,
    })
    h.service.Broadcast(event)
    return nil
}
```

```go
// app/ai/bridge.go — prepending page context to prompts
func buildPageContextPrefix(ctx *PageContext) string {
    if ctx == nil || ctx.URL == "" { return "" }
    var b strings.Builder
    b.WriteString("<page-context>\n")
    b.WriteString(fmt.Sprintf("URL: %s\n", ctx.URL))
    if ctx.MainContent != "" {
        content := ctx.MainContent
        if len(content) > 8000 { content = content[:8000] + "..." }
        b.WriteString(fmt.Sprintf("Page content:\n%s\n", content))
    }
    b.WriteString("</page-context>\n\n")
    return b.String()
}
```

### Source Files
- `orch-ref/app/websocket/browser_handler.go` — BrowserHandler (stores page context, broadcasts updates)
- `orch-ref/app/ai/bridge.go` — `buildPageContextPrefix()` (formats context for Claude)
- `orch-ref/app/ai/types.go` — `PageContext` struct

---

## Appendix A: Settings Sync (WebSocket ↔ SyncClient)

### Protocol
WebSocket commands for local reads/writes, with async push to cloud via SyncClient.

### Commands
| Command | Payload | Action |
|---------|---------|--------|
| `settings.set` | `{key, value}` | Save locally + async push to server |
| `settings.get` | `{key}` | Read from local SQLite store |
| `settings.get_all` | `{}` | Read all key-value pairs |

### Flow
1. Frontend sends `settings.set` command via WebSocket.
2. `SettingsHandler` writes to local `settings.Store` (SQLite).
3. Responds immediately with `{key, saved: true}`.
4. In a goroutine, pushes the setting to the cloud via `syncclient.Client.Push()` as an entity of type `"setting"`.
5. Entity ID is deterministic: `uuid.NewSHA1(uuid.NameSpaceURL, "setting:"+key)`.

### Source Files
- `orch-ref/app/websocket/settings_handler.go` — SettingsHandler

---

## Appendix B: Entity Sync (WebSocket ↔ SyncClient)

### Protocol
WebSocket commands for bidirectional entity sync across desktop windows and cloud.

### Supported Entities
`note`, `project`, `integration`, `ai_session`

### Commands (Client -> Server)
```
sync.note.upsert      { id, title, content, tags, pinned, version }
sync.note.delete       { id, version }
sync.project.upsert   { id, name, slug, description, path, version }
sync.project.delete    { id, version }
sync.integration.upsert { id, provider, access_token, ..., version }
sync.integration.delete { id, provider, version }
sync.ai_session.upsert { id, name, model, workspace, pinned, icon, color, version }
sync.ai_session.delete { id, version }
```

### Events (Server -> Client)
```
sync:note:upsert / sync:note:delete
sync:project:upsert / sync:project:delete
sync:integration:upsert / sync:integration:delete
sync:ai_session:upsert / sync:ai_session:delete
```

### Optimistic Broadcast
When a client sends a sync command:
1. An idempotency key is generated: `entityType:action:id:version`.
2. A `PushRecord` is created and pushed to the server asynchronously (goroutine).
3. The change is immediately broadcast to all other local WebSocket clients (optimistic update).
4. The `_sender_id` field lets clients skip applying their own echo.

### Pull Broadcast
When the pull loop receives remote changes, `SyncHandler.ApplyPullRecords()` broadcasts each record as a `sync:<entity>:<action>` event to all local WebSocket clients.

### Source Files
- `orch-ref/app/websocket/sync_handler.go` — SyncHandler
- `orch-ref/bootstrap/desktop.go` lines 338-359 — pull handler callback

---

## Appendix C: Settings Sync Bridge (Desktop ↔ Chrome Extension)

### Protocol
In-process Go callbacks to prevent infinite loops during bidirectional settings sync.

### Problem
Both the desktop settings store and the Chrome extension settings store need to stay in sync. A naive approach causes infinite loops: desktop change -> Chrome update -> desktop change -> ...

### Solution
`SyncBridge` uses a `syncing` boolean guard:

```go
// app/settings/bridge.go
func (b *SyncBridge) OnDesktopChange(key string, value any) {
    b.mu.Lock()
    if b.syncing { b.mu.Unlock(); return } // break the loop
    b.syncing = true
    b.mu.Unlock()
    defer func() { b.mu.Lock(); b.syncing = false; b.mu.Unlock() }()
    b.chrome.Set(key, value) // forward to Chrome
}

func (b *SyncBridge) OnChromeChange(key string, value any) {
    b.mu.Lock()
    if b.syncing { b.mu.Unlock(); return } // break the loop
    b.syncing = true
    b.mu.Unlock()
    defer func() { b.mu.Lock(); b.syncing = false; b.mu.Unlock() }()
    b.desktop.Set(key, value) // forward to Desktop
}
```

### Source Files
- `orch-ref/app/settings/bridge.go` — SyncBridge

---

## Appendix D: Event Bus Integration

### Protocol
The `Integration` struct bridges the WebSocket service with an in-process event bus (`events.Bus`).

### Bidirectional Flow
1. **Bus -> WebSocket:** Subscribes to specified event types. When events fire, transforms them to WebSocket messages and broadcasts (to all clients or user-specific).
2. **WebSocket -> Bus:** Registers handlers for incoming WebSocket events and commands. Publishes them to the event bus as `events.Event` (commands get a `command.` prefix).

### Configuration
```go
config := &IntegrationConfig{
    BroadcastEvents: []string{"entity.created", "entity.updated"},
    UserEvents:      []string{"notification.personal"},
    TransformEvent:  customTransformer, // or nil for default
}
integration := NewIntegration(logger, wsService, eventBus, config)
integration.Start()
```

### Source Files
- `orch-ref/app/websocket/integration.go` — WebSocket-EventBus integration

---

## Appendix E: Service Auto-Start Sequence

The `autoStartServices()` method in `bootstrap/desktop.go` orchestrates the service startup order after a 500ms delay (allowing Wails to register default service entries):

```
1. registerServiceCallbacks()
   ├── settings-api   → Already running on :19191, mark as Running
   ├── websocket      → Register start/stop callbacks for :8765
   └── grpc-engine    → Register start/stop callbacks for :50051

2. For each {websocket, grpc-engine}:
   ├── Re-probe port (another process may have bound it)
   ├── If port open → mark Running (reuse existing)
   └── If port closed → call reg.Start(id) → invoke start callback

3. Start Discord bot (unless PARENT_WEBSOCKET_PORT is set)

4. Wire WebSocket handlers:
   ├── SyncHandler.Register()
   ├── AIBridgeAdapter.Start()
   ├── PreviewCoordinator.Start()
   └── BrowserHandler.Register()
```

### Subprocess Mode
When `PARENT_WEBSOCKET_PORT` is set, the desktop app runs as a child process. It skips starting WebSocket, gRPC engine, and Discord bot — all three are managed by the parent process. Services are marked as `Running` to indicate they exist on the parent's ports.

### Source Files
- `orch-ref/bootstrap/desktop.go` lines 684-723 — `autoStartServices()`
- `orch-ref/bootstrap/desktop.go` lines 727-818 — `registerServiceCallbacks()`

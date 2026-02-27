# Architecture Overview — Orchestra Reference

> High-level system architecture, startup sequence, and proposed decoupled rebuild strategy.
> This is the entry-point document. Read this first, then dive into specific artifacts.

---

## 1. What Orchestra Is

Orchestra MCP is an AI-agentic IDE targeting 5 platforms: Desktop (Wails v3), Chrome Extension, Mobile iOS, Mobile Android, and Web Dashboard. It provides **186 MCP tools** across **34 categories** for project management, sprint planning, PRD authoring, multi-integration sync, and AI-assisted development workflows.

The system is built from three technology pillars:

| Pillar | Stack | Role |
|--------|-------|------|
| **Go Backend** | Fiber v3, GORM, pure-Go SQLite | Settings API, WebSocket server, MCP stdio server, AI bridge, sync client, integrations |
| **Rust Engine** | Tonic gRPC, Tree-sitter, Tantivy, rusqlite | Code parsing (14 languages), full-text search indexing, vector-based memory (RAG) |
| **React Frontends** | TypeScript, Zustand, Vite, Turborepo | 19 shared packages, desktop shell, Chrome extension, dashboard |

Key numbers at a glance:

| Metric | Count |
|--------|-------|
| MCP tools | 186 |
| Tool categories | 34 |
| HTTP REST endpoints | 134 |
| WebSocket event types | 6 |
| gRPC services | 4 |
| Frontend packages | 19 |
| UI components | ~200+ |
| External integrations | 10 |
| Desktop subsystems | 23 |
| Plugin capabilities | 8 interfaces |
| Workflow states | 13 |

---

## 2. Current Architecture (Monolith)

### System Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                     Desktop App (Wails v3)                               │
│  ┌────────────┐  ┌────────────┐  ┌──────────────┐  ┌───────────────┐   │
│  │  Tray Mgr  │  │  Mode Mgr  │  │ Window Mgr   │  │ Service Reg   │   │
│  │ (sys tray) │  │ (3 modes)  │  │ (multi-win)  │  │ (start/stop)  │   │
│  └────────────┘  └────────────┘  └──────────────┘  └───────────────┘   │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │              React Frontend (19 packages)                        │   │
│  │  AI Chat | Tasks | DevTools | Settings | Search | Editor | ...  │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│           │                    │                    │                    │
│     ┌─────▼──────┐      ┌─────▼──────┐      ┌─────▼──────┐            │
│     │Settings API│      │  WebSocket │      │ MCP Bridge │            │
│     │  :19191    │      │   :8765    │      │(in-process)│            │
│     └─────┬──────┘      └─────┬──────┘      └─────┬──────┘            │
│           │                   │                    │                    │
└───────────┼───────────────────┼────────────────────┼────────────────────┘
            │                   │                    │
     ┌──────▼──────┐           │              ┌─────▼──────┐
     │  AI Bridge  │           │              │ MCP Server │
     │ (Claude CLI │           │              │  (stdio)   │
     │  subprocess)│           │              │ 186 tools  │
     └──────┬──────┘           │              └─────┬──────┘
            │ stdin/stdout     │                    │ gRPC (optional)
            │ stream-json      │              ┌─────▼──────┐
     ┌──────▼──────┐    ┌─────▼──────┐       │Rust Engine │
     │  Claude CLI │    │ Sync Client│       │  :50051    │
     └─────────────┘    │  (polling) │       │ Parse+Search│
                        └─────┬──────┘       │  +Memory   │
                              │ HTTPS        └────────────┘
                        ┌─────▼──────┐
                        │Orchestra   │
                        │Web (Laravel)│
                        └────────────┘

   ┌────────────┐    ┌─────────────┐    ┌────────────┐
   │ Discord Bot│    │ Chrome Ext  │    │ Integrations│
   │ (gateway)  │    │  (WS :8765) │    │ GitHub,Jira │
   └────────────┘    └─────────────┘    │ Linear,etc  │
                                        └────────────┘
```

### The Problem

The `bootstrap/desktop.go` monolith (~1188 lines) wires approximately 20 subsystems in a single `RunDesktopWithIcons()` function with tight coupling. Specific pain points:

1. **Single-process dependency chain.** The Settings API, WebSocket server, Rust engine, AI bridge, sync client, Discord bot, all integrations hubs, DevTools manager, LSP proxy, MCP bridge, and notification service all boot within one Go binary. A crash in any subsystem can take down the entire application.

2. **Cannot run services independently.** You cannot start the Settings API without pulling in the desktop bootstrap. You cannot run the WebSocket server without the full Wails lifecycle. Testing any subsystem in isolation requires mocking the entire bootstrap chain.

3. **Port conflict workarounds.** The dual-port strategy (19191 primary, 19192 fallback proxy) is a symptom of multiple processes competing for ownership of the same service. The `lazyOpener` pattern -- polling every 100ms for up to 15 seconds waiting for the desktop app to initialize -- papers over a fundamental ordering problem.

4. **Integration hubs are in-process singletons.** GitHub, Jira, Linear, Notion, and Figma OAuth hubs are created directly in bootstrap and passed via setter injection. There is no way to restart a single integration without restarting the entire app.

5. **Dev mode vs. desktop mode are divergent paths.** The `orch --dev` startup sequence differs from the desktop startup in non-obvious ways (vite dev server, file watchers, subprocess ownership). This dual-path logic complicates testing and introduces subtle bugs.

6. **Frontend is tightly coupled to the desktop shell.** The 19 React packages assume a Wails webview environment with specific APIs (window management, TTS, notifications, screenshot). Running the same UI in a browser or Chrome extension requires shimming these native APIs.

### Startup Sequence (Desktop Mode)

The current 16-step boot sequence when `Orchestra MCP.app` launches:

```
 1. Plugin bootstrap        Discover -> resolve deps -> load -> boot -> register
 2. Settings Store          Open SQLite at ~/Library/Application Support/Orchestra/settings.db
 3. Session Store           Open SQLite for AI chat session persistence
 4. AI Bridge               Create bridge with settings getter for API keys
 5. Integration hubs        GitHub, Notion, Jira, Linear, Figma OAuth (non-blocking)
 6. MCP Bridge              Register all 186 MCP tools for HTTP access
 7. Settings API (:19191)   Start HTTP server; if port in use, proxy on :19192
 8. Sync Client             Start pull polling + outbox worker (background)
 9. DevTools                Terminal, SSH, database, logs session manager
10. LSP Proxy               Monaco editor language intelligence
11. Notification service    Resolve bundled sounds
12. Desktop window          Launch Wails app on main goroutine
 --- (after 500ms delay, in goroutine) ---
13. WebSocket Server (:8765)  Start Fiber WebSocket server
14. Rust Engine (:50051)      Start gRPC engine subprocess
15. Discord Bot               Connect gateway, register commands (if configured)
16. Auto-updater              Check for updates (if not dev mode)
```

### Services at a Glance

| # | Service | Protocol | Port | Binary | Status |
|---|---------|----------|------|--------|--------|
| 1 | MCP Server | JSON-RPC 2.0 (stdio) | -- | `bin/orchestra` | Production |
| 2 | Settings API | HTTP REST | 19191 / 19192 | embedded | Production |
| 3 | WebSocket Server | WebSocket (Fiber v3) | 8765 | embedded | Production |
| 4 | Rust Engine | gRPC (Tonic) | 50051 | `bin/orchestra-engine` | Production |
| 5 | Desktop App | Wails v3 GUI | -- | `Orchestra MCP.app` | Production |
| 6 | Discord Bot | Discord Gateway WS | -- | embedded | Optional |
| 7 | Sync Client | HTTP polling | -- | embedded | Background |

---

## 3. Proposed Decoupled Architecture

### Design Principles

1. **Each service is a standalone binary** that can run, test, and deploy independently.
2. **Services communicate via well-defined protocols** -- HTTP, WebSocket, gRPC, stdio -- never via in-process function calls between service boundaries.
3. **The MCP CLI is the backbone.** It works standalone with zero dependencies. All other services are additive layers.
4. **Desktop is a thin shell** that discovers and connects to running services, but does not own them.
5. **Each integration is a separate plugin** with its own lifecycle, credentials, and failure domain.
6. **Graceful degradation everywhere.** Missing services reduce functionality; they never crash the system.

### Decoupled System Diagram

```
┌───────────────────────────────────────────────────────────┐
│                    Service Discovery                       │
│          (port files in ~/.orchestra/services/)            │
└────┬──────────┬──────────┬──────────┬──────────┬──────────┘
     │          │          │          │          │
┌────▼────┐ ┌──▼───┐ ┌────▼────┐ ┌──▼───┐ ┌────▼────┐
│  orch   │ │orch- │ │orch-ws  │ │orch- │ │orch-    │
│  CLI    │ │settn │ │WebSocket│ │engine│ │desktop  │
│         │ │gs    │ │ :8765   │ │:50051│ │(thin    │
│ 186 MCP │ │:19191│ │         │ │      │ │ shell)  │
│ tools   │ │      │ │ Fiber   │ │ gRPC │ │         │
│ stdio   │ │SQLite│ │ 6 hdlrs │ │Parse │ │ Wails   │
│ + SSE   │ │      │ │         │ │Search│ │ Window  │
│         │ │AI Brd│ │         │ │Memory│ │ Tray    │
│ ZERO    │ │OAuth │ │         │ │      │ │ Mode    │
│ deps    │ │hubs  │ │         │ │ ZERO │ │         │
└─────────┘ └──────┘ └─────────┘ │ deps │ │discovers│
                                  └──────┘ │services │
┌─────────┐ ┌──────────┐ ┌─────────────┐  └─────────┘
│orch-sync│ │orch-     │ │orch-discord │
│background│ │chrome    │ │discord bot  │
│daemon   │ │extension │ │standalone   │
│         │ │          │ │             │
│outbox + │ │connects  │ │connects to  │
│pull poll│ │to WS     │ │settings API │
└─────────┘ └──────────┘ └─────────────┘
```

### Standalone Apps / Services

| App | What It Does | Depends On | Key Port |
|-----|-------------|------------|----------|
| `orch` | MCP tools + stdio/SSE server | Nothing (fully standalone) | -- (stdio) |
| `orch-settings` | Settings API, AI bridge, OAuth hubs, MCP HTTP bridge | SQLite (local file) | 19191 |
| `orch-ws` | WebSocket server with 6 handler types | `orch-settings` (reads config) | 8765 |
| `orch-engine` | Rust gRPC: parse, search, memory | Nothing (fully standalone) | 50051 |
| `orch-desktop` | Wails thin shell: window manager, tray, mode cycling | Discovers settings/ws/engine | -- |
| `orch-chrome` | Chrome extension UI | Connects to `orch-ws` | -- |
| `orch-sync` | Background sync daemon with outbox | `orch-settings` (auth token) | -- |
| `orch-discord` | Discord bot (gateway + slash commands) | `orch-settings` (AI bridge) | -- |

### Service Discovery

Services register themselves by writing a JSON port file to `~/.orchestra/services/`:

```
~/.orchestra/services/
  settings.json     {"port": 19191, "pid": 12345, "started_at": "..."}
  websocket.json    {"port": 8765,  "pid": 12346, "started_at": "..."}
  engine.json       {"port": 50051, "pid": 12347, "started_at": "..."}
```

Discovery protocol:
1. On startup, a service writes its port file.
2. On shutdown, it removes the file.
3. Consumers read the file to discover the address. If the file is stale (process not running), they clean it up.
4. Health checks via `GET /health` (HTTP/WS) or gRPC health service confirm liveness.
5. Fallback: environment variables (`ORCHESTRA_SETTINGS_PORT`, `ORCHESTRA_WS_PORT`, `ORCHESTRA_ENGINE_PORT`) override file-based discovery.

### Shared Contracts

All services share these contracts to stay in sync without compile-time coupling:

| Contract | Format | Location | Consumers |
|----------|--------|----------|-----------|
| Data model types | Go package | `pkg/types/` | All Go services |
| Proto definitions | `.proto` files | `proto/` | Go backend, Rust engine, TS frontends |
| WebSocket message protocol | JSON envelope spec | `pkg/wsprotocol/` | orch-ws, orch-desktop, orch-chrome |
| MCP tool schemas | JSON-RPC tool definitions | Generated from `app/tools/` | orch CLI, orch-settings (HTTP bridge) |
| Settings keys | Go constants | `pkg/settings/keys.go` | All services reading settings |
| Workflow state machine | Go package | `pkg/workflow/` | MCP tools, WebSocket broadcaster |

---

## 4. Build-First Priorities

### Phase 1: Foundation (Week 1-2)

Goal: A fully working CLI that an AI agent can use standalone.

| # | Deliverable | Description | Already Done? |
|---|------------|-------------|---------------|
| 1 | `orch` CLI | 186 MCP tools over stdio, SSE transport for browsers | Yes |
| 2 | `pkg/types/` | Extract shared Go types from `app/types/` into a standalone package | No |
| 3 | `orch-settings` | Extract Settings API into a standalone binary with SQLite, AI bridge, OAuth hubs | No |
| 4 | Service discovery | Port file mechanism in `~/.orchestra/services/` | No |

**Exit criteria:** `orch` runs with no other process. `orch-settings` runs independently and serves all 134 HTTP endpoints. Both can discover each other via port files.

### Phase 2: Desktop Shell (Week 3-4)

Goal: A thin desktop app that connects to already-running services.

| # | Deliverable | Description |
|---|------------|-------------|
| 5 | `orch-desktop` | Thin Wails shell: window manager, tray, mode manager. Discovers settings/ws/engine via port files. |
| 6 | React frontend (core) | AI chat + task management pages using the existing 19 packages |
| 7 | `orch-ws` | Standalone WebSocket server with all 6 handler types |

**Exit criteria:** Desktop launches and connects to independently running settings and WebSocket services. AI chat works end-to-end.

### Phase 3: Intelligence (Week 5-6)

Goal: Full AI capabilities with memory and search.

| # | Deliverable | Description |
|---|------------|-------------|
| 8 | `orch-engine` | Already standalone Rust binary. Add port-file registration. |
| 9 | AI bridge refactor | Move AI bridge into `orch-settings` cleanly, with cross-session memory via engine gRPC |
| 10 | Context injection | Active file context, current task context, memory retrieval integrated into AI prompts |

**Exit criteria:** AI chat has full memory. Rust engine auto-discovered by all consumers.

### Phase 4: Collaboration (Week 7-8)

Goal: Multi-device sync and real-time updates.

| # | Deliverable | Description |
|---|------------|-------------|
| 11 | `orch-sync` | Standalone background daemon: outbox worker + pull polling |
| 12 | Cloud sync | End-to-end sync with Orchestra Web API (Laravel backend) |
| 13 | Chrome extension | `orch-chrome` connects to `orch-ws` for real-time page context |

**Exit criteria:** Changes sync across desktop, Chrome, and web. Outbox survives restarts.

### Phase 5: Integrations (Week 9+)

Goal: Each integration is a standalone plugin.

| # | Deliverable | Description |
|---|------------|-------------|
| 14 | `orch-discord` | Standalone Discord bot binary |
| 15 | GitHub plugin | Standalone plugin with 17 MCP tools, OAuth, issue/PR sync |
| 16 | Jira plugin | OAuth + issue sync as standalone plugin |
| 17 | Linear plugin | OAuth + issue/cycle sync as standalone plugin |
| 18 | Notion plugin | OAuth + page push as standalone plugin |
| 19 | Figma plugin | OAuth + file/node/component pull as standalone plugin |

**Exit criteria:** Each integration runs in its own process with isolated failure. Disabling one has zero effect on others.

---

## 5. Migration Path

| Old Reference Path | New Standalone App | Key Change |
|-------------------|-------------------|------------|
| `bootstrap/desktop.go` (1188 lines) | Deleted -- replaced by thin shell | No more monolith wiring |
| `app/settings/` (20+ files) | `orch-settings` binary | Independent process, own main() |
| `app/websocket/` (12 files) | `orch-ws` binary | Independent process |
| `app/tools/` (34 categories) | `orch` CLI (already done) | Already standalone |
| `engine/` (Rust) | `orch-engine` binary | Already standalone, add port file |
| `app/desktop/` (54 files) | `orch-desktop` (thin shell) | Only Wails + window/tray/mode management |
| `app/discord/` + `app/bot/` | `orch-discord` binary | Standalone, connects to settings API |
| `app/syncclient/` | `orch-sync` daemon | Background process with outbox |
| `app/github/` | `github-plugin` | Standalone plugin with own credential store |
| `app/jira/` | `jira-plugin` | Standalone plugin |
| `app/linear/` | `linear-plugin` | Standalone plugin |
| `app/notion/` | `notion-plugin` | Standalone plugin |
| `app/figma/` | `figma-plugin` | Standalone plugin |
| `app/ai/` | Embedded in `orch-settings` | AI bridge stays co-located with settings (needs API keys) |
| `app/types/` | `pkg/types/` | Shared Go module, imported by all services |
| `app/plugins/` (23 files) | `pkg/plugins/` | Shared plugin runtime, imported by host apps |
| `app/transport/` | Stays in `orch` CLI | MCP stdio/SSE transport |

### What Gets Deleted

- `bootstrap/desktop.go` -- the monolith orchestrator
- `bootstrap/app.go` -- plugin bootstrap (moves into each binary's own main)
- `cmd/desktop/main.go` -- replaced by thin-shell `orch-desktop`
- The `lazyOpener` pattern -- no longer needed when services are independent
- Dual-port proxy (19192) -- no longer needed when only one process owns settings

### What Stays the Same

- All 186 MCP tools and their implementations
- The TOON/YAML file-based storage in `.projects/`
- The 13-state workflow state machine
- The WebSocket message protocol (JSON envelope)
- The gRPC proto definitions
- All 19 React frontend packages
- The plugin system's 4-phase lifecycle (Load, Boot, Register, Shutdown)

---

## 6. Key Architectural Decisions

### File-Based Storage as Default

The MCP tools store all project data as YAML frontmatter + Markdown body files in `.projects/`. This is deliberate:

- Works offline with zero infrastructure
- Human-readable and git-trackable
- The Rust engine provides optional vector search and indexing on top
- PostgreSQL (cloud) and SQLite (local desktop) are additive layers for sync, not replacements

### Fallback-First Design

Every optional dependency degrades gracefully:

| Component | When Missing | Fallback |
|-----------|-------------|----------|
| Rust engine | Binary not found | TOON/markdown file storage for memory |
| AI bridge | No API keys configured | Tools still work, chat disabled |
| Sync client | No auth token | Local-only operation |
| Discord bot | Not configured | No Discord presence |
| Any integration | OAuth not completed | Integration tools return "not connected" |

### Plugin System

The plugin architecture (8 capability interfaces, 4-phase lifecycle, topological dependency sort, DI container) is designed so that every feature -- including core features -- is a plugin. This means:

- Features can be disabled at runtime via feature flags
- Third-party plugins follow the same contracts as first-party code
- The host app (desktop, CLI, etc.) is just a plugin host with a main loop

---

## 7. Communication Patterns Summary

All inter-service communication uses one of four patterns:

| Pattern | Used Between | Details |
|---------|-------------|---------|
| **JSON-RPC 2.0 over stdio** | AI agents <-> MCP Server | Line-delimited, 10MB max message, bidirectional |
| **HTTP REST** | Frontend <-> Settings API, Sync client <-> Orchestra Web | 134 endpoints on :19191, CORS enabled |
| **WebSocket** | Frontend <-> WS Server, Chrome ext <-> WS Server | JSON envelope with 6 message types, ping/pong keepalive |
| **gRPC** | Go backend <-> Rust Engine | 4 services (Health, Parse, Search, Memory), protobuf |

Additionally, the AI bridge communicates with Claude CLI via **stdin/stdout stream-json** format, and the Discord bot connects to the **Discord Gateway WebSocket (v10)**.

---

## 8. Data Model Summary

Data is stored across three tiers:

| Tier | Technology | Location | Purpose |
|------|-----------|----------|---------|
| **File-based** | TOON/YAML + Markdown | `.projects/{slug}/` | Projects, epics, stories, tasks, sprints, notes, PRDs, templates, memory chunks, session logs |
| **Local SQLite** | pure-Go sqlite | `~/Library/Application Support/Orchestra/` | Settings (key-value), AI chat sessions, sync outbox, integration credentials (AES-256-GCM encrypted) |
| **Cloud PostgreSQL** | pgvector, JSONB, tsvector | Orchestra Web API | Source of truth for sync, user accounts, team data |

Key entity counts in the data model:

| Category | Entity Types |
|----------|-------------|
| Project management | ProjectStatus, IssueData (epic/story/task/bug), Sprint, Template, Retrospective |
| PRD system | PrdSession, PrdAnswer, PrdQuestion, PrdTemplate, AgentBriefing |
| Memory / sessions | MemoryChunk, MemoryIndex, SessionLog, SessionEvent |
| Notes | Note (YAML frontmatter + markdown body) |
| Notifications | Notification, NotificationAction, NotificationPreferences, FilterRule |
| Integrations | GitHub (issues, PRs, CI), Jira (issues, transitions), Linear (issues, teams, cycles), Notion (pages), Figma (files, nodes, components) |
| Sync protocol | PullRecord, PushRecord, OutboxRecord, DeviceRegistration |
| gRPC messages | MemorySession, Observation, SearchResult, ParseResult |

---

## 9. Integration Map Summary

10 external integrations, all using AES-256-GCM encrypted credential storage:

| Integration | Auth | Sync Direction | MCP Tools |
|------------|------|---------------|-----------|
| GitHub | OAuth 2.0 / PAT | Bidirectional (issues, PRs, CI) | 17 tools |
| Jira | OAuth 2.0 (3LO) | Bidirectional (issues) | -- |
| Linear | OAuth 2.0 | Bidirectional (issues, teams, cycles) | -- |
| Notion | OAuth 2.0 | Push only (pages) | -- |
| Figma | OAuth 2.0 PKCE / PAT | Pull only (files, nodes) | Via MCP bridge |
| Discord | Bot token | Push only (notifications) | Accessible via bot commands |
| Slack | Bot token + App token | Push only (notifications) | Accessible via bot commands |
| Firebase | Service account JSON | Push only (notifications, analytics) | -- |
| Orchestra Web | Laravel Sanctum | Bidirectional (sync protocol) | -- |
| Apple Notes | AppleScript (macOS) | Push only (notes) | -- |

---

## 10. Artifact Index

| # | Artifact | Lines | What It Covers |
|---|---------|-------|---------------|
| **00** | **Architecture Overview** (this document) | -- | System diagram, monolith problems, decoupled rebuild strategy, migration path |
| **01** | [Service Catalog](01-service-catalog.md) | 707 | 7 services: MCP server, Settings API, WebSocket, Rust engine, Desktop app, Discord bot, Sync client. Protocols, ports, startup sequences, source files |
| **02** | [Data Model](02-data-model.md) | 2510 | Every entity type across file-based (TOON/YAML), SQLite, PostgreSQL, and gRPC. Fields, relationships, storage locations, workflow states |
| **03** | [API Surface](03-api-surface.md) | 739 | All 134 HTTP REST endpoints, WebSocket commands/events, gRPC RPCs, MCP JSON-RPC methods. Request/response schemas |
| **04** | [MCP Tool Catalog](04-mcp-tool-catalog.md) | 2534 | All 186 MCP tools across 34 categories with full parameter detail. Plus 3 resources and 3 prompts |
| **05** | [Integration Map](05-integration-map.md) | 860 | 10 external integrations: OAuth flows, credential storage, API clients, data exchange, MCP tools per integration |
| **06** | [Communication Patterns](06-communication-patterns.md) | 998 | How every service talks to every other service. 9 communication patterns with code examples, message formats, failure modes |
| **07** | [Desktop Subsystems](07-desktop-subsystems.md) | 1212 | 23 desktop subsystems across 54 Go files. Window manager, mode manager, tray, spirit/bubble windows, hotkeys, permissions, screenshots |
| **08** | [Plugin System](08-plugin-system.md) | 1167 | 4-phase lifecycle, manifest schema, 8 capability interfaces, dependency resolution (Kahn's algorithm), IPC over Unix sockets, event bus, DI container |
| **09** | [UI Inventory](09-ui-inventory.md) | 1174 | 19 frontend packages: ~200+ components, stores, hooks, dependencies. AI chat (50+ event cards), editor, tasks, devtools, themes, icons |
| **10** | [Feature Matrix](10-feature-matrix.md) | 586 | Feature-to-service mapping table. 5 complexity tiers. 7-phase build priority recommendation |

**Total reference material: ~12,487 lines across 10 artifacts.**

---

## 11. Architecture Artifacts

The following artifacts document the architecture decisions:

| # | Artifact | What It Covers |
|---|---------|---------------|
| **11** | [QUIC Mesh Architecture](11-quic-mesh-architecture.md) | Event-driven QUIC mesh design: Protobuf envelope, connection protocol, event flow, service discovery, outbox, ordering |
| **12** | [Monorepo Structure](12-monorepo-structure.md) | Complete directory tree, proto generation for 6 languages, Go workspace, Makefile, Docker Compose |
| **13** | [Rebuild Decisions](13-rebuild-decisions.md) | Final stack, all architecture decisions, what changes from reference, what we keep, build order |
| **14** | [Plugin Host Architecture](14-plugin-host-architecture.md) | Plugin host pattern, Protobuf contract, orchestrator design, star topology |
| **15** | [Phase 1 Implementation](15-phase1-implementation.md) | Step-by-step build plan: QUIC, proto, SDK, orchestrator, plugins |
| **16** | [Feature-Driven Workflow](16-feature-driven-workflow.md) | Feature docs replace Scrum, doc-first TDD, human review gate |
| **17** | [AI Reasoning & Memory](17-ai-reasoning-memory.md) | Project inception, guided reasoning, per-project RAG, memory system |

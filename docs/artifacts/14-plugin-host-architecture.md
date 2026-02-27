# Plugin Host Architecture — Orchestra

> The orchestrator pattern. A minimal host routes everything. Every feature is an out-of-process plugin.
> Plugins communicate via QUIC + Protobuf. Any language can be a plugin.
> The orchestrator is the ONLY router — every message goes through it.

---

## 1. Core Concept

```
═══════════════════════════════════════════════════════════════
                    THE FACTORY PATTERN
═══════════════════════════════════════════════════════════════

There is no "MCP server" or "API server" or "sync daemon."

There is ONE thing: the Orchestrator.
Everything else is a Plugin.

The Orchestrator:
  1. Reads manifests to know which plugins exist
  2. Starts plugins (or connects to running ones)
  3. Routes ALL messages between plugins via QUIC
  4. That's it. ~500 lines of Go.

A Plugin:
  1. Is a standalone binary (any language)
  2. Implements the PluginHost Protobuf service
  3. Declares what it PROVIDES and what it NEEDS
  4. Communicates ONLY with the orchestrator, never directly with other plugins

═══════════════════════════════════════════════════════════════
```

---

## 2. System Diagram

```
                        ┌─────────────────────┐
                        │    ORCHESTRATOR      │
                        │                     │
                        │  - Plugin loader    │
                        │  - Message router   │
                        │  - Event dispatcher │
                        │  - Lifecycle mgr    │
                        │  - QUIC listener    │
                        │                     │
                        │  (~500 lines of Go) │
                        └──────────┬──────────┘
                                   │
                    QUIC + Protobuf │ (every connection)
                                   │
         ┌─────────────────────────┼─────────────────────────┐
         │              │          │          │               │
    ┌────▼────┐   ┌────▼────┐ ┌──▼───┐ ┌───▼────┐   ┌─────▼─────┐
    │transport│   │transport│ │tools │ │tools   │   │ storage   │
    │.stdio   │   │.web     │ │.feat │ │.memory │   │.markdown  │
    │         │   │transport│ │ures  │ │        │   │           │
    │ JSON-RPC│   │.web     │ │ ~35  │ │ 6 tools│   │ Proto+MD  │
    │ stdin/  │   │WebTrans │ │tools │ │        │   │ files     │
    │ stdout  │   │         │ │      │ │        │   │           │
    │  (Go)   │   │  (Go)   │ │ (Go) │ │  (Go)  │   │   (Go)    │
    └─────────┘   └─────────┘ └──────┘ └────────┘   └───────────┘

    ┌─────────┐   ┌─────────┐ ┌──────────┐ ┌──────────────────┐
    │ storage │   │ storage │ │ engine   │ │ engine           │
    │.sqlite  │   │.postgres│ │.parse    │ │.search           │
    │         │   │         │ │          │ │                  │
    │ local   │   │ cloud   │ │Tree-sit  │ │ Tantivy          │
    │ SQLite  │   │ pgvector│ │14 langs  │ │ full-text        │
    │  (Go)   │   │  (Go)   │ │ (Rust)   │ │ (Rust)           │
    └─────────┘   └─────────┘ └──────────┘ └──────────────────┘

    ┌─────────┐   ┌─────────┐ ┌──────────┐ ┌──────────────────┐
    │ engine  │   │ai.claude│ │ai.openai │ │ integration      │
    │.vectors │   │         │ │          │ │ .github          │
    │         │   │ Claude  │ │ OpenAI   │ │                  │
    │ LanceDB │   │ API     │ │ API      │ │ OAuth + 17 tools │
    │ (Rust)  │   │  (Go)   │ │  (Go)    │ │   (Go)           │
    └─────────┘   └─────────┘ └──────────┘ └──────────────────┘

    ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐
    │ integration │  │ integration │  │ ui.macos            │
    │ .jira       │  │ .linear     │  │                     │
    │             │  │             │  │ SwiftUI native app  │
    │  (Go)       │  │  (Go)       │  │   (Swift)           │
    └─────────────┘  └─────────────┘  └─────────────────────┘
```

**Key insight**: The orchestrator is a star topology, not a mesh. Every arrow goes through the center. This gives total observability, tracing, and control.

---

## 3. Protobuf Plugin Contract

This is the single most important file in the entire system. Every plugin in any language implements this service.

```protobuf
syntax = "proto3";
package orchestra.plugin.v1;

import "google/protobuf/any.proto";
import "google/protobuf/timestamp.proto";
import "google/protobuf/struct.proto";

// ================================================================
// PluginHost — the service every plugin implements
// ================================================================

service PluginHost {
  // Lifecycle
  rpc Register (PluginManifest) returns (RegistrationResult);
  rpc Boot (BootRequest) returns (BootResult);
  rpc Shutdown (ShutdownRequest) returns (ShutdownResult);
  rpc Health (HealthRequest) returns (HealthResult);

  // Tool handling (for plugins that provide tools)
  rpc HandleToolCall (ToolRequest) returns (ToolResponse);
  rpc ListTools (ListToolsRequest) returns (ListToolsResponse);

  // Event handling (for plugins that subscribe to events)
  rpc HandleEvent (EventEnvelope) returns (EventAck);

  // Storage operations (for storage plugins)
  rpc StorageRead (StorageReadRequest) returns (StorageReadResponse);
  rpc StorageWrite (StorageWriteRequest) returns (StorageWriteResponse);
  rpc StorageDelete (StorageDeleteRequest) returns (StorageDeleteResponse);
  rpc StorageList (StorageListRequest) returns (StorageListResponse);
  rpc StorageQuery (StorageQueryRequest) returns (StorageQueryResponse);

  // AI operations (for AI provider plugins)
  rpc AIChat (AIChatRequest) returns (stream AIChatChunk);
  rpc AIEmbed (AIEmbedRequest) returns (AIEmbedResponse);

  // Transport operations (for transport plugins)
  rpc TransportSend (TransportSendRequest) returns (TransportSendResponse);
  rpc TransportReceive (stream TransportReceiveRequest) returns (stream TransportReceiveResponse);
}

// ================================================================
// Plugin Manifest — declares provides/needs
// ================================================================

message PluginManifest {
  string id = 1;                          // "tools.features", "storage.markdown", "ai.claude"
  string version = 2;                     // semver
  string language = 3;                    // "go", "rust", "swift", "kotlin", "csharp", "typescript"

  // What this plugin provides
  repeated string provides_tools = 4;     // tool names: ["create_feature", "list_features", ...]
  repeated string provides_events = 5;    // event topics: ["feature.created", "feature.updated", ...]
  repeated string provides_storage = 6;   // storage types: ["markdown", "sqlite", "postgres"]
  repeated string provides_transport = 7; // transport types: ["stdio", "webtransport", "quic"]
  repeated string provides_ai = 8;        // AI providers: ["claude", "openai"]

  // What this plugin needs from other plugins (via orchestrator)
  repeated string needs_storage = 9;      // ["read", "write", "list", "delete", "query"]
  repeated string needs_events = 10;      // ["publish", "subscribe"]
  repeated string needs_ai = 11;          // ["chat", "embed"]
  repeated string needs_tools = 12;       // tool names from other plugins it wants to call

  // Plugin metadata
  string description = 13;
  string author = 14;
  string binary = 15;                     // path to binary (for auto-start)
  repeated string args = 16;              // CLI arguments
  map<string, string> env = 17;           // environment variables
}

message RegistrationResult {
  bool accepted = 1;
  string reject_reason = 2;
  string assigned_port = 3;              // QUIC port assigned by orchestrator
}

// ================================================================
// Tool Call Flow
// ================================================================

message ToolRequest {
  string request_id = 1;                 // UUIDv7
  string tool_name = 2;                  // "create_feature"
  google.protobuf.Struct arguments = 3;  // tool arguments as JSON-like struct
  string caller_plugin = 4;             // which plugin initiated the call
  string trace_parent = 5;              // W3C traceparent
}

message ToolResponse {
  string request_id = 1;
  bool success = 2;
  google.protobuf.Struct result = 3;    // tool result
  string error_code = 4;
  string error_message = 5;
}

message ListToolsRequest {}

message ListToolsResponse {
  repeated ToolDefinition tools = 1;
}

message ToolDefinition {
  string name = 1;
  string description = 2;
  google.protobuf.Struct input_schema = 3;  // JSON Schema as Struct
}

// ================================================================
// Event Flow
// ================================================================

message EventEnvelope {
  string id = 1;                          // UUIDv7
  string topic = 2;                       // "feature.created"
  string source_plugin = 3;              // "tools.features"
  google.protobuf.Timestamp timestamp = 4;
  google.protobuf.Any payload = 5;
  map<string, string> metadata = 6;
  string trace_parent = 7;
}

message EventAck {
  bool processed = 1;
  string error = 2;
}

// ================================================================
// Storage Operations
// ================================================================

message StorageReadRequest {
  string path = 1;                       // "projects/my-app/tasks/task-001.md"
  string storage_type = 2;              // "markdown", "sqlite", "postgres" (empty = default)
}

message StorageReadResponse {
  bytes content = 1;
  map<string, string> metadata = 2;     // Protobuf-encoded structured fields as key-value
  int64 version = 3;
}

message StorageWriteRequest {
  string path = 1;
  bytes content = 2;
  map<string, string> metadata = 3;
  int64 expected_version = 4;           // 0 = create new, >0 = update with CAS
  string storage_type = 5;
}

message StorageWriteResponse {
  bool success = 1;
  int64 new_version = 2;
  string error = 3;
}

message StorageDeleteRequest {
  string path = 1;
  string storage_type = 2;
}

message StorageDeleteResponse {
  bool success = 1;
}

message StorageListRequest {
  string prefix = 1;                    // "projects/my-app/tasks/"
  string pattern = 2;                   // "*.md"
  string storage_type = 3;
}

message StorageListResponse {
  repeated StorageEntry entries = 1;
}

message StorageEntry {
  string path = 1;
  int64 size = 2;
  int64 version = 3;
  google.protobuf.Timestamp modified_at = 4;
}

message StorageQueryRequest {
  string storage_type = 1;             // must be "sqlite" or "postgres"
  string query = 2;                    // SQL or structured query
  repeated google.protobuf.Value params = 3;
}

message StorageQueryResponse {
  repeated google.protobuf.Struct rows = 1;
  int64 affected_rows = 2;
}

// ================================================================
// AI Operations
// ================================================================

message AIChatRequest {
  string provider = 1;                  // "claude", "openai"
  string model = 2;                     // "claude-opus-4-6", "gpt-4o"
  repeated ChatMessage messages = 3;
  repeated ToolDefinition tools = 4;    // tools available to the AI
  int32 max_tokens = 5;
  float temperature = 6;
  string system_prompt = 7;
}

message ChatMessage {
  string role = 1;                      // "user", "assistant", "system", "tool"
  string content = 2;
  string tool_call_id = 3;
  string tool_name = 4;
}

message AIChatChunk {
  string type = 1;                      // "text", "tool_use", "done", "error"
  string content = 2;
  string tool_name = 3;
  string tool_call_id = 4;
  google.protobuf.Struct tool_arguments = 5;
}

message AIEmbedRequest {
  string provider = 1;
  string model = 2;
  repeated string texts = 3;
}

message AIEmbedResponse {
  repeated Embedding embeddings = 1;
}

message Embedding {
  repeated float values = 1;
  int32 dimensions = 2;
}

// ================================================================
// Transport Operations
// ================================================================

message TransportSendRequest {
  string transport_type = 1;           // "stdio", "webtransport"
  bytes data = 2;                      // raw bytes to send
  string destination = 3;             // client ID or empty for broadcast
}

message TransportSendResponse {
  bool sent = 1;
}

message TransportReceiveRequest {
  string transport_type = 1;
}

message TransportReceiveResponse {
  bytes data = 1;
  string source = 2;                  // client ID or transport identifier
}

// ================================================================
// Lifecycle
// ================================================================

message BootRequest {
  map<string, string> config = 1;      // runtime config from orchestrator
}

message BootResult {
  bool ready = 1;
  string error = 2;
}

message ShutdownRequest {
  int32 timeout_seconds = 1;
}

message ShutdownResult {
  bool clean = 1;
}

message HealthRequest {}

message HealthResult {
  ServiceStatus status = 1;
  string message = 2;
  map<string, string> details = 3;
}

enum ServiceStatus {
  UNKNOWN = 0;
  HEALTHY = 1;
  DEGRADED = 2;
  UNHEALTHY = 3;
}
```

---

## 4. Message Flow: create_feature End-to-End

```
Agent (Claude)
  │
  │ stdin: {"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_feature",...}}
  │
  ▼
┌─────────────────┐
│ transport.stdio  │  (Go plugin, receives stdin)
│                 │
│ Parses JSON-RPC │
│ Sends to orch:  │
│ TransportReceive│ ──────────────────────────────────────┐
└─────────────────┘                                        │
                                                           ▼
                                              ┌─────────────────────┐
                                              │   ORCHESTRATOR      │
                                              │                     │
                                              │ 1. Receives from    │
                                              │    transport.stdio  │
                                              │                     │
                                              │ 2. Parses: tool     │
                                              │    "create_feature" │
                                              │    belongs to       │
                                              │    tools.features   │
                                              │                     │
                                              │ 3. Routes:          │
                                              │    HandleToolCall   │
                                              │    → tools.features │
                                              └──────────┬──────────┘
                                                         │
                                                         ▼
                                              ┌─────────────────────┐
                                              │   tools.features    │
                                              │                     │
                                              │ 4. Builds feature   │
                                              │    doc data         │
                                              │                     │
                                              │ 5. Needs to write   │
                                              │    file. Calls back │
                                              │    → orchestrator:  │
                                              │    StorageWrite     │
                                              └──────────┬──────────┘
                                                         │
                                                         ▼
                                              ┌─────────────────────┐
                                              │   ORCHESTRATOR      │
                                              │                     │
                                              │ 6. Routes           │
                                              │    StorageWrite     │
                                              │    → storage.markdown   │
                                              └──────────┬──────────┘
                                                         │
                                                         ▼
                                              ┌─────────────────────┐
                                              │   storage.markdown      │
                                              │                     │
                                              │ 7. Writes Proto+MD  │
                                              │    file to disk     │
                                              │                     │
                                              │ 8. Returns OK +     │
                                              │    new version      │
                                              └──────────┬──────────┘
                                                         │
                                                         ▼
                                              ┌─────────────────────┐
                                              │   ORCHESTRATOR      │
                                              │                     │
                                              │ 9. Returns storage  │
                                              │    result to        │
                                              │    tools.features      │
                                              └──────────┬──────────┘
                                                         │
                                                         ▼
                                              ┌─────────────────────┐
                                              │   tools.features    │
                                              │                     │
                                              │ 10. Wants to publish│
                                              │  "feature.created"  │
                                              │     Calls back      │
                                              │     → orchestrator  │
                                              └──────────┬──────────┘
                                                         │
                                                         ▼
                                              ┌─────────────────────┐
                                              │   ORCHESTRATOR      │
                                              │                     │
                                              │ 11. Broadcasts      │
                                              │  "feature.created"  │
                                              │     to all plugins  │
                                              │     subscribed to   │
                                              │     "feature.*":    │
                                              │     - engine.search │
                                              │     - storage.sync  │
                                              │     - ui.macos      │
                                              └──────────┬──────────┘
                                                         │
                                              (parallel fan-out)
                                              ┌──────┬───┴───┬──────┐
                                              ▼      ▼       ▼      │
                                           search  sync   macos     │
                                           indexes queues  updates   │
                                           feature for     UI        │
                                                   cloud             │
                                                                     │
                                              ┌──────────────────────┘
                                              │
                                              ▼
                                              ┌─────────────────────┐
                                              │   tools.features    │
                                              │                     │
                                              │ 12. Returns tool    │
                                              │     result to       │
                                              │     orchestrator    │
                                              └──────────┬──────────┘
                                                         │
                                                         ▼
                                              ┌─────────────────────┐
                                              │   ORCHESTRATOR      │
                                              │                     │
                                              │ 13. Routes result   │
                                              │     back to         │
                                              │     transport.stdio │
                                              └──────────┬──────────┘
                                                         │
                                                         ▼
                                              ┌─────────────────────┐
                                              │ transport.stdio      │
                                              │                     │
                                              │ 14. Formats JSON-RPC│
                                              │     response        │
                                              │     Writes stdout   │
                                              └─────────────────────┘
                                                         │
                                                         ▼
                                                  Agent receives result
```

---

## 5. Plugin Categories

### Transport Plugins (how data enters/exits the system)

| Plugin | Language | What It Does |
|--------|----------|-------------|
| `transport.stdio` | Go | MCP JSON-RPC over stdin/stdout (for AI agents) |
| `transport.sse` | Go | Server-Sent Events for browser MCP clients |
| `transport.webtransport` | Go | WebTransport/HTTP/3 for web apps |

### Storage Plugins (where data lives)

| Plugin | Language | What It Does |
|--------|----------|-------------|
| `storage.markdown` | Go | Protobuf metadata + Markdown body files in `.projects/` |
| `storage.sqlite` | Go | Local SQLite for settings, sessions, outbox |
| `storage.postgres` | Go | Cloud PostgreSQL + pgvector for sync |
| `storage.redis` | Go | Redis Streams for pub/sub, cache, rate limiting |

### Tool Plugins (what the system can do)

| Plugin | Language | Tools |
|--------|----------|-------|
| `tools.features` | Go | ~58 tools: feature CRUD, workflow, agentops, time, deps, git, notifications, quality |
| `tools.memory` | Go | 6 tools: save_memory, search_memory, get_context, sessions |
| `tools.notes` | Go | 6 tools: notes CRUD + search |
| `tools.prd` | Go | 14 tools: PRD authoring with templates |
| `tools.devtools` | Go | 25 tools: terminal, SSH, database, logs, debugger |
| `tools.claude` | Go | 6 tools: skills, agents, hooks management |
| `tools.usage` | Go | 3 tools: record/get/reset usage |

### Engine Plugins (heavy computation, Rust)

| Plugin | Language | What It Does |
|--------|----------|-------------|
| `engine.parse` | Rust | Tree-sitter code parsing (14 languages) |
| `engine.search` | Rust | Tantivy full-text search indexing |
| `engine.vectors` | Rust | LanceDB vector store for embeddings/RAG |

### AI Plugins (LLM providers)

| Plugin | Language | What It Does |
|--------|----------|-------------|
| `ai.claude` | Go | Anthropic Claude API (chat + tools) |
| `ai.openai` | Go | OpenAI API (chat + embeddings) |

### Integration Plugins (external services)

| Plugin | Language | What It Does |
|--------|----------|-------------|
| `integration.github` | Go | OAuth + 17 MCP tools (issues, PRs, repos, CI) |
| `integration.jira` | Go | OAuth 3LO + issue CRUD + transitions |
| `integration.linear` | Go | OAuth + GraphQL (issues, teams, cycles) |
| `integration.notion` | Go | OAuth + page push (markdown-to-blocks) |
| `integration.figma` | Go | OAuth PKCE + 4 MCP tools |
| `integration.discord` | Go | Bot gateway + slash commands + AI chat |
| `integration.slack` | Go | Bot + app token, Socket Mode |

### UI Plugins (native frontends)

| Plugin | Language | What It Does |
|--------|----------|-------------|
| `ui.macos` | Swift | SwiftUI macOS native app |
| `ui.ios` | Swift | SwiftUI iOS app (shares OrchestraKit) |
| `ui.android` | Kotlin | Compose Android app |
| `ui.windows` | C# | WinUI 3 Windows app |
| `ui.linux` | Go | GTK4 Linux app |
| `ui.web` | TS/React | Web dashboard (connects via transport.webtransport) |
| `ui.chrome` | TS/React | Chrome extension |

---

## 6. Orchestrator Design

The orchestrator is intentionally tiny. It does 4 things:

### 6.1 Plugin Loader

```
1. Read ~/.orchestra/plugins.yaml (or manifest per plugin)
2. For each enabled plugin:
   a. Start the binary (if auto_start)
   b. Open QUIC connection
   c. Call Register(PluginManifest)
   d. Validate provides/needs can be satisfied
   e. Call Boot(config)
3. Build routing table:
   - tool_name → plugin_id
   - event_topic → [plugin_ids]
   - storage_type → plugin_id
   - transport_type → plugin_id
   - ai_provider → plugin_id
```

### 6.2 Message Router

Every message type has a routing rule:

| Message Type | Routing Logic |
|-------------|---------------|
| `HandleToolCall(name)` | Look up `name` in tool routing table → send to owning plugin |
| `StorageRead/Write/Delete/List` | Look up `storage_type` → send to storage plugin |
| `StorageQuery` | Look up `storage_type` (must be sqlite/postgres) → send to storage plugin |
| `EventEnvelope(topic)` | Match `topic` against all subscriptions → fan-out to all matching plugins |
| `AIChatRequest(provider)` | Look up `provider` → send to AI plugin |
| `TransportSend(type)` | Look up `transport_type` → send to transport plugin |

### 6.3 Event Dispatcher

When a plugin publishes an event:
1. Plugin calls orchestrator with `EventEnvelope`
2. Orchestrator matches `topic` against all plugin subscriptions (wildcard matching)
3. Fan-out: send `HandleEvent(envelope)` to each matching plugin in parallel
4. Collect acks (non-blocking — publisher doesn't wait)

### 6.4 Lifecycle Manager

```
Boot order (topological sort by dependencies):
  1. storage.markdown      (no deps — provides storage)
  2. storage.sqlite    (no deps — provides storage)
  3. transport.stdio   (no deps — provides transport)
  4. tools.features   (needs: storage)
  5. tools.memory      (needs: storage)
  6. engine.parse      (no deps)
  7. engine.search     (no deps)
  8. ai.claude         (no deps)
  9. integration.github (needs: storage, events)
  10. ui.macos          (needs: tools, events)

Shutdown: reverse order
```

---

## 7. What This Changes from the Previous Design

| Before (artifact 11-12) | Now |
|-------------------------|-----|
| 5 hardcoded services (orch-mcp, orch-server, orch-engine, orch-sync, orch-gateway) | 1 orchestrator + N plugins |
| Services talk directly to each other (mesh) | All communication goes through orchestrator |
| Each service is a monolith within its domain | Each capability is a separate plugin/binary |
| Adding a new integration = code in orch-server | Adding a new integration = new plugin binary |
| Replacing storage = rewrite service internals | Replacing storage = load different storage plugin |
| Star topology emerged from mesh | Star topology is the design |

---

## 8. Configuration: `~/.orchestra/plugins.yaml`

```yaml
# Which plugins to load and how
orchestrator:
  quic_port: 0                    # 0 = OS-assigned
  log_level: info
  trace_enabled: true

plugins:
  # Transport
  transport.stdio:
    enabled: true
    binary: orch-transport-stdio
    auto_start: true

  transport.webtransport:
    enabled: true
    binary: orch-transport-web
    auto_start: true
    config:
      port: 4433
      tls_cert: ~/.orchestra/certs/gateway.crt
      tls_key: ~/.orchestra/certs/gateway.key

  # Storage
  storage.markdown:
    enabled: true
    binary: orch-storage-markdown
    auto_start: true
    config:
      workspace: .

  storage.sqlite:
    enabled: true
    binary: orch-storage-sqlite
    auto_start: true
    config:
      path: ~/.orchestra/data/orchestra.db

  storage.postgres:
    enabled: false                 # enable when cloud sync needed
    binary: orch-storage-postgres
    config:
      dsn: postgres://user:pass@host/db

  # Tools
  tools.features:
    enabled: true
    binary: orch-tools-features
    auto_start: true

  tools.memory:
    enabled: true
    binary: orch-tools-memory
    auto_start: true

  # Engine (Rust)
  engine.parse:
    enabled: true
    binary: orch-engine-parse
    auto_start: true

  engine.search:
    enabled: true
    binary: orch-engine-search
    auto_start: true

  engine.vectors:
    enabled: true
    binary: orch-engine-vectors
    auto_start: true
    config:
      db_path: ~/.orchestra/data/vectors

  # AI
  ai.claude:
    enabled: true
    binary: orch-ai-claude
    auto_start: true
    config:
      api_key_env: ANTHROPIC_API_KEY

  # Integrations (disabled by default, enable per-user)
  integration.github:
    enabled: false
    binary: orch-integration-github

  integration.jira:
    enabled: false
    binary: orch-integration-jira

  # UI (platform-specific)
  ui.macos:
    enabled: true
    binary: Orchestra.app/Contents/MacOS/Orchestra
    auto_start: false              # user launches manually
```

---

## 9. Build Order (Updated)

### Phase 0: The Foundation
1. **Protobuf contract** (`proto/orchestra/plugin/v1/plugin.proto`) — the PluginHost service
2. **Orchestrator** (`services/orchestrator/`) — ~500 lines, plugin loader + message router
3. **libs/go/plugin/** — shared Go library for building plugins (QUIC client, manifest, lifecycle)

### Phase 1: Minimal Working System
4. **storage.markdown** — first storage plugin (Protobuf metadata + Markdown files)
5. **transport.stdio** — first transport plugin (MCP JSON-RPC)
6. **tools.features** — first tool plugin (~58 feature-driven tools)

**Exit criteria**: Agent can call `create_feature` via stdio, orchestrator routes to tools.features, tools.features writes via storage.markdown (Protobuf metadata + Markdown body). End-to-end working.

### Phase 2: Web Access
7. **transport.webtransport** — WebTransport for browsers
8. **ui.web** — React web dashboard

### Phase 3: Intelligence
9. **engine.parse** (Rust) — Tree-sitter
10. **engine.search** (Rust) — Tantivy
11. **engine.vectors** (Rust) — LanceDB
12. **tools.memory** — memory tools (uses engine.vectors)
13. **ai.claude** — Claude API provider

### Phase 4: macOS Native
14. **ui.macos** — SwiftUI native app

### Phase 5: Integrations
15. **integration.github** — first integration plugin
16. **storage.sqlite** — local structured storage
17. **storage.postgres** — cloud sync storage
18. More integrations as needed

### Phase 6: More Platforms
19. **ui.ios**, **ui.android**, **ui.windows**, **ui.linux**

---

## 10. Why This Is Better

1. **True isolation**: One plugin crashes, others keep running. The orchestrator restarts it.
2. **Any language**: Rust plugins, Swift plugins, Go plugins — all implement the same Protobuf service.
3. **Hot-swappable**: Stop `storage.markdown`, start `storage.postgres`. No code changes needed.
4. **Testable**: Test any plugin in isolation by mocking the orchestrator.
5. **Observable**: Every message goes through the orchestrator. Full tracing. Full logging.
6. **Scalable**: Today it's 15 plugins on a laptop. Tomorrow it's 50 plugins across a cluster.
7. **Open**: Third-party plugins are first-class citizens. Same contract as core plugins.

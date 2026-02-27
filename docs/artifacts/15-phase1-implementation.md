# Phase 1 Implementation — Smallest Usable Win

> Build the minimum system an AI agent can connect to via MCP stdio and manage features.
> Ship something usable fast. Then push capabilities module by module.
> Each step produces a runnable artifact you can test.
> QUIC from day one. No shortcuts on the stack.

---

## The Goal

```
Agent (Claude, Cursor, etc.)
  │
  │ stdin: {"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_feature",...}}
  │
  ▼
transport.stdio  →  orchestrator  →  tools.features  →  orchestrator  →  storage.markdown
                        (QUIC)           (QUIC)             (QUIC)
                                                                           │
                                                                     .projects/my-app/
                                                                       features/FEAT-001.md
```

End-to-end: agent sends `create_feature`, gets a result, feature doc exists on disk. That's the small win.

---

## Transport: QUIC From Day One

The stack says QUIC. We build QUIC. No gRPC-first, no migration path, no shortcuts.

| Component | Library | Notes |
|-----------|---------|-------|
| Orchestrator | `quic-go/quic-go` | QUIC listener, accepts plugin connections |
| Go plugins | `quic-go/quic-go` | Connect to orchestrator over QUIC |
| Rust plugins (Phase 3) | `quinn` | Same Protobuf over QUIC |
| Swift plugins (Phase 4) | `Network.framework` | Apple's native QUIC |
| Browser (Phase 2) | WebTransport API | Via `quic-go/webtransport-go` |

### Phase 1 TLS: Self-signed localhost

QUIC mandates TLS 1.3. For Phase 1 (single machine):
- Auto-generate a CA at `~/.orchestra/certs/ca.crt` + `ca.key` on first run
- Each plugin gets a signed cert: `{plugin-id}.crt` + `{plugin-id}.key`
- `tls.Config.ClientAuth = tls.RequireAndVerifyClientCert` (mTLS)
- All on localhost — the `libs/go/plugin/` package handles this transparently

This is ~100 lines in `libs/go/plugin/certs.go`. We build it once.

### QUIC stream protocol

Messages on QUIC streams use length-delimited Protobuf:
```
[4 bytes: big-endian uint32 length][N bytes: Protobuf message]
```

Each RPC call = one bidirectional QUIC stream (open → request → response → close). QUIC streams are cheap (few bytes of framing). No head-of-line blocking.

---

## Step 1: Protobuf Contract

**What**: Define the minimum PluginHost service. Lifecycle + tools + storage.

**Exit criteria**: `buf lint` passes, `buf generate` produces Go code, all types compile.

### Files

```
proto/
  buf.yaml
  buf.gen.yaml
  orchestra/plugin/v1/
    plugin.proto
gen/go/
  go.mod
  orchestra/plugin/v1/
    plugin.pb.go          (generated)
```

### plugin.proto (Phase 1 subset)

```protobuf
syntax = "proto3";
package orchestra.plugin.v1;

import "google/protobuf/struct.proto";
import "google/protobuf/timestamp.proto";

// ================================================================
// PluginHost — every plugin implements this
// Phase 1: lifecycle + tools + storage
// Phase 2 adds: HandleEvent, AIChat, AIEmbed, TransportSend/Receive
// ================================================================
// NOTE: This is NOT a gRPC service. Messages are sent over QUIC streams
// using length-delimited Protobuf. The service definition is used
// for code generation and documentation only.
// ================================================================

// Each RPC maps to a message pair sent on a QUIC bidirectional stream:
//   Client opens stream → writes Request → reads Response → stream closes

// ---- Request envelope (wraps all RPCs) ----

message PluginRequest {
  string request_id = 1;  // UUIDv7
  oneof request {
    PluginManifest register = 10;
    BootRequest boot = 11;
    ShutdownRequest shutdown = 12;
    HealthRequest health = 13;
    ToolRequest tool_call = 20;
    ListToolsRequest list_tools = 21;
    StorageReadRequest storage_read = 30;
    StorageWriteRequest storage_write = 31;
    StorageDeleteRequest storage_delete = 32;
    StorageListRequest storage_list = 33;
  }
}

message PluginResponse {
  string request_id = 1;
  oneof response {
    RegistrationResult register = 10;
    BootResult boot = 11;
    ShutdownResult shutdown = 12;
    HealthResult health = 13;
    ToolResponse tool_call = 20;
    ListToolsResponse list_tools = 21;
    StorageReadResponse storage_read = 30;
    StorageWriteResponse storage_write = 31;
    StorageDeleteResponse storage_delete = 32;
    StorageListResponse storage_list = 33;
  }
}

// ---- Lifecycle ----

message PluginManifest {
  string id = 1;                        // "tools.features", "storage.markdown"
  string version = 2;                   // semver
  string language = 3;                  // "go", "rust", "swift", etc.
  repeated string provides_tools = 4;
  repeated string provides_events = 5;
  repeated string provides_storage = 6; // ["markdown", "sqlite", "postgres"]
  repeated string provides_transport = 7;
  repeated string provides_ai = 8;
  repeated string needs_storage = 9;
  repeated string needs_events = 10;
  repeated string needs_ai = 11;
  repeated string needs_tools = 12;
  string description = 13;
  string author = 14;
  string binary = 15;
  repeated string args = 16;
  map<string, string> env = 17;
}

message RegistrationResult {
  bool accepted = 1;
  string reject_reason = 2;
}

message BootRequest {
  map<string, string> config = 1;
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
  enum Status {
    UNKNOWN = 0;
    HEALTHY = 1;
    DEGRADED = 2;
    UNHEALTHY = 3;
  }
  Status status = 1;
  string message = 2;
  map<string, string> details = 3;
}

// ---- Tools ----

message ToolRequest {
  string request_id = 1;               // UUIDv7
  string tool_name = 2;
  google.protobuf.Struct arguments = 3;
  string caller_plugin = 4;
  string trace_parent = 5;
}

message ToolResponse {
  string request_id = 1;
  bool success = 2;
  google.protobuf.Struct result = 3;
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

// ---- Storage ----

message StorageReadRequest {
  string path = 1;                      // "projects/my-app/features/FEAT-001.md"
  string storage_type = 2;             // "markdown" (default), "sqlite", "postgres"
}

message StorageReadResponse {
  bytes content = 1;                    // Markdown body
  google.protobuf.Struct metadata = 2; // Structured fields (status, priority, etc.)
  int64 version = 3;
}

message StorageWriteRequest {
  string path = 1;
  bytes content = 2;                    // Markdown body
  google.protobuf.Struct metadata = 3; // Structured fields
  int64 expected_version = 4;          // 0 = create, >0 = CAS update
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
  string prefix = 1;                   // "projects/my-app/features/"
  string pattern = 2;                  // "*.md"
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
```

### Key design: PluginRequest/PluginResponse envelope

Instead of gRPC services, we use a single request/response envelope over QUIC streams. One bidirectional QUIC stream per RPC call. The `oneof` field determines which operation is being called. This keeps the Protobuf definitions clean while using raw QUIC instead of gRPC.

### buf.gen.yaml (Go only for Phase 1)

```yaml
version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/orchestra-agents/gen/go
plugins:
  - remote: buf.build/protocolbuffers/go
    out: ../gen/go
    opt: [paths=source_relative]
```

Note: No gRPC plugin — we're using raw QUIC, not gRPC. Only Protobuf serialization is generated.

### Dependencies

- `buf` CLI (`brew install bufbuild/buf/buf`)
- gen/go/go.mod: `google.golang.org/protobuf`

### Test

```bash
cd proto && buf lint
buf generate
cd ../gen/go && go build ./...
```

---

## Step 2: Plugin SDK (`libs/go/`)

**What**: Shared library every Go plugin imports. QUIC transport, mTLS cert management, request/response framing, orchestrator client, manifest builder, lifecycle hooks, domain types, helpers.

**Exit criteria**: A test plugin connects to a test orchestrator over QUIC with mTLS, registers, and responds to ListTools + HandleToolCall.

### Files

```
libs/go/
  go.mod
  plugin/
    plugin.go           # Plugin struct: QUIC server + orch client + tools
    server.go           # QUIC listener, accept streams, dispatch requests
    client.go           # OrchestratorClient: QUIC connection + send/receive
    framing.go          # Length-delimited Protobuf read/write on QUIC streams
    certs.go            # Auto CA generation, cert signing, mTLS config
    manifest.go         # Fluent manifest builder
    lifecycle.go        # LifecycleHooks interface
    plugin_test.go
  types/
    feature.go          # FeatureData, FeatureStatus, ReviewEntry
    project.go          # ProjectStatus, ProjectConfig
    workflow.go         # Feature workflow states, transitions, gate evidence
  helpers/
    args.go             # GetString, GetFloat64, GetStringSlice from Struct
    results.go          # TextResult, JSONResult, ErrorResult
    paths.go            # ProjectsDir, FeaturesDir, FileExists
    strings.go          # Slugify, ContainsCI, Now
    validate.go         # ValidateArgs
  protocol/
    jsonrpc.go          # JSON-RPC 2.0 request/response types
    mcp.go              # MCP constants, InitializeResult, ServerCapabilities
```

### Plugin API

```go
p := plugin.New("tools.features", "1.0.0").
    Description("Feature-driven workflow — 35 tools").
    ProvidesTools("create_feature", "get_feature", "advance_feature" /* ... */).
    NeedsStorage("read", "write", "list", "delete").
    OnBoot(func(ctx context.Context, cfg map[string]string) error {
        return nil
    }).
    RegisterTool("create_feature", createFeatureDef, createFeatureHandler).
    Build()

p.Serve(ctx, listenAddr, orchestratorAddr)
// listenAddr: QUIC port for orchestrator to connect
// orchestratorAddr: QUIC address to call back for storage
```

### OrchestratorClient (QUIC-based)

```go
type OrchestratorClient struct {
    conn quic.Connection  // QUIC connection to orchestrator
}

func (c *OrchestratorClient) Send(ctx context.Context, req *pluginv1.PluginRequest) (*pluginv1.PluginResponse, error) {
    stream, _ := c.conn.OpenStreamSync(ctx)
    defer stream.Close()
    framing.WriteMessage(stream, req)
    return framing.ReadMessage[pluginv1.PluginResponse](stream)
}

// Convenience methods:
func (c *OrchestratorClient) StorageRead(ctx context.Context, path string) ([]byte, *structpb.Struct, int64, error)
func (c *OrchestratorClient) StorageWrite(ctx context.Context, path string, content []byte, meta *structpb.Struct, version int64) (int64, error)
func (c *OrchestratorClient) StorageList(ctx context.Context, prefix, pattern string) ([]*pluginv1.StorageEntry, error)
func (c *OrchestratorClient) StorageDelete(ctx context.Context, path string) error
```

### framing.go (QUIC message framing)

```go
func WriteMessage(stream quic.Stream, msg proto.Message) error {
    data, _ := proto.Marshal(msg)
    length := make([]byte, 4)
    binary.BigEndian.PutUint32(length, uint32(len(data)))
    stream.Write(length)
    stream.Write(data)
    return nil
}

func ReadMessage[T proto.Message](stream quic.Stream) (*T, error) {
    length := make([]byte, 4)
    io.ReadFull(stream, length)
    size := binary.BigEndian.Uint32(length)
    data := make([]byte, size)
    io.ReadFull(stream, data)
    var msg T
    proto.Unmarshal(data, &msg)
    return &msg, nil
}
```

### certs.go (auto mTLS)

```go
func EnsureCerts(pluginID string) (*tls.Config, error) {
    certsDir := filepath.Join(os.Getenv("HOME"), ".orchestra", "certs")
    caPath := filepath.Join(certsDir, "ca.crt")

    if !fileExists(caPath) {
        generateCA(certsDir)  // ~40 lines: RSA 4096, self-signed CA
    }

    certPath := filepath.Join(certsDir, "services", pluginID+".crt")
    if !fileExists(certPath) {
        signCert(certsDir, pluginID)  // ~30 lines: sign with CA
    }

    // Return tls.Config with mTLS
    return &tls.Config{
        Certificates: []tls.Certificate{loadCert(certPath, keyPath)},
        ClientCAs:    loadCA(caPath),
        RootCAs:      loadCA(caPath),
        ClientAuth:   tls.RequireAndVerifyClientCert,
        NextProtos:   []string{"orch-plugin/1"},
    }, nil
}
```

### Dependencies

```
github.com/quic-go/quic-go
google.golang.org/protobuf
github.com/google/uuid
github.com/orchestra-agents/gen/go  (via go.work)
```

### Test

```bash
cd libs/go && go test ./plugin/... -v
```

Start QUIC server in goroutine, connect with QUIC client (mTLS), verify Register → Boot → ListTools → HandleToolCall over QUIC streams.

---

## Step 3: Orchestrator

**What**: The central hub. ~500-600 lines of Go. QUIC listener. Reads `plugins.yaml`, starts plugin binaries, connects over QUIC+mTLS, builds routing tables, forwards messages.

**Exit criteria**: Starts, loads config, launches plugin binaries, QUIC connections established, routes HandleToolCall and StorageWrite correctly.

### Files

```
services/orchestrator/
  go.mod
  cmd/main.go
  internal/
    config.go           # Parse plugins.yaml
    loader.go           # Start plugin binaries, QUIC connect, register, boot
    router.go           # Routing tables: tool→plugin, storage→plugin
    server.go           # QUIC listener (for plugin callbacks)
    orchestrator.go     # Main struct
    orchestrator_test.go
```

### Architecture

```
┌─────────────────────────────────────────────┐
│  ORCHESTRATOR                               │
│                                             │
│  QUIC Listener (plugins call back here):    │
│    StorageRead/Write/Delete/List            │
│      → routes to storage.markdown plugin    │
│    HandleToolCall                           │
│      → routes to tool-owning plugin         │
│    ListTools                                │
│      → aggregates from all tool plugins     │
│                                             │
│  Plugin Connections (QUIC clients):         │
│    tools.features  → quic.Connection        │
│    storage.markdown → quic.Connection       │
│    transport.stdio → quic.Connection        │
│                                             │
│  Routing Tables:                            │
│    toolName     → pluginConn                │
│    storageType  → pluginConn                │
└─────────────────────────────────────────────┘
```

### Plugin startup protocol

1. Orchestrator starts plugin binary as child process:
   `--orchestrator-addr=localhost:{port} --listen-addr=:0`
2. Plugin starts QUIC listener (mTLS), prints `READY <address>` to stderr
3. Orchestrator reads stderr, opens QUIC connection (mTLS)
4. Opens stream → sends `PluginRequest{register: manifest}` → reads `PluginResponse`
5. Opens stream → sends `PluginRequest{boot: config}` → reads `PluginResponse`
6. Plugin is live

### Router

```go
type Router struct {
    tools      map[string]*PluginConn   // "create_feature" → tools.features
    storage    map[string]*PluginConn   // "markdown"        → storage.markdown
    allTools   []*pluginv1.ToolDefinition
}

func (r *Router) RouteToolCall(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error)
func (r *Router) RouteStorageRead(ctx context.Context, req *pluginv1.StorageReadRequest) (*pluginv1.StorageReadResponse, error)
func (r *Router) RouteStorageWrite(ctx context.Context, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error)
func (r *Router) ListAllTools() []*pluginv1.ToolDefinition
```

### plugins.yaml

```yaml
orchestrator:
  quic_port: 50100
  log_level: info
  certs_dir: ~/.orchestra/certs

plugins:
  storage.markdown:
    enabled: true
    binary: orch-storage-markdown
    auto_start: true
    config:
      workspace: "."

  tools.features:
    enabled: true
    binary: orch-tools-features
    auto_start: true

  transport.stdio:
    enabled: true
    binary: orch-transport-stdio
    auto_start: true
```

### Dependencies

```
github.com/quic-go/quic-go
google.golang.org/protobuf
gopkg.in/yaml.v3
github.com/orchestra-agents/gen/go
github.com/orchestra-agents/libs/go
```

---

## Step 4: `storage.markdown` Plugin

**What**: First real plugin. Reads/writes Protobuf metadata + Markdown body files to `.projects/`.

**Exit criteria**: Connects to orchestrator over QUIC, handles StorageRead/Write/Delete/List. Files on disk are human-readable Markdown with JSON metadata block.

### Files

```
plugins/storage-markdown/
  go.mod
  cmd/main.go
  internal/
    storage.go          # StoragePlugin: handles storage requests
    reader.go           # Parse .md files: JSON metadata block + Markdown body
    writer.go           # Write JSON metadata block + Markdown body
    paths.go            # Workspace-relative path resolution
    storage_test.go
```

### On-disk file format

```markdown
<!-- META
{
  "id": "FEAT-001",
  "type": "feature",
  "status": "in-progress",
  "priority": "high",
  "assignee": "claude-opus",
  "labels": ["backend", "auth"],
  "version": 2,
  "created_at": "2026-02-26T10:00:00Z"
}
META -->

# FEAT-001: User authentication with JWT

## Context
The API server needs authentication...

## Definition of Done
- [ ] POST /login returns access + refresh token
...
```

**Why JSON in HTML comments**:
- Pure `.md` file — renders correctly in any Markdown viewer
- JSON block is hidden when rendered (HTML comment)
- `google.protobuf.Struct` serializes naturally to/from JSON
- Human-readable, git-trackable, single file per entity

### Dependencies

```
github.com/quic-go/quic-go
google.golang.org/protobuf
github.com/orchestra-agents/gen/go
github.com/orchestra-agents/libs/go
```

---

## Step 5: `tools.features` Plugin

**What**: The feature-driven workflow engine. ~35 tools. Does NOT touch filesystem — calls orchestrator storage over QUIC.

**Exit criteria**: All tools registered, respond to HandleToolCall, correctly use storage via orchestrator. Feature lifecycle works end-to-end.

### Files

```
plugins/tools-features/
  go.mod
  cmd/main.go
  internal/
    features.go                 # FeaturesPlugin: tool registration
    storage/
      client.go                 # Feature-specific storage wrapper
      scanner.go                # ScanAllFeatures via storage client
    tools/
      project.go                # 4 tools: create/list/delete/status
      feature.go                # 5 tools: create/get/update/list/delete + search
      workflow.go               # 5 tools: get_next/set_current/advance/reject/status
      review.go                 # 3 tools: request_review/submit_review/pending
      dependency.go             # 3 tools: add/remove/graph
      wip.go                    # 3 tools: set/get/check limits
      split.go                  # 1 tool: split_feature (context-aware sizing)
      reporting.go              # 3 tools: progress/blocked/review_queue
      metadata.go               # 5 tools: labels, assign, estimate, notes
    features_test.go
```

### Feature-specific storage wrapper

```go
type FeatureStorage struct {
    orch *plugin.OrchestratorClient
}

func (s *FeatureStorage) ReadFeature(ctx context.Context, project, featureID string) (*types.FeatureData, string, error) {
    path := fmt.Sprintf("projects/%s/features/%s.md", project, featureID)
    content, meta, version, err := s.orch.StorageRead(ctx, path)
    // Parse meta Struct into FeatureData
    // Return body as Markdown content
}

func (s *FeatureStorage) WriteFeature(ctx context.Context, project string, feature *types.FeatureData, body string) error {
    path := fmt.Sprintf("projects/%s/features/%s.md", project, feature.ID)
    meta := featureToStruct(feature)
    _, err := s.orch.StorageWrite(ctx, path, []byte(body), meta, feature.Version)
    return err
}

func (s *FeatureStorage) ListFeatures(ctx context.Context, project string) ([]types.FeatureData, error) {
    entries, _ := s.orch.StorageList(ctx, fmt.Sprintf("projects/%s/features/", project), "*.md")
    // Read each, collect
}
```

### Tool count: ~35

| File | Tools | Count |
|------|-------|-------|
| project.go | create_project, list_projects, delete_project, get_project_status | 4 |
| feature.go | create_feature, get_feature, update_feature, list_features, delete_feature, search | 6 |
| workflow.go | get_next_feature, set_current_feature, advance_feature, reject_feature, get_workflow_status | 5 |
| review.go | request_review, submit_review, get_pending_reviews | 3 |
| dependency.go | add_dependency, remove_dependency, get_dependency_graph | 3 |
| wip.go | set_wip_limits, get_wip_limits, check_wip_limit | 3 |
| split.go | split_feature | 1 |
| reporting.go | get_progress, get_blocked_features, get_review_queue | 3 |
| metadata.go | add_labels, remove_labels, assign_feature, unassign_feature, set_estimate, save_note, list_notes | 7 |
| **Total** | | **35** |

### Dependencies

```
github.com/quic-go/quic-go
google.golang.org/protobuf
github.com/orchestra-agents/gen/go
github.com/orchestra-agents/libs/go
```

---

## Step 6: `transport.stdio` Plugin

**What**: MCP JSON-RPC over stdin/stdout. How AI agents talk to Orchestra.

**Exit criteria**: `echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | orch-transport-stdio` returns all 35 tools.

### Files

```
plugins/transport-stdio/
  go.mod
  cmd/main.go
  internal/
    transport.go        # Read stdin, write stdout loop
    handler.go          # MCP protocol: initialize, tools/list, tools/call, ping
    translator.go       # JSON-RPC ↔ Protobuf translation
    transport_test.go
```

### Architecture

```
stdin → scanner → JSON-RPC parse → MCP handler → orchestrator (QUIC) → JSON-RPC response → stdout
```

This plugin has NO tools of its own. It's a bridge:
- `tools/list` → calls orchestrator `ListTools` (aggregated from all plugins) over QUIC
- `tools/call` → calls orchestrator `HandleToolCall` over QUIC → routed to owning plugin

### MCP protocol support (Phase 1)

| Method | Status |
|--------|--------|
| `initialize` | Supported |
| `notifications/initialized` | Supported |
| `tools/list` | Supported |
| `tools/call` | Supported |
| `ping` | Supported |
| `resources/list` | Phase 2 |
| `prompts/list` | Phase 2 |

### Dependencies

```
github.com/quic-go/quic-go
google.golang.org/protobuf
github.com/orchestra-agents/gen/go
github.com/orchestra-agents/libs/go
```

---

## Step 7: Wire It All Together

**What**: go.work, Makefile, default config, end-to-end test.

**Exit criteria**: `make build && make test-e2e` — agent sends create_feature, gets result, feature doc on disk.

### go.work

```go
go 1.23

use (
    ./gen/go
    ./libs/go
    ./services/orchestrator
    ./plugins/storage-markdown
    ./plugins/tools-features
    ./plugins/transport-stdio
)
```

### Makefile

```makefile
BIN_DIR := $(CURDIR)/bin

.PHONY: proto build test test-e2e clean

proto:
	cd proto && buf lint && buf generate

build: build-orchestrator build-storage-markdown build-tools-features build-transport-stdio

build-orchestrator:
	cd services/orchestrator && go build -o $(BIN_DIR)/orch-orchestrator ./cmd/

build-storage-markdown:
	cd plugins/storage-markdown && go build -o $(BIN_DIR)/orch-storage-markdown ./cmd/

build-tools-features:
	cd plugins/tools-features && go build -o $(BIN_DIR)/orch-tools-features ./cmd/

build-transport-stdio:
	cd plugins/transport-stdio && go build -o $(BIN_DIR)/orch-transport-stdio ./cmd/

test:
	cd libs/go && go test ./... -v
	cd services/orchestrator && go test ./... -v
	cd plugins/storage-markdown && go test ./... -v
	cd plugins/tools-features && go test ./... -v
	cd plugins/transport-stdio && go test ./... -v

test-e2e: build
	./scripts/test-e2e.sh

clean:
	rm -rf $(BIN_DIR)
```

### End-to-end test

```bash
#!/bin/bash
set -e

# 1. Start orchestrator (auto-generates certs, starts plugins)
bin/orch-orchestrator --config plugins.yaml &
ORCH_PID=$!
sleep 3

# 2. tools/list — verify create_feature exists
RESULT=$(echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | \
  timeout 5 bin/orch-transport-stdio --orchestrator-addr=localhost:50100)
echo "$RESULT" | grep -q "create_feature" || (echo "FAIL: create_feature not found"; exit 1)

# 3. create_project
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"create_project","arguments":{"name":"Test Project"}}}' | \
  timeout 5 bin/orch-transport-stdio --orchestrator-addr=localhost:50100

# 4. create_feature
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"create_feature","arguments":{"project":"test-project","title":"User authentication","priority":"high"}}}' | \
  timeout 5 bin/orch-transport-stdio --orchestrator-addr=localhost:50100

# 5. Verify feature doc exists
test -f .projects/test-project/features/FEAT-*.md || (echo "FAIL: feature file not created"; exit 1)

# Cleanup
kill $ORCH_PID
rm -rf .projects/test-project
echo "ALL TESTS PASSED"
```

---

## Build Summary

| Step | What | New Files | Lines (est.) | Exit Criteria |
|------|------|-----------|-------------|---------------|
| 1 | Proto contract | 3 | ~250 proto | `buf lint` passes, Go compiles |
| 2 | Plugin SDK (QUIC + mTLS) | 14 | ~1,000 | Test plugin registers over QUIC + mTLS |
| 3 | Orchestrator | 6 | ~500 | Starts, QUIC connects plugins, routes messages |
| 4 | storage.markdown | 6 | ~400 | Read/write/list/delete .md files over QUIC |
| 5 | tools.features | 12 | ~1,200 | 35 tools respond correctly |
| 6 | transport.stdio | 5 | ~300 | JSON-RPC stdin → MCP response stdout |
| 7 | Integration | 4 | ~100 | `make build && make test-e2e` passes |
| **Total** | | **~50 files** | **~3,750** | Agent can manage features via MCP over QUIC |

---

## What This Foundation Enables

After Phase 1, adding any new capability is: **build a binary, add to plugins.yaml**.

| Future Plugin | What Changes |
|--------------|-------------|
| `transport.sse` | New binary. Zero orchestrator changes. |
| `transport.webtransport` | New binary. Zero orchestrator changes. |
| `tools.memory` | New binary. Zero orchestrator changes. |
| `tools.notes` | New binary. Zero orchestrator changes. |
| `engine.search` (Rust) | New Rust binary via `quinn`. Zero orchestrator changes. |
| `storage.sqlite` | New binary. Zero orchestrator changes. |
| `storage.postgres` | New binary. Zero orchestrator changes. |
| `integration.github` | New binary with 17 tools. Zero orchestrator changes. |
| `ai.claude` | New binary. Add AI routing to orchestrator (~50 lines). |
| `ui.macos` (Swift) | Swift binary via `Network.framework`. Zero orchestrator changes. |
| mDNS discovery | Add to orchestrator `loader.go`. Zero plugin code changes. |

---

## Risk Mitigations

| Risk | Mitigation |
|------|-----------|
| QUIC + TLS complexity | `libs/go/plugin/certs.go` auto-generates CA + signs certs. ~100 lines, done once. |
| Plugin startup coordination | Each plugin prints `READY <addr>` to stderr. Orchestrator watches with 10s timeout. |
| Storage path conventions | Constants in `libs/go/helpers/paths.go`. Used by tools and storage plugins. |
| Metadata type fidelity | `google.protobuf.Struct` preserves arrays, numbers, booleans. |
| Large feature docs | `context_budget` field ensures each feature fits the agent's context window. |

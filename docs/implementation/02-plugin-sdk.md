# Step 2: Plugin SDK (`libs/go/`)

## Status: Complete

## What Was Built

The shared Go library that every Go plugin imports. Provides QUIC transport, mTLS certificate management, Protobuf framing, a fluent plugin builder API, domain types, and helper utilities.

## Module

`github.com/orchestrated-mcp/framework/libs/go`

Dependencies: `quic-go/quic-go`, `google.golang.org/protobuf`, `github.com/google/uuid`

## Package Overview

### `plugin/` — Core QUIC Transport + mTLS

| File | Purpose | Key API |
|------|---------|---------|
| `framing.go` | Length-delimited Protobuf framing | `WriteMessage(w, msg)`, `ReadMessage(r, msg)` |
| `certs.go` | Auto mTLS certificate management | `EnsureCA()`, `GenerateCert()`, `ServerTLSConfig()`, `ClientTLSConfig()` |
| `server.go` | QUIC server (accept streams, dispatch) | `NewServer()`, `ListenAndServe()`, `RegisterTool()` |
| `client.go` | QUIC client (connect to orchestrator) | `NewOrchestratorClient()`, `Send(req)` |
| `manifest.go` | Fluent manifest builder | `NewManifest(id).ProvidesTools().Build()` |
| `lifecycle.go` | Lifecycle hooks interface | `OnBoot(config)`, `OnShutdown()` |
| `plugin.go` | Main plugin struct, ties everything together | `New(id).RegisterTool().Build().Run(ctx)` |
| `plugin_test.go` | Integration tests (11 tests, all pass) | Framing, certs, QUIC flow |

### `types/` — Domain Types

| File | Purpose | Key Types |
|------|---------|-----------|
| `feature.go` | Feature entity and status enum | `FeatureData`, `FeatureStatus`, `ReviewEntry` |
| `project.go` | Project entity | `ProjectData` |
| `workflow.go` | State machine transitions | `CanTransition()`, `NextStatuses()`, `ValidTransitions` |

### `helpers/` — Utilities

| File | Purpose | Key Functions |
|------|---------|---------------|
| `args.go` | Extract typed values from `structpb.Struct` | `GetString()`, `GetInt()`, `GetBool()`, `GetStringSlice()` |
| `results.go` | Build `ToolResponse` helpers | `TextResult()`, `JSONResult()`, `ErrorResult()` |
| `paths.go` | Path constants and builders | `FeaturePath()`, `ProjectPath()` |
| `strings.go` | String utilities | `Slugify()`, `NowISO()`, `NewUUID()`, `NewFeatureID()` |
| `validate.go` | Validation helpers | `ValidateRequired()`, `ValidateOneOf()` |

### `protocol/` — JSON-RPC + MCP Types

| File | Purpose | Key Types |
|------|---------|-----------|
| `jsonrpc.go` | JSON-RPC 2.0 types | `JSONRPCRequest`, `JSONRPCResponse`, `JSONRPCError` |
| `mcp.go` | MCP protocol types | `MCPInitializeResult`, `MCPToolDefinition`, `MCPToolResult` |

## Plugin Lifecycle

```
1. New("tools.features")
     .Version("1.0.0")
     .RegisterTool("create_feature", desc, schema, handler)
     .BuildWithTools()

2. plugin.ParseFlags()        // --orchestrator-addr, --listen-addr, --certs-dir

3. plugin.Run(ctx)
   a. EnsureCA + GenerateCert     → auto mTLS certs
   b. Start QUIC server           → accept connections
   c. Print "READY <addr>"        → stderr signal to orchestrator
   d. Connect to orchestrator     → QUIC client
   e. Send Register(manifest)     → orchestrator accepts
   f. Serve tool calls            → until context cancelled
   g. OnShutdown + close          → cleanup
```

## Certificate Management

```
~/.orchestra/certs/
├── ca.crt              # Auto-generated ed25519 CA (1-year validity)
├── ca.key              # CA private key (0600 permissions)
├── {plugin-name}.crt   # Plugin cert signed by CA
└── {plugin-name}.key   # Plugin private key
```

- `EnsureCA()` is idempotent — generates on first call, loads from disk after
- `ServerTLSConfig()` requires client certs (mTLS)
- `ClientTLSConfig()` trusts the same CA
- TLS 1.3 minimum, ALPN protocol `"orchestra-plugin"`

## Feature Workflow State Machine

```
backlog → todo → in-progress → ready-for-testing → in-testing
                     ↑                                  │
                     │                            pass/fail
                     │                                  │
              needs-edits ← in-review ← documented ← ready-for-docs ← (pass)
                     │          │                              ↑
                     └──────────┘ (approved → done)     (fail → in-progress)
```

The cycle: `needs-edits → in-progress` is the heartbeat of iterative delivery.

## Tests

11 tests covering:
- Framing roundtrip (write + read)
- Oversized message rejection (16MB limit)
- CA generation and reload from disk
- Server/Client TLS config creation
- Manifest builder fluent API
- Full QUIC integration: Register, ListTools, ToolCall, ToolCallNotFound, Health

```bash
cd libs/go && go test ./plugin/... -v
```

## Verification

```bash
cd libs/go && go build ./...                  # All packages compile
cd libs/go && go test ./plugin/... -v         # 11/11 tests pass
```

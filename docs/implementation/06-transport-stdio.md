# Step 6: transport.stdio Plugin

## Status: Complete

## What Was Built

The MCP JSON-RPC bridge — reads JSON-RPC requests from stdin, translates to Protobuf, sends over QUIC to the orchestrator, translates back to JSON-RPC, writes to stdout. This is the entry point for AI agents (Claude, GPT, etc.) to interact with the system.

## Module

`github.com/orchestrated-mcp/framework/plugins/transport-stdio`

## Files

| File | Purpose |
|------|---------|
| `cmd/main.go` | Entry point — parse flags, connect QUIC, start stdio loop |
| `internal/transport.go` | StdioTransport: stdin → parse → dispatch → stdout loop |
| `internal/handler.go` | MCP protocol handlers: initialize, tools/list, tools/call, ping |
| `internal/translator.go` | JSON-RPC ↔ Protobuf ToolRequest/ToolResponse conversion |
| `internal/transport_test.go` | 20 tests |

## MCP Methods (Phase 1)

| Method | Description |
|--------|-------------|
| `initialize` | Returns server capabilities (tools support) |
| `notifications/initialized` | Client notification, acknowledged silently |
| `tools/list` | Queries orchestrator for all registered tools |
| `tools/call` | Forwards tool invocation to orchestrator |
| `ping` | Returns empty result |

## Architecture

```
stdin (JSON-RPC) → StdioTransport → Handler → Sender (QUIC) → Orchestrator
                                                                    ↓
stdout (JSON-RPC) ← StdioTransport ← Handler ← Response ←─────────┘
```

### Sender Interface

```go
type Sender interface {
    Send(ctx context.Context, req *pluginv1.PluginRequest) (*pluginv1.PluginResponse, error)
}
```

Abstraction over `OrchestratorClient` for testability. Tests use `mockSender` that returns canned responses without QUIC.

### Transport Details

- 10MB scanner buffer for large JSON-RPC messages
- Empty lines skipped (robustness)
- Notifications (no `id` field) handled silently
- Invalid JSON returns JSON-RPC parse error (-32700)
- Unknown methods return method-not-found error (-32601)

## Tests (20 pass)

| Test | Coverage |
|------|----------|
| TestInitialize | MCP initialize handshake |
| TestPing | Ping/pong |
| TestToolsList | List tools from mock sender |
| TestToolsCall | Successful tool invocation |
| TestToolsCallError | Tool returns error response |
| TestToolsCallNetworkError | QUIC send failure |
| TestToolsCallMissingName | Missing tool name validation |
| TestToolsListNetworkError | List tools QUIC failure |
| TestNotification | Notification ignored (no response) |
| TestMethodNotFound | Unknown method error |
| TestParseError | Invalid JSON handling |
| TestMultipleRequests | Sequential request processing |
| TestEmptyLines | Empty line skipping |
| TestStructToMap | Protobuf Struct → map conversion |
| TestStructToMapNil | Nil struct handling |
| TestMapToStruct | Map → Protobuf Struct conversion |
| TestToolDefinitionToMCP | Proto ToolDefinition → MCP format |
| TestToolResponseToMCPSuccess | Successful response translation |
| TestToolResponseToMCPError | Error response translation |
| TestToolResponseToMCPFallback | Missing result fallback |

```bash
cd plugins/transport-stdio && go test ./... -v
```

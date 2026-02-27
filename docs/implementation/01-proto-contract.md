# Step 1: Proto Contract

## Status: Complete

## What Was Built

The Protobuf contract that defines the entire plugin wire protocol. This is the single source of truth for all message types exchanged between the orchestrator and plugins over QUIC.

## Files

| File | Purpose |
|------|---------|
| `proto/buf.yaml` | Buf module config (`buf.build/orchestrated-mcp/framework`) |
| `proto/buf.gen.yaml` | Code generation config (Go output to `gen/go/`) |
| `proto/orchestra/plugin/v1/plugin.proto` | The Protobuf schema (~200 lines) |
| `gen/go/go.mod` | Go module for generated types |
| `gen/go/orchestra/plugin/v1/plugin.pb.go` | Generated Go types (auto-generated, gitignored) |

## Proto Schema Overview

### Message Envelope

All communication uses a single request/response pair with `oneof` for polymorphism:

```protobuf
message PluginRequest {
  string request_id = 1;  // UUIDv7
  oneof request {
    PluginManifest register = 10;        // Plugin registration
    BootRequest boot = 11;               // Orchestrator boots plugin
    ShutdownRequest shutdown = 12;       // Graceful shutdown
    HealthRequest health = 13;           // Health check
    ToolRequest tool_call = 20;          // Invoke a tool
    ListToolsRequest list_tools = 21;    // List available tools
    StorageReadRequest storage_read = 30;
    StorageWriteRequest storage_write = 31;
    StorageDeleteRequest storage_delete = 32;
    StorageListRequest storage_list = 33;
  }
}
```

### Message Categories

1. **Lifecycle** (field numbers 10-13): Register, Boot, Shutdown, Health
2. **Tools** (field numbers 20-21): ToolCall, ListTools
3. **Storage** (field numbers 30-33): Read, Write, Delete, List

### Wire Protocol

NOT gRPC. Raw QUIC streams with length-delimited framing:
```
[4 bytes: big-endian uint32 length][N bytes: Protobuf message]
```

Each RPC = one bidirectional QUIC stream.

## Verification

```bash
cd proto && buf lint          # Passes STANDARD lint rules
cd proto && buf generate      # Generates plugin.pb.go
cd gen/go && go build ./...   # Compiles cleanly
```

## Module Path

All Go modules use `github.com/orchestrated-mcp/framework` as the base module path.

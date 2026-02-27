# Step 3: Orchestrator (`services/orchestrator/`)

## Status: Complete

## What Was Built

The central hub that loads plugin binaries, manages their lifecycle, and routes messages between them. Star topology — all inter-plugin communication goes through the orchestrator.

## Module

`github.com/orchestrated-mcp/framework/services/orchestrator`

## Files

| File | Purpose |
|------|---------|
| `cmd/main.go` | Entry point — parse flags, load config, start orchestrator, handle signals |
| `internal/config.go` | Parse `plugins.yaml` — plugin IDs, binaries, args, env, enabled flag |
| `internal/orchestrator.go` | Main struct — ties server, router, loader together |
| `internal/loader.go` | Plugin binary launcher — start process, read READY, QUIC connect, Register, Boot |
| `internal/router.go` | Message routing — toolName→plugin, storageType→plugin |
| `internal/server.go` | QUIC server for plugin callbacks — proxy storage/tool requests |
| `internal/orchestrator_test.go` | 5 test functions, 10 sub-tests |

## Architecture

```
                    ┌──────────────────────────┐
                    │      Orchestrator         │
                    │  ┌────────┐  ┌────────┐  │
  Plugin A ←─QUIC──│  │ Server │  │ Router │  │──QUIC─→ Plugin B
  (tools)          │  └────────┘  └────────┘  │         (storage)
                    │       ┌────────┐         │
                    │       │ Loader │         │
                    │       └────────┘         │
                    └──────────────────────────┘
```

## Plugin Startup Protocol

```
1. Orchestrator starts: ./plugin --orchestrator-addr=host:port --listen-addr=localhost:0 --certs-dir=dir
2. Plugin starts QUIC listener, prints "READY <addr>" to stderr
3. Orchestrator reads stderr, connects via QUIC (mTLS)
4. Sends Register(manifest) → verifies accepted
5. Sends Boot(config) → verifies ready
6. Queries ListTools → populates routing table
7. Plugin is live
```

## Routing

- **Tool routing**: `toolName → RunningPlugin` — when a request arrives for "create_feature", route to the plugin that provides it
- **Storage routing**: `storageType → RunningPlugin` — when a plugin requests StorageRead with type "markdown", route to storage.markdown plugin
- **Proxy pattern**: Plugin A → orchestrator (StorageRead) → Plugin B (storage) → response back through orchestrator → Plugin A

## Config Format (`plugins.yaml`)

```yaml
listen_addr: "localhost:50100"
certs_dir: "~/.orchestra/certs"
plugins:
  - id: storage.markdown
    binary: ./bin/storage-markdown
    enabled: true
    config:
      workspace: .
  - id: tools.features
    binary: ./bin/tools-features
    enabled: true
  - id: transport.stdio
    binary: ./bin/transport-stdio
    enabled: true
```

## Tests (10/10 pass)

| Test | Coverage |
|------|----------|
| TestLoadConfig | YAML parsing with all fields |
| TestLoadConfigDefaults | Default listen_addr and certs_dir |
| TestRouter (4 sub-tests) | Route echo to A, greet to B, not-found error, unregister |
| TestRouterListAllTools | Aggregate tools from multiple plugins |
| TestOrchestratorServerDispatch (3 sub-tests) | Health, ToolCall, ListTools through proxy server |

```bash
cd services/orchestrator && go test ./... -v
```

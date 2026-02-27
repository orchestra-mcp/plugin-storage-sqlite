# Orchestra Framework

An AI-agentic IDE framework built on a plugin host architecture. Every capability is a plugin — storage, tools, transport, AI — communicating over QUIC with mTLS and Protobuf messages.

## Architecture

```
Agent (Claude, GPT, etc.)
  │ JSON-RPC (stdin/stdout)
  ▼
transport.stdio
  │ QUIC + Protobuf
  ▼
orchestrator ──────────────────┐
  │ QUIC + Protobuf            │ QUIC + Protobuf
  ▼                            ▼
tools.features (34 tools)    storage.markdown (disk)
```

**Star topology** — the orchestrator is the only router. Plugins never talk directly to each other. Any language that speaks QUIC + Protobuf can be a plugin.

## Quick Start

```bash
# Install dependencies
make install

# Build all binaries
make build

# Run all tests (66 unit + 1 E2E)
make test
make test-e2e

# Start the system
./bin/orchestrator --config plugins.yaml
# In another terminal:
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"0.1.0"}}}' | ./bin/transport-stdio --orchestrator-addr=localhost:50100
```

## Binaries

| Binary | Description |
|--------|-------------|
| `bin/orchestrator` | Central hub — starts plugins, routes QUIC messages |
| `bin/storage-markdown` | File storage with `<!-- META {json} META -->` + Markdown format |
| `bin/tools-features` | 34 feature-driven workflow tools |
| `bin/transport-stdio` | MCP JSON-RPC bridge (stdin/stdout to QUIC) |

## Project Structure

```
orchestra-agents/
├── proto/                          # Protobuf contract (plugin.proto)
├── gen/go/                         # Generated Go code from proto
├── libs/go/                        # Plugin SDK
│   ├── plugin/                     #   QUIC server/client, mTLS, framing
│   ├── types/                      #   Feature, Project, Workflow types
│   ├── helpers/                    #   Args, results, paths, validation
│   └── protocol/                   #   JSON-RPC + MCP types
├── services/orchestrator/          # Orchestrator service
│   ├── cmd/main.go                 #   Entry point
│   └── internal/                   #   Config, loader, router, server
├── plugins/
│   ├── storage-markdown/           # Markdown file storage plugin
│   ├── tools-features/             # 34 workflow tools plugin
│   └── transport-stdio/            # MCP stdio bridge plugin
├── docs/
│   ├── artifacts/                  # 17 design artifacts
│   └── implementation/             # Step-by-step implementation docs (01-07)
├── scripts/test-e2e.sh             # End-to-end integration test
├── plugins.yaml                    # Default plugin configuration
├── go.work                         # Go workspace (6 modules)
└── Makefile                        # Build, test, clean, proto targets
```

## Makefile Targets

```bash
make proto              # Lint + generate proto code
make build              # Build all 4 binaries to bin/
make test               # Run all unit tests (66 tests)
make test-e2e           # Build + run end-to-end integration test
make clean              # Remove build artifacts and certs
```

## Plugin SDK

Build a plugin in ~50 lines of Go:

```go
package main

import (
    "context"
    "github.com/orchestrated-mcp/framework/libs/go/plugin"
    "google.golang.org/protobuf/types/known/structpb"
)

func main() {
    p := plugin.New("my.plugin").
        Version("0.1.0").
        Description("My custom plugin").
        ProvidesTools("hello").
        BuildWithTools()

    p.Server().RegisterTool("hello", "Say hello", nil,
        func(ctx context.Context, args *structpb.Struct) (string, error) {
            return `{"message": "Hello from my plugin!"}`, nil
        })

    p.ParseFlags()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    p.Run(ctx)
}
```

## Feature Workflow (34 Tools)

The tools.features plugin implements an 11-state feature lifecycle:

```
backlog → todo → in-progress → ready-for-review → in-review → done
                      ↑                                  │
                      └──── needs-edits ◄────────────────┘
```

### Tool Categories

| Category | Tools | Description |
|----------|-------|-------------|
| Project | 4 | create, list, delete, get_status |
| Feature | 6 | create, get, update, list, delete, search |
| Workflow | 5 | advance, reject, get_next, set_current, get_status |
| Review | 3 | request, submit, get_pending |
| Dependencies | 3 | add, remove, get_graph |
| WIP Limits | 3 | set, get, check |
| Reporting | 3 | progress, blocked, review_queue |
| Metadata | 7 | labels, assignment, estimate, notes |

## Protocol

All communication uses length-delimited Protobuf over QUIC with mTLS:

```
[4 bytes big-endian uint32 length][N bytes Protobuf PluginRequest/Response]
```

- **mTLS**: Auto-generated ed25519 CA + per-plugin certificates at `~/.orchestra/certs/`
- **Plugin startup**: Binary prints `READY <addr>` to stderr, orchestrator connects and sends Register + Boot
- **Storage format**: `<!-- META {json} META -->` + blank line + Markdown body

## Tech Stack

- **Go**: quic-go, protobuf, uuid (orchestrator + plugins)
- **Proto**: buf lint + buf generate
- **Transport**: QUIC with mTLS (no gRPC)
- **Storage**: Markdown files with JSON metadata headers

## Module Path

```
github.com/orchestrated-mcp/framework
```

## License

MIT

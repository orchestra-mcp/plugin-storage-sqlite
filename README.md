# Orchestra Framework

An AI-agentic IDE framework built on a plugin host architecture. Every capability is a plugin — storage, tools, transport, AI — communicating over QUIC with mTLS and Protobuf messages.

## Install

```bash
# macOS / Linux — one-line install
curl -fsSL https://raw.githubusercontent.com/orchestra-mcp/framework/master/scripts/install.sh | sh

# Homebrew
brew install orchestra-mcp/tap/orchestra

# npm
npx @orchestra-mcp/cli init

# From source
git clone https://github.com/orchestra-mcp/framework.git
cd framework && make install
```

## Quick Start

```bash
# 1. Initialize Orchestra in your project (auto-detects your IDE)
cd your-project
orchestra init

# 2. That's it — your AI IDE now has 34 project management tools
```

`orchestra init` detects your IDE and writes the correct MCP config. Supported IDEs:

| IDE | Config File |
|-----|-------------|
| Claude Code | `.mcp.json` |
| Cursor | `.cursor/mcp.json` |
| VS Code / Copilot | `.vscode/mcp.json` |
| Cline | `.vscode/mcp.json` |
| Windsurf | `~/.codeium/windsurf/mcp_config.json` |
| Codex (OpenAI) | `.codex/config.toml` |
| Gemini Code Assist | `.gemini/settings.json` |
| Zed | `.zed/settings.json` |
| Continue.dev | `.continue/mcpServers/orchestra.yaml` |

```bash
# Configure a specific IDE
orchestra init --ide=cursor

# Configure ALL supported IDEs at once
orchestra init --all
```

## Architecture

```
Agent (Claude, GPT, Gemini, etc.)
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

## CLI Commands

```bash
orchestra init      # Initialize MCP configs for your IDE(s)
orchestra serve     # Start the MCP stdio server (used by IDE configs automatically)
orchestra version   # Print version info
orchestra help      # Show usage
```

### `orchestra init`

```
--workspace=DIR   Project directory (default: current directory)
--ide=NAME        Target IDE: claude, cursor, vscode, cline, windsurf, codex, gemini, zed, continue
--all             Generate configs for all supported IDEs
```

### `orchestra serve`

```
--workspace=DIR   Project workspace directory (default: current directory)
--certs-dir=DIR   mTLS certificates directory (default: ~/.orchestra/certs)
--log=FILE        Log file path (default: .orchestra-mcp.log)
```

## Binaries

| Binary | Description |
|--------|-------------|
| `orchestra` | CLI — init, serve, version |
| `orchestrator` | Central hub — starts plugins, routes QUIC messages |
| `storage-markdown` | File storage with `<!-- META {json} META -->` + Markdown body |
| `tools-features` | 34 feature-driven workflow tools |
| `transport-stdio` | MCP JSON-RPC bridge (stdin/stdout to QUIC) |

All 5 binaries are co-located in the same directory. `orchestra serve` finds the other 4 as siblings.

## Feature Workflow (34 Tools)

The tools.features plugin implements an 11-state feature lifecycle:

```
backlog → todo → in-progress → ready-for-testing → in-testing →
  ready-for-docs → in-docs → documented → in-review → done
                                              │
                      needs-edits ◄───────────┘
```

### Tool Categories

| Category | Count | Tools |
|----------|-------|-------|
| **Project** | 4 | `create_project`, `list_projects`, `delete_project`, `get_project_status` |
| **Feature** | 6 | `create_feature`, `get_feature`, `update_feature`, `list_features`, `delete_feature`, `search_features` |
| **Workflow** | 5 | `advance_feature`, `reject_feature`, `get_next_feature`, `set_current_feature`, `get_workflow_status` |
| **Review** | 3 | `request_review`, `submit_review`, `get_pending_reviews` |
| **Dependencies** | 3 | `add_dependency`, `remove_dependency`, `get_dependency_graph` |
| **WIP Limits** | 3 | `set_wip_limits`, `get_wip_limits`, `check_wip_limit` |
| **Reporting** | 3 | `get_progress`, `get_blocked_features`, `get_review_queue` |
| **Metadata** | 7 | `add_labels`, `remove_labels`, `assign_feature`, `unassign_feature`, `set_estimate`, `save_note`, `list_notes` |

All tools return formatted Markdown (tables, headings, bold text).

## Project Structure

```
orchestra-agents/
├── cmd/orchestra/                     # CLI binary (init, serve, version)
│   ├── main.go
│   └── internal/
│       ├── serve.go                   #   MCP server launcher
│       ├── initcmd.go                 #   IDE config generator
│       ├── ide.go                     #   9 IDE format definitions
│       ├── detect.go                  #   Project name/type detection
│       └── version.go                 #   Version info
├── proto/                             # Protobuf contract (plugin.proto)
├── gen/go/                            # Generated Go code from proto
├── libs/go/                           # Plugin SDK
│   ├── plugin/                        #   QUIC server/client, mTLS, framing
│   ├── types/                         #   Feature, Project, Workflow types
│   ├── helpers/                       #   Args, results, paths, validation, formatters
│   └── protocol/                      #   JSON-RPC + MCP types
├── services/orchestrator/             # Orchestrator service
│   ├── cmd/main.go
│   └── internal/                      #   Config, loader, router, server
├── plugins/
│   ├── storage-markdown/              # Markdown file storage plugin
│   ├── tools-features/                # 34 workflow tools plugin
│   └── transport-stdio/               # MCP stdio bridge plugin
├── packaging/
│   ├── homebrew/orchestra.rb          # Homebrew formula template
│   └── npm/                           # npm wrapper package
├── scripts/
│   ├── install.sh                     # curl | sh installer
│   ├── mcp-serve.sh                   # Legacy MCP launcher (use orchestra serve instead)
│   └── test-e2e.sh                    # End-to-end integration test
├── .github/workflows/
│   ├── ci.yml                         # CI: build + test on push/PR
│   └── release.yml                    # Release: cross-compile on tag push
├── docs/
│   ├── artifacts/                     # 17 design artifacts
│   └── implementation/                # Step-by-step build docs (01-08)
├── plugins.yaml                       # Default plugin configuration
├── go.work                            # Go workspace (7 modules)
└── Makefile                           # Build, test, install, release, clean
```

## Makefile Targets

```bash
make build              # Build all 5 binaries to bin/
make test               # Run all unit tests (66 tests)
make test-e2e           # Build + run end-to-end integration test
make install            # Install all 5 binaries to /usr/local/bin
make release            # Cross-compile for darwin/linux × amd64/arm64
make clean              # Remove build artifacts and certs
make proto              # Lint + generate proto code
```

## Plugin SDK

Build a plugin in ~50 lines of Go:

```go
package main

import (
    "context"
    "github.com/orchestra-mcp/sdk-go/plugin"
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

## Module Paths

```
github.com/orchestra-mcp/gen-go                    # Generated protobuf code
github.com/orchestra-mcp/sdk-go                    # Plugin SDK
github.com/orchestra-mcp/orchestrator              # Central hub
github.com/orchestra-mcp/plugin-storage-markdown   # Storage plugin
github.com/orchestra-mcp/plugin-tools-features     # Tools plugin
github.com/orchestra-mcp/plugin-transport-stdio    # Transport plugin
github.com/orchestra-mcp/cli                       # CLI binary
```

## License

MIT

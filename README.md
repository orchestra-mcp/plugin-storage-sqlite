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

# 2. That's it — your AI IDE now has 36 project management tools
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
tools.features (36 tools)    storage.markdown (disk)
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
| `storage-markdown` | File storage with YAML frontmatter metadata + Markdown body |
| `tools-features` | 36 feature-driven workflow tools |
| `transport-stdio` | MCP JSON-RPC bridge (stdin/stdout to QUIC) |

All 5 binaries are co-located in the same directory. `orchestra serve` finds the other 4 as siblings.

## Feature Workflow (36 Tools)

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
| **Dependencies** | 4 | `add_dependency`, `remove_dependency`, `get_dependency_graph`, `get_blocked_features` |
| **WIP Limits** | 3 | `set_wip_limits`, `get_wip_limits`, `check_wip_limit` |
| **Reporting** | 3 | `get_progress`, `get_review_queue`, `get_blocked_features` |
| **Metadata** | 8 | `add_labels`, `remove_labels`, `assign_feature`, `unassign_feature`, `set_estimate`, `save_note`, `list_notes` |

All tools return formatted Markdown (tables, headings, bold text).

## Dependency Management

Orchestra uses `orchestra.json` and `orchestra.lock` (similar to `composer.json`/`composer.lock`) to manage packages:

```bash
# orchestra.json — declares required packages with version constraints
{
    "require": {
        "orchestra-mcp/sdk-go": "^0.1.0",
        "orchestra-mcp/plugin-tools-features": "^0.1.0"
    }
}

# orchestra.lock — pins exact resolved versions (committed to repo)
```

Each plugin also has its own `orchestra.json` declaring its identity and dependencies. CI workflows, sync scripts, and release scripts all read from `orchestra.lock` as the single source of truth.

## Project Structure

```
framework/
├── libs/                              # All packages (each is a separate GitHub repo)
│   ├── proto/                         #   Protobuf definitions (plugin.proto)
│   ├── gen-go/                        #   Generated Go code from proto
│   ├── sdk-go/                        #   Plugin SDK (QUIC, mTLS, framing, helpers)
│   ├── orchestrator/                  #   Central hub (config, loader, router, server)
│   ├── plugin-storage-markdown/       #   Markdown storage with YAML frontmatter
│   ├── plugin-tools-features/         #   36 feature workflow tools
│   ├── plugin-transport-stdio/        #   MCP JSON-RPC stdin/stdout bridge
│   └── cli/                           #   CLI binary (init, serve, version)
├── scripts/
│   ├── install.sh                     # curl | sh installer
│   ├── new-plugin.sh                  # Plugin generator (tools/storage/transport)
│   ├── sync-repos.sh                  # Push libs/ to individual GitHub repos
│   ├── release.sh                     # Sync + tag + create GitHub releases
│   ├── orchestra-fmt.sh               # Format & validate orchestra.json files
│   └── test-e2e.sh                    # End-to-end integration test
├── packaging/
│   ├── homebrew/orchestra.rb          # Homebrew formula template
│   └── npm/                           # npm wrapper package
├── .github/workflows/
│   ├── ci.yml                         # CI: build + test + vet on push/PR
│   └── release.yml                    # Release: cross-compile on tag push
├── docs/
│   ├── artifacts/                     # 17 design artifacts
│   └── implementation/                # Step-by-step build docs (01-07)
├── orchestra.json                     # Package manifest (like composer.json)
├── orchestra.lock                     # Pinned versions (like composer.lock)
├── plugins.yaml                       # Default plugin runtime configuration
├── go.work                            # Go workspace (7 modules)
└── Makefile                           # Build, test, install, release, clean
```

## Makefile Targets

```bash
make build              # Build all 5 binaries to bin/
make test               # Run all unit tests (62 tests)
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

### Scaffold a New Plugin

```bash
./scripts/new-plugin.sh my-plugin tools       # Tools plugin (provides MCP tools)
./scripts/new-plugin.sh my-plugin storage     # Storage plugin (provides data backend)
./scripts/new-plugin.sh my-plugin transport   # Transport plugin (bridges protocols)
```

The generator creates all files (`cmd/main.go`, `internal/`, `go.mod`, `orchestra.json`, docs, CI), updates `go.work`, `Makefile`, `orchestra.json`, and `orchestra.lock` automatically.

## Scripts

```bash
# Sync local libs/ to individual GitHub repos
./scripts/sync-repos.sh                     # Sync all 8 packages
./scripts/sync-repos.sh sdk-go cli          # Sync specific packages
./scripts/sync-repos.sh --dry-run           # Preview without pushing

# Release all packages
./scripts/release.sh v0.2.0                 # Sync + tag + release all
./scripts/release.sh v0.2.0 --force         # Re-tag existing version
./scripts/release.sh v0.2.0 --dry-run       # Preview

# Format & validate orchestra.json files
./scripts/orchestra-fmt.sh                  # Format all
./scripts/orchestra-fmt.sh --check          # CI check (exit 1 if unformatted)
./scripts/orchestra-fmt.sh --validate       # Validate cross-references
```

## Protocol

All communication uses length-delimited Protobuf over QUIC with mTLS:

```
[4 bytes big-endian uint32 length][N bytes Protobuf PluginRequest/Response]
```

- **mTLS**: Auto-generated ed25519 CA + per-plugin certificates at `~/.orchestra/certs/`
- **Plugin startup**: Binary prints `READY <addr>` to stderr, orchestrator connects and sends Register + Boot
- **Storage format**: YAML frontmatter (`---` delimiters) + Markdown body

## Module Paths

All packages are published as independent Go modules:

```
github.com/orchestra-mcp/proto                     # Protobuf definitions
github.com/orchestra-mcp/gen-go                    # Generated protobuf code
github.com/orchestra-mcp/sdk-go                    # Plugin SDK
github.com/orchestra-mcp/orchestrator              # Central hub
github.com/orchestra-mcp/plugin-storage-markdown   # Storage plugin
github.com/orchestra-mcp/plugin-tools-features     # Tools plugin
github.com/orchestra-mcp/plugin-transport-stdio    # Transport plugin
github.com/orchestra-mcp/cli                       # CLI binary
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, coding standards, and PR process.

## License

MIT

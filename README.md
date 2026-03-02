# Orchestra Framework

An AI-agentic IDE framework with 290 MCP tools across 36 plugins. Single-process in-process architecture — 4 core plugins bundled, 32 optional plugins installable separately.

## Install

```bash
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/orchestra-mcp/framework/master/scripts/install.sh | sh

# Windows (Git Bash)
curl -fsSL --ssl-revoke-best-effort https://raw.githubusercontent.com/orchestra-mcp/framework/master/scripts/install.sh | sh

# From source
git clone https://github.com/orchestra-mcp/framework.git
cd framework && make install
```

**Supported platforms**: macOS (amd64, arm64), Linux (amd64, arm64), Windows (amd64 via Git Bash)

## Quick Start

```bash
# 1. Initialize Orchestra in your project (auto-detects your IDE)
cd your-project
orchestra init

# 2. That's it — your AI IDE now has 290 tools + 5 prompts
```

`orchestra init` detects your IDE and writes the correct MCP config:

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
orchestra serve (single process)
  ├── storage.markdown        (in-process, disk storage)
  ├── tools.features          (in-process, 34 workflow tools)
  ├── tools.marketplace       (in-process, 15 tools + 5 prompts)
  ├── transport.stdio         (in-process, JSON-RPC bridge)
  │
  ├── TCP server :50101       (for desktop apps — Swift, Windows, Linux)
  │
  └── External plugins        (QUIC + mTLS, optional)
      ├── bridge.claude       (5 tools — AI provider)
      ├── bridge.openai       (5 tools — AI provider)
      ├── bridge.gemini       (5 tools — AI provider)
      ├── agent.orchestrator  (20 tools — multi-agent)
      ├── devtools.*          (80+ tools — git, docker, ssh, terminal, etc.)
      └── 27 more plugins...
```

**Single-process** — 4 core plugins run in-process via direct Go function calls. Optional plugins connect over QUIC with mTLS and Protobuf framing.

## CLI Commands

```bash
orchestra init                          # Initialize MCP configs for your IDE(s)
orchestra serve                         # Start the MCP server
orchestra version                       # Print version info
orchestra update                        # Self-update to latest release
orchestra plugin install <name>         # Install an optional plugin
orchestra plugin remove <name>          # Remove a plugin
orchestra plugin list                   # List installed plugins
orchestra plugin search <query>         # Search available plugins
orchestra pack install <repo>[@ver]     # Install a content pack (skills/agents/hooks)
orchestra pack remove <name>            # Remove an installed pack
orchestra pack list                     # List installed packs
orchestra pack search <query>           # Search available packs
orchestra pack recommend                # Recommend packs for your project
```

## Plugins (36 total)

### Core Plugins (bundled, always available)

| Plugin | Tools | Description |
|--------|-------|-------------|
| `storage.markdown` | — | File storage with YAML frontmatter + Markdown body |
| `tools.features` | 34 | Feature-driven workflow (11-state lifecycle) |
| `tools.marketplace` | 15+5p | Pack management and marketplace |
| `transport.stdio` | — | MCP JSON-RPC bridge (stdin/stdout) |

### Optional Plugins (install separately)

| Category | Plugins | Tools |
|----------|---------|-------|
| **AI Bridges** | bridge.claude, bridge.openai, bridge.gemini, bridge.ollama, bridge.firecrawl | 25 |
| **Agent** | agent.orchestrator | 20 |
| **DevTools** | devtools.git, devtools.docker, devtools.terminal, devtools.ssh, devtools.file-explorer, devtools.database, devtools.debugger, devtools.test-runner, devtools.log-viewer, devtools.services, devtools.devops, devtools.components | 110+ |
| **AI Awareness** | ai.screenshot, ai.vision, ai.browser-context, ai.screen-reader | 25 |
| **Tools** | tools.agentops, tools.sessions, tools.workspace, tools.notes, tools.docs, tools.markdown, tools.extension-generator | 56 |
| **Services** | services.voice, services.notifications | 16 |
| **Integration** | integration.figma | 6 |
| **Transport** | transport.quic-bridge, transport.webtransport | — |

```bash
# Install a plugin
orchestra plugin install bridge-claude

# Search plugins
orchestra plugin search "ai"
```

## Content Packs

Packs are installable bundles of skills (slash commands), agents (specialized sub-agents), and hooks (shell scripts). 17 official packs:

| Pack | Stacks | Contents |
|------|--------|----------|
| `pack-essentials` | all | project-manager, qa-testing, docs, scrum-master, devops |
| `pack-go-backend` | go | go-backend skill, go-architect + qa-go agents |
| `pack-rust-engine` | rust | rust-engine skill, rust-engineer + qa-rust agents |
| `pack-react-frontend` | react, ts | typescript-react, ui-design, tailwind skills |
| `pack-database` | all | database-sync skill, dba + postgres/sqlite/redis agents |
| `pack-ai` | all | ai-agentic skill, ai-engineer + lancedb agents |
| `pack-proto` | go, rust | proto-grpc skill, quic-protocol agent |
| `pack-desktop` | go, swift | wails, macos, native-widgets skills |
| `pack-mobile` | react | react-native skill, mobile-dev agent |
| `pack-chrome` | typescript | chrome-extension skill |
| `pack-infra` | all | gcp-infrastructure skill, devops agent |
| `pack-extensions` | all | native-extensions, raycast, vscode skills |
| `pack-analytics` | all | clickhouse-engineer agent |
| `pack-native-swift` | swift | swift-plugin agent |
| `pack-native-kotlin` | kotlin | kotlin-plugin agent |
| `pack-native-csharp` | csharp | csharp-plugin agent |
| `pack-native-gtk` | c | gtk-plugin agent |

```bash
orchestra pack install github.com/orchestra-mcp/pack-go-backend
orchestra pack recommend    # Auto-detect stacks and suggest packs
orchestra pack search "react"
```

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
| **Dependencies** | 4 | `add_dependency`, `remove_dependency`, `get_dependency_graph`, `get_blocked_features` |
| **WIP Limits** | 3 | `set_wip_limits`, `get_wip_limits`, `check_wip_limit` |
| **Reporting** | 3 | `get_progress`, `get_review_queue`, `get_blocked_features` |
| **Metadata** | 8 | `add_labels`, `remove_labels`, `assign_feature`, `unassign_feature`, `set_estimate`, `save_note`, `list_notes` |

## Project Structure

```
framework/
├── libs/                              # All packages (each is a separate GitHub repo)
│   ├── proto/                         #   Protobuf definitions (plugin.proto)
│   ├── gen-go/                        #   Generated Go code from proto
│   ├── sdk-go/                        #   Plugin SDK (QUIC, mTLS, framing, helpers)
│   ├── orchestrator/                  #   Central hub (config, loader, router, server)
│   ├── cli/                           #   CLI binary (init, serve, version, pack, plugin)
│   ├── plugin-storage-markdown/       #   Core: Markdown storage
│   ├── plugin-tools-features/         #   Core: 34 feature workflow tools
│   ├── plugin-tools-marketplace/      #   Core: pack management + marketplace
│   ├── plugin-transport-stdio/        #   Core: MCP JSON-RPC bridge
│   ├── plugin-bridge-claude/          #   Optional: Claude AI provider
│   ├── plugin-bridge-openai/          #   Optional: OpenAI provider
│   ├── plugin-bridge-gemini/          #   Optional: Gemini provider
│   ├── plugin-agent-orchestrator/     #   Optional: Multi-agent orchestration
│   ├── plugin-devtools-*/             #   Optional: 12 devtools plugins
│   ├── plugin-ai-*/                   #   Optional: 4 AI awareness plugins
│   └── ... (36 plugins total)
├── packs/                             # 17 installable content packs
├── scripts/
│   ├── install.sh                     # curl | sh installer (macOS/Linux/Windows)
│   ├── new-plugin.sh                  # Plugin generator
│   ├── sync-repos.sh                  # Push libs/ to individual GitHub repos
│   ├── release.sh                     # Tag + create GitHub releases
│   ├── ship.sh                        # Full ship pipeline (build, test, sync, release)
│   └── test-e2e.sh                    # End-to-end integration test
├── .github/workflows/
│   ├── ci.yml                         # CI: build + test + vet on push/PR
│   └── release.yml                    # Release: 5-platform cross-compile on tag push
├── orchestra.json                     # Package manifest
├── orchestra.lock                     # Pinned versions (44 packages)
├── go.work                            # Go workspace
└── Makefile                           # Build, test, install, release
```

## Makefile Targets

```bash
make build              # Build all binaries to bin/
make test               # Run all unit tests
make test-e2e           # Build + run end-to-end integration test
make install            # Install binaries to /usr/local/bin
make release            # Cross-compile for darwin/linux × amd64/arm64 + windows/amd64
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
./scripts/new-plugin.sh my-plugin tools       # Tools plugin
./scripts/new-plugin.sh my-plugin storage     # Storage plugin
./scripts/new-plugin.sh my-plugin transport   # Transport plugin
```

## Scripts

```bash
# Ship a release (build, test, sync, tag, release)
./scripts/ship.sh v1.0.0
./scripts/ship.sh v1.0.0 --skip-pr --message "Release v1.0.0"
./scripts/ship.sh v1.0.0 --dry-run

# Sync local libs/ to individual GitHub repos
./scripts/sync-repos.sh                     # Sync all 44 packages
./scripts/sync-repos.sh sdk-go cli          # Sync specific packages

# Release all packages
./scripts/release.sh v1.0.0                 # Tag + release all 44 repos
./scripts/release.sh v1.0.0 --force         # Re-tag existing version
```

## Protocol

All inter-plugin communication uses length-delimited Protobuf over QUIC with mTLS:

```
[4 bytes big-endian uint32 length][N bytes Protobuf PluginRequest/Response]
```

- **mTLS**: Auto-generated ed25519 CA + per-plugin certificates at `~/.orchestra/certs/`
- **In-process**: Core plugins bypass QUIC entirely — direct Go function calls
- **TCP**: Desktop apps (Swift, Windows, Linux) connect via length-delimited Protobuf on port 50101
- **Storage format**: YAML frontmatter (`---` delimiters) + Markdown body

## Module Paths

All 44 packages are published as independent Go modules under `github.com/orchestra-mcp/`:

```
proto, gen-go, sdk-go, orchestrator, cli,
plugin-storage-markdown, plugin-tools-features, plugin-tools-marketplace,
plugin-transport-stdio, plugin-bridge-claude, plugin-bridge-openai,
plugin-bridge-gemini, plugin-bridge-ollama, plugin-bridge-firecrawl,
plugin-agent-orchestrator, plugin-tools-agentops, plugin-tools-sessions,
plugin-tools-workspace, plugin-tools-notes, plugin-tools-docs,
plugin-tools-markdown, plugin-tools-extension-generator,
plugin-devtools-git, plugin-devtools-docker, plugin-devtools-terminal,
plugin-devtools-ssh, plugin-devtools-file-explorer, plugin-devtools-database,
plugin-devtools-debugger, plugin-devtools-test-runner, plugin-devtools-log-viewer,
plugin-devtools-services, plugin-devtools-devops, plugin-devtools-components,
plugin-ai-screenshot, plugin-ai-vision, plugin-ai-browser-context,
plugin-ai-screen-reader, plugin-services-voice, plugin-services-notifications,
plugin-integration-figma, plugin-transport-quic-bridge,
plugin-transport-webtransport, plugin-engine-rag
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, coding standards, and PR process.

## License

MIT

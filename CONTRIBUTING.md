# Contributing to Orchestra MCP Framework

Thank you for your interest in contributing!

## Prerequisites

- **Go 1.23+** — [install](https://go.dev/dl/)
- **buf CLI** — [install](https://buf.build/docs/installation/) (only for proto changes)
- **git**

## Development Setup

```bash
# Clone with all modules
git clone https://github.com/orchestra-mcp/framework.git
cd framework

# Build all binaries
make build

# Run all tests
make test

# Run end-to-end test
make test-e2e
```

## Project Structure

```
libs/
├── proto/                       # Protobuf definitions
├── gen-go/                      # Generated Go code
├── sdk-go/                      # Plugin SDK
├── orchestrator/                # Central hub service
├── plugin-storage-markdown/     # File storage plugin
├── plugin-tools-features/       # 34 workflow tools
├── plugin-transport-stdio/      # MCP JSON-RPC bridge
└── cli/                         # orchestra CLI binary
```

Each directory under `libs/` is an independent Go module with its own repo under `github.com/orchestra-mcp/`.

## Creating a New Plugin

```bash
./scripts/new-plugin.sh my-plugin tools       # Tools plugin
./scripts/new-plugin.sh my-plugin storage     # Storage plugin
./scripts/new-plugin.sh my-plugin transport   # Transport plugin
```

## Coding Standards

- Run `gofmt` on all Go files
- Run `go vet ./...` before committing
- Write tests for new functionality
- Use `context.Context` through the call chain
- Never use `replace` directives in `go.mod` — the `go.work` file handles local development

## Pull Request Process

1. Fork the repository
2. Create a feature branch from `master`
3. Make your changes with tests
4. Ensure `make build` and `make test` pass
5. Submit a pull request

## Releasing

```bash
./scripts/release.sh v0.2.0              # Sync + tag all repos
./scripts/release.sh v0.2.0 --dry-run    # Preview first
```

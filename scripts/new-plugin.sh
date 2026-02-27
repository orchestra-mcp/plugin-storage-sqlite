#!/usr/bin/env bash
# new-plugin.sh — Scaffold a new Orchestra plugin from template.
#
# Usage:
#   ./scripts/new-plugin.sh my-plugin tools       # Tools plugin (provides MCP tools)
#   ./scripts/new-plugin.sh my-plugin storage     # Storage plugin (provides data backend)
#   ./scripts/new-plugin.sh my-plugin transport   # Transport plugin (bridges protocols)

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ORG="orchestra-mcp"
MODULE_PREFIX="github.com/${ORG}"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
info()  { printf "${CYAN}[new]${NC}   %s\n" "$*"; }
ok()    { printf "${GREEN}[ok]${NC}    %s\n" "$*"; }
warn()  { printf "${YELLOW}[warn]${NC}  %s\n" "$*"; }
fail()  { printf "${RED}[error]${NC} %s\n" "$*" >&2; exit 1; }

# ──────────────────────────────────────────────────────
# Parse arguments
# ──────────────────────────────────────────────────────

if [[ $# -lt 2 ]]; then
  echo "Usage: $0 <name> <type>"
  echo ""
  echo "  name    Plugin name (e.g. my-plugin -> libs/plugin-my-plugin/)"
  echo "  type    Plugin type: tools, storage, transport"
  echo ""
  echo "Examples:"
  echo "  $0 my-plugin tools"
  echo "  $0 my-storage storage"
  echo "  $0 my-bridge transport"
  exit 1
fi

NAME="$1"
TYPE="$2"
PLUGIN_DIR="${ROOT}/libs/plugin-${NAME}"
MODULE="${MODULE_PREFIX}/plugin-${NAME}"
BINARY="${NAME}"
PLUGIN_ID="${TYPE}.${NAME}"

# Validate type
case "$TYPE" in
  tools|storage|transport) ;;
  *) fail "Invalid type: $TYPE (must be tools, storage, or transport)" ;;
esac

# Check not already exists
if [[ -d "$PLUGIN_DIR" ]]; then
  fail "Directory already exists: libs/plugin-${NAME}/"
fi

info "Scaffolding ${TYPE} plugin: plugin-${NAME}"
echo ""

# ──────────────────────────────────────────────────────
# Create directory structure
# ──────────────────────────────────────────────────────

mkdir -p "${PLUGIN_DIR}"/{cmd,.github/ISSUE_TEMPLATE,.github/workflows,docs,internal}

case "$TYPE" in
  tools)
    mkdir -p "${PLUGIN_DIR}/internal/tools"
    mkdir -p "${PLUGIN_DIR}/internal/storage"
    ;;
esac

# ──────────────────────────────────────────────────────
# Common files (all types)
# ──────────────────────────────────────────────────────

# go.mod
cat > "${PLUGIN_DIR}/go.mod" <<EOF
module ${MODULE}

go 1.23

require (
	github.com/orchestra-mcp/gen-go v0.1.0
	github.com/orchestra-mcp/sdk-go v0.1.0
	google.golang.org/protobuf v1.36.11
)
EOF

# .gitignore
cat > "${PLUGIN_DIR}/.gitignore" <<EOF
# Binaries
${BINARY}
bin/

# Go
*.test
*.out
vendor/

# Test artifacts
.projects/

# OS
.DS_Store
Thumbs.db

# IDE
.idea/
.vscode/
*.swp
*.swo
EOF

# LICENSE
cat > "${PLUGIN_DIR}/LICENSE" <<'EOF'
MIT License

Copyright (c) 2026 Orchestra MCP Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF

# CODE_OF_CONDUCT.md
cat > "${PLUGIN_DIR}/CODE_OF_CONDUCT.md" <<'EOF'
# Contributor Covenant Code of Conduct

## Our Pledge

We as members, contributors, and leaders pledge to make participation in our
community a harassment-free experience for everyone.

## Our Standards

Examples of behavior that contributes to a positive environment:

* Demonstrating empathy and kindness toward other people
* Being respectful of differing opinions, viewpoints, and experiences
* Giving and gracefully accepting constructive feedback

Examples of unacceptable behavior:

* Trolling, insulting or derogatory comments, and personal or political attacks
* Public or private harassment
* Publishing others' private information without their explicit permission

## Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be
reported to the community leaders.

## Attribution

This Code of Conduct is adapted from the [Contributor Covenant](https://www.contributor-covenant.org), version 2.0.
EOF

# SECURITY.md
cat > "${PLUGIN_DIR}/SECURITY.md" <<'EOF'
# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in this project, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please use [GitHub's private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing-information-about-vulnerabilities/privately-reporting-a-security-vulnerability).

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest  | Yes       |

## Response Time

We aim to acknowledge security reports within 48 hours and provide a fix within 7 days for critical issues.
EOF

# CHANGELOG.md
cat > "${PLUGIN_DIR}/CHANGELOG.md" <<'EOF'
# Changelog

## [0.1.0] - Initial Release

- Scaffolded from `scripts/new-plugin.sh`
EOF

# .github/workflows/ci.yml
cat > "${PLUGIN_DIR}/.github/workflows/ci.yml" <<'EOF'
name: CI

on:
  push:
    branches: [master]
  pull_request:
    branches: [master]

permissions:
  contents: read

jobs:
  test:
    name: Test & Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"

      - name: Download dependencies
        run: go mod download

      - name: Vet
        run: go vet ./...

      - name: Test
        run: go test ./...

      - name: Build
        run: go build ./cmd/
EOF

# .github/ISSUE_TEMPLATE/bug_report.md
cat > "${PLUGIN_DIR}/.github/ISSUE_TEMPLATE/bug_report.md" <<EOF
---
name: Bug Report
about: Report a bug in plugin-${NAME}
title: "[Bug] "
labels: bug
assignees: ''
---

## Description
A clear description of the bug.

## Steps to Reproduce
1. ...
2. ...

## Expected Behavior
What you expected to happen.

## Actual Behavior
What actually happened.

## Environment
- OS:
- Go version:
- Plugin version:
EOF

# .github/ISSUE_TEMPLATE/feature_request.md
cat > "${PLUGIN_DIR}/.github/ISSUE_TEMPLATE/feature_request.md" <<EOF
---
name: Feature Request
about: Suggest a feature for plugin-${NAME}
title: "[Feature] "
labels: enhancement
assignees: ''
---

## Description
A clear description of the feature you'd like.

## Use Case
Why is this feature needed?

## Proposed Solution
How should it work?
EOF

# .github/ISSUE_TEMPLATE/config.yml
cat > "${PLUGIN_DIR}/.github/ISSUE_TEMPLATE/config.yml" <<'EOF'
blank_issues_enabled: true
EOF

# docs/CONTRIBUTING.md
cat > "${PLUGIN_DIR}/docs/CONTRIBUTING.md" <<EOF
# Contributing to plugin-${NAME}

## Prerequisites

- Go 1.23+
- git

## Development

\`\`\`bash
# Build
go build ./cmd/

# Test
go test ./... -v

# Vet
go vet ./...
\`\`\`

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: \`go test ./...\`
5. Submit a pull request
EOF

# ──────────────────────────────────────────────────────
# Type-specific files
# ──────────────────────────────────────────────────────

case "$TYPE" in

# =====================================================
# TOOLS TYPE
# =====================================================
tools)

# cmd/main.go
cat > "${PLUGIN_DIR}/cmd/main.go" <<GOEOF
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pluginv1 "${MODULE_PREFIX}/gen-go/orchestra/plugin/v1"
	"${MODULE_PREFIX}/sdk-go/plugin"
	"${MODULE}/internal"
	"${MODULE}/internal/storage"
)

func main() {
	builder := plugin.New("${PLUGIN_ID}").
		Version("0.1.0").
		Description("${NAME} tools plugin").
		Author("Orchestra").
		Binary("${BINARY}").
		NeedsStorage("markdown")

	adapter := &clientAdapter{}
	store := storage.NewDataStorage(adapter)

	tp := &internal.ToolsPlugin{Storage: store}
	tp.RegisterTools(builder)

	p := builder.BuildWithTools()
	p.ParseFlags()
	adapter.plugin = p

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	if err := p.Run(ctx); err != nil {
		log.Fatalf("${PLUGIN_ID}: %v", err)
	}
}

type clientAdapter struct {
	plugin *plugin.Plugin
}

func (a *clientAdapter) Send(ctx context.Context, req *pluginv1.PluginRequest) (*pluginv1.PluginResponse, error) {
	client := a.plugin.OrchestratorClient()
	if client == nil {
		return nil, fmt.Errorf("orchestrator client not connected")
	}
	return client.Send(ctx, req)
}
GOEOF

# internal/plugin.go
cat > "${PLUGIN_DIR}/internal/plugin.go" <<GOEOF
package internal

import (
	"${MODULE_PREFIX}/sdk-go/plugin"
	"${MODULE}/internal/storage"
	"${MODULE}/internal/tools"
)

// ToolsPlugin holds the storage reference and registers all tools.
type ToolsPlugin struct {
	Storage *storage.DataStorage
}

// RegisterTools registers all tools with the plugin builder.
func (tp *ToolsPlugin) RegisterTools(builder *plugin.PluginBuilder) {
	s := tp.Storage

	builder.RegisterTool("hello",
		"Say hello to someone",
		tools.HelloSchema(), tools.Hello(s))
}
GOEOF

# internal/storage/client.go
cat > "${PLUGIN_DIR}/internal/storage/client.go" <<GOEOF
package storage

import (
	"context"

	pluginv1 "${MODULE_PREFIX}/gen-go/orchestra/plugin/v1"
)

// StorageClient sends requests to the orchestrator for storage operations.
type StorageClient interface {
	Send(ctx context.Context, req *pluginv1.PluginRequest) (*pluginv1.PluginResponse, error)
}

// DataStorage wraps the storage client for tool handlers.
type DataStorage struct {
	client StorageClient
}

// NewDataStorage creates a new DataStorage with the given client.
func NewDataStorage(client StorageClient) *DataStorage {
	return &DataStorage{client: client}
}
GOEOF

# internal/tools/hello.go
cat > "${PLUGIN_DIR}/internal/tools/hello.go" <<GOEOF
package tools

import (
	"context"
	"fmt"

	pluginv1 "${MODULE_PREFIX}/gen-go/orchestra/plugin/v1"
	"${MODULE_PREFIX}/sdk-go/helpers"
	"${MODULE}/internal/storage"
	"google.golang.org/protobuf/types/known/structpb"
)

// HelloSchema returns the JSON Schema for the hello tool.
func HelloSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "Name to greet",
			},
		},
		"required": []any{"name"},
	})
	return s
}

// Hello returns a tool handler that greets someone by name.
func Hello(_ *storage.DataStorage) func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "name"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		name := helpers.GetString(req.Arguments, "name")
		return helpers.TextResult(fmt.Sprintf("Hello, %s!", name)), nil
	}
}
GOEOF

# README.md (tools)
cat > "${PLUGIN_DIR}/README.md" <<EOF
# Orchestra Plugin: ${NAME}

A tools plugin for the [Orchestra MCP](https://github.com/${ORG}/framework) framework.

## Install

\`\`\`bash
go install ${MODULE}/cmd@latest
\`\`\`

## Usage

Add to your \`plugins.yaml\`:

\`\`\`yaml
- id: ${PLUGIN_ID}
  binary: ./bin/${BINARY}
  enabled: true
\`\`\`

## Tools

| Tool | Description |
|------|-------------|
| \`hello\` | Say hello to someone |

## Related Packages

- [sdk-go](https://github.com/${ORG}/sdk-go) — Plugin SDK
- [gen-go](https://github.com/${ORG}/gen-go) — Generated Protobuf types
EOF

# docs/TOOLS_REFERENCE.md
cat > "${PLUGIN_DIR}/docs/TOOLS_REFERENCE.md" <<EOF
# Tools Reference

## hello

Say hello to someone.

### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| \`name\` | string | Yes | Name to greet |

### Example

\`\`\`json
{
  "name": "hello",
  "arguments": {"name": "World"}
}
\`\`\`

Response: \`"Hello, World!"\`
EOF

;; # end tools

# =====================================================
# STORAGE TYPE
# =====================================================
storage)

# cmd/main.go
cat > "${PLUGIN_DIR}/cmd/main.go" <<GOEOF
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"${MODULE_PREFIX}/sdk-go/plugin"
	"${MODULE}/internal"
)

func main() {
	workspace := flag.String("workspace", ".", "Root workspace directory")

	storage := internal.NewStoragePlugin(*workspace)

	p := plugin.New("${PLUGIN_ID}").
		Version("0.1.0").
		Description("${NAME} storage plugin").
		Author("Orchestra").
		Binary("${BINARY}").
		ProvidesStorage("${NAME}").
		SetStorageHandler(storage).
		BuildWithTools()

	p.ParseFlags()

	// Re-read workspace after flags are parsed.
	storage = internal.NewStoragePlugin(*workspace)
	p.Server().SetStorageHandler(storage)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	if err := p.Run(ctx); err != nil {
		log.Fatalf("${PLUGIN_ID}: %v", err)
	}
}
GOEOF

# internal/storage.go
cat > "${PLUGIN_DIR}/internal/storage.go" <<GOEOF
package internal

import (
	"context"
	"fmt"

	pluginv1 "${MODULE_PREFIX}/gen-go/orchestra/plugin/v1"
)

// StoragePlugin implements the storage handler interface.
type StoragePlugin struct {
	workspace string
}

// NewStoragePlugin creates a new StoragePlugin rooted at workspace.
func NewStoragePlugin(workspace string) *StoragePlugin {
	return &StoragePlugin{workspace: workspace}
}

// Read retrieves a document by path.
func (s *StoragePlugin) Read(ctx context.Context, req *pluginv1.StorageReadRequest) (*pluginv1.StorageReadResponse, error) {
	return nil, fmt.Errorf("not implemented: Read %s", req.GetPath())
}

// Write stores a document at path.
func (s *StoragePlugin) Write(ctx context.Context, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	return nil, fmt.Errorf("not implemented: Write %s", req.GetPath())
}

// Delete removes a document at path.
func (s *StoragePlugin) Delete(ctx context.Context, req *pluginv1.StorageDeleteRequest) (*pluginv1.StorageDeleteResponse, error) {
	return nil, fmt.Errorf("not implemented: Delete %s", req.GetPath())
}

// List returns documents matching a prefix.
func (s *StoragePlugin) List(ctx context.Context, req *pluginv1.StorageListRequest) (*pluginv1.StorageListResponse, error) {
	return nil, fmt.Errorf("not implemented: List %s", req.GetPrefix())
}
GOEOF

# internal/reader.go
cat > "${PLUGIN_DIR}/internal/reader.go" <<GOEOF
package internal

import (
	"google.golang.org/protobuf/types/known/structpb"
)

// ParseFile parses raw file data into metadata and body.
func ParseFile(data []byte) (metadata *structpb.Struct, body []byte, err error) {
	// TODO: implement parsing logic
	return nil, data, nil
}
GOEOF

# internal/writer.go
cat > "${PLUGIN_DIR}/internal/writer.go" <<GOEOF
package internal

import (
	"google.golang.org/protobuf/types/known/structpb"
)

// FormatFile serializes metadata and body into raw file data.
func FormatFile(metadata *structpb.Struct, body []byte) ([]byte, error) {
	// TODO: implement formatting logic
	return body, nil
}
GOEOF

# README.md (storage)
cat > "${PLUGIN_DIR}/README.md" <<EOF
# Orchestra Plugin: ${NAME}

A storage plugin for the [Orchestra MCP](https://github.com/${ORG}/framework) framework.

## Install

\`\`\`bash
go install ${MODULE}/cmd@latest
\`\`\`

## Usage

Add to your \`plugins.yaml\`:

\`\`\`yaml
- id: ${PLUGIN_ID}
  binary: ./bin/${BINARY}
  enabled: true
  provides_storage:
    - ${NAME}
  config:
    workspace: .
\`\`\`

## Supported Operations

| Operation | Description |
|-----------|-------------|
| Read | Retrieve a document by path |
| Write | Store a document at path |
| Delete | Remove a document |
| List | List documents by prefix |

## Related Packages

- [sdk-go](https://github.com/${ORG}/sdk-go) — Plugin SDK
- [gen-go](https://github.com/${ORG}/gen-go) — Generated Protobuf types
EOF

# docs/STORAGE_FORMAT.md
cat > "${PLUGIN_DIR}/docs/STORAGE_FORMAT.md" <<EOF
# Storage Format

## Overview

This plugin stores data in the workspace directory.

## File Layout

\`\`\`
{workspace}/.projects/{project_id}/
  features/{feature_id}.ext
\`\`\`

## Format Details

TODO: Document the storage format for this plugin.
EOF

;; # end storage

# =====================================================
# TRANSPORT TYPE
# =====================================================
transport)

# cmd/main.go
cat > "${PLUGIN_DIR}/cmd/main.go" <<GOEOF
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"${MODULE_PREFIX}/sdk-go/plugin"
	"${MODULE}/internal"
)

func main() {
	orchestratorAddr := flag.String("orchestrator-addr", "localhost:9100", "Address of the orchestrator")
	certsDir := flag.String("certs-dir", plugin.DefaultCertsDir, "Directory for mTLS certificates")
	flag.Parse()

	if *orchestratorAddr == "" {
		log.Fatal("--orchestrator-addr is required")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	resolvedCertsDir := plugin.ResolveCertsDir(*certsDir)

	clientTLS, err := plugin.ClientTLSConfig(resolvedCertsDir, "${PLUGIN_ID}-client")
	if err != nil {
		log.Fatalf("client TLS config: %v", err)
	}

	client, err := plugin.NewOrchestratorClient(ctx, *orchestratorAddr, clientTLS)
	if err != nil {
		log.Fatalf("connect to orchestrator at %s: %v", *orchestratorAddr, err)
	}
	defer client.Close()

	fmt.Fprintf(os.Stderr, "${PLUGIN_ID}: connected to orchestrator at %s\n", *orchestratorAddr)

	transport := internal.NewTransport(client, os.Stdin, os.Stdout)
	if err := transport.Run(ctx); err != nil {
		if ctx.Err() != nil {
			fmt.Fprintf(os.Stderr, "${PLUGIN_ID}: shutting down\n")
			return
		}
		log.Fatalf("${PLUGIN_ID}: %v", err)
	}
}
GOEOF

# internal/transport.go
cat > "${PLUGIN_DIR}/internal/transport.go" <<GOEOF
package internal

import (
	"bufio"
	"context"
	"io"
	"sync"

	pluginv1 "${MODULE_PREFIX}/gen-go/orchestra/plugin/v1"
)

// Sender sends requests to the orchestrator.
type Sender interface {
	Send(ctx context.Context, req *pluginv1.PluginRequest) (*pluginv1.PluginResponse, error)
}

// Transport bridges an external protocol to the orchestrator.
type Transport struct {
	sender Sender
	reader *bufio.Scanner
	writer io.Writer
	mu     sync.Mutex
}

// NewTransport creates a new Transport.
func NewTransport(sender Sender, in io.Reader, out io.Writer) *Transport {
	return &Transport{
		sender: sender,
		reader: bufio.NewScanner(in),
		writer: out,
	}
}

// Run starts the transport read/write loop.
func (t *Transport) Run(ctx context.Context) error {
	// TODO: implement protocol read loop
	<-ctx.Done()
	return ctx.Err()
}
GOEOF

# internal/handler.go
cat > "${PLUGIN_DIR}/internal/handler.go" <<GOEOF
package internal

import (
	"fmt"
)

// dispatch routes an incoming request method to the appropriate handler.
func (t *Transport) dispatch(method string, params []byte) ([]byte, error) {
	switch method {
	case "initialize":
		return t.handleInitialize(params)
	case "ping":
		return []byte("{}"), nil
	default:
		return nil, fmt.Errorf("method not found: %s", method)
	}
}

func (t *Transport) handleInitialize(params []byte) ([]byte, error) {
	// TODO: implement initialization handshake
	return []byte("{}"), nil
}
GOEOF

# README.md (transport)
cat > "${PLUGIN_DIR}/README.md" <<EOF
# Orchestra Plugin: ${NAME}

A transport plugin for the [Orchestra MCP](https://github.com/${ORG}/framework) framework.

## Install

\`\`\`bash
go install ${MODULE}/cmd@latest
\`\`\`

## Usage

\`\`\`bash
${BINARY} --orchestrator-addr localhost:9100 --certs-dir ~/.orchestra/certs
\`\`\`

## How It Works

This transport plugin connects to the Orchestra orchestrator as a client and
bridges an external protocol to the internal QUIC mesh.

## Related Packages

- [sdk-go](https://github.com/${ORG}/sdk-go) — Plugin SDK
- [gen-go](https://github.com/${ORG}/gen-go) — Generated Protobuf types
EOF

# docs/PROTOCOL.md
cat > "${PLUGIN_DIR}/docs/PROTOCOL.md" <<EOF
# Protocol

## Overview

This transport plugin bridges an external protocol to the Orchestra orchestrator.

## Supported Methods

| Method | Description |
|--------|-------------|
| \`initialize\` | Protocol handshake |
| \`ping\` | Health check |

## Message Flow

\`\`\`
External Client <-> Transport Plugin <-> Orchestrator <-> Plugins
\`\`\`

## Details

TODO: Document the external protocol supported by this transport.
EOF

;; # end transport
esac

ok "Generated libs/plugin-${NAME}/ (${TYPE} type)"

# ──────────────────────────────────────────────────────
# Update go.work
# ──────────────────────────────────────────────────────

info "Adding to go.work..."
sed -i '' '/^)/i\
	./libs/plugin-'"${NAME}"'
' "${ROOT}/go.work"
ok "go.work updated"

# ──────────────────────────────────────────────────────
# Update Makefile
# ──────────────────────────────────────────────────────

info "Adding to Makefile..."

# Convert name to uppercase variable name (my-plugin -> MY_PLUGIN)
VAR_NAME="$(echo "${NAME}" | tr '[:lower:]-' '[:upper:]_')"

# Add binary variable after last existing binary variable
sed -i '' "/^ORCHESTRA_CLI/a\\
${VAR_NAME} := \$(BIN_DIR)/${BINARY}
" "${ROOT}/Makefile"

# Add to build: dependencies
sed -i '' "s/^build: \(.*\)/build: \1 build-${NAME}/" "${ROOT}/Makefile"

# Add to .PHONY
sed -i '' "s/^\.PHONY: \(.*\)/\.PHONY: \1 build-${NAME}/" "${ROOT}/Makefile"

# Add to BINARIES list
sed -i '' "s/^BINARIES := \(.*\)/BINARIES := \1 ${BINARY}/" "${ROOT}/Makefile"

# Add build target at end (before clean section)
cat >> "${ROOT}/Makefile" <<EOF

build-${NAME}:
	@mkdir -p \$(BIN_DIR)
	go build -o \$(${VAR_NAME}) ./libs/plugin-${NAME}/cmd/
EOF

# Add test line to test-unit
sed -i '' "/go test .\/libs\/plugin-transport-stdio/a\\
	go test ./libs/plugin-${NAME}/... -v
" "${ROOT}/Makefile"

ok "Makefile updated"

# ──────────────────────────────────────────────────────
# Generate orchestra.json for the plugin
# ──────────────────────────────────────────────────────

info "Generating orchestra.json..."
cat > "${PLUGIN_DIR}/orchestra.json" <<OJEOF
{
    "name": "${ORG}/plugin-${NAME}",
    "description": "${NAME} ${TYPE} plugin for Orchestra MCP.",
    "type": "${TYPE}",
    "license": "MIT",
    "homepage": "https://orchestra-mcp.dev",
    "support": {
        "issues": "https://github.com/${ORG}/plugin-${NAME}/issues",
        "source": "https://github.com/${ORG}/plugin-${NAME}"
    },
    "require": {
        "${ORG}/gen-go": "^0.1.0",
        "${ORG}/sdk-go": "^0.1.0"
    },
    "bin": "${BINARY}",
    "extra": {
        "plugin-id": "${PLUGIN_ID}",
        "provides": ["${TYPE}"],
        "branch-alias": {
            "dev-master": "0.1.x-dev"
        }
    }
}
OJEOF
ok "orchestra.json created"

# ──────────────────────────────────────────────────────
# Update root orchestra.json and orchestra.lock
# ──────────────────────────────────────────────────────

info "Adding to orchestra.json and orchestra.lock..."
python3 -c "
import json

# Update orchestra.json — add to require
with open('${ROOT}/orchestra.json') as f:
    manifest = json.load(f)
manifest['require']['${ORG}/plugin-${NAME}'] = '^0.1.0'
order = manifest.get('extra', {}).get('install-order', [])
if 'plugin-${NAME}' not in order:
    # Insert before 'cli' (last item) if it exists, otherwise append
    if 'cli' in order:
        idx = order.index('cli')
        order.insert(idx, 'plugin-${NAME}')
    else:
        order.append('plugin-${NAME}')
    manifest.setdefault('extra', {})['install-order'] = order
with open('${ROOT}/orchestra.json', 'w') as f:
    json.dump(manifest, f, indent=4)
    f.write('\n')

# Update orchestra.lock — add to packages
with open('${ROOT}/orchestra.lock') as f:
    lock = json.load(f)
new_pkg = {
    'name': '${ORG}/plugin-${NAME}',
    'version': 'v0.1.0',
    'source': {
        'type': 'git',
        'url': 'https://github.com/${ORG}/plugin-${NAME}.git',
        'reference': 'master'
    },
    'type': '${TYPE}',
    'description': '${NAME} ${TYPE} plugin for Orchestra MCP.',
    'path': 'libs/plugin-${NAME}',
    'binary': '${BINARY}',
    'require': {
        '${ORG}/gen-go': '^0.1.0',
        '${ORG}/sdk-go': '^0.1.0'
    }
}
# Insert before cli (last) if possible
pkgs = lock['packages']
cli_idx = next((i for i, p in enumerate(pkgs) if p['name'].endswith('/cli')), len(pkgs))
pkgs.insert(cli_idx, new_pkg)
with open('${ROOT}/orchestra.lock', 'w') as f:
    json.dump(lock, f, indent=4)
    f.write('\n')
"
ok "orchestra.json and orchestra.lock updated"

# ──────────────────────────────────────────────────────
# Run go mod tidy
# ──────────────────────────────────────────────────────

info "Running go mod tidy..."
cd "${PLUGIN_DIR}" && go mod tidy 2>&1
ok "Dependencies resolved"

# ──────────────────────────────────────────────────────
# Initialize git repo
# ──────────────────────────────────────────────────────

info "Initializing git repo..."
cd "${PLUGIN_DIR}"
git init -b master --quiet
git add -A
git commit -m "Initial scaffold from new-plugin.sh" --quiet
ok "Git initialized"

# ──────────────────────────────────────────────────────
# Summary
# ──────────────────────────────────────────────────────

echo ""
echo "────────────────────────────────────────────────────────"
ok "Plugin scaffolded: libs/plugin-${NAME}/ (${TYPE})"
echo ""
info "Next steps:"
info "  1. cd libs/plugin-${NAME}"
info "  2. Edit internal/ to implement your ${TYPE} logic"
info "  3. Test:  go test ./libs/plugin-${NAME}/..."
info "  4. Build: make build-${NAME}"
info "  5. To publish separately, create repo: ${ORG}/plugin-${NAME}"
info "     Then run: ./scripts/sync-repos.sh plugin-${NAME}"
echo ""

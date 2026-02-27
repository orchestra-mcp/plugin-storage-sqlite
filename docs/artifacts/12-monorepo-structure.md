# Monorepo Structure — Orchestra Rebuild

> Complete directory layout for the microservices rebuild. Each service is independently buildable.
> Proto definitions shared across 6 languages. Go workspace for local development.

---

## 1. Complete Directory Tree

```
orchestra-agents/
├── .github/
│   └── workflows/
│       ├── ci-go.yml                    # Go lint + test for all Go services
│       ├── ci-rust.yml                  # Cargo clippy + test for orch-engine
│       ├── ci-web.yml                   # TypeScript lint + vitest + Playwright
│       ├── ci-proto.yml                 # Buf lint + breaking change detection
│       ├── release-mcp.yml             # GoReleaser for orch-mcp binary
│       ├── release-server.yml          # Cloud Run deploy for orch-server
│       ├── release-engine.yml          # Cargo release for orch-engine
│       ├── release-desktop-macos.yml   # SwiftUI archive + notarize
│       ├── release-desktop-windows.yml # WinUI3 MSIX packaging
│       ├── release-android.yml         # Kotlin/Compose AAB upload
│       ├── release-ios.yml             # SwiftUI IPA via Xcode Cloud
│       └── release-chrome.yml          # Chrome Web Store publish
│
├── proto/                               # SINGLE SOURCE OF TRUTH for all .proto files
│   ├── buf.yaml                         # Buf module config
│   ├── buf.gen.yaml                     # Codegen for Go + Rust + TS + Swift + Kotlin + C#
│   ├── buf.lock
│   ├── orchestra/
│   │   ├── common/v1/
│   │   │   ├── common.proto             # Shared: UUID, Timestamp, Pagination, Error
│   │   │   └── sync.proto               # SyncEnvelope, VersionVector, ConflictResolution
│   │   ├── engine/v1/
│   │   │   ├── health.proto             # HealthService (ping, readiness, liveness)
│   │   │   ├── parse.proto              # ParseService (Tree-sitter: 14 languages)
│   │   │   ├── search.proto             # SearchService (Tantivy: index, query)
│   │   │   ├── memory.proto             # MemoryService (LanceDB: save, search)
│   │   │   └── component.proto          # ComponentBundlerService
│   │   ├── mcp/v1/
│   │   │   ├── project.proto            # Project, Epic, Story, Task CRUD
│   │   │   ├── workflow.proto           # WorkflowState enum, transitions
│   │   │   ├── sprint.proto             # Sprint, Burndown, Velocity, Standup
│   │   │   ├── prd.proto                # PrdSession, Questions, Answers, Templates
│   │   │   ├── memory.proto             # MemoryChunk, SessionLog (MCP-level)
│   │   │   └── integration.proto        # Integration config, credentials, sync records
│   │   ├── mesh/v1/
│   │   │   └── envelope.proto           # Envelope, ControlMessage, Handshake, Heartbeat
│   │   ├── server/v1/
│   │   │   ├── auth.proto               # AuthService (login, tokens, devices)
│   │   │   ├── user.proto               # UserService (profile, preferences, teams)
│   │   │   └── sync.proto               # SyncService (pull, push, conflicts)
│   │   └── gateway/v1/
│   │       ├── mesh.proto               # MeshService (discovery, topology, health)
│   │       └── events.proto             # EventService (publish, subscribe, stream)
│   └── third_party/                     # Vendored external protos (google/api, etc.)
│
├── gen/                                 # ALL GENERATED CODE
│   ├── go/orchestra/                    # Go protobuf (buf generate)
│   ├── rust/                            # Rust protobuf (tonic-build in engine)
│   ├── ts/orchestra/                    # TypeScript (connect-es)
│   ├── swift/Orchestra/                 # Swift (grpc-swift)
│   ├── kotlin/com/orchestra/            # Kotlin (grpc-kotlin)
│   └── csharp/Orchestra/               # C# (Grpc.Tools)
│
├── libs/                                # SHARED LIBRARIES
│   ├── go/                              # Shared Go packages
│   │   ├── go.mod                       # module github.com/orchestra-agents/libs/go
│   │   ├── quic/                        # QUIC mesh library
│   │   │   ├── mesh.go                  # Service mesh: discovery, connection pool
│   │   │   ├── transport.go             # QUIC transport with mTLS
│   │   │   ├── discovery.go             # mDNS + seed file discovery
│   │   │   ├── pubsub.go               # Event pub/sub over QUIC streams
│   │   │   └── health.go               # Health check protocol
│   │   ├── types/                       # Shared domain types
│   │   │   ├── issue.go                 # IssueData (epic/story/task/bug)
│   │   │   ├── project.go              # ProjectStatus, ProjectConfig
│   │   │   ├── workflow.go             # 13 WorkflowState constants, transitions
│   │   │   ├── sprint.go               # Sprint, SprintMetrics, Velocity
│   │   │   ├── prd.go                  # PrdSession, PrdQuestion, PrdAnswer
│   │   │   ├── memory.go              # MemoryChunk, SessionLog
│   │   │   ├── notification.go        # Notification, NotificationAction
│   │   │   └── sync.go                # SyncRecord, VersionVector
│   │   ├── markdown/                    # Protobuf metadata + Markdown content storage
│   │   │   ├── reader.go               # Read .md files + .pb sidecar metadata
│   │   │   ├── writer.go               # Write Protobuf metadata + Markdown body
│   │   │   └── reader_test.go
│   │   ├── helpers/                     # Shared utilities
│   │   │   ├── args.go                  # Argument validation
│   │   │   ├── paths.go                # Workspace path resolution
│   │   │   ├── results.go             # MCP result formatting
│   │   │   ├── strings.go             # String manipulation
│   │   │   └── validate.go            # Input validation
│   │   └── protocol/                   # Shared protocol definitions
│   │       ├── jsonrpc.go              # JSON-RPC 2.0 envelope types
│   │       └── mcp.go                  # MCP protocol constants
│   ├── c/                               # Shared C library (~1000 lines)
│   │   ├── CMakeLists.txt
│   │   ├── include/orch/
│   │   │   ├── pty.h                   # PTY allocation and I/O
│   │   │   ├── io.h                    # High-performance I/O primitives
│   │   │   └── crypto.h               # AES-256-GCM, HKDF key derivation
│   │   └── src/
│   │       ├── pty.c                    # POSIX + Windows ConPTY
│   │       ├── io.c                     # splice, sendfile, io_uring hints
│   │       └── crypto.c                # OS-level crypto primitives
│   └── ts/                              # Shared TypeScript packages
│       ├── package.json
│       ├── tsconfig.base.json
│       ├── shared/                      # @orchestra/shared (types, stores, hooks, API)
│       │   ├── package.json
│       │   └── src/
│       │       ├── types/
│       │       ├── stores/              # Zustand stores
│       │       ├── hooks/
│       │       ├── api/                 # Typed API client
│       │       └── index.ts
│       └── ui/                          # @orchestra/ui (shadcn/ui components)
│           ├── package.json
│           └── src/
│               ├── components/
│               ├── layouts/
│               ├── theme/               # Tailwind CSS v4 tokens
│               └── index.ts
│
├── services/                            # ALL BACKEND SERVICES
│   ├── orch-mcp/                        # MCP stdio server — the brain
│   │   ├── go.mod                       # Standalone Go module
│   │   ├── Makefile
│   │   ├── Dockerfile
│   │   ├── cmd/main.go                 # CLI: orch-mcp --workspace . | orch-mcp init
│   │   ├── internal/
│   │   │   ├── transport/              # MCP stdio + SSE JSON-RPC server
│   │   │   │   ├── server.go           # Read loop, request dispatch
│   │   │   │   ├── writer.go           # ResponseWriter + StdioWriter
│   │   │   │   └── sse.go             # SSE transport for browsers
│   │   │   ├── tools/                  # All MCP tools by category
│   │   │   │   ├── project.go          # 7 tools
│   │   │   │   ├── epic.go             # 5 tools
│   │   │   │   ├── story.go            # 5 tools
│   │   │   │   ├── task.go             # 5 tools
│   │   │   │   ├── workflow.go         # 5 tools
│   │   │   │   ├── lifecycle.go        # 2 tools (advance/reject)
│   │   │   │   ├── sprint.go           # 6 tools
│   │   │   │   ├── scrum.go            # 5 tools
│   │   │   │   ├── prd.go              # 14 tools
│   │   │   │   ├── prd_templates.go    # 3 tools
│   │   │   │   ├── dependency.go       # 3 tools
│   │   │   │   ├── metadata.go         # 7 tools
│   │   │   │   ├── memory.go           # 6 tools
│   │   │   │   ├── notes.go            # 6 tools
│   │   │   │   ├── artifacts.go        # 5 tools
│   │   │   │   ├── usage.go            # 3 tools
│   │   │   │   ├── claude.go           # 6 tools
│   │   │   │   ├── team.go
│   │   │   │   ├── bugfix.go           # 2 tools
│   │   │   │   ├── notification.go
│   │   │   │   ├── readme.go
│   │   │   │   ├── desktop.go
│   │   │   │   ├── devtools.go
│   │   │   │   ├── component.go
│   │   │   │   ├── figma.go
│   │   │   │   ├── github.go
│   │   │   │   ├── jira.go
│   │   │   │   ├── linear.go
│   │   │   │   ├── analytics.go
│   │   │   │   └── prompts.go
│   │   │   ├── engine/                 # Rust engine client + fallback
│   │   │   │   ├── bridge.go
│   │   │   │   ├── client.go
│   │   │   │   └── manager.go
│   │   │   └── bootstrap/init.go      # Workspace init (orch-mcp init)
│   │   ├── resources/                   # Bundled skills, agents, docs, hooks
│   │   │   ├── agents/                 # 16 agent .md files
│   │   │   ├── skills/                 # 20 skill directories
│   │   │   ├── hooks/
│   │   │   ├── docs/
│   │   │   └── embed.go               # go:embed directives
│   │   └── testdata/
│   │
│   ├── orch-engine/                     # Rust data engine
│   │   ├── Cargo.toml                   # Standalone Rust binary
│   │   ├── Makefile
│   │   ├── Dockerfile
│   │   ├── build.rs                     # tonic-build (reads ../../proto/)
│   │   └── src/
│   │       ├── main.rs
│   │       ├── lib.rs
│   │       ├── config.rs
│   │       ├── services/               # gRPC implementations
│   │       │   ├── health.rs
│   │       │   ├── parse.rs            # Tree-sitter (14 languages)
│   │       │   ├── search.rs           # Tantivy full-text
│   │       │   ├── memory.rs           # Vector embeddings
│   │       │   └── component.rs
│   │       ├── parser/                 # Tree-sitter internals
│   │       │   ├── registry.rs         # Language grammar registry
│   │       │   ├── symbols.rs
│   │       │   ├── navigation.rs
│   │       │   ├── highlight.rs
│   │       │   └── wrapper.rs
│   │       ├── index/                  # Tantivy search
│   │       │   ├── schema.rs
│   │       │   ├── writer.rs
│   │       │   ├── reader.rs
│   │       │   ├── pipeline.rs
│   │       │   └── config.rs
│   │       ├── memory/                 # Vector memory
│   │       │   ├── storage.rs
│   │       │   ├── embeddings.rs
│   │       │   ├── search.rs
│   │       │   ├── sessions.rs
│   │       │   ├── observations.rs
│   │       │   ├── summaries.rs
│   │       │   └── schema.rs
│   │       ├── lancedb/                # LanceDB vector store
│   │       │   ├── store.rs
│   │       │   └── schema.rs
│   │       ├── db/                     # SQLite management
│   │       │   ├── pool.rs
│   │       │   └── schema.rs
│   │       └── markdown/
│   │           ├── renderer.rs
│   │           └── types.rs
│   │
│   ├── orch-server/                     # Go API server (HTTP/3)
│   │   ├── go.mod                       # Standalone Go module
│   │   ├── Makefile
│   │   ├── Dockerfile
│   │   ├── cmd/main.go
│   │   └── internal/
│   │       ├── handlers/               # Fiber v3 HTTP handlers
│   │       │   ├── auth_handler.go
│   │       │   ├── settings_handler.go
│   │       │   ├── sync_handler.go
│   │       │   ├── ai_handler.go
│   │       │   ├── mcp_handler.go      # MCP tool HTTP bridge
│   │       │   ├── integration_handler.go
│   │       │   ├── ws_handler.go
│   │       │   └── health_handler.go
│   │       ├── services/
│   │       │   ├── auth_service.go
│   │       │   ├── sync_service.go
│   │       │   ├── ai_service.go
│   │       │   └── export_service.go
│   │       ├── repositories/           # GORM data access
│   │       │   ├── user_repo.go
│   │       │   ├── settings_repo.go
│   │       │   └── sync_repo.go
│   │       ├── middleware/
│   │       │   ├── auth.go
│   │       │   ├── cors.go
│   │       │   ├── ratelimit.go
│   │       │   └── logging.go
│   │       ├── websocket/
│   │       │   ├── server.go
│   │       │   ├── protocol.go
│   │       │   ├── sync_ws.go
│   │       │   ├── ai_bridge_ws.go
│   │       │   └── workflow_ws.go
│   │       ├── integrations/           # External service clients
│   │       │   ├── github/             # OAuth, issues, PRs, credentials
│   │       │   ├── jira/               # OAuth, JQL, issues
│   │       │   ├── linear/             # OAuth, GraphQL, issues
│   │       │   ├── notion/             # OAuth, markdown-to-blocks
│   │       │   ├── figma/              # OAuth PKCE, files/nodes
│   │       │   ├── slack/              # Bot, handlers
│   │       │   └── discord/            # Bot, gateway, handlers
│   │       └── database/migrations/    # PostgreSQL SQL migrations
│   │
│   ├── orch-sync/                       # Background sync daemon
│   │   ├── go.mod
│   │   ├── Makefile
│   │   ├── Dockerfile
│   │   ├── cmd/main.go
│   │   └── internal/
│   │       ├── client.go               # Pull polling (30s, exponential backoff)
│   │       ├── outbox.go               # SQLite outbox store
│   │       ├── outbox_worker.go        # Background flush (50/batch, 10 retries)
│   │       ├── migration.go
│   │       └── device.go
│   │
│   └── orch-gateway/                    # QUIC-to-WebTransport bridge
│       ├── go.mod
│       ├── Makefile
│       ├── Dockerfile
│       ├── cmd/main.go
│       └── internal/
│           ├── webtransport.go         # HTTP/3 WebTransport server
│           ├── proxy.go                # QUIC ↔ WebTransport proxying
│           ├── auth.go                 # mTLS + JWT
│           └── session.go              # Client session management
│
├── apps/                                # ALL NATIVE CLIENTS
│   ├── macos/                           # SwiftUI macOS app
│   │   ├── Orchestra.xcodeproj/
│   │   ├── Orchestra/
│   │   │   ├── OrchestraApp.swift
│   │   │   ├── Views/
│   │   │   ├── Models/
│   │   │   ├── Services/              # QUIC client, gRPC, local storage
│   │   │   ├── Extensions/            # Spotlight, Keychain, iCloud, Notifications
│   │   │   └── Resources/
│   │   ├── OrchestraWidget/            # macOS WidgetKit
│   │   ├── OrchestraIntents/           # Siri Shortcuts
│   │   └── Packages/OrchestraKit/      # Shared Swift package (macOS + iOS)
│   │       └── Sources/OrchestraKit/
│   │           ├── Transport/          # QUIC client (Network.framework)
│   │           ├── Proto/              # gen/swift/ generated code
│   │           └── Storage/            # SQLite + LanceDB local
│   │
│   ├── ios/                             # SwiftUI iOS (shares OrchestraKit)
│   │   ├── Orchestra-iOS.xcodeproj/
│   │   ├── Orchestra-iOS/
│   │   │   ├── OrchestraApp.swift
│   │   │   ├── Views/
│   │   │   ├── Services/
│   │   │   └── Resources/
│   │   └── OrchestraWidget-iOS/
│   │
│   ├── android/                         # Kotlin/Compose
│   │   ├── build.gradle.kts
│   │   ├── app/src/main/kotlin/com/orchestra/
│   │   │   ├── OrchestraApp.kt
│   │   │   ├── ui/                     # Compose screens
│   │   │   ├── data/                   # Room DB, DataStore
│   │   │   ├── domain/                 # Use cases
│   │   │   └── service/               # QUIC client, background sync
│   │   └── lib-orchestra/             # Shared Android library
│   │
│   ├── windows/                         # C#/WinUI3
│   │   ├── Orchestra.sln
│   │   ├── Orchestra/                  # WinUI3 project
│   │   │   ├── App.xaml.cs
│   │   │   ├── Views/
│   │   │   ├── ViewModels/
│   │   │   ├── Services/
│   │   │   └── Assets/
│   │   ├── Orchestra.Core/            # Shared .NET library
│   │   └── Orchestra.Widgets/         # Adaptive Cards widgets
│   │
│   ├── linux/                           # Go + GTK4
│   │   ├── go.mod
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── ui/                     # GTK4 components (gotk4)
│   │   │   ├── services/
│   │   │   └── tray/
│   │   ├── flatpak/
│   │   └── snap/
│   │
│   ├── web/                             # React/TypeScript web dashboard
│   │   ├── package.json                # @orchestra/web
│   │   ├── vite.config.ts
│   │   ├── src/
│   │   │   ├── main.tsx
│   │   │   ├── App.tsx
│   │   │   ├── pages/
│   │   │   ├── components/
│   │   │   ├── stores/
│   │   │   ├── hooks/
│   │   │   └── lib/webtransport.ts    # WebTransport client
│   │   └── public/
│   │
│   └── chrome/                          # Chrome extension
│       ├── package.json                # @orchestra/chrome
│       ├── manifest.json               # Manifest V3
│       ├── vite.config.ts
│       └── src/
│           ├── background/             # Service worker
│           ├── content/                # Content scripts
│           ├── popup/
│           ├── sidepanel/
│           └── lib/webtransport.ts
│
├── deploy/                              # Deployment
│   ├── docker/
│   │   ├── docker-compose.yml          # Dev infra (Postgres, Redis, ClickHouse)
│   │   ├── docker-compose.test.yml     # CI integration tests
│   │   └── docker-compose.prod.yml     # Production
│   ├── k8s/                            # Kubernetes manifests
│   └── terraform/                      # GCP infrastructure
│
├── scripts/
│   ├── proto-gen.sh                    # Run buf generate for all languages
│   ├── install-deps.sh
│   ├── check-ports.sh
│   ├── bundle-macos.sh
│   ├── codesign-macos.sh
│   ├── build-msix.sh
│   └── build-flatpak.sh
│
├── tools/
│   ├── buf/                            # Buf config overrides
│   └── golangci/.golangci.yml          # Shared Go lint config
│
├── docs/
│   ├── architecture/                   # ADRs
│   ├── api/                            # API docs
│   ├── guides/                         # Developer guides
│   └── artifacts/                      # Reference artifacts (this directory)
│
├── Makefile                             # Top-level orchestrator
├── go.work                              # Go workspace (all Go modules)
├── go.work.sum
├── pnpm-workspace.yaml                 # Frontend workspace
├── package.json
├── turbo.json                           # Turborepo config
├── .gitignore
├── .editorconfig
├── CLAUDE.md
├── AGENTS.md
├── CONTEXT.md
├── README.md
└── LICENSE
```

---

## 2. Proto Generation Strategy

### buf.gen.yaml (6 languages)

```yaml
version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/orchestra-agents/gen/go
plugins:
  # Go
  - remote: buf.build/protocolbuffers/go
    out: ../gen/go
    opt: [paths=source_relative]
  - remote: buf.build/grpc/go
    out: ../gen/go
    opt: [paths=source_relative]

  # TypeScript (connect-es for WebTransport)
  - remote: buf.build/connectrpc/es
    out: ../gen/ts
    opt: [target=ts]

  # Swift
  - local: protoc-gen-grpc-swift
    out: ../gen/swift
    opt: [Visibility=Public]
  - local: protoc-gen-swift
    out: ../gen/swift
    opt: [Visibility=Public]

  # Kotlin
  - local: protoc-gen-grpc-kotlin
    out: ../gen/kotlin
  - remote: buf.build/protocolbuffers/java
    out: ../gen/kotlin

  # C#
  - local: protoc-gen-csharp
    out: ../gen/csharp
  - local: protoc-gen-grpc-csharp
    out: ../gen/csharp
```

### Rust (exception — uses tonic-build)

```rust
// services/orch-engine/build.rs
fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::configure()
        .build_server(true)
        .build_client(false)
        .compile_protos(
            &[
                "../../proto/orchestra/engine/v1/health.proto",
                "../../proto/orchestra/engine/v1/parse.proto",
                "../../proto/orchestra/engine/v1/search.proto",
                "../../proto/orchestra/engine/v1/memory.proto",
                "../../proto/orchestra/engine/v1/component.proto",
            ],
            &["../../proto"],
        )?;
    Ok(())
}
```

### How Each Platform Imports Generated Code

| Language | Location | Import Path |
|----------|----------|-------------|
| Go | `gen/go/orchestra/` | `github.com/orchestra-agents/gen/go/orchestra/engine/v1` |
| Rust | Built into `target/` | `include!(concat!(env!("OUT_DIR"), "/orchestra.engine.v1.rs"))` |
| TypeScript | `gen/ts/orchestra/` | `@orchestra/gen` via pnpm workspace |
| Swift | `gen/swift/Orchestra/` | Local package in `OrchestraKit/Package.swift` |
| Kotlin | `gen/kotlin/com/orchestra/` | Gradle source set |
| C# | `gen/csharp/Orchestra/` | Project reference in `.csproj` |

---

## 3. Go Workspace

### go.work

```go
go 1.25

use (
    ./libs/go
    ./gen/go
    ./services/orch-mcp
    ./services/orch-server
    ./services/orch-sync
    ./services/orch-gateway
    ./apps/linux
)
```

Each service has its own `go.mod` and is independently buildable. The workspace resolves local paths during development. CI builds each service independently.

### Import pattern

```go
import (
    "github.com/orchestra-agents/libs/go/types"
    "github.com/orchestra-agents/libs/go/markdown"
    "github.com/orchestra-agents/libs/go/quic"
    "github.com/orchestra-agents/libs/go/helpers"
    enginev1 "github.com/orchestra-agents/gen/go/orchestra/engine/v1"
    mcpv1 "github.com/orchestra-agents/gen/go/orchestra/mcp/v1"
)
```

---

## 4. Top-Level Makefile

```makefile
.PHONY: help install proto build test lint clean dev

VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT   ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE     ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BIN_DIR  := $(CURDIR)/bin

# Setup
install:
	@for svc in services/orch-mcp services/orch-server services/orch-sync services/orch-gateway; do \
		(cd $$svc && go mod download); done
	cd services/orch-engine && cargo fetch
	pnpm install

proto:
	cd proto && buf lint && buf generate

lint-proto:
	cd proto && buf lint && buf breaking --against '.git#branch=main'

# Build (all)
build: build-go build-rust build-frontend
build-go: build-mcp build-server build-sync build-gateway

build-mcp:
	@mkdir -p $(BIN_DIR)
	cd services/orch-mcp && go build -ldflags "-s -w \
		-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)" \
		-o $(BIN_DIR)/orch-mcp ./cmd/

build-engine build-rust:
	cd services/orch-engine && cargo build --release
	@mkdir -p $(BIN_DIR) && cp services/orch-engine/target/release/orchestra-engine $(BIN_DIR)/orch-engine

build-server:
	@mkdir -p $(BIN_DIR)
	cd services/orch-server && go build -ldflags "-s -w -X main.Version=$(VERSION)" -o $(BIN_DIR)/orch-server ./cmd/

build-sync:
	@mkdir -p $(BIN_DIR)
	cd services/orch-sync && go build -o $(BIN_DIR)/orch-sync ./cmd/

build-gateway:
	@mkdir -p $(BIN_DIR)
	cd services/orch-gateway && go build -o $(BIN_DIR)/orch-gateway ./cmd/

build-frontend: build-web build-chrome
build-web:
	cd apps/web && pnpm build
build-chrome:
	cd apps/chrome && pnpm build

# Dev
dev:
	$(MAKE) -j4 dev-engine dev-server dev-web dev-mcp
dev-mcp:
	cd services/orch-mcp && go run ./cmd/ --workspace $(CURDIR)
dev-engine:
	cd services/orch-engine && cargo watch -x run
dev-server:
	cd services/orch-server && air
dev-web:
	cd apps/web && pnpm dev

# Test
test: test-go test-rust test-frontend
test-go:
	@for svc in services/orch-mcp services/orch-server services/orch-sync services/orch-gateway libs/go; do \
		(cd $$svc && go test -race -v ./...); done
test-rust:
	cd services/orch-engine && cargo test
test-frontend:
	pnpm -r test

# Lint
lint: lint-proto
	@for svc in services/orch-mcp services/orch-server services/orch-sync services/orch-gateway libs/go; do \
		(cd $$svc && golangci-lint run ./...); done
	cd services/orch-engine && cargo clippy -- -D warnings
	pnpm -r lint

# Docker (infra only)
docker-up:
	docker compose -f deploy/docker/docker-compose.yml up -d
docker-down:
	docker compose -f deploy/docker/docker-compose.yml down

# Clean
clean:
	rm -rf $(BIN_DIR)
	rm -rf services/orch-engine/target
	pnpm -r clean
```

---

## 5. Docker Compose (Dev Infrastructure)

Only databases run in Docker. Services run natively for hot-reload.

```yaml
# deploy/docker/docker-compose.yml
services:
  postgres:
    image: pgvector/pgvector:pg16
    container_name: orch-postgres
    ports: ["5432:5432"]
    environment:
      POSTGRES_DB: orchestra
      POSTGRES_USER: orchestra
      POSTGRES_PASSWORD: orchestra_dev
    volumes: [postgres_data:/var/lib/postgresql/data]
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U orchestra"]
      interval: 5s

  redis:
    image: redis:7-alpine
    container_name: orch-redis
    ports: ["6379:6379"]
    volumes: [redis_data:/data]
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s

  clickhouse:
    image: clickhouse/clickhouse-server:24
    container_name: orch-clickhouse
    ports: ["8123:8123", "9000:9000"]
    volumes: [clickhouse_data:/var/lib/clickhouse]
    environment:
      CLICKHOUSE_DB: orchestra_analytics

  adminer:
    image: adminer:latest
    container_name: orch-adminer
    ports: ["8080:8080"]
    depends_on: [postgres]

volumes:
  postgres_data:
  redis_data:
  clickhouse_data:
```

---

## 6. Design Rationale

| Decision | Why |
|----------|-----|
| `services/` not `cmd/` | Each service is self-contained with own `cmd/`, `internal/`, `go.mod`. Explicit boundary. |
| `libs/` not `pkg/` | Internal shared code for monorepo. Not meant for external import. |
| `gen/` at repo root | Single output location for all 6 languages. One `buf generate` call. |
| Integrations in `orch-server` | Need HTTP for OAuth callbacks, DB for credentials. Extract to own service if they grow. |
| C library separate | Used by Go (CGo), Rust (FFI), and native apps. CMake build shared by all. |
| `pnpm-workspace.yaml` covers `libs/ts` + `apps` | Shared UI components used by web + chrome without publishing. |
| Go workspace (`go.work`) | Each service has locked deps, workspace resolves locally. CI builds independently. |

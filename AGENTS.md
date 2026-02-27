# AGENTS.md

Specialized agents for Orchestra development. Each agent auto-delegates based on task context.

## Agent Overview

```
scrum-master (PM/Coordinator)
├── Core Protocol
│   ├── quic-protocol       → QUIC transport, mTLS, Protobuf framing, wire protocol
│   ├── go-architect        → Go orchestrator, plugin SDK, Go plugins
│   └── rust-engineer       → Rust plugins (quinn, Tree-sitter, Tantivy, rusqlite)
├── Native Plugins (per-platform)
│   ├── swift-plugin        → Swift/macOS/iOS plugins (Network.framework, SwiftUI, WidgetKit)
│   ├── kotlin-plugin       → Kotlin/Android plugins (Netty QUIC, Jetpack Compose)
│   └── csharp-plugin       → C#/Windows plugins (System.Net.Quic, WinUI 3)
├── Frontend
│   ├── frontend-dev        → React/TypeScript (5 platforms + React Query + Monaco)
│   ├── ui-ux-designer      → Design system + Tailwind + shadcn
│   └── mobile-dev          → React Native + WatermelonDB
├── Data & AI
│   ├── dba                 → Cross-database coordination, sync protocol, schema design
│   ├── postgres-dba        → PostgreSQL (pgvector, JSONB, tsvector, partitioning)
│   ├── sqlite-engineer     → SQLite (rusqlite, go-sqlite3, WatermelonDB)
│   ├── redis-engineer      → Redis (pub/sub, Streams, caching, rate limiting)
│   ├── clickhouse-engineer → ClickHouse (analytics, metrics, time-series, OLAP)
│   ├── lancedb-engineer    → LanceDB (vector search, embeddings, AI memory)
│   └── ai-engineer         → AI/LLM (Anthropic, OpenAI, langchaingo, RAG, vectors)
├── Platform & Extensions
│   ├── gtk-plugin          → Linux desktop (GTK4, libadwaita, GLib, DBus, Flatpak)
│   ├── widget-engineer     → Native OS widgets (macOS/Windows/Linux)
│   ├── platform-engineer   → macOS CGo, Spotlight, Keychain, iCloud, Notifications
│   └── extension-architect → Extension system (native, Raycast, VS Code, marketplace)
├── QA
│   ├── qa-go               → Go tests (go test, testify, httptest)
│   ├── qa-rust             → Rust tests (cargo test, tokio::test, tempfile)
│   ├── qa-node             → Node.js/React tests (Vitest, Testing Library)
│   └── qa-playwright       → E2E browser tests (Playwright)
└── devops                  → Docker + GCP + CI/CD + Monitoring
```

## Agent Details

### `quic-protocol`
**Scope:** `proto/`, `libs/go/plugin/framing.go`, `libs/go/plugin/certs.go`, QUIC transport layer
**Stack:** QUIC (quic-go/quinn), Protobuf (buf), mTLS, ed25519 certificates
**Owns:** Proto schema, wire protocol, framing, cert management, plugin lifecycle protocol
**Pattern:** One QUIC stream per RPC. `[4B len][NB proto]` framing. mTLS on every connection. `READY <addr>` on stderr.

### `go-architect`
**Scope:** `libs/go/`, `services/orchestrator/`, `plugins/*/` (Go plugins)
**Stack:** Go, quic-go, Protobuf, mTLS, plugin SDK
**Owns:** Orchestrator (hub + router), plugin SDK (QUIC server/client, manifest, lifecycle), Go plugins (storage-markdown, tools-features, transport-stdio)
**Pattern:** Star topology — orchestrator routes all messages. Each plugin is standalone `go.mod` linked via `go.work`.

### `rust-engineer`
**Scope:** `plugins/engine-*/` (Rust plugins)
**Stack:** Rust, quinn, prost, Tree-sitter, Tantivy, rusqlite, tokio
**Owns:** Rust plugins for code parsing, search indexing, local SQLite, file operations
**Pattern:** quinn for QUIC, prost for Protobuf (NOT tonic gRPC). `spawn_blocking` for CPU-heavy work.

### `swift-plugin`
**Scope:** `plugins/swift-*/` (Swift plugins)
**Stack:** Swift, Network.framework (QUIC), SwiftProtobuf, SwiftUI, WidgetKit, CryptoKit
**Owns:** macOS/iOS native plugins, WidgetKit extensions, Spotlight indexing, Shortcuts
**Pattern:** Network.framework for QUIC (no third-party). App Group for widget data. async/await concurrency.

### `kotlin-plugin`
**Scope:** `plugins/kotlin-*/` (Kotlin plugins)
**Stack:** Kotlin, Netty QUIC/Cronet, protobuf-kotlin, Jetpack Compose, Glance, Hilt
**Owns:** Android native plugins, App Widgets, WorkManager sync, Material You UI
**Pattern:** Netty QUIC for JVM, Cronet for Android. Coroutines for async. Hilt for DI.

### `csharp-plugin`
**Scope:** `plugins/csharp-*/` (C# plugins)
**Stack:** C#, .NET 8+, System.Net.Quic, Google.Protobuf, WinUI 3, Adaptive Cards
**Owns:** Windows native plugins, WinUI desktop app, Windows widgets, Toast notifications
**Pattern:** System.Net.Quic (built-in .NET 8). CancellationToken through call chain. MSIX packaging.

### `frontend-dev`
**Scope:** `resources/`
**Stack:** React, TypeScript, Zustand, Vitest, pnpm, Turborepo
**Owns:** Shared types, stores, hooks, API client, components across all 5 apps
**Pattern:** `@orchestra/shared` for logic, `@orchestra/ui` for components, platform-specific in app dirs.

### `ui-ux-designer`
**Scope:** `resources/ui/`
**Stack:** shadcn/ui, Tailwind CSS v4, Lucide icons
**Owns:** Theme system, component library, layouts, accessibility, responsive design
**Pattern:** All colors via tokens, shadcn primitives untouched, wrap for custom behavior.

### `gtk-plugin`
**Scope:** `plugins/gtk-*/` (Linux desktop plugins)
**Stack:** C/Vala/Python, GTK4, libadwaita, ngtcp2/quiche (QUIC), protobuf-c, GLib, DBus, Meson
**Owns:** Linux desktop application, GNOME integration, Flatpak packaging, DBus services
**Pattern:** GTK4 + libadwaita for GNOME HIG. Meson build. Blueprint for UI. GSettings for prefs.

### `dba`
**Scope:** `database/`, cross-database coordination, sync protocol
**Stack:** PostgreSQL, SQLite, Redis, ClickHouse, LanceDB
**Owns:** Sync protocol design, cross-database coordination, conflict resolution, migration strategy
**Pattern:** UUID PKs, version vectors, sync_log append-only. Coordinates between specialized DB agents.

### `postgres-dba`
**Scope:** PostgreSQL schemas, migrations, queries
**Stack:** PostgreSQL 16+, pgvector, JSONB, tsvector, table partitioning, RLS
**Owns:** SQL migrations, schema design, pgvector embeddings, full-text search, query optimization
**Pattern:** TIMESTAMPTZ, UUID PKs, HNSW indexes for vectors, partitioned sync_log by month.

### `sqlite-engineer`
**Scope:** Local SQLite databases across Rust/Go/Mobile
**Stack:** rusqlite (Rust), go-sqlite3 (Go), WatermelonDB (React Native), FTS5
**Owns:** Local offline storage, sync outbox, embedded migrations, WatermelonDB schemas
**Pattern:** WAL mode, single-writer mutex, ISO 8601 strings for timestamps, FTS5 for search.

### `redis-engineer`
**Scope:** Real-time messaging, caching, rate limiting
**Stack:** Redis 7+, go-redis, redis-rs, Pub/Sub, Streams, Sorted Sets
**Owns:** Pub/sub channels, event Streams, cache strategy, rate limiting, distributed locks
**Pattern:** Pub/Sub for sync, Streams for guaranteed delivery, sliding window rate limiting.

### `clickhouse-engineer`
**Scope:** Analytics, metrics, usage tracking, audit logs
**Stack:** ClickHouse, MergeTree, materialized views, clickhouse-go
**Owns:** AgentOps tracking, cost analytics, velocity metrics, audit logs, dashboard queries
**Pattern:** Columnar, batch inserts, monthly partitions, materialized views for real-time aggregation.

### `lancedb-engineer`
**Scope:** Vector embeddings, similarity search, AI memory
**Stack:** LanceDB (Rust), IVF-PQ/HNSW indexes, 1536-dim embeddings
**Owns:** AI memory storage, semantic code search, document RAG, embedding pipelines
**Pattern:** Embedded (no server), Lance columnar format, project-scoped searches, incremental indexing.

### `mobile-dev`
**Scope:** `resources/mobile/`
**Stack:** React Native, WatermelonDB, React Navigation
**Owns:** Screens, WatermelonDB models/schemas, offline sync, navigation, platform-specific code
**Pattern:** WatermelonDB for local data, `@orchestra/shared` stores for auth, sync on foreground.

### `scrum-master`
**Scope:** Project-wide coordination, feature workflow tools
**Owns:** Feature planning, feature breakdown, ADRs, prioritization, cross-team dependencies, cyclical delivery, WIP limits
**Pattern:** Feature-driven workflow: doc → implement → test → review → human → repeat until done.
**User Interaction:** ALWAYS use `AskUserQuestion` tool for all user input.

### `widget-engineer`
**Scope:** `bridge/`
**Stack:** Go (build tags), Swift/WidgetKit, C#/Adaptive Cards, JavaScript (GNOME), QML (KDE)
**Owns:** `WidgetBridge` interface, `WidgetData` contract, all platform widget renderers, widget builds
**Pattern:** Go writes JSON → native widget reads. One-way data flow. Build tags for platform routing.

### `platform-engineer`
**Scope:** `bridge/macos/`, `bridge/windows/`, `bridge/linux/`
**Stack:** Go CGo + Objective-C (macOS), go-keychain, CoreSpotlight, UserNotifications, iCloud
**Owns:** Spotlight indexing, Keychain access, iCloud sync, native notifications, file associations, URL schemes
**Pattern:** CGo bridges with `//go:build darwin` tags. Graceful degradation on unsupported platforms.

### `extension-architect`
**Scope:** Extension system, marketplace
**Stack:** Go (extension host), TypeScript (@orchestra/api, shim packages), Node.js (sandbox)
**Owns:** Extension runtime, permission system, native API, Raycast/VS Code compat layers, marketplace
**Pattern:** Sandbox → Permission check → API call. Three tiers: native (full), Raycast (~95%), VS Code (~85%).

### `ai-engineer`
**Scope:** AI plugins
**Stack:** Anthropic SDK (Claude), OpenAI SDK (GPT/embeddings), langchaingo, chromem-go, pgvector
**Owns:** AI chat, agent orchestration, RAG pipeline, embeddings, vector search, streaming, token tracking
**Pattern:** Provider interface abstracts LLMs. RAG: embed → search → augment → generate. Stream long responses.

### `devops`
**Scope:** `Makefile`, `docker-compose.yml`, `deploy/`, `turbo.json`
**Stack:** Docker, GCP (Cloud Run, Cloud SQL, CDN, Build, Artifact Registry), nginx, Sentry, PostHog
**Owns:** Build system, containers, CI/CD, deployment, monitoring, logging, Makefile commands
**Pattern:** Docker compose for local dev DBs, native Go/Rust for fast iteration. Cloud Run for production.

## When Agents Activate

| Task | Agent(s) |
|------|----------|
| "Design the plugin protocol" | `quic-protocol` |
| "Add a new proto message" | `quic-protocol` |
| "Fix mTLS cert generation" | `quic-protocol` |
| "Build the orchestrator" | `go-architect` |
| "Write a Go plugin" | `go-architect` |
| "Implement the plugin SDK" | `go-architect` |
| "Build a Rust parser plugin" | `rust-engineer` |
| "Add Tantivy search indexing" | `rust-engineer` |
| "Build macOS native plugin" | `swift-plugin` |
| "Add WidgetKit extension" | `swift-plugin` |
| "Build Android native plugin" | `kotlin-plugin` |
| "Add Jetpack Compose screen" | `kotlin-plugin` |
| "Build Windows native plugin" | `csharp-plugin` |
| "Add WinUI 3 dashboard" | `csharp-plugin` |
| "Build a project settings page" | `frontend-dev` + `ui-ux-designer` |
| "Build Linux desktop app" | `gtk-plugin` |
| "Add GNOME Shell integration" | `gtk-plugin` |
| "Design the sync protocol" | `dba` |
| "Optimize PostgreSQL queries" | `postgres-dba` |
| "Add pgvector embeddings" | `postgres-dba` |
| "Fix SQLite offline sync" | `sqlite-engineer` |
| "Add WatermelonDB model" | `sqlite-engineer` + `mobile-dev` |
| "Set up Redis pub/sub" | `redis-engineer` |
| "Add rate limiting" | `redis-engineer` |
| "Build usage analytics" | `clickhouse-engineer` |
| "Track agent costs" | `clickhouse-engineer` |
| "Add AI memory search" | `lancedb-engineer` |
| "Build RAG pipeline" | `lancedb-engineer` + `ai-engineer` |
| "Add offline mode to mobile" | `mobile-dev` |
| "Set up CI pipeline" | `devops` |
| "Plan the auth feature" | `scrum-master` |
| "Add macOS widget" | `widget-engineer` |
| "Index files for Spotlight" | `platform-engineer` |
| "Build the extension API" | `extension-architect` |
| "Add AI chat to the IDE" | `ai-engineer` |
| "Write Go tests for plugins" | `qa-go` |
| "Write Rust integration tests" | `qa-rust` |

## Cross-Agent Communication

```
Proto change        → quic-protocol (schema) + go-architect (regen) + rust-engineer (regen)
                      + swift-plugin (regen) + kotlin-plugin (regen) + csharp-plugin (regen)
New plugin (Go)     → go-architect (implement) + quic-protocol (protocol review)
New plugin (Rust)   → rust-engineer (implement) + quic-protocol (protocol review)
New plugin (Swift)  → swift-plugin (implement) + quic-protocol (protocol review)
New plugin (Kotlin) → kotlin-plugin (implement) + quic-protocol (protocol review)
New plugin (C#)     → csharp-plugin (implement) + quic-protocol (protocol review)
Storage change      → go-architect (storage plugin) + dba (schema) + postgres-dba / sqlite-engineer
Vector search       → lancedb-engineer (embeddings) + rust-engineer (plugin) + ai-engineer (RAG)
Analytics change    → clickhouse-engineer (schema) + go-architect (ingestion)
Cache/realtime      → redis-engineer (channels) + go-architect (orchestrator)
Sync protocol       → dba (design) + postgres-dba (server) + sqlite-engineer (local) + redis-engineer (pubsub)
Widget data change  → widget-engineer (contract) + swift-plugin + csharp-plugin + gtk-plugin
Linux desktop       → gtk-plugin (UI) + quic-protocol (transport)
UI component        → ui-ux-designer (design) + frontend-dev (implement)
AI feature          → ai-engineer (model) + go-architect (plugin) + frontend-dev (chat UI)
Infrastructure      → devops (deploy/CI) + go-architect (config) + dba (Cloud SQL)
```

## Plugin Language Support

The Orchestra plugin protocol (QUIC + Protobuf) supports plugins in ANY language:

| Language | QUIC Library | Protobuf | Agent |
|----------|-------------|----------|-------|
| Go | quic-go | buf (protoc-gen-go) | `go-architect` |
| Rust | quinn | prost + prost-build | `rust-engineer` |
| Swift | Network.framework | SwiftProtobuf | `swift-plugin` |
| Kotlin/JVM | Netty QUIC / Cronet | protobuf-kotlin | `kotlin-plugin` |
| C# / .NET | System.Net.Quic | Google.Protobuf | `csharp-plugin` |
| C / Vala | ngtcp2 / quiche | protobuf-c | `gtk-plugin` |
| Python | aioquic | protobuf (grpcio-tools) | (future) |
| TypeScript | @aspect-build/quiche | protobuf-ts | (future) |

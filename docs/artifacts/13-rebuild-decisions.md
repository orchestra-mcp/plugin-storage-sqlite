# Rebuild Decisions — Orchestra

> Every architectural decision made during the planning session for the ground-up rebuild.
> This is the single source of truth for "why we chose X over Y."

---

## 1. The Final Stack

```
═══════════════════════════════════════════════════════
              ORCHESTRA MCP — FINAL STACK
═══════════════════════════════════════════════════════

CORE PHILOSOPHY:
  Agents speak Markdown. Orchestra MCP translates to everything.
  MCP is the universal adapter between AI and infrastructure.
  Doc-first, TDD-driven, feature-driven workflow. Cyclical until human approves.
  The doc IS the spec. The doc IS the memory. The doc IS the deliverable.

LANGUAGES:
  Go              → server, agent, MCP, integrations, Linux app
  Rust            → data engine, RAG, artifacts
  C               → PTY, IO, crypto (~1000 lines)
  Swift/SwiftUI   → macOS + iOS
  Kotlin/Compose  → Android
  C# / WinUI 3    → Windows
  Go + GTK4       → Linux
  React/TS        → Web + Chrome extension

PROTOCOL:
  QUIC            → everything
  WebTransport    → browsers
  Protobuf        → message format
  mTLS            → auth
  MCP             → AI ↔ Orchestra (Markdown in, Markdown out)

DATABASE:
  PostgreSQL + pgvector → server
  Redis + Streams       → server
  ClickHouse            → web analytics
  SQLite                → local
  LanceDB               → local vectors

WEB SERVER:
  Caddy                 → reverse proxy, auto TLS, HTTP/3
═══════════════════════════════════════════════════════
```

---

## 2. Architecture Decisions

### 2.1 Microservices, Event-Driven

**Decision**: Rebuild as standalone microservices communicating via events over QUIC.

**Why**: The reference codebase (`orch-ref/`) grew into a tightly-coupled monolith. The `bootstrap/desktop.go` file alone was 1188 lines wiring ~20 subsystems. A crash in any subsystem took down everything. Services couldn't be tested independently. Port conflicts required workaround hacks (dual-port proxy, lazyOpener polling).

**What this means**: Each service is a standalone binary with its own `go.mod` or `Cargo.toml`. Services discover each other, connect directly over QUIC, and exchange Protobuf events. No service imports another service's code.

### 2.2 QUIC Everywhere (No HTTP/REST, No WebSocket)

**Decision**: QUIC as the universal transport. WebTransport for browsers.

**Why**: The reference used HTTP/REST (Fiber v3, 134 routes) + WebSocket (18 events) + gRPC (5 services) + stdio (MCP). Four different protocols for what is essentially the same thing: sending typed messages between services. QUIC unifies this:
- Multiplexed streams (no head-of-line blocking)
- 0-RTT reconnect (fast service restarts)
- Built-in TLS 1.3 (mTLS natural)
- Works for both local IPC and network

**Exception**: MCP stdio stays JSON-RPC per the MCP specification (agents expect JSON over stdin/stdout).

### 2.3 Peer-to-Peer Mesh (No Central Broker)

**Decision**: Services connect directly to each other. No NATS, no Redis Streams for event routing.

**Why**: With 5 services on a developer's laptop, a central broker adds complexity and a single point of failure for no benefit. 5 choose 2 = 10 connections max. Each connection is a single UDP socket. The mesh library handles subscription matching at the application layer.

**Trade-off**: If we scale to 50+ services in the future, we'd need a broker. For v1 with 5 services, direct mesh is simpler.

### 2.4 Service Discovery: File + mDNS

**Decision**: Seed file (`~/.orchestra/mesh.json`) for same-machine. mDNS/Bonjour for LAN.

**Why**:
- File-based is fast and simple for the common case (all services on one machine)
- mDNS enables team development (services on different machines on the same network)
- Both work without any infrastructure (no consul, no etcd)

### 2.5 Protobuf Everywhere

**Decision**: All inter-service messages are Protobuf-encoded. Proto files are the contract.

**Why**:
- Typed and versioned (breaking changes caught by `buf breaking`)
- Fast serialization/deserialization
- Generated code for all 6 languages from one source
- Self-documenting via `.proto` files

### 2.6 Protobuf + Markdown Replaces TOON

**Decision**: Drop the TOON file format (YAML frontmatter + Markdown body). Use Protobuf for structured metadata and plain Markdown for content.

**Why**: We already have Protobuf as the universal message format. Having a second serialization format (YAML) adds complexity with no benefit. Protobuf handles structured fields (status, priority, assignee, dates, labels) and Markdown handles the human-readable content body. The storage plugin reads/writes `.md` files with a companion `.pb` sidecar (or a single file with a binary header). Files remain git-trackable and human-readable since the Markdown body is plain text.

**What this means**: No `libs/go/toon/` package. The `storage.markdown` plugin handles file I/O. Structured fields are Protobuf-encoded, content is Markdown. The `StorageReadResponse.metadata` map carries decoded Protobuf fields as key-value pairs for cross-plugin use.

### 2.7 mTLS Between All Services

**Decision**: Every service authenticates to every other service via mutual TLS.

**Why**: Even on localhost, mTLS prevents rogue processes from joining the mesh. The local CA is auto-generated on first run. No manual certificate management needed.

---

## 3. What Changes from the Reference

| Aspect | Reference (orch-ref) | Rebuild |
|--------|---------------------|-------------|
| Architecture | Monolith (1 binary, 20 subsystems) | 5 standalone microservices |
| Transport | HTTP/REST + WebSocket + gRPC + stdio | QUIC + WebTransport + stdio |
| Message format | JSON everywhere (except gRPC) | Protobuf everywhere (except MCP stdio) |
| Service communication | In-process function calls | Event-driven over QUIC mesh |
| Auth between services | None (same process) | mTLS |
| Discovery | Port files + polling | Seed file + mDNS |
| Desktop UI | Wails v3 (Go + React webview) | SwiftUI (native macOS) |
| Mobile | React Native | Swift (iOS) + Kotlin (Android) |
| Windows | Not supported | C#/WinUI 3 |
| Linux | Not supported | Go + GTK4 |
| Web | React/TS (in Wails webview) | React/TS (standalone, WebTransport) |
| Local vectors | Custom Rust memory store | LanceDB |
| Cloud DB | Laravel (PHP) | Go server + PostgreSQL |
| Analytics | None | ClickHouse |
| Reverse proxy | Not used | Caddy (HTTP/3, auto TLS) |
| Build system | Single Makefile | Per-service Makefile + top-level orchestrator |
| Go modules | Single go.mod | go.work workspace with per-service modules |

---

## 4. What We Keep from the Reference

These concepts are proven and carry over (logic rewritten, not code copied):

| What | Why It's Proven |
|------|----------------|
| MCP tools (feature-driven) | ~58 feature tools (CRUD, workflow, agentops, time, deps, git, notifications, quality, assignment) + memory + inception |
| Gated workflow machine | backlog → todo → in-progress → ... → human-review → done with evidence gates |
| Protobuf + Markdown storage | Structured data in Protobuf, content body in Markdown, human-readable, git-trackable |
| Cyclical delivery loop | doc → implement → test → review → human → update → repeat until approved |
| PRD system (adapted) | Guided inception interview, feature discovery with reasoning |
| Plugin architecture | Load/Boot/Register/Shutdown lifecycle |
| WIP limits + dependencies | Constraint enforcement for task management |
| Tree-sitter parsing (14 languages) | Rust engine code intelligence |
| Tantivy full-text search | Rust engine search indexing |
| 10 integrations (GitHub, Jira, Linear, etc.) | OAuth flows, API clients, sync patterns |

---

## 5. What We Delete (Not Carried Over)

| What | Why |
|------|-----|
| `bootstrap/desktop.go` (1188 lines) | The monolith wiring. Replaced by independent service mains. |
| `bootstrap/app.go` | Plugin bootstrap. Each binary does its own init. |
| Wails v3 desktop shell | Replaced by native SwiftUI/Kotlin/WinUI3/GTK4 apps. |
| All 19 React packages (267 components) | Web/Chrome keep React. Desktop uses native UI. |
| `lazyOpener` pattern | No longer needed when services are independent. |
| Dual-port proxy (19191/19192) | No port conflicts with OS-assigned ports + mDNS. |
| Fiber v3 HTTP REST server | Replaced by QUIC. Web access via orch-gateway. |
| Gorilla WebSocket server | Replaced by QUIC streams / WebTransport. |
| Laravel sync backend (PHP) | Replaced by Go `orch-server` + PostgreSQL. |
| `app/types/*.go` (23 files) | Rewritten in `libs/go/types/` with cleaner structure. |
| `app/tools/*.go` (87 files) | Rewritten in `services/orch-mcp/internal/tools/` from scratch. |
| TOON file format (`libs/go/toon/`) | Replaced by Protobuf metadata + Markdown body. No YAML frontmatter. |

---

## 6. Build Order

### Phase 0: Foundation
1. Set up monorepo structure (proto/, libs/, services/, apps/)
2. Define all Protobuf schemas (mesh envelope, MCP types, engine services)
3. Generate code for all 6 languages
4. Build `libs/go/quic/` mesh library
5. Build `libs/go/types/` and `libs/go/markdown/` shared packages

### Phase 1: MCP Server (First Service)
6. Build `orch-mcp` with feature-driven tools (~58 tools first):
   - Project CRUD (4 tools)
   - Feature CRUD + context-aware sizing (6 tools)
   - Workflow gates with evidence (8 tools)
   - Dependencies + WIP limits (7 tools)
   - AgentOps — cost & token tracking (4 tools)
   - Time management — estimation & cycle time (3 tools)
   - Git workflow — branches, merge (3 tools)
   - Notifications — human alerts (3 tools)
   - Conflict resolution — file locks, overlap detection (4 tools)
   - Automated quality gates — lint, test, coverage, security (4 tools)
   - Reporting (3 tools)
   - Assignment & metadata + notes (7 tools + 3 notes)
7. MCP stdio transport (JSON-RPC)
8. Protobuf + Markdown file storage
9. QUIC mesh participant (publish events)

### Phase 2: Web Client (MCP over HTTP)
10. Build `orch-gateway` (WebTransport bridge)
11. Build `apps/web` React dashboard
12. MCP tools accessible from browser via WebTransport → gateway → mesh → orch-mcp

### Phase 3: Rust Engine
13. Build `orch-engine` with Tree-sitter, Tantivy, LanceDB
14. QUIC mesh participant
15. gRPC services (Parse, Search, Memory)

### Phase 4: API Server
16. Build `orch-server` (HTTP/3, auth, settings, AI bridge, integrations)
17. PostgreSQL migrations
18. OAuth flows for GitHub, Jira, Linear, Notion, Figma

### Phase 5: macOS App (First Native Client)
19. Build `apps/macos` SwiftUI app
20. QUIC client via Network.framework
21. AI Chat screen
22. Task board
23. Notes editor
24. Settings

### Phase 6: More Platforms
25. `apps/ios` (shares OrchestraKit with macOS)
26. `apps/android` (Kotlin/Compose)
27. `apps/windows` (C#/WinUI 3)
28. `apps/linux` (Go + GTK4)
29. `apps/chrome` (Chrome extension)

### Phase 7: Sync + Integrations
30. Build `orch-sync` daemon
31. Cloud sync with PostgreSQL
32. Discord bot
33. Slack bot
34. Remaining integrations

---

## 7. Key Technical Choices

### QUIC Libraries

| Language | Library | Notes |
|----------|---------|-------|
| Go | `quic-go/quic-go` | Production-ready, RFC 9000 compliant |
| Rust | `quinn` | 86M+ downloads, pairs with `tonic-h3` |
| Swift | `Network.framework` | Apple's native QUIC, built into iOS/macOS |
| Kotlin | `cloudflare/quiche` | Via Android NDK bindings |
| C# | `System.Net.Quic` | .NET 7+ native QUIC |
| TypeScript | WebTransport API | Browser-native, Chrome + Firefox |

### mDNS Library

`grandcat/zeroconf` — Full RFC 6762/6763, better IPv6 than `hashicorp/mdns`, channel-based API.

### Local Vector DB

LanceDB — Purpose-built for embeddings, columnar storage, Rust-native. Replaces custom vector store in Rust engine.

### Proto Generation

Buf (`buf.build`) for Go, TS, Swift, Kotlin, C#. Rust uses `tonic-build` in `build.rs` (Buf doesn't support Rust codegen well).

---

## 8. What to Build First: Feature-Driven Tools (~58 Tools)

The first set of tools to implement in `orch-mcp`, proving the architecture end-to-end:

| Category | Tools | Count |
|----------|-------|-------|
| Project CRUD | create_project, list_projects, delete_project, get_project_status | 4 |
| Feature CRUD | create_feature, get_feature, update_feature, list_features, delete_feature, search | 6 |
| Workflow | get_next_feature, set_current_feature, advance_feature, reject_feature, get_workflow_status | 5 |
| Review (Human Cycle) | request_review, submit_review, get_pending_reviews | 3 |
| Dependencies + WIP | add_dependency, remove_dependency, get_dependency_graph, set_wip_limits, get_wip_limits, check_wip_limit, split_feature | 7 |
| AgentOps | record_usage, get_usage, get_project_cost, set_budget | 4 |
| Time Management | set_estimate, get_estimate_accuracy, get_cycle_time | 3 |
| Git | create_branch, get_branch_status, merge_feature | 3 |
| Notifications | send_notification, get_notification_config, set_notification_config | 3 |
| Conflicts | lock_files, unlock_files, check_conflicts, resolve_conflict | 4 |
| Quality Gates | run_quality_check, get_quality_config, set_quality_config, get_quality_report | 4 |
| Reporting | get_progress, get_blocked_features, get_review_queue | 3 |
| Assignment + Metadata | assign_feature, unassign_feature, add_labels, remove_labels, save_note, list_notes, search_notes | 7 |
| **Total** | | **~56** |

---

## 9. Artifact Index

All artifacts in `docs/artifacts/`:

| # | Artifact | Lines | What It Covers |
|---|---------|-------|---------------|
| 00 | `00-architecture-overview.md` | 444 | Reference system diagram, monolith problems, original decoupling proposal |
| 01 | `01-service-catalog.md` | 707 | 7 reference services: protocols, ports, startup |
| 02 | `02-data-model.md` | 2,510 | Every entity type: TOON, SQLite, PostgreSQL, proto |
| 03 | `03-api-surface.md` | 739 | 134 REST routes, 18 WS events, 5 gRPC services |
| 04 | `04-mcp-tool-catalog.md` | 2,534 | 186 MCP tools, 34 categories, full param tables |
| 05 | `05-integration-map.md` | 860 | 10 integrations: auth, APIs, sync |
| 06 | `06-communication-patterns.md` | 998 | 12 inter-service patterns |
| 07 | `07-desktop-subsystems.md` | 1,212 | 23 desktop subsystems |
| 08 | `08-plugin-system.md` | 1,167 | Plugin lifecycle, 8 capabilities |
| 09 | `09-ui-inventory.md` | 1,174 | 19 packages, ~267 components |
| 10 | `10-feature-matrix.md` | 586 | Feature-to-service mapping, build priority |
| **11** | **`11-quic-mesh-architecture.md`** | **~350** | **QUIC mesh design, Protobuf envelope, event flow** |
| **12** | **`12-monorepo-structure.md`** | **~500** | **Directory tree, proto gen, Makefile, Docker** |
| **13** | **`13-rebuild-decisions.md`** | **this file** | **All decisions, what changes, build order** |

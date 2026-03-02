# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Orchestra MCP is an AI-agentic IDE targeting 5 platforms: Desktop (Wails), Chrome Extension, Mobile iOS, Mobile Android, and Web Dashboard. Built with Go (Fiber v3 + GORM), Rust (Tonic gRPC + Tree-sitter + Tantivy), and React/TypeScript (pnpm + Turborepo + Zustand).

The old Laravel codebase is preserved at `old-ref/` for reference. All new development happens in the repo root.

## Key Commands

```bash
# Development
make dev                 # Start Go backend + Rust engine + all frontends
make dev-go              # Go server only (air hot-reload)
make dev-rust            # Rust engine only (cargo watch)
make dev-frontend        # All frontends via Turborepo

# Build
make build               # Build everything (Go + Rust + MCP + frontends)
make build-go            # Go server binary -> bin/server
make build-rust          # Rust engine binary
make build-mcp           # MCP plugin binary -> bin/orchestra-mcp
make build-frontend      # All frontend apps

# MCP Plugin
make mcp-build           # Build MCP plugin binary
make mcp-init            # Initialize MCP in current project
make mcp-start           # Start MCP stdio server

# Install & Test
make install             # Install all deps (Go + Rust + pnpm)
make test                # All tests (Go + MCP + Rust + Frontend)
make clean               # Remove build artifacts

# Proto (generate Go + Rust + TypeScript from .proto files)
make proto

# Add a shadcn component
cd resources/ui && npx shadcn@latest add {component}
```

## Project Structure

```
orchestra-mcp/
├── app/                      # Go backend (Fiber + GORM)
│   ├── handlers/             # HTTP handlers (controllers)
│   ├── models/               # GORM models
│   ├── services/             # Business logic
│   ├── repositories/         # Data access
│   ├── middleware/            # Fiber middleware
│   ├── routes/               # Route registration
│   ├── plugins/              # Plugin runtime (8 files — the foundation)
│   │   ├── contracts.go      # Plugin interface + 15 Has* capability interfaces
│   │   ├── manager.go        # PluginManager with topological sort
│   │   ├── context.go        # PluginContext with DI
│   │   ├── manifest.go       # PluginManifest
│   │   ├── registry.go       # ServiceRegistry (thread-safe DI)
│   │   ├── contributes.go    # ContributesRegistry
│   │   ├── features.go       # FeatureManager (feature flags)
│   │   └── loader.go         # PluginLoader (auto-discovery)
│   └── gen/proto/            # Generated protobuf Go code
├── config/                   # Go configuration
│   └── plugins.go            # Plugin registry config
├── plugins/                  # ALL PLUGINS (each is standalone)
│   └── mcp/                  # MCP Plugin — first plugin (85 tools)
│       ├── go.mod            # Standalone module
│       ├── config/mcp.go     # McpConfig
│       ├── providers/        # Plugin registration (bridges to app/plugins)
│       ├── src/
│       │   ├── cmd/main.go   # CLI entry -> orchestra-mcp binary
│       │   ├── types/        # Type definitions (5 files)
│       │   ├── toon/         # TOON/YAML file parser
│       │   ├── workflow/     # State machine transitions
│       │   ├── helpers/      # Shared utilities (5 files)
│       │   ├── transport/    # MCP stdio JSON-RPC server
│       │   ├── tools/        # All 85 MCP tools
│       │   └── bootstrap/    # Workspace init command
│       └── resources/        # Bundled skills + agents
├── cmd/server/main.go        # Go HTTP server entry point
├── engine/                   # Rust engine (gRPC)
│   └── src/                  # Tree-sitter, Tantivy, tower-lsp, rusqlite
├── proto/                    # Shared protobuf definitions
├── database/migrations/      # PostgreSQL SQL migrations
├── resources/                # All frontends (pnpm monorepo)
│   ├── shared/               # @orchestra/shared
│   ├── ui/                   # @orchestra/ui (shadcn/ui)
│   ├── extension/            # Chrome Extension
│   ├── dashboard/            # Web Dashboard
│   ├── desktop/              # Wails Desktop UI
│   └── mobile/               # React Native
├── old-ref/                  # Old Laravel codebase (reference only)
├── Makefile                  # Central command runner
├── go.mod                    # Root Go module
└── pnpm-workspace.yaml       # Frontend workspace config
```

## Architecture

### Plugin System (Component-First)

Everything is a plugin. The plugin runtime at `app/plugins/` provides:
- **Plugin interface** with capability contracts (Has* interfaces)
- **PluginManager** with topological dependency sort and boot sequence
- **FeatureManager** for runtime feature flags — disable any feature by turning off its plugin
- **ServiceRegistry** for plugin-scoped dependency injection
- **ContributesRegistry** for VS Code-style contributions (commands, menus, settings)

Each plugin is a standalone Go module with its own `go.mod`, pushable as a separate GitHub repo. Plugin folder convention: `config/`, `providers/`, `src/`, `resources/`, `README.md`.

### MCP Plugin (Pure Go, 85 tools)

The first plugin at `plugins/mcp/`. Provides project management tools via MCP protocol:
- **Build**: `cd plugins/mcp && go build -o orchestra-mcp ./src/cmd/`
- **Run**: `orchestra-mcp --workspace .` (stdio JSON-RPC)
- **Init**: `orchestra-mcp init --workspace .` (creates .mcp.json, .projects/, .claude/, CLAUDE.md, AGENTS.md, CONTEXT.md)
- **Packages**: `types/`, `toon/`, `workflow/`, `helpers/`, `transport/`, `tools/`, `engine/`, `bootstrap/`
- **Workflow**: 13-state lifecycle (backlog → todo → in-progress → ready-for-testing → in-testing → ready-for-docs → in-docs → documented → in-review → done)
- **Multi-Audience PRD**: 4 audience types (business/product/technical/qa) with conditional follow-up questions, validation, agent briefings, auto-backlog generation, and reusable templates
- **Sprint Management**: Create/start/end sprints with auto task promotion (backlog→todo), velocity tracking, burndown charts, standup summaries, retrospectives
- **Parallel Agents**: `get_next_task` supports `epic_id`, `story_id`, `assignee`, `label` filters for scoped agent work
- **WIP Limits**: Configurable max in-progress tasks (global + per-assignee), enforced on `set_current_task`
- **Dependencies**: Task dependency graph with blocker/blocked-by relationships
- **Engine**: Optional Rust gRPC engine for vector search memory (auto-starts/stops, TOON fallback)
- **Extensible**: Other plugins push tools via `RegisterExternalTools()` — appears in stdio + REST

### Three-Layer Database

- **PostgreSQL** (cloud) — Source of truth. pgvector for embeddings, JSONB for settings, tsvector for full-text search, partitioned sync_log
- **SQLite** (local) — Offline support on Desktop and Mobile. Managed by Rust engine (rusqlite) and WatermelonDB (React Native)
- **Redis** — Real-time pub/sub for sync, session cache, rate limiting

### Sync System

All syncable entities use UUID primary keys and include `version`, `created_at`, `updated_at`, `deleted_at`. Changes are logged to `sync_log` and published via Redis pub/sub. Clients push local changes and pull remote changes via WebSocket. Conflict resolution: last-write-wins with version vectors.

### Go Backend (Fiber v3 + GORM)

REST API, WebSocket sync hub, job queue, auth (JWT). Architecture: Handlers → Services → Repositories. All data mutations go through SyncService to log changes.

### Rust Engine (Tonic gRPC)

CPU-intensive operations: Tree-sitter parsing, Tantivy search indexing, file diffing, content hashing, zstd compression, AES-256-GCM encryption, local SQLite management. Go communicates with Rust via gRPC.

### React Frontends (pnpm + Turborepo + Zustand)

Five apps share `@orchestra/shared` (types, stores, hooks, API client) and `@orchestra/ui` (shadcn/ui components, Tailwind CSS v4 theme). Platform-specific code stays in each app directory.

## Skills (Slash Commands)

Every skill is both auto-activated by context AND available as a `/command`. Use `/skill-name` to manually load a skill's patterns and conventions.

| Command | Domain | Technologies |
|---------|--------|-------------|
| `/go-backend` | Go API layer | Fiber v3, GORM, JWT, asynq, gocron, stripe-go, zerolog, go-mail, validator |
| `/rust-engine` | Rust engine | Tonic gRPC, Tree-sitter, Tantivy, tower-lsp, ropey, dashmap, ring, rusqlite |
| `/typescript-react` | Frontend | React, TypeScript, Zustand, React Query, Axios, React Router, Monaco, xterm.js, Vite |
| `/ui-design` | Design system | shadcn/ui, Tailwind CSS v4, Lucide icons, themes, responsive, accessibility |
| `/database-sync` | Data layer | PostgreSQL, pgvector, SQLite, Redis, sync protocol, migrations |
| `/proto-grpc` | Contracts | Protobuf, Buf, tonic-build, Go/Rust code generation |
| `/chrome-extension` | Browser | Chrome Manifest V3, service worker, content scripts, side panel |
| `/wails-desktop` | Desktop | Wails v3, Go-React bindings, system tray, window management |
| `/react-native-mobile` | Mobile | React Native, WatermelonDB, React Navigation, offline sync |
| `/native-widgets` | OS Widgets | macOS WidgetKit, Windows Adaptive Cards, Linux GNOME/KDE |
| `/macos-integration` | macOS | CGo, Spotlight, Keychain, iCloud, Notifications, file associations |
| `/native-extensions` | Extension API | Lifecycle, commands, editor, AI, filesystem, UI, permissions, sandbox |
| `/raycast-compat` | Raycast shim | List/Detail/Form/Action components, ~95% compatibility |
| `/vscode-compat` | VS Code shim | LSP/DAP, themes, snippets, grammars, ~85% compatibility |
| `/extension-marketplace` | Marketplace | Publishing, search, CLI, versioning, reviews, auto-updates |
| `/ai-agentic` | AI/LLM | Anthropic SDK, OpenAI SDK, langchaingo, chromem-go, pgvector, RAG |
| `/gcp-infrastructure` | Infrastructure | Cloud Run, Cloud SQL, CDN, Cloud Build, Docker, nginx, Sentry, PostHog |
| `/project-manager` | Process | Sprint planning, feature breakdown, ADRs, cross-team coordination |
| `/docs` | Documentation | Architecture, plugin system, API references, package relationships |
| `/qa-testing` | QA/Testing | Multi-agent: go test, cargo test, vitest, Playwright, coverage, CI |
| `/plugin-generator` | Plugin scaffolding | new-plugin.sh, tools/storage/transport templates, SDK patterns |

## Agents

Specialized agents in `.claude/agents/` auto-delegate based on task context. See [AGENTS.md](AGENTS.md) for full details.

| Agent | Role |
|-------|------|
| `quic-protocol` | QUIC transport, mTLS, Protobuf framing, wire protocol |
| `go-architect` | Go orchestrator, plugin SDK, Go plugins (quic-go) |
| `rust-engineer` | Rust plugins (quinn, Tree-sitter, Tantivy, rusqlite) |
| `swift-plugin` | Swift/macOS/iOS plugins (Network.framework, SwiftUI, WidgetKit) |
| `kotlin-plugin` | Kotlin/Android plugins (Netty QUIC, Jetpack Compose) |
| `csharp-plugin` | C#/Windows plugins (System.Net.Quic, WinUI 3) |
| `frontend-dev` | React/TypeScript across all 5 platforms |
| `ui-ux-designer` | shadcn/ui, Tailwind, accessibility, responsive |
| `dba` | Cross-database coordination, sync protocol |
| `postgres-dba` | PostgreSQL (pgvector, JSONB, tsvector, partitioning) |
| `sqlite-engineer` | SQLite (rusqlite, go-sqlite3, WatermelonDB) |
| `redis-engineer` | Redis (pub/sub, Streams, caching, rate limiting) |
| `clickhouse-engineer` | ClickHouse (analytics, metrics, OLAP) |
| `lancedb-engineer` | LanceDB (vector search, embeddings, AI memory) |
| `gtk-plugin` | Linux desktop (GTK4, libadwaita, Flatpak) |
| `mobile-dev` | React Native, WatermelonDB, offline sync |
| `scrum-master` | Feature planning, cyclical delivery, WIP limits, ADRs, coordination |
| `widget-engineer` | Native OS widgets (Swift/C#/JS/QML) |
| `platform-engineer` | macOS CGo, Spotlight, Keychain, iCloud |
| `extension-architect` | Extension system (native, Raycast, VS Code, marketplace) |
| `ai-engineer` | AI chat, RAG, agents, embeddings, vector search |
| `devops` | Docker, GCP, CI/CD, monitoring, deployment |
| `qa-go` | Go testing (go test, testify, httptest, plugin tests) |
| `qa-rust` | Rust testing (cargo test, tokio::test, tempfile) |
| `qa-node` | Node/React testing (vitest, @testing-library, component/store tests) |
| `qa-playwright` | E2E browser testing (Playwright, page objects, visual regression) |

## User Interaction Rule (MANDATORY)

**ALWAYS use the `AskUserQuestion` tool when you need user input.** Never print questions as plain text and wait for a response. The scrum-master agent and project-manager skill must use `AskUserQuestion` for:
- PRD session questions (present MCP question via `AskUserQuestion`, then pass answer to `answer_prd_question`)
- Sprint planning decisions (sprint goal, dates, scope)
- Architecture and design choices
- Priority and scope decisions
- Review approval (Gate 4 requires human approval via `AskUserQuestion`)
- Any clarification or confirmation needed from the user

## Mandatory Workflow Rule

**ALL work MUST go through Orchestra MCP tools.** When the user asks you to do ANY task — build, fix, test, refactor, document, investigate, or change anything:

1. `search_features` / `list_features` — check for existing feature
2. `create_feature` — create one if needed (with `kind`: feature/bug/hotfix/chore)
3. `set_current_feature` — start work (moves to in-progress)
4. Do the work
5. `advance_feature` — pass gates with structured evidence
6. `request_review` + `AskUserQuestion` — get user approval
7. `submit_review` — complete

**Never do any work without an active feature.** This includes running tests, writing docs, investigating bugs, and refactoring. The MCP enforces gated transitions — you cannot advance without evidence.

### Feature Kinds

Every feature has a `kind` field: `feature` (default), `bug`, `hotfix`, or `chore`.

- **feature** — New functionality or enhancement
- **bug** — Defect report (Gate 3/docs skipped automatically)
- **hotfix** — Urgent fix (Gate 3/docs skipped automatically)
- **chore** — Maintenance, refactoring, CI work

Use `create_bug_report` as a shortcut for bugs — it sets kind=bug, default priority=P1, and optionally links to the feature that caused the regression via `related_feature`.

### Plan-First for Large Tasks (MANDATORY)

When a user request would result in **3 or more features**, you MUST create a plan before implementation:

1. `create_plan` — Create the plan in `draft` status with title and description
2. Present the plan to the user via `AskUserQuestion` for approval
3. `approve_plan` — Move from draft → approved
4. `breakdown_plan` — Break the plan into features with dependencies (pass a JSON array of feature definitions). This auto-creates all features with `plan:{plan_id}` labels and sets up dependency chains. Plan moves to `in-progress`.
5. Work each feature through the full lifecycle (in order of dependencies)
6. `complete_plan` — After all linked features are `done`, mark the plan as completed

**Do NOT skip the plan step for large tasks.** The plan is stored via MCP and provides traceability.

### User Request Queue

When the user sends a new request while you are busy working on a feature:

1. `create_request` — Save it to the queue with kind (feature/hotfix/bug) and priority
2. Continue working on the current feature
3. After the current feature reaches `done`, call `get_next_request` to pick up the next queued request
4. `convert_request` — Convert it into a feature (auto-creates with correct kind/priority)
5. Work the new feature through the full lifecycle

Use `list_requests` to see the queue and `dismiss_request` to discard irrelevant requests.

### Bug Reporting

When a completed feature causes a regression or breakage:

1. `create_bug_report` — Creates a feature with kind=bug, links to the original feature via `related_feature` param
2. The bug follows the same workflow but **Gate 3 (docs) is auto-skipped** for bugs and hotfixes
3. Work the bug through: backlog → todo → in-progress → testing → review → done

### Enforced Gates (MCP validates evidence)

The MCP **rejects** `advance_feature` if evidence is missing or malformed at gated transitions. Evidence must be markdown with `## Section` headers, each with at least 10 characters of content. **Sections marked with (files) must contain actual file paths** — not just prose.

| Gate | Transition | Required Sections | Tool | Skippable |
|------|-----------|-------------------|------|-----------|
| 1 | in-progress → ready-for-testing | `## Summary`, `## Changes` **(files)**, `## Verification` | `advance_feature` | No |
| 2 | in-testing → ready-for-docs | `## Summary`, `## Results`, `## Coverage` | `advance_feature` | No |
| 3 | in-docs → documented | `## Summary`, `## Location` **(files)** | `advance_feature` | **Yes** (bug, hotfix) |
| 4 | documented → in-review | `## Summary`, `## Quality`, `## Checklist` **(files)** | `request_review` | No |
| 5 | in-review → done | User approval via `AskUserQuestion` | `submit_review` | No |

**Gate evidence format:**
```
evidence: "## Summary\n<what was done>\n\n## Changes\n- libs/foo/bar.go (added validation)\n- libs/baz/qux.go (new file)\n\n## Verification\n<how to test>"
```

Call `get_gate_requirements` to see what's needed for the next transition.

### Free Transitions (no gate)

These transitions can be done without evidence:
- backlog → todo, todo → in-progress, ready-for-testing → in-testing, ready-for-docs → in-docs, needs-edits → in-progress

### Review Flow (Gate 4-5)

1. Call `request_review` with self-review evidence (sections: `## Summary`, `## Quality`, `## Checklist`)
2. MCP moves feature to `in-review` and instructs you to ask the user
3. Use `AskUserQuestion` to present the review to the user with options: "Approve" / "Needs Edits"
4. Call `submit_review` with the user's decision (`status: "approved"` or `status: "needs-edits"`)

**Do NOT call `submit_review` without user approval.** `advance_feature` is blocked from `in-review` — you must use `submit_review`.

## Sub-Agent Orchestration Rules

Sub-agents (launched via the `Task` tool) do **NOT** have access to MCP tools. They cannot call `advance_feature`, `set_current_feature`, or any workflow tools. The main agent must own the full feature lifecycle.

### Rules

1. **Sub-agents are for code writing ONLY** — Use sub-agents only during the `in-progress` phase to write code. They return code results, nothing more.
2. **Main agent owns the lifecycle** — The main agent (you) must handle ALL gate transitions: test, document, review. Never delegate gate work to a sub-agent that can't call MCP tools.
3. **One feature at a time** — Work one feature through its FULL lifecycle (in-progress → done) before starting the next. Never batch multiple features in parallel through gates.
4. **Summarize sub-agent results** — After a sub-agent returns, summarize what it built to the user before advancing. The user must see what happened.
5. **Never mark done without gates** — After a sub-agent writes code, YOU must: run tests (Gate 1), verify coverage (Gate 2), write docs (Gate 3), get user review (Gate 4-5). Each gate needs structured evidence with `## Section` headers.

### Correct Pattern

```
1. set_current_feature(feature_id)              → in-progress
2. Delegate code writing to sub-agent (Task tool)
3. Sub-agent returns → summarize results to user
4. Run tests yourself or delegate to qa-* agent
5. advance_feature(evidence="## Summary\n...\n\n## Changes\n...\n\n## Verification\n...")
                                                 → ready-for-testing [GATE 1]
6. advance_feature                               → in-testing (free)
7. Verify coverage and edge cases
8. advance_feature(evidence="## Summary\n...\n\n## Results\n...\n\n## Coverage\n...")
                                                 → ready-for-docs [GATE 2]
9. advance_feature                               → in-docs (free)
10. Write documentation yourself
11. advance_feature(evidence="## Summary\n...\n\n## Location\n...")
                                                 → documented [GATE 3]
12. request_review(evidence="## Summary\n...\n\n## Quality\n...\n\n## Checklist\n...")
                                                 → in-review [GATE 4]
13. AskUserQuestion → present review to user for approval
14. submit_review(status="approved")             → done [GATE 5]
15. Move to next feature
```

### Anti-Patterns (NEVER DO)

- Spawning 5 sub-agents in parallel, then batch-advancing all 5 features to done
- Letting a sub-agent "handle everything" including testing and docs
- Advancing through gates without providing structured evidence (MCP will reject it)
- Calling `advance_feature` from `in-review` (must use `submit_review`)
- Calling `submit_review` without asking the user via `AskUserQuestion` first
- Starting the next feature before the current one reaches done
- Advancing through multiple gates in rapid succession without doing real work between them
- **Using `sleep`, `wait`, or any delay command to bypass gate cooldowns** — the MCP enforces escalating cooldowns that double for each rapid gate passage, and evidence uniqueness checks that reject copy-pasted content
- Writing fake/boilerplate evidence (e.g., "All tests passed" without actually running tests)
- Copying or templating evidence across gates — each gate requires unique evidence reflecting specific work done at that stage (MCP rejects evidence >60% similar to prior gates)
- Requesting review for one feature, then immediately starting work on another before the review is resolved

### Programmatic Guardrails (MCP-Enforced)

These rules are enforced at the MCP tool level — violation attempts will return errors:

1. **One feature at a time per assignee** — `set_current_feature` checks for any active feature (in-progress through in-review) with the same assignee. Different assignees (parallel agents) can each work on their own feature. MCP returns `wip_violation` error if violated.

2. **Escalating gate cooldown** — Gated transitions require at least 30 seconds since the last status change. If multiple gates are passed within a 5-minute window, the cooldown **doubles for each additional gate** (30s → 60s → 120s → 240s, capped at 10 minutes). This makes `sleep`-based bypass exponentially harder. MCP returns `gate_cooldown` error if violated. **NEVER use `sleep`, `wait`, or any delay command to bypass gate cooldowns.** Do real work between gates instead.

3. **Evidence uniqueness** — New gate evidence is compared against all prior gate evidence in the feature body using Jaccard similarity. If evidence is more than 60% similar to any previous gate's evidence, MCP returns `evidence_duplicate` error. Each gate requires unique, specific evidence reflecting the actual work performed at that stage.

4. **Evidence substance requirements** — Gate evidence must meet minimum content thresholds: at least 20 characters per section, at least 100 characters total (Gate 1/2), 80 characters (Gate 3), and 120 characters (Gate 4/review). Sections requiring file paths must contain at least 1 distinct file path. This prevents minimal/boilerplate evidence.

5. **Timestamped audit trail** — Every transition appends an ISO-8601 timestamp to the feature body. Post-hoc review can detect if gates were passed unrealistically fast.

6. **Review requires user approval** — `advance_feature` is blocked from `in-review`. Only `submit_review` can move to `done`, and it requires `AskUserQuestion` first.

7. **Model capability check** — `set_current_feature` accepts a `model` parameter. If provided and the feature has an estimate, MCP validates the model can handle that size. Tier 1 (Haiku/Flash/GPT-3.5) → S only. Tier 2 (Sonnet/GPT-4o/Gemini Pro) → S, M. Tier 3 (Opus/GPT-4/Gemini Ultra) → S, M, L, XL. Returns `model_capability` error if the model is too small — break the feature down or use a bigger model.

8. **File path evidence** — Gate 1 (Changes), Gate 3 (Location), and Gate 4 (Checklist) sections must contain actual file paths. The MCP rejects evidence that is pure prose without referencing real files.

## Conventions

### Go
- Handler methods: `Index`, `Show`, `Store`, `Update`, `Delete`
- Services contain business logic; repositories are pure data access
- All entities use UUID primary keys with `SyncModel` base
- Error responses: `{"error": "code", "message": "...", "details": {}}`
- Always pass `context.Context` through the call chain
- Use interfaces for services (testability)

### Rust
- Use `thiserror` for typed errors, `anyhow` for application errors
- Never use `unwrap()` in production — use `?` operator
- Use `tokio::task::spawn_blocking` for CPU-heavy synchronous work
- Proto code via `tonic-build` in `build.rs` (not buf for Rust)
- Logging via `tracing` crate

### TypeScript/React
- Import types with `type` keyword
- Zustand stores: separate `State` and `Actions` interfaces
- Use `@orchestra/*` aliases, never relative `../../../` cross-package
- All API responses typed with `ApiResponse<T>`
- Functional components only, `FC` for typing

### Database
- All syncable entities: UUID PK + version + timestamps + soft delete
- PostgreSQL: `TIMESTAMPTZ`; SQLite: ISO 8601 strings
- JSONB for flexible metadata, never for queried fields
- Never store file contents in DB — use content_hash + object storage

# 20 — Plugin Expansion: New Plugins for Desktop App

> Reference document for all new plugins needed before the Swift desktop app.
> Covers: infrastructure changes, 22+ new plugins, markdown parser, AI awareness, voice/TTS/STT, notifications, docs/wiki, devtools, and integrations.

---

## Table of Contents

1. [Current State](#1-current-state)
2. [Infrastructure Changes](#2-infrastructure-changes)
3. [Markdown Parser Plugin](#3-markdown-parser-plugin)
4. [Docs / Wiki Plugin](#4-docs--wiki-plugin)
5. [Notes Plugin](#5-notes-plugin)
6. [AI Awareness Plugins](#6-ai-awareness-plugins)
7. [Chrome Extension Generator Plugin](#7-chrome-extension-generator-plugin)
8. [Voice Service Plugin](#8-voice-service-plugin)
9. [Notification System Plugin](#9-notification-system-plugin)
10. [DevTools Plugins](#10-devtools-plugins)
11. [Integration Plugins](#11-integration-plugins)
12. [Existing Markdown Frontend Components](#12-existing-markdown-frontend-components)
13. [Existing Chrome Extension Reference](#13-existing-chrome-extension-reference)
14. [Existing Screenshot Infrastructure](#14-existing-screenshot-infrastructure)
15. [Engine RAG Knowledge Infrastructure](#15-engine-rag-knowledge-infrastructure)
16. [Plugin Template Pattern](#16-plugin-template-pattern)
17. [Distribution Strategy](#17-distribution-strategy)

---

## 1. Current State

**11 binaries, 90 MCP tools:**

| Binary | Plugin ID | Tools | Description |
|--------|-----------|-------|-------------|
| `orchestra` | CLI | — | CLI commands (serve, init, version, pack) |
| `orchestrator` | orchestrator | — | Plugin loader, message router, lifecycle |
| `transport-stdio` | transport.stdio | — | MCP stdio JSON-RPC bridge |
| `transport-quic-bridge` | transport.quic-bridge | — | Desktop/web QUIC client bridge |
| `storage-markdown` | storage.markdown | — | Protobuf metadata + Markdown body storage |
| `tools-features` | tools.features | 34 | Feature CRUD, workflow, sprints, reviews, deps |
| `tools-marketplace` | tools.marketplace | 15 | Pack install/remove/search, stacks, config |
| `tools-agentops` | tools.agentops | 8 | Agent operation management |
| `tools-sessions` | tools.sessions | 6 | Claude Code session management |
| `bridge-claude` | bridge.claude | 5 | Claude CLI bridge |
| `engine-rag` | engine.rag | 22 | Rust: Tree-sitter, Tantivy, SQLite memory |

---

## 2. Infrastructure Changes

### 2.1 Streaming (Proto + SDK)

**Problem:** AI chat buffers entire response. `bridge-claude/process.go` accumulates all JSONL into `strings.Builder`, returns only after `cmd.Wait()`. QUIC bridge sends one response per stream.

**Solution:** Add streaming messages to proto:

```protobuf
// Request
StreamStart  { stream_id, tool_name, arguments }
StreamCancel { stream_id }

// Response
StreamChunk  { stream_id, data (bytes), content_type }
StreamEnd    { stream_id, success, error_message }
```

SDK changes:
- `RegisterStreamingTool(name, desc, schema, handler)` on builder
- `StreamingToolHandler = func(ctx, req, stream chan<- []byte) error`
- `SendStream(ctx, req) (<-chan *PluginResponse, error)` on client

Fix 4 break points:
1. `bridge-claude/process.go` — yield StreamChunk per JSONL line
2. `sdk-go/plugin/server.go` — dispatch StreamStart, run handler in goroutine
3. `sdk-go/plugin/client.go` — read multiple responses from one QUIC stream
4. `transport-quic-bridge/bridge.go` — forward StreamChunk as JSON-RPC notifications

### 2.2 Events (Proto + Orchestrator)

**Problem:** Proto has `provides_events`/`needs_events` manifest fields but ZERO runtime implementation.

**Solution:** Add event messages to proto:

```protobuf
// Request
Subscribe   { subscription_id, topic, filters }
Unsubscribe { subscription_id }
Publish     { topic, event_type, payload, source_plugin }

// Response
EventDelivery { subscription_id, topic, event_type, payload, source_plugin, timestamp }
```

Orchestrator changes:
- `router.go` — `eventSubscriptions map[string][]subscription`, fan-out on Publish
- `handler.go` — handle Subscribe/Unsubscribe/Publish message types
- `loader.go` — parse `provides_events`/`needs_events` for auto-subscriptions at boot

---

## 3. Markdown Parser Plugin

**Plugin ID:** `tools.markdown`
**Repo:** `github.com/orchestra-mcp/plugin-tools-markdown`

**Purpose:** Backend markdown parser returning typed AST JSON. Single source of truth for all platforms (Swift, React, React Native). The frontend already has a custom zero-dependency TS parser (`parseMarkdown.ts`, ~400 lines) but it lacks Mermaid, KaTeX, footnotes, and embeds. This Go plugin replaces per-platform parser duplication.

**Tools:**

| Tool | Description |
|------|-------------|
| `md_parse` | Parse markdown string → structured AST JSON (blocks + inline nodes) |
| `md_parse_file` | Parse markdown file from path → AST |
| `md_parse_frontmatter` | Extract YAML frontmatter metadata from markdown |
| `md_render_html` | AST → HTML (optional, for email/export) |
| `md_render_plaintext` | AST → plain text (for search indexing, summaries) |
| `md_toc` | Extract table of contents from headings |
| `md_lint` | Validate markdown against style rules |
| `md_transform` | Apply transformations (resolve relative links, inject anchors, etc.) |

**AST Node Types:**

```json
{
  "type": "document",
  "children": [
    { "type": "heading", "level": 1, "children": [{ "type": "text", "value": "Title" }] },
    { "type": "paragraph", "children": [
      { "type": "text", "value": "Hello " },
      { "type": "strong", "children": [{ "type": "text", "value": "world" }] }
    ]},
    { "type": "code_block", "language": "go", "value": "func main() {}" },
    { "type": "table", "align": ["left", "center"], "header": [...], "rows": [...] },
    { "type": "mermaid", "value": "graph TD; A-->B;" },
    { "type": "math_block", "value": "E = mc^2" },
    { "type": "task_list", "children": [
      { "type": "task_item", "checked": true, "children": [...] }
    ]},
    { "type": "footnote_ref", "id": "1" },
    { "type": "footnote_def", "id": "1", "children": [...] },
    { "type": "embed", "url": "https://...", "embed_type": "youtube" }
  ],
  "frontmatter": { "title": "...", "tags": [...] }
}
```

**Supported extensions (beyond GFM):**
- Mermaid diagrams (`mermaid` code blocks)
- KaTeX/math (`$inline$` and `$$block$$`)
- Footnotes (`[^1]` references and definitions)
- Embeds (YouTube, Twitter, Figma URLs → typed embed nodes)
- Task lists with checkbox state
- YAML frontmatter extraction
- Table of contents generation
- Heading anchors (auto-slugified IDs)

**Go library:** `github.com/yuin/goldmark` with extensions (goldmark-meta, goldmark-mermaid, goldmark-katex). Custom AST walker to produce JSON output.

---

## 4. Docs / Wiki Plugin

**Plugin ID:** `tools.docs`
**Repo:** `github.com/orchestra-mcp/plugin-tools-docs`

**Purpose:** Project wiki and documentation system. Thin orchestration layer over engine-rag (Tantivy search, Tree-sitter parsing, SQLite memories). Uses orchestrator storage for wiki pages.

**Tools:**

| Tool | Description |
|------|-------------|
| `doc_create` | Create wiki page (title, body, category, tags, parent_id) |
| `doc_get` | Get page by ID or slug |
| `doc_update` | Update page content/metadata |
| `doc_delete` | Delete page |
| `doc_list` | List pages (project, category, tree view) |
| `doc_search` | Hybrid search — delegates to engine-rag `search` + `search_memory` |
| `doc_generate` | Auto-generate docs from code via engine-rag `parse_file`/`get_symbols` |
| `doc_index` | Index wiki pages into engine-rag for search |
| `doc_tree` | Full wiki tree structure (nested pages) |
| `doc_export` | Export wiki as static markdown/HTML |

**Cross-plugin calls (via orchestrator):**
- `parse_file`, `get_symbols`, `get_imports` → engine-rag
- `search`, `search_memory`, `get_context` → engine-rag
- `index_file` → engine-rag
- `get_project_summary` → engine-rag
- `md_parse` → tools.markdown (for rendering)

**Storage:** `.projects/{project}/docs/{slug}.md` with YAML frontmatter.
**Categories:** `api-reference`, `guide`, `architecture`, `tutorial`, `changelog`, `decision-record`

---

## 5. Notes Plugin

**Plugin ID:** `tools.notes`
**Repo:** `github.com/orchestra-mcp/plugin-tools-notes`

**Purpose:** Standalone note CRUD. Replaces `save_note`/`list_notes` from tools-features.

**Tools:**

| Tool | Description |
|------|-------------|
| `create_note` | Create note (title, body, tags, icon, color, project) |
| `get_note` | Get by ID |
| `update_note` | Update fields |
| `delete_note` | Soft-delete |
| `list_notes` | Filter by project, tags, pinned, search |
| `search_notes` | Full-text search |
| `pin_note` | Toggle pin |
| `tag_note` | Add/remove tags |

**Storage:** `.projects/{project}/notes/{id}.md`
**Cleanup:** Remove `save_note`/`list_notes` from `tools-features/internal/tools/metadata.go`

---

## 6. AI Awareness Plugins

Four separate plugins for AI visual and contextual awareness.

### 6.1 ai.screenshot

**Plugin ID:** `ai.screenshot`
**Repo:** `github.com/orchestra-mcp/plugin-ai-screenshot`

**Purpose:** Screen capture and basic annotation. Native capture on each platform.

| Tool | Description |
|------|-------------|
| `capture_screen` | Full screen capture → base64 PNG |
| `capture_region` | Region capture (x, y, width, height) → base64 PNG |
| `capture_window` | Capture specific window by title/PID |
| `capture_interactive` | OS-native selection tool (screencapture -i on macOS) |
| `annotate_screenshot` | Add rectangles, arrows, text to captured image |
| `list_captures` | List recent captures with metadata |

**macOS:** `screencapture` CLI or ScreenCaptureKit (Swift helper binary)
**Linux:** `gnome-screenshot`, `scrot`, or `grim` (Wayland)

### 6.2 ai.vision

**Plugin ID:** `ai.vision`
**Repo:** `github.com/orchestra-mcp/plugin-ai-vision`

**Purpose:** Image analysis via Claude Vision API. Sends screenshots/images to Claude for understanding.

| Tool | Description |
|------|-------------|
| `analyze_image` | Send image to Claude Vision → structured description |
| `extract_text` | OCR: extract all text from image |
| `find_elements` | Detect UI elements (buttons, inputs, labels, menus) with bounding boxes |
| `compare_images` | Visual diff between two images |
| `describe_screen` | High-level description of what's on screen (app state, content) |
| `extract_data` | Extract structured data from image (tables, forms, charts) |

**Uses:** Anthropic SDK (Claude Vision) or OpenAI Vision as fallback. Takes base64 image input (from ai.screenshot or any source).

### 6.3 ai.browser-context

**Plugin ID:** `ai.browser-context`
**Repo:** `github.com/orchestra-mcp/plugin-ai-browser-context`

**Purpose:** Chrome extension browser awareness. Extracts page content, DOM, and visual context from active browser tabs.

| Tool | Description |
|------|-------------|
| `get_page_content` | Extract current page text, headings, metadata |
| `get_page_dom` | Get simplified DOM tree (interactive elements, landmarks) |
| `get_selected_text` | Get user's text selection in browser |
| `get_open_tabs` | List all open tabs with titles and URLs |
| `get_page_screenshot` | Capture visible page as image (via Chrome extension) |
| `navigate_to` | Navigate browser tab to URL |
| `execute_script` | Run JavaScript in page context (sandboxed) |

**Communication:** WebSocket to Chrome extension service worker. Extension sends `page:request_context`, content script extracts and returns `PageContent`.

**Reference:** `orch-ref/resources/chrome/src/background/service-worker.ts` and `orch-ref/resources/chrome/src/content/extract.ts`

### 6.4 ai.screen-reader

**Plugin ID:** `ai.screen-reader`
**Repo:** `github.com/orchestra-mcp/plugin-ai-screen-reader`

**Purpose:** Accessibility tree and element hierarchy. Uses OS accessibility APIs to understand screen structure without screenshots.

| Tool | Description |
|------|-------------|
| `get_accessibility_tree` | Get full accessibility tree of focused app |
| `get_focused_element` | Get currently focused UI element details |
| `find_element` | Find element by label, role, or value |
| `get_element_hierarchy` | Get parent/child element chain |
| `list_windows` | List all windows with app name, title, position |
| `get_window_elements` | Get all interactive elements in a window |

**macOS:** Accessibility API via CGo (AXUIElement)
**Linux:** AT-SPI2 (D-Bus)

---

## 7. Chrome Extension Generator Plugin

**Plugin ID:** `tools.extension-generator`
**Repo:** `github.com/orchestra-mcp/plugin-tools-extension-generator`

**Purpose:** NOT a manager — an orchestrator for developing Chrome extensions that integrate with the Orchestra ecosystem. Scaffolds, builds, and manages Chrome extension projects.

| Tool | Description |
|------|-------------|
| `ext_scaffold` | Generate Chrome extension from template (Manifest V3, service worker, content script) |
| `ext_add_feature` | Add feature to extension (sidebar, popup, content script, background handler) |
| `ext_connect_orchestra` | Add Orchestra WebSocket/QUIC bridge to extension |
| `ext_add_content_script` | Generate content script for specific domain/pattern |
| `ext_build` | Build extension (bundle, minify, zip for Chrome Web Store) |
| `ext_validate` | Validate manifest.json and permissions |
| `ext_list_projects` | List extension projects in workspace |
| `ext_publish` | Prepare for Chrome Web Store publishing |

**Templates:** Manifest V3 boilerplate, Orchestra SDK connector, content script patterns, service worker patterns.

---

## 8. Voice Service Plugin

**Plugin ID:** `services.voice`
**Repo:** `github.com/orchestra-mcp/plugin-services-voice`

**Purpose:** Text-to-speech and speech-to-text using OS native APIs or external providers. Essential for voice commands, audio feedback, meeting transcription, and accessibility.

**Tools:**

| Tool | Description |
|------|-------------|
| `tts_speak` | Text → speech using OS TTS (macOS NSSpeechSynthesizer, Linux espeak/festival) |
| `tts_speak_provider` | Text → speech via provider (ElevenLabs, OpenAI TTS, Google Cloud TTS) |
| `tts_list_voices` | List available voices (OS native + configured providers) |
| `tts_stop` | Stop current speech playback |
| `stt_listen` | Start microphone listening → returns transcribed text (streaming) |
| `stt_transcribe_file` | Audio file → transcribed text |
| `stt_list_models` | List available STT models (OS + Whisper + provider models) |
| `voice_config` | Configure default voice, provider, language, speed, pitch |

**OS TTS:**
- macOS: `say` CLI / NSSpeechSynthesizer via CGo (50+ built-in voices)
- Linux: `espeak-ng` / `festival` / `piper` (neural TTS)

**OS STT:**
- macOS: Speech Recognition framework via CGo (SFSpeechRecognizer)
- Linux: Vosk (offline) / PocketSphinx

**Provider TTS:**
- ElevenLabs API (high-quality neural voices)
- OpenAI TTS API (tts-1, tts-1-hd models)
- Google Cloud Text-to-Speech

**Provider STT:**
- OpenAI Whisper API (transcription + translation)
- Google Cloud Speech-to-Text
- Deepgram API

**Streaming:** `RegisterStreamingTool` for `stt_listen` (real-time transcription chunks as user speaks)

**Reference:** The old Chrome extension had Moonshine STT WASM for Google Meet transcription (`orch-ref/resources/chrome/src/offscreen/offscreen.ts`).

---

## 9. Notification System Plugin

**Plugin ID:** `services.notifications`
**Repo:** `github.com/orchestra-mcp/plugin-services-notifications`

**Purpose:** OS-native push notifications from any plugin or tool. Essential for build completion alerts, test results, deployment status, AI response ready, timer/reminder, and any async operation feedback.

**Tools:**

| Tool | Description |
|------|-------------|
| `notify_send` | Send OS notification (title, body, subtitle, icon, sound, actions) |
| `notify_schedule` | Schedule notification for future time (cron or one-shot) |
| `notify_cancel` | Cancel scheduled notification by ID |
| `notify_list_pending` | List all scheduled/pending notifications |
| `notify_badge` | Set app badge count (macOS dock badge, Linux desktop badge) |
| `notify_config` | Configure preferences (DND hours, channels, sound per category) |
| `notify_history` | Get notification history (last N sent notifications) |
| `notify_create_channel` | Create notification channel/category (build, test, deploy, ai, reminder) |

**macOS:**
- `NSUserNotificationCenter` (legacy) / `UNUserNotificationCenter` (modern) via CGo
- Supports: title, subtitle, body, sound, actions (buttons), reply field
- Badge count via `NSDockTile.setBadgeLabel`

**Linux:**
- D-Bus `org.freedesktop.Notifications` (GNOME/KDE/XFCE)
- Supports: title, body, icon, actions, urgency levels

**Notification channels (pre-configured):**
- `build` — Build started/completed/failed
- `test` — Test suite results
- `deploy` — Deployment status
- `ai` — AI response ready, streaming complete
- `reminder` — Scheduled reminders
- `system` — System alerts (disk, memory, service down)
- `git` — PR merged, review requested, CI status

**Actions:** Clickable notification buttons that can trigger MCP tool calls (e.g., "View Results" → opens test results, "Retry" → re-runs failed build).

**Events integration:** Subscribes to orchestrator events and can auto-notify based on event type rules.

---

## 10. DevTools Plugins

All DevTools plugins are separate GitHub repos from day 1.

### 10.1 devtools.file-explorer (Full IDE + LSP)
**Repo:** `github.com/orchestra-mcp/plugin-devtools-file-explorer`

Not just a file browser — a **full IDE backend** with LSP code intelligence via the Rust engine. Provides file operations, namespace navigation, go-to-definition, autocomplete, hover, and diagnostics.

**File Operation Tools:**
`list_directory`, `read_file`, `write_file`, `move_file`, `delete_file`, `file_info`, `file_search`

**IDE / Code Intelligence Tools (delegates to engine-rag):**

| Tool | Description |
|------|-------------|
| `code_symbols` | Document symbols (functions, classes, structs) — wraps `get_symbols` |
| `code_goto_definition` | Go-to-definition at file:line:col via symbol resolution graph |
| `code_find_references` | Find all references to symbol across workspace |
| `code_hover` | Type, signature, docstring at position |
| `code_complete` | Scope-aware autocomplete at cursor |
| `code_diagnostics` | Syntax errors + lint warnings |
| `code_actions` | Available actions at position (add import, extract fn) |
| `code_workspace_symbols` | Search symbols across workspace — wraps `search_symbols` |
| `code_namespace` | Full namespace/module tree (Go pkgs, Rust mods, TS exports) |
| `code_imports` | Resolve imports for a file |

**Rust Engine LSP Enhancements Required** (new engine-rag tools):

| New Engine Tool | Description |
|-----------------|-------------|
| `lsp_open_document` | Track open doc, cache parse tree |
| `lsp_close_document` | Remove from cache |
| `lsp_update_document` | Incremental edit via `parse_incremental()` |
| `lsp_goto_definition` | Resolve symbol → target file:line:col |
| `lsp_find_references` | All refs via reverse symbol index |
| `lsp_hover` | Type + signature + docstring |
| `lsp_complete` | Scope-aware suggestions |
| `lsp_diagnostics` | Tree-sitter errors + lint |
| `lsp_workspace_symbols` | Tantivy search with type/kind filters |
| `lsp_build_index` | Build/rebuild symbol resolution graph |

**Existing Rust engine capabilities (already built):**
- Tree-sitter: 14 grammars, `parse_incremental()`, 19 symbol kinds
- `CodeSymbol` struct: name, kind, range (line:col), detail (signature), children
- Tantivy: full-text search over symbol names
- Import extraction, `.gitignore`-aware indexing

**New Rust engine additions needed:**
- Symbol resolution graph: cross-file import→definition in SQLite (`.orchestra/lsp.db`)
- Document state cache: open docs + cached parse trees
- Docstring extraction: Tree-sitter queries for JSDoc/rustdoc/pydoc
- Scope analysis: AST parent walk for available symbols at cursor
- Module resolution: per-language import path resolver (Go, Rust, TS)

### 10.2 devtools.terminal
**Repo:** `github.com/orchestra-mcp/plugin-devtools-terminal`
**Tools:** `create_terminal`, `send_input`, `get_output`, `resize_terminal`, `list_terminals`, `close_terminal`
**Deps:** `github.com/creack/pty`

### 10.3 devtools.ssh
**Repo:** `github.com/orchestra-mcp/plugin-devtools-ssh`
**Tools:** `ssh_connect`, `ssh_exec`, `ssh_disconnect`, `ssh_list_sessions`, `ssh_upload`, `ssh_download`, `ssh_list_remote`
**Deps:** `golang.org/x/crypto/ssh`

### 10.4 devtools.services
**Repo:** `github.com/orchestra-mcp/plugin-devtools-services`
**Tools:** `list_services`, `start_service`, `stop_service`, `restart_service`, `service_logs`, `service_info`
**Platform:** macOS `launchctl`, Linux `systemctl`

### 10.5 devtools.docker
**Repo:** `github.com/orchestra-mcp/plugin-devtools-docker`
**Tools:** `docker_list_containers`, `docker_start`, `docker_stop`, `docker_restart`, `docker_logs`, `docker_exec`, `docker_list_images`, `docker_compose_up`, `docker_compose_down`, `docker_inspect`
**Deps:** `github.com/docker/docker/client`

### 10.6 devtools.debugger
**Repo:** `github.com/orchestra-mcp/plugin-devtools-debugger`
**Tools:** `debug_start`, `debug_stop`, `debug_set_breakpoint`, `debug_remove_breakpoint`, `debug_continue`, `debug_step_over`, `debug_step_into`, `debug_evaluate`
**Protocol:** DAP (Debug Adapter Protocol) — delve, lldb, node-debug

### 10.7 devtools.test-runner
**Repo:** `github.com/orchestra-mcp/plugin-devtools-test-runner`
**Tools:** `test_discover`, `test_run`, `test_run_suite`, `test_results`, `test_coverage`, `test_watch`
**Frameworks:** Playwright (primary), go test, cargo test, vitest

### 10.8 devtools.log-viewer
**Repo:** `github.com/orchestra-mcp/plugin-devtools-log-viewer`
**Tools:** `log_watch`, `log_search`, `log_list_sources`, `log_tail`, `log_parse`

### 10.9 devtools.database
**Repo:** `github.com/orchestra-mcp/plugin-devtools-database`
**Tools:** `db_connect`, `db_disconnect`, `db_query`, `db_list_tables`, `db_describe_table`, `db_list_connections`, `db_export`, `db_import`
**Deps:** `lib/pq`, `go-sql-driver/mysql`, `mattn/go-sqlite3`

### 10.10 devtools.devops
**Repo:** `github.com/orchestra-mcp/plugin-devtools-devops`
**Tools:** `devops_list_pipelines`, `devops_trigger_pipeline`, `devops_pipeline_status`, `devops_pipeline_logs`, `devops_list_deployments`, `devops_deploy`, `devops_rollback`, `devops_env_vars`
**Deps:** `google/go-github`

---

## 11. Integration Plugins

### 11.1 integration.figma
**Repo:** `github.com/orchestra-mcp/plugin-integration-figma`
**Tools:** `figma_get_file`, `figma_get_components`, `figma_get_styles`, `figma_export_node`, `figma_get_node`, `figma_sync_tokens`
**Deps:** Figma REST API

### 11.2 devtools.components
**Repo:** `github.com/orchestra-mcp/plugin-devtools-components`
**Tools:** `component_list`, `component_preview`, `component_create`, `component_inspect`, `component_sync_figma`, `component_library`
**Cross-plugin:** integration.figma + engine-rag

---

## 12. Existing Markdown Frontend Components

**Location:** `packages/@orchestra-mcp/editor/`

| Component | File | Purpose |
|-----------|------|---------|
| `MarkdownRenderer` | `src/MarkdownRenderer/MarkdownRenderer.tsx` | Read-only markdown display |
| `parseMarkdown` | `src/MarkdownRenderer/parseMarkdown.ts` | Custom zero-dependency parser (~400 lines) |
| `inlineFormat` | `src/MarkdownRenderer/inlineFormat.ts` | Bold, italic, code, links, images |
| `MarkdownEditor` | `src/MarkdownEditor/MarkdownEditor.tsx` | Split-pane WYSIWYG with live preview |
| `CodeBlock` | `src/CodeBlock/CodeBlock.tsx` | Syntax highlighting (regex-based, 20+ langs) |
| `ChatMarkdown` | `@orchestra-mcp/ai/src/ChatMarkdown/` | Chat-optimized markdown wrapper |
| `DataTable` | `@orchestra-mcp/widgets/src/DataTable/` | Table rendering with sort/export |

**Current parser supports:** Headings, paragraphs, bold, italic, strikethrough, inline code, code blocks, tables, links, images, lists (ul/ol), task lists, blockquotes, horizontal rules, heading anchors.

**Missing from current parser:** Mermaid, KaTeX/math, footnotes, embeds, custom HTML, AST output.

**No external markdown libraries used** (no remark, marked, markdown-it, react-markdown, Prism, highlight.js).

---

## 13. Existing Chrome Extension Reference

**Location:** `orch-ref/resources/chrome/src/`

| Component | File | Purpose |
|-----------|------|---------|
| Service Worker | `background/service-worker.ts` | WebSocket to desktop, tab mgmt, page relay |
| Content Script | `content/extract.ts` | Zero-dep page extraction (title, headings, content, meta) |
| Offscreen Doc | `offscreen/offscreen.ts` | Audio capture (Moonshine STT for Google Meet) |

**PageContent interface:**
```typescript
interface PageContent {
  title: string
  url: string
  selectedText: string
  mainContent: string       // max 10,000 chars
  metaDescription: string
  headings: string[]
  language: string
  isArticle: boolean
}
```

**Communication:** WebSocket (`ws://localhost:8765`) between extension and desktop app.

---

## 14. Existing Screenshot Infrastructure

**Location:** `orch-ref/app/settings/`

| File | Platform | Method |
|------|----------|--------|
| `screenshot_darwin.go` | macOS | `screencapture -i -s -x` → PNG bytes |
| `screenshot_other.go` | Linux | `gnome-screenshot`, `scrot`, `grim`, `ImageMagick` |

**HTTP endpoints (NOT MCP tools):**
- `POST /api/screenshot/capture` — full screen
- `POST /api/screenshot/region` — crop to x,y,w,h
- `POST /api/screenshot/start` — interactive selection

**Response:** `{ "ok": true, "image": "<base64 PNG>" }`

---

## 15. Engine RAG Knowledge Infrastructure

**22 MCP tools** the docs plugin can leverage:

| Service | Tools | What It Does |
|---------|-------|--------------|
| Parse | `parse_file`, `get_symbols`, `get_imports` | Tree-sitter code parsing (14 grammars) |
| Search | `index_file`, `index_directory`, `search`, `search_symbols`, `delete_from_index`, `clear_index`, `get_index_stats` | Tantivy full-text indexing |
| Memory | `save_memory`, `get_memory`, `update_memory`, `delete_memory`, `list_memories`, `search_memory`, `get_context` | SQLite + cosine similarity |
| Session | `start_session`, `end_session`, `save_observation`, `get_project_summary` | Agent session tracking |

**Memory categories:** `project-context`, `decisions`, `patterns`, `knowledge`, `feedback`, `sessions`

**SQLite tables:** sessions, observations, summaries, memories, embeddings

---

## 16. Plugin Template Pattern

Generated by `scripts/new-plugin.sh`:

```
libs/plugin-{name}/
├── cmd/main.go              # clientAdapter + plugin.New() builder
├── internal/
│   ├── plugin.go            # RegisterTools(builder)
│   ├── tools/*.go           # Schema() + Handler() per tool
│   └── storage/client.go    # orchestrator storage access
├── go.mod                   # github.com/orchestra-mcp/plugin-{name}
├── orchestra.json           # plugin manifest
├── README.md, LICENSE, CI, docs/
```

**Builder API:**
```go
plugin.New("tools.notes").
    Version("0.1.0").
    NeedsStorage("markdown").
    RegisterTool(name, desc, schema, handler).
    BuildWithTools().
    Run(ctx)
```

**SDK helpers:** `GetString`, `GetInt`, `GetBool`, `ValidateRequired`, `TextResult`, `JSONResult`, `ErrorResult`

---

## 17. Distribution Strategy

**All new plugins as separate GitHub repos from day 1:**

- Build locally in `libs/plugin-{name}/` during development
- Each has own `go.mod` → `github.com/orchestra-mcp/plugin-{name}`
- Use `scripts/sync-repos.sh plugin-{name}` to push to separate repo
- Use `scripts/release.sh v0.x.0 plugin-{name}` for releases
- GitHub Actions CI per repo (from `new-plugin.sh` template)
- Same pattern as marketplace packs (each integration = separate repo)

**Total new plugins: 22**
- Infrastructure: 6 changes (proto, SDK, orchestrator, bridge-claude, quic-bridge)
- Core: 3 (tools.markdown, tools.docs, tools.notes)
- AI Awareness: 4 (ai.screenshot, ai.vision, ai.browser-context, ai.screen-reader)
- System Services: 3 (services.voice, services.notifications, tools.extension-generator)
- DevTools: 10 (file-explorer, terminal, ssh, services, docker, debugger, test-runner, log-viewer, database, devops)
- Integration: 2 (integration.figma, devtools.components)

**New tool count: ~180 tools** → Total: ~270 MCP tools

# 18 — Swift Universal App: Architecture & Implementation Guide

> Comprehensive reference for building the Orchestra MCP universal Apple app in Swift/SwiftUI.
> Supports macOS, iOS, iPadOS, watchOS, tvOS, and visionOS from a single codebase.
> Plugin-based architecture mirroring the Go framework — every feature is a plugin.
> Compiled from screenshots, plugin expansion (artifact 20), multi-agent orchestrator (artifact 19), and architecture decisions.

---

## Table of Contents

1. [Vision & Context](#1-vision--context)
2. [Architecture](#2-architecture)
3. [App Structure](#3-app-structure)
4. [Swift Plugin System](#4-swift-plugin-system)
5. [Navigation & Screens](#5-navigation--screens)
6. [Design System](#6-design-system)
7. [Window Management](#7-window-management)
8. [QUIC Transport](#8-quic-transport)
9. [Data Layer](#9-data-layer)
10. [AI Chat](#10-ai-chat)
11. [Project Management](#11-project-management)
12. [Notes & Docs](#12-notes--docs)
13. [Developer Tools](#13-developer-tools)
14. [AI Awareness](#14-ai-awareness)
15. [Voice & Notifications](#15-voice--notifications)
16. [Settings & Integrations](#16-settings--integrations)
17. [Native Apple Features](#17-native-apple-features)
18. [Build Phases](#18-build-phases)

---

## 1. Vision & Context

### What We're Building

A universal Apple app for Orchestra MCP — an AI-agentic IDE that manages projects, features, sprints, and AI chat sessions across all Apple platforms. The app connects to the Orchestra plugin ecosystem via QUIC and exposes 270+ MCP tools through a native SwiftUI interface.

### Platform Matrix

| Platform | Experience | Key Features |
|----------|-----------|--------------|
| **macOS** | Full IDE | Multi-window, system tray, global hotkeys, DevTools, Spirit floating chat, ScreenCaptureKit |
| **iOS/iPadOS** | Project management + AI chat | Split view (iPad), tab navigation (iPhone), offline cache, push notifications |
| **watchOS** | Quick status + notifications | Complications, project progress glance, notification actions |
| **tvOS** | Dashboard + monitoring | Project dashboards, sprint burndown, real-time metrics display |
| **visionOS** | Spatial workspace | Immersive code review, multi-window spatial layout, spatial AI chat |
| **CarPlay** | Voice-first AI assistant | Voice commands, project status read-aloud, hands-free AI chat, notification summaries |

### Plugin Ecosystem (Current State)

**33 plugins, 270+ MCP tools** (11 current + 22 planned from plugin expansion):

| Category | Plugins | Tools |
|----------|---------|-------|
| Core | orchestrator, transport.stdio, transport.quic-bridge, storage.markdown | — |
| Features | tools.features, tools.marketplace | 49 |
| AI Bridges | bridge.claude, bridge.openai, bridge.gemini, bridge.ollama | 20 |
| Agent Ops | tools.agentops, tools.sessions | 14 |
| Engine | engine.rag (Rust) | 22 |
| Content | tools.markdown, tools.docs, tools.notes | 26 |
| AI Awareness | ai.screenshot, ai.vision, ai.browser-context, ai.screen-reader | 25 |
| Services | services.voice, services.notifications, tools.extension-generator | 24 |
| DevTools | devtools.file-explorer, .terminal, .ssh, .services, .docker, .debugger, .test-runner, .log-viewer, .database, .devops | 80 |
| Integrations | integration.figma, devtools.components | 12 |
| Multi-Agent | agent-orchestrator (ADK) | 20 |

### The Old App (Screenshots Reference)

The old Wails-based app had these sections:

| Section | Icon | Description |
|---------|------|-------------|
| **Chats** | Speech bubble | AI chat sessions with Claude, multi-model, conversation history, tool call results, streaming |
| **Projects** | Grid | Project list with task counts, completion %, status badges |
| **Notes** | Calendar | Note list with search, rich markdown editor |
| **Developer Tools** | Terminal | File Explorer, Terminal, Database, SSH, Log Viewer, Services, Debugger |
| **Components** | Layers | Component browser with preview |
| **Integrations** | Settings | AI Providers (7 providers), service integrations (Discord, Slack, GitHub, etc.) |
| **Settings** | Gear | General, Appearance, Notifications, Windows, AI, Voice, Sync & Account |

### Key UI Patterns from Screenshots

- **Dark theme** with deep navy/black backgrounds
- **Purple accent** (#a900ff) for active states, brand elements, CTAs
- **Left sidebar** with icon rail (56px) for section switching
- **Session/item sidebar** (280px) for lists within each section
- **Main content area** taking remaining space
- **User profile** at bottom-left with avatar, name, email, status dot
- **Green "Live" indicators** for active connections
- **Badge counts** on tool results
- **Model badges** in purple pills

---

## 2. Architecture

### Plugin Host Pattern

The Swift app is a consumer in the star-topology orchestrator architecture. It connects to independently running plugins via QUIC.

```
                      ORCHESTRATOR
                    (Go, ~500 lines)
                         |
            QUIC + Protobuf (mTLS)
                         |
    ┌───────────┬────────┴────────┬───────────┐
transport    tools.features    storage     engine.rag
 .stdio       (34 tools)     .markdown      (Rust)
  (Go)          (Go)            (Go)       22 tools

    ┌───────────┬────────────────┬───────────┐
bridge.claude  bridge.openai  bridge.gemini  bridge.ollama
  (5 tools)     (5 tools)      (5 tools)     (5 tools)

    ┌───────────┬────────────────┬───────────┐
tools.agentops tools.sessions  tools.marketplace
  (8 tools)     (6 tools)       (15 tools)

    ┌───────────┬────────────────┐
agent-orchestrator  tools.markdown  tools.docs  tools.notes
   (20 tools)        (8 tools)      (10 tools)   (8 tools)

    ┌───────────┐
services.voice  services.notifications
  (8 tools)       (8 tools)

    ┌───────────────────────────────────────────┐
devtools.file-explorer  .terminal  .ssh  .docker  ...
              (80 tools across 10 plugins)

                                    orchestra-swift ← THIS APP
                                    (Universal Apple)
```

### Communication Contract

- **Transport**: QUIC via `Network.framework` (Apple native)
- **Auth**: mTLS — CA at `~/.orchestra/certs/ca.crt`, app cert signed by CA
- **Wire format**: Length-delimited Protobuf (4-byte big-endian uint32 length + Protobuf bytes)
- **Proto contract**: Generated Swift code from `swift-protobuf`
- **Message routing**: All messages go through orchestrator — never direct plugin-to-plugin
- **Streaming**: StreamStart → StreamChunk* → StreamEnd (for AI chat, long-running ops)
- **Events**: Subscribe/Publish/EventDelivery (for real-time updates)
- **Tool calls**: Send `ToolRequest` → receive `ToolResponse` (with optional `provider` field for AI routing)

### Platform Availability

| Framework | macOS | iOS | watchOS | tvOS | visionOS | CarPlay |
|-----------|-------|-----|---------|------|----------|---------|
| Network.framework (QUIC) | 14+ | 17+ | 10+ | 17+ | 1.0+ | 17+ |
| SwiftUI | 14+ | 17+ | 10+ | 17+ | 1.0+ | — |
| CarPlay.framework | — | 17+ | — | — | — | 17+ |
| ScreenCaptureKit | 14+ | — | — | — | — | — |
| UserNotifications | 14+ | 17+ | 10+ | 17+ | 1.0+ | 17+ |
| Speech (STT) | 14+ | 17+ | — | — | 1.0+ | 17+ |
| AVSpeechSynthesizer (TTS) | 14+ | 17+ | — | — | 1.0+ | 17+ |

### Graceful Degradation

If the orchestrator is not running, the app should:
1. Show a "Not Connected" status in the status bar
2. Allow browsing locally cached data
3. Retry connection with exponential backoff (1s → 30s max)
4. Auto-reconnect when orchestrator becomes available

---

## 3. App Structure

### Project Layout

```
apps/swift/                                # github.com/orchestra-mcp/orchestra-swift
├── Package.swift                          # SPM multi-platform package
├── OrchestraKit/
│   └── Sources/OrchestraKit/
│       ├── Transport/
│       │   ├── QUICConnection.swift       # Network.framework QUIC client
│       │   ├── StreamFramer.swift         # Length-delimited Protobuf framing
│       │   └── MTLSConfig.swift           # mTLS cert loading
│       ├── Proto/
│       │   └── Messages.swift             # Proto message types (placeholder until buf gen)
│       ├── Models/
│       │   ├── AppState.swift             # ObservableObject root state
│       │   ├── Project.swift              # Project model
│       │   ├── Feature.swift              # Feature/task model
│       │   ├── Note.swift                 # Note model
│       │   ├── ChatSession.swift          # Chat session model
│       │   └── ChatMessage.swift          # Chat message model
│       ├── Plugins/
│       │   ├── OrchestraPlugin.swift      # Plugin protocol + AppSection
│       │   └── PluginRegistry.swift       # Plugin registry + discovery
│       └── Services/
│           ├── OrchestraClient.swift      # High-level orchestrator client
│           ├── ToolService.swift          # MCP tool call proxy
│           └── ConnectionState.swift      # Connection status enum
├── Shared/
│   └── Sources/Shared/
│       ├── App/
│       │   └── ContentView.swift          # Root view (plugin-driven, platform-adaptive)
│       ├── Plugins/                       # Built-in plugins (each is self-contained)
│       │   ├── ChatPlugin/
│       │   │   ├── ChatPlugin.swift       # Plugin registration
│       │   │   ├── ChatView.swift         # Chat layout
│       │   │   ├── ChatSessionList.swift  # Session sidebar
│       │   │   └── ChatMessageView.swift  # Message bubble
│       │   ├── ProjectsPlugin/
│       │   │   ├── ProjectsPlugin.swift
│       │   │   ├── ProjectsView.swift
│       │   │   └── ProjectDetailView.swift
│       │   ├── NotesPlugin/
│       │   │   ├── NotesPlugin.swift
│       │   │   └── NotesView.swift
│       │   ├── DevToolsPlugin/
│       │   │   ├── DevToolsPlugin.swift
│       │   │   └── DevToolsView.swift
│       │   └── SettingsPlugin/
│       │       ├── SettingsPlugin.swift
│       │       └── SettingsView.swift
│       ├── Components/
│       │   ├── StatusBadge.swift
│       │   ├── ConnectionIndicator.swift
│       │   └── EmptyStateView.swift
│       └── Theme/
│           ├── OrchestraTheme.swift       # Color tokens
│           └── OrchestraTypography.swift   # Font definitions
├── macOS/
│   └── Sources/macOS/
│       ├── OrchestraApp+macOS.swift       # macOS entry with MenuBarExtra
│       └── MacCommands.swift              # Keyboard shortcuts
├── iOS/
│   └── Sources/iOS/
│       └── OrchestraApp+iOS.swift         # iOS entry with TabView
├── watchOS/
│   └── Sources/watchOS/
│       └── OrchestraApp+watchOS.swift     # watchOS stub
├── tvOS/
│   └── Sources/tvOS/
│       └── OrchestraApp+tvOS.swift        # tvOS stub
├── visionOS/
│   └── Sources/visionOS/
│       └── OrchestraApp+visionOS.swift    # visionOS stub
├── CarPlay/
│   └── Sources/CarPlay/
│       └── OrchestraCarPlay.swift         # CarPlay scene + voice-first interface
├── scripts/
│   └── new-swift-plugin.sh               # Plugin creator script
└── README.md
```

### App Lifecycle (macOS)

```swift
@main
struct OrchestraApp: App {
    @NSApplicationDelegateAdaptor(AppDelegate.self) var appDelegate
    @StateObject private var appState = AppState()
    @StateObject private var pluginRegistry = PluginRegistry()

    init() {
        // Register built-in plugins
        pluginRegistry.register(ChatPlugin())
        pluginRegistry.register(ProjectsPlugin())
        pluginRegistry.register(NotesPlugin())
        pluginRegistry.register(DevToolsPlugin())
        pluginRegistry.register(SettingsPlugin())
    }

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(appState)
                .environmentObject(pluginRegistry)
        }
        .defaultSize(width: 1280, height: 860)
        .commands { MacCommands() }

        Window("Spirit", id: "spirit") {
            SpiritWindow()
                .environmentObject(appState)
        }
        .windowStyle(.hiddenTitleBar)
        .defaultSize(width: 420, height: 640)

        MenuBarExtra("Orchestra", systemImage: "waveform") {
            TrayMenu()
                .environmentObject(appState)
        }

        Settings {
            SettingsView()
                .environmentObject(appState)
        }
    }
}
```

### App Lifecycle (iOS)

```swift
@main
struct OrchestraApp: App {
    @StateObject private var appState = AppState()
    @StateObject private var pluginRegistry = PluginRegistry()

    init() {
        pluginRegistry.register(ChatPlugin())
        pluginRegistry.register(ProjectsPlugin())
        pluginRegistry.register(NotesPlugin())
        pluginRegistry.register(SettingsPlugin())
        // No DevTools on iOS
    }

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(appState)
                .environmentObject(pluginRegistry)
        }
    }
}
```

### Startup Sequence

```
1. Load settings from platform-appropriate location
2. Initialize theme (color theme + component variant)
3. Register built-in plugins with PluginRegistry
4. Setup UI (main window on macOS, tab view on iOS)
5. Connect to orchestrator via QUIC
6. Subscribe to events (feature.*, workflow.*, sprint.*, note.*)
7. Load cached data (projects, sessions, notes)
8. Check for updates (macOS: 5s delay, then every 6 hours)
9. Request platform permissions (first run only)
10. Restore window state (macOS: embedded/floating/bubble)
```

---

## 4. Swift Plugin System

### Design Principles

The app mirrors the Go framework's plugin architecture. Every screen/feature is a Swift plugin that registers itself with the PluginRegistry. This enables:
- **Incremental development** — build one plugin at a time for small wins
- **Easy extensibility** — `new-swift-plugin.sh` scaffolds new plugins
- **Feature isolation** — each plugin is self-contained with its own views and logic
- **Platform adaptation** — plugins declare which platforms they support

### Plugin Protocol

```swift
/// Section where the plugin appears in the UI.
enum AppSection: String, CaseIterable {
    case sidebar     // Main navigation (Chat, Projects, Notes)
    case devtools    // Developer Tools sub-section
    case settings    // Settings sub-section
}

/// Every feature in the app conforms to this protocol.
protocol OrchestraPlugin: Identifiable {
    var id: String { get }
    var name: String { get }
    var icon: String { get }             // SF Symbol name
    var section: AppSection { get }
    var order: Int { get }               // Sort order within section
    var supportedPlatforms: Set<Platform> { get }

    @ViewBuilder func makeView() -> AnyView
    func onActivate()
    func onDeactivate()
}

enum Platform {
    case macOS, iOS, iPadOS, watchOS, tvOS, visionOS, carPlay
}
```

### Plugin Registry

```swift
@MainActor
class PluginRegistry: ObservableObject {
    @Published private(set) var plugins: [any OrchestraPlugin] = []

    func register(_ plugin: any OrchestraPlugin) {
        plugins.append(plugin)
        plugins.sort { $0.order < $1.order }
    }

    func plugin(for id: String) -> (any OrchestraPlugin)? {
        plugins.first { $0.id == id }
    }

    var sidebarPlugins: [any OrchestraPlugin] {
        plugins.filter { $0.section == .sidebar && $0.supportedPlatforms.contains(currentPlatform) }
    }

    var devToolPlugins: [any OrchestraPlugin] {
        plugins.filter { $0.section == .devtools && $0.supportedPlatforms.contains(currentPlatform) }
    }

    var settingsPlugins: [any OrchestraPlugin] {
        plugins.filter { $0.section == .settings && $0.supportedPlatforms.contains(currentPlatform) }
    }
}
```

### Built-in Plugins (Phase 1)

| Plugin | ID | Section | Platforms |
|--------|----|---------|-----------|
| ChatPlugin | `chat` | sidebar | All |
| ProjectsPlugin | `projects` | sidebar | macOS, iOS, iPadOS, tvOS, visionOS |
| NotesPlugin | `notes` | sidebar | macOS, iOS, iPadOS, visionOS |
| DevToolsPlugin | `devtools` | sidebar | macOS, visionOS |
| SettingsPlugin | `settings` | sidebar | All |

### Plugin Creator Script

`scripts/new-swift-plugin.sh` generates:

```bash
./scripts/new-swift-plugin.sh my-feature sidebar
# Creates:
# Shared/Sources/Shared/Plugins/MyFeaturePlugin/
#   ├── MyFeaturePlugin.swift   # Plugin registration
#   └── MyFeatureView.swift     # Main view
# Prints registration line to add to app entry point
```

---

## 5. Navigation & Screens

### Platform-Adaptive Navigation

```swift
struct ContentView: View {
    @EnvironmentObject var registry: PluginRegistry
    @State private var selectedPlugin: String? = "chat"

    var body: some View {
        #if os(macOS) || os(visionOS)
        NavigationSplitView {
            SidebarView(selected: $selectedPlugin)
        } detail: {
            if let id = selectedPlugin, let plugin = registry.plugin(for: id) {
                plugin.makeView()
            }
        }
        #elseif os(iOS)
        if UIDevice.current.userInterfaceIdiom == .pad {
            NavigationSplitView { ... } detail: { ... }
        } else {
            TabView(selection: $selectedPlugin) {
                ForEach(registry.sidebarPlugins) { plugin in
                    plugin.makeView()
                        .tabItem {
                            Label(plugin.name, systemImage: plugin.icon)
                        }
                        .tag(plugin.id)
                }
            }
        }
        #elseif os(watchOS)
        NavigationStack {
            List(registry.sidebarPlugins) { plugin in ... }
        }
        #elseif os(tvOS)
        TabView { ... }
        #endif
    }
}
```

### Sidebar (Icon Rail — macOS/iPadOS)

| # | Icon | Section | Route |
|---|------|---------|-------|
| 1 | `bubble.left.and.bubble.right` | Chat | `/chat` |
| 2 | `square.grid.2x2` | Projects | `/projects` |
| 3 | `note.text` | Notes | `/notes` |
| 4 | `terminal` | Developer Tools | `/devtools` |
| — | *(spacer)* | | |
| 5 | `gearshape` | Settings | `/settings` |

**Active state**: Purple accent background, white icon.
**Inactive state**: Muted gray icon, transparent background.
**Brand logo**: 36x36 purple rounded square at top.
**User profile**: Avatar at bottom with context menu.

### Cross-Platform Feature Matrix

| Feature | macOS | iOS/iPadOS | watchOS | tvOS | visionOS | CarPlay |
|---------|-------|-----------|---------|------|----------|---------|
| Chat | Full | Full | Quick reply | View | Spatial | Voice |
| Projects | Full | Full | Status | Dashboard | Full | Status |
| Notes | Full | Full | — | — | Full | — |
| DevTools | Full (10 tools) | — | — | — | Full | — |
| System Tray | Yes | — | — | — | — | — |
| Window Modes | 3 modes | — | — | — | Volumes | — |
| Screenshot | Yes | — | — | — | — | — |
| Voice STT/TTS | Yes | Yes | — | — | Yes | Yes |
| Widgets | Yes | Yes | Complications | — | — | — |
| Global Hotkey | Yes | — | — | — | — | — |
| Push Notifications | Yes | Yes | Yes | — | Yes | Yes |

---

## 6. Design System

### Color Tokens

```swift
extension Color {
    // Backgrounds
    static let background       = Color(hex: "#0a0d14")  // --color-bg
    static let surface          = Color(hex: "#111520")  // --color-bg-alt
    static let surfaceContrast  = Color(hex: "#080a10")  // --color-bg-contrast
    static let surfaceActive    = Color(hex: "#1a1f2e")  // --color-bg-active
    static let surfaceSelection = Color(hex: "#a900ff").opacity(0.15)

    // Text
    static let textPrimary      = Color(hex: "#e8ecf4")  // --color-fg
    static let textMuted        = Color(hex: "#8892a8")  // --color-fg-muted
    static let textDim          = Color(hex: "#4a5268")  // --color-fg-dim
    static let textBright       = Color(hex: "#f8fafc")  // --color-fg-bright

    // Structure
    static let border           = Color(hex: "#1e2436")  // --color-border
    static let accent           = Color(hex: "#a900ff")  // --color-accent

    // Semantic
    static let success          = Color(hex: "#22c55e")
    static let warning          = Color(hex: "#f59e0b")
    static let error            = Color(hex: "#ef4444")
    static let info             = Color(hex: "#3b82f6")

    // Syntax (code editor)
    static let syntaxBlue       = Color(hex: "#82aaff")
    static let syntaxCyan       = Color(hex: "#89ddff")
    static let syntaxGreen      = Color(hex: "#c3e88d")
    static let syntaxYellow     = Color(hex: "#ffcb6b")
    static let syntaxOrange     = Color(hex: "#f78c6c")
    static let syntaxRed        = Color(hex: "#ff5370")
    static let syntaxPurple     = Color(hex: "#c792ea")
}
```

### Typography

```swift
extension Font {
    static let bodyDefault      = Font.system(size: 14)
    static let bodySecondary    = Font.system(size: 13)
    static let label            = Font.system(size: 12, weight: .semibold)
    static let caption          = Font.system(size: 11)
    static let sectionTitle     = Font.system(size: 16, weight: .semibold)
    static let code             = Font.system(.body, design: .monospaced)
}
```

### Spacing & Geometry

| Element | Default | Compact | Modern |
|---------|---------|---------|--------|
| Sidebar width | 56pt | 48pt | 64pt |
| Sidebar item | 36x36pt | 32x32pt | 40x40pt |
| Topbar height | 44pt | 36pt | 48pt |
| Status bar height | 24pt | 20pt | 28pt |
| Corner radius (sm) | 6pt | 3pt | 4pt |
| Corner radius (md) | 10pt | 5pt | 8pt |
| Corner radius (lg) | 14pt | 8pt | 12pt |

### 25 Color Themes

All 25 themes from `@orchestra-mcp/theme` are supported. Default: `orchestra` (deep navy).

| ID | Background | Accent | Light? |
|----|-----------|--------|--------|
| `orchestra` | #0a0d14 | #a900ff | No |
| `dracula` | #282a36 | #bd93f9 | No |
| `github-dark` | #0d1117 | #58a6ff | No |
| `github-light` | #ffffff | #0366d6 | Yes |
| `one-dark` | #282c34 | #528bff | No |
| `monokai-pro` | #2d2a2e | #ffd866 | No |
| `synthwave-84` | #262335 | #ff7edb | No |

---

## 7. Window Management

### Three Window Modes (macOS Only)

| Mode | Window | Size | Behavior |
|------|--------|------|----------|
| **Embedded** | Main window | 1280x860 | Full IDE, all sections |
| **Floating** | Spirit window | 420x640 | Frameless, always-on-top, semi-transparent, mini chat |
| **Bubble** | Bubble window | 56x56 | Non-resizable, always-on-top, circular overlay |

**Cycling**: `Cmd+Shift+O` global hotkey cycles modes.

### Spirit Window (Floating Mini Chat)

```swift
class SpiritPanel: NSPanel {
    override var canBecomeKey: Bool { true }

    init() {
        super.init(contentRect: NSRect(x: 0, y: 0, width: 420, height: 640),
                   styleMask: [.titled, .closable, .resizable, .fullSizeContentView],
                   backing: .buffered, defer: false)

        titlebarAppearsTransparent = true
        titleVisibility = .hidden
        isMovableByWindowBackground = true
        level = .floating
        collectionBehavior = [.canJoinAllSpaces, .fullScreenAuxiliary, .ignoresCycle]
        backgroundColor = NSColor(red: 10/255, green: 10/255, blue: 20/255, alpha: 0.9)
        isOpaque = false
        hasShadow = true
        minSize = NSSize(width: 320, height: 400)
    }
}
```

---

## 8. QUIC Transport

### Network.framework QUIC Client

```swift
import Network

class QUICConnection: ObservableObject {
    @Published var state: ConnectionState = .disconnected
    private var connection: NWConnection?
    private let queue = DispatchQueue(label: "orchestra.quic")

    func connect(to host: String = "localhost", port: UInt16 = 50100, pluginID: String = "ui.swift") {
        let tlsOptions = try MTLSConfig.create(pluginID: pluginID)

        let quicOptions = NWProtocolQUIC.Options()
        quicOptions.alpn = ["orchestra-plugin"]

        let params = NWParameters(tls: tlsOptions)
        params.defaultProtocolStack.transportProtocol = quicOptions

        let endpoint = NWEndpoint.hostPort(host: .init(host), port: .init(rawValue: port)!)
        connection = NWConnection(to: endpoint, using: params)
        connection?.stateUpdateHandler = handleStateUpdate
        connection?.start(queue: queue)
    }
}
```

### Length-Delimited Protobuf Framing

```swift
struct StreamFramer {
    static let maxMessageSize: UInt32 = 16 * 1024 * 1024  // 16 MB

    static func write(_ data: Data, to connection: NWConnection) async throws {
        var length = UInt32(data.count).bigEndian
        let header = Data(bytes: &length, count: 4)
        try await connection.send(content: header + data)
    }

    static func read(from connection: NWConnection) async throws -> Data {
        let header = try await connection.receive(length: 4)
        let size = header.withUnsafeBytes { $0.load(as: UInt32.self).bigEndian }
        guard size <= maxMessageSize else { throw TransportError.messageTooLarge(size) }
        return try await connection.receive(length: Int(size))
    }
}
```

### mTLS Configuration

```swift
struct MTLSConfig {
    static func create(pluginID: String) throws -> NWProtocolTLS.Options {
        let certsDir: URL
        #if os(macOS)
        certsDir = FileManager.default.homeDirectoryForCurrentUser
            .appendingPathComponent(".orchestra/certs")
        #else
        certsDir = Bundle.main.bundleURL.appendingPathComponent("Certs")
        #endif

        let caCert = try loadCertificate(from: certsDir.appendingPathComponent("ca.crt"))
        let pluginCert = try loadCertificate(from: certsDir.appendingPathComponent("\(pluginID).crt"))
        let pluginKey = try loadPrivateKey(from: certsDir.appendingPathComponent("\(pluginID).key"))

        let tlsOptions = NWProtocolTLS.Options()
        // Configure mTLS identity and CA verification
        return tlsOptions
    }
}
```

---

## 9. Data Layer

### MCP Tool Proxy

All data operations go through MCP tool calls routed via the orchestrator.

```swift
class ToolService {
    private let client: OrchestraClient

    func callTool(name: String, arguments: [String: Any]) async throws -> ToolResponse {
        return try await client.sendToolCall(name: name, arguments: arguments)
    }

    // Convenience wrappers
    func listProjects() async throws -> [Project] { ... }
    func getProjectStatus(slug: String) async throws -> ProjectStatus { ... }
    func createFeature(project: String, title: String, ...) async throws -> Feature { ... }
    func advanceFeature(id: String, evidence: String) async throws -> Feature { ... }

    // Multi-LLM AI calls
    func aiPrompt(prompt: String, provider: String, model: String) async throws -> String { ... }
    func spawnSession(id: String, prompt: String, provider: String) async throws -> ChatResponse { ... }
}
```

### Local Cache

```swift
// ~/Library/Application Support/Orchestra/cache.db (macOS)
// App container (iOS)
struct LocalCache {
    func cacheProjects(_ projects: [Project])
    func cachedProjects() -> [Project]
    func cacheNotes(_ notes: [Note])
    func cachedNotes() -> [Note]
    func cacheSessions(_ sessions: [ChatSession])
    func invalidateCache(for entity: String)
}
```

---

## 10. AI Chat

### Multi-LLM Chat Architecture

The chat interface supports 4 AI providers (Claude, OpenAI, Gemini, Ollama) via the bridge plugins.

```
ChatView
├── ChatSessionList (280px sidebar)
│   ├── Search bar
│   ├── New Chat button (+)
│   └── Session list (name, provider badge, model badge, date)
└── ChatBox
    ├── ChatHeader (session name, provider, model, Live dot)
    ├── ChatBody (scrollable messages)
    │   ├── ChatMessageView (user: right, purple)
    │   ├── ChatMessageView (assistant: left, surface)
    │   ├── EventCards (tool call results)
    │   └── TypingIndicator
    ├── StatusLine (typing status + elapsed timer)
    └── ChatInput
        ├── TextArea (auto-resize, 1-6 lines)
        └── Tray (provider, model, tools, attach, send/stop)
```

### Provider & Model Selection

```swift
struct ModelSelector: View {
    let providers: [AIProvider] = [
        AIProvider(id: "claude", name: "Anthropic", models: [
            "claude-opus-4-6", "claude-sonnet-4-6", "claude-haiku-4-5"
        ]),
        AIProvider(id: "openai", name: "OpenAI", models: [
            "gpt-4o", "gpt-4o-mini", "o1", "o1-mini"
        ]),
        AIProvider(id: "gemini", name: "Google", models: [
            "gemini-2.5-pro", "gemini-2.5-flash", "gemini-2.0-flash"
        ]),
        AIProvider(id: "ollama", name: "Ollama", models: [
            "llama3", "codellama", "mistral"
        ]),
    ]
}
```

### Message Types

```swift
struct ChatMessage: Identifiable {
    let id: String
    let role: MessageRole      // .user, .assistant, .system
    let content: String
    let timestamp: Date
    var streaming: Bool = false
    var thinking: String?
    var events: [ToolEvent] = []
    var provider: String?      // "claude", "openai", "gemini", "ollama"
    var model: String?
    var attachments: [Attachment] = []
}
```

### Streaming Glow Effect

When assistant is streaming, the bubble gets a rotating gradient border:

```swift
.overlay(
    RoundedRectangle(cornerRadius: 16)
        .strokeBorder(
            AngularGradient(
                gradient: Gradient(colors: [.blue, .purple, .red, .pink, .blue]),
                center: .center,
                angle: .degrees(glowAngle)
            ),
            lineWidth: 2
        )
        .opacity(message.streaming ? 1 : 0)
)
```

### Event Cards (Tool Results)

| Card Type | Renders |
|-----------|---------|
| `BashCard` | Terminal output with command + exit code |
| `ReadCard` | File content with line numbers |
| `EditCard` | Diff view (old/new) |
| `TaskCard` | Task summary with status/priority badges |
| `ProjectStatusCard` | Progress bar + stats |
| `SprintCard` | Sprint overview |
| `QuestionCard` | Interactive question with option buttons |

---

## 11. Project Management

### Projects Section

```
ProjectsView
├── ProjectSidebar (280px)
│   ├── Search bar
│   ├── New Project button (+)
│   └── Project list (icon, name, task count, %, Active badge)
└── ProjectDetail
    ├── Header (icon, name, description)
    ├── Status card (Total Tasks, Completed, Completion %)
    │   └── Progress bar (purple, animated)
    ├── Status breakdown
    └── BacklogTree (collapsible epics > stories > tasks)
```

### Workflow States (13-State Machine)

```swift
enum WorkflowState: String, CaseIterable {
    case backlog, todo, inProgress, readyForTesting, inTesting
    case readyForDocs, inDocs, documented, inReview, done
    case blocked, rejected, cancelled

    var color: Color {
        switch self {
        case .backlog, .cancelled: return Color(hex: "#6b7280")
        case .todo: return Color(hex: "#3b82f6")
        case .inProgress: return .accent
        case .readyForTesting: return Color(hex: "#eab308")
        case .inTesting: return Color(hex: "#f97316")
        case .readyForDocs: return Color(hex: "#06b6d4")
        case .inDocs: return Color(hex: "#14b8a6")
        case .documented, .done: return .success
        case .inReview: return Color(hex: "#6366f1")
        case .blocked, .rejected: return .error
        }
    }
}
```

---

## 12. Notes & Docs

### Notes Plugin

Uses `tools.notes` plugin (8 tools): `create_note`, `get_note`, `update_note`, `delete_note`, `list_notes`, `search_notes`, `pin_note`, `tag_note`.

```
NotesView
├── NotesSidebar (280px)
│   ├── Search bar
│   ├── New Note button (+)
│   ├── Pinned Notes section
│   └── Other Notes section
└── NoteEditor
    ├── Toolbar (back, pin toggle, save, delete)
    ├── Title input (large, borderless)
    ├── Tags bar (add/remove tags)
    └── Content editor (monospace, markdown)
```

### Docs/Wiki Plugin

Uses `tools.docs` plugin (10 tools): `doc_create`, `doc_get`, `doc_update`, `doc_delete`, `doc_list`, `doc_search`, `doc_generate`, `doc_index`, `doc_tree`, `doc_export`.

Categories: `api-reference`, `guide`, `architecture`, `tutorial`, `changelog`, `decision-record`

### Markdown Parser Plugin

Uses `tools.markdown` plugin (8 tools): `md_parse`, `md_parse_file`, `md_parse_frontmatter`, `md_render_html`, `md_render_plaintext`, `md_toc`, `md_lint`, `md_transform`.

Supports: Mermaid, KaTeX, footnotes, embeds, task lists, frontmatter, heading anchors.

---

## 13. Developer Tools

### DevTools Plugin Container

The DevTools section hosts sub-plugins for each tool type. Each is a separate Go plugin on the backend, exposed through the DevToolsPlugin container.

| Tool | Plugin ID | Tools | Description |
|------|-----------|-------|-------------|
| File Explorer | devtools.file-explorer | 17 | Full IDE: file ops + LSP code intelligence |
| Terminal | devtools.terminal | 6 | PTY terminal sessions |
| SSH | devtools.ssh | 7 | Remote SSH + SFTP |
| Services | devtools.services | 6 | Service manager (launchctl/systemctl) |
| Docker | devtools.docker | 10 | Container management |
| Debugger | devtools.debugger | 8 | DAP protocol debugger |
| Test Runner | devtools.test-runner | 6 | Multi-framework test runner |
| Log Viewer | devtools.log-viewer | 5 | Log streaming + search |
| Database | devtools.database | 8 | SQL query editor + schema browser |
| DevOps | devtools.devops | 8 | CI/CD pipeline management |

### File Explorer (Full IDE + LSP)

Not just a file browser — a full IDE backend with LSP code intelligence via the Rust engine:

**File Tools:** `list_directory`, `read_file`, `write_file`, `move_file`, `delete_file`, `file_info`, `file_search`

**Code Intelligence Tools:** `code_symbols`, `code_goto_definition`, `code_find_references`, `code_hover`, `code_complete`, `code_diagnostics`, `code_actions`, `code_workspace_symbols`, `code_namespace`, `code_imports`

---

## 14. AI Awareness

Four AI awareness plugins provide visual and contextual understanding.

### ai.screenshot (6 tools)
`capture_screen`, `capture_region`, `capture_window`, `capture_interactive`, `annotate_screenshot`, `list_captures`

macOS: ScreenCaptureKit. iOS: UIScreen snapshot.

### ai.vision (6 tools)
`analyze_image`, `extract_text`, `find_elements`, `compare_images`, `describe_screen`, `extract_data`

Uses Claude Vision API or OpenAI Vision as fallback.

### ai.browser-context (7 tools)
`get_page_content`, `get_page_dom`, `get_selected_text`, `get_open_tabs`, `get_page_screenshot`, `navigate_to`, `execute_script`

Communicates via WebSocket to Chrome extension.

### ai.screen-reader (6 tools)
`get_accessibility_tree`, `get_focused_element`, `find_element`, `get_element_hierarchy`, `list_windows`, `get_window_elements`

macOS: Accessibility API via AXUIElement.

---

## 15. Voice & Notifications

### Voice Plugin (services.voice — 8 tools)
`tts_speak`, `tts_speak_provider`, `tts_list_voices`, `tts_stop`, `stt_listen`, `stt_transcribe_file`, `stt_list_models`, `voice_config`

**OS TTS:** macOS NSSpeechSynthesizer, iOS AVSpeechSynthesizer
**OS STT:** macOS SFSpeechRecognizer, iOS SFSpeechRecognizer
**Provider TTS:** ElevenLabs, OpenAI TTS, Google Cloud TTS
**Provider STT:** OpenAI Whisper, Google Cloud Speech, Deepgram

### Notifications Plugin (services.notifications — 8 tools)
`notify_send`, `notify_schedule`, `notify_cancel`, `notify_list_pending`, `notify_badge`, `notify_config`, `notify_history`, `notify_create_channel`

**Channels:** `build`, `test`, `deploy`, `ai`, `reminder`, `system`, `git`
**Actions:** Clickable buttons that trigger MCP tool calls.

---

## 16. Settings & Integrations

### Settings Navigation

| Section | Settings |
|---------|----------|
| **General** | Timezone, language |
| **Appearance** | Color theme (25 themes), component variant (3 options) |
| **Notifications** | Permission status, channels, DND hours |
| **Windows** | Default window mode, spirit bounds, bubble position (macOS) |
| **AI** | Default provider + model, auto-approve toggle |
| **Voice** | STT engine, TTS voice, language |
| **Sync & Account** | Sync status, connected devices, API tokens |

### AI Providers (Built-in Multi-LLM)

| Provider | Models | Auth |
|----------|--------|------|
| Anthropic | Claude Opus 4.6, Sonnet 4.6, Haiku 4.5 | Built-in / API key |
| OpenAI | GPT-4o, GPT-4o-mini, o1, o1-mini | API key |
| Google Gemini | Gemini 2.5 Pro, 2.5 Flash, 2.0 Flash | API key |
| Ollama | Llama 3, CodeLlama, Mistral, etc. | Local (no key) |

### Service Integrations

Discord, Slack, GitHub, Jira, Linear, Notion, Apple, Figma — each shows connection status and configuration.

---

## 17. Native Apple Features

### Global Hotkey (macOS)

```swift
NSEvent.addGlobalMonitorForEvents(matching: .keyDown) { event in
    if event.modifierFlags.contains([.command, .shift]) && event.keyCode == 31 {
        modeManager.cycleMode()
    }
}
```

### Permissions

```swift
struct PermissionsService {
    func requestNotificationPermission() async -> Bool
    func checkScreenRecordingPermission() -> Bool    // macOS only
    func checkMicrophonePermission() -> AVAuthorizationStatus
    func checkAccessibilityPermission() -> Bool      // macOS only
}
```

### Keychain (Credential Storage)

```swift
class KeychainService {
    func saveAPIKey(provider: String, key: String) throws
    func loadAPIKey(provider: String) throws -> String?
    func deleteAPIKey(provider: String) throws
}
```

### Auto-Updater (macOS)

```swift
class UpdaterService {
    let repo = "orchestra-mcp/orchestra-swift"
    let checkInterval: TimeInterval = 6 * 3600

    func checkForUpdate() async throws -> UpdateInfo?
    func applyUpdate(_ info: UpdateInfo) async throws
}
```

---

## 18. Build Phases

### Phase 1: Shell + Chat (MVP) — macOS + iOS

1. SPM project setup (`apps/swift/Package.swift`)
2. `OrchestraKit` Swift Package with QUIC transport + plugin system
3. Main window with plugin-driven sidebar navigation
4. System tray (macOS MenuBarExtra)
5. AI chat plugin: session list + multi-LLM conversation + streaming
6. Settings plugin: appearance (theme picker)
7. Connection status indicator

**Exit criteria**: Launch app, see plugin-driven sidebar, create chat session with any provider, see streaming response.

### Phase 2: Projects + Notes — macOS + iOS

1. Projects plugin: list, detail, backlog tree, workflow states
2. Notes plugin: list, editor, pin/unpin, icon/color picker
3. Search spotlight (Cmd+K on macOS)
4. Data caching in local SQLite
5. iOS: Tab bar + split view layouts

### Phase 3: Developer Tools — macOS + visionOS

1. File Explorer plugin (with LSP code intelligence)
2. Terminal plugin (PTY)
3. Additional DevTools plugins (Database, SSH, Log Viewer)
4. Component browser plugin

### Phase 4: Window Modes + Native Features — macOS

1. Spirit window (floating mini chat)
2. Bubble window (always-on-top overlay)
3. Global hotkey (Cmd+Shift+O)
4. Screenshot capture (ScreenCaptureKit)
5. macOS permissions flow
6. Auto-updater

### Phase 5: AI Awareness + Voice — macOS + iOS

1. AI Vision plugin (Claude Vision API)
2. Screenshot plugin (ScreenCaptureKit)
3. Voice STT/TTS plugin
4. Notifications plugin
5. Browser context plugin (Chrome extension bridge)

### Phase 6: Extended Platforms

1. watchOS: Complications, project status glance
2. tvOS: Dashboard with sprint burndown
3. visionOS: Spatial workspace, multi-window layout
4. CarPlay: Voice-first AI assistant, project status, hands-free chat
5. WidgetKit extensions (macOS + iOS)
6. Siri Shortcuts

---

## Appendix A: MCP Tools for Swift App

Key MCP tools the app calls (via orchestrator):

| Category | Tools |
|----------|-------|
| **Projects** | `list_projects`, `get_project_status`, `create_project`, `get_progress` |
| **Features** | `list_features`, `create_feature`, `advance_feature`, `get_next_feature`, `set_current_feature`, `search_features`, `get_blocked_features`, `get_dependency_graph` |
| **Notes** | `create_note`, `get_note`, `update_note`, `delete_note`, `list_notes`, `search_notes`, `pin_note`, `tag_note` |
| **AI Chat** | `ai_prompt`, `spawn_session`, `kill_session`, `session_status`, `list_active` (per provider) |
| **Sessions** | `create_session`, `send_message`, `list_sessions`, `get_session`, `delete_session`, `pause_session` |
| **Agent Ops** | `create_account`, `list_accounts`, `get_account_env`, `set_budget`, `check_budget`, `report_usage` |
| **Memory** | `search_memory`, `get_context`, `save_observation`, `get_project_summary`, `list_memories` |
| **Search** | `search`, `search_symbols`, `index_directory`, `get_index_stats` |
| **Parse** | `parse_file`, `get_symbols`, `get_imports` |
| **Docs** | `doc_create`, `doc_get`, `doc_update`, `doc_list`, `doc_search`, `doc_generate` |
| **Markdown** | `md_parse`, `md_render_html`, `md_toc`, `md_lint` |
| **Marketplace** | `search_packs`, `recommend_packs`, `install_pack`, `list_packs` |
| **Multi-Agent** | `define_agent`, `run_workflow`, `list_runs`, `compare_providers` |

## Appendix B: Color Parity

| Mobile | Desktop Theme | Note |
|--------|--------------|------|
| `#0f0f0f` (bg) | `#0a0d14` (--color-bg) | Desktop is slightly blue-tinted |
| `#171717` (surface) | `#111520` (--color-bg-alt) | Desktop is blue-tinted |
| `#a900ff` (brandPurple) | `#a900ff` (--color-accent) | Same |
| `#f5f5f5` (textPrimary) | `#e8ecf4` (--color-fg) | Desktop slightly cooler |

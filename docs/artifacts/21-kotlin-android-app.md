# 21 — Kotlin Universal App: Architecture & Implementation Guide

> Comprehensive reference for building the Orchestra MCP universal Android app in Kotlin/Jetpack Compose.
> Supports Android Phone, Android Tablet, ChromeOS, Wear OS, Android TV, and Android Auto from a single codebase.
> Plugin-based architecture mirroring the Go framework — every feature is a plugin.
> Companion to artifact 18 (Swift Universal App), artifact 19 (Multi-Agent Orchestrator), and artifact 20 (Plugin Expansion).

---

## Table of Contents

1. [Vision & Context](#1-vision--context)
2. [Architecture](#2-architecture)
3. [App Structure](#3-app-structure)
4. [Kotlin Plugin System](#4-kotlin-plugin-system)
5. [Navigation & Screens](#5-navigation--screens)
6. [Design System](#6-design-system)
7. [Window Management](#7-window-management) ← ChromeOS freeform, Crostini bridge, keyboard shortcuts
8. [QUIC Transport](#8-quic-transport)
9. [Data Layer](#9-data-layer)
10. [AI Chat](#10-ai-chat)
11. [Project Management](#11-project-management)
12. [Notes & Docs](#12-notes--docs)
13. [Developer Tools](#13-developer-tools)
14. [AI Awareness](#14-ai-awareness)
15. [Voice & Notifications](#15-voice--notifications)
16. [Settings & Integrations](#16-settings--integrations)
17. [Native Android Features](#17-native-android-features)
18. [Build Phases](#18-build-phases)
- [Appendix A: MCP Tools](#appendix-a-mcp-tools-for-kotlin-app)
- [Appendix B: Dependency Catalog](#appendix-b-dependency-catalog-libsversionstoml)
- [Appendix C: Color Parity](#appendix-c-color-parity-android--desktop)
- [Appendix D: Swift ↔ Kotlin Mapping](#appendix-d-swift--kotlin-mapping)
- [Appendix E: Platform Decision Matrix](#appendix-e-platform-decision-matrix) ← ChromeOS decisions
- [Appendix F: ChromeOS Manifest](#appendix-f-chromeos-specific-manifest)

---

## 1. Vision & Context

### What We're Building

A universal Android app for Orchestra MCP — an AI-agentic IDE that manages projects, features, sprints, and AI chat sessions across all Android platforms. The app connects to the Orchestra plugin ecosystem via QUIC and exposes 270+ MCP tools through a native Jetpack Compose interface.

### Platform Matrix

| Platform | Experience | Key Features |
|----------|-----------|--------------|
| **Android Phone** | Project management + AI chat | Bottom nav, offline cache, push notifications, PiP AI chat |
| **Android Tablet** | Full IDE | Two-pane layout, keyboard shortcuts, multi-window, stylus support |
| **ChromeOS** | Full desktop IDE | Freeform windows, taskbar integration, keyboard/mouse first, Linux container bridge, file system access |
| **Wear OS** | Quick status + notifications | Tiles, complications, project progress glance, notification actions |
| **Android TV** | Dashboard + monitoring | Leanback, project dashboards, sprint burndown, real-time metrics |
| **Android Auto** | Voice-first AI assistant | Voice commands, project status read-aloud, hands-free AI chat |

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

Same sections as the Swift app — Chat, Projects, Notes, Developer Tools, Components, Integrations, Settings. See artifact 18 section 1 for full screenshot breakdown.

### Key UI Patterns

- **Dark theme** with deep navy/black backgrounds (Material You dark)
- **Purple accent** (#a900ff) for active states, brand elements, CTAs
- **Navigation rail** (80dp) on tablet, bottom nav on phone
- **List/detail** pattern within each section
- **Main content area** taking remaining space
- **Green "Live" indicators** for active connections
- **Badge counts** on tool results
- **Model badges** in purple pills

---

## 2. Architecture

### Plugin Host Pattern

The Kotlin app is a consumer in the star-topology orchestrator architecture. It connects to independently running plugins via QUIC — identical to the Swift app.

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

                                    orchestra-kotlin ← THIS APP
                                    (Universal Android)
```

### Communication Contract

- **Transport**: QUIC via `Netty QUIC` (io.netty.incubator:netty-incubator-codec-quic)
- **Auth**: mTLS — CA at `~/.orchestra/certs/ca.crt`, app cert signed by CA
- **Wire format**: Length-delimited Protobuf (4-byte big-endian uint32 length + Protobuf bytes)
- **Proto contract**: Generated Kotlin code from `protobuf-kotlin` + `buf`
- **Message routing**: All messages go through orchestrator — never direct plugin-to-plugin
- **Streaming**: StreamStart → StreamChunk* → StreamEnd (for AI chat, long-running ops)
- **Events**: Subscribe/Publish/EventDelivery (for real-time updates)
- **Tool calls**: Send `ToolRequest` → receive `ToolResponse` (with optional `provider` field for AI routing)

### Platform Availability

| Framework | Phone | Tablet | ChromeOS | Wear OS | TV | Auto |
|-----------|-------|--------|----------|---------|-----|------|
| Netty QUIC | API 24+ | API 24+ | API 24+ | API 30+ | API 24+ | API 28+ |
| Jetpack Compose | API 24+ | API 24+ | API 24+ | API 26+ | API 24+ | — |
| Compose for TV (Leanback) | — | — | — | — | API 24+ | — |
| Compose for Wear | — | — | — | API 26+ | — | — |
| Android Auto (Car App Library) | — | — | — | — | — | API 23+ |
| WorkManager | API 24+ | API 24+ | API 24+ | API 26+ | API 24+ | API 23+ |
| SpeechRecognizer (STT) | API 24+ | API 24+ | API 24+ | API 30+ | — | API 23+ |
| TextToSpeech (TTS) | API 24+ | API 24+ | API 24+ | API 26+ | API 24+ | API 23+ |
| Freeform multi-window | — | — | Yes | — | — | — |
| Linux container (Crostini) | — | — | Yes | — | — | — |
| Chrome extension bridge | — | — | Yes | — | — | — |

### Graceful Degradation

If the orchestrator is not running, the app should:
1. Show a "Not Connected" status in the top app bar
2. Allow browsing locally cached data
3. Retry connection with exponential backoff (1s → 30s max)
4. Auto-reconnect when orchestrator becomes available

---

## 3. App Structure

### Project Layout

```
apps/kotlin/                               # github.com/orchestra-mcp/orchestra-kotlin
├── build.gradle.kts                       # Root build config (version catalog)
├── gradle/
│   └── libs.versions.toml                 # Version catalog
├── settings.gradle.kts                    # Module declarations
├── orchestra-kit/                         # Shared SDK module
│   ├── build.gradle.kts
│   └── src/main/kotlin/dev/orchestra/kit/
│       ├── transport/
│       │   ├── QUICConnection.kt          # Netty QUIC client
│       │   ├── StreamFramer.kt            # Length-delimited Protobuf framing
│       │   └── MTLSConfig.kt             # mTLS cert loading (BouncyCastle)
│       ├── proto/
│       │   └── Messages.kt               # Proto message types (placeholder until buf gen)
│       ├── models/
│       │   ├── Project.kt                # Project model
│       │   ├── Feature.kt                # Feature/task model
│       │   ├── Note.kt                   # Note model
│       │   ├── ChatSession.kt            # Chat session model
│       │   └── ChatMessage.kt            # Chat message model
│       ├── plugins/
│       │   ├── OrchestraPlugin.kt        # Plugin interface + AppSection
│       │   └── PluginRegistry.kt         # Plugin registry + discovery
│       └── services/
│           ├── OrchestraClient.kt        # High-level orchestrator client
│           ├── ToolService.kt            # MCP tool call proxy
│           └── ConnectionState.kt        # Connection status sealed class
├── shared/                               # Shared UI module
│   ├── build.gradle.kts
│   └── src/main/kotlin/dev/orchestra/shared/
│       ├── app/
│       │   └── OrchestraContent.kt       # Root composable (plugin-driven, adaptive)
│       ├── plugins/                      # Built-in plugins (each is self-contained)
│       │   ├── chat/
│       │   │   ├── ChatPlugin.kt         # Plugin registration
│       │   │   ├── ChatScreen.kt         # Chat layout
│       │   │   ├── ChatSessionList.kt    # Session sidebar
│       │   │   └── ChatMessageItem.kt    # Message bubble
│       │   ├── projects/
│       │   │   ├── ProjectsPlugin.kt
│       │   │   ├── ProjectsScreen.kt
│       │   │   └── ProjectDetailScreen.kt
│       │   ├── notes/
│       │   │   ├── NotesPlugin.kt
│       │   │   └── NotesScreen.kt
│       │   ├── devtools/
│       │   │   ├── DevToolsPlugin.kt
│       │   │   └── DevToolsScreen.kt
│       │   └── settings/
│       │       ├── SettingsPlugin.kt
│       │       └── SettingsScreen.kt
│       ├── components/
│       │   ├── StatusBadge.kt
│       │   ├── ConnectionIndicator.kt
│       │   └── EmptyStateView.kt
│       └── theme/
│           ├── OrchestraTheme.kt         # Material 3 theme
│           ├── Color.kt                  # Color tokens
│           └── Type.kt                   # Typography definitions
├── app/                                  # Phone + Tablet app module
│   ├── build.gradle.kts
│   └── src/main/kotlin/dev/orchestra/app/
│       ├── OrchestraApplication.kt       # Application class + Hilt
│       ├── MainActivity.kt              # Single Activity, Compose host
│       └── di/                           # Hilt dependency injection
│           └── AppModule.kt
├── wear/                                 # Wear OS module
│   ├── build.gradle.kts
│   └── src/main/kotlin/dev/orchestra/wear/
│       ├── WearActivity.kt              # Wear OS entry
│       ├── WearApp.kt                   # Wear Compose root
│       ├── tiles/
│       │   └── ProjectStatusTile.kt     # Wear OS tile service
│       └── complications/
│           └── TaskCountComplication.kt  # Complication data source
├── tv/                                   # Android TV module
│   ├── build.gradle.kts
│   └── src/main/kotlin/dev/orchestra/tv/
│       ├── TvActivity.kt               # TV entry
│       ├── TvApp.kt                    # Compose for TV root
│       └── screens/
│           ├── DashboardScreen.kt       # Project dashboards
│           └── BurndownScreen.kt        # Sprint burndown
├── chromeos/                             # ChromeOS-specific module
│   ├── build.gradle.kts
│   └── src/main/kotlin/dev/orchestra/chromeos/
│       ├── ChromeOSCompat.kt            # Runtime detection (ARC feature flag)
│       ├── ChromeOSFileAccess.kt        # Crostini + MyFiles + Drive integration
│       ├── ChromeOSShortcuts.kt         # Desktop keyboard shortcut map
│       ├── CrostiniBridge.kt            # QUIC connect to Crostini orchestrator
│       └── ChromeExtensionBridge.kt     # WebSocket bridge to Chrome extension
├── auto/                                 # Android Auto module
│   ├── build.gradle.kts
│   └── src/main/kotlin/dev/orchestra/auto/
│       ├── OrchestraCarService.kt       # Car App Library service
│       ├── OrchestraSession.kt          # Car session
│       └── screens/
│           ├── VoiceChatScreen.kt       # Voice-first AI chat
│           └── ProjectStatusScreen.kt   # Project status read-aloud
├── scripts/
│   └── new-kotlin-plugin.sh             # Plugin creator script
└── README.md
```

### App Lifecycle (Phone/Tablet)

```kotlin
@HiltAndroidApp
class OrchestraApplication : Application() {
    @Inject lateinit var pluginRegistry: PluginRegistry
    @Inject lateinit var orchestraClient: OrchestraClient

    override fun onCreate() {
        super.onCreate()
        // Register built-in plugins
        pluginRegistry.register(ChatPlugin())
        pluginRegistry.register(ProjectsPlugin())
        pluginRegistry.register(NotesPlugin())
        pluginRegistry.register(DevToolsPlugin())
        pluginRegistry.register(SettingsPlugin())
    }
}

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            OrchestraTheme {
                OrchestraContent()
            }
        }
    }
}
```

### App Lifecycle (Wear OS)

```kotlin
class WearActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            OrchestraWearTheme {
                WearApp()
            }
        }
    }
}

@Composable
fun WearApp() {
    val pluginRegistry = remember { PluginRegistry() }
    pluginRegistry.register(ChatPlugin())    // Quick reply only
    pluginRegistry.register(ProjectsPlugin()) // Status only
    pluginRegistry.register(SettingsPlugin())

    SwipeDismissableNavHost(
        navController = rememberSwipeDismissableNavController(),
        startDestination = "home"
    ) {
        composable("home") { ProjectStatusGlance() }
        composable("chat") { QuickReplyChat() }
    }
}
```

### App Lifecycle (Android TV)

```kotlin
class TvActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            OrchestraTvTheme {
                TvApp()
            }
        }
    }
}

@Composable
fun TvApp() {
    val pluginRegistry = remember { PluginRegistry() }
    pluginRegistry.register(ProjectsPlugin())  // Dashboard
    pluginRegistry.register(ChatPlugin())       // View-only
    pluginRegistry.register(SettingsPlugin())

    // Compose for TV TabRow + content
    TvTabNavigation(pluginRegistry)
}
```

### App Lifecycle (ChromeOS)

ChromeOS runs the same `app` module — no separate entry point needed. The `MainActivity` detects the ARC runtime at startup and applies desktop-class configuration:

```kotlin
@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()

        // ChromeOS: apply freeform window config, keyboard shortcuts
        if (ChromeOSCompat.isChromeOS(this)) {
            configureForChromeOS()
            // Connect to Crostini orchestrator if available
            lifecycle.addObserver(CrostiniBridge(this))
        }

        setContent {
            OrchestraTheme {
                OrchestraContent(
                    // ChromeOS gets full DevTools, otherwise tablet/phone
                    includeDevTools = ChromeOSCompat.isDesktopMode(this)
                )
            }
        }
    }

    override fun onConfigurationChanged(newConfig: Configuration) {
        super.onConfigurationChanged(newConfig)
        // Handle ChromeOS window resize without recreating activity
    }
}
```

### App Lifecycle (Android Auto)

```kotlin
class OrchestraCarService : CarAppService() {
    override fun createHostValidator() = HostValidator.ALLOW_ALL_HOSTS_VALIDATOR

    override fun onCreateSession() = OrchestraSession()
}

class OrchestraSession : Session() {
    override fun onCreateScreen(intent: Intent): Screen {
        return VoiceChatScreen(carContext)
    }
}
```

### Startup Sequence

```
1. Hilt dependency injection initialization
2. Detect platform: phone / tablet / ChromeOS / Wear / TV / Auto
3. Load settings from DataStore (proto)
4. Initialize theme (Material You dynamic colors + Orchestra overrides)
5. Register built-in plugins with PluginRegistry (platform-filtered)
6. [ChromeOS only] Apply freeform window config, register keyboard shortcuts
7. Setup UI (adaptive layout based on window size class)
8. Connect to orchestrator via QUIC (background coroutine)
   - [ChromeOS] Try Crostini localhost first, then remote fallback
   - [Other]    Try configured host, then offline cache
9. Subscribe to events (feature.*, workflow.*, sprint.*, note.*)
10. Load cached data from Room DB (projects, sessions, notes)
11. Check for updates (6-hour interval, via Play Store or GitHub)
12. Request runtime permissions (first run only)
13. [ChromeOS only] Connect Chrome extension bridge (WebSocket)
14. Restore navigation state (SavedStateHandle)
```

---

## 4. Kotlin Plugin System

### Design Principles

The app mirrors the Go framework's plugin architecture. Every screen/feature is a Kotlin plugin that registers itself with the PluginRegistry. This enables:
- **Incremental development** — build one plugin at a time for small wins
- **Easy extensibility** — `new-kotlin-plugin.sh` scaffolds new plugins
- **Feature isolation** — each plugin is self-contained with its own screens and logic
- **Platform adaptation** — plugins declare which platforms they support

### Plugin Interface

```kotlin
/**
 * Section where the plugin appears in the UI.
 */
enum class AppSection {
    Sidebar,    // Main navigation (Chat, Projects, Notes)
    DevTools,   // Developer Tools sub-section
    Settings,   // Settings sub-section
}

/**
 * Platforms the plugin can run on.
 */
enum class Platform {
    Phone, Tablet, ChromeOS, WearOS, TV, Auto
}

/**
 * Every feature in the app implements this interface.
 */
interface OrchestraPlugin {
    val id: String
    val name: String
    val icon: ImageVector            // Material icon
    val section: AppSection
    val order: Int                   // Sort order within section
    val supportedPlatforms: Set<Platform>

    @Composable
    fun Content(modifier: Modifier)

    fun onActivate() {}
    fun onDeactivate() {}
}
```

### Plugin Registry

```kotlin
class PluginRegistry {
    private val _plugins = mutableStateListOf<OrchestraPlugin>()
    val plugins: List<OrchestraPlugin> get() = _plugins

    fun register(plugin: OrchestraPlugin) {
        _plugins.add(plugin)
        _plugins.sortBy { it.order }
    }

    fun plugin(id: String): OrchestraPlugin? =
        _plugins.firstOrNull { it.id == id }

    val sidebarPlugins: List<OrchestraPlugin>
        get() = _plugins.filter {
            it.section == AppSection.Sidebar && it.supportedPlatforms.contains(currentPlatform)
        }

    val devToolPlugins: List<OrchestraPlugin>
        get() = _plugins.filter {
            it.section == AppSection.DevTools && it.supportedPlatforms.contains(currentPlatform)
        }

    val settingsPlugins: List<OrchestraPlugin>
        get() = _plugins.filter {
            it.section == AppSection.Settings && it.supportedPlatforms.contains(currentPlatform)
        }
}
```

### Built-in Plugins (Phase 1)

| Plugin | ID | Section | Platforms |
|--------|----|---------|-----------|
| ChatPlugin | `chat` | Sidebar | All |
| ProjectsPlugin | `projects` | Sidebar | Phone, Tablet, ChromeOS, TV |
| NotesPlugin | `notes` | Sidebar | Phone, Tablet, ChromeOS |
| DevToolsPlugin | `devtools` | Sidebar | Tablet, ChromeOS |
| SettingsPlugin | `settings` | Sidebar | All |

### Plugin Creator Script

`scripts/new-kotlin-plugin.sh` generates:

```bash
./scripts/new-kotlin-plugin.sh my-feature sidebar
# Creates:
# shared/src/main/kotlin/dev/orchestra/shared/plugins/myfeature/
#   ├── MyFeaturePlugin.kt    # Plugin registration
#   └── MyFeatureScreen.kt    # Main composable
# Prints registration line to add to Application class
```

---

## 5. Navigation & Screens

### Adaptive Navigation (Window Size Classes)

```kotlin
@Composable
fun OrchestraContent() {
    val pluginRegistry = LocalPluginRegistry.current
    var selectedPlugin by rememberSaveable { mutableStateOf("chat") }
    val windowSizeClass = currentWindowAdaptiveInfo().windowSizeClass

    when {
        // Tablet / ChromeOS / Large screen — Navigation Rail + Detail
        windowSizeClass.windowWidthSizeClass == WindowWidthSizeClass.EXPANDED -> {
            Row(Modifier.fillMaxSize()) {
                NavigationRail(
                    selectedPlugin = selectedPlugin,
                    onSelect = { selectedPlugin = it },
                    plugins = pluginRegistry.sidebarPlugins
                )
                pluginRegistry.plugin(selectedPlugin)?.Content(Modifier.weight(1f))
            }
        }
        // Phone — Bottom Navigation
        else -> {
            Scaffold(
                bottomBar = {
                    NavigationBar {
                        pluginRegistry.sidebarPlugins.forEach { plugin ->
                            NavigationBarItem(
                                selected = selectedPlugin == plugin.id,
                                onClick = { selectedPlugin = plugin.id },
                                icon = { Icon(plugin.icon, plugin.name) },
                                label = { Text(plugin.name) }
                            )
                        }
                    }
                }
            ) { padding ->
                pluginRegistry.plugin(selectedPlugin)
                    ?.Content(Modifier.padding(padding))
            }
        }
    }
}
```

### Navigation Rail (Tablet — 80dp)

| # | Icon | Section | Route |
|---|------|---------|-------|
| 1 | `Icons.Default.Chat` | Chat | `/chat` |
| 2 | `Icons.Default.GridView` | Projects | `/projects` |
| 3 | `Icons.Default.StickyNote2` | Notes | `/notes` |
| 4 | `Icons.Default.Terminal` | Developer Tools | `/devtools` |
| — | *(spacer)* | | |
| 5 | `Icons.Default.Settings` | Settings | `/settings` |

**Active state**: Purple accent background, white icon.
**Inactive state**: Muted gray icon, transparent background.
**Brand logo**: 36x36dp purple rounded square at top.
**User avatar**: At bottom with dropdown menu.

### Cross-Platform Feature Matrix

| Feature | Phone | Tablet | ChromeOS | Wear OS | TV | Auto |
|---------|-------|--------|----------|---------|-----|------|
| Chat | Full | Full | Full | Quick reply | View | Voice |
| Projects | Full | Full | Full | Status tile | Dashboard | Status |
| Notes | Full | Full | Full | — | — | — |
| DevTools | — | Full (10 tools) | Full (10 tools) | — | — | — |
| PiP Chat | Yes | Yes | Yes | — | — | — |
| Multi-Window | — | Yes | Freeform | — | — | — |
| Screenshot | Yes | Yes | Yes | — | — | — |
| Voice STT/TTS | Yes | Yes | Yes | Yes | — | Yes |
| Widgets | Yes | Yes | — | Tiles + Complications | — | — |
| Keyboard Shortcuts | — | Yes | Full (desktop-grade) | — | D-pad | — |
| Push Notifications | Yes | Yes | Yes | Yes | — | Yes |
| Linux Terminal | — | — | Yes (Crostini) | — | — | — |
| Chrome Extension | — | — | Yes (bridge) | — | — | — |
| File System Access | — | — | Full (ChromeOS Files) | — | — | — |
| Taskbar Pin | — | — | Yes | — | — | — |

---

## 6. Design System

### Color Tokens

```kotlin
object OrchestraColors {
    // Backgrounds
    val Background       = Color(0xFF0A0D14)  // --color-bg
    val Surface          = Color(0xFF111520)  // --color-bg-alt
    val SurfaceContrast  = Color(0xFF080A10)  // --color-bg-contrast
    val SurfaceActive    = Color(0xFF1A1F2E)  // --color-bg-active
    val SurfaceSelection = Color(0xFFA900FF).copy(alpha = 0.15f)

    // Text
    val TextPrimary      = Color(0xFFE8ECF4)  // --color-fg
    val TextMuted        = Color(0xFF8892A8)  // --color-fg-muted
    val TextDim          = Color(0xFF4A5268)  // --color-fg-dim
    val TextBright       = Color(0xFFF8FAFC)  // --color-fg-bright

    // Structure
    val Border           = Color(0xFF1E2436)  // --color-border
    val Accent           = Color(0xFFA900FF)  // --color-accent

    // Semantic
    val Success          = Color(0xFF22C55E)
    val Warning          = Color(0xFFF59E0B)
    val Error            = Color(0xFFEF4444)
    val Info             = Color(0xFF3B82F6)

    // Syntax (code viewer)
    val SyntaxBlue       = Color(0xFF82AAFF)
    val SyntaxCyan       = Color(0xFF89DDFF)
    val SyntaxGreen      = Color(0xFFC3E88D)
    val SyntaxYellow     = Color(0xFFFFCB6B)
    val SyntaxOrange     = Color(0xFFF78C6C)
    val SyntaxRed        = Color(0xFFFF5370)
    val SyntaxPurple     = Color(0xFFC792EA)
}
```

### Material 3 Theme Integration

```kotlin
@Composable
fun OrchestraTheme(
    themeId: String = "orchestra",
    content: @Composable () -> Unit
) {
    val colors = themeColors(themeId)

    val colorScheme = darkColorScheme(
        primary = colors.accent,
        onPrimary = Color.White,
        primaryContainer = colors.accent.copy(alpha = 0.2f),
        secondary = colors.accent,
        background = colors.background,
        onBackground = colors.textPrimary,
        surface = colors.surface,
        onSurface = colors.textPrimary,
        surfaceVariant = colors.surfaceActive,
        onSurfaceVariant = colors.textMuted,
        outline = colors.border,
        error = colors.error,
        onError = Color.White,
    )

    MaterialTheme(
        colorScheme = colorScheme,
        typography = OrchestraTypography,
        shapes = OrchestraShapes,
        content = content
    )
}
```

### Typography

```kotlin
val OrchestraTypography = Typography(
    bodyLarge = TextStyle(fontSize = 16.sp, lineHeight = 24.sp),
    bodyMedium = TextStyle(fontSize = 14.sp, lineHeight = 20.sp),
    bodySmall = TextStyle(fontSize = 12.sp, lineHeight = 16.sp),
    labelLarge = TextStyle(fontSize = 14.sp, fontWeight = FontWeight.SemiBold),
    labelMedium = TextStyle(fontSize = 12.sp, fontWeight = FontWeight.SemiBold),
    labelSmall = TextStyle(fontSize = 11.sp, fontWeight = FontWeight.Medium),
    titleLarge = TextStyle(fontSize = 22.sp, fontWeight = FontWeight.Bold),
    titleMedium = TextStyle(fontSize = 16.sp, fontWeight = FontWeight.SemiBold),
    titleSmall = TextStyle(fontSize = 14.sp, fontWeight = FontWeight.SemiBold),
)

val CodeFont = FontFamily(Font(R.font.jetbrains_mono))
```

### Spacing & Geometry

| Element | Default (Phone) | Compact | Expanded (Tablet) |
|---------|-----------------|---------|---------------------|
| Nav rail width | — | — | 80dp |
| Bottom nav height | 80dp | 64dp | — |
| Nav rail item | — | — | 56x56dp |
| Top app bar height | 64dp | 56dp | 64dp |
| Status bar height | 24dp | 24dp | 24dp |
| Corner radius (sm) | 8dp | 4dp | 8dp |
| Corner radius (md) | 12dp | 8dp | 12dp |
| Corner radius (lg) | 16dp | 12dp | 16dp |

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

### Adaptive Window Modes

Unlike macOS's three window modes, Android leverages platform-native multi-window:

| Mode | Platform | Behavior |
|------|----------|----------|
| **Full screen** | Phone, Tablet | Standard single-activity mode |
| **Split screen** | Tablet, Foldable | Android multi-window split (two-pane) |
| **Picture-in-Picture** | Phone, Tablet, ChromeOS | Mini AI chat overlay (PiP) |
| **Freeform** | ChromeOS, Samsung DeX | Desktop-like resizable floating windows |
| **Pop-up** | Foldable | Flex mode with top/bottom panels |
| **Maximized** | ChromeOS | True maximized window, taskbar visible |
| **Tiled** | ChromeOS | Snap-to-half like Windows/macOS tiling |

### Picture-in-Picture (Mini Chat)

```kotlin
class MainActivity : ComponentActivity() {

    override fun onUserLeaveHint() {
        super.onUserLeaveHint()
        if (chatViewModel.isStreaming) {
            enterPictureInPictureMode(
                PictureInPictureParams.Builder()
                    .setAspectRatio(Rational(9, 16))
                    .setAutoEnterEnabled(true)
                    .setSeamlessResizeEnabled(true)
                    .build()
            )
        }
    }
}
```

### Foldable Support (Jetpack WindowManager)

```kotlin
@Composable
fun AdaptiveChatLayout() {
    val windowInfo = currentWindowAdaptiveInfo()
    val foldingFeature = windowInfo.windowPosture

    when {
        foldingFeature.isTableTop -> {
            // Top half: chat messages, Bottom half: input
            Column(Modifier.fillMaxSize()) {
                ChatMessages(Modifier.weight(1f))
                HorizontalDivider()
                ChatInput(Modifier.weight(1f))
            }
        }
        foldingFeature.isBook -> {
            // Left: session list, Right: chat
            Row(Modifier.fillMaxSize()) {
                ChatSessionList(Modifier.weight(1f))
                VerticalDivider()
                ChatBox(Modifier.weight(1f))
            }
        }
        else -> {
            // Standard layout
            ChatScreen()
        }
    }
}
```

### ChromeOS Desktop Mode

On ChromeOS, the app runs as a full desktop-class application with freeform windows, keyboard/mouse input, and file system integration.

```kotlin
/**
 * Detect ChromeOS at runtime.
 */
object ChromeOSCompat {
    fun isChromeOS(context: Context): Boolean =
        context.packageManager.hasSystemFeature("org.chromium.arc")

    fun isDesktopMode(context: Context): Boolean =
        isChromeOS(context) || isDeXMode(context)

    private fun isDeXMode(context: Context): Boolean {
        val config = context.resources.configuration
        return try {
            val semDesktopMode = config.javaClass.getField("SEM_DESKTOP_MODE_ENABLED")
            config.javaClass.getField("semDesktopModeEnabled").getInt(config) ==
                semDesktopMode.getInt(null)
        } catch (_: Exception) { false }
    }
}

/**
 * ChromeOS-specific window configuration.
 */
fun Activity.configureForChromeOS() {
    if (ChromeOSCompat.isChromeOS(this)) {
        // Request freeform window mode (resizable, like a desktop app)
        window.setDecorFitsSystemWindows(false)

        // Set initial desktop-friendly size
        val display = windowManager.defaultDisplay
        val metrics = DisplayMetrics().also { display.getMetrics(it) }
        val width = (metrics.widthPixels * 0.75).toInt()
        val height = (metrics.heightPixels * 0.85).toInt()
        window.setLayout(width, height)

        // Minimum window size for usability
        window.attributes = window.attributes.apply {
            // 800x600 minimum
            if (Build.VERSION.SDK_INT >= 24) {
                window.setSustainedPerformanceMode(true)
            }
        }
    }
}
```

### ChromeOS File System Integration

```kotlin
/**
 * Access ChromeOS file system — MyFiles, Google Drive, Linux files (Crostini).
 */
class ChromeOSFileAccess(private val context: Context) {
    private val storageManager = context.getSystemService<StorageManager>()

    /** Open ChromeOS file picker with all volumes visible. */
    fun openFilePicker(): Intent =
        Intent(Intent.ACTION_OPEN_DOCUMENT).apply {
            addCategory(Intent.CATEGORY_OPENABLE)
            type = "*/*"
            putExtra(Intent.EXTRA_ALLOW_MULTIPLE, true)
        }

    /** Check if Linux (Crostini) files are accessible. */
    fun hasLinuxFiles(): Boolean =
        File("/mnt/chromeos/LinuxFiles").exists() ||
            context.packageManager.hasSystemFeature("org.chromium.arc.crostini")

    /** Get path to Linux container home directory (for orchestrator access). */
    fun linuxHomePath(): String? =
        if (hasLinuxFiles()) "/mnt/chromeos/LinuxFiles" else null
}
```

### ChromeOS + Crostini Bridge

On ChromeOS, the orchestrator can run inside the Linux (Crostini) container while the Kotlin app runs as an Android app. They communicate over QUIC via `localhost` (Crostini and Android share the same network namespace).

```
┌─────────────────────────────────────────────────────┐
│                     ChromeOS                         │
│                                                      │
│  ┌──────────────────┐    ┌────────────────────────┐  │
│  │   Android (ARC)  │    │  Linux (Crostini)      │  │
│  │                  │    │                        │  │
│  │  Orchestra App   │◄──►│  Orchestrator (Go)     │  │
│  │  (Kotlin/Compose)│QUIC│  Plugins (Go/Rust)     │  │
│  │                  │    │  Engine RAG (Rust)      │  │
│  │  Port 50100 ─────┼────┤  Listening :50100      │  │
│  └──────────────────┘    └────────────────────────┘  │
│                                                      │
│  ┌──────────────────┐                                │
│  │  Chrome Browser  │                                │
│  │  Extension ◄─────┼── WebSocket bridge             │
│  └──────────────────┘                                │
└─────────────────────────────────────────────────────┘
```

### ChromeOS Keyboard Shortcuts

```kotlin
/**
 * Desktop-grade keyboard shortcuts for ChromeOS / physical keyboard.
 */
@Composable
fun ChromeOSShortcuts(
    onNewChat: () -> Unit,
    onSearch: () -> Unit,
    onToggleDevTools: () -> Unit,
    onSwitchTab: (Int) -> Unit,
) {
    val keyboardModifiers = Modifier
        .onPreviewKeyEvent { event ->
            if (event.type != KeyEventType.KeyDown) return@onPreviewKeyEvent false
            val ctrl = event.isCtrlPressed
            val shift = event.isShiftPressed

            when {
                ctrl && event.key == Key.N -> { onNewChat(); true }
                ctrl && event.key == Key.K -> { onSearch(); true }
                ctrl && event.key == Key.Backtick -> { onToggleDevTools(); true }
                ctrl && event.key == Key.One -> { onSwitchTab(0); true }
                ctrl && event.key == Key.Two -> { onSwitchTab(1); true }
                ctrl && event.key == Key.Three -> { onSwitchTab(2); true }
                ctrl && event.key == Key.Four -> { onSwitchTab(3); true }
                ctrl && shift && event.key == Key.P -> { onSearch(); true } // Command palette
                else -> false
            }
        }
}
```

### ChromeOS Taskbar Integration

```xml
<!-- AndroidManifest.xml — ChromeOS taskbar pinning + window defaults -->
<activity
    android:name=".MainActivity"
    android:resizeableActivity="true"
    android:supportsPictureInPicture="true"
    android:configChanges="screenSize|smallestScreenSize|screenLayout|orientation|keyboard|keyboardHidden">

    <!-- ChromeOS: request freeform window support -->
    <meta-data android:name="WindowManagerPreference:FreeformWindowSize"
        android:value="phone" />
    <meta-data android:name="WindowManagerPreference:FreeformWindowOrientation"
        android:value="landscape" />

    <!-- Launch in maximized mode on ChromeOS -->
    <layout
        android:defaultHeight="85%"
        android:defaultWidth="75%"
        android:gravity="center"
        android:minHeight="600dp"
        android:minWidth="800dp" />
</activity>
```

---

## 8. QUIC Transport

### Netty QUIC Client

```kotlin
class QUICConnection(
    private val scope: CoroutineScope
) {
    private val _state = MutableStateFlow<ConnectionState>(ConnectionState.Disconnected)
    val state: StateFlow<ConnectionState> = _state.asStateFlow()

    private var channel: QuicChannel? = null
    private val eventLoopGroup = NioEventLoopGroup(1)

    suspend fun connect(
        host: String = "localhost",
        port: Int = 50100,
        pluginID: String = "ui.kotlin"
    ) {
        _state.value = ConnectionState.Connecting

        try {
            val sslContext = MTLSConfig.create(pluginID)

            val codec = QuicClientCodecBuilder()
                .sslContext(sslContext)
                .applicationProtocol("orchestra-plugin")
                .build()

            val bootstrap = Bootstrap()
                .group(eventLoopGroup)
                .channel(NioDatagramChannel::class.java)
                .handler(codec)

            channel = QuicChannel.newBootstrap(bootstrap)
                .handler(OrchestraStreamHandler())
                .remoteAddress(InetSocketAddress(host, port))
                .connect()
                .await()

            _state.value = ConnectionState.Connected
            startReceiveLoop()
        } catch (e: Exception) {
            _state.value = ConnectionState.Error(e.message ?: "Connection failed")
            scheduleReconnect()
        }
    }

    private fun scheduleReconnect() {
        scope.launch {
            var delay = 1000L
            while (_state.value !is ConnectionState.Connected) {
                delay(delay)
                try { connect() } catch (_: Exception) {}
                delay = (delay * 2).coerceAtMost(30_000L)
            }
        }
    }
}
```

### Length-Delimited Protobuf Framing

```kotlin
object StreamFramer {
    private const val MAX_MESSAGE_SIZE = 16 * 1024 * 1024 // 16 MB

    suspend fun write(data: ByteArray, channel: QuicStreamChannel) {
        val header = ByteBuffer.allocate(4).apply {
            putInt(data.size)
            flip()
        }
        channel.writeAndFlush(
            Unpooled.wrappedBuffer(header.array(), data)
        ).await()
    }

    suspend fun read(channel: QuicStreamChannel): ByteArray {
        val header = channel.readBytes(4)
        val size = ByteBuffer.wrap(header).int
        require(size <= MAX_MESSAGE_SIZE) { "Message too large: $size" }
        return channel.readBytes(size)
    }
}
```

### mTLS Configuration

```kotlin
object MTLSConfig {
    fun create(pluginID: String): QuicSslContext {
        val certsDir = File(System.getProperty("user.home"), ".orchestra/certs")

        val caCert = loadCertificate(File(certsDir, "ca.crt"))
        val pluginCert = loadCertificate(File(certsDir, "$pluginID.crt"))
        val pluginKey = loadPrivateKey(File(certsDir, "$pluginID.key"))

        return QuicSslContextBuilder.forClient()
            .trustManager(caCert)
            .keyManager(pluginKey, pluginCert)
            .applicationProtocols("orchestra-plugin")
            .build()
    }

    // On Android, load from app-internal storage or bundled assets
    fun createForAndroid(context: Context, pluginID: String): QuicSslContext {
        val certsDir = File(context.filesDir, "orchestra/certs")
        // Same pattern but Android-friendly paths
        return QuicSslContextBuilder.forClient()
            .trustManager(File(certsDir, "ca.crt"))
            .keyManager(
                loadPrivateKey(File(certsDir, "$pluginID.key")),
                loadCertificate(File(certsDir, "$pluginID.crt"))
            )
            .applicationProtocols("orchestra-plugin")
            .build()
    }
}
```

---

## 9. Data Layer

### MCP Tool Proxy

All data operations go through MCP tool calls routed via the orchestrator.

```kotlin
class ToolService(
    private val client: OrchestraClient
) {
    suspend fun callTool(name: String, arguments: Map<String, Any>): ToolResponse {
        return client.sendToolCall(name = name, arguments = arguments)
    }

    // Convenience wrappers
    suspend fun listProjects(): List<Project> { ... }
    suspend fun getProjectStatus(slug: String): ProjectStatus { ... }
    suspend fun createFeature(project: String, title: String, ...): Feature { ... }
    suspend fun advanceFeature(id: String, evidence: String): Feature { ... }

    // Multi-LLM AI calls
    suspend fun aiPrompt(prompt: String, provider: String, model: String): String { ... }
    suspend fun spawnSession(id: String, prompt: String, provider: String): ChatResponse { ... }
}
```

### Local Cache (Room)

```kotlin
@Database(
    entities = [ProjectEntity::class, NoteEntity::class, SessionEntity::class, MessageEntity::class],
    version = 1
)
abstract class OrchestraDatabase : RoomDatabase() {
    abstract fun projectDao(): ProjectDao
    abstract fun noteDao(): NoteDao
    abstract fun sessionDao(): SessionDao
    abstract fun messageDao(): MessageDao
}

@Dao
interface ProjectDao {
    @Query("SELECT * FROM projects ORDER BY updated_at DESC")
    fun observeAll(): Flow<List<ProjectEntity>>

    @Upsert
    suspend fun upsert(projects: List<ProjectEntity>)

    @Query("DELETE FROM projects WHERE id = :id")
    suspend fun delete(id: String)
}
```

### Repository Pattern (Offline-First)

```kotlin
class ProjectRepository(
    private val toolService: ToolService,
    private val projectDao: ProjectDao,
    private val connectionState: StateFlow<ConnectionState>
) {
    fun observeProjects(): Flow<List<Project>> =
        projectDao.observeAll().map { entities ->
            entities.map { it.toDomain() }
        }

    suspend fun refreshProjects() {
        if (connectionState.value is ConnectionState.Connected) {
            val remote = toolService.listProjects()
            projectDao.upsert(remote.map { it.toEntity() })
        }
    }

    suspend fun createProject(name: String, description: String): Project {
        val project = toolService.createProject(name, description)
        projectDao.upsert(listOf(project.toEntity()))
        return project
    }
}
```

---

## 10. AI Chat

### Multi-LLM Chat Architecture

The chat interface supports 4 AI providers (Claude, OpenAI, Gemini, Ollama) via the bridge plugins.

```
ChatScreen
├── ChatSessionList (280dp sidebar on tablet, drawer on phone)
│   ├── Search bar
│   ├── New Chat FAB
│   └── Session list (name, provider badge, model badge, date)
└── ChatBox
    ├── ChatTopBar (session name, provider, model, Live dot)
    ├── ChatMessages (LazyColumn, scrollable)
    │   ├── ChatMessageItem (user: end-aligned, purple)
    │   ├── ChatMessageItem (assistant: start-aligned, surface)
    │   ├── EventCards (tool call results)
    │   └── TypingIndicator
    ├── StatusLine (typing status + elapsed timer)
    └── ChatInput
        ├── TextField (auto-resize, 1-6 lines)
        └── ActionRow (provider, model, tools, attach, send/stop)
```

### Provider & Model Selection

```kotlin
data class AIProvider(
    val id: String,
    val name: String,
    val models: List<String>
)

val providers = listOf(
    AIProvider("claude", "Anthropic", listOf(
        "claude-opus-4-6", "claude-sonnet-4-6", "claude-haiku-4-5"
    )),
    AIProvider("openai", "OpenAI", listOf(
        "gpt-4o", "gpt-4o-mini", "o1", "o1-mini"
    )),
    AIProvider("gemini", "Google", listOf(
        "gemini-2.5-pro", "gemini-2.5-flash", "gemini-2.0-flash"
    )),
    AIProvider("ollama", "Ollama", listOf(
        "llama3", "codellama", "mistral"
    )),
)
```

### Message Types

```kotlin
data class ChatMessage(
    val id: String,
    val role: MessageRole,          // User, Assistant, System
    val content: String,
    val timestamp: Instant,
    val streaming: Boolean = false,
    val thinking: String? = null,
    val events: List<ToolEvent> = emptyList(),
    val provider: String? = null,   // "claude", "openai", "gemini", "ollama"
    val model: String? = null,
    val attachments: List<Attachment> = emptyList(),
)

enum class MessageRole { User, Assistant, System }
```

### Streaming Glow Effect

When the assistant is streaming, the bubble gets an animated gradient border:

```kotlin
@Composable
fun StreamingBorder(streaming: Boolean) {
    val infiniteTransition = rememberInfiniteTransition()
    val angle by infiniteTransition.animateFloat(
        initialValue = 0f,
        targetValue = 360f,
        animationSpec = infiniteRepeatable(
            animation = tween(2000, easing = LinearEasing)
        )
    )

    val brush = Brush.sweepGradient(
        colors = listOf(
            Color.Blue, Color(0xFFA900FF), Color.Red, Color(0xFFEC4899), Color.Blue
        ),
        center = Offset.Unspecified
    )

    Box(
        Modifier
            .then(if (streaming) Modifier.border(2.dp, brush, RoundedCornerShape(16.dp)) else Modifier)
    )
}
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
ProjectsScreen
├── ProjectSidebar (280dp on tablet, full screen list on phone)
│   ├── Search bar
│   ├── New Project FAB
│   └── Project list (icon, name, task count, %, Active badge)
└── ProjectDetail
    ├── Header (icon, name, description)
    ├── Status card (Total Tasks, Completed, Completion %)
    │   └── LinearProgressIndicator (purple, animated)
    ├── Status breakdown chips
    └── BacklogTree (expandable epics > stories > tasks)
```

### Workflow States (13-State Machine)

```kotlin
enum class WorkflowState(val label: String, val color: Color) {
    Backlog("Backlog", Color(0xFF6B7280)),
    Todo("To Do", Color(0xFF3B82F6)),
    InProgress("In Progress", OrchestraColors.Accent),
    ReadyForTesting("Ready for Testing", Color(0xFFEAB308)),
    InTesting("In Testing", Color(0xFFF97316)),
    ReadyForDocs("Ready for Docs", Color(0xFF06B6D4)),
    InDocs("In Docs", Color(0xFF14B8A6)),
    Documented("Documented", OrchestraColors.Success),
    InReview("In Review", Color(0xFF6366F1)),
    Done("Done", OrchestraColors.Success),
    Blocked("Blocked", OrchestraColors.Error),
    Rejected("Rejected", OrchestraColors.Error),
    Cancelled("Cancelled", Color(0xFF6B7280)),
}
```

---

## 12. Notes & Docs

### Notes Plugin

Uses `tools.notes` plugin (8 tools): `create_note`, `get_note`, `update_note`, `delete_note`, `list_notes`, `search_notes`, `pin_note`, `tag_note`.

```
NotesScreen
├── NotesSidebar (280dp on tablet, full screen list on phone)
│   ├── Search bar
│   ├── New Note FAB
│   ├── Pinned Notes section
│   └── Other Notes section
└── NoteEditor
    ├── TopAppBar (back, pin toggle, save, delete)
    ├── Title input (large, borderless)
    ├── Tags bar (FlowRow of chips, add/remove)
    └── Content editor (monospace, markdown)
```

### Docs/Wiki Plugin

Uses `tools.docs` plugin (10 tools): `doc_create`, `doc_get`, `doc_update`, `doc_delete`, `doc_list`, `doc_search`, `doc_generate`, `doc_index`, `doc_tree`, `doc_export`.

Categories: `api-reference`, `guide`, `architecture`, `tutorial`, `changelog`, `decision-record`

### Markdown Renderer

Uses `tools.markdown` plugin (8 tools) plus client-side rendering with `Markwon` library (Android's best markdown renderer):

```kotlin
// Markwon for rendering markdown in Compose
@Composable
fun MarkdownText(
    markdown: String,
    modifier: Modifier = Modifier
) {
    val markwon = remember {
        Markwon.builder(context)
            .usePlugin(SyntaxHighlightPlugin.create(Prism4jThemeDarkula.create()))
            .usePlugin(TablePlugin.create(context))
            .usePlugin(TaskListPlugin.create(context))
            .usePlugin(StrikethroughPlugin.create())
            .usePlugin(LatexPlugin.create())
            .build()
    }
    // Render to AnnotatedString for Compose
    AndroidView(
        factory = { ctx ->
            TextView(ctx).also { markwon.setMarkdown(it, markdown) }
        },
        modifier = modifier
    )
}
```

---

## 13. Developer Tools

### DevTools Plugin Container (Tablet + ChromeOS)

The DevTools section hosts sub-plugins for each tool type. Each is a separate Go plugin on the backend, exposed through the DevToolsPlugin container on tablets and ChromeOS.

| Tool | Plugin ID | Tools | Description |
|------|-----------|-------|-------------|
| File Explorer | devtools.file-explorer | 17 | File ops + code intelligence |
| Terminal | devtools.terminal | 6 | Terminal sessions (via QUIC) |
| SSH | devtools.ssh | 7 | Remote SSH + SFTP |
| Services | devtools.services | 6 | Service manager |
| Docker | devtools.docker | 10 | Container management |
| Debugger | devtools.debugger | 8 | DAP protocol debugger |
| Test Runner | devtools.test-runner | 6 | Multi-framework test runner |
| Log Viewer | devtools.log-viewer | 5 | Log streaming + search |
| Database | devtools.database | 8 | SQL query editor + schema browser |
| DevOps | devtools.devops | 8 | CI/CD pipeline management |

### File Explorer (Tablet + ChromeOS IDE)

On tablets and ChromeOS with keyboard/mouse (Samsung DeX, Chromebooks), the File Explorer becomes a full code editor. On ChromeOS, it additionally integrates with the Linux (Crostini) file system for direct access to project files:

**File Tools:** `list_directory`, `read_file`, `write_file`, `move_file`, `delete_file`, `file_info`, `file_search`

**Code Intelligence Tools:** `code_symbols`, `code_goto_definition`, `code_find_references`, `code_hover`, `code_complete`, `code_diagnostics`, `code_actions`, `code_workspace_symbols`, `code_namespace`, `code_imports`

---

## 14. AI Awareness

Four AI awareness plugins provide visual and contextual understanding.

### ai.screenshot (6 tools)
`capture_screen`, `capture_region`, `capture_window`, `capture_interactive`, `annotate_screenshot`, `list_captures`

Android: MediaProjection API for screen capture.

### ai.vision (6 tools)
`analyze_image`, `extract_text`, `find_elements`, `compare_images`, `describe_screen`, `extract_data`

Uses Claude Vision API or OpenAI Vision as fallback. Also CameraX for live capture.

### ai.browser-context (7 tools)
`get_page_content`, `get_page_dom`, `get_selected_text`, `get_open_tabs`, `get_page_screenshot`, `navigate_to`, `execute_script`

Communicates via Chrome Custom Tabs / Android accessibility service.

### ai.screen-reader (6 tools)
`get_accessibility_tree`, `get_focused_element`, `find_element`, `get_element_hierarchy`, `list_windows`, `get_window_elements`

Android: AccessibilityService API.

---

## 15. Voice & Notifications

### Voice Plugin (services.voice — 8 tools)
`tts_speak`, `tts_speak_provider`, `tts_list_voices`, `tts_stop`, `stt_listen`, `stt_transcribe_file`, `stt_list_models`, `voice_config`

**OS TTS:** Android `TextToSpeech` (all platforms)
**OS STT:** Android `SpeechRecognizer` (phone, tablet, auto)
**Provider TTS:** ElevenLabs, OpenAI TTS, Google Cloud TTS
**Provider STT:** OpenAI Whisper, Google Cloud Speech, Deepgram

### Notifications Plugin (services.notifications — 8 tools)
`notify_send`, `notify_schedule`, `notify_cancel`, `notify_list_pending`, `notify_badge`, `notify_config`, `notify_history`, `notify_create_channel`

**Channels:** `build`, `test`, `deploy`, `ai`, `reminder`, `system`, `git`
**Actions:** Clickable buttons that trigger MCP tool calls via pending intents.

```kotlin
object NotificationChannels {
    val channels = mapOf(
        "build" to NotificationChannel("build", "Build", NotificationManager.IMPORTANCE_DEFAULT),
        "test" to NotificationChannel("test", "Tests", NotificationManager.IMPORTANCE_DEFAULT),
        "deploy" to NotificationChannel("deploy", "Deploy", NotificationManager.IMPORTANCE_HIGH),
        "ai" to NotificationChannel("ai", "AI", NotificationManager.IMPORTANCE_DEFAULT),
        "reminder" to NotificationChannel("reminder", "Reminders", NotificationManager.IMPORTANCE_HIGH),
        "system" to NotificationChannel("system", "System", NotificationManager.IMPORTANCE_LOW),
        "git" to NotificationChannel("git", "Git", NotificationManager.IMPORTANCE_DEFAULT),
    )
}
```

---

## 16. Settings & Integrations

### Settings Navigation

| Section | Settings | ChromeOS only? |
|---------|----------|----------------|
| **General** | Timezone, language | No |
| **Appearance** | Color theme (25 themes), Material You toggle | No |
| **Notifications** | Permission status, channels, DND hours | No |
| **Display** | Default layout mode, font size, code font | No |
| **AI** | Default provider + model, auto-approve toggle | No |
| **Voice** | STT engine, TTS voice, language | No |
| **Sync & Account** | Sync status, connected devices, API tokens | No |
| **ChromeOS** | Crostini orchestrator path, Linux files toggle, window size preset, Chrome extension URL | Yes |
| **About** | Version, licenses, update check | No |

### AI Providers (Built-in Multi-LLM)

| Provider | Models | Auth |
|----------|--------|------|
| Anthropic | Claude Opus 4.6, Sonnet 4.6, Haiku 4.5 | Built-in / API key |
| OpenAI | GPT-4o, GPT-4o-mini, o1, o1-mini | API key |
| Google Gemini | Gemini 2.5 Pro, 2.5 Flash, 2.0 Flash | API key |
| Ollama | Llama 3, CodeLlama, Mistral, etc. | Local (no key) |

### Service Integrations

Discord, Slack, GitHub, Jira, Linear, Notion, Google, Figma — each shows connection status and configuration.

---

## 17. Native Android Features

### App Shortcuts

```xml
<shortcuts xmlns:android="http://schemas.android.com/apk/res/android">
    <shortcut android:shortcutId="new_chat"
        android:enabled="true"
        android:icon="@drawable/ic_chat"
        android:shortcutShortLabel="@string/new_chat"
        android:shortcutLongLabel="@string/new_chat_long">
        <intent
            android:action="dev.orchestra.NEW_CHAT"
            android:targetPackage="dev.orchestra.app"
            android:targetClass="dev.orchestra.app.MainActivity" />
    </shortcut>
    <shortcut android:shortcutId="projects"
        android:enabled="true"
        android:icon="@drawable/ic_projects"
        android:shortcutShortLabel="@string/projects">
        <intent
            android:action="dev.orchestra.PROJECTS"
            android:targetPackage="dev.orchestra.app"
            android:targetClass="dev.orchestra.app.MainActivity" />
    </shortcut>
</shortcuts>
```

### Permissions

```kotlin
object PermissionsHelper {
    val requiredPermissions = listOf(
        Manifest.permission.POST_NOTIFICATIONS,     // Android 13+
        Manifest.permission.RECORD_AUDIO,            // Voice STT
    )

    val optionalPermissions = listOf(
        Manifest.permission.CAMERA,                  // Screenshot/vision
        Manifest.permission.READ_EXTERNAL_STORAGE,   // File access
    )

    // MediaProjection for screen capture (runtime request)
    fun requestScreenCapture(activity: Activity) {
        val projectionManager = activity.getSystemService<MediaProjectionManager>()
        val intent = projectionManager.createScreenCaptureIntent()
        activity.startActivityForResult(intent, REQUEST_SCREEN_CAPTURE)
    }
}
```

### Credential Storage (EncryptedSharedPreferences)

```kotlin
class CredentialStore(context: Context) {
    private val prefs = EncryptedSharedPreferences.create(
        context,
        "orchestra_credentials",
        MasterKey.Builder(context)
            .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
            .build(),
        EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
        EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
    )

    fun saveAPIKey(provider: String, key: String) {
        prefs.edit { putString("api_key_$provider", key) }
    }

    fun loadAPIKey(provider: String): String? =
        prefs.getString("api_key_$provider", null)

    fun deleteAPIKey(provider: String) {
        prefs.edit { remove("api_key_$provider") }
    }
}
```

### Wear OS Tiles

```kotlin
class ProjectStatusTileService : TileService() {
    override fun onTileRequest(requestParams: RequestBuilders.TileRequest): ListenableFuture<Tile> {
        return Futures.immediateFuture(
            Tile.Builder()
                .setResourcesVersion("1")
                .setTileTimeline(
                    Timeline.Builder().addTimelineEntry(
                        TimelineEntry.Builder().setLayout(
                            Layout.Builder().setRoot(
                                Column.Builder()
                                    .addContent(projectStatusLayout())
                                    .addContent(taskCountLayout())
                                    .build()
                            ).build()
                        ).build()
                    ).build()
                ).build()
        )
    }
}
```

### Wear OS Complications

```kotlin
class TaskCountComplicationService : SuspendingComplicationDataSourceService() {
    override suspend fun onComplicationRequest(request: ComplicationRequest): ComplicationData {
        val taskCount = toolService.getActiveTaskCount()
        return ShortTextComplicationData.Builder(
            text = PlainComplicationText.Builder("$taskCount").build(),
            contentDescription = PlainComplicationText.Builder("$taskCount active tasks").build()
        )
            .setMonochromaticImage(
                MonochromaticImage.Builder(
                    Icon.createWithResource(this, R.drawable.ic_tasks)
                ).build()
            )
            .build()
    }
}
```

### Android TV Leanback

```kotlin
@Composable
fun TvDashboard(projects: List<Project>) {
    TvLazyColumn {
        item {
            Text(
                "Orchestra Dashboard",
                style = MaterialTheme.typography.headlineLarge,
                modifier = Modifier.padding(48.dp, 32.dp)
            )
        }
        item {
            TvLazyRow(
                contentPadding = PaddingValues(horizontal = 48.dp),
                horizontalArrangement = Arrangement.spacedBy(24.dp)
            ) {
                items(projects) { project ->
                    ProjectDashboardCard(
                        project = project,
                        modifier = Modifier.width(320.dp)
                    )
                }
            }
        }
        item {
            SprintBurndownChart(
                modifier = Modifier
                    .padding(48.dp)
                    .fillMaxWidth()
                    .height(300.dp)
            )
        }
    }
}
```

### Glance App Widgets (Home Screen)

```kotlin
class ProjectStatusWidget : GlanceAppWidget() {
    override suspend fun provideGlance(context: Context, id: GlanceId) {
        provideContent {
            OrchestraWidgetTheme {
                Column(
                    modifier = GlanceModifier
                        .fillMaxSize()
                        .background(OrchestraColors.Surface)
                        .padding(16.dp)
                        .cornerRadius(16.dp)
                ) {
                    Text(
                        "Active Project",
                        style = TextStyle(fontWeight = FontWeight.Bold, color = ColorProvider(OrchestraColors.TextPrimary))
                    )
                    Spacer(GlanceModifier.height(8.dp))
                    Text(
                        projectName,
                        style = TextStyle(color = ColorProvider(OrchestraColors.TextMuted))
                    )
                    Spacer(GlanceModifier.height(8.dp))
                    LinearProgressIndicator(
                        progress = completionPercent,
                        color = ColorProvider(OrchestraColors.Accent),
                        backgroundColor = ColorProvider(OrchestraColors.Border)
                    )
                    Spacer(GlanceModifier.height(4.dp))
                    Text(
                        "${(completionPercent * 100).toInt()}% complete",
                        style = TextStyle(fontSize = 12.sp, color = ColorProvider(OrchestraColors.TextDim))
                    )
                }
            }
        }
    }
}

class ProjectStatusWidgetReceiver : GlanceAppWidgetReceiver() {
    override val glanceAppWidget = ProjectStatusWidget()
}
```

### Auto-Updater

```kotlin
class UpdateChecker(
    private val context: Context,
    private val scope: CoroutineScope
) {
    private val repo = "orchestra-mcp/orchestra-kotlin"
    private val checkInterval = 6.hours

    fun startPeriodicCheck() {
        scope.launch {
            while (true) {
                delay(checkInterval)
                checkForUpdate()
            }
        }
    }

    private suspend fun checkForUpdate(): UpdateInfo? {
        // Check GitHub releases API or Play Store
        // Return update info if newer version available
        return null
    }
}
```

---

## 18. Build Phases

### Phase 1: Shell + Chat (MVP) — Phone + Tablet + ChromeOS

1. Multi-module Gradle project setup (`apps/kotlin/`)
2. `orchestra-kit` module with QUIC transport + plugin system
3. Adaptive layout with window size classes (phone: bottom nav, tablet/ChromeOS: nav rail)
4. AI chat plugin: session list + multi-LLM conversation + streaming
5. Settings plugin: appearance (theme picker, Material You toggle)
6. Connection status indicator
7. Hilt DI + Room DB + DataStore
8. ChromeOS: freeform window support, `resizeableActivity=true`, min size 800x600dp

**Exit criteria**: Launch app on phone, tablet, and Chromebook, see plugin-driven navigation, create chat session with any provider, see streaming response. On ChromeOS: resizable freeform window with nav rail.

### Phase 2: Projects + Notes — Phone + Tablet + ChromeOS

1. Projects plugin: list, detail, backlog tree, workflow states
2. Notes plugin: list, editor, pin/unpin, icon/color picker
3. Search (Ctrl+K on tablet/ChromeOS keyboards)
4. Data caching in Room DB
5. Offline-first repository pattern

### Phase 3: Developer Tools — Tablet + ChromeOS

1. File Explorer plugin (with code intelligence via engine.rag)
2. Terminal plugin (QUIC-tunneled)
3. Additional DevTools plugins (Database, SSH, Log Viewer)
4. Keyboard shortcut support for physical keyboards
5. ChromeOS: full keyboard shortcut map (Ctrl+N, Ctrl+K, Ctrl+`, Ctrl+1-4)

### Phase 4: Multi-Window + Native Features — All

1. Picture-in-Picture for AI chat streaming
2. Foldable device support (WindowManager)
3. Split-screen / multi-window on tablets
4. ChromeOS freeform windows, tiling, taskbar pinning
5. ChromeOS + Crostini bridge (connect to orchestrator in Linux container)
6. ChromeOS file system access (MyFiles, Google Drive, Linux files)
7. App shortcuts (long-press launcher icon)
8. Glance app widgets (home screen)
9. Runtime permissions flow
10. Samsung DeX desktop mode support

### Phase 5: AI Awareness + Voice — Phone + Tablet + ChromeOS

1. AI Vision plugin (Claude Vision API + CameraX)
2. Screenshot plugin (MediaProjection API)
3. Voice STT/TTS plugin
4. Notifications plugin (channels, actions)
5. Accessibility service integration
6. ChromeOS: Chrome extension bridge via WebSocket (browser context awareness)

### Phase 6: Extended Platforms

1. Wear OS: Tiles, complications, project status glance
2. Android TV: Leanback dashboard with sprint burndown
3. Android Auto: Voice-first AI chat, project status
4. Samsung Galaxy Watch integration
5. Google Assistant Actions

---

## Appendix A: MCP Tools for Kotlin App

Key MCP tools the app calls (via orchestrator) — identical to Swift app:

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

## Appendix B: Dependency Catalog (libs.versions.toml)

```toml
[versions]
kotlin = "2.1.0"
agp = "8.7.0"
compose-bom = "2025.02.00"
compose-compiler = "2.1.0"
hilt = "2.53"
room = "2.7.0"
datastore = "1.1.2"
netty-quic = "0.0.68.Final"
protobuf = "4.29.0"
protobuf-kotlin = "4.29.0"
coroutines = "1.10.0"
lifecycle = "2.8.7"
navigation = "2.8.5"
wear-compose = "1.5.0"
wear-tiles = "1.5.0"
horologist = "0.6.24"
tv-compose = "1.0.0-beta02"
car-app = "1.7.0"
glance = "1.1.1"
markwon = "4.6.2"
coil = "3.0.4"
bouncycastle = "1.78.1"

[libraries]
# Compose
compose-bom = { group = "androidx.compose", name = "compose-bom", version.ref = "compose-bom" }
compose-material3 = { group = "androidx.compose.material3", name = "material3" }
compose-material3-adaptive = { group = "androidx.compose.material3.adaptive", name = "adaptive" }
compose-material3-window = { group = "androidx.compose.material3", name = "material3-window-size-class" }
compose-ui-tooling = { group = "androidx.compose.ui", name = "ui-tooling" }

# Architecture
hilt-android = { group = "com.google.dagger", name = "hilt-android", version.ref = "hilt" }
hilt-compiler = { group = "com.google.dagger", name = "hilt-android-compiler", version.ref = "hilt" }
room-runtime = { group = "androidx.room", name = "room-runtime", version.ref = "room" }
room-ktx = { group = "androidx.room", name = "room-ktx", version.ref = "room" }
room-compiler = { group = "androidx.room", name = "room-compiler", version.ref = "room" }
datastore-proto = { group = "androidx.datastore", name = "datastore", version.ref = "datastore" }

# Transport
netty-quic = { group = "io.netty.incubator", name = "netty-incubator-codec-quic", version.ref = "netty-quic" }
protobuf-kotlin = { group = "com.google.protobuf", name = "protobuf-kotlin", version.ref = "protobuf-kotlin" }
bouncycastle = { group = "org.bouncycastle", name = "bcpkix-jdk18on", version.ref = "bouncycastle" }

# Lifecycle + Navigation
lifecycle-viewmodel = { group = "androidx.lifecycle", name = "lifecycle-viewmodel-compose", version.ref = "lifecycle" }
lifecycle-runtime = { group = "androidx.lifecycle", name = "lifecycle-runtime-compose", version.ref = "lifecycle" }
navigation-compose = { group = "androidx.navigation", name = "navigation-compose", version.ref = "navigation" }

# Wear OS
wear-compose-material = { group = "androidx.wear.compose", name = "compose-material3", version.ref = "wear-compose" }
wear-compose-navigation = { group = "androidx.wear.compose", name = "compose-navigation", version.ref = "wear-compose" }
wear-tiles = { group = "androidx.wear.tiles", name = "tiles", version.ref = "wear-tiles" }
horologist-compose-layout = { group = "com.google.android.horologist", name = "horologist-compose-layout", version.ref = "horologist" }

# TV
tv-compose-material = { group = "androidx.tv", name = "tv-material", version.ref = "tv-compose" }

# Auto
car-app = { group = "androidx.car.app", name = "app", version.ref = "car-app" }

# Widgets
glance = { group = "androidx.glance", name = "glance-appwidget", version.ref = "glance" }
glance-material = { group = "androidx.glance", name = "glance-material3", version.ref = "glance" }

# Content
markwon-core = { group = "io.noties.markwon", name = "core", version.ref = "markwon" }
markwon-syntax = { group = "io.noties.markwon", name = "syntax-highlight", version.ref = "markwon" }
markwon-tables = { group = "io.noties.markwon", name = "ext-tables", version.ref = "markwon" }
coil-compose = { group = "io.coil-kt.coil3", name = "coil-compose", version.ref = "coil" }

# Coroutines
coroutines-core = { group = "org.jetbrains.kotlinx", name = "kotlinx-coroutines-core", version.ref = "coroutines" }
coroutines-android = { group = "org.jetbrains.kotlinx", name = "kotlinx-coroutines-android", version.ref = "coroutines" }

[plugins]
android-application = { id = "com.android.application", version.ref = "agp" }
android-library = { id = "com.android.library", version.ref = "agp" }
kotlin-android = { id = "org.jetbrains.kotlin.android", version.ref = "kotlin" }
kotlin-compose = { id = "org.jetbrains.kotlin.plugin.compose", version.ref = "kotlin" }
hilt = { id = "com.google.dagger.hilt.android", version.ref = "hilt" }
ksp = { id = "com.google.devtools.ksp", version = "2.1.0-1.0.29" }
protobuf = { id = "com.google.protobuf", version = "0.9.4" }
```

## Appendix C: Color Parity (Android ↔ Desktop)

| Android | Desktop Theme | Note |
|---------|--------------|------|
| `0xFF0A0D14` (bg) | `#0a0d14` (--color-bg) | Exact match |
| `0xFF111520` (surface) | `#111520` (--color-bg-alt) | Exact match |
| `0xFFA900FF` (accent) | `#a900ff` (--color-accent) | Exact match |
| `0xFFE8ECF4` (textPrimary) | `#e8ecf4` (--color-fg) | Exact match |
| Material You dynamic | n/a | Android-exclusive, optional overlay |

## Appendix D: Swift ↔ Kotlin Mapping

| Swift Concept | Kotlin Equivalent |
|---------------|-------------------|
| SwiftUI | Jetpack Compose |
| `@StateObject` / `@ObservedObject` | `ViewModel` + `StateFlow` |
| `@EnvironmentObject` | Hilt `@Inject` / `CompositionLocal` |
| `@Published` | `MutableStateFlow` / `mutableStateOf` |
| `NavigationSplitView` | `NavigationRail` + adaptive layouts |
| `TabView` | `NavigationBar` (Material 3) |
| SF Symbols | Material Icons |
| `Network.framework` (QUIC) | Netty QUIC (incubator) |
| Core Data / SwiftData | Room DB |
| Keychain | EncryptedSharedPreferences |
| UserDefaults | DataStore (Proto) |
| WidgetKit | Glance App Widgets |
| Complications | Wear OS Complications |
| CarPlay | Android Auto (Car App Library) |
| SPM (Swift Package Manager) | Gradle multi-module |
| `async/await` | Kotlin Coroutines (`suspend`) |
| `AsyncSequence` | Kotlin `Flow` |
| Combine | Kotlin `Flow` |
| `@MainActor` | `Dispatchers.Main` |
| `DispatchQueue` | `CoroutineScope` / `Dispatchers` |

## Appendix E: Platform Decision Matrix

When a feature needs a platform-specific implementation, use this decision table:

| Question | Answer | Decision |
|----------|--------|----------|
| Does the user have a physical keyboard? | Phone: No / Tablet: Maybe / ChromeOS: Always / TV: No / Auto: No | Enable keyboard shortcuts only on Tablet + ChromeOS |
| Can the app run in a floating window? | ChromeOS: Yes (freeform) / Others: PiP only | Full multi-window only on ChromeOS |
| Can it access Linux filesystem? | ChromeOS + Crostini: Yes | Crostini bridge only on ChromeOS |
| Does it show a desktop taskbar? | ChromeOS: Yes | Taskbar pinning only on ChromeOS |
| Is there a Chrome browser on the same device? | ChromeOS: Yes (native) | Extension bridge only on ChromeOS (+ debug mode on Android) |
| Can it run the Go orchestrator locally? | ChromeOS (Crostini): Yes / Desktop (macOS/Linux/Windows): Yes / Android phone: No | Crostini auto-connect on ChromeOS |
| Does a home screen exist? | Phone + Tablet: Yes / ChromeOS: Launcher only | Glance widgets on Phone + Tablet only |
| Is it a mouse-first device? | ChromeOS: Yes / DeX: Yes / Others: No | Context menus, hover tooltips, right-click on ChromeOS + DeX |

## Appendix F: ChromeOS-Specific Manifest

```xml
<!-- Full AndroidManifest.xml additions for ChromeOS -->
<uses-feature android:name="android.hardware.touchscreen"
    android:required="false" />
<uses-feature android:name="android.hardware.faketouch"
    android:required="false" />

<!-- Declare keyboard/mouse input as non-required (ChromeOS provides them) -->
<uses-feature android:name="android.hardware.type.pc"
    android:required="false" />

<application
    android:resizeableActivity="true">

    <activity android:name=".MainActivity"
        android:resizeableActivity="true"
        android:supportsPictureInPicture="true"
        android:configChanges="screenSize|smallestScreenSize|screenLayout|orientation|keyboard|keyboardHidden|navigation">

        <!-- ChromeOS freeform window default size + min constraints -->
        <layout
            android:defaultHeight="85%"
            android:defaultWidth="75%"
            android:gravity="center"
            android:minHeight="600dp"
            android:minWidth="800dp" />

        <!-- ChromeOS window behavior hints -->
        <meta-data android:name="WindowManagerPreference:FreeformWindowSize"
            android:value="desktop" />
        <meta-data android:name="WindowManagerPreference:FreeformWindowOrientation"
            android:value="landscape" />
    </activity>
</application>
```

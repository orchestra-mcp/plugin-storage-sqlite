# 21 — Windows Desktop App: Architecture & Implementation Guide

> Comprehensive reference for building the Orchestra MCP universal Windows app in C#/WinUI 3.
> Supports Windows 10/11 Desktop, Xbox, HoloLens (Mixed Reality), and Windows IoT from a single codebase.
> Plugin-based architecture mirroring the Go framework — every feature is a plugin.
> Windows counterpart to artifact 18 (Swift Universal App). Compiled from the C# plugin agent, plugin expansion (artifact 20), multi-agent orchestrator (artifact 19), and architecture decisions.

---

## Table of Contents

1. [Vision & Context](#1-vision--context)
2. [Architecture](#2-architecture)
3. [App Structure](#3-app-structure)
4. [C# Plugin System](#4-c-plugin-system)
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
17. [Native Windows Features](#17-native-windows-features)
18. [Packaging & Distribution](#18-packaging--distribution)
19. [Build Phases](#19-build-phases)

---

## 1. Vision & Context

### What We're Building

A universal Windows app for Orchestra MCP — an AI-agentic IDE that manages projects, features, sprints, and AI chat sessions across the Microsoft platform family. The app connects to the Orchestra plugin ecosystem via QUIC and exposes 270+ MCP tools through a native WinUI 3 interface.

### Platform Matrix

| Platform | Experience | Key Features |
|----------|-----------|--------------|
| **Windows 10/11 Desktop** | Full IDE | Multi-window, system tray, global hotkeys, DevTools, Spirit floating chat, Windows.Graphics.Capture |
| **Xbox** | Dashboard + monitoring | Gamepad navigation, project dashboards, sprint burndown, real-time metrics display |
| **HoloLens / Mixed Reality** | Spatial workspace | Immersive code review, multi-window holographic layout, spatial AI chat |
| **Windows IoT** | Headless + dashboard | Kiosk mode, project dashboards, CI/CD monitoring, notification relay |

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

### The Old App (Reference)

The old Wails-based app had these sections (Windows will replicate all of them natively):

| Section | Icon | Description |
|---------|------|-------------|
| **Chats** | Chat bubble | AI chat sessions with Claude, multi-model, conversation history, tool call results, streaming |
| **Projects** | Grid | Project list with task counts, completion %, status badges |
| **Notes** | Edit | Note list with search, rich markdown editor |
| **Developer Tools** | Terminal | File Explorer, Terminal, Database, SSH, Log Viewer, Services, Debugger |
| **Components** | Layers | Component browser with preview |
| **Integrations** | Settings | AI Providers (7 providers), service integrations (Discord, Slack, GitHub, etc.) |
| **Settings** | Gear | General, Appearance, Notifications, Windows, AI, Voice, Sync & Account |

### Key UI Patterns

- **Dark theme** with deep navy/black backgrounds
- **Purple accent** (#A900FF) for active states, brand elements, CTAs
- **Left sidebar** with icon rail (56px) for section switching
- **Session/item sidebar** (280px) for lists within each section
- **Main content area** taking remaining space
- **User profile** at bottom-left with PersonPicture, name, email, status dot
- **Green "Live" indicators** for active connections
- **Badge counts** on tool results
- **Model badges** in purple pills
- **Fluent Design** (Mica/Acrylic) material throughout

### Apple-to-Windows Technology Map

| Apple (Artifact 18) | Windows (This Document) |
|----------------------|------------------------|
| SwiftUI | WinUI 3 (Windows App SDK) |
| Network.framework (QUIC) | System.Net.Quic (.NET 8+) |
| swift-protobuf | Google.Protobuf NuGet |
| SPM (Package.swift) | .NET Solution (.sln + .csproj) |
| Keychain | Windows Credential Manager |
| ScreenCaptureKit | Windows.Graphics.Capture |
| AXUIElement (Accessibility) | UI Automation (UIAutomationClient) |
| NSSpeechSynthesizer / SFSpeechRecognizer | Windows.Media.SpeechSynthesis / SpeechRecognition |
| UNUserNotificationCenter | Windows Toast Notifications (AppNotificationManager) |
| MenuBarExtra (system tray) | NotifyIcon (system tray via H.NotifyIcon.WinUI) |
| NSPanel (floating window) | Win32 WS_EX_TOPMOST + WS_EX_TOOLWINDOW |
| SF Symbols | Segoe Fluent Icons / Fluent System Icons |
| WidgetKit | Windows Widgets (Adaptive Cards) |
| Siri Shortcuts | Cortana / Windows Search integration |
| App Store / notarization | MSIX / Microsoft Store / winget |
| CarPlay | N/A (no automotive equivalent) |
| visionOS | HoloLens / Mixed Reality |
| watchOS | N/A (no wearable equivalent) |
| tvOS | Xbox |

---

## 2. Architecture

### Plugin Host Pattern

The Windows app is a consumer in the star-topology orchestrator architecture. It connects to independently running plugins via QUIC.

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

                                    orchestra-windows ← THIS APP
                                    (Universal Windows)
```

### Communication Contract

- **Transport**: QUIC via `System.Net.Quic` (.NET 8+ built-in, no third-party deps)
- **Auth**: mTLS — CA at `%USERPROFILE%\.orchestra\certs\ca.crt`, app cert signed by CA
- **Wire format**: Length-delimited Protobuf (4-byte big-endian uint32 length + Protobuf bytes)
- **Proto contract**: Generated C# code from `Google.Protobuf` (buf.build)
- **Message routing**: All messages go through orchestrator — never direct plugin-to-plugin
- **Streaming**: StreamStart → StreamChunk* → StreamEnd (for AI chat, long-running ops)
- **Events**: Subscribe/Publish/EventDelivery (for real-time updates)
- **Tool calls**: Send `ToolRequest` → receive `ToolResponse` (with optional `provider` field for AI routing)

### Platform Availability

| Framework | Windows 10 | Windows 11 | Xbox | HoloLens | IoT |
|-----------|-----------|-----------|------|----------|-----|
| System.Net.Quic (.NET 8+) | 19041+ | Yes | Yes | Yes | Yes |
| WinUI 3 (Windows App SDK) | 19041+ | Yes | Yes* | Yes* | No |
| Windows.Graphics.Capture | 19041+ | Yes | No | No | No |
| AppNotification (Toast) | 19041+ | Yes | Yes | Yes | No |
| Windows.Media.Speech | 19041+ | Yes | Yes | Yes | No |
| UI Automation | Yes | Yes | No | Yes | No |
| Adaptive Cards Widgets | No | 22H2+ | No | No | No |

*Xbox and HoloLens WinUI support requires Windows App SDK 1.5+.

### Graceful Degradation

If the orchestrator is not running, the app should:
1. Show a "Not Connected" status in the title bar
2. Allow browsing locally cached data
3. Retry connection with exponential backoff (1s → 30s max)
4. Auto-reconnect when orchestrator becomes available

---

## 3. App Structure

### Project Layout

```
apps/windows/                              # github.com/orchestra-mcp/orchestra-windows
├── Orchestra.sln                          # Solution file
├── src/
│   ├── Orchestra.Core/                    # Shared library (.NET 8 class library)
│   │   ├── Orchestra.Core.csproj
│   │   ├── Transport/
│   │   │   ├── QUICConnection.cs          # System.Net.Quic client
│   │   │   ├── StreamFramer.cs            # Length-delimited Protobuf framing
│   │   │   └── MTLSConfig.cs              # mTLS cert loading
│   │   ├── Proto/
│   │   │   └── Generated/                 # buf-generated C# Protobuf types
│   │   ├── Models/
│   │   │   ├── AppState.cs                # Observable root state
│   │   │   ├── Project.cs                 # Project model
│   │   │   ├── Feature.cs                 # Feature/task model
│   │   │   ├── Note.cs                    # Note model
│   │   │   ├── ChatSession.cs             # Chat session model
│   │   │   └── ChatMessage.cs             # Chat message model
│   │   ├── Plugins/
│   │   │   ├── IOrchestraPlugin.cs        # Plugin interface + AppSection
│   │   │   └── PluginRegistry.cs          # Plugin registry + discovery
│   │   └── Services/
│   │       ├── OrchestraClient.cs         # High-level orchestrator client
│   │       ├── ToolService.cs             # MCP tool call proxy
│   │       └── ConnectionState.cs         # Connection status enum
│   ├── Orchestra.Desktop/                 # WinUI 3 Desktop app
│   │   ├── Orchestra.Desktop.csproj
│   │   ├── App.xaml / App.xaml.cs         # Application entry point
│   │   ├── MainWindow.xaml / .cs          # Main window with NavigationView
│   │   ├── Plugins/                       # Built-in plugins (each is self-contained)
│   │   │   ├── ChatPlugin/
│   │   │   │   ├── ChatPlugin.cs          # Plugin registration
│   │   │   │   ├── ChatPage.xaml / .cs    # Chat layout
│   │   │   │   ├── ChatSessionList.xaml   # Session sidebar
│   │   │   │   └── ChatMessageControl.xaml # Message bubble
│   │   │   ├── ProjectsPlugin/
│   │   │   │   ├── ProjectsPlugin.cs
│   │   │   │   ├── ProjectsPage.xaml / .cs
│   │   │   │   └── ProjectDetailPage.xaml
│   │   │   ├── NotesPlugin/
│   │   │   │   ├── NotesPlugin.cs
│   │   │   │   └── NotesPage.xaml / .cs
│   │   │   ├── DevToolsPlugin/
│   │   │   │   ├── DevToolsPlugin.cs
│   │   │   │   └── DevToolsPage.xaml / .cs
│   │   │   └── SettingsPlugin/
│   │   │       ├── SettingsPlugin.cs
│   │   │       └── SettingsPage.xaml / .cs
│   │   ├── Controls/
│   │   │   ├── StatusBadge.xaml
│   │   │   ├── ConnectionIndicator.xaml
│   │   │   └── EmptyStateControl.xaml
│   │   ├── Windows/
│   │   │   ├── SpiritWindow.xaml / .cs    # Floating mini chat
│   │   │   └── BubbleWindow.xaml / .cs    # Always-on-top bubble
│   │   ├── Theme/
│   │   │   ├── OrchestraTheme.xaml        # Color tokens as ResourceDictionary
│   │   │   ├── ThemeManager.cs            # 25 theme switcher
│   │   │   └── Fonts.xaml                 # Typography definitions
│   │   ├── TrayIcon/
│   │   │   └── TrayIconService.cs         # System tray via H.NotifyIcon.WinUI
│   │   ├── Assets/
│   │   │   ├── orchestra-logo.png
│   │   │   └── Icons/
│   │   └── Package.appxmanifest           # MSIX manifest
│   ├── Orchestra.Widgets/                 # Windows Widget Provider
│   │   ├── Orchestra.Widgets.csproj
│   │   ├── ProjectWidget.cs              # Adaptive Cards widget
│   │   └── SprintWidget.cs
│   └── Orchestra.Xbox/                    # Xbox-specific (if applicable)
│       └── Orchestra.Xbox.csproj
├── tests/
│   ├── Orchestra.Core.Tests/
│   │   ├── Orchestra.Core.Tests.csproj
│   │   └── Transport/
│   └── Orchestra.Desktop.Tests/
│       └── Orchestra.Desktop.Tests.csproj
├── scripts/
│   └── new-windows-plugin.ps1             # Plugin creator script
└── README.md
```

### App Lifecycle (Windows Desktop)

```csharp
// App.xaml.cs
public partial class App : Application
{
    private readonly IHost _host;

    public App()
    {
        _host = Host.CreateDefaultBuilder()
            .ConfigureServices((_, services) =>
            {
                // Core services
                services.AddSingleton<AppState>();
                services.AddSingleton<PluginRegistry>();
                services.AddSingleton<QUICConnection>();
                services.AddSingleton<OrchestraClient>();
                services.AddSingleton<ToolService>();
                services.AddSingleton<ThemeManager>();
                services.AddSingleton<TrayIconService>();

                // Built-in plugins
                services.AddSingleton<IOrchestraPlugin, ChatPlugin>();
                services.AddSingleton<IOrchestraPlugin, ProjectsPlugin>();
                services.AddSingleton<IOrchestraPlugin, NotesPlugin>();
                services.AddSingleton<IOrchestraPlugin, DevToolsPlugin>();
                services.AddSingleton<IOrchestraPlugin, SettingsPlugin>();

                // Windows
                services.AddSingleton<MainWindow>();
                services.AddTransient<SpiritWindow>();
                services.AddTransient<BubbleWindow>();
            })
            .Build();

        InitializeComponent();
    }

    protected override void OnLaunched(LaunchActivatedEventArgs args)
    {
        var registry = _host.Services.GetRequiredService<PluginRegistry>();
        foreach (var plugin in _host.Services.GetServices<IOrchestraPlugin>())
        {
            registry.Register(plugin);
        }

        var mainWindow = _host.Services.GetRequiredService<MainWindow>();
        mainWindow.Activate();

        // Connect to orchestrator
        _ = _host.Services.GetRequiredService<QUICConnection>()
            .ConnectAsync("localhost", 50100, "ui.windows");
    }

    public static T GetService<T>() where T : class
        => (Current as App)?._host.Services.GetRequiredService<T>()!;
}
```

### Startup Sequence

```
1. Load settings from %LOCALAPPDATA%\Orchestra\settings.json
2. Initialize theme (color theme + Fluent Design variant)
3. Register built-in plugins with PluginRegistry via DI
4. Setup UI (main window with NavigationView)
5. Connect to orchestrator via QUIC (System.Net.Quic)
6. Subscribe to events (feature.*, workflow.*, sprint.*, note.*)
7. Load cached data (projects, sessions, notes) from local SQLite
8. Check for updates (5s delay, then every 6 hours)
9. Request platform permissions (first run only)
10. Restore window state (embedded/floating/bubble)
11. Show system tray icon
```

---

## 4. C# Plugin System

### Design Principles

The app mirrors the Go framework's plugin architecture. Every screen/feature is a C# plugin that registers itself with the PluginRegistry. This enables:
- **Incremental development** — build one plugin at a time for small wins
- **Easy extensibility** — `new-windows-plugin.ps1` scaffolds new plugins
- **Feature isolation** — each plugin is self-contained with its own pages and logic
- **Platform adaptation** — plugins declare which platforms they support

### Plugin Interface

```csharp
/// Section where the plugin appears in the UI.
public enum AppSection
{
    Sidebar,    // Main navigation (Chat, Projects, Notes)
    DevTools,   // Developer Tools sub-section
    Settings    // Settings sub-section
}

/// Windows platform targets.
[Flags]
public enum WindowsPlatform
{
    Desktop    = 1 << 0,   // Windows 10/11 Desktop
    Xbox       = 1 << 1,   // Xbox console
    HoloLens   = 1 << 2,   // HoloLens / Mixed Reality
    IoT        = 1 << 3,   // Windows IoT Core
    All        = Desktop | Xbox | HoloLens | IoT
}

/// Every feature in the app implements this interface.
public interface IOrchestraPlugin
{
    /// Unique plugin identifier (e.g., "chat", "projects", "notes").
    string Id { get; }

    /// Human-readable name shown in navigation.
    string Name { get; }

    /// Segoe Fluent Icon glyph (e.g., "\uE8BD" for Chat).
    string IconGlyph { get; }

    /// Which section of the UI this plugin belongs to.
    AppSection Section { get; }

    /// Sort order within the section (lower = higher in list).
    int Order { get; }

    /// Which Windows platforms this plugin supports.
    WindowsPlatform SupportedPlatforms { get; }

    /// Create the main page for this plugin.
    Type PageType { get; }

    /// Called when the plugin becomes the active section.
    void OnActivate();

    /// Called when the user navigates away from this plugin.
    void OnDeactivate();
}
```

### Abstract Base Class

```csharp
public abstract class OrchestraPluginBase : IOrchestraPlugin
{
    public abstract string Id { get; }
    public abstract string Name { get; }
    public abstract string IconGlyph { get; }
    public abstract AppSection Section { get; }
    public abstract int Order { get; }
    public abstract Type PageType { get; }

    public virtual WindowsPlatform SupportedPlatforms =>
        WindowsPlatform.Desktop | WindowsPlatform.HoloLens;

    public virtual void OnActivate() { }
    public virtual void OnDeactivate() { }
}
```

### Plugin Registry

```csharp
public sealed class PluginRegistry
{
    private readonly List<IOrchestraPlugin> _plugins = [];

    public IReadOnlyList<IOrchestraPlugin> Plugins => _plugins;

    public void Register(IOrchestraPlugin plugin)
    {
        _plugins.Add(plugin);
        _plugins.Sort((a, b) => a.Order.CompareTo(b.Order));
    }

    public IOrchestraPlugin? GetPlugin(string id) =>
        _plugins.FirstOrDefault(p => p.Id == id);

    public IEnumerable<IOrchestraPlugin> SidebarPlugins =>
        _plugins.Where(p => p.Section == AppSection.Sidebar && SupportsCurrentPlatform(p));

    public IEnumerable<IOrchestraPlugin> DevToolPlugins =>
        _plugins.Where(p => p.Section == AppSection.DevTools && SupportsCurrentPlatform(p));

    public IEnumerable<IOrchestraPlugin> SettingsPlugins =>
        _plugins.Where(p => p.Section == AppSection.Settings && SupportsCurrentPlatform(p));

    private static bool SupportsCurrentPlatform(IOrchestraPlugin plugin)
    {
        var current = GetCurrentPlatform();
        return plugin.SupportedPlatforms.HasFlag(current);
    }

    private static WindowsPlatform GetCurrentPlatform()
    {
        // Detect at runtime via AnalyticsInfo.VersionInfo.DeviceFamily
        var family = Windows.System.Profile.AnalyticsInfo.VersionInfo.DeviceFamily;
        return family switch
        {
            "Windows.Desktop" => WindowsPlatform.Desktop,
            "Windows.Xbox" => WindowsPlatform.Xbox,
            "Windows.Holographic" => WindowsPlatform.HoloLens,
            "Windows.IoT" => WindowsPlatform.IoT,
            _ => WindowsPlatform.Desktop
        };
    }
}
```

### Built-in Plugins (Phase 1)

| Plugin | ID | Section | Platforms |
|--------|----|---------|-----------|
| ChatPlugin | `chat` | Sidebar | All |
| ProjectsPlugin | `projects` | Sidebar | Desktop, Xbox, HoloLens |
| NotesPlugin | `notes` | Sidebar | Desktop, HoloLens |
| DevToolsPlugin | `devtools` | Sidebar | Desktop, HoloLens |
| SettingsPlugin | `settings` | Sidebar | All |

### Example Plugin Registration

```csharp
public sealed class ChatPlugin : OrchestraPluginBase
{
    public override string Id => "chat";
    public override string Name => "Chat";
    public override string IconGlyph => "\uE8BD";    // Segoe Fluent: Chat
    public override AppSection Section => AppSection.Sidebar;
    public override int Order => 1;
    public override Type PageType => typeof(ChatPage);
    public override WindowsPlatform SupportedPlatforms => WindowsPlatform.All;
}
```

### Plugin Creator Script

`scripts/new-windows-plugin.ps1` generates:

```powershell
.\scripts\new-windows-plugin.ps1 -Name "MyFeature" -Section "Sidebar"
# Creates:
# src/Orchestra.Desktop/Plugins/MyFeaturePlugin/
#   ├── MyFeaturePlugin.cs     # Plugin registration
#   └── MyFeaturePage.xaml/cs  # Main page
# Prints DI registration line to add to App.xaml.cs
```

---

## 5. Navigation & Screens

### Platform-Adaptive Navigation

```csharp
// MainWindow.xaml.cs
public sealed partial class MainWindow : Window
{
    private readonly PluginRegistry _registry;

    public MainWindow(PluginRegistry registry, AppState appState)
    {
        _registry = registry;
        InitializeComponent();

        // Build navigation items from plugins
        foreach (var plugin in _registry.SidebarPlugins)
        {
            NavView.MenuItems.Add(new NavigationViewItem
            {
                Content = plugin.Name,
                Icon = new FontIcon { Glyph = plugin.IconGlyph },
                Tag = plugin.Id
            });
        }

        NavView.SelectedItem = NavView.MenuItems[0];
    }

    private void NavView_SelectionChanged(NavigationView sender,
        NavigationViewSelectionChangedEventArgs args)
    {
        if (args.SelectedItem is NavigationViewItem item &&
            item.Tag is string pluginId)
        {
            var plugin = _registry.GetPlugin(pluginId);
            if (plugin != null)
            {
                ContentFrame.Navigate(plugin.PageType);
                plugin.OnActivate();
            }
        }
    }
}
```

### MainWindow XAML (NavigationView with Icon Rail)

```xml
<!-- MainWindow.xaml -->
<Window x:Class="Orchestra.Desktop.MainWindow"
        xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"
        xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml"
        Title="Orchestra" Width="1280" Height="860">
    <Grid>
        <NavigationView x:Name="NavView"
                        PaneDisplayMode="LeftCompact"
                        IsBackButtonVisible="Collapsed"
                        IsSettingsVisible="False"
                        SelectionChanged="NavView_SelectionChanged"
                        OpenPaneLength="280"
                        CompactPaneLength="56">

            <NavigationView.PaneHeader>
                <StackPanel Orientation="Horizontal" Padding="12,8">
                    <Image Source="ms-appx:///Assets/orchestra-logo.png"
                           Width="36" Height="36" />
                    <TextBlock Text="Orchestra" VerticalAlignment="Center"
                               Margin="12,0,0,0" FontWeight="SemiBold" />
                </StackPanel>
            </NavigationView.PaneHeader>

            <NavigationView.PaneFooter>
                <controls:ConnectionIndicator />
            </NavigationView.PaneFooter>

            <Frame x:Name="ContentFrame" />
        </NavigationView>
    </Grid>
</Window>
```

### Sidebar (Icon Rail — Desktop)

| # | Glyph | Section | Route |
|---|-------|---------|-------|
| 1 | `\uE8BD` (Chat) | Chat | `/chat` |
| 2 | `\uF0E2` (Grid) | Projects | `/projects` |
| 3 | `\uE70B` (Edit) | Notes | `/notes` |
| 4 | `\uE756` (DeveloperTools) | Developer Tools | `/devtools` |
| — | *(spacer)* | | |
| 5 | `\uE713` (Settings) | Settings | `/settings` |

**Active state**: Purple accent background, white icon.
**Inactive state**: Muted gray icon, transparent background.
**Brand logo**: 36x36 purple rounded square at top.
**User profile**: PersonPicture at bottom with context flyout.

### Cross-Platform Feature Matrix

| Feature | Desktop | Xbox | HoloLens | IoT |
|---------|---------|------|----------|-----|
| Chat | Full | View | Spatial | Headless relay |
| Projects | Full | Dashboard | Full | Dashboard |
| Notes | Full | — | Full | — |
| DevTools | Full (10 tools) | — | Full | — |
| System Tray | Yes | — | — | — |
| Window Modes | 3 modes | — | — | — |
| Screenshot | Yes | — | — | — |
| Voice STT/TTS | Yes | Yes | Yes | — |
| Widgets | Win 11 22H2+ | — | — | — |
| Global Hotkey | Yes | — | — | — |
| Toast Notifications | Yes | Yes | Yes | — |
| Jump Lists | Yes | — | — | — |
| Share Target | Yes | — | — | — |

---

## 6. Design System

### Color Tokens (ResourceDictionary)

```xml
<!-- Theme/OrchestraTheme.xaml -->
<ResourceDictionary xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"
                    xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml">

    <!-- Backgrounds -->
    <Color x:Key="BackgroundColor">#0A0D14</Color>
    <Color x:Key="SurfaceColor">#111520</Color>
    <Color x:Key="SurfaceContrastColor">#080A10</Color>
    <Color x:Key="SurfaceActiveColor">#1A1F2E</Color>
    <Color x:Key="SurfaceSelectionColor">#26A900FF</Color>

    <SolidColorBrush x:Key="BackgroundBrush" Color="{StaticResource BackgroundColor}" />
    <SolidColorBrush x:Key="SurfaceBrush" Color="{StaticResource SurfaceColor}" />
    <SolidColorBrush x:Key="SurfaceContrastBrush" Color="{StaticResource SurfaceContrastColor}" />
    <SolidColorBrush x:Key="SurfaceActiveBrush" Color="{StaticResource SurfaceActiveColor}" />
    <SolidColorBrush x:Key="SurfaceSelectionBrush" Color="{StaticResource SurfaceSelectionColor}" />

    <!-- Text -->
    <Color x:Key="TextPrimaryColor">#E8ECF4</Color>
    <Color x:Key="TextMutedColor">#8892A8</Color>
    <Color x:Key="TextDimColor">#4A5268</Color>
    <Color x:Key="TextBrightColor">#F8FAFC</Color>

    <SolidColorBrush x:Key="TextPrimaryBrush" Color="{StaticResource TextPrimaryColor}" />
    <SolidColorBrush x:Key="TextMutedBrush" Color="{StaticResource TextMutedColor}" />
    <SolidColorBrush x:Key="TextDimBrush" Color="{StaticResource TextDimColor}" />
    <SolidColorBrush x:Key="TextBrightBrush" Color="{StaticResource TextBrightColor}" />

    <!-- Structure -->
    <Color x:Key="BorderColor">#1E2436</Color>
    <Color x:Key="AccentColor">#A900FF</Color>

    <SolidColorBrush x:Key="BorderBrush" Color="{StaticResource BorderColor}" />
    <SolidColorBrush x:Key="AccentBrush" Color="{StaticResource AccentColor}" />

    <!-- Semantic -->
    <SolidColorBrush x:Key="SuccessBrush" Color="#22C55E" />
    <SolidColorBrush x:Key="WarningBrush" Color="#F59E0B" />
    <SolidColorBrush x:Key="ErrorBrush" Color="#EF4444" />
    <SolidColorBrush x:Key="InfoBrush" Color="#3B82F6" />

    <!-- Syntax (code editor) -->
    <SolidColorBrush x:Key="SyntaxBlueBrush" Color="#82AAFF" />
    <SolidColorBrush x:Key="SyntaxCyanBrush" Color="#89DDFF" />
    <SolidColorBrush x:Key="SyntaxGreenBrush" Color="#C3E88D" />
    <SolidColorBrush x:Key="SyntaxYellowBrush" Color="#FFCB6B" />
    <SolidColorBrush x:Key="SyntaxOrangeBrush" Color="#F78C6C" />
    <SolidColorBrush x:Key="SyntaxRedBrush" Color="#FF5370" />
    <SolidColorBrush x:Key="SyntaxPurpleBrush" Color="#C792EA" />

</ResourceDictionary>
```

### Typography

```xml
<!-- Theme/Fonts.xaml -->
<ResourceDictionary>
    <x:Double x:Key="BodyFontSize">14</x:Double>
    <x:Double x:Key="BodySecondaryFontSize">13</x:Double>
    <x:Double x:Key="LabelFontSize">12</x:Double>
    <x:Double x:Key="CaptionFontSize">11</x:Double>
    <x:Double x:Key="SectionTitleFontSize">16</x:Double>

    <FontFamily x:Key="DefaultFontFamily">Segoe UI Variable</FontFamily>
    <FontFamily x:Key="MonospaceFontFamily">Cascadia Code</FontFamily>
</ResourceDictionary>
```

### Spacing & Geometry

| Element | Default | Compact | Modern |
|---------|---------|---------|--------|
| NavigationView CompactPaneLength | 56px | 48px | 64px |
| Navigation item | 36x36px | 32x32px | 40x40px |
| Title bar height | 44px | 36px | 48px |
| Status bar height | 24px | 20px | 28px |
| Corner radius (sm) | 6px | 3px | 4px |
| Corner radius (md) | 10px | 5px | 8px |
| Corner radius (lg) | 14px | 8px | 12px |

### 25 Color Themes

All 25 themes from `@orchestra-mcp/theme` are supported. Default: `orchestra` (deep navy).

```csharp
public sealed class ThemeManager
{
    private static readonly Dictionary<string, ThemeDefinition> Themes = new()
    {
        ["orchestra"]    = new("#0A0D14", "#A900FF", false),
        ["dracula"]      = new("#282A36", "#BD93F9", false),
        ["github-dark"]  = new("#0D1117", "#58A6FF", false),
        ["github-light"] = new("#FFFFFF", "#0366D6", true),
        ["one-dark"]     = new("#282C34", "#528BFF", false),
        ["monokai-pro"]  = new("#2D2A2E", "#FFD866", false),
        ["synthwave-84"] = new("#262335", "#FF7EDB", false),
        // ... 18 more themes
    };

    public void ApplyTheme(string themeId)
    {
        if (!Themes.TryGetValue(themeId, out var theme)) return;

        var resources = Application.Current.Resources;
        resources["BackgroundColor"] = ColorFromHex(theme.Background);
        resources["AccentColor"] = ColorFromHex(theme.Accent);
        // ... update all tokens
    }

    private record ThemeDefinition(string Background, string Accent, bool IsLight);
}
```

### Fluent Design Integration

WinUI 3 natively supports Fluent Design System features:

| Feature | Usage |
|---------|-------|
| **Acrylic** | Navigation pane background, Spirit window |
| **Mica** | Main window background (Windows 11) |
| **Reveal Highlight** | Navigation items on hover |
| **Connected Animations** | Page transitions |
| **Depth/Shadow** | Floating panels, dialogs |
| **Motion** | 150ms ease-out for page transitions |

---

## 7. Window Management

### Three Window Modes (Desktop Only)

| Mode | Window | Size | Behavior |
|------|--------|------|----------|
| **Embedded** | Main window | 1280x860 | Full IDE, all sections |
| **Floating** | Spirit window | 420x640 | Frameless, always-on-top, semi-transparent, mini chat |
| **Bubble** | Bubble window | 56x56 | Non-resizable, always-on-top, circular overlay |

**Cycling**: `Win+Shift+O` global hotkey cycles modes.

### Spirit Window (Floating Mini Chat)

```csharp
// Windows/SpiritWindow.xaml.cs
public sealed partial class SpiritWindow : Window
{
    public SpiritWindow()
    {
        InitializeComponent();

        // Get native window handle
        var hwnd = WinRT.Interop.WindowNative.GetWindowHandle(this);
        var windowId = Win32Interop.GetWindowIdFromWindow(hwnd);
        var appWindow = AppWindow.GetFromWindowId(windowId);

        // Configure as floating panel
        appWindow.Title = "Spirit";
        appWindow.Resize(new SizeInt32(420, 640));
        appWindow.SetPresenter(AppWindowPresenterKind.CompactOverlay);

        // Always on top + semi-transparent
        var presenter = appWindow.Presenter as CompactOverlayPresenter;

        // Set extended style for tool window (no taskbar button)
        SetWindowStyle(hwnd);

        // Apply acrylic background
        if (Content is FrameworkElement root)
        {
            root.RequestedTheme = ElementTheme.Dark;
        }
    }

    private static void SetWindowStyle(IntPtr hwnd)
    {
        const int GWL_EXSTYLE = -20;
        const int WS_EX_TOOLWINDOW = 0x00000080;

        var style = PInvoke.GetWindowLong(hwnd, GWL_EXSTYLE);
        PInvoke.SetWindowLong(hwnd, GWL_EXSTYLE, style | WS_EX_TOOLWINDOW);
    }
}
```

### Spirit Window XAML

```xml
<!-- Windows/SpiritWindow.xaml -->
<Window x:Class="Orchestra.Desktop.Windows.SpiritWindow"
        xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"
        xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml">
    <Grid Background="{ThemeResource BackgroundBrush}" Opacity="0.95"
          CornerRadius="12" Padding="0">
        <Grid.RowDefinitions>
            <RowDefinition Height="32" />  <!-- Drag title bar -->
            <RowDefinition Height="*" />   <!-- Chat content -->
            <RowDefinition Height="Auto" /> <!-- Input -->
        </Grid.RowDefinitions>

        <!-- Draggable title bar -->
        <Grid Grid.Row="0" Background="Transparent"
              PointerPressed="TitleBar_PointerPressed">
            <StackPanel Orientation="Horizontal" Padding="12,4">
                <Ellipse Width="8" Height="8"
                         Fill="{StaticResource SuccessBrush}" />
                <TextBlock Text="Spirit" Margin="8,0,0,0"
                           Foreground="{StaticResource TextMutedBrush}"
                           FontSize="{StaticResource CaptionFontSize}" />
            </StackPanel>
        </Grid>

        <!-- Inline chat view -->
        <plugins:SpiritChatView Grid.Row="1" />

        <!-- Chat input -->
        <plugins:ChatInputControl Grid.Row="2" IsCompact="True" />
    </Grid>
</Window>
```

### Bubble Window

```csharp
// Windows/BubbleWindow.xaml.cs
public sealed partial class BubbleWindow : Window
{
    public BubbleWindow()
    {
        InitializeComponent();

        var hwnd = WinRT.Interop.WindowNative.GetWindowHandle(this);
        var windowId = Win32Interop.GetWindowIdFromWindow(hwnd);
        var appWindow = AppWindow.GetFromWindowId(windowId);

        // Fixed 56x56, non-resizable, always on top
        appWindow.Resize(new SizeInt32(56, 56));
        appWindow.SetPresenter(AppWindowPresenterKind.CompactOverlay);

        // Remove title bar, make circular
        var titleBar = appWindow.TitleBar;
        titleBar.ExtendsContentIntoTitleBar = true;

        // Position bottom-right
        var display = DisplayArea.GetFromWindowId(windowId, DisplayAreaFallback.Primary);
        var x = display.WorkArea.Width - 80;
        var y = display.WorkArea.Height - 80;
        appWindow.Move(new PointInt32(x, y));
    }
}
```

### Window Mode Manager

```csharp
public sealed class WindowModeManager
{
    public enum Mode { Embedded, Floating, Bubble }

    private Mode _current = Mode.Embedded;
    private readonly MainWindow _mainWindow;
    private SpiritWindow? _spiritWindow;
    private BubbleWindow? _bubbleWindow;

    public void CycleMode()
    {
        _current = _current switch
        {
            Mode.Embedded => Mode.Floating,
            Mode.Floating => Mode.Bubble,
            Mode.Bubble => Mode.Embedded,
            _ => Mode.Embedded
        };
        ApplyMode();
    }

    private void ApplyMode()
    {
        switch (_current)
        {
            case Mode.Embedded:
                _spiritWindow?.Close(); _spiritWindow = null;
                _bubbleWindow?.Close(); _bubbleWindow = null;
                _mainWindow.Activate();
                break;
            case Mode.Floating:
                _mainWindow.AppWindow.Hide();
                _bubbleWindow?.Close(); _bubbleWindow = null;
                _spiritWindow = App.GetService<SpiritWindow>();
                _spiritWindow.Activate();
                break;
            case Mode.Bubble:
                _mainWindow.AppWindow.Hide();
                _spiritWindow?.Close(); _spiritWindow = null;
                _bubbleWindow = App.GetService<BubbleWindow>();
                _bubbleWindow.Activate();
                break;
        }
    }
}
```

---

## 8. QUIC Transport

### System.Net.Quic Client (.NET 8+)

```csharp
using System.Net.Quic;
using System.Net.Security;
using System.Security.Cryptography.X509Certificates;

public sealed class QUICConnection : IAsyncDisposable
{
    private QuicConnection? _connection;
    private QuicStream? _stream;
    private ConnectionState _state = ConnectionState.Disconnected;
    private CancellationTokenSource? _cts;

    public event Action<ConnectionState>? StateChanged;

    public async Task ConnectAsync(string host, int port, string pluginId,
        CancellationToken ct = default)
    {
        _cts = CancellationTokenSource.CreateLinkedTokenSource(ct);
        var (clientCert, caCert) = MTLSConfig.LoadCerts(pluginId);

        try
        {
            _state = ConnectionState.Connecting;
            StateChanged?.Invoke(_state);

            _connection = await QuicConnection.ConnectAsync(
                new QuicClientConnectionOptions
                {
                    RemoteEndPoint = new IPEndPoint(
                        IPAddress.Parse(host), port),
                    DefaultStreamErrorCode = 0,
                    DefaultCloseErrorCode = 0,
                    ClientAuthenticationOptions = new SslClientAuthenticationOptions
                    {
                        TargetHost = "orchestrator",
                        ApplicationProtocols = [
                            new SslApplicationProtocol("orchestra-plugin")
                        ],
                        ClientCertificates = new X509CertificateCollection
                        {
                            clientCert
                        },
                        RemoteCertificateValidationCallback = (_, cert, _, _) =>
                            ValidateServerCert(cert, caCert)
                    }
                },
                _cts.Token);

            _stream = await _connection.OpenOutboundStreamAsync(
                QuicStreamType.Bidirectional, _cts.Token);

            _state = ConnectionState.Connected;
            StateChanged?.Invoke(_state);

            // Start receive loop
            _ = ReceiveLoopAsync(_cts.Token);
        }
        catch (Exception)
        {
            _state = ConnectionState.Disconnected;
            StateChanged?.Invoke(_state);
            _ = ReconnectWithBackoffAsync();
        }
    }

    private async Task ReconnectWithBackoffAsync()
    {
        var delay = TimeSpan.FromSeconds(1);
        var maxDelay = TimeSpan.FromSeconds(30);

        while (_state == ConnectionState.Disconnected && _cts?.IsCancellationRequested != true)
        {
            await Task.Delay(delay, _cts?.Token ?? default);
            delay = TimeSpan.FromSeconds(Math.Min(delay.TotalSeconds * 2, maxDelay.TotalSeconds));

            try
            {
                await ConnectAsync("localhost", 50100, "ui.windows", _cts?.Token ?? default);
                return;
            }
            catch { /* retry */ }
        }
    }

    private static bool ValidateServerCert(X509Certificate? cert,
        X509Certificate2 caCert)
    {
        if (cert is not X509Certificate2 serverCert) return false;
        using var chain = new X509Chain();
        chain.ChainPolicy.CustomTrustStore.Add(caCert);
        chain.ChainPolicy.TrustMode = X509ChainTrustMode.CustomRootTrust;
        return chain.Build(serverCert);
    }

    public async ValueTask DisposeAsync()
    {
        _cts?.Cancel();
        if (_stream != null) await _stream.DisposeAsync();
        if (_connection != null) await _connection.DisposeAsync();
    }
}
```

### Length-Delimited Protobuf Framing

Matches `libs/sdk-go/plugin/framing.go` byte-for-byte.

```csharp
using System.Buffers.Binary;
using Google.Protobuf;
using Orchestra.Plugin.V1;

public static class StreamFramer
{
    public const int MaxMessageSize = 16 * 1024 * 1024; // 16 MB

    /// Write a length-delimited Protobuf message.
    /// Format: [4 bytes big-endian uint32 length][N bytes Protobuf payload]
    public static async Task WriteAsync(IMessage message, Stream stream,
        CancellationToken ct = default)
    {
        var data = message.ToByteArray();
        if (data.Length > MaxMessageSize)
            throw new InvalidOperationException(
                $"Message size {data.Length} exceeds maximum {MaxMessageSize}");

        var header = new byte[4];
        BinaryPrimitives.WriteUInt32BigEndian(header, (uint)data.Length);

        await stream.WriteAsync(header, ct);
        await stream.WriteAsync(data, ct);
        await stream.FlushAsync(ct);
    }

    /// Read a length-delimited Protobuf message.
    public static async Task<T> ReadAsync<T>(Stream stream,
        CancellationToken ct = default) where T : IMessage<T>, new()
    {
        var header = new byte[4];
        await stream.ReadExactlyAsync(header, ct);

        var size = BinaryPrimitives.ReadUInt32BigEndian(header);
        if (size > MaxMessageSize)
            throw new InvalidOperationException(
                $"Message size {size} exceeds maximum {MaxMessageSize}");

        var data = new byte[size];
        await stream.ReadExactlyAsync(data, ct);

        var parser = new MessageParser<T>(() => new T());
        return parser.ParseFrom(data);
    }
}
```

### mTLS Configuration

```csharp
public static class MTLSConfig
{
    public static (X509Certificate2 clientCert, X509Certificate2 caCert) LoadCerts(
        string pluginId)
    {
        var certsDir = Path.Combine(
            Environment.GetFolderPath(Environment.SpecialFolder.UserProfile),
            ".orchestra", "certs");

        var caCert = new X509Certificate2(
            Path.Combine(certsDir, "ca.crt"));
        var clientCert = X509Certificate2.CreateFromPemFile(
            Path.Combine(certsDir, $"{pluginId}.crt"),
            Path.Combine(certsDir, $"{pluginId}.key"));

        return (clientCert, caCert);
    }
}
```

---

## 9. Data Layer

### MCP Tool Proxy

All data operations go through MCP tool calls routed via the orchestrator.

```csharp
public sealed class ToolService
{
    private readonly OrchestraClient _client;

    public ToolService(OrchestraClient client) => _client = client;

    public async Task<ToolResponse> CallToolAsync(string name,
        IDictionary<string, object> arguments, CancellationToken ct = default)
    {
        return await _client.SendToolCallAsync(name, arguments, ct);
    }

    // Convenience wrappers
    public async Task<IReadOnlyList<Project>> ListProjectsAsync(CancellationToken ct = default)
    {
        var response = await CallToolAsync("list_projects", new Dictionary<string, object>(), ct);
        return ParseProjects(response.Result);
    }

    public async Task<ProjectStatus> GetProjectStatusAsync(string slug,
        CancellationToken ct = default)
    {
        var args = new Dictionary<string, object> { ["slug"] = slug };
        var response = await CallToolAsync("get_project_status", args, ct);
        return ParseProjectStatus(response.Result);
    }

    public async Task<Feature> CreateFeatureAsync(string project, string title,
        string? description = null, CancellationToken ct = default)
    {
        var args = new Dictionary<string, object>
        {
            ["project"] = project,
            ["title"] = title
        };
        if (description != null) args["description"] = description;

        var response = await CallToolAsync("create_feature", args, ct);
        return ParseFeature(response.Result);
    }

    public async Task<Feature> AdvanceFeatureAsync(string id, string evidence,
        CancellationToken ct = default)
    {
        var args = new Dictionary<string, object>
        {
            ["feature_id"] = id,
            ["evidence"] = evidence
        };
        var response = await CallToolAsync("advance_feature", args, ct);
        return ParseFeature(response.Result);
    }

    // Multi-LLM AI calls
    public async Task<string> AiPromptAsync(string prompt, string provider,
        string model, CancellationToken ct = default)
    {
        var args = new Dictionary<string, object>
        {
            ["prompt"] = prompt,
            ["provider"] = provider,
            ["model"] = model
        };
        var response = await CallToolAsync("ai_prompt", args, ct);
        return response.Result.Fields["text"].StringValue;
    }
}
```

### Local Cache (Microsoft.Data.Sqlite)

```csharp
// %LOCALAPPDATA%\Orchestra\cache.db
public sealed class LocalCache : IAsyncDisposable
{
    private readonly SqliteConnection _db;

    public LocalCache()
    {
        var dbPath = Path.Combine(
            Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
            "Orchestra", "cache.db");
        Directory.CreateDirectory(Path.GetDirectoryName(dbPath)!);

        _db = new SqliteConnection($"Data Source={dbPath}");
        _db.Open();
        InitializeSchema();
    }

    private void InitializeSchema()
    {
        using var cmd = _db.CreateCommand();
        cmd.CommandText = """
            CREATE TABLE IF NOT EXISTS projects (
                id TEXT PRIMARY KEY,
                data TEXT NOT NULL,
                updated_at TEXT NOT NULL
            );
            CREATE TABLE IF NOT EXISTS notes (
                id TEXT PRIMARY KEY,
                data TEXT NOT NULL,
                updated_at TEXT NOT NULL
            );
            CREATE TABLE IF NOT EXISTS chat_sessions (
                id TEXT PRIMARY KEY,
                data TEXT NOT NULL,
                updated_at TEXT NOT NULL
            );
            """;
        cmd.ExecuteNonQuery();
    }

    public async Task CacheProjectsAsync(IEnumerable<Project> projects)
    {
        foreach (var project in projects)
        {
            using var cmd = _db.CreateCommand();
            cmd.CommandText = """
                INSERT OR REPLACE INTO projects (id, data, updated_at)
                VALUES ($id, $data, $updated_at)
                """;
            cmd.Parameters.AddWithValue("$id", project.Id);
            cmd.Parameters.AddWithValue("$data", JsonSerializer.Serialize(project));
            cmd.Parameters.AddWithValue("$updated_at", DateTimeOffset.UtcNow.ToString("O"));
            await cmd.ExecuteNonQueryAsync();
        }
    }

    public async Task<IReadOnlyList<Project>> GetCachedProjectsAsync()
    {
        using var cmd = _db.CreateCommand();
        cmd.CommandText = "SELECT data FROM projects ORDER BY updated_at DESC";
        using var reader = await cmd.ExecuteReaderAsync();

        var projects = new List<Project>();
        while (await reader.ReadAsync())
        {
            projects.Add(JsonSerializer.Deserialize<Project>(reader.GetString(0))!);
        }
        return projects;
    }

    public async ValueTask DisposeAsync()
    {
        await _db.DisposeAsync();
    }
}
```

---

## 10. AI Chat

### Multi-LLM Chat Architecture

The chat interface supports 4 AI providers (Claude, OpenAI, Gemini, Ollama) via the bridge plugins.

```
ChatPage
├── ChatSessionList (280px sidebar panel)
│   ├── AutoSuggestBox (search)
│   ├── Button (+ New Chat)
│   └── ListView (session name, provider badge, model badge, date)
└── ChatBox
    ├── ChatHeader (session name, provider, model, Live dot)
    ├── ChatBody (ScrollViewer with ItemsRepeater)
    │   ├── ChatMessageControl (user: right, purple)
    │   ├── ChatMessageControl (assistant: left, surface)
    │   ├── EventCards (tool call results)
    │   └── ProgressRing (typing indicator)
    ├── StatusLine (typing status + elapsed timer)
    └── ChatInput
        ├── TextBox (AcceptsReturn, max 6 lines)
        └── CommandBar (provider, model, tools, attach, send/stop)
```

### Provider & Model Selection

```csharp
public sealed record AIProvider(string Id, string Name, IReadOnlyList<string> Models);

public static class AIProviders
{
    public static readonly IReadOnlyList<AIProvider> All =
    [
        new("claude", "Anthropic", [
            "claude-opus-4-6", "claude-sonnet-4-6", "claude-haiku-4-5"
        ]),
        new("openai", "OpenAI", [
            "gpt-4o", "gpt-4o-mini", "o1", "o1-mini"
        ]),
        new("gemini", "Google", [
            "gemini-2.5-pro", "gemini-2.5-flash", "gemini-2.0-flash"
        ]),
        new("ollama", "Ollama", [
            "llama3", "codellama", "mistral"
        ]),
    ];
}
```

### Message Types

```csharp
public enum MessageRole { User, Assistant, System }

public sealed class ChatMessage
{
    public required string Id { get; init; }
    public required MessageRole Role { get; init; }
    public required string Content { get; set; }
    public required DateTimeOffset Timestamp { get; init; }
    public bool Streaming { get; set; }
    public string? Thinking { get; set; }
    public List<ToolEvent> Events { get; } = [];
    public string? Provider { get; init; }   // "claude", "openai", "gemini", "ollama"
    public string? Model { get; init; }
    public List<Attachment> Attachments { get; } = [];
}
```

### Streaming Glow Effect

When assistant is streaming, the message border gets a rotating gradient:

```xml
<!-- ChatMessageControl.xaml -->
<Border x:Name="GlowBorder" CornerRadius="16" Padding="1"
        Visibility="{x:Bind Message.Streaming, Mode=OneWay}">
    <Border.Background>
        <LinearGradientBrush x:Name="GlowGradient">
            <GradientStop Color="#3B82F6" Offset="0.0" />
            <GradientStop Color="#A900FF" Offset="0.25" />
            <GradientStop Color="#EF4444" Offset="0.5" />
            <GradientStop Color="#EC4899" Offset="0.75" />
            <GradientStop Color="#3B82F6" Offset="1.0" />
        </LinearGradientBrush>
    </Border.Background>
    <Border Background="{StaticResource SurfaceBrush}" CornerRadius="15"
            Padding="16">
        <!-- Message content -->
    </Border>
</Border>
```

```csharp
// Rotate the gradient via animation
private void StartGlowAnimation()
{
    var storyboard = new Storyboard();
    var animation = new DoubleAnimation
    {
        From = 0, To = 360,
        Duration = new Duration(TimeSpan.FromSeconds(3)),
        RepeatBehavior = RepeatBehavior.Forever
    };
    Storyboard.SetTarget(animation, GlowRotation);
    Storyboard.SetTargetProperty(animation, "Angle");
    storyboard.Children.Add(animation);
    storyboard.Begin();
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
ProjectsPage
├── ProjectSidebar (280px, NavigationView secondary pane)
│   ├── AutoSuggestBox (search)
│   ├── Button (+ New Project)
│   └── ListView (icon, name, task count, %, Active badge)
└── ProjectDetail
    ├── Header (icon, name, description)
    ├── Status card (Total Tasks, Completed, Completion %)
    │   └── ProgressBar (purple, animated)
    ├── Status breakdown (InfoBar per state)
    └── BacklogTree (TreeView — epics > stories > tasks)
```

### Workflow States (13-State Machine)

```csharp
public enum WorkflowState
{
    Backlog, Todo, InProgress, ReadyForTesting, InTesting,
    ReadyForDocs, InDocs, Documented, InReview, Done,
    Blocked, Rejected, Cancelled
}

public static class WorkflowStateExtensions
{
    public static Color GetColor(this WorkflowState state) => state switch
    {
        WorkflowState.Backlog or WorkflowState.Cancelled => ColorFromHex("#6B7280"),
        WorkflowState.Todo => ColorFromHex("#3B82F6"),
        WorkflowState.InProgress => ColorFromHex("#A900FF"),
        WorkflowState.ReadyForTesting => ColorFromHex("#EAB308"),
        WorkflowState.InTesting => ColorFromHex("#F97316"),
        WorkflowState.ReadyForDocs => ColorFromHex("#06B6D4"),
        WorkflowState.InDocs => ColorFromHex("#14B8A6"),
        WorkflowState.Documented or WorkflowState.Done => ColorFromHex("#22C55E"),
        WorkflowState.InReview => ColorFromHex("#6366F1"),
        WorkflowState.Blocked or WorkflowState.Rejected => ColorFromHex("#EF4444"),
        _ => ColorFromHex("#6B7280")
    };

    public static string GetDisplayName(this WorkflowState state) => state switch
    {
        WorkflowState.InProgress => "In Progress",
        WorkflowState.ReadyForTesting => "Ready for Testing",
        WorkflowState.InTesting => "In Testing",
        WorkflowState.ReadyForDocs => "Ready for Docs",
        WorkflowState.InDocs => "In Docs",
        WorkflowState.InReview => "In Review",
        _ => state.ToString()
    };
}
```

---

## 12. Notes & Docs

### Notes Plugin

Uses `tools.notes` plugin (8 tools): `create_note`, `get_note`, `update_note`, `delete_note`, `list_notes`, `search_notes`, `pin_note`, `tag_note`.

```
NotesPage
├── NotesSidebar (280px)
│   ├── AutoSuggestBox (search)
│   ├── Button (+ New Note)
│   ├── Pinned Notes section (Expander)
│   └── Other Notes section (Expander)
└── NoteEditor
    ├── CommandBar (back, pin toggle, save, delete)
    ├── TextBox (title — large, borderless)
    ├── ItemsRepeater (tags — add/remove chips)
    └── TextBox (content — monospace, AcceptsReturn, markdown)
```

### Docs/Wiki Plugin

Uses `tools.docs` plugin (10 tools): `doc_create`, `doc_get`, `doc_update`, `doc_delete`, `doc_list`, `doc_search`, `doc_generate`, `doc_index`, `doc_tree`, `doc_export`.

Categories: `api-reference`, `guide`, `architecture`, `tutorial`, `changelog`, `decision-record`

### Markdown Parser Plugin

Uses `tools.markdown` plugin (8 tools): `md_parse`, `md_parse_file`, `md_parse_frontmatter`, `md_render_html`, `md_render_plaintext`, `md_toc`, `md_lint`, `md_transform`.

Supports: Mermaid, KaTeX, footnotes, embeds, task lists, frontmatter, heading anchors.

Rendered in WinUI via `WebView2` control for HTML output:

```csharp
public sealed partial class MarkdownPreview : UserControl
{
    private readonly ToolService _tools;

    public async Task RenderAsync(string markdown)
    {
        var response = await _tools.CallToolAsync("md_render_html",
            new Dictionary<string, object> { ["markdown"] = markdown });

        var html = response.Result.Fields["html"].StringValue;
        await MarkdownWebView.EnsureCoreWebView2Async();
        MarkdownWebView.NavigateToString(WrapWithTheme(html));
    }

    private string WrapWithTheme(string html) => $"""
        <!DOCTYPE html>
        <html><head>
        <style>
            body {{ background: #111520; color: #E8ECF4; font-family: 'Segoe UI Variable'; }}
            pre {{ background: #0A0D14; padding: 16px; border-radius: 8px; }}
            code {{ font-family: 'Cascadia Code'; }}
            a {{ color: #A900FF; }}
        </style>
        </head><body>{html}</body></html>
        """;
}
```

---

## 13. Developer Tools

### DevTools Plugin Container

The DevTools section hosts sub-plugins for each tool type. Each is a separate Go plugin on the backend, exposed through the DevToolsPlugin container.

| Tool | Plugin ID | Tools | Description |
|------|-----------|-------|-------------|
| File Explorer | devtools.file-explorer | 17 | Full IDE: file ops + LSP code intelligence |
| Terminal | devtools.terminal | 6 | PTY terminal sessions (PowerShell, cmd, WSL) |
| SSH | devtools.ssh | 7 | Remote SSH + SFTP |
| Services | devtools.services | 6 | Service manager (sc.exe, Get-Service) |
| Docker | devtools.docker | 10 | Container management |
| Debugger | devtools.debugger | 8 | DAP protocol debugger |
| Test Runner | devtools.test-runner | 6 | Multi-framework test runner |
| Log Viewer | devtools.log-viewer | 5 | Log streaming + search (Event Viewer) |
| Database | devtools.database | 8 | SQL query editor + schema browser |
| DevOps | devtools.devops | 8 | CI/CD pipeline management |

### Terminal Plugin (Windows-Specific)

Windows terminal supports multiple shells via Windows Terminal integration:

```csharp
public sealed class TerminalService
{
    public enum ShellType { PowerShell, Cmd, Wsl, GitBash }

    public string GetShellPath(ShellType shell) => shell switch
    {
        ShellType.PowerShell => "pwsh.exe",  // PowerShell 7+
        ShellType.Cmd => "cmd.exe",
        ShellType.Wsl => "wsl.exe",
        ShellType.GitBash => @"C:\Program Files\Git\bin\bash.exe",
        _ => "pwsh.exe"
    };
}
```

### Services Plugin (Windows-Specific)

Uses `sc.exe` and PowerShell `Get-Service` instead of `launchctl`/`systemctl`:

| macOS (launchctl) | Windows Equivalent |
|-------------------|--------------------|
| `launchctl list` | `Get-Service` / `sc.exe query` |
| `launchctl start <svc>` | `Start-Service` / `sc.exe start` |
| `launchctl stop <svc>` | `Stop-Service` / `sc.exe stop` |
| `launchctl load <plist>` | `New-Service` / `sc.exe create` |

### File Explorer (Full IDE + LSP)

Same tools as Swift but rendered in WinUI:

**File Tools:** `list_directory`, `read_file`, `write_file`, `move_file`, `delete_file`, `file_info`, `file_search`

**Code Intelligence Tools:** `code_symbols`, `code_goto_definition`, `code_find_references`, `code_hover`, `code_complete`, `code_diagnostics`, `code_actions`, `code_workspace_symbols`, `code_namespace`, `code_imports`

Code editor rendered via `WebView2` + Monaco Editor (same engine as VS Code):

```csharp
public sealed partial class CodeEditorControl : UserControl
{
    public async Task LoadFileAsync(string path, string content, string language)
    {
        await EditorWebView.EnsureCoreWebView2Async();
        await EditorWebView.CoreWebView2.ExecuteScriptAsync($"""
            editor.setValue({JsonSerializer.Serialize(content)});
            monaco.editor.setModelLanguage(editor.getModel(), '{language}');
            """);
    }
}
```

---

## 14. AI Awareness

Four AI awareness plugins provide visual and contextual understanding.

### ai.screenshot (6 tools)
`capture_screen`, `capture_region`, `capture_window`, `capture_interactive`, `annotate_screenshot`, `list_captures`

Windows: **Windows.Graphics.Capture** (Win10 19041+).

```csharp
using Windows.Graphics.Capture;
using Windows.Graphics.DirectX;
using Windows.Graphics.DirectX.Direct3D11;

public sealed class ScreenCaptureService
{
    public async Task<SoftwareBitmap> CaptureScreenAsync()
    {
        var picker = new GraphicsCapturePicker();
        var item = await picker.CreateForWindowAsync(
            WinRT.Interop.WindowNative.GetWindowHandle(App.MainWindow));

        if (item == null) return null!;

        var framePool = Direct3D11CaptureFramePool.Create(
            _device, DirectXPixelFormat.B8G8R8A8UIntNormalized, 1,
            item.Size);

        var session = framePool.CreateCaptureSession(item);
        session.StartCapture();

        var frame = await GetNextFrameAsync(framePool);
        session.Dispose();
        return ConvertToSoftwareBitmap(frame);
    }

    public async Task<SoftwareBitmap> CaptureWindowAsync(IntPtr hwnd)
    {
        var item = GraphicsCaptureItem.CreateFromWindowHandle(hwnd);
        // ... same capture flow
    }
}
```

### ai.vision (6 tools)
`analyze_image`, `extract_text`, `find_elements`, `compare_images`, `describe_screen`, `extract_data`

Uses Claude Vision API or OpenAI Vision as fallback. Additionally, Windows provides built-in OCR:

```csharp
using Windows.Media.Ocr;

public async Task<string> ExtractTextAsync(SoftwareBitmap bitmap)
{
    var engine = OcrEngine.TryCreateFromUserProfileLanguages();
    var result = await engine.RecognizeAsync(bitmap);
    return result.Text;
}
```

### ai.browser-context (7 tools)
`get_page_content`, `get_page_dom`, `get_selected_text`, `get_open_tabs`, `get_page_screenshot`, `navigate_to`, `execute_script`

Communicates via WebSocket to Chrome extension (same as macOS).

### ai.screen-reader (6 tools)
`get_accessibility_tree`, `get_focused_element`, `find_element`, `get_element_hierarchy`, `list_windows`, `get_window_elements`

Windows: **UI Automation** (UIAutomationClient COM, or `System.Windows.Automation`).

```csharp
using System.Windows.Automation;

public sealed class UIAutomationService
{
    public AutomationElement GetFocusedElement() =>
        AutomationElement.FocusedElement;

    public IReadOnlyList<WindowInfo> ListWindows()
    {
        var root = AutomationElement.RootElement;
        var windows = root.FindAll(TreeScope.Children,
            new PropertyCondition(AutomationElement.ControlTypeProperty,
                ControlType.Window));

        return windows.Cast<AutomationElement>()
            .Select(w => new WindowInfo
            {
                Title = w.Current.Name,
                ProcessId = w.Current.ProcessId,
                BoundingRect = w.Current.BoundingRectangle,
                Handle = new IntPtr(w.Current.NativeWindowHandle)
            })
            .ToList();
    }

    public AccessibilityTree GetAccessibilityTree(IntPtr hwnd)
    {
        var element = AutomationElement.FromHandle(hwnd);
        return BuildTree(element, maxDepth: 5);
    }
}
```

---

## 15. Voice & Notifications

### Voice Plugin (services.voice — 8 tools)
`tts_speak`, `tts_speak_provider`, `tts_list_voices`, `tts_stop`, `stt_listen`, `stt_transcribe_file`, `stt_list_models`, `voice_config`

**OS TTS:** `Windows.Media.SpeechSynthesis`
**OS STT:** `Windows.Media.SpeechRecognition`
**Provider TTS:** ElevenLabs, OpenAI TTS, Google Cloud TTS
**Provider STT:** OpenAI Whisper, Google Cloud Speech, Deepgram

```csharp
using Windows.Media.SpeechSynthesis;
using Windows.Media.SpeechRecognition;

public sealed class VoiceService
{
    private readonly SpeechSynthesizer _synthesizer = new();
    private SpeechRecognizer? _recognizer;

    // TTS
    public async Task SpeakAsync(string text)
    {
        var stream = await _synthesizer.SynthesizeTextToStreamAsync(text);
        var player = new MediaPlayer();
        player.Source = MediaSource.CreateFromStream(stream, stream.ContentType);
        player.Play();
    }

    public IReadOnlyList<VoiceInformation> ListVoices() =>
        SpeechSynthesizer.AllVoices.ToList();

    // STT
    public async Task<string> ListenAsync(TimeSpan? timeout = null)
    {
        _recognizer ??= new SpeechRecognizer();
        await _recognizer.CompileConstraintsAsync();

        var result = await _recognizer.RecognizeAsync();
        return result.Status == SpeechRecognitionResultStatus.Success
            ? result.Text
            : string.Empty;
    }
}
```

### Notifications Plugin (services.notifications — 8 tools)
`notify_send`, `notify_schedule`, `notify_cancel`, `notify_list_pending`, `notify_badge`, `notify_config`, `notify_history`, `notify_create_channel`

**Channels:** `build`, `test`, `deploy`, `ai`, `reminder`, `system`, `git`

Windows Toast Notifications via `Microsoft.Windows.AppNotifications`:

```csharp
using Microsoft.Windows.AppNotifications;
using Microsoft.Windows.AppNotifications.Builder;

public sealed class NotificationService
{
    public void Initialize()
    {
        var manager = AppNotificationManager.Default;
        manager.NotificationInvoked += OnNotificationInvoked;
        manager.Register();
    }

    public void SendToast(string title, string body, string channel,
        IReadOnlyList<NotificationAction>? actions = null)
    {
        var builder = new AppNotificationBuilder()
            .AddArgument("channel", channel)
            .AddText(title)
            .AddText(body)
            .SetGroup(channel);

        if (actions != null)
        {
            foreach (var action in actions)
            {
                builder.AddButton(new AppNotificationButton(action.Label)
                    .AddArgument("action", action.ToolName)
                    .AddArgument("args", action.ToolArgs));
            }
        }

        AppNotificationManager.Default.Show(builder.BuildNotification());
    }

    public void ScheduleToast(string title, string body, DateTimeOffset when)
    {
        // Windows scheduled notifications via ScheduledToastNotification
        // Requires Windows.UI.Notifications for scheduled support
    }

    private void OnNotificationInvoked(AppNotificationManager sender,
        AppNotificationActivatedEventArgs args)
    {
        // Parse action arguments and trigger MCP tool call
        if (args.Arguments.TryGetValue("action", out var toolName))
        {
            _ = App.GetService<ToolService>().CallToolAsync(toolName,
                new Dictionary<string, object>());
        }
    }
}

public sealed record NotificationAction(string Label, string ToolName, string ToolArgs);
```

---

## 16. Settings & Integrations

### Settings Navigation

| Section | Settings |
|---------|----------|
| **General** | Timezone, language, startup behavior |
| **Appearance** | Color theme (25 themes), component variant (3 options), Mica/Acrylic toggle |
| **Notifications** | Permission status, channels, Focus Assist integration |
| **Windows** | Default window mode, Spirit bounds, Bubble position |
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

Discord, Slack, GitHub, Jira, Linear, Notion, Microsoft 365, Figma — each shows connection status and configuration.

### Settings Storage

```csharp
// %LOCALAPPDATA%\Orchestra\settings.json
public sealed class SettingsService
{
    private readonly string _path = Path.Combine(
        Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
        "Orchestra", "settings.json");

    public AppSettings Load()
    {
        if (!File.Exists(_path)) return new AppSettings();
        var json = File.ReadAllText(_path);
        return JsonSerializer.Deserialize<AppSettings>(json) ?? new AppSettings();
    }

    public void Save(AppSettings settings)
    {
        Directory.CreateDirectory(Path.GetDirectoryName(_path)!);
        File.WriteAllText(_path, JsonSerializer.Serialize(settings,
            new JsonSerializerOptions { WriteIndented = true }));
    }
}

public sealed class AppSettings
{
    public string Theme { get; set; } = "orchestra";
    public string ComponentVariant { get; set; } = "default";
    public string DefaultProvider { get; set; } = "claude";
    public string DefaultModel { get; set; } = "claude-sonnet-4-6";
    public string WindowMode { get; set; } = "embedded";
    public bool StartMinimized { get; set; }
    public bool LaunchAtStartup { get; set; }
    public string Language { get; set; } = "en";
    public Dictionary<string, bool> NotificationChannels { get; set; } = new()
    {
        ["build"] = true, ["test"] = true, ["deploy"] = true,
        ["ai"] = true, ["reminder"] = true, ["system"] = true, ["git"] = true
    };
}
```

---

## 17. Native Windows Features

### Global Hotkey (RegisterHotKey — Win32)

```csharp
using System.Runtime.InteropServices;

public sealed partial class GlobalHotkeyService : IDisposable
{
    private const int WM_HOTKEY = 0x0312;
    private const int MOD_WIN = 0x0008;
    private const int MOD_SHIFT = 0x0004;
    private const int HOTKEY_ID = 9001;

    [LibraryImport("user32.dll")]
    [return: MarshalAs(UnmanagedType.Bool)]
    private static partial bool RegisterHotKey(IntPtr hWnd, int id, uint fsModifiers, uint vk);

    [LibraryImport("user32.dll")]
    [return: MarshalAs(UnmanagedType.Bool)]
    private static partial bool UnregisterHotKey(IntPtr hWnd, int id);

    private readonly IntPtr _hwnd;

    public GlobalHotkeyService(IntPtr hwnd)
    {
        _hwnd = hwnd;
        // Win+Shift+O  (O = 0x4F)
        RegisterHotKey(_hwnd, HOTKEY_ID, MOD_WIN | MOD_SHIFT, 0x4F);
    }

    public void HandleMessage(uint msg, IntPtr wParam)
    {
        if (msg == WM_HOTKEY && wParam == HOTKEY_ID)
        {
            App.GetService<WindowModeManager>().CycleMode();
        }
    }

    public void Dispose()
    {
        UnregisterHotKey(_hwnd, HOTKEY_ID);
    }
}
```

### Windows Credential Manager

```csharp
using Windows.Security.Credentials;

public sealed class CredentialService
{
    private const string ResourcePrefix = "Orchestra_";
    private readonly PasswordVault _vault = new();

    public void SaveAPIKey(string provider, string key)
    {
        var credential = new PasswordCredential(
            $"{ResourcePrefix}{provider}", "api_key", key);
        _vault.Add(credential);
    }

    public string? LoadAPIKey(string provider)
    {
        try
        {
            var credential = _vault.Retrieve($"{ResourcePrefix}{provider}", "api_key");
            credential.RetrievePassword();
            return credential.Password;
        }
        catch (Exception)
        {
            return null;
        }
    }

    public void DeleteAPIKey(string provider)
    {
        try
        {
            var credential = _vault.Retrieve($"{ResourcePrefix}{provider}", "api_key");
            _vault.Remove(credential);
        }
        catch { /* not found, ignore */ }
    }
}
```

### Auto-Updater (MSIX + GitHub Releases)

```csharp
using Windows.Management.Deployment;

public sealed class UpdaterService
{
    private const string Repo = "orchestra-mcp/orchestra-windows";
    private static readonly TimeSpan CheckInterval = TimeSpan.FromHours(6);

    public async Task<UpdateInfo?> CheckForUpdateAsync(CancellationToken ct = default)
    {
        using var http = new HttpClient();
        http.DefaultRequestHeaders.UserAgent.ParseAdd("Orchestra-Windows/1.0");

        var json = await http.GetStringAsync(
            $"https://api.github.com/repos/{Repo}/releases/latest", ct);
        var release = JsonSerializer.Deserialize<GitHubRelease>(json);

        if (release == null) return null;

        var currentVersion = Package.Current.Id.Version;
        var latestVersion = Version.Parse(release.TagName.TrimStart('v'));

        if (latestVersion <= new Version(
            currentVersion.Major, currentVersion.Minor,
            currentVersion.Build, currentVersion.Revision))
            return null;

        var msixAsset = release.Assets
            .FirstOrDefault(a => a.Name.EndsWith(".msix", StringComparison.OrdinalIgnoreCase));

        return msixAsset == null ? null : new UpdateInfo
        {
            Version = release.TagName,
            DownloadUrl = msixAsset.BrowserDownloadUrl,
            ReleaseNotes = release.Body
        };
    }

    public async Task ApplyUpdateAsync(UpdateInfo info, CancellationToken ct = default)
    {
        // Download MSIX
        using var http = new HttpClient();
        var msixPath = Path.Combine(Path.GetTempPath(), "orchestra-update.msix");
        var bytes = await http.GetByteArrayAsync(info.DownloadUrl, ct);
        await File.WriteAllBytesAsync(msixPath, bytes, ct);

        // Install via PackageManager
        var pm = new PackageManager();
        await pm.AddPackageAsync(new Uri(msixPath),
            null, DeploymentOptions.ForceApplicationShutdown);
    }
}
```

### Jump Lists

```csharp
using Windows.UI.StartScreen;

public static class JumpListService
{
    public static async Task UpdateJumpListAsync(IReadOnlyList<Project> recentProjects)
    {
        var jumpList = await JumpList.LoadCurrentAsync();
        jumpList.Items.Clear();

        // Recent projects
        foreach (var project in recentProjects.Take(5))
        {
            var item = JumpListItem.CreateWithArguments(
                $"--project {project.Slug}",
                project.Name);
            item.GroupName = "Recent Projects";
            item.Logo = new Uri("ms-appx:///Assets/project-icon.png");
            jumpList.Items.Add(item);
        }

        // Quick actions
        var newChat = JumpListItem.CreateWithArguments("--new-chat", "New AI Chat");
        newChat.GroupName = "Actions";
        jumpList.Items.Add(newChat);

        var spirit = JumpListItem.CreateWithArguments("--spirit", "Open Spirit");
        spirit.GroupName = "Actions";
        jumpList.Items.Add(spirit);

        await jumpList.SaveAsync();
    }
}
```

### System Tray (H.NotifyIcon.WinUI)

```csharp
// NuGet: H.NotifyIcon.WinUI
using H.NotifyIcon;

public sealed class TrayIconService : IDisposable
{
    private TaskbarIcon? _trayIcon;

    public void Initialize(Window mainWindow)
    {
        _trayIcon = new TaskbarIcon
        {
            IconSource = new BitmapIconSource
            {
                UriSource = new Uri("ms-appx:///Assets/tray-icon.ico")
            },
            ToolTipText = "Orchestra MCP"
        };

        var menu = new MenuFlyout();

        menu.Items.Add(new MenuFlyoutItem
        {
            Text = "Open Orchestra",
            Command = new RelayCommand(() => mainWindow.Activate())
        });

        menu.Items.Add(new MenuFlyoutItem
        {
            Text = "Spirit Mode",
            Command = new RelayCommand(() =>
                App.GetService<WindowModeManager>().SetMode(WindowModeManager.Mode.Floating))
        });

        menu.Items.Add(new MenuFlyoutSeparator());

        var statusItem = new MenuFlyoutItem
        {
            Text = "Connected",
            Icon = new FontIcon { Glyph = "\uF13D", Foreground = new SolidColorBrush(Colors.Green) },
            IsEnabled = false
        };
        menu.Items.Add(statusItem);

        menu.Items.Add(new MenuFlyoutSeparator());

        menu.Items.Add(new MenuFlyoutItem
        {
            Text = "Quit",
            Command = new RelayCommand(() => Application.Current.Exit())
        });

        _trayIcon.ContextFlyout = menu;
    }

    public void Dispose() => _trayIcon?.Dispose();
}
```

### Share Target

Register as a Share Target to receive content from other apps:

```xml
<!-- Package.appxmanifest -->
<Extensions>
    <uap:Extension Category="windows.shareTarget">
        <uap:ShareTarget>
            <uap:SupportedFileTypes>
                <uap:SupportsAnyFileType />
            </uap:SupportedFileTypes>
            <uap:DataFormat>Text</uap:DataFormat>
            <uap:DataFormat>URI</uap:DataFormat>
        </uap:ShareTarget>
    </uap:Extension>
</Extensions>
```

### Windows Hello (Biometric Auth)

```csharp
using Windows.Security.Credentials;

public sealed class WindowsHelloService
{
    public async Task<bool> IsAvailableAsync()
    {
        return await KeyCredentialManager.IsSupportedAsync();
    }

    public async Task<bool> AuthenticateAsync()
    {
        var result = await KeyCredentialManager.RequestCreateAsync(
            "Orchestra", KeyCredentialCreationOption.ReplaceExisting);
        return result.Status == KeyCredentialStatus.Success;
    }
}
```

### Permissions

```csharp
public sealed class PermissionsService
{
    // Notification permissions
    public async Task<bool> RequestNotificationPermissionAsync()
    {
        var manager = AppNotificationManager.Default;
        var setting = manager.Setting;
        return setting == AppNotificationSetting.Enabled;
    }

    // Microphone (for STT)
    public async Task<bool> RequestMicrophonePermissionAsync()
    {
        var status = await AppCapability.Create("microphone").RequestAccessAsync();
        return status == AppCapabilityAccessStatus.Allowed;
    }

    // Screen capture
    public bool CheckScreenCapturePermission()
    {
        // Windows.Graphics.Capture requires user consent via picker
        // No upfront permission required
        return GraphicsCaptureSession.IsSupported;
    }

    // Location (for timezone auto-detect)
    public async Task<bool> RequestLocationPermissionAsync()
    {
        var status = await Geolocator.RequestAccessAsync();
        return status == GeolocationAccessStatus.Allowed;
    }
}
```

---

## 18. Packaging & Distribution

### MSIX (Primary — Microsoft Store + Sideload)

MSIX is the modern Windows packaging format — signed, auto-updating, clean uninstall.

```xml
<!-- Package.appxmanifest -->
<Package
  xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10"
  xmlns:uap="http://schemas.microsoft.com/appx/manifest/uap/windows10"
  xmlns:rescap="http://schemas.microsoft.com/appx/manifest/foundation/windows10/restrictedcapabilities"
  IgnorableNamespaces="uap rescap">

  <Identity
    Name="OrchestraMCP.OrchestraWindows"
    Publisher="CN=Orchestra MCP"
    Version="1.0.0.0"
    ProcessorArchitecture="x64" />

  <Properties>
    <DisplayName>Orchestra MCP</DisplayName>
    <PublisherDisplayName>Orchestra MCP</PublisherDisplayName>
    <Logo>Assets\StoreLogo.png</Logo>
  </Properties>

  <Dependencies>
    <TargetDeviceFamily Name="Windows.Desktop"
      MinVersion="10.0.19041.0" MaxVersionTested="10.0.22621.0" />
  </Dependencies>

  <Resources>
    <Resource Language="x-generate"/>
  </Resources>

  <Applications>
    <Application Id="App"
      Executable="$targetnametoken$.exe"
      EntryPoint="$targetentrypoint$">
      <uap:VisualElements
        DisplayName="Orchestra MCP"
        Description="AI-agentic IDE for project and feature management"
        BackgroundColor="transparent"
        Square150x150Logo="Assets\Square150x150Logo.png"
        Square44x44Logo="Assets\Square44x44Logo.png">
        <uap:DefaultTile Wide310x150Logo="Assets\Wide310x150Logo.png"/>
        <uap:SplashScreen Image="Assets\SplashScreen.png"/>
      </uap:VisualElements>
    </Application>
  </Applications>

  <Capabilities>
    <rescap:Capability Name="runFullTrust"/>
    <Capability Name="internetClient"/>
    <Capability Name="microphone"/>
    <uap:Capability Name="userAccountInformation"/>
    <rescap:Capability Name="screenDuplication"/>
  </Capabilities>
</Package>
```

### Build & Sign

```powershell
# Build for publishing (Release, self-contained)
dotnet publish src/Orchestra.Desktop/Orchestra.Desktop.csproj `
  -c Release -r win-x64 --self-contained `
  -p:PublishSingleFile=false `
  -p:WindowsAppSDKSelfContained=true `
  -p:GenerateAppxPackageOnBuild=true `
  -p:AppxPackageSigningEnabled=true `
  -p:PackageCertificateThumbprint=$env:CERT_THUMBPRINT

# Or via Makefile
make build-windows
make sign-windows
```

### Release Matrix

| Channel | Format | Target | Update Mechanism |
|---------|--------|--------|-----------------|
| **Microsoft Store** | MSIX | End users | Store auto-update |
| **winget** | MSIX installer | Developers | `winget upgrade orchestra` |
| **Sideload** | MSIX (signed) | Enterprise | MSIX auto-update (PackageManager API) |
| **Direct ZIP** | Portable EXE | Advanced users | Manual / in-app check |
| **Chocolatey** | nupkg (wraps MSIX) | DevOps | `choco upgrade orchestra` |

### winget Manifest

```yaml
# manifests/o/OrchestraMCP/Orchestra/1.0.0/OrchestraMCP.Orchestra.installer.yaml
PackageIdentifier: OrchestraMCP.Orchestra
PackageVersion: 1.0.0
MinimumOSVersion: 10.0.19041.0
InstallerType: msix
Installers:
  - Architecture: x64
    InstallerUrl: https://github.com/orchestra-mcp/orchestra-windows/releases/download/v1.0.0/Orchestra-1.0.0-x64.msix
    InstallerSha256: <sha256>
    Scope: user
```

```powershell
# Users install via:
winget install OrchestraMCP.Orchestra

# Or upgrade:
winget upgrade orchestra
```

### GitHub Actions Release Pipeline

```yaml
# .github/workflows/release.yml
name: Release Windows App

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup .NET 8
        uses: actions/setup-dotnet@v4
        with:
          dotnet-version: '8.0.x'

      - name: Install Windows App SDK
        run: |
          dotnet workload install windowsappsdk

      - name: Decode certificate
        run: |
          $certBytes = [System.Convert]::FromBase64String("${{ secrets.CERT_BASE64 }}")
          $certPath = "${{ runner.temp }}\orchestra.pfx"
          [IO.File]::WriteAllBytes($certPath, $certBytes)

      - name: Build MSIX
        run: |
          dotnet publish src/Orchestra.Desktop/Orchestra.Desktop.csproj `
            -c Release -r win-x64 `
            -p:WindowsAppSDKSelfContained=true `
            -p:GenerateAppxPackageOnBuild=true `
            -p:AppxPackageSigningEnabled=true `
            -p:PackageCertificateKeyFile="${{ runner.temp }}\orchestra.pfx" `
            -p:PackageCertificatePassword="${{ secrets.CERT_PASSWORD }}"

      - name: Upload MSIX to GitHub Release
        uses: softprops/action-gh-release@v2
        with:
          files: |
            src/Orchestra.Desktop/bin/Release/**/*.msix
            src/Orchestra.Desktop/bin/Release/**/*.appxbundle
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Submit winget manifest
        run: |
          wingetcreate update OrchestraMCP.Orchestra \
            --version ${{ github.ref_name }} \
            --urls https://github.com/orchestra-mcp/orchestra-windows/releases/download/${{ github.ref_name }}/Orchestra-${{ github.ref_name }}-x64.msix \
            --submit
        env:
          WINGETCREATE_GITHUB_TOKEN: ${{ secrets.WINGET_TOKEN }}
```

### Signing Strategy

| Environment | Certificate | Authority |
|-------------|-------------|-----------|
| Development | Self-signed (dev cert) | Local machine |
| CI | Code signing cert (PFX in GitHub Secrets) | DigiCert / Sectigo |
| Microsoft Store | Store-managed | Microsoft |

### App Startup (First Run)

```csharp
// Check if running from MSIX package
public static bool IsRunningAsPackage =>
    Windows.ApplicationModel.Package.Current is not null;

// Check for update on first run
public static async Task CheckFirstRunAsync()
{
    var settings = App.GetService<SettingsService>().Load();
    if (!settings.HasRunBefore)
    {
        settings.HasRunBefore = true;
        App.GetService<SettingsService>().Save(settings);
        // Show onboarding
    }
}
```

### Xbox & HoloLens Builds

```xml
<!-- Xbox-specific target family -->
<TargetDeviceFamily Name="Windows.Xbox"
  MinVersion="10.0.19041.0" MaxVersionTested="10.0.22621.0" />

<!-- HoloLens-specific target family -->
<TargetDeviceFamily Name="Windows.Holographic"
  MinVersion="10.0.19041.0" MaxVersionTested="10.0.22621.0" />
```

Xbox and HoloLens apps are distributed through Xbox Developer Mode and HoloLens Device Portal respectively (not Microsoft Store for initial releases).

### Makefile Targets

```makefile
build-windows:
	cd apps/windows && dotnet publish src/Orchestra.Desktop \
	  -c Release -r win-x64 --self-contained \
	  -o ../../bin/windows

sign-windows:
	cd apps/windows && dotnet publish src/Orchestra.Desktop \
	  -c Release -r win-x64 \
	  -p:GenerateAppxPackageOnBuild=true \
	  -p:AppxPackageSigningEnabled=true

test-windows:
	cd apps/windows && dotnet test

clean-windows:
	cd apps/windows && dotnet clean && rm -rf bin/windows
```

---

## 19. Build Phases

### Phase 1: Shell + Chat (MVP) — Windows Desktop

1. .NET Solution setup (`apps/windows/Orchestra.sln`)
2. `Orchestra.Core` class library with QUIC transport + plugin system
3. Main window with plugin-driven NavigationView
4. System tray (H.NotifyIcon.WinUI)
5. AI chat plugin: session list + multi-LLM conversation + streaming
6. Settings plugin: appearance (theme picker)
7. Connection status indicator

**Exit criteria**: Launch app, see plugin-driven NavigationView, create chat session with any provider, see streaming response.

### Phase 2: Projects + Notes — Windows Desktop

1. Projects plugin: list, detail, backlog tree (TreeView), workflow states
2. Notes plugin: list, editor, pin/unpin
3. Search spotlight (Ctrl+K command palette)
4. Data caching in local SQLite (Microsoft.Data.Sqlite)
5. Jump Lists for recent projects

### Phase 3: Developer Tools — Windows Desktop

1. File Explorer plugin (with LSP code intelligence via WebView2 + Monaco)
2. Terminal plugin (PowerShell PTY via ConPTY)
3. Additional DevTools plugins (Database, SSH, Log Viewer)
4. Component browser plugin

### Phase 4: Window Modes + Native Features — Windows Desktop

1. Spirit window (floating mini chat, CompactOverlay presenter)
2. Bubble window (always-on-top circular overlay)
3. Global hotkey (Win+Shift+O via RegisterHotKey)
4. Screenshot capture (Windows.Graphics.Capture)
5. Windows permissions flow
6. Auto-updater (MSIX + GitHub Releases)
7. Windows Hello integration

### Phase 5: AI Awareness + Voice — Windows Desktop

1. AI Vision plugin (Claude Vision API + Windows OCR)
2. Screenshot plugin (Windows.Graphics.Capture)
3. Voice STT/TTS plugin (Windows.Media.Speech)
4. Toast Notifications plugin (AppNotificationManager)
5. Browser context plugin (Chrome extension bridge)
6. UI Automation plugin (screen reader)

### Phase 6: Extended Platforms

1. Windows 11 Widgets: Adaptive Cards for project status, sprint progress
2. Share Target: receive text/URLs from other apps into AI chat
3. Xbox: Dashboard with gamepad navigation, sprint burndown
4. HoloLens: Spatial workspace, holographic multi-window layout
5. MSIX distribution via Microsoft Store and winget
6. Cortana / Windows Search integration

---

## Appendix A: MCP Tools for Windows App

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

## Appendix B: NuGet Package Reference

| Package | Version | Purpose |
|---------|---------|---------|
| `Microsoft.WindowsAppSDK` | 1.5+ | WinUI 3 framework |
| `Google.Protobuf` | 3.25+ | Protobuf serialization |
| `Microsoft.Data.Sqlite` | 8.0+ | Local SQLite cache |
| `H.NotifyIcon.WinUI` | 2.1+ | System tray icon |
| `CommunityToolkit.WinUI.UI` | 8.0+ | UI helpers and controls |
| `CommunityToolkit.Mvvm` | 8.2+ | MVVM toolkit (ObservableObject, RelayCommand) |
| `Microsoft.Extensions.DependencyInjection` | 8.0+ | DI container |
| `Microsoft.Extensions.Hosting` | 8.0+ | Generic host |
| `xUnit` | 2.7+ | Unit testing |
| `Moq` | 4.20+ | Mocking framework |

## Appendix C: Color Parity

| Swift App | Windows App | Note |
|-----------|------------|------|
| `Color(hex: "#0a0d14")` | `#0A0D14` (BackgroundColor) | Identical |
| `Color(hex: "#111520")` | `#111520` (SurfaceColor) | Identical |
| `Color(hex: "#a900ff")` | `#A900FF` (AccentColor) | Identical |
| `Color(hex: "#e8ecf4")` | `#E8ECF4` (TextPrimaryColor) | Identical |
| SF Symbols | Segoe Fluent Icons | Different icon set, same semantics |
| `.monospaced` system font | Cascadia Code | Platform-native monospace |
| Segoe UI Variable (default) | Segoe UI Variable | Same on Windows |

## Appendix D: Keyboard Shortcuts

| Action | macOS (Artifact 18) | Windows (This Document) |
|--------|---------------------|------------------------|
| Cycle window mode | `Cmd+Shift+O` | `Win+Shift+O` |
| Command palette | `Cmd+K` | `Ctrl+K` |
| New chat | `Cmd+N` | `Ctrl+N` |
| Settings | `Cmd+,` | `Ctrl+,` |
| Search | `Cmd+F` | `Ctrl+F` |
| Close window | `Cmd+W` | `Ctrl+W` / `Alt+F4` |
| Quit | `Cmd+Q` | `Alt+F4` (main window) |
| Toggle sidebar | `Cmd+B` | `Ctrl+B` |

## Appendix E: Three-Platform Parity

Cross-reference of how each major feature is implemented across all three native desktop targets.

| Feature | Swift / macOS (Art. 18) | Windows (This Doc.) | Linux / GNOME (Art. 21) |
|---------|------------------------|---------------------|------------------------|
| **UI Framework** | SwiftUI | WinUI 3 | GTK4 + libadwaita |
| **Language** | Swift 5.9+ | C# / .NET 8 | Vala |
| **QUIC** | Network.framework | System.Net.Quic | ngtcp2 (C) |
| **Protobuf** | swift-protobuf | Google.Protobuf | protobuf-c |
| **Plugin interface** | `OrchestraPlugin` (protocol) | `IOrchestraPlugin` (interface) | `OrchestraPlugin` (interface) |
| **Plugin registry** | `PluginRegistry` (@MainActor) | `PluginRegistry` (singleton) | `PluginRegistry` (singleton) |
| **Navigation** | NavigationSplitView | NavigationView | AdwNavigationView |
| **Theme tokens** | `extension Color` | `ResourceDictionary` | CSS custom properties |
| **Floating window** | NSPanel (.floating level) | Win32 WS_EX_TOOLWINDOW | GtkWindow (floating hint) |
| **System tray** | MenuBarExtra | H.NotifyIcon.WinUI | StatusNotifierItem (libdbusmenu) |
| **Global hotkey** | NSEvent global monitor | RegisterHotKey (Win32) | keybinder-3.0 |
| **Credential storage** | Keychain | Windows Credential Manager | libsecret (Keyring/KWallet) |
| **Screenshot** | ScreenCaptureKit | Windows.Graphics.Capture | xdg-desktop-portal / PipeWire |
| **Accessibility** | AXUIElement | UI Automation | AT-SPI2 |
| **TTS/STT** | AVSpeechSynthesizer / SFSpeechRecognizer | Windows.Media.Speech | speech-dispatcher + Vosk |
| **Notifications** | UNUserNotificationCenter | AppNotificationManager | libnotify + GNotification |
| **Packaging** | App Store / notarization | MSIX / Microsoft Store / winget | Flatpak / Snap / deb / rpm |
| **Auto-update** | UpdaterService (zip download) | PackageManager (MSIX) | Flatpak auto-update |
| **Local cache** | SQLite (local) | Microsoft.Data.Sqlite | SQLite (local) |
| **Markdown render** | WKWebView | WebView2 + Monaco | WebKitGTK |
| **Mono font** | SF Mono / system mono | Cascadia Code | JetBrains Mono / system mono |
| **DI** | SwiftUI EnvironmentObject | Microsoft.Extensions.DI | Manual / GObject signals |
| **Testing** | XCTest / Swift Testing | xUnit / Moq | GTest / Vala test |
| **CI** | GitHub Actions (macos-latest) | GitHub Actions (windows-latest) | GitHub Actions (ubuntu-latest) |

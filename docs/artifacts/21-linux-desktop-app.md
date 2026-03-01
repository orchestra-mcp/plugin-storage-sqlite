# 21 — Linux Desktop App: Architecture & Implementation Guide

> Comprehensive reference for building the Orchestra MCP Linux desktop app in GTK4/libadwaita + Vala.
> Supports GNOME, KDE Plasma, XFCE, Sway/Hyprland, and headless from a single codebase.
> Plugin-based architecture mirroring the Go framework and the Swift universal app (artifact 18).
> Compiled from the GTK pack agent, plugin expansion (artifact 20), multi-agent orchestrator (artifact 19), old Linux stubs, and architecture decisions.

---

## Table of Contents

1. [Vision & Context](#1-vision--context)
2. [Architecture](#2-architecture)
3. [App Structure](#3-app-structure)
4. [Vala Plugin System](#4-vala-plugin-system)
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
17. [Native Linux Features](#17-native-linux-features)
18. [Packaging & Distribution](#18-packaging--distribution)
19. [Build Phases](#19-build-phases)

---

## 1. Vision & Context

### What We're Building

A native Linux desktop app for Orchestra MCP — an AI-agentic IDE that manages projects, features, sprints, and AI chat sessions. The app connects to the Orchestra plugin ecosystem via QUIC and exposes 270+ MCP tools through a native GTK4/libadwaita interface. It is the Linux counterpart to the Swift universal Apple app (artifact 18).

### Why Vala + GTK4

| Factor | Vala | C + GTK4 | Python + PyGObject |
|--------|------|----------|-------------------|
| **GObject native** | Compiles to C/GObject | Manual boilerplate | Runtime binding |
| **Type safety** | Strong (C#-like syntax) | Manual | Dynamic |
| **Performance** | Native binary (~same as C) | Native binary | Interpreted |
| **Dev speed** | Fast (modern syntax) | Slow (verbose) | Fastest |
| **GNOME ecosystem** | First-class citizen | First-class citizen | First-class citizen |
| **Binary size** | ~5-15MB | ~3-10MB | ~50MB+ (bundled) |
| **Packaging** | Flatpak/deb/rpm | Flatpak/deb/rpm | Flatpak/pip |

**Decision**: Vala as primary language. It compiles to C via `valac`, produces native GObject binaries, has modern syntax (closures, generics, async/await, null safety), and is the language of choice for GNOME apps (GNOME Builder, Geary, Shotwell).

### Desktop Environment Matrix

| Environment | Toolkit | Wayland | X11 | Tray | Notifications | Secrets |
|-------------|---------|---------|-----|------|---------------|---------|
| **GNOME 45+** | GTK4/libadwaita | Native | XWayland | SNI (extension) | GNotification | GNOME Keyring |
| **KDE Plasma 6** | GTK4 (compat) | Native | XWayland | StatusNotifierItem | D-Bus | KDE Wallet |
| **XFCE 4.18+** | GTK4 (native) | Partial | Native | System tray | D-Bus | GNOME Keyring |
| **Sway/Hyprland** | GTK4 | Native | — | Waybar module | D-Bus (mako/dunst) | libsecret |
| **COSMIC** | GTK4 (compat) | Native | XWayland | applets | D-Bus | libsecret |
| **Headless/SSH** | — | — | — | — | — | — |

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

### Key UI Patterns (Parity with Swift App)

- **Dark theme** with deep navy/black backgrounds (GNOME dark preference)
- **Purple accent** (#a900ff) for active states, brand elements, CTAs
- **Left sidebar** with `AdwNavigationSplitView` for section switching
- **List/detail** pattern with `AdwNavigationPage` for sub-navigation
- **User profile** at bottom-left with avatar, name, status
- **Green "Live" indicators** for active connections
- **Badge counts** on tool results
- **Model badges** in purple pills

---

## 2. Architecture

### Plugin Host Pattern

The Linux app is a consumer in the star-topology orchestrator architecture. Identical to the Swift app — connects to independently running plugins via QUIC.

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

                                    orchestra-linux ← THIS APP
                                    (GTK4/libadwaita)
```

### Communication Contract

- **Transport**: QUIC via `ngtcp2` (reference C QUIC) or `quiche` (Cloudflare C API)
- **Auth**: mTLS — CA at `~/.orchestra/certs/ca.crt`, app cert signed by CA
- **Wire format**: Length-delimited Protobuf (4-byte big-endian uint32 length + Protobuf bytes)
- **Proto contract**: Generated Vala bindings from `protobuf-c` + Vala VAPI
- **Message routing**: All messages go through orchestrator — never direct plugin-to-plugin
- **Streaming**: StreamStart → StreamChunk* → StreamEnd (for AI chat, long-running ops)
- **Events**: Subscribe/Publish/EventDelivery (for real-time updates)
- **Tool calls**: Send `ToolRequest` → receive `ToolResponse` (with optional `provider` field for AI routing)

### Graceful Degradation

If the orchestrator is not running, the app should:
1. Show a "Not Connected" status in the header bar
2. Allow browsing locally cached data (SQLite)
3. Retry connection with exponential backoff (1s → 30s max)
4. Auto-reconnect when orchestrator becomes available
5. Display toast notification on reconnection

---

## 3. App Structure

### Project Layout

```
apps/linux/                                # github.com/orchestra-mcp/orchestra-linux
├── meson.build                           # Meson root build file
├── meson_options.txt                     # Build options (profile, quic backend)
├── orchestra-kit/                        # Core library (transport, models, services)
│   ├── meson.build
│   └── src/
│       ├── transport/
│       │   ├── quic-connection.vala      # QUIC client (ngtcp2 or quiche)
│       │   ├── stream-framer.vala        # Length-delimited Protobuf framing
│       │   └── mtls-config.vala          # mTLS cert loading via GnuTLS/OpenSSL
│       ├── proto/
│       │   ├── messages.vapi            # Vala bindings for protobuf-c types
│       │   └── generated/              # protoc-c generated .c/.h files
│       ├── models/
│       │   ├── app-state.vala           # GObject root state
│       │   ├── project.vala             # Project model
│       │   ├── feature.vala             # Feature/task model
│       │   ├── note.vala                # Note model
│       │   ├── chat-session.vala        # Chat session model
│       │   └── chat-message.vala        # Chat message model
│       ├── plugins/
│       │   ├── orchestra-plugin.vala    # Plugin interface + AppSection
│       │   └── plugin-registry.vala     # Plugin registry + discovery
│       └── services/
│           ├── orchestra-client.vala    # High-level orchestrator client
│           ├── tool-service.vala        # MCP tool call proxy
│           └── connection-state.vala    # Connection status enum
├── shared/                              # Shared UI components
│   ├── meson.build
│   └── src/
│       ├── app/
│       │   └── content-view.vala        # Root view (plugin-driven, DE-adaptive)
│       ├── plugins/                     # Built-in plugins (each self-contained)
│       │   ├── chat-plugin/
│       │   │   ├── chat-plugin.vala     # Plugin registration
│       │   │   ├── chat-view.vala       # Chat layout
│       │   │   ├── chat-session-list.vala
│       │   │   └── chat-message-row.vala
│       │   ├── projects-plugin/
│       │   │   ├── projects-plugin.vala
│       │   │   ├── projects-view.vala
│       │   │   └── project-detail-view.vala
│       │   ├── notes-plugin/
│       │   │   ├── notes-plugin.vala
│       │   │   └── notes-view.vala
│       │   ├── devtools-plugin/
│       │   │   ├── devtools-plugin.vala
│       │   │   └── devtools-view.vala
│       │   └── settings-plugin/
│       │       ├── settings-plugin.vala
│       │       └── settings-view.vala
│       ├── widgets/
│       │   ├── status-badge.vala
│       │   ├── connection-indicator.vala
│       │   └── empty-state-view.vala
│       └── theme/
│           ├── orchestra-theme.vala     # Color tokens + CSS provider
│           └── style.css               # GTK4 CSS overrides
├── desktop/                            # GNOME desktop app entry point
│   ├── meson.build
│   └── src/
│       └── main.vala                   # AdwApplication + window setup
├── data/
│   ├── dev.orchestra.desktop.desktop.in  # .desktop entry
│   ├── dev.orchestra.desktop.metainfo.xml.in  # AppStream metadata
│   ├── dev.orchestra.desktop.gschema.xml     # GSettings schema
│   └── icons/
│       ├── hicolor/
│       │   ├── scalable/apps/dev.orchestra.desktop.svg
│       │   └── symbolic/apps/dev.orchestra.desktop-symbolic.svg
│       └── logo-36.svg
├── po/                                # Translations (gettext)
│   ├── POTFILES
│   ├── LINGUAS
│   └── *.po
├── flatpak/
│   ├── dev.orchestra.desktop.yml      # Flatpak manifest
│   └── dev.orchestra.desktop.Devel.yml
├── snap/
│   └── snapcraft.yaml                 # Snap package
├── debian/
│   ├── control
│   ├── rules
│   └── changelog
├── scripts/
│   ├── new-linux-plugin.sh           # Plugin creator script
│   └── build-release.sh             # Release build helper
└── README.md
```

### App Lifecycle

```vala
public class Orchestra.Application : Adw.Application {
    private Orchestra.PluginRegistry registry;
    private Orchestra.AppState state;

    public Application () {
        Object (
            application_id: "dev.orchestra.desktop",
            flags: ApplicationFlags.DEFAULT_FLAGS
        );
    }

    construct {
        // Register actions
        ActionEntry[] action_entries = {
            { "about", on_about_action },
            { "preferences", on_preferences_action },
            { "quit", quit },
            { "cycle-window", on_cycle_window },
            { "search", on_search },
        };
        add_action_entries (action_entries, this);

        // Keyboard shortcuts
        set_accels_for_action ("app.quit", { "<Control>q" });
        set_accels_for_action ("app.preferences", { "<Control>comma" });
        set_accels_for_action ("app.cycle-window", { "<Control><Shift>o" });
        set_accels_for_action ("app.search", { "<Control>k" });
    }

    protected override void startup () {
        base.startup ();

        // Initialize state
        state = new AppState ();

        // Register built-in plugins
        registry = new PluginRegistry ();
        registry.register (new ChatPlugin ());
        registry.register (new ProjectsPlugin ());
        registry.register (new NotesPlugin ());
        registry.register (new DevToolsPlugin ());
        registry.register (new SettingsPlugin ());

        // Load theme
        load_orchestra_theme ();
    }

    protected override void activate () {
        var window = active_window ?? new Orchestra.Window (this, state, registry);
        window.present ();

        // Connect to orchestrator (async)
        state.client.connect_async.begin ("localhost", 50100, "ui.linux", (obj, res) => {
            try {
                state.client.connect_async.end (res);
            } catch (Error e) {
                warning ("Orchestrator connection failed: %s", e.message);
            }
        });
    }
}

int main (string[] args) {
    var app = new Orchestra.Application ();
    return app.run (args);
}
```

### Startup Sequence

```
1. Load GSettings (theme, window geometry, last plugin)
2. Initialize theme (CSS provider + color tokens)
3. Register built-in plugins with PluginRegistry
4. Setup UI (AdwApplicationWindow + NavigationSplitView)
5. Connect to orchestrator via QUIC
6. Subscribe to events (feature.*, workflow.*, sprint.*, note.*)
7. Load cached data from SQLite (~/.local/share/orchestra/cache.db)
8. Check for updates (Flatpak: automatic, deb/rpm: 6-hour check)
9. Request portal permissions (first run: notifications, background)
10. Restore window geometry (GSettings: width, height, maximized)
```

---

## 4. Vala Plugin System

### Design Principles

Mirrors the Go framework's plugin architecture and the Swift app's plugin system. Every screen/feature is a Vala plugin that registers with the PluginRegistry. This enables:
- **Incremental development** — build one plugin at a time
- **Easy extensibility** — `new-linux-plugin.sh` scaffolds new plugins
- **Feature isolation** — each plugin is self-contained with its own views
- **DE adaptation** — plugins can query the desktop environment

### Plugin Interface

```vala
/** Section where the plugin appears in the UI. */
public enum Orchestra.AppSection {
    SIDEBAR,    // Main navigation (Chat, Projects, Notes)
    DEVTOOLS,   // Developer Tools sub-section
    SETTINGS;   // Settings sub-section
}

/** Every feature in the app conforms to this interface. */
public interface Orchestra.Plugin : Object {
    public abstract string id { get; }
    public abstract string name { get; }
    public abstract string icon_name { get; }      // Freedesktop icon name
    public abstract AppSection section { get; }
    public abstract int order { get; }

    public abstract Gtk.Widget create_view ();
    public abstract void on_activate ();
    public abstract void on_deactivate ();
}
```

### Plugin Registry

```vala
public class Orchestra.PluginRegistry : Object {
    private GenericArray<Plugin> _plugins = new GenericArray<Plugin> ();

    public signal void plugin_added (Plugin plugin);

    public void register (Plugin plugin) {
        _plugins.add (plugin);
        _plugins.sort ((a, b) => a.order - b.order);
        plugin_added (plugin);
    }

    public Plugin? get_plugin (string id) {
        for (int i = 0; i < _plugins.length; i++) {
            if (_plugins[i].id == id) return _plugins[i];
        }
        return null;
    }

    public GenericArray<Plugin> sidebar_plugins {
        owned get {
            var result = new GenericArray<Plugin> ();
            for (int i = 0; i < _plugins.length; i++) {
                if (_plugins[i].section == AppSection.SIDEBAR)
                    result.add (_plugins[i]);
            }
            return result;
        }
    }

    public GenericArray<Plugin> devtools_plugins {
        owned get {
            var result = new GenericArray<Plugin> ();
            for (int i = 0; i < _plugins.length; i++) {
                if (_plugins[i].section == AppSection.DEVTOOLS)
                    result.add (_plugins[i]);
            }
            return result;
        }
    }

    public GenericArray<Plugin> settings_plugins {
        owned get {
            var result = new GenericArray<Plugin> ();
            for (int i = 0; i < _plugins.length; i++) {
                if (_plugins[i].section == AppSection.SETTINGS)
                    result.add (_plugins[i]);
            }
            return result;
        }
    }
}
```

### Built-in Plugins (Phase 1)

| Plugin | ID | Section | Icon |
|--------|----|---------|------|
| ChatPlugin | `chat` | sidebar | `chat-message-new-symbolic` |
| ProjectsPlugin | `projects` | sidebar | `view-grid-symbolic` |
| NotesPlugin | `notes` | sidebar | `accessories-text-editor-symbolic` |
| DevToolsPlugin | `devtools` | sidebar | `utilities-terminal-symbolic` |
| SettingsPlugin | `settings` | sidebar | `preferences-system-symbolic` |

### Plugin Creator Script

`scripts/new-linux-plugin.sh` generates:

```bash
./scripts/new-linux-plugin.sh my-feature sidebar
# Creates:
# shared/src/plugins/my-feature-plugin/
#   ├── my-feature-plugin.vala    # Plugin registration
#   └── my-feature-view.vala      # Main view
# Updates meson.build sources
# Prints registration line to add to main.vala
```

---

## 5. Navigation & Screens

### AdwNavigationSplitView Layout

```vala
public class Orchestra.Window : Adw.ApplicationWindow {
    private Adw.NavigationSplitView split_view;
    private Gtk.ListBox sidebar_list;
    private Adw.NavigationPage content_page;
    private PluginRegistry registry;
    private string? selected_plugin_id = "chat";

    public Window (Adw.Application app, AppState state, PluginRegistry registry) {
        Object (application: app, default_width: 1280, default_height: 860);
        this.registry = registry;

        // Build split view
        split_view = new Adw.NavigationSplitView ();

        // Sidebar
        var sidebar_page = new Adw.NavigationPage.with_tag ("sidebar", _("Orchestra"));
        var sidebar_toolbar = new Adw.ToolbarView ();
        var header = new Adw.HeaderBar ();
        sidebar_toolbar.add_top_bar (header);

        sidebar_list = new Gtk.ListBox ();
        sidebar_list.set_selection_mode (Gtk.SelectionMode.SINGLE);
        sidebar_list.add_css_class ("navigation-sidebar");
        sidebar_list.row_selected.connect (on_sidebar_row_selected);
        populate_sidebar ();
        sidebar_toolbar.set_content (sidebar_list);
        sidebar_page.set_child (sidebar_toolbar);

        // Content
        content_page = new Adw.NavigationPage.with_tag ("content", _("Chat"));
        split_view.set_sidebar (sidebar_page);
        split_view.set_content (content_page);
        set_content (split_view);

        // Show initial plugin
        show_plugin ("chat");
    }

    private void populate_sidebar () {
        var plugins = registry.sidebar_plugins;
        for (int i = 0; i < plugins.length; i++) {
            var row = create_sidebar_row (plugins[i]);
            sidebar_list.append (row);
        }
    }

    private Gtk.Widget create_sidebar_row (Plugin plugin) {
        var row = new Adw.ActionRow ();
        row.set_title (plugin.name);
        var icon = new Gtk.Image.from_icon_name (plugin.icon_name);
        row.add_prefix (icon);
        row.set_data<string> ("plugin-id", plugin.id);
        return row;
    }

    private void show_plugin (string id) {
        var plugin = registry.get_plugin (id);
        if (plugin == null) return;

        // Deactivate old
        if (selected_plugin_id != null) {
            var old = registry.get_plugin (selected_plugin_id);
            if (old != null) old.on_deactivate ();
        }

        selected_plugin_id = id;
        plugin.on_activate ();

        var toolbar = new Adw.ToolbarView ();
        toolbar.add_top_bar (new Adw.HeaderBar ());
        toolbar.set_content (plugin.create_view ());
        content_page.set_child (toolbar);
        content_page.set_title (plugin.name);
    }
}
```

### Sidebar (Icon Rail Style)

| # | Icon | Section | Route |
|---|------|---------|-------|
| 1 | `chat-message-new-symbolic` | Chat | `/chat` |
| 2 | `view-grid-symbolic` | Projects | `/projects` |
| 3 | `accessories-text-editor-symbolic` | Notes | `/notes` |
| 4 | `utilities-terminal-symbolic` | Developer Tools | `/devtools` |
| — | *(spacer)* | | |
| 5 | `preferences-system-symbolic` | Settings | `/settings` |

**Active state**: Purple accent background (`@accent_bg_color`), white icon.
**Inactive state**: `@window_fg_color`, transparent background.
**Brand logo**: 36x36 SVG at top of sidebar.
**Connection indicator**: Green/red dot in header bar.
**User profile**: Avatar + name at bottom via `AdwAvatar`.

### Cross-Platform Feature Matrix

| Feature | GNOME | KDE | XFCE | Sway/Hyprland | Headless |
|---------|-------|-----|------|---------------|----------|
| Chat | Full | Full | Full | Full | CLI-only |
| Projects | Full | Full | Full | Full | CLI-only |
| Notes | Full | Full | Full | Full | CLI-only |
| DevTools | Full (10 tools) | Full | Full | Full | CLI-only |
| System Tray | SNI extension | StatusNotifierItem | System tray | Waybar module | — |
| Window Modes | 3 modes | 3 modes | 2 modes | 2 modes | — |
| Screenshot | GNOME Screenshot / xdg-portal | Spectacle / xdg-portal | xfce4-screenshooter | grim + slurp | — |
| Voice STT/TTS | PipeWire | PipeWire | PulseAudio | PipeWire | — |
| Widgets | GNOME Extensions | Plasma Widgets | — | — | — |
| Global Hotkey | GNOME shortcut / xdg-portal | KDE shortcut | xfce4-hotkey | swaymsg | — |
| Notifications | GNotification | D-Bus | D-Bus | mako/dunst | — |

---

## 6. Design System

### Color Tokens (GTK4 CSS)

```css
/* Orchestra theme — GTK4 CSS */
@define-color orchestra_bg         #0a0d14;
@define-color orchestra_surface    #111520;
@define-color orchestra_contrast   #080a10;
@define-color orchestra_active     #1a1f2e;
@define-color orchestra_selection  rgba(169, 0, 255, 0.15);

@define-color orchestra_fg         #e8ecf4;
@define-color orchestra_fg_muted   #8892a8;
@define-color orchestra_fg_dim     #4a5268;
@define-color orchestra_fg_bright  #f8fafc;

@define-color orchestra_border     #1e2436;
@define-color orchestra_accent     #a900ff;

@define-color orchestra_success    #22c55e;
@define-color orchestra_warning    #f59e0b;
@define-color orchestra_error      #ef4444;
@define-color orchestra_info       #3b82f6;

/* Syntax highlighting (code views) */
@define-color syntax_blue          #82aaff;
@define-color syntax_cyan          #89ddff;
@define-color syntax_green         #c3e88d;
@define-color syntax_yellow        #ffcb6b;
@define-color syntax_orange        #f78c6c;
@define-color syntax_red           #ff5370;
@define-color syntax_purple        #c792ea;
```

### Vala Theme Loader

```vala
public class Orchestra.Theme : Object {
    private Gtk.CssProvider provider;
    private string current_theme = "orchestra";

    public void load (Adw.Application app) {
        provider = new Gtk.CssProvider ();
        provider.load_from_resource ("/dev/orchestra/desktop/style.css");

        Gtk.StyleContext.add_provider_for_display (
            Gdk.Display.get_default (),
            provider,
            Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION
        );

        // Force dark color scheme
        var style_manager = Adw.StyleManager.get_default ();
        style_manager.set_color_scheme (Adw.ColorScheme.FORCE_DARK);
    }

    public void set_accent (string hex_color) {
        // libadwaita 1.6+ supports named accent colors
        // For older versions, inject via CSS
        var css = "@define-color accent_bg_color %s;\n".printf (hex_color);
        css += "@define-color accent_color %s;\n".printf (hex_color);
        var accent_provider = new Gtk.CssProvider ();
        accent_provider.load_from_string (css);
        Gtk.StyleContext.add_provider_for_display (
            Gdk.Display.get_default (),
            accent_provider,
            Gtk.STYLE_PROVIDER_PRIORITY_USER
        );
    }
}
```

### Typography

```css
/* Orchestra typography */
.body-default { font-size: 14px; }
.body-secondary { font-size: 13px; }
.label-text { font-size: 12px; font-weight: 600; }
.caption-text { font-size: 11px; }
.section-title { font-size: 16px; font-weight: 600; }
.monospace { font-family: "JetBrains Mono", "Fira Code", "Source Code Pro", monospace; }
```

### Spacing & Geometry

| Element | Default | Compact |
|---------|---------|---------|
| Sidebar width | 240px (split) | 200px |
| Sidebar icon | 24x24px | 20x20px |
| Header bar height | 47px (libadwaita) | 47px |
| Content padding | 24px | 16px |
| Corner radius (sm) | 6px | 4px |
| Corner radius (md) | 12px | 8px |
| Corner radius (lg) | 16px | 12px |

### 25 Color Themes

All 25 themes from `@orchestra-mcp/theme` are supported via GTK4 CSS providers. Default: `orchestra` (deep navy).

| ID | Background | Accent | Dark? |
|----|-----------|--------|-------|
| `orchestra` | #0a0d14 | #a900ff | Yes |
| `dracula` | #282a36 | #bd93f9 | Yes |
| `github-dark` | #0d1117 | #58a6ff | Yes |
| `github-light` | #ffffff | #0366d6 | No |
| `one-dark` | #282c34 | #528bff | Yes |
| `monokai-pro` | #2d2a2e | #ffd866 | Yes |
| `synthwave-84` | #262335 | #ff7edb | Yes |

Theme switching loads a new `Gtk.CssProvider` per theme and updates `Adw.StyleManager.color_scheme` for light themes.

---

## 7. Window Management

### Three Window Modes

| Mode | Window | Size | Behavior |
|------|--------|------|----------|
| **Main** | Primary window | 1280x860 | Full IDE, all sections |
| **Floating** | Spirit window | 420x640 | Always-on-top via `GDK_TOPLEVEL_STATE_ABOVE`, mini chat |
| **Bubble** | Bubble overlay | 56x56 | Layer-shell overlay (Wayland) or `_NET_WM_STATE_ABOVE` (X11) |

**Cycling**: `Ctrl+Shift+O` global shortcut cycles modes.

### Spirit Window (Floating Mini Chat)

```vala
public class Orchestra.SpiritWindow : Adw.Window {
    public SpiritWindow (Adw.Application app) {
        Object (application: app, default_width: 420, default_height: 640);

        // Semi-transparent background
        add_css_class ("spirit-window");

        // Title bar
        titlebar = new Adw.HeaderBar () {
            show_title = false
        };
        var close_btn = new Gtk.Button.from_icon_name ("window-close-symbolic");
        close_btn.clicked.connect (() => hide ());
        ((Adw.HeaderBar) titlebar).pack_end (close_btn);

        // Chat view (compact mode)
        var chat = new ChatView.compact ();
        var toolbar = new Adw.ToolbarView ();
        toolbar.add_top_bar ((Adw.HeaderBar) titlebar);
        toolbar.set_content (chat);
        set_content (toolbar);
    }

    public override void realize () {
        base.realize ();

        // Make always-on-top
        var surface = get_surface () as Gdk.ToplevelSurface;
        if (surface != null) {
            // Wayland: layer-shell (via gtk4-layer-shell)
            // X11: _NET_WM_STATE_ABOVE via GDK
        }
    }
}
```

### Spirit CSS

```css
.spirit-window {
    background-color: alpha(@orchestra_bg, 0.92);
    border: 1px solid @orchestra_border;
    border-radius: 16px;
}
```

### Wayland Layer Shell (for Bubble/Overlay)

```vala
// Uses gtk4-layer-shell for Wayland compositors
// Fallback to X11 window hints on X11 sessions

#if HAVE_LAYER_SHELL
using GtkLayerShell;

public void setup_bubble_layer (Gtk.Window window) {
    GtkLayerShell.init_for_window (window);
    GtkLayerShell.set_layer (window, GtkLayerShell.Layer.OVERLAY);
    GtkLayerShell.set_anchor (window, GtkLayerShell.Edge.BOTTOM, true);
    GtkLayerShell.set_anchor (window, GtkLayerShell.Edge.RIGHT, true);
    GtkLayerShell.set_margin (window, GtkLayerShell.Edge.BOTTOM, 20);
    GtkLayerShell.set_margin (window, GtkLayerShell.Edge.RIGHT, 20);
    GtkLayerShell.set_exclusive_zone (window, -1);
}
#endif
```

---

## 8. QUIC Transport

### QUIC Client via ngtcp2

```vala
/**
 * QUIC connection to orchestrator using ngtcp2 (C library).
 * Vala calls C functions directly via VAPI binding.
 */
public class Orchestra.QUICConnection : Object {
    public signal void state_changed (ConnectionState state);
    public signal void message_received (uint8[] data);

    private ConnectionState _state = ConnectionState.DISCONNECTED;
    public ConnectionState state {
        get { return _state; }
        private set {
            _state = value;
            state_changed (value);
        }
    }

    /**
     * Connect to orchestrator.
     * @param host Orchestrator hostname
     * @param port QUIC port (default 50100)
     * @param plugin_id This app's plugin ID
     */
    public async void connect_async (string host, uint16 port, string plugin_id) throws Error {
        state = ConnectionState.CONNECTING;

        // Load mTLS certificates
        var certs_dir = Path.build_filename (
            Environment.get_home_dir (), ".orchestra", "certs"
        );
        var ca_cert = Path.build_filename (certs_dir, "ca.crt");
        var plugin_cert = Path.build_filename (certs_dir, "%s.crt".printf (plugin_id));
        var plugin_key = Path.build_filename (certs_dir, "%s.key".printf (plugin_id));

        // Initialize ngtcp2 via C binding (see quic-native.c)
        var result = yield quic_native_connect (host, port, ca_cert, plugin_cert, plugin_key);
        if (!result) {
            state = ConnectionState.ERROR;
            throw new TransportError.CONNECTION_FAILED ("QUIC connection to %s:%u failed", host, port);
        }

        state = ConnectionState.CONNECTED;

        // Start read loop
        read_loop.begin ();
    }

    public async void send (uint8[] data) throws Error {
        if (state != ConnectionState.CONNECTED) {
            throw new TransportError.NOT_CONNECTED ("Not connected to orchestrator");
        }
        yield stream_framer_write (data);
    }

    private async void read_loop () {
        while (state == ConnectionState.CONNECTED) {
            try {
                var data = yield stream_framer_read ();
                message_received (data);
            } catch (Error e) {
                warning ("Read error: %s", e.message);
                state = ConnectionState.DISCONNECTED;
                reconnect_with_backoff.begin ();
                break;
            }
        }
    }

    private async void reconnect_with_backoff () {
        uint delay = 1000; // 1s
        while (state == ConnectionState.DISCONNECTED) {
            state = ConnectionState.RECONNECTING;
            try {
                yield connect_async (_host, _port, _plugin_id);
                return;
            } catch {
                yield async_sleep (delay);
                delay = uint.min (delay * 2, 30000); // Cap at 30s
            }
        }
    }
}
```

### Length-Delimited Protobuf Framing

```vala
/**
 * Frame protocol: [4-byte big-endian uint32 length][N-byte protobuf payload]
 * Max message size: 16 MB
 */
namespace Orchestra.StreamFramer {
    private const uint32 MAX_MESSAGE_SIZE = 16 * 1024 * 1024;

    public async void write (uint8[] data, OutputStream stream) throws Error {
        // Write 4-byte length header (big-endian)
        uint32 length = (uint32) data.length;
        uint8 header[4];
        header[0] = (uint8) ((length >> 24) & 0xFF);
        header[1] = (uint8) ((length >> 16) & 0xFF);
        header[2] = (uint8) ((length >> 8) & 0xFF);
        header[3] = (uint8) (length & 0xFF);
        yield stream.write_all_async (header, Priority.DEFAULT, null, null);
        yield stream.write_all_async (data, Priority.DEFAULT, null, null);
    }

    public async uint8[] read (InputStream stream) throws Error {
        // Read 4-byte length header
        uint8 header[4];
        yield stream.read_all_async (header, Priority.DEFAULT, null, null);
        uint32 size = (header[0] << 24) | (header[1] << 16) | (header[2] << 8) | header[3];

        if (size > MAX_MESSAGE_SIZE) {
            throw new TransportError.MESSAGE_TOO_LARGE ("Message size %u exceeds 16MB", size);
        }

        // Read payload
        var buf = new uint8[size];
        yield stream.read_all_async (buf, Priority.DEFAULT, null, null);
        return buf;
    }
}
```

### mTLS Configuration

```vala
public class Orchestra.MTLSConfig : Object {
    public string ca_cert_path { get; construct; }
    public string plugin_cert_path { get; construct; }
    public string plugin_key_path { get; construct; }

    public MTLSConfig.for_plugin (string plugin_id) {
        var certs_dir = Path.build_filename (
            Environment.get_home_dir (), ".orchestra", "certs"
        );
        Object (
            ca_cert_path: Path.build_filename (certs_dir, "ca.crt"),
            plugin_cert_path: Path.build_filename (certs_dir, "%s.crt".printf (plugin_id)),
            plugin_key_path: Path.build_filename (certs_dir, "%s.key".printf (plugin_id))
        );
    }

    public GnuTLS.Certificate create_credentials () throws Error {
        // Load CA cert for peer verification
        // Load client cert + key for mTLS identity
        // Returns configured GnuTLS credential set
    }
}
```

### Native C QUIC Bridge

For QUIC transport, a small C file wraps ngtcp2 and exposes simple async functions to Vala:

```c
/* quic-native.c — C bridge for ngtcp2 QUIC transport */
#include <ngtcp2/ngtcp2.h>
#include <ngtcp2/ngtcp2_crypto_gnutls.h>
#include <glib.h>

typedef struct {
    ngtcp2_conn *conn;
    GIOChannel *socket_channel;
    GMainContext *context;
    guint source_id;
    /* callbacks back to Vala */
    void (*on_data)(const uint8_t *data, size_t len, void *user_data);
    void *user_data;
} OrchestraQuicConn;

gboolean orchestra_quic_connect (
    OrchestraQuicConn *self,
    const char *host, uint16_t port,
    const char *ca_path, const char *cert_path, const char *key_path,
    GError **error
);

gboolean orchestra_quic_send (OrchestraQuicConn *self, const uint8_t *data, size_t len, GError **error);
void orchestra_quic_close (OrchestraQuicConn *self);
```

---

## 9. Data Layer

### MCP Tool Proxy

All data operations go through MCP tool calls routed via the orchestrator.

```vala
public class Orchestra.ToolService : Object {
    private OrchestraClient client;

    public ToolService (OrchestraClient client) {
        this.client = client;
    }

    public async Json.Object call_tool (string name, Json.Object arguments) throws Error {
        return yield client.send_tool_call (name, arguments);
    }

    // Convenience wrappers
    public async GenericArray<Project> list_projects () throws Error {
        var result = yield call_tool ("list_projects", new Json.Object ());
        return Project.from_json_array (result);
    }

    public async ProjectStatus get_project_status (string slug) throws Error {
        var args = new Json.Object ();
        args.set_string_member ("project", slug);
        var result = yield call_tool ("get_project_status", args);
        return new ProjectStatus.from_json (result);
    }

    public async Feature create_feature (string project, string title, string type = "task") throws Error {
        var args = new Json.Object ();
        args.set_string_member ("project", project);
        args.set_string_member ("title", title);
        args.set_string_member ("type", type);
        var result = yield call_tool ("create_feature", args);
        return new Feature.from_json (result);
    }

    // Multi-LLM AI calls
    public async string ai_prompt (string prompt, string provider, string model) throws Error {
        var args = new Json.Object ();
        args.set_string_member ("prompt", prompt);
        args.set_string_member ("provider", provider);
        args.set_string_member ("model", model);
        var result = yield call_tool ("ai_prompt", args);
        return result.get_string_member ("content");
    }

    public async ChatResponse spawn_session (string id, string prompt, string provider) throws Error {
        var args = new Json.Object ();
        args.set_string_member ("session_id", id);
        args.set_string_member ("prompt", prompt);
        args.set_string_member ("provider", provider);
        var result = yield call_tool ("spawn_session", args);
        return new ChatResponse.from_json (result);
    }
}
```

### Local Cache (SQLite)

```vala
// ~/.local/share/orchestra/cache.db
public class Orchestra.LocalCache : Object {
    private Sqlite.Database db;

    public LocalCache () throws Error {
        var data_dir = Path.build_filename (
            Environment.get_user_data_dir (), "orchestra"
        );
        DirUtils.create_with_parents (data_dir, 0755);
        var db_path = Path.build_filename (data_dir, "cache.db");

        Sqlite.Database.open (db_path, out db);
        create_tables ();
    }

    private void create_tables () {
        db.exec ("""
            CREATE TABLE IF NOT EXISTS projects (
                slug TEXT PRIMARY KEY,
                data TEXT NOT NULL,
                updated_at INTEGER NOT NULL
            );
            CREATE TABLE IF NOT EXISTS notes (
                id TEXT PRIMARY KEY,
                data TEXT NOT NULL,
                updated_at INTEGER NOT NULL
            );
            CREATE TABLE IF NOT EXISTS sessions (
                id TEXT PRIMARY KEY,
                data TEXT NOT NULL,
                updated_at INTEGER NOT NULL
            );
        """);
    }

    public void cache_projects (GenericArray<Project> projects) { /* ... */ }
    public GenericArray<Project> cached_projects () { /* ... */ }
    public void cache_notes (GenericArray<Note> notes) { /* ... */ }
    public GenericArray<Note> cached_notes () { /* ... */ }
    public void invalidate (string entity) { /* ... */ }
}
```

---

## 10. AI Chat

### Multi-LLM Chat Architecture

The chat interface supports 4 AI providers (Claude, OpenAI, Gemini, Ollama) via the bridge plugins.

```
ChatView (AdwNavigationSplitView)
├── ChatSessionList (240px sidebar)
│   ├── GtkSearchEntry
│   ├── New Chat button (GtkButton)
│   └── GtkListBox (session rows: name, provider badge, model badge, date)
└── ChatBox (AdwNavigationPage)
    ├── ChatHeader (AdwHeaderBar: session name, provider pill, model pill, Live dot)
    ├── ChatBody (GtkListBox: scrollable messages)
    │   ├── ChatMessageRow (user: right-aligned, purple bg)
    │   ├── ChatMessageRow (assistant: left-aligned, surface bg)
    │   ├── EventCards (tool call results in AdwExpanderRow)
    │   └── TypingIndicator (animated dots)
    ├── StatusLine (GtkLabel: typing status + elapsed timer)
    └── ChatInput (AdwClamp)
        ├── GtkTextView (auto-resize, 1-6 lines)
        └── GtkBox (provider picker, model picker, tools toggle, attach, send/stop)
```

### Provider & Model Selection

```vala
public struct Orchestra.AIProvider {
    public string id;
    public string name;
    public string[] models;
}

public const AIProvider[] PROVIDERS = {
    { "claude", "Anthropic", { "claude-opus-4-6", "claude-sonnet-4-6", "claude-haiku-4-5" } },
    { "openai", "OpenAI", { "gpt-4o", "gpt-4o-mini", "o1", "o1-mini" } },
    { "gemini", "Google", { "gemini-2.5-pro", "gemini-2.5-flash", "gemini-2.0-flash" } },
    { "ollama", "Ollama", { "llama3", "codellama", "mistral" } },
};
```

### Message Model

```vala
public enum Orchestra.MessageRole {
    USER,
    ASSISTANT,
    SYSTEM;
}

public class Orchestra.ChatMessage : Object {
    public string id { get; construct; }
    public MessageRole role { get; set; }
    public string content { get; set; }
    public DateTime timestamp { get; construct; }
    public bool streaming { get; set; default = false; }
    public string? thinking { get; set; }
    public GenericArray<ToolEvent> events { get; set; }
    public string? provider { get; set; }
    public string? model { get; set; }
    public GenericArray<Attachment> attachments { get; set; }
}
```

### Streaming Glow Effect (GTK4 CSS)

```css
/* Applied to message row during streaming */
@keyframes glow-rotate {
    0% { border-color: @syntax_blue; }
    25% { border-color: @orchestra_accent; }
    50% { border-color: @syntax_red; }
    75% { border-color: @syntax_purple; }
    100% { border-color: @syntax_blue; }
}

.message-streaming {
    border: 2px solid @orchestra_accent;
    animation: glow-rotate 3s linear infinite;
    border-radius: 16px;
}
```

### Event Cards (Tool Results)

| Card Type | Widget | Renders |
|-----------|--------|---------|
| `BashCard` | `AdwExpanderRow` | Terminal output with command + exit code |
| `ReadCard` | `GtkSourceView` | File content with line numbers |
| `EditCard` | `GtkSourceView` | Diff view (old/new) with syntax highlight |
| `TaskCard` | `AdwActionRow` | Task summary with status/priority badges |
| `ProjectStatusCard` | `GtkLevelBar` | Progress bar + stats |
| `SprintCard` | `AdwPreferencesGroup` | Sprint overview |
| `QuestionCard` | `AdwActionRow` + `GtkButton` | Interactive question with option buttons |

---

## 11. Project Management

### Projects Section

```
ProjectsView (AdwNavigationSplitView)
├── ProjectSidebar (240px)
│   ├── GtkSearchEntry
│   ├── New Project button (+)
│   └── GtkListBox (AdwActionRow per project: icon, name, task count, %, Active badge)
└── ProjectDetail (AdwNavigationPage)
    ├── Header (AdwBanner: icon, name, description)
    ├── Status group (AdwPreferencesGroup)
    │   ├── Total Tasks (AdwActionRow)
    │   ├── Completed (AdwActionRow)
    │   └── Completion % (GtkLevelBar: purple, animated)
    ├── Status breakdown (AdwPreferencesGroup per state)
    └── BacklogTree (GtkTreeListModel: collapsible epics > stories > tasks)
```

### Workflow States (13-State Machine)

```vala
public enum Orchestra.WorkflowState {
    BACKLOG,
    TODO,
    IN_PROGRESS,
    READY_FOR_TESTING,
    IN_TESTING,
    READY_FOR_DOCS,
    IN_DOCS,
    DOCUMENTED,
    IN_REVIEW,
    DONE,
    BLOCKED,
    REJECTED,
    CANCELLED;

    public string to_label () {
        switch (this) {
            case BACKLOG: return _("Backlog");
            case TODO: return _("To Do");
            case IN_PROGRESS: return _("In Progress");
            case READY_FOR_TESTING: return _("Ready for Testing");
            case IN_TESTING: return _("In Testing");
            case READY_FOR_DOCS: return _("Ready for Docs");
            case IN_DOCS: return _("In Docs");
            case DOCUMENTED: return _("Documented");
            case IN_REVIEW: return _("In Review");
            case DONE: return _("Done");
            case BLOCKED: return _("Blocked");
            case REJECTED: return _("Rejected");
            case CANCELLED: return _("Cancelled");
            default: return "";
        }
    }

    public string to_css_class () {
        switch (this) {
            case BACKLOG: case CANCELLED: return "dim-label";
            case TODO: return "accent";
            case IN_PROGRESS: return "orchestra-purple";
            case READY_FOR_TESTING: return "warning";
            case IN_TESTING: return "warning";
            case READY_FOR_DOCS: return "orchestra-cyan";
            case IN_DOCS: return "orchestra-teal";
            case DOCUMENTED: case DONE: return "success";
            case IN_REVIEW: return "orchestra-indigo";
            case BLOCKED: case REJECTED: return "error";
            default: return "";
        }
    }
}
```

---

## 12. Notes & Docs

### Notes Plugin

Uses `tools.notes` plugin (8 tools): `create_note`, `get_note`, `update_note`, `delete_note`, `list_notes`, `search_notes`, `pin_note`, `tag_note`.

```
NotesView (AdwNavigationSplitView)
├── NotesSidebar (240px)
│   ├── GtkSearchEntry
│   ├── New Note button (+)
│   ├── Pinned Notes section (GtkListBox)
│   └── Other Notes section (GtkListBox)
└── NoteEditor (AdwNavigationPage)
    ├── AdwHeaderBar (back, pin toggle, save, delete)
    ├── Title input (GtkEntry: large, borderless)
    ├── Tags bar (GtkFlowBox: add/remove tags)
    └── Content editor (GtkSourceView: monospace, markdown syntax)
```

### GtkSourceView for Markdown

```vala
private GtkSource.View create_editor () {
    var buffer = new GtkSource.Buffer (null);
    var lang_manager = GtkSource.LanguageManager.get_default ();
    buffer.set_language (lang_manager.get_language ("markdown"));

    var scheme_manager = GtkSource.StyleSchemeManager.get_default ();
    buffer.set_style_scheme (scheme_manager.get_scheme ("orchestra-dark"));

    var view = new GtkSource.View.with_buffer (buffer);
    view.set_monospace (true);
    view.set_show_line_numbers (false);
    view.set_wrap_mode (Gtk.WrapMode.WORD_CHAR);
    view.set_top_margin (16);
    view.set_bottom_margin (16);
    view.set_left_margin (24);
    view.set_right_margin (24);
    return view;
}
```

### Docs/Wiki Plugin

Uses `tools.docs` plugin (10 tools): `doc_create`, `doc_get`, `doc_update`, `doc_delete`, `doc_list`, `doc_search`, `doc_generate`, `doc_index`, `doc_tree`, `doc_export`.

Categories: `api-reference`, `guide`, `architecture`, `tutorial`, `changelog`, `decision-record`

---

## 13. Developer Tools

### DevTools Plugin Container

The DevTools section hosts sub-plugins for each tool type. Each is a separate Go plugin on the backend.

| Tool | Plugin ID | Tools | Description |
|------|-----------|-------|-------------|
| File Explorer | devtools.file-explorer | 17 | Full IDE: file ops + LSP code intelligence |
| Terminal | devtools.terminal | 6 | VTE terminal sessions |
| SSH | devtools.ssh | 7 | Remote SSH + SFTP |
| Services | devtools.services | 6 | systemctl service manager |
| Docker | devtools.docker | 10 | Container management |
| Debugger | devtools.debugger | 8 | DAP protocol debugger |
| Test Runner | devtools.test-runner | 6 | Multi-framework test runner |
| Log Viewer | devtools.log-viewer | 5 | journalctl + log streaming |
| Database | devtools.database | 8 | SQL query editor + schema browser |
| DevOps | devtools.devops | 8 | CI/CD pipeline management |

### Terminal (VTE)

Linux terminal uses VTE (GNOME Terminal library), the standard GTK terminal emulator:

```vala
private Vte.Terminal create_terminal () {
    var terminal = new Vte.Terminal ();
    terminal.set_font (Pango.FontDescription.from_string ("JetBrains Mono 13"));
    terminal.set_color_background (Gdk.RGBA () { red = 0.04, green = 0.05, blue = 0.08, alpha = 1.0 });
    terminal.set_color_foreground (Gdk.RGBA () { red = 0.91, green = 0.93, blue = 0.96, alpha = 1.0 });
    terminal.set_scrollback_lines (10000);

    // Spawn shell
    terminal.spawn_async (
        Vte.PtyFlags.DEFAULT,
        Environment.get_home_dir (),
        { Environment.get_variable ("SHELL") ?? "/bin/bash" },
        null, // env
        SpawnFlags.DEFAULT,
        null, null, -1, null, null
    );

    return terminal;
}
```

### File Explorer (GtkSourceView)

```vala
// File tree uses GtkTreeListModel + GtkColumnView
private Gtk.Widget create_file_explorer () {
    var root = Gio.File.new_for_path (workspace_path);
    var model = new Gtk.TreeListModel (
        create_directory_model (root),
        false, true,
        (item) => create_directory_model (((FileInfo) item).file)
    );

    var selection = new Gtk.SingleSelection (model);
    var view = new Gtk.ListView (selection, create_file_row_factory ());
    return view;
}
```

### Service Manager (systemd)

Linux-specific — uses systemd D-Bus API:

```vala
public async GenericArray<ServiceInfo> list_services () throws Error {
    var bus = yield Bus.get (BusType.SYSTEM);
    var reply = yield bus.call (
        "org.freedesktop.systemd1",
        "/org/freedesktop/systemd1",
        "org.freedesktop.systemd1.Manager",
        "ListUnits",
        null, null, DBusCallFlags.NONE, -1
    );
    // Parse unit list into ServiceInfo objects
}
```

---

## 14. AI Awareness

Four AI awareness plugins provide visual and contextual understanding.

### ai.screenshot (6 tools)

`capture_screen`, `capture_region`, `capture_window`, `capture_interactive`, `annotate_screenshot`, `list_captures`

**Implementation strategy per session type:**

| Session | Tool | Fallback |
|---------|------|----------|
| GNOME Wayland | xdg-desktop-portal Screenshot | gnome-screenshot |
| KDE Wayland | xdg-desktop-portal Screenshot | spectacle |
| Sway/Hyprland | grim + slurp | — |
| X11 (any DE) | xdg-desktop-portal Screenshot | scrot / import (ImageMagick) |

```vala
public async Bytes capture_screen () throws Error {
    // Prefer xdg-desktop-portal (works everywhere, respects Wayland security)
    try {
        return yield portal_screenshot ();
    } catch {
        // Fallback to CLI tools
        return yield cli_screenshot ();
    }
}

private async Bytes portal_screenshot () throws Error {
    var bus = yield Bus.get (BusType.SESSION);
    var reply = yield bus.call (
        "org.freedesktop.portal.Desktop",
        "/org/freedesktop/portal/desktop",
        "org.freedesktop.portal.Screenshot",
        "Screenshot",
        new Variant ("(sa{sv})", "", null),
        null, DBusCallFlags.NONE, -1
    );
    var uri = reply.get_child_value (0).get_string ();
    var file = File.new_for_uri (uri);
    return yield file.load_bytes_async (null, null);
}
```

### ai.screen-reader (6 tools — AT-SPI2)

`get_accessibility_tree`, `get_focused_element`, `find_element`, `get_element_hierarchy`, `list_windows`, `get_window_elements`

**Linux uses AT-SPI2** (Assistive Technology Service Provider Interface) via D-Bus:

```vala
public async AccessibilityNode get_accessibility_tree () throws Error {
    var bus = yield Bus.get (BusType.SESSION);
    // AT-SPI2 registry at org.a11y.atspi.Registry
    var reply = yield bus.call (
        "org.a11y.atspi.Registry",
        "/org/a11y/atspi/accessible/root",
        "org.a11y.atspi.Accessible",
        "GetChildren",
        null, null, DBusCallFlags.NONE, -1
    );
    // Recursively build accessibility tree
}
```

### ai.vision (6 tools)

Uses Claude Vision API or OpenAI Vision — platform-independent, same as Swift app.

### ai.browser-context (7 tools)

Communicates via WebSocket to Chrome extension — platform-independent, same as Swift app.

---

## 15. Voice & Notifications

### Voice Plugin (services.voice — 8 tools)

`tts_speak`, `tts_speak_provider`, `tts_list_voices`, `tts_stop`, `stt_listen`, `stt_transcribe_file`, `stt_list_models`, `voice_config`

**Linux TTS:**

| Engine | Library | Quality | Offline |
|--------|---------|---------|---------|
| espeak-ng | CLI / libespeak-ng | Low (robotic) | Yes |
| Festival | CLI | Medium | Yes |
| Piper | CLI | High (neural) | Yes |
| Speech Dispatcher | libspeechd | Varies (meta-engine) | Yes |

**Linux STT:**

| Engine | Library | Quality | Offline |
|--------|---------|---------|---------|
| Vosk | libvosk | High | Yes |
| PocketSphinx | libpocketsphinx | Medium | Yes |
| Whisper.cpp | CLI | Very High | Yes |

**Audio backend**: PipeWire (modern) or PulseAudio (legacy) via GStreamer.

```vala
// TTS via Speech Dispatcher (wraps espeak/festival/piper)
public async void tts_speak (string text) throws Error {
    var conn = new Spd.Connection ("orchestra", null, null, Spd.Mode.THREADED);
    conn.say (Spd.Priority.TEXT, text);
}

// STT via Vosk (offline speech recognition)
public async string stt_listen () throws Error {
    // Open PipeWire/PulseAudio mic stream via GStreamer
    // Feed audio to Vosk model
    // Return transcribed text
}
```

### Notifications Plugin (services.notifications — 8 tools)

`notify_send`, `notify_schedule`, `notify_cancel`, `notify_list_pending`, `notify_badge`, `notify_config`, `notify_history`, `notify_create_channel`

**Linux notifications via D-Bus:**

```vala
public async uint32 notify_send (string title, string body, string? icon = null,
                                  string urgency = "normal") throws Error {
    var bus = yield Bus.get (BusType.SESSION);

    var hints = new HashTable<string, Variant> (str_hash, str_equal);
    switch (urgency) {
        case "low": hints["urgency"] = new Variant.byte (0); break;
        case "normal": hints["urgency"] = new Variant.byte (1); break;
        case "critical": hints["urgency"] = new Variant.byte (2); break;
    }

    // Desktop entry hint for proper icon/grouping
    hints["desktop-entry"] = new Variant.string ("dev.orchestra.desktop");

    var reply = yield bus.call (
        "org.freedesktop.Notifications",
        "/org/freedesktop/Notifications",
        "org.freedesktop.Notifications",
        "Notify",
        new Variant ("(susssasa{sv}i)",
            "Orchestra",           // app_name
            0,                     // replaces_id
            icon ?? "dev.orchestra.desktop",  // icon
            title,                 // summary
            body,                  // body
            new Variant.array (new VariantType ("s"), {}),  // actions
            hints,                 // hints
            5000                   // timeout (ms)
        ),
        null, DBusCallFlags.NONE, -1
    );

    return reply.get_child_value (0).get_uint32 ();
}
```

**GNotification (GNOME-native alternative):**

```vala
// For GNOME environments, GNotification provides richer integration
public void notify_gnome (string title, string body) {
    var notification = new Notification (title);
    notification.set_body (body);
    notification.set_icon (new ThemedIcon ("dev.orchestra.desktop"));
    notification.set_priority (NotificationPriority.NORMAL);

    // Add actions
    notification.add_button (_("View"), "app.show-result");
    notification.add_button (_("Dismiss"), "app.dismiss-notification");

    ((Application) application).send_notification ("orchestra-notify", notification);
}
```

**Notification channels** (pre-configured): `build`, `test`, `deploy`, `ai`, `reminder`, `system`, `git`

---

## 16. Settings & Integrations

### Settings Navigation

Uses `AdwPreferencesWindow` (standard libadwaita settings pattern):

```vala
public class Orchestra.SettingsWindow : Adw.PreferencesWindow {
    construct {
        set_title (_("Settings"));
        set_search_enabled (true);

        add (create_general_page ());
        add (create_appearance_page ());
        add (create_notifications_page ());
        add (create_windows_page ());
        add (create_ai_page ());
        add (create_voice_page ());
        add (create_sync_page ());
    }
}
```

| Section | Settings |
|---------|----------|
| **General** | Timezone, language, autostart |
| **Appearance** | Color theme (25 themes), dark/light/system, accent color |
| **Notifications** | Permission status, channels, DND hours |
| **Windows** | Default window mode, spirit position, bubble anchor |
| **AI** | Default provider + model, auto-approve toggle |
| **Voice** | STT engine (Vosk/Whisper), TTS engine (Piper/espeak), language |
| **Sync & Account** | Sync status, connected devices, API tokens |

### AI Providers (Built-in Multi-LLM)

| Provider | Models | Auth |
|----------|--------|------|
| Anthropic | Claude Opus 4.6, Sonnet 4.6, Haiku 4.5 | Built-in / API key |
| OpenAI | GPT-4o, GPT-4o-mini, o1, o1-mini | API key |
| Google Gemini | Gemini 2.5 Pro, 2.5 Flash, 2.0 Flash | API key |
| Ollama | Llama 3, CodeLlama, Mistral, etc. | Local (no key) |

### GSettings Schema

```xml
<?xml version="1.0" encoding="UTF-8"?>
<schemalist>
  <schema id="dev.orchestra.desktop" path="/dev/orchestra/desktop/">
    <key name="color-theme" type="s">
      <default>'orchestra'</default>
      <summary>Color theme</summary>
    </key>
    <key name="color-scheme" type="s">
      <default>'dark'</default>
      <summary>Color scheme (dark, light, system)</summary>
    </key>
    <key name="window-width" type="i">
      <default>1280</default>
    </key>
    <key name="window-height" type="i">
      <default>860</default>
    </key>
    <key name="window-maximized" type="b">
      <default>false</default>
    </key>
    <key name="default-provider" type="s">
      <default>'claude'</default>
    </key>
    <key name="default-model" type="s">
      <default>'claude-sonnet-4-6'</default>
    </key>
    <key name="window-mode" type="s">
      <default>'main'</default>
      <summary>Window mode (main, floating, bubble)</summary>
    </key>
    <key name="autostart" type="b">
      <default>false</default>
    </key>
    <key name="last-plugin" type="s">
      <default>'chat'</default>
    </key>
  </schema>
</schemalist>
```

---

## 17. Native Linux Features

### Global Hotkey

```vala
// Option 1: xdg-desktop-portal GlobalShortcuts (Wayland-safe, GNOME 45+)
public async void register_global_shortcut () throws Error {
    var bus = yield Bus.get (BusType.SESSION);
    yield bus.call (
        "org.freedesktop.portal.Desktop",
        "/org/freedesktop/portal/desktop",
        "org.freedesktop.portal.GlobalShortcuts",
        "BindShortcuts",
        new Variant ("(sa{sv}a(sa{sv})s)",
            "", null,
            {
                { "cycle-window", {
                    { "description", new Variant.string (_("Cycle Orchestra window mode")) },
                    { "preferred-trigger", new Variant.string ("<Control><Shift>o") }
                }}
            },
            ""
        ),
        null, DBusCallFlags.NONE, -1
    );
}

// Option 2: X11 XGrabKey fallback (for X11 sessions)
// Uses Gdk.X11.Display for X11-specific key grabbing
```

### Credentials (libsecret)

```vala
public class Orchestra.SecretService : Object {
    private const Secret.Schema SCHEMA = {
        "dev.orchestra.desktop",
        Secret.SchemaFlags.NONE,
        {
            { "provider", Secret.SchemaAttributeType.STRING },
        }
    };

    public async void save_api_key (string provider, string key) throws Error {
        yield Secret.password_store (
            SCHEMA, null,
            "Orchestra %s API Key".printf (provider),
            key,
            null,
            "provider", provider
        );
    }

    public async string? load_api_key (string provider) throws Error {
        return yield Secret.password_lookup (
            SCHEMA, null,
            "provider", provider
        );
    }

    public async void delete_api_key (string provider) throws Error {
        yield Secret.password_clear (
            SCHEMA, null,
            "provider", provider
        );
    }
}
```

### System Tray (StatusNotifierItem)

```vala
// StatusNotifierItem via D-Bus (works on KDE, XFCE, Sway/Waybar)
// GNOME requires AppIndicator extension
public class Orchestra.TrayIcon : Object {
    private AppIndicator.Indicator indicator;

    public TrayIcon () {
        indicator = new AppIndicator.Indicator (
            "orchestra",
            "dev.orchestra.desktop",
            AppIndicator.IndicatorCategory.APPLICATION_STATUS
        );
        indicator.set_status (AppIndicator.IndicatorStatus.ACTIVE);

        var menu = new Gtk.Menu ();
        var show_item = new Gtk.MenuItem.with_label (_("Show Orchestra"));
        show_item.activate.connect (() => activate_action ("app.show", null));
        menu.append (show_item);

        var quit_item = new Gtk.MenuItem.with_label (_("Quit"));
        quit_item.activate.connect (() => activate_action ("app.quit", null));
        menu.append (quit_item);

        menu.show_all ();
        indicator.set_menu (menu);
    }
}
```

### XDG Paths

```vala
// Standard Linux paths
namespace Orchestra.Paths {
    public string config_dir () {
        // ~/.config/orchestra/
        return Path.build_filename (Environment.get_user_config_dir (), "orchestra");
    }

    public string data_dir () {
        // ~/.local/share/orchestra/
        return Path.build_filename (Environment.get_user_data_dir (), "orchestra");
    }

    public string cache_dir () {
        // ~/.cache/orchestra/
        return Path.build_filename (Environment.get_user_cache_dir (), "orchestra");
    }

    public string certs_dir () {
        // ~/.orchestra/certs/ (shared with Go plugins)
        return Path.build_filename (Environment.get_home_dir (), ".orchestra", "certs");
    }

    public string runtime_dir () {
        // /run/user/$UID/orchestra/ (session-scoped, tmpfs)
        return Path.build_filename (Environment.get_user_runtime_dir (), "orchestra");
    }
}
```

### Autostart

```vala
// XDG autostart: ~/.config/autostart/dev.orchestra.desktop.desktop
public void enable_autostart () {
    var autostart_dir = Path.build_filename (
        Environment.get_user_config_dir (), "autostart"
    );
    DirUtils.create_with_parents (autostart_dir, 0755);

    var desktop_file = "[Desktop Entry]\nType=Application\nName=Orchestra\nExec=orchestra-desktop --background\nIcon=dev.orchestra.desktop\nX-GNOME-Autostart-enabled=true\n";

    var path = Path.build_filename (autostart_dir, "dev.orchestra.desktop.desktop");
    FileUtils.set_contents (path, desktop_file);
}
```

### xdg-desktop-portal Integration

| Portal | Purpose | Used For |
|--------|---------|----------|
| `org.freedesktop.portal.Screenshot` | Screen capture | ai.screenshot |
| `org.freedesktop.portal.FileChooser` | File open/save | File operations |
| `org.freedesktop.portal.Notification` | Notifications | services.notifications |
| `org.freedesktop.portal.GlobalShortcuts` | Hotkeys | Global shortcut (Ctrl+Shift+O) |
| `org.freedesktop.portal.Background` | Background running | Keep alive when closed |
| `org.freedesktop.portal.Secret` | Credential storage | API key storage |
| `org.freedesktop.portal.OpenURI` | URL opening | External links |

### Auto-Updater

```vala
// Flatpak: automatic via GNOME Software / flatpak update
// deb/rpm: Check GitHub releases API
public class Orchestra.UpdaterService : Object {
    private const string REPO = "orchestra-mcp/orchestra-linux";
    private const int64 CHECK_INTERVAL = 6 * 3600; // 6 hours

    public async UpdateInfo? check_for_update () throws Error {
        var session = new Soup.Session ();
        var uri = "https://api.github.com/repos/%s/releases/latest".printf (REPO);
        var msg = new Soup.Message ("GET", uri);
        var body = yield session.send_and_read_async (msg, Priority.DEFAULT, null);

        var parser = new Json.Parser ();
        parser.load_from_data ((string) body.get_data ());
        var root = parser.get_root ().get_object ();
        var latest = root.get_string_member ("tag_name");

        if (latest != CURRENT_VERSION) {
            return UpdateInfo () {
                version = latest,
                download_url = root.get_string_member ("html_url"),
                release_notes = root.get_string_member ("body")
            };
        }
        return null;
    }
}
```

---

## 18. Packaging & Distribution

### Flatpak (Primary — Universal Linux)

```yaml
# flatpak/dev.orchestra.desktop.yml
app-id: dev.orchestra.desktop
runtime: org.gnome.Platform
runtime-version: '46'
sdk: org.gnome.Sdk
command: orchestra-desktop

finish-args:
  - --share=ipc
  - --share=network
  - --socket=fallback-x11
  - --socket=wayland
  - --socket=pulseaudio
  - --device=dri
  - --talk-name=org.freedesktop.Notifications
  - --talk-name=org.freedesktop.secrets
  - --talk-name=org.a11y.Bus
  - --filesystem=home/.orchestra:ro
  - --filesystem=xdg-run/pipewire-0
  # Portal access (screenshot, file picker, background, global shortcuts)
  - --talk-name=org.freedesktop.portal.Desktop

modules:
  # ngtcp2 (QUIC transport)
  - name: ngtcp2
    buildsystem: cmake-ninja
    config-opts:
      - -DCMAKE_BUILD_TYPE=Release
      - -DENABLE_GNUTLS=ON
    sources:
      - type: archive
        url: https://github.com/ngtcp2/ngtcp2/releases/download/v1.4.0/ngtcp2-1.4.0.tar.xz

  # protobuf-c
  - name: protobuf-c
    buildsystem: autotools
    sources:
      - type: archive
        url: https://github.com/protobuf-c/protobuf-c/releases/download/v1.5.0/protobuf-c-1.5.0.tar.gz

  # libsecret (credential storage)
  - name: libsecret
    buildsystem: meson
    sources:
      - type: archive
        url: https://download.gnome.org/sources/libsecret/0.21/libsecret-0.21.4.tar.xz

  # GtkSourceView 5 (code editor)
  - name: gtksourceview
    buildsystem: meson
    sources:
      - type: archive
        url: https://download.gnome.org/sources/gtksourceview/5.12/gtksourceview-5.12.0.tar.xz

  # VTE (terminal emulator)
  - name: vte
    buildsystem: meson
    config-opts:
      - -Dgtk4=true
      - -Dgtk3=false
    sources:
      - type: archive
        url: https://download.gnome.org/sources/vte/0.76/vte-0.76.0.tar.xz

  # Orchestra Desktop
  - name: orchestra-desktop
    buildsystem: meson
    sources:
      - type: dir
        path: ..
```

### Snap (Ubuntu/Canonical)

```yaml
# snap/snapcraft.yaml
name: orchestra-desktop
base: core24
version: '0.1.0'
summary: Orchestra MCP — AI-agentic IDE for Linux
description: |
  Native GTK4/libadwaita desktop app for Orchestra MCP.
  Manages projects, features, sprints, and AI chat sessions.
confinement: strict
grade: stable

apps:
  orchestra-desktop:
    command: usr/bin/orchestra-desktop
    extensions: [gnome]
    plugs:
      - home
      - network
      - desktop
      - desktop-legacy
      - wayland
      - x11
      - audio-playback
      - audio-record
      - password-manager-service
```

### Debian/Ubuntu (.deb)

```
debian/
├── control
│   Package: orchestra-desktop
│   Version: 0.1.0
│   Architecture: amd64
│   Depends: libgtk-4-1 (>= 4.12), libadwaita-1-0 (>= 1.4),
│            libsecret-1-0, libvte-2.91-gtk4-0, libgtksourceview-5-0
│   Description: Orchestra MCP — AI-agentic IDE for Linux
├── rules
│   %:
│       dh $@ --buildsystem=meson
└── changelog
```

### RPM (Fedora/RHEL)

```spec
Name:           orchestra-desktop
Version:        0.1.0
Release:        1%{?dist}
Summary:        Orchestra MCP — AI-agentic IDE for Linux

License:        MIT
URL:            https://github.com/orchestra-mcp/orchestra-linux
Source0:        %{name}-%{version}.tar.xz

BuildRequires:  meson vala gcc
BuildRequires:  pkgconfig(gtk4) >= 4.12
BuildRequires:  pkgconfig(libadwaita-1) >= 1.4
BuildRequires:  pkgconfig(gtksourceview-5)
BuildRequires:  pkgconfig(vte-2.91-gtk4)
BuildRequires:  pkgconfig(libsecret-1)
BuildRequires:  pkgconfig(ngtcp2)
BuildRequires:  pkgconfig(libprotobuf-c)

Requires:       gtk4 >= 4.12
Requires:       libadwaita >= 1.4
```

### AppImage

```bash
# Build AppImage using linuxdeploy
linuxdeploy \
  --appdir AppDir \
  --executable build/src/orchestra-desktop \
  --desktop-file data/dev.orchestra.desktop.desktop \
  --icon-file data/icons/hicolor/scalable/apps/dev.orchestra.desktop.svg \
  --plugin gtk \
  --output appimage
```

### Distribution Matrix

| Format | Sandbox | Auto-Update | DE Integration | Target |
|--------|---------|-------------|----------------|--------|
| **Flatpak** | Yes (portal) | Yes (Flathub) | Full (portals) | Primary — all distros |
| **Snap** | Yes (AppArmor) | Yes (snapd) | Good (extensions) | Ubuntu |
| **.deb** | No | Manual/PPA | Native | Debian/Ubuntu |
| **.rpm** | No | Manual/COPR | Native | Fedora/RHEL |
| **AppImage** | No | Manual | Limited | Portable |
| **AUR** | No | Yes (AUR helpers) | Native | Arch Linux |

---

## 19. Build Phases

### Phase 1: Shell + Chat (MVP) — GNOME

1. Meson project setup (`apps/linux/meson.build`)
2. `orchestra-kit` Vala library with QUIC transport + plugin system
3. Main window with `AdwNavigationSplitView` plugin-driven sidebar
4. AI chat plugin: session list + multi-LLM conversation + streaming
5. Settings plugin: appearance (theme picker via `AdwPreferencesWindow`)
6. Connection status indicator in header bar
7. Flatpak manifest for development builds

**Exit criteria**: Launch app on GNOME, see plugin-driven sidebar, create chat session with any provider, see streaming response.

### Phase 2: Projects + Notes — GNOME + KDE

1. Projects plugin: list, detail, backlog tree, workflow states
2. Notes plugin: list, GtkSourceView editor, pin/unpin, tags
3. Search spotlight (Ctrl+K via `GtkSearchBar`)
4. Data caching in local SQLite (`~/.local/share/orchestra/cache.db`)
5. KDE Plasma compatibility testing + fixes
6. libsecret credential storage for API keys

### Phase 3: Developer Tools — All DEs

1. Terminal plugin (VTE 4)
2. File Explorer plugin (GtkTreeListModel + GtkSourceView)
3. Service manager (systemd D-Bus integration)
4. Additional DevTools plugins (Database, SSH, Log Viewer)
5. XFCE + Sway/Hyprland compatibility testing

### Phase 4: Window Modes + Native Features

1. Spirit window (floating mini chat, always-on-top)
2. Bubble window (gtk4-layer-shell on Wayland, X11 hints on X11)
3. Global hotkey (xdg-desktop-portal GlobalShortcuts + X11 fallback)
4. System tray (StatusNotifierItem / AppIndicator)
5. XDG autostart
6. Auto-updater (Flatpak: automatic, others: GitHub releases check)

### Phase 5: AI Awareness + Voice

1. Screenshot plugin (xdg-desktop-portal Screenshot + CLI fallbacks)
2. Screen reader plugin (AT-SPI2 D-Bus)
3. Voice STT/TTS plugin (Piper/espeak-ng + Vosk/Whisper.cpp via GStreamer)
4. Notifications plugin (D-Bus org.freedesktop.Notifications + GNotification)
5. Browser context plugin (Chrome extension WebSocket bridge)

### Phase 6: Packaging + Distribution

1. Flatpak on Flathub (primary distribution)
2. Snap on Snapcraft
3. .deb packages (Debian/Ubuntu PPA)
4. .rpm packages (Fedora COPR)
5. AppImage for portable use
6. AUR package for Arch Linux
7. CI/CD: GitHub Actions for all package formats

---

## Appendix A: MCP Tools for Linux App

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

## Appendix B: Swift ↔ Linux Parity Map

| Swift (artifact 18) | Linux (this artifact) | Notes |
|----------------------|----------------------|-------|
| SwiftUI | GTK4 / libadwaita | Both declarative, different paradigms |
| Network.framework QUIC | ngtcp2 (C) | Both native QUIC |
| swift-protobuf | protobuf-c | Both code-generated from same .proto |
| NavigationSplitView | Adw.NavigationSplitView | Nearly identical API/concept |
| MenuBarExtra | StatusNotifierItem/AppIndicator | Tray icon |
| NSPanel (Spirit) | GtkWindow + layer-shell | Floating window |
| ScreenCaptureKit | xdg-desktop-portal Screenshot | Screen capture |
| NSSpeechSynthesizer | espeak-ng / Piper | TTS |
| SFSpeechRecognizer | Vosk / Whisper.cpp | STT |
| Keychain | libsecret (GNOME Keyring / KDE Wallet) | Credential storage |
| UserNotifications | D-Bus org.freedesktop.Notifications | Push notifications |
| AXUIElement | AT-SPI2 (D-Bus) | Accessibility tree |
| GSettings | GSettings | Same library (GLib) |
| SPM (Package.swift) | Meson (meson.build) | Build system |
| .app bundle | Flatpak/deb/rpm/AppImage/snap | Packaging |
| App Store | Flathub / Snapcraft / repos | Distribution |

## Appendix C: Dependencies

| Library | Version | Purpose |
|---------|---------|---------|
| GTK4 | >= 4.12 | UI toolkit |
| libadwaita | >= 1.4 | GNOME adaptive widgets |
| GtkSourceView 5 | >= 5.10 | Code editor with syntax highlighting |
| VTE 4 | >= 0.74 | Terminal emulator |
| libsecret | >= 0.20 | Credential storage |
| ngtcp2 | >= 1.2 | QUIC transport |
| GnuTLS | >= 3.8 | mTLS for QUIC |
| protobuf-c | >= 1.5 | Protobuf serialization |
| json-glib | >= 1.8 | JSON parsing |
| libsoup 3 | >= 3.4 | HTTP client (updates, API calls) |
| sqlite3 | >= 3.40 | Local cache |
| gtk4-layer-shell | >= 1.0 | Wayland overlay windows |
| libappindicator3 | >= 12.10 | System tray (optional) |
| GStreamer 1.0 | >= 1.22 | Audio pipeline (voice) |

## Appendix D: Freedesktop Icon Names

| Plugin | Icon Name | Symbolic |
|--------|-----------|----------|
| Chat | `chat-message-new` | `chat-message-new-symbolic` |
| Projects | `view-grid` | `view-grid-symbolic` |
| Notes | `accessories-text-editor` | `accessories-text-editor-symbolic` |
| DevTools | `utilities-terminal` | `utilities-terminal-symbolic` |
| Settings | `preferences-system` | `preferences-system-symbolic` |
| File Explorer | `system-file-manager` | `system-file-manager-symbolic` |
| SSH | `network-server` | `network-server-symbolic` |
| Docker | `application-x-container` | `application-x-container-symbolic` |
| Database | `accessories-database` | `accessories-database-symbolic` |
| Services | `preferences-system-network` | `preferences-system-network-symbolic` |

# Desktop Subsystems -- Orchestra Reference

> Every desktop-specific subsystem: interface, platform variants, configuration. Extracted from `orch-ref/app/desktop/`.

## Overview

The desktop app (Wails v3) is a macOS/Linux/Windows GUI that embeds a React frontend and orchestrates all backend services. The `app/desktop/` package contains 20+ subsystems across 54 Go source files. Platform-specific behavior is isolated using Go build tags (`darwin`, `linux`, `windows`, `other`).

| # | Subsystem | Key Files | Platform-Specific | Purpose |
|---|-----------|-----------|-------------------|---------|
| 1 | App | `app.go`, `app_menu.go` | No | Main app lifecycle, menu bar |
| 2 | Window Manager | `window_manager.go` | No | Multi-window registry and control |
| 3 | Window API | `window_api.go` | No | MCP/HTTP-triggered window operations |
| 4 | Mode Manager | `mode_manager.go` | No | embedded/floating/bubble mode cycling |
| 5 | Spirit Window | `spirit_window.go`, `spirit_window_*.go` | Yes (darwin/linux/win/other) | Floating mini chat window |
| 6 | Bubble Window | `bubble_window.go`, `bubble_window_*.go` | Yes (darwin/linux/win/other) | Always-on-top circular overlay |
| 7 | Tray Manager | `tray_manager.go`, `tray_service.go`, `tray_actions.go`, `tray_helpers.go`, `tray_types.go` | Partial (observer is darwin) | System tray icon + dynamic menu |
| 8 | Tray Observer | `tray_observer_darwin.go`, `tray_observer_other.go` | Yes (darwin/other) | macOS menu template image observer |
| 9 | Panel Manager | `panel_manager.go`, `panel_service.go`, `panel_types.go`, `panel_helpers.go` | No | Plugin-contributed UI panels |
| 10 | Service Registry | `service_registry.go`, `service_lifecycle.go`, `service_proc_*.go`, `service_tray.go` | Yes (unix/windows proc) | Start/stop/status for all services |
| 11 | Notification | `notification.go` | No (runtime dispatch) | OS-native notifications + sound |
| 12 | TTS Service | `tts_service.go` | No (runtime dispatch) | Text-to-speech via native APIs |
| 13 | Screenshot | `screenshot_darwin.go`, `screenshot_other.go`, `screenshot_window.go`, `screenshot_window_*.go` | Yes (darwin/other) | Screen capture + region overlay |
| 14 | Hotkey | `hotkey_darwin.go`, `hotkey_linux.go`, `hotkey_windows.go`, `hotkey_other.go` | Yes (all platforms) | Global keyboard shortcuts |
| 15 | Permissions | `permissions_darwin.go`, `permissions_other.go` | Yes (darwin/other) | macOS permission prompts + status |
| 16 | Boot Collector | `boot_collector.go`, `boot_collector_binding.go` | No | Gathers plugin UI contributions |
| 17 | Settings Service | `settings_service.go` | No | Persistent key-value settings |
| 18 | Workspace Init | `workspace_init.go` | No | Auto-init MCP on workspace switch |
| 19 | Icon Resolver | `icon_resolver.go` | No (runtime dispatch) | Platform icon path resolution |
| 20 | UI Manifest | `ui_manifest.go` | No | Aggregated plugin UI definitions |
| 21 | Window Level | `window_level_darwin.go`, `window_level_other.go` | Yes (darwin/other) | Always-on-top level helpers |
| 22 | Titlebar Zoom | `titlebar_zoom_darwin.go`, `titlebar_zoom_other.go` | Yes (darwin/other) | macOS double-click zoom handler |
| 23 | Updater | `app/updater/` (separate package) | Yes (darwin restart) | GitHub release auto-updater |

---

## 1. App (`app.go`, `app_menu.go`)

The central orchestrator. Creates the Wails application, initializes all subsystems, and manages the main window lifecycle.

### Struct

```go
type App struct {
    logger             *slog.Logger
    wails              *application.App
    ctx                context.Context
    assets             fs.FS
    trayIcon           []byte
    appIcon            []byte
    windowManager      *WindowManager
    trayService        *TrayService
    trayManager        *TrayManager
    panelService       *PanelService
    panelManager       *PanelManager
    settingsService    *SettingsService
    nativeTTSService   *NativeTTSService
    serviceRegistry    *ServiceRegistry
    modeManager        *ModeManager
    pendingSettings    []pluginSettingsEntry
    windowPinned       bool
    onWorkspaceChange  func(path string)
    recentProjectsMenu *application.Menu
}
```

### Lifecycle

| Method | Signature | Description |
|--------|-----------|-------------|
| `NewApp` | `(logger *slog.Logger, ctx context.Context, assets fs.FS) *App` | Constructor. Creates ServiceRegistry immediately. |
| `Start` | `() error` | Initializes all subsystems in order: SettingsService, NativeTTSService, Wails app, WindowManager, TrayService/TrayManager, PanelService/PanelManager, ModeManager. Installs titlebar zoom, requests permissions, sets up dock reopen handler, builds menu, opens main window, restores window mode. Calls `a.wails.Run()` (blocking). |
| `Shutdown` | `() error` | Stops all services, closes all panels, saves settings, removes hotkey/titlebar zoom/menu observer, destroys tray, closes all windows. |

### Start Initialization Order

1. SettingsService -- load persisted settings from `~/Library/Application Support/Orchestra/settings.json`
2. Register default settings (appearance, general, updates, spirit/bubble bounds, window mode)
3. Flush pending plugin settings
4. NativeTTSService -- create TTS service instance
5. Wails app -- create with SPA asset handler, register SettingsService and NativeTTSService as Wails services
6. WindowManager -- multi-window registry
7. TrayService + TrayManager -- system tray with icon, default items (Open Workspace, Pin Window, Check for Updates, Quit)
8. ServiceRegistry onChange -> TrayManager.Rebuild
9. PanelService + PanelManager -- plugin-contributed panels
10. ModeManager -- restores persisted window mode
11. InstallTitlebarZoom -- macOS double-click zoom
12. RequestPermissions -- macOS notification/screen recording/microphone prompts
13. Dock reopen handler -- show/recreate main window
14. Application menu bar (File, Edit, View, Window, Help)
15. Open main window at `/` (1280x860, hidden titlebar)
16. Restore spirit/bubble window if last mode was floating/bubble

### Key Public Methods

```go
func (a *App) WindowManager() *WindowManager
func (a *App) TrayService() *TrayService
func (a *App) TrayManager() *TrayManager
func (a *App) PanelService() *PanelService
func (a *App) PanelManager() *PanelManager
func (a *App) SettingsService() *SettingsService
func (a *App) ModeManager() *ModeManager
func (a *App) ServiceRegistry() *ServiceRegistry
func (a *App) GetWails() *application.App

func (a *App) GetWindowMode() string
func (a *App) SetWindowMode(mode string)
func (a *App) CycleWindowMode()

func (a *App) CaptureScreen() ([]byte, error)
func (a *App) CaptureScreenRegion(x, y, width, height int) ([]byte, error)

func (a *App) OpenURL(url string) error
func (a *App) OnWorkspaceChange(fn func(path string))

func (a *App) SetAppIcon(iconData []byte)
func (a *App) SetTrayIcon(iconData []byte)

func (a *App) RequestOSPermissions()
func (a *App) GetOSPermissionStatus() settings.PermissionStatus
func (a *App) RequestOSPermission(kind string)
func (a *App) OpenOSPermissionSettings(kind string)

func (a *App) RegisterPluginSettings(pluginID string, p plugins.HasSettings)
```

### Application Menu (`app_menu.go`)

Builds a native macOS/Windows menu bar:

- **File** -- Open Workspace (`Cmd+O`), Recent Projects (dynamic submenu from settings API), Close Window
- **Edit** -- Standard Undo/Redo/Cut/Copy/Paste (Wails built-in)
- **View** -- Reload, Zoom, Fullscreen (Wails built-in)
- **Window** -- Minimize, Zoom, Front (Wails built-in)
- **Help** -- Orchestra Documentation link, Report Issue link

Recent Projects are fetched from `http://127.0.0.1:19191/api/workspace/recent` and include a "Clear Recent Projects" action.

### SPA Handler

The `spaHandler()` function serves embedded assets with React Router fallback -- any URL that does not map to a real file serves `index.html`.

---

## 2. Window Manager (`window_manager.go`)

Thread-safe registry for all Wails webview windows. Every window (main, spirit, bubble, panels, screenshot overlay) is tracked here by name.

### Struct

```go
type WindowManager struct {
    app     *application.App
    logger  *slog.Logger
    windows map[string]*application.WebviewWindow
    mu      sync.RWMutex
}

type WindowOptions struct {
    Name             string
    Title            string
    Width, Height    int
    MinWidth, MinHeight, MaxWidth, MaxHeight int
    AlwaysOnTop      bool
    URL              string
    BackgroundColour application.RGBA
    Frameless        bool
    Hidden           bool
}
```

### Public Methods

```go
func NewWindowManager(app *application.App, logger *slog.Logger) *WindowManager
func (wm *WindowManager) CreateWindow(opts WindowOptions) (*application.WebviewWindow, error)
func (wm *WindowManager) GetWindow(name string) (*application.WebviewWindow, error)
func (wm *WindowManager) CloseWindow(name string) error
func (wm *WindowManager) TrackWindow(name string, win *application.WebviewWindow)
func (wm *WindowManager) RemoveWindow(name string)
func (wm *WindowManager) FocusWindow(name string) error
func (wm *WindowManager) SetWindowSize(name string, width, height int) error
func (wm *WindowManager) SetWindowPosition(name string, x, y int) error
func (wm *WindowManager) SetWindowTitle(name, title string) error
func (wm *WindowManager) HideWindow(name string) error
func (wm *WindowManager) ShowWindow(name string) error
func (wm *WindowManager) MinimizeWindow(name string) error
func (wm *WindowManager) MaximizeWindow(name string) error
func (wm *WindowManager) ListWindows() []string
func (wm *WindowManager) CloseAll()
```

### Window Defaults

- Default size: 1200x800
- macOS: invisible titlebar height 38, translucent backdrop, hidden inset titlebar

---

## 3. Window API (`window_api.go`)

HTTP-triggered window operations. Implements `settings.WindowOpener` so the Go HTTP server can open/close windows when the MCP binary POSTs to `/api/windows/open`.

### Public Methods

```go
func (a *App) OpenWindow(req settings.OpenWindowRequest) error
func (a *App) CloseWindow(name string) error
func (a *App) ShowSaveDialog(filename string) (string, error)
func (a *App) ShowOpenDirectoryDialog(message string) (string, error)
func (a *App) ShowOpenFileDialog(message string, filters []settings.FileFilter) (string, error)
```

### Behavior

- `OpenWindow`: waits up to 10s for WindowManager init, reuses existing window if same name (reload + focus), otherwise creates new always-on-top window. Routes query params from `data` map onto the URL.
- `ShowSaveDialog`: native Save As dialog, defaults to `~/Downloads`
- `ShowOpenDirectoryDialog`: native folder picker
- `ShowOpenFileDialog`: native file picker with optional type filters

### URL Builder

The `buildURL(route string, data map[string]any) string` helper appends scalar values (string, float64, int, bool) as query params. Complex types (slices, maps) are stored server-side and fetched via `/api/windows/data/{name}`.

---

## 4. Mode Manager (`mode_manager.go`)

Manages transitions between three window modes. State is persisted to settings so the last mode is restored on app restart.

### Constants

```go
const (
    ModeEmbedded = "embedded"  // Full main window
    ModeFloating = "floating"  // Spirit mini window
    ModeBubble   = "bubble"    // Always-on-top bubble overlay
    modeSettingsKey = "spirit.window.mode"
)
```

### Struct

```go
type ModeManager struct {
    mu   sync.Mutex
    mode string
    app  *App
}
```

### Public Methods

```go
func NewModeManager(app *App) *ModeManager
func (m *ModeManager) CurrentMode() string
func (m *ModeManager) SetMode(mode string)
func (m *ModeManager) CycleMode()  // embedded -> floating -> bubble -> embedded
```

### Behavior

- `SetMode`: closes the old window (CloseSpiritWindow / CloseBubbleWindow), opens the new one (OpenSpiritWindow / OpenBubbleWindow), persists to settings.
- `CycleMode`: embedded -> floating -> bubble -> embedded.
- On startup, `restore()` reads the persisted mode from settings and sets `m.mode` (no window is opened -- that happens in `App.Start`).

---

## 5. Spirit Window (`spirit_window.go` + platform variants)

A floating mini chat window. Frameless, always-on-top, semi-transparent background, resizable with min constraints. Position and size are persisted to settings.

### Constants

```go
const (
    spiritWindowName  = "spirit"
    spiritWidth       = 420
    spiritHeight      = 640
    spiritMinWidth    = 320
    spiritMinHeight   = 400
    spiritSettingsKey = "spirit.window.bounds"
)
```

### Public Methods (on App)

```go
func (a *App) OpenSpiritWindow()   // Create or focus
func (a *App) CloseSpiritWindow()  // Save bounds and close
func (a *App) ToggleSpiritWindow() // Open if closed, close if open
```

### Behavior

- Loads at `/panels/spirit` route
- Frameless with RGBA background `{10, 10, 20, 230}`
- Persists bounds (x, y, width, height) as JSON in settings on move, resize, and close

### Platform Variants

`spiritWindowMacOptions()` returns platform-specific `application.MacWindow`:

| Platform | Build Tag | Behavior |
|----------|-----------|----------|
| **macOS** | `darwin` | Translucent backdrop, hidden inset titlebar (28px drag region), `CanJoinAllSpaces + FullScreenAuxiliary + IgnoresCycle` collection behavior |
| **Linux** | `linux` | Empty MacWindow. Wails handles AlwaysOnTop via `_NET_WM_STATE_ABOVE` (X11) or layer-shell (Wayland) |
| **Windows** | `windows` | Empty MacWindow. Wails handles AlwaysOnTop via `WS_EX_TOPMOST`, Frameless via `WS_POPUP` |
| **Other** | `!darwin && !windows && !linux` | Empty MacWindow (fallback) |

---

## 6. Bubble Window (`bubble_window.go` + platform variants)

A 56x56 pixel always-on-top circular overlay. Fixed size (no resize), transparent background. Position only (no size) is persisted.

### Constants

```go
const (
    bubbleWindowName  = "bubble"
    bubbleSize        = 56
    bubbleSettingsKey = "bubble.window.bounds"
)
```

### Public Methods (on App)

```go
func (a *App) OpenBubbleWindow()   // Create or focus
func (a *App) CloseBubbleWindow()  // Save position and close
func (a *App) ToggleBubbleWindow() // Open if closed, close if open
```

### Behavior

- Loads at `/panels/bubble` route
- Fixed 56x56 (min=max=56), frameless, non-resizable
- RGBA background `{0, 0, 0, 1}` (essentially transparent)
- Persists position (x, y) as JSON in settings on move and close

### Platform Variants

`bubbleWindowMacOptions()`:

| Platform | Behavior |
|----------|----------|
| **macOS** | Translucent backdrop, fully hidden titlebar, `CanJoinAllSpaces + FullScreenAuxiliary + IgnoresCycle` |
| **Linux/Windows/Other** | Empty MacWindow |

---

## 7. Tray Manager (`tray_manager.go`, `tray_service.go`, `tray_actions.go`, `tray_helpers.go`, `tray_types.go`)

A multi-layer system tray implementation with service/manager/type separation and a pluggable action handler.

### Types (`tray_types.go`)

```go
type ActionType string  // "ipc" | "url" | "command" | "plugin"

type TrayAction struct {
    Type    ActionType
    Payload map[string]interface{}
}

type TrayItem struct {
    ID, Label, Icon, Color, Group, Hotkey, Tooltip, PluginID string
    Action      *TrayAction
    Order       int
    Visible, Enabled, Checked, IsSeparator, IsCheckbox bool
}

type TrayGroup struct { Name, Label string; Order int; Items []*TrayItem }

type TrayConfig struct {
    Icon, Title, DefaultGroup string
    Groups       []*TrayGroup
    MaxItems     int
    SortByOrder, ShowIcons, ShowTooltips, EnableHotkeys bool
}
```

### TrayService (`tray_service.go`)

Thread-safe registry of tray items and groups. Default config: max 50 items, sorted by order, icons/tooltips/hotkeys enabled.

```go
func NewTrayService(logger *slog.Logger) *TrayService
func (s *TrayService) RegisterItem(item *TrayItem) error
func (s *TrayService) UnregisterItem(id string) error
func (s *TrayService) GetItem(id string) (*TrayItem, error)
func (s *TrayService) ListItems() []*TrayItem
func (s *TrayService) ListItemsByGroup(group string) []*TrayItem
func (s *TrayService) RegisterGroup(group *TrayGroup) error
func (s *TrayService) GetGroup(name string) (*TrayGroup, error)
func (s *TrayService) UpdateItem(item *TrayItem) error
func (s *TrayService) SetItemVisibility(id string, visible bool) error
func (s *TrayService) SetItemEnabled(id string, enabled bool) error
func (s *TrayService) GetConfig() *TrayConfig
func (s *TrayService) UpdateConfig(config *TrayConfig)
func (s *TrayService) Count() int
func (s *TrayService) Clear()
```

### TrayManager (`tray_manager.go`)

Bridges TrayService with Wails SystemTray. Creates the native tray, rebuilds the menu from TrayService state.

```go
func NewTrayManager(app *application.App, service *TrayService, desktopApp *App, logger *slog.Logger) *TrayManager
func (tm *TrayManager) Initialize() error
func (tm *TrayManager) InitializeWithIcon(iconData []byte) error
func (tm *TrayManager) Rebuild() error          // Reconstructs entire menu from TrayService
func (tm *TrayManager) UpdateItem(id string) error
func (tm *TrayManager) SetIcon(iconData []byte)
func (tm *TrayManager) SetLabel(label string)   // macOS only
func (tm *TrayManager) Destroy()
func (tm *TrayManager) SetServiceRegistry(r *ServiceRegistry)
```

On macOS, the tray icon is set as a template icon (`SetTemplateIcon`) for automatic light/dark mode adaptation.

During `Rebuild()`, the "Services" submenu is injected before the quit separator via `BuildServicesSubmenu()`.

### TrayHelpers (`tray_helpers.go`)

Builder-pattern helpers for creating tray items and actions:

```go
func NewTrayItem(id, label string) *TrayItem
func NewSeparator(id string) *TrayItem
func (t *TrayItem) WithIcon/WithColor/WithGroup/WithHotkey/WithAction/WithOrder/WithTooltip/WithPluginID/WithCheckbox(...)

func NewIPCAction(event string, data map[string]interface{}) *TrayAction
func NewURLAction(url string) *TrayAction
func NewCommandAction(command string, args []string) *TrayAction
func NewPluginAction(pluginID, method string, params map[string]interface{}) *TrayAction
```

### ActionHandler (`tray_actions.go`)

Dispatches tray click actions by type:

| ActionType | Handler | Behavior |
|------------|---------|----------|
| `ipc` | `handleIPC` | Routes to built-in events or emits custom Wails events |
| `url` | `handleURL` | Opens URL in default browser |
| `command` | `handleCommand` | Executes shell command |
| `plugin` | `handlePlugin` | Emits `plugin:invoke` Wails event |

Built-in IPC events: `settings:open`, `chat:open`, `workspace:open`, `spirit:toggle`, `bubble:toggle`, `window:pin`, `updater:check`, `app:quit`.

### Default Tray Items (registered in `App.registerDefaultTrayItems`)

| ID | Label | Action |
|----|-------|--------|
| `open-workspace` | Open Workspace... | IPC `workspace:open` |
| `pin-window` | Pin Window (checkbox) | IPC `window:pin` |
| `updater` | Check for Updates... | IPC `updater:check` |
| `sep-1` | (separator) | -- |
| `quit` | Quit Orchestra MCP | IPC `app:quit` |

### Service Tray (`service_tray.go`)

Builds a "Services (N/M)" submenu inside the tray from the ServiceRegistry:

```go
func BuildServicesSubmenu(menu *application.Menu, registry *ServiceRegistry, onAction func()) *application.Menu
```

Each service gets a label with status icon (`●` running, `○` stopped/starting, `✖` error) and a submenu with Start/Stop/Restart actions plus port info. Bulk "Start All" and "Stop All" actions are appended.

---

## 8. Tray Observer (`tray_observer_darwin.go`, `tray_observer_other.go`)

macOS-only NSMenu observer that marks menu item images as template images for proper dark mode support.

### macOS Implementation (CGo + Objective-C)

Creates an `OrchestraMenuObserver` that listens for `NSMenuDidBeginTrackingNotification`. When a menu opens, it iterates all items and submenus, setting `[image setTemplate:YES]` and sizing to 16x16.

```go
func InstallMenuTemplateObserver(logger *slog.Logger)  // darwin: installs NSNotificationCenter observer
func RemoveMenuTemplateObserver(logger *slog.Logger)   // darwin: removes observer
```

### Non-macOS

Both functions are no-ops.

**Note**: Currently disabled in `App.Start()` for freeze debugging.

---

## 9. Panel Manager (`panel_manager.go`, `panel_service.go`, `panel_types.go`, `panel_helpers.go`)

Plugin-contributed UI panels that can be displayed as separate windows, modals, embedded views, or sidebars.

### Panel Types

```go
type PanelType string     // "standard" | "modal" | "window" | "sidebar"
type PanelPosition string // "center" | "left" | "right" | "top" | "bottom"

type Panel struct {
    ID, Title, Route    string
    Type                PanelType
    Position            PanelPosition
    Icon, PluginID      string
    Width, Height, MinWidth, MinHeight int
    Resizable, Closable, Singleton, Visible bool
    WindowID            string
    Props               PanelProps
    Permissions         []string
}

type PanelProps struct {
    ShowHeader, ShowFooter, Scrollable, BackdropBlur, CloseOnEscape, CloseOnOutside bool
    Padding   int
    Custom    map[string]interface{}
}

type PanelConfig struct {
    MaxPanels          int           // default: 50
    DefaultWidth       int           // default: 800
    DefaultHeight      int           // default: 600
    DefaultType        PanelType     // default: "standard"
    DefaultPosition    PanelPosition // default: "center"
    EnableTransitions  bool          // default: true
    TransitionDuration int           // default: 300ms
    PersistState       bool          // default: true
    AllowMultipleModal bool          // default: false
}

type PanelEventType string // "opened" | "closed" | "focused" | "blurred" | "resized" | "moved"
type PanelEvent struct { Type PanelEventType; PanelID, WindowID string; Data map[string]interface{} }
```

### PanelService (`panel_service.go`)

Thread-safe panel registry with event channel (buffered 100). Enforces singleton constraint and modal exclusivity.

```go
func NewPanelService(logger *slog.Logger) *PanelService
func (s *PanelService) RegisterPanel(panel *Panel) error
func (s *PanelService) UnregisterPanel(id string) error
func (s *PanelService) GetPanel(id string) (*Panel, error)
func (s *PanelService) ListPanels() []*Panel
func (s *PanelService) ListPanelsByPlugin(pluginID string) []*Panel
func (s *PanelService) ListPanelsByType(panelType PanelType) []*Panel
func (s *PanelService) GetOpenPanels() []*Panel
func (s *PanelService) GetPanelByRoute(route string) (*Panel, error)
func (s *PanelService) OpenPanel(id, windowID string) error
func (s *PanelService) ClosePanel(id string) error
func (s *PanelService) UpdatePanel(panel *Panel) error
func (s *PanelService) RegisterGroup(group *PanelGroup) error
func (s *PanelService) GetGroup(name string) (*PanelGroup, error)
func (s *PanelService) GetConfig() *PanelConfig
func (s *PanelService) UpdateConfig(config *PanelConfig)
func (s *PanelService) Count() int
func (s *PanelService) CountOpen() int
func (s *PanelService) Clear()
func (s *PanelService) Events() <-chan *PanelEvent
func (s *PanelService) Close()
```

### PanelManager (`panel_manager.go`)

Integrates PanelService with WindowManager. Handles opening panels in the correct mode:

```go
func NewPanelManager(app *application.App, service *PanelService, windowManager *WindowManager, logger *slog.Logger) *PanelManager
func (pm *PanelManager) OpenPanel(panelID string) error
func (pm *PanelManager) ClosePanel(panelID string) error
func (pm *PanelManager) TogglePanel(panelID string) error
func (pm *PanelManager) GetOpenPanels() []*Panel
func (pm *PanelManager) CloseAllPanels() error
func (pm *PanelManager) OpenPanelByRoute(route string) error
func (pm *PanelManager) ResizePanel(panelID string, width, height int) error
func (pm *PanelManager) SetPanelPosition(panelID string, x, y int) error
```

Panel opening behavior by type:

| PanelType | Behavior |
|-----------|----------|
| `window` | Creates a new Wails window named `panel-{id}` via WindowManager |
| `modal` | Opens in the main window, emits navigation event |
| `standard` / `sidebar` | Opens embedded in the main window |

### Panel Helpers (`panel_helpers.go`)

Builder-pattern constructors:

```go
func NewPanel(id, title, route string) *Panel              // standard, 800x600
func NewModalPanel(id, title, route string) *Panel          // modal, backdrop blur, close on escape/outside
func NewWindowPanel(id, title, route string) *Panel          // separate window, 1000x700
func NewSidebarPanel(id, title, route string, position PanelPosition) *Panel  // sidebar, 300px wide

func (p *Panel) WithType/WithPosition/WithIcon/WithSize/WithMinSize/WithResizable/WithClosable/WithPluginID/WithSingleton/WithProps/WithPermissions(...)
```

---

## 10. Service Registry (`service_registry.go`, `service_lifecycle.go`, `service_proc_unix.go`, `service_proc_windows.go`, `service_tray.go`)

Tracks all running services (Go HTTP, WebSocket, gRPC, MCP) with lifecycle management. Services can be started/stopped via callbacks or external commands.

### Types

```go
type ServiceStatus string // "stopped" | "starting" | "running" | "stopping" | "error"
type ServiceType string   // "http" | "websocket" | "grpc" | "mcp"

type ServiceEntry struct {
    ID        string
    Label     string
    Type      ServiceType
    Host      string
    Port      int
    Status    ServiceStatus
    Error     string
    PID       int
    StartCmd  string
    StartArgs []string
    startFn   func(ctx context.Context) error  // internal callback
    stopFn    func() error                     // internal callback
}

type ServiceRegistry struct {
    logger   *slog.Logger
    mu       sync.RWMutex
    services map[string]*ServiceEntry
    onChange  func()
}
```

### ServiceRegistry Methods

```go
func NewServiceRegistry(logger *slog.Logger) *ServiceRegistry
func (r *ServiceRegistry) OnChange(fn func())
func (r *ServiceRegistry) Register(entry *ServiceEntry)
func (r *ServiceRegistry) SetCallbacks(id string, startFn func(ctx context.Context) error, stopFn func() error)
func (r *ServiceRegistry) SetStatus(id string, status ServiceStatus, errMsg string)
func (r *ServiceRegistry) SetPID(id string, pid int)
func (r *ServiceRegistry) Get(id string) (*ServiceEntry, bool)
func (r *ServiceRegistry) List() []*ServiceEntry    // sorted by type then label
func (r *ServiceRegistry) RunningCount() int
func (e *ServiceEntry) CanControl() bool
func (e *ServiceEntry) Addr() string
```

### Lifecycle (`service_lifecycle.go`)

```go
func (r *ServiceRegistry) Start(id string) error    // Uses startFn or StartCmd
func (r *ServiceRegistry) Stop(id string) error      // Uses stopFn or kills PID
func (r *ServiceRegistry) Restart(id string) error   // Stop, wait 500ms, Start
func (r *ServiceRegistry) StopAll()
func IsPortOpen(host string, port int) bool
```

Start flow:
1. If `startFn` is set, call it with 30s timeout context
2. If `StartCmd` is set, exec the command in a new process group
3. If port is configured, poll with `waitForPort()` (10s timeout, 500ms intervals)

### Platform Process Management

| Platform | `setSvcProcAttr` | `killSvcProcess` |
|----------|-----------------|------------------|
| **Unix** (darwin/linux) | Sets `Setpgid: true` for process group | Sends `SIGTERM` to process group, then `SIGKILL` after 2s |
| **Windows** | No-op | `proc.Kill()` |

---

## 11. Notification Service (`notification.go`)

Cross-platform desktop notifications and sound playback. All errors are best-effort (logged but never returned).

### Struct

```go
type NotificationService struct {
    logger    *slog.Logger
    soundsDir string
}
```

### Public Methods

```go
func NewNotificationService(logger *slog.Logger, soundsDir string) *NotificationService
func (s *NotificationService) Send(req settings.NotifyRequest) error  // implements settings.Notifier
func (s *NotificationService) PlaySound(name string) bool
```

### Notification Delivery by Platform

| Platform | Method | Mechanism |
|----------|--------|-----------|
| **macOS** | `sendOsascript` | AppleScript `display notification` (works in unsigned dev builds) |
| **Linux** | `sendNotifySend` | `notify-send --app-name="Orchestra MCP"` (libnotify) |
| **Windows** | `sendPowerShellToast` | PowerShell `Windows.UI.Notifications.ToastNotification` |

### Sound Playback by Platform

| Platform | Command |
|----------|---------|
| **macOS** | `afplay -v 0.80 {path}` |
| **Linux** | `paplay` (PulseAudio) or `aplay` (ALSA) |
| **Windows** | PowerShell `System.Media.SoundPlayer` |

Sound files are `.mp3` in the configured `soundsDir`. The child process is reaped asynchronously.

---

## 12. TTS Service (`tts_service.go`)

Native text-to-speech exposed to the Wails frontend as a bound service. Supports macOS `say` and Windows SAPI. Linux has no native TTS support (frontend can use Web Speech API).

### Types

```go
type TTSBackend string  // "native" | "webspeech"

type NativeTTSVoice struct {
    Name   string
    Lang   string
    Sample string
}

type NativeTTSService struct {
    mu      sync.Mutex
    current *exec.Cmd
}
```

### Public Methods (Wails Bindings)

```go
func NewNativeTTSService() *NativeTTSService
func (s *NativeTTSService) ListVoices() []NativeTTSVoice
func (s *NativeTTSService) Speak(text, voiceName string) error
func (s *NativeTTSService) Stop()
func (s *NativeTTSService) StatusJSON() map[string]interface{}
```

### Platform Implementation

| Platform | List Voices | Speak |
|----------|-------------|-------|
| **macOS** | Parses `say -v ?` output, filters novelty voices (Bad News, Bahh, Bells, etc.), deduplicates Siri variants | `say [-v voiceName] text` |
| **Windows** | PowerShell `System.Speech.Synthesis.SpeechSynthesizer.GetInstalledVoices()` | PowerShell `SpeechSynthesizer.Speak()` |
| **Linux** | Returns `nil` (no native TTS) | No-op |

The `Stop()` method kills the current speech process immediately. Only one speech can be active at a time.

---

## 13. Screenshot (`screenshot_darwin.go`, `screenshot_other.go`, `screenshot_window.go`, `screenshot_window_darwin.go`, `screenshot_window_other.go`)

Screen capture using macOS ScreenCaptureKit via CGo. Includes a fullscreen transparent overlay window for region selection.

### Capture Functions

```go
func captureScreen() ([]byte, error)                                    // Full screen PNG
func captureScreenRegion(x, y, width, height int) ([]byte, error)      // Cropped region PNG
```

### macOS Implementation (CGo + Objective-C)

Uses `SCScreenshotManager` from ScreenCaptureKit framework:
1. `SCShareableContent.getShareableContentWithCompletionHandler` to get display list
2. `SCContentFilter` for the main display
3. `SCStreamConfiguration` matching display dimensions, cursor hidden
4. `SCScreenshotManager.captureImageWithFilter` to get `CGImageRef`
5. `CGImageDestinationCreateWithData` to encode as PNG

Region capture: captures full screen, decodes PNG, crops with `image.SubImage`, re-encodes.

**Requires**: Screen Recording permission in System Settings > Privacy & Security.

### Non-macOS

Both functions return `errors.New("screenshot capture not supported on this platform")`.

### Screenshot Overlay Window (`screenshot_window.go`)

A fullscreen transparent window for interactive region selection:

```go
func (a *App) OpenScreenshotOverlay()  // Creates "screenshot" window at /panels/screenshot, goes fullscreen
func (a *App) CloseScreenshotOverlay()
```

- 1920x1080 initial size, then `.Fullscreen()`
- Frameless, non-resizable, always-on-top
- RGBA `{0, 0, 0, 1}` (transparent)
- macOS: translucent backdrop, hidden titlebar, CanJoinAllSpaces

---

## 14. Hotkey (`hotkey_darwin.go`, `hotkey_linux.go`, `hotkey_windows.go`, `hotkey_other.go`)

Global keyboard shortcut registration. Currently only implemented on macOS.

### Functions

```go
func InstallGlobalHotkey(logger *slog.Logger, callback func())
func RemoveGlobalHotkey(logger *slog.Logger)
```

### macOS Implementation (CGo + Objective-C)

Registers `Cmd+Shift+O` as a global hotkey using `NSEvent.addGlobalMonitorForEventsMatchingMask`. When triggered, calls the Go callback via `//export goHotkeyCallback`.

The callback is typically `ModeManager.CycleMode()` to cycle through embedded/floating/bubble modes.

### Linux / Windows

Stub implementations that log "not yet implemented". TODOs reference:
- Linux: `XGrabKey` (X11) or dbus shortcut (Wayland)
- Windows: `RegisterHotKey` Win32 API

### Other

No-op functions.

**Note**: Currently disabled in `App.Start()` for freeze debugging.

---

## 15. Permissions (`permissions_darwin.go`, `permissions_other.go`)

macOS system permission management for 4 capabilities: notifications, screen recording, microphone, and accessibility.

### PermissionStatus Struct

```go
type PermissionStatus struct {
    Notifications   string  // "not_determined" | "granted" | "denied" | "provisional"
    ScreenRecording string  // "granted" | "denied"
    Microphone      string  // "not_determined" | "restricted" | "denied" | "granted"
    Accessibility   string  // "granted" | "denied"
}
```

### Functions

```go
func GetPermissionStatus() PermissionStatus
func RequestPermissions(logger *slog.Logger)        // Request all at once (startup)
func RequestOSPermission(kind string)               // Request specific: "screen_recording" | "notifications" | "microphone" | "accessibility"
func OpenSystemSettings(kind string)                // Open System Settings to relevant pane
```

### macOS Implementation (CGo + Objective-C)

| Permission | Check API | Request API | Settings URL |
|------------|-----------|-------------|--------------|
| Notifications | `UNUserNotificationCenter.getNotificationSettings` | `requestAuthorizationWithOptions` (alert+sound+badge) | `x-apple.systempreferences:com.apple.preference.notifications` |
| Screen Recording | `CGPreflightScreenCaptureAccess()` | `CGRequestScreenCaptureAccess()` | `com.apple.preference.security?Privacy_ScreenCapture` |
| Microphone | `AVCaptureDevice.authorizationStatusForMediaType(AVMediaTypeAudio)` | `AVCaptureDevice.requestAccessForMediaType` | `com.apple.preference.security?Privacy_Microphone` |
| Accessibility | `AXIsProcessTrusted()` | `AXIsProcessTrustedWithOptions(prompt:YES)` | `com.apple.preference.security?Privacy_Accessibility` |

Notification permission is skipped when not running as a `.app` bundle (dev mode detection via `NSBundle.mainBundle.bundleIdentifier`).

### Non-macOS

All functions are no-ops. `GetPermissionStatus()` returns `"granted"` for everything.

---

## 16. Boot Collector (`boot_collector.go`, `boot_collector_binding.go`)

Gathers UI contributions from all loaded plugins at startup. Iterates plugin capability interfaces and builds a `UIManifest`.

### BootCollector

```go
type BootCollector struct {
    logger   *slog.Logger
    manifest *UIManifest
}

func NewBootCollector(logger *slog.Logger) *BootCollector
func (c *BootCollector) Collect(pluginMap map[string]*plugins.PluginInfo) *UIManifest
```

Checks each plugin for these capability interfaces:
- `HasPanels` -> `GetPanels()` -> `manifest.Panels`
- `HasTrayItems` -> `GetTrayItems()` -> `manifest.TrayItems`
- `HasContextMenus` -> `GetContextMenus()` -> `manifest.ContextMenus`
- `HasIPCHandlers` -> `GetIPCHandlers()` -> `manifest.IPCHandlers`
- `HasSidebarViews` -> `GetSidebarViews()` -> `manifest.SidebarViews`
- `HasTabs` -> `GetTabs()` -> `manifest.Tabs`

All collections are sorted by `Order` field after collection.

### UIManifestService (Wails Binding)

```go
type UIManifestService struct {
    logger   *slog.Logger
    manifest *UIManifest
    mu       sync.RWMutex
}

func NewUIManifestService(logger *slog.Logger) *UIManifestService
func (s *UIManifestService) CollectFromPlugins(pluginMap map[string]*plugins.PluginInfo)
func (s *UIManifestService) GetUIManifest() *UIManifest        // Wails binding
func (s *UIManifestService) GetPanels() []plugins.PanelDef
func (s *UIManifestService) GetSidebarViews() []plugins.SidebarViewDef
func (s *UIManifestService) GetTabs() []plugins.TabDef
func (s *UIManifestService) Refresh(pluginMap map[string]*plugins.PluginInfo)
```

The frontend calls `GetUIManifest()` once at boot to populate its panel registry, tray state, context menus, IPC channels, sidebar views, and tabs.

---

## 17. Settings Service (`settings_service.go`)

Persistent key-value settings store. Persists to a JSON file at `~/.config/Orchestra/settings.json` (or platform equivalent). Exposed to the frontend as a Wails binding.

### Types

```go
type SettingsService struct {
    logger   *slog.Logger
    settings map[string]*SettingEntry
    groups   map[string]*SettingGroup
    filePath string
    onChange func(key string, value any)
    mu       sync.RWMutex
}

type SettingGroup struct { ID, Label, Description, Icon string; Order int }

type SettingEntry struct {
    Key, Label, Type, Group, PluginID string
    Default, Value any
    Order int
}

type SettingsState struct {
    Groups   []SettingGroup
    Settings []SettingEntry
}
```

### Public Methods

```go
func NewSettingsService(logger *slog.Logger, configDir string) *SettingsService
func (s *SettingsService) RegisterGroup(g SettingGroup)
func (s *SettingsService) RegisterSetting(e SettingEntry) error
func (s *SettingsService) Get(key string) (any, error)
func (s *SettingsService) Set(key string, value any) error   // Persists to disk, fires onChange
func (s *SettingsService) SetOnChange(fn func(key string, value any))
func (s *SettingsService) GetState() SettingsState           // Full snapshot for frontend
func (s *SettingsService) Save() error                       // Persist to JSON file
func (s *SettingsService) Load() error                       // Read from JSON file
func (s *SettingsService) SaveFileDialog(filename, content, mimeType string) (string, error)
func (s *SettingsService) SaveBinaryFile(filename string, data string) (string, error)  // base64 input
```

### Default Settings Groups

| Group ID | Label | Icon |
|----------|-------|------|
| `appearance` | Appearance | palette |
| `general` | General | settings |
| `updates` | Updates | download |

### Default Settings

| Key | Type | Default | Group |
|-----|------|---------|-------|
| `appearance.colorTheme` | string | `"orchestra"` | appearance |
| `appearance.componentVariant` | string | `"default"` | appearance |
| `updates.autoCheck` | boolean | `true` | updates |
| `spirit.window.bounds` | string | `""` | general |
| `bubble.window.bounds` | string | `""` | general |
| `spirit.window.mode` | string | `"embedded"` | general |

### File Save Helpers

`SaveFileDialog` and `SaveBinaryFile` write to `~/Downloads` with automatic counter-based deduplication if the filename already exists.

---

## 18. Workspace Init (`workspace_init.go`)

Auto-initializes Orchestra MCP when switching to a new workspace.

### Functions

```go
func needsOrchestraInit(workspacePath string) bool
func (a *App) switchWorkspace(path string)
```

### Behavior

1. `needsOrchestraInit` checks for `.projects/` directory or `.mcp.json` file
2. `switchWorkspace` runs in a background goroutine:
   - Auto-initializes MCP if needed via `mcp.RunInit(path)`
   - POSTs to `http://127.0.0.1:19191/api/workspace/open` to persist workspace
   - Emits `workspace:changed` Wails event for frontend
   - Calls `onWorkspaceChange` callback for WebSocket broadcast
   - Uses `application.InvokeAsync` for UI operations (focus/menu rebuild) to avoid AppKit crashes from background goroutines

---

## 19. Icon Resolver (`icon_resolver.go`)

Resolves boxicon names to platform-specific file paths.

### Struct

```go
type IconResolver struct {
    baseDir string  // {resourceDir}/os-icons/
}
```

### Public Methods

```go
func NewIconResolver(resourceDir string) *IconResolver
func (r *IconResolver) Resolve(iconName string) (string, error)
func (r *IconResolver) ResolveWithFallback(iconName, fallback string) string
```

### Platform Icon Paths

| Platform | Path Pattern |
|----------|-------------|
| **macOS** | `os-icons/darwin/{name}.png` (@2x auto-discovered by OS) |
| **Windows** | `os-icons/windows/{name}-32.png` |
| **Linux** | `os-icons/linux/{name}.png` |

---

## 20. UI Manifest (`ui_manifest.go`)

Aggregated plugin UI definitions. Populated by BootCollector, consumed by the frontend at boot.

```go
type UIManifest struct {
    Panels       []plugins.PanelDef
    TrayItems    []plugins.TrayItemDef
    ContextMenus []plugins.ContextMenuDef
    IPCHandlers  []plugins.IPCHandlerDef
    SidebarViews []plugins.SidebarViewDef
    Tabs         []plugins.TabDef
}

func NewUIManifest() *UIManifest  // initializes all slices to empty (not nil)
```

---

## 21. Window Level (`window_level_darwin.go`, `window_level_other.go`)

Helpers to set window z-order level. Thin wrappers around `win.SetAlwaysOnTop()`.

```go
func SetWindowModalLevel(win *application.WebviewWindow)   // Sets always-on-top true
func SetWindowNormalLevel(win *application.WebviewWindow)   // Sets always-on-top false
```

Both platforms use the same Wails API. The darwin-specific file exists for potential future NSWindow level manipulation via CGo.

---

## 22. Titlebar Zoom (`titlebar_zoom_darwin.go`, `titlebar_zoom_other.go`)

macOS-specific handler that detects double-clicks in the invisible titlebar region and triggers the standard window zoom (maximize).

### macOS Implementation (CGo + Objective-C)

Installs an `NSEvent` local monitor for `NSEventMaskLeftMouseDown`. On double-click:
1. Checks if the window has `NSWindowStyleMaskFullSizeContentView` (hidden titlebar)
2. Checks if click Y position is in the top 52px (titlebar zone)
3. Calls `[window zoom:nil]` and consumes the event

```go
func InstallTitlebarZoom(logger *slog.Logger)   // Installs the event monitor
func RemoveTitlebarZoom(logger *slog.Logger)    // Removes the event monitor
```

### Non-macOS

Both functions are no-ops.

---

## 23. Updater (`app/updater/`)

Separate package (`package updater`) that checks GitHub releases for updates using `go-selfupdate`.

### Types

```go
type Info struct {
    Version, ReleaseURL, ReleaseNotes, AssetURL, PublishedAt string
}

type Updater struct {
    logger      *slog.Logger
    repo        string          // "orchestra-mcp/framework"
    current     string          // current semver
    interval    time.Duration   // default: 6 hours
    onAvailable func(Info)
    checking    bool
    cancel      context.CancelFunc
}
```

### Public Methods

```go
func New(logger *slog.Logger, repo, currentVersion string) *Updater
func (u *Updater) OnAvailable(fn func(Info))
func (u *Updater) Start(ctx context.Context)    // First check after 5s, then every 6h
func (u *Updater) Stop()
func (u *Updater) CheckNow(ctx context.Context) (*Info, error)
func (u *Updater) Apply(ctx context.Context) error   // Download + replace binary
func (u *Updater) Restart() error                    // Relaunch after update
```

### Update Flow

1. `CheckNow`: Uses `go-selfupdate` GitHubSource to detect latest release. Compares with current version using semver. Returns `nil` if up-to-date.
2. `Apply`: Downloads the release asset, validates checksum (`checksums.txt`), replaces the current executable.
3. `Restart`:
   - **macOS**: If inside `.app` bundle (`Contents/MacOS/binary`), uses `open -n {appPath}` to relaunch
   - **Other**: Re-exec the binary with same args

### Periodic Check Loop

`Start()` spawns a goroutine that checks after 5 seconds, then every 6 hours. Skips concurrent checks (mutex-guarded `checking` flag).

---

## Platform CGo Summary

The following subsystems use CGo with Objective-C on macOS:

| Subsystem | Frameworks | Purpose |
|-----------|------------|---------|
| Hotkey | Cocoa | `NSEvent.addGlobalMonitorForEventsMatchingMask` |
| Permissions | Foundation, UserNotifications, CoreGraphics, AppKit, AVFoundation, ApplicationServices | Permission APIs (UNUserNotificationCenter, CGPreflightScreenCaptureAccess, AVCaptureDevice, AXIsProcessTrusted) |
| Screenshot | ScreenCaptureKit, CoreGraphics, ImageIO, UniformTypeIdentifiers, Foundation | SCScreenshotManager screen capture |
| Tray Observer | Cocoa | NSNotificationCenter menu observer |
| Titlebar Zoom | Cocoa | NSEvent local monitor for double-click |

All CGo features have no-op stubs on non-macOS platforms via build tags.

---

## Wiring: How Subsystems Connect

```
App.Start()
  |
  +-- SettingsService        (settings.json persistence)
  +-- NativeTTSService       (Wails service binding)
  +-- Wails App              (SPA asset handler, service bindings)
  +-- WindowManager          (tracks all windows by name)
  +-- TrayService            (item/group registry)
  +-- TrayManager            (bridges TrayService -> Wails SystemTray)
  |     +-- ActionHandler    (IPC/URL/command/plugin dispatch)
  |     +-- ServiceRegistry  (onChange -> TrayManager.Rebuild)
  +-- PanelService           (panel/group registry, event channel)
  +-- PanelManager           (bridges PanelService -> WindowManager)
  +-- ModeManager            (embedded/floating/bubble state machine)
  +-- InstallTitlebarZoom    (macOS event monitor)
  +-- RequestPermissions     (macOS one-time prompts)
  +-- openMainWindow         (1280x860 at /, hidden titlebar)
  +-- Restore mode           (spirit/bubble if last mode != embedded)

App.Shutdown()
  |
  +-- ServiceRegistry.StopAll()
  +-- PanelManager.CloseAllPanels()
  +-- PanelService.Close()
  +-- SettingsService.Save()
  +-- RemoveGlobalHotkey()
  +-- RemoveTitlebarZoom()
  +-- RemoveMenuTemplateObserver()
  +-- TrayManager.Destroy()
  +-- WindowManager.CloseAll()
```

### Plugin Integration Points

- **Tray Items**: Plugins implement `HasTrayItems` -> BootCollector collects -> TrayService registers -> TrayManager rebuilds
- **Panels**: Plugins implement `HasPanels` -> BootCollector collects -> PanelService registers -> PanelManager opens
- **Settings**: Plugins implement `HasSettings` -> `App.RegisterPluginSettings()` -> SettingsService registers groups and entries
- **Services**: Backend services register via `ServiceRegistry.Register()` + `SetCallbacks()` -> appear in tray Services submenu
- **Context Menus / IPC / Sidebar / Tabs**: Collected by BootCollector into UIManifest -> frontend reads via Wails binding

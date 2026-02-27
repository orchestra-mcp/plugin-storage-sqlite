# Plugin System -- Orchestra Reference

> The complete plugin architecture: lifecycle, manifest, capabilities, IPC, events, dependency resolution, DI container, and discovery. Extracted from `orch-ref/app/plugins/`.

---

## Overview

The plugin system is the foundational architecture of Orchestra. Every feature -- MCP tools, notifications, Discord integration, DevTools sessions -- is delivered as a plugin. The system provides:

- A **4-phase lifecycle** (Load, Boot, Register, Shutdown) orchestrated by the `LifecycleManager`
- A **manifest schema** (`manifest.json`) with validation, permissions, and contribution declarations
- **8 capability interfaces** (`Has*`) that plugins implement to contribute UI elements, settings, IPC handlers, and more
- A **dependency resolver** using Kahn's algorithm (topological sort) with cycle detection
- **IPC communication** over Unix domain sockets with JSON-RPC-style messages
- An **event bus** with wildcard pattern matching for publish/subscribe
- A **DI container** with transient and singleton service lifetimes

Source files: 23 Go files (15 implementation + 3 test files + 5 capability files) in `orch-ref/app/plugins/`.

---

## Plugin Interface

The core contract every plugin must implement:

```go
// Plugin represents a plugin that can be loaded by the framework
type Plugin interface {
    // Load is called when the plugin is first discovered
    // Plugins should parse their manifest and prepare for initialization
    Load(ctx *Context) error

    // Boot is called after all dependencies are loaded
    // Plugins should register services and subscribe to events
    Boot(ctx *Context) error

    // Register is called after all plugins are booted
    // Plugins should register UI contributions (commands, menus, panels)
    Register(ctx *Context) error

    // Shutdown is called when the application is shutting down
    // Plugins should cleanup resources and close connections
    Shutdown(ctx *Context) error
}
```

**Phase purposes:**

| Phase | When Called | Purpose |
|-------|-----------|---------|
| `Load` | During discovery, before dependencies are resolved | Parse manifest, validate configuration, prepare internal state |
| `Boot` | After all dependencies have been loaded | Register services in the DI container, subscribe to events |
| `Register` | After all plugins are booted | Register UI contributions (commands, menus, panels, themes) |
| `Shutdown` | During application teardown | Close connections, release resources, persist state |

---

## Plugin Context

Every lifecycle method receives a `Context` providing framework access:

```go
type Context struct {
    // Manifest is the plugin's parsed manifest.json
    Manifest *Manifest

    // Logger is a structured logger scoped to this plugin
    Logger *slog.Logger

    // Context for cancellation and timeouts
    Ctx context.Context

    // PluginDir is the absolute path to the plugin directory
    PluginDir string

    // TODO: Future additions
    // EventBus   EventBus
    // Container  Container
    // Config     map[string]interface{}
}
```

The logger is automatically scoped with the plugin ID via `slog.With("plugin", id)`, so all log entries from a plugin include its identity.

---

## Plugin States

Plugins move through a strict state machine:

```
StateUnloaded --> StateLoaded --> StateBooted --> StateRegistered --> StateShutdown
                      |               |                |
                      v               v                v
                  StateError      StateError       StateError
```

```go
type PluginState int

const (
    StateUnloaded    PluginState = iota  // Not loaded yet
    StateLoaded                          // Load() completed successfully
    StateBooted                          // Boot() completed successfully
    StateRegistered                      // Register() completed successfully
    StateShutdown                        // Shutdown() was called
    StateError                           // Plugin encountered an error
)
```

State transitions are enforced by the `LifecycleManager`:
- `Load` only runs on plugins in `StateUnloaded`
- `Boot` only runs on plugins in `StateLoaded` (skips `StateError`)
- `RegisterAll` only runs on plugins in `StateBooted` (skips `StateError`)
- `Shutdown` runs on any active state (`StateLoaded`, `StateBooted`, or `StateRegistered`)

---

## Plugin Info

Runtime metadata tracked for each loaded plugin:

```go
type PluginInfo struct {
    Instance Plugin      // The plugin implementation
    Manifest *Manifest   // The parsed manifest
    State    PluginState // Current lifecycle state
    Error    error       // Last error if State is StateError
    Dir      string      // Absolute path to the plugin directory
}
```

---

## Manifest Schema

Every plugin requires a `manifest.json` file. The full schema:

```go
type Manifest struct {
    ID           string             `json:"id"`                    // Required. Lowercase, numbers, hyphens only
    Name         string             `json:"name"`                  // Required. Display name
    Version      string             `json:"version"`               // Required. Semantic versioning (e.g., 1.0.0)
    Description  string             `json:"description,omitempty"` // Optional
    Author       string             `json:"author,omitempty"`      // Optional
    License      string             `json:"license,omitempty"`     // Optional
    Homepage     string             `json:"homepage,omitempty"`    // Optional
    Repository   string             `json:"repository,omitempty"` // Optional
    Main         string             `json:"main"`                  // Required. Entry point binary/file
    Dependencies []string           `json:"dependencies,omitempty"` // Plugin IDs this plugin depends on
    Services     map[string]Service `json:"services,omitempty"`    // Services provided by this plugin
    Contributes  *Contributes       `json:"contributes,omitempty"` // UI/functionality contributions
    Permissions  []string           `json:"permissions,omitempty"` // Required permissions
    Engines      map[string]string  `json:"engines,omitempty"`     // Engine version constraints
}
```

### Validation Rules

1. **ID**: Required. Must match `^[a-z0-9-]+$` (lowercase letters, numbers, hyphens only).
2. **Name**: Required. Any string.
3. **Version**: Required. Must match semantic versioning: `^\d+\.\d+\.\d+(-[a-zA-Z0-9.-]+)?(\+[a-zA-Z0-9.-]+)?$`
4. **Main**: Required. Path to the plugin entry point.
5. **Permissions**: Each must be one of the valid permission strings (see below).
6. **Panel locations**: Must be `sidebar`, `panel`, or `modal`.
7. **Theme types**: Must be `dark` or `light`.

### Valid Permissions

| Permission | Description |
|-----------|-------------|
| `filesystem:read` | Read files from disk |
| `filesystem:write` | Write files to disk |
| `network:http` | Make HTTP requests |
| `network:websocket` | Open WebSocket connections |
| `clipboard:read` | Read from clipboard |
| `clipboard:write` | Write to clipboard |
| `shell:execute` | Execute shell commands |
| `database:read` | Read from database |
| `database:write` | Write to database |
| `ai:api` | Access AI/LLM APIs |
| `system:notifications` | Send system notifications |

### Service Definition

```go
type Service struct {
    Interface string `json:"interface"` // Interface name this service implements
    Singleton bool   `json:"singleton,omitempty"` // Whether to share a single instance
}
```

### Contributes Definition

VS Code-style contributions for UI and functionality:

```go
type Contributes struct {
    Commands []Command `json:"commands,omitempty"`
    Menus    []Menu    `json:"menus,omitempty"`
    Panels   []Panel   `json:"panels,omitempty"`
    Themes   []Theme   `json:"themes,omitempty"`
}
```

**Command:**
```go
type Command struct {
    ID         string `json:"id"`
    Title      string `json:"title"`
    Category   string `json:"category,omitempty"`
    Icon       string `json:"icon,omitempty"`
    Keybinding string `json:"keybinding,omitempty"`
}
```

**Menu:**
```go
type Menu struct {
    ID      string `json:"id"`
    Label   string `json:"label"`
    Command string `json:"command,omitempty"` // Command ID to execute
    Group   string `json:"group,omitempty"`
    When    string `json:"when,omitempty"`    // Conditional expression
}
```

**Panel:**
```go
type Panel struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    Component string `json:"component"`
    Icon      string `json:"icon,omitempty"`
    Location  string `json:"location,omitempty"` // "sidebar" | "panel" | "modal"
}
```

**Theme:**
```go
type Theme struct {
    ID    string `json:"id"`
    Label string `json:"label"`
    Path  string `json:"path"`
    Type  string `json:"type,omitempty"` // "dark" | "light"
}
```

### Example manifest.json

```json
{
  "id": "orchestra-discord",
  "name": "Discord Notifications",
  "version": "1.0.0",
  "description": "Send workflow notifications to Discord",
  "main": "discord",
  "dependencies": ["orchestra-notifications"],
  "permissions": ["network:http", "network:websocket"],
  "services": {
    "discord-notifier": {
      "interface": "NotificationChannel",
      "singleton": true
    }
  },
  "contributes": {
    "commands": [
      {
        "id": "discord.send-test",
        "title": "Send Test Notification",
        "category": "Discord"
      }
    ],
    "panels": [
      {
        "id": "discord-settings",
        "title": "Discord Settings",
        "component": "DiscordSettingsPanel",
        "location": "modal"
      }
    ]
  }
}
```

---

## Capability Interfaces

Plugins implement `Has*` interfaces to declare capabilities. The boot collector calls the corresponding `Get*` method during startup. There are **8 capability interfaces**:

### HasSettings

Contributes settings groups and individual settings to the application settings panel.

```go
type HasSettings interface {
    GetSettingGroups() []SettingGroupDef
    GetSettings() []SettingDef
}
```

**SettingGroupDef** -- defines a collapsible group in the settings UI:
```go
type SettingGroupDef struct {
    ID          string
    Label       string
    Description string
    Icon        string
    Order       int
    Collapsible bool
}
```

**SettingDef** -- defines a single setting with type, default, and validation:
```go
type SettingDef struct {
    Key         string
    Label       string
    Description string
    Placeholder string
    Type        string // string, number, boolean, select, multi-select, color, range
    Default     any
    Group       string
    Order       int
    Options     []SettingOptionDef
    Validation  *SettingValidationDef
}
```

**SettingValidationDef** -- validation constraints:
```go
type SettingValidationDef struct {
    Required  bool
    Min       float64
    Max       float64
    Step      float64
    Pattern   string
    MinLength int
    MaxLength int
}
```

**Real-world example** (notifications plugin):
```go
func (p *NotificationsPlugin) GetSettings() []plugins.SettingDef {
    return []plugins.SettingDef{
        {Key: "notifications.enabled", Label: "Enable Notifications", Type: "boolean", Default: true},
        {Key: "notifications.sound", Label: "Notification Sound", Type: "boolean", Default: true},
        {Key: "notifications.position", Label: "Position", Type: "select", Default: "top-right",
            Options: []plugins.SettingOptionDef{
                {Label: "Top Right", Value: "top-right"},
                {Label: "Bottom Right", Value: "bottom-right"},
            }},
    }
}
```

### HasIPCHandlers

Contributes IPC handlers for inter-process communication between app and plugins.

```go
type HasIPCHandlers interface {
    GetIPCHandlers() []IPCHandlerDef
}
```

```go
type IPCHandlerDef struct {
    Channel     string     `json:"channel"`               // e.g., "mcp:call", "discord:send"
    Description string     `json:"description,omitempty"`
    Pattern     IPCPattern `json:"pattern"`                // request | stream | bistream
    Permissions []string   `json:"permissions,omitempty"`
}
```

**IPC Patterns:**
```go
const (
    IPCPatternRequest  IPCPattern = "request"  // Request-response
    IPCPatternStream   IPCPattern = "stream"   // Server-sent events
    IPCPatternBiStream IPCPattern = "bistream" // Bidirectional stream
)
```

### HasSidebarViews

Contributes views to the sidebar (Chrome extension side panel, desktop sidebar).

```go
type HasSidebarViews interface {
    GetSidebarViews() []SidebarViewDef
}
```

```go
type SidebarViewDef struct {
    ID        string             `json:"id"`
    Title     string             `json:"title"`
    Icon      string             `json:"icon"`
    Route     string             `json:"route"`
    Order     int                `json:"order"`
    Visible   bool               `json:"visible"`
    Badge     string             `json:"badge,omitempty"`
    Actions   []SidebarActionDef `json:"actions,omitempty"`
    HasSearch bool               `json:"hasSearch"`
}

type SidebarActionDef struct {
    ID      string `json:"id"`
    Icon    string `json:"icon"`
    Tooltip string `json:"tooltip"`
    Action  string `json:"action"`
}
```

### HasTrayItems

Contributes items to the system tray menu (desktop only).

```go
type HasTrayItems interface {
    GetTrayItems() []TrayItemDef
}
```

```go
type TrayItemDef struct {
    ID          string         `json:"id"`
    Label       string         `json:"label"`
    Icon        string         `json:"icon"`
    Color       string         `json:"color,omitempty"`
    Group       string         `json:"group"`
    Hotkey      string         `json:"hotkey,omitempty"`
    ActionType  TrayActionType `json:"actionType"`  // ipc | link | panel | callback
    Action      string         `json:"action"`
    Order       int            `json:"order"`
    Visible     bool           `json:"visible"`
    Enabled     bool           `json:"enabled"`
    Tooltip     string         `json:"tooltip,omitempty"`
    IsSeparator bool           `json:"isSeparator"`
    Children    []TrayItemDef  `json:"children,omitempty"` // Nested submenu
}
```

**Tray Action Types:**
```go
const (
    TrayActionIPC      TrayActionType = "ipc"      // Send an IPC message
    TrayActionLink     TrayActionType = "link"      // Open a URL
    TrayActionPanel    TrayActionType = "panel"     // Open a panel window
    TrayActionCallback TrayActionType = "callback"  // Call a Go function
)
```

### HasTabs

Contributes tab definitions to the tab bar (Chrome extension, desktop).

```go
type HasTabs interface {
    GetTabs() []TabDef
}
```

```go
type TabDef struct {
    ID       string `json:"id"`
    Title    string `json:"title"`
    Icon     string `json:"icon"`
    Route    string `json:"route"`
    Closable bool   `json:"closable"`
    Order    int    `json:"order"`
    Pinned   bool   `json:"pinned"`
    Visible  bool   `json:"visible"`
}
```

### HasPanels

Contributes panel windows (separate window views).

```go
type HasPanels interface {
    GetPanels() []PanelDef
}
```

```go
type PanelDef struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Route       string    `json:"route"`
    Icon        string    `json:"icon"`
    Type        PanelType `json:"type"`       // default | frameless | modal
    Width       int       `json:"width"`
    Height      int       `json:"height"`
    MinWidth    int       `json:"minWidth"`
    MinHeight   int       `json:"minHeight"`
    Resizable   bool      `json:"resizable"`
    Closable    bool      `json:"closable"`
    Singleton   bool      `json:"singleton"`  // Only one instance allowed
    Visible     bool      `json:"visible"`
    LazyLoad    bool      `json:"lazyLoad"`
    Order       int       `json:"order"`
    Permissions []string  `json:"permissions,omitempty"`
}
```

**Panel Types:**
```go
const (
    PanelTypeDefault   PanelType = "default"   // Standard window chrome
    PanelTypeFrameless PanelType = "frameless"  // No window chrome
    PanelTypeModal     PanelType = "modal"      // Modal dialog
)
```

### HasDevToolsSessions

Contributes dev tool session types (terminal, database, SSH, etc.).

```go
type HasDevToolsSessions interface {
    GetSessionProviders() []SessionProviderDef
}
```

```go
type SessionProviderDef struct {
    Type        string `json:"type"`        // e.g., "terminal", "database", "ssh"
    Name        string `json:"name"`        // Display name
    Icon        string `json:"icon"`        // Boxicon name
    Description string `json:"description"` // Short description
    Order       int    `json:"order"`       // Sort order in picker
}
```

### HasContextMenus

Contributes context menu items to various targets in the UI.

```go
type HasContextMenus interface {
    GetContextMenus() []ContextMenuDef
}
```

```go
type ContextMenuDef struct {
    ID       string            `json:"id"`
    Label    string            `json:"label"`
    Icon     string            `json:"icon,omitempty"`
    Hotkey   string            `json:"hotkey,omitempty"`
    Target   ContextMenuTarget `json:"target"`   // Where the menu appears
    Group    string            `json:"group,omitempty"`
    Action   string            `json:"action"`
    Order    int               `json:"order"`
    Visible  bool              `json:"visible"`
    Enabled  bool              `json:"enabled"`
    Children []ContextMenuDef  `json:"children,omitempty"` // Nested submenu
}
```

**Context Menu Targets:**
```go
const (
    ContextMenuEditor  ContextMenuTarget = "editor"  // Code editor area
    ContextMenuSidebar ContextMenuTarget = "sidebar" // Sidebar items
    ContextMenuTab     ContextMenuTarget = "tab"     // Tab bar items
    ContextMenuPanel   ContextMenuTarget = "panel"   // Panel windows
    ContextMenuTray    ContextMenuTarget = "tray"    // System tray
    ContextMenuGlobal  ContextMenuTarget = "global"  // Available everywhere
)
```

---

## Lifecycle Manager

The `LifecycleManager` orchestrates plugin phases and tracks state:

```go
type LifecycleManager struct {
    plugins map[string]*PluginInfo
    logger  *slog.Logger
}
```

### Methods

| Method | Purpose |
|--------|---------|
| `Register(id, instance, manifest, dir)` | Add a plugin to the manager (must not already exist) |
| `Load(ctx)` | Run `Load()` on all `StateUnloaded` plugins |
| `Boot(ctx)` | Run `Boot()` on all `StateLoaded` plugins, skip `StateError` |
| `RegisterAll(ctx)` | Run `Register()` on all `StateBooted` plugins, skip `StateError` |
| `Shutdown(ctx)` | Run `Shutdown()` on active plugins in **reverse order** |
| `GetPlugin(id)` | Return `PluginInfo` by ID |
| `ListPlugins()` | Return all registered plugins (defensive copy) |
| `GetPluginsByState(state)` | Return plugins in a specific state |

### Error Handling

- If any plugin fails during `Load`, `Boot`, or `RegisterAll`, its state is set to `StateError` and the error is returned immediately (fail-fast).
- During `Shutdown`, errors are logged but execution continues to shut down remaining plugins. The last error is returned.
- Failed plugins (`StateError`) are silently skipped in subsequent phases.

### Shutdown Order

Shutdown runs in **reverse registration order** -- plugins registered last are shut down first. This ensures that plugins depending on other plugins are cleaned up before their dependencies.

---

## Dependency Resolution

The `DependencyResolver` determines safe load order using **Kahn's algorithm** for topological sorting:

```go
type DependencyResolver struct {
    plugins map[string]*DiscoveredPlugin
}
```

### Algorithm (Kahn's Topological Sort)

1. Build an in-degree map: for each plugin, count how many dependencies it has.
2. Seed a queue with all plugins that have zero in-degree (no dependencies).
3. Dequeue plugins one at a time, appending to the sorted result.
4. For each dequeued plugin, decrement the in-degree of all plugins that depend on it.
5. If any plugin reaches zero in-degree, add it to the queue.
6. After the queue is empty, if `len(sorted) != len(plugins)` then a **circular dependency** exists.

### Methods

| Method | Purpose |
|--------|---------|
| `AddPlugin(plugin)` | Add a discovered plugin (validates non-nil manifest and ID) |
| `Resolve()` | Topological sort -- returns ordered list or error for missing deps/cycles |
| `ValidateDependencies()` | Check that all declared dependencies exist in the resolver |
| `HasCircularDependency()` | Returns `true` if `Resolve()` detects a cycle |

### Error Cases

- **Missing dependency**: `"missing dependency: <dep-id>"` -- a plugin declares a dependency that was not discovered.
- **Circular dependency**: `"circular dependency detected"` -- plugins form a dependency cycle.

---

## Plugin Discovery

The `DiscoveryScanner` scans the filesystem for plugins:

```go
type DiscoveryScanner struct {
    pluginDir string
}

type DiscoveredPlugin struct {
    Manifest *Manifest
    Path     string // Absolute path to the plugin directory
}
```

### How Discovery Works

1. Verify the `pluginDir` exists.
2. Recursively walk the directory tree using `filepath.Walk`.
3. For each file named `manifest.json`, parse and validate it with `ParseManifest`.
4. If the manifest is valid, create a `DiscoveredPlugin` with the manifest and the parent directory path.
5. Invalid manifests are silently skipped (discovery continues).

### Methods

| Method | Purpose |
|--------|---------|
| `Discover()` | Scan and return all valid plugins |
| `DiscoverOne(pluginID)` | Find a specific plugin by ID |
| `Count()` | Return the number of discovered plugins |

---

## IPC Protocol

Inter-process communication uses JSON messages over Unix domain sockets with newline delimiters.

### Message Format

```go
type IPCMessage struct {
    Type    string          `json:"type"`    // "call", "response", "event", "error"
    ID      string          `json:"id"`      // Message ID for request/response matching
    Method  string          `json:"method"`  // Method name for calls
    Params  json.RawMessage `json:"params"`  // Parameters (JSON encoded)
    Result  json.RawMessage `json:"result"`  // Result for responses
    Error   *IPCError       `json:"error"`   // Error for error responses
    Event   string          `json:"event"`   // Event name for events
    Payload json.RawMessage `json:"payload"` // Event payload
}

type IPCError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    string `json:"data,omitempty"`
}
```

### Message Types

| Type | Fields Used | Purpose |
|------|------------|---------|
| `call` | `ID`, `Method`, `Params` | Request from app to plugin (or vice versa) |
| `response` | `ID`, `Result` | Success response matched by `ID` |
| `error` | `ID`, `Error` | Error response matched by `ID` |
| `event` | `Event`, `Payload` | Fire-and-forget notification |

### Error Codes (JSON-RPC-compatible)

| Code | Constant | Meaning |
|------|----------|---------|
| -32700 | `ErrParseError` | Invalid JSON |
| -32600 | `ErrInvalidRequest` | Invalid request object |
| -32601 | `ErrMethodNotFound` | Method not found |
| -32602 | `ErrInvalidParams` | Invalid method parameters |
| -32603 | `ErrInternalError` | Internal error |
| -32000 | `ErrTimeout` | Request timeout |
| -32001 | `ErrPluginCrashed` | Plugin process crashed |

### Message Constructors

```go
func NewCallMessage(id, method string, params interface{}) (*IPCMessage, error)
func NewResponseMessage(id string, result interface{}) (*IPCMessage, error)
func NewErrorMessage(id string, code int, message string, data string) *IPCMessage
func NewEventMessage(event string, payload interface{}) (*IPCMessage, error)
```

### Convenience Accessors

```go
func (m *IPCMessage) GetParams(v interface{}) error   // Unmarshal params into struct
func (m *IPCMessage) GetResult(v interface{}) error   // Unmarshal result into struct
func (m *IPCMessage) GetPayload(v interface{}) error  // Unmarshal event payload into struct
func (m *IPCMessage) Marshal() ([]byte, error)        // Serialize to JSON
func (m *IPCMessage) Unmarshal(data []byte) error     // Deserialize from JSON
```

---

## IPC Bridge

The `IPCBridge` manages communication with a single plugin process over a Unix domain socket:

```go
type IPCBridge struct {
    pluginID   string
    socketPath string          // /tmp/orchestra-<pluginID>.sock
    process    *exec.Cmd
    conn       net.Conn
    handlers   map[string]IPCHandler
    pending    map[string]chan *IPCMessage  // Request/response matching
    // ... mutex, context, cancel, started flag
}

type IPCHandler func(params json.RawMessage) (interface{}, error)
```

### Connection Lifecycle

1. **Start**: Creates a Unix domain socket at `/tmp/orchestra-<pluginID>.sock`.
2. Spawns the plugin process with `ORCHESTRA_IPC_SOCKET` environment variable pointing to the socket.
3. Waits up to **10 seconds** for the plugin to connect.
4. Starts a message handler goroutine reading newline-delimited JSON.

### Stop (Graceful Shutdown)

1. Cancels the context.
2. Closes the socket connection.
3. Sends `SIGINT` to the plugin process.
4. Waits up to **5 seconds** for graceful exit.
5. Force-kills with `SIGKILL` if the process does not exit.
6. Removes the socket file.

### Request/Response Flow

1. Generate a unique message ID (nanosecond timestamp).
2. Create a response channel and register it in `pending[msgID]`.
3. Send the call message over the socket.
4. Wait for a response on the channel with a configurable timeout.
5. On response, remove from `pending` and return the result.
6. On timeout, clean up and return a timeout error.

### Message Handling Loop

The bridge reads newline-delimited JSON from the socket and dispatches based on `Type`:

- **call**: Look up the handler by `Method`, invoke it, and send back a `response` or `error`.
- **response/error**: Match by `ID` to a pending request channel.
- **event**: Log the event (extensible for custom event handlers).

### Buffer Size

The scanner uses a **1 MB buffer** (`1024 * 1024` bytes) for reading messages.

---

## IPC Manager

The `IPCManager` manages all IPC bridges for all plugin processes:

```go
type IPCManager struct {
    bridges map[string]*IPCBridge
    logger  *slog.Logger
    mu      sync.RWMutex
    ctx     context.Context
    cancel  context.CancelFunc
}
```

### Methods

| Method | Purpose |
|--------|---------|
| `StartPlugin(pluginID, pluginPath, args)` | Start a plugin process and establish IPC |
| `StopPlugin(pluginID)` | Stop a specific plugin process |
| `StopAll()` | Stop all plugin processes |
| `Call(pluginID, method, params, timeout)` | Call a method on a specific plugin |
| `SendEvent(pluginID, event, payload)` | Send an event to a specific plugin |
| `Broadcast(event, payload)` | Send an event to **all** plugins |
| `GetStats()` | Return running status for all plugins |

### Standard Handlers

Every bridge gets three standard handlers registered automatically:

| Handler | Response |
|---------|----------|
| `ping` | `{"pong": "ok"}` |
| `getInfo` | `{"pluginID": "<id>", "status": "running"}` |
| `shutdown` | `{"status": "shutting down"}` (triggers async stop after 100ms) |

### Plugin Stats

```go
type PluginStats struct {
    PluginID string
    Running  bool
    PID      int
}
```

---

## Event Bus

An in-process publish/subscribe system with wildcard pattern matching:

### Interface

```go
type EventBus interface {
    Publish(event Event) error
    Subscribe(eventType string, handler EventHandler) error
    Unsubscribe(eventType string, handler EventHandler) error
    UnsubscribeAll(subscriber interface{}) error
}

type Event struct {
    Type      string      // e.g., "app.ready", "plugin.loaded", "file.created"
    Payload   interface{} // Event data (any type)
    Source    string       // Plugin ID that published the event
    Timestamp time.Time   // When the event was published
}

type EventHandler func(event Event) error
```

### Implementation: InProcessEventBus

Goroutine-safe implementation using `sync.RWMutex`:

```go
type InProcessEventBus struct {
    subscribers map[string][]EventHandler
    mu          sync.RWMutex
    logger      *slog.Logger
}
```

### Pattern Matching

The event bus supports three matching modes:

| Pattern | Example | Matches |
|---------|---------|---------|
| **Exact** | `"file.created"` | Only `"file.created"` |
| **Wildcard prefix** | `"file.*"` | `"file.created"`, `"file.deleted"`, `"file.updated"` |
| **Global wildcard** | `"*"` | Every event |

### Additional Methods

| Method | Purpose |
|--------|---------|
| `SubscriberCount(eventType)` | Count handlers matching an event type (includes wildcard matches) |
| `Clear()` | Remove all subscriptions |

### Error Handling

- If a handler returns an error, it is logged but **does not stop other handlers** from being called.
- The last handler error is returned from `Publish()`.

---

## DI Container

A goroutine-safe dependency injection container with two service lifetimes:

### Interface

```go
type Container interface {
    Register(name string, factory ServiceFactory) error       // Transient
    RegisterSingleton(name string, instance interface{}) error // Singleton
    Resolve(name string) (interface{}, error)
    Has(name string) bool
    Remove(name string) error
    Clear()
}

type ServiceFactory func(c Container) (interface{}, error)
```

### Service Lifetimes

```go
type ServiceLifetime int

const (
    Transient ServiceLifetime = iota // New instance each time Resolve is called
    Singleton                        // Same instance returned every time
)
```

### Implementation: SimpleContainer

```go
type SimpleContainer struct {
    services  map[string]*ServiceRegistration
    resolving map[string]bool // Circular dependency detection
    mu        sync.RWMutex
}
```

### Service Registration

```go
type ServiceRegistration struct {
    Name     string
    Lifetime ServiceLifetime
    Factory  ServiceFactory    // For transient services
    Instance interface{}       // For singleton services
}
```

### Circular Dependency Detection

When resolving a transient service, the container tracks services currently being resolved in the `resolving` map. If a factory tries to resolve a service that is already being resolved, the container returns: `"circular dependency detected: <name>"`.

### Thread Safety

All operations are protected by `sync.RWMutex`:
- `Resolve`, `Has`: Read lock (`RLock`)
- `Register`, `RegisterSingleton`, `Remove`, `Clear`: Write lock (`Lock`)
- Transient factory invocation: Briefly locks for circular dependency tracking, then releases

---

## Bootstrap Sequence

The `App.Bootstrap()` method in `bootstrap/app.go` executes the full plugin initialization:

```go
type App struct {
    Logger    *slog.Logger
    Lifecycle *plugins.LifecycleManager
    PluginDir string
}
```

### Full Bootstrap Flow

```
App.Bootstrap(ctx)
  |
  +--> Step 1: discoverPlugins()
  |      Uses DiscoveryScanner to scan PluginDir
  |      Walks filesystem looking for manifest.json files
  |      Returns []*DiscoveredPlugin
  |
  +--> Step 2: resolveDependencies(discovered)
  |      Creates DependencyResolver
  |      Adds all discovered plugins
  |      Validates all dependencies exist
  |      Checks for circular dependencies
  |      Runs topological sort (Kahn's algorithm)
  |      Returns plugins in safe load order
  |
  +--> Step 3: registerPlugins(ordered)
  |      For each plugin in dependency order:
  |        Load the plugin implementation
  |        Call Lifecycle.Register(id, instance, manifest, path)
  |
  +--> Step 4: initializePlugins(ctx)
         |
         +--> Phase 1: Lifecycle.Load(ctx)
         |      Each plugin: Load(ctx) --> StateLoaded
         |
         +--> Phase 2: Lifecycle.Boot(ctx)
         |      Each plugin: Boot(ctx) --> StateBooted
         |
         +--> Phase 3: Lifecycle.RegisterAll(ctx)
                Each plugin: Register(ctx) --> StateRegistered
```

### Desktop Bootstrap Extension

The `DesktopApp` wraps `App` and adds:

1. Call `App.Bootstrap(ctx)` to load all plugins.
2. Initialize settings store (SQLite-backed HTTP server on port 19191).
3. Wire AI bridge, GitHub/Notion/Jira/Linear/Figma hubs.
4. Start MCP bridge for project management tools.
5. Initialize DevTools session manager.
6. Start auto-updater.
7. Auto-start services (WebSocket on 8765, Rust gRPC engine on 50051).
8. Start Wails desktop application on the main goroutine.

### Shutdown Flow

```
App.Shutdown(ctx)
  |
  +--> Lifecycle.Shutdown(ctx)
         Iterates plugins in REVERSE registration order
         Each plugin: Shutdown(ctx) --> StateShutdown
         Errors are logged but do not stop other shutdowns
         Returns the last error (if any)
```

---

## Real-World Plugin Examples

### MCP Plugin (`app/providers/mcp/plugin.go`)

Implements: `Plugin`, `HasConfig`, `HasCommands`, `HasFeatureFlag`, `HasMcpTools`, `HasMcpResources`, `HasMcpPrompts`, `HasRoutes`.

```go
type McpPlugin struct {
    active            bool
    workspace         string
    externalTools     []plugins.McpToolDefinition
    externalResources []plugins.McpResourceDefinition
    externalPrompts   []plugins.McpPromptDefinition
    // ...
}

func (p *McpPlugin) ID() string             { return "orchestra/mcp" }
func (p *McpPlugin) Dependencies() []string { return []string{"orchestra/workspace"} }
```

Key features:
- Depends on `orchestra/workspace` for resolving the current workspace path.
- Provides `RegisterExternalTools()`, `RegisterExternalResources()`, and `RegisterExternalPrompts()` so other plugins can push tools into the MCP server.
- Exposes `mcp:start` and `mcp:init` commands.

### Discord Plugin (`app/providers/discord/providers/plugin.go`)

Implements: `Plugin`, `HasConfig`, `HasFeatureFlag`.

```go
func (p *DiscordPlugin) ID() string             { return "orchestra/discord" }
func (p *DiscordPlugin) Dependencies() []string { return []string{"orchestra/notifications"} }
```

Key features:
- Depends on `orchestra/notifications` for the notification channel registry.
- Loads configuration from the settings service (not environment variables).
- Registers as a notification channel so other plugins can route notifications through Discord.
- Registers a settings panel via `Contributes.Register("settings.panel", ...)`.

### Notifications Plugin (`app/providers/notifications/plugin.go`)

Implements: `HasSettings`.

Contributes a settings group and four settings (enabled, sound, do not disturb, position) with select options and boolean toggles.

---

## Testing Coverage

The codebase includes comprehensive tests for three subsystems:

### IPC Protocol Tests (`ipc_protocol_test.go`)
- Message construction for all 4 types (call, response, error, event)
- Marshal/unmarshal round-trip fidelity
- Typed accessor methods (`GetParams`, `GetResult`, `GetPayload`)
- Error handling for missing fields
- All 7 error codes validated

### IPC Bridge Tests (`ipc_bridge_test.go`)
- Bridge creation and socket path generation
- Handler registration and invocation
- Message sending over mock `net.Pipe` connections
- Call handling with successful responses
- Method-not-found error responses
- Response matching via pending channel
- Event sending and receiving
- Error on operations before bridge is started
- Graceful stop when not started
- Multiple handler registration

### IPC Manager Tests (`ipc_manager_test.go`)
- Manager creation
- Empty stats retrieval
- Broadcast with no bridges
- Error handling for nonexistent plugins (stop, call, send event)
- Standard handler registration (ping, getInfo, shutdown)
- Standard handler invocation with response validation
- Multiple plugin management
- Concurrent access safety (100 goroutine reads)
- Bridge lifecycle (add, verify, stop, verify removed)
- Duplicate start prevention
- Custom handler invocation alongside standard handlers

---

## File Reference

| File | Lines | Purpose |
|------|-------|---------|
| `plugin.go` | 107 | Core `Plugin` interface, `Context`, `PluginState`, `PluginInfo` |
| `manifest.go` | 177 | `Manifest` struct, `Contributes`, validation, permissions |
| `lifecycle.go` | 197 | `LifecycleManager` -- Load/Boot/RegisterAll/Shutdown orchestration |
| `container.go` | 47 | `Container` interface, `ServiceFactory`, `ServiceLifetime` |
| `container_impl.go` | 150 | `SimpleContainer` -- thread-safe DI with circular dependency detection |
| `ipc_protocol.go` | 135 | `IPCMessage`, `IPCError`, error codes, message constructors |
| `ipc_bridge.go` | 330 | `IPCBridge` -- Unix socket IPC with process management |
| `ipc_manager.go` | 193 | `IPCManager` -- multi-plugin IPC orchestration |
| `eventbus.go` | 47 | `EventBus` interface, `Event`, `EventHandler` |
| `eventbus_impl.go` | 170 | `InProcessEventBus` -- wildcard pattern matching pub/sub |
| `resolver.go` | 119 | `DependencyResolver` -- Kahn's algorithm topological sort |
| `discovery.go` | 106 | `DiscoveryScanner` -- filesystem plugin discovery |
| `capabilities.go` | 54 | `HasSettings` capability + setting definitions |
| `capabilities_ipc.go` | 25 | `HasIPCHandlers` capability |
| `capabilities_sidebar.go` | 30 | `HasSidebarViews` capability |
| `capabilities_tray.go` | 36 | `HasTrayItems` capability |
| `capabilities_tabs.go` | 21 | `HasTabs` capability |
| `capabilities_panels.go` | 37 | `HasPanels` capability |
| `capabilities_devtools.go` | 18 | `HasDevToolsSessions` capability |
| `capabilities_context_menu.go` | 35 | `HasContextMenus` capability |
| `ipc_protocol_test.go` | 335 | Protocol message tests |
| `ipc_bridge_test.go` | 381 | Bridge communication tests |
| `ipc_manager_test.go` | 355 | Manager orchestration tests |

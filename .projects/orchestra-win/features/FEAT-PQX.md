---
created_at: "2026-02-28T02:56:01Z"
description: |-
    Implement `Orchestra.Desktop/MainWindow.xaml/.cs` — the root WinUI 3 window with plugin-driven `NavigationView`.

    **Layout:**
    - `NavigationView` with `PaneDisplayMode=LeftCompact`, `CompactPaneLength=56`, `OpenPaneLength=280`
    - Pane header: 36×36 purple logo + "Orchestra" wordmark (hidden in compact)
    - Menu items built dynamically from `PluginRegistry.SidebarPlugins` (Glyph + Name)
    - Pane footer: `ConnectionIndicator` control (green dot when connected)
    - Content area: `Frame` navigates to `plugin.PageType` on selection

    **Title bar:** Custom `ExtendsContentIntoTitleBar = true`, draggable region, Mica backdrop on Windows 11

    **App.xaml.cs:** Generic Host with DI container — registers all singleton services + built-in plugins, connects `QUICConnection` on `OnLaunched`

    **Startup sequence:** load settings → init theme → register plugins → show window → connect QUIC → subscribe events → load cache → check updates → show tray icon
id: FEAT-PQX
priority: P0
project_id: orchestra-win
status: backlog
title: MainWindow + NavigationView shell
updated_at: "2026-02-28T02:56:01Z"
version: 0
---

# MainWindow + NavigationView shell

Implement `Orchestra.Desktop/MainWindow.xaml/.cs` — the root WinUI 3 window with plugin-driven `NavigationView`.

**Layout:**
- `NavigationView` with `PaneDisplayMode=LeftCompact`, `CompactPaneLength=56`, `OpenPaneLength=280`
- Pane header: 36×36 purple logo + "Orchestra" wordmark (hidden in compact)
- Menu items built dynamically from `PluginRegistry.SidebarPlugins` (Glyph + Name)
- Pane footer: `ConnectionIndicator` control (green dot when connected)
- Content area: `Frame` navigates to `plugin.PageType` on selection

**Title bar:** Custom `ExtendsContentIntoTitleBar = true`, draggable region, Mica backdrop on Windows 11

**App.xaml.cs:** Generic Host with DI container — registers all singleton services + built-in plugins, connects `QUICConnection` on `OnLaunched`

**Startup sequence:** load settings → init theme → register plugins → show window → connect QUIC → subscribe events → load cache → check updates → show tray icon

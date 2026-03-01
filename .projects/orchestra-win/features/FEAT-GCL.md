---
created_at: "2026-02-28T02:56:45Z"
description: |-
    Implement `Orchestra.Desktop/TrayIcon/TrayIconService.cs` — system tray icon with context menu using `H.NotifyIcon.WinUI` NuGet package.

    **Tray icon:** `orchestra-tray.ico` (16/32/48px variants), ToolTip "Orchestra MCP"

    **Context menu (`MenuFlyout`):**
    - "Open Orchestra" → `mainWindow.Activate()`
    - "Spirit Mode" → `WindowModeManager.SetMode(Floating)`
    - "Bubble Mode" → `WindowModeManager.SetMode(Bubble)`
    - Separator
    - "⬤ Connected" / "○ Disconnected" (non-interactive status, updates on `QUICConnection.StateChanged`)
    - Separator
    - "Check for Updates"
    - "Quit" → `Application.Current.Exit()`

    **Double-click:** restore main window

    **Badge:** show pending notification count as tray badge (Windows 11)

    **Platform:** Desktop only (`WindowsPlatform.Desktop`)
id: FEAT-GCL
priority: P0
project_id: orchestra-win
status: backlog
title: System tray — H.NotifyIcon.WinUI
updated_at: "2026-02-28T02:56:45Z"
version: 0
---

# System tray — H.NotifyIcon.WinUI

Implement `Orchestra.Desktop/TrayIcon/TrayIconService.cs` — system tray icon with context menu using `H.NotifyIcon.WinUI` NuGet package.

**Tray icon:** `orchestra-tray.ico` (16/32/48px variants), ToolTip "Orchestra MCP"

**Context menu (`MenuFlyout`):**
- "Open Orchestra" → `mainWindow.Activate()`
- "Spirit Mode" → `WindowModeManager.SetMode(Floating)`
- "Bubble Mode" → `WindowModeManager.SetMode(Bubble)`
- Separator
- "⬤ Connected" / "○ Disconnected" (non-interactive status, updates on `QUICConnection.StateChanged`)
- Separator
- "Check for Updates"
- "Quit" → `Application.Current.Exit()`

**Double-click:** restore main window

**Badge:** show pending notification count as tray badge (Windows 11)

**Platform:** Desktop only (`WindowsPlatform.Desktop`)

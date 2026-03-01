---
created_at: "2026-02-28T02:55:50Z"
description: |-
    Implement the plugin abstraction layer in `Orchestra.Core/Plugins/` that mirrors the Go framework's plugin architecture — every screen/feature is a plugin.

    **`IOrchestraPlugin` interface:**
    - `Id`, `Name`, `IconGlyph` (Segoe Fluent glyph e.g. `\uE8BD`), `Section` (Sidebar/DevTools/Settings), `Order`, `SupportedPlatforms` (`[Flags] WindowsPlatform` enum: Desktop/Xbox/HoloLens/IoT)
    - `Type PageType` — the WinUI Page type to navigate to
    - `void OnActivate()`, `void OnDeactivate()`

    **`OrchestraPluginBase`** — abstract base with default `SupportedPlatforms = Desktop | HoloLens`

    **`PluginRegistry`** — singleton via DI:
    - `Register(IOrchestraPlugin)` with sorted insert
    - `SidebarPlugins`, `DevToolPlugins`, `SettingsPlugins` (filtered by current platform via `AnalyticsInfo.VersionInfo.DeviceFamily`)
    - `GetPlugin(id)`

    **`scripts/new-windows-plugin.ps1`** — PowerShell scaffold script: creates `Plugins/{Name}Plugin/`, `{Name}Plugin.cs`, `{Name}Page.xaml/.cs`, prints DI registration line
id: FEAT-WSP
priority: P0
project_id: orchestra-win
status: backlog
title: C# plugin system — IOrchestraPlugin + PluginRegistry
updated_at: "2026-02-28T02:55:50Z"
version: 0
---

# C# plugin system — IOrchestraPlugin + PluginRegistry

Implement the plugin abstraction layer in `Orchestra.Core/Plugins/` that mirrors the Go framework's plugin architecture — every screen/feature is a plugin.

**`IOrchestraPlugin` interface:**
- `Id`, `Name`, `IconGlyph` (Segoe Fluent glyph e.g. `\uE8BD`), `Section` (Sidebar/DevTools/Settings), `Order`, `SupportedPlatforms` (`[Flags] WindowsPlatform` enum: Desktop/Xbox/HoloLens/IoT)
- `Type PageType` — the WinUI Page type to navigate to
- `void OnActivate()`, `void OnDeactivate()`

**`OrchestraPluginBase`** — abstract base with default `SupportedPlatforms = Desktop | HoloLens`

**`PluginRegistry`** — singleton via DI:
- `Register(IOrchestraPlugin)` with sorted insert
- `SidebarPlugins`, `DevToolPlugins`, `SettingsPlugins` (filtered by current platform via `AnalyticsInfo.VersionInfo.DeviceFamily`)
- `GetPlugin(id)`

**`scripts/new-windows-plugin.ps1`** — PowerShell scaffold script: creates `Plugins/{Name}Plugin/`, `{Name}Plugin.cs`, `{Name}Page.xaml/.cs`, prints DI registration line

---
created_at: "2026-02-28T03:16:40Z"
description: |-
    Implement Xbox-specific adaptations — gamepad-first navigation and simplified project dashboard layout.

    **Navigation changes:**
    - Replace `NavigationView` with fullscreen `TabView` on Xbox device family
    - `XY focus navigation` enabled: gamepad D-pad = navigate, A = select, B = back, triggers = scroll
    - `GamepadNavigation.IsGamepadNavigationEnabled = true`

    **`XboxDashboardPage.xaml`:**
    - Full-screen 2×3 grid of project cards
    - Each card: project name, progress ring, feature count, last activity timestamp
    - Sprint burndown chart (`CartesianChart` via LiveCharts2 NuGet)
    - Real-time bottom ticker: active AI sessions, in-progress features, blocked count

    **Plugin set (Xbox):**
    - Chat: full (voice-first via STT, gamepad Back button = dictate)
    - Projects: dashboard + read-only backlog
    - Settings: AI provider, theme, voice only
    - DevTools / Notes / Docs: hidden on Xbox

    **`WindowsPlatform.Xbox`** detection:** `AnalyticsInfo.VersionInfo.DeviceFamily == "Windows.Xbox"`

    **`GamepadInputService.cs`:** `Windows.Gaming.Input.Gamepad.GamepadAdded/Removed`, button event polling on UI thread

    **Platform:** Xbox (Windows App SDK 1.5+ with Xbox support)
id: FEAT-JRG
priority: P3
project_id: orchestra-win
status: backlog
title: Xbox platform — gamepad navigation + dashboard
updated_at: "2026-02-28T03:16:40Z"
version: 0
---

# Xbox platform — gamepad navigation + dashboard

Implement Xbox-specific adaptations — gamepad-first navigation and simplified project dashboard layout.

**Navigation changes:**
- Replace `NavigationView` with fullscreen `TabView` on Xbox device family
- `XY focus navigation` enabled: gamepad D-pad = navigate, A = select, B = back, triggers = scroll
- `GamepadNavigation.IsGamepadNavigationEnabled = true`

**`XboxDashboardPage.xaml`:**
- Full-screen 2×3 grid of project cards
- Each card: project name, progress ring, feature count, last activity timestamp
- Sprint burndown chart (`CartesianChart` via LiveCharts2 NuGet)
- Real-time bottom ticker: active AI sessions, in-progress features, blocked count

**Plugin set (Xbox):**
- Chat: full (voice-first via STT, gamepad Back button = dictate)
- Projects: dashboard + read-only backlog
- Settings: AI provider, theme, voice only
- DevTools / Notes / Docs: hidden on Xbox

**`WindowsPlatform.Xbox`** detection:** `AnalyticsInfo.VersionInfo.DeviceFamily == "Windows.Xbox"`

**`GamepadInputService.cs`:** `Windows.Gaming.Input.Gamepad.GamepadAdded/Removed`, button event polling on UI thread

**Platform:** Xbox (Windows App SDK 1.5+ with Xbox support)

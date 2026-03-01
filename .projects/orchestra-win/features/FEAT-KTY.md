---
created_at: "2026-02-28T03:16:06Z"
description: "Implement `Orchestra.Widgets/` project — Windows 11 widget provider using the Widgets Board API and Adaptive Cards.\n\n**`Orchestra.Widgets.csproj`:** separate MSIX component, references `Microsoft.Windows.Widgets.Providers` from Windows App SDK.\n\n**Widgets (2 initial):**\n\n**`ProjectWidget.cs`** — shows active project status:\n```\n┌────────────────────────────┐\n│ \U0001F7E3 orchestra-win            │\n│ ████████░░ 73% complete     │\n│ 12 in-progress · 4 blocked  │\n│ [Open] [New Feature]        │\n└────────────────────────────┘\n```\nAdaptive Card JSON with `ProgressBar`, `FactSet` (status breakdown), action buttons\n\n**`SprintWidget.cs`** — current sprint burndown:\n```\n┌────────────────────────────┐\n│ Sprint 3 · 8 days left      │\n│ ▓▓▓▓▓░░░░░ 52% done         │\n│ 18/35 tasks complete        │\n│ [View Sprint] [Standup]     │\n└────────────────────────────┘\n```\n\n**`WidgetProvider` lifecycle:** `CreateWidget`, `DeleteWidget`, `OnWidgetContextChanged`, `OnActionInvoked` → dispatches to `ToolService` via named pipe IPC to main app\n\n**Refresh:** `PeriodicTimer` every 5 minutes calls `WidgetManager.GetDefault().UpdateWidget(...)`\n\n**Registration:** `Package.appxmanifest` `windows.widget` extension entry\n\n**Platform:** Windows 11 22H2+ (Widgets Board). Gracefully absent on Windows 10."
id: FEAT-KTY
priority: P2
project_id: orchestra-win
status: backlog
title: Windows 11 Widgets — Adaptive Cards (project status + sprint)
updated_at: "2026-02-28T03:16:06Z"
version: 0
---

# Windows 11 Widgets — Adaptive Cards (project status + sprint)

Implement `Orchestra.Widgets/` project — Windows 11 widget provider using the Widgets Board API and Adaptive Cards.

**`Orchestra.Widgets.csproj`:** separate MSIX component, references `Microsoft.Windows.Widgets.Providers` from Windows App SDK.

**Widgets (2 initial):**

**`ProjectWidget.cs`** — shows active project status:
```
┌────────────────────────────┐
│ 🟣 orchestra-win            │
│ ████████░░ 73% complete     │
│ 12 in-progress · 4 blocked  │
│ [Open] [New Feature]        │
└────────────────────────────┘
```
Adaptive Card JSON with `ProgressBar`, `FactSet` (status breakdown), action buttons

**`SprintWidget.cs`** — current sprint burndown:
```
┌────────────────────────────┐
│ Sprint 3 · 8 days left      │
│ ▓▓▓▓▓░░░░░ 52% done         │
│ 18/35 tasks complete        │
│ [View Sprint] [Standup]     │
└────────────────────────────┘
```

**`WidgetProvider` lifecycle:** `CreateWidget`, `DeleteWidget`, `OnWidgetContextChanged`, `OnActionInvoked` → dispatches to `ToolService` via named pipe IPC to main app

**Refresh:** `PeriodicTimer` every 5 minutes calls `WidgetManager.GetDefault().UpdateWidget(...)`

**Registration:** `Package.appxmanifest` `windows.widget` extension entry

**Platform:** Windows 11 22H2+ (Widgets Board). Gracefully absent on Windows 10.

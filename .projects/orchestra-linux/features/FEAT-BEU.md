---
created_at: "2026-02-28T02:54:03Z"
description: 'WindowModeManager cycling between three modes: Main (primary AdwApplicationWindow 1280x860), Floating (SpiritWindow 420x640), Bubble (BubbleWindow 56x56). Ctrl+Shift+O global shortcut via xdg-desktop-portal GlobalShortcuts portal (Wayland-safe) with X11 XGrabKey fallback. app.cycle-window action triggers mode change. Current mode persisted to GSettings window-mode key. Mode transitions: hide current window, show next. Main→Floating→Bubble→Main cycle.'
id: FEAT-BEU
priority: P1
project_id: orchestra-linux
status: backlog
title: Window mode cycling (Ctrl+Shift+O)
updated_at: "2026-02-28T02:54:03Z"
version: 0
---

# Window mode cycling (Ctrl+Shift+O)

WindowModeManager cycling between three modes: Main (primary AdwApplicationWindow 1280x860), Floating (SpiritWindow 420x640), Bubble (BubbleWindow 56x56). Ctrl+Shift+O global shortcut via xdg-desktop-portal GlobalShortcuts portal (Wayland-safe) with X11 XGrabKey fallback. app.cycle-window action triggers mode change. Current mode persisted to GSettings window-mode key. Mode transitions: hide current window, show next. Main→Floating→Bubble→Main cycle.

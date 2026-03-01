---
created_at: "2026-02-28T02:54:12Z"
description: 'Register global keyboard shortcut Ctrl+Shift+O via xdg-desktop-portal GlobalShortcuts interface (org.freedesktop.portal.GlobalShortcuts). On X11 sessions without portal support: fallback to GDK event filter with XGrabKey via GDK X11 display APIs. Gracefully handle portal unavailability (GNOME without portal, headless). Show first-run dialog explaining shortcut registration requires portal permission. Trigger: app.cycle-window action to cycle Main→Floating→Bubble modes.'
id: FEAT-QZI
priority: P1
project_id: orchestra-linux
status: backlog
title: Global hotkey registration (xdg-portal + X11 fallback)
updated_at: "2026-02-28T02:54:12Z"
version: 0
---

# Global hotkey registration (xdg-portal + X11 fallback)

Register global keyboard shortcut Ctrl+Shift+O via xdg-desktop-portal GlobalShortcuts interface (org.freedesktop.portal.GlobalShortcuts). On X11 sessions without portal support: fallback to GDK event filter with XGrabKey via GDK X11 display APIs. Gracefully handle portal unavailability (GNOME without portal, headless). Show first-run dialog explaining shortcut registration requires portal permission. Trigger: app.cycle-window action to cycle Main→Floating→Bubble modes.

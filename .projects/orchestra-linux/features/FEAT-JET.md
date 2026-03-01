---
created_at: "2026-02-28T02:54:31Z"
description: 'Implement ai.screenshot tool calls on the Linux frontend side. Primary: xdg-desktop-portal Screenshot interface (org.freedesktop.portal.Screenshot) — works on all Wayland compositors and X11. Returns URI of captured PNG, load via Gio.File.load_bytes_async(), encode as base64. Interactive region selection: portal PickColor or grim+slurp on Wayland, scrot -s on X11. Fallbacks by DE: GNOME (gnome-screenshot), KDE (spectacle --nonotify), XFCE (xfce4-screenshooter), Sway/Hyprland (grim+slurp). Send base64 image to ai.vision tools for analysis.'
id: FEAT-JET
priority: P2
project_id: orchestra-linux
status: backlog
title: Screenshot capture (xdg-portal + CLI fallback)
updated_at: "2026-02-28T02:54:31Z"
version: 0
---

# Screenshot capture (xdg-portal + CLI fallback)

Implement ai.screenshot tool calls on the Linux frontend side. Primary: xdg-desktop-portal Screenshot interface (org.freedesktop.portal.Screenshot) — works on all Wayland compositors and X11. Returns URI of captured PNG, load via Gio.File.load_bytes_async(), encode as base64. Interactive region selection: portal PickColor or grim+slurp on Wayland, scrot -s on X11. Fallbacks by DE: GNOME (gnome-screenshot), KDE (spectacle --nonotify), XFCE (xfce4-screenshooter), Sway/Hyprland (grim+slurp). Send base64 image to ai.vision tools for analysis.

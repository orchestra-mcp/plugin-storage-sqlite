---
created_at: "2026-02-28T03:07:00Z"
description: |-
    Implement `Orchestra.Desktop/Plugins/DevToolsPlugin/Terminal/` — full PTY terminal using Windows ConPTY API.

    **`TerminalPage.xaml`** — tab-based multi-terminal:
    - `TabView` (WinUI 3) for multiple sessions, "+ New Tab" with shell picker dropdown
    - `TerminalControl` (`WebView2` + xterm.js for rendering, ConPTY for PTY backend)
    - Shell options: PowerShell 7 (`pwsh.exe`), Command Prompt (`cmd.exe`), WSL (`wsl.exe`), Git Bash

    **ConPTY backend (`TerminalService.cs`):**
    - `CreatePseudoConsole` Win32 P/Invoke (Windows 10 1809+)
    - `CreateProcess` with ConPTY handle attached
    - Bidirectional pipe: stdin → PTY, PTY stdout → xterm.js via `WebView2.PostWebMessageAsString`
    - `ResizePseudoConsole` on resize events

    **Features:** copy/paste, 10k-line scrollback, Ctrl+scroll font zoom, hyperlink detection, 256-color + true-color ANSI

    **MCP tools called:** `terminal_create`, `terminal_write`, `terminal_read`, `terminal_resize`, `terminal_close`, `terminal_list`

    **Platform:** Desktop only
id: FEAT-DBY
priority: P1
project_id: orchestra-win
status: backlog
title: Terminal sub-plugin — PowerShell/cmd/WSL via ConPTY
updated_at: "2026-02-28T03:07:00Z"
version: 0
---

# Terminal sub-plugin — PowerShell/cmd/WSL via ConPTY

Implement `Orchestra.Desktop/Plugins/DevToolsPlugin/Terminal/` — full PTY terminal using Windows ConPTY API.

**`TerminalPage.xaml`** — tab-based multi-terminal:
- `TabView` (WinUI 3) for multiple sessions, "+ New Tab" with shell picker dropdown
- `TerminalControl` (`WebView2` + xterm.js for rendering, ConPTY for PTY backend)
- Shell options: PowerShell 7 (`pwsh.exe`), Command Prompt (`cmd.exe`), WSL (`wsl.exe`), Git Bash

**ConPTY backend (`TerminalService.cs`):**
- `CreatePseudoConsole` Win32 P/Invoke (Windows 10 1809+)
- `CreateProcess` with ConPTY handle attached
- Bidirectional pipe: stdin → PTY, PTY stdout → xterm.js via `WebView2.PostWebMessageAsString`
- `ResizePseudoConsole` on resize events

**Features:** copy/paste, 10k-line scrollback, Ctrl+scroll font zoom, hyperlink detection, 256-color + true-color ANSI

**MCP tools called:** `terminal_create`, `terminal_write`, `terminal_read`, `terminal_resize`, `terminal_close`, `terminal_list`

**Platform:** Desktop only

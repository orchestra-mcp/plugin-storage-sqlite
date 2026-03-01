---
created_at: "2026-02-28T03:12:29Z"
description: |-
    Implement `Orchestra.Desktop/Plugins/DevToolsPlugin/Docker/` — Docker Desktop integration.

    **`DockerPage.xaml`** — `TabView` with 5 tabs:
    - **Containers:** `ListView` (name, image, status dot, ports, CPU%, Memory), toolbar (Start/Stop/Restart/Remove/Logs/Shell), detail flyout (env vars, mounts, network, inspect JSON)
    - **Images:** `ListView` (name, tag, size, created) + Pull/Push/Remove/Run actions
    - **Volumes:** `ListView` + inspect flyout
    - **Networks:** `ListView` + inspect flyout
    - **Compose:** YAML editor (`WebView2` + Monaco) + Up/Down/Restart actions

    **`DockerService.cs`** — wraps Docker CLI (`docker.exe`) or Docker Engine REST API via named pipe (`npipe://./pipe/docker_engine`). Detects Docker Desktop vs. Rancher Desktop.

    **Log streaming:** `docker logs -f` → live `TextBox` via `Process.OutputDataReceived`

    **MCP tools called:** `docker_list`, `docker_start`, `docker_stop`, `docker_restart`, `docker_remove`, `docker_logs`, `docker_exec`, `docker_pull`, `docker_images`, `docker_compose`

    **Platform:** Desktop only
id: FEAT-XEU
priority: P2
project_id: orchestra-win
status: backlog
title: Docker sub-plugin — container + image management
updated_at: "2026-02-28T03:12:29Z"
version: 0
---

# Docker sub-plugin — container + image management

Implement `Orchestra.Desktop/Plugins/DevToolsPlugin/Docker/` — Docker Desktop integration.

**`DockerPage.xaml`** — `TabView` with 5 tabs:
- **Containers:** `ListView` (name, image, status dot, ports, CPU%, Memory), toolbar (Start/Stop/Restart/Remove/Logs/Shell), detail flyout (env vars, mounts, network, inspect JSON)
- **Images:** `ListView` (name, tag, size, created) + Pull/Push/Remove/Run actions
- **Volumes:** `ListView` + inspect flyout
- **Networks:** `ListView` + inspect flyout
- **Compose:** YAML editor (`WebView2` + Monaco) + Up/Down/Restart actions

**`DockerService.cs`** — wraps Docker CLI (`docker.exe`) or Docker Engine REST API via named pipe (`npipe://./pipe/docker_engine`). Detects Docker Desktop vs. Rancher Desktop.

**Log streaming:** `docker logs -f` → live `TextBox` via `Process.OutputDataReceived`

**MCP tools called:** `docker_list`, `docker_start`, `docker_stop`, `docker_restart`, `docker_remove`, `docker_logs`, `docker_exec`, `docker_pull`, `docker_images`, `docker_compose`

**Platform:** Desktop only

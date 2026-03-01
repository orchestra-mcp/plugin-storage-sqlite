---
created_at: "2026-02-28T02:53:26Z"
description: 'Services sub-tool showing systemd units. GtkColumnView with columns: service name, description, status (active/inactive/failed), enabled/disabled toggle. Color code status: green=active, red=failed, gray=inactive. Start/Stop/Restart buttons in row suffix. Service detail popover: full description, loaded state, PID, memory usage, log tail. Search bar to filter services. Refresh button. Calls devtools.services tools: list_services, start_service, stop_service, restart_service, service_logs, service_info. Uses systemd D-Bus API on Linux via tools plugin.'
id: FEAT-NOZ
priority: P2
project_id: orchestra-linux
status: backlog
title: Services manager sub-tool (systemd)
updated_at: "2026-02-28T02:53:26Z"
version: 0
---

# Services manager sub-tool (systemd)

Services sub-tool showing systemd units. GtkColumnView with columns: service name, description, status (active/inactive/failed), enabled/disabled toggle. Color code status: green=active, red=failed, gray=inactive. Start/Stop/Restart buttons in row suffix. Service detail popover: full description, loaded state, PID, memory usage, log tail. Search bar to filter services. Refresh button. Calls devtools.services tools: list_services, start_service, stop_service, restart_service, service_logs, service_info. Uses systemd D-Bus API on Linux via tools plugin.

---
created_at: "2026-02-28T02:55:44Z"
description: 'scripts/new-linux-plugin.sh — plugin creator script mirroring new-swift-plugin.sh. Usage: ./scripts/new-linux-plugin.sh my-feature sidebar. Creates: shared/src/plugins/my-feature-plugin/my-feature-plugin.vala (OrchestraPlugin implementation stub with id, name, icon_name, section, order, create_view(), on_activate(), on_deactivate()) and my-feature-view.vala (GtkWidget placeholder). Updates shared/meson.build to add new .vala files to sources list. Prints registration line to add to desktop/src/main.vala. Validates section arg (sidebar/devtools/settings). Makes script executable (chmod +x).'
id: FEAT-BUG
priority: P1
project_id: orchestra-linux
status: backlog
title: new-linux-plugin.sh scaffolding script
updated_at: "2026-02-28T02:55:44Z"
version: 0
---

# new-linux-plugin.sh scaffolding script

scripts/new-linux-plugin.sh — plugin creator script mirroring new-swift-plugin.sh. Usage: ./scripts/new-linux-plugin.sh my-feature sidebar. Creates: shared/src/plugins/my-feature-plugin/my-feature-plugin.vala (OrchestraPlugin implementation stub with id, name, icon_name, section, order, create_view(), on_activate(), on_deactivate()) and my-feature-view.vala (GtkWidget placeholder). Updates shared/meson.build to add new .vala files to sources list. Prints registration line to add to desktop/src/main.vala. Validates section arg (sidebar/devtools/settings). Makes script executable (chmod +x).

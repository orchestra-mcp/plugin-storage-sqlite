---
created_at: "2026-02-28T02:12:25Z"
description: 'Add all 22 new plugin entries to plugins.yaml. Update orchestra.json require + install-order. Remove old save_note/list_notes from tools-features. Final integration testing: orchestra serve with all plugins, verify ~290 tools in tools/list.'
id: FEAT-TBQ
labels:
    - phase-8
    - wiring
priority: P0
project_id: orchestra-tools
status: done
title: 'Wiring: plugins.yaml + orchestra.json + cleanup'
updated_at: "2026-02-28T05:25:00Z"
version: 0
---

# Wiring: plugins.yaml + orchestra.json + cleanup

Add all 22 new plugin entries to plugins.yaml. Update orchestra.json require + install-order. Remove old save_note/list_notes from tools-features. Final integration testing: orchestra serve with all plugins, verify ~290 tools in tools/list.

---
**backlog -> done**: Added tools.markdown, tools.notes, tools.docs (depends_on: tools.markdown), devtools.git to plugins.yaml. All 4 binaries built (bin/tools-markdown, bin/tools-notes, bin/tools-docs, bin/devtools-git). Makefile already had all build targets. Total plugins: 15 (was 11).

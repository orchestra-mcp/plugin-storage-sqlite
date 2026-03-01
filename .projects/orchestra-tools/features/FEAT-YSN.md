---
created_at: "2026-02-28T02:11:37Z"
description: 'Replace save_note/list_notes from tools-features. Tools: create_note, get_note, update_note, delete_note, list_notes, search_notes, pin_note, tag_note. Storage: .projects/{project}/notes/{id}.md. Cleanup: remove old note tools from tools-features/internal/tools/metadata.go.'
id: FEAT-YSN
labels:
    - phase-2
    - core-plugin
    - notes
priority: P1
project_id: orchestra-tools
status: done
title: Standalone notes plugin (tools.notes)
updated_at: "2026-02-28T05:25:00Z"
version: 0
---

# Standalone notes plugin (tools.notes)

Replace save_note/list_notes from tools-features. Tools: create_note, get_note, update_note, delete_note, list_notes, search_notes, pin_note, tag_note. Storage: .projects/{project}/notes/{id}.md. Cleanup: remove old note tools from tools-features/internal/tools/metadata.go.

---
**in-progress -> done**: 13 tests passing (create_note, get_note, list_notes, delete_note, search_notes, tag/pin validation, storage helpers). In-memory mock StorageClient used for testing. Binary built to bin/tools-notes. Wired into plugins.yaml.

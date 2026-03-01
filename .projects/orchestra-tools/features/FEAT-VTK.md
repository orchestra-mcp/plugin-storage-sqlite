---
created_at: "2026-02-28T02:11:35Z"
depends_on:
    - FEAT-YPX
description: 'Thin orchestration over engine-rag. Tools: doc_create, doc_get, doc_update, doc_delete, doc_list, doc_search, doc_generate, doc_index, doc_tree, doc_export. Storage: .projects/{project}/docs/{slug}.md. Cross-plugin calls to engine-rag for parse_file, get_symbols, search, search_memory, index_file. Depends on PLUGIN-MARKDOWN.'
id: FEAT-VTK
labels:
    - phase-2
    - core-plugin
    - docs
priority: P1
project_id: orchestra-tools
status: done
title: Wiki / documentation plugin (tools.docs)
updated_at: "2026-02-28T05:25:00Z"
version: 0
---

# Wiki / documentation plugin (tools.docs)

Thin orchestration over engine-rag. Tools: doc_create, doc_get, doc_update, doc_delete, doc_list, doc_search, doc_generate, doc_index, doc_tree, doc_export. Storage: .projects/{project}/docs/{slug}.md. Cross-plugin calls to engine-rag for parse_file, get_symbols, search, search_memory, index_file. Depends on PLUGIN-MARKDOWN.

---
**in-progress -> done**: 14 tests passing (doc_create with category/tags, doc_get, doc_list with category filter, doc_delete, doc_search). In-memory mock StorageClient. Binary built to bin/tools-docs. Wired into plugins.yaml with depends_on: tools.markdown.

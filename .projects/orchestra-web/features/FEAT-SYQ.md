---
blocks:
    - FEAT-OMM
created_at: "2026-02-28T03:27:23Z"
depends_on:
    - FEAT-TFH
description: |-
    Documentation browser backed by plugin-tools-docs, rendered with @orchestra-mcp/editor.

    File: `apps/web/src/pages/docs.tsx` + `apps/web/src/pages/docs/[slug].tsx`

    Docs index (`docs.tsx`):
    - Two-panel layout: left nav tree (240px) + right content area
    - Nav tree: nested list of doc sections via mcp.callTool("list_docs")
    - Each item: title, section badge, click → navigate to /docs/:slug
    - Search Input at top of nav → filter docs list
    - Skeleton from @orchestra-mcp/ui during load

    Doc viewer (`docs/[slug].tsx`):
    - Fetch content via mcp.callTool("get_doc", {slug})
    - Render with MarkdownRenderer from @orchestra-mcp/editor (full markdown + syntax highlighting)
    - Table of contents auto-generated from headings (h2/h3), sticky on right side
    - Breadcrumb from @orchestra-mcp/ui at top
    - Prev/Next navigation buttons at bottom (from sibling docs list)
    - "Edit" Button (if doc is editable) → switches MarkdownRenderer → MarkdownEditor inline, saves via mcp.callTool("update_doc")

    Sidebar nav item: "Docs" with bx-book icon

    Acceptance: docs list renders with nav tree, clicking loads doc with full markdown rendering, TOC highlights active heading on scroll, search filters nav
id: FEAT-SYQ
priority: P2
project_id: orchestra-web
status: backlog
title: Docs Viewer Page
updated_at: "2026-02-28T03:28:20Z"
version: 0
---

# Docs Viewer Page

Documentation browser backed by plugin-tools-docs, rendered with @orchestra-mcp/editor.

File: `apps/web/src/pages/docs.tsx` + `apps/web/src/pages/docs/[slug].tsx`

Docs index (`docs.tsx`):
- Two-panel layout: left nav tree (240px) + right content area
- Nav tree: nested list of doc sections via mcp.callTool("list_docs")
- Each item: title, section badge, click → navigate to /docs/:slug
- Search Input at top of nav → filter docs list
- Skeleton from @orchestra-mcp/ui during load

Doc viewer (`docs/[slug].tsx`):
- Fetch content via mcp.callTool("get_doc", {slug})
- Render with MarkdownRenderer from @orchestra-mcp/editor (full markdown + syntax highlighting)
- Table of contents auto-generated from headings (h2/h3), sticky on right side
- Breadcrumb from @orchestra-mcp/ui at top
- Prev/Next navigation buttons at bottom (from sibling docs list)
- "Edit" Button (if doc is editable) → switches MarkdownRenderer → MarkdownEditor inline, saves via mcp.callTool("update_doc")

Sidebar nav item: "Docs" with bx-book icon

Acceptance: docs list renders with nav tree, clicking loads doc with full markdown rendering, TOC highlights active heading on scroll, search filters nav

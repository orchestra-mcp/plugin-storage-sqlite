---
blocks:
    - FEAT-VTK
created_at: "2026-02-28T02:11:23Z"
description: 'Go plugin returning typed AST JSON. Tools: md_parse, md_parse_file, md_parse_frontmatter, md_render_html, md_render_plaintext, md_toc, md_lint, md_transform. Supports GFM + Mermaid + KaTeX + footnotes + embeds. Uses goldmark. Single source of truth for Swift/React/RN.'
id: FEAT-YPX
labels:
    - phase-2
    - core-plugin
    - markdown
priority: P0
project_id: orchestra-tools
status: done
title: Markdown AST parser plugin (tools.markdown)
updated_at: "2026-02-28T05:25:00Z"
version: 0
---

# Markdown AST parser plugin (tools.markdown)

Go plugin returning typed AST JSON. Tools: md_parse, md_parse_file, md_parse_frontmatter, md_render_html, md_render_plaintext, md_toc, md_lint, md_transform. Supports GFM + Mermaid + KaTeX + footnotes + embeds. Uses goldmark. Single source of truth for Swift/React/RN.


---
**in-progress -> ready-for-testing**: Implementation complete: 8 tools (md_parse, md_parse_file, md_parse_frontmatter, md_render_html, md_render_plaintext, md_toc, md_lint, md_transform) built with goldmark v1.7.16. Full AST-to-JSON converter with GFM + Footnote + Typographer support. Lint rules: trailing whitespace, missing final newline, consecutive blank lines, heading level jumps, missing image alt. Transforms: add_heading_ids, resolve_relative_links, strip_html, normalize_whitespace. go build + go vet pass clean.

**ready-for-testing -> done**: 29 tests passing across parser (8) and tools (21). All 8 tools tested: md_parse, md_render_html, md_render_plaintext, md_parse_frontmatter, md_toc, md_lint (4 lint rules), md_transform (6 transforms). Binary built to bin/tools-markdown. Wired into plugins.yaml.

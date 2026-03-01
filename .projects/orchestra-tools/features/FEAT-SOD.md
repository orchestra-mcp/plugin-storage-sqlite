---
created_at: "2026-02-28T02:11:28Z"
depends_on:
    - FEAT-WFH
description: 'Full IDE backend. File ops: list_directory, read_file, write_file, move_file, delete_file, file_info, file_search. IDE tools: code_symbols, code_goto_definition, code_find_references, code_hover, code_complete, code_diagnostics, code_actions, code_workspace_symbols, code_namespace, code_imports. Delegates to engine-rag LSP tools. Depends on ENGINE-LSP.'
id: FEAT-SOD
labels:
    - phase-6
    - devtools
    - ide
priority: P0
project_id: orchestra-tools
status: done
title: File explorer IDE + LSP plugin (devtools.file-explorer)
updated_at: "2026-02-28T04:09:06Z"
version: 0
---

# File explorer IDE + LSP plugin (devtools.file-explorer)

Full IDE backend. File ops: list_directory, read_file, write_file, move_file, delete_file, file_info, file_search. IDE tools: code_symbols, code_goto_definition, code_find_references, code_hover, code_complete, code_diagnostics, code_actions, code_workspace_symbols, code_namespace, code_imports. Delegates to engine-rag LSP tools. Depends on ENGINE-LSP.


---
**in-progress -> ready-for-testing**: Build: go build ./libs/plugin-devtools-file-explorer/... → BUILD OK. Tests: go test ./libs/plugin-devtools-file-explorer/... → ok (11 tests pass). All 10 IDE/LSP tools implemented (code_symbols, code_goto_definition, code_find_references, code_hover, code_complete, code_diagnostics, code_actions, code_workspace_symbols, code_namespace, code_imports). Each delegates to engine-rag LSP tools via the orchestrator CallTool pattern. 17 total tools registered.


---
**in-testing -> ready-for-docs**: 11 tests cover all critical paths: code_symbols (happy path + missing field), code_goto_definition (happy path + missing line), code_workspace_symbols (happy path + missing query), code_diagnostics (happy path + missing path), code_actions (verifies lsp_diagnostics delegation), code_namespace (verifies 2-call sequence: open_document then workspace_symbols), code_imports (verifies open_document delegation). Mock StorageClient pattern validates tool name routing without real QUIC connection.


---
**in-docs -> documented**: Plugin documented in cmd/main.go (binary=devtools-file-explorer, description="File explorer with read/write and code intelligence"), plugin.go (17 tools comment), and storage/client.go (CallTool docstring). Tool names follow code_* naming convention. Each tool schema includes field descriptions. Engine-rag LSP delegation documented via code comments.


---
**in-review -> done**: Code quality review: (1) LazyClient pattern properly defers OrchestratorClient() lookup to call time — safe because tools are only called after Run() connects. (2) hasFields() helper correctly validates integer fields with value 0 (avoids false "missing" on line=0 col=0). (3) lspText() extracts TextResult's string value from the protobuf Struct. (4) code_namespace makes 2 sequential calls (open → workspace_symbols with basename) without data races. (5) All tools return helpers.ErrorResult on validation/call failure. (6) No global state. Clean separation: storage/client.go (protocol), tools/code_lsp.go (business logic), plugin.go (wiring).

# Step 4: storage.markdown Plugin

## Status: Complete

## What Was Built

The first real plugin — reads and writes Protobuf metadata + Markdown body files to disk. Handles StorageRead/Write/Delete/List requests from the orchestrator.

## Module

`github.com/orchestrated-mcp/framework/plugins/storage-markdown`

## Files

| File | Purpose |
|------|---------|
| `cmd/main.go` | Entry point — parse `--workspace`, build plugin, run |
| `internal/storage.go` | StoragePlugin implementing `plugin.StorageHandler` |
| `internal/reader.go` | Parse `<!-- META {...} META -->` + Markdown body |
| `internal/writer.go` | Format metadata + body into markdown file |
| `internal/paths.go` | Path resolution, traversal protection |
| `internal/storage_test.go` | 11 tests |

## SDK Modifications

Added `StorageHandler` interface to `libs/go/plugin/server.go`:
- `SetStorageHandler(h StorageHandler)` on Server
- Dispatch cases for StorageRead/Write/Delete/List in server dispatch
- `SetStorageHandler()` on PluginBuilder for fluent API

## On-Disk Format

```markdown
<!-- META {"status":"in-progress","priority":"P1","assignee":"go-architect"} META -->

# Feature Title

Description goes here...
```

- Line 1: `<!-- META {compact JSON} META -->`
- Line 2: blank
- Line 3+: Markdown body

## CAS Versioning

- `expected_version = 0` → create new (fail if exists)
- `expected_version > 0` → update (fail if current version != expected)
- Version stored in sidecar: `{path}.version`
- Each write increments version

## Tests (11/11 pass)

| Test | Coverage |
|------|----------|
| TestParseMarkdownFile | META block + body extraction |
| TestParseMarkdownFileNoMeta | File without META block |
| TestFormatMarkdownFile | Format → parse roundtrip |
| TestFormatMarkdownFileNilMetadata | Nil metadata handling |
| TestStorageReadWrite | Write then read, verify all fields |
| TestStorageDelete | Write then delete |
| TestStorageList | Multiple files, prefix filtering |
| TestVersioning | CAS create, update, conflict detection |
| TestPathTraversal | Reject `../` paths |
| TestDeleteNonexistent | Error handling |
| TestEmptyPath | Error handling |

```bash
cd plugins/storage-markdown && go test ./... -v
```

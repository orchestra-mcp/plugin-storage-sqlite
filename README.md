# Orchestra Storage SQLite Plugin

SQLite-based storage plugin that persists structured data to `~/.orchestra/db/` with async markdown export for git visibility.

## Install

```bash
go get github.com/orchestra-mcp/plugin-storage-sqlite
```

## Usage

```bash
# Build
go build -o bin/storage-sqlite ./cmd/

# Run (started automatically by the orchestrator)
bin/storage-sqlite --workspace /path/to/project --orchestrator-addr localhost:9100
```

## Architecture

```
SQLite (source of truth)           Markdown (git visibility)
~/.orchestra/db/<hash>.db    →     {workspace}/.projects/
                              async dual-write
```

- **SQLite primary**: All reads/writes go to SQLite with WAL mode for concurrent access
- **Markdown export**: Async dual-write to `.projects/` so features, persons, plans stay git-trackable
- **Pure Go**: Uses `modernc.org/sqlite` (no CGo) for cross-platform binary distribution
- **Auto-migration**: Imports existing `.projects/` markdown data on first run

## Database Location

Each workspace gets its own database at `~/.orchestra/db/<hash>.db` where `<hash>` is the first 16 hex chars of SHA-256(workspace path). A global index at `~/.orchestra/db/index.json` maps hashes to workspace paths for cross-workspace queries.

## Schema

| Table | Description |
|-------|-------------|
| `projects` | Project metadata (slug, name, description) |
| `features` | Features with status, priority, kind, labels, dependencies |
| `persons` | Team member profiles (name, email, role, integrations) |
| `plans` | Plans for breaking down large tasks |
| `requests` | Queued user requests |
| `assignment_rules` | Auto-assignment rules (kind → person) |
| `notes` | Notes with tags, pinning, soft-delete |
| `wip_config` | WIP limits per project |
| `sessions` | AI chat sessions |
| `session_turns` | Individual message turns within sessions |
| `packs` | Installed marketplace packs |
| `stacks` | Detected technology stacks |
| `kv_store` | Generic key-value fallback for untyped paths |

## Path Routing

Storage paths are automatically routed to the correct SQL table:

| Path Pattern | Table |
|-------------|-------|
| `{project}/features/FEAT-ABC.md` | `features` |
| `{project}/persons/PERS-ABC.md` | `persons` |
| `{project}/plans/PLAN-ABC.md` | `plans` |
| `{project}/requests/REQ-ABC.md` | `requests` |
| `{project}/assignment-rules/RULE-ABC.md` | `assignment_rules` |
| `{project}/notes/note-abc.md` | `notes` |
| `{project}/project.json` | `projects` |
| `{project}/wip.json` | `wip_config` |
| `bridge/sessions/{uuid}.md` | `sessions` |
| `bridge/sessions/{uuid}/turn-NNN.md` | `session_turns` |
| `.packs/*` | `packs` |
| `.stacks/*` | `stacks` |
| anything else | `kv_store` |

## Optimistic Concurrency (CAS)

| `expected_version` | Behavior |
|---|---|
| `0` | **Create**: Fails if the entity already exists |
| `-1` | **Upsert**: Write unconditionally |
| `> 0` | **Update**: Fails if the current version does not match |

## Supported Operations

| Operation | Description |
|-----------|-------------|
| **StorageRead** | Read metadata + body from SQLite |
| **StorageWrite** | Write with CAS versioning + async markdown export |
| **StorageDelete** | Delete from SQLite + remove markdown file |
| **StorageList** | List entities by prefix with typed table queries |

## Related Packages

| Package | Description |
|---------|-------------|
| [sdk-go](https://github.com/orchestra-mcp/sdk-go) | Plugin SDK this plugin is built on |
| [orchestrator](https://github.com/orchestra-mcp/orchestrator) | Central hub that loads this plugin |
| [plugin-storage-markdown](https://github.com/orchestra-mcp/plugin-storage-markdown) | File-based storage (fallback) |
| [plugin-tools-features](https://github.com/orchestra-mcp/plugin-tools-features) | Feature tools that use this storage |

## License

[MIT](LICENSE)

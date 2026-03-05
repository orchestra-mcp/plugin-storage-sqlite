# Contributing to plugin-storage-sqlite

## Prerequisites

- Go 1.23+
- `gofmt`, `go vet`

## Development Setup

```bash
git clone https://github.com/orchestra-mcp/plugin-storage-sqlite.git
cd plugin-storage-sqlite
go mod download
go build ./cmd/...
```

## Running Locally

```bash
go build -o storage-sqlite ./cmd/
./storage-sqlite --workspace=. --listen-addr=localhost:0 --certs-dir=~/.orchestra/certs
```

The plugin prints `READY <addr>` to stderr when it is ready to accept QUIC connections.

## Running Tests

```bash
go test ./...
```

Tests use temporary directories and do not require a running orchestrator.

## Code Style

- Run `gofmt` on all files.
- Run `go vet ./...` before committing.
- All exported functions and types must have doc comments.
- Error handling: wrap errors with context via `fmt.Errorf`.

## Key Implementation Details

- **Thread safety**: A `sync.Mutex` guards concurrent writes to prevent race conditions during version checking.
- **Path routing**: Storage paths are parsed and routed to typed SQL tables. Unknown paths fall back to `kv_store`.
- **CAS versioning**: Compare-and-swap with `expected_version` (0=create, -1=upsert, >0=must match).
- **Dual-write**: SQLite is source of truth; markdown files are exported asynchronously for git visibility.
- **Pure Go SQLite**: Uses `modernc.org/sqlite` — no CGo required, cross-platform.
- **WAL mode**: Write-Ahead Logging for concurrent read access during writes.

## Pull Request Process

1. Fork the repository and create a feature branch from `main`.
2. Write or update tests for your changes.
3. Run `go test ./...` and `go vet ./...`.
4. Update `docs/SCHEMA.md` if changing the database schema.

## Related Repositories

- [orchestra-mcp/proto](https://github.com/orchestra-mcp/proto) -- Protobuf schema
- [orchestra-mcp/sdk-go](https://github.com/orchestra-mcp/sdk-go) -- Go Plugin SDK
- [orchestra-mcp/orchestrator](https://github.com/orchestra-mcp/orchestrator) -- Central hub
- [orchestra-mcp/plugin-storage-markdown](https://github.com/orchestra-mcp/plugin-storage-markdown) -- File-based storage
- [orchestra-mcp](https://github.com/orchestra-mcp) -- Organization home

# Step 7: Wire It All Together

## Status: Complete

## What Was Built

The final wiring step — go.work, Makefile, default plugins.yaml, and end-to-end test script. This step validates the complete Phase 1 flow:

```
Agent → transport.stdio (JSON-RPC) → orchestrator (QUIC) → tools.features (QUIC)
                                                                    ↓
Disk ← storage.markdown (QUIC) ← orchestrator ← tools.features (StorageWrite)
```

## Files

| File | Purpose |
|------|---------|
| `go.work` | Go workspace linking all 6 modules |
| `Makefile` | Build, test, clean, proto, e2e targets |
| `plugins.yaml` | Default config: storage.markdown + tools.features |
| `scripts/test-e2e.sh` | Full integration test (7 steps) |

## Go Workspace (`go.work`)

Links 6 modules so `go build`, `go test`, and IDE navigation work across the entire codebase with local module references:

```
gen/go                     (proto generated code)
libs/go                    (plugin SDK)
services/orchestrator      (central hub)
plugins/storage-markdown   (disk storage)
plugins/tools-features     (34 workflow tools)
plugins/transport-stdio    (MCP JSON-RPC bridge)
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make proto` | Lint + generate proto code |
| `make build` | Build all 4 binaries to `bin/` |
| `make build-orchestrator` | Build orchestrator only |
| `make build-storage-markdown` | Build storage plugin only |
| `make build-tools-features` | Build tools plugin only |
| `make build-transport-stdio` | Build transport plugin only |
| `make test` | Run all unit tests (66 tests) |
| `make test-e2e` | Build + run end-to-end integration test |
| `make clean` | Remove `bin/` and certs |

## Binaries (4)

| Binary | Path | Description |
|--------|------|-------------|
| `orchestrator` | `bin/orchestrator` | Central hub — starts plugins, routes messages |
| `storage-markdown` | `bin/storage-markdown` | Disk storage with META+Markdown format |
| `tools-features` | `bin/tools-features` | 34 feature workflow tools |
| `transport-stdio` | `bin/transport-stdio` | MCP stdin/stdout bridge |

## E2E Test (`scripts/test-e2e.sh`)

7-step integration test using temporary workspace:

| Step | What It Tests |
|------|---------------|
| 1 | Start orchestrator → boots storage.markdown + tools.features |
| 2 | MCP `initialize` → returns protocol version + capabilities |
| 3 | MCP `tools/list` → returns all 34 tools including create_project, create_feature |
| 4 | `create_project` → project created on disk |
| 5 | `create_feature` → FEAT-XXX.md file created with META block |
| 6 | Verify on disk → file exists, contains `<!-- META {...} META -->` |
| 7 | `get_feature` → reads feature back via QUIC, verifies title |

## Bug Fixes During Wiring

### 1. Plugin READY Address (server.go)

Plugins printed `READY localhost:0` instead of the actual bound address. Fixed by adding `ActualAddr()` to Server that blocks until the listener binds, then returns the real address.

### 2. Storage Route Registration (config.go, loader.go)

The orchestrator's loader built a fresh manifest from scratch and only populated `ProvidesTools` from ListTools queries. It never learned storage capabilities. Fixed by adding `provides_storage` to the plugin config YAML and populating the manifest from config in the loader.

## Verification

```bash
# Build all binaries
make build

# Run all 66 unit tests
make test

# Run end-to-end integration test
make test-e2e
```

## Phase 1 Summary

| Step | Component | Tests |
|------|-----------|-------|
| 1 | Proto contract | buf lint + go build |
| 2 | Plugin SDK | 11 tests |
| 3 | Orchestrator | 10 tests |
| 4 | storage.markdown | 11 tests |
| 5 | tools.features | 24 tests |
| 6 | transport.stdio | 20 tests |
| 7 | Wire together | 7-step E2E test |
| **Total** | **6 modules, 4 binaries** | **66 unit + 1 E2E** |

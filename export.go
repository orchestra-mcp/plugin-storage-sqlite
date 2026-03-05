package storagesqlite

import (
	"github.com/orchestra-mcp/plugin-storage-sqlite/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// NewStorage creates the SQLite storage handler for the given workspace directory.
// The actual database file is stored at ~/.orchestra/db/<hash>.db, not inside the project.
// Markdown files are also written to .projects/ for git visibility (dual-write).
func NewStorage(workspace string) plugin.StorageHandler {
	return internal.NewStoragePlugin(workspace)
}

package internal

import (
	"database/sql"
	"log"
	"sync"

	"github.com/orchestra-mcp/sdk-go/plugin"
)

// StoragePlugin implements plugin.StorageHandler backed by SQLite.
// SQLite is the source of truth; markdown files are also written for git visibility.
type StoragePlugin struct {
	workspace string
	db        *sql.DB
	mu        sync.Mutex
}

// NewStoragePlugin creates a new SQLite storage plugin for the given workspace.
func NewStoragePlugin(workspace string) plugin.StorageHandler {
	dbPath := DBPath(workspace)

	db, err := OpenDB(dbPath)
	if err != nil {
		panic("storage-sqlite: open database: " + err.Error())
	}

	if err := InitSchema(db); err != nil {
		panic("storage-sqlite: init schema: " + err.Error())
	}

	// Register workspace in global index for cross-workspace queries.
	_ = RegisterWorkspace(workspace)

	// Auto-migrate existing .projects/ data on first run.
	if err := MigrateFromMarkdown(db, workspace); err != nil {
		log.Printf("storage-sqlite: migration warning: %v", err)
	}

	return &StoragePlugin{
		workspace: workspace,
		db:        db,
	}
}

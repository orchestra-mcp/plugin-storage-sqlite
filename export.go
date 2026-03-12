package storagesqlite

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/orchestra-mcp/plugin-storage-sqlite/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// NewStorage creates the SQLite storage handler for the given workspace directory.
// The actual database file is stored at ~/.orchestra/db/<hash>.db, not inside the project.
// Markdown files are also written to .projects/ for git visibility (dual-write).
func NewStorage(workspace string) plugin.StorageHandler {
	return internal.NewStoragePlugin(workspace)
}

// BootstrapProject creates a project record in the workspace's SQLite database
// and writes a matching project.json to .projects/ for git visibility.
// Safe to call multiple times — skips if the project already exists.
// This is meant to be called from `orchestra init` before the server starts.
func BootstrapProject(workspace, slug, name string) error {
	dbPath := internal.DBPath(workspace)

	db, err := internal.OpenDB(dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	if err := internal.InitSchema(db); err != nil {
		return fmt.Errorf("init schema: %w", err)
	}

	_ = internal.RegisterWorkspace(workspace)

	now := time.Now().UTC().Format(time.RFC3339)

	_, err = db.Exec(`INSERT OR IGNORE INTO projects (slug, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)`, slug, name, now, now)
	if err != nil {
		return fmt.Errorf("insert project: %w", err)
	}

	// Dual-write: create .projects/<slug>/project.json for git visibility.
	projDir := filepath.Join(workspace, ".projects", slug)
	os.MkdirAll(projDir, 0755)
	projJSON := fmt.Sprintf(`{"slug":%q,"name":%q,"created_at":%q}`, slug, name, now)
	_ = os.WriteFile(filepath.Join(projDir, "project.json"), []byte(projJSON), 0644)

	return nil
}

package storagesqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/orchestra-mcp/plugin-storage-sqlite/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// NewStorage creates the SQLite storage handler for the given workspace directory.
// The actual database file is stored at ~/.orchestra/db/<hash>.db, not inside the project.
func NewStorage(workspace string) plugin.StorageHandler {
	return internal.NewStoragePlugin(workspace)
}

// BootstrapProject creates a project record in the workspace's SQLite database.
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

	return nil
}

// openWorkspaceDB opens the SQLite database for a workspace (read-only queries).
func openWorkspaceDB(workspace string) (*sql.DB, error) {
	dbPath := internal.DBPath(workspace)
	return internal.OpenDB(dbPath)
}

// QuerySkillNames returns sorted skill names (slugs) from the workspace database.
func QuerySkillNames(workspace string) []string {
	db, err := openWorkspaceDB(workspace)
	if err != nil {
		return nil
	}
	rows, err := db.Query(`SELECT slug FROM skills ORDER BY slug`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			names = append(names, name)
		}
	}
	return names
}

// QueryAgentNames returns sorted agent names (slugs) from the workspace database.
func QueryAgentNames(workspace string) []string {
	db, err := openWorkspaceDB(workspace)
	if err != nil {
		return nil
	}
	rows, err := db.Query(`SELECT slug FROM agents ORDER BY slug`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			names = append(names, name)
		}
	}
	return names
}

// QueryProjectSkillNames returns sorted skill slugs included in a project.
func QueryProjectSkillNames(workspace, projectID string) []string {
	db, err := openWorkspaceDB(workspace)
	if err != nil {
		return nil
	}
	rows, err := db.Query(`SELECT s.slug FROM skills s
		JOIN project_skills ps ON s.id = ps.skill_id
		WHERE ps.project_id = ? AND ps.enabled = 1
		ORDER BY s.slug`, projectID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			names = append(names, name)
		}
	}
	return names
}

// QueryProjectAgentNames returns sorted agent slugs included in a project.
func QueryProjectAgentNames(workspace, projectID string) []string {
	db, err := openWorkspaceDB(workspace)
	if err != nil {
		return nil
	}
	rows, err := db.Query(`SELECT a.slug FROM agents a
		JOIN project_agents pa ON a.id = pa.agent_id
		WHERE pa.project_id = ? AND pa.enabled = 1
		ORDER BY a.slug`, projectID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			names = append(names, name)
		}
	}
	return names
}

// ExportToMarkdown exports all SQLite data to .projects/ as markdown files.
// This is an on-demand operation — only called when the user explicitly requests export.
// Returns the number of entities exported.
func ExportToMarkdown(workspace string) (int, error) {
	return internal.ExportToMarkdown(workspace)
}

// ChangeLogEntry represents a single change tracked for sync.
type ChangeLogEntry = internal.ChangeLogEntry

// GetUnsyncedChanges returns all unsynced change_log entries for the workspace.
func GetUnsyncedChanges(workspace string, limit int) ([]ChangeLogEntry, error) {
	return internal.GetUnsyncedChanges(workspace, limit)
}

// MarkChangesSynced marks entries as synced up to (and including) maxID.
func MarkChangesSynced(workspace string, maxID int64) error {
	return internal.MarkChangesSynced(workspace, maxID)
}

// GetSyncCursor returns the last sync cursor value for the given key.
func GetSyncCursor(workspace, key string) (string, error) {
	return internal.GetSyncCursor(workspace, key)
}

// SetSyncCursor stores a sync cursor value for the given key.
func SetSyncCursor(workspace, key, value string) error {
	return internal.SetSyncCursor(workspace, key, value)
}


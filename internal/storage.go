package internal

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/orchestra-mcp/sdk-go/plugin"
)

// StoragePlugin implements plugin.StorageHandler backed by SQLite.
// SQLite is the sole source of truth. Use ExportToMarkdown for on-demand export.
type StoragePlugin struct {
	workspace string
	db        *sql.DB
	mu        sync.Mutex
}

// recordChange appends an entry to the change_log table for sync tracking.
// Must be called while s.mu is held (inside Write/Delete).
func (s *StoragePlugin) recordChange(entityType, entityID, action string, version int64) {
	if entityType == "" || entityType == "kv_store" {
		return // skip KV entries — not syncable
	}
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(`INSERT INTO change_log (entity_type, entity_id, action, version, timestamp)
		VALUES (?, ?, ?, ?, ?)`, entityType, entityID, action, version, now)
	if err != nil {
		log.Printf("[storage-sqlite] change_log write error: %v", err)
	}
}

// ChangeLogEntry represents a single change tracked for sync.
type ChangeLogEntry struct {
	ID         int64  `json:"id"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Action     string `json:"action"`
	Version    int64  `json:"version"`
	Timestamp  string `json:"timestamp"`
	Synced     bool   `json:"synced"`
}

// GetUnsyncedChanges returns all change_log entries not yet synced.
func GetUnsyncedChanges(workspace string, limit int) ([]ChangeLogEntry, error) {
	db, err := OpenDB(DBPath(workspace))
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 1000
	}
	rows, err := db.Query(`SELECT id, entity_type, entity_id, action, version, timestamp
		FROM change_log WHERE synced = 0 ORDER BY id ASC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []ChangeLogEntry
	for rows.Next() {
		var e ChangeLogEntry
		if err := rows.Scan(&e.ID, &e.EntityType, &e.EntityID, &e.Action, &e.Version, &e.Timestamp); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// MarkChangesSynced marks change_log entries as synced up to (and including) maxID.
func MarkChangesSynced(workspace string, maxID int64) error {
	db, err := OpenDB(DBPath(workspace))
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE change_log SET synced = 1 WHERE id <= ? AND synced = 0`, maxID)
	return err
}

// GetSyncCursor returns the last sync cursor value for the given key.
func GetSyncCursor(workspace, key string) (string, error) {
	db, err := OpenDB(DBPath(workspace))
	if err != nil {
		return "", err
	}
	var val string
	err = db.QueryRow(`SELECT value FROM sync_state WHERE key = ?`, key).Scan(&val)
	if err != nil {
		return "", nil // not found is ok
	}
	return val, nil
}

// SetSyncCursor stores a sync cursor value for the given key.
func SetSyncCursor(workspace, key, value string) error {
	db, err := OpenDB(DBPath(workspace))
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO sync_state (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`, key, value)
	return err
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

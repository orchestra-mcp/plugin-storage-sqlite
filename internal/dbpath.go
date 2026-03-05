package internal

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// dbDir returns the global database directory at ~/.orchestra/db/.
func dbDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".orchestra", "db")
}

// workspaceHash returns the first 16 hex chars of SHA-256(absWorkspace).
func workspaceHash(absWorkspace string) string {
	h := sha256.Sum256([]byte(absWorkspace))
	return fmt.Sprintf("%x", h[:8]) // 8 bytes = 16 hex chars
}

// DBPath returns the SQLite database file path for a workspace.
// Format: ~/.orchestra/db/<sha256(absWorkspace)[:16]>.db
func DBPath(absWorkspace string) string {
	return filepath.Join(dbDir(), workspaceHash(absWorkspace)+".db")
}

// GlobalDBPath returns the global database path for non-workspace config.
// Format: ~/.orchestra/db/global.db
func GlobalDBPath() string {
	return filepath.Join(dbDir(), "global.db")
}

// indexFile is the workspace index at ~/.orchestra/db/index.json.
// Maps hash → absolute workspace path for cross-workspace queries.
var indexFile = filepath.Join(dbDir(), "index.json")

var indexMu sync.Mutex

// RegisterWorkspace adds a workspace to the global index so it can be
// discovered for cross-workspace queries.
func RegisterWorkspace(absWorkspace string) error {
	indexMu.Lock()
	defer indexMu.Unlock()

	index := loadIndex()
	hash := workspaceHash(absWorkspace)
	index[hash] = absWorkspace

	return saveIndex(index)
}

// ListWorkspaces returns all registered workspace paths from the index.
func ListWorkspaces() map[string]string {
	indexMu.Lock()
	defer indexMu.Unlock()
	return loadIndex()
}

func loadIndex() map[string]string {
	index := make(map[string]string)
	data, err := os.ReadFile(filepath.Join(dbDir(), "index.json"))
	if err != nil {
		return index
	}
	_ = json.Unmarshal(data, &index)
	return index
}

func saveIndex(index map[string]string) error {
	dir := dbDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "index.json"), data, 0644)
}

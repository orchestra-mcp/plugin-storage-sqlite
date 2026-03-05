package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	dbCache   = make(map[string]*sql.DB)
	dbCacheMu sync.Mutex
)

// OpenDB opens (or reuses) a SQLite database at the given path with
// production pragmas: WAL mode, busy timeout, foreign keys, synchronous=NORMAL.
func OpenDB(dbPath string) (*sql.DB, error) {
	dbCacheMu.Lock()
	defer dbCacheMu.Unlock()

	if db, ok := dbCache[dbPath]; ok {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		delete(dbCache, dbPath)
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Set pragmas for performance and correctness.
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA foreign_keys=ON",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA cache_size=-20000", // 20MB cache
		"PRAGMA temp_store=MEMORY",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			db.Close()
			return nil, fmt.Errorf("set pragma %q: %w", p, err)
		}
	}

	db.SetMaxOpenConns(1) // SQLite is single-writer

	dbCache[dbPath] = db
	return db, nil
}

// CloseAll closes all cached database connections. Call on shutdown.
func CloseAll() {
	dbCacheMu.Lock()
	defer dbCacheMu.Unlock()

	for path, db := range dbCache {
		db.Close()
		delete(dbCache, path)
	}
}

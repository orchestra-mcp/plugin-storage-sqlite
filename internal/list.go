package internal

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoragePlugin) List(_ context.Context, req *pluginv1.StorageListRequest) (*pluginv1.StorageListResponse, error) {
	prefix := req.Prefix
	if prefix == "" {
		prefix = "."
	}

	// Root prefix with project.json pattern → list all projects.
	cleanPrefix := strings.Trim(filepath.ToSlash(filepath.Clean(prefix)), "/")
	if cleanPrefix == "." && req.Pattern == "project.json" {
		return s.listAllProjects()
	}

	entity, projectID := isListPrefix(prefix)

	switch entity {
	case entityFeature:
		return s.listFeatures(projectID)
	case entityPerson:
		return s.listPersons(projectID)
	case entityPlan:
		return s.listPlans(projectID)
	case entityRequest:
		return s.listRequests(projectID)
	case entityAssignmentRule:
		return s.listAssignmentRules(projectID)
	case entityNote:
		return s.listNotes(projectID)
	case entitySession:
		return s.listSessions()
	case entitySessionTurn:
		return s.listSessionTurns(projectID) // projectID is sessionID here
	case entityPack:
		return s.listPacks()
	case entityHookEvent:
		return s.listHookEvents()
	default:
		return s.listKV(prefix, req.Pattern)
	}
}

func (s *StoragePlugin) listAllProjects() (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT slug, '' AS project_id, version, updated_at FROM projects`)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return scanEntries(rows, func(slug, _ string) string {
		return filepath.Join(slug, "project.json")
	})
}

func (s *StoragePlugin) listFeatures(projectID string) (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT id, project_id, version, updated_at FROM features WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list features: %w", err)
	}
	return scanEntries(rows, func(id, pid string) string {
		return filepath.Join(pid, "features", id+".md")
	})
}

func (s *StoragePlugin) listPersons(projectID string) (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT id, project_id, version, updated_at FROM persons WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list persons: %w", err)
	}
	return scanEntries(rows, func(id, pid string) string {
		return filepath.Join(pid, "persons", id+".md")
	})
}

func (s *StoragePlugin) listPlans(projectID string) (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT id, project_id, version, updated_at FROM plans WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list plans: %w", err)
	}
	return scanEntries(rows, func(id, pid string) string {
		return filepath.Join(pid, "plans", id+".md")
	})
}

func (s *StoragePlugin) listRequests(projectID string) (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT id, project_id, version, updated_at FROM requests WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list requests: %w", err)
	}
	return scanEntries(rows, func(id, pid string) string {
		return filepath.Join(pid, "requests", id+".md")
	})
}

func (s *StoragePlugin) listAssignmentRules(projectID string) (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT id, project_id, version, updated_at FROM assignment_rules WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list assignment rules: %w", err)
	}
	return scanEntries(rows, func(id, pid string) string {
		return filepath.Join(pid, "assignment-rules", id+".md")
	})
}

func (s *StoragePlugin) listNotes(projectID string) (*pluginv1.StorageListResponse, error) {
	// .global notes may be stored with project_id="" or ".global" — match both.
	var rows *sql.Rows
	var err error
	if projectID == ".global" {
		rows, err = s.db.Query(`SELECT id, project_id, version, updated_at FROM notes WHERE project_id IN ('', '.global') AND deleted = 0`)
	} else {
		rows, err = s.db.Query(`SELECT id, project_id, version, updated_at FROM notes WHERE project_id = ? AND deleted = 0`, projectID)
	}
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	return scanEntries(rows, func(id, pid string) string {
		if pid == "" {
			pid = ".global"
		}
		return filepath.Join(pid, "notes", id+".md")
	})
}

func (s *StoragePlugin) listSessions() (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT id, '' AS project_id, 0 AS version, created_at AS updated_at FROM sessions`)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	return scanEntries(rows, func(id, _ string) string {
		return filepath.Join("bridge", "sessions", id+".md")
	})
}

func (s *StoragePlugin) listSessionTurns(sessionID string) (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT
		printf('turn-%03d', turn_number) AS id,
		? AS project_id,
		0 AS version,
		timestamp AS updated_at
		FROM session_turns WHERE session_id = ? ORDER BY turn_number`, sessionID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list session turns: %w", err)
	}
	return scanEntries(rows, func(turnID, sid string) string {
		return filepath.Join("bridge", "sessions", sid, turnID+".md")
	})
}

func (s *StoragePlugin) listPacks() (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT name, '' AS project_id, 0 AS version, installed_at AS updated_at FROM packs`)
	if err != nil {
		return nil, fmt.Errorf("list packs: %w", err)
	}
	return scanEntries(rows, func(name, _ string) string {
		return filepath.Join(".packs", name+".json")
	})
}

func (s *StoragePlugin) listHookEvents() (*pluginv1.StorageListResponse, error) {
	// Hook events are append-only, just return the single file entry.
	var count int64
	s.db.QueryRow(`SELECT COUNT(*) FROM hook_events`).Scan(&count)

	entries := []*pluginv1.StorageEntry{{
		Path:       ".events/hook-events.toon",
		Size:       count,
		Version:    0,
		ModifiedAt: timestamppb.Now(),
	}}
	return &pluginv1.StorageListResponse{Entries: entries}, nil
}

func (s *StoragePlugin) listKV(prefix, pattern string) (*pluginv1.StorageListResponse, error) {
	// List KV entries matching prefix.
	normalizedPrefix := strings.Trim(filepath.ToSlash(filepath.Clean(prefix)), "/")

	rows, err := s.db.Query(`SELECT path, length(content), version, updated_at
		FROM kv_store WHERE path LIKE ?`, normalizedPrefix+"%")
	if err != nil {
		return nil, fmt.Errorf("list kv: %w", err)
	}
	defer rows.Close()

	var entries []*pluginv1.StorageEntry
	for rows.Next() {
		var path string
		var size, version int64
		var updatedAt string
		if err := rows.Scan(&path, &size, &version, &updatedAt); err != nil {
			continue
		}

		// Apply pattern filter if specified.
		if pattern != "" && pattern != "*" {
			matched, _ := filepath.Match(pattern, filepath.Base(path))
			if !matched {
				continue
			}
		}

		t, _ := time.Parse(time.RFC3339, updatedAt)
		entries = append(entries, &pluginv1.StorageEntry{
			Path:       path,
			Size:       size,
			Version:    version,
			ModifiedAt: timestamppb.New(t),
		})
	}

	return &pluginv1.StorageListResponse{Entries: entries}, nil
}

// scanEntries is a helper that scans rows of (id, project_id, version, updated_at)
// and builds StorageEntry list with paths constructed by pathFunc.
func scanEntries(rows *sql.Rows, pathFunc func(id, projectID string) string) (*pluginv1.StorageListResponse, error) {
	defer rows.Close()

	var entries []*pluginv1.StorageEntry
	for rows.Next() {
		var id, projectID, updatedAt string
		var version int64
		if err := rows.Scan(&id, &projectID, &version, &updatedAt); err != nil {
			continue
		}

		t, _ := time.Parse(time.RFC3339, updatedAt)
		entries = append(entries, &pluginv1.StorageEntry{
			Path:       pathFunc(id, projectID),
			Size:       0,
			Version:    version,
			ModifiedAt: timestamppb.New(t),
		})
	}

	return &pluginv1.StorageListResponse{Entries: entries}, nil
}

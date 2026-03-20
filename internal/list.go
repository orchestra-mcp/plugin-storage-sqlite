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

	// Root prefix with *.md pattern → aggregate all .md entities across all projects.
	if cleanPrefix == "." && req.Pattern == "*.md" {
		return s.listAllEntities()
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
	case entityDelegation:
		return s.listDelegations(projectID)
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
	case entityDoc:
		return s.listDocs(projectID)
	case entitySkill:
		return s.listSkills()
	case entityAgent:
		return s.listAgents()
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

// listAllEntities aggregates all .md entity rows across all projects.
// This is called when the sync engine does List(prefix="", pattern="*.md")
// to discover all entities for push.
func (s *StoragePlugin) listAllEntities() (*pluginv1.StorageListResponse, error) {
	var allEntries []*pluginv1.StorageEntry

	// Each query selects (id, project_id, version, updated_at) and uses a path builder.
	type entityQuery struct {
		query    string
		pathFunc func(id, pid string) string
	}

	queries := []entityQuery{
		{
			query: `SELECT id, project_id, version, updated_at FROM features`,
			pathFunc: func(id, pid string) string {
				return filepath.Join(pid, "features", id+".md")
			},
		},
		{
			query: `SELECT id, project_id, version, updated_at FROM persons`,
			pathFunc: func(id, pid string) string {
				return filepath.Join(pid, "persons", id+".md")
			},
		},
		{
			query: `SELECT id, project_id, version, updated_at FROM plans`,
			pathFunc: func(id, pid string) string {
				return filepath.Join(pid, "plans", id+".md")
			},
		},
		{
			query: `SELECT id, project_id, version, updated_at FROM requests`,
			pathFunc: func(id, pid string) string {
				return filepath.Join(pid, "requests", id+".md")
			},
		},
		{
			query: `SELECT id, project_id, version, updated_at FROM delegations`,
			pathFunc: func(id, pid string) string {
				return filepath.Join(pid, "delegations", id+".md")
			},
		},
		{
			query: `SELECT id, project_id, version, updated_at FROM assignment_rules`,
			pathFunc: func(id, pid string) string {
				return filepath.Join(pid, "assignment-rules", id+".md")
			},
		},
		{
			query: `SELECT id, project_id, version, updated_at FROM notes WHERE deleted = 0`,
			pathFunc: func(id, pid string) string {
				if pid == "" {
					pid = ".global"
				}
				return filepath.Join(pid, "notes", id+".md")
			},
		},
		{
			query: `SELECT id, project_id, version, updated_at FROM docs`,
			pathFunc: func(id, pid string) string {
				return filepath.Join(pid, "docs", id+".md")
			},
		},
		{
			query: `SELECT id, '' AS project_id, version, updated_at FROM skills`,
			pathFunc: func(id, _ string) string {
				return filepath.Join(".skills", id+".md")
			},
		},
		{
			query: `SELECT id, '' AS project_id, version, updated_at FROM agents`,
			pathFunc: func(id, _ string) string {
				return filepath.Join(".agents", id+".md")
			},
		},
		{
			query: `SELECT id, '' AS project_id, 0 AS version, created_at AS updated_at FROM sessions`,
			pathFunc: func(id, _ string) string {
				return filepath.Join("bridge", "sessions", id+".md")
			},
		},
	}

	for _, eq := range queries {
		rows, err := s.db.Query(eq.query)
		if err != nil {
			continue // table may not exist yet
		}
		resp, err := scanEntries(rows, eq.pathFunc)
		if err != nil {
			continue
		}
		allEntries = append(allEntries, resp.Entries...)
	}

	return &pluginv1.StorageListResponse{Entries: allEntries}, nil
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

func (s *StoragePlugin) listDelegations(projectID string) (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT id, project_id, version, updated_at FROM delegations WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list delegations: %w", err)
	}
	return scanEntries(rows, func(id, pid string) string {
		return filepath.Join(pid, "delegations", id+".md")
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

func (s *StoragePlugin) listDocs(projectID string) (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT id, project_id, version, updated_at FROM docs WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list docs: %w", err)
	}
	return scanEntries(rows, func(id, pid string) string {
		return filepath.Join(pid, "docs", id+".md")
	})
}

func (s *StoragePlugin) listSkills() (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT id, '' AS project_id, version, updated_at FROM skills`)
	if err != nil {
		return nil, fmt.Errorf("list skills: %w", err)
	}
	return scanEntries(rows, func(id, _ string) string {
		return filepath.Join(".skills", id+".md")
	})
}

func (s *StoragePlugin) listAgents() (*pluginv1.StorageListResponse, error) {
	rows, err := s.db.Query(`SELECT id, '' AS project_id, version, updated_at FROM agents`)
	if err != nil {
		return nil, fmt.Errorf("list agents: %w", err)
	}
	return scanEntries(rows, func(id, _ string) string {
		return filepath.Join(".agents", id+".md")
	})
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

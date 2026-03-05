package internal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
)

// MigrateFromMarkdown imports existing .projects/ data into SQLite on first run.
// Idempotent: skips if tables already have data.
func MigrateFromMarkdown(db *sql.DB, workspace string) error {
	projectsBase := filepath.Join(workspace, projectsDir)
	if _, err := os.Stat(projectsBase); os.IsNotExist(err) {
		return nil // No .projects/ directory, nothing to migrate.
	}

	// Check if DB already has data (skip if so).
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM features`).Scan(&count)
	if count > 0 {
		log.Println("[storage-sqlite] migration skipped: database already has data")
		return nil
	}

	log.Printf("[storage-sqlite] migrating data from %s", projectsBase)

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin migration tx: %w", err)
	}
	defer tx.Rollback()

	migrated := 0
	skipped := 0

	err = filepath.Walk(projectsBase, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".version") {
			return nil
		}

		relPath, err := filepath.Rel(projectsBase, path)
		if err != nil {
			return nil
		}
		relPath = filepath.ToSlash(relPath)

		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("[storage-sqlite] migration: skip %s: %v", relPath, err)
			skipped++
			return nil
		}

		meta, body, err := parseMarkdownFile(data)
		if err != nil {
			log.Printf("[storage-sqlite] migration: skip %s: %v", relPath, err)
			skipped++
			return nil
		}

		route := routePath(relPath)

		switch route.Entity {
		case entityFeature:
			err = migrateFeature(tx, route, meta, body)
		case entityPerson:
			err = migratePerson(tx, route, meta, body)
		case entityPlan:
			err = migratePlan(tx, route, meta, body)
		case entityRequest:
			err = migrateRequest(tx, route, meta, body)
		case entityAssignmentRule:
			err = migrateAssignmentRule(tx, route, meta, body)
		case entityNote:
			err = migrateNote(tx, route, meta, body)
		case entityProject:
			err = migrateProject(tx, route, meta, body)
		case entityWIPConfig:
			err = migrateWIPConfig(tx, route, data)
		case entitySession:
			err = migrateSession(tx, route, meta, body)
		case entitySessionTurn:
			err = migrateSessionTurn(tx, route, meta, body)
		default:
			err = migrateKV(tx, relPath, data, meta)
		}

		if err != nil {
			log.Printf("[storage-sqlite] migration: skip %s: %v", relPath, err)
			skipped++
		} else {
			migrated++
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("walk projects: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration: %w", err)
	}

	log.Printf("[storage-sqlite] migration complete: %d migrated, %d skipped", migrated, skipped)
	return nil
}

func migrateFeature(tx *sql.Tx, r routedPath, meta *structpb.Struct, body []byte) error {
	m := metaMap(meta)
	r.ProjectID = strDef(m, "project_id", r.ProjectID)
	ensureProjectTx(tx, r.ProjectID)

	_, err := tx.Exec(`INSERT OR IGNORE INTO features (id, project_id, title, description, status, priority, kind,
		assignee, estimate, labels, depends_on, body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.EntityID, r.ProjectID, str(m, "title"), str(m, "description"),
		strDef(m, "status", "backlog"), strDef(m, "priority", "P2"), strDef(m, "kind", "feature"),
		str(m, "assignee"), str(m, "estimate"),
		jsonArr(m, "labels"), jsonArr(m, "depends_on"),
		string(body), intVal(m, "version"),
		strDef(m, "created_at", "2024-01-01T00:00:00Z"), strDef(m, "updated_at", "2024-01-01T00:00:00Z"))
	return err
}

func migratePerson(tx *sql.Tx, r routedPath, meta *structpb.Struct, body []byte) error {
	m := metaMap(meta)
	r.ProjectID = strDef(m, "project_id", r.ProjectID)
	ensureProjectTx(tx, r.ProjectID)

	_, err := tx.Exec(`INSERT OR IGNORE INTO persons (id, project_id, name, email, role, status, bio,
		github_email, integrations, labels, body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.EntityID, r.ProjectID, str(m, "name"), str(m, "email"),
		strDef(m, "role", "developer"), strDef(m, "status", "active"), str(m, "bio"),
		str(m, "github_email"), jsonMap(m, "integrations"), jsonArr(m, "labels"),
		string(body), intVal(m, "version"),
		strDef(m, "created_at", "2024-01-01T00:00:00Z"), strDef(m, "updated_at", "2024-01-01T00:00:00Z"))
	return err
}

func migratePlan(tx *sql.Tx, r routedPath, meta *structpb.Struct, body []byte) error {
	m := metaMap(meta)
	_, err := tx.Exec(`INSERT OR IGNORE INTO plans (id, project_id, title, description, status,
		features, body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.EntityID, str(m, "project_id"), str(m, "title"), str(m, "description"),
		strDef(m, "status", "draft"), jsonArr(m, "features"),
		string(body), intVal(m, "version"),
		strDef(m, "created_at", "2024-01-01T00:00:00Z"), strDef(m, "updated_at", "2024-01-01T00:00:00Z"))
	return err
}

func migrateRequest(tx *sql.Tx, r routedPath, meta *structpb.Struct, body []byte) error {
	m := metaMap(meta)
	_, err := tx.Exec(`INSERT OR IGNORE INTO requests (id, project_id, title, description, kind, status, priority,
		body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.EntityID, str(m, "project_id"), str(m, "title"), str(m, "description"),
		strDef(m, "kind", "feature"), strDef(m, "status", "pending"),
		strDef(m, "priority", "P2"), string(body), intVal(m, "version"),
		strDef(m, "created_at", "2024-01-01T00:00:00Z"), strDef(m, "updated_at", "2024-01-01T00:00:00Z"))
	return err
}

func migrateAssignmentRule(tx *sql.Tx, r routedPath, meta *structpb.Struct, body []byte) error {
	m := metaMap(meta)
	_, err := tx.Exec(`INSERT OR IGNORE INTO assignment_rules (id, project_id, kind, person_id, body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		r.EntityID, str(m, "project_id"), str(m, "kind"), str(m, "person_id"),
		string(body), intVal(m, "version"),
		strDef(m, "created_at", "2024-01-01T00:00:00Z"), strDef(m, "updated_at", "2024-01-01T00:00:00Z"))
	return err
}

func migrateNote(tx *sql.Tx, r routedPath, meta *structpb.Struct, body []byte) error {
	m := metaMap(meta)
	pinned := 0
	if boolVal(m, "pinned") {
		pinned = 1
	}
	deleted := 0
	if boolVal(m, "deleted") {
		deleted = 1
	}

	_, err := tx.Exec(`INSERT OR IGNORE INTO notes (id, project_id, title, body, pinned, deleted,
		tags, icon, color, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.EntityID, str(m, "project_id"), str(m, "title"), string(body),
		pinned, deleted, jsonArr(m, "tags"), str(m, "icon"), str(m, "color"),
		intVal(m, "version"),
		strDef(m, "created_at", "2024-01-01T00:00:00Z"), strDef(m, "updated_at", "2024-01-01T00:00:00Z"))
	return err
}

func migrateProject(tx *sql.Tx, r routedPath, meta *structpb.Struct, body []byte) error {
	m := metaMap(meta)
	_, err := tx.Exec(`INSERT OR IGNORE INTO projects (slug, name, description, metadata, body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		r.EntityID, strDef(m, "name", r.EntityID), str(m, "description"),
		jsonMap(m, "metadata"), string(body), intVal(m, "version"),
		strDef(m, "created_at", "2024-01-01T00:00:00Z"), strDef(m, "updated_at", "2024-01-01T00:00:00Z"))
	return err
}

func migrateWIPConfig(tx *sql.Tx, r routedPath, data []byte) error {
	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}
	maxIP := intVal(config, "max_in_progress")
	_, err := tx.Exec(`INSERT OR IGNORE INTO wip_config (project_id, max_in_progress, version, updated_at)
		VALUES (?, ?, 1, datetime('now'))`, r.EntityID, maxIP)
	return err
}

func migrateSession(tx *sql.Tx, r routedPath, meta *structpb.Struct, body []byte) error {
	m := metaMap(meta)
	_, err := tx.Exec(`INSERT OR IGNORE INTO sessions (id, account_id, name, workspace, model, permission_mode,
		allowed_tools, max_budget, system_prompt, status, message_count,
		total_tokens_in, total_tokens_out, total_cost_usd, claude_session_id,
		last_message_at, body, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.EntityID, str(m, "account_id"), str(m, "name"), str(m, "workspace"),
		str(m, "model"), strDef(m, "permission_mode", "plan"),
		jsonArr(m, "allowed_tools"), floatVal(m, "max_budget"),
		str(m, "system_prompt"), strDef(m, "status", "active"),
		intVal(m, "message_count"), intVal(m, "total_tokens_in"),
		intVal(m, "total_tokens_out"), floatVal(m, "total_cost_usd"),
		str(m, "claude_session_id"), str(m, "last_message_at"),
		string(body), strDef(m, "created_at", "2024-01-01T00:00:00Z"))
	return err
}

func migrateSessionTurn(tx *sql.Tx, r routedPath, meta *structpb.Struct, body []byte) error {
	m := metaMap(meta)
	turnNum := intVal(m, "number")
	if turnNum == 0 {
		turnNum = int64(extractTurnNumber(r.SubID))
	}

	_, err := tx.Exec(`INSERT OR IGNORE INTO session_turns (session_id, turn_number, user_prompt, response,
		tokens_in, tokens_out, cost_usd, model, duration_ms, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.EntityID, turnNum, str(m, "user_prompt"), string(body),
		intVal(m, "tokens_in"), intVal(m, "tokens_out"),
		floatVal(m, "cost_usd"), str(m, "model"),
		intVal(m, "duration_ms"), strDef(m, "timestamp", "2024-01-01T00:00:00Z"))
	return err
}

func migrateKV(tx *sql.Tx, path string, data []byte, meta *structpb.Struct) error {
	metaJSON := "{}"
	if meta != nil {
		if b, err := json.Marshal(meta.AsMap()); err == nil {
			metaJSON = string(b)
		}
	}

	_, err := tx.Exec(`INSERT OR IGNORE INTO kv_store (path, content, metadata, version, created_at, updated_at)
		VALUES (?, ?, ?, 1, datetime('now'), datetime('now'))`,
		path, data, metaJSON)
	return err
}

func ensureProjectTx(tx *sql.Tx, slug string) {
	if slug == "" {
		return
	}
	tx.Exec(`INSERT OR IGNORE INTO projects (slug, name, created_at, updated_at)
		VALUES (?, ?, datetime('now'), datetime('now'))`, slug, slug)
}

// parseMarkdownFile parses YAML frontmatter + body from a markdown file.
// Duplicated from plugin-storage-markdown to avoid import cycle.
func parseMarkdownFile(data []byte) (*structpb.Struct, []byte, error) {
	delim := []byte("---")

	if !bytes.HasPrefix(data, delim) {
		return nil, data, nil
	}

	firstNewline := bytes.IndexByte(data, '\n')
	if firstNewline == -1 {
		return nil, data, nil
	}

	firstLine := bytes.TrimRight(data[:firstNewline], "\r")
	if !bytes.Equal(firstLine, delim) {
		return nil, data, nil
	}

	rest := data[firstNewline+1:]
	closingIdx := -1
	offset := 0
	for offset < len(rest) {
		lineEnd := bytes.IndexByte(rest[offset:], '\n')
		var line []byte
		if lineEnd == -1 {
			line = rest[offset:]
		} else {
			line = rest[offset : offset+lineEnd]
		}
		if bytes.Equal(bytes.TrimRight(line, "\r"), delim) {
			closingIdx = offset
			break
		}
		if lineEnd == -1 {
			break
		}
		offset += lineEnd + 1
	}

	if closingIdx == -1 {
		return nil, nil, fmt.Errorf("missing closing ---")
	}

	yamlData := rest[:closingIdx]
	var m map[string]any
	if err := yaml.Unmarshal(yamlData, &m); err != nil {
		return nil, nil, fmt.Errorf("parse YAML: %w", err)
	}

	meta, err := structpb.NewStruct(m)
	if err != nil {
		return nil, nil, fmt.Errorf("convert to structpb: %w", err)
	}

	bodyStart := closingIdx + len(delim)
	if bodyStart < len(rest) && rest[bodyStart] == '\r' {
		bodyStart++
	}
	if bodyStart < len(rest) && rest[bodyStart] == '\n' {
		bodyStart++
	}
	if bodyStart < len(rest) && rest[bodyStart] == '\n' {
		bodyStart++
	} else if bodyStart+1 < len(rest) && rest[bodyStart] == '\r' && rest[bodyStart+1] == '\n' {
		bodyStart += 2
	}

	return meta, rest[bodyStart:], nil
}

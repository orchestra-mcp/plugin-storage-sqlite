package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *StoragePlugin) Write(ctx context.Context, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	route := routePath(req.Path)

	s.mu.Lock()
	defer s.mu.Unlock()

	var resp *pluginv1.StorageWriteResponse
	var err error

	switch route.Entity {
	case entityFeature:
		resp, err = s.writeFeature(route, req)
	case entityPerson:
		resp, err = s.writePerson(route, req)
	case entityPlan:
		resp, err = s.writePlan(route, req)
	case entityRequest:
		resp, err = s.writeRequest(route, req)
	case entityAssignmentRule:
		resp, err = s.writeAssignmentRule(route, req)
	case entityNote:
		resp, err = s.writeNote(route, req)
	case entityProject:
		resp, err = s.writeProject(route, req)
	case entityWIPConfig:
		resp, err = s.writeWIPConfig(route, req)
	case entitySession:
		resp, err = s.writeSession(route, req)
	case entitySessionTurn:
		resp, err = s.writeSessionTurn(route, req)
	case entityPack:
		resp, err = s.writePacks(route, req)
	case entityStack:
		resp, err = s.writeStacks(route, req)
	default:
		resp, err = s.writeKV(route, req)
	}

	if err != nil {
		return &pluginv1.StorageWriteResponse{Success: false, Error: err.Error()}, nil
	}

	// Dual-write: async markdown export for git visibility.
	if resp.Success {
		go s.exportMarkdown(req.Path, req.Metadata, req.Content)
	}

	return resp, nil
}

func (s *StoragePlugin) writeFeature(r routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)
	now := time.Now().UTC().Format(time.RFC3339)

	newVer, err := s.casCheck("features", "id", r.EntityID, req.ExpectedVersion)
	if err != nil {
		return &pluginv1.StorageWriteResponse{Success: false, Error: err.Error()}, nil
	}

	_, err = s.db.Exec(`INSERT INTO features (id, project_id, title, description, status, priority, kind,
		assignee, estimate, labels, depends_on, body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title=excluded.title, description=excluded.description, status=excluded.status,
			priority=excluded.priority, kind=excluded.kind, assignee=excluded.assignee,
			estimate=excluded.estimate, labels=excluded.labels, depends_on=excluded.depends_on,
			body=excluded.body, version=excluded.version, updated_at=excluded.updated_at`,
		r.EntityID, str(m, "project_id"), str(m, "title"), str(m, "description"),
		str(m, "status"), strDef(m, "priority", "P2"), strDef(m, "kind", "feature"),
		str(m, "assignee"), str(m, "estimate"),
		jsonArr(m, "labels"), jsonArr(m, "depends_on"),
		string(req.Content), newVer, strDef(m, "created_at", now), now)
	if err != nil {
		return nil, fmt.Errorf("write feature: %w", err)
	}

	// Ensure project exists.
	s.ensureProject(r.ProjectID)

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: newVer}, nil
}

func (s *StoragePlugin) writePerson(r routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)
	now := time.Now().UTC().Format(time.RFC3339)

	newVer, err := s.casCheck("persons", "id", r.EntityID, req.ExpectedVersion)
	if err != nil {
		return &pluginv1.StorageWriteResponse{Success: false, Error: err.Error()}, nil
	}

	_, err = s.db.Exec(`INSERT INTO persons (id, project_id, name, email, role, status, bio,
		github_email, integrations, labels, body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name=excluded.name, email=excluded.email, role=excluded.role,
			status=excluded.status, bio=excluded.bio, github_email=excluded.github_email,
			integrations=excluded.integrations, labels=excluded.labels,
			body=excluded.body, version=excluded.version, updated_at=excluded.updated_at`,
		r.EntityID, str(m, "project_id"), str(m, "name"), str(m, "email"),
		strDef(m, "role", "developer"), strDef(m, "status", "active"), str(m, "bio"),
		str(m, "github_email"), jsonMap(m, "integrations"), jsonArr(m, "labels"),
		string(req.Content), newVer, strDef(m, "created_at", now), now)
	if err != nil {
		return nil, fmt.Errorf("write person: %w", err)
	}

	s.ensureProject(r.ProjectID)
	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: newVer}, nil
}

func (s *StoragePlugin) writePlan(r routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)
	now := time.Now().UTC().Format(time.RFC3339)

	newVer, err := s.casCheck("plans", "id", r.EntityID, req.ExpectedVersion)
	if err != nil {
		return &pluginv1.StorageWriteResponse{Success: false, Error: err.Error()}, nil
	}

	_, err = s.db.Exec(`INSERT INTO plans (id, project_id, title, description, status,
		features, body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title=excluded.title, description=excluded.description, status=excluded.status,
			features=excluded.features, body=excluded.body, version=excluded.version,
			updated_at=excluded.updated_at`,
		r.EntityID, str(m, "project_id"), str(m, "title"), str(m, "description"),
		strDef(m, "status", "draft"), jsonArr(m, "features"),
		string(req.Content), newVer, strDef(m, "created_at", now), now)
	if err != nil {
		return nil, fmt.Errorf("write plan: %w", err)
	}

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: newVer}, nil
}

func (s *StoragePlugin) writeRequest(r routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)
	now := time.Now().UTC().Format(time.RFC3339)

	newVer, err := s.casCheck("requests", "id", r.EntityID, req.ExpectedVersion)
	if err != nil {
		return &pluginv1.StorageWriteResponse{Success: false, Error: err.Error()}, nil
	}

	_, err = s.db.Exec(`INSERT INTO requests (id, project_id, title, description, kind, status, priority,
		body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title=excluded.title, description=excluded.description, kind=excluded.kind,
			status=excluded.status, priority=excluded.priority, body=excluded.body,
			version=excluded.version, updated_at=excluded.updated_at`,
		r.EntityID, str(m, "project_id"), str(m, "title"), str(m, "description"),
		strDef(m, "kind", "feature"), strDef(m, "status", "pending"),
		strDef(m, "priority", "P2"), string(req.Content), newVer,
		strDef(m, "created_at", now), now)
	if err != nil {
		return nil, fmt.Errorf("write request: %w", err)
	}

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: newVer}, nil
}

func (s *StoragePlugin) writeAssignmentRule(r routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)
	now := time.Now().UTC().Format(time.RFC3339)

	newVer, err := s.casCheck("assignment_rules", "id", r.EntityID, req.ExpectedVersion)
	if err != nil {
		return &pluginv1.StorageWriteResponse{Success: false, Error: err.Error()}, nil
	}

	_, err = s.db.Exec(`INSERT INTO assignment_rules (id, project_id, kind, person_id, body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			kind=excluded.kind, person_id=excluded.person_id, body=excluded.body,
			version=excluded.version, updated_at=excluded.updated_at`,
		r.EntityID, str(m, "project_id"), str(m, "kind"), str(m, "person_id"),
		string(req.Content), newVer, strDef(m, "created_at", now), now)
	if err != nil {
		return nil, fmt.Errorf("write assignment rule: %w", err)
	}

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: newVer}, nil
}

func (s *StoragePlugin) writeNote(r routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)
	now := time.Now().UTC().Format(time.RFC3339)

	newVer, err := s.casCheck("notes", "id", r.EntityID, req.ExpectedVersion)
	if err != nil {
		return &pluginv1.StorageWriteResponse{Success: false, Error: err.Error()}, nil
	}

	pinned := 0
	if boolVal(m, "pinned") {
		pinned = 1
	}
	deleted := 0
	if boolVal(m, "deleted") {
		deleted = 1
	}

	_, err = s.db.Exec(`INSERT INTO notes (id, project_id, title, body, pinned, deleted,
		tags, icon, color, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title=excluded.title, body=excluded.body, pinned=excluded.pinned,
			deleted=excluded.deleted, tags=excluded.tags, icon=excluded.icon,
			color=excluded.color, version=excluded.version, updated_at=excluded.updated_at`,
		r.EntityID, str(m, "project_id"), str(m, "title"), string(req.Content),
		pinned, deleted, jsonArr(m, "tags"), str(m, "icon"), str(m, "color"),
		newVer, strDef(m, "created_at", now), now)
	if err != nil {
		return nil, fmt.Errorf("write note: %w", err)
	}

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: newVer}, nil
}

func (s *StoragePlugin) writeProject(r routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)
	now := time.Now().UTC().Format(time.RFC3339)

	newVer, err := s.casCheck("projects", "slug", r.EntityID, req.ExpectedVersion)
	if err != nil {
		return &pluginv1.StorageWriteResponse{Success: false, Error: err.Error()}, nil
	}

	_, err = s.db.Exec(`INSERT INTO projects (slug, name, description, metadata, body, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			name=excluded.name, description=excluded.description, metadata=excluded.metadata,
			body=excluded.body, version=excluded.version, updated_at=excluded.updated_at`,
		r.EntityID, strDef(m, "name", r.EntityID), str(m, "description"),
		jsonMap(m, "metadata"), string(req.Content), newVer,
		strDef(m, "created_at", now), now)
	if err != nil {
		return nil, fmt.Errorf("write project: %w", err)
	}

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: newVer}, nil
}

func (s *StoragePlugin) writeWIPConfig(r routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)
	now := time.Now().UTC().Format(time.RFC3339)

	newVer, err := s.casCheck("wip_config", "project_id", r.EntityID, req.ExpectedVersion)
	if err != nil {
		return &pluginv1.StorageWriteResponse{Success: false, Error: err.Error()}, nil
	}

	_, err = s.db.Exec(`INSERT INTO wip_config (project_id, max_in_progress, version, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(project_id) DO UPDATE SET
			max_in_progress=excluded.max_in_progress, version=excluded.version,
			updated_at=excluded.updated_at`,
		r.EntityID, intVal(m, "max_in_progress"), newVer, now)
	if err != nil {
		return nil, fmt.Errorf("write wip config: %w", err)
	}

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: newVer}, nil
}

func (s *StoragePlugin) writeSession(r routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := s.db.Exec(`INSERT INTO sessions (id, account_id, name, workspace, model, permission_mode,
		allowed_tools, max_budget, system_prompt, status, message_count,
		total_tokens_in, total_tokens_out, total_cost_usd, claude_session_id,
		last_message_at, body, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name=excluded.name, workspace=excluded.workspace, model=excluded.model,
			permission_mode=excluded.permission_mode, allowed_tools=excluded.allowed_tools,
			max_budget=excluded.max_budget, system_prompt=excluded.system_prompt,
			status=excluded.status, message_count=excluded.message_count,
			total_tokens_in=excluded.total_tokens_in, total_tokens_out=excluded.total_tokens_out,
			total_cost_usd=excluded.total_cost_usd, claude_session_id=excluded.claude_session_id,
			last_message_at=excluded.last_message_at, body=excluded.body`,
		r.EntityID, str(m, "account_id"), str(m, "name"), str(m, "workspace"),
		str(m, "model"), strDef(m, "permission_mode", "plan"),
		jsonArr(m, "allowed_tools"), floatVal(m, "max_budget"),
		str(m, "system_prompt"), strDef(m, "status", "active"),
		intVal(m, "message_count"), intVal(m, "total_tokens_in"),
		intVal(m, "total_tokens_out"), floatVal(m, "total_cost_usd"),
		str(m, "claude_session_id"), strDef(m, "last_message_at", now),
		string(req.Content), strDef(m, "created_at", now))
	if err != nil {
		return nil, fmt.Errorf("write session: %w", err)
	}

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: 1}, nil
}

func (s *StoragePlugin) writeSessionTurn(r routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)
	now := time.Now().UTC().Format(time.RFC3339)

	turnNum := intVal(m, "number")
	if turnNum == 0 {
		turnNum = int64(extractTurnNumber(r.SubID))
	}

	_, err := s.db.Exec(`INSERT INTO session_turns (session_id, turn_number, user_prompt, response,
		tokens_in, tokens_out, cost_usd, model, duration_ms, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(session_id, turn_number) DO UPDATE SET
			user_prompt=excluded.user_prompt, response=excluded.response,
			tokens_in=excluded.tokens_in, tokens_out=excluded.tokens_out,
			cost_usd=excluded.cost_usd, model=excluded.model,
			duration_ms=excluded.duration_ms`,
		r.EntityID, turnNum, str(m, "user_prompt"), string(req.Content),
		intVal(m, "tokens_in"), intVal(m, "tokens_out"),
		floatVal(m, "cost_usd"), str(m, "model"),
		intVal(m, "duration_ms"), strDef(m, "timestamp", now))
	if err != nil {
		return nil, fmt.Errorf("write session turn: %w", err)
	}

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: 1}, nil
}

func (s *StoragePlugin) writePacks(_ routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)
	now := time.Now().UTC().Format(time.RFC3339)

	// Packs registry is typically the full set written at once.
	// Handle single pack upsert if name is present, or KV fallback.
	name := str(m, "name")
	if name == "" {
		return s.writeKV(routedPath{Raw: ".packs/registry.json"}, req)
	}

	_, err := s.db.Exec(`INSERT INTO packs (name, version, repo, installed_at, metadata)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			version=excluded.version, repo=excluded.repo, metadata=excluded.metadata`,
		name, str(m, "version"), str(m, "repo"),
		strDef(m, "installed_at", now), jsonMap(m, "metadata"))
	if err != nil {
		return nil, fmt.Errorf("write pack: %w", err)
	}

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: 1}, nil
}

func (s *StoragePlugin) writeStacks(_ routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	m := metaMap(req.Metadata)

	_, err := s.db.Exec(`INSERT INTO stacks (id, stacks, version) VALUES (1, ?, ?)
		ON CONFLICT(id) DO UPDATE SET stacks=excluded.stacks, version=excluded.version`,
		jsonArr(m, "stacks"), intVal(m, "version")+1)
	if err != nil {
		return nil, fmt.Errorf("write stacks: %w", err)
	}

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: intVal(m, "version") + 1}, nil
}

func (s *StoragePlugin) writeKV(r routedPath, req *pluginv1.StorageWriteRequest) (*pluginv1.StorageWriteResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	metaJSON := "{}"
	if req.Metadata != nil {
		if b, err := json.Marshal(req.Metadata.AsMap()); err == nil {
			metaJSON = string(b)
		}
	}

	var currentVer int64
	_ = s.db.QueryRow(`SELECT version FROM kv_store WHERE path = ?`, r.Raw).Scan(&currentVer)

	if req.ExpectedVersion > 0 && currentVer != req.ExpectedVersion {
		return &pluginv1.StorageWriteResponse{
			Success: false,
			Error:   fmt.Sprintf("version conflict: expected %d, current %d", req.ExpectedVersion, currentVer),
		}, nil
	}
	if req.ExpectedVersion == 0 && currentVer > 0 {
		return &pluginv1.StorageWriteResponse{
			Success: false,
			Error:   "entry already exists (expected_version=0 means create new)",
		}, nil
	}

	newVer := currentVer + 1
	_, err := s.db.Exec(`INSERT INTO kv_store (path, content, metadata, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			content=excluded.content, metadata=excluded.metadata,
			version=excluded.version, updated_at=excluded.updated_at`,
		r.Raw, req.Content, metaJSON, newVer, now, now)
	if err != nil {
		return nil, fmt.Errorf("write kv: %w", err)
	}

	return &pluginv1.StorageWriteResponse{Success: true, NewVersion: newVer}, nil
}

// --- CAS helpers ---

// casCheck implements compare-and-swap versioning for typed entities.
// Returns the new version number or an error if the check fails.
func (s *StoragePlugin) casCheck(table, pkCol, pkVal string, expectedVersion int64) (int64, error) {
	var currentVer int64
	err := s.db.QueryRow(
		fmt.Sprintf(`SELECT version FROM %s WHERE %s = ?`, table, pkCol), pkVal,
	).Scan(&currentVer)

	exists := err == nil

	if expectedVersion == -1 {
		// Upsert: write unconditionally.
		if exists {
			return currentVer + 1, nil
		}
		return 1, nil
	}

	if expectedVersion == 0 {
		// Create new: fail if exists.
		if exists {
			return 0, fmt.Errorf("already exists (expected_version=0 means create new)")
		}
		return 1, nil
	}

	// Update: expected version must match.
	if !exists {
		return 0, fmt.Errorf("not found")
	}
	if currentVer != expectedVersion {
		return 0, fmt.Errorf("version conflict: expected %d, current %d", expectedVersion, currentVer)
	}
	return currentVer + 1, nil
}

// ensureProject creates a minimal project row if it doesn't exist.
func (s *StoragePlugin) ensureProject(slug string) {
	if slug == "" {
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	s.db.Exec(`INSERT OR IGNORE INTO projects (slug, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)`, slug, slug, now, now)
}

// --- metadata extraction helpers ---

func metaMap(meta *structpb.Struct) map[string]any {
	if meta == nil {
		return map[string]any{}
	}
	return meta.AsMap()
}

func str(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func strDef(m map[string]any, key, def string) string {
	s := str(m, key)
	if s == "" {
		return def
	}
	return s
}

func intVal(m map[string]any, key string) int64 {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return int64(n)
		case int64:
			return n
		case int:
			return int64(n)
		}
	}
	return 0
}

func floatVal(m map[string]any, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

func boolVal(m map[string]any, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func jsonArr(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if arr, ok := v.([]any); ok && len(arr) > 0 {
			b, _ := json.Marshal(arr)
			return string(b)
		}
	}
	return "[]"
}

func jsonMap(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if obj, ok := v.(map[string]any); ok && len(obj) > 0 {
			b, _ := json.Marshal(obj)
			return string(b)
		}
	}
	return "{}"
}

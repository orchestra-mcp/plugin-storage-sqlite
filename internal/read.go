package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *StoragePlugin) Read(_ context.Context, req *pluginv1.StorageReadRequest) (*pluginv1.StorageReadResponse, error) {
	route := routePath(req.Path)

	switch route.Entity {
	case entityFeature:
		return s.readFeature(route)
	case entityPerson:
		return s.readPerson(route)
	case entityPlan:
		return s.readPlan(route)
	case entityRequest:
		return s.readRequest(route)
	case entityAssignmentRule:
		return s.readAssignmentRule(route)
	case entityNote:
		return s.readNote(route)
	case entityProject:
		return s.readProject(route)
	case entityWIPConfig:
		return s.readWIPConfig(route)
	case entitySession:
		return s.readSession(route)
	case entitySessionTurn:
		return s.readSessionTurn(route)
	case entityPack:
		return s.readPacks(route)
	case entityStack:
		return s.readStacks(route)
	case entityHookEvent:
		return s.readKV(route)
	default:
		return s.readKV(route)
	}
}

func (s *StoragePlugin) readFeature(r routedPath) (*pluginv1.StorageReadResponse, error) {
	var (
		id, projectID, title, description, status, priority, kind string
		assignee, estimate, labelsJSON, dependsOnJSON, body       string
		version                                                   int64
		createdAt, updatedAt                                      string
	)
	err := s.db.QueryRow(`SELECT id, project_id, title, description, status, priority, kind,
		assignee, estimate, labels, depends_on, body, version, created_at, updated_at
		FROM features WHERE id = ?`, r.EntityID).Scan(
		&id, &projectID, &title, &description, &status, &priority, &kind,
		&assignee, &estimate, &labelsJSON, &dependsOnJSON, &body,
		&version, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("read feature %s: %w", r.EntityID, err)
	}

	m := map[string]any{
		"id": id, "project_id": projectID, "title": title,
		"description": description, "status": status, "priority": priority,
		"version": float64(version), "created_at": createdAt, "updated_at": updatedAt,
	}
	if kind != "" {
		m["kind"] = kind
	}
	if assignee != "" {
		m["assignee"] = assignee
	}
	if estimate != "" {
		m["estimate"] = estimate
	}
	setJSONArray(m, "labels", labelsJSON)
	setJSONArray(m, "depends_on", dependsOnJSON)

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Content: []byte(body), Metadata: meta, Version: version}, nil
}

func (s *StoragePlugin) readPerson(r routedPath) (*pluginv1.StorageReadResponse, error) {
	var (
		id, projectID, name, email, role, status, bio    string
		githubEmail, integrationsJSON, labelsJSON, body   string
		version                                           int64
		createdAt, updatedAt                              string
	)
	err := s.db.QueryRow(`SELECT id, project_id, name, email, role, status, bio,
		github_email, integrations, labels, body, version, created_at, updated_at
		FROM persons WHERE id = ?`, r.EntityID).Scan(
		&id, &projectID, &name, &email, &role, &status, &bio,
		&githubEmail, &integrationsJSON, &labelsJSON, &body,
		&version, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("read person %s: %w", r.EntityID, err)
	}

	m := map[string]any{
		"id": id, "project_id": projectID, "name": name,
		"role": role, "status": status, "version": float64(version),
		"created_at": createdAt, "updated_at": updatedAt,
	}
	if email != "" {
		m["email"] = email
	}
	if bio != "" {
		m["bio"] = bio
	}
	if githubEmail != "" {
		m["github_email"] = githubEmail
	}
	setJSONMap(m, "integrations", integrationsJSON)
	setJSONArray(m, "labels", labelsJSON)

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Content: []byte(body), Metadata: meta, Version: version}, nil
}

func (s *StoragePlugin) readPlan(r routedPath) (*pluginv1.StorageReadResponse, error) {
	var (
		id, projectID, title, description, status string
		featuresJSON, body                        string
		version                                   int64
		createdAt, updatedAt                      string
	)
	err := s.db.QueryRow(`SELECT id, project_id, title, description, status,
		features, body, version, created_at, updated_at
		FROM plans WHERE id = ?`, r.EntityID).Scan(
		&id, &projectID, &title, &description, &status,
		&featuresJSON, &body, &version, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("read plan %s: %w", r.EntityID, err)
	}

	m := map[string]any{
		"id": id, "project_id": projectID, "title": title,
		"description": description, "status": status,
		"version": float64(version), "created_at": createdAt, "updated_at": updatedAt,
	}
	setJSONArray(m, "features", featuresJSON)

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Content: []byte(body), Metadata: meta, Version: version}, nil
}

func (s *StoragePlugin) readRequest(r routedPath) (*pluginv1.StorageReadResponse, error) {
	var (
		id, projectID, title, description, kind, status, priority string
		body                                                      string
		version                                                   int64
		createdAt, updatedAt                                      string
	)
	err := s.db.QueryRow(`SELECT id, project_id, title, description, kind, status, priority,
		body, version, created_at, updated_at
		FROM requests WHERE id = ?`, r.EntityID).Scan(
		&id, &projectID, &title, &description, &kind, &status, &priority,
		&body, &version, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("read request %s: %w", r.EntityID, err)
	}

	m := map[string]any{
		"id": id, "project_id": projectID, "title": title,
		"description": description, "kind": kind, "status": status,
		"priority": priority, "version": float64(version),
		"created_at": createdAt, "updated_at": updatedAt,
	}

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Content: []byte(body), Metadata: meta, Version: version}, nil
}

func (s *StoragePlugin) readAssignmentRule(r routedPath) (*pluginv1.StorageReadResponse, error) {
	var (
		id, projectID, kind, personID, body string
		version                             int64
		createdAt, updatedAt                string
	)
	err := s.db.QueryRow(`SELECT id, project_id, kind, person_id, body, version,
		created_at, updated_at FROM assignment_rules WHERE id = ?`, r.EntityID).Scan(
		&id, &projectID, &kind, &personID, &body, &version, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("read assignment rule %s: %w", r.EntityID, err)
	}

	m := map[string]any{
		"id": id, "project_id": projectID, "kind": kind,
		"person_id": personID, "version": float64(version),
		"created_at": createdAt, "updated_at": updatedAt,
	}

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Content: []byte(body), Metadata: meta, Version: version}, nil
}

func (s *StoragePlugin) readNote(r routedPath) (*pluginv1.StorageReadResponse, error) {
	var (
		id, projectID, title, body      string
		pinned, deleted                 int
		tagsJSON, icon, color           string
		version                         int64
		createdAt, updatedAt            string
	)
	err := s.db.QueryRow(`SELECT id, project_id, title, body, pinned, deleted,
		tags, icon, color, version, created_at, updated_at
		FROM notes WHERE id = ?`, r.EntityID).Scan(
		&id, &projectID, &title, &body, &pinned, &deleted,
		&tagsJSON, &icon, &color, &version, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("read note %s: %w", r.EntityID, err)
	}

	m := map[string]any{
		"id": id, "project_id": projectID, "title": title,
		"pinned": pinned != 0, "deleted": deleted != 0,
		"version": float64(version), "created_at": createdAt, "updated_at": updatedAt,
	}
	if icon != "" {
		m["icon"] = icon
	}
	if color != "" {
		m["color"] = color
	}
	setJSONArray(m, "tags", tagsJSON)

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Content: []byte(body), Metadata: meta, Version: version}, nil
}

func (s *StoragePlugin) readProject(r routedPath) (*pluginv1.StorageReadResponse, error) {
	var (
		slug, name, description, metadataJSON, body string
		version                                     int64
		createdAt, updatedAt                        string
	)
	err := s.db.QueryRow(`SELECT slug, name, description, metadata, body, version,
		created_at, updated_at FROM projects WHERE slug = ?`, r.EntityID).Scan(
		&slug, &name, &description, &metadataJSON, &body, &version, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("read project %s: %w", r.EntityID, err)
	}

	m := map[string]any{
		"slug": slug, "name": name, "description": description,
		"version": float64(version), "created_at": createdAt, "updated_at": updatedAt,
	}
	setJSONMap(m, "metadata", metadataJSON)

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Content: []byte(body), Metadata: meta, Version: version}, nil
}

func (s *StoragePlugin) readWIPConfig(r routedPath) (*pluginv1.StorageReadResponse, error) {
	var (
		projectID   string
		maxIP       int
		version     int64
		updatedAt   string
	)
	err := s.db.QueryRow(`SELECT project_id, max_in_progress, version, updated_at
		FROM wip_config WHERE project_id = ?`, r.EntityID).Scan(
		&projectID, &maxIP, &version, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			m := map[string]any{"project_id": r.EntityID, "max_in_progress": float64(0), "version": float64(0)}
			meta, _ := structpb.NewStruct(m)
			return &pluginv1.StorageReadResponse{Metadata: meta, Version: 0}, nil
		}
		return nil, fmt.Errorf("read wip config %s: %w", r.EntityID, err)
	}

	m := map[string]any{
		"project_id":      projectID,
		"max_in_progress": float64(maxIP),
		"version":         float64(version),
		"updated_at":      updatedAt,
	}

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Metadata: meta, Version: version}, nil
}

func (s *StoragePlugin) readSession(r routedPath) (*pluginv1.StorageReadResponse, error) {
	var (
		id, accountID, name, workspace, model, permMode string
		allowedToolsJSON, systemPrompt, status, body    string
		claudeSessionID, lastMessageAt, createdAt       string
		maxBudget, totalCost                            float64
		msgCount, tokensIn, tokensOut                   int64
	)
	err := s.db.QueryRow(`SELECT id, account_id, name, workspace, model, permission_mode,
		allowed_tools, max_budget, system_prompt, status, message_count,
		total_tokens_in, total_tokens_out, total_cost_usd, claude_session_id,
		last_message_at, body, created_at
		FROM sessions WHERE id = ?`, r.EntityID).Scan(
		&id, &accountID, &name, &workspace, &model, &permMode,
		&allowedToolsJSON, &maxBudget, &systemPrompt, &status, &msgCount,
		&tokensIn, &tokensOut, &totalCost, &claudeSessionID,
		&lastMessageAt, &body, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("read session %s: %w", r.EntityID, err)
	}

	m := map[string]any{
		"id": id, "account_id": accountID, "name": name,
		"workspace": workspace, "model": model, "permission_mode": permMode,
		"max_budget": maxBudget, "system_prompt": systemPrompt,
		"status": status, "message_count": float64(msgCount),
		"total_tokens_in": float64(tokensIn), "total_tokens_out": float64(tokensOut),
		"total_cost_usd": totalCost, "created_at": createdAt,
		"last_message_at": lastMessageAt,
	}
	if claudeSessionID != "" {
		m["claude_session_id"] = claudeSessionID
	}
	setJSONArray(m, "allowed_tools", allowedToolsJSON)

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Content: []byte(body), Metadata: meta, Version: 0}, nil
}

func (s *StoragePlugin) readSessionTurn(r routedPath) (*pluginv1.StorageReadResponse, error) {
	// SubID is "turn-NNN" — extract number.
	turnNum := extractTurnNumber(r.SubID)

	var (
		userPrompt, response, model, timestamp string
		tokensIn, tokensOut, durationMS        int64
		costUSD                                float64
	)
	err := s.db.QueryRow(`SELECT user_prompt, response, tokens_in, tokens_out,
		cost_usd, model, duration_ms, timestamp
		FROM session_turns WHERE session_id = ? AND turn_number = ?`,
		r.EntityID, turnNum).Scan(
		&userPrompt, &response, &tokensIn, &tokensOut,
		&costUSD, &model, &durationMS, &timestamp)
	if err != nil {
		return nil, fmt.Errorf("read turn %s/%s: %w", r.EntityID, r.SubID, err)
	}

	m := map[string]any{
		"number":      float64(turnNum),
		"user_prompt": userPrompt,
		"tokens_in":   float64(tokensIn),
		"tokens_out":  float64(tokensOut),
		"cost_usd":    costUSD,
		"model":       model,
		"duration_ms": float64(durationMS),
		"timestamp":   timestamp,
	}

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Content: []byte(response), Metadata: meta, Version: 0}, nil
}

func (s *StoragePlugin) readPacks(_ routedPath) (*pluginv1.StorageReadResponse, error) {
	rows, err := s.db.Query(`SELECT name, version, repo, installed_at, metadata FROM packs`)
	if err != nil {
		return nil, fmt.Errorf("read packs: %w", err)
	}
	defer rows.Close()

	var packs []any
	for rows.Next() {
		var name, version, repo, installedAt, metaJSON string
		if err := rows.Scan(&name, &version, &repo, &installedAt, &metaJSON); err != nil {
			continue
		}
		p := map[string]any{
			"name": name, "version": version, "repo": repo,
			"installed_at": installedAt,
		}
		var meta map[string]any
		if json.Unmarshal([]byte(metaJSON), &meta) == nil {
			p["metadata"] = meta
		}
		packs = append(packs, p)
	}

	m := map[string]any{"packs": packs}
	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Metadata: meta, Version: 0}, nil
}

func (s *StoragePlugin) readStacks(_ routedPath) (*pluginv1.StorageReadResponse, error) {
	var stacksJSON string
	var version int64
	err := s.db.QueryRow(`SELECT stacks, version FROM stacks WHERE id = 1`).Scan(&stacksJSON, &version)
	if err != nil {
		if err == sql.ErrNoRows {
			m := map[string]any{"stacks": []any{}, "version": float64(0)}
			meta, _ := structpb.NewStruct(m)
			return &pluginv1.StorageReadResponse{Metadata: meta, Version: 0}, nil
		}
		return nil, fmt.Errorf("read stacks: %w", err)
	}

	m := map[string]any{"version": float64(version)}
	setJSONArray(m, "stacks", stacksJSON)

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Metadata: meta, Version: version}, nil
}

func (s *StoragePlugin) readKV(r routedPath) (*pluginv1.StorageReadResponse, error) {
	var content []byte
	var metaJSON string
	var version int64
	err := s.db.QueryRow(`SELECT content, metadata, version FROM kv_store WHERE path = ?`,
		r.Raw).Scan(&content, &metaJSON, &version)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", r.Raw, err)
	}

	m := make(map[string]any)
	_ = json.Unmarshal([]byte(metaJSON), &m)
	m["version"] = float64(version)

	meta, _ := structpb.NewStruct(m)
	return &pluginv1.StorageReadResponse{Content: content, Metadata: meta, Version: version}, nil
}

// --- helpers ---

func setJSONArray(m map[string]any, key, jsonStr string) {
	if jsonStr == "" || jsonStr == "[]" {
		return
	}
	var arr []any
	if json.Unmarshal([]byte(jsonStr), &arr) == nil && len(arr) > 0 {
		m[key] = arr
	}
}

func setJSONMap(m map[string]any, key, jsonStr string) {
	if jsonStr == "" || jsonStr == "{}" {
		return
	}
	var obj map[string]any
	if json.Unmarshal([]byte(jsonStr), &obj) == nil && len(obj) > 0 {
		m[key] = obj
	}
}

func extractTurnNumber(subID string) int {
	// "turn-001" → 1
	n := 0
	fmt.Sscanf(subID, "turn-%d", &n)
	return n
}

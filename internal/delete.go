package internal

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
)

func (s *StoragePlugin) Delete(_ context.Context, req *pluginv1.StorageDeleteRequest) (*pluginv1.StorageDeleteResponse, error) {
	route := routePath(req.Path)

	s.mu.Lock()
	defer s.mu.Unlock()

	var table, pkCol, pkVal string

	switch route.Entity {
	case entityFeature:
		table, pkCol, pkVal = "features", "id", route.EntityID
	case entityPerson:
		table, pkCol, pkVal = "persons", "id", route.EntityID
	case entityPlan:
		table, pkCol, pkVal = "plans", "id", route.EntityID
	case entityRequest:
		table, pkCol, pkVal = "requests", "id", route.EntityID
	case entityAssignmentRule:
		table, pkCol, pkVal = "assignment_rules", "id", route.EntityID
	case entityNote:
		table, pkCol, pkVal = "notes", "id", route.EntityID
	case entityProject:
		table, pkCol, pkVal = "projects", "slug", route.EntityID
	case entityWIPConfig:
		table, pkCol, pkVal = "wip_config", "project_id", route.EntityID
	case entitySession:
		table, pkCol, pkVal = "sessions", "id", route.EntityID
	case entitySessionTurn:
		// Delete specific turn.
		turnNum := extractTurnNumber(route.SubID)
		result, err := s.db.Exec(`DELETE FROM session_turns WHERE session_id = ? AND turn_number = ?`,
			route.EntityID, turnNum)
		if err != nil {
			return &pluginv1.StorageDeleteResponse{Success: false}, fmt.Errorf("delete turn: %w", err)
		}
		n, _ := result.RowsAffected()
		if n == 0 {
			return &pluginv1.StorageDeleteResponse{Success: false}, fmt.Errorf("turn not found: %s/%s", route.EntityID, route.SubID)
		}
		go s.deleteMarkdown(req.Path)
		return &pluginv1.StorageDeleteResponse{Success: true}, nil
	case entityPack:
		// Delete all packs (registry reset).
		_, err := s.db.Exec(`DELETE FROM packs`)
		if err != nil {
			return &pluginv1.StorageDeleteResponse{Success: false}, fmt.Errorf("delete packs: %w", err)
		}
		go s.deleteMarkdown(req.Path)
		return &pluginv1.StorageDeleteResponse{Success: true}, nil
	default:
		// KV fallback.
		result, err := s.db.Exec(`DELETE FROM kv_store WHERE path = ?`, route.Raw)
		if err != nil {
			return &pluginv1.StorageDeleteResponse{Success: false}, fmt.Errorf("delete kv: %w", err)
		}
		n, _ := result.RowsAffected()
		if n == 0 {
			return &pluginv1.StorageDeleteResponse{Success: false}, fmt.Errorf("not found: %s", req.Path)
		}
		go s.deleteMarkdown(req.Path)
		return &pluginv1.StorageDeleteResponse{Success: true}, nil
	}

	result, err := s.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE %s = ?`, table, pkCol), pkVal)
	if err != nil {
		return &pluginv1.StorageDeleteResponse{Success: false}, fmt.Errorf("delete from %s: %w", table, err)
	}

	n, _ := result.RowsAffected()
	if n == 0 {
		return &pluginv1.StorageDeleteResponse{Success: false}, fmt.Errorf("not found: %s", req.Path)
	}

	go s.deleteMarkdown(req.Path)
	return &pluginv1.StorageDeleteResponse{Success: true}, nil
}

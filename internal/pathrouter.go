package internal

import (
	"path/filepath"
	"strings"
)

// entityType identifies which SQL table a storage path maps to.
type entityType int

const (
	entityUnknown entityType = iota
	entityProject
	entityFeature
	entityPerson
	entityPlan
	entityRequest
	entityAssignmentRule
	entityNote
	entityWIPConfig
	entitySession
	entitySessionTurn
	entityPack
	entityStack
	entityDelegation
	entityDoc
	entitySkill
	entityAgent
	entityProjectSkill
	entityProjectAgent
	entityKV // fallback key-value store
)

// routedPath contains the parsed result of routing a storage path.
type routedPath struct {
	Entity    entityType
	Table     string
	ProjectID string // project slug (empty for non-project entities)
	EntityID  string // e.g. FEAT-ABC, PERS-XYZ
	SubID     string // e.g. turn number for session turns
	Raw       string // original path
}

// routePath translates a storage path into a table + entity ID.
//
// Examples:
//
//	"my-project/features/FEAT-ABC.md"             → features, project=my-project, id=FEAT-ABC
//	"my-project/persons/PERS-ABC.md"              → persons
//	"my-project/plans/PLAN-ABC.md"                → plans
//	"my-project/requests/REQ-ABC.md"              → requests
//	"my-project/assignment-rules/RULE-ABC.md"     → assignment_rules
//	"my-project/notes/note-abc.md"                → notes
//	"my-project/project.json"                     → projects (config)
//	"my-project/wip.json"                         → wip_config
//	"bridge/sessions/<uuid>.md"                   → sessions
//	"bridge/sessions/<uuid>/turn-001.md"          → session_turns
//	".packs/registry.json"                        → packs (special)
//	".stacks/stacks.json"                         → stacks
func routePath(storagePath string) routedPath {
	// Normalize: strip leading/trailing slashes, clean path.
	p := strings.Trim(filepath.ToSlash(filepath.Clean(storagePath)), "/")

	result := routedPath{Raw: storagePath}
	parts := strings.Split(p, "/")

	if len(parts) == 0 {
		result.Entity = entityKV
		result.Table = "kv_store"
		return result
	}

	// Handle dot-prefixed special paths.
	switch {
	case parts[0] == ".packs":
		result.Entity = entityPack
		result.Table = "packs"
		return result
	case parts[0] == ".stacks":
		result.Entity = entityStack
		result.Table = "stacks"
		return result
	case parts[0] == ".skills":
		result.Entity = entitySkill
		result.Table = "skills"
		if len(parts) >= 2 {
			result.EntityID = stripExt(parts[1])
		}
		return result
	case parts[0] == ".agents":
		result.Entity = entityAgent
		result.Table = "agents"
		if len(parts) >= 2 {
			result.EntityID = stripExt(parts[1])
		}
		return result
	}

	// Bridge paths: bridge/sessions/<id>.md or bridge/sessions/<id>/turn-NNN.md
	if parts[0] == "bridge" && len(parts) >= 3 && parts[1] == "sessions" {
		sessionID := stripExt(parts[2])
		if len(parts) == 3 {
			result.Entity = entitySession
			result.Table = "sessions"
			result.EntityID = sessionID
			return result
		}
		if len(parts) == 4 {
			result.Entity = entitySessionTurn
			result.Table = "session_turns"
			result.EntityID = sessionID
			result.SubID = stripExt(parts[3]) // e.g. "turn-001"
			return result
		}
	}

	// Project-scoped paths: <project>/<type>/<id>.md or <project>/<file>
	if len(parts) >= 1 {
		result.ProjectID = parts[0]
	}

	if len(parts) == 2 {
		// <project>/project.json or <project>/wip.json
		switch parts[1] {
		case "project.json":
			result.Entity = entityProject
			result.Table = "projects"
			result.EntityID = result.ProjectID
			return result
		case "wip.json":
			result.Entity = entityWIPConfig
			result.Table = "wip_config"
			result.EntityID = result.ProjectID
			return result
		}
	}

	if len(parts) == 3 {
		entityID := stripExt(parts[2])
		result.EntityID = entityID

		switch parts[1] {
		case "features":
			result.Entity = entityFeature
			result.Table = "features"
			return result
		case "persons":
			result.Entity = entityPerson
			result.Table = "persons"
			return result
		case "plans":
			result.Entity = entityPlan
			result.Table = "plans"
			return result
		case "requests":
			result.Entity = entityRequest
			result.Table = "requests"
			return result
		case "delegations":
			result.Entity = entityDelegation
			result.Table = "delegations"
			return result
		case "assignment-rules":
			result.Entity = entityAssignmentRule
			result.Table = "assignment_rules"
			return result
		case "notes":
			result.Entity = entityNote
			result.Table = "notes"
			return result
		case "docs":
			result.Entity = entityDoc
			result.Table = "docs"
			return result
		}
	}

	// Fallback: store in kv_store.
	result.Entity = entityKV
	result.Table = "kv_store"
	return result
}

// stripExt removes the file extension from a filename.
func stripExt(name string) string {
	ext := filepath.Ext(name)
	if ext == "" {
		return name
	}
	return name[:len(name)-len(ext)]
}

// isListPrefix checks if a prefix matches a known entity directory.
// Returns the entity type and project ID for list queries.
func isListPrefix(prefix string) (entityType, string) {
	p := strings.Trim(filepath.ToSlash(filepath.Clean(prefix)), "/")
	parts := strings.Split(p, "/")

	if len(parts) == 0 {
		return entityKV, ""
	}

	switch {
	case parts[0] == ".packs":
		return entityPack, ""
	case parts[0] == ".stacks":
		return entityStack, ""
	case parts[0] == ".skills":
		return entitySkill, ""
	case parts[0] == ".agents":
		return entityAgent, ""
	case parts[0] == "bridge" && len(parts) >= 2 && parts[1] == "sessions":
		if len(parts) == 3 {
			// bridge/sessions/<id>/ — list turns
			return entitySessionTurn, stripExt(parts[2])
		}
		return entitySession, ""
	}

	if len(parts) == 2 {
		projectID := parts[0]
		switch parts[1] {
		case "features":
			return entityFeature, projectID
		case "persons":
			return entityPerson, projectID
		case "plans":
			return entityPlan, projectID
		case "requests":
			return entityRequest, projectID
		case "delegations":
			return entityDelegation, projectID
		case "assignment-rules":
			return entityAssignmentRule, projectID
		case "notes":
			return entityNote, projectID
		case "docs":
			return entityDoc, projectID
		}
	}

	return entityKV, ""
}

package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestrated-mcp/framework/gen/go/orchestra/plugin/v1"
	"github.com/orchestrated-mcp/framework/libs/go/helpers"
	"github.com/orchestrated-mcp/framework/libs/go/types"
	"github.com/orchestrated-mcp/framework/plugins/tools-features/internal/storage"
	"google.golang.org/protobuf/types/known/structpb"
)

// ---------- Schemas ----------

func AdvanceFeatureSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"project_id": map[string]any{"type": "string", "description": "Project slug"},
			"feature_id": map[string]any{"type": "string", "description": "Feature ID"},
			"evidence":   map[string]any{"type": "string", "description": "Evidence for the transition (optional)"},
		},
		"required": []any{"project_id", "feature_id"},
	})
	return s
}

func RejectFeatureSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"project_id": map[string]any{"type": "string", "description": "Project slug"},
			"feature_id": map[string]any{"type": "string", "description": "Feature ID"},
			"reason":     map[string]any{"type": "string", "description": "Reason for rejection"},
		},
		"required": []any{"project_id", "feature_id", "reason"},
	})
	return s
}

func GetNextFeatureSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"project_id": map[string]any{"type": "string", "description": "Project slug"},
			"status":     map[string]any{"type": "string", "description": "Filter by status (optional)"},
			"assignee":   map[string]any{"type": "string", "description": "Filter by assignee (optional)"},
		},
		"required": []any{"project_id"},
	})
	return s
}

func SetCurrentFeatureSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"project_id": map[string]any{"type": "string", "description": "Project slug"},
			"feature_id": map[string]any{"type": "string", "description": "Feature ID"},
		},
		"required": []any{"project_id", "feature_id"},
	})
	return s
}

func GetWorkflowStatusSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"project_id": map[string]any{"type": "string", "description": "Project slug"},
		},
		"required": []any{"project_id"},
	})
	return s
}

// ---------- Handlers ----------

// AdvanceFeature advances a feature to the next valid status in the workflow.
func AdvanceFeature(store *storage.FeatureStorage) ToolHandler {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "project_id", "feature_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		projectID := helpers.GetString(req.Arguments, "project_id")
		featureID := helpers.GetString(req.Arguments, "feature_id")
		evidence := helpers.GetString(req.Arguments, "evidence")

		feat, body, version, err := store.ReadFeature(ctx, projectID, featureID)
		if err != nil {
			return helpers.ErrorResult("not_found", err.Error()), nil
		}

		nextStatuses := types.NextStatuses(feat.Status)
		if len(nextStatuses) == 0 {
			return helpers.ErrorResult("workflow_error",
				fmt.Sprintf("feature %s is in terminal status %q and cannot be advanced", featureID, feat.Status)), nil
		}

		// Take the first valid transition (the "happy path").
		newStatus := nextStatuses[0]
		oldStatus := feat.Status
		feat.Status = newStatus
		feat.UpdatedAt = helpers.NowISO()

		// Append evidence to body if provided.
		if evidence != "" {
			body += fmt.Sprintf("\n\n---\n**%s -> %s**: %s\n", oldStatus, newStatus, evidence)
		}

		_, err = store.WriteFeature(ctx, projectID, featureID, feat, body, version)
		if err != nil {
			return helpers.ErrorResult("storage_error", err.Error()), nil
		}

		return helpers.JSONResult(map[string]any{
			"feature_id": featureID,
			"from":       string(oldStatus),
			"to":         string(newStatus),
			"message":    fmt.Sprintf("Advanced %s from %s to %s", featureID, oldStatus, newStatus),
		})
	}
}

// RejectFeature sets a feature's status to needs-edits.
func RejectFeature(store *storage.FeatureStorage) ToolHandler {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "project_id", "feature_id", "reason"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		projectID := helpers.GetString(req.Arguments, "project_id")
		featureID := helpers.GetString(req.Arguments, "feature_id")
		reason := helpers.GetString(req.Arguments, "reason")

		feat, body, version, err := store.ReadFeature(ctx, projectID, featureID)
		if err != nil {
			return helpers.ErrorResult("not_found", err.Error()), nil
		}

		if !types.CanTransition(feat.Status, types.StatusNeedsEdits) {
			return helpers.ErrorResult("workflow_error",
				fmt.Sprintf("cannot reject feature from status %q", feat.Status)), nil
		}

		oldStatus := feat.Status
		feat.Status = types.StatusNeedsEdits
		feat.UpdatedAt = helpers.NowISO()

		body += fmt.Sprintf("\n\n---\n**Rejected (%s -> needs-edits)**: %s\n", oldStatus, reason)

		_, err = store.WriteFeature(ctx, projectID, featureID, feat, body, version)
		if err != nil {
			return helpers.ErrorResult("storage_error", err.Error()), nil
		}

		return helpers.JSONResult(map[string]any{
			"feature_id": featureID,
			"from":       string(oldStatus),
			"to":         string(types.StatusNeedsEdits),
			"reason":     reason,
			"message":    fmt.Sprintf("Rejected %s: %s", featureID, reason),
		})
	}
}

// GetNextFeature returns the next feature to work on based on filters.
func GetNextFeature(store *storage.FeatureStorage) ToolHandler {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "project_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		projectID := helpers.GetString(req.Arguments, "project_id")
		statusFilter := helpers.GetString(req.Arguments, "status")
		assigneeFilter := helpers.GetString(req.Arguments, "assignee")

		features, err := store.ListFeatures(ctx, projectID)
		if err != nil {
			return helpers.ErrorResult("storage_error", err.Error()), nil
		}

		// Default: find features in "todo" status.
		if statusFilter == "" {
			statusFilter = string(types.StatusTodo)
		}

		// Priority order: P0 > P1 > P2 > P3.
		priorityRank := map[string]int{"P0": 0, "P1": 1, "P2": 2, "P3": 3}

		var best *types.FeatureData
		bestRank := 999

		for _, f := range features {
			if string(f.Status) != statusFilter {
				continue
			}
			if assigneeFilter != "" && f.Assignee != assigneeFilter {
				continue
			}
			rank, ok := priorityRank[f.Priority]
			if !ok {
				rank = 99
			}
			if best == nil || rank < bestRank {
				best = f
				bestRank = rank
			}
		}

		if best == nil {
			return helpers.TextResult("No features found matching the criteria."), nil
		}

		return helpers.JSONResult(map[string]any{
			"feature": best,
			"message": fmt.Sprintf("Next feature: %s - %s", best.ID, best.Title),
		})
	}
}

// SetCurrentFeature sets a feature's status to in-progress.
func SetCurrentFeature(store *storage.FeatureStorage) ToolHandler {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "project_id", "feature_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		projectID := helpers.GetString(req.Arguments, "project_id")
		featureID := helpers.GetString(req.Arguments, "feature_id")

		feat, body, version, err := store.ReadFeature(ctx, projectID, featureID)
		if err != nil {
			return helpers.ErrorResult("not_found", err.Error()), nil
		}

		if !types.CanTransition(feat.Status, types.StatusInProgress) {
			return helpers.ErrorResult("workflow_error",
				fmt.Sprintf("cannot set to in-progress from status %q", feat.Status)), nil
		}

		oldStatus := feat.Status
		feat.Status = types.StatusInProgress
		feat.UpdatedAt = helpers.NowISO()

		_, err = store.WriteFeature(ctx, projectID, featureID, feat, body, version)
		if err != nil {
			return helpers.ErrorResult("storage_error", err.Error()), nil
		}

		return helpers.JSONResult(map[string]any{
			"feature_id": featureID,
			"from":       string(oldStatus),
			"to":         string(types.StatusInProgress),
			"message":    fmt.Sprintf("Set %s to in-progress", featureID),
		})
	}
}

// GetWorkflowStatus returns feature counts per status for a project.
func GetWorkflowStatus(store *storage.FeatureStorage) ToolHandler {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "project_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		projectID := helpers.GetString(req.Arguments, "project_id")

		features, err := store.ListFeatures(ctx, projectID)
		if err != nil {
			return helpers.ErrorResult("storage_error", err.Error()), nil
		}

		statusCounts := make(map[string]int)
		for _, f := range features {
			statusCounts[string(f.Status)]++
		}

		return helpers.JSONResult(map[string]any{
			"project_id":     projectID,
			"total_features": len(features),
			"by_status":      statusCounts,
		})
	}
}

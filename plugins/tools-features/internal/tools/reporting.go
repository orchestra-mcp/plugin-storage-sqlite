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

func GetProgressSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"project_id": map[string]any{"type": "string", "description": "Project slug"},
		},
		"required": []any{"project_id"},
	})
	return s
}

func GetBlockedFeaturesSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"project_id": map[string]any{"type": "string", "description": "Project slug"},
		},
		"required": []any{"project_id"},
	})
	return s
}

func GetReviewQueueSchema() *structpb.Struct {
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

// GetProgress returns completion percentage and feature counts by status.
func GetProgress(store *storage.FeatureStorage) ToolHandler {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "project_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		projectID := helpers.GetString(req.Arguments, "project_id")

		features, err := store.ListFeatures(ctx, projectID)
		if err != nil {
			return helpers.ErrorResult("storage_error", err.Error()), nil
		}

		total := len(features)
		done := 0
		statusCounts := make(map[string]int)

		for _, f := range features {
			statusCounts[string(f.Status)]++
			if f.Status == types.StatusDone {
				done++
			}
		}

		var pctDone float64
		if total > 0 {
			pctDone = float64(done) / float64(total) * 100
		}

		return helpers.JSONResult(map[string]any{
			"project_id":      projectID,
			"total_features":  total,
			"done":            done,
			"percent_done":    fmt.Sprintf("%.1f%%", pctDone),
			"by_status":       statusCounts,
		})
	}
}

// GetBlockedFeatures returns features that are blocked by unfinished dependencies.
func GetBlockedFeatures(store *storage.FeatureStorage) ToolHandler {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "project_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		projectID := helpers.GetString(req.Arguments, "project_id")

		features, err := store.ListFeatures(ctx, projectID)
		if err != nil {
			return helpers.ErrorResult("storage_error", err.Error()), nil
		}

		// Build a status map for quick lookup.
		statusMap := make(map[string]types.FeatureStatus)
		for _, f := range features {
			statusMap[f.ID] = f.Status
		}

		type blockedInfo struct {
			Feature     *types.FeatureData `json:"feature"`
			BlockedBy   []string           `json:"blocked_by"`
		}

		var blocked []blockedInfo
		for _, f := range features {
			if len(f.DependsOn) == 0 {
				continue
			}
			var unblockers []string
			for _, depID := range f.DependsOn {
				depStatus, exists := statusMap[depID]
				if !exists || depStatus != types.StatusDone {
					unblockers = append(unblockers, depID)
				}
			}
			if len(unblockers) > 0 {
				blocked = append(blocked, blockedInfo{
					Feature:   f,
					BlockedBy: unblockers,
				})
			}
		}

		if blocked == nil {
			blocked = []blockedInfo{}
		}

		return helpers.JSONResult(map[string]any{
			"blocked":  blocked,
			"count":    len(blocked),
		})
	}
}

// GetReviewQueue returns all features currently awaiting review.
func GetReviewQueue(store *storage.FeatureStorage) ToolHandler {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "project_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		projectID := helpers.GetString(req.Arguments, "project_id")

		features, err := store.ListFeatures(ctx, projectID)
		if err != nil {
			return helpers.ErrorResult("storage_error", err.Error()), nil
		}

		var inReview []*types.FeatureData
		for _, f := range features {
			if f.Status == types.StatusInReview {
				inReview = append(inReview, f)
			}
		}

		if inReview == nil {
			inReview = []*types.FeatureData{}
		}

		return helpers.JSONResult(map[string]any{
			"features": inReview,
			"count":    len(inReview),
		})
	}
}

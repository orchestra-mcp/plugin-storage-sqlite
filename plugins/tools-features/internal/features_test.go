package internal

import (
	"context"
	"encoding/json"
	"testing"

	pluginv1 "github.com/orchestrated-mcp/framework/gen/go/orchestra/plugin/v1"
	"github.com/orchestrated-mcp/framework/plugins/tools-features/internal/storage"
	"github.com/orchestrated-mcp/framework/plugins/tools-features/internal/tools"
	"google.golang.org/protobuf/types/known/structpb"
)

// testEnv sets up an in-memory storage backend and returns the feature storage
// wrapper plus a context.
func testEnv() (*storage.FeatureStorage, context.Context) {
	mem := storage.NewInMemoryStorage()
	store := storage.NewFeatureStorage(mem)
	return store, context.Background()
}

// callTool is a test helper that invokes a tool handler with the given
// arguments and returns the response.
func callTool(t *testing.T, handler tools.ToolHandler, args map[string]any) *pluginv1.ToolResponse {
	t.Helper()
	s, err := structpb.NewStruct(args)
	if err != nil {
		t.Fatalf("failed to create args struct: %v", err)
	}
	ctx := context.Background()
	resp, err := handler(ctx, &pluginv1.ToolRequest{
		ToolName:  "test",
		Arguments: s,
	})
	if err != nil {
		t.Fatalf("tool handler returned error: %v", err)
	}
	return resp
}

// resultMap extracts the result struct as a map from a ToolResponse.
func resultMap(t *testing.T, resp *pluginv1.ToolResponse) map[string]any {
	t.Helper()
	if resp.Result == nil {
		t.Fatal("response has no result")
	}
	return resp.Result.AsMap()
}

// createTestProject creates a project and returns the store.
func createTestProject(t *testing.T, store *storage.FeatureStorage, name string) {
	t.Helper()
	resp := callTool(t, tools.CreateProject(store), map[string]any{
		"name": name,
	})
	if !resp.Success {
		t.Fatalf("create_project failed: %s", resp.ErrorMessage)
	}
}

// createTestFeature creates a feature and returns the feature ID.
func createTestFeature(t *testing.T, store *storage.FeatureStorage, projectID, title string) string {
	t.Helper()
	resp := callTool(t, tools.CreateFeature(store), map[string]any{
		"project_id": projectID,
		"title":      title,
	})
	if !resp.Success {
		t.Fatalf("create_feature failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	id, ok := m["id"].(string)
	if !ok || id == "" {
		t.Fatal("create_feature did not return an id")
	}
	return id
}

func TestCreateAndGetProject(t *testing.T) {
	store, _ := testEnv()

	// Create project.
	resp := callTool(t, tools.CreateProject(store), map[string]any{
		"name":        "My App",
		"description": "A test project",
	})
	if !resp.Success {
		t.Fatalf("create_project failed: %s", resp.ErrorMessage)
	}

	// Get project status.
	resp = callTool(t, tools.GetProjectStatus(store), map[string]any{
		"project_id": "my-app",
	})
	if !resp.Success {
		t.Fatalf("get_project_status failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)

	proj, ok := m["project"].(map[string]any)
	if !ok {
		t.Fatal("expected project in result")
	}
	if proj["name"] != "My App" {
		t.Errorf("project name: got %q, want %q", proj["name"], "My App")
	}
	if proj["slug"] != "my-app" {
		t.Errorf("project slug: got %q, want %q", proj["slug"], "my-app")
	}

	// List projects.
	resp = callTool(t, tools.ListProjects(store), map[string]any{})
	if !resp.Success {
		t.Fatalf("list_projects failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	count, _ := m["count"].(float64)
	if count != 1 {
		t.Errorf("expected 1 project, got %v", count)
	}
}

func TestCreateAndGetFeature(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Test Project")

	// Create feature.
	resp := callTool(t, tools.CreateFeature(store), map[string]any{
		"project_id":  "test-project",
		"title":       "User Authentication",
		"description": "Implement OAuth2 login flow",
		"priority":    "P0",
	})
	if !resp.Success {
		t.Fatalf("create_feature failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	featureID := m["id"].(string)

	if m["status"] != "backlog" {
		t.Errorf("expected status backlog, got %v", m["status"])
	}
	if m["priority"] != "P0" {
		t.Errorf("expected priority P0, got %v", m["priority"])
	}

	// Get feature.
	resp = callTool(t, tools.GetFeature(store), map[string]any{
		"project_id": "test-project",
		"feature_id": featureID,
	})
	if !resp.Success {
		t.Fatalf("get_feature failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	feat, _ := m["feature"].(map[string]any)
	if feat["title"] != "User Authentication" {
		t.Errorf("feature title: got %q, want %q", feat["title"], "User Authentication")
	}
	body, _ := m["body"].(string)
	if body == "" {
		t.Error("expected non-empty body")
	}
}

func TestListFeatures(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "List Test")

	createTestFeature(t, store, "list-test", "Feature A")
	createTestFeature(t, store, "list-test", "Feature B")
	createTestFeature(t, store, "list-test", "Feature C")

	resp := callTool(t, tools.ListFeatures(store), map[string]any{
		"project_id": "list-test",
	})
	if !resp.Success {
		t.Fatalf("list_features failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	count, _ := m["count"].(float64)
	if count != 3 {
		t.Errorf("expected 3 features, got %v", count)
	}

	// Test status filter.
	resp = callTool(t, tools.ListFeatures(store), map[string]any{
		"project_id": "list-test",
		"status":     "in-progress",
	})
	if !resp.Success {
		t.Fatalf("list_features with status filter failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	count, _ = m["count"].(float64)
	if count != 0 {
		t.Errorf("expected 0 in-progress features, got %v", count)
	}
}

func TestWorkflowAdvance(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Workflow Test")
	featureID := createTestFeature(t, store, "workflow-test", "Auth Feature")

	// Expected workflow path: backlog -> todo -> in-progress -> ready-for-testing
	// -> in-testing -> ready-for-docs -> in-docs -> documented -> in-review -> done
	expectedTransitions := []struct {
		from string
		to   string
	}{
		{"backlog", "todo"},
		{"todo", "in-progress"},
		{"in-progress", "ready-for-testing"},
		{"ready-for-testing", "in-testing"},
		{"in-testing", "ready-for-docs"},
		{"ready-for-docs", "in-docs"},
		{"in-docs", "documented"},
		{"documented", "in-review"},
		{"in-review", "done"},
	}

	for _, tt := range expectedTransitions {
		resp := callTool(t, tools.AdvanceFeature(store), map[string]any{
			"project_id": "workflow-test",
			"feature_id": featureID,
			"evidence":   "Completed " + tt.from,
		})
		if !resp.Success {
			t.Fatalf("advance from %s failed: %s", tt.from, resp.ErrorMessage)
		}
		m := resultMap(t, resp)
		if m["from"] != tt.from {
			t.Errorf("expected from %q, got %q", tt.from, m["from"])
		}
		if m["to"] != tt.to {
			t.Errorf("expected to %q, got %q", tt.to, m["to"])
		}
	}

	// Advancing from done should fail.
	resp := callTool(t, tools.AdvanceFeature(store), map[string]any{
		"project_id": "workflow-test",
		"feature_id": featureID,
	})
	if resp.Success {
		t.Error("expected advance from done to fail")
	}
	if resp.ErrorCode != "workflow_error" {
		t.Errorf("expected workflow_error, got %q", resp.ErrorCode)
	}
}

func TestWorkflowReject(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Reject Test")
	featureID := createTestFeature(t, store, "reject-test", "Rejectable Feature")

	// Advance to in-review: backlog -> todo -> in-progress -> ready-for-testing ->
	// in-testing -> ready-for-docs -> in-docs -> documented -> in-review
	for i := 0; i < 8; i++ {
		resp := callTool(t, tools.AdvanceFeature(store), map[string]any{
			"project_id": "reject-test",
			"feature_id": featureID,
		})
		if !resp.Success {
			t.Fatalf("advance step %d failed: %s", i, resp.ErrorMessage)
		}
	}

	// Reject from in-review.
	resp := callTool(t, tools.RejectFeature(store), map[string]any{
		"project_id": "reject-test",
		"feature_id": featureID,
		"reason":     "Missing error handling",
	})
	if !resp.Success {
		t.Fatalf("reject_feature failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	if m["from"] != "in-review" {
		t.Errorf("expected from in-review, got %v", m["from"])
	}
	if m["to"] != "needs-edits" {
		t.Errorf("expected to needs-edits, got %v", m["to"])
	}

	// From needs-edits, can go back to in-progress.
	resp = callTool(t, tools.SetCurrentFeature(store), map[string]any{
		"project_id": "reject-test",
		"feature_id": featureID,
	})
	if !resp.Success {
		t.Fatalf("set_current_feature from needs-edits failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	if m["to"] != "in-progress" {
		t.Errorf("expected to in-progress, got %v", m["to"])
	}
}

func TestDependencies(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Deps Test")
	feat1 := createTestFeature(t, store, "deps-test", "Foundation")
	feat2 := createTestFeature(t, store, "deps-test", "Build on Foundation")

	// Add dependency: feat2 depends on feat1.
	resp := callTool(t, tools.AddDependency(store), map[string]any{
		"project_id":    "deps-test",
		"feature_id":    feat2,
		"depends_on_id": feat1,
	})
	if !resp.Success {
		t.Fatalf("add_dependency failed: %s", resp.ErrorMessage)
	}

	// Get dependency graph.
	resp = callTool(t, tools.GetDependencyGraph(store), map[string]any{
		"project_id": "deps-test",
	})
	if !resp.Success {
		t.Fatalf("get_dependency_graph failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	edgesRaw, _ := m["edges"].([]any)
	if len(edgesRaw) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edgesRaw))
	}
	edge, _ := edgesRaw[0].(map[string]any)
	if edge["from"] != feat2 {
		t.Errorf("edge from: got %q, want %q", edge["from"], feat2)
	}
	if edge["to"] != feat1 {
		t.Errorf("edge to: got %q, want %q", edge["to"], feat1)
	}

	// Verify feat2 shows as blocked.
	resp = callTool(t, tools.GetBlockedFeatures(store), map[string]any{
		"project_id": "deps-test",
	})
	if !resp.Success {
		t.Fatalf("get_blocked_features failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	blockedCount, _ := m["count"].(float64)
	if blockedCount != 1 {
		t.Errorf("expected 1 blocked feature, got %v", blockedCount)
	}

	// Remove dependency.
	resp = callTool(t, tools.RemoveDependency(store), map[string]any{
		"project_id":    "deps-test",
		"feature_id":    feat2,
		"depends_on_id": feat1,
	})
	if !resp.Success {
		t.Fatalf("remove_dependency failed: %s", resp.ErrorMessage)
	}

	// Verify no edges remain.
	resp = callTool(t, tools.GetDependencyGraph(store), map[string]any{
		"project_id": "deps-test",
	})
	if !resp.Success {
		t.Fatalf("get_dependency_graph after remove failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	edgesRaw, _ = m["edges"].([]any)
	if len(edgesRaw) != 0 {
		t.Errorf("expected 0 edges after removal, got %d", len(edgesRaw))
	}
}

func TestLabels(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Labels Test")
	featureID := createTestFeature(t, store, "labels-test", "Labeled Feature")

	// Add labels.
	resp := callTool(t, tools.AddLabels(store), map[string]any{
		"project_id": "labels-test",
		"feature_id": featureID,
		"labels":     []any{"backend", "urgent"},
	})
	if !resp.Success {
		t.Fatalf("add_labels failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	labels, _ := m["labels"].([]any)
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}

	// Add duplicate label (should not duplicate).
	resp = callTool(t, tools.AddLabels(store), map[string]any{
		"project_id": "labels-test",
		"feature_id": featureID,
		"labels":     []any{"backend", "frontend"},
	})
	if !resp.Success {
		t.Fatalf("add_labels (with duplicate) failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	labels, _ = m["labels"].([]any)
	if len(labels) != 3 {
		t.Errorf("expected 3 labels (no duplicate), got %d", len(labels))
	}

	// Remove a label.
	resp = callTool(t, tools.RemoveLabels(store), map[string]any{
		"project_id": "labels-test",
		"feature_id": featureID,
		"labels":     []any{"urgent"},
	})
	if !resp.Success {
		t.Fatalf("remove_labels failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	labels, _ = m["labels"].([]any)
	if len(labels) != 2 {
		t.Errorf("expected 2 labels after removal, got %d", len(labels))
	}
}

func TestAssignment(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Assign Test")
	featureID := createTestFeature(t, store, "assign-test", "Assignable Feature")

	// Assign feature.
	resp := callTool(t, tools.AssignFeature(store), map[string]any{
		"project_id": "assign-test",
		"feature_id": featureID,
		"assignee":   "alice",
	})
	if !resp.Success {
		t.Fatalf("assign_feature failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	if m["assignee"] != "alice" {
		t.Errorf("expected assignee alice, got %v", m["assignee"])
	}

	// Verify via get_feature.
	resp = callTool(t, tools.GetFeature(store), map[string]any{
		"project_id": "assign-test",
		"feature_id": featureID,
	})
	if !resp.Success {
		t.Fatalf("get_feature failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	feat, _ := m["feature"].(map[string]any)
	if feat["assignee"] != "alice" {
		t.Errorf("expected assignee alice from get_feature, got %v", feat["assignee"])
	}

	// Unassign feature.
	resp = callTool(t, tools.UnassignFeature(store), map[string]any{
		"project_id": "assign-test",
		"feature_id": featureID,
	})
	if !resp.Success {
		t.Fatalf("unassign_feature failed: %s", resp.ErrorMessage)
	}

	// Verify unassigned. When assignee is empty, it may be omitted from JSON
	// serialization (nil in the map) or set to empty string.
	resp = callTool(t, tools.GetFeature(store), map[string]any{
		"project_id": "assign-test",
		"feature_id": featureID,
	})
	m = resultMap(t, resp)
	feat, _ = m["feature"].(map[string]any)
	assignee, _ := feat["assignee"].(string)
	if assignee != "" {
		t.Errorf("expected empty assignee after unassign, got %v", feat["assignee"])
	}
}

func TestSearch(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Search Test")

	createTestFeature(t, store, "search-test", "User Authentication")
	createTestFeature(t, store, "search-test", "Database Migration")
	createTestFeature(t, store, "search-test", "User Profile Page")

	// Search for "user".
	resp := callTool(t, tools.SearchFeatures(store), map[string]any{
		"project_id": "search-test",
		"query":      "user",
	})
	if !resp.Success {
		t.Fatalf("search_features failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	count, _ := m["count"].(float64)
	if count != 2 {
		t.Errorf("expected 2 search results for 'user', got %v", count)
	}

	// Search for something that does not exist.
	resp = callTool(t, tools.SearchFeatures(store), map[string]any{
		"project_id": "search-test",
		"query":      "nonexistent",
	})
	if !resp.Success {
		t.Fatalf("search_features (no results) failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	count, _ = m["count"].(float64)
	if count != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %v", count)
	}
}

func TestEstimate(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Estimate Test")
	featureID := createTestFeature(t, store, "estimate-test", "Estimable Feature")

	resp := callTool(t, tools.SetEstimate(store), map[string]any{
		"project_id": "estimate-test",
		"feature_id": featureID,
		"estimate":   "L",
	})
	if !resp.Success {
		t.Fatalf("set_estimate failed: %s", resp.ErrorMessage)
	}

	// Verify.
	resp = callTool(t, tools.GetFeature(store), map[string]any{
		"project_id": "estimate-test",
		"feature_id": featureID,
	})
	m := resultMap(t, resp)
	feat, _ := m["feature"].(map[string]any)
	if feat["estimate"] != "L" {
		t.Errorf("expected estimate L, got %v", feat["estimate"])
	}

	// Invalid estimate.
	resp = callTool(t, tools.SetEstimate(store), map[string]any{
		"project_id": "estimate-test",
		"feature_id": featureID,
		"estimate":   "XXL",
	})
	if resp.Success {
		t.Error("expected set_estimate with invalid value to fail")
	}
}

func TestNotes(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Notes Test")
	featureID := createTestFeature(t, store, "notes-test", "Notable Feature")

	resp := callTool(t, tools.SaveNote(store), map[string]any{
		"project_id": "notes-test",
		"feature_id": featureID,
		"note":       "This is an important note.",
	})
	if !resp.Success {
		t.Fatalf("save_note failed: %s", resp.ErrorMessage)
	}

	resp = callTool(t, tools.ListNotes(store), map[string]any{
		"project_id": "notes-test",
		"feature_id": featureID,
	})
	if !resp.Success {
		t.Fatalf("list_notes failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	body, _ := m["body"].(string)
	if body == "" {
		t.Error("expected non-empty body with notes")
	}
}

func TestUpdateFeature(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Update Test")
	featureID := createTestFeature(t, store, "update-test", "Original Title")

	resp := callTool(t, tools.UpdateFeature(store), map[string]any{
		"project_id": "update-test",
		"feature_id": featureID,
		"title":      "Updated Title",
		"priority":   "P1",
	})
	if !resp.Success {
		t.Fatalf("update_feature failed: %s", resp.ErrorMessage)
	}

	resp = callTool(t, tools.GetFeature(store), map[string]any{
		"project_id": "update-test",
		"feature_id": featureID,
	})
	m := resultMap(t, resp)
	feat, _ := m["feature"].(map[string]any)
	if feat["title"] != "Updated Title" {
		t.Errorf("expected Updated Title, got %v", feat["title"])
	}
	if feat["priority"] != "P1" {
		t.Errorf("expected P1, got %v", feat["priority"])
	}
}

func TestDeleteProject(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Delete Me")
	createTestFeature(t, store, "delete-me", "Feature in deleted project")

	resp := callTool(t, tools.DeleteProject(store), map[string]any{
		"project_id": "delete-me",
	})
	if !resp.Success {
		t.Fatalf("delete_project failed: %s", resp.ErrorMessage)
	}

	// Verify project is gone.
	resp = callTool(t, tools.GetProjectStatus(store), map[string]any{
		"project_id": "delete-me",
	})
	if resp.Success {
		t.Error("expected get_project_status to fail after deletion")
	}
}

func TestDeleteFeature(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Del Feature Test")
	featureID := createTestFeature(t, store, "del-feature-test", "To Delete")

	resp := callTool(t, tools.DeleteFeature(store), map[string]any{
		"project_id": "del-feature-test",
		"feature_id": featureID,
	})
	if !resp.Success {
		t.Fatalf("delete_feature failed: %s", resp.ErrorMessage)
	}

	// Verify feature is gone.
	resp = callTool(t, tools.GetFeature(store), map[string]any{
		"project_id": "del-feature-test",
		"feature_id": featureID,
	})
	if resp.Success {
		t.Error("expected get_feature to fail after deletion")
	}
}

func TestGetNextFeature(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Next Test")
	id1 := createTestFeature(t, store, "next-test", "Low Priority Feature")
	id2 := createTestFeature(t, store, "next-test", "High Priority Feature")

	// Set one feature to P0.
	callTool(t, tools.UpdateFeature(store), map[string]any{
		"project_id": "next-test",
		"feature_id": id2,
		"priority":   "P0",
	})

	// Advance both to todo (from backlog).
	callTool(t, tools.AdvanceFeature(store), map[string]any{
		"project_id": "next-test",
		"feature_id": id1,
	})
	callTool(t, tools.AdvanceFeature(store), map[string]any{
		"project_id": "next-test",
		"feature_id": id2,
	})

	// Get next feature should return the P0 one.
	resp := callTool(t, tools.GetNextFeature(store), map[string]any{
		"project_id": "next-test",
	})
	if !resp.Success {
		t.Fatalf("get_next_feature failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	feat, _ := m["feature"].(map[string]any)
	if feat["id"] != id2 {
		t.Errorf("expected P0 feature %s, got %v", id2, feat["id"])
	}
}

func TestReviewWorkflow(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Review Test")
	featureID := createTestFeature(t, store, "review-test", "Reviewable Feature")

	// Advance to documented (8 transitions).
	for i := 0; i < 7; i++ {
		resp := callTool(t, tools.AdvanceFeature(store), map[string]any{
			"project_id": "review-test",
			"feature_id": featureID,
		})
		if !resp.Success {
			t.Fatalf("advance step %d failed: %s", i, resp.ErrorMessage)
		}
	}

	// Request review.
	resp := callTool(t, tools.RequestReview(store), map[string]any{
		"project_id": "review-test",
		"feature_id": featureID,
	})
	if !resp.Success {
		t.Fatalf("request_review failed: %s", resp.ErrorMessage)
	}

	// Check pending reviews.
	resp = callTool(t, tools.GetPendingReviews(store), map[string]any{
		"project_id": "review-test",
	})
	if !resp.Success {
		t.Fatalf("get_pending_reviews failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	reviewCount, _ := m["count"].(float64)
	if reviewCount != 1 {
		t.Errorf("expected 1 pending review, got %v", reviewCount)
	}

	// Submit review as needs-edits.
	resp = callTool(t, tools.SubmitReview(store), map[string]any{
		"project_id": "review-test",
		"feature_id": featureID,
		"status":     "needs-edits",
		"comment":    "Needs more tests",
	})
	if !resp.Success {
		t.Fatalf("submit_review failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	if m["to"] != "needs-edits" {
		t.Errorf("expected to needs-edits, got %v", m["to"])
	}
}

func TestProgress(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "Progress Test")
	createTestFeature(t, store, "progress-test", "Feature 1")
	createTestFeature(t, store, "progress-test", "Feature 2")

	resp := callTool(t, tools.GetProgress(store), map[string]any{
		"project_id": "progress-test",
	})
	if !resp.Success {
		t.Fatalf("get_progress failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	total, _ := m["total_features"].(float64)
	if total != 2 {
		t.Errorf("expected 2 total features, got %v", total)
	}
	done, _ := m["done"].(float64)
	if done != 0 {
		t.Errorf("expected 0 done, got %v", done)
	}
}

func TestWIPLimits(t *testing.T) {
	store, _ := testEnv()
	createTestProject(t, store, "WIP Test")

	// Set WIP limit.
	resp := callTool(t, tools.SetWIPLimits(store), map[string]any{
		"project_id":      "wip-test",
		"max_in_progress": 2.0,
	})
	if !resp.Success {
		t.Fatalf("set_wip_limits failed: %s", resp.ErrorMessage)
	}

	// Get WIP limits.
	resp = callTool(t, tools.GetWIPLimits(store), map[string]any{
		"project_id": "wip-test",
	})
	if !resp.Success {
		t.Fatalf("get_wip_limits failed: %s", resp.ErrorMessage)
	}
	m := resultMap(t, resp)
	limit, _ := m["max_in_progress"].(float64)
	if limit != 2 {
		t.Errorf("expected WIP limit 2, got %v", limit)
	}

	// Check WIP limit (should be within since nothing is in-progress).
	resp = callTool(t, tools.CheckWIPLimit(store), map[string]any{
		"project_id": "wip-test",
	})
	if !resp.Success {
		t.Fatalf("check_wip_limit failed: %s", resp.ErrorMessage)
	}
	m = resultMap(t, resp)
	withinLimit, _ := m["within_limit"].(bool)
	if !withinLimit {
		t.Error("expected within_limit to be true with 0 in-progress")
	}
}

// TestValidation verifies that tools reject missing required arguments.
func TestValidation(t *testing.T) {
	store, _ := testEnv()

	tests := []struct {
		name    string
		handler tools.ToolHandler
		args    map[string]any
	}{
		{"create_project missing name", tools.CreateProject(store), map[string]any{}},
		{"create_feature missing project_id", tools.CreateFeature(store), map[string]any{"title": "x"}},
		{"create_feature missing title", tools.CreateFeature(store), map[string]any{"project_id": "x"}},
		{"get_feature missing args", tools.GetFeature(store), map[string]any{}},
		{"assign missing assignee", tools.AssignFeature(store), map[string]any{"project_id": "x", "feature_id": "y"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := callTool(t, tt.handler, tt.args)
			if resp.Success {
				t.Error("expected validation failure")
			}
			if resp.ErrorCode != "validation_error" {
				t.Errorf("expected validation_error, got %q", resp.ErrorCode)
			}
		})
	}
}

// Ensure json import is used (for future tests or if needed).
var _ = json.Marshal

# Step 5: tools.features Plugin

## Status: Complete

## What Was Built

The feature-driven workflow engine — 34 tools for project/feature CRUD, workflow state machine, reviews, dependencies, WIP limits, and reporting. All storage goes through the orchestrator QUIC connection (no direct disk access).

## Module

`github.com/orchestrated-mcp/framework/plugins/tools-features`

## Tools (34 total)

### Project (4)
| Tool | Args | Description |
|------|------|-------------|
| `create_project` | name, description | Create project with slug |
| `list_projects` | — | List all projects |
| `delete_project` | project_id | Delete project + all features |
| `get_project_status` | project_id | Feature counts by status |

### Feature (6)
| Tool | Args | Description |
|------|------|-------------|
| `create_feature` | project_id, title, description, priority? | Create FEAT-XXX |
| `get_feature` | project_id, feature_id | Get data + body |
| `update_feature` | project_id, feature_id, title?, description?, priority? | Update fields |
| `list_features` | project_id, status? | List with optional filter |
| `delete_feature` | project_id, feature_id | Delete feature |
| `search_features` | project_id, query | Text search titles/descriptions |

### Workflow (5)
| Tool | Args | Description |
|------|------|-------------|
| `advance_feature` | project_id, feature_id, evidence? | Next valid status |
| `reject_feature` | project_id, feature_id, reason | → needs-edits |
| `get_next_feature` | project_id, status?, assignee? | Priority-ordered next |
| `set_current_feature` | project_id, feature_id | → in-progress |
| `get_workflow_status` | project_id | Counts per status |

### Review (3)
| Tool | Args | Description |
|------|------|-------------|
| `request_review` | project_id, feature_id | → in-review |
| `submit_review` | project_id, feature_id, status, comment | approved→done / needs-edits |
| `get_pending_reviews` | project_id | List in-review features |

### Dependency (3)
| Tool | Args | Description |
|------|------|-------------|
| `add_dependency` | project_id, feature_id, depends_on_id | Bidirectional |
| `remove_dependency` | project_id, feature_id, depends_on_id | Remove link |
| `get_dependency_graph` | project_id | All edges |

### WIP (3)
| Tool | Args | Description |
|------|------|-------------|
| `set_wip_limits` | project_id, max_in_progress | Set limit |
| `get_wip_limits` | project_id | Get limit |
| `check_wip_limit` | project_id | Current vs max |

### Reporting (3)
| Tool | Args | Description |
|------|------|-------------|
| `get_progress` | project_id | % done, counts |
| `get_blocked_features` | project_id | Blocked by deps |
| `get_review_queue` | project_id | Awaiting review |

### Metadata (7)
| Tool | Args | Description |
|------|------|-------------|
| `add_labels` | project_id, feature_id, labels[] | Add labels |
| `remove_labels` | project_id, feature_id, labels[] | Remove labels |
| `assign_feature` | project_id, feature_id, assignee | Set assignee |
| `unassign_feature` | project_id, feature_id | Clear assignee |
| `set_estimate` | project_id, feature_id, estimate | S/M/L/XL |
| `save_note` | project_id, feature_id, note | Append to body |
| `list_notes` | project_id, feature_id | Return body |

## Architecture

```
Tool Handler → FeatureStorage → StorageClient → Orchestrator (QUIC) → storage.markdown
```

- `StorageClient` interface with single `Send()` method
- `InMemoryStorage` for testing (no QUIC needed)
- `FeatureStorage` provides high-level read/write for features and projects

## Tests (24 pass)

19 test functions + 5 validation sub-tests covering all tools, workflow state machine, dependencies, labels, assignment, search, WIP limits, and input validation.

```bash
cd plugins/tools-features && go test ./... -v
```
